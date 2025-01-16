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
	test_case.StructTransWithCHeader()
	fmt.Println("=====================================")
	test_case.UnionTrans()
}
