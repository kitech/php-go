#include <assert.h>
#include <stdio.h>
#include <stdlib.h>
#include <stdint.h>
#include <assert.h>

/**
 *  PHP includes
 */
#include <php.h>
#include <zend_exceptions.h>
#include <zend_interfaces.h>
#include <zend_ini.h>
#include <SAPI.h>
#include <zend_hash.h>
#include <zend_API.h>

#ifdef ZTS
#include "TSRM.h"
#endif

#include "_cgo_export.h"

// #include "sztypes.h"
#include "../zend/goapi.h"
#include "../zend/clog.h"
#include "objectmap.h"
#include "class.h"

#define GLOBAL_VCLASS_NAME "_PHPGO_GLOBAL_"
// #define GLOBAL_VCLASS_ID 0
#define MAX_ARG_NUM 10
#define MFN 128  // MAX_FUNC_NUM
#define MCN 128  // MAX_CLASS_NUM

// TODO PHP7支持
static zend_module_entry g_entry = {0};
static phpgo_object_map *g_class_map = NULL;

struct phpgo_callback_signature {
    char argtys[10];
    int retys;
    int varidict;
};

static char *phpgo_argtys[MCN*MFN] = {0};
static int  phpgo_retys[MCN*MFN] = {0};
static struct phpgo_callback_signature phpgo_cbsigs[MCN*MFN] = {0};

// ZEND_DECLARE_MODULE_GLOBALS(phpgo);
// static void init_globals(zend_phpgo_globals *globals) {}

static int(*phpgo_module_startup_cbfunc)(int, int) = 0;
static int(*phpgo_module_shutdown_cbfunc)(int, int) = 0;
static int(*phpgo_request_startup_cbfunc)(int, int) = 0;
static int(*phpgo_request_shutdown_cbfunc)(int, int) = 0;

int phpgo_module_startup_func(int type, int module_number TSRMLS_DC)
{
    if (phpgo_module_startup_cbfunc) {
        return call_golang_function(phpgo_module_startup_cbfunc, type, module_number, 0, 0, 0, 0, 0, 0, 0, 0);
    }
    return 0;
}

int phpgo_module_shutdown_func(int type, int module_number TSRMLS_DC)
{
    if (phpgo_module_shutdown_cbfunc) {
        return call_golang_function(phpgo_module_shutdown_cbfunc, type, module_number, 0, 0, 0, 0, 0, 0, 0, 0);
    }
    return 0;
}

int phpgo_request_startup_func(int type, int module_number TSRMLS_DC)
{
    if (phpgo_request_startup_cbfunc) {
        return call_golang_function(phpgo_request_startup_cbfunc, type, module_number, 0, 0, 0, 0, 0, 0, 0, 0);
    }
    return 0;
}

int phpgo_request_shutdown_func(int type, int module_number TSRMLS_DC)
{
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
    char *cname = (char*)GLOBAL_VCLASS_NAME;
    phpgo_class_entry* pce = (phpgo_class_entry*)phpgo_object_map_get(g_class_map, cname);
    zend_function_entry* funcs = phpgo_class_get_funcs(pce);

    zend_module_entry te = {
        STANDARD_MODULE_HEADER,
        name, // "phpgo",
        funcs, // g_funcs[GLOBAL_VCLASS_ID],
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

int phpgo_get_module_number() {
    return g_entry.module_number;
}


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
#ifdef ZEND_ENGINE_3
    case _IS_BOOL:
#else
    case IS_BOOL:
#endif
        return "bool";
    case IS_LONG:
        return "long";
#ifdef ZEND_ENGINE_3
    case IS_UNDEF:
#endif
    case IS_NULL:
        return "void";
    case IS_OBJECT:
        return "object";
    case IS_RESOURCE:
        return "pointer";
    case IS_CALLABLE:
        return "callable";
    case IS_LEXICAL_VAR:
        return "var";
    default:
        return "unknown";
    }
    return NULL;
 }

