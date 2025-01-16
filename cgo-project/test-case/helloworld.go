package test_case

/*
#include <stdio.h>
*/
import "C"

// go tool cgo helloworld.go
func HelloWorld() {
	println("Hello, World!")
	C.puts(C.CString("Hello, World!\n"))
}
