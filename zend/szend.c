
#include <stdlib.h>
#include <stdint.h>

#include "sztypes.h"

#ifdef HAVE_CONFIG_H
#include "config.h"
#endif
#include "zend_API.h"

#ifdef ZEND_ENGINE_3
int gozend_call_user_function(zval *object, char *func_name, zval *retval_ptr, int argc, zval *params)
{
    zval function_name;
    ZVAL_NULL(&function_name);
    ZVAL_STRING(&function_name, func_name);

    assert(retval_ptr != NULL);
    if (call_user_function(CG(function_table), object , &function_name,
                           retval_ptr, argc, params TSRMLS_CC) == SUCCESS) {
        return 0;
    }
    return -1;
}
#else /* ZEND_ENGINE_2 */
int gozend_call_user_function(zval **object, char *func_name, zval *retval_ptr, int argc, zval **params)
{
  zval function_name;
  INIT_ZVAL(function_name);
  ZVAL_STRING(&function_name, func_name, 1);

  assert(retval_ptr != NULL);
  if (call_user_function(CG(function_table), object , &function_name,
                         retval_ptr, argc, params TSRMLS_CC) == SUCCESS) {
    return 0;
  }
  return -1;
}
#endif


#ifdef ZEND_ENGINE_3
static int gozend_call_user_function_string_ex(char *func_name, char *str, char **retstr, zval *object)
{
  zval args[9];
  ZVAL_NULL(&args[0]);
  ZVAL_STRING(&args[0], str);

  zval retval_ptr;
  ZVAL_NULL(&retval_ptr);

    int ret = gozend_call_user_function(object, func_name, &retval_ptr, 1, args);
    if (ret == 0) {
    char *cstr = NULL;
    int cslen = 0;

    if (Z_TYPE_P(&retval_ptr) == IS_STRING) {
      cslen = Z_STRLEN_P(&retval_ptr);
      cstr = estrndup(Z_STRVAL_P(&retval_ptr), cslen);
      *retstr = cstr;
    }
  }
  return ret;
}
#else /* ZEND_ENGINE_2 */
static int gozend_call_user_function_string_ex(char *func_name, char *str, char **retstr, zval **object)
{
  zval *args[9];
  MAKE_STD_ZVAL(args[0]);
  ZVAL_STRING(args[0], str, 1);

  zval retval_ptr;
  ZVAL_NULL(&retval_ptr)

    int ret = gozend_call_user_function(object, func_name, &retval_ptr, 1, args);
  if (ret == 0) {
    char *cstr = NULL;
    int cslen = 0;

    if (Z_TYPE_P(&retval_ptr) == IS_STRING) {
      cslen = Z_STRLEN_P(&retval_ptr);
      cstr = estrndup(Z_STRVAL_P(&retval_ptr), cslen);
      *retstr = cstr;
    }
  }
  return ret;
}
#endif


char *gozend_call_user_function_string(char *func_name, char *str)
{
    char *res = NULL;
    int ret = gozend_call_user_function_string_ex(func_name, str, &res, NULL);
    return res;
}

int gozend_iszts() {
#ifdef ZTS
    return 1;
#else
    return 0;
#endif
}

int gozend_zend_version_no() {
#ifdef ZEND_ENGINE_3
     return 3;
#else
     return 2;
#endif
}

char* gozend_zend_version() {
    return get_zend_version();
}

void gozend_efree(void *ptr) {
    efree(ptr);
}

char *gozend_estrdup(char *str) {
    return estrdup(str);
}

char *gozend_estrndup(char *str, unsigned int length) {
    return estrndup(str, length);
}

void *gozend_emalloc(size_t size) {
    return emalloc(size);
}

void *gozend_ecalloc(size_t nmemb, size_t size) {
    return ecalloc(nmemb, size);
}

void *gozend_erealloc(void *ptr, size_t size) {
    return erealloc(ptr, size);
}

char gozend_eval_string(char *code)
{
    zval retval;
    #ifdef ZEND_ENGINE_3
    ZVAL_NULL(&retval);
    #else
    INIT_ZVAL(retval);
    #endif

    int ret = zend_eval_string(code, &retval, (char*)"" TSRMLS_CC);

    // zval_ptr_dtor(&retval);
    zval_dtor(&retval);
    return ret == FAILURE;
}


void call_user_function_callback(char *data)
{
    uint32_t argc = 0;
    zval retval_ptr;
    char *func_name = "say_hello_123";
    zval function_name;
#ifdef ZEND_ENGINE_3
    ZVAL_NULL(&retval_ptr);
    ZVAL_NULL(&function_name);
    ZVAL_STRING(&function_name, func_name);
#else
    INIT_ZVAL(retval_ptr);
    INIT_ZVAL(function_name);
    ZVAL_STRING(&function_name, func_name, 1);
#endif

    void *cobj = NULL; /* no object */
    if (call_user_function(CG(function_table), cobj , &function_name,
                           &retval_ptr, argc, NULL TSRMLS_CC) == SUCCESS) {
        /* do something with retval_ptr here if you like */
    }

    /* don't forget to free the zvals */
    // zval_ptr_dtor(&retval_ptr);
    zval_dtor(&function_name);
}

// find function name in HashTable
static int _gozend_function_exists_ht(char *fname, HashTable* ht) {
#ifdef ZEND_ENGINE_3
	zend_string *name;
	zend_function *func;
	zend_string *lcname;

    name = zend_string_init(fname, strlen(fname), 0);

	if (ZSTR_VAL(name)[0] == '\\') {
		/* Ignore leading "\" */
		lcname = zend_string_alloc(ZSTR_LEN(name) - 1, 0);
		zend_str_tolower_copy(ZSTR_VAL(lcname), ZSTR_VAL(name) + 1, ZSTR_LEN(name) - 1);
	} else {
		lcname = zend_string_tolower(name);
	}

	func = zend_hash_find_ptr(ht, lcname);
	zend_string_release(lcname);
    zend_string_release(name);

	/*
	 * A bit of a hack, but not a bad one: we see if the handler of the function
	 * is actually one that displays "function is disabled" message.
	 */
	if (func && (func->type != ZEND_INTERNAL_FUNCTION ||
                 func->internal_function.handler != zif_display_disabled_function)) {
        return 1;
    } else {
        return 0;
    }
#else
	char *name;
	int name_len;
	zend_function *func;
    int retval = 0;

    name = fname;
    name_len = strlen(fname);
	/* Ignore leading "\" */
	if (name[0] == '\\') {
		name = &name[1];
		name_len--;
	}

	retval = (zend_hash_find(ht, name, name_len+1, (void **)&func) == SUCCESS);

	/*
	 * A bit of a hack, but not a bad one: we see if the handler of the function
	 * is actually one that displays "function is disabled" message.
	 */
	if (retval && func->type == ZEND_INTERNAL_FUNCTION &&
		func->internal_function.handler == zif_display_disabled_function) {
		retval = 0;
	}

    return retval;
#endif
    return 0;
}

// seems only can be used when vm executing
int gozend_function_exists(char *fname) {
    return _gozend_function_exists_ht(fname, executor_globals.function_table);
}

int gozend_function_registered(char *fname) {
    // zif_function_exists(); // this call method not compatible on MacOS
    // zif_function_existsccc();

    return _gozend_function_exists_ht(fname, compiler_globals.function_table);
}

