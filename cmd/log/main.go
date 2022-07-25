// Copyright (C) liasica. 2022-present.
//
// Created at 2022-07-25
// Based on aurservd by liasica, magicrolan@qq.com.

package main

import (
    "github.com/auroraride/aurservd/cmd/log/internal"
    "github.com/spf13/cobra"
    "log"
    "os"
    "os/signal"
    "syscall"
)

func main() {
    c := make(chan os.Signal, 1)
    signal.Notify(c, os.Interrupt, syscall.SIGTERM, os.Kill)
    go func() {
        <-c
        log.Printf("程序退出\n\n")
        os.Exit(0)
    }()

    cmd := &cobra.Command{
        Use:               "aurlog",
        Short:             "极光出行日志管理",
        CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
    }
    cmd.AddCommand(
        internal.Install(),
        internal.Run(),
    )
    err := cmd.Execute()
    if err != nil {
        log.Fatal(err)
    }
}
