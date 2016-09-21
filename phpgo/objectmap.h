#ifndef _OBJECT_HASH_H_
#define _OBJECT_HASH_H_

struct _zend_function_entry;
struct _zend_class_entry;
struct _zend_module_entry;

void phpgo_function_map_add(const char *class_name, const char *func_name, struct _zend_function_entry *fe);

struct _zend_function_entry* phpgo_function_map_get(const char *class_name, const char *func_name);

void phpgo_class_map_add(const char *class_name, struct _zend_class_entry *ce);

struct _zend_class_entry* phpgo_class_map_get(const char *class_name);

void phpgo_callback_map_add(const char *class_name, const char *func_name, int cbid);

int phpgo_callback_map_get(const char *class_name, const char *func_name);

#endif

