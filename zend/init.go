package zend

/*
#include "szend.h"
#include "clog.h"
*/
import "C"

// import "unsafe"
import "runtime"

import "time"
import "math/rand"

// import "os"
import "fmt"
import "log"

const (
	SINGLE_THREAD = 1
)

//
func init() {
	initSingleThread()
	initCLog()
	initRawLog()
}

func initSingleThread() {
	if C.gozend_iszts() == 0 {
		omp := runtime.GOMAXPROCS(0)
		if omp > SINGLE_THREAD {
			runtime.GOMAXPROCS(SINGLE_THREAD)
			fmt.Printf("Adjust GOMAXPROCS %d => %d\n", omp, SINGLE_THREAD)
		}
	}

	rand.Seed(time.Now().UnixNano())
}

func initRawLog() {
	prefix := "phpgo"
	// log.SetFlags(log.Llongfile | log.LstdFlags)
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	oldPrefix := log.Prefix()
	log.SetPrefix(fmt.Sprintf("[%s ] ", prefix))
	if len(oldPrefix) > 0 {
		log.Printf("Switch log prefix: %s => [%s ]\n", oldPrefix, prefix)
	}
}

func initCLog() {
	C.clog_init_fd(C.STDOUT_FILENO, C.STDOUT_FILENO)
	C.clog_init_fd(C.STDERR_FILENO, C.STDERR_FILENO)
	C.dlog_set_level(C.STDOUT_FILENO, 0)
	C.dlog_set_level(C.STDERR_FILENO, 0)
}

const (
	LOG_NONE  = int(-1)
	LOG_ERROR = int(C.CLOG_ERROR)
	LOG_WARN  = int(C.CLOG_WARN)
	LOG_INFO  = int(C.CLOG_INFO)
	LOG_DEBUG = int(C.CLOG_DEBUG)
)

func LogInitFD(fd int) {
	C.clog_init_fd(C.int(fd), C.int(fd))
}

func LogSetLevel(fd int, level int) {
	C.dlog_set_level(C.int(fd), C.int(level))
}
