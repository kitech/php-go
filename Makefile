
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

export CGO_CFLAGS := $(CFLAGS) $(CGO_CFLAGS)
export CGO_LDFLAGS := $(LDFLAGS) $(CGO_LDFLAGS)

GOOS := $(shell go env GOOS)
GOARCH := $(shell go env GOARCH)

all:
	go install ./zend
	go install ./phpgo
	go build -v -buildmode=c-shared -o hello.so examples/hello.go
	# $(PHPEXE) -d extension=./hello.so examples/hello.php

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
	rm -f $(GOPATH)/pkg/$(GOOS)_$(GOARCH)/zend.a
	rm -f $(GOPATH)/pkg/$(GOOS)_$(GOARCH)/phpgo.a
	rm -f hello.so
