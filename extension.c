#include <assert.h>
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
#define MCN 128
static zend_function_entry g_funcs[MCN][MFN] = {0};
static zend_class_entry g_centries[MCN] = {0};
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
        g_funcs[0],
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

int zend_add_function(int cidx, int fidx, int cbid, char *name)
{
    // TODO thread?
    if (init_phpgo_handlers_indicator == 0) {
        init_phpgo_handlers_indicator = 1;
        init_phpgo_handlers();
    }

    printf("add func %s at %d:%d=%d, %p\n", name, cidx, fidx, cbid, phpgo_handlers[cbid]);

    if (phpgo_handlers[cbid] == NULL) {
        assert(phpgo_handlers[cbid] != NULL);
        return -1;
    }

    zend_function_entry *fe = &g_funcs[cidx][fidx];
    zend_function_entry e = {strdup(name), phpgo_handlers[cbid], NULL, 0, 0};
    memcpy(fe, &e, sizeof(e));

    return 0;
}

int zend_add_class(int cidx, char *cname)
{
    cname = strdup(cname);
    zend_class_entry *ce = &g_centries[cidx];
    INIT_CLASS_ENTRY_EX((*ce), cname, strlen(cname), g_funcs[cidx]);
    zend_register_internal_class(ce TSRMLS_CC);

    return 0;
}

int zend_add_method(int cidx, int fidx, int cbid, char *cname, char *mname)
{
    printf("add mth %s::%s at %d:%d=%d, %p\n", cname, mname, cidx, fidx, cbid, phpgo_handlers[cbid]);

    if (phpgo_handlers[cbid] == NULL) {
        assert(phpgo_handlers[cbid] != NULL);
        return -1;
    }

    zend_function_entry *fe = &g_funcs[cidx][fidx];
    zend_function_entry e = {mname, phpgo_handlers[cbid], NULL, 0,
                             0 | ZEND_ACC_PUBLIC};
    memcpy(fe, &e, sizeof(e));

    return 0;
}