// 计算要转换的参数个数，有效的参数个数
// php提供的参数个数要>=go函数需要的参数个数
// 如果go是变长参数个数， 则计算出 最小 值，并把php提供的所有参数全部传递
static int phpgo_function_num_args(int cbid, int supply_num_args)
{
    int num_args = phpgo_argtys[cbid] != NULL ? strlen(phpgo_argtys[cbid]) : 0;

    return num_args;
}

// PHP函数参数转换为go类型的参数
static void* phpgo_function_conv_arg(int cbid, int idx, char ch, int zty, zval *zarg)
{
    int prmty = zty;
    void *rv = NULL;

#ifdef ZEND_ENGINE_3
    zval *conv_zarg = zarg;
    zval *macro_zarg = zarg;
#else
    zval **conv_zarg = &zarg;
    zval *macro_zarg = zarg;
#endif

    if (ch == 's') {
        if (prmty != IS_STRING) {
            zend_error(E_WARNING, "Parameters type not match, need (%c), given: %s",
                       ch, type2name(prmty));
        }

        convert_to_string_ex(conv_zarg);
        char *arg = Z_STRVAL_P(macro_zarg);
        printf("arg%d(%c), %s(%d)\n", idx, ch, arg, Z_STRLEN_P(macro_zarg));
        goapi_new_value(GT_String, (uint64_t)arg, &rv);
    } else if (ch == 'l') {
        if (prmty != IS_LONG) {
            zend_error(E_WARNING, "Parameters type not match, need (%c), given: %s",
                       ch, type2name(prmty));
        }

        convert_to_long_ex(conv_zarg);
        long arg = (long)Z_LVAL_P(macro_zarg);
        printf("arg%d(%c), %d\n", idx, ch, arg);
        goapi_new_value(GT_Int64, (uint64_t)arg, &rv);
    } else if (ch == 'b') {
        convert_to_boolean_ex(conv_zarg);
        zend_bool arg = (zend_bool)Z_BVAL_P(macro_zarg);
        goapi_new_value(GT_Bool, (uint64_t)arg, &rv);
    } else if (ch == 'd') {
        convert_to_double_ex(conv_zarg);
        double arg = (double)Z_DVAL_P(macro_zarg);
        double* parg = calloc(1, sizeof(double));
        *parg = (double)Z_DVAL_P(macro_zarg);
        // argv[idx] = (void*)(ulong)arg;  // error
        // memcpy(&argv[idx], &arg, sizeof(argv[idx])); // ok
        // float32(uintptr(ax))
        goapi_new_value(GT_Float64, (uint64_t)parg, &rv);
    } else if (ch == 'a') {
        if (Z_TYPE_P(macro_zarg) == IS_STRING) {
            char *arg = Z_STRVAL_P(macro_zarg);
            goapi_new_value(GT_String, (uint64_t)arg, &rv);
#ifdef ZEND_ENGINE_3
        } else if (Z_TYPE_P(macro_zarg) == _IS_BOOL) {
#else
        } else if (Z_TYPE_P(macro_zarg) == IS_BOOL) {
#endif
            zend_bool arg = (zend_bool)Z_BVAL_P(macro_zarg);
            goapi_new_value(GT_Bool, (uint64_t)arg, &rv);
        } else if (Z_TYPE_P(macro_zarg) == IS_DOUBLE) {
            double* parg = calloc(1, sizeof(double));
            *parg = (double)Z_DVAL_P(macro_zarg);
            goapi_new_value(GT_Float64, (uint64_t)parg, &rv);
        } else if (Z_TYPE_P(macro_zarg) == IS_LONG) {
            long arg = (long)Z_LVAL_P(macro_zarg);
            goapi_new_value(GT_Int64, (uint64_t)arg, &rv);
        } else if (Z_TYPE_P(zarg) == IS_OBJECT) {
            void *parg = (void*)macro_zarg;
            goapi_new_value(GT_UnsafePointer, (uint64_t)parg, &rv);
        } else if (Z_TYPE_P(zarg) == IS_RESOURCE) {
            void *parg = (void*)macro_zarg;
            goapi_new_value(GT_UnsafePointer, (uint64_t)parg, &rv);
        } else {
            printf("arg%d(%c), %s(%d) unsported to Any\n", idx, ch,
                   type2name(Z_TYPE_P(macro_zarg)), Z_TYPE_P(macro_zarg));
        }

        // TODO array/vector convert
    } else if (ch == 'v') {
        HashTable *arr_hash = Z_ARRVAL_P(macro_zarg);
        HashPosition pos;
        zval *edata = NULL;

        for (zend_hash_internal_pointer_reset_ex(arr_hash, &pos);
#ifdef ZEND_ENGINE_3
             (edata = zend_hash_get_current_data_ex(arr_hash, &pos)) != NULL;
#else
             zend_hash_get_current_data_ex(arr_hash, (void**)&edata, &pos) == SUCCESS;
#endif
             zend_hash_move_forward_ex(arr_hash, &pos)) {

#ifdef ZEND_ENGINE_3
            zval *conv_edata = edata;
#else
            zval **conv_edata = &edata;
#endif

            if (rv == NULL) {
                switch (Z_TYPE_P(edata)) {
                case IS_LONG:
                    goapi_array_new(GT_Int64, &rv);
                    break;
                default:
                    goapi_array_new(GT_String, &rv);
                    break;
                }
            }

            switch (Z_TYPE_P(edata)) {
            case IS_LONG:
                goapi_array_push(rv, (void*)Z_LVAL_P(edata), &rv);
                break;
            default:
                convert_to_string_ex(conv_edata);
                printf("array idx(%d)=T(%d=%s, val=%s)\n", pos,
                       Z_TYPE_P(edata), type2name(Z_TYPE_P(edata)), Z_STRVAL_P(edata));
                goapi_array_push(rv, Z_STRVAL_P(edata), &rv);
                break;
            }
        }

        // TODO array/map convert
    } else if (ch == 'm') {
        // 简化成string=>string映射吧
        goapi_map_new(&rv);

        HashTable *arr_hash = Z_ARRVAL_P(macro_zarg);
        HashPosition pos;
        zval *edata;

        for (zend_hash_internal_pointer_reset_ex(arr_hash, &pos);
#ifdef ZEND_ENGINE_3
             (edata = zend_hash_get_current_data_ex(arr_hash, &pos)) != NULL;
#else
             zend_hash_get_current_data_ex(arr_hash, (void**)&edata, &pos) == SUCCESS;
#endif
             zend_hash_move_forward_ex(arr_hash, &pos)) {

            zval vkey = {0};
            zval *ekey = &vkey;
            zend_hash_get_current_key_zval_ex(arr_hash, &vkey, &pos);

#ifdef ZEND_ENGINE_3
            zval *conv_ekey = ekey;
            zval *conv_edata = edata;
#else
            zval **conv_ekey = &ekey;
            zval **conv_edata = &edata;
#endif

            switch (Z_TYPE_P(edata)) {
            default:
                convert_to_string_ex(conv_ekey);
                convert_to_string_ex(conv_edata);
                printf("array idx(%d)=K(%s)=T(%d=%s, val=%s)\n", pos, Z_STRVAL_P(ekey),
                       Z_TYPE_P(edata), type2name(Z_TYPE_P(edata)), Z_STRVAL_P(edata));
                goapi_map_add(rv, Z_STRVAL_P(ekey), Z_STRVAL_P(edata));
                break;
            }
        }

    } else {
        printf("arg%d(%c), %s(%d)\n", idx, ch, type2name(prmty), prmty);
        zend_error(E_WARNING, "Parameters type not match, need (%c), given: %s",
                   ch, type2name(prmty));
        return NULL;
    }

    return rv;
}

