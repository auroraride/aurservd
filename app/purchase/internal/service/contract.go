package service

import (
	"context"
	"errors"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/liasica/edocseal/pb"
	"go.uber.org/zap"

	"github.com/auroraride/aurservd/app/biz/definition"
	"github.com/auroraride/aurservd/app/model"
	mp "github.com/auroraride/aurservd/app/purchase/internal/model"
	"github.com/auroraride/aurservd/app/rpc"
	"github.com/auroraride/aurservd/app/service"
	"github.com/auroraride/aurservd/internal/ar"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/purchaseorder"
	"github.com/auroraride/aurservd/internal/ent/purchasepayment"
	"github.com/auroraride/aurservd/pkg/tools"
	"github.com/auroraride/aurservd/pkg/utils"
)

type contractService struct {
	orm *ent.PurchaseOrderClient
}

func NewContract() *contractService {
	return &contractService{
		orm: ent.Database.PurchaseOrder,
	}
}

// Sign 签约
func (s *contractService) Sign(ctx context.Context, r *ent.Rider, req *mp.ContractSignNewReq) (res *definition.ContractSignNewRes, err error) {
	// 查找订单
	o, _ := ent.Database.PurchaseOrder.QueryNotDeleted().
		Where(
			purchaseorder.ID(req.OrderID),
			purchaseorder.StatusEQ(purchaseorder.StatusPending),
		).WithStore(func(query *ent.StoreQuery) {
		query.WithCity(func(query *ent.CityQuery) {
			query.WithParent()
		})
	}).First(ctx)
	if o == nil {
		return nil, errors.New("未找到订单信息")
	}

	person, _ := r.QueryPerson().First(ctx)
	if person == nil {
		return nil, errors.New("未找到骑手信息")
	}

	var city, province string
	if o.Edges.Store != nil && o.Edges.Store.Edges.City != nil {
		city = o.Edges.Store.Edges.City.Name
		if o.Edges.Store.Edges.City.Edges.Parent != nil {
			province = o.Edges.Store.Edges.City.Edges.Parent.Name
		}
	} else {
		return nil, errors.New("未找到城市信息")
	}

	// 获取模版id
	cfg := ar.Config.Contract
	// 请求签署合同
	url, err := rpc.Sgin(ctx, &pb.ContractSignRequest{
		DocId:    req.DocId,
		Image:    req.Seal,
		Name:     person.Name,
		Province: province,
		City:     city,
		Address:  person.FaceVerifyResult.Address,
		Phone:    r.Phone,
		Idcard:   person.IDCardNumber,
	}, cfg.Address)
	if err != nil {
		zap.L().Error("签署合同失败", zap.Error(err))
		return nil, err
	}

	if url == "" {
		zap.L().Error("签署合同失败", zap.String("url", url))
		return nil, errors.New("签署合同失败")
	}

	now := time.Now()
	billingDates := o.InstallmentPlan.BillingDates(now)

	//  更新合同状态
	err = o.Update().
		SetContractURL(url).
		Exec(ctx)
	if err != nil {
		zap.L().Error("更新合同状态失败", zap.Error(err))
		return nil, err
	}
	if strings.HasPrefix(url, "https://c.auroraride.com/") {
		hash := strings.TrimPrefix(url, "https://c.auroraride.com/")
		// 发送短信
		_, err = service.NewSms().SendSignSuccess(now, "时光驹电动车购买合同", hash, r.Phone)
		if err != nil {
			return nil, err
		}
	}
	// 更新订单
	_ = o.Update().SetNextDate(now).Exec(ctx)
	// 更新分期计划开始时间
	for k, v := range billingDates {
		_ = ent.Database.PurchasePayment.Update().
			Where(
				purchasepayment.OrderID(o.ID),
				purchasepayment.Index(k),
			).
			SetBillingDate(v).
			Exec(ctx)
	}
	return &definition.ContractSignNewRes{
		Link: url,
	}, nil
}

