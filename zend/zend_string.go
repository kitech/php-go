package zend

/*
#include <zend_API.h>
#include <zend_string.h>

#ifdef ZEND_ENGINE_3
zend_string *zend_string_new_ownered(char *str, size_t len, int persistent) {
	zend_string *ret = zend_string_alloc(len, persistent);

	// memcpy(ZSTR_VAL(ret), str, len);
    memcpy(&ZSTR_VAL(ret), &str, sizeof(&str));
	ZSTR_VAL(ret)[len] = '\0';
	return ret;
}
char *zend_string_to_char(void *zs) { return &((zend_string*)zs)->val[0];}
#else
char *zend_string_to_char(void *zs) { return (char*)zs;}
#endif

*/
import "C"
import "unsafe"
import "reflect"
import "log"

// type ZString *C.zend_string

/*
func fromZString(s *C.zend_string) string {
	r := C.GoStringN(&s.val[0], C.int(s.len))
	return r
}
*/

func fromZString(s interface{}) string {
	if s == nil {
		return ""
	}

	sv := reflect.ValueOf(s)
	// log.Println(s, sv.Type().String(), sv.Type().Kind(), sv.Pointer())
	if sv.Type().Kind() != reflect.Ptr {
		log.Println("not supported type:", sv.Type().Kind())
		return ""
	}
	if sv.IsNil() {
		return ""
	}

	switch sv.Type().String() {
	case "*zend._Ctype_struct__zend_string":
		s_c := C.zend_string_to_char(unsafe.Pointer(sv.Pointer()))
		s_go := C.GoString(s_c)
		return s_go
	case "*zend._Ctype_char":
		s_go := C.GoString(s.(*C.char))
		return s_go
	}
	return ""
}

/*
func toZString(s string) *C.zend_string {
	persistent := 0
	s_ := C.CString(s)
	// defer C.free(s_)
	// r := C.zend_string_init(s_, C.size_t(len(s)), C.int(persistent))
	r := C.zend_string_new_ownered(s_, C.size_t(len(s)), C.int(persistent))
	return r
}
*/
