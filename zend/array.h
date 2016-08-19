#include <zend_API.h>

//create zval
zval* php_array_create_zval();

//create php7 zval
void php7_array_init(zval* zv);

//$arr[int] = int;
void php_array_add_index_long(void* arr, ulong idx, long n);
//$arr[int] = 'string';
void php_array_add_index_string(void* arr, ulong idx, char* value);
//$arr[int] = 3.14;
void php_array_add_index_double(void* arr, ulong idx, double d);


//$arr['string'] = int;
void php_array_add_assoc_long(void* arr, char* key, long n);
//$arr['string'] = 'string';
void php_array_add_assoc_string(void* arr, char* key, char* value);
//$arr['string'] = 3.14;
void php_array_add_assoc_double(void* arr, char* key, double d);



//$array[int] = array[]
void php_array_add_index_zval(void* arr, ulong index, void* value);
//$array['string'] = array[]
void php_array_add_assoc_zval(void* arr, char* key, void* value);
//$array[] = array[]
void php_array_add_next_index_zval(void* arr, void* zvalue);