// Create 添加合同
func (s *contractService) Create(ctx context.Context, r *ent.Rider, req *mp.ContractCreateReq) (*definition.ContractCreateRes, error) {
	var link, docId string
	o, _ := s.orm.QueryNotDeleted().WithStore(func(query *ent.StoreQuery) {
		query.WithCity(func(query *ent.CityQuery) {
			query.WithParent()
		})
	}).Where(purchaseorder.ID(req.OrderID), purchaseorder.ContractURLIsNil()).First(ctx)

	if o == nil {
		return nil, errors.New("未找到订单信息")
	}

	if o.Edges.Store == nil || o.Edges.Store.Edges.City == nil {
		return nil, errors.New("未找到城市信息")
	}

	ec := o.Edges.Store.Edges.City

	// 定义变量
	var (
		m            = make(ar.Map)
		p            = service.NewPerson().GetNormalAuthedPerson(r)
		isEnterprise = r.EnterpriseID != nil
		err          error
		// 当前日期
		now = time.Now().Format("2006年01月02日")
	)
	// 判断是否需要补充身份信息
	if p.FaceVerifyResult == nil || p.FaceVerifyResult != nil && p.FaceVerifyResult.Address == "" || p.Name == "" || p.IDCardNumber == "" || r.Phone == "" {
		return &definition.ContractCreateRes{NeedRealName: true}, nil
	}

	sn := tools.NewUnique().NewSN()
	// 填充公共变量
	// 合同编号
	m["sn"] = sn
	// 骑手姓名
	m["name"] = p.Name
	// 身份证号
	m["idcard"] = p.IDCardNumber
	// 户口地址
	m["address"] = p.FaceVerifyResult.Address
	// 骑手电话
	m["phone"] = r.Phone
	// 限制城市
	m["city"] = ec.Name
	// 骑手签字
	m["riderSign"] = p.Name
	// 紧急联系人
	m["riderContact"] = r.Contact.String()
	// 企业签署日期
	m["aurDate"] = now
	// 骑手签署日期
	m["riderDate"] = now

	var un *model.ContractSignUniversal

	if un == nil {
		return nil, errors.New("合同信息错误")
	}

	m["payMonth"] = un.Month

	// 单电方案
	m["schemaBattery"] = true
	// 电池规格
	// 单电方案起租日
	m["batteryStart"] = now
	// 单电方案结束日
	m["batteryStop"] = un.Stop
	// 单电月租金
	m["batteryPrice"] = un.Price
	// 单电方案首次缴纳月数
	m["batteryPayMonth"] = un.Month
	// 单电方案首次缴纳租金
	m["batteryPayTotal"] = un.Total

	// 获取模版id
	cfg := ar.Config.Contract
	templateID := cfg.Template.Personal
	if isEnterprise {
		templateID = cfg.Template.Enterprise
	}

	values := make(map[string]*pb.ContractFromField)

	for k, v := range m {
		switch value := v.(type) {
		case bool:
			values[k] = &pb.ContractFromField{
				Value: &pb.ContractFromField_Checkbox{Checkbox: value},
			}
		case string:
			values[k] = &pb.ContractFromField{
				Value: &pb.ContractFromField_Text{Text: value},
			}
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
			var str string
			str, err = numberToString(value)
			if err != nil {
				return nil, err
			}
			values[k] = &pb.ContractFromField{
				Value: &pb.ContractFromField_Text{Text: str},
			}
		}
	}

	expiresAt := time.Now().Add(model.ContractExpiration * time.Minute)
	contractCreateResponse, err := rpc.Create(context.Background(), values, &definition.ContractCreateRPCReq{
		TemplateId: templateID,
		Addr:       cfg.Address,
		ExpiresAt:  expiresAt.Unix(),
		Idcard:     r.IDCardNumber,
	})
	if err != nil {
		return nil, err
	}

	link = contractCreateResponse.Url
	docId = contractCreateResponse.DocId

	_ = o.Update().SetContractURL(link).SetDocID(docId).Exec(ctx)

	if err != nil {
		return nil, err
	}

	return &definition.ContractCreateRes{
		Link:  link,
		DocId: docId,
	}, nil
}

func numberToString(num interface{}) (string, error) {
	switch v := num.(type) {
	case int, int8, int16, int32, int64:
		return strconv.FormatInt(reflect.ValueOf(num).Int(), 10), nil
	case uint, uint8, uint16, uint32, uint64:
		return strconv.FormatUint(reflect.ValueOf(num).Uint(), 10), nil
	case float32:
		return strconv.FormatFloat(float64(v), 'f', 2, 32), nil
	case float64:
		return strconv.FormatFloat(v, 'f', 2, 64), nil
	}
	return "", errors.New("类型错误")
}

// Detail 查看合同
func (s *contractService) Detail(ctx context.Context, r *ent.Rider, req *definition.ContractDetailReq) (*definition.ContractDetailRes, error) {
	cfg := ar.Config.Contract.EncryptKey
	docId, err := utils.DecryptAES([]byte(cfg), req.DocId)
	if err != nil || docId == "" {
		zap.L().Error("解密失败", zap.Error(err), zap.String("docId", req.DocId))
		return nil, errors.New("解密失败")
	}
	o, _ := s.orm.QueryNotDeleted().Where(purchaseorder.DocID(docId), purchaseorder.RiderID(r.ID)).First(ctx)
	if o == nil {
		return nil, errors.New("未找到合同信息")
	}
	return &definition.ContractDetailRes{
		Link: o.ContractURL,
	}, nil
}
