#include <stdlib.h>

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

#define MFN 128
static zend_function_entry g_funcs[MFN] = {0};
static zend_module_entry g_entry = {0};


// ZEND_DECLARE_MODULE_GLOBALS(phpgo);
// static void init_globals(zend_phpgo_globals *globals) {}

void phpgo_module_startup_func(int type, int module_number TSRMLS_DC)
{
    
}

void phpgo_module_shutdown_func(int type, int module_number TSRMLS_DC)
{
    
}

void phpgo_request_startup_func(int type, int module_number TSRMLS_DC)
{
    
}

void phpgo_request_shutdown_func(int type, int module_number TSRMLS_DC)
{
    
}

void *get_module_impl() {
    zend_module_entry te = {
        STANDARD_MODULE_HEADER,
        "phpgo",
        g_funcs,
        phpgo_module_startup_func,
        phpgo_module_shutdown_func,
        phpgo_request_startup_func,		/* Replace with NULL if there's nothing to do at request start */
        phpgo_request_shutdown_func,	/* Replace with NULL if there's nothing to do at request end */
        NULL,
        "1.0",
        STANDARD_MODULE_PROPERTIES
    };

    memcpy(&g_entry, &te, sizeof(zend_module_entry));
    return &g_entry;
}

void phpgo_function_handler(int ht, zval *return_value, zval **return_value_ptr,
                            zval *this_ptr, int return_value_used TSRMLS_DC)
{
    printf("handler called.\n");
}

void zend_add_function(char *name)
{
    for (int i = 0; i < MFN; i ++) {
        zend_function_entry *fe = &g_funcs[i];

        if (fe->fname == NULL) {
            zend_function_entry e = {name, phpgo_function_handler, NULL, 0, 0};
            memcpy(fe, &e, sizeof(e));
            break;
        } else {
            if (strcmp(name, fe->fname) == 0) {
                // dup function, override the old one
                break;
            }
        }
    }
}
