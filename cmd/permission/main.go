// Copyright (C) liasica. 2022-present.
//
// Created at 2022-07-05
// Based on aurservd by liasica, magicrolan@qq.com.

package main

import (
    "github.com/auroraride/aurservd/app/permission"
    "github.com/auroraride/aurservd/pkg/utils"
    "io/ioutil"
    "os"
    "path/filepath"
    "regexp"
    "sort"
    "strings"
)

var (
    skipper = map[string]bool{
        "selection": true,
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
        if skipper[name] {
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
            method := sub[2]
            sn := sub[3]
            desc := sub[4]
            pg.Permissions = append(pg.Permissions, permission.Item{
                Key:    utils.Md5String(method + api),
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
}
