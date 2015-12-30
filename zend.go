package zend

// 用go封装Zend的C API函数。

import (
	"errors"
	// "fmt"
)

/*
#include <stdlib.h>
#cgo CFLAGS: -I/usr/include/php -I/usr/include/php/Zend -I/usr/include/php/TSRM
#cgo LDFLAGS: -lphp5

extern char* gozend_call_user_function_string(char *func_name, char *arg);
extern void call_user_function_callback(char *arg);
*/
import "C"
import "unsafe"

////export call_user_function_string
func Call_user_function_string(func_name string, arg string) (string, error) {
	cfname := C.CString(func_name)
	carg := C.CString(arg)
	defer C.free(unsafe.Pointer(cfname))
	defer C.free(unsafe.Pointer(carg))

	ret := C.gozend_call_user_function_string(cfname, carg)
	if ret == nil {
		return "", errors.New("call error")
	}

	defer C.free(unsafe.Pointer(ret))
	return C.GoString(ret), nil
}

func Call_user_function(func_name string, args ...interface{}) {
	argc := len(args)
	if argc > 0 {
	}

	// TODO args to zval**
}
