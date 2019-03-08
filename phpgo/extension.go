package phpgo

/*
#include "extension.h"
#include "../zend/compat.h"
#include "../zend/szend.h"

#include <php.h>
#include <zend_exceptions.h>
#include <zend_interfaces.h>
#include <zend_ini.h>
#include <zend_constants.h>
#include <SAPI.h>

*/
import "C"
import "unsafe"
import "reflect"
import "errors"
import "fmt"
import "log"
import "os"
import "strings"

import "github.com/kitech/php-go/zend"

// 一个程序只能创建一个扩展
// 所以使用全局变量也没有问题。
var (
	ExtName string = ""
	ExtVer  string = "1.0"
)

//export InitExtension
func InitExtension(name string, version string) int {
	ExtName = name
	if len(version) > 0 {
		ExtVer = version
	}

	return 0
}

var gext = NewExtension()

type FuncEntry struct {
	class  string
	method string
	fn     interface{}
	isctor bool
	isdtor bool
}

func NewFuncEntry(class string, method string, fn interface{}) *FuncEntry {
	return &FuncEntry{class, method, fn, false, false}
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

func (this *FuncEntry) IsDtor() bool {
	return !this.IsGlobal() && this.isdtor
}

func (this *FuncEntry) IsMethod() bool {
	return !this.IsGlobal() && !this.isctor && !this.isdtor
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

	objs   map[uintptr]interface{}        // php's this => go's this
	objs_p map[unsafe.Pointer]interface{} // php's this => go's this

	// phpgo init function
	module_startup_func   func(int, int) int
	module_shutdown_func  func(int, int) int
	request_startup_func  func(int, int) int
	request_shutdown_func func(int, int) int

	inis *zend.IniEntries // ini entries

	//
	me *C.zend_module_entry
	fe *C.zend_function_entry
}

// 与C中的同名结构体对应
type phpgo_callback_signature struct {
	argtys   [10]int8
	rety     int
	varidict int
}

// TODO 把entry位置与cbid分开，这样cbfunc就能够更紧凑了
func NewExtension() *Extension {
	syms := make(map[string]int, 0)
	classes := make(map[string]int, 0)
	cbs := make(map[int]*FuncEntry, 0)
	objs := make(map[uintptr]interface{}, 0)
	objs_p := make(map[unsafe.Pointer]interface{}, 0)

	classes["global"] = 0 // 可以看作内置函数的类

	this := &Extension{syms: syms, classes: classes, cbs: cbs,
		objs: objs, objs_p: objs_p}
	this.inis = zend.NewIniEntries()

	return this
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

		argtys := zend.ArgTypes2Php(f)
		var cargtys *C.char = nil
		if argtys != nil {
			cargtys = C.CString(*argtys)
		}
		rety := zend.RetType2Php(f)

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
// @param ctor 是该类的构造函数，原型 func NewClass(...) *Class
func AddClass(name string, ctor interface{}) error {

	if _, has := gext.classes[name]; !has {
		cidx := len(gext.classes)
		var n C.int = 0

		// must add begin add class
		if int(n) == 0 {
			addCtor(cidx, name, ctor)
			addDtor(cidx, name, ctor)
			addMethods(cidx, name, ctor)
		}

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

func addDtor(cidx int, cname string, ctor interface{}) {
	mname := "__destruct"
	fidx := 1
	addMethod(ctor, cidx, fidx, cname, mname, ctor, false, true)
}

func addCtor(cidx int, cname string, ctor interface{}) {
	mname := "__construct"
	fidx := 0
	addMethod(ctor, cidx, fidx, cname, mname, ctor, true, false)
}

func addMethods(cidx int, cname string, ctor interface{}) {
	fty := reflect.TypeOf(ctor)
	cls := fty.Out(0)

	for idx := 0; idx < cls.NumMethod(); idx++ {
		mth := cls.Method(idx)
		addMethod(ctor, cidx, idx+2, cname, mth.Name, mth.Func.Interface(), false, false)
	}
}

func addMethod(ctor interface{}, cidx int, fidx int, cname string, mname string, fn interface{}, isctor, isdtor bool) {
	// cidx := gext.classes[cname]
	// cbid := gencbid(cidx, fidx)
	cbid := nxtcbid()

	fe := NewFuncEntry(cname, mname, fn)
	fe.isctor = isctor
	fe.isdtor = isdtor

	argtys := zend.ArgTypes2Php(fn)
	var cargtys *C.char = nil
	if argtys != nil {
		cargtys = C.CString(*argtys)
	}

	isSelf := false
    methodRetType := reflect.TypeOf(fn)
    if methodRetType.NumOut() > 0 {
        classType := reflect.TypeOf(ctor).Out(0)
        isSelf = classType == methodRetType.Out(0)
    }

	var rety int
    if !isSelf {
    	rety = zend.RetType2Php(fn)
    } else {
        rety = zend.PHPTY_IS_SELF
    }

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

func validFunc(fn interface{}) bool {
	fty := reflect.TypeOf(fn)
	if fty.Kind() != reflect.Func {
		log.Panicln("What's that?", fty.Kind().String())
	}

	if fty.IsVariadic() {
		log.Panicln("Can't support variadic func.", fty.Kind().String())
	}

	for idx := 0; idx < fty.NumIn(); idx++ {
		switch fty.In(idx).Kind() {
		case reflect.Func:
			fallthrough
		case reflect.Array:
			fallthrough
		case reflect.Slice:
			fallthrough
		case reflect.Chan:
			fallthrough
		case reflect.Map:
			fallthrough
		default:
			log.Panicln("Can't support arg type:", idx, fty.In(idx).Kind().String())
		}
	}

	return true
}

/*
* @param namespace string optional
 */
func AddConstant(name string, val interface{}, namespace interface{}) error {
	if len(name) == 0 {
		return nil
	}

	if namespace != nil {
		log.Println("Warning, namespace parameter not supported now. omited.")
	}

	module_number := C.phpgo_get_module_number()
	modname := C.CString(strings.ToUpper(name))
	defer C.free(unsafe.Pointer(modname))

	if val != nil {
		valty := reflect.TypeOf(val)

		switch valty.Kind() {
		case reflect.String:
			v := val.(string)
			modval := C.CString(v)
			defer C.free(unsafe.Pointer(modval))

			C.zend_register_stringl_constant_compat(modname, C.size_t(len(name)), modval, C.size_t(len(v)),
				C.CONST_CS|C.CONST_PERSISTENT, C.int(module_number))
		case reflect.Int, reflect.Int32, reflect.Uint32, reflect.Int64, reflect.Uint64,
			reflect.Int8, reflect.Uint8:
			iv := reflect.ValueOf(val).Convert(reflect.TypeOf(int64(1))).Interface()
			C.zend_register_long_constant_compat(modname, C.size_t(len(name)), C.zend_long(iv.(int64)),
				C.CONST_CS|C.CONST_PERSISTENT, C.int(module_number))
		case reflect.Float32, reflect.Float64:
			fv := reflect.ValueOf(val).Convert(reflect.TypeOf(float64(1.0))).Interface()
			C.zend_register_double_constant_compat(modname, C.size_t(len(name)), C.double(fv.(float64)),
				C.CONST_CS|C.CONST_PERSISTENT, C.int(module_number))
		case reflect.Bool:
			v := val.(bool)
			var bv int8 = 1
			if v == false {
				bv = 0
			}
			C.zend_register_bool_constant_compat(modname, C.size_t(len(name)), C.zend_bool(bv),
				C.CONST_CS|C.CONST_PERSISTENT, C.int(module_number))
		default:
			err := fmt.Errorf("Warning, unsported constant value type: %v", valty.Kind().String())
			log.Println(err)
			return err
		}
	} else {
		C.zend_register_null_constant_compat(modname, C.size_t(len(name)),
			C.CONST_CS|C.CONST_PERSISTENT, C.int(module_number))
	}

	return nil
}

func AddIniVar(name string, value interface{}, modifiable bool,
	onModifier func(*zend.IniEntry, string, int) int,
	displayer func(*zend.IniEntry, int)) {
	ie := zend.NewIniEntryDef()
	ie.Fill(name, value, modifiable, nil, nil)
	ie.SetModifier(onModifier)
	ie.SetDisplayer(displayer)
	gext.inis.Add(ie)
}

// TODO 如果比较多的话，可以摘出来，放在builtin.go中
// 内置函数注册，内置类注册。
func addBuiltins() {
	// nice fix exit crash bug.
	var iret C.int = 0
	if iret = C.gozend_function_registered(C.CString("GoExit")); iret == C.int(0) {
		AddFunc("GoExit", func(code int) { os.Exit(code) })
	}
	if iret = C.gozend_function_registered(C.CString("GoGo")); iret == C.int(0) {
		AddFunc("GoGo", func(fn interface{}) { log.Println(fn) })
	}
	if iret = C.gozend_function_registered(C.CString("GoPanic")); iret == C.int(0) {
		AddFunc("GoPanic", func() { panic("got") })
	}
	if iret = C.gozend_function_registered(C.CString("GoRecover")); iret == C.int(0) {
		AddFunc("GoRecover", func() { recover() })
	}
	if iret = C.gozend_function_registered(C.CString("GoPrintln")); iret == C.int(0) {
		AddFunc("GoPrintln", func(p0 int, v interface{}) { log.Println(v, 123333) })
	}
}

// TODO init func with go's //export
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

	C.phpgo_register_init_functions(tocip(go_module_startup_func), tocip(gext.module_shutdown_func),
		tocip(gext.request_startup_func), tocip(gext.request_shutdown_func))
}

// the module_startup_func proxy
func go_module_startup_func(a0 int, a1 int) int {
	// for test begin
	modifier := func(ie *zend.IniEntry, newValue string, stage int) int {
		// log.Println(ie.Name(), newValue, stage)
		return 0
	}
	displayer := func(ie *zend.IniEntry, itype int) {
		// log.Println(ie.Name(), itype)
	}
	AddIniVar("phpgo.hehe_int", 123, true, modifier, displayer)
	AddIniVar("phpgo.hehe_bool", true, true, modifier, displayer)
	AddIniVar("phpgo.hehe_long", 123, true, modifier, displayer)
	AddIniVar("phpgo.hehe_string", "strval123", true, modifier, displayer)
	// for test end

	gext.inis.Register(a1)

	return gext.module_startup_func(a0, a1)
}

//
// 以整数类型传递go值类型的实现的回调方式
//export on_phpgo_function_callback
func on_phpgo_function_callback(cbid int, phpthis uintptr,
	a0 uintptr, a1 uintptr, a2 uintptr, a3 uintptr, a4 uintptr,
	a5 uintptr, a6 uintptr, a7 uintptr, a8 uintptr, a9 uintptr) uintptr {

	args := []uintptr{a0, a1, a2, a3, a4, a5, a6, a7, a8, a9}
	if len(args) > 0 {
	}

	log.Println("go callback called:", cbid, phpthis, gext.cbs[cbid])
	log.Println("go callback called:", args)

	fe := gext.cbs[cbid]
	// fe.fn.(func())()

	// 根据方法原型中的参数个数与类型，从当前函数中的a0-a9中提取正确的值出来
	fval := reflect.ValueOf(fe.fn)
	argv := zend.ArgValuesFromPhp(fe.fn, args)

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
	ret := zend.RetValue2Php(fe.fn, outs)
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

//
// 以指针类型传递go值类型的实现的回调方式
//export on_phpgo_function_callback_p
func on_phpgo_function_callback_p(cbid int, phpthis unsafe.Pointer,
	a0 unsafe.Pointer, a1 unsafe.Pointer, a2 unsafe.Pointer, a3 unsafe.Pointer, a4 unsafe.Pointer,
	a5 unsafe.Pointer, a6 unsafe.Pointer, a7 unsafe.Pointer, a8 unsafe.Pointer, a9 unsafe.Pointer,
	retpp *unsafe.Pointer, op unsafe.Pointer) {

	args := []unsafe.Pointer{a0, a1, a2, a3, a4, a5, a6, a7, a8, a9}
	if len(args) > 0 {
	}

	log.Println("go callback called:", cbid, phpthis, gext.cbs[cbid], op)
	log.Println("go callback called:", args)

	fe := gext.cbs[cbid]
	// fe.fn.(func())()

	if op == nil && !fe.IsGlobal() {
		panic("is not a class or a function")
	}

	// 根据方法原型中的参数个数与类型，从当前函数中的a0-a9中提取正确的值出来
	fval := reflect.ValueOf(fe.fn)
	argv := zend.ArgValuesFromPhp_p(fe.fn, args, fe.IsMethod())

	if fe.IsMethod() {
		zend.CHKNILEXIT(phpthis, "wtf")
		gothis, has := gext.objs_p[op]
		if !has {
			panic("wtf")
		}
		// argv = append([]reflect.Value{reflect.ValueOf(gothis)}, argv...)
		argv[0] = reflect.ValueOf(gothis)
	}

	outs := fval.Call(argv)
	ret := zend.RetValue2Php_p(fe.fn, outs)
	log.Println("meta call ret:", outs, ret)

	if fe.IsCtor() {
		zend.CHKNILEXIT(phpthis, "wtf")
		if _, has := gext.objs_p[op]; has {
			panic("wtf")
		}
		gext.objs_p[op] = outs[0].Interface()
	}

	if fe.IsDtor() {
		zend.CHKNILEXIT(phpthis, "wtf")
		if _, has := gext.objs_p[op]; !has {
			panic("wtf")
		}
		delete(gext.objs_p, op)
	}

	*retpp = ret
	// return ret
}

//
// 比较通用的在C中调用go任意函数的方法
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
	argv := zend.ArgValuesFromPhp(fval.Interface(), args)
	if len(argv) > 0 {
	}

	outs := fval.Call(argv)
	ret := zend.RetValue2Php(fval.Interface(), outs)

	return ret
}

//
// 比较通用的在C中调用go任意函数的方法（但参数是都指针形式的）
// 该函数直接根据函数指定fp函数指针对应的函数。
//export call_golang_function_p
func call_golang_function_p(fp unsafe.Pointer, a0 unsafe.Pointer, a1 unsafe.Pointer, a2 unsafe.Pointer,
	a3 unsafe.Pointer, a4 unsafe.Pointer, a5 unsafe.Pointer, a6 unsafe.Pointer,
	a7 unsafe.Pointer, a8 unsafe.Pointer, a9 unsafe.Pointer) unsafe.Pointer {

	return nil
}
