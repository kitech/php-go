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
func goapi_array_new(kind int, retpp *unsafe.Pointer) {
	sty := FROMCIP(goapi_type_r(kind)).(reflect.Type)
	log.Println(sty.Kind().String(), sty)
	arrval := reflect.MakeSlice(reflect.SliceOf(sty), 0, 0)
	*retpp = TOCIP(arrval.Interface())
}

//export goapi_array_push
func goapi_array_push(arrp unsafe.Pointer, elm unsafe.Pointer, retpp *unsafe.Pointer) {
	// *arr = append(*arr, s)
	arr := FROMCIP(arrp)
	vty := reflect.TypeOf(arr)
	if vty != nil {
	}

	switch vty.Elem().Kind() {
	case reflect.Int64:
		elmval := (int64)(C.int64_t(uintptr(elm)))
		narr := append(arr.([]int64), elmval)
		*retpp = TOCIP(narr)
	case reflect.String:
		// only support string element
		elmval := C.GoString((*C.char)(elm))
		narr := append(arr.([]string), elmval)
		*retpp = TOCIP(narr)
	}
	return
}

//export goapi_map_new
func goapi_map_new(retpp *unsafe.Pointer) {
	// m := make(map[string]string, 0)
	kty := FROMCIP(goapi_type_r(int(reflect.String))).(reflect.Type)
	vty := FROMCIP(goapi_type_r(int(reflect.String))).(reflect.Type)
	mv := reflect.MakeMap(reflect.MapOf(kty, vty))
	m := mv.Interface().(map[string]string)

	*retpp = TOCIP(m)
}

//export goapi_map_add
func goapi_map_add(mp unsafe.Pointer, keyp unsafe.Pointer, valuep unsafe.Pointer) {
	m := FROMCIP(mp).(map[string]string)
	key := C.GoString((*C.char)(keyp))
	value := C.GoString((*C.char)(valuep))

	m[key] = value
}

//export goapi_map_get
func goapi_map_get(mp unsafe.Pointer, keyp unsafe.Pointer, retpp *unsafe.Pointer) {
	m := FROMCIP(mp).(map[string]string)
	key := C.GoString((*C.char)(keyp))

	if _, has := m[key]; has {
		v := m[key]
		*retpp = unsafe.Pointer(C.CString(v))
	}

	return
}

//export goapi_map_del
func goapi_map_del(mp unsafe.Pointer, keyp unsafe.Pointer) {
	m := FROMCIP(mp).(map[string]string)
	key := C.GoString((*C.char)(keyp))

	if _, has := m[key]; has {
		delete(m, key)
	}
}

//export goapi_map_has
func goapi_map_has(mp unsafe.Pointer, keyp unsafe.Pointer) bool {
	m := FROMCIP(mp).(map[string]string)
	key := C.GoString((*C.char)(keyp))

	if _, has := m[key]; has {
		return true
	}

	return false
}

//export goapi_chan_new
func goapi_chan_new(kind int, buffer int, retpp *unsafe.Pointer) {
	cty := FROMCIP(goapi_type_r(kind)).(reflect.Type)
	chval := reflect.MakeChan(cty, buffer)

	*retpp = TOCIP(chval.Interface())
}

// TODO
//export goapi_chan_read
func goapi_chan_read(chp unsafe.Pointer, retpp *unsafe.Pointer) {
	ch := FROMCIP(chp)
	chv := reflect.ValueOf(ch)

	rv, ok := chv.Recv()
	if ok {
		*retpp = TOCIP(rv.Interface())
	}
}

// TODO
//export goapi_chan_write
func goapi_chan_write(chp unsafe.Pointer, elm unsafe.Pointer) {
	ch := FROMCIP(chp)
	chv := reflect.ValueOf(ch)

	chv.Send(reflect.ValueOf(FROMCIP(elm)))
}

// TODO
//export goapi_chan_close
func goapi_chan_close(chp unsafe.Pointer) {
	ch := FROMCIP(chp)
	chv := reflect.ValueOf(ch)

	chv.Close()
}

func goapi_type_r(kind int) unsafe.Pointer {
	var retpp unsafe.Pointer
	goapi_type(kind, &retpp)
	return retpp
}

//export goapi_type
func goapi_type(kind int, retpp *unsafe.Pointer) {
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

	*retpp = TOCIP(refty)
}

//export goapi_typeof
func goapi_typeof(v unsafe.Pointer, retpp *unsafe.Pointer) {
	gv := FROMCIP(v)
	*retpp = TOCIP(reflect.TypeOf(gv))
}

//export goapi_typeid
func goapi_typeid(v unsafe.Pointer) int {
	gv := FROMCIP(v)
	return int(reflect.TypeOf(gv).Kind())
}

//export goapi_new
func goapi_new(kind int, retpp *unsafe.Pointer) {
	refty := goapi_type_r(kind)
	refval := reflect.New(FROMCIP(refty).(reflect.Type))

	*retpp = TOCIP(refval)
}
func goapi_new_kind(kind reflect.Kind, retpp *unsafe.Pointer) {
	goapi_new(int(kind), retpp)
}
func goapi_new_type(kty reflect.Type, retpp *unsafe.Pointer) {
	goapi_new_kind(kty.Kind(), retpp)
}

//export goapi_new_value
func goapi_new_value(kind int, v uintptr, retpp *unsafe.Pointer) {
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

	*retpp = TOCIP(refval)
}

//export goapi_set_value
func goapi_set_value(gv unsafe.Pointer, v uintptr, retpp *unsafe.Pointer) {
	giv := FROMCIP(gv)
	gvty := reflect.TypeOf(giv)
	var nv unsafe.Pointer
	goapi_new_value(int(gvty.Kind()), v, &nv)
	reflect.ValueOf(giv).Set(reflect.ValueOf(FROMCIP(nv)))

	if *retpp != gv {
		*retpp = gv
	}
}

//export goapi_get_value
func goapi_get_value(gv unsafe.Pointer) uintptr {
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
	rvty := FROMCIP(goapi_type_r(int(reflect.Uintptr))).(reflect.Type)
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
