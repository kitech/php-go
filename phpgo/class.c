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
#include "uthash.h"
#include "objectmap.h"

#include "class.h"

struct _phpgo_class_entry {
    zend_class_entry *ce;
    // phpgo_function_entry **fes;
    const char *class_name;

    phpgo_object_map *fmap;
    zend_function_entry *fes;
};

phpgo_class_entry* phpgo_class_new(const char *class_name) {
    phpgo_class_entry* pce = (phpgo_class_entry*)calloc(1, sizeof(phpgo_class_entry));
    memset(pce, 0, sizeof(phpgo_class_entry));
    pce->ce = (zend_class_entry*)calloc(1, sizeof(zend_class_entry));
    pce->class_name = class_name;
    pce->fmap = NULL;

    return pce;
}

void phpgo_class_method_add(phpgo_class_entry* pce, const char *func_name) {
    phpgo_function_entry* pfe = phpgo_function_new(func_name);
    phpgo_object_map_add(&pce->fmap, func_name, pfe);

    int count = phpgo_object_map_count(pce->fmap);
    if (count == 1) {
        pce->fes = (zend_function_entry*)calloc(count + 1, sizeof(zend_function_entry));
    } else {
        pce->fes = (zend_function_entry*)realloc(pce->fes, (count + 1) * sizeof(zend_function_entry));
    }
    memset(&pce->fes[count], 0, sizeof(zend_function_entry));
    memcpy(&pce->fes[count-1], phpgo_function_get(pfe), sizeof(zend_function_entry));
    pce->ce->info.internal.builtin_functions = pce->fes;
}

zend_class_entry* phpgo_class_get(phpgo_class_entry* pce) {
    return pce->ce;
}

zend_function_entry* phpgo_class_get_funcs(phpgo_class_entry* pce) {
    return (zend_function_entry*)pce->ce->info.internal.builtin_functions;
}

phpgo_function_entry* phpgo_class_method_get(phpgo_class_entry* pce, const char *func_name) {
    phpgo_function_entry* pfe = (phpgo_function_entry*)phpgo_object_map_get(pce->fmap, func_name);
    if (pfe != NULL) {
        return pfe;
    }
    return NULL;
}

int phpgo_class_method_count(phpgo_class_entry* pce) {
    return phpgo_object_map_count(pce->fmap);
}

// function operations
struct _phpgo_function_entry {
    const char *func_name;
    zend_function_entry *fe;
};

#ifdef ZEND_ENGINE_3
extern void phpgo_function_handler(zend_execute_data *execute_data, zval *return_value);
#else
extern void phpgo_function_handler(int ht, zval *return_value, zval **return_value_ptr,
                                   zval *this_ptr, int return_value_used TSRMLS_DC);
#endif

phpgo_function_entry *phpgo_function_new(const char *func_name) {
    phpgo_function_entry *pce = (phpgo_function_entry*)calloc(1, sizeof(phpgo_function_entry));
    pce->func_name = func_name;
    pce->fe = (zend_function_entry*)calloc(1, sizeof(zend_function_entry));

    zend_function_entry *fe = pce->fe;
    zend_function_entry e = {strdup(func_name), phpgo_function_handler, NULL, 0, 0};
    memcpy(fe, &e, sizeof(e));

    return pce;
}

int phpgo_function_delete(phpgo_function_entry *pfe) {
    free((void*)(pfe->fe->fname));
    free(pfe->fe);
    free(pfe);
    return 0;
}

zend_function_entry* phpgo_function_get(phpgo_function_entry* pfe) {
    return pfe->fe;
}

//
#define MAX_ARG_NUM 10
struct _phpgo_callback_info {
    char arg_types[MAX_ARG_NUM];
    int ret_type;
    int varidict;
};

phpgo_callback_info* phpgo_callback_info_new(char *arg_types, int ret_type) {
    phpgo_callback_info* cbi = (phpgo_callback_info*)calloc(1, sizeof(phpgo_callback_info));
    strncpy(cbi->arg_types, arg_types, sizeof(cbi->arg_types));
    cbi->ret_type = ret_type;
    return cbi;
}

char* phpgo_callback_info_get_arg_types(phpgo_callback_info* cbi) {
    return cbi->arg_types;
}

int phpgo_callback_info_get_ret_type(phpgo_callback_info* cbi) {
    return cbi->ret_type;
}
