package main

import (
	"log"

	"zhuhai_travel_backend/config"
	"zhuhai_travel_backend/database"
	"zhuhai_travel_backend/routes"
)

func main() {
	cfg := config.Load()

	// 初始化数据库
	database.Init(cfg)

	// 设置路由
	r := routes.SetupRouter()

	log.Printf("服务启动 -> http://0.0.0.0:%s", cfg.ServerPort)
	if err := r.Run("0.0.0.0:" + cfg.ServerPort); err != nil {
		log.Fatalf("服务启动失败: %v", err)
	}
}
