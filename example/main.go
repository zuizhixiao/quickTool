package main

import (
	"TEMPLATE/config"
	"TEMPLATE/initialize"
)

func main() {
	config.SetUp()

	//加载数据库
	initialize.Mysql()

	//加载redis
	initialize.Redis()

	initialize.RunServer()
}
