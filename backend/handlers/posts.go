package handlers

import (
	"archive/zip"
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"math"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"bitlog/db"
	"bitlog/models"
	"bitlog/utils"
)

type CreatePostRequest struct {
	Title         string `json:"title" binding:"required"`
	Slug          string `json:"slug" binding:"required"`
	Content       string `json:"content"`
	ContentBase64 string `json:"contentBase64"`
}

type UpdatePostRequest struct {
	Title         *string `json:"title"`
	Slug          *string `json:"slug"`
	Content       *string `json:"content"`
	ContentBase64 *string `json:"contentBase64"`
}

// GetPosts 获取文章列表 (支持分页与模糊搜索)
func GetPosts(c *gin.Context) {
	query := c.Query("q")
	pageStr := c.Query("page")
	limitStr := c.Query("limit")

	dbQuery := db.DB.Model(&models.Post{}).Order("createdAt desc")

	if query != "" {
		likePattern := "%" + query + "%"
		dbQuery = dbQuery.Where("title LIKE ? OR content LIKE ?", likePattern, likePattern)
	}

	// 如果没有 page 参数，返回完整列表（兼容原 Next.js API）
	if pageStr == "" {
		var posts []models.Post
		if err := dbQuery.Find(&posts).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch posts"})
			return
		}
		c.JSON(http.StatusOK, posts)
		return
	}

	// 执行分页逻辑
	page := 1
	if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
		page = p
	}
	limit := 10
	if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
		limit = l
	}

	var total int64
	dbQuery.Count(&total)

	var posts []models.Post
	offset := (page - 1) * limit
	if err := dbQuery.Offset(offset).Limit(limit).Find(&posts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch posts"})
		return
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	c.JSON(http.StatusOK, gin.H{
		"posts":      posts,
		"total":      total,
		"page":       page,
		"totalPages": totalPages,
	})
}

// GetPostDetail 获取单篇文章详情
func GetPostDetail(c *gin.Context) {
	slug := c.Param("slug")
	var post models.Post
	if err := db.DB.Where("slug = ?", slug).First(&post).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		return
	}
	c.JSON(http.StatusOK, post)
}

// CreatePost 撰写新文章
func CreatePost(c *gin.Context) {
	var req CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "标题和固定链接不能为空"})
		return
	}

	// 解码 Base64 避免传输转义 Bug
	content := req.Content
	if req.ContentBase64 != "" {
		decoded, err := base64.StdEncoding.DecodeString(req.ContentBase64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的 Base64 内容"})
			return
		}
		content = string(decoded)
	}

	// 唯一性校验
	var count int64
	db.DB.Model(&models.Post{}).Where("slug = ?", req.Slug).Count(&count)
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "该固定链接已被使用"})
		return
	}

	post := models.Post{
		Title:     req.Title,
		Slug:      req.Slug,
		Content:   content,
		Published: true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := db.DB.Create(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建博客文章失败"})
		return
	}

	c.JSON(http.StatusOK, post)
}

// UpdatePost 修改博客
func UpdatePost(c *gin.Context) {
	slug := c.Param("slug")
	var post models.Post
	if err := db.DB.Where("slug = ?", slug).First(&post).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		return
	}

	var req UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的参数内容"})
		return
	}

	var content *string
	if req.ContentBase64 != nil {
		decoded, err := base64.StdEncoding.DecodeString(*req.ContentBase64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的 Base64 内容"})
			return
		}
		cStr := string(decoded)
		content = &cStr
	} else if req.Content != nil {
		content = req.Content
	}

	// 如果修改了 URL Slug，校验新 Slug 冲突
	if req.Slug != nil && *req.Slug != "" && *req.Slug != slug {
		var count int64
		db.DB.Model(&models.Post{}).Where("slug = ?", *req.Slug).Count(&count)
		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "该固定链接已被使用"})
			return
		}
		post.Slug = *req.Slug
	}

	if req.Title != nil && *req.Title != "" {
		post.Title = *req.Title
	}
	if content != nil {
		post.Content = *content
	}

	if err := db.DB.Save(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存更改失败"})
		return
	}

	c.JSON(http.StatusOK, post)
}

// DeletePost 删除单篇文章
func DeletePost(c *gin.Context) {
	slug := c.Param("slug")
	var post models.Post
	if err := db.DB.Where("slug = ?", slug).First(&post).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		return
	}

	if err := db.DB.Delete(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

// DeleteAllPosts 清空全部文章 (危险操作)
func DeleteAllPosts(c *gin.Context) {
	tx := db.DB.Begin()
	var posts []models.Post
	if err := tx.Find(&posts).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取文章失败"})
		return
	}

	if len(posts) > 0 {
		if err := tx.Exec("DELETE FROM Post").Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
			return
		}
	}
	tx.Commit()

	c.JSON(http.StatusOK, gin.H{"deleted": len(posts)})
}

