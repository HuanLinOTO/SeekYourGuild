package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"seekyourguild/internal/config"
	"seekyourguild/internal/database"
	"seekyourguild/internal/handlers"
	"seekyourguild/internal/middleware"
	"seekyourguild/internal/static"

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

	// 静态文件服务（前端）
	if static.HasStaticFiles() {
		log.Println("Serving embedded static files")
		
		// 托管静态资源文件
		r.GET("/assets/*filepath", func(c *gin.Context) {
			filePath := strings.TrimPrefix(c.Request.URL.Path, "/")
			data, contentType, err := static.ReadFile(filePath)
			if err != nil {
				c.Status(http.StatusNotFound)
				return
			}
			c.Data(http.StatusOK, contentType, data)
		})
		
		// 根路径直接返回 index.html
		r.GET("/", func(c *gin.Context) {
			data, contentType, err := static.ReadFile("index.html")
			if err != nil {
				c.String(http.StatusInternalServerError, "index.html not found")
				return
			}
			c.Data(http.StatusOK, contentType, data)
		})
		
		// 处理其他路由 - SPA fallback
		r.NoRoute(func(c *gin.Context) {
			// API 路由返回 404
			if strings.HasPrefix(c.Request.URL.Path, "/api") {
				c.JSON(http.StatusNotFound, gin.H{
					"code":    404,
					"message": "API not found",
				})
				return
			}
			
			// 尝试提供静态文件
			path := strings.TrimPrefix(c.Request.URL.Path, "/")
			
			// 检查文件是否存在
			if data, contentType, err := static.ReadFile(path); err == nil {
				c.Data(http.StatusOK, contentType, data)
				return
			}
			
			// SPA fallback - 返回 index.html
			data, contentType, _ := static.ReadFile("index.html")
			c.Data(http.StatusOK, contentType, data)
		})
	} else {
		log.Println("No embedded static files found, running in API-only mode")
	}

	// 启动服务器
	addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	log.Printf("Server starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
