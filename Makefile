# go option
GO         ?= go
PKG        := dep ensure -v
LDFLAGS    := -w -s
GOFLAGS    :=
TAGS       := linux libsqlite3
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

.PHONY: lint
lint:
	@echo " ===> running linter ... <=== "
	@revive -config default.toml -formatter stylish

.PHONY: build
build: lint
	@echo " ===> building releases in ./bin/... <=== "
	GOBIN=$(BINDIR) $(GO) install $(GOFLAGS) -compiler gccgo -gccgoflags '$(GCCGOGLAGS)' -tags '$(TAGS)' -ldflags '$(LDFLAGS)' github.com/tritonmedia/ignis/...

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