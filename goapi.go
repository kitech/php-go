package zend

/*
#include <stdlib.h>
*/
import "C"
import "reflect"
import "unsafe"
import "log"

////
//export goapi_array_new
func goapi_array_new(kind int) unsafe.Pointer {
	sty := FROMCIP(goapi_type(kind)).(reflect.Type)
	log.Println(sty.Kind().String(), sty)
	arrval := reflect.MakeSlice(reflect.SliceOf(sty), 0, 0)
	return TOCIP(arrval.Interface())
}

//export goapi_array_push
func goapi_array_push(arrp unsafe.Pointer, elm unsafe.Pointer) (retarr unsafe.Pointer) {
	// *arr = append(*arr, s)
	arr := FROMCIP(arrp)
	vty := reflect.TypeOf(arr)
	if vty != nil {
	}

	switch vty.Elem().Kind() {
	case reflect.Int64:
		elmval := (int64)(C.int64_t(uintptr(elm)))
		narr := append(arr.([]int64), elmval)
		retarr = TOCIP(narr)
	case reflect.String:
		// only support string element
		elmval := C.GoString((*C.char)(elm))
		narr := append(arr.([]string), elmval)
		retarr = TOCIP(narr)
	}
	return
}

//export goapi_map_new
func goapi_map_new() *map[string]string {
	m := make(map[string]string, 0)
	kty := FROMCIP(goapi_type(int(reflect.String))).(reflect.Type)
	vty := FROMCIP(goapi_type(int(reflect.Int))).(reflect.Type)
	reflect.MakeMap(reflect.MapOf(kty, vty))

	return &m
}

//export goapi_map_add
func goapi_map_add(m *map[string]string, key string, value string) {
	(*m)[key] = value
}

//export goapi_map_get
func goapi_map_get(m *map[string]string, key string) *string {
	if _, has := (*m)[key]; has {
		v := (*m)[key]
		return &v
	}
	return nil
}

//export goapi_chan_new
func goapi_chan_new(kind int, buffer int) interface{} {
	cty := FROMCIP(goapi_type(kind)).(reflect.Type)
	chval := reflect.MakeChan(cty, buffer)

	return chval.Interface()
}

//export goapi_type
func goapi_type(kind int) unsafe.Pointer {
	wkind := (reflect.Kind)(kind)

	var refty reflect.Type

	switch wkind {
	case reflect.Invalid:
	case reflect.Bool:
		refty = reflect.TypeOf(true)
	case reflect.Int:
		refty = reflect.TypeOf(int(0))
	case reflect.Int8:
		refty = reflect.TypeOf(int8(0))
	case reflect.Int16:
		refty = reflect.TypeOf(int16(0))
	case reflect.Int32:
		refty = reflect.TypeOf(int32(0))
	case reflect.Int64:
		refty = reflect.TypeOf(int64(0))
	case reflect.Uint:
		refty = reflect.TypeOf(uint(0))
	case reflect.Uint8:
		refty = reflect.TypeOf(uint8(0))
	case reflect.Uint16:
		refty = reflect.TypeOf(uint16(0))
	case reflect.Uint32:
		refty = reflect.TypeOf(uint32(0))
	case reflect.Uint64:
		refty = reflect.TypeOf(uint64(0))
	case reflect.Uintptr:
		refty = reflect.TypeOf(uintptr(0))
	case reflect.Float32:
		refty = reflect.TypeOf(float32(1.0))
	case reflect.Float64:
		refty = reflect.TypeOf(float64(1.0))
	case reflect.Complex64:
	case reflect.Complex128:
	case reflect.Array:
		refty = reflect.TypeOf([]interface{}{})
	case reflect.Chan:
	case reflect.Func:
	case reflect.Interface:
	case reflect.Map:
		refty = reflect.TypeOf(map[interface{}]interface{}{})
	case reflect.Ptr:
	case reflect.Slice:
	case reflect.String:
		refty = reflect.TypeOf("")
	case reflect.Struct:
	case reflect.UnsafePointer:
		refty = reflect.TypeOf(unsafe.Pointer(uintptr(0)))
	}

	return TOCIP(refty)
}

//export goapi_typeof
func goapi_typeof(v unsafe.Pointer) unsafe.Pointer {
	gv := FROMCIP(v)
	return TOCIP(reflect.TypeOf(gv))
}

//export goapi_typeid
func goapi_typeid(v unsafe.Pointer) int {
	gv := FROMCIP(v)
	return int(reflect.TypeOf(gv).Kind())
}

