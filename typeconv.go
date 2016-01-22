package zend

// go类型转换为PHP类型

import (
	"fmt"
	"reflect"
)

const (
	PHPTY_IS_NULL         = 0
	PHPTY_IS_LONG         = 1
	PHPTY_IS_DOUBLE       = 2
	PHPTY_IS_BOOL         = 3
	PHPTY_IS_ARRAY        = 4
	PHPTY_IS_OBJECT       = 5
	PHPTY_IS_STRING       = 6
	PHPTY_IS_RESOURCE     = 7
	PHPTY_IS_CONSTANT     = 8
	PHPTY_IS_CONSTANT_AST = 9
	PHPTY_IS_CALLABLE     = 10

	PHPTY_IS_CONSTANT_TYPE_MASK    = 0x00f
	PHPTY_IS_CONSTANT_UNQUALIFIED  = 0x010
	PHPTY_IS_LEXICAL_VAR           = 0x020
	PHPTY_IS_LEXICAL_REF           = 0x040
	PHPTY_IS_CONSTANT_IN_NAMESPACE = 0x100
)

func ArgsGo2Php(f interface{}) (ptfs *string) {
	fty := reflect.TypeOf(f)
	if fty.Kind() != reflect.Func {
		return
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
		default:
			fmt.Println("wtf", fty.In(idx).String())
		}
	}
	fmt.Println("===", tfs, fty.NumIn())
	return &tfs
}

func RetGo2Php(f interface{}) (rety int) {
	fty := reflect.TypeOf(f)
	if fty.Kind() != reflect.Func {
		return
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
		}
	}

	return
}
