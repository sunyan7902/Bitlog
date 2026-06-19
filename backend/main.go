package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"bitlog/db"
	"bitlog/handlers"
	"bitlog/middleware"
)

func main() {
	// 1. 初始化 SQLite 数据库连接
	db.InitDB()

	// 2. 区分运行模式 (生产模式禁用 gin debug 日志)
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// 3. 全局跨域处理中间件
	r.Use(middleware.CORS())

	// 4. 静态资源托管
	// 托管本地上传的图片/文件
	r.Static("/uploads", "./data/uploads")
	// 托管编译后的前端静态资源
	r.Static("/assets", "./dist/assets")
	r.StaticFile("/favicon.ico", "./dist/favicon.ico")

	// 5. API 路由分组规划
	api := r.Group("/api")
	{
		// --- 身份验证接口 ---
		auth := api.Group("/auth")
		{
			// 登录接口附带 IP 防爆破限流保护
			auth.POST("/login", middleware.IPRateLimiter(), handlers.Login)
			auth.POST("/logout", handlers.Logout)
		}

		// --- 站点全局配置接口 (公开读取) ---
		api.GET("/site-config", handlers.GetSiteConfig)

		// --- 前台公开只读博客接口 (受站点可见性 private 中间件保护) ---
		api.GET("/posts", middleware.CheckVisibility(), handlers.GetPosts)
		api.GET("/posts/:slug", middleware.CheckVisibility(), handlers.GetPostDetail)

		// --- 需要管理员 JWT 授权的高危管理接口 ---
		protected := api.Group("/")
		protected.Use(middleware.JWTAuth())
		{
			// 获取当前登录用户会话信息
			protected.GET("/auth/me", handlers.Me)

			// 修改管理员用户名/密码
			protected.POST("/auth/settings", handlers.UpdateCredentials)

			// 修改站点配置
			protected.PUT("/site-config", handlers.UpdateSiteConfig)

			// 文章写作与修改
			protected.POST("/posts", handlers.CreatePost)
			protected.PUT("/posts/:slug", handlers.UpdatePost)
			protected.DELETE("/posts/:slug", handlers.DeletePost)

			// 图片及附件安全上传
			protected.POST("/upload", handlers.UploadFile)

			// 数据库一键管理 (清除/导入/打包导出)
			protected.DELETE("/posts/delete-all", handlers.DeleteAllPosts)
			protected.GET("/posts/export", handlers.ExportPosts)
			protected.POST("/posts/import", handlers.ImportPosts)
		}
	}

	// 6. 兼容 React Router SPA 前端路由刷新的 Fallback 处理
	r.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path
		// 如果是未定义的 API 接口，返回 JSON 404
		if strings.HasPrefix(path, "/api/") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Not Found"})
			return
		}
		// 检查在 ./dist/ 下是否存在该物理文件 (防止目录遍历漏洞)
		cleanPath := filepath.Clean(path)
		if !strings.HasPrefix(cleanPath, "..") {
			distFilePath := filepath.Join("./dist", cleanPath)
			if fileInfo, err := os.Stat(distFilePath); err == nil && !fileInfo.IsDir() {
				c.File(distFilePath)
				return
			}
		}
		// 其它路径（如前台路由 /blog, /new, /settings）一律返回 index.html
		c.File("./dist/index.html")
	})

	// 7. 端口读取与启动
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000" // 保持原 Docker 映射的 3000 端口
	}

	log.Printf("LimBlog Go Server is running on port %s...", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Server startup failed: %v", err)
	}
}
