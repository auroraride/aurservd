// Created at 2024-01-11

package biz

import (
	"strconv"

	"github.com/auroraride/aurservd/app/biz/definition"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/tencent"
	"github.com/auroraride/aurservd/pkg/tools"
)

type personBiz struct {
	orm *ent.PersonClient
}

func NewPerson() *personBiz {
	return &personBiz{
		orm: ent.Database.Person,
	}
}

// CertificationOcr 获取人身核验OCR参数
func (b *personBiz) CertificationOcr(r *ent.Rider) (res *definition.PersonCertificationOcrRes, err error) {
	w := tencent.NewWbFace()

	userId := strconv.FormatUint(r.ID, 10)
	orderNo := tools.NewUnique().Rand(32)

	res = &definition.PersonCertificationOcrRes{
		AppID:   w.AppId(),
		UserId:  userId,
		OrderNo: orderNo,
		Version: w.Version(),
	}

	var ticket string
	ticket, err = w.NonceTicket(userId)
	if err != nil {
		return
	}

	res.Sign, res.Nonce = w.Sign(userId, ticket)
	return
}

// CertificationFace 获取人身核验参数
func (b *personBiz) CertificationFace(r *ent.Rider, orderNo string) (res *definition.PersonCertificationOcrRes, err error) {
	return
}
