#ifndef _PHPGO_EXTENSION_H_
#define _PHPGO_EXTENSION_H_

extern void* get_module_impl();
extern int zend_add_function(char *name);
extern int zend_add_class(char *name);

#endif
