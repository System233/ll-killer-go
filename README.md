# ll-killer-go - 玲珑杀手Go <!-- omit from toc -->

[![Build](https://github.com/System233/ll-killer-go/actions/workflows/build.yaml/badge.svg)](https://github.com/System233/ll-killer-go/actions/workflows/build.yaml) 

## 项目简介 <!-- omit from toc -->

本项目是[Linglong Killer Self-Service (ll-killer 玲珑杀手)](https://github.com/System233/linglong-killer-self-service)的重写版本，去除了构建阶段的shell脚本，全部采用go实现，并添加了一些增强功能。

`ll-killer-go` 是一款专为解决玲珑容器应用构建问题而设计的命令行工具。它帮助开发者快速创建、构建和生成玲珑容器应用项目，同时提供一整套辅助构建与调试功能。使用该工具，用户可以在玲珑容器中获得与传统容器（如 `docker`）类似的构建体验，包括对特权命令的支持和 `apt` 安装等功能，免去了手动解压软件包（如 `deb`）和修复依赖库查找路径的繁琐操作。`ll-killer-go` 通过重建自由的容器环境，确保构建过程的一致性与可靠性，极大地提高了开发效率和可维护性。

### ⚠️文档会随着更新而落后，最新功能请查看 [更新日志](CHANGELOG.md) 或命令帮助


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
  - [2. 基于Makefile的自动打包](#2-基于makefile的自动打包)
    - [创建打包环境](#创建打包环境)
    - [配置软件源](#配置软件源)
    - [全自动构建GIMP](#全自动构建gimp)
      - [打包本地deb](#打包本地deb)
      - [调整构建内容](#调整构建内容)
    - [应用测试](#应用测试)
  - [3. 手动打包](#3-手动打包)
    - [创建打包环境](#创建打包环境-1)
    - [配置软件源](#配置软件源-1)
    - [创建 linglong.yaml 项目配置](#创建-linglongyaml-项目配置)
    - [手动构建GIMP应用](#手动构建gimp应用)
    - [测试构建结果](#测试构建结果)
    - [安装构建结果](#安装构建结果)
    - [清理构建缓存（可选）](#清理构建缓存可选)
  - [3. Makefile配置说明](#3-makefile配置说明)
    - [配置变量](#配置变量)
    - [自定义流程文件](#自定义流程文件)
    - [其他特殊文件](#其他特殊文件)
- [疑难解答](#疑难解答)
    - [build环境内使用apt安装出现chown权限问题](#build环境内使用apt安装出现chown权限问题)
    - [apt安装某些包时总是提示依赖不满足](#apt安装某些包时总是提示依赖不满足)
    - [出现`fork/exec`操作不允许等问题](#出现forkexec操作不允许等问题)
    - [玲珑1.7.11以上直接使用`ll-killer build`提示挂载文件失败](#玲珑1711以上直接使用ll-killer-build提示挂载文件失败)
    - [某些应用在安装时出现 invalid mode 0104755 with bits 04000](#某些应用在安装时出现-invalid-mode-0104755-with-bits-04000)
    - [应用运行出现 cannot execute: required file not found](#应用运行出现-cannot-execute-required-file-not-found)
- [贡献与维护](#贡献与维护)
- [许可](#许可)


## 快速教程

教程中的示例应用是GIMP，本章节提供两个版本的打包教程，一个是基于 **Makefile** 的全流程自动化打包，另一个是直接使用 **ll-killer** 或shell脚本手动管理打包流程。

推荐使用 **Makefile** 自动管理项目。

---

### 1. 获取与配置 `ll-killer`

首先，确保系统已安装必要的依赖：

```sh
sudo apt install apt-file erofs-utils
```
所有可选依赖及其功能：
* **apt-file**:用于查找和安装缺失的库
* **ostree**: 用于取代`ll-builder`的base下载功能。
* **erofs-utils**: `layer build`子命令所需依赖，用于取代`ll-builder`的打包功能，以加速打包，减少磁盘复制和占用。
* **erofsfuse**: layer挂载相关功能需要，用于取代自动化测试时的`ll-cli`依赖。
* **xvfb xdotool scrot**: 运行自动化测试所需依赖。


然后，下载 `ll-killer` 并赋予执行权限：

```sh
wget https://github.com/System233/ll-killer-go/releases/latest/download/ll-killer-amd64 -O ll-killer
chmod +x ll-killer
```

如果上面的地址无法下载，可尝试更换这个地址重试：

```sh
wget https://ll-killer.win/releases/latest/download/ll-killer-amd64 -O ll-killer
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

接下来的步骤默认 `ll-killer` 已加入 `PATH`。

---

### 2. 基于Makefile的自动打包

**ll-killer** 提供了一个 **Makefile**，可以自动计算参数并调用 **ll-killer** 的相关功能，实现项目构建的自动管理。

#### 创建打包环境
```sh
# 创建工作目录并进入
mkdir gimp && cd gimp

# 初始化ll-killer项目脚本，如果脚本行为不符合你的预期，一切都可以修改。
ll-killer init
```

#### 配置软件源

创建 `sources.list` 并填充软件源，以下为 Deepin V23 的APT软件源，注意添加`[trusted=yes]`选项以忽略仓库签名：

```sh
cat >sources.list <<EOF
 deb [trusted=yes] https://community-packages.deepin.com/deepin/beige beige main commercial community
 deb [trusted=yes] https://com-store-packages.uniontech.com/appstorev23 beige appstore
 deb [trusted=yes] https://community-packages.deepin.com/driver-23/ driver non-free
EOF
```
#### 全自动构建GIMP

```sh
make config PKG=gimp
make
```
整个过程自动执行ll-killer指令，最后输出layer文件。
>[!NOTE]
>项目目录下会生成一个`config.mk`文件，该文件配置了make运行时用到的参数，你可以使用`make config KEY=VAL`或者直接修改来编辑配置。
>Makefile中提供了丰富的配置指令，可以运行`make help`查看，或查看[3. Makefile配置说明](#3-makefile配置说明)。

##### 打包本地deb
* 如果你需要打包本地deb文件，将参数中的`gimp`替换为`./本地deb路径`即可。

##### 调整构建内容

使用以下命令进入容器shell，`make build`
```sh
make build
```

#### 应用测试

Makefile还封装了自动化测试功能，需要安装`xvfb xdotool scrot`三个必要依赖。
```sh
make test
```
运行完成后，在tests文件夹中查看应用输出日志以及屏幕截图。

---

### 3. 手动打包

以下步骤使用shell命令演示如何直接使用ll-killer命令进行打包。

#### 创建打包环境

```sh
# 创建工作目录并进入
mkdir gimp && cd gimp

# 初始化ll-killer项目脚本，如果脚本行为不符合你的预期，一切都可以修改。
ll-killer init
```

#### 配置软件源

创建 `sources.list` 并填充软件源，以下为 Deepin V23 的APT软件源，注意添加`[trusted=yes]`选项以忽略仓库签名：

```sh
cat >sources.list <<EOF
 deb [trusted=yes] https://community-packages.deepin.com/deepin/beige beige main commercial community
 deb [trusted=yes] https://com-store-packages.uniontech.com/appstorev23 beige appstore
 deb [trusted=yes] https://community-packages.deepin.com/driver-23/ driver non-free
EOF
```

可以使用记事本创建此文件。

#### 创建 linglong.yaml 项目配置

```sh
echo "[获取包元数据]"
ll-killer apt -- apt show gimp > pkg.info
echo "[创建玲珑项目]"
ll-killer create --from pkg.info --base org.deepin.base/23.1.0
```

#### 手动构建GIMP应用

```sh
echo "[将包安装至构建环境]"
ll-killer build --ptrace -- apt install -y gimp

echo "[检查缺失库]"
#缺失库输出到了ldd-check.log文件
ll-killer build -- build-aux/ldd-check.sh > ldd-check.log 

echo "[搜索缺失库所在包]"
#找到的库所在deb包名输出到ldd-found.log，找不到的库输出到ldd-notfound.log 
ll-killer apt -- build-aux/ldd-search.sh ldd-check.log ldd-found.log ldd-notfound.log 

echo "[安装找到的缺失库]"
ll-killer build -- apt install -y $(cat ldd-found.log)

echo "[提交文件到玲珑容器中]"
# 玲珑1.7版本以下需要去掉--skip-output-check --skip-strip-symbols参数
ll-killer commit -- --skip-output-check --skip-strip-symbols
echo "[导出layer文件]"
ll-killer export -- --layer
```
构建完成后，最终会生成两个 `.layer` 文件，仅需 `*_binary.layer`。  

>[!NOTE]
>提示：`ll-killer commit`和`ll-killer export`是`ll-builder build`和`ll-builder export`的别名命令。
>你可以使用新的`ll-killer layer build`命令取代它们。


#### 测试构建结果
此命令进入容器环境的shell，请结合`/usr/share/applications/`路径内desktop文件中的Exec指令来测试程序。
ll-killer统一将应用快捷方式放置在`/usr/share/applications/`中， 你可以通过命令`grep Exec /usr/share/applications/*.desktop`快速查看desktop快捷方式中的启动指令。

```sh
ll-builder run
```
>[!NOTE]
> run命令默认执行`/opt/apps/APPID/files/entrypoint.sh`入口点。如需进入原始玲珑容器环境，请运行：
>```sh
>ll-builder run --exec bash
>```
#### 安装构建结果

```sh
# 玲珑 >= 1.7.x 版本需要 sudo 权限
ll-cli install *_binary.layer
```

#### 清理构建缓存（可选）

如果需要清理构建过程中产生的缓存，可执行：

```sh
# 清除构建日志和layer输出
rm *.log *.layer
# 删除linglong文件夹
rm -rf linglong

# 删除ll-builder缓存，如果不删除，你的磁盘将被ll-builder耗尽。
# 你可以使用ll-killer layer build来取代ll-builder，避免相关问题。
ll-builder list|grep ":gimp.linyaps/"|xargs -r ll-builder remove 

```

---

### 3. Makefile配置说明

Makefile中的某些选项需要安装相应的依赖，请仔细查看说明，所有依赖已在 [获取与配置 `ll-killer`](#1-获取与配置-ll-killer) 章节中给出。

#### 配置变量
使用 `KEY=VALUE` 的格式设置变量，具体如下：

| 变量名                | 说明                                                                                            | 默认值              |
| --------------------- | ----------------------------------------------------------------------------------------------- | ------------------- |
| **PKG**               | 要打包的 Debian 软件包名称                                                                      | `app`               |
| **PKGID**             | 从`PKG`输入中解析的deb包名                                                                      |
| **APPID**             | 玲珑应用 ID                                                                                     | `PKGID.linyaps`     |
| **ENABLE_LDD_CHECK**  | 是否启用依赖检查 (`0` 关闭, `1` 启用)                                                           | `1`                 |
| **ENABLE_PTRACE**     | 是否启用 `ptrace` (`0` 关闭, `1` 启用)                                                          | `1`                 |
| **ENABLE_INSTALL**    | 是否启用自动安装依赖，关闭以完全使用自定义构建脚本 (`0` 关闭, `1` 启用)                         | `1`                 |
| **ENABLE_OSTREE**     | 是否启用自定义 `ostree` 仓库，彻底移除 `ll-builder` 依赖，需要安装`ostree` (`0` 关闭, `1` 启用) | `0`                 |
| **ENABLE_RM_DESKTOP** | 是否自动删除不属于主包的快捷方式 (`0` 关闭, `1` 启用)                                           | `$(ENABLE_INSTALL)` |
| **ENABLE_TEST_NOCLI** | 是否启用无 `ll-cli` 测试，需先启用 `ENABLE_OSTREE`，需要安装`erofsfuse` (`0` 关闭, `1` 启用)    | `0`                 |
| **LDD_CHECK_MODE**    | 依赖检查模式，快速模式下仅检查应用目录的依赖情况 (`fast` = 快速, `full` = 全量)                 | `fast`              |
| **CREATE_ARGS**       | 传递给 `ll-killer create` 的额外参数，用于调整自动生成的linglong.yaml                           | *(空)*              |
| **BUILD_ARGS**        | 传递给 `ll-killer build` 的额外参数，用于调整构建环境                                           | *(空)*              |
| **FILTER_LIST**       | 自动化测试-启动项测试排除项目，空格分隔                                                         | NoDisplay Hidden    |

可选值（适用于 `FILTER_LIST`）：`NoDisplay`、`Hidden`、`Terminal`

#### 自定义流程文件

你可以在项目目录下创建这些文件来在相应阶段执行自定义操作。

| 名称              | 说明                                                                          |
| ----------------- | ----------------------------------------------------------------------------- |
| **deps.list**     | 额外依赖包名，一行一个，其中填写的包名将在安装阶段一起安装                    |
| **build.sh**      | 自定义构建脚本，在安装阶段之前执行，注意文件有无执行权限以及`#!/bin/bash`行   |
| **post-build.sh** | 自定义后构建脚本，在安装阶段之后执行，注意文件有无执行权限以及`#!/bin/bash`行 |
| **sources.list**  | 安装依赖时使用的 APT 源配置                                                   |

#### 其他特殊文件

| 名称               | 说明                                             |
| ------------------ | ------------------------------------------------ |
| **apt.conf.d**     | APT 配置文件夹                                   |
| **auth.conf.d**    | APT 授权配置文件夹（遇到 APT 源返回 401 时使用） |
| **build-aux**      | 辅助构建脚本目录                                 |
| **sources.list.d** | 使用的 APT 源配置文件夹                          |
| **linglong.yaml**  | 项目配置文件                                     |
| **ll-killer**      | `ll-killer` 可执行文件【可选】                   |
| **Makefile**       | 基于 `Makefile` 的玲珑项目构建规则               |
| **\*.log**         | 各阶段构建日志                                   |
| **\*.layer**       | 生成的 layer 文件                                |

---


## 疑难解答

#### build环境内使用apt安装出现chown权限问题
 可以添加`--ptrace`参数解决。

#### apt安装某些包时总是提示依赖不满足 
请确保`sources.list`中列出的源与玲珑`base`兼容，或更换`aptitude`命令来安装。

#### 出现`fork/exec`操作不允许等问题

如果系统没有安装过玲珑，需要手动启用非特权命名空间功能
```sh
sudo sysctl -w kernel.unprivileged_userns_clone=1
sudo sysctl -w user.max_user_namespaces=28633
sudo sysctl -w kernel.apparmor_restrict_unprivileged_userns=0
sudo sysctl -w kernel.apparmor_restrict_unprivileged_unconfined=0
```

<details>
<summary>
更多信息
</summary>

**以下内容适用于专业用户**

如果问题仍旧出现，需要具体分析，可以设置环境变量`KILLER_DEBUG=1`来启用ll-killer的调试信息输出。
`entrypoint.sh`相关问题可以通过删除脚本中的`--rootfs`参数来在shell中观察`/run/app.rootf`内的文件是否正常。

</details>

#### 玲珑1.7.11以上直接使用`ll-killer build`提示挂载文件失败

该版本玲珑修改了`/tmp`文件系统，**ll-killer** 依赖 **tmpfs** 与主机通信。

**解决办法1:** 升级ll-killer版本到1.5.4或更高。

**解决办法2:** 使用`Makefile`管理项目流程，并启用 `ENABLE_OSTREE=1` 功能，直接移除 `ll-builder` 依赖。

**ll-builder** 在 **ll-killer** 环境中本身只是个base下载器，启用`ENABLE_OSTREE`选项，用 `build-aux/ostree.mk` 取代它。

#### 某些应用在安装时出现 invalid mode 0104755 with bits 04000
应用内某些文件设置了SUID/SGID所致，需要在构建时确保此类权限已经全部删除，可用使用`chmod a-s -R $PREFIX`命令递归的删除所有SUID/SGID权限。

* **ll-killer v1.5.4** 起自动在准备文件系统阶段删除此类权限。


#### 应用运行出现 cannot execute: required file not found

这是某些目录被错误覆盖所致，ll-killer中，`$PREFIX`目录下的所有文件将叠加至容器根目录，请确保叠加方式正确，或`$PREFIX`下没有多余的文件。

比如：玲珑版本1.7.x下ll-killer构建的应用可能出现bash: /entrypoint.sh: cannot execute: required file not found，原因是`ll-builer build`命令默认`--skip-strip-symbols`启用符号剔除，而该功能会在`$PREFIX`下创建lib目录放置某些调试文件，此lib文件夹会在ll-killer在准备根文件系统时覆盖掉/lib符号链接，造成所有应用无法启动。

**解决办法:** 使用`ll-killer layer build`或在使用`ll-builder build`时添加`--skip-strip-symbols`选项来禁用符号剔除。

**此问题的下一步计划:** 1.调整叠加目录的位置或结构，避免被lib之类的文件夹意外覆盖根目录；2. 重新设计叠加目录结构，增加无overlayfs和无mergefs的启动模式支持（无ll-killer模式）。

## 贡献与维护

欢迎对 `ll-killer-go` 提交问题报告、特性请求和贡献代码。请遵循项目的贡献指南来提交你的代码或文档改进。

## 许可

`ll-killer-go` 项目采用 [MIT License](LICENSE) 进行开源。