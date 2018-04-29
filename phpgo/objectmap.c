
#include <stdio.h>

#include "objectmap.h"
#include "uthash.h"

#include <zend_API.h>

// static zend_function_entry g_funcs[MCN][MFN] = {0};
// static zend_class_entry g_centries[MCN] = {0};
// static zend_module_entry g_entry = {0};

typedef struct _phpgo_function_map {
    const char *name;
    zend_function_entry *fe;
    int id;
    UT_hash_handle hh;
} phpgo_function_map;

typedef struct _phpgo_class_map {
    const char *name;
    zend_class_entry *ce;
    int id;
    UT_hash_handle hh;
} phpgo_class_map;

typedef struct _phpgo_module_map {
    const char *name;
    zend_module_entry *me;
    int id;
    UT_hash_handle hh;
} phpgo_module_map;

typedef struct _phpgo_callback_map {
    const char *name;
    int cbid;
    UT_hash_handle hh;
} phpgo_callback_map;

struct _phpgo_object_map {
    const char *name;
    void *obj;
    UT_hash_handle hh;
};

static phpgo_function_map *g_funcs_map = NULL;
static phpgo_class_map *g_classes_map = NULL;
static phpgo_callback_map *g_callbacks_map = NULL;

void phpgo_function_map_add(const char *class_name, const char *func_name, zend_function_entry *fe)
{
    int id = UTHASH_CNT(hh, g_funcs_map);
    phpgo_function_map *m = (phpgo_function_map*)malloc(sizeof(phpgo_function_map));
    char *key = (char*)calloc(1, class_name == NULL ? 0 : strlen(class_name) + strlen(func_name) + 2 + 1);
    sprintf(key, "%s::%s%c", class_name == NULL ? "" : class_name, func_name, '\0');
    m->name = key;
    m->fe = fe;
    m->id = id;
    UTHASH_ADD_KEYPTR(hh, g_funcs_map, key, strlen(key), m);
}

zend_function_entry* phpgo_function_map_get(const char *class_name, const char *func_name)
{
    phpgo_function_map *m = NULL;
    char key[(class_name == NULL ? 0 : strlen(class_name)) + strlen(func_name) + 2 + 1];
    sprintf(key, "%s::%s%c", (class_name == NULL ? "" : class_name), func_name, '\0');
    UTHASH_FIND_STR(g_funcs_map, key, m);
    if (m == NULL) {
        return NULL;
    }
    return m->fe;
}

void phpgo_class_map_add(const char *class_name, zend_class_entry *ce)
{
    int id = UTHASH_CNT(hh, g_classes_map);
    phpgo_class_map *m = (phpgo_class_map*)malloc(sizeof(phpgo_class_map));
    m->name = class_name;
    m->ce = ce;
    m->id = id;
    UTHASH_ADD_KEYPTR(hh, g_classes_map, class_name, strlen(class_name), m);
}

zend_class_entry* phpgo_class_map_get(const char *class_name)
{
    phpgo_class_map *m = NULL;
    UTHASH_FIND_STR(g_classes_map, class_name, m);
    if (m == NULL) {
        return NULL;
    }

    return m->ce;
}

void phpgo_callback_map_add(const char *class_name, const char *func_name, int cbid)
{
    int id = UTHASH_CNT(hh, g_callbacks_map);
    phpgo_callback_map *m = (phpgo_callback_map*)malloc(sizeof(phpgo_callback_map));
    char *key = (char*)calloc(1, class_name == NULL ? 0 : strlen(class_name) + strlen(func_name) + 2 + 1);
    sprintf(key, "%s::%s%c", class_name == NULL ? "" : class_name, func_name, '\0');
    m->name = key;
    m->cbid = cbid;
    UTHASH_ADD_KEYPTR(hh, g_callbacks_map, key, strlen(key), m);
}

int phpgo_callback_map_get(const char *class_name, const char *func_name)
{
    phpgo_callback_map *m = NULL;
    char key[(class_name == NULL ? 0 : strlen(class_name)) + strlen(func_name) + 2 + 1];
    sprintf(key, "%s::%s%c", (class_name == NULL ? "" : class_name), func_name, '\0');
    UTHASH_FIND_STR(g_callbacks_map, key, m);
    if (m == NULL) {
        return -1;
    }
    return m->cbid;
}


phpgo_object_map* phpgo_object_map_new() {
    phpgo_object_map *om = (phpgo_object_map*)malloc(sizeof(phpgo_object_map));
    return om;
}

void phpgo_object_map_add(phpgo_object_map** om, const char *name, void* obj)
{
    phpgo_object_map *m = (phpgo_object_map*)malloc(sizeof(phpgo_object_map));
    char *key = strdup(name);
    m->name = key;
    m->obj = obj;
    UTHASH_ADD_KEYPTR(hh, *om, key, strlen(key), m);
}

void* phpgo_object_map_get(phpgo_object_map* om, const char *name)
{
    phpgo_object_map *m = NULL;
    const char *key = name;
    UTHASH_FIND_STR(om, key, m);
    if (m == NULL) {
        return NULL;
    }
    return m->obj;
}

int phpgo_object_map_count(phpgo_object_map* om) {
    return UTHASH_COUNT(om);
}

