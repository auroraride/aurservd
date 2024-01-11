// Created at 2024-01-11

package biz

import (
	cloudauth "github.com/alibabacloud-go/cloudauth-20190307/v3/client"

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
	passed := *data.Passed == "T"
	return &definition.PersonCertificationResultRes{Passed: passed}, nil
}
