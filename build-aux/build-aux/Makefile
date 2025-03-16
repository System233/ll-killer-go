# Copyright (c) 2025 System233
# 
# This software is released under the MIT License.
# https://opensource.org/licenses/MIT
# 

SHELL=/bin/bash -o pipefail
CONFIG_TARGET ?= config.mk
-include $(CONFIG_TARGET)

PKG ?= app
APPID ?= $(PKG).linyaps

KILLER ?= $(if $(wildcard ll-killer),./ll-killer,ll-killer)
BUILD_ARGS ?=
CREATE_ARGS ?=
ENABLE_LDD_CHECK ?= 1
ENABLE_INSTALL ?=1
ENABLE_OSTREE ?=0
LDD_CHECK_MODE ?= fast

SOURCES_LIST = sources.list
UPDATE_TARGET ?= apt-update.log
INSTALL_TARGET ?= apt-install.log
BUILD_TARGET ?= build.log
POST_BUILD_TARGET ?= post-build.log
PKG_INFO ?= pkg.info
YAML_CONFIG ?= linglong.yaml
EXTRA_DEPS ?= deps.list
EXTRA_BUILD ?= build.sh
EXTRA_POST_BUILD ?= post-build.sh
LDD_CHECK_TARGET ?= ldd-check.log
LDD_FOUND_TARGET ?= ldd-found.log
LDD_NOTFOUND_TARGET ?= ldd-notfound.log
LDD_SEARCH_TARGET ?= "$(LDD_FOUND_TARGET)" "$(LDD_NOTFOUND_TARGET)"
LDD_INSTALL_TARGET ?= apt-install-extra.log

