package test_case

/*
#include <stdint.h>
union SayHello {
 int Say;
 float Hello;
};
union SayHello init_sayhello(){
    union SayHello us;
    us.Say = 100;
    return us;
}
int SayHello_Say_helper(union SayHello * us){
    return us->Say;
}
*/
import "C"
import (
	"encoding/binary"
	"fmt"
	"unsafe"
)

func UnionTrans() {
	SayHello := C.init_sayhello()                              //在C语言中，union的大小是其最大成员的大小，如果float改成double，byte[4]就变成了byte[8]
	fmt.Println("C-helper ", C.SayHello_Say_helper(&SayHello)) // 通过C辅助函数

	buff := C.GoBytes(unsafe.Pointer(&SayHello), 4)
	Say2 := binary.LittleEndian.Uint32(buff)
	fmt.Println("binary ", Say2) // 从内存直接解码一个int32

	fmt.Println("unsafe modify ", *(*C.int)(unsafe.Pointer(&SayHello))) // 强制类型转换
}
