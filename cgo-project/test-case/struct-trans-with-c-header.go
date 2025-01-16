package test_case

/*
#cgo CFLAGS: -I../cpp
#cgo LDFLAGS: -L../cpp -lstruct_trans
#include "struct_trans.h"
*/
import "C"
import "fmt"

/*
#cgo CFLAGS 用于指定编译 C 代码时的标志：在 ../cpp 目录中查找头文件
#cgo LDFLAGS 用于指定链接 C 代码时的标志：在 ../cpp 目录中查找库文件, -lstruct_trans：指定要链接的库文件。-lstruct_trans 表示
链接名为 libstruct_trans.a（静态库）或 libstruct_trans.so（共享库）的文件
cd cpp	编译 cpp 文件夹中的 C 代码为静态库或共享库
gcc -c struct_trans.c -o struct_trans.o
ar rcs libstruct_trans.a struct_trans.o
*/
func StructTransWithCHeader() {
	test := C.Test{}
	fmt.Println(test.a)
	fmt.Println(test.b)
	fmt.Println(test._type)
	// fmt.Println(test.size) // 位数据
	fmt.Println(C.get_size(&test))
	fmt.Println(test.arr1[0])
	// fmt.Println(test.arr) // 零长数组无法直接访问
	fmt.Println(C.Test_arr2_helper(&test, 1))

	// test2 := C.struct_Test2{}
	// fmt.Println(test2.c) // 指定了特殊对齐规则的结构体，无法在 CGO 中访问

	color := C.Color(C.RED)
	switch color {
	case C.RED:
		fmt.Println("RED")
	case C.GREEN:
		fmt.Println("GREEN")
	}
}
