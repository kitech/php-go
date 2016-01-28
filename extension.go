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
import "os"

// 一个程序只能创建一个扩展
// 所以使用全局变量也没有问题。
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

type FuncEntry struct {
	class  string
	method string
	fn     interface{}
	isctor bool
}

func NewFuncEntry(class string, method string, fn interface{}) *FuncEntry {
	return &FuncEntry{class, method, fn, false}
}

func (this *FuncEntry) Name() string {
	return this.class + "_" + this.method
}

func (this *FuncEntry) IsGlobal() bool {
	return this.class == "global"
}

func (this *FuncEntry) IsCtor() bool {
	return !this.IsGlobal() && this.isctor
}

func (this *FuncEntry) IsMethod() bool {
	return !this.IsGlobal() && !this.isctor
}

// 支持的函数类型为，
// 至少要是个函数或者方法
// 最多只能返回一个值
// 参数个数小于等于10
// 返回值类型，必须是以下4类，string, intx, floatx, bool
func (this *FuncEntry) IsSupported() bool {
	return true
}

type Extension struct {
	syms    map[string]int
	classes map[string]int
	cbs     map[int]*FuncEntry // cbid => callable callbak

	fidx int // = 0

	objs map[uintptr]interface{} // php's this => go's this

	// phpgo init function
	module_startup_func   func(int, int) int
	module_shutdown_func  func(int, int) int
	request_startup_func  func(int, int) int
	request_shutdown_func func(int, int) int

	me *C.zend_module_entry
	fe *C.zend_function_entry
}

// TODO 把entry位置与cbid分开，这样cbfunc就能够更紧凑了
func NewExtension() *Extension {
	syms := make(map[string]int, 0)
	classes := make(map[string]int, 0)
	cbs := make(map[int]*FuncEntry, 0)
	objs := make(map[uintptr]interface{}, 0)

	classes["global"] = 0 // 可以看作内置函数的类
	return &Extension{syms: syms, classes: classes, cbs: cbs, objs: objs}
}

// depcreated
func gencbid(cidx int, fidx int) int {
	return cidx*128 + fidx
}

func nxtcbid() int {
	return len(gext.syms)
}

func AddFunc(name string, f interface{}) error {
	fe := NewFuncEntry("global", name, f)
	sname := fe.Name()

	if _, has := gext.syms[sname]; !has {
		// TODO check f type

		cidx := 0
		fidx := gext.fidx
		// cbid := gencbid(0, fidx)
		cbid := nxtcbid()

		argtys := ArgTypes2Php(f)
		var cargtys *C.char = nil
		if argtys != nil {
			cargtys = C.CString(*argtys)
		}
		rety := RetType2Php(f)

		cname := C.CString(name)
		n := C.zend_add_function(C.int(cidx), C.int(fidx), C.int(cbid), cname, cargtys, C.int(rety))
		C.free(unsafe.Pointer(cname))
		if argtys != nil {
			C.free(unsafe.Pointer(cargtys))
		}

		if int(n) == 0 {
			gext.syms[sname] = cbid
			gext.cbs[cbid] = fe
			gext.fidx += 1
			return nil
		}
	}

	return errors.New("add func error.")
}

// 添加新类的时候，可以把类的公共方法全部导出吧
// 不用逐个方法添加，简单多了。
func AddClass(name string, ctor interface{}) error {

	if _, has := gext.classes[name]; !has {
		cidx := len(gext.classes)

		var n C.int = 0
		if int(n) == 0 {
			addCtor(cidx, name, ctor)
			addMethods(cidx, name, ctor)
		}

		fmt.Println("add class:", name, cidx)
		cname := C.CString(name)
		n = C.zend_add_class(C.int(cidx), cname)
		C.free(unsafe.Pointer(cname))

		if int(n) == 0 {
			gext.classes[name] = cidx
		}
		return nil
	}

	return errors.New("add class error.")
}

func addCtor(cidx int, cname string, ctor interface{}) {
	mname := "__construct"
	fidx := 0
	addMethod(cidx, fidx, cname, mname, ctor, true)
}

func addMethods(cidx int, cname string, ctor interface{}) {
	fty := reflect.TypeOf(ctor)
	cls := fty.Out(0)

	for idx := 0; idx < cls.NumMethod(); idx++ {
		mth := cls.Method(idx)
		addMethod(cidx, idx+1, cname, mth.Name, mth.Func.Interface(), false)
	}
}

