
#define CLOG_MAIN
#include "clog.h"

// only for inline compile

void dlog_set_level(int id, enum clog_level level)
{
#ifdef LOGLEVEL
    if (strcasecmp(LOGLEVEL, "error") == 0) {
        clog_set_level(id, CLOG_ERROR);
    } else if (strcasecmp(LOGLEVEL, "info") == 0) {
        clog_set_level(id, CLOG_INFO);
    } else if (strcasecmp(LOGLEVEL, "warn") == 0) {
        clog_set_level(id, CLOG_WARN);
    } else {
        clog_set_level(id, CLOG_DEBUG);
    }
    #warning "defined log level"
#else
// #warning "not defined log level"
    clog_set_level(id, CLOG_DEBUG);
#endif
}
