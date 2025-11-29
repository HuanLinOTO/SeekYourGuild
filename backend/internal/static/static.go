package static

import (
	"embed"
	"io/fs"
	"mime"
	"net/http"
	"path/filepath"
)

//go:embed dist/*
var staticFiles embed.FS

// GetFileSystem 返回静态文件系统
func GetFileSystem() http.FileSystem {
	fsys, err := fs.Sub(staticFiles, "dist")
	if err != nil {
		panic(err)
	}
	return http.FS(fsys)
}

// HasStaticFiles 检查是否有嵌入的静态文件
func HasStaticFiles() bool {
	entries, err := staticFiles.ReadDir("dist")
	if err != nil {
		return false
	}
	return len(entries) > 0
}

// ReadFile 读取嵌入的文件并返回内容和 MIME 类型
func ReadFile(path string) ([]byte, string, error) {
	data, err := staticFiles.ReadFile("dist/" + path)
	if err != nil {
		return nil, "", err
	}
	
	// 根据文件扩展名确定 MIME 类型
	ext := filepath.Ext(path)
	contentType := mime.TypeByExtension(ext)
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	
	return data, contentType, nil
}
