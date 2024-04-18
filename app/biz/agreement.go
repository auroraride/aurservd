package biz

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"strings"

	"github.com/golang-module/carbon/v2"

	"github.com/auroraride/aurservd/app/biz/definition"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/assets"
	"github.com/auroraride/aurservd/internal/ali"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/agreement"
	"github.com/auroraride/aurservd/internal/ent/enterpriseprice"
	"github.com/auroraride/aurservd/internal/ent/plan"
)

type agreementBiz struct {
	orm      *ent.AgreementClient
	ctx      context.Context
	modifier *model.Modifier
}

func NewAgreement() *agreementBiz {
	return &agreementBiz{
		orm: ent.Database.Agreement,
		ctx: context.Background(),
	}
}

func NewAgreementWithModifierBiz(m *model.Modifier) *agreementBiz {
	s := NewAgreement()
	if m != nil {
		s.ctx = context.WithValue(s.ctx, model.CtxModifierKey{}, m)
		s.modifier = m
	}
	return s
}

// Create 创建协议
func (a *agreementBiz) Create(req *definition.AgreementCreateReq) error {
	if req.IsDefault != nil && *req.IsDefault {
		// 查询是否有默认协议 如果有则修改为非默认
		ag, _ := a.orm.QueryNotDeleted().Where(
			agreement.IsDefault(true),
			agreement.UserType(req.UserType),
		).First(a.ctx)
		if ag != nil {
			// 修改为非默认
			_, err := a.orm.UpdateOneID(ag.ID).SetIsDefault(false).Save(a.ctx)
			if err != nil {
				return err
			}
		}
	}

	url, hash := a.GenerateAgreementFile(req.Name, req.Content)
	_, err := a.orm.Create().
		SetName(req.Name).
		SetContent(req.Content).
		SetUserType(req.UserType).
		SetForceReadTime(req.ForceReadTime).
		SetNillableIsDefault(req.IsDefault).
		SetURL(url).
		SetHash(hash).
		Save(a.ctx)
	if err != nil {
		return err
	}

	return nil
}

// Modify 修改协议
func (a *agreementBiz) Modify(req *definition.AgreementModifyReq) error {
	if req.IsDefault != nil && *req.IsDefault {
		// 查询是否有默认协议 如果有则修改为非默认
		ag, _ := a.orm.QueryNotDeleted().Where(
			agreement.IsDefault(true),
			agreement.UserType(req.UserType),
			agreement.Or(
				agreement.IDNEQ(req.ID),
			),
		).First(a.ctx)
		if ag != nil {
			// 修改为非默认
			_, err := a.orm.UpdateOneID(ag.ID).SetIsDefault(false).Save(a.ctx)
			if err != nil {
				return err
			}
		}
	}

	q := a.orm.UpdateOneID(req.ID).
		SetName(req.Name).
		SetContent(req.Content).
		SetUserType(req.UserType).
		SetForceReadTime(req.ForceReadTime).
		SetNillableIsDefault(req.IsDefault)

	if req.Content != "" || req.Name != "" {
		url, hash := a.GenerateAgreementFile(req.Name, req.Content)
		q.
			SetURL(url).
			SetHash(hash)
	}

	_, err := q.Save(a.ctx)
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}
	return nil
}

// List 协议列表
func (a *agreementBiz) List() (res []*definition.AgreementDetail, err error) {
	res = make([]*definition.AgreementDetail, 0)
	items, _ := a.orm.QueryNotDeleted().All(a.ctx)
	if len(items) == 0 {
		return res, nil
	}

	for _, v := range items {
		res = append(res, &definition.AgreementDetail{
			ID:            v.ID,
			Name:          v.Name,
			Content:       v.Content,
			UserType:      v.UserType,
			ForceReadTime: v.ForceReadTime,
			IsDefault:     &v.IsDefault,
			CreatedAt:     v.CreatedAt.Format(carbon.DateTimeLayout),
			URL:           v.URL,
			Hash:          v.Hash,
		})
	}
	return
}

