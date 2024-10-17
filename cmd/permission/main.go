// Copyright (C) liasica. 2022-present.
//
// Created at 2022-07-05
// Based on aurservd by liasica, magicrolan@qq.com.

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/povsister/scp"

	"github.com/auroraride/aurservd/app/permission"
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
	ds := []string{
		"./app/controller/v1/mapi",
		"./app/controller/v2/mapi",
		"./app/purchase/controller/mapi",
	}
	m := make(map[string]*permission.Group)
	for _, d := range ds {
		files, _ := os.ReadDir(d)
		for _, f := range files {
			if f.IsDir() || f.Name() == "mapi.go" {
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

			doc, _ := os.ReadFile(filepath.Join(d, f.Name()))
			re := regexp.MustCompile(`(?m)// @Router\s+(.*) \[(.*)][\S\s]*?@Summary\s+(.*)\s+(.*)`)
			bs := re.FindAllStringSubmatch(string(doc), -1)
			for _, sub := range bs {
				api := sub[1]
				if apiSkipper[api] {
					continue
				}
				method := sub[2]
				desc := sub[3]
				pg.Permissions = append(pg.Permissions, permission.Item{
					Key:    permission.GetKey(method, api),
					Method: method,
					Api:    api,
					Desc:   desc,
				})
			}
		}
	}
	permission.Save(m)

	addrs := []string{
		"118.116.4.16:26610",
		"118.116.4.16:26611",
	}

	for _, addr := range addrs {
		upload(addr)
	}
}

func upload(addr string) {
	privPEM, _ := os.ReadFile("~/.ssh/id_rsa")
	sshConf, _ := scp.NewSSHConfigFromPrivateKey("root", privPEM)
	client, err := scp.NewClient(addr, sshConf, &scp.ClientOption{})
	if err != nil {
		fmt.Println("ssh connect error ", err)
		return
	}
	defer func(client *scp.Client) {
		_ = client.Close()
	}(client)

	err = client.CopyFileToRemote(permission.PermFile, "/var/www/aurservd/config/permission.yaml", &scp.FileTransferOption{})
	if err != nil {
		fmt.Println("[api] Error while copying file ", err)
		return
	}
}
