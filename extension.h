#ifndef _PHPGO_EXTENSION_H_
#define _PHPGO_EXTENSION_H_

extern void* get_module_impl();
extern int zend_add_function(int cidx, int fidx, int cbid, char *name, char *atys, int rety);
extern int zend_add_class(int cidx, char *name);
extern int zend_add_method(int cidx, int fidx, int cbid, char *name, char *mname, char *atys, int rety);

#endif
