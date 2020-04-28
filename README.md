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
创建虚拟机(返回虚拟机ID)-->获取mysql信息(返回IP，mac)-->绑定浮动IP(返回浮动IP)-->SSH远程连接虚拟机-->执行安装脚本安装mysql-->初始化mysql密码(返回密码)-->传递文件和脚本-->执行数据库优化和主备配置(暂无)-->交付用户

其中绑定浮动IP没有API，需转发http请求进行资源创建，分为以下几步：
1.获取token(返回token字符串)
2.通过mac获取port-id(返回port id) /v2.0/ports?mac_address=xxxxxx&fields=id
3.创建floating ip并绑定port(返回floating ip)

# my.cnf.master & my.cnf.slave
用于mysql数据库主备配置的文件
