#include <assert.h>
#include <stdio.h>
#include <stdlib.h>
#include <stdint.h>

/**
 *  PHP includes
 */
#include <php.h>
#include <zend_exceptions.h>
#include <zend_interfaces.h>
#include <zend_ini.h>
#include <SAPI.h>

#ifdef ZTS
#include "TSRM.h"
#endif

#include "_cgo_export.h"

// #include "sztypes.h"
#include "goapi.h"


#define GLOBAL_VCLASS_ID 0
#define MAX_ARG_NUM 10
#define MFN 128  // MAX_FUNC_NUM
#define MCN 128  // MAX_CLASS_NUM

// TODO PHP7支持
static zend_function_entry g_funcs[MCN][MFN] = {0};
static zend_class_entry g_centries[MCN] = {0};
static zend_module_entry g_entry = {0};

// ZEND_DECLARE_MODULE_GLOBALS(phpgo);
// static void init_globals(zend_phpgo_globals *globals) {}

static int(*phpgo_module_startup_cbfunc)(int, int) = 0;
static int(*phpgo_module_shutdown_cbfunc)(int, int) = 0;
static int(*phpgo_request_startup_cbfunc)(int, int) = 0;
static int(*phpgo_request_shutdown_cbfunc)(int, int) = 0;

int phpgo_module_startup_func(int type, int module_number TSRMLS_DC)
{
    printf("%s:%d %s\n", __FILE__, __LINE__, __FUNCTION__);
    if (phpgo_module_startup_cbfunc) {
        return call_golang_function(phpgo_module_startup_cbfunc, type, module_number, 0, 0, 0, 0, 0, 0, 0, 0);
    }
    return 0;
}

int phpgo_module_shutdown_func(int type, int module_number TSRMLS_DC)
{
    printf("%s:%d %s\n", __FILE__, __LINE__, __FUNCTION__);
    if (phpgo_module_shutdown_cbfunc) {
        return call_golang_function(phpgo_module_shutdown_cbfunc, type, module_number, 0, 0, 0, 0, 0, 0, 0, 0);
    }
    return 0;
}

int phpgo_request_startup_func(int type, int module_number TSRMLS_DC)
{
    printf("%s:%d %s\n", __FILE__, __LINE__, __FUNCTION__);
    if (phpgo_request_startup_cbfunc) {
        return call_golang_function(phpgo_request_startup_cbfunc, type, module_number, 0, 0, 0, 0, 0, 0, 0, 0);
    }
    return 0;
}

int phpgo_request_shutdown_func(int type, int module_number TSRMLS_DC)
{
    printf("%s:%d %s\n", __FILE__, __LINE__, __FUNCTION__);
    if (phpgo_request_shutdown_cbfunc) {
        return call_golang_function(phpgo_request_shutdown_cbfunc, type, module_number, 0, 0, 0, 0, 0, 0, 0, 0);
    }
    return 0;
}

void phpgo_register_init_functions(void *module_startup_func, void *module_shutdown_func,
                                   void *request_startup_func, void *request_shutdown_func)
{
    phpgo_module_startup_cbfunc = (int (*)(int, int))module_startup_func;
    phpgo_module_shutdown_cbfunc = (int (*)(int, int))module_shutdown_func;
    phpgo_request_startup_cbfunc = (int (*)(int, int))request_startup_func;
    phpgo_request_shutdown_cbfunc = (int (*)(int, int))request_shutdown_func;
    return;
}

// TODO module name from args
// TODO module version from args
// TODO module startup/shutdown... function from args, or set seperator
void *phpgo_get_module(char *name, char *version) {
    zend_module_entry te = {
        STANDARD_MODULE_HEADER,
        name, // "phpgo",
        g_funcs[GLOBAL_VCLASS_ID],
        phpgo_module_startup_func,
        phpgo_module_shutdown_func,
        phpgo_request_startup_func,		/* Replace with NULL if there's nothing to do at request start */
        phpgo_request_shutdown_func,	/* Replace with NULL if there's nothing to do at request end */
        NULL,
        version, // "1.0",
        STANDARD_MODULE_PROPERTIES
    };

    memcpy(&g_entry, &te, sizeof(zend_module_entry));
    return &g_entry;
}

static char *phpgo_argtys[MCN*MFN] = {0};
static int  phpgo_retys[MCN*MFN] = {0};

