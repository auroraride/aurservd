// Copyright (C) liasica. 2021-present.
//
// Created at 2021/12/13
// Based on aurservd by liasica, magicrolan@qq.com.

package model

const (
	PersonUnauthenticated      PersonAuthStatus = iota // 未认证
	PersonAuthPending                                  // 认证中
	PersonAuthenticated                                // 已认证
	PersonAuthenticationFailed                         // 认证失败
)

type PersonAuthStatus uint8

func (s PersonAuthStatus) String() string {
	switch s {
	case PersonUnauthenticated:
		return "未认证"
	case PersonAuthPending:
		return "认证中"
	case PersonAuthenticated:
		return "已认证"
	default:
		return "认证失败"
	}
}

func (s PersonAuthStatus) Value() uint8 {
	return uint8(s)
}

type BaiduFaceVerifyResult struct {
	Birthday       string  `json:"birthday"`       // 生日
	IssueAuthority string  `json:"issueAuthority"` // 签发机关
	Address        string  `json:"address"`        // 地址
	Gender         string  `json:"gender"`         // 性别
	Nation         string  `json:"nation"`         // 民族
	ExpireTime     string  `json:"expireTime"`     // 身份证失效日期
	Name           string  `json:"name"`           // 姓名
	IssueTime      string  `json:"issueTime"`      // 身份证生效日期
	IdCardNumber   string  `json:"idCardNumber"`   // 身份证号
	Score          float64 `json:"score"`          // 人脸实名认证得分
	LivenessScore  float64 `json:"livenessScore"`  // 活体检测分数：在线图片活体：活体验证通过时返回活体分数，不通过则返回0。数字/动作/视频活体：活体通过/不通过均会返回活体分数
	Spoofing       float64 `json:"spoofing"`       // 合成图分数：若未进行合成图检测，则返回0；若进行活体检测，则返回合成图检测分值
}

type PersonFaceVerifyResult struct {
	Name            string  `json:"name"`            // 姓名
	Sex             string  `json:"sex"`             // 性别
	Nation          string  `json:"nation"`          // 民族
	Birth           string  `json:"birth"`           // 出生日期（例：19920320）
	Address         string  `json:"address"`         // 地址
	Idcard          string  `json:"idcard"`          // 身份证号
	ValidStartDate  string  `json:"validStartDate"`  // 证件生效日期
	ValidExpireDate string  `json:"validExpireDate"` // 证件失效日期
	Authority       string  `json:"authority"`       // 发证机关
	FrontClarity    string  `json:"frontClarity"`    // 正面图片清晰度
	BackClarity     string  `json:"backClarity"`     // 反面图片清晰度
	LiveRate        float64 `json:"liveRate"`        // 活体检测得分
	Similarity      float64 `json:"similarity"`      // 人脸比对得分
}

// RequireAuth 是否需要认证
func (s PersonAuthStatus) RequireAuth() bool {
	return s != PersonAuthenticated
}

// PersonBanReq 封禁或解封骑手身份
type PersonBanReq struct {
	ID  uint64 `json:"id" ` // 骑手ID
	Ban bool   `json:"ban"` // `true`封禁 `false`解封
}

type Person struct {
	// 证件号码
	IDCardNumber string `json:"idCardNumber,omitempty"`
	// 证件人像面
	IDCardPortrait string `json:"idCardPortrait,omitempty"`
	// 证件国徽面
	IDCardNational string `json:"idCardNational,omitempty"`
	// 证件人像照片
	IDCardHead string `json:"idCardHead,omitempty"`
	// 实名认证人脸照片
	AuthFace string `json:"authFace,omitempty"`
}
