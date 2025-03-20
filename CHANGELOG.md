
<a name="v1.4.26"></a>
## [v1.4.26](https://github.com/System233/ll-killer-go/compare/v1.4.25...v1.4.26) (2025-03-20)

### 新增功能

* 增加xpm图标的修复和测试功能

### 错误修复

* 确保入口可执行文件具有执行和读权限


<a name="v1.4.25"></a>
## [v1.4.25](https://github.com/System233/ll-killer-go/compare/v1.4.24...v1.4.25) (2025-03-20)

### 代码重构

* 规范更新功能的下载地址并显示更新内容

### 错误修复

* 修复init进程阻塞，以及fuse递归爆栈


<a name="v1.4.24"></a>
## [v1.4.24](https://github.com/System233/ll-killer-go/compare/v1.4.23...v1.4.24) (2025-03-20)

### 新增功能

* make中添加rebuild、ldd-check、ldd-search和check任务，详情查看make help

### 错误修复

* 调整fast模式下ldd-check的搜索范围


<a name="v1.4.23"></a>
## [v1.4.23](https://github.com/System233/ll-killer-go/compare/v1.4.21...v1.4.23) (2025-03-20)

### 新增功能

* 新增update子命令，支持通过update命令更新ll-killer
* 将自身作为init进程阻塞容器并等待所有后台进程退出 。 本功能用于解决某些应用拉起daemon进程后退出，导致容器销毁的问题。 为此，exec和layer exec子命令添加了wait-timeout参数，可以控制等待的最大时长，若指定了此参数，达到超时时间将杀死所有进程。

### 构建系统

* 哈希校验升级至sha256，校验文件更名为SHA256SUMS，增加SHA256SUMS.asc签名文件


<a name="v1.4.21"></a>
## [v1.4.21](https://github.com/System233/ll-killer-go/compare/v1.4.20...v1.4.21) (2025-03-18)

### 兼容性更改

* clean-fs时同时删除filesystem和output

### 新增功能

* 允许通过FILTER_LIST变量设置make test的排除列表

### 错误修复

* 避免make help等只读操作生成config.mk


<a name="v1.4.20"></a>
## [v1.4.20](https://github.com/System233/ll-killer-go/compare/v1.4.19...v1.4.20) (2025-03-17)

### 错误修复

* 修复make clean-all命令
* ostree当id为空时，不要搜索依赖


<a name="v1.4.19"></a>
## [v1.4.19](https://github.com/System233/ll-killer-go/compare/v1.4.18...v1.4.19) (2025-03-17)

### 兼容性更改

* 自动化测试时，仅避免测试失败时退出

### 新增功能

* 允许make PKG直接传入本地deb文件，需要./开头


<a name="v1.4.18"></a>
## [v1.4.18](https://github.com/System233/ll-killer-go/compare/v1.4.17...v1.4.18) (2025-03-17)

### 兼容性更改

* 自动化测试失败时不再提前终止，每次运行完整的测试
* 减少Makefile中的clean命令

### 错误修复

* 处理/usr/local/share中的快捷方式


<a name="v1.4.17"></a>
## [v1.4.17](https://github.com/System233/ll-killer-go/compare/v1.4.15...v1.4.17) (2025-03-17)

### 新增功能

* 从主机复制所有sources.list
* 添加make build用于进入当前配置下的构建环境
* build支持指定rootfs/runtime; 新增layer exec子命令; Makefile新增ENABLE_OSTREE选项，彻底移除ll-builder依赖
* 增加基于Makefile的ostree管理器[实验性功能]
* layer build添加runtime挂载支持
* 添加仅参数打印功能，用于CLI获取相关信息

### 构建系统

* 修复构建系统日志生成

### 错误修复

* 准备文件系统阶段过滤特殊文件(socket/fifo/block)
* 改进ostree仓库管理；Makefile添加ostree/ptrace功能开关
* 修正Makefile的ostree功能
* 退出前卸载目录以避免孤儿进程
* 修复layer build的根子目录挂载
* layer build自定义rootfs


