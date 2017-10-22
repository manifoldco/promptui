LINTERS=\
	gofmt \
	golint \
	gosimple \
	vet \
	misspell \
	ineffassign \
	deadcode

ci: $(LINTERS) test

.PHONY: ci

#################################################
# Bootstrapping for base golang package deps
#################################################

CMD_PKGS=\
	github.com/golang/lint/golint \
	honnef.co/go/tools/cmd/gosimple \
	github.com/client9/misspell/cmd/misspell \
	github.com/gordonklaus/ineffassign \
	github.com/tsenart/deadcode \
	github.com/alecthomas/gometalinter

define VENDOR_BIN_TMPL
vendor/bin/$(notdir $(1)): vendor
	go build -o $$@ ./vendor/$(1)
VENDOR_BINS += vendor/bin/$(notdir $(1))
endef

$(foreach cmd_pkg,$(CMD_PKGS),$(eval $(call VENDOR_BIN_TMPL,$(cmd_pkg))))
$(patsubst %,%-bin,$(filter-out gofmt vet,$(LINTERS))): %-bin: vendor/bin/%
gofmt-bin vet-bin:

bootstrap:
	which dep || go get -u github.com/golang/dep/cmd/dep

vendor: Gopkg.lock
	dep ensure

.PHONY: bootstrap $(CMD_PKGS)

#################################################
# Test and linting
#################################################

test: vendor
	@CGO_ENABLED=0 go test -v $$(go list ./... | grep -v vendor)

$(LINTERS): %: vendor/bin/gometalinter %-bin vendor
	PATH=`pwd`/vendor/bin:$$PATH gometalinter --tests --disable-all --vendor \
	    --deadline=5m -s data ./... --enable $@

.PHONY: $(LINTERS) test
