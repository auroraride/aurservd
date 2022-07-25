// Copyright (C) liasica. 2022-present.
//
// Created at 2022-07-25
// Based on aurservd by liasica, magicrolan@qq.com.

package internal

import (
    "github.com/aliyun/aliyun-oss-go-sdk/oss"
    "github.com/auroraride/aurservd/pkg/snag"
)

func alioss() *oss.Bucket {
    client, err := oss.New(cfg.Endpoint, cfg.AccessKeyId, cfg.AccessKeySecret)
    if err != nil {
        snag.Panic(err)
    }
    var bucket *oss.Bucket
    bucket, err = client.Bucket(cfg.Bucket)
    if err != nil {
        snag.Panic(err)
    }
    return bucket
}
