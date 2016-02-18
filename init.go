package zend

/*
#include "extension.h"
*/
import "C"
import "unsafe"
import "runtime"

import "time"
import "math/rand"

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

func init() {
	omp := runtime.GOMAXPROCS(1)
	fmt.Println("run gozend init...", omp)
	omp = runtime.GOMAXPROCS(1)
	fmt.Println("run gozend init...", omp)

	rand.Seed(time.Now().UnixNano())

}
