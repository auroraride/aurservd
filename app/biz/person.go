// Created at 2024-01-11

package biz

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/golang-module/carbon/v2"
	faceid "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/faceid/v20180301"

	"github.com/auroraride/aurservd/app/biz/definition"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/app/service"
	"github.com/auroraride/aurservd/internal/ali"
	"github.com/auroraride/aurservd/internal/ar"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/person"
	"github.com/auroraride/aurservd/internal/ent/rider"
	"github.com/auroraride/aurservd/internal/tencent"
	"github.com/auroraride/aurservd/pkg/silk"
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
	if service.NewRider().IsAuthed(r) {
		return nil, errors.New("当前已认证，无法重复认证")
	}

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

// 解析ocr识别结果
// 通过腾讯OCR订单号获取到的result中，包含订单号`OrderNo`
// 通过腾讯身份证识别及信息核验识别的身份证result中，不包含`OrderNo`
func (b *personBiz) ocrResult(creator *ent.PersonCreate, identity *definition.PersonIdentity, result *tencent.OcrResult, faceOrderNo string) {
	// 异步上传照片到阿里云OSS
	portrait, national, head := b.uploadOcrFiles(result)

	url := ar.Config.Aliyun.Oss.Url
	if !strings.HasSuffix(url, "/") {
		url += "/"
	}

	// 保存人像面
	if portrait != "" {
		creator.SetIDCardPortrait(url + portrait)
	}

	// 保存国徽面
	if national != "" {
		creator.SetIDCardNational(url + national)
	}

	// 解析有效期
	var start, expire string
	if result.ValidDate != "" {
		arr := strings.Split(result.ValidDate, "-")
		if len(arr) > 0 {
			start = arr[0]
		}
		if len(arr) > 1 {
			expire = arr[1]
		}
	}

	// 解析清晰度
	var fc, bc float64
	if result.FrontClarity != "" {
		fc, _ = strconv.ParseFloat(result.FrontClarity, 64)
	}
	if result.BackClarity != "" {
		bc, _ = strconv.ParseFloat(result.BackClarity, 64)
	}

	mfvr := &model.PersonFaceVerifyResult{
		Name:            identity.Name,
		Sex:             result.Sex,
		Nation:          result.Nation,
		Birth:           identity.IDCardNumber[7:14],
		Address:         result.Address,
		IDCardNumber:    identity.IDCardNumber,
		ValidStartDate:  start,
		ValidExpireDate: expire,
		Authority:       result.Authority,
		FrontClarity:    fc,
		BackClarity:     bc,
		FaceOrderNo:     faceOrderNo,
		OcrOrderNo:      result.OrderNo,
		Head:            head,
	}
	if head != "" {
		mfvr.Head = url + head
	}

	creator.SetFaceVerifyResult(mfvr)
}

// 异步上传ocr照片到阿里云OSS
// TODO: 阿里云OSS以前的图片移动目录
func (b *personBiz) uploadOcrFiles(result *tencent.OcrResult) (portrait, national, head string) {
	prefix := "__rider_assets/faceverify/" + result.Idcard + "/ocr-" + result.OrderNo + "-"

	if result.FrontCrop != "" {
		portrait = prefix + "portrait.jpg"
	}
	if result.BackCrop != "" {
		national = prefix + "national.jpg"
	}
	if result.HeadPhoto != "" {
		head = prefix + "head.jpg"
	}

	go func() {
		oss := ali.NewOss()
		if portrait != "" {
			oss.UploadBase64(portrait, result.FrontCrop)
		}
		if national != "" {
			oss.UploadBase64(national, result.BackCrop)
		}
		if head != "" {
			oss.UploadBase64(head, result.HeadPhoto)
		}
	}()

	return
}

// 上传人身核验图片和视频
func (b *personBiz) uploadFaceVerifyFiles(idCardNumber string, result *tencent.FaceVerifyResult) (photo, video string) {
	prefix := "__rider_assets/faceverify/" + idCardNumber + "/faceverify-" + result.OrderNo + "-"

	if result.Photo != "" {
		photo = prefix + "photo.jpg"
	}

	if result.Video != "" {
		video = prefix + "video.mp4"
	}

	go func() {
		oss := ali.NewOss()
		if photo != "" {
			oss.UploadBase64(photo, result.Photo)
		}
		if video != "" {
			oss.UploadBase64(video, result.Video)
		}
	}()

	return
}

