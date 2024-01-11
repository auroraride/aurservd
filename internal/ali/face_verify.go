// Copyright (C) liasica. 2024-present.
//
// Created at 2024-01-10
// Based on aurservd by liasica, magicrolan@qq.com.

package ali

import (
	"errors"

	cloudauth "github.com/alibabacloud-go/cloudauth-20190307/v3/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/credentials-go/credentials"
	"github.com/lithammer/shortuuid/v4"

	"github.com/auroraride/aurservd/internal/ar"
)

const (
	faceVerifySuccess = "success"
)

var faceVerifyCodes = map[string]string{
	"200": faceVerifySuccess,
	"201": "姓名和身份证不一致",
	"202": "查询不到身份信息",
	"203": "查询不到照片或照片不可用",
	"204": "人脸比对不一致",
	"205": "活体检测存在风险",
	"206": "业务策略限制",
	"207": "人脸与身份证人脸比对不一致",
	"209": "权威比对源异常",
	"400": "参数不能为空",
	"401": "参数非法",
	"402": "应用配置不存在",
	"403": "认证未完成",
	"404": "认证场景配置不存在",
	"406": "无效的CertifyId",
	"410": "未开通服务",
	"411": "RAM无权限",
	"412": "欠费中",
	"414": "设备类型不支持",
	"415": "SDK版本不支持",
	"416": "系统版本不支持",
	"417": "无法使用刷脸服务",
	"418": "刷脸失败次数过多",
	"424": "身份认证记录不存在",
	"427": "业务策略限制",
	"500": "系统错误",
}

type FaceVerify struct {
	SceneId     *int64
	CallbackUrl *string

	*cloudauth.Client
}

func NewFaceVerify() (ca *FaceVerify, err error) {
	cfg := ar.Config.Aliyun.FaceVerify

	var credential credentials.Credential
	credential, err = credentials.NewCredential(&credentials.Config{
		Type:            tea.String("access_key"),
		AccessKeyId:     tea.String(cfg.AccessKeyId),
		AccessKeySecret: tea.String(cfg.AccessKeySecret),
	})
	if err != nil {
		return
	}

	config := &openapi.Config{
		Credential: credential,
		Endpoint:   tea.String(cfg.Endpoint),
	}

	var client *cloudauth.Client
	client, err = cloudauth.NewClient(config)
	if err != nil {
		return
	}

	return &FaceVerify{
		SceneId:     tea.Int64(cfg.SceneId),
		Client:      client,
		CallbackUrl: tea.String(cfg.Callback),
	}, nil
}

type RequestCertifyIdParams struct {
	Name         *string // 姓名
	IDCardNumber *string // 身份证号
	MetaInfo     *string // 设备信息
}

func (c *FaceVerify) RequestCertifyId(params RequestCertifyIdParams) (string, error) {
	request := &cloudauth.InitFaceVerifyRequest{
		SceneId:       c.SceneId,
		ProductCode:   tea.String("ID_PLUS"),
		OuterOrderNo:  tea.String(shortuuid.NewWithNamespace("cloudauth")),
		Model:         tea.String("PHOTINUS_LIVENESS"), // 眨眼动作活体+炫彩活体双重检测
		CallbackUrl:   c.CallbackUrl,
		CallbackToken: tea.String(shortuuid.NewWithNamespace("cloudauth/callback")),
	}

	// 判定请求类型
	if params.Name != nil && params.IDCardNumber != nil {
		request.Model = tea.String("ID_PRO")
		request.CertType = tea.String("IDENTITY_CARD")
		request.CertName = params.Name
		request.CertNo = params.IDCardNumber
	}

	response, err := c.InitFaceVerify(request)
	if err != nil {
		return "", err
	}

	res := response.Body
	if res == nil || res.Code == nil {
		return "", errors.New("实名认证请求失败")
	}

	code := *res.Code
	message := faceVerifyCodes[code]
	if message == "" {
		message = "实名认证请求失败：" + code
	}
	if message != faceVerifySuccess {
		return "", errors.New(message)
	}

	// 返回结果
	return *res.ResultObject.CertifyId, nil
}

// FaceVerifyResult 实名认证结果
// https://help.aliyun.com/zh/id-verification/financial-grade-id-verification/describefaceverify?spm=a2c4g.11186623.0.0.3514251flCmXAg#table-e6f-aq7-04a
type FaceVerifyResult struct {
	FaceAge            int    `json:"faceAge,omitempty"`
	FaceAttack         string `json:"faceAttack,omitempty"`
	FaceOcclusion      string `json:"faceOcclusion,omitempty"`
	FacialPictureFront struct {
		FaceAttackScore   float64 `json:"faceAttackScore,omitempty"`
		Gender            string  `json:"gender,omitempty"`
		IdCardVerifyScore float64 `json:"idCardVerifyScore,omitempty"`
		OssBucketName     string  `json:"ossBucketName,omitempty"`
		OssObjectName     string  `json:"ossObjectName,omitempty"`
		PictureUrl        string  `json:"pictureUrl,omitempty"`
		QualityScore      float64 `json:"qualityScore,omitempty"`
		VerifyScore       float64 `json:"verifyScore,omitempty"`
	} `json:"facialPictureFront,omitempty"`
	OcrIdCardInfo struct {
		Address     string `json:"address,omitempty"`
		Authority   string `json:"authority,omitempty"`
		Birth       string `json:"birth,omitempty"`
		CertName    string `json:"certName,omitempty"`
		CertNo      string `json:"certNo,omitempty"`
		EndDate     string `json:"endDate,omitempty"`
		Nationality string `json:"nationality,omitempty"`
		Sex         string `json:"sex,omitempty"`
		StartDate   string `json:"startDate,omitempty"`
	} `json:"ocrIdCardInfo,omitempty"`
	OcrPictureFront struct {
		OssBucketName                 string `json:"ossBucketName,omitempty"`
		OssIdFaceObjectName           string `json:"ossIdFaceObjectName,omitempty"`
		OssIdFaceUrl                  string `json:"ossIdFaceUrl,omitempty"`
		OssIdNationalEmblemObjectName string `json:"ossIdNationalEmblemObjectName,omitempty"`
		OssIdNationalEmblemUrl        string `json:"ossIdNationalEmblemUrl,omitempty"`
	} `json:"ocrPictureFront,omitempty"`
}

// Describe 获取验证结果
func (c *FaceVerify) Describe(certifyId string) (*cloudauth.DescribeFaceVerifyResponseBodyResultObject, error) {
	response, err := c.DescribeFaceVerify(&cloudauth.DescribeFaceVerifyRequest{
		CertifyId: tea.String(certifyId),
		SceneId:   c.SceneId,
	})

	if err != nil {
		return nil, err
	}

	res := response.Body
	if res == nil || res.Code == nil {
		return nil, errors.New("实名认证结果请求失败")
	}

	err = c.describeCode(res.Code)
	if err != nil {
		return nil, err
	}

	err = c.describeCode(res.ResultObject.SubCode)
	if err != nil {
		return nil, err
	}

	return res.ResultObject, nil
}

func (c *FaceVerify) describeCode(code *string) error {
	if code == nil {
		return errors.New("实名认证结果解析失败")
	}

	message := faceVerifyCodes[*code]
	if message == "" {
		message = "实名认证结果请求失败：" + *code
	}
	if message != faceVerifySuccess {
		return errors.New(message)
	}

	return nil
}
