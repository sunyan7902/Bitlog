package utils

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/gosimple/slug"
	"gopkg.in/yaml.v3"
)

// Frontmatter YAML 头部结构
type Frontmatter struct {
	Title     string    `yaml:"title"`
	Slug      string    `yaml:"slug"`
	Date      time.Time `yaml:"date"`
	Published bool      `yaml:"published"`
}

// ParseMarkdown 提取 YAML Frontmatter 和正文
func ParseMarkdown(content string) (Frontmatter, string, error) {
	var fm Frontmatter
	// 兼容 \r\n 换行符
	normalized := strings.ReplaceAll(content, "\r\n", "\n")

	if !strings.HasPrefix(normalized, "---\n") {
		return fm, strings.TrimSpace(normalized), nil
	}

	endIdx := strings.Index(normalized[4:], "\n---")
	if endIdx == -1 {
		return fm, strings.TrimSpace(normalized), nil
	}

	// 计算绝对结束位置
	endIdx += 4
	frontBlock := normalized[4:endIdx]
	body := strings.TrimSpace(normalized[endIdx+4:])

	// 解析 YAML
	err := yaml.Unmarshal([]byte(frontBlock), &fm)
	if err != nil {
		return fm, body, err
	}

	return fm, body, nil
}

// BuildMarkdown 拼接带 Frontmatter 的 Markdown 正文
func BuildMarkdown(title, slug string, createdAt time.Time, published bool, content string) string {
	safeTitle := strings.ReplaceAll(title, "\"", "\\\"")
	fm := fmt.Sprintf("---\ntitle: \"%s\"\nslug: \"%s\"\ndate: \"%s\"\npublished: %t\n---\n\n%s",
		safeTitle, slug, createdAt.Format(time.RFC3339), published, content)
	return fm
}

// CleanFilename 过滤操作系统非法字符防止导出文件名错误
func CleanFilename(name string) string {
	reg := regexp.MustCompile(`[\\/:*?"<>|]`)
	cleaned := reg.ReplaceAllString(name, "_")
	return strings.TrimSpace(cleaned)
}

// MakeSafeSlug 将中文标题转换为拼音/URL友好的 Slug
func MakeSafeSlug(text string) string {
	s := slug.Make(text)
	if s == "" {
		s = fmt.Sprintf("post-%d", time.Now().UnixNano())
	}
	return s
}

// ZipItem 单篇 Markdown 结构
type ZipItem struct {
	Filename  string
	Title     string
	Slug      string
	Content   string
	CreatedAt time.Time
	Published bool
}

// ReadZipFile 读取 Zip 压缩包内的 Markdown
func ReadZipFile(zipData []byte) ([]ZipItem, error) {
	reader, err := zip.NewReader(bytes.NewReader(zipData), int64(len(zipData)))
	if err != nil {
		return nil, err
	}

	var items []ZipItem
	for _, f := range reader.File {
		if f.FileInfo().IsDir() || !strings.HasSuffix(strings.ToLower(f.Name), ".md") {
			continue
		}

		rc, err := f.Open()
		if err != nil {
			return nil, err
		}

		contentBytes, err := io.ReadAll(rc)
		rc.Close()
		if err != nil {
			return nil, err
		}

		fm, body, err := ParseMarkdown(string(contentBytes))
		if err != nil {
			// 解析失败直接跳过
			continue
		}

		basename := strings.TrimSuffix(filepath.Base(f.Name), filepath.Ext(f.Name))
		title := fm.Title
		if title == "" {
			title = strings.ReplaceAll(basename, "-", " ")
		}

		slugVal := fm.Slug
		if slugVal == "" {
			slugVal = MakeSafeSlug(basename)
		}

		createdAt := fm.Date
		if createdAt.IsZero() {
			createdAt = time.Now()
		}

		items = append(items, ZipItem{
			Filename:  f.Name,
			Title:     title,
			Slug:      slugVal,
			Content:   body,
			CreatedAt: createdAt,
			Published: fm.Published,
		})
	}

	return items, nil
}