static void phpgo_function_conv_args(int cbid, int supply_num_args, void *argv[])
{
    int num_args = phpgo_function_num_args(cbid, supply_num_args);
    // int supply_num_args = ZEND_NUM_ARGS();

#ifdef ZEND_ENGINE_3
    zval args[MAX_ARG_NUM] = {0};
#else
    zval **args[MAX_ARG_NUM] = {0};
#endif
    if (zend_get_parameters_array_ex(num_args, args) == FAILURE) {
        printf("param count error: %d\n", num_args);
        WRONG_PARAM_COUNT;
        return;
    }

    // void *argv[MAX_ARG_NUM] = {0};
    dlog_debug("parse params: %d\n", num_args);
    // function has not this arg, don't -1
    for (int idx = 0; idx < num_args; idx ++) {
#ifdef ZEND_ENGINE_3
        zval *zarg = &args[idx];
#else
        zval *zarg = *(args[idx]);
#endif
        printf("arg%d, type=%d\n", idx, Z_TYPE_P(zarg));
        int prmty = Z_TYPE_P(zarg);
        char ch = phpgo_argtys[cbid][idx];

        argv[idx] = phpgo_function_conv_arg(cbid, idx, ch, prmty, zarg);
    }

    return;
}

static void phpgo_method_conv_args(int cbid, int supply_num_args, void *argv[])
{
    int num_args = phpgo_function_num_args(cbid, supply_num_args);
    // int supply_num_args = ZEND_NUM_ARGS();

#ifdef ZEND_ENGINE_3
    zval args[MAX_ARG_NUM] = {0};
#else
    zval **args[MAX_ARG_NUM] = {0};
#endif

    if (zend_get_parameters_array_ex(num_args-1, args) == FAILURE) {
        printf("param count error: %d\n", num_args);
        WRONG_PARAM_COUNT;
        return;
    }

    // void *argv[MAX_ARG_NUM] = {0};
    dlog_debug("parse params: %d\n", num_args);
    for (int idx = 0; idx < num_args-1; idx ++) {
#ifdef ZEND_ENGINE_3
        zval *zarg = &args[idx];
#else
        zval *zarg = *(args[idx]);
#endif

        printf("arg%d, type=%d\n", idx, Z_TYPE_P(zarg));
        int prmty = Z_TYPE_P(zarg);
        char ch = phpgo_argtys[cbid][idx+1];

        argv[idx] = phpgo_function_conv_arg(cbid, idx, ch, prmty, zarg);
    }

    return;
}

