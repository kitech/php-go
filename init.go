package zend

/*
#include "extension.h"
*/
import "C"
import "unsafe"
import "runtime"
import "fmt"

//export get_module
func get_module() unsafe.Pointer {
	if len(ExtName) == 0 {
		panic("ext name not set.")
	}

	mod := C.get_module_impl()
	fmt.Printf("mod=%p\n", mod)

	return unsafe.Pointer(mod)
}

func init() {
	runtime.GOMAXPROCS(1)
}
