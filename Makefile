GO  ?= go
GOARCH ?= $(shell $(GO) env GOARCH)
TARGET ?= $(shell $(uname -m)-$(uname -s | tr '[:upper:]' '[:lower:]')-gnu)
LDFLAGS_STATIC := -extldflags -static
TRIMPATH := -trimpath

VERSION := $(shell git describe --tags --always)
BUILDTIME := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS := -X "main.Version=$(VERSION)" -X "main.BuildTime=$(BUILDTIME)"

GO_BUILD := $(GO) build $(TRIMPATH) $(GO_BUILDMODE_STATIC) \
	$(EXTRA_FLAGS) -ldflags "$(LDFLAGS) $(LDFLAGS_STATIC) $(EXTRA_LDFLAGS)"

ll-killer: *.go build-aux/* ptrace/*.go apt.conf.d/* libfuse-overlayfs.a libgnu.a
	$(GO_BUILD) -o $@ .

fuse-overlayfs/Makefile: fuse-overlayfs/configure.ac fuse-overlayfs/Makefile.am
	cd fuse-overlayfs;\
	git apply --check ../patches/fuse-overlayfs.patch -q && git apply ../patches/fuse-overlayfs.patch;\
	./autogen.sh;\
	LIBS="-ldl" LDFLAGS="-static" ./configure --host=$(TARGET);

libfuse-overlayfs.a libgnu.a: fuse-overlayfs/Makefile fuse-overlayfs/*.c fuse-overlayfs/*.h
	make -C fuse-overlayfs
	cp fuse-overlayfs/lib/libgnu.a \
	   fuse-overlayfs/libfuse-overlayfs.a .