// ExportPosts 打包导出所有文章
func ExportPosts(c *gin.Context) {
	var posts []models.Post
	if err := db.DB.Order("createdAt desc").Find(&posts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "导出数据失败"})
		return
	}

	var buf bytes.Buffer
	zipWriter := zip.NewWriter(&buf)

	for _, post := range posts {
		// 生成带元数据的 markdown 文本
		mdText := utils.BuildMarkdown(post.Title, post.Slug, post.CreatedAt, post.Published, post.Content)
		
		// 转换安全文件名
		safeTitle := utils.CleanFilename(post.Title)
		if safeTitle == "" {
			safeTitle = post.Slug
		}
		filename := fmt.Sprintf("%s.md", safeTitle)

		f, err := zipWriter.Create(filename)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "创建 Zip 文件失败"})
			return
		}
		if _, err := f.Write([]byte(mdText)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "写入 Zip 失败"})
			return
		}
	}

	if err := zipWriter.Close(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成 Zip 归档失败"})
		return
	}

	todayStr := time.Now().Format("2006-01-02")
	c.Header("Content-Type", "application/zip")
	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="limblog-export-%s.zip"`, todayStr))
	c.Data(http.StatusOK, "application/zip", buf.Bytes())
}

// ImportPosts 导入文章 (支持 .zip 压缩包或单篇 .md 文本)
func ImportPosts(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请选择导入文件"})
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法读取上传文件"})
		return
	}
	defer file.Close()

	results := gin.H{
		"imported": 0,
		"skipped":  0,
		"errors":   []string{},
	}

	importedCount := 0
	skippedCount := 0
	var errorLogs []string

	filenameLower := strings.ToLower(fileHeader.Filename)

	if strings.HasSuffix(filenameLower, ".zip") {
		// ZIP 导入逻辑
		zipBytes, err := io.ReadAll(file)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "解压文件失败"})
			return
		}

		items, err := utils.ReadZipFile(zipBytes)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "读取 ZIP 格式异常"})
			return
		}

		for _, item := range items {
			// 校验冲突
			var count int64
			db.DB.Model(&models.Post{}).Where("slug = ?", item.Slug).Count(&count)
			if count > 0 {
				skippedCount++
				errorLogs = append(errorLogs, fmt.Sprintf("跳过 \"%s\": slug \"%s\" 已存在", item.Title, item.Slug))
				continue
			}

			// 入库
			post := models.Post{
				Title:     item.Title,
				Slug:      item.Slug,
				Content:   item.Content,
				Published: item.Published,
				CreatedAt: item.CreatedAt,
				UpdatedAt: time.Now(),
			}
			if err := db.DB.Create(&post).Error; err != nil {
				errorLogs = append(errorLogs, fmt.Sprintf("导入 \"%s\" 失败: %v", item.Title, err))
			} else {
				importedCount++
			}
		}
	} else if strings.HasSuffix(filenameLower, ".md") {
		// 单文件 .md 导入
		mdBytes, err := io.ReadAll(file)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "读取 Markdown 失败"})
			return
		}

		fm, body, err := utils.ParseMarkdown(string(mdBytes))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "解析 Markdown YAML Frontmatter 失败"})
			return
		}

		basename := strings.TrimSuffix(filepath.Base(fileHeader.Filename), filepath.Ext(fileHeader.Filename))
		title := fm.Title
		if title == "" {
			title = strings.ReplaceAll(basename, "-", " ")
		}

		slugVal := fm.Slug
		if slugVal == "" {
			slugVal = utils.MakeSafeSlug(basename)
		}

		createdAt := fm.Date
		if createdAt.IsZero() {
			createdAt = time.Now()
		}

		var count int64
		db.DB.Model(&models.Post{}).Where("slug = ?", slugVal).Count(&count)
		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("该固定链接 (slug: %s) 已存在", slugVal)})
			return
		}

		post := models.Post{
			Title:     title,
			Slug:      slugVal,
			Content:   body,
			Published: fm.Published,
			CreatedAt: createdAt,
			UpdatedAt: time.Now(),
		}
		if err := db.DB.Create(&post).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("导入失败: %v", err)})
			return
		}
		importedCount++
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "仅支持 .zip 或 .md 文件格式"})
		return
	}

	results["imported"] = importedCount
	results["skipped"] = skippedCount
	results["errors"] = errorLogs

	c.JSON(http.StatusOK, results)
}