// Delete 删除协议
func (a *agreementBiz) Delete(id uint64) (err error) {
	first, _ := a.orm.QueryNotDeleted().Where(agreement.ID(id)).First(a.ctx)
	if first == nil {
		return errors.New("协议不存在")
	}
	// 判定是否有正在使用的协议
	var exist bool
	if first.UserType == definition.AgreementUserTypePersonal.Value() {
		exist, err = ent.Database.Plan.QueryNotDeleted().Where(plan.AgreementID(id)).Exist(a.ctx)
	} else {
		exist, err = ent.Database.EnterprisePrice.QueryNotDeleted().Where(enterpriseprice.AgreementID(id)).Exist(a.ctx)
	}

	if exist {
		return errors.New("该协议正在使用中，无法删除")
	}
	_, err = a.orm.SoftDeleteOneID(id).Save(a.ctx)
	if err != nil {
		return err
	}
	return nil
}

// AgreementSelection 协议选择
func (a *agreementBiz) AgreementSelection(req *definition.AgreementSearchReq) (res []definition.AgreementSelectionRes) {
	res = make([]definition.AgreementSelectionRes, 0)
	q := ent.Database.Agreement.QueryNotDeleted().Order(ent.Desc(agreement.FieldCreatedAt))
	if req.UserType != nil {
		q.Where(agreement.UserType(*req.UserType))
	}
	items, _ := q.All(a.ctx)
	if len(items) == 0 {
		return res
	}

	for _, v := range items {
		res = append(res, definition.AgreementSelectionRes{
			ID:        v.ID,
			Name:      v.Name,
			IsDefault: v.IsDefault,
			Hash:      v.Hash,
		})
	}
	return
}

// Detail 协议详情
func (a *agreementBiz) Detail(id uint64) (res *definition.AgreementDetail, err error) {
	item, _ := a.orm.QueryNotDeleted().Where(agreement.ID(id)).First(a.ctx)
	if item == nil {
		return nil, err
	}
	return &definition.AgreementDetail{
		ID:            item.ID,
		Name:          item.Name,
		Content:       item.Content,
		UserType:      item.UserType,
		ForceReadTime: item.ForceReadTime,
		IsDefault:     &item.IsDefault,
		CreatedAt:     item.CreatedAt.Format(carbon.DateTimeLayout),
		URL:           item.URL,
		Hash:          item.Hash,
	}, nil
}

// GenerateAgreementFile 生成协议文件
func (a *agreementBiz) GenerateAgreementFile(title string, content string) (url string, hashStr string) {
	prefix := "agreement/"

	data := strings.Replace(assets.LegalTemplate, "{{- .Title -}}", title, 1)
	data = strings.Replace(data, "{{- .Content -}}", content, 1)

	dataByte := []byte(data)

	hash := md5.Sum([]byte(data))
	hashStr = hex.EncodeToString(hash[:])

	name := prefix + hashStr + ".html"
	ali.NewOss().UploadBytes(name, dataByte)
	return name, hashStr
}

// QueryAgreementByEnterprisePriceID 通过价格ID获取协议信息
func (a *agreementBiz) QueryAgreementByEnterprisePriceID(id uint64) *definition.AgreementDetail {
	res := new(definition.AgreementDetail)
	item, _ := ent.Database.EnterprisePrice.QueryNotDeleted().Where(enterpriseprice.ID(id)).WithAgreement().First(a.ctx)
	if item == nil {
		return res
	}

	// 查询个签默认协议
	var defaultAgreement *ent.Agreement
	defaultAgreement, _ = ent.Database.Agreement.QueryNotDeleted().
		Where(
			agreement.UserType(model.AgreementUserTypePersonal.Value()),
			agreement.IsDefault(true),
		).First(a.ctx)

	if item.Edges.Agreement != nil {
		res.ID = item.Edges.Agreement.ID
		res.Name = item.Edges.Agreement.Name
		res.URL = item.Edges.Agreement.URL
		res.Hash = item.Edges.Agreement.Hash
		res.ForceReadTime = item.Edges.Agreement.ForceReadTime
	} else if defaultAgreement != nil {
		res.ID = defaultAgreement.ID
		res.Name = defaultAgreement.Name
		res.URL = defaultAgreement.URL
		res.Hash = defaultAgreement.Hash
		res.ForceReadTime = defaultAgreement.ForceReadTime
	}
	return res
}