<a name="v1.4.15"></a>
## [v1.4.15](https://github.com/System233/ll-killer-go/compare/v1.4.14...v1.4.15) (2025-03-15)

### 错误修复

* 去除测试中的env.sh依赖


<a name="v1.4.14"></a>
## [v1.4.14](https://github.com/System233/ll-killer-go/compare/v1.4.13...v1.4.14) (2025-03-15)

### 代码调整

* 优化测试输出信息

### 兼容性更改

* 取消将ll-killer复制到项目位置

### 性能改进

* 增加快速依赖检查模式，详情查看make help

### 新增功能

* Makefile添加ENABLE_INSTALL选项，可用于关闭默认的apt安装功能
* Makefile增加可选build.sh/post-build.sh步骤

### 构建系统

* 优化Makefile构建依赖

### 错误修复

* 修复并优化自动化测试
* 测试时排除隐藏的启动项 feat: 测试服务单元和右键菜单
* 消除低版本scrot不支持of选项的提示


<a name="v1.4.13"></a>
## [v1.4.13](https://github.com/System233/ll-killer-go/compare/v1.4.12...v1.4.13) (2025-03-15)

### 错误修复

* 允许测试搜索系统默认图标


<a name="v1.4.12"></a>
## [v1.4.12](https://github.com/System233/ll-killer-go/compare/v1.4.11...v1.4.12) (2025-03-15)

### 错误修复

* 修正安装缺失库时传递的build选项


<a name="v1.4.11"></a>
## [v1.4.11](https://github.com/System233/ll-killer-go/compare/v1.4.10...v1.4.11) (2025-03-15)

### 代码调整

* 为init/create子命令添加短标志
* 版本信息中增加架构信息

### 兼容性更改

* 默认启用NEVM补丁

### 新增功能

* 新增自动化测试功能, 使用make test即可体验

### 错误修复

* 修复非严格模式下的build环境启动
* 切换到fuse-overlayfs主分支以消除lazytime提示


<a name="v1.4.10"></a>
## [v1.4.10](https://github.com/System233/ll-killer-go/compare/v1.4.9...v1.4.10) (2025-03-14)

### 构建系统

* 编译nevm版本时在版本号中添加nevm标记
* 添加EVM补丁，使用ENABLE_NO_EVM=yes启用 当fuse读写有问题时可用此补丁
* 确保补丁已应用


<a name="v1.4.9"></a>
## [v1.4.9](https://github.com/System233/ll-killer-go/compare/v1.4.8...v1.4.9) (2025-03-13)

### 错误修复

* 更正Makefile中安装的缺失库文件列表


<a name="v1.4.8"></a>
## [v1.4.8](https://github.com/System233/ll-killer-go/compare/v1.4.7...v1.4.8) (2025-03-13)

### 错误修复

* 修复某些叠加目录的权限不正确


<a name="v1.4.7"></a>
## [v1.4.7](https://github.com/System233/ll-killer-go/compare/v1.4.6...v1.4.7) (2025-03-13)

### 错误修复

* 修复某些apt提示找不到status文件


<a name="v1.4.6"></a>
## [v1.4.6](https://github.com/System233/ll-killer-go/compare/v1.4.5...v1.4.6) (2025-03-13)

### 代码重构

* 切换到 gopkg.in/yaml.v2

### 功能改进

* 完善基于Makefile的玲珑项目管理

### 错误修复

* 在apt环境中屏蔽来自主机的包缓存


<a name="v1.4.5"></a>
## [v1.4.5](https://github.com/System233/ll-killer-go/compare/v1.4.4...v1.4.5) (2025-03-12)

### 代码重构

* 调整layer子命令帮助文案
* 分离build-aux资源文件

### 兼容性更改

* 将包名从本地迁移到GitHub

### 构建系统

* CI自动测试

### 错误修复

* 修复内置版本号显示 build: 添加单元测试


<a name="v1.4.4"></a>
## [v1.4.4](https://github.com/System233/ll-killer-go/compare/v1.4.3...v1.4.4) (2025-03-12)

