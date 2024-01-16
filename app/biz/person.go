// Created at 2024-01-11

package biz

import (
	cloudauth "github.com/alibabacloud-go/cloudauth-20190307/v3/client"
	jsoniter "github.com/json-iterator/go"

	"github.com/auroraride/aurservd/app/biz/definition"
	"github.com/auroraride/aurservd/internal/ali"
	"github.com/auroraride/aurservd/internal/ent"
)

type personBiz struct {
	orm *ent.PersonClient
}

func NewPerson() *personBiz {
	return &personBiz{
		orm: ent.Database.Person,
	}
}

// Certification 发起实名认证
func (s *personBiz) Certification(req *definition.PersonCertificationReq) (*definition.PersonCertification, error) {
	client, err := ali.NewFaceVerify()
	if err != nil {
		return nil, err
	}
	var id string
	id, err = client.RequestCertifyId(ali.RequestCertifyIdParams{
		Name:         req.Name,
		IDCardNumber: req.IDCardNumber,
		MetaInfo:     req.MetaInfo,
	})
	if err != nil {
		return nil, err
	}

	return &definition.PersonCertification{CertifyId: id}, nil
}

// CertificationResult 获取实名认证结果
func (s *personBiz) CertificationResult(req *definition.PersonCertification) (*definition.PersonCertificationResultRes, error) {
	client, err := ali.NewFaceVerify()
	if err != nil {
		return nil, err
	}

	var data *cloudauth.DescribeFaceVerifyResponseBodyResultObject
	data, err = client.Describe(req.CertifyId)
	if err != nil {
		return nil, err
	}
	var passed bool
	if data.MaterialInfo != nil && data.Passed != nil {
		passed = *data.Passed == "T"

		// 解析返回
		result := new(ali.FaceVerifyResult)
		err = jsoniter.Unmarshal([]byte(*data.MaterialInfo), result)
		if err != nil {
			return nil, err
		}

		// 存储实名信息
		// 面容
		// result.FacialPictureFront.OssObjectName

		// PRO中 OcrIdCardInfo / OcrPictureFront 均为空

		// faceVerifyResult := model.FaceVerifyResult{
		// 	Birthday:       "",
		// 	IssueAuthority: "",
		// 	Address:        "",
		// 	Gender:         "",
		// 	Nation:         "",
		// 	ExpireTime:     "",
		// 	Name:           "",
		// 	IssueTime:      "",
		// 	IdCardNumber:   "",
		// 	Score:          0,
		// 	LivenessScore:  0,
		// 	Spoofing:       0,
		// }
	}

	return &definition.PersonCertificationResultRes{Passed: passed}, nil
}
