package main

import (
	"fmt"
	"log"

	"seekyourguild/internal/config"
	"seekyourguild/internal/database"
	"seekyourguild/internal/handlers"
	"seekyourguild/internal/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置
	cfg := config.Load()

	// 初始化数据库
	if err := database.Init(&cfg.Database); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// 创建 Gin 引擎
	r := gin.Default()

	// CORS 配置
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// API 路由
	api := r.Group("/api")
	{
		// 公开接口 - 获取群列表
		api.GET("/groups", handlers.GetGroups)
		
		// 上报接口 - 需要认证
		report := api.Group("/report")
		report.Use(middleware.AuthMiddleware(cfg.API.Key))
		{
			report.POST("/member", handlers.ReportMemberChange)
			report.POST("/message", handlers.ReportMessage)
			report.POST("/heartbeat", handlers.ReportHeartbeat)
		}

		// 统计接口 - 公开访问
		stats := api.Group("/stats")
		{
			stats.GET("/overview", handlers.GetOverview)
			stats.GET("/member-changes", handlers.GetMemberChanges)
			stats.GET("/messages", handlers.GetMessageStats)
			stats.GET("/active-users", handlers.GetActiveUsers)
		}
	}

	// 启动服务器
	addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	log.Printf("Server starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
