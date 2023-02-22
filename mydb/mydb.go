package mydb

import (
	"fmt"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

var Db *sql.DB

func Init_db() {
	var err error
	Db, err = sql.Open("mysql", "root:123456@/douyin")
	if err != nil{
		panic(err)
	}

	//defer Db.Close()
	if err:=Db.Ping();err!=nil{
		fmt.Println("连接失败")
		panic(err)
	}
	fmt.Println("连接成功")
}
