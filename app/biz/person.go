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
func (p *personBiz) CertificationOcr(r *ent.Rider) (res *definition.PersonCertificationOcrRes, err error) {
	userId := strconv.FormatUint(r.ID, 10)

	var ticket string
	w := tencent.NewWbFace()
	ticket, err = w.NonceTicket(userId)
	if err != nil {
		return
	}

	return &definition.PersonCertificationOcrRes{
		AppID:   w.AppId(),
		UserId:  userId,
		OrderNo: tools.NewUnique().NewSN28(),
		Ticket:  ticket,
	}, nil
}
