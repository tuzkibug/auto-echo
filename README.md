# 告知  
## auto-openstack整体已更新为echo框架，并迁移至本仓库，请知，谢谢！  

# 说明  
## base  
组成服务的基本功能，包括各项openstack操作（openstack.go），mysql操作（mysql.go）以及其他基础操作。  

## controllers  
调用base的服务，构造echo架构的控制器，匹配请求的路由。目前有创建，mysql密码初始化，IP获取，SSH操作等各项控制器。  

## main.go  
主文件，构造web服务器，定义路由及对应的控制器，完成自动化功能。  

## base_master.txt & base_slave.txt    
用于mysql数据库主备节点启动的文件模板，其中需替换root密码，数据库密码和主节点IP  

## master.txt & slave.txt  
替换用户名密码等参数后生成的新的配置文件  

## hosts_base  
用于CDH虚拟机记录主机节点信息的文件模板，需替换IP和主机名  

## hosts  
替换完成的hosts文件，上传到/etc/hosts完成替换  

## trash_can  
废弃的文件  

# 为什么选择echo  
高性能，编码极简，高扩展性，极轻量级  

# 自动化业务逻辑  
## mysql  
由财哥制作的新镜像，更新启动逻辑  
编辑主mysql节点启动文件-->拉起主mysql节点，等待其安装配置完成-->获取其IP-->编辑备mysql节点启动文件-->拉起备用mysql节点，等待其安装配置完成-->获取其IP-->返回节点内部IP-->(可选)配置浮动IP给主/备节点   

## CDH  
由俊哥制作CDH镜像，启动逻辑如下：  
启动server-->启动若干agent-->获取主机名和IP-->编辑hosts文件-->配置浮动IP-->上传hosts文件-->执行安装脚本-->服务启动完成  

## Floating IP
绑定浮动IP没有现成可用的API，所以手动构造新的http请求进行资源创建，分为以下几步：  
1.获取openstack token(返回token字符串)---完成  
2.通过mac获取port-id(返回port id) /v2.0/ports?mac_address=xxxxxx&fields=id ---完成  
3.创建floating ip并绑定port(返回floating ip) ---完成  
