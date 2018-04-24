package zend

/*
#include <zend_API.h>
#include <zend_ini.h>

#ifdef ZEND_ENGINE_2
typedef zend_ini_entry zend_ini_entry_def;
typedef char zend_string;  // 为了让gozend_ini_modifier7在php5下能够编译过
#endif

extern int gozend_ini_modifier7(zend_ini_entry *entry, zend_string *new_value, void *mh_arg1, void *mh_arg2, void *mh_arg3, int stage);
extern int gozend_ini_modifier5(zend_ini_entry *entry, char *new_value, uint new_value_length, void *mh_arg1, void *mh_arg2, void *mh_arg3, int stage);
extern void gozend_ini_displayer(zend_ini_entry *ini_entry, int type);
*/
import "C"
import "unsafe"

import (
	"fmt"
	// "reflect"
	"log"
	"runtime"
)

type IniEntries struct {
	zies []C.zend_ini_entry_def
}

var zend_ini_entry_def_zero C.zend_ini_entry_def

func NewIniEntries() *IniEntries {
	this := &IniEntries{}
	this.zies = make([]C.zend_ini_entry_def, 1)
	this.zies[0] = zend_ini_entry_def_zero
	return this
}

func (this *IniEntries) Register(module_number int) int {
	r := C.zend_register_ini_entries(&this.zies[0], C.int(module_number))
	log.Println(r)
	return int(r)
}
func (this *IniEntryDef) Unregister(module_number int) {
	C.zend_unregister_ini_entries(C.int(module_number))
}

func (this *IniEntries) Add(ie *IniEntryDef) {
	this.zies[len(this.zies)-1] = ie.zie
	this.zies = append(this.zies, zend_ini_entry_def_zero)
}

type IniEntry struct {
	zie *C.zend_ini_entry
}

func newZendIniEntryFrom(ie *C.zend_ini_entry) *IniEntry {
	return &IniEntry{ie}
}
func (this *IniEntry) Name() string      { return fromZString(this.zie.name) }
func (this *IniEntry) Value() string     { return fromZString(this.zie.value) }
func (this *IniEntry) OrigValue() string { return fromZString(this.zie.orig_value) }

const (
	INI_USER   = int(C.ZEND_INI_USER)
	INI_PERDIR = int(C.ZEND_INI_PERDIR)
	INI_SYSTEM = int(C.ZEND_INI_SYSTEM)
)

type IniEntryDef struct {
	zie C.zend_ini_entry_def

	onModify  func(ie *IniEntry, newValue string, stage int) int
	onDisplay func(ie *IniEntry, itype int)
}

func NewIniEntryDef() *IniEntryDef {
	this := &IniEntryDef{}
	// this.zie = (*C.zend_ini_entry_def)(C.calloc(1, C.sizeof_zend_ini_entry_def))
	runtime.SetFinalizer(this, zendIniEntryDefFree)
	return this
}

func zendIniEntryDefFree(this *IniEntryDef) {
	if _, ok := iniNameEntries[C.GoString(this.zie.name)]; ok {
		delete(iniNameEntries, C.GoString(this.zie.name))
	}

	if this.zie.name != nil {
		C.free(unsafe.Pointer(this.zie.name))
	}
	if this.zie.value != nil {
		C.free(unsafe.Pointer(this.zie.value))
	}
}

