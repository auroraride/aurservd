// Copyright (C) liasica. 2024-present.
//
// Created at 2024-02-28
// Based on aurservd by liasica, magicrolan@qq.com.

package ali

import (
	"encoding/hex"
	"time"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	openapiutil "github.com/alibabacloud-go/openapi-util/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/sts"

	"github.com/auroraride/aurservd/internal/ar"
	"github.com/auroraride/aurservd/pkg/utils"
)

const (
	OcrActionRecognizeIdcard = "RecognizeIdcard"
	OcrVersion               = "2021-07-07"
	OcrContentType           = "application/octet-stream"
	signedHeaders            = "content-type;host;x-acs-action;x-acs-content-sha256;x-acs-date;x-acs-security-token;x-acs-signature-nonce;x-acs-version"
)

type OcrClient struct {
	*openapi.Client

	EndPoint           *string
	SignatureAlgorithm *string
}

var _ocrClient *OcrClient

// NewOcr Ocr客户端单例模式
func NewOcr() *OcrClient {
	if _ocrClient == nil {
		cfg := ar.Config.Aliyun.Ocr
		config := &openapi.Config{
			AccessKeyId:     &cfg.AccessKeyId,
			AccessKeySecret: &cfg.AccessKeySecret,
		}
		config.Endpoint = &cfg.Endpoint
		_client, _ := openapi.NewClient(config)
		_ocrClient = &OcrClient{
			Client:             _client,
			EndPoint:           &cfg.Endpoint,
			SignatureAlgorithm: tea.String("ACS3-HMAC-SHA256"),
		}
	}
	return _ocrClient
}

type OcrParams struct {
	ContentType   string
	Action        string
	Date          string
	Token         string
	Nonce         string
	Version       string
	Authorization string
}

// Signature 获取阿里云Ocr签名
func (c *OcrClient) Signature(hash string) (params *OcrParams, err error) {
	cfg := ar.Config.Aliyun.Ocr

	var credentials *sts.Credentials
	credentials, err = stsToken(cfg.RegionId, cfg.AccessKeyId, cfg.AccessKeySecret, cfg.Arn, cfg.RamRole)
	if err != nil {
		return
	}

	params = &OcrParams{
		ContentType: OcrContentType,
		Action:      OcrActionRecognizeIdcard,
		Date:        time.Now().UTC().Format("2006-01-02T15:04:05Z"),
		Token:       credentials.SecurityToken,
		Nonce:       utils.RandStr(32),
		Version:     OcrVersion,
	}

	canonicalRequest := "POST\n" +
		"/\n" +
		"OutputFigure=true&OutputQualityInfo=true\n" + // CanonicalQueryString，即使为空也需要填写一个换行符
		"content-type:" + OcrContentType + "\n" +
		"host:" + cfg.Endpoint + "\n" +
		"x-acs-action:" + OcrActionRecognizeIdcard + "\n" +
		"x-acs-content-sha256:" + hash + "\n" +
		"x-acs-date:" + params.Date + "\n" +
		"x-acs-security-token:" + params.Token + "\n" +
		"x-acs-signature-nonce:" + params.Nonce + "\n" +
		"x-acs-version:" + OcrVersion + "\n" +
		"\n" + // 消息头后面需要加一个换行符
		signedHeaders + "\n" +
		hash

	signType := *c.SignatureAlgorithm
	stringToSign := signType + "\n" + hex.EncodeToString(openapiutil.Hash([]byte(canonicalRequest), c.SignatureAlgorithm))
	signature := hex.EncodeToString(openapiutil.SignatureMethod(credentials.AccessKeySecret, stringToSign, signType))
	params.Authorization = signType + " Credential=" + credentials.AccessKeyId + ",SignedHeaders=" + signedHeaders + ",Signature=" + signature

	return
}
