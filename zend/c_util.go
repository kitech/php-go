package zend

/*
 */
import "C"
import "unsafe"

// import "reflect"

func c2goBool(ok C.int) bool {
	if ok == 1 {
		return true
	}
	return false
}

func go2cBool(ok bool) C.int {
	if ok {
		return 1
	}
	return 0
}

//
type go2cfnty *[0]byte

// 参数怎么传递
func go2cfnp(fn unsafe.Pointer) *[0]byte {
	return go2cfnty(fn)
}
func go2cfn(fn interface{}) *[0]byte {
	// assert(reflect.TypeOf(fn).Kind == reflect.Ptrx)
	return go2cfnp(fn.(unsafe.Pointer))
}
