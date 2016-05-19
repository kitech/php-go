
PHPDIR=/usr
PHPCONFIG=$(PHPDIR)/bin/php-config56

export PATH := $(PHPDIR)/bin:$(PATH)
export CFLAGS := $(shell export PATH=$(PATH) && $(PHPCONFIG) --includes)
export LDFLAGS := -L$(shell export PATH=$(PATH) && $(PHPCONFIG) --prefix)/lib/

export GOPATH := $(PWD):$(GOPATH)
export CGO_CFLAGS := $(CFLAGS) $(CGO_CFLAGS)
export CGO_LDFLAGS := $(LDFLAGS) $(CGO_LDFLAGS)

all:
	go install ./zend
	go install ./phpgo
	go build -v -linkshared -buildmode=c-shared -o hello.so examples/hello.go
	# php56 -d extension=./hello.so examples/hello.php

clean:
	rm -f ../../pkg/linux_amd64/zend.a
	rm -f ../../pkg/linux_amd64/phpgo.a

