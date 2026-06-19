package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"bitlog/db"
	"bitlog/models"
)

// GetSiteConfig 获取站点全部 KV 配置
func GetSiteConfig(c *gin.Context) {
	var configs []models.SiteConfig
	result := db.DB.Find(&configs)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{}) // 兼容原 Node 实现，出错时返回空对象
		return
	}

	resultMap := make(map[string]string)
	for _, conf := range configs {
		resultMap[conf.Key] = conf.Value
	}

	c.JSON(http.StatusOK, resultMap)
}

// UpdateSiteConfig 批量更新/保存站点配置
func UpdateSiteConfig(c *gin.Context) {
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的配置数据格式"})
		return
	}

	tx := db.DB.Begin()
	for k, v := range req {
		valStr := fmt.Sprintf("%v", v)
		config := models.SiteConfig{
			Key:   k,
			Value: valStr,
		}
		// GORM Save 会自动执行 Upsert (有则更新，无则插入)
		if err := tx.Save(&config).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "更新配置失败"})
			return
		}
	}
	tx.Commit()

	c.JSON(http.StatusOK, gin.H{"success": true})
}
