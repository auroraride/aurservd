// Copyright (C) liasica. 2022-present.
//
// Created at 2022-08-25
// Based on aurservd by liasica, magicrolan@qq.com.

package script

import (
    "github.com/auroraride/aurservd/internal/esign"
    "github.com/spf13/cobra"
)

func personCommand() *cobra.Command {
    cmd := &cobra.Command{
        Use:   "person",
        Short: "骑手信息工具",
    }
    cmd.AddCommand(
        modifyPersonCommand(),
    )
    return cmd
}

func modifyPersonCommand() *cobra.Command {
    var (
        accountID string
        name      string
    )

    cmd := &cobra.Command{
        Use:   "esign",
        Short: "修改骑手E签宝信息",
        Run: func(_ *cobra.Command, _ []string) {
            esign.New().ModifyAccount(accountID, esign.PersonAccountReq{Name: name})
        },
    }

    cmd.Flags().StringVar(&accountID, "account", "", "E签宝accountID")
    cmd.Flags().StringVar(&name, "name", "", "骑手姓名")

    _ = cmd.MarkFlagRequired("account")
    _ = cmd.MarkFlagRequired("name")

    return cmd
}
