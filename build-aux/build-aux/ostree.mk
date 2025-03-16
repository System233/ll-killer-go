
TARGET ?= ost.mk
-include $(TARGET)
OSTREE_ROOT ?= $(shell realpath ~/.cache/linglong-builder/repo/)
OSTREE_NAME ?= stable
OSTREE_REMOTE ?= https://mirror-repo-linglong.deepin.com/repos/stable
OSTREE=ostree --repo=$(OSTREE_ROOT)

ID ?= 
MODULE ?= binary
ifneq ($(ID),)
	REF ?= $(ID)/$(MODULE)
	REF_FILTER = $(shell echo "$(REF)" | sed -E -e 's\#:\#/\#' -e 's:(\/[0-9]+\.[0-9]+\.[0-9]+)/:\1.*/:')
	REF_REMOTE_NAME ?= $(shell $(OSTREE) remote refs $(OSTREE_NAME)|grep "$(REF_FILTER)"|tail -n1)
	REF_NAME ?= $(shell echo "$(REF_REMOTE_NAME)"|cut -d: -f2)
	OSTREE_TARGET ?= $(OSTREE_ROOT)/layers/$(REF_NAME)
	OSTREE_TARGET_FILES ?= $(OSTREE_TARGET)/files
else
	TARGET=
endif

$(OSTREE_ROOT):
	mkdir -p $(OSTREE_ROOT)
	$(OSTREE) init --mode=bare-user-only 
	$(OSTREE) remote add $(OSTREE_NAME) $(OSTREE_REMOTE) --no-gpg-verify 

$(OSTREE_TARGET_FILES):
	mkdir -p $(OSTREE_TARGET)
	rm -df $(OSTREE_TARGET)
	$(OSTREE) checkout $(REF_REMOTE_NAME) $(OSTREE_TARGET) || true
	test -d $(OSTREE_TARGET_FILES)

$(TARGET): $(OSTREE_TARGET_FILES)
	@echo "REF ?= $(REF)" >$@
	@echo "REF_FILTER ?= $(REF_FILTER)" >>$@
	@echo "REF_REMOTE_NAME ?= $(REF_REMOTE_NAME)" >>$@
	@echo "REF_NAME ?= $(REF_NAME)" >>$@
	@echo "OSTREE_TARGET ?= $(OSTREE_TARGET)" >>$@
	@echo "OSTREE_ROOT ?= $(OSTREE_ROOT)" >>$@
	@echo "OSTREE_NAME ?= $(OSTREE_NAME)" >>$@
	@echo "OSTREE_REMOTE ?= $(OSTREE_REMOTE)" >>$@

clean:
	rm -rf $(OSTREE_TARGET)

show: 
	@echo $(OSTREE_TARGET_FILES)

all: $(TARGET)

.PHONY: all show
.DEFAULT_GOAL := all 