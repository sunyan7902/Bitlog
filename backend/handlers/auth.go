package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"bitlog/db"
	"bitlog/middleware"
	"bitlog/models"
	"bitlog/utils"
)

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UpdateCredentialsRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Login 登录控制器
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户名和密码不能为空"})
		return
	}

	var user models.User
	result := db.DB.Where("username = ?", req.Username).First(&user)
	if result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	// 密码哈希校验
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	// 签发 JWT
	token, err := utils.SignToken(user.ID, user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器内部错误"})
		return
	}

	// 登录成功，清除该 IP 的登录防爆破限制
	clientIP := middleware.GetClientIP(c)
	middleware.ClearRateLimit(clientIP)

	// 写入会话 Cookie
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("bitlog_session", token, 7*24*60*60, "/", "", false, true) // 7天有效, Secure: false, HttpOnly: true

	c.JSON(http.StatusOK, gin.H{"message": "登录成功"})
}

// Logout 登出控制器
func Logout(c *gin.Context) {
	// 删除会话 Cookie (通过设置 maxAge = -1)
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("bitlog_session", "", -1, "/", "", false, true)

	// 重定向至主页，兼容原接口逻辑
	c.Redirect(http.StatusFound, "/")
}

// UpdateCredentials 更改凭据
func UpdateCredentials(c *gin.Context) {
	// 从 JWT 中文解析出的 userId 
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req UpdateCredentialsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求格式"})
		return
	}

	if req.Username == "" && req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "没有提供需要更新的数据"})
		return
	}

	updates := make(map[string]interface{})
	if req.Username != "" {
		updates["username"] = req.Username
	}
	if req.Password != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器内部错误"})
			return
		}
		updates["password"] = string(hashed)
	}

	var user models.User
	result := db.DB.Model(&user).Where("id = ?", userID).Updates(updates)
	if result.Error != nil {
		// 校验唯一性索引报错 P2002/UniqueIndex
		if result.Error.Error() == "UNIQUE constraint failed: User.username" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "用户名已存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器内部错误"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "凭据更新成功"})
}

// Me 获取当前登录用户会话信息
func Me(c *gin.Context) {
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"username": username})
}
