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
import "unsafe"
import "reflect"
import "errors"
import "fmt"

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
	fnames  map[string]interface{}
	fnos    map[int]interface{}
	classes map[string]int
	ctors   map[string]interface{}

	me *C.zend_module_entry
	fe *C.zend_function_entry
}

func NewExtension() *Extension {
	fnames := make(map[string]interface{}, 0)
	fnos := make(map[int]interface{}, 0)
	classes := make(map[string]int, 0)
	ctors := make(map[string]interface{}, 0)

	return &Extension{fnames: fnames, fnos: fnos, classes: classes, ctors: ctors}
}

func AddFunc(name string, f interface{}) error {
	if _, has := gext.fnames[name]; !has {
		// TODO check f type

		cname := C.CString(name)
		n := C.zend_add_function(cname)
		C.free(unsafe.Pointer(cname))

		if int(n) >= 0 {
			gext.fnos[int(n)] = f
			gext.fnames[name] = f
			return nil
		}
	}

	return errors.New("add func error.")
}

// 添加新类的时候，可以把类的公共方法全部导出吧
// 不用逐个方法添加，简单多了。
func AddClass(name string, ctor interface{}) error {

	if _, has := gext.classes[name]; !has {

		cname := C.CString(name)
		n := C.zend_add_class(cname)
		C.free(unsafe.Pointer(cname))

		if int(n) >= 0 {
			gext.classes[name] = int(n)
			gext.ctors[name] = ctor

			addMethods(name, ctor)
			return nil
		}
	}

	return errors.New("add class error.")
}

func addCtor(cname string, ctor interface{}) {
	mname := "__construct"
	addMethod(cname, mname)
}

func addMethods(cname string, ctor interface{}) {
	fty := reflect.TypeOf(ctor)
	cls := fty.Out(0)

	for i := 0; i < cls.NumMethod(); i++ {
		fmt.Println(i, cname, cls.Method(i).Name)
		addMethod(cname, cls.Method(i).Name)
	}
}

func addMethod(cname string, mname string) {
	ccname := C.CString(cname)
	cmname := C.CString(mname)

	mn := C.zend_add_method(ccname, cmname)
	C.free(unsafe.Pointer(ccname))
	C.free(unsafe.Pointer(cmname))

	if mn >= 0 {

	}
}

// 内置函数注册，内置类注册。
func addBuiltins() {
	// nice fix exit crash bug.
	// AddFunc("GoExit", func() { os.Exit(0) })
}

//export on_phpgo_function_callback
func on_phpgo_function_callback(no int) {
	fmt.Println("go callback called:", no, gext.fnos[no])
	f := gext.fnos[no]
	f.(func())()
}
