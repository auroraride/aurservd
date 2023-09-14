// Copyright (C) liasica. 2022-present.
//
// Created at 2022-07-25
// Based on aurservd by liasica, magicrolan@qq.com.

package internal

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/golang-module/carbon/v2"
	"github.com/spf13/cobra"
)

func Run() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "run",
		Short: "启动 AurLog 日志监控",
		Run: func(cmd *cobra.Command, args []string) {
			ParseConfig()

			run()
		},
	}
	cmd.Flags().StringVarP(&configFile, "config", "c", "/etc/aurlog/aurlog.yaml", "配置路径")

	return cmd
}

func run() {
	log.Println("开始监听")
	for {
		doLogs()
		time.Sleep(12 * time.Hour)
	}
}

func doLogs() {
	root := cfg.LogPath
	if !strings.HasSuffix(root, "/") {
		root += "/"
	}
	_ = filepath.Walk(root, func(path string, info fs.FileInfo, _ error) error {
		if info.IsDir() {
			return nil
		}
		if carbon.CreateFromStdTime(info.ModTime()).IsToday() {
			return nil
		}

		key := strings.Replace(path, root, "", 1)
		if !strings.Contains(key, "/") {
			key = fmt.Sprintf("app/%s", key)
		}

		log.Printf("即将处理: %s", key)

		// 上传阿里云
		start := time.Now()
		err := alioss().PutObjectFromFile(key, path)
		if err != nil {
			log.Printf("%s 处理失败: %s", key, err)
			return nil
		}

		log.Printf("%s 处理成功: %.2fMB - %.4fs", key, float64(info.Size())/1024.0/1024.0, time.Since(start).Seconds())

		// 删除
		_ = os.Remove(path)

		return nil
	})
}
