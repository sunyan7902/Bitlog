package utils

import (
	"bytes"
	"errors"
	"image"
	_ "image/gif" // 注册 gif 解码器
	"image/jpeg"
	_ "image/png" // 注册 png 解码器
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"

	"golang.org/x/image/draw"
	_ "golang.org/x/image/webp" // 注册 webp 解码器
)

// Magic Bytes 映射
var magicTable = map[string][]byte{
	"image/jpeg": {0xFF, 0xD8, 0xFF},
	"image/png":  {0x89, 0x50, 0x4E, 0x47},
	"image/gif":  {0x47, 0x49, 0x46, 0x38},
}

// ValidateImageSecure 校验图片格式（MIME + Magic Bytes 幻数）
func ValidateImageSecure(fileHeader *multipart.FileHeader, file io.Reader) (string, error) {
	// 1. 基础 MIME 类型校验
	mimeType := fileHeader.Header.Get("Content-Type")
	allowedTypes := []string{"image/jpeg", "image/png", "image/gif", "image/webp"}
	isAllowed := false
	for _, t := range allowedTypes {
		if t == mimeType {
			isAllowed = true
			break
		}
	}
	if !isAllowed {
		return "", errors.New("Invalid file type. Allowed types: JPEG, PNG, GIF, WebP")
	}

	// 2. 读取前几字节进行 Magic Bytes 校验
	headBytes := make([]byte, 12)
	n, err := file.Read(headBytes)
	if err != nil && err != io.EOF {
		return "", err
	}

	// 将指针重置以便后续读取
	if seeker, ok := file.(io.ReadSeeker); ok {
		seeker.Seek(0, io.SeekStart)
	}

	// WebP 特殊魔数校验 (RIFF .... WEBP)
	if mimeType == "image/webp" {
		if n >= 12 && string(headBytes[0:4]) == "RIFF" && string(headBytes[8:12]) == "WEBP" {
			return mimeType, nil
		}
		return "", errors.New("File content does not match declared WebP type")
	}

	// 其他图片魔数校验
	expectedMagic, exists := magicTable[mimeType]
	if exists {
		if n < len(expectedMagic) || !bytes.HasPrefix(headBytes, expectedMagic) {
			return "", errors.New("File content does not match declared MIME type")
		}
	}

	return mimeType, nil
}

// OptimizeAndResizeImage 缩放图片并转换为 JPEG（80% 质量），GIF 和失败情况回退为原图
func OptimizeAndResizeImage(file io.Reader, mimeType string) ([]byte, string, error) {
	// 读取所有字节
	originalData, err := io.ReadAll(file)
	if err != nil {
		return nil, "", err
	}

	// 针对 GIF，我们不做任何压缩和修改，以保留其动态帧
	if mimeType == "image/gif" {
		return originalData, ".gif", nil
	}

	// 解码图片
	img, formatName, err := image.Decode(bytes.NewReader(originalData))
	if err != nil {
		// 解码失败时（例如未知格式），回退到原图以防内容丢失
		return originalData, getExtensionFromMime(mimeType), nil
	}

	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// 设定最大宽度为 1200 像素
	maxWidth := 1200
	if width <= maxWidth {
		// 原图宽度小于等于 1200，不做缩放，仅进行 80% 质量的 JPEG 格式压缩
		var outBuf bytes.Buffer
		err = jpeg.Encode(&outBuf, img, &jpeg.Options{Quality: 80})
		if err != nil {
			return originalData, getExtensionFromMime(mimeType), nil
		}
		return outBuf.Bytes(), ".jpg", nil
	}

	// 等比例计算高度
	newWidth := maxWidth
	newHeight := int(float64(height) * (float64(newWidth) / float64(width)))

	// 创建目标画布
	rect := image.Rect(0, 0, newWidth, newHeight)
	resizedImg := image.NewRGBA(rect)

	// 使用更平滑的 CatmullRom 插值算法进行画质缩放
	draw.CatmullRom.Scale(resizedImg, rect, img, bounds, draw.Over, nil)

	// 编码为 JPEG 80% 质量
	var outBuf bytes.Buffer
	err = jpeg.Encode(&outBuf, resizedImg, &jpeg.Options{Quality: 80})
	if err != nil {
		// 缩放失败回退原图
		return originalData, getExtensionFromMime(mimeType), nil
	}

	_ = formatName // 忽略未使用的格式名称
	return outBuf.Bytes(), ".jpg", nil
}

func getExtensionFromMime(mimeType string) string {
	switch mimeType {
	case "image/jpeg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/gif":
		return ".gif"
	case "image/webp":
		return ".webp"
	default:
		return ".jpg"
	}
}

// GetSafeFilename 生成安全的文件名
func GetSafeFilename(originalName string) string {
	ext := filepath.Ext(originalName)
	base := strings.TrimSuffix(originalName, ext)

	// 移除非英文字母和数字的字符，替换为 "-"
	safeBase := ""
	for _, r := range base {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			safeBase += string(r)
		} else {
			safeBase += "-"
		}
	}
	// 压缩连续的 "-"
	for strings.Contains(safeBase, "--") {
		safeBase = strings.ReplaceAll(safeBase, "--", "-")
	}
	safeBase = strings.Trim(safeBase, "-")
	if safeBase == "" {
		safeBase = "upload"
	}
	return safeBase
}
