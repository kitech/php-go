package zend

// 用go封装Zend的C API函数。

import (
	"errors"
	// "fmt"
)

/*
// #cgo CFLAGS: -I/usr/include/php -I/usr/include/php/Zend -I/usr/include/php/TSRM
// #cgo LDFLAGS: -L/home/dev/php5/lib -lphp5
#cgo CFLAGS: -g -O2 -std=c99 -D_GNU_SOURCE
// #cgo LDFLAGS: -lphp5
#cgo linux LDFLAGS: -Wl,--warn-unresolved-symbols -Wl,--unresolved-symbols=ignore-all
#cgo darwin LDFLAGS: -undefined dynamic_lookup

#include <stdlib.h>
#include "sztypes.h"
#include "szend.h"

*/
import "C"
import "unsafe"
import "fmt"

const (
	ZEND_ENGINE_3 = 3
	ZEND_ENGINE_2 = 2
)

var ZEND_ENGINE = int(C.gozend_zend_version_no())

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

	fmt.Printf("retp: %p\n", ret)
	fmt.Println("ret:", C.GoString(ret))
	defer C.gozend_efree(unsafe.Pointer(ret))
	return C.GoString(ret), nil
}

func Call_user_function(func_name string, args ...interface{}) {
	argc := len(args)
	if argc > 0 {
	}

	// TODO args to zval**
}

func Eval_string(code string) bool {
	ccode := C.CString(code)
	defer C.free(unsafe.Pointer(ccode))

	ret := C.gozend_eval_string(ccode)
	return ret == 1
}
