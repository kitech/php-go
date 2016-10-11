#ifndef _PHPGO_CLASS_H_
#define _PHPGO_CLASS_H_

#define GLOBAL_VCLASS_NAME "_PHPGO_GLOBAL_"
#define MAX_ARG_NUM 10

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
zend_class_entry* phpgo_class_get(phpgo_class_entry* pce);
zend_function_entry* phpgo_class_get_funcs(phpgo_class_entry* pce);
int phpgo_class_method_count(phpgo_class_entry* pce);

//
typedef struct _phpgo_callback_info phpgo_callback_info;
phpgo_callback_info* phpgo_callback_info_new(char *arg_types, int ret_type);
char* phpgo_callback_info_get_arg_types(phpgo_callback_info* cbi);
int phpgo_callback_info_get_ret_type(phpgo_callback_info* cbi);

#endif

