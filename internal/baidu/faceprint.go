// Copyright (C) liasica. 2021-present.
//
// Created at 2021/12/13
// Based on aurservd by liasica, magicrolan@qq.com.

package baidu

import (
	"fmt"

	"github.com/auroraride/adapter/log"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"

	"github.com/auroraride/aurservd/internal/ar"
	"github.com/auroraride/aurservd/pkg/snag"
	"github.com/auroraride/aurservd/pkg/utils"
)

const (
	// TypeAuth 实名认证
	TypeAuth = "authenticator"
	// TypeFace 人脸校验
	TypeFace = "face"
)

const (
	// faceprintTokenUrl 获取人脸识别 verify_token 接口URL
	faceprintTokenUrl = `https://aip.baidubce.com/rpc/2.0/brain/solution/faceprint/verifyToken/generate?access_token=%s`

	// faceprintAuthUrl 获取人脸识别认证接口URL
	faceprintAuthUrl = `https://brain.baidu.com/face/print/?token=%s&successUrl=%s&failedUrl=%s`

	// faceprintSimpleUrl 获取人脸图片接口URL
	faceprintSimpleUrl = `https://aip.baidubce.com/rpc/2.0/brain/solution/faceprint/result/simple?access_token=%s`

	// faceprintResultUrl 获取人脸识别认证结果详细信息接口URL
	faceprintResultUrl = `https://aip.baidubce.com/rpc/2.0/brain/solution/faceprint/result/detail?access_token=%s`

	// faceprintSubmitUrl 指定用户信息上报接口URL
	faceprintSubmitUrl = `https://brain.baidu.com/solution/faceprint/idcard/submit?access_token=%s`
)

// faceprintTokenResp 人脸识别 verify_token 返回体
type faceprintTokenResp struct {
	Success bool `json:"success"`
	Result  struct {
		VerifyToken string `json:"verify_token"`
	} `json:"result"`
	LogId string `json:"log_id"`
}

// faceprintFaceResp 人脸认证结果
type faceprintFaceResp struct {
	Success bool `json:"success"`
	Result  struct {
		Image string `json:"image"`
	} `json:"result"`
	LogId string `json:"log_id"`
}

// faceprintDetailResp 人脸识别详细结果
type faceprintDetailResp struct {
	Success bool   `json:"success"`
	LogId   string `json:"log_id"`
	Result  struct {
		FaceImg      string `json:"faceImg,omitempty"`
		VerifyResult struct {
			Score         float64 `json:"score"`
			LivenessScore float64 `json:"liveness_score"`
			Spoofing      float64 `json:"spoofing"`
		} `json:"verify_result"`
		IdcardOcrResult struct {
			Birthday       string `json:"birthday"`
			IssueAuthority string `json:"issue_authority"`
			Address        string `json:"address"`
			Gender         string `json:"gender"`
			Nation         string `json:"nation"`
			ExpireTime     string `json:"expire_time"`
			Name           string `json:"name"`
			IssueTime      string `json:"issue_time"`
			IdCardNumber   string `json:"id_card_number"`
		} `json:"idcard_ocr_result"`
		IdcardImages struct {
			FrontBase64 string `json:"front_base64"`
			BackBase64  string `json:"back_base64"`
		} `json:"idcard_images"`
		IdcardConfirm struct {
			IdcardNumber string `json:"idcard_number"`
			Name         string `json:"name"`
		} `json:"idcard_confirm"`
	} `json:"result"`
}

// faceprintSubmitResp 指定用户信息上报返回
type faceprintSubmitResp struct {
	VerifyToken     string `json:"verify_token"`
	IdName          string `json:"id_name"`
	IdNo            string `json:"id_no"`
	CertificateType int    `json:"certificate_type"`
}

// getFaceprintUrl 获取人身核验URL
func (b *faceClient) getFaceprintUrl(typ string) (url string, token string) {
	var err error
	var planId string
	cfg := ar.Config.Baidu.Face
	switch typ {
	case TypeAuth:
		planId = cfg.AuthPlanId
	case TypeFace:
		planId = cfg.FacePlanId
	}

	var r *resty.Response
	// 获取 verify_token
	res := new(faceprintTokenResp)
	r, err = resty.New().R().
		SetResult(res).
		SetBody(map[string]string{"plan_id": planId}).
		Post(fmt.Sprintf(faceprintTokenUrl, b.accessToken))
	if err != nil {
		snag.Panic(err)
	}
	zap.L().Info("获取人身核验URL", log.ResponseBody(r.Body()))
	if !res.Success {
		snag.Panic("实名认证请求失败")
	}
	token = res.Result.VerifyToken
	str := fmt.Sprintf("%s?type=%s&token=%s&state=", cfg.Redirect, typ, token)
	url = fmt.Sprintf(
		faceprintAuthUrl,
		token,
		utils.EncodeURIComponent(str+"success"),
		utils.EncodeURIComponent(str+"failed"),
	)
	return
}

// GetFaceUrl 获取人脸校验URL
func (b *faceClient) GetFaceUrl(name, icNum string) (uri string, token string) {
	uri, token = b.getFaceprintUrl(TypeFace)
	res := new(faceprintSubmitResp)
	r, err := resty.New().R().
		SetResult(res).
		SetBody(ar.Map{
			"verify_token":     token,
			"id_name":          name,
			"id_no":            icNum,
			"certificate_type": 0,
		}).
		Post(fmt.Sprintf(faceprintSubmitUrl, b.accessToken))
	if err != nil {
		snag.Panic(err)
	}
	zap.L().Info("获取人脸校验URL", log.ResponseBody(r.Body()))
	return
}

// GetAuthenticatorUrl 实名认证
func (b *faceClient) GetAuthenticatorUrl() (string, string) {
	return b.getFaceprintUrl(TypeAuth)
}

// FaceResult 获取人脸照片
func (b *faceClient) FaceResult(token string) (res *faceprintFaceResp, err error) {
	res = new(faceprintFaceResp)
	var r *resty.Response
	r, err = resty.New().R().
		SetResult(res).
		SetBody(ar.Map{"verify_token": token}).
		Post(fmt.Sprintf(faceprintSimpleUrl, b.accessToken))
	if err != nil {
		return
	}
	zap.L().Info("获取人脸照片结果", log.ResponseBody(r.Body()))
	return
}

// AuthenticatorResult 获取实名认证结果
func (b *faceClient) AuthenticatorResult(token string) (res *faceprintDetailResp, err error) {
	var simple *faceprintFaceResp
	simple, err = b.FaceResult(token)
	if err != nil {
		return
	}

	res = new(faceprintDetailResp)
	var r *resty.Response
	r, err = resty.New().R().
		SetResult(res).
		SetBody(map[string]string{"verify_token": token}).
		Post(fmt.Sprintf(faceprintResultUrl, b.accessToken))
	if err != nil {
		return
	}
	zap.L().Info("获取实名认证结果", log.ResponseBody(r.Body()))
	res.Result.FaceImg = simple.Result.Image
	return
}
