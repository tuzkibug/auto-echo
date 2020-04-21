package base

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql" // 引入包，不使用，使其调用init函数注册mysql
)

func ConnectToMysql(newpassword string) {
	db, err := sql.Open("mysql", "root@tcp("+"192.168.56.109"+":3306)/mysql?charset=utf8mb4")
	if err != nil {
		fmt.Println("创建数据库对象失败")
		return
	}
	defer db.Close() // 延迟关闭 db对象创建成功后才可以调用close方法

	// 实际去尝试连接数据库
	for {
		err = db.Ping()
		if err != nil {
			fmt.Println("连接数据库失败")
			return
		} else {
			break
		}
	}

	fmt.Println("连接数据库成功")

	//	pass1 := "12345"
	//	fmt.Println("请设置你的mysql数据库root密码")
	//	fmt.Scanln(&pass1)
	sqlStr := "alter user 'root'@'%' identified by '" + newpassword + "'"
	// 非常重要：确保QueryRow之后调用Scan方法，否则持有的数据库链接不会被释放
	err = db.QueryRow(sqlStr).Scan()
	if err != nil {
		fmt.Printf("failed, err:%v\n", err)
		return
	}
	fmt.Println("root password is changed ! New password is root!")
}
