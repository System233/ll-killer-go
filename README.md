# ll-killer-go - 玲珑杀手Go <!-- omit from toc -->

[![Build](https://github.com/System233/ll-killer-go/actions/workflows/build.yaml/badge.svg)](https://github.com/System233/ll-killer-go/actions/workflows/build.yaml) 

## 项目简介 <!-- omit from toc -->

本项目是[Linglong Killer Self-Service (ll-killer 玲珑杀手)](https://github.com/System233/linglong-killer-self-service)的重写版本，去除了构建阶段的shell脚本，全部采用go实现，并添加了一些增强功能。

`ll-killer-go` 是一款专为解决玲珑容器应用构建问题而设计的命令行工具。它帮助开发者快速创建、构建和生成玲珑容器应用项目，同时提供一整套辅助构建与调试功能。使用该工具，用户可以在玲珑容器中获得与传统容器（如 `docker`）类似的构建体验，包括对特权命令的支持和 `apt` 安装等功能，免去了手动解压软件包（如 `deb`）和修复依赖库查找路径的繁琐操作。`ll-killer-go` 通过重建自由的容器环境，确保构建过程的一致性与可靠性，极大地提高了开发效率和可维护性。

### ⚠️文档滞后，新功能请查看 [更新日志](CHANGELOG.md) 或命令帮助

### 功能概述
- **创建项目**：自动创建项目所需的配置文件和辅助脚本，初始化项目环境。
- **构建项目**：进入构建环境，支持以 root 权限执行各种命令。
- **提交内容**：将构建后的文件提交到玲珑容器，确保文件路径和图标等修复。
- **辅助构建工具**：包括缺失库检查、符号链接修复、图标文件修复等多种实用功能。
- **运行与调试**：提供一个一致的运行时环境，支持容器应用的调试和测试。

### 特色功能

- **隔离的 APT 环境**：提供独立的 APT 环境，让你可以轻松查找和安装所需的软件包，而无需担心宿主机环境的干扰和权限问题。
- **可重用的构建环境**：构建环境在命令执行结束后依然保留，你可以多次进入容器进行调整，直到完全满足需求，避免重复创建环境的麻烦。
- **root 权限模拟**：尽管在构建环境中能够执行特权命令（如 `apt install`），但无需担心安全问题。工具使用宿主机用户权限进行操作，避免了实际 root 权限的风险。
- **全局可写容器环境**：在构建环境中，你可以自由修改任何位置，提供灵活的操作空间，方便调试和调整。
- **自动修复图标和快捷方式**：`build-aux` 工具集内置了多种自动修复脚本，帮助你快速修复常见问题，如图标和快捷方式的更新。
- **自动检查缺失库**：使用 `ll-killer script build-aux/build-and-check.sh` 命令可以自动检测并修复构建过程中缺失的运行时库，确保容器应用能顺利运行。
  

## 目录 <!-- omit from toc -->


- [快速教程](#快速教程)
  - [1. 获取与配置 `ll-killer`](#1-获取与配置-ll-killer)
    - [全局安装 `ll-killer`（可选）](#全局安装-ll-killer可选)
  - [2. 开始打包一个应用](#2-开始打包一个应用)
    - [创建打包环境](#创建打包环境)
    - [配置软件源](#配置软件源)
    - [创建 `config.mk` 配置构建参数](#创建-configmk-配置构建参数)
    - [指定依赖（可选）](#指定依赖可选)
    - [开始构建](#开始构建)
    - [测试构建结果](#测试构建结果)
    - [安装构建结果](#安装构建结果)
    - [清理构建缓存（可选）](#清理构建缓存可选)
  - [3. 挂载模式说明](#3-挂载模式说明)
    - [`fuse-overlayfs` 模式](#fuse-overlayfs-模式)
    - [`merge` 模式（默认回退）](#merge-模式默认回退)
    - [模式对比](#模式对比)
    - [构建过程中遇到目录写入问题？](#构建过程中遇到目录写入问题)
- [安装与配置](#安装与配置)
  - [获取 ll-killer](#获取-ll-killer)
  - [获取静态fuse-overlayfs（可选）](#获取静态fuse-overlayfs可选)
  - [环境要求](#环境要求)
- [命令概览](#命令概览)
- [各命令详细介绍](#各命令详细介绍)
  - [1. `apt` — 进入隔离的 APT 环境](#1-apt--进入隔离的-apt-环境)
  - [2. `create` — 创建玲珑项目](#2-create--创建玲珑项目)
    - [从APT元数据创建项目](#从apt元数据创建项目)
  - [2. `build` — 构建或进入构建环境](#2-build--构建或进入构建环境)
  - [3. `exec` — 进入运行时环境](#3-exec--进入运行时环境)
  - [4. `run` — 启动容器](#4-run--启动容器)
  - [5. `commit` — 提交构建内容](#5-commit--提交构建内容)
  - [6. `clean` — 清除构建内容](#6-clean--清除构建内容)
  - [7. `build-aux` — 创建辅助构建脚本](#7-build-aux--创建辅助构建脚本)
  - [8. `ptrace` — 修正系统调用](#8-ptrace--修正系统调用)
  - [9. `script` — 执行自定义构建脚本](#9-script--执行自定义构建脚本)
  - [注意事项](#注意事项)
- [挂载相关功能](#挂载相关功能)
    - [1. 进入运行时环境并挂载文件系统](#1-进入运行时环境并挂载文件系统)
    - [例 1: 使用 `merge` 合并文件系统](#例-1-使用-merge-合并文件系统)
    - [例 2: 挂载源路径到目标路径，指定用户和组](#例-2-挂载源路径到目标路径指定用户和组)
    - [例 3: 使用自定义的 Unix 套接字和合并根目录路径](#例-3-使用自定义的-unix-套接字和合并根目录路径)
    - [例 4: 使用不同的挂载选项](#例-4-使用不同的挂载选项)
  - [2. 挂载选项详解](#2-挂载选项详解)
    - [基本语法](#基本语法)
    - [支持的挂载标志](#支持的挂载标志)
    - [合并挂载（Merge Mount）](#合并挂载merge-mount)
  - [例：使用合并挂载](#例使用合并挂载)
  - [3. 注意事项](#3-注意事项)
- [高级进阶玩法](#高级进阶玩法)
  - [1. 脱离ll-builder/ll-box启动一个容器](#1-脱离ll-builderll-box启动一个容器)
    - [1.1 具有Root的玲珑开发环境](#11-具有root的玲珑开发环境)
    - [1.2 Rootless的玲珑开发环境](#12-rootless的玲珑开发环境)
    - [1.3 基于Ptrace+Rootless的玲珑开发环境](#13-基于ptracerootless的玲珑开发环境)
- [疑难解答](#疑难解答)
- [贡献与维护](#贡献与维护)
- [许可](#许可)


## 快速教程

本章节提供一个基于`Makefile`的简明的打包教程，结合实际案例帮助快速掌握流程。

---

### 1. 获取与配置 `ll-killer`

首先，确保系统已安装必要的依赖：

```sh
sudo apt install apt-file fuse-overlayfs
```

然后，下载 `ll-killer` 并赋予执行权限：

```sh
wget https://github.com/System233/ll-killer-go/releases/latest/download/ll-killer-amd64 -O ll-killer
chmod +x ll-killer
```

#### 全局安装 `ll-killer`（可选）

为了方便使用，可将 `ll-killer` 安装到 `~/.local/bin`：

```sh
mkdir -p ~/.local/bin
mv ll-killer ~/.local/bin
```

如果 `~/.local/bin` 未添加至 `PATH`，请执行：

```sh
echo 'export PATH=$HOME/.local/bin:$PATH' >>~/.bashrc
source ~/.bashrc
```

---

### 2. 开始打包一个应用

接下来的步骤默认 `ll-killer` 已加入 `PATH`。

#### 创建打包环境

```sh
# 创建工作目录并进入
mkdir gimp && cd gimp

# 生成辅助构建文件
ll-killer build-aux
```

#### 配置软件源

创建 `sources.list` 并填充 Deepin V23 的软件源：

```sh
cat >sources.list <<EOF
 deb [trusted=yes] https://community-packages.deepin.com/deepin/beige beige main commercial community
 deb [trusted=yes] https://com-store-packages.uniontech.com/appstorev23 beige appstore
 deb [trusted=yes] https://community-packages.deepin.com/driver-23/ driver non-free
EOF
```

#### 创建 `config.mk` 配置构建参数

```sh
make config PKG="gimp" KILLER_CREATE_ARGS="--base org.deepin.base/23.1.0"
```

**参数说明：**

- `PKG`：Deb 包名。
- `APPID`：玲珑包名，默认值为 `${PKG}.linyaps`。
- `KILLER_CREATE_ARGS`：传递给 `ll-killer create` 的参数，应确保 `base` 与 `sources.list` 匹配。
- `KILLER_BUILD_ARGS`：传递给 `ll-killer build` 的参数，可根据需要调整。

#### 指定依赖（可选）

可以创建 `deps.list`，并在其中列出依赖包名，每行一个。

#### 开始构建

```sh
make
```

构建完成后，最终会生成两个 `.layer` 文件，仅需 `*_binary.layer`。  

**提示:** 运行`make -n`来了解`make`构建流程，输出示例：
```sh
echo "[获取包元数据]"
ll-killer apt -- apt show gimp > pkg.info
echo "[创建玲珑项目]"
ll-killer create --from pkg.info --id "gimp.linyaps"  
touch deps.list
echo "[将包安装至构建环境]"
ll-killer build --ptrace -- apt install -y gimp  | tee apt-install.log
echo "[检查缺失库]"
ll-killer build -- build-aux/ldd-check.sh > ldd-check.log
echo "[搜索缺失库所在包]"
ll-killer apt -- build-aux/ldd-search.sh ldd-check.log ldd-found.log ldd-notfound.log;
echo "[安装找到的缺失库]"
ll-killer build -- apt install -y  | tee apt-install-extra.log;
echo "[提交文件到容器中]"
ll-killer commit -- --skip-output-check --skip-strip-symbols
echo "[导出layer文件]"
ll-killer export -- --layer
```

#### 测试构建结果
此命令创建容器环境的shell，请结合desktop中的Exec指令来测试程序。

```sh
ll-builder run
```

run命令默认执行`/opt/apps/APPID/files/entrypoint.sh`入口点。
如需进入原始玲珑容器环境，请运行：
```sh
ll-builder run --exec bash
```
#### 安装构建结果

```sh
# 玲珑 >= 1.7.x 版本需要 sudo 权限
ll-cli install *_binary.layer
```

#### 清理构建缓存（可选）

如果需要清理构建过程中产生的缓存，可执行：

```sh
make clean-ll   # 删除玲珑缓存
make clean-fs   # 删除构建缓存
make clean-apt  # 删除 apt 缓存
make clean      # 删除所有缓存
```

---

### 3. 挂载模式说明

`ll-killer` 运行时支持 **`fuse-overlayfs`** 和 **`merge`** 两种根文件系统挂载模式。

#### `fuse-overlayfs` 模式

- 需确保 `fuse-overlayfs` 为 **静态编译** 或能直接在容器中运行。
- 需要将 `fuse-overlayfs` 二进制文件复制到 `build-aux` 目录。
- 启动速度快，但要求宿主机支持 `fuse` 内核模块。

#### `merge` 模式（默认回退）

- 若 `fuse-overlayfs` 不可用或运行失败，则自动回退至 `merge` 模式。
- 需要计算合并目录，目录冲突越多，合并时间越长。
- 兼容性好，作为兜底方案。

#### 模式对比

| 模式            | 启动速度 | 依赖 `fuse` | 兼容性 |
|---------------|--------|-----------|------|
| `fuse-overlayfs` | 快     | 是         | 低   |
| `merge`         | 慢     | 否         | 高   |

#### 构建过程中遇到目录写入问题？

请参考 [疑难解答](#疑难解答) 章节，并根据需要调整 `make` 传递的 `KILLER_BUILD_ARGS` 参数。

---



## 安装与配置

### 获取 ll-killer

你可以下载预编译的二进制，或手动编译本项目。

**使用预编译二进制**

在项目 [Release](https://github.com/System233/ll-killer-go/releases) 页下载预编译的二进制文件，一般使用amd64版本，下载后改名为`ll-killer`，并添加执行权限。   
上述步骤可使用以下命令一键完成：
```bash
wget https://github.com/System233/ll-killer-go/releases/latest/download/ll-killer-amd64 -O ll-killer
chmod +x ll-killer

./ll-killer -h
```

你可以将命令安装至`~/.local/bin`，以便随时使用。
```bash
mkdir -p ~/.local/bin
mv ./ll-killer ~/.local/bin
```

**手动编译**

1. 克隆或下载项目源码。  
   ```sh
   git clone https://github.com/System233/ll-killer-go.git
   cd ll-killer-go
   ```
2. 安装 Golang 环境。  
   ```sh
   sudo apt install golang
   ```
3. 使用 `make` 命令编译并生成可执行文件，默认生成主机架构的二进制。


### 获取静态fuse-overlayfs（可选）

如需启用fuse-overlayfs模式，需要使用以下命令静态编译[fuse-overlayfs](https://github.com/containers/fuse-overlayfs.git)。

```sh
git clone https://github.com/containers/fuse-overlayfs.git
cd fuse-overlayfs
autoreconf
automake --add-missing
autoreconf
LIBS="-ldl" LDFLAGS="-static" ./configure 
make -j8
```

### 环境要求
- Linux 系统（支持多种发行版）
- Go 编译环境
- 主机必须安装 `fuse-overlayfs`/`linglong-bin`
- 主机可选安装：`apt`/`apt-file`，用于apt相关和依赖查找功能。

## 命令概览

`ll-killer-go` 提供了一些核心命令来帮助构建和管理玲珑容器应用：

- `apt`：进入隔离的 APT 环境。
- `build`：构建或进入构建环境。
- `exec`：进入运行时环境。
- `run`：启动容器，执行应用。
- `commit`：提交构建内容到玲珑容器。
- `clean`：清除构建内容。
- `build-aux`：创建辅助构建脚本。
- `ptrace`：修正系统调用（目前仅支持 `chown` 调用）。
- `script`：执行自定义构建流程。
- `help`：显示帮助信息。
  

## 各命令详细介绍

### 1. `apt` — 进入隔离的 APT 环境

**应用场景**：此命令用于在构建环境中隔离 APT 操作，确保不会污染宿主机环境。在隔离环境中，你可以使用 `apt-file` 或内置工具（如 `ldd-search.sh`）来查找包依赖。

**用法**：
```bash
ll-killer apt -- <command> [arguments...]
```

例如：
```bash
ll-killer apt -- bash
```

**注意事项**：
- 当前目录下的 `apt.conf`, `apt.conf.d`, `sources.list`, 和 `sources.list.d` 文件会被挂载到 `/etc` 目录，可以在这些文件中自定义 APT 配置。
- 隔离环境中的 APT 缓存会被构建容器重用。


### 2. `create` — 创建玲珑项目

**应用场景**：`create` 命令用于创建新的玲珑容器应用项目。你可以通过此命令生成一个基础项目框架，并根据需要自定义项目的构建命令、应用描述、运行时环境等。此命令还支持从APT Package元数据创建项目，并自动生成相关配置，省去手动查找和配置依赖包的麻烦。

用法：

```bash
ll-killer create [flags]
```

示例：

```bash
ll-killer create --name "MyApp" --version "1.0.0""
```

#### 从APT元数据创建项目

`ll-killer create` 支持通过`apt show`输出的包元数据创建项目，以下是创建GIMP图像处理应用项目的示例：
```bash
# 将 gimp 包的信息保存到 package.info
apt show gimp > package.info

# 从 package.info 提取信息并在当前目录创建玲珑项目
ll-killer create --name "MyApp" --version "1.0.0" --from "package.info"
```
该命令会根据 `apt show <pkg>` 输出的元数据创建应用项目，自动填充包名、描述、版本等信息。

**注意事项**

- 使用 `--from` 参数时，可以通过 `apt show <pkg>` 输出的元数据自动生成项目的基础配置信息。这对于从现有的APT包创建容器应用项目非常有用。
- 默认情况下，`ll-killer create`会在创建项目时自动构建一次。
- 若指定 `--no-build`，则只会生成项目的基本模板，而不会自动进行构建步骤。您可以手动根据需求运行构建命令。  
- 在从未对项目进行构建(如`ll-builder build`)的情况下，使用`ll-killer build`时必须关闭严格模式。

### 2. `build` — 构建或进入构建环境

**应用场景**：构建项目的核心命令。此命令可以进入构建环境，执行各种构建任务，确保以 root 权限执行命令，如 `apt` 安装等。你还可以自定义构建命令。

**用法**：
```bash
ll-killer build [flags] -- <build command> [arguments...]
```

例如：
```bash
ll-killer build -- make install
```

**注意事项**：
- 构建环境在退出后不会消失，你可以多次进入环境进行调整。
- 使用 `--strict` 标志可以启动严格模式，确保与运行时环境一致，但此模式下不包括开发工具（如 gcc）。
- 使用 `--strict` 标志启用严格模式时，项目必须至少使用`ll-builder build`构建一次
- 构建环境会挂载主机根文件系统，也支持临时根目录和自定义 `fuse-overlayfs` 配置。
- 主机必须安装`fuse-overlayfs`


### 3. `exec` — 进入运行时环境

**应用场景**：在项目构建完成后，此命令用于启动和测试容器环境。你可以根据需要进行调试，甚至自定义挂载和权限设置。

**用法**：
```bash
ll-killer exec [flags]
```

例如：
```bash
ll-killer exec -- bash
```

**挂载选项**：
- 支持多种挂载选项，例如：`--mount /source:/target:options`，其中 `options` 包括支持的挂载标志（如 `rw`, `nosuid` 等）和文件系统类型。
- 额外支持 `merge` 挂载类型，用于在无overlayfs支持下的堆叠文件系统。
- 额外支持 `fuse-overlayfs` 挂载类型，用于支持基于fuse的堆叠文件系统。
- 当应用运行在不支持fuse的环境中，或fuse-overlayfs不可用时，自动回退merge挂载模式。
详细信息参见[挂载相关功能](#挂载相关功能)


**注意事项**：
- 你可以使用 `--rootfs` 指定不同的根目录路径，支持合并多个目录。
- `--uid` 和 `--gid` 可用来指定使用的用户和组 ID。
- 此命令在`ll-box`容器/`ll-builder run/build`内部使用，你不应该在主机上运行此命令。

### 4. `run` — 启动容器

**应用场景**：启动构建后的容器进行应用的测试和运行。

**用法**：
```bash
ll-killer run [flags]
```

此命令等同于 `ll-builder run`，可用于提供一致的CLI体验。

### 5. `commit` — 提交构建内容

**应用场景**：将构建完成的应用提交到玲珑容器中，确保文件复制、快捷方式和图标修复等操作已完成。

**用法**：
```bash
ll-killer commit [flags]
```

此命令等同于 `ll-builder build`，可在 `linglong.yaml` 文件中对构建指令进行调整。

**注意事项**：
- 你可以在项目根目录提供一个`fuse-overlayfs`文件以启用`fuse-overlayfs`挂载模式

### 6. `clean` — 清除构建内容

**应用场景**：用于清除当前项目的构建环境内容，确保后续的构建不会受到先前内容的影响。

**用法**：
```bash
ll-killer clean [flags]
```

### 7. `build-aux` — 创建辅助构建脚本

**应用场景**：用于生成一系列辅助构建脚本，帮助开发者修复符号链接、图标路径、依赖问题等。

**用法**：
```bash
ll-killer build-aux [flags]
```

**包含的辅助脚本**：
- `entrypoint.sh`：玲珑容器的入口点。
- `env.sh`：环境配置脚本。
- `ldd-check.sh`：检查容器内缺失的库。
- `ldd-search.sh`：查找缺失库的所在包。
- `relink.sh`：修复容器内的符号链接。
- `setup-desktop.sh`：修复桌面文件中的图标和执行指令。
- `setup-filesystem.sh`：从构建环境复制文件到目标目录。
- `setup-icon.sh`：将图标文件按尺寸放置到指定位置。
- `setup.sh`：调用上述工具，将构建环境内容复制到玲珑容器中，并执行相关补丁。

### 8. `ptrace` — 修正系统调用

**应用场景**：此命令用于拦截并修正构建过程中的系统调用，目前仅支持修正 `chown` 调用，避免权限问题。

**用法**：
```bash
ll-killer ptrace -- <command> [arguments...]
```

例如：
```bash
ll-killer ptrace -- chown root:root /path/to/file
```

无论chown设置所有者为谁，都只会设置成自己，以确保该命令成功。  
此功能用于忽略apt安装某些库过程中可能出现的权限问题。

**注意事项**：
- 此模式将极大降低目标的性能



### 9. `script` — 执行自定义构建脚本

**应用场景**：用于执行自定义的构建流程，确保环境能够找到 `ll-killer` 二进制文件。

**用法**：
```bash
ll-killer script -- <build script> [arguments...]
```

此命令等同于：
```bash
KILLER_EXEC=path/to/ll-killer <build script> [arguments...]
```

### 注意事项
- 使用 `--help` 参数可以查看每个子命令的详细帮助信息。
- 确保在项目创建和构建阶段，`linglong.yaml` 配置文件已经正确设置。
- 严格模式下，构建环境中不会包含开发工具，如编译器等，仅包含运行时环境所需的最小工具集。

## 挂载相关功能

`ll-killer` 提供了强大的挂载功能，使你能够灵活地管理和调整容器中的文件系统。通过挂载选项，你可以将主机上的目录或文件与容器内的目录进行绑定、合并或覆盖。本文将详细介绍如何使用 `ll-killer` 的挂载功能以及常见的用法和注意事项。

`ll-killer`的`--mount`支持常见的`mount`挂载标志，支持的文件系统包括`overlay`、`tmpfs`、`devpts`，支持额外的`merge`和`fuse-overlayfs`文件系统挂载类型。

#### 1. 进入运行时环境并挂载文件系统

你可以使用 `--mount` 参数来挂载主机文件系统到容器内。以下是几个常见的挂载用法示例。

#### 例 1: 使用 `merge` 合并文件系统

在运行时环境中，使用当前用户的 UID 和 GID，合并挂载文件系统并使用 `/myrootfs` 的内容覆盖容器根文件系统。合并文件系统使你能够在没有内核 `overlayfs` 或 `fuse` 模块的支持下，堆叠多个源目录。

```bash
ll-killer exec --mount /+/myrootfs:/rootfs::merge -- bash
```

**解释**：
- `/+/myrootfs`: 这是合并挂载的源目录。多个源目录使用 `+` 符号分隔。
- `/rootfs`: 这是目标目录，即容器内的根文件系统将会被合并到此目录。
- `::merge`: 指定合并挂载类型。

#### 例 2: 挂载源路径到目标路径，指定用户和组

将主机上的 `/path/to/source` 目录挂载到容器内的 `/path/to/target` 目录，同时指定容器中运行进程的 UID 和 GID。

```bash
ll-killer exec --mount /path/to/source:/path/to/target --uid 2000 --gid 2000 --chroot=false
```

**解释**：
- `/path/to/source:/path/to/target`: 源路径与目标路径的映射关系。
- `--uid 2000 --gid 2000`: 指定容器内进程的 UID 和 GID。
- `--chroot=false`: 禁用 `chroot` 环境。

#### 例 3: 使用自定义的 Unix 套接字和合并根目录路径

挂载自定义的 Unix 套接字文件和合并根文件系统路径。该命令支持将 `/etc` 目录挂载到容器内。

```bash
ll-killer exec --socket /tmp/myapp.sock --rootfs /tmp/myapp.rootfs --mount /etc:/etc
```

**解释**：
- `--socket /tmp/myapp.sock`: 指定自定义的 Unix 套接字文件。
- `--rootfs /tmp/myapp.rootfs`: 指定合并后的根文件系统路径。
- `--mount /etc:/etc`: 挂载主机的 `/etc` 目录到容器的 `/etc`。

#### 例 4: 使用不同的挂载选项

你可以为挂载操作指定多个选项，例如读写权限、禁用某些功能等。

```bash
ll-killer exec --mount /data:/data:rw+nosuid --uid 1000 --gid 1000 --rootfs /var/run/myapp.rootfs
```

**解释**：
- `--mount /data:/data:rw+nosuid`: 将 `/data` 目录挂载为读写模式，并禁用 `setuid` 位。
- `--uid 1000 --gid 1000`: 指定 UID 和 GID。
- `--rootfs /var/run/myapp.rootfs`: 使用自定义根文件系统路径。

### 2. 挂载选项详解

#### 基本语法
挂载选项使用 `--mount` 参数，语法如下：

```bash
--mount 源目录:目标目录:挂载标志:文件系统类型:额外数据
```

- **源目录**：挂载的源目录，支持多个源目录时，使用 `+` 符号分隔。
- **目标目录**：容器内的目标目录，文件系统将被挂载到此位置。
- **挂载标志**：类似于系统的 `mount` 命令，支持多种挂载选项（如 `bind`、`rbind` 等）。
- **文件系统类型**：默认为 `none`，可以指定特定的文件系统类型。
- **额外数据**：可以指定排除某些目录等附加信息。

#### 支持的挂载标志

`ll-killer` 支持以下挂载标志，具体的标志根据你的需求选择使用：

- **bind**、**rbind**：常用的挂载标志，表示源目录和目标目录的绑定。
- **noexec**：禁止在挂载的文件系统上执行二进制文件。
- **nosuid**：禁止设置 `setuid` 位。
- **rw**、**ro**：指定文件系统的读写模式。
- **async**、**sync**：指定挂载时的异步或同步操作。
- **private**、**shared**：指定挂载的共享类型。

具体可以参考 Linux `mount` 命令的挂载标志，`ll-killer` 提供了大部分常见的挂载选项。

#### 合并挂载（Merge Mount）

`ll-killer` 提供了一种特殊的挂载类型——合并挂载（`merge`）。这使你能够将多个目录堆叠在一起，在不依赖内核的 `overlayfs` 或 `fuse` 模块的情况下实现文件系统层叠。

**语法**：
```bash
--mount 源目录1+源目录2+源目录N:目标目录:挂载标志:merge:排除目录1+排除目录2+排除目录N
```

- **源目录**：支持多个源目录，多个源目录之间用 `+` 分隔。
- **目标目录**：合并文件系统挂载的目标位置。
- **挂载标志**：指定挂载类型，默认为 `bind` 或 `rbind`。
- **文件系统类型**：必须为 `merge`。
- **排除目录**：防止合并特定目录，默认为 `/proc`、`/dev`、`/tmp`、`/run` 等目录。

### 例：使用合并挂载

```bash
ll-killer exec --mount /dir1+/dir2:/mnt/merged:merge
```

**解释**：
- `/dir1+/dir2`: 将两个目录合并在一起，`dir2` 中的文件将覆盖 `dir1` 中的冲突文件。
- `/mnt/merged`: 合并后的文件系统将挂载到这个目标目录。

### 3. 注意事项

- 在进行合并挂载时，文件系统的优先级按从右到左的顺序决定，越靠后的目录优先级越高，后挂载的目录中的冲突文件会覆盖前面挂载的文件。
- 如果某目录只在一个源目录中出现，它将直接绑定到目标目录。如果该目录是只读的，那么挂载后的文件夹也会保持只读属性。
- 合并挂载对于没有内核 `overlayfs` 或 `fuse` 模块支持的系统仍然有效，适用于跨平台的使用场景。

通过 `ll-killer` 的挂载功能，你可以灵活地管理容器环境中的文件系统，使得容器的构建和运行更加高效和隔离。

## 高级进阶玩法

### 1. 脱离ll-builder/ll-box启动一个容器
在最新的ll-killer中，可以使用`ll-killer exec`子命令直接启动容器环境，从而绕过玲珑容器的一系列限制，同时支持root和rootless权限。

此功能与传统的容器工具（如`runc`、`podman`等）有异曲同工之妙，但`ll-killer`不依赖`OCI`标准配置，直接通过命令行传递参数，避免了shell处理JSON配置文件难的问题，也避免了显式配置`CAP`，允许从当前环境继承各类权限。


#### 1.1 具有Root的玲珑开发环境
此容器中的root具有和宿主机相同的权限，主要用于解决文件系统权限相关的问题，在基于内核的idmap挂载实现之前，不失为一种解决办法。
```bash
# 填写玲珑base的存储地址，在这个目录中：~/.cache/linglong-builder/layers/main/中
# 格式： {BASE名称}/{BASE版本}/{架构}/{binary运行时版本/develop开发版本}
# export KILLER_DEBUG=1 # ll-killer调试模式
BASE_DIR="/home/uos/.cache/linglong-builder/layers/main/org.deepin.base/23.1.0.2/arm64/binary/files/"
CWD=$(pwd)
FS="$CWD/linglong/filesystem"
ROOTFS=$FS/merged
mkdir -p $FS/{merged,diff,work}
ll-killer exec \
--mount "overlay:$FS/merged::overlay:lowerdir=$BASE_DIR,upperdir=$FS/diff,workdir=$FS/work" \
--mount "/proc:$ROOTFS/proc:rbind" \
--mount "/dev:$ROOTFS/dev:rbind" \
--mount "/tmp:$ROOTFS/tmp:rbind" \
--mount "/sys:$ROOTFS/sys:rbind" \
--mount "/etc/resolv.conf:$ROOTFS/etc/resolv.conf:rbind" \
--rootfs $ROOTFS \
--root \
-- ${@:-sh}

```
使用root权限执行上述命令即可获得一个具有真实root权限的玲珑容器运行时环境，你可以在其中安装所需的软件，文件系统差异由overlay文件系统自动捕获到`$FS/diff`中。  

若已在当前项目使用ll-killer创建并初始化玲珑项目，可以直接使用`ll-killer commit`命令打包本应用。

**注意事项**：
* 此模式下会产生root所属的文件。


#### 1.2 Rootless的玲珑开发环境
`Rootless`环境中的用户身份仍然是root，但该root不是真实的root，它的权限与运行命令时用户的权限相同。  
使用非root用户执行之前的shell脚本即可进入rootless玲珑开发环境，可以在环境中正常运行任何不影响文件权限的命令，包括apt安装等，只是一旦应用使用的chown等命令，则会遇到权限问题。


**注意事项**：
* 由于尚未实现内核`idmap`挂载，存在文件权限问题，这个问题请看下一节。

#### 1.3 基于Ptrace+Rootless的玲珑开发环境
在上一个章节中已经讲述了如何启动`rootless`开发环境，但是存在一些权限问题，本章节提供了一种权限问题的缓解办法。

`ll-killer`提供了一个`ptrace`子命令，该命令在`amd64`、`arm64`和`loong64(未测试)`架构上可用。
该命令可以拦截子命令的`chown`系统调用，使其重定向到当前用户，可以避免出现权限问题。

**示例：**
在`rootless`容器中执行此命令，可规避`apt install`时出现的chown错误。
```bash
ll-killer ptrace -- apt install <deb>

# 即便当前UID不是1234，也没有ROOT权限，此命令也不会错误。
ll-killer ptrace -- chown 1234:1234 test

# 直接进入一个屏蔽chown错误的命令行环境。
ll-killer ptrace -- bash
```

**注意事项**：
 * 此模式将极大降低程序性能，请仅在开发环境使用。

## 疑难解答

**build环境内使用apt安装出现chown权限问题**  
 可以添加`--ptrace`参数解决。

**build环境内某些目录写入出现`Permission denied`**  
 是因为系统仓库内置的`fuse-overlayfs`有bug所致，需要重新编译最新版本，参考[获取静态fuse-overlayfs](#获取静态fuse-overlayfs)章节。

**build环境内某些文件写入出现`Operation not permitted`，但可读取**  
 因为内核没有启用EVM，需要启用内核EVM，或使用`make ENABLE_NO_EVM=yes`重新编译`ll-killer-go`。

**使用root模式构建后，ll-killer commit/ll-builder build出现错误**  
删除linglong文件夹，或`chown -R $(id -u):$(id -g) linglong`将文件所有者改为自己。

**apt安装某些包时总是提示依赖不满足**  
请确保`sources.list`中列出的源与玲珑`base`兼容。

## 贡献与维护

欢迎对 `ll-killer-go` 提交问题报告、特性请求和贡献代码。请遵循项目的贡献指南来提交你的代码或文档改进。

## 许可

`ll-killer-go` 项目采用 [MIT License](LICENSE) 进行开源。