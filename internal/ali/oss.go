// Copyright (C) liasica. 2021-present.
//
// Created at 2021/12/13
// Based on aurservd by liasica, magicrolan@qq.com.

package ali

import (
    "bytes"
    "encoding/base64"
    "github.com/aliyun/alibaba-cloud-sdk-go/services/sts"
    "github.com/aliyun/aliyun-oss-go-sdk/oss"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/auroraride/aurservd/pkg/snag"
    "github.com/go-resty/resty/v2"
    "strings"
)

type ossClient struct {
    *oss.Client

    url    string
    Bucket *oss.Bucket
}

func NewOss() *ossClient {
    cfg := ar.Config.Aliyun.Oss
    client, err := oss.New(cfg.Endpoint, cfg.AccessKeyId, cfg.AccessKeySecret)
    if err != nil {
        snag.Panic(err)
    }
    var bucket *oss.Bucket
    bucket, err = client.Bucket(cfg.Bucket)
    if err != nil {
        snag.Panic(err)
    }
    return &ossClient{
        Client: client,
        Bucket: bucket,
        url:    cfg.Url,
    }
}

// UploadUrlFile 从URL获取资源并上传
func (c *ossClient) UploadUrlFile(name string, url string) string {
    r, err := resty.New().R().Get(url)
    if err != nil {
        snag.Panic(err)
    }
    return c.UploadBytes(name, r.Body())
}

// UploadBase64ImageJpeg 上传jpg图片
func (c *ossClient) UploadBase64ImageJpeg(name string, b64 string) string {
    b, err := base64.StdEncoding.DecodeString(b64)
    if err != nil {
        snag.Panic(err)
    }
    return c.UploadBytes(name, b)
}

// UploadBytes 上传文件
func (c *ossClient) UploadBytes(name string, b []byte) string {
    err := c.Bucket.PutObject(name, bytes.NewReader(b))
    if err != nil {
        snag.Panic(err)
    }
    url := c.url
    if !strings.HasSuffix(url, "/") {
        url += "/"
    }
    return url + name
}

// StsToken 获取临时访问token
// @doc https://help.aliyun.com/document_detail/383950.html
func (c *ossClient) StsToken() *model.AliyunOssStsRes {
    cfg := ar.Config.Aliyun.Oss
    client, err := sts.NewClientWithAccessKey(cfg.RegionId, cfg.AccessKeyId, cfg.AccessKeySecret)
    if err != nil {
        snag.Panic(err)
    }

    // 构建请求对象。
    request := sts.CreateAssumeRoleRequest()
    request.Scheme = "https"

    request.RoleArn = cfg.Arn
    request.RoleSessionName = cfg.RamRole

    // 发起请求，并得到响应。
    response, err := client.AssumeRole(request)
    if err != nil {
        snag.Panic(err)
    }

    res := response.Credentials
    return &model.AliyunOssStsRes{
        AccessKeySecret: res.AccessKeySecret,
        Expiration:      res.Expiration,
        AccessKeyId:     res.AccessKeyId,
        StsToken:        res.SecurityToken,
        Bucket:          cfg.Bucket,
    }
}
