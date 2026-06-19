package db

import (
	"log"
	"os"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"bitlog/models"
)

// DB 全局数据库句柄
var DB *gorm.DB

// InitDB 初始化数据库连接
func InitDB() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "data/bitlog.db" // 默认路径
	}

	// 兼容 prisma 的 file: 协议前缀
	dbPath := strings.TrimPrefix(dbURL, "file:")

	// 自动创建数据库文件所在的父级目录
	if lastSlash := strings.LastIndex(dbPath, "/"); lastSlash != -1 {
		dir := dbPath[:lastSlash]
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Fatalf("Failed to create database parent directory: %v", err)
		}
	} else if lastBackslash := strings.LastIndex(dbPath, "\\"); lastBackslash != -1 {
		dir := dbPath[:lastBackslash]
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Fatalf("Failed to create database parent directory: %v", err)
		}
	}

	var err error
	// 使用 busy_timeout 优化 SQLite 并发读写锁冲突
	DB, err = gorm.Open(sqlite.Open(dbPath+"?_pragma=busy_timeout(5000)"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("Failed to connect SQLite database: %v", err)
	}

	log.Printf("Successfully connected to SQLite database at: %s", dbPath)

	// 自动迁移字段与建表
	err = DB.AutoMigrate(&models.User{}, &models.Post{}, &models.SiteConfig{})
	if err != nil {
		log.Fatalf("Database auto-migration failed: %v", err)
	}

	// 初始化默认管理员数据
	seedAdminUser()
}

func seedAdminUser() {
	var count int64
	DB.Model(&models.User{}).Count(&count)
	if count == 0 {
		username := os.Getenv("ADMIN_USERNAME")
		if username == "" {
			username = "admin"
		}
		password := os.Getenv("ADMIN_PASSWORD")
		if password == "" {
			password = "123456"
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
		if err != nil {
			log.Fatalf("Failed to hash default password: %v", err)
		}

		admin := models.User{
			Username:  username,
			Password:  string(hashedPassword),
			CreatedAt: time.Now(),
		}

		if err := DB.Create(&admin).Error; err != nil {
			log.Fatalf("Failed to create default admin user: %v", err)
		}
		log.Printf("Default admin user created. Username: %s, Password: %s", username, password)
	}
}
