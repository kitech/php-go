#include <stdlib.h>
#include <stdio.h>
#include <dlfcn.h>

int main(int argc, char **argv)
{
    void *soh = dlopen("./hello.so", RTLD_NOW);
    void *(*get_mod)();
    get_mod = dlsym(soh, "get_module");
    void *v = get_mod();
    printf("soh: %p, %p, %p\n", soh, get_mod, v);
    return 0;
}
