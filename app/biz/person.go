// Created at 2024-01-11

package biz

import (
	"context"
	crand "crypto/rand"
	"encoding/base64"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/auroraride/adapter/rpc/pb"
	"github.com/golang-module/carbon/v2"
	"github.com/lithammer/shortuuid/v4"
	"google.golang.org/protobuf/proto"

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

const (
	personAesKeySize            = 32
	personAesIvSize             = 16
	personAesKeyIvEncryptedSize = 512
)

type personBiz struct {
	orm *ent.PersonClient
}

func NewPerson() *personBiz {
	return &personBiz{
		orm: ent.Database.Person,
	}
}

// 加密身份信息
func encryptPersonIdentity(identity *pb.PersonIdentity) string {
	src, _ := proto.Marshal(identity)
	rsa := ar.PersonRsa()

	// 生成随机AES密钥和向量
	key := make([]byte, personAesKeySize)
	iv := make([]byte, personAesIvSize)
	_, _ = crand.Read(key)
	_, _ = crand.Read(iv)

	// AES加密
	aes := tools.NewAesCrypto(iv, key)
	encrypted, _ := aes.CBCEncrypt(src)

	// 加密AES密钥和向量
	encryptedKeyIv, _ := rsa.Encrypt(append(key, iv...))

	// 合并后base64编码
	return base64.StdEncoding.EncodeToString(append(encryptedKeyIv, encrypted...))
}

// 解密身份信息
func decryptPersonIdentity(str string) (identity *pb.PersonIdentity, err error) {
	rsa := ar.PersonRsa()

	var src []byte
	src, err = base64.StdEncoding.DecodeString(str)
	if err != nil {
		return
	}

	// 获取AES密钥和向量
	ikEnctyped := src[:personAesKeyIvEncryptedSize]
	var ik []byte
	ik, err = rsa.Decrypt(ikEnctyped)
	if err != nil {
		return
	}
	if len(ik) != personAesKeySize+personAesIvSize {
		return nil, errors.New("密钥长度错误")
	}
	key, iv := ik[:personAesKeySize], ik[personAesKeySize:]

	// 解密身份信息
	aes := tools.NewAesCrypto(iv, key)
	var b []byte
	b, err = aes.CBCDecrypt(src[personAesKeyIvEncryptedSize:])
	if err != nil {
		return
	}

	identity = new(pb.PersonIdentity)
	err = proto.Unmarshal(b, identity)
	return
}

// CertificationOcrClient 获取人身核验OCR参数
func (b *personBiz) CertificationOcrClient(r *ent.Rider) (res *definition.PersonCertificationOcrClientRes, err error) {
	if service.NewRider().IsAuthed(r) {
		return nil, errors.New("当前已认证，无法重复认证")
	}

	w := tencent.NewWbFace()

	userId := strconv.FormatUint(r.ID, 10)
	orderNo := tools.NewUnique().Rand(32)

	res = &definition.PersonCertificationOcrClientRes{
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
func (b *personBiz) ocrResult(creator *ent.PersonCreate, identity *pb.PersonIdentity, ocrOrderNo, faceOrderNo string) {
	result := identity.OcrResult

	// 异步上传照片到阿里云OSS
	portrait, national := b.uploadOcrFiles(result)

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

	mfvr := &model.PersonFaceVerifyResult{
		Name:            identity.Name,
		Sex:             result.Sex,
		Nation:          result.Nation,
		Birth:           identity.IdCardNumber[7:14],
		Address:         result.Address,
		IDCardNumber:    identity.IdCardNumber,
		ValidStartDate:  result.ValidStartDate,
		ValidExpireDate: result.ValidExpireDate,
		Authority:       result.Authority,
		PortraitClarity: result.PortraitClarity,
		NationalClarity: result.NationalClarity,
		FaceOrderNo:     faceOrderNo,
		OcrOrderNo:      ocrOrderNo,
	}

	creator.SetFaceVerifyResult(mfvr)
}

// 异步上传ocr照片到阿里云OSS
// TODO: 阿里云OSS以前的图片移动目录
func (b *personBiz) uploadOcrFiles(result *pb.PersonIdentityOcrResult) (portrait, national string) {
	prefix := "__rider_assets/faceverify/" + result.IdCardNumber + "/ocr-" + shortuuid.New() + "-"

	if result.PortraitCrop != "" {
		portrait = prefix + "portrait.jpg"
	}
	if result.NationalCrop != "" {
		national = prefix + "national.jpg"
	}

	go func() {
		oss := ali.NewOss()
		if portrait != "" {
			oss.UploadBase64(portrait, result.PortraitCrop)
		}
		if national != "" {
			oss.UploadBase64(national, result.NationalCrop)
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

	var identity *pb.PersonIdentity
	// 解密身份信息
	identity, err = decryptPersonIdentity(req.Identity)
	if err != nil {
		return
	}

	// 获取生日
	birth := identity.IdCardNumber[7:14]
	birthday := carbon.Parse(birth).StdTime().AddDate(18, 0, 0)

	// 未年满18岁认证标记为失败
	if birthday.After(time.Now()) {
		return nil, errors.New("未年满18周岁")
	}

	w := tencent.NewWbFace()

	// 判定今日实名认证次数
	// TODO: 后台设定次数
	times := w.GetTimes(identity.IdCardNumber)
	if times >= 5 {
		return nil, errors.New("实名次数过于频繁，请明天再试")
	}

	// 判定是否绑定其他账号
	p, _ := ent.Database.Person.
		QueryNotDeleted().
		Where(person.IDCardNumber(identity.IdCardNumber)).
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
		SetIDCardNumber(identity.IdCardNumber).
		SetName(identity.Name)

	// 查询客户端OCR结果
	// if req.OrderNo != "" {
	// 	// 查询OCR结果
	// 	err, result = w.OcrResult(req.OrderNo)
	// 	if err != nil {
	// 		return
	// 	}
	// 	b.ocrResult(creator, identity, result, faceOrderNo)
	// }

	if identity.OcrResult != nil {
		// 判定证件有效期
		if identity.OcrResult.ValidExpireDate != "" {
			expireDate := carbon.Parse(identity.OcrResult.ValidExpireDate).StdTime()
			if expireDate.Before(time.Now()) {
				return nil, errors.New("证件已过期")
			}
		}
		b.ocrResult(creator, identity, req.OrderNo, faceOrderNo)
	} else {
		creator.SetFaceVerifyResult(&model.PersonFaceVerifyResult{
			IDCardNumber: identity.IdCardNumber,
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
	err = r.Update().
		SetPersonID(id).
		SetName(identity.Name).
		SetIDCardNumber(identity.IdCardNumber).
		Exec(context.Background())
	if err != nil {
		return
	}

	// 获取人脸核身参数
	userId := strconv.FormatUint(r.ID, 10)

	var faceId, sign, nonce string
	faceId, sign, nonce, err = w.GetFaceId(&tencent.FaceIdReq{
		OrderNo: faceOrderNo,
		Name:    identity.Name,
		IdNo:    identity.IdCardNumber,
		UserId:  userId,
	})
	if err != nil {
		return
	}

	return &definition.PersonCertificationFaceRes{
		PersonCertificationOcrClientRes: definition.PersonCertificationOcrClientRes{
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

// CertificationOcrCloud 获取阿里云Ocr签名
func (b *personBiz) CertificationOcrCloud(hash string) (res *definition.PersonCertificationOcrCloudRes, err error) {
	var params *ali.OcrParams
	params, err = ali.NewOcr().Signature(hash)
	if err != nil {
		return
	}

	return &definition.PersonCertificationOcrCloudRes{
		ContentType:   params.ContentType,
		Action:        params.Action,
		Date:          params.Date,
		Token:         params.Token,
		Nonce:         params.Nonce,
		Version:       params.Version,
		Authorization: params.Authorization,
		Url:           params.Url,
	}, nil
}
