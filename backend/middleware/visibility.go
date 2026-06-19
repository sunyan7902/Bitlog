package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"bitlog/db"
	"bitlog/models"
	"bitlog/utils"
)

// CheckVisibility 拦截私密模式下的游客访问
func CheckVisibility() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 如果有有效的管理员 Token，则直接允许通过（管理员可见）
		token, err := c.Cookie("bitlog_session")
		if err == nil && token != "" {
			if _, err := utils.VerifyToken(token); err == nil {
				c.Next()
				return
			}
		}

		// 2. 检查站点可见性配置
		var config models.SiteConfig
		result := db.DB.Where("key = ?", "site_visibility").First(&config)
		if result.Error == nil && config.Value == "private" {
			// 如果是私密模式，直接阻断并返回 404 状态以伪装隐身
			c.JSON(http.StatusNotFound, gin.H{"error": "Not Found"})
			c.Abort()
			return
		}

		c.Next()
	}
}