static void phpgo_function_reutrn_php_array(void *p0, zval *return_value) {
    array_init(return_value);

    goapi_set_php_array(p0, &return_value);
}

// go类型的返回值转换为PHP类型的变量值
static int phpgo_function_conv_ret(int cbid, void *p0, zval *return_value)
{
    RETVAL_LONG(5);

    printf("convert ret: %p, type: %d(%s)\n",
           p0, phpgo_retys[cbid], type2name(phpgo_retys[cbid]));

    uint64_t rv = (uint64_t)goapi_get_value(p0);
    // 返回值解析转换
    switch (phpgo_retys[cbid]) {
    case IS_STRING:
#ifdef ZEND_ENGINE_3
        RETVAL_STRINGL((char*)rv, strlen((char*)rv));
#else
        RETVAL_STRINGL((char*)rv, strlen((char*)rv), 1);
#endif
        free((char*)rv);
        break;
    case IS_DOUBLE:
        RETVAL_DOUBLE(*(double*)rv);
        free((double*)rv);
        break;
#ifdef ZEND_ENGINE_3
    case _IS_BOOL:
#else
    case IS_BOOL:
#endif
        RETVAL_BOOL(rv);
        break;
    case IS_LONG:
        RETVAL_LONG(rv);
        break;
    case IS_NULL:
        RETVAL_NULL();
        break;
#ifdef ZEND_ENGINE_3
    case IS_UNDEF:
        RETVAL_NULL();
        break;
#endif
    case IS_RESOURCE:
        RETVAL_NULL();
        break;
    case IS_ARRAY:
        phpgo_function_reutrn_php_array(p0, return_value);
        break;
    default:
        // wtf?
        zend_error(E_WARNING, "unrecognized return value: %d.", phpgo_retys[cbid]);
        break;
    }

    return 0;
}


