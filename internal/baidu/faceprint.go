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

// getFaceprintUrl 获取人脸核验URL
// TODO 验证页面随便跳转到那个页面, flutter前端进行webview路由拦截获取token请求验证结果, 成功或失败进行弹窗提示
func (b *baiduClient) getFaceprintUrl(typ string) (url string, token string) {
    var err error
    var planId string
    cfg := ar.Config.Baidu.Face
    switch typ {
    case TypeAuth:
        planId = cfg.AuthPlanId
    case TypeFace:
        planId = cfg.FacePlanId
    }

    // 获取 verify_token
    res := new(faceprintTokenResp)
    _, err = resty.New().R().
        SetResult(res).
        SetBody(map[string]string{"plan_id": planId}).
        Post(fmt.Sprintf(faceprintTokenUrl, b.accessToken))
    if err != nil {
        panic(response.NewError(err))
    }
    if !res.Success {
        panic(response.NewError("实名认证请求失败"))
    }
    token = res.Result.VerifyToken
    str := fmt.Sprintf("%s?type=%s&token=%s&state=", cfg.Callback, typ, token)
    url = fmt.Sprintf(
        faceprintAuthUrl,
        token,
        utils.EncodeURIComponent(str+"success"),
        utils.EncodeURIComponent(str+"failed"),
    )
    return
}

// GetFaceUrl 获取人脸校验URL
// TODO 缓存token和用户对应关系, token只能请求一次
func (b *baiduClient) GetFaceUrl(name, icNum string) string {
    uri, token := b.getFaceprintUrl(TypeFace)
    res := new(faceprintSubmitResp)
    _, err := resty.New().R().
        SetResult(res).
        SetBody(ar.Map{
            "verify_token":     token,
            "id_name":          name,
            "id_no":            icNum,
            "certificate_type": 0,
        }).
        Post(fmt.Sprintf(faceprintSubmitUrl, b.accessToken))
    if err != nil {
        panic(response.NewError(err))
    }
    return uri
}

// GetAuthenticatorUrl 实名认证
// TODO 缓存token和用户对应关系, token只能请求一次
func (b *baiduClient) GetAuthenticatorUrl() string {
    uri, _ := b.getFaceprintUrl(TypeAuth)
    return uri
}

// FaceResult 获取人脸照片
func (b *baiduClient) FaceResult(token string) (res *faceprintFaceResp, err error) {
    res = new(faceprintFaceResp)
    _, err = resty.New().R().
        SetResult(res).
        SetBody(ar.Map{"verify_token": token}).
        Post(fmt.Sprintf(faceprintSimpleUrl, b.accessToken))
    if err != nil {
        return
    }
    return
}

// AuthenticatorResult 获取实名认证结果
func (b *baiduClient) AuthenticatorResult(token string) (res *faceprintDetailResp, err error) {
    simple := new(faceprintFaceResp)
    simple, err = b.FaceResult(token)
    if err != nil {
        return
    }

    res = new(faceprintDetailResp)
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
