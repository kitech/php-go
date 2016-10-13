
# Usage:
# PHPCFG=/usr/bin/php-config55 make
# PHPCFG=php-config55 make
# PHPCFG=/usr/bin/php-config make
# PHPCFG=php-config make
# make

ifeq ($(PHPCFG),)
	PHPCFG=/usr/bin/php-config
endif

PHPEXE := $(shell $(PHPCFG) --php-binary)
PHPDIR := $(shell $(PHPCFG) --prefix)

export PATH := $(PHPDIR)/bin:$(PATH)
export CFLAGS := $(shell $(PHPCFG) --includes)
export LDFLAGS := -L$(shell $(PHPCFG) --prefix)/lib/

export GOPATH := $(PWD):$(GOPATH)
export CGO_CFLAGS := $(CFLAGS) $(CGO_CFLAGS)
export CGO_LDFLAGS := $(LDFLAGS) $(CGO_LDFLAGS)

all:
	go install ./zend
	go install ./phpgo
	go build -v -buildmode=c-shared -o hello.so examples/hello.go
	# $(PHPEXE) -d extension=./hello.so examples/hello.php

clean:
	rm -f $(GOPATH)/pkg/linux_amd64/zend.a
	rm -f $(GOPATH)/pkg/linux_amd64/phpgo.a
	rm -f hello.so
