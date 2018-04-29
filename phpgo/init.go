package phpgo

/*
#include "extension.h"
*/
import "C"
import "unsafe"

//export get_module
func get_module() unsafe.Pointer {
	if len(ExtName) == 0 {
		panic("ext name not set.")
	}

	addBuiltins()

	mod := C.phpgo_get_module(C.CString(ExtName), C.CString(ExtVer))

	return unsafe.Pointer(mod)
}
