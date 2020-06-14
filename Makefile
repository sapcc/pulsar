IMAGE   ?= sapcc/pulsar
GOOS    ?= $(shell go env | grep GOOS | cut -d'"' -f2)
BINARY  := pulsar

SRCDIRS  := cmd pkg
PACKAGES := $(shell find $(SRCDIRS) -type d)
GO_PKG	 := github.com/sapcc/pulsar
GOFILES  := $(addsuffix /*.go,$(PACKAGES))
GOFILES  := $(wildcard $(GOFILES))

.PHONY: all
all: bin/$(GOOS)/$(BINARY)

bin/%/$(BINARY): BUILD_DATE= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
bin/%/$(BINARY): GIT_REVISION= $(shell git rev-parse --short HEAD)
bin/%/$(BINARY): GIT_BRANCH=$(shell git rev-parse --abbrev-ref HEAD)
bin/%/$(BINARY): VERSION=$(shell cat VERSION)
bin/%/$(BINARY): $(GOFILES) Makefile
	go build -mod vendor -ldflags "-s -w -X github.com/sapcc/pulsar/pkg/version.Revision=$(GIT_REVISION) -X github.com/sapcc/pulsar/pkg/version.Branch=$(GIT_BRANCH) -X github.com/sapcc/pulsar/pkg/version.BuildDate=$(BUILD_DATE) -X github.com/sapcc/pulsar/pkg/version.Version=$(VERSION)" -o bin/$*/$(BINARY) main.go

build: VERSION=$(shell cat VERSION)
build: bin/linux/$(BINARY)
	docker build -t $(IMAGE):$(VERSION) .

.PHONY: tests
tests:
	@if s="$$(gofmt -s -l *.go pkg 2>/dev/null)" && test -n "$$s"; then printf ' => %s\n%s\n' gofmt  "$$s"; false; fi
	DEBUG=1 && go test -v github.com/sapcc/pulsar/pkg/... | grep -v "no test files"

push: VERSION=$(shell cat VERSION)
push: build
	docker push $(IMAGE):$(VERSION)

.PHONY: clean
clean:
	rm -rf bin/*

.PHONY: vendor
vendor:
	go mod vendor

git-push-tag: VERSION=$(shell cat VERSION)
git-push-tag:
	git push origin ${VERSION}

git-tag-release: VERSION=$(shell cat VERSION)
git-tag-release: check-release-version
	git tag --annotate ${VERSION} --message "autoscaler ${VERSION}"

check-release-version: VERSION=$(shell cat VERSION)
check-release-version:
	if test x$$(git tag --list ${VERSION}) != x; \
	then \
		echo "Tag [${VERSION}] already exists. Please check the working copy."; git diff . ; exit 1;\
	fi

relase: VERSION=$(shell cat VERSION)
release: git-tag-release git-push-tag push
