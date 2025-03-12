GO  ?= go
GOARCH ?= $(shell $(GO) env GOARCH)
MODULE := $(shell $(GO) list -m)
TARGET ?= $(shell echo $$(uname -m)-$$(uname -s | tr '[:upper:]' '[:lower:]')-gnu)
LDFLAGS_STATIC := -extldflags -static
TRIMPATH := -trimpath

VERSION := $(shell git describe --tags --always)
BUILDTIME := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS := -X "$(MODULE)/config.Version=$(VERSION)" -X "$(MODULE)/config.BuildTime=$(BUILDTIME)"

FUSE_LIBS := libfuse-overlayfs.a libgnu.a
FUSE_DIR := fuse-overlayfs
FUSE_PROJECT_DEPS := configure.ac Makefile.am
FUSE_PROJECT_SRC=Makefile *.c *h
FUSE_PROJECT := $(foreach file, $(FUSE_PROJECT_DEPS), $(FUSE_DIR)/$(file))
FUSE_SRCS := $(foreach file, $(FUSE_PROJECT_SRC), $(FUSE_DIR)/$(file))

RES_DIRS := build-aux/build-aux build-aux/apt.conf.d
SRC_DIRS := build-aux config apps layer pty utils types
GO_SOURCES:= $(wildcard *.go) $(shell find $(SRC_DIRS) -name '*.go')
GO_RESOURCES:= $(foreach dir, $(RES_DIRS), $(dir)/*)
GO_TEST_SOURCES:= $(wildcard *_test.go) $(shell find $(SRC_DIRS) -name '*_test.go')
GO_BUILD := $(GO) build $(TRIMPATH) $(GO_BUILDMODE_STATIC) \
	$(EXTRA_FLAGS) -ldflags "$(LDFLAGS) $(LDFLAGS_STATIC) $(EXTRA_LDFLAGS)"
GO_TEST_DIRS := $(shell echo $(GO_TEST_SOURCES)|xargs -r dirname|sort -u)
GO_TEST_DIRS:= $(foreach dir, $(GO_TEST_DIRS), ./$(dir))

ll-killer: $(GO_SOURCES) $(GO_RESOURCES) $(FUSE_LIBS)
	$(GO_BUILD) -o $@ .

$(FUSE_DIR)/Makefile: $(FUSE_PROJECT)
	cd $(FUSE_DIR);\
	git apply --check ../patches/fuse-overlayfs.patch -q && git apply ../patches/fuse-overlayfs.patch;\
	./autogen.sh;\
	LIBS="-ldl" LDFLAGS="-static" ./configure --host=$(TARGET);

$(FUSE_LIBS): $(FUSE_SRCS)
	make -C $(FUSE_DIR)
	cp $(FUSE_DIR)/lib/libgnu.a \
	   $(FUSE_DIR)/libfuse-overlayfs.a .

test: ll-killer $(GO_TEST_SOURCES)
	$(GO) test -ldflags "$(LDFLAGS) $(LDFLAGS_STATIC) $(EXTRA_LDFLAGS)" $(GO_TEST_DIRS)

all: ll-killer test

.PHONY: all test