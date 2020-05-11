#!/bin/bash
# Program:
#first release
#    1. User input CM server hostname or CM agent hostname, program will use it to config:
#    2. Obtain IP address and netmask, calculate the net IP address
#    3. Config HOSTNAME in local and network
#    4. Add IP Mapping in /etc/hosts
#    5. Modify NTP configuration, including NTP Server and crontab task in Client
#    6. Change CM server hostname in CM Agent configuration
#    7. Start NTP service, MySQL service, CM Server service, CM Agent service in CM Server
#    8. Start CM Agent service in CM Agent
#2nd release change
#    1. Server IP can't obtained by ping server hostname, so change to manual input.
#    2. Add restart system proceed between finish configuration change and start service. 
#    3. Remove systemctl start "${1}" in startservice

#History:
#2020/04/28  Zhujianjun First release
#2020/05/08  2nd release, update proceed flow.

PATH=/bin:/sbin:/usr/bin:/usr/sbin:/usr/local/bin:/usr/local/sbin:~/bin
export PATH

#获取本机IP地址
ipaddr=$(ifconfig -a|grep inet|grep -v 127.0.0.1|grep -v inet6|awk '{print $2}'|tr -d "addr:")
netmask=$(ifconfig -a|grep inet|grep -v 127.0.0.1|grep -v inet6|awk '{print $4}')
#获取网络号
if [ "${ipaddr}" == "" ];then
	echo "No ip address obtined, please check the network..."
	exit 1
else
	netip='' #网络号
	for index in {1..4}; do
	    si=$(echo $ipaddr | cut -d "." -f $index)
	    sm=$(echo $netmask | cut -d "." -f $index)
	    if [ $index -ne 1 ]
	    then
	        netip="$netip."
	    fi
	    netip="$netip$[$si&$sm]"
	done
fi

echo -e "\nLocal IP Address is normal\nIP:${ipaddr}\nNETMASK:${netmask}\nNET SEGMENT:${netip}\n"

function hostnamecheck(){
	#主机名是否有效确认，只由小写字母及数字构成且第一位为小写字母
	[ "${1}" == "" ] && echo "Hostname is blank, please input the right hostname" && exit 1
	namecheck1=$(echo "${1}" |grep '[^[:lower:]^[:digit:]]')
	[ "${namecheck1}" != "" ] && echo "Hostname only can contain lower character and number, please input again..."  && exit 1
	namecheck2=$(echo "${1}" |grep '^[[:lower:]*]')
	[ "${namecheck2}" == "" ] && echo "Hostname only can start with lower character, please input again..."  && exit 1
	echo "Hostname check Passed."
}

function is_valid_ip_format(){  
		#输入的IP格式是否有效
    if [[ "$1" =~ ^[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}$ ]] ;then  
        echo "CM Server IP valid..."
    else  
      echo "CM Server IP format invalid, please check CM server IP address again.." && exit 1  
    fi  
}

function backupfile(){
	#备份配置文件
	echo -e "Start backup configuration..."
	$(cp /etc/hosts /etc/hosts.bak)	
  $(cp /etc/hostname /etc/hostname.bak)			
  $(cp /etc/sysconfig/network /etc/sysconfig/network.bak)	
  $(cp /etc/cloudera-scm-agent/config.ini /etc/cloudera-scm-agent/config.ini.bak)	
	if [ "${1}" == "server" ];then
		$(cp /etc/ntp.conf /etc/ntp.conf.bak)	
	else 
	  $(cp /var/spool/cron/root /var/spool/cron/root.bak)	
	echo -e "Backup configuration file Finished"
	fi
}

function execfile(){
	#执行恢复配置文件过程
	echo -e "Start restore backup configuration file:${1}"
	cp "${1}" "${1}".new
	cp "${1}".bak "${1}"
	rm -rf "${1}".bak
	echo -e "Restore backup configuration file:${1} Finished"
}

