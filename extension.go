package zend

/*
#include "extension.h"

#include <php.h>
#include <zend_exceptions.h>
#include <zend_interfaces.h>
#include <zend_ini.h>
#include <SAPI.h>

*/
import "C"

var (
	ExtName string = ""
	ExtVer  string = "1.0"
)

func InitExtension(name string, version ...string) {
	ExtName = name
	if len(version) > 0 {
		ExtVer = version[0]
	}
}

var gext = NewExtension()

type Extension struct {
	fns map[string]interface{}
	me  *C.zend_module_entry
	fe  *C.zend_function_entry
}

func NewExtension() *Extension {
	fns := make(map[string]interface{}, 0)

	return &Extension{fns: fns}
}

func AddFunc(name string, f interface{}) {
	if _, has := gext.fns[name]; !has {
		// TODO check f type
		gext.fns[name] = f
		C.zend_add_function(C.CString(name))
	}
}
