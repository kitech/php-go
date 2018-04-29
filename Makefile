# Usage:
#	Default config:
#		make
#	Use a specific php-config:
#		PHPCFG=/usr/bin/php-config72 make
#	Disable zend log (<error):
#		CGO_CFLAGS="-DLOGLEVEL=error" make

ifeq ($(PHPCFG),)
	PHPCFG=/usr/bin/php-config
endif

ifeq ($(APP),)
	APP=examples
endif

PHPEXE := $(shell $(PHPCFG) --php-binary)
PHPDIR := $(shell $(PHPCFG) --prefix)

export PATH := $(PHPDIR)/bin:$(PATH)
export CFLAGS := $(shell $(PHPCFG) --includes)
export LDFLAGS := -L$(shell $(PHPCFG) --prefix)/lib/

export CGO_CFLAGS := $(CFLAGS) $(CGO_CFLAGS)
export CGO_LDFLAGS := $(LDFLAGS) $(CGO_LDFLAGS)

GOOS := $(shell go env GOOS)
GOARCH := $(shell go env GOARCH)

all:
	go install ./zend
	go install ./phpgo
	go build -buildmode=c-shared -o $(APP).so $(APP)/main.go
	# $(PHPEXE) -d extension=./$(APP).so $(APP)/main.php

zdlib:
	go install -v -x ./zend

pglib:
	go install -v -x ./phpgo

test:
	$(PHPEXE) -d extension=./hello.so examples/hello.php

# quickly build check syntax errors
quickbc:
	go build -v -x -buildmode=c-archive -o php-grpc.c.a examples/hello.go

clean:
	rm -f $(GOPATH)/pkg/$(GOOS)_$(GOARCH)/github.com/kitech/php-go/zend.a
	rm -f $(GOPATH)/pkg/$(GOOS)_$(GOARCH)/github.com/kitech/php-go/phpgo.a
	rm -f hello.so
