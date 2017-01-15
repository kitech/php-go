#ifndef _SZEND_H_
#define _SZEND_H_

#include <stdlib.h>
#include <stdint.h>

extern char* gozend_call_user_function_string(char *func_name, char *arg);
extern int gozend_iszts();
extern void gozend_efree(void *ptr);
extern char *gozend_estrdup(char *str);
extern char *gozend_estrndup(char *str, unsigned int lenght);
extern void *gozend_emalloc(size_t size);
extern void *gozend_ecalloc(size_t nmemb, size_t size);
extern void *gozend_erealloc(void *ptr, size_t size);
extern char gozend_eval_string(char *code);
extern void call_user_function_callback(char *arg);
extern int gozend_function_exists(char *fname);
extern int gozend_function_registered(char *fname);
extern int gozend_zend_version_no();

#endif
