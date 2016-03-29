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

//
func init() {
	omp := runtime.GOMAXPROCS(1)
	fmt.Println("run gozend init...", omp)
	omp = runtime.GOMAXPROCS(1)
	fmt.Println("run gozend init...", omp)

	rand.Seed(time.Now().UnixNano())

}
