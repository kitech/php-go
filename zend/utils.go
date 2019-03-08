package zend

import "reflect"
import "unsafe"
import "log"

func CHKNILOMIT(v interface{}) {}
func CHKNILEXIT(v interface{}, args ...interface{}) {
	if v == nil {
		if len(args) > 0 {
			log.Panicln(args...)
		}
		panic(v)
	}
}

// 看来不可能实现，需要在调用处展开的
func CHKNILRET(v interface{}) (dfn func()) {

	dfn = func() {
		if r := recover(); r != nil {
		}
	}

	if v == nil || v.(bool) == false {
		panic(nil)
	}

	return
}

// 看来不可能实现，需要在调用处展开的
// 也许两个连用，接近实现return到上层函数的功能，还是很啰嗦。
func RETURN_IF_DECL(v interface{}) {
	if v == nil || v.(bool) == false {
		panic(nil)
	}
}

// defer RETURN_IF_EXEC()()
func RETURN_IF_EXEC(v interface{}) func() {
	return func() {
		if r := recover(); r != nil {
		}
	}
}

type Any interface{}

func init() {
	// log.SetFlags(log.Llongfile | log.LstdFlags)
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	log.SetPrefix("[phpgo] ")
}

func TOCIP(v interface{}) unsafe.Pointer {
	return unsafe.Pointer(&v)
}

func FROMCIP(p unsafe.Pointer) interface{} {
	if p == nil {
		return nil
	}

	rp := (*interface{})(p)
	if rp == nil || *rp == nil {
        return nil
    }
    
	return reflect.ValueOf(*rp).Interface()
}
