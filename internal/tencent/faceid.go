// Copyright (C) liasica. 2024-present.
//
// Created at 2024-02-27
// Based on aurservd by liasica, magicrolan@qq.com.

package tencent

import (
	"errors"
	"log"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	faceid "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/faceid/v20180301"
	"go.uber.org/zap"

	"github.com/auroraride/aurservd/internal/ar"
	"github.com/auroraride/aurservd/pkg/silk"
)

var _faceId *FaceId

type FaceId struct {
	*faceid.Client
}

func NewFaceId() *FaceId {
	return _faceId
}

func BootFaceId() {
	cfg := ar.Config.Tencent.FaceId
	credential := common.NewCredential(
		cfg.SecretId,
		cfg.SecretKey,
	)
	client, err := faceid.NewClient(credential, "", profile.NewClientProfile())
	if err != nil {
		log.Fatalf("腾讯faceid初始化失败: %v", err)
	}

	_faceId = &FaceId{client}
}

func (f *FaceId) IdCardOCR(image string) (result *faceid.IdCardOCRVerificationResponseParams, err error) {
	var res *faceid.IdCardOCRVerificationResponse
	res, err = f.IdCardOCRVerification(&faceid.IdCardOCRVerificationRequest{ImageBase64: &image})
	if err != nil {
		zap.L().Info("身份证OCR识别失败", zap.Error(err))
		return
	}

	result = res.Response
	if result.Result == nil {
		return nil, errors.New("身份证识别失败")
	}

	empty := silk.String("")

	// 性别
	if result.Sex != nil && *result.Sex == "null" {
		result.Sex = empty
	}

	// 民族
	if result.Nation != nil && *result.Nation == "null" {
		result.Nation = empty
	}

	// 生日
	if result.Birth != nil && *result.Birth == "null" {
		result.Birth = empty
	}

	// 地址
	if result.Address != nil && *result.Address == "null" {
		result.Address = empty
	}

	switch *result.Result {
	case "0":
		return
	case "-1":
		return nil, errors.New("姓名和身份证号不一致")
	case "-2":
		return nil, errors.New("非法身份证号")
	case "-3":
		return nil, errors.New("非法姓名")
	case "-4":
		return nil, errors.New("证件库服务异常")
	case "-5":
		return nil, errors.New("证件库中无此身份证记录")
	case "-6":
		return nil, errors.New("权威比对系统升级中，请稍后再试")
	case "-7":
		return nil, errors.New("认证次数超过当日限制")
	}
	return
}
