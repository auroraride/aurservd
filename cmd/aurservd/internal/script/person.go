// Copyright (C) liasica. 2022-present.
//
// Created at 2022-08-25
// Based on aurservd by liasica, magicrolan@qq.com.

package script

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/auroraride/aurservd/internal/esign"
)

func personCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "person",
		Short: "骑手信息工具",
	}
	cmd.AddCommand(
		esignCommand(),
	)
	return cmd
}

func esignCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "esign",
		Short: "E签宝服务",
	}
	cmd.AddCommand(
		modifyPersonCommand(),
		queryPersonCommand(),
		createPersonCommand(),
	)
	return cmd
}

func modifyPersonCommand() *cobra.Command {
	var (
		accountID string
		name      string
		idNumber  string
	)

	cmd := &cobra.Command{
		Use:   "modify",
		Short: "修改骑手E签宝信息",
		Run: func(_ *cobra.Command, _ []string) {
			esign.New().ModifyAccount(accountID, esign.PersonAccountReq{Name: name, IdNumber: idNumber})
		},
	}

	cmd.Flags().StringVar(&accountID, "account", "", "E签宝accountID")
	cmd.Flags().StringVar(&name, "name", "", "骑手姓名")
	cmd.Flags().StringVar(&idNumber, "idNumber", "", "骑手身份证信息")

	_ = cmd.MarkFlagRequired("account")
	_ = cmd.MarkFlagRequired("name")
	_ = cmd.MarkFlagRequired("idNumber")

	return cmd
}

func queryPersonCommand() *cobra.Command {
	var (
		accountID string
	)

	cmd := &cobra.Command{
		Use:   "query",
		Short: "查询骑手E签宝信息",
		Run: func(_ *cobra.Command, _ []string) {
			esign.New().QueryAccount(accountID)
		},
	}

	cmd.Flags().StringVar(&accountID, "account", "", "E签宝accountID")

	_ = cmd.MarkFlagRequired("account")

	return cmd
}

func createPersonCommand() *cobra.Command {
	var (
		name     string
		idNumber string
		phone    string
	)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "创建E签宝账户",
		Run: func(_ *cobra.Command, _ []string) {
			accountId := esign.New().CreatePersonAccount(esign.PersonAccountReq{
				ThirdPartyUserId: idNumber,
				Name:             name,
				IdType:           "CRED_PSN_CH_IDCARD",
				IdNumber:         idNumber,
				Mobile:           phone,
			})
			fmt.Printf("创建成功，accountId = %s\n\n", accountId)
		},
	}

	cmd.Flags().StringVar(&phone, "phone", "", "骑手电话号码")
	cmd.Flags().StringVar(&name, "name", "", "骑手姓名")
	cmd.Flags().StringVar(&idNumber, "idNumber", "", "骑手身份证信息")

	_ = cmd.MarkFlagRequired("phone")
	_ = cmd.MarkFlagRequired("name")
	_ = cmd.MarkFlagRequired("idNumber")

	return cmd
}
