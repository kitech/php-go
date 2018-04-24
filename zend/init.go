package zend

/*
#include "szend.h"
#include "clog.h"
*/
import "C"

// import "unsafe"
import (
	"io/ioutil"
	"log"
	"math/rand"
	"runtime"
	"time"
)

const (
	SINGLE_THREAD = 1
)

func init() {
	initRawLog()
	initCLog()
	initSingleThread()
}

func initRawLog() {
	// To enable use true:
	if false {
		log.SetFlags(log.Lshortfile | log.LstdFlags)
		log.SetPrefix("[phpgo] ")
	} else {
		log.SetFlags(0)
		log.SetOutput(ioutil.Discard)
	}
}

func initCLog() {
	C.clog_init_fd(C.STDOUT_FILENO, C.STDOUT_FILENO)
	C.clog_init_fd(C.STDERR_FILENO, C.STDERR_FILENO)
	C.dlog_set_level(C.STDOUT_FILENO, 0)
	C.dlog_set_level(C.STDERR_FILENO, 0)
}

func initSingleThread() {
	if C.gozend_iszts() == 0 {
		omp := runtime.GOMAXPROCS(0)
		if omp > SINGLE_THREAD {
			runtime.GOMAXPROCS(SINGLE_THREAD)
			log.Printf("Adjust GOMAXPROCS %d => %d\n", omp, SINGLE_THREAD)
		}
	}

	rand.Seed(time.Now().UnixNano())
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
