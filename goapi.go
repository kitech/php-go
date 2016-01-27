package zend

/*
 */
import "C"
import "reflect"
import "unsafe"

////
//export goapi_array_new
func goapi_array_new(kind int) interface{} {
	wkind := (reflect.Kind)(kind)

	var arr interface{} = nil

	switch wkind {
	case reflect.Invalid:
	case reflect.Bool:
		arr = make([]bool, 0)
	case reflect.Int:
		arr = make([]int, 0)
	case reflect.Int8:
		arr = make([]int8, 0)
	case reflect.Int16:
		arr = make([]int16, 0)
	case reflect.Int32:
		arr = make([]int32, 0)
	case reflect.Int64:
		arr = make([]int64, 0)
	case reflect.Uint:
		arr = make([]uint, 0)
	case reflect.Uint8:
		arr = make([]uint8, 0)
	case reflect.Uint16:
		arr = make([]uint16, 0)
	case reflect.Uint32:
		arr = make([]uint32, 0)
	case reflect.Uint64:
		arr = make([]uint64, 0)
	case reflect.Uintptr:
		arr = make([]uintptr, 0)
	case reflect.Float32:
		arr = make([]float32, 0)
	case reflect.Float64:
		arr = make([]float64, 0)
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
		arr = make([]string, 0)
	case reflect.Struct:
	case reflect.UnsafePointer:
		arr = make([]unsafe.Pointer, 0)
	}

	// arr := make([]string, 0)
	if 1 == reflect.String {

	}
	return reflect.ValueOf(arr).Interface()
}

//export goapi_array_push
func goapi_array_push(arr interface{}, elem interface{}) {
	// *arr = append(*arr, s)
}

//export goapi_map_new
func goapi_map_new() *map[string]string {
	m := make(map[string]string, 0)
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

func get_type_by_kind(kind reflect.Kind) reflect.Type {
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

	return refty
}

func new_value_by_kind(kind reflect.Kind) interface{} {
	refty := get_type_by_kind(kind)
	refval := reflect.New(refty)

	return refval
}
