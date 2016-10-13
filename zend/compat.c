#include "compat.h"

ZEND_API void zend_register_stringl_constant_compat(const char *name, size_t name_len, char *strval, size_t strlen, int flags, int module_number)
{
    zend_register_stringl_constant(name, name_len, strval, strlen, flags, module_number);
}

ZEND_API void zend_register_long_constant_compat(const char *name, size_t name_len, zend_long lval, int flags, int module_number)
{
    zend_register_long_constant(name, name_len, lval, flags, module_number);
}

ZEND_API void zend_register_double_constant_compat(const char *name, size_t name_len, double dval, int flags, int module_number)
{
    zend_register_double_constant(name, name_len, dval, flags, module_number);
}

ZEND_API void zend_register_bool_constant_compat(const char *name, size_t name_len, zend_bool bval, int flags, int module_number)
{
#ifdef ZEND_ENGINE_3
    zend_register_bool_constant(name, name_len, bval, flags, module_number);
#else
    zend_register_long_constant(name, name_len, bval, flags, module_number);
#endif
}

ZEND_API void zend_register_null_constant_compat(const char *name, size_t name_len, int flags, int module_number)
{
#ifdef ZEND_ENGINE_3
    zend_register_null_constant(name, name_len, flags, module_number);
#else
    zend_register_stringl_constant(name, name_len, NULL, 0, flags, module_number);
#endif
}

