// Copyright (C) liasica. 2021-present.
//
// Created at 2021/12/10
// Based on aurservd by liasica, magicrolan@qq.com.

package internal

import (
    "github.com/spf13/cobra"
    "io/ioutil"
    "os"
    "path/filepath"
)

const (
    defaultEntPath = "./internal/ent"
)

type FileInfoType uint8

const (
    FileInfoTypeFile = iota
    FileInfoTypeDir
)

var (
    removeNotable = map[string]FileInfoType{
        "db.go":    FileInfoTypeFile,
        "internal": FileInfoTypeDir,
        "schema":   FileInfoTypeDir,
    }
)

func CleanCmd() *cobra.Command {
    return &cobra.Command{
        Use:   "clean [path]",
        Short: "clean generated go code for the ent directory",
        Run: func(cmd *cobra.Command, path []string) {
            p := defaultEntPath
            if len(path) > 0 {
                p = path[0]
            }
            fs, _ := ioutil.ReadDir(p)
            for _, f := range fs {
                name := f.Name()
                if _, ok := removeNotable[name]; !ok {
                    _ = os.RemoveAll(filepath.Join(p, name))
                }
            }
        },
    }
}
