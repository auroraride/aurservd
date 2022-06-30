// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-20
// Based on aurservd by liasica, magicrolan@qq.com.

package logging

import (
    "github.com/alibabacloud-go/tea/tea"
    sls "github.com/aliyun/aliyun-log-go-sdk"
    "github.com/auroraride/aurservd/internal/ali"
    "github.com/auroraride/aurservd/internal/ar"
    jsoniter "github.com/json-iterator/go"
    log "github.com/sirupsen/logrus"
    "reflect"
    "strconv"
    "time"
)

var (
    indexToken = []string{`,`, ` `, `'`, `"`, `;`, `=`, `(`, `)`, `[`, `]`, `{`, `}`, `?`, `@`, `&`, `<`, `>`, `/`, `:`, `\n`, `\t`, `\r`}
)

type Logger interface {
    GetLogstoreName() string
    Send()
}

func Boot() {
    cfg := ar.Config.Aliyun.Sls
    bootLogStore(cfg.Project, cfg.CabinetLog, CabinetLog{})
    bootLogStore(cfg.Project, cfg.DoorLog, DoorOperateLog{})
    bootLogStore(cfg.Project, cfg.OperateLog, OperateLog{})
    bootLogStore(cfg.Project, cfg.ExchangeLog, ExchangeLog{})
    bootLogStore(cfg.Project, cfg.HealthLog, HealthLog{})
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

// PutLog 提交日志
func PutLog(ptr Logger) {
    go func() {
        err := ali.NewSls().PutLogs(ar.Config.Aliyun.Sls.Project, ptr.GetLogstoreName(), &sls.LogGroup{
            Logs: []*sls.Log{{
                Time:     tea.Uint32(uint32(time.Now().Unix())),
                Contents: GenerateLogContent(ptr),
            }},
        })
        if err != nil {
            log.Error(err)
            return
        }
    }()
}

func GetCount(logstore string, query string, from time.Time) (total int) {
    var cnt []struct {
        Count string `json:"count"`
    }
    cfg := ar.Config.Aliyun.Sls
    response, err := ali.NewSls().GetLogsV2(cfg.Project, logstore, &sls.GetLogRequest{
        From:    from.Unix(),
        To:      time.Now().Unix(),
        Reverse: true,
        Query:   query + " | SELECT COUNT(*) as count",
    })
    if err != nil {
        log.Error(err)
        return
    }
    var b []byte
    b, err = jsoniter.Marshal(response.Logs)
    if err != nil {
        return
    }
    err = jsoniter.Unmarshal(b, &cnt)
    if err != nil {
        return
    }
    if len(cnt) < 1 {
        return
    }
    total, _ = strconv.Atoi(cnt[0].Count)
    return
}
