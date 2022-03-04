// Copyright (C) liasica. 2022-present.
//
// Created at 2022-03-04
// Based on aurservd by liasica, magicrolan@qq.com.

package script

import (
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/app/service"
    log "github.com/sirupsen/logrus"
    "github.com/spf13/cobra"
)

var managerCmd = &cobra.Command{
    Use:   "manager",
    Short: "管理员操作",
}

func managerAddCmd() *cobra.Command {
    var (
        name     string
        phone    string
        password string
    )
    cmd := &cobra.Command{
        Use:   "add",
        Short: "添加管理员",
        Run: func(cmd *cobra.Command, args []string) {
            req := &model.ManagerAddReq{
                Name: name,
            }
            req.Phone = phone
            req.Password = password
            err := service.NewManager().Add(req)
            if err != nil {
                log.Errorf("添加管理员: %s %s 失败: %v", name, phone, err)
                return
            }
            log.Printf("添加管理员: %s %s 成功", name, phone)
        },
    }
    cmd.Flags().StringVarP(&name, "name", "n", "", "管理员姓名")
    cmd.Flags().StringVarP(&phone, "phone", "u", "", "管理员手机号")
    cmd.Flags().StringVarP(&password, "password", "p", "", "密码")

    _ = cmd.MarkFlagRequired("name")
    _ = cmd.MarkFlagRequired("phone")
    _ = cmd.MarkFlagRequired("password")
    return cmd
}

func init() {
    managerCmd.AddCommand(managerAddCmd())
}
