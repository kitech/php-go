package zend

// go类型转换为PHP类型

/*
#include <stdlib.h>

#include <zend.h>
#include "../zend/compat.h"
#include "../zend/customize.h"
*/
import "C"
import "unsafe"

import (
	"log"
	"reflect"
)

const (
	PHPTY_IS_NULL         = C.IS_NULL
	PHPTY_IS_LONG         = C.IS_LONG
	PHPTY_IS_DOUBLE       = C.IS_DOUBLE
	PHPTY_IS_BOOL         = C._IS_BOOL
	PHPTY_IS_TRUE         = C.IS_TRUE
	PHPTY_IS_FALSE        = C.IS_FALSE
	PHPTY_IS_ARRAY        = C.IS_ARRAY
	PHPTY_IS_OBJECT       = C.IS_OBJECT
	PHPTY_IS_STRING       = C.IS_STRING
	PHPTY_IS_RESOURCE     = C.IS_RESOURCE
	PHPTY_IS_CONSTANT     = C.IS_CONSTANT
	PHPTY_IS_CONSTANT_AST = C.IS_CONSTANT_AST
	PHPTY_IS_CALLABLE     = C.IS_CALLABLE
    PHPTY_IS_SELF         = C.IS_SELF
    PHPTY_IS_INTERFACE    = C.IS_ZVAL

	PHPTY_IS_CONSTANT_TYPE_MASK    = 0x00f
	PHPTY_IS_CONSTANT_UNQUALIFIED  = 0x010
	PHPTY_IS_LEXICAL_VAR           = 0x020
	PHPTY_IS_LEXICAL_REF           = 0x040
	PHPTY_IS_CONSTANT_IN_NAMESPACE = 0x100
)

func ArgTypes2Php(fn interface{}) (ptfs *string) {
	fty := reflect.TypeOf(fn)
	if fty.Kind() != reflect.Func {
		log.Panicln("What's that?", fty.Kind().String())
	}

	var tfs string

	for idx := 0; idx < fty.NumIn(); idx++ {
		switch fty.In(idx).Kind() {
		case reflect.String:
			tfs = tfs + "s"
		case reflect.Float64:
			fallthrough
		case reflect.Float32:
			tfs = tfs + "d"
		case reflect.Bool:
			tfs = tfs + "b"
		case reflect.Int64:
			fallthrough
		case reflect.Uint64:
			fallthrough
		case reflect.Int32:
			fallthrough
		case reflect.Uint32:
			fallthrough
		case reflect.Int:
			fallthrough
		case reflect.Uint:
			fallthrough
		case reflect.Int16:
			fallthrough
		case reflect.Uint16:
			fallthrough
		case reflect.Int8:
			fallthrough
		case reflect.Uint8:
			tfs = tfs + "l"
		case reflect.Ptr:
			tfs = tfs + "p"
		case reflect.Interface:
			tfs = tfs + "a" // Any/interface
		case reflect.Slice:
			tfs = tfs + "v" // vector
		case reflect.Map:
			tfs = tfs + "m"
		default:
			log.Panicln("Unsupported kind:", "wtf", fty.In(idx).String(),
				fty.In(idx).Kind(), reflect.TypeOf(fn).IsVariadic())
		}
	}
	log.Println("===", tfs, fty.NumIn())
	return &tfs
}

func RetType2Php(fn interface{}) (rety int) {
	fty := reflect.TypeOf(fn)
	if fty.Kind() != reflect.Func {
		log.Panicln("What's that?", fty.Kind().String())
	}

	if fty.NumOut() > 0 {
		switch fty.Out(0).Kind() {
		case reflect.String:
			rety = PHPTY_IS_STRING
		case reflect.Float64:
			fallthrough
		case reflect.Float32:
			rety = PHPTY_IS_DOUBLE
		case reflect.Bool:
			rety = PHPTY_IS_BOOL
		case reflect.Int64:
			fallthrough
		case reflect.Uint64:
			fallthrough
		case reflect.Int32:
			fallthrough
		case reflect.Uint32:
			fallthrough
		case reflect.Int:
			fallthrough
		case reflect.Uint:
			fallthrough
		case reflect.Int16:
			fallthrough
		case reflect.Uint16:
			fallthrough
		case reflect.Int8:
			fallthrough
		case reflect.Uint8:
			rety = PHPTY_IS_LONG
		case reflect.Ptr:
			rety = PHPTY_IS_RESOURCE
		case reflect.Slice, reflect.Array, reflect.Map:
			rety = PHPTY_IS_ARRAY
        case reflect.Interface:
            // 尝试支持返回interface
            rety = PHPTY_IS_INTERFACE
		default:
			log.Panicln("wtf", fty.Out(0).String(), fty.Out(0).Kind().String())
		}
	}

	return
}

