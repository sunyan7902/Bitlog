package handlers

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"bitlog/utils"
)

// UploadFile 处理单图上传
func UploadFile(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "未选择上传文件"})
		return
	}

	// 1. 文件大小校验 (5MB 限额)
	const maxSize = 5 * 1024 * 1024
	if fileHeader.Size > maxSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文件过大，单张图片不能超过 5MB"})
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法读取上传文件"})
		return
	}
	defer file.Close()

	// 2. 真实 Magic Bytes 安全校验
	mimeType, err := utils.ValidateImageSecure(fileHeader, file)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 3. 图像压缩与尺寸优化
	finalData, ext, err := utils.OptimizeAndResizeImage(file, mimeType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "图片优化失败"})
		return
	}

	// 4. 生成安全文件名 (基础名 + 时间戳 + 高强度随机数)
	safeBase := utils.GetSafeFilename(fileHeader.Filename)
	randSuffix, _ := rand.Int(rand.Reader, big.NewInt(1e9))
	uniqueFilename := fmt.Sprintf("%s-%d-%d%s", safeBase, time.Now().UnixMilli(), randSuffix.Int64(), ext)

	// 5. 保存至宿主机持久化目录
	uploadDir := "./data/uploads"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建上传文件夹失败"})
		return
	}

	destPath := filepath.Join(uploadDir, uniqueFilename)
	if err := os.WriteFile(destPath, finalData, 0644); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存文件失败"})
		return
	}

	// 返回图片访问的静态路由地址
	fileURL := fmt.Sprintf("/uploads/%s", uniqueFilename)
	c.JSON(http.StatusOK, gin.H{"url": fileURL})
}
