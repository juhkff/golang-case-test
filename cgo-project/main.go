package main

import (
	test_case "cgo-project/test-case"
	"fmt"
)

func main() {
	test_case.HelloWorld()
	fmt.Println("=====================================")
	test_case.StructTrans()
	fmt.Println("=====================================")
	/*
		执行以下代码需要先将该的方法所需要的C文件通过注释中的指令编译
		具体见该StructTransWithCHeader()方法的注释
	*/
	// test_case.StructTransWithCHeader()
	fmt.Println("=====================================")
	test_case.UnionTrans()
}