// CertificationFace 提交身份信息并获取人脸核身参数
func (b *personBiz) CertificationFace(r *ent.Rider, req *definition.PersonCertificationFaceReq) (res *definition.PersonCertificationFaceRes, err error) {
	if service.NewRider().IsAuthed(r) {
		return nil, errors.New("当前已认证，无法重复认证")
	}

	var result *tencent.OcrResult

	identity := new(definition.PersonIdentity)

	// 通过OCR获取身份信息
	if req.PortraitImage != nil {
		var params *faceid.IdCardOCRVerificationResponseParams
		params, err = tencent.NewFaceId().IdCardOCR(*req.PortraitImage)
		if err != nil {
			return
		}

		result = &tencent.OcrResult{
			Name:      *params.Name,
			Sex:       *params.Sex,
			Nation:    *params.Nation,
			Birth:     *params.Birth,
			Address:   *params.Address,
			Idcard:    *params.IdCard,
			FrontCrop: *req.PortraitImage,
			BackCrop:  *req.NationalImage,
		}
	}

	// 通过加密字符串获取身份信息
	if req.Identity != nil {
		err = identity.UnPack(*req.Identity)
		if err != nil {
			return
		}
	}

	// 获取生日
	birth := identity.IDCardNumber[7:14]
	birthday := carbon.Parse(birth).StdTime().AddDate(18, 0, 0)

	// 未年满18岁认证标记为失败
	if birthday.After(time.Now()) {
		return nil, errors.New("未年满18周岁")
	}

	w := tencent.NewWbFace()

	// 判定今日实名认证次数
	// TODO: 后台设定次数
	times := w.GetTimes(identity.IDCardNumber)
	if times >= 5 {
		return nil, errors.New("实名次数过于频繁，请明天再试")
	}

	// 判定是否绑定其他账号
	p, _ := ent.Database.Person.
		QueryNotDeleted().
		Where(person.IDCardNumber(identity.IDCardNumber)).
		WithRiders(func(query *ent.RiderQuery) {
			query.Where(rider.DeletedAtIsNil(), rider.IDNotIn(r.ID))
		}).
		First(context.Background())
	if p != nil && len(p.Edges.Riders) > 0 {
		phone := p.Edges.Riders[0].Phone
		phone = phone[:3] + strings.Repeat("*", 5) + phone[8:]
		return &definition.PersonCertificationFaceRes{BindedPhone: phone}, nil
	}

	// 生成人脸核身订单号
	faceOrderNo := tools.NewUnique().Rand(32)
	// if req.OrderNo != "" {
	// 	faceOrderNo = req.OrderNo
	// }

	// 保存或更新实人表
	creator := ent.Database.Person.Create().
		SetStatus(model.PersonAuthPending.Value()).
		SetIDCardNumber(identity.IDCardNumber).
		SetName(identity.Name)

	if req.OrderNo != "" {
		// 查询OCR结果
		err, result = w.OcrResult(req.OrderNo)
		if err != nil {
			return
		}
		b.ocrResult(creator, identity, result, faceOrderNo)
	}

	if result != nil {
		b.ocrResult(creator, identity, result, faceOrderNo)
	} else {
		creator.SetFaceVerifyResult(&model.PersonFaceVerifyResult{
			IDCardNumber: identity.IDCardNumber,
			Name:         identity.Name,
			Birth:        birth,
			FaceOrderNo:  faceOrderNo,
		})
	}

	// 保存实人并返回ID
	var id uint64
	id, err = creator.OnConflictColumns(person.FieldIDCardNumber).
		UpdateNewValues().
		ID(context.Background())
	if err != nil {
		return
	}

	// 判断ID是否为实名认证的ID, 如果不是, 则删除
	if r.PersonID != nil && id != *r.PersonID {
		err = ent.Database.Person.DeleteOneID(*r.PersonID).Exec(context.Background())
		if err != nil {
			return
		}
	}

	// 更新骑手表
	err = r.Update().SetPersonID(id).SetIDCardNumber(identity.IDCardNumber).Exec(context.Background())
	if err != nil {
		return
	}

	// 获取人脸核身参数
	userId := strconv.FormatUint(r.ID, 10)

	var faceId, sign, nonce string
	faceId, sign, nonce, err = w.GetFaceId(&tencent.FaceIdReq{
		OrderNo: faceOrderNo,
		Name:    identity.Name,
		IdNo:    identity.IDCardNumber,
		UserId:  userId,
	})
	if err != nil {
		return
	}

	return &definition.PersonCertificationFaceRes{
		PersonCertificationOcrRes: definition.PersonCertificationOcrRes{
			AppID:   w.AppId(),
			UserId:  userId,
			OrderNo: faceOrderNo,
			Version: w.Version(),
			Nonce:   nonce,
			Sign:    sign,
		},
		FaceId:  faceId,
		Licence: w.Licence(),
	}, nil
}

// CertificationFaceResult 获取人脸核身结果
func (b *personBiz) CertificationFaceResult(r *ent.Rider, req *definition.PersonCertificationFaceResultReq) (res *definition.PersonCertificationFaceResultRes, err error) {
	rp := r.Edges.Person
	if rp == nil || rp.FaceVerifyResult == nil {
		return nil, errors.New("未找到实人信息")
	}

	mfvr := rp.FaceVerifyResult
	if mfvr.FaceOrderNo != req.OrderNo {
		return nil, errors.New("实人信息不匹配")
	}

	url := ar.Config.Aliyun.Oss.Url
	if !strings.HasSuffix(url, "/") {
		url += "/"
	}

	var (
		result   *tencent.FaceVerifyResult
		authface *string
	)

	w := tencent.NewWbFace()
	result, err = w.FaceVerifyResult(req.OrderNo)
	res = new(definition.PersonCertificationFaceResultRes)

	defer func() {
		status := model.PersonAuthenticated
		if err != nil {
			status = model.PersonAuthenticationFailed
		}
		ctx := context.Background()
		err = ent.Database.Person.UpdateOne(rp).
			SetFaceVerifyResult(mfvr).
			SetStatus(status.Value()).
			SetAuthAt(time.Now()).
			SetNillableAuthFace(authface).
			Exec(ctx)
		if err == nil {
			service.NewPromotionReferralsService().RiderBindReferrals(r)
		}
	}()

	if err != nil {
		return
	}

	mfvr.LiveRate, _ = strconv.ParseFloat(result.LiveRate, 64)
	mfvr.Similarity, _ = strconv.ParseFloat(result.Similarity, 64)

	// 上传文件到阿里云OSS
	photo, video := b.uploadFaceVerifyFiles(rp.IDCardNumber, result)

	if photo != "" {
		mfvr.Photo = url + photo
		authface = silk.String(mfvr.Photo)
	}

	if video != "" {
		mfvr.Video = url + video
	}

	res.Success = true

	return
}
