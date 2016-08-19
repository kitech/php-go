
ifeq ($(PHPDIR),)
	PHPDIR=/usr
endif

PHPCONFIG=$(PHPDIR)/bin/php-config

export PATH := $(PHPDIR)/bin:$(PATH)
export CFLAGS := $(shell export PATH=$(PATH) && $(PHPCONFIG) --includes)
export LDFLAGS := -L$(shell export PATH=$(PATH) && $(PHPCONFIG) --prefix)/lib/

export GOPATH := $(PWD):$(GOPATH)
export CGO_CFLAGS := $(CFLAGS) $(CGO_CFLAGS)
export CGO_LDFLAGS := $(LDFLAGS) $(CGO_LDFLAGS)

all:
	# go install ./zend
	# CGO_ENABLE=1 go install ./phpgo
	CGO_ENABLE=1 go build -gcflags="-newexport=0 -ssa=0" -v -linkshared -buildmode=c-shared -o hello.so miniphpgo.go
	# php -d extension=./hello.so miniphpgo.php

clean:
	rm -f ../../pkg/linux_amd64/zend.a
	rm -f ../../pkg/linux_amd64/phpgo.a
	rm -f hello.so
