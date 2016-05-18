package main

import "log"
import "reflect"
import "math/rand"
import "time"
import "phpgo"

type PGDemo struct {
}

func NewPGDemo() *PGDemo {
	log.Println("NewPGDemo...")
	return &PGDemo{}
}

func (this *PGDemo) Hello1() int {
	rn := rand.Intn(65536)
	log.Println("called in go scope, hello1 method", rn)
	return rn
}

func (this *PGDemo) Hello2() {

}

func (this *PGDemo) Hello3() {

}

func (this *PGDemo) Hello4() {

}

func (this *PGDemo) Hello5() {

}

func (this *PGDemo) Hello6() {

}

func phpgo_hello() {
	log.Println("ext user func called")
}

func phpgo_hello3(a0 int, a1 string, a2 float64) {
	log.Println("ext user func called222", a0, a1, a2)
}

func phpgo_hello7(names []int64) {
	log.Println(names)
}

func phpgo_hello8(names map[string]string) {
	log.Println(names)
}

func module_startup(ptype int, module_number int) int {
	println("module_startup", ptype, module_number)
	return rand.Int()
}
func module_shutdown(ptype int, module_number int) int {
	println("module_shutdown", ptype, module_number)
	return rand.Int()
}
func request_startup(ptype int, module_number int) int {
	println("request_startup", ptype, module_number)
	return rand.Int()
}
func request_shutdown(ptype int, module_number int) int {
	println("request_shutdown", ptype, module_number)
	return rand.Int()
}

func init() {
	log.Println("run us init...")
	rand.Seed(time.Now().UnixNano())

	f2 := func() int {
		log.Println("ext user func called222 closure")
		return 567
	}

	phpgo.InitExtension("pg0", "")
	phpgo.RegisterInitFunctions(module_startup, module_shutdown, request_startup, request_shutdown)

	phpgo.AddFunc("foo_hello", phpgo_hello)
	// phpgo.AddFunc("foo_hello2", phpgo_hello2)
	phpgo.AddFunc("foo_hello2", f2)
	phpgo.AddFunc("foo_hello3", phpgo_hello3)
	phpgo.AddFunc("foo_hello7", phpgo_hello7)
	phpgo.AddFunc("foo_hello8", phpgo_hello8)

	//
	dm := NewPGDemo()
	log.Println("method fn:", dm.Hello1, reflect.TypeOf(dm.Hello1), NewPGDemo)

	phpgo.AddClass("PGDemo", NewPGDemo)
	// zend.AddMethod("PGDemo", "hello1", PGDemo.hello1)
}

func main() {
	panic("wtf")
	log.Println("run here for what?")
}
