package zend

/*
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
	initRawLog()
}

func initSingleThread() {
	omp := runtime.GOMAXPROCS(0)
	if omp > SINGLE_THREAD {
		runtime.GOMAXPROCS(SINGLE_THREAD)
		fmt.Printf("Adjust GOMAXPROCS %d => %d\n", omp, SINGLE_THREAD)
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
