// Copyright (C) liasica. 2022-present.
//
// Created at 2022-07-25
// Based on aurservd by liasica, magicrolan@qq.com.

package internal

import (
    _ "embed"
    "fmt"
    "github.com/spf13/cobra"
    "io/ioutil"
    "os"
    "os/exec"
)

var (
    //go:embed aurlog.service
    assetService string

    //go:embed aurlog.yaml
    assetConfig []byte
)

func Install() *cobra.Command {
    return &cobra.Command{
        Use:   "install",
        Short: "安装 AurLog 服务",
        Run: func(cmd *cobra.Command, args []string) {
            installLog()
            installService()
            installConfig()

            _, err := exec.Command("systemctl", "daemon-reload").Output()
            if err != nil {
                fmt.Printf("已完成安装, 执行 systemctl daemon-reload 失败: %s, 请手动执行\n", err)
            }
        },
    }
}

func installLog() {
    p := "/var/log/aurlog/aurlog.pid"
    err := createDirIfNotExists(p)
    if err != nil {
        FmtFatalln(err)
    }
}

func installService() {
    p := "/etc/systemd/system/aurlog.service"
    e, err := os.Executable()
    if err != nil {
        FmtFatalln(err)
    }
    assetService = fmt.Sprintf(assetService, e)
    err = ioutil.WriteFile(p, []byte(assetService), 0644)
    if err != nil {
        FmtFatalln(err)
    }
}

func installConfig() {
    p := "/etc/aurlog/aurlog.yaml"
    err := createDirIfNotExists(p)
    if err != nil {
        FmtFatalln(err)
    }

    err = ioutil.WriteFile(p, assetConfig, 0644)
    if err != nil {
        FmtFatalln(err)
    }
}
