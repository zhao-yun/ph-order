package main

import (
	"demo/util/postgres"
	"demo/util/redis"
	"demo/util/socket"
)

func main() {

	// 数据库初始化
	postgres.Init()

	// Redis初始化
	redis.Init()

	// Socket初始化
	go socket.Init()

	// 路由配置
	routeInit()

}
