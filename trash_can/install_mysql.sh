#!/bin/bash

`rm -rf /usr/local/mysql* || (echo "rm mysql failed";exit)` \
&& `rm -rf /etc/my.cnf || (echo "rm my.cnf failed";exit)` \
&& `cp /root/my.cnf /etc/my.cnf || (echo "cp failed";exit)` \
&& `groupadd mysql || (echo "groupadd failed";exit)` \
&& `useradd -g mysql mysql || (echo "useradd failed";exit)` \
&& `mkdir /usr/local/mysql/data || (echo "mkdir data failed";exit)` \
&& `/usr/local/mysql/bin/mysqld --initialize-insecure --user=mysql --basedir=/usr/local/mysql/ --datadir=/usr/local/mysql/data/ || (echo "initial failed";exit)` \
&& `cp /usr/local/mysql/support-files/mysql.server /etc/init.d/mysqld || (echo "cp service failed";exit)` \
&& `ln -s /usr/local/mysql/bin/mysql /usr/bin/mysql || (echo "ln failed";exit)` \
&& `service mysqld start >/dev/null || (echo "service start failed";exit)` \
&& `mysql -u root -e "update mysql.user set host='%' where user='root';" || (echo "remote mysql failed";exit)` \
&& `service mysqld restart >/dev/null || (echo "service start failed";exit)` \
&& `systemctl enable mysqld`
