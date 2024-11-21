package main

import (
	"fmt"
	"os"
	"syscall"
	"time"
	"unsafe"
)

var (
	user32         = syscall.NewLazyDLL("user32.dll")
	procKeybdEvent = user32.NewProc("keybd_event")
)

const (
	VK_F            = 0x46 // F 键的虚拟键码
	KEYEVENTF_KEYUP = 0x0002
)

func pressKey(keyCode byte) {
	// 按下键
	procKeybdEvent.Call(uintptr(keyCode), 0, 0, 0)
	// 松开键
	procKeybdEvent.Call(uintptr(keyCode), 0, KEYEVENTF_KEYUP, 0)
}

func isAdmin() bool {
	_, err := os.Open("\\\\.\\PHYSICALDRIVE0")
	return err == nil
}

func runAsAdmin() {
	verb := "runas"
	exe, err := os.Executable()
	if err != nil {
		fmt.Println("Failed to get executable path:", err)
		return
	}
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Failed to get current directory:", err)
		return
	}
	argv, _ := syscall.UTF16PtrFromString(exe)
	cwdPtr, _ := syscall.UTF16PtrFromString(cwd)
	verbPtr, _ := syscall.UTF16PtrFromString(verb)

	var showCmd int32 = 1 // SW_NORMAL
	ret, _, err := syscall.NewLazyDLL("shell32.dll").NewProc("ShellExecuteW").Call(
		0,
		uintptr(unsafe.Pointer(verbPtr)),
		uintptr(unsafe.Pointer(argv)),
		0,
		uintptr(unsafe.Pointer(cwdPtr)),
		uintptr(showCmd),
	)
	if ret <= 32 {
		fmt.Println("Failed to start as admin:", err)
	}
	os.Exit(0)
}

func main() {
	if !isAdmin() {
		runAsAdmin()
		return
	}

	for {
		pressKey(VK_F)
		time.Sleep(500 * time.Millisecond)
	}
}
