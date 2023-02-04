// Copyright (C) liasica. 2022-present.
//
// Created at 2022-07-05
// Based on aurservd by liasica, magicrolan@qq.com.

package permission

import (
    "github.com/auroraride/aurservd/pkg/utils"
    "github.com/fsnotify/fsnotify"
    "github.com/spf13/viper"
    "go.uber.org/zap"
    "golang.org/x/exp/slices"
    "gopkg.in/yaml.v3"
    "os"
)

const (
    PermFile = "config/permission.yaml"
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
    v.SetConfigFile(PermFile)

    load()

    v.OnConfigChange(func(e fsnotify.Event) {
        _ = read()
    })

    v.WatchConfig()
}

func GetKey(method, api string) string {
    return utils.Md5String(method + api)
}

func Contains(method, api string, perms []string) bool {
    key := GetKey(method, api)

    // 无需校验的时候直接返回true
    if _, ok := Items[key]; !ok {
        return true
    }

    return slices.Contains(perms, key)
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
    f := utils.NewFile(PermFile)
    if !f.IsExist() {
        err := f.CreateDirectoryIfNotExist()
        if err != nil {
            zap.L().Fatal("权限目录创建失败", zap.Error(err))
            return
        }
        err = os.WriteFile(PermFile, []byte(""), 0644)
        if err != nil {
            zap.L().Fatal("默认权限保存失败", zap.Error(err))
            return
        }
    }

    err := read()
    if err != nil {
        zap.L().Fatal("权限读取失败: %v", zap.Error(err))
    }
}

func Save(m map[string]*Group) {
    b, err := yaml.Marshal(m)
    if err != nil {
        zap.L().Fatal("权限配置保存失败: %#v", zap.Error(err))
    }
    _ = os.WriteFile(PermFile, b, 0644)
}
