# CHANGELOG

## v1.6.6

- Support sqlite3

## v1.6.5

- fix page stock
- memorize page# when refresh

## v1.6.4

- Dark mode for frondend
- Localhost bypass

## v1.6.3

- From the repo `https://github.com/ghostong/gocron`
  - 新增`定时任务`搜索栏`命令`搜索条件
  - 新增`定时任务->查看日志`页面`手动执行`按钮
  - 新增`任务依赖`页面,以树状图形式展示`子任务`
  - 新增`任务日志`中`1天`,`7天`,`1个月`之前的日志按钮
  - 修改`定时任务`页任务`标签`搜索改为选择框
  - 修改`定时任务`输入框记录搜索条件
  - 修改`任务日志->清空日志` 由 `DELETE` 改为 `TRUNCATE`

- Upgrade to Vue 3
- Typo fix: creater => creator

## v1.6

- 优化整体**界面样式与布局**，包括界面色系，列表，详情，按钮组，分页等
- 调整权限等级，增加**超级管理员**，可以管理所有任务；**管理员**调整为管理自己的任务和查看其他任务和日志，普通用户与原有权限一致，仅可查看所有任务和日志
- 任务详情页增加快捷选择crontab按钮组
- 任务详情页支持更改任务状态
- 任务列表支持**标签**,**命令**搜索

## v1.5

- 前端使用Vue+ElementUI重构
- 任务通知
  - 新增WebHook通知
  - 自定义通知模板
  - 匹配任务执行结果关键字发送通知
- 任务列表页显示任务下次执行时间

## v1.4

- HTTP任务支持POST请求
- 后台手动停止运行中的shell任务
- 任务执行失败重试间隔时间支持用户自定义
- 修复API接口调用报403错误

## v1.3

- 支持多用户登录
- 增加用户权限控制

## v1.2.2

- 用户登录页增加图形验证码
- 支持从旧版本升级
- 任务批量开启、关闭、删除
- 调度器与任务节点支持HTTPS双向认证
- 修复任务列表页总记录数显示错误

## v1.1

- 任务可同时在多个节点上运行
- *nix平台默认禁止以root用户运行任务节点
- 子任务命令中增加预定义占位符, 子任务可根据主任务运行结果执行相应操作
- 删除守护进程模块
- Web访问日志输出到终端
