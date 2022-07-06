// Copyright (C) liasica. 2022-present.
//
// Created at 2022-07-05
// Based on aurservd by liasica, magicrolan@qq.com.

package permission

import (
    "github.com/auroraride/aurservd/pkg/utils"
    "github.com/fsnotify/fsnotify"
    log "github.com/sirupsen/logrus"
    "github.com/spf13/viper"
    "gopkg.in/yaml.v3"
    "os"
)

const (
    permFile = "config/permission.yaml"
)

var Groups map[string]*Group
var Items map[string]Item
var Keys []string

type Group struct {
    Name        string `json:"name"`  // 权限分组名称
    Desc        string `json:"desc"`  // 权限分组描述
    Permissions []Item `json:"items"` // 权限列表
}

type Item struct {
    Key    string `json:"key"`    // 权限KEY: MD5(METHOD + PATH)
    Method string `json:"method"` // api请求method
    Api    string `json:"api"`    // api请求path
    Desc   string `json:"desc"`   // 接口描述
    SN     string `json:"sn"`     // 接口编号
}

var v *viper.Viper

func init() {
    v = viper.New()
    v.SetConfigFile(permFile)

    load()

    v.OnConfigChange(func(e fsnotify.Event) {
        log.Infof("权限文件已改动: %s, 重载权限: %v", e.Name, read())
    })

    v.WatchConfig()
}

func GetKey(method, api string) string {
    return utils.Md5String(method + api)
}

func read() error {
    // 读取权限
    err := v.ReadInConfig()
    if err != nil {
        return err
    }
    Groups = make(map[string]*Group)
    err = v.Unmarshal(&Groups)
    Items = make(map[string]Item)
    Keys = make([]string, 0)
    for _, group := range Groups {
        for _, permission := range group.Permissions {
            Items[permission.Key] = permission
            Keys = append(Keys, permission.Key)
        }
    }
    return err
}

func load() {
    // 判断文件是否存在
    f := utils.NewFile(permFile)
    if !f.IsExist() {
        err := f.CreateDirectoryIfNotExist()
        if err != nil {
            log.Fatalf("权限目录创建失败: %v", err)
            return
        }
        err = os.WriteFile(permFile, []byte(""), 0644)
        if err != nil {
            log.Fatalf("默认权限保存失败: %v", err)
            return
        }
    }

    err := read()
    if err != nil {
        log.Fatalf("权限读取失败: %v", err)
    }
}

func Save(m map[string]*Group) {
    b, err := yaml.Marshal(m)
    if err != nil {
        log.Fatalf("权限配置保存失败: %#v", err)
    }
    _ = os.WriteFile(permFile, b, 0644)
}
