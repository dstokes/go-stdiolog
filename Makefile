GOOS ?= linux
GOARCH ?= amd64
VERSION ?= 0.0.1

GOBIN := $(GOPATH)/bin
DEBFILE := stdiolog_$(VERSION)-0_$(GOARCH).deb

all: build

build: fpm $(DEBFILE)

$(DEBFILE):
	@GOOS=$(GOOS) GOARCH=$(GOARCH) go install
	@fpm -a all -s dir -t deb -n stdiolog -v $(VERSION) -p $(DEBFILE) $(GOBIN)/$(GOOS)_$(GOARCH)/go-stdiolog=/usr/bin/stdiolog

fpm:
ifeq ($(shell which fpm),)
	$(error Install fpm first (https://github.com/jordansissel/fpm))
endif

clean:
	@rm $(DEBFILE)

.PHONY: fpm clean