function restorefile(){
	#备份配置文件
	echo -e "\nStart restore configuration..."
	$(test -e /etc/hosts.bak) && execfile /etc/hosts || echo "/etc/hosts Backup file not exist, Keep the configuration"
	$(test -e /etc/hostname.bak) && execfile /etc/hostname || echo "/etc/hostname Backup file not exist, Keep the configuration"
	$(test -e /etc/sysconfig/network.bak) && execfile /etc/sysconfig/network || echo "/etc/sysconfig/network Backup file not exist, Keep the configuration"
	$(test -e /etc/cloudera-scm-agent/config.ini.bak) && execfile /etc/cloudera-scm-agent/config.ini || echo "/etc/cloudera-scm-agent/config Backup file not exist, Keep the configuration"
	$(test -e /etc/ntp.conf.bak) && execfile /etc/ntp.conf || echo "/etc/ntp.conf Backup file not exist, Keep the configuration"
	$(test -e /var/spool/cron/root.bak) && execfile /var/spool/cron/root || echo "/var/spool/cron/root Backup file not exist, Keep the configuration"
	echo -e "Restore configuration Finished"
}

function startservice(){
	#启动服务
	servicecheck=$(systemctl status "${1}")
	if [ servicecheck == "Unit ${1}.service could not be found." ];then
		echo -e "Service:${1} not be found, Can't Start.."
	else
		systemctl enable "${1}"
		#systemctl start "${1}"   #service will auto start after linux system restart
		systemctl status "${1}"
	fi
}

function stopservice(){
	#关闭服务
	echo -e "Start to stop service:${1}"
	checkservice=$(systemctl status "${1}")
	if [ checkservice == "Unit ${1}.service could not be found." ];then
		echo -e "Service:${1} not be found, no need to stop.."
	else
		systemctl disable "${1}"
		systemctl stop "${1}"
		systemctl status "${1}"
		echo -e "\nService:${1} Stopped"
	fi
}

function restartsystem(){
	  echo -e "The Linux system will restart now to make configuration valid... "
		sync
	  sync
	  sync
	  shutdown -r now
}

read -p "Please make choice: 1. Install CM; 2. Remove CM\nplease input 1 or 2:" option

