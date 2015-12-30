
#include <stdlib.h>

#ifdef HAVE_CONFIG_H
#include "config.h"
#endif
#include "zend_API.h"


int gozend_call_user_function(zval **object, char *func_name, zval *retval_ptr, int argc, zval **params)
{
    zval function_name;
    INIT_ZVAL(function_name);
    ZVAL_STRING(&function_name, func_name, 1);

    assert(retval_ptr != NULL);
    if (call_user_function(CG(function_table), object , &function_name,
                           retval_ptr, argc, params TSRMLS_CC) == SUCCESS) {
        return 0;
    }
    return -1;
}

int gozend_call_user_function_string_ex(char *func_name, char *str, char **retstr, zval **object)
{
    zval *args[9];
    MAKE_STD_ZVAL(args[0]);
    ZVAL_STRING(args[0], str, 1);

    zval retval_ptr;
    INIT_ZVAL(retval_ptr);

    int ret = gozend_call_user_function(object, func_name, &retval_ptr, 1, args);
    if (ret == 0) {
        char *cstr = NULL;
        int cslen = 0;

        if (Z_TYPE_P(&retval_ptr) == IS_STRING) {
            cslen = Z_STRLEN_P(&retval_ptr);
            cstr = estrndup(Z_STRVAL_P(&retval_ptr), cslen);
            *retstr = cstr;
        }
    }
    return ret;
}


char *gozend_call_user_function_string(char *func_name, char *str)
{
    char *res = NULL;
    int ret = gozend_call_user_function_string_ex(func_name, str, &res, NULL);
    return res;
}


void call_user_function_callback(char *data)
{
    zval *args[9];
    MAKE_STD_ZVAL(args[0]);
    ZVAL_STRING(args[0], data, 1);
    zend_uint argc = 0;

    zval retval_ptr;
    INIT_ZVAL(retval_ptr);

    char *func_name = "say_hello_123";
    zval function_name;
    INIT_ZVAL(function_name);
    ZVAL_STRING(&function_name, func_name, 1);

    void *cobj = NULL; /* no object */
    if (call_user_function(CG(function_table), cobj , &function_name,
                           &retval_ptr, argc, NULL TSRMLS_CC) == SUCCESS) {
        /* do something with retval_ptr here if you like */
    }

    /* don't forget to free the zvals */
    // zval_ptr_dtor(&retval_ptr);
    zval_dtor(&function_name);
}
