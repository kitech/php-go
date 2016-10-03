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
    phpgo_function_entry **fes;
    const char *class_name;
    int func_count;

    phpgo_object_map *fmap;
};

phpgo_class_entry* phpgo_class_new(const char *class_name) {
    phpgo_class_entry* pce = (phpgo_class_entry*)calloc(1, sizeof(phpgo_class_entry));
    memset(pce, 0, sizeof(phpgo_class_entry));
    pce->ce = (zend_class_entry*)calloc(1, sizeof(zend_class_entry));
    pce->fes = NULL;
    pce->class_name = class_name;
    pce->fmap = NULL;
    return pce;
}

void phpgo_class_method_add(phpgo_class_entry* pce, const char *func_name) {
    int curr_index = -1;

    if (pce->fes == NULL) {
        pce->fes = (phpgo_function_entry**)calloc(1, 1 * sizeof(phpgo_function_entry*));
        pce->func_count = 1;
    } else {
        pce->func_count += 1;
        pce->fes = (phpgo_function_entry**)realloc(pce->fes, pce->func_count * sizeof(phpgo_function_entry*));
    }

    curr_index = pce->func_count - 1;
    pce->fes[curr_index] = phpgo_function_new(func_name);

    phpgo_object_map_add(pce->fmap, func_name, pce->fes[curr_index]);
}

phpgo_function_entry* phpgo_class_method_get(phpgo_class_entry* pce, const char *func_name) {
    phpgo_function_entry* pfe = (phpgo_function_entry*)phpgo_object_map_get(pce->fmap, func_name);
    if (pfe != NULL) {
        return pfe;
    }
    return NULL;
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
    free(pfe);
    return 0;
}

zend_function_entry* phpgo_function_get(phpgo_function_entry* pfe) {
    return pfe->fe;
}