func ArgValuesFromPhp(fn interface{}, args []uintptr) (argv []reflect.Value) {
	fty := reflect.TypeOf(fn)
	if fty.Kind() != reflect.Func {
		log.Panicln("What's that?", fty.Kind().String())
	}

	argv = make([]reflect.Value, 0)
	for idx := 0; idx < fty.NumIn(); idx++ {
		switch fty.In(idx).Kind() {
		case reflect.String:
			var arg = C.GoString((*C.char)(unsafe.Pointer(args[idx])))
			var v = reflect.ValueOf(arg).Convert(fty.In(idx))
			argv = append(argv, v)
		case reflect.Float64:
			fallthrough
		case reflect.Float32:
			var arg = (*C.double)(unsafe.Pointer(args[idx]))
			var v = reflect.ValueOf(*arg).Convert(fty.In(idx))
			argv = append(argv, v)
			C.free(unsafe.Pointer(args[idx]))
		case reflect.Bool:
			var arg = (C.int)(args[idx])
			if arg == 1 {
				argv = append(argv, reflect.ValueOf(true))
			} else {
				argv = append(argv, reflect.ValueOf(false))
			}
		case reflect.Int64:
			fallthrough
		case reflect.Uint64:
			fallthrough
		case reflect.Int32:
			fallthrough
		case reflect.Uint32:
			fallthrough
		case reflect.Int:
			fallthrough
		case reflect.Uint:
			fallthrough
		case reflect.Int16:
			fallthrough
		case reflect.Uint16:
			fallthrough
		case reflect.Int8:
			fallthrough
		case reflect.Uint8:
			var arg = (C.ulonglong)(args[idx])
			var v = reflect.ValueOf(arg).Convert(fty.In(idx))
			argv = append(argv, v)
		case reflect.Ptr:
			// 有可能是go类的this指针
			if idx == 0 {
				// 这里仅是设置一个点位符号，这个gothis指针的位置
				argv = append(argv, reflect.ValueOf(0))
			} else {
				// 不支持其他非this参数的指针
				panic("wtf")
			}
		case reflect.Interface:
			log.Println("Unsupported interface kind:",
				fty.In(idx).Kind().String(), fty.In(idx).String())
		default:
			log.Panicln("Unsupported kind:",
				fty.In(idx).Kind().String(), fty.In(idx).String())
		}
	}

	if len(argv) != fty.NumIn() {
		panic("wtf")
	}
	return
}

func ArgValuesFromPhp_p(fn interface{}, args []unsafe.Pointer, ismth bool) (argv []reflect.Value) {
	fty := reflect.TypeOf(fn)
	if fty.Kind() != reflect.Func {
		log.Panicln("What's that?", fty.Kind().String())
	}

	argv = make([]reflect.Value, fty.NumIn())

	aidx := 0
	for idx := 0; idx < fty.NumIn(); idx++ {
		if ismth && idx == 0 {
			continue
		}

		aty := fty.In(idx)
		aiv := FROMCIP(args[aidx])
		if aiv != nil {
			if !reflect.TypeOf(aiv).ConvertibleTo(aty) {
				log.Panicln("can't convert ", reflect.TypeOf(aiv).Kind().String(), aty.Kind().String(),
					aty.Elem().Kind(), reflect.TypeOf(aiv).AssignableTo(aty), reflect.TypeOf(aiv).Elem().Kind())
			}
			argv[idx] = reflect.ValueOf(aiv).Convert(aty)
		} else {
			log.Panicln("nil arg", idx, ismth, fty.Name())
		}

		aidx += 1
	}

	if len(argv) != fty.NumIn() {
		panic("wtf")
	}
	return
}

// TODO 多值返回的支持？
func RetValue2Php(fn interface{}, rvs []reflect.Value) (retv uintptr) {
	fty := reflect.TypeOf(fn)
	if fty.Kind() != reflect.Func {
		log.Panicln("What's that?", fty.Kind().String())
	}
	retv = 0

	if fty.NumOut() > 0 {
		switch fty.Out(0).Kind() {
		case reflect.String:
			// 需要reciever 释放内存
			retv = uintptr(unsafe.Pointer(C.CString(rvs[0].Interface().(string))))
		case reflect.Float64:
			fallthrough
		case reflect.Float32:
			// 需要reciever 释放内存
			var pdv *C.double = (*C.double)(C.malloc(8))
			*pdv = (C.double)(rvs[0].Interface().(float64))
			retv = uintptr(unsafe.Pointer(pdv))
		case reflect.Bool:
			var bv = rvs[0].Interface().(bool)
			if bv {
				retv = uintptr(1)
			} else {
				retv = uintptr(0)
			}
		case reflect.Int64:
			fallthrough
		case reflect.Uint64:
			fallthrough
		case reflect.Int32:
			fallthrough
		case reflect.Uint32:
			fallthrough
		case reflect.Int:
			fallthrough
		case reflect.Uint:
			fallthrough
		case reflect.Int16:
			fallthrough
		case reflect.Uint16:
			fallthrough
		case reflect.Int8:
			fallthrough
		case reflect.Uint8:
			var dty = reflect.TypeOf(uint64(0))
			var nv = rvs[0].Convert(dty).Interface().(uint64)
			retv = uintptr(nv)
		case reflect.Ptr:
			var nv = rvs[0].Pointer()
			retv = uintptr(nv)
		default:
			log.Panicln("Unsupported kind:", fty.Out(0).Kind().String())
		}
	}

	return
}

// TODO 多值返回的支持？
func RetValue2Php_p(fn interface{}, rvs []reflect.Value) (retv unsafe.Pointer) {
	fty := reflect.TypeOf(fn)
	if fty.Kind() != reflect.Func {
		log.Panicln("What's that?", fty.Kind().String())
	}

	if fty.NumOut() > 0 {
		rvr := rvs[0].Interface()
		return TOCIP(rvr)
	}
	return nil
}