// TODO 支持php callable类型，方便实现回调？
// TODO 支持php array类型？有点复杂。
char *type2name(int type)
{
    switch (type) {
    case IS_ARRAY:
        return "array";
    case IS_STRING:
        return "string";
    case IS_DOUBLE:
        return "double";
    case IS_BOOL:
        return "bool";
    case IS_LONG:
        return "long";
    case IS_NULL:
        return "void";
    case IS_RESOURCE:
        return "pointer";
    default:
        return "unknown";
    }
    return NULL;
 }

// 计算要转换的参数个数，有效的参数个数
static int phpgo_function_num_args(int cbid)
{
    return 0;
}

// PHP函数参数转换为go类型的参数
static int phpgo_function_conv_args()
{

}

// go类型的返回值转换为PHP类型的变量值
static int phpgo_function_conv_ret()
{

}

void phpgo_function_handler(int cbid, int ht, zval *return_value, zval **return_value_ptr,
                            zval *this_ptr, int return_value_used TSRMLS_DC)
{
    printf("function handler called.%d, this=%p\n", cbid, this_ptr);

    int num_args = phpgo_argtys[cbid] != NULL ? strlen(phpgo_argtys[cbid]) : 0;
    int supply_num_args = ZEND_NUM_ARGS();

    zval **args[10];
    if (zend_get_parameters_array_ex(num_args, args) == FAILURE) {
        printf("param count error: %d\n", num_args);
        WRONG_PARAM_COUNT;
        return;
    }

    void *argv[10] = {0};
    printf("parse params: %d\n", num_args);
    for (int idx = 0; idx < num_args; idx ++) {
        printf("arg%d, type=%d\n", idx, Z_TYPE_PP(args[idx]));
        int prmty = Z_TYPE_PP(args[idx]);
        char ch = phpgo_argtys[cbid][idx];
        zval **zarg = args[idx];

        if (ch == 's') {
            if (prmty != IS_STRING) {
                zend_error(E_WARNING, "Parameters type not match, need (%c), given: %s",
                           ch, type2name(prmty));
            }

            convert_to_string_ex(zarg);
            char *arg = Z_STRVAL_PP(zarg);
            argv[idx] = arg;
            printf("arg%d(%c), %s(%d)\n", idx, ch, arg, Z_STRLEN_PP(zarg));
            void *gv = goapi_new_value(GT_String, argv[idx]);
        } else if (ch == 'l') {
            if (prmty != IS_LONG) {
                zend_error(E_WARNING, "Parameters type not match, need (%c), given: %s",
                           ch, type2name(prmty));
            }

            convert_to_long_ex(zarg);
            long arg = (long)Z_LVAL_PP(zarg);
            printf("arg%d(%c), %d\n", idx, ch, arg);
            argv[idx] = (void*)arg;
            void *gv = goapi_new_value(GT_Int64, argv[idx]);
        } else if (ch == 'b') {
            convert_to_boolean_ex(zarg);
            zend_bool arg = (zend_bool)Z_BVAL_PP(zarg);
            argv[idx] = (void*)(long)arg;
            void *gv = goapi_new_value(GT_Bool, argv[idx]);
        } else if (ch == 'd') {
            convert_to_double_ex(zarg);
            double arg = (double)Z_DVAL_PP(zarg);
            double* parg = calloc(1, sizeof(double));
            *parg = (double)Z_DVAL_PP(zarg);
            argv[idx] = parg;
            // argv[idx] = (void*)(ulong)arg;  // error
            // memcpy(&argv[idx], &arg, sizeof(argv[idx])); // ok
            // float32(uintptr(ax))
            void *gv = goapi_new_value(GT_Float64, argv[idx]);
        } else {
            printf("arg%d(%c), unknown\n", idx, ch);
            zend_error(E_WARNING, "Parameters type not match, need (%c), given: %s",
                       ch, type2name(prmty));
            return;
        }
    }

    uint64_t rv = on_phpgo_function_callback(cbid, (uint64_t)this_ptr,
                                             (uint64_t)argv[0], (uint64_t)argv[1],
                                             (uint64_t)argv[2], (uint64_t)argv[3],
                                             (uint64_t)argv[4], (uint64_t)argv[5],
                                             (uint64_t)argv[6], (uint64_t)argv[7],
                                             (uint64_t)argv[8], (uint64_t)argv[9]);

    // 返回值解析转换
    switch (phpgo_retys[cbid]) {
    case IS_STRING:
        RETVAL_STRINGL((char*)rv, strlen((char*)rv), 1);
        free((char*)rv);
        break;
    case IS_DOUBLE:
        RETVAL_DOUBLE(*(double*)rv);
        free((double*)rv);
        break;
    case IS_BOOL:
        RETVAL_BOOL(rv);
        break;
    case IS_LONG:
        RETVAL_LONG(rv);
        break;
    case IS_NULL:
        RETVAL_NULL();
        break;
    case IS_RESOURCE:
        RETVAL_NULL();
        break;
    default:
        // wtf?
        zend_error(E_WARNING, "unrecognized return value: %d.", phpgo_retys[cbid]);
        break;
    }
    // TODO cleanup?
}

