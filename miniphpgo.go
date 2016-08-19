package main

/*
// #include "extension.h"
#include <stdio.h>
#include <assert.h>
#include <stdio.h>
#include <stdlib.h>
#include <stdint.h>

#include <php.h>
	#include <zend_exceptions.h>
	#include <zend_interfaces.h>
	#include <zend_ini.h>
	#include <SAPI.h>
	#include <zend_hash.h>

#ifdef ZTS
#include "TSRM.h"
#endif

// #include "_cgo_export.h"

#define GLOBAL_VCLASS_ID 0
#define MAX_ARG_NUM 10
#define MFN 128  // MAX_FUNC_NUM
#define MCN 128  // MAX_CLASS_NUM

static int phpgo_module_startup_func(int type, int module_number TSRMLS_DC)
{
    return 0;
}

static int phpgo_module_shutdown_func(int type, int module_number TSRMLS_DC)
{
    return 0;
}

static int phpgo_request_startup_func(int type, int module_number TSRMLS_DC)
{
    return 0;
}

static int phpgo_request_shutdown_func(int type, int module_number TSRMLS_DC)
{
    return 0;
}

// TODO PHP7支持
static zend_function_entry g_funcs[MCN][MFN] = {0};
static zend_module_entry g_entry = {0};
static void *phpgo_get_module(char *name, char *version)
{
    zend_module_entry te = {
        STANDARD_MODULE_HEADER,
        name, // "phpgo",
        g_funcs[GLOBAL_VCLASS_ID],
        phpgo_module_startup_func,
        phpgo_module_shutdown_func,
        phpgo_request_startup_func,
        phpgo_request_shutdown_func,
        NULL,
        version, // "1.0",
        STANDARD_MODULE_PROPERTIES
    };

    memcpy(&g_entry, &te, sizeof(zend_module_entry));
    return &g_entry;
}
*/
import "C"
import "unsafe"

// import "runtime"

// import "os"
import "fmt"
import _ "runtime/cgo"

//
// 给php加载扩展时调用的扩展入口函数。
//export get_module
func get_module() unsafe.Pointer {
	fmt.Println("hhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhh")
	var ExtName = "hehoo"
	var ExtVer = "1.0"

	mod := C.phpgo_get_module(C.CString(ExtName), C.CString(ExtVer))
	fmt.Printf("mod=%p\n", mod)

	return unsafe.Pointer(mod)
}

func main() {}