//export goapi_new
func goapi_new(kind int) unsafe.Pointer {
	refty := goapi_type(kind)
	refval := reflect.New(FROMCIP(refty).(reflect.Type))

	return TOCIP(refval)
}
func goapi_new_kind(kind reflect.Kind) unsafe.Pointer {
	return goapi_new(int(kind))
}
func goapi_new_type(kty reflect.Type) unsafe.Pointer {
	return goapi_new_kind(kty.Kind())
}

//export goapi_new_value
func goapi_new_value(kind int, v uintptr) unsafe.Pointer {
	wkind := reflect.Kind(kind)

	var refval interface{}
	switch wkind {
	case reflect.Invalid:
	case reflect.Bool:
		if v == 0 {
			refval = false
		} else {
			refval = true
		}
	case reflect.Int:
		refval = int(v)
	case reflect.Int8:
		refval = int8(v)
	case reflect.Int16:
		refval = int16(v)
	case reflect.Int32:
		refval = int32(v)
	case reflect.Int64:
		refval = int64(v)
	case reflect.Uint:
		refval = uint(v)
	case reflect.Uint8:
		refval = uint8(v)
	case reflect.Uint16:
		refval = uint16(v)
	case reflect.Uint32:
		refval = uint32(v)
	case reflect.Uint64:
		refval = uint64(v)
	case reflect.Uintptr:
		refval = v
	case reflect.Float32:
		refval = float32(*(*C.float)(unsafe.Pointer(v)))
		C.free(unsafe.Pointer(v))
	case reflect.Float64:
		refval = float64(*(*C.double)(unsafe.Pointer(v)))
		C.free(unsafe.Pointer(v))
	case reflect.Complex64:
	case reflect.Complex128:
	case reflect.Array:
	case reflect.Chan:
	case reflect.Func:
	case reflect.Interface:
	case reflect.Map:
	case reflect.Ptr:
	case reflect.Slice:
	case reflect.String:
		refval = C.GoString((*C.char)(unsafe.Pointer(v)))
	case reflect.Struct:
	case reflect.UnsafePointer:
		refval = unsafe.Pointer(v)
	}

	return TOCIP(refval)
}

//export goapi_set_value
func goapi_set_value(gv unsafe.Pointer, v uintptr) unsafe.Pointer {
	giv := FROMCIP(gv)
	gvty := reflect.TypeOf(giv)
	nv := goapi_new_value(int(gvty.Kind()), v)
	reflect.ValueOf(giv).Set(reflect.ValueOf(FROMCIP(nv)))

	return gv
}

//export goapi_get
func goapi_get(gv unsafe.Pointer) uintptr {
	giv := FROMCIP(gv)
	gvty := reflect.TypeOf(giv)
	if gvty == nil {
		return 0
	}

	var rv uintptr

	switch gvty.Kind() {
	case reflect.Invalid:
	case reflect.Bool:
		if giv.(bool) {
			rv = 1
		} else {
			rv = 0
		}
	case reflect.Int:
		rv = (uintptr)(giv.(int))
	case reflect.Int8:
		rv = (uintptr)(giv.(int8))
	case reflect.Int16:
		rv = (uintptr)(giv.(int16))
	case reflect.Int32:
		rv = (uintptr)(giv.(int32))
	case reflect.Int64:
		rv = (uintptr)(giv.(int64))
	case reflect.Uint:
		rv = (uintptr)(giv.(uint))
	case reflect.Uint8:
		rv = (uintptr)(giv.(uint8))
	case reflect.Uint16:
		rv = (uintptr)(giv.(uint16))
	case reflect.Uint32:
		rv = (uintptr)(giv.(uint32))
	case reflect.Uint64:
		rv = (uintptr)(giv.(uint64))
	case reflect.Uintptr:
		rv = giv.(uintptr)
	case reflect.Float32:
	case reflect.Float64:
	case reflect.Complex64:
	case reflect.Complex128:
	case reflect.Array:
	case reflect.Chan:
	case reflect.Func:
	case reflect.Interface:
	case reflect.Map:
	case reflect.Ptr:
	case reflect.Slice:
	case reflect.String:
		rv = uintptr(unsafe.Pointer(C.CString(giv.(string))))
	case reflect.Struct:
	case reflect.UnsafePointer:
		rv = (uintptr)(giv.(unsafe.Pointer))
	}

	// 简洁方式
	rvty := FROMCIP(goapi_type(int(reflect.Uintptr))).(reflect.Type)
	if gvty.ConvertibleTo(rvty) {
		rv = reflect.ValueOf(giv).Convert(rvty).Interface().(uintptr)
	} else {
		switch gvty.Kind() {
		case reflect.Ptr:
			rv = reflect.ValueOf(giv).Pointer()
		default:
			log.Panicln("can not convert:", giv, gvty.Kind(), gvty, rvty)
		}
	}
	return rv
}
