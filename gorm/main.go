package main

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var dsn = "root:root@tcp(localhost:3306)/example?charset=utf8&parseTime=True&loc=Local&readTimeout=1s&timeout=1s&writeTimeout=3s"

func main() {
	db, err := gorm.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	db.LogMode(true)
	if err := db.DB().Ping(); err != nil {
		panic(err)
	}
	fmt.Println("success")
}
