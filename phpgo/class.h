#ifndef _PHPGO_CLASS_H_
#define _PHPGO_CLASS_H_

// struct _phpgo_function_entry;
typedef struct _phpgo_function_entry phpgo_function_entry;

phpgo_function_entry *phpgo_function_new(const char *func_name);
int phpgo_function_delete(phpgo_function_entry *pfe);
zend_function_entry* phpgo_function_get(phpgo_function_entry* pfe);

///////
// struct _phpgo_class_entry;
typedef struct _phpgo_class_entry phpgo_class_entry;

phpgo_class_entry* phpgo_class_new(const char *class_name);
void phpgo_class_method_add(phpgo_class_entry* pce, const char *func_name);
phpgo_function_entry* phpgo_class_method_get(phpgo_class_entry* pce, const char *func_name);

#endif

