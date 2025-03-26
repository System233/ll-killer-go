
LAYER_ROOT ?= $(HOME)/.cache/linglong-builder
OSTREE_ROOT ?= $(LAYER_ROOT)/repo
OSTREE_NAME ?= stable
OSTREE_REMOTE ?= https://mirror-repo-linglong.deepin.com/repos/stable
OSTREE=ostree --repo=$(OSTREE_ROOT)

ID ?= 
MODULE ?= binary

REF_ID = $(shell echo "$(ID)" | sed -E -e 's@:@/@')
REF_FILTER = $(shell echo "$(REF_ID)/($(MODULE)|runtime)" | sed -E -e 's:(\/[0-9]+\.[0-9]+\.[0-9]+)/:\1.*/:')
REF_REMOTE_NAME = $(OSTREE) remote refs $(OSTREE_NAME)|grep -P "$(REF_FILTER)"|tail -n1

OSTREE_REF_DIR = $(OSTREE_ROOT)/../layers/$(REF_ID)
OSTREE_MODULE_DIR = $(OSTREE_REF_DIR)/$(MODULE)
OSTREE_TARGET_DIR = $(OSTREE_MODULE_DIR)/files

$(OSTREE_ROOT):
	mkdir -p $(OSTREE_ROOT)
	$(OSTREE) init --mode=bare-user-only 
	$(OSTREE) remote add $(OSTREE_NAME) $(OSTREE_REMOTE) --no-gpg-verify 

$(OSTREE_MODULE_DIR): $(if $(wildcard $(OSTREE_ROOT)),,$(OSTREE_ROOT))
	mkdir -p $(OSTREE_REF_DIR)
	REF_NAME=$$($(REF_REMOTE_NAME))&&\
	$(OSTREE) pull "$${REF_NAME}"&&\
	$(OSTREE) checkout "$${REF_NAME}" $(OSTREE_MODULE_DIR)

clean:
	rm -rf $(OSTREE_TARGET)

show: $(OSTREE_MODULE_DIR)
	@echo $(OSTREE_TARGET_DIR)

all: $(OSTREE_MODULE_DIR)

.PHONY: all show clean
.DEFAULT_GOAL := all