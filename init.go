package phpgo

/*
#include "extension.h"
*/
import "C"
import "unsafe"

// import "runtime"

// import "os"
import "fmt"

//
// 给php加载扩展时调用的扩展入口函数。
//export get_module
func get_module() unsafe.Pointer {
	if len(ExtName) == 0 {
		panic("ext name not set.")
	}

	addBuiltins()

	mod := C.phpgo_get_module(C.CString(ExtName), C.CString(ExtVer))
	fmt.Printf("mod=%p\n", mod)

	return unsafe.Pointer(mod)
}