#ifdef ZEND_ENGINE_3
void phpgo_function_handler7(int cbid, zend_execute_data *execute_data, zval *return_value)
{
    zval *this_ptr = &execute_data->This;  // always not null for php7
    zend_object* op = NULL;
    if (NULL != this_ptr && IS_OBJECT == execute_data->This.u1.v.type) {
        op = execute_data->This.value.obj;
    }
    zend_string *func_name = execute_data->func->common.function_name;
    char *class_name = ZEND_FN_SCOPE_NAME(execute_data->func);
    printf("function handler called.%d, this=%p, atys=%s, op=%p, func=%s, class=%s\n",
           cbid, this_ptr, phpgo_argtys[cbid], op, ZSTR_VAL(func_name), class_name);

    zend_function_entry *fe = phpgo_function_map_get(class_name, ZSTR_VAL(func_name));
    int cbid_dyn = phpgo_callback_map_get(class_name, ZSTR_VAL(func_name));
    dlog_debug("mapget: %p, %d=?%d\n", fe, cbid_dyn, cbid);
    cbid = cbid_dyn; // for transition

    void *argv[MAX_ARG_NUM] = {0};
    if (op == NULL) {
        phpgo_function_conv_args(cbid, (ZEND_NUM_ARGS()), argv);
    } else {
        phpgo_method_conv_args(cbid, (ZEND_NUM_ARGS()), argv);
    }

    void* rv = NULL;
    on_phpgo_function_callback_p(cbid, this_ptr, argv[0], argv[1],
                                 argv[2], argv[3], argv[4], argv[5],
                                 argv[6], argv[7], argv[8], argv[9], &rv, (void*) op);
    dlog_debug("inout ret:%p\n", rv);
    phpgo_function_conv_ret(cbid, rv, return_value);
}

void phpgo_function_handler(zend_execute_data *execute_data, zval *return_value)
{
    zend_string *func_name = execute_data->func->common.function_name;
    char *class_name = ZEND_FN_SCOPE_NAME(execute_data->func);

    int cbid = phpgo_callback_map_get(class_name, ZSTR_VAL(func_name));

    phpgo_function_handler7(cbid, execute_data, return_value);
}

#else  // php < 7.0
void phpgo_function_handler5(int cbid, int ht, zval *return_value, zval **return_value_ptr,
                            zval *this_ptr, int return_value_used TSRMLS_DC)
{
    void* op = NULL;
    if (NULL != this_ptr && IS_OBJECT == this_ptr->type) {
        zend_object_handle handle = this_ptr->value.obj.handle;
        struct _store_object *obj;
        obj = &EG(objects_store).object_buckets[handle].bucket.obj;
        op = &obj->object;
    }
    const char *func_name = get_active_function_name();
    const char *class_name = NULL;
    if (NULL != this_ptr) {
        zend_class_entry *ce = zend_get_class_entry(this_ptr);
        class_name = ce->name;
    }
    printf("function handler called.%d, this=%p, atys=%s, op=%p, func=%s, class=%s\n",
           cbid, this_ptr, phpgo_argtys[cbid], op, func_name, class_name);

    void *argv[MAX_ARG_NUM] = {0};
    if (this_ptr == NULL) {
        phpgo_function_conv_args(cbid, (ZEND_NUM_ARGS()), argv);
    } else {
        phpgo_method_conv_args(cbid, (ZEND_NUM_ARGS()), argv);
    }

    void* rv = NULL;
    on_phpgo_function_callback_p(cbid, this_ptr, argv[0], argv[1],
                                 argv[2], argv[3], argv[4], argv[5],
                                 argv[6], argv[7], argv[8], argv[9], &rv, (void*) op);
    dlog_debug("inout ret:%p\n", rv);
    phpgo_function_conv_ret(cbid, rv, return_value);
}

