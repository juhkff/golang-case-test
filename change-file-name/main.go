package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	dir, err := os.Executable()
	if err != nil {
		fmt.Println("获取目录失败:", err)
		return
	}
	dir = filepath.Dir(dir)
	change_name(dir)
	pause()
}

func change_name(dir string) {
	// 获取可执行文件所在的目录
	execPath, err := os.Executable()
	if err != nil {
		fmt.Println("获取可执行文件路径失败:", err)
		return
	}

	// 遍历目录下的所有文件
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// 获取执行文件自身的名字
		if filepath.Join(dir, info.Name()) == execPath || info.Name() == filepath.Base(dir) {
			// 跳过自身这次循环以及当前目录
			return nil
		}
		// 只处理文件，跳过目录
		if !info.IsDir() {
			// 获取文件名和扩展名
			fileName := info.Name()
			ext := filepath.Ext(fileName)
			baseName := strings.TrimSuffix(fileName, ext)

			// 新文件名
			newName := "icon-" + baseName

			// 构建完整的新路径
			newPath := filepath.Join(dir, newName)

			// 重命名文件
			err := os.Rename(path, newPath)
			if err != nil {
				fmt.Println("重命名文件失败:", err)
			} else {
				fmt.Println("重命名文件:", fileName, "->", newName)
			}
		} else {
			//递归遍历目录
			change_name(filepath.Join(dir, info.Name()))
		}
		return nil
	})

	if err != nil {
		fmt.Println("遍历目录失败:", err)
	}
}

func pause() {
	fmt.Println("按回车键退出...")
	fmt.Scanln()
}
