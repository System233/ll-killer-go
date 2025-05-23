GO  ?= go
GOARCH ?= $(shell $(GO) env GOARCH)
MODULE := $(shell $(GO) list -m)
TARGET ?= $(shell echo $$(uname -m)-$$(uname -s | tr '[:upper:]' '[:lower:]')-gnu)
ENABLE_NO_EVM ?= yes
LDFLAGS_STATIC := -extldflags -static
TRIMPATH := -trimpath

VERSION := $(shell git describe --tags --always)
BUILDTIME := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
ifeq ($(ENABLE_NO_EVM),yes)
	EXTRA_TAG := -nevm
else
	EXTRA_TAG := 
endif
LDFLAGS := -X "$(MODULE)/config.Version=$(VERSION)-$(GOARCH)$(EXTRA_TAG)" \
		   -X "$(MODULE)/config.BuildTime=$(BUILDTIME)" \
		   -X "$(MODULE)/config.Tag=$(VERSION)" \
		   -X "$(MODULE)/config.Variant=$(GOARCH)" \

FUSE_LIBS := libfuse-overlayfs.a libgnu.a
FUSE_DIR := fuse-overlayfs
FUSE_PROJECT_DEPS := configure.ac Makefile.am
FUSE_PROJECT_SRC= *.c *h
FUSE_PROJECT := $(foreach file, $(FUSE_PROJECT_DEPS), $(FUSE_DIR)/$(file))
FUSE_SRCS := $(foreach file, $(FUSE_PROJECT_SRC), $(FUSE_DIR)/$(file))

RES_DIRS := build-aux/build-aux build-aux/apt.conf.d
SRC_DIRS := build-aux config apps layer pty utils updater reexec
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
	if test "$(ENABLE_NO_EVM)" = "yes" ;then \
		git -C $(FUSE_DIR)  apply --check ../patches/fuse-overlayfs-nevm.patch -R -q || git -C $(FUSE_DIR) apply ../patches/fuse-overlayfs-nevm.patch; \
	else \
		git -C $(FUSE_DIR)  apply --check ../patches/fuse-overlayfs-nevm.patch -q || git -C $(FUSE_DIR) apply ../patches/fuse-overlayfs-nevm.patch -R; \
	fi
	git -C $(FUSE_DIR) apply --check ../patches/fuse-overlayfs.patch -R -q || git -C $(FUSE_DIR)  apply ../patches/fuse-overlayfs.patch; 
	cd $(FUSE_DIR) && ./autogen.sh && LIBS="-ldl" LDFLAGS="-static" ./configure --host=$(TARGET);

$(FUSE_LIBS): $(FUSE_DIR)/Makefile $(FUSE_SRCS)
	make -C $(FUSE_DIR)
	cp $(FUSE_DIR)/lib/libgnu.a \
	   $(FUSE_DIR)/libfuse-overlayfs.a .

test: ll-killer $(GO_TEST_SOURCES)
	$(GO) test -ldflags "$(LDFLAGS) $(LDFLAGS_STATIC) $(EXTRA_LDFLAGS)" $(GO_TEST_DIRS)

all: ll-killer test

.PHONY: all test