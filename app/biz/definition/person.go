// Copyright (C) liasica. 2024-present.
//
// Created at 2024-01-11
// Based on aurservd by liasica, magicrolan@qq.com.

package definition

import (
	"encoding/base64"
	"encoding/binary"
	"errors"
	"time"

	"github.com/auroraride/aurservd/pkg/utils"
)

// PersonIdentity 用户身份信息
// 加密方法：字符串[10位时间戳大端序存储uint64][18位身份证][姓名]转为字节组后遍历，当前位(字节-序号)^序号，获得结果后进行base64转码为字符串
type PersonIdentity struct {
	IDCardNumber string `json:"idCardNumber"`
	Name         string `json:"name"`
	Time         int64  `json:"-"`
}

func NewPersonIdentity(idCardNumber, name string) *PersonIdentity {
	return &PersonIdentity{
		IDCardNumber: idCardNumber,
		Name:         name,
		Time:         time.Now().Unix(),
	}
}

// Pack 数据打包
func (p *PersonIdentity) Pack() string {
	src := make([]byte, 8)
	binary.BigEndian.PutUint64(src, uint64(p.Time))
	src = append(src, []byte(p.IDCardNumber+p.Name)...)
	for i := 0; i < len(src); i++ {
		src[i] = byte((int(src[i]) - i) ^ i)
	}
	return base64.StdEncoding.EncodeToString(src)
}

// UnPack 数据解包
func (p *PersonIdentity) UnPack(str string) error {
	decoded, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return err
	}

	if len(decoded) < 30 {
		return errors.New("数据格式错误")
	}

	for i := 0; i < len(decoded); i++ {
		decoded[i] = byte(int(decoded[i]) ^ i + i)
	}

	p.Time = int64(binary.BigEndian.Uint64(decoded[:8]))
	p.IDCardNumber = string(decoded[8:26])
	if !utils.NewRegex().MatchIDCardNumber(p.IDCardNumber) {
		return errors.New("身份证号校验失败")
	}

	p.Name = string(decoded[26:])

	return nil
}

// PersonCertificationOcrRes 腾讯人身核验OCR参数
type PersonCertificationOcrRes struct {
	AppID   string `json:"appId"`   // WBAppid
	UserId  string `json:"userId"`  // 用户唯一标识
	OrderNo string `json:"orderNo"` // 订单号
	Version string `json:"version"` // 版本号
	Nonce   string `json:"nonce"`   // 随机字符串
	Sign    string `json:"sign"`    // 签名
}

// PersonCertificationFaceReq 实名认证人脸核身请求参数
type PersonCertificationFaceReq struct {
	OrderNo  string `json:"orderNo" query:"orderNo"`                       // 订单号，用户使用OCR识别时不为空
	Identity string `json:"identity" query:"identity" validate:"required"` // 身份信息
}

// PersonCertificationFaceRes 实名认证人脸核身参数
type PersonCertificationFaceRes struct {
	PersonCertificationOcrRes
	FaceId  string `json:"faceId"`
	Licence string `json:"licence"`
}

// PersonCertificationFaceResultReq 实名认证人脸核身结果请求参数
type PersonCertificationFaceResultReq struct {
	OrderNo string `json:"orderNo" query:"orderNo" validate:"required"` // 订单号
}

// PersonCertificationFaceResultRes 实名认证人脸核身结果
type PersonCertificationFaceResultRes struct {
	Success bool `json:"success"` // 是否成功
}
