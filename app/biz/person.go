// Created at 2024-01-11

package biz

import (
	"strconv"

	"github.com/auroraride/aurservd/app/biz/definition"
	"github.com/auroraride/aurservd/internal/ali"
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

	res.Sign, res.Nonce = w.Sign(ticket, userId)
	return
}

// CertificationFace 获取人身核验参数
func (b *personBiz) CertificationFace(r *ent.Rider, orderNo string) (res *definition.PersonCertificationOcrRes, err error) {
	// 查询OCR结果
	w := tencent.NewWbFace()
	var result *tencent.OcrResult
	err, result = w.OcrResult(orderNo)
	if err != nil {
		return
	}

	// 异步上传照片到阿里云OSS
	b.uploadFiles(result)

	// url := ar.Config.Aliyun.Oss.Url
	// if !strings.HasSuffix(url, "/") {
	// 	url += "/"
	// }

	// 保存或更新实人表
	return
}

// 异步上传照片到阿里云OSS
// TODO: 阿里云OSS以前的图片移动到faceverify目录下
func (b *personBiz) uploadFiles(result *tencent.OcrResult) (portrait, national, head string) {
	prefix := "faceverify/" + result.Idcard + "/" + result.OrderNo + "-"

	portrait = prefix + "portrait.jpg"
	national = prefix + "national.jpg"
	head = prefix + "head.jpg"

	go func() {
		oss := ali.NewOss()
		if result.FrontCrop != "" {
			oss.UploadBase64ImageJpeg(portrait, result.FrontCrop)
		}
		if result.BackCrop != "" {
			oss.UploadBase64ImageJpeg(national, result.BackCrop)
		}
		if result.HeadPhoto != "" {
			oss.UploadBase64ImageJpeg(head, result.HeadPhoto)
		}
	}()

	return
}
