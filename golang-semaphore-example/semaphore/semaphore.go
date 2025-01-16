package semaphore

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"syscall"
	"unsafe"

	"gopkg.in/yaml.v3"
)

/*
#include <sys/sem.h>
typedef struct sembuf sembuf;
*/
import "C"

// LockKey 默认信号量key，会被配置文件的值覆盖
var LockKey = 2216

// ConcurrentNum 默认并发量，会被配置文件的值覆盖
var ConcurrentNum = 30

// LockFilePath 默认锁文件路径，会被配置文件的值覆盖
var LockFilePath = "/usr/local/semaphore/lockFile"

var configFilePath string

var projectPath, _ = filepath.Abs("../")

type Config struct {
	LockKey       int    `yaml:"lockKey"`
	ConcurrentNum int    `yaml:"concurrentNum"`
	LockFilePath  string `yaml:"lockFilePath"`
}

func init() {
	// 配置文件路径在项目根目录下
	configFilePath = filepath.Join(projectPath, "config.yaml")
	configFile, err := ReadConfig(configFilePath)
	if err != nil {
		fmt.Printf("初始化: 读取配置失败: %v\n", err)
	} else {
		LockKey = configFile.LockKey
		ConcurrentNum = configFile.ConcurrentNum
		LockFilePath = configFile.LockFilePath
		LockFilePath = filepath.Join(projectPath, LockFilePath, "lockFile")
		// log.Logger(fmt.Sprintf("读取设定并发量: %d\n", ConcurrentNum))
	}
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

func SetSemaphore() (r1 uintptr, r2 uintptr, err syscall.Errno) {
	r1, r2, err = syscall.Syscall(syscall.SYS_SEMGET, uintptr(LockKey), uintptr(1), uintptr(C.IPC_CREAT|00666))
	if int(r1) < 0 {
		return
	}

	// 准备SETVAL命令的参数
	cmd := uintptr(C.SETVAL)
	val := uintptr(ConcurrentNum)

	// 调用SYS_SEMCTL设置信号量值
	_, _, err = syscall.Syscall6(syscall.SYS_SEMCTL, r1, 0, cmd, val, 0, 0)
	if err != 0 {
		return
	}
	return
}

func GetLockFile() (file *os.File, err error) {
	if !fileExists(LockFilePath) {
		dirPath := filepath.Dir(LockFilePath)
		if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
			panic(err)
		}
		file, err = os.OpenFile(LockFilePath, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
		if err != nil {
			if os.IsExist(err) {
				//lock_file已被其它进程创建
				file, err = os.OpenFile(LockFilePath, os.O_RDWR|os.O_CREATE, 0666)
			}
		}
	} else {
		file, err = os.OpenFile(LockFilePath, os.O_RDWR|os.O_CREATE, 0666)
	}
	return
}

// SemGet 获取信号量组ID
// return : r1-semId, err-syscall.Errno, err2-error
func SemGet() (r1 uintptr, err syscall.Errno, err2 error) {
	r1, _, err = syscall.Syscall(syscall.SYS_SEMGET, uintptr(LockKey), uintptr(1), uintptr(00666))
	var file *os.File
	//第一次运行需要初始化信号量
	if int(r1) < 0 {
		//创建或获取文件用于锁
		file, err2 = GetLockFile()
		defer func() {
			_ = file.Close() //忽略错误
		}()
		if err2 != nil {
			fmt.Printf("获取锁文件失败: %v\n", err2)
			return
		}
		//获取锁
		err2 = syscall.Flock(int(file.Fd()), syscall.LOCK_EX)
		if err2 != nil {
			fmt.Printf("文件锁Lock失败: %v\n", err2)
			return
		}
		//确保释放锁
		defer func() {
			err2 = syscall.Flock(int(file.Fd()), syscall.LOCK_UN)
			if err2 != nil {
				fmt.Printf("文件锁Unlock失败: %v\n", err2)
				return
			}
		}()
		//二次验证
		r1, _, err = syscall.Syscall(syscall.SYS_SEMGET, uintptr(LockKey), uintptr(1), uintptr(00666))
		if int(r1) < 0 {
			//初始化信号量
			r1, _, err = SetSemaphore()
			if int(r1) < 0 {
				fmt.Printf("信号量初始化失败: %v\n", err)
				return
			}
			fmt.Printf("初始化信号量成功: %d\n", SemShow(int(r1)))
		}
	}
	return
}

func SemLock(semId int, index int) (file *os.File, err syscall.Errno, err2 error) {
	stSemBuf := C.sembuf{
		sem_num: 0,
		sem_op:  -1,
		sem_flg: C.SEM_UNDO,
	}
	//获取读锁
	file, err2 = GetLockFile()
	if err2 != nil {
		fmt.Printf("获取读锁文件失败: %v\n", err2)
		return
	}
	err2 = syscall.Flock(int(file.Fd()), syscall.LOCK_SH)
	if err2 != nil {
		fmt.Printf("读锁Lock失败: %v\n", err2)
		return
	}
	var r1 uintptr
	//信号量-1
	for {
		r1, _, err = syscall.Syscall(syscall.SYS_SEMOP, uintptr(semId), uintptr(unsafe.Pointer(&stSemBuf)), 1)
		if err == syscall.EINTR {
			//系统调用被中断，重试
			fmt.Printf("线程%v系统调用被中断，重试\n", index)
			continue
		} else {
			break
		}
	}
	if int(r1) < 0 {
		fmt.Printf("请求信号量出错: %v\n", err)
		err2 = syscall.Flock(int(file.Fd()), syscall.LOCK_UN)
		return
	}
	return
}

func SemRelease(semId int, file *os.File, index int) (err syscall.Errno, err2 error) {
	stSemBuf := C.sembuf{
		sem_num: 0,
		sem_op:  1,
		sem_flg: C.SEM_UNDO,
	}
	defer func() {
		_ = file.Close()
	}()
	defer func() {
		//释放读锁
		if err2 != nil {
			fmt.Printf("获取读锁文件失败: %v\n", err2)
			return
		}
		err2 = syscall.Flock(int(file.Fd()), syscall.LOCK_UN)
		if err2 != nil {
			fmt.Printf("读锁释放失败: %v\n", err2)
			return
		}
	}()
	//信号量+1
	isSuccess, _, err := syscall.Syscall(syscall.SYS_SEMOP, uintptr(semId), uintptr(unsafe.Pointer(&stSemBuf)), 1)
	if int(isSuccess) < 0 {
		fmt.Printf("线程%v释放信号量出错: %v\n", index, err)
	}
	return
}

func SemShow(semId int) int {
	val, r2, err := syscall.Syscall(syscall.SYS_SEMCTL, uintptr(semId), 0, uintptr(C.GETVAL))
	if int(val) < 0 {
		fmt.Printf("信号量(id: %d)读值出错: %v, %v, %v\n", semId, val, r2, err)
	}
	return int(val)
}

func ReadConfig(filePath string) (config Config, err error) {
	configFile, err := os.Open(filePath)
	if err != nil {
		//log.Logger(fmt.Errorf("无法打开配置文件: %v\n", err))
		return
	}
	defer func() {
		_ = configFile.Close()
	}()
	configData, err := io.ReadAll(configFile)
	if err != nil {
		//log.Logger(fmt.Errorf("无法读取配置文件: %v\n", err))
		return
	}
	err = yaml.Unmarshal(configData, &config)
	if err != nil {
		//return config, fmt.Errorf("无法解析配置文件: %v", err)
		return
	}
	//加载配置文件
	return config, nil
}
