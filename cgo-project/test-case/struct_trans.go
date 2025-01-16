package test_case

/*
struct Test {
    int a;
    float b;
    double type;
    int size:10;
    int arr1[10];
    int arr2[];
};

int Test_arr2_helper(struct Test * tm ,int pos){
    return tm->arr2[pos];
}

// 辅助函数读取位字段
int get_size(struct Test *tm) {
    return tm->size;
}

// 辅助函数设置位字段
void set_size(struct Test *tm, int value) {
    tm->size = value;
}

enum Color{
	RED,
	GREEN,
	BLUE
};

// 指示编译器按1字节对齐方式对struct Test2进行内存对齐
#pragma  pack(1)
struct Test2 {
    float a;
    char b;
    int c;
};
*/
import "C"
import "fmt"

func StructTrans() {
	test := C.struct_Test{}
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

	color := C.enum_Color(C.RED)
	switch color {
	case C.RED:
		fmt.Println("RED")
	case C.GREEN:
		fmt.Println("GREEN")
	}
}
