package models

import (
	"time"

	"github.com/lucsky/cuid"
	"gorm.io/gorm"
)

// User 映射原 Prisma User 表
type User struct {
	ID        string    `gorm:"column:id;primaryKey" json:"id"`
	Username  string    `gorm:"column:username;uniqueIndex" json:"username"`
	Password  string    `gorm:"column:password" json:"-"`
	CreatedAt time.Time `gorm:"column:createdAt;default:CURRENT_TIMESTAMP" json:"createdAt"`
	Posts     []Post    `gorm:"foreignKey:AuthorID" json:"posts,omitempty"`
}

// TableName 指定原 User 表名 (区分大小写)
func (User) TableName() string {
	return "User"
}

// BeforeCreate 钩子：插入前自动生成 CUID
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == "" {
		u.ID = cuid.New()
	}
	return nil
}

// Post 映射原 Prisma Post 表
type Post struct {
	ID        string    `gorm:"column:id;primaryKey" json:"id"`
	Slug      string    `gorm:"column:slug;uniqueIndex" json:"slug"`
	Title     string    `gorm:"column:title" json:"title"`
	Content   string    `gorm:"column:content" json:"content"`
	Published bool      `gorm:"column:published;default:true" json:"published"`
	CreatedAt time.Time `gorm:"column:createdAt;default:CURRENT_TIMESTAMP" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updatedAt" json:"updatedAt"`
	AuthorID  *string   `gorm:"column:authorId" json:"authorId,omitempty"`
	Author    *User     `gorm:"foreignKey:AuthorID" json:"author,omitempty"`
}

// TableName 指定原 Post 表名
func (Post) TableName() string {
	return "Post"
}

// BeforeCreate 钩子：自动生成 CUID 并填充 UpdatedAt
func (p *Post) BeforeCreate(tx *gorm.DB) (err error) {
	if p.ID == "" {
		p.ID = cuid.New()
	}
	p.UpdatedAt = time.Now()
	return nil
}

// BeforeUpdate 钩子：更新时刷新 UpdatedAt 时间
func (p *Post) BeforeUpdate(tx *gorm.DB) (err error) {
	p.UpdatedAt = time.Now()
	return nil
}

// SiteConfig 映射原 Prisma SiteConfig 表
type SiteConfig struct {
	Key   string `gorm:"column:key;primaryKey" json:"key"`
	Value string `gorm:"column:value" json:"value"`
}

// TableName 指定原 SiteConfig 表名
func (SiteConfig) TableName() string {
	return "SiteConfig"
}
