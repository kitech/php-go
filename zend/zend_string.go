package zend

/*
#include <zend_API.h>
#include <zend_string.h>

zend_string *zend_string_new_ownered(char *str, size_t len, int persistent) {
	zend_string *ret = zend_string_alloc(len, persistent);

	// memcpy(ZSTR_VAL(ret), str, len);
    memcpy(&ZSTR_VAL(ret), &str, sizeof(&str));
	ZSTR_VAL(ret)[len] = '\0';
	return ret;
}
*/
import "C"

type ZString *C.zend_string

func fromZString(s *C.zend_string) string {
	r := C.GoStringN(&s.val[0], C.int(s.len))
	return r
}

func toZString(s string) *C.zend_string {
	persistent := 0
	s_ := C.CString(s)
	// defer C.free(s_)
	// r := C.zend_string_init(s_, C.size_t(len(s)), C.int(persistent))
	r := C.zend_string_new_ownered(s_, C.size_t(len(s)), C.int(persistent))
	return r
}