if [ "${option}" == "1" ];then
	echo -e "Please define the Hostname,\\e[1;31mMUST NOT\e[0m use the same HOSTNAME with others. \nProgram will use the hostname to finish configuration. \n"
	read -p "Local host is the CM Server or not? (y/n):" flag
	if [ "${flag}" == "Y" -o "${flag}" == "y" ];then
	  read -p "Please input the Server hostname(Only small character and Number): " ServerHostName
	  #判断主机名是否只包含小写字母和数字
	  hostnamecheck ${ServerHostName}
	  
	  read -p "ServerHostName is ${ServerHostName}, Whether continue? (y/n)" doublecheck
	  [ "${doublecheck}" != "y" -a  "${doublecheck}" != "Y" ] && echo "Configuration Abort, Bye..." && exit 0
	  
	  #备份文件
	  backupfile server
	  
	  #添加主机名映射及使主机名配置生效
	  echo -e "\n00.Change HOSTNAME in /etc/hostname"
	  sed -i '1c '${ServerHostName}'' /etc/hostname

	  echo -e "01.Add HOSTNAME in /etc/sysconfig/network"
	  echo "HOSTNAME=${ServerHostName}">> /etc/sysconfig/network
	  source /etc/sysconfig/network
	  
	  #添加ip地址映射
	  #echo -e "02.Add IP mapping in /etc/hosts"
	  #echo "${ipaddr} ${ServerHostName}" >>  /etc/hosts
	    
	  #追加ntp配置信息
	  echo -e "03.Add NTP setting in /etc/ntp.conf"
	  sed -i 's/#restrict 192.168.1.0 mask 255.255.255.0/restrict '${netip}' mask '${netmask}'/g' /etc/ntp.conf
	  
	  #启动NTP服务
	  echo -e "04.Set NTP service autostart when reboot"
	  startservice ntpd

	  #启动MySQL
	  echo -e "05.Set MySQL service autostart when reboot"
	  startservice mysqld

	  #修改CDH配置文件
	  echo -e "06.Change server hostname in CM Agent setting"
	  sed -i 's/server_host=localhost/server_host='${ServerHostName}'/g' /etc/cloudera-scm-agent/config.ini
	  
	  #启动CM Server
	  echo -e "07.Set CM Server service autostart when reboot"
	  startservice cloudera-scm-server
	   
	  #启动CM Agent
	  echo -e "08.Set CM Agent service autostart when reboot"
	  startservice cloudera-scm-agent
	  echo -e "Finish Configuration in below files:\n /etc/hostname \n /etc/sysconfig/network \n /etc/hosts \n /etc/ntp.conf \n /etc/cloudera-scm-agent/config.ini\n"
	  
	  #重启系统
	  restartsystem
	    
	elif [ "${flag}" == "N" -o "${flag}" == "n" ];then
		read -p "Please input the Server hostname(Only small character and Number): " ServerHostName
	  #判断主机名是否只包含小写字母和数字
	  hostnamecheck ${ServerHostName}
		
		read -p "Please define the local hostname(Only small character and Number): " Agenthostname
		#判断主机名是否只包含小写字母和数字
		hostnamecheck ${Agenthostname}
		
		#获取服务端IP地址
		read -p "please input CM Server IP address: " ServerIP
		is_valid_ip_format ${ServerIP}
		
		read -p "ServerHostName is ${ServerHostName}, ServerIP is ${ServerIP}, Agenthostname is ${Agenthostname}, Whether continue? (y/n)" doublecheck
	  [ "${doublecheck}" != "y" -a  "${doublecheck}" != "Y" ] && echo "Configuration Abort, Bye..." && exit 0
		
		#备份文件
	  backupfile agent
	  	
	  #添加主机名映射及使主机名配置生效
	  echo -e "\n00.Change HOSTNAME in /etc/hostname"
	  sed -i '1c '${Agenthostname}'' /etc/hostname
	  
	  echo -e "01.Add HOSTNAME in /etc/sysconfig/network"
	  echo "HOSTNAME=${Agenthostname}">> /etc/sysconfig/network
	  source /etc/sysconfig/network
	  
	  #添加ip地址映射
	  #echo -e "02.Add IP mapping in /etc/hosts"
	  #echo "${ServerIP} ${ServerHostName}" >>  /etc/hosts
	  #echo "${ipaddr} ${Agenthostname}" >>  /etc/hosts
	  	  
	  #添加5分钟定时同步任务
	  echo -e "03.Add NTP setting in crontab task.."
	  echo "*/5 * * * * /usr/sbin/ntpdate ${ServerHostName}  >/dev/null 2>&1;/sbin/hwclock -w" >> /var/spool/cron/root
	  
	  #修改CDH配置文件
	  echo -e "04.Change server hostname in CM Agent setting"
	  sed -i 's/server_host=localhost/server_host='${ServerHostName}'/g' /etc/cloudera-scm-agent/config.ini
	  
	  #启动CM Agent
	  echo -e "05.Set CM Agent service autostart when reboot"
	  startservice cloudera-scm-agent
	  
	  echo -e "Finish Configuration in below files:\n /etc/hostname \n /etc/sysconfig/network \n /etc/hosts \n /var/spool/cron/root \n /etc/cloudera-scm-agent/config.ini\n" 
	  
	  #重启系统
	  restartsystem
	  
	else
	  echo "Please input Y or N, exit...."
	  exit 1
	fi
elif [ "${option}" == "2" ];then
	stopservice ntpd
	stopservice mysqld
	stopservice cloudera-scm-server
	stopservice cloudera-scm-agent
  restorefile
else
	echo "Please input 1 or 2, exit...."
	exit 1 
fi