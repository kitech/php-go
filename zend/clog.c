
#define CLOG_MAIN
#include "clog.h"

// only for inline compile

#define STRINGIZE_NX(A) #A
#define STRINGIZE(A) STRINGIZE_NX(A)

void dlog_set_level(int id, int level)
{
#ifdef LOGLEVEL
    if (strcasecmp(STRINGIZE(LOGLEVEL), "error") == 0) {
        clog_set_level(id, CLOG_ERROR);
    } else if (strcasecmp(STRINGIZE(LOGLEVEL), "info") == 0) {
        clog_set_level(id, CLOG_INFO);
    } else if (strcasecmp(STRINGIZE(LOGLEVEL), "warn") == 0) {
        clog_set_level(id, CLOG_WARN);
    } else {
        clog_set_level(id, CLOG_DEBUG);
    }
    // #warning "defined log level"
#else
    // #warning "not defined log level"
    clog_set_level(id, CLOG_DEBUG);
#endif
}
