// Copyright (C) liasica. 2021-present.
//
// Created at 2021-12-10
// Based on aurservd by liasica, magicrolan@qq.com.

package utils

import (
	"net/http"
	"os"
	"path/filepath"
)

type file struct {
	Path string
}

func NewFile(p string) *file {
	return &file{p}
}

// IsExist 文件是否存在
func (f *file) IsExist() bool {
	if _, err := os.Stat(f.Path); os.IsNotExist(err) {
		return false
	}
	return true
}

// CreateDirectoryIfNotExist 若目录不存在则创建
func (f *file) CreateDirectoryIfNotExist() error {
	d := filepath.Dir(f.Path)
	if _, err := os.Stat(d); os.IsNotExist(err) {
		return os.MkdirAll(d, 0755)
	}
	return nil
}

func GetFileContentType(out *os.File) (string, error) {

	// Only the first 512 bytes are used to sniff the content type.
	buffer := make([]byte, 512)

	_, err := out.Read(buffer)
	if err != nil {
		return "", err
	}

	// Use the net/http package's handy DectectContentType function. Always returns a valid
	// content-type by returning "application/octet-stream" if no others seemed to match.
	contentType := http.DetectContentType(buffer)

	return contentType, nil
}
