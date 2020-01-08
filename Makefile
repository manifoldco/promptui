LINTERS=$(shell grep "// lint" tools/tools.go | awk '{gsub(/\"/, "", $$1); print $$1}' | awk -F / '{print $$NF}') \
	gofmt \
	vet

ci: $(LINTERS) cover

.PHONY: ci

#################################################
# Bootstrapping for base golang package and tool deps
#################################################

CMD_PKGS=$(shell grep '	"' tools/tools.go | awk -F '"' '{print $$2}')

define VENDOR_BIN_TMPL
tools/vendor/bin/$(notdir $(1)): tools/vendor/$(1) | vendor tools/vendor
	GOBIN=`pwd`/tools/vendor/bin sh -c 'cd tools && go install ./vendor/$(1)'
VENDOR_BINS += tools/vendor/bin/$(notdir $(1))
tools/vendor/$(1): tools/vendor
endef

$(foreach cmd_pkg,$(CMD_PKGS),$(eval $(call VENDOR_BIN_TMPL,$(cmd_pkg))))

$(patsubst %,%-bin,$(filter-out gofmt vet,$(LINTERS))): %-bin: tools/vendor/bin/%
gofmt-bin vet-bin:

tools/vendor: tools/go.sum
	GO111MODULE=on sh -c 'cd tools && go mod vendor'

vendor: go.sum
	GO111MODULE=on go mod vendor

mod-update:
	GO111MODULE=on go get -u -m
	GO111MODULE=on go mod tidy

mod-tidy:
	GO111MODULE=on go mod tidy

.PHONY: $(CMD_PKGS)
.PHONY: mod-update mod-tidy

#################################################
# Test and linting
#################################################

test: vendor
	CGO_ENABLED=0 go test $$(go list ./... | grep -v generated)

$(LINTERS): %: tools/vendor/bin/gometalinter %-bin tools/vendor
	PATH=`pwd`/tools/vendor/bin:$$PATH gometalinter --tests --disable-all --vendor \
		--deadline=5m -s data --enable $@ ./...

COVER_TEST_PKGS:=$(shell find . -type f -name '*_test.go' | grep -v vendor | rev | cut -d "/" -f 2- | rev | grep -v generated | sort -u)
$(COVER_TEST_PKGS:=-cover): %-cover: all-cover.txt
	@CGO_ENABLED=0 go test -v -coverprofile=$@.out -covermode=atomic ./$*
	@if [ -f $@.out ]; then \
		grep -v "mode: atomic" < $@.out >> all-cover.txt; \
		rm $@.out; \
	fi

all-cover.txt:
	echo "mode: atomic" > all-cover.txt

cover: vendor all-cover.txt $(COVER_TEST_PKGS:=-cover)

.PHONY: $(LINTERS) test
.PHONY: cover all-cover.txt
