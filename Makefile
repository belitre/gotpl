# makefile based on the one from helm: https://github.com/kubernetes/helm/blob/master/Makefile

GO        ?= go
BINDIR    := $(CURDIR)/bin
LDFLAGS   := -w -s
TESTS     := ./...
TESTFLAGS :=

TARGETS   ?= darwin/amd64 linux/amd64 windows/amd64
DIST_DIRS = find * -type d -exec

# Required for globs to work correctly
SHELL=/bin/bash

.PHONY: build
build:
	GOBIN=$(BINDIR) $(GO) install -ldflags '$(LDFLAGS)' github.com/belitre/gotpl/...

# usage: make clean build-cross dist VERSION=v0.2-alpha
.PHONY: build-cross
build-cross: LDFLAGS += -extldflags "-static"
build-cross:
	CGO_ENABLED=0 gox -parallel=3 -output="_dist/{{.OS}}-{{.Arch}}/{{.Dir}}" -osarch='$(TARGETS)' -ldflags '$(LDFLAGS)' github.com/belitre/gotpl

.PHONY: dist
dist:
	( \
		cd _dist && \
		$(DIST_DIRS) tar -zcf gotpl-${VERSION}-{}.tar.gz {} \; && \
		$(DIST_DIRS) zip -r gotpl-${VERSION}-{}.zip {} \; \
	)

.PHONY: test
test: build
test: TESTFLAGS += -race -v
test: test-unit

.PHONY: test-unit
test-unit:
	@echo
	@echo "==> Running unit tests <=="
	$(GO) test $(TESTS) $(GOFLAGS) $(TESTFLAGS)

.PHONY: clean
clean:
	@rm -rf $(BINDIR) ./_dist

HAS_GODEP := $(shell command -v dep;)
HAS_GOX := $(shell command -v gox;)
HAS_GIT := $(shell command -v git;)

.PHONY: bootstrap
bootstrap:
ifndef HAS_GODEP
	$(error You must install dep: https://github.com/golang/dep)
endif
ifndef HAS_GOX
	go get -u github.com/mitchellh/gox
endif

ifndef HAS_GIT
	$(error You must install Git)
endif
	dep ensure

include versioning.mk