// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-20
// Based on aurservd by liasica, magicrolan@qq.com.

package logging

import (
    sls "github.com/aliyun/aliyun-log-go-sdk"
    "github.com/auroraride/aurservd/internal/ali"
    "github.com/auroraride/aurservd/internal/ar"
    log "github.com/sirupsen/logrus"
    "reflect"
)

var (
    indexToken = []string{`,`, ` `, `'`, `"`, `;`, `=`, `(`, `)`, `[`, `]`, `{`, `}`, `?`, `@`, `&`, `<`, `>`, `/`, `:`, `\n`, `\t`, `\r`}
)

func Boot() {
    cfg := ar.Config.Aliyun.Sls
    bootLogStore(cfg.Project, cfg.CabinetLog, CabinetLog{})
    bootLogStore(cfg.Project, cfg.DoorLog, DoorOperateLog{})
    bootLogStore(cfg.Project, cfg.OperateLog, OperateLog{})
}

// bootLogStore 自动创建logstore
// SDK 参考 https://help.aliyun.com/document_detail/286951.html
// 查询语法: https://help.aliyun.com/document_detail/29060.htm?spm=a2c4g.11186623.0.0.4e902a73u4EmXO#concept-tnd-1jq-zdb
func bootLogStore(project, logstore string, typ any) {
    var err error
    slsc := ali.NewSls()
    _, err = slsc.GetLogStore(project, logstore)
    if err != nil {
        e := err.(*sls.Error)
        if e.Code == sls.LOGSTORE_NOT_EXIST {
            // 创建 logstore
            err = slsc.CreateLogStoreV2(project, &sls.LogStore{
                Name:       logstore,
                TTL:        3650,
                ShardCount: 2,
            })
            if err != nil {
                log.Fatal(err)
            }

            // 创建index
            ikm := make(map[string]sls.IndexKey)
            t := reflect.TypeOf(typ)
            for i := 0; i < t.NumField(); i++ {
                f := t.Field(i)
                idx, ok := f.Tag.Lookup("index")
                if !ok {
                    continue
                }
                // tag := f.Tag.Get("sls")
                json := f.Tag.Get("json")
                fk := f.Type.Kind()
                if isString, ok := f.Tag.Lookup("string"); ok && isString == "true" {
                    fk = reflect.String
                }

                var ik sls.IndexKey
                switch fk {
                case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
                    reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
                    reflect.Uintptr:
                    ik = sls.IndexKey{
                        Type: "long",
                    }
                    break
                case reflect.Float32, reflect.Float64:
                    ik = sls.IndexKey{
                        Type: "double",
                    }
                    break
                default:
                    ik = sls.IndexKey{
                        Token: indexToken,
                        Type:  "text",
                        Chn:   true,
                    }
                    break
                }
                if idx == "doc" {
                    ik.DocValue = true
                }
                ik.Alias = json
                ikm[json] = ik
            }
            err = slsc.CreateIndex(project, logstore, sls.Index{
                Line: &sls.IndexLine{Token: indexToken, Chn: true},
                Keys: ikm,
            })
            if err != nil {
                log.Fatal(err)
            }
        }
    }
}