APT_FILE_TARGET := apt.conf.d/apt-file.conf
BUILD_AUX_DEPS := build-aux/*
APPEND_LOG = 2>&1 | tee $@~
MOVE_LOG = @mv -f $@~ $@


ifeq ($(wildcard $(YAML_CONFIG)),$(YAML_CONFIG))
	LAYER_NAME := $(shell $(KILLER) layer build --print-layer-name)
else
	LAYER_NAME := $(APPID).layer
endif


ifeq ($(wildcard $(EXTRA_BUILD)),)
    EXTRA_BUILD_TARGET =
else
    EXTRA_BUILD_TARGET = $(BUILD_TARGET)
endif

OSTREE_MAKE_ARGS ?=
OSTREE_MAKE = $(MAKE) --no-print-directory -f build-aux/ostree.mk $(OSTREE_MAKE_ARGS)
OSTREE_BASE ?= base.mk
OSTREE_RUNTIME ?= runtime.mk

ifeq ($(ENABLE_OSTREE),1)
	APP_BASE_DIR ?=
	APP_RUNTIME_DIR ?=
	APP_BASE=$(shell $(KILLER) layer build --print base)
	APP_RUNTIME=$(shell $(KILLER) layer build --print runtime)
	APP_MODULE=$(shell $(KILLER) layer build --print module)

	OSTREE_MAKE_BASE = $(OSTREE_MAKE) TARGET=$(OSTREE_BASE) ID=$(APP_BASE) MODLUE=$(APP_MODULE)
	OSTREE_MAKE_RUNTIME = $(OSTREE_MAKE) TARGET=$(OSTREE_RUNTIME) ID=$(APP_RUNTIME) MODLUE=$(APP_MODULE)

	APP_BASE_DIR=$(shell $(OSTREE_MAKE_BASE) show)
	APP_RUNTIME_DIR=$(shell $(OSTREE_MAKE_RUNTIME) show)

	OVERLAY_ROOT=$(PWD)/linglong/filesystem/
	OVERLAY_MERGED=$(OVERLAY_ROOT)/merged
	OVERLAY_DIFF=$(OVERLAY_ROOT)/diff
	OVERLAY_OVERWRITE=$(OVERLAY_ROOT)/overwrite
	OVERLAY_WORK=$(OVERLAY_ROOT)/work

	OSTREE_CREATE_ARGS = --no-build
	OSTREE_BUILD_ARGS = --rootfs "$(APP_BASE_DIR)" --runtime "$(APP_RUNTIME_DIR)"
	EXTRA_BUILD_TARGET += $(OSTREE_BASE) $(OSTREE_RUNTIME)
endif

$(OSTREE_BASE): $(YAML_CONFIG)
	@echo "[拉取Base]"
	@$(OSTREE_MAKE_BASE)  $(OSTREE_BASE)
$(OSTREE_RUNTIME): $(YAML_CONFIG)
	@echo "[拉取Runtime]"
	@$(OSTREE_MAKE_RUNTIME) $(OSTREE_RUNTIME) 


ifeq ($(wildcard $(EXTRA_POST_BUILD)),)
    EXTRA_POST_BUILD_TARGET =
else
    EXTRA_POST_BUILD_TARGET = $(POST_BUILD_TARGET)
endif

ifeq ($(ENABLE_INSTALL),1)
	INSTALL_DEP := $(INSTALL_TARGET)
else
	INSTALL_DEP :=$(YAML_CONFIG) $(EXTRA_BUILD_TARGET)
endif

ifeq ($(ENABLE_LDD_CHECK),1)
	LAYER_DEP := $(LDD_INSTALL_TARGET)
else
	LAYER_DEP := $(INSTALL_DEP)
endif

KILLER_BUILD_ARGS = $(BUILD_ARGS) $(OSTREE_BUILD_ARGS)
KILLER_CREATE_ARGS = $(CREATE_ARGS) $(OSTREE_CREATE_ARGS)

$(CONFIG_TARGET):
	@echo "PKG ?= $(PKG)" > $(CONFIG_TARGET)
	@echo "APPID ?= $(APPID)" >> $(CONFIG_TARGET)
	@echo "KILLER_EXEC ?= $(KILLER_EXEC)" >> $(CONFIG_TARGET)
	@echo "CREATE_ARGS ?= $(CREATE_ARGS)" >> $(CONFIG_TARGET)
	@echo "BUILD_ARGS ?= $(BUILD_ARGS)" >> $(CONFIG_TARGET)
	@echo "ENABLE_LDD_CHECK ?= $(ENABLE_LDD_CHECK)" >> $(CONFIG_TARGET)
	@echo "ENABLE_INSTALL ?= $(ENABLE_INSTALL)" >> $(CONFIG_TARGET)
	@echo "LDD_CHECK_MODE ?= $(LDD_CHECK_MODE)" >> $(CONFIG_TARGET)
	@echo "UPDATE_TARGET ?= $(UPDATE_TARGET)" >> $(CONFIG_TARGET)
	@echo "INSTALL_TARGET ?= $(INSTALL_TARGET)" >> $(CONFIG_TARGET)
	@echo "EXTRA_BUILD ?= $(EXTRA_BUILD)" >> $(CONFIG_TARGET)
	@echo "EXTRA_POST_BUILD ?= $(EXTRA_POST_BUILD)" >> $(CONFIG_TARGET)
	@echo "BUILD_TARGET ?= $(BUILD_TARGET)" >> $(CONFIG_TARGET)
	@echo "POST_BUILD_TARGET ?= $(POST_BUILD_TARGET)" >> $(CONFIG_TARGET)
	@echo "PKG_INFO ?= $(PKG_INFO)" >> $(CONFIG_TARGET)
	@echo "YAML_CONFIG ?= $(YAML_CONFIG)" >> $(CONFIG_TARGET)
	@echo "EXTRA_DEPS ?= $(EXTRA_DEPS)" >> $(CONFIG_TARGET)
	@echo "LDD_CHECK_TARGET ?= $(LDD_CHECK_TARGET)" >> $(CONFIG_TARGET)
	@echo "LDD_NOTFOUND_TARGET ?= $(LDD_NOTFOUND_TARGET)" >> $(CONFIG_TARGET)
	@echo "LDD_FOUND_TARGET ?= $(LDD_FOUND_TARGET)" >> $(CONFIG_TARGET)
	@echo "LDD_SEARCH_TARGET ?= $(LDD_SEARCH_TARGET)" >> $(CONFIG_TARGET)
	@echo "LDD_INSTALL_TARGET ?= $(LDD_INSTALL_TARGET)" >> $(CONFIG_TARGET)
	

$(SOURCES_LIST):
	@echo "[从主机复制sources.list]"
	cp /etc/apt/sources.list $(SOURCES_LIST)

$(UPDATE_TARGET): $(SOURCES_LIST)
	@echo "[更新APT缓存]"
	@test "$(ENABLE_LDD_CHECK)" != "1" && rm -vf $(APT_FILE_TARGET) || true
	$(KILLER) apt -- apt update -y $(APPEND_LOG)
	$(MOVE_LOG)

$(PKG_INFO): $(UPDATE_TARGET)
	@echo "[获取包元数据]"
	$(KILLER) apt -- apt show $(PKG) > $(PKG_INFO)~
	$(MOVE_LOG)

$(EXTRA_DEPS):
	touch $(EXTRA_DEPS)

$(YAML_CONFIG): $(PKG_INFO)
	@echo "[创建玲珑项目]"
	$(KILLER) create --from $(PKG_INFO) --id $(APPID) $(KILLER_CREATE_ARGS) 

$(BUILD_TARGET): $(YAML_CONFIG) $(EXTRA_BUILD)
	@echo "[执行自定义构建]"
	$(KILLER) build $(KILLER_BUILD_ARGS) --ptrace -- $(EXTRA_BUILD) $(APPEND_LOG)
	$(MOVE_LOG)

$(INSTALL_TARGET): $(YAML_CONFIG) $(EXTRA_DEPS) $(EXTRA_BUILD_TARGET)
	@echo "[将包安装至构建环境]"
	$(KILLER) build $(KILLER_BUILD_ARGS) --ptrace -- apt install -y $(PKG) $(shell cat $(EXTRA_DEPS)) $(APPEND_LOG)
	$(MOVE_LOG)

$(POST_BUILD_TARGET): $(INSTALL_DEP) $(EXTRA_POST_BUILD)
	@echo "[执行自定义后构建]"
	$(KILLER) build $(KILLER_BUILD_ARGS) --ptrace -- $(EXTRA_POST_BUILD) $(APPEND_LOG)
	$(MOVE_LOG)

$(LDD_CHECK_TARGET):  $(INSTALL_DEP) $(EXTRA_POST_BUILD_TARGET) build-aux/ldd-check.sh
	@echo "[检查缺失库]"
	$(KILLER) build $(KILLER_BUILD_ARGS) -- build-aux/ldd-check.sh $(LDD_CHECK_MODE) > $(LDD_CHECK_TARGET)~
	$(MOVE_LOG)

.INTERMEDIATE: ldd_search
$(LDD_FOUND_TARGET) $(LDD_NOTFOUND_TARGET): ldd_search

ldd_search: $(LDD_CHECK_TARGET) build-aux/ldd-search.sh
	@echo "[搜索缺失库所在包]"
	$(KILLER) apt -- build-aux/ldd-search.sh $(LDD_CHECK_TARGET) $(LDD_FOUND_TARGET)~ $(LDD_NOTFOUND_TARGET)~
	@mv -f $(LDD_FOUND_TARGET)~ $(LDD_FOUND_TARGET)
	@mv -f $(LDD_NOTFOUND_TARGET)~ $(LDD_NOTFOUND_TARGET)

$(LDD_INSTALL_TARGET): $(LDD_FOUND_TARGET)
	@echo "[安装找到的缺失库]"
	$(KILLER) build $(KILLER_BUILD_ARGS) --ptrace -- apt install -y $(shell cat $(LDD_FOUND_TARGET)) $(APPEND_LOG)
	$(MOVE_LOG)

$(LAYER_NAME): $(LAYER_DEP) $(EXTRA_POST_BUILD_TARGET) $(BUILD_AUX_DEPS)
	@echo "[导出layer文件]"
	$(KILLER) layer build -o $(LAYER_NAME)

build: $(YAML_CONFIG) $(EXTRA_DEPS) $(EXTRA_BUILD_TARGET)
	@$(KILLER) build $(KILLER_BUILD_ARGS) --ptrace -- bash

clean-apt:
	rm -rf linglong/apt
clean-fs:
	rm -rf linglong/filesystem
clean-test:
	rm -rf tests
clean-ll:
	rm -rf linglong/output
	ll-builder list|grep ":$(APPID)/"|xargs -r ll-builder remove || true
clean-log:
	rm -rf *.log 
clean-layer:
	rm -rf *.layer
clean: clean-ll clean-fs clean-test clean-log clean-layer
clean-mk:
	rm -rf $(CONFIG_TARGET)
clean-all: clean clean-apt

test-dynamic: $(LAYER_NAME)
	@echo "[执行自动化测试]"
	build-aux/test.sh $(APPID) $(LAYER_NAME) 2>&1 | tee test.log
test: $(YAML_CONFIG)
	@$(MAKE) --no-print-directory test-dynamic

help:
	@echo "  ll-killer - 基于Makefile的玲珑项目管理工具"
	@echo ""
	@echo "  用法:"
	@echo "    make [目标] [变量=值 ...]"
	@echo ""
	@echo "  目标:"
	@echo "    all           构建layer文件（默认目标）"
	@echo "    config        生成或更新配置文件 $(CONFIG_TARGET)"
	@echo "    init          初始化项目并创建 $(YAML_CONFIG)"
	@echo "    build         进入当前配置下的构建环境"
	@echo "    update        更新APT软件源缓存"
	@echo "    install       安装主包和额外依赖到构建环境"
	@echo "    layer         构建layer文件"
	@echo "    clean         清除所有生成的文件、测试输出和layer"
	@echo "    clean-apt     清除APT相关缓存文件"
	@echo "    clean-fs      清除构建文件系统内容"
	@echo "    clean-ll      清除玲珑构建输出"
	@echo "    clean-layer   清除layer文件"
	@echo "    clean-log     清除log文件"
	@echo "    clean-test    清除tests文件"
	@echo "    clean-mk      清除$(CONFIG_TARGET)配置文件"
	@echo "    clean-all     执行clean并清除clean-apt"
	@echo "    test          执行自动化测试"
	@echo "    help          显示本帮助信息"
	@echo ""
	@echo "  配置变量:  KEY=VALUE 说明 (当前值)"
	@echo "    PKG=name                 要打包的Debian软件包名称 ($(PKG))"
	@echo "    APPID=name.linyaps       玲珑应用ID ($(APPID))"
	@echo "    KILLER=ll-killer         ll-killer工具路径 ($(KILLER))"
	@echo "    ENABLE_LDD_CHECK=0|1     是否启用依赖检查 ($(ENABLE_LDD_CHECK))"
	@echo "    ENABLE_INSTALL=0|1       是否启用自动安装依赖 ($(ENABLE_INSTALL))"
	@echo "    ENABLE_OSTREE=0|1        是否启用自定义ostree仓库，彻底移除ll-builder依赖 ($(ENABLE_OSTREE))"
	@echo "    LDD_CHECK_MODE=fast|full 依赖检查模式 fast=快速/full=全量 ($(LDD_CHECK_MODE))"
	@echo "    CREATE_ARGS=             传递给ll-killer create的额外参数 ($(CREATE_ARGS))"
	@echo "    BUILD_ARGS=              传递给ll-killer build的额外参数 ($(BUILD_ARGS))"
	@echo "    OSTREE_MAKE_ARGS=        传递给ostree.mk的额外参数 ($(OSTREE_MAKE_ARGS))"
	@echo ""
	@echo "  典型工作流:"
	@echo "    0. 准备自定义的sources.list或留空，自动从主机复制"
	@echo "    1. make config PKG=your-package  # 初始化配置"
	@echo "    2. make init                     # 创建项目配置"
	@echo "    3. make layer                    # 构建最终layer文件"
	@echo ""
	@echo "  示例:"
	@echo "    # 构建GIMP的玲珑包 (默认从主机复制sources.list)"
	@echo "    make PKG=gimp APPID=gimp.linyaps"
	@echo "    "
	@echo "    # 指定Base等项目创建参数"
	@echo "    make PKG=gimp CREATE_ARGS=\"--base=org.deepin.base/23.1.0\""
	@echo "    "
	@echo "    # 禁用LDD检查"
	@echo "    make ENABLE_LDD_CHECK=0"
	@echo ""
	@echo "  特殊文件:"
	@echo "    apt.conf.d      APT配置文件夹"
	@echo "    auth.conf.d     APT授权配置文件夹（遇到APT源返回401时使用）"
	@echo "    build-aux       辅助构建脚本目录"
	@echo "    build.sh        自定义构建脚本"
	@echo "    post-build.sh   自定义后构建脚本"
	@echo "    build-aux       辅助构建脚本目录"
	@echo "    sources.list    使用的APT源配置"
	@echo "    sources.list.d  使用的APT源配置文件夹"
	@echo "    linglong.yaml   项目配置文件"
	@echo "    ll-killer       ll-killer可执行文件"
	@echo "    Makefile        基于Makefile的玲珑项目构建规则"
	@echo "    *.log           各阶段构建日志"
	@echo "    *.layer         生成的layer文件"

.DEFAULT:
	@echo "错误: 不支持选定的目标，请运行 $(MAKE) help 查看帮助信息。" >&2

init: $(YAML_CONFIG)
update: $(UPDATE_TARGET)
install: $(BUILD_TARGET)
config: $(CONFIG_TARGET)
layer-dynamic: $(LAYER_NAME)
layer: $(YAML_CONFIG)
	@$(MAKE) --no-print-directory layer-dynamic
all: layer

.PHONY: all build config clean-mk init update install layer layer-dynamic clean clean-all clean-apt clean-test clean-log clean-layer clean-ll clean-fs help test test-dynamic
.DEFAULT_GOAL := all 
