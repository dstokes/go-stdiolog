GOOS ?= linux
GOARCH ?= amd64
VERSION ?= $(shell awk '/Version/ {gsub("\"", "", $$4); print $$4}' version.go)

XC_OS ?= darwin linux
XC_ARCH ?= 386 amd64

GOBIN := $(GOPATH)/bin
DEBFILE := stdiolog_$(VERSION)-0_$(GOARCH).deb

all: build

build: $(BUILDS)

deb: fpm $(DEBFILE)

fpm:
ifeq ($(shell which fpm),)
	$(error Install fpm first (https://github.com/jordansissel/fpm))
endif

$(DEBFILE):
	@GOOS=$(GOOS) GOARCH=$(GOARCH) go install
	@fpm -a all -s dir -t deb -n stdiolog -v $(VERSION) -p $(DEBFILE) $(GOBIN)/$(GOOS)_$(GOARCH)/go-stdiolog=/usr/bin/stdiolog

# Create a github release
release:
	@gox -os="$(XC_OS)" -arch="$(XC_ARCH)" -output "dist/{{.OS}}_{{.Arch}}/{{.Dir}}"
	@for platform in $$(find ./dist -mindepth 1 -maxdepth 1 -type d); do \
		pushd $$platform >/dev/null; \
		zip ../$$(basename $$platform).zip ./* >/dev/null; \
		popd >/dev/null; \
	done
	@ghr -u dstokes $(VERSION) dist/

clean:
	@rm -rf dist/

.PHONY: fpm clean release version
