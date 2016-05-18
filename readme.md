

### php-go

Write PHP extension using go/golang. Zend API wrapper for go/golang. 

Simple, easy, fun to write PHP extensions.


### Features

* function support
* struct/class support
* constant support
* primitive data type as parameters and return values, (u)int*/float*/string/bool
* complex data type as parameters, map/slice/array
* all can be done by programitic


### Environment

* Modern Linux/Unix system
* PHP 5.5+/7.x
* go version 1.4+


### Build & Install

go get:

```
go get github.com/kitech/php-go
cd $GOPATH/src/github.com/kitech/php-go
make
```


manual:

    mkdir -p $GOPATH/src/github.com/kitech
    git clone https://github.com/kitech/php-go.git $GOPATH/src/kitech/php-go
    cd $GOPATH/src/kitech/php-go
    make
    ls -lh php-go/hello.so
    php56 -d extension=./hello.so examples/hello.php


note: php7 support not complete, use php5 for test.


### Examples

    import "phpgo"
    
    func module_startup(ptype int, module_number int) int {
        return 0
    }
    func module_shutdown(ptype int, module_number int) int {
        return 0
    }
    func request_startup(ptype int, module_number int) int {
        return 0
    }
    func request_shutdown(ptype int, module_number int) int {
        return 0
    }
    
    function foo_in_go() {
    }
    
    type Bar struct{}
    func NewBar() *Bar{
        return &Bar{}
    }
    
    func init() {
	    phpgo.InitExtension("mymod", "1.0")
	    phpgo.RegisterInitFunctions(module_startup, module_shutdown, request_startup, request_shutdown)
        
        phpgo.AddFunc("foo_in_php", foo_in_go)
        phpgo.AddClass("bar_in_php", NewBar)
    }
    

### TODO

- [ ] install with go get 
- [ ] improve php7 support
- [ ] namespace support


Contributing
------------
1. Fork it
2. Create your feature branch (``git checkout -b my-new-feature``)
3. Commit your changes (``git commit -am 'Add some feature'``)
4. Push to the branch (``git push origin my-new-feature``)
5. Create new Pull Request

