
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

