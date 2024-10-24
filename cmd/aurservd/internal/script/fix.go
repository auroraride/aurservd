// Copyright (C) liasica. 2022-present.
//
// Created at 2022-08-18
// Based on aurservd by liasica, magicrolan@qq.com.

package script

import (
	"github.com/spf13/cobra"

	"github.com/auroraride/aurservd/cmd/aurservd/internal/script/fix"
)

func fixCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "fix",
		Short: "修复指令",
	}

	cmd.AddCommand(
		fix.Commission(),
		fix.Reminder(),
		fix.Assistance(),
		fix.Contract(),
		fix.Sql(),
		fix.Asset(),
		fix.AssetSN(),
	)

	return cmd
}