func (this *IniEntryDef) Fill3(name string, defaultValue interface{}, modifiable bool,
	onModify func(), arg1, arg2, arg3 interface{}, displayer func()) {
	this.zie.name = C.CString(name)
	this.zie.modifiable = go2cBool(modifiable)
	// this.zie.orig_modifiable = go2cBool(modifiable)
	if ZEND_ENGINE == ZEND_ENGINE_3 {
		this.zie.on_modify = go2cfn(C.gozend_ini_modifier7)
	} else {
		this.zie.on_modify = go2cfn(C.gozend_ini_modifier5)
	}
	this.zie.displayer = go2cfn(C.gozend_ini_displayer)

	value := fmt.Sprintf("%v", defaultValue)
	this.zie.value = C.CString(value)

	if arg1 == nil {
		this.zie.mh_arg1 = nil
	}
	if arg2 == nil {
		this.zie.mh_arg2 = nil
	}
	if arg3 == nil {
		this.zie.mh_arg3 = nil
	}

	if ZEND_ENGINE == ZEND_ENGINE_3 {
		this.zie.name_length = C.uint32_t(len(name))
		this.zie.value_length = C.uint32_t(len(value))
	} else {
		// why need +1 for php5?
		// if not, zend_alter_ini_entry_ex:280行会出现zend_hash_find无结果失败
		this.zie.name_length = C.uint32_t(len(name) + 1)
		this.zie.value_length = C.uint32_t(len(value) + 1)
	}
	log.Println(name, len(name))

	iniNameEntries[name] = this
}

func (this *IniEntryDef) Fill2(name string, defaultValue interface{}, modifiable bool,
	onModify func(), arg1, arg2 interface{}, displayer func()) {
	this.Fill3(name, defaultValue, modifiable, onModify, arg1, arg2, nil, displayer)
}

func (this *IniEntryDef) Fill1(name string, defaultValue interface{}, modifiable bool,
	onModify func(), arg1 interface{}, displayer func()) {
	this.Fill3(name, defaultValue, modifiable, onModify, arg1, nil, nil, displayer)
}

func (this *IniEntryDef) Fill(name string, defaultValue interface{}, modifiable bool,
	onModify func(), displayer func()) {
	this.Fill3(name, defaultValue, modifiable, onModify, nil, nil, nil, displayer)
}

func (this *IniEntryDef) SetModifier(modifier func(ie *IniEntry, newValue string, state int) int) {
	this.onModify = modifier
}

func (this *IniEntryDef) SetDisplayer(displayer func(ie *IniEntry, itype int)) {
	this.onDisplay = displayer
}

var iniNameEntries = make(map[string]*IniEntryDef, 0)

// the new_value is really not *C.char, it's *C.zend_string
//export gozend_ini_modifier7
func gozend_ini_modifier7(ie *C.zend_ini_entry, new_value *C.zend_string, mh_arg1 unsafe.Pointer, mh_arg2 unsafe.Pointer, mh_arg3 unsafe.Pointer, stage C.int) C.int {
	// log.Println(ie, "//", new_value, stage, ie.modifiable)
	// log.Println(ie.orig_modifiable, ie.modified, fromZString(ie.orig_value))
	// log.Println(fromZString(new_value), fromZString(ie.name), fromZString(ie.value))
	if iedef, ok := iniNameEntries[fromZString(ie.name)]; ok {
		iedef.onModify(newZendIniEntryFrom(ie), fromZString(new_value), int(stage))
	} else {
		log.Println("wtf", fromZString(ie.name))
	}
	return 0
}

//export gozend_ini_modifier5
func gozend_ini_modifier5(ie *C.zend_ini_entry, new_value *C.char, new_value_length C.uint, mh_arg1 unsafe.Pointer, mh_arg2 unsafe.Pointer, mh_arg3 unsafe.Pointer, stage C.int) C.int {
	// log.Println(ie, "//", new_value, new_value_length, stage, ie.modifiable)
	// log.Println(ie.orig_modifiable, ie.modified, fromZString(ie.orig_value))
	// log.Println(fromZString(new_value), fromZString(ie.name), fromZString(ie.value))
	if iedef, ok := iniNameEntries[fromZString(ie.name)]; ok {
		iedef.onModify(newZendIniEntryFrom(ie), fromZString(new_value), int(stage))
	} else {
		log.Println("wtf", fromZString(ie.name))
	}
	return 0
}

//export gozend_ini_displayer
func gozend_ini_displayer(ie *C.zend_ini_entry, itype C.int) {
	log.Println(ie, itype)
	if iedef, ok := iniNameEntries[fromZString(ie.name)]; ok {
		iedef.onDisplay(newZendIniEntryFrom(ie), int(itype))
	} else {
		log.Println("wtf", fromZString(ie.name))
	}
}
