#include "array.h"

#ifdef ZEND_ENGINE_3
#include <zend_long.h>
#endif

//create zval
//php7 zval在栈 无法返回 故返回nil 调用方处理
zval* php_array_create_zval() {
#ifdef ZEND_ENGINE_3
    return NULL;
#else
    zval* arr;
    ALLOC_INIT_ZVAL(arr);
    array_init(arr);
    return arr;
#endif
}

//create php7 zval
void php7_array_init(zval* zv) {
    array_init(zv);
}

//$arr[int] = int;
void php_array_add_index_long(void* arr, ulong idx, long n) {
#ifdef ZEND_ENGINE_3
    add_index_long((zval*)arr, (zend_ulong)idx, (zend_long)n);
#else
    add_index_long((zval*)arr, idx, n);
#endif
}
//$arr[int] = 'string';
void php_array_add_index_string(void* arr, ulong idx, char* value) {
#ifdef ZEND_ENGINE_3
    add_index_string((zval*)arr, (zend_ulong)idx, value);
#else
    add_index_string((zval*)arr, idx, value, 1);
#endif
}
//$arr[int] = 3.14;
void php_array_add_index_double(void* arr, ulong idx, double d) {
#ifdef ZEND_ENGINE_3
    add_index_double((zval*)arr, (zend_ulong)idx, d);
#else
    add_index_double((zval*)arr, idx, d);
#endif
}


//$arr['string'] = int;
void php_array_add_assoc_long(void* arr, char* key, long n) {
#ifdef ZEND_ENGINE_3
    add_assoc_long((zval*)arr, key, (zend_long)n);
#else
    add_assoc_long((zval*)arr, key, n);
#endif
}
//$arr['string'] = 'string';
void php_array_add_assoc_string(void* arr, char* key, char* value) {
#ifdef ZEND_ENGINE_3
    add_assoc_string((zval*)arr, key, value);
#else
    add_assoc_string((zval*)arr, key, value, 1);
#endif
}
//$arr['string'] = 3.14
void php_array_add_assoc_double(void* arr, char* key, double d) {
    add_assoc_double((zval*)arr, key, d);
}


//$array[int] = array[]
void php_array_add_index_zval(void* arr, ulong index, void* value) {
#ifdef ZEND_ENGINE_3
    add_index_zval((zval*)arr, (zend_ulong)index, (zval*)value);
#else
    add_index_zval((zval*)arr, index, (zval*)value);
#endif
}
//$array['string'] = array[]
void php_array_add_assoc_zval(void* arr, char* key, void* value) {
    add_assoc_zval((zval*)arr, key, (zval*)value);
}
//$array[] = array[]
void php_array_add_next_index_zval(void* arr, void* zvalue) {
    add_next_index_zval((zval*)arr, (zval*)zvalue);
}
