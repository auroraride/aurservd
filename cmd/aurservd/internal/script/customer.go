// Copyright (C) liasica. 2022-present.
//
// Created at 2022-07-27
// Based on aurservd by liasica, magicrolan@qq.com.

package script

import (
    "bytes"
    "fmt"
    "github.com/auroraride/aurservd/internal/ali"
    "github.com/golang-module/carbon/v2"
    "github.com/spf13/cobra"
    "io/ioutil"
    "strings"
    "time"
)

func customerCommand() *cobra.Command {
    cmd := &cobra.Command{
        Use:   "customer",
        Short: "极光出行客服工具",
    }

    cmd.AddCommand(
        customerAppNotice(),
    )

    return cmd
}

func customerAppNotice() *cobra.Command {
    var (
        path string
        tel  string
        tmpl string
    )
    cmd := &cobra.Command{
        Use:   "appnotice",
        Short: "切换APP通知",
        Run: func(cmd *cobra.Command, args []string) {
            var buf bytes.Buffer

            vms := ali.NewVms()

            b, err := ioutil.ReadFile(path)
            if err != nil {
                fmt.Println(err)
                return
            }
            arr := strings.Split(string(b), "\n")
            for _, s := range arr {
                s = strings.ReplaceAll(s, " ", "")
                if s == "" {
                    continue
                }
                state := "成功"
                success := vms.SendVoiceMessageByTts(&s, nil, &tel, &tmpl)
                if !success {
                    state = "失败"
                }
                buf.WriteString(s)
                buf.WriteString("  ")
                buf.WriteString(state)
                buf.WriteString("\n")
            }

            err = ioutil.WriteFile(fmt.Sprintf("runtime/appnotice-%s.txt", time.Now().Format(carbon.ShortDateTimeLayout)), buf.Bytes(), 0755)
            if err != nil {
                fmt.Println(err)
            }
        },
    }

    cmd.Flags().StringVarP(&path, "path", "p", "", "文件路径")
    cmd.Flags().StringVarP(&tel, "tel", "e", "02863804608", "呼叫电话")
    cmd.Flags().StringVarP(&tmpl, "tmpl", "l", "TTS_246450093", "呼叫模板")
    _ = cmd.MarkFlagRequired("path")
    return cmd
}
