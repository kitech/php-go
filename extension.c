#include <stdio.h>
#include <stdlib.h>
#include <stdint.h>
#include "sztypes.h"

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

#define MFN 128
static zend_function_entry g_funcs[MFN] = {0};
static zend_module_entry g_entry = {0};

// ZEND_DECLARE_MODULE_GLOBALS(phpgo);
// static void init_globals(zend_phpgo_globals *globals) {}

int phpgo_module_startup_func(int type, int module_number TSRMLS_DC)
{
    printf("%s:%d %s\n", __FILE__, __LINE__, __FUNCTION__);
    return 0;
}

int phpgo_module_shutdown_func(int type, int module_number TSRMLS_DC)
{
    printf("%s:%d %s\n", __FILE__, __LINE__, __FUNCTION__);
    return 0;
}

int phpgo_request_startup_func(int type, int module_number TSRMLS_DC)
{
    printf("%s:%d %s\n", __FILE__, __LINE__, __FUNCTION__);
    return 0;
}

int phpgo_request_shutdown_func(int type, int module_number TSRMLS_DC)
{
    printf("%s:%d %s\n", __FILE__, __LINE__, __FUNCTION__);
    return 0;
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

void phpgo_function_handler(int fno, int ht, zval *return_value, zval **return_value_ptr,
                            zval *this_ptr, int return_value_used TSRMLS_DC)
{
    printf("function handler called.%d\n", fno);
    on_phpgo_function_callback(fno);
}

static void (*phpgo_handlers[MFN])(INTERNAL_FUNCTION_PARAMETERS) = {0};

#define PHPGO_FH_DECL(no) \
    void phpgo_function_handler_##no(INTERNAL_FUNCTION_PARAMETERS) {     \
        phpgo_function_handler(no, ht, return_value, return_value_ptr, this_ptr, return_value_used); \
    }
#define PHPGO_FH_ADD(no) phpgo_handlers[no] = phpgo_function_handler_##no

PHPGO_FH_DECL(0); PHPGO_FH_DECL(1); PHPGO_FH_DECL(2); PHPGO_FH_DECL(3); PHPGO_FH_DECL(4);
static int init_phpgo_handlers() {
    // phpgo_handlers[0] = phpgo_function_handler_0;
    PHPGO_FH_ADD(0); PHPGO_FH_ADD(1); PHPGO_FH_ADD(2); PHPGO_FH_ADD(3); PHPGO_FH_ADD(4);
    return 1;
}
static int init_phpgo_handlers_indicator = 0;

int zend_add_function(char *name)
{
    // TODO thread?
    if (init_phpgo_handlers_indicator == 0) {
        init_phpgo_handlers_indicator = 1;
        init_phpgo_handlers();
    }

    for (int i = 0; i < MFN; i ++) {
        zend_function_entry *fe = &g_funcs[i];

        if (fe->fname == NULL) {
            zend_function_entry e = {strdup(name), phpgo_handlers[i], NULL, 0, 0};
            memcpy(fe, &e, sizeof(e));
            printf("add func at %d\n", i);
            return i;
        } else {
            if (strcmp(name, fe->fname) == 0) {
                // dup function, override the old one
                return i;
            }
        }
    }

    // nospace
    return -1;
}

static zend_class_entry g_centries[MFN] = {0};
static zend_function_entry g_cmths[MFN][MFN] = {0};

void phpgo_method_handler(int ht, zval *return_value, zval **return_value_ptr,
                          zval *this_ptr, int return_value_used TSRMLS_DC)
{
    printf("method handler called.\n");
}

int zend_add_class(char *cname)
{
    cname = strdup(cname);
    int npos = -1;
    for (int i = 0; i < MFN; i++) {
        if (g_centries[i].name == NULL) {
            npos = i;
            break;
        }
    }

    if (npos == -1) {
        // nospace
        return npos;
    }
    printf("clsno: %d\n", npos);
    zend_class_entry *ce = &g_centries[npos];
    INIT_CLASS_ENTRY_EX((*ce), cname, strlen(cname), g_cmths[npos]);
    zend_register_internal_class(ce TSRMLS_CC);

    return npos;
}

int zend_add_method(char *cname, char *mname)
{
    int cpos = -1;
    for (int i = 0; i < MFN; i ++) {
        if (strcmp(g_centries[i].name, cname) == 0) {
            cpos = i;
            break;
        }
    }
    if (cpos == -1) {
        return -1;
    }

    for (int i = 0; i < MFN; i ++) {
        zend_function_entry *fe = &g_cmths[cpos][i];

        if (fe->fname == NULL) {
            zend_function_entry e = {mname, phpgo_method_handler, NULL, 0,
                                     0 | ZEND_ACC_PUBLIC};
            memcpy(fe, &e, sizeof(e));
            return i;
        } else {
            if (strcmp(mname, fe->fname) == 0) {
                // dup method, override the old one
                return i;
            }
        }
    }

    // nospace
    return -1;
}

