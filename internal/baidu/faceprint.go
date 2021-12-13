// Copyright (C) liasica. 2021-present.
//
// Created at 2021/12/13
// Based on aurservd by liasica, magicrolan@qq.com.

package baidu

import (
    "fmt"
    "github.com/auroraride/aurservd/app/response"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/auroraride/aurservd/pkg/utils"
    "github.com/go-resty/resty/v2"
)

const (
    // faceprintTokenUrl 获取人脸识别 verify_token 接口URL
    faceprintTokenUrl = `https://aip.baidubce.com/rpc/2.0/brain/solution/faceprint/verifyToken/generate?access_token=%s`

    // faceprintAuthUrl 获取人脸识别认证 URL
    faceprintAuthUrl = `https://brain.baidu.com/face/print/?token=%s&successUrl=%s&failedUrl=%s`

    // faceprintSimpleUrl 获取人脸图片
    faceprintSimpleUrl = `https://aip.baidubce.com/rpc/2.0/brain/solution/faceprint/result/simple?access_token=%s`

    // faceprintResultUrl 获取人脸识别认证结果详细信息
    faceprintResultUrl = `https://aip.baidubce.com/rpc/2.0/brain/solution/faceprint/result/detail?access_token=%s`
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

// Faceprint 人脸识别
func (b *baiduClient) Faceprint() (url string) {
    var err error
    cfg := ar.Config.Baidu.Face
    // 获取 verify_token
    res := new(faceprintTokenResp)
    _, err = resty.New().R().
        SetResult(res).
        SetBody(map[string]string{"plan_id": cfg.PlanId}).
        Post(fmt.Sprintf(faceprintTokenUrl, b.accessToken))
    if err != nil {
        panic(response.NewError(err))
    }
    if !res.Success {
        panic(response.NewError("人脸识别请求失败"))
    }
    token := res.Result.VerifyToken
    url = fmt.Sprintf(
        faceprintAuthUrl,
        token,
        utils.EncodeURIComponent(cfg.SuccessUrl+"?token="+token),
        utils.EncodeURIComponent(cfg.FailedUrl+"?token="+token),
    )
    return
}

// faceprintFace 获取人脸照片
func (b *baiduClient) faceprintFace(token string) (res *faceprintFaceResp, err error) {
    res = new(faceprintFaceResp)
    _, err = resty.New().R().
        SetResult(res).
        SetBody(map[string]string{"verify_token": token}).
        Post(fmt.Sprintf(faceprintSimpleUrl, b.accessToken))
    return
}

// FaceprintResult 获取人脸识别验证结果
func (b *baiduClient) FaceprintResult(token string) (res *faceprintDetailResp, err error) {
    simple := new(faceprintFaceResp)
    res = new(faceprintDetailResp)
    simple, err = b.faceprintFace(token)
    if err != nil {
        return
    }
    if !simple.Success {
        return
    }
    _, err = resty.New().R().
        SetResult(res).
        SetBody(map[string]string{"verify_token": token}).
        Post(fmt.Sprintf(faceprintResultUrl, b.accessToken))
    if err != nil {
        return
    }
    res.Result.FaceImg = simple.Result.Image
    return
}
