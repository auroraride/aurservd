// Copyright (C) liasica. 2022-present.
//
// Created at 2022-07-05
// Based on aurservd by liasica, magicrolan@qq.com.

package main

import (
    "fmt"
    "github.com/auroraride/aurservd/app/permission"
    "github.com/povsister/scp"
    "io/ioutil"
    "os"
    "path/filepath"
    "regexp"
    "sort"
    "strings"
)

var (
    nameSkipper = map[string]bool{
        "selection": true,
    }
    apiSkipper = map[string]bool{
        "/manager/v1/permission":    true,
        "/manager/v1/user/signin":   true,
        "/manager/v1/battery/model": true,
        "/manager/v1/city":          true,
    }
)

func main() {
    d := "./app/controller/v1/mapi"
    files, _ := os.ReadDir(d)
    m := make(map[string]*permission.Group)
    for _, f := range files {
        if f.IsDir() {
            continue
        }
        name := strings.TrimSuffix(f.Name(), filepath.Ext(f.Name()))
        if nameSkipper[name] {
            continue
        }

        pg, ok := m[name]
        if !ok {
            pg = &permission.Group{
                Name: name,
                Desc: "",
            }
            m[name] = pg
        }

        // 已经保存的desc获取
        if x, ok := permission.Groups[name]; ok {
            pg.Desc = x.Desc
        }

        doc, _ := ioutil.ReadFile(filepath.Join(d, f.Name()))
        re := regexp.MustCompile(`(?m)// @Router\s+(.*) \[(.*)][\S\s]*?@Summary\s+(.*)? (.*)`)
        bs := re.FindAllStringSubmatch(string(doc), -1)
        for _, sub := range bs {
            api := sub[1]
            if apiSkipper[api] {
                continue
            }
            method := sub[2]
            sn := sub[3]
            desc := sub[4]
            pg.Permissions = append(pg.Permissions, permission.Item{
                Key:    permission.GetKey(method, api),
                Method: method,
                Api:    api,
                Desc:   desc,
                SN:     sn,
            })
        }

        sort.Slice(pg.Permissions, func(i, j int) bool {
            return strings.Compare(pg.Permissions[i].SN, pg.Permissions[j].SN) < 0
        })
    }

    permission.Save(m)

    upload()
}

func upload() {
    privPEM, err := ioutil.ReadFile("/Users/liasica/.ssh/id_rsa")
    sshConf, _ := scp.NewSSHConfigFromPrivateKey("root", privPEM)
    client, err := scp.NewClient("39.106.77.239", sshConf, &scp.ClientOption{})
    if err != nil {
        fmt.Println("ssh connect error ", err)
        return
    }
    defer func(client *scp.Client) {
        _ = client.Close()
    }(client)

    err = client.CopyFileToRemote(permission.PermFile, "/var/www/api.auroraride.com/config/permission.yaml", &scp.FileTransferOption{})
    if err != nil {
        fmt.Println("[api] Error while copying file ", err)
        return
    }

    err = client.CopyFileToRemote(permission.PermFile, "/var/www/next-api.auroraride.com/config/permission.yaml", &scp.FileTransferOption{})
    if err != nil {
        fmt.Println("[next-api] Error while copying file ", err)
        return
    }
}