static void (*phpgo_handlers[MCN*MFN])(INTERNAL_FUNCTION_PARAMETERS) = {0};


#define PHPGO_FH_DECL(no) \
    void phpgo_function_handler_##no(INTERNAL_FUNCTION_PARAMETERS) {     \
        phpgo_function_handler(no, ht, return_value, return_value_ptr, this_ptr, return_value_used); \
    }
#define PHPGO_FH_ADD(no) phpgo_handlers[no] = phpgo_function_handler_##no

// PHPGO_FH_DECL(0); PHPGO_FH_DECL(1); PHPGO_FH_DECL(2); PHPGO_FH_DECL(3); PHPGO_FH_DECL(4);
// PHPGO_FH_DECL(128); PHPGO_FH_DECL(129); PHPGO_FH_DECL(130); PHPGO_FH_DECL(131); PHPGO_FH_DECL(132);
#include "extension_syms.txt"
static int init_phpgo_handlers() {
    // phpgo_handlers[0] = phpgo_function_handler_0;
    // PHPGO_FH_ADD(0); PHPGO_FH_ADD(1); PHPGO_FH_ADD(2); PHPGO_FH_ADD(3); PHPGO_FH_ADD(4);
    // PHPGO_FH_ADD(128); PHPGO_FH_ADD(129); PHPGO_FH_ADD(130); PHPGO_FH_ADD(131); PHPGO_FH_ADD(132);
    #include "extension_adds.txt"
    return 1;
}
static int init_phpgo_handlers_indicator = 0;

int zend_add_function(int cidx, int fidx, int cbid, char *name, char *atys, int rety)
{
    // TODO thread?
    if (init_phpgo_handlers_indicator == 0) {
        init_phpgo_handlers_indicator = 1;
        init_phpgo_handlers();
    }

    printf("add func %s at %d:%d=%d, %p, atys=%s, rety=%d\n",
           name, cidx, fidx, cbid, phpgo_handlers[cbid], atys, rety);

    if (phpgo_handlers[cbid] == NULL) {
        assert(phpgo_handlers[cbid] != NULL);
        return -1;
    }

    zend_function_entry *fe = &g_funcs[cidx][fidx];
    zend_function_entry e = {strdup(name), phpgo_handlers[cbid], NULL, 0, 0};
    memcpy(fe, &e, sizeof(e));

    phpgo_argtys[cbid] = atys == NULL ? NULL : strdup(atys);
    phpgo_retys[cbid] = rety;

    return 0;
}

int zend_add_class(int cidx, char *cname)
{
    zend_class_entry *ce = &g_centries[cidx];
    INIT_CLASS_ENTRY_EX((*ce), cname, strlen(cname), g_funcs[cidx]);
    zend_register_internal_class(ce TSRMLS_CC);

    return 0;
}

int zend_add_method(int cidx, int fidx, int cbid, char *cname, char *mname, char *atys, int rety)
{
    printf("add mth %s::%s at %d:%d=%d, %p, atys=%s, rety=%d\n",
           cname, mname, cidx, fidx, cbid, phpgo_handlers[cbid], atys, rety);

    if (phpgo_handlers[cbid] == NULL) {
        assert(phpgo_handlers[cbid] != NULL);
        return -1;
    }

    zend_function_entry *fe = &g_funcs[cidx][fidx];
    zend_function_entry e = {strdup(mname), phpgo_handlers[cbid], NULL, 0,
                             0 | ZEND_ACC_PUBLIC};
    memcpy(fe, &e, sizeof(e));

    phpgo_argtys[cbid] = atys == NULL ? NULL : strdup(atys);
    phpgo_retys[cbid] = rety;

    return 0;
}

// TODO PHP class instance destroyed