void phpgo_function_handler(int ht, zval *return_value, zval **return_value_ptr,
                             zval *this_ptr, int return_value_used TSRMLS_DC)
{
    const char *func_name = get_active_function_name(TSRMLS_C);
    const char *class_name = NULL;
    if (NULL != this_ptr) {
        zend_class_entry *ce = zend_get_class_entry(this_ptr);
        class_name = ce->name;
    }

    // object maybe destruct after script execute finished, and out of executed file scope
    if (!zend_is_executing()) {
        if (func_name == NULL && class_name != NULL) {
            func_name = "__destruct";
        }
    }
    int cbid = phpgo_callback_map_get(class_name, func_name);

    phpgo_function_handler5(cbid, ht, return_value, return_value_ptr, this_ptr, return_value_used);
}

void phpgo_function_handler_dep(int cbid, int ht, zval *return_value, zval **return_value_ptr,
                            zval *this_ptr, int return_value_used TSRMLS_DC)
{
    printf("function handler called.%d, this=%p, atys=%s\n",
           cbid, this_ptr, phpgo_argtys[cbid]);

    int num_args = phpgo_argtys[cbid] != NULL ? strlen(phpgo_argtys[cbid]) : 0;
    int supply_num_args = ZEND_NUM_ARGS();

    zval **args[10];
    if (zend_get_parameters_array_ex(num_args, args) == FAILURE) {
        printf("param count error: %d\n", num_args);
        WRONG_PARAM_COUNT;
        return;
    }

    dlog_debug("parse params: %d\n", num_args);
    void *argv[10] = {0};
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
        } else if (ch == 'l') {
            if (prmty != IS_LONG) {
                zend_error(E_WARNING, "Parameters type not match, need (%c), given: %s",
                           ch, type2name(prmty));
            }

            convert_to_long_ex(zarg);
            long arg = (long)Z_LVAL_PP(zarg);
            printf("arg%d(%c), %d\n", idx, ch, arg);
            argv[idx] = (void*)arg;
        } else if (ch == 'b') {
            convert_to_boolean_ex(zarg);
            zend_bool arg = (zend_bool)Z_BVAL_PP(zarg);
            argv[idx] = (void*)(long)arg;
        } else if (ch == 'd') {
            convert_to_double_ex(zarg);
            double arg = (double)Z_DVAL_PP(zarg);
            double* parg = calloc(1, sizeof(double));
            *parg = (double)Z_DVAL_PP(zarg);
            argv[idx] = parg;
            // argv[idx] = (void*)(ulong)arg;  // error
            // memcpy(&argv[idx], &arg, sizeof(argv[idx])); // ok
            // float32(uintptr(ax))
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
#ifdef ZEND_ENGINE_3
        RETVAL_STRINGL((char*)rv, strlen((char*)rv));
#else
        RETVAL_STRINGL((char*)rv, strlen((char*)rv), 1);
#endif
        free((char*)rv);
        break;
    case IS_DOUBLE:
        RETVAL_DOUBLE(*(double*)rv);
        free((double*)rv);
        break;
#ifdef ZEND_ENGINE_3
    case _IS_BOOL:
#else
    case IS_BOOL:
#endif
        RETVAL_BOOL(rv);
        break;
    case IS_LONG:
        RETVAL_LONG(rv);
        break;
    case IS_NULL:
        RETVAL_NULL();
        break;
#ifdef ZEND_ENGINE_3
    case IS_UNDEF:
        RETVAL_NULL();
        break;
#endif
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
#endif  // #ifdef ZEND_ENGINE_3


int zend_add_class(int cidx, char *cname)
{
    if (phpgo_object_map_get(g_class_map, cname) == NULL) {
        dlog_debug("Class %s not added.", cname);
        return -1;
    }

    phpgo_class_entry* pce = phpgo_object_map_get(g_class_map, cname);
    zend_class_entry *ce = (zend_class_entry*)phpgo_class_get(pce);
    INIT_CLASS_ENTRY_EX((*ce), cname, strlen(cname), phpgo_class_get_funcs(pce));
    zend_register_internal_class(ce TSRMLS_CC);

    phpgo_class_map_add(cname, ce);

    return 0;
}

int zend_add_class_not_register(int cidx, char *cname)
{
    if (phpgo_object_map_get(g_class_map, cname) != NULL) {
        dlog_debug("Class %s already added.", cname);
        return -1;
    }

    // TODO thread safe?
    phpgo_class_entry* pce = phpgo_class_new(cname);
    phpgo_object_map_add(&g_class_map, cname, pce);
    int class_count = phpgo_object_map_count(g_class_map);
    assert(cidx == class_count);

    zend_class_entry *ce = (zend_class_entry*)phpgo_class_get(pce);
    phpgo_class_map_add(cname, ce);

    return 0;
}

int zend_add_function(int cidx, int fidx, int cbid, char *name, char *atys, int rety)
{
    printf("add func %s at %d:%d=%d, atys=%s, rety=%d\n",
           name, cidx, fidx, cbid, atys, rety);

    int cnlen = strlen(GLOBAL_VCLASS_NAME);
    char *cname = GLOBAL_VCLASS_NAME;
    char *mname = name;

    phpgo_class_entry* pce = (phpgo_class_entry*)phpgo_object_map_get(g_class_map, cname);
    if (pce == NULL) {
        dlog_debug("pce empty: %d for %s", phpgo_object_map_count(g_class_map), cname);
        zend_add_class_not_register(cidx, cname);
        pce = (phpgo_class_entry*)phpgo_object_map_get(g_class_map, cname);
    }
    phpgo_class_method_add(pce, mname);
    phpgo_function_entry* pfe = phpgo_class_method_get(pce, mname);
    zend_function_entry* fe = (zend_function_entry*)phpgo_function_get(pfe);

    zend_function_entry e = {strdup(name), phpgo_function_handler, NULL, 0, 0};
    memcpy(fe, &e, sizeof(e));
    (&g_entry)->functions = phpgo_class_get_funcs(pce);

    phpgo_argtys[cbid] = atys == NULL ? NULL : strdup(atys);
    phpgo_retys[cbid] = rety;

    phpgo_function_map_add(NULL, name, fe);
    phpgo_callback_map_add(NULL, name, cbid);

    return 0;
}


int zend_add_method(int cidx, int fidx, int cbid, char *cname, char *mname, char *atys, int rety)
{
    printf("add mth %s::%s at %d:%d=%d, atys=%s, rety=%d\n",
           cname, mname, cidx, fidx, cbid, atys, rety);

    phpgo_class_entry* pce = (phpgo_class_entry*)phpgo_object_map_get(g_class_map, cname);
    if (pce == NULL) {
        dlog_debug("pce empty: %d", phpgo_object_map_count(g_class_map));
        zend_add_class_not_register(cidx, cname);
        pce = (phpgo_class_entry*)phpgo_object_map_get(g_class_map, cname);
    }
    phpgo_class_method_add(pce, mname);
    phpgo_function_entry* pfe = phpgo_class_method_get(pce, mname);
    zend_function_entry* fe = (zend_function_entry*)phpgo_function_get(pfe);

    phpgo_argtys[cbid] = atys == NULL ? NULL : strdup(atys);
    phpgo_retys[cbid] = rety;

    phpgo_function_map_add(cname, mname, fe);
    phpgo_callback_map_add(cname, mname, cbid);

    return 0;
}


