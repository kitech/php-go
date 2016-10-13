#ifndef _PHP_COMPAT_H_
#define _PHP_COMPAT_H_

#include <zend.h>
#include <zend_constants.h>
#include <zend_modules.h>

#ifdef ZEND_ENGINE_3
// #define Z_BVAL_P(zv) Z_LVAL_P(zv)
/*
#define Z_BVAL_PP(zv) Z_LVAL_P(*(zv))
#define Z_LVAL_PP(zv) Z_LVAL_P(*(zv))
#define Z_DVAL_PP(zv) Z_DVAL_P(*(zv))
#define Z_STRVAL_PP(zv) Z_STRVAL_P(*(zv))
#define Z_STRLEN_PP(zv) Z_STRLEN_P(*(zv))
#define Z_ARRVAL_PP(zv) Z_ARRVAL_P(*(zv))
#define Z_TYPE_PP(zv) Z_TYPE_P(*(zv))
*/
#endif

#ifdef ZEND_ENGINE_2
typedef long zend_long;
#define _IS_BOOL IS_BOOL
#define IS_TRUE (IS_CALLABLE+1)
#define IS_FALSE (IS_CALLABLE+2)
#define IS_UNDEF (IS_CALLABLE+3)

#if !defined(IS_CONSTANT_AST)
#if ZEND_MODULE_API_NO < 20121212  // < PHP 5.5
#error "Maybe you are using not supported zend < 2.5.0 (PHP < 5.5.0)"
#else
#define IS_CONSTANT_AST (IS_CALLABLE+4)
#endif
#endif

#endif

ZEND_API void zend_register_stringl_constant_compat(const char *name, size_t name_len, char *strval, size_t strlen, int flags, int module_number);
ZEND_API void zend_register_long_constant_compat(const char *name, size_t name_len, zend_long lval, int flags, int module_number);
ZEND_API void zend_register_double_constant_compat(const char *name, size_t name_len, double dval, int flags, int module_number);
ZEND_API void zend_register_bool_constant_compat(const char *name, size_t name_len, zend_bool bval, int flags, int module_number);
ZEND_API void zend_register_null_constant_compat(const char *name, size_t name_len, int flags, int module_number);


#endif
