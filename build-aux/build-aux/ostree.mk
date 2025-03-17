OSTREE_ROOT ?= $(shell realpath ~/.cache/linglong-builder/repo)
OSTREE_NAME ?= stable
OSTREE_REMOTE ?= https://mirror-repo-linglong.deepin.com/repos/stable
OSTREE=ostree --repo=$(OSTREE_ROOT)

ID ?= 
MODULE ?= binary
ifneq ($(ID),)
REF_ID = $(shell echo "$(ID)" | sed -E -e 's@:@/@')
REF_FILTER = $(shell echo "$(REF_ID)/($(MODULE)|runtime)" | sed -E -e 's:(\/[0-9]+\.[0-9]+\.[0-9]+)/:\1.*/:')
REF_REMOTE_NAME = $(shell $(OSTREE) remote refs $(OSTREE_NAME)|grep -P "$(REF_FILTER)"|tail -n1)

ifneq ($(REF_REMOTE_NAME),)
OSTREE_REF_DIR = $(OSTREE_ROOT)/../layers/$(REF_ID)
OSTREE_MODULE_DIR = $(OSTREE_REF_DIR)/$(MODULE)
OSTREE_TARGET_DIR = $(OSTREE_MODULE_DIR)/files
else
$(error "找不到指定的依赖:ID=$(ID) MODULE=$(MODULE)")
endif

endif


$(OSTREE_ROOT):
	mkdir -p $(OSTREE_ROOT)
	$(OSTREE) init --mode=bare-user-only 
	$(OSTREE) remote add $(OSTREE_NAME) $(OSTREE_REMOTE) --no-gpg-verify 

$(OSTREE_MODULE_DIR): $(if($(wildcard $(OSTREE_ROOT))),,$(OSTREE_ROOT))
	mkdir -p $(OSTREE_REF_DIR)
	$(eval REF_NAME := $(REF_REMOTE_NAME))
	$(OSTREE) pull $(REF_NAME)
	$(OSTREE) checkout $(REF_NAME) $(OSTREE_MODULE_DIR)
	touch $(OSTREE_MODULE_DIR)

clean:
	rm -rf $(OSTREE_TARGET)

show: $(OSTREE_MODULE_DIR)
	@echo $(OSTREE_TARGET_DIR)

all: $(OSTREE_MODULE_DIR)

.PHONY: all show clean
.DEFAULT_GOAL := all