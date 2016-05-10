#include <php.h>

char* gozend_php_version() {
    return PHP_VERSION;
}

int gozend_php_version_id() {
    return PHP_VERSION_ID;
}