### 错误修复

* 防止覆盖resolv.conf/localtime/timezone/machine-id


<a name="v1.4.3"></a>
## [v1.4.3](https://github.com/System233/ll-killer-go/compare/v1.4.2...v1.4.3) (2025-03-11)

### 代码调整

* 调整layer子命令帮助

### 兼容性更改

* 弃用commit/export组合命令，请考虑使用 layer build 子命令来生成layer文件。


<a name="v1.4.2"></a>
## [v1.4.2](https://github.com/System233/ll-killer-go/compare/v1.4.1...v1.4.2) (2025-03-11)

### 错误修复

* killer打包环境变量名更正为KILLER_PACKER


<a name="v1.4.1"></a>
## [v1.4.1](https://github.com/System233/ll-killer-go/compare/v1.4.0...v1.4.1) (2025-03-11)

### 新增功能

* 为layer build子命令添加后处理支持

### 构建系统

* 更新make依赖

### 错误修复

* 使用KILLER_PICKER标识killer layer build环境


<a name="v1.4.0"></a>
## [v1.4.0](https://github.com/System233/ll-killer-go/compare/v1.3.1...v1.4.0) (2025-03-11)

### 新增功能

* 添加layer build子命令，无需ll-builder即可构建并输出layer


<a name="v1.3.1"></a>
## [v1.3.1](https://github.com/System233/ll-killer-go/compare/v1.3.0...v1.3.1) (2025-03-10)

### 错误修复

* 修复build-aux子命令创建文件


<a name="v1.3.0"></a>
## [v1.3.0](https://github.com/System233/ll-killer-go/compare/v1.2.1...v1.3.0) (2025-03-10)

### 代码重构

* 调整代码结构

### 兼容性更改

* 禁用run命令别名上的参数解析

### 新增功能

* 入口点添加wait后台进程支持，避免后台进程运行时容器被销毁
* exec子命令添加wait选项等待后台进程全部退出
* exec子命令添加nsenter选项 [因权限问题暂时不可用]
* 添加nsenter子命令
* 新增layer系列子命令
* 增加Dbus/右键菜单补丁支持，添加更多systemd查找位置

### 错误修复

* 处理可能的主线程被替换的情况
* build-aux强制覆盖选项和避免覆盖自身


<a name="v1.2.1"></a>
## [v1.2.1](https://github.com/System233/ll-killer-go/compare/v1.2.0...v1.2.1) (2025-03-09)

### 错误修复

* ptrace处理signaled终止信号


<a name="v1.2.0"></a>
## [v1.2.0](https://github.com/System233/ll-killer-go/compare/v1.1.4...v1.2.0) (2025-03-09)

### 代码调整

* 禁用commit/export上的参数解析，现在无需双横线分割

### 性能改进

* 默认使用内置fuse/ifovl挂载，提升性能

### 新增功能

* 添加内置fuse-overlayfs挂载模式: ifovl，无需再提供外部二进制
* 添加内置overlay命令

### 构建系统

* 移除changelog更新
* 更新构建系统Changelog条件

### 错误修复

* 避免ifovl模式下进程进入后台


<a name="v1.1.4"></a>
## [v1.1.4](https://github.com/System233/ll-killer-go/compare/v1.1.3...v1.1.4) (2025-03-08)

### 代码调整

* 调整命令行文本

### 新增功能

* 添加systemd服务单元支持

### 构建系统

* 添加changelog生成


<a name="v1.1.3"></a>
## [v1.1.3](https://github.com/System233/ll-killer-go/compare/v1.1.2...v1.1.3) (2025-03-08)

### 错误修复

* 处理进程信号
* create命令仅当指定from时读取元数据


<a name="v1.1.2"></a>
## [v1.1.2](https://github.com/System233/ll-killer-go/compare/v1.1.1...v1.1.2) (2025-03-07)

### 错误修复

* 正确识别远程返回值


