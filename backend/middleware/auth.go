package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"bitlog/utils"
)

// JWTAuth 拦截未登录或过期请求
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从 Cookie 中获取会话令牌
		token, err := c.Cookie("bitlog_session")
		if err != nil || token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// 校验令牌合法性
		claims, err := utils.VerifyToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// 将授权的用户信息暂存入上下文
		c.Set("userId", claims.UserID)
		c.Set("username", claims.Username)

		c.Next()
	}
}