func addMethod(cidx int, fidx int, cname string, mname string, fn interface{}, isctor bool) {
	// cidx := gext.classes[cname]
	// cbid := gencbid(cidx, fidx)
	cbid := nxtcbid()

	fe := NewFuncEntry(cname, mname, fn)
	fe.isctor = isctor

	argtys := ArgTypes2Php(fn)
	var cargtys *C.char = nil
	if argtys != nil {
		cargtys = C.CString(*argtys)
	}
	rety := RetType2Php(fn)

	ccname := C.CString(cname)
	cmname := C.CString(mname)

	mn := C.zend_add_method(C.int(cidx), C.int(fidx), C.int(cbid), ccname, cmname, cargtys, C.int(rety))
	C.free(unsafe.Pointer(ccname))
	C.free(unsafe.Pointer(cmname))
	if argtys != nil {
		C.free(unsafe.Pointer(cargtys))
	}

	if mn == 0 {
		gext.cbs[cbid] = fe
		gext.syms[fe.Name()] = cbid
	}
}

// TODO 如果比较多的话，可以摘出来，放在builtin.go中
// 内置函数注册，内置类注册。
func addBuiltins() {
	// nice fix exit crash bug.
	AddFunc("GoExit", func() { os.Exit(0) })
	AddFunc("GoGo", func() {})
	AddFunc("GoPanic", func() { panic("got") })
	AddFunc("GoRecover", func() { recover() })
}

// 注册php module 初始化函数
func RegisterInitFunctions(module_startup_func func(int, int) int,
	module_shutdown_func func(int, int) int,
	request_startup_func func(int, int) int,
	request_shutdown_func func(int, int) int) {

	gext.module_startup_func = module_startup_func
	gext.module_shutdown_func = module_shutdown_func
	gext.request_startup_func = request_startup_func
	gext.request_shutdown_func = request_shutdown_func

	tocip := func(f interface{}) unsafe.Pointer {
		return unsafe.Pointer(&f)
	}

	C.phpgo_register_init_functions(tocip(gext.module_startup_func), tocip(gext.module_shutdown_func),
		tocip(gext.request_startup_func), tocip(gext.request_shutdown_func))
}

//export on_phpgo_function_callback
func on_phpgo_function_callback(cbid int, phpthis uintptr,
	a0 uintptr, a1 uintptr, a2 uintptr, a3 uintptr, a4 uintptr,
	a5 uintptr, a6 uintptr, a7 uintptr, a8 uintptr, a9 uintptr) uintptr {

	args := []uintptr{a0, a1, a2, a3, a4, a5, a6, a7, a8, a9}
	if len(args) > 0 {
	}

	fmt.Println("go callback called:", cbid, phpthis, gext.cbs[cbid])
	fmt.Println("go callback called:", args, C.GoString((*C.char)(unsafe.Pointer(a1))))

	fe := gext.cbs[cbid]
	// fe.fn.(func())()

	// 根据方法原型中的参数个数与类型，从当前函数中的a0-a9中提取正确的值出来
	fval := reflect.ValueOf(fe.fn)
	argv := ArgValuesFromPhp(fe.fn, args)

	if fe.IsMethod() {
		if phpthis == 0 {
			panic("wtf")
		}
		if _, has := gext.objs[phpthis]; !has {
			panic("wtf")
		}
		gothis := gext.objs[phpthis]
		// argv = append([]reflect.Value{reflect.ValueOf(gothis)}, argv...)
		argv[0] = reflect.ValueOf(gothis)
	}

	outs := fval.Call(argv)
	ret := RetValue2Php(fe.fn, outs)
	fmt.Println("meta call ret:", outs, ret)

	if fe.IsCtor() {
		if phpthis == 0 {
			panic("wtf")
		}
		if _, has := gext.objs[phpthis]; has {
			panic("wtf")
		}
		gext.objs[phpthis] = outs[0].Interface()
	}

	return ret
}

// 比较通用的在C中调用函数的方法
// on_phpgo_function_callback是根据cbid来确定如何调用函数
// 该函数直接根据函数指定fp函数指针对应的函数。
//export call_golang_function
func call_golang_function(fp unsafe.Pointer, a0 uintptr, a1 uintptr, a2 uintptr, a3 uintptr, a4 uintptr,
	a5 uintptr, a6 uintptr, a7 uintptr, a8 uintptr, a9 uintptr) uintptr {

	fval := reflect.ValueOf(*(*interface{})(fp))
	if fval.Interface() == nil {
		panic("wtf")
	}

	args := []uintptr{a0, a1, a2, a3, a4, a5, a6, a7, a8, a9}
	if len(args) > 0 {
	}
	argv := ArgValuesFromPhp(fval.Interface(), args)
	if len(argv) > 0 {
	}

	outs := fval.Call(argv)
	ret := RetValue2Php(fval.Interface(), outs)

	return ret
}