<a name="v1.1.1"></a>
## [v1.1.1](https://github.com/System233/ll-killer-go/compare/v1.1.0...v1.1.1) (2025-03-06)

### 构建系统

* 构建结果添加sha1校验


<a name="v1.1.0"></a>
## [v1.1.0](https://github.com/System233/ll-killer-go/compare/v1.0.13...v1.1.0) (2025-03-06)

### 错误修复

* 处理build退出代码


<a name="v1.0.13"></a>
## [v1.0.13](https://github.com/System233/ll-killer-go/compare/v1.0.12...v1.0.13) (2025-03-05)

### 错误修复

* 重复进入shell


<a name="v1.0.12"></a>
## [v1.0.12](https://github.com/System233/ll-killer-go/compare/v1.0.11...v1.0.12) (2025-03-05)

### 错误修复

* fuse挂载模式绑定目录
* 分离fuse参数


<a name="v1.0.11"></a>
## [v1.0.11](https://github.com/System233/ll-killer-go/compare/v1.0.10...v1.0.11) (2025-03-05)

### 错误修复

* auth.conf.d绑定


<a name="v1.0.10"></a>
## [v1.0.10](https://github.com/System233/ll-killer-go/compare/v1.0.8...v1.0.10) (2025-03-05)

### 代码调整

* 改进错误信息
* 改进退出信息
* 调整pty日志
* 调整帮助信息

### 代码重构

* 创建项目时复制二进制
* 改进错误输出

### 兼容性更改

* 绑定主机dev，防止合并proc/run/sys/tmp/home/root/opt

### 新增功能

* 添加Makefile支持
* 允许无shebang的shell启动命令
* 添加强制覆盖选项
* 挂载pts/shm/mqueue设备
* 切换到pts终端
* 允许通过.killer-debug启用debug
* 在项目目录创建ll-killer副本
* exec添加no-fail标志
* 合并share目录

### 构建系统

* 调整构建依赖

### 错误修复

* 更正amd64系统调用寄存器
* 创建新文件时截断文件
* 合并opt目录
* 修复自动退出时机
* 修复自身查找路径
* 移除create调试信息
* 适配1.7.x
* 版本号规范化
* fchownat/lchown系统调用
* fchownat/lchown系统调用
* 添加SIGTERM信号和版本号
* ptrace命令解析终止符


<a name="v1.0.8"></a>
## [v1.0.8](https://github.com/System233/ll-killer-go/compare/v1.0.7...v1.0.8) (2025-03-01)

### 错误修复

* 添加auth.conf.d和sys挂载


<a name="v1.0.7"></a>
## [v1.0.7](https://github.com/System233/ll-killer-go/compare/v1.0.6...v1.0.7) (2025-03-01)

### 错误修复

* 创建应用时版本号至少保留一位0


<a name="v1.0.6"></a>
## [v1.0.6](https://github.com/System233/ll-killer-go/compare/v1.0.5...v1.0.6) (2025-03-01)

### 错误修复

* 入口点名称
* build时初始化apt目录


<a name="v1.0.5"></a>
## [v1.0.5](https://github.com/System233/ll-killer-go/compare/v1.0.4...v1.0.5) (2025-03-01)

### 新增功能

* 脱离ll-builder/ll-box

### 错误修复

* 当FSType时不执行绑定挂载


<a name="v1.0.4"></a>
## [v1.0.4](https://github.com/System233/ll-killer-go/compare/v1.0.3...v1.0.4) (2025-02-28)

### 新增功能

* 添加arm64和loong64的ptrace支持


<a name="v1.0.3"></a>
## [v1.0.3](https://github.com/System233/ll-killer-go/compare/v1.0.2...v1.0.3) (2025-02-28)

### 错误修复

* 移除不支持环境的ptrace参数


<a name="v1.0.2"></a>
## [v1.0.2](https://github.com/System233/ll-killer-go/compare/v1.0.1...v1.0.2) (2025-02-28)

### 错误修复

* 重新添加create命令


<a name="v1.0.1"></a>
## v1.0.1 (2025-02-28)

