# 特别告知
## auto-openstack整体已更新为echo框架，并迁移至本仓库，请知，谢谢！

# 说明
## base
组成服务的基本功能，包括各项openstack操作（openstack.go），mysql初始化（mysql.go），ssh远程操作（ssh.go）。

## controllers
调用base的服务，构造echo架构的控制器，匹配请求的路由。目前有创建，mysql密码初始化，IP获取，SSH操作等各项控制器。

## main.go
主文件，构造web服务器，定义路由及对应的控制器，完成自动化功能。

# 为什么选择echo
高性能，编码极简，高扩展性，极轻量级
