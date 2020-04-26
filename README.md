# 特别告知
## auto-openstack整体已更新为echo框架，并迁移至本仓库，请知，谢谢！

# 说明
## base
组成服务的基本功能，包括各项openstack操作（openstack.go），mysql初始化（mysql.go），ssh远程操作（ssh.go）。

## controllers
调用base的服务，构造echo架构的控制器，匹配请求的路由。目前有创建，mysql密码初始化，IP获取，SSH操作等各项控制器。

## main.go
主文件，构造web服务器，定义路由及对应的控制器，完成自动化功能。

## mysql_config
mysql部署和主备配置时可能用到的脚本或配置

# 为什么选择echo
高性能，编码极简，高扩展性，极轻量级

# 自动化业务逻辑
创建虚拟机-->SSH远程连接虚拟机-->安装mysql-->获取mysql信息（IP）-->初始化mysql密码-->数据库优化和主备配置(暂无)-->交付用户

# my.cnf.master & my.cnf.slave
用于mysql数据库主备配置的文件
