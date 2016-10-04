#ifndef _OBJECT_HASH_H_
#define _OBJECT_HASH_H_

struct _zend_function_entry;
struct _zend_class_entry;
struct _zend_module_entry;

typedef struct _phpgo_object_map phpgo_object_map;

void phpgo_function_map_add(const char *class_name, const char *func_name, struct _zend_function_entry *fe);

struct _zend_function_entry* phpgo_function_map_get(const char *class_name, const char *func_name);

void phpgo_class_map_add(const char *class_name, struct _zend_class_entry *ce);

struct _zend_class_entry* phpgo_class_map_get(const char *class_name);

void phpgo_callback_map_add(const char *class_name, const char *func_name, int cbid);

int phpgo_callback_map_get(const char *class_name, const char *func_name);

phpgo_object_map *phpgo_object_map_new();
void phpgo_object_map_add(phpgo_object_map** om, const char *name, void* obj);
void* phpgo_object_map_get(phpgo_object_map* om, const char *name);
int phpgo_object_map_count(phpgo_object_map* om);

#endif

