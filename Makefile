# go option
GO         ?= go
PKG        := dep ensure -v
LDFLAGS    := -w -s
GOFLAGS    :=
TAGS       := 
GCCGOGLAGS := -fgo-optimize-allocs -O3 -march=native
CFLAGS     := -fgo-optimize-allocs -O3 -march=native
BINDIR     := $(CURDIR)/bin

# Required for globs to work correctly
SHELL=/bin/bash


.PHONY: all
all: build
.PHONY: dep
dep:
	@echo " ===> Installing dependencies via '$$(echo $(PKG) | awk '{ print $$1 }')' <=== "
	@$(PKG)

.PHONY: build
build:
	@echo " ===> building releases in ./bin/... <=== "
	CGO_ENABLED=1 GOBIN=$(BINDIR) $(GO) install -v $(GOFLAGS) -tags '$(TAGS)' -ldflags '$(LDFLAGS)' github.com/tritonmedia/ignis/...

.PHONY: release
release: lint build
	@echo " ===> building release <=== "
	@strip bin/**
	@tar -cvzf ignis.tar.gz bin/ignis
	@mkdir -p "_dist"
	@mv ignis.tar.gz _dist
	@make checksum


.PHONY: checksum
checksum:
	for f in _dist/*.tar.gz ; do \
		shasum -a 256 "$${f}"  | awk '{print $$1}' > "$${f}.sha256" ; \
	done

.PHONY: clean
clean:
	@rm -rf $(BINDIR)
	@echo Going to delete db, giving you 5s to ^C
	@sleep 5
	@rm ignis.db