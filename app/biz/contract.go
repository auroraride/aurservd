package biz

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
	"github.com/auroraride/aurservd/app/rpc"
	"github.com/auroraride/aurservd/app/service"
	"github.com/auroraride/aurservd/internal/ar"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/allocate"
	"github.com/auroraride/aurservd/internal/ent/contract"
	"github.com/auroraride/aurservd/internal/ent/subscribe"
	"github.com/auroraride/aurservd/pkg/tools"
)

type Contract struct {
	orm *ent.ContractClient
	ctx context.Context
}

func NewContract() *Contract {
	return &Contract{
		orm: ent.Database.Contract,
		ctx: context.Background(),
	}
}

// Sign 签约
func (s *Contract) Sign(r *ent.Rider, req *definition.ContractSignNewReq) (res *definition.ContractSignNewRes, err error) {
	// 查找订阅
	sub, _ := ent.Database.Subscribe.QueryNotDeleted().
		Where(
			subscribe.ID(req.SubscribeID),
			subscribe.Status(model.SubscribeStatusInactive),
		).WithCity(func(query *ent.CityQuery) {
		query.WithParent()
	}).First(s.ctx)
	if sub == nil {
		return nil, errors.New("未找到骑士卡")
	}

	// 是否免签或已签约
	if !service.NewSubscribe().NeedContract(sub) {
		return nil, errors.New("当前订阅无需签约")
	}

	if sub.BrandID == nil && sub.EbikeID != nil {
		return nil, errors.New("当前订阅错误")
	}

	// 查找分配信息
	allo, _ := service.NewAllocate().QueryEffectiveSubscribeID(sub.ID)
	if allo == nil {
		return nil, errors.New("未找到分配信息")
	}

	if sub.BrandID != nil && allo.StoreID == nil && allo.StationID == nil {
		return nil, errors.New("电车必须由门店或站点分配")
	}

	// 城市
	ec := sub.Edges.City
	if ec == nil {
		return nil, errors.New("未找到城市信息")
	}

	// 判定非智能套餐门店库存
	if allo.StoreID != nil && allo.BatteryID == nil && !service.NewStock().CheckStore(*allo.StoreID, sub.Model, 1) {
		return nil, errors.New("库存不足")
	}

	person, _ := r.QueryPerson().First(s.ctx)
	if person == nil {
		return nil, errors.New("未找到骑手信息")
	}

	var city, province string
	if sub.Edges.City != nil {
		city = sub.Edges.City.Name
		// 获取省份
		if sub.Edges.City.Edges.Parent != nil {
			province = sub.Edges.City.Edges.Parent.Name
		}
	}

	cont, _ := ent.Database.Contract.QueryNotDeleted().Where(contract.AllocateID(allo.ID), contract.SubscribeID(req.SubscribeID)).First(s.ctx)
	if cont == nil {
		return nil, errors.New("未找到合同信息")
	}
	// 获取模版id
	cfg := ar.Config.Contract
	// 请求签署合同
	url, err := rpc.Sgin(s.ctx, &pb.ContractSignRequest{
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

	var files []string
	files = append(files, url)

	now := time.Now()

	//  更新合同状态
	err = cont.Update().
		SetStatus(model.ContractStatusSuccess.Value()).
		SetFiles(files).
		SetSignedAt(now).
		SetEffective(true).
		Exec(s.ctx)
	if err != nil {
		zap.L().Error("更新合同状态失败", zap.Error(err))
		return nil, err
	}

	err = service.NewContract().Update(cont)
	if err != nil {
		zap.L().Error("更新合同状态失败", zap.Error(err))
		return nil, err
	}

	if strings.HasPrefix(url, "https://c.auroraride.com/") {
		hash := strings.TrimPrefix(url, "https://c.auroraride.com/")
		// 发送短信
		_, err = service.NewSms().SendSignSuccess(now, "时光驹电动车电池租赁合同", hash, r.Phone)
		if err != nil {
			return nil, err
		}
	}

	return &definition.ContractSignNewRes{
		Link: url,
	}, nil
}

// Create 骑手添加合同
func (s *Contract) Create(r *ent.Rider, req *definition.ContractCreateReq) (*definition.ContractCreateRes, error) {
	sub, _ := ent.Database.Subscribe.QueryNotDeleted().Where(subscribe.ID(req.SubscribeID), subscribe.Status(model.SubscribeStatusInactive)).WithCity().WithPlan().First(s.ctx)
	if sub == nil {
		return nil, errors.New("未找到骑士卡")
	}

	// 是否免签或已签约
	if !service.NewSubscribe().NeedContract(sub) {
		return &definition.ContractCreateRes{Effective: true}, nil
	}

	if sub.BrandID == nil && sub.EbikeID != nil {
		return nil, errors.New("当前订阅错误")
	}

	var link, docId string
	skip := false
	co, _ := s.orm.QueryNotDeleted().Where(contract.SubscribeID(sub.ID), contract.LinkNotNil(), contract.Status(model.ContractStatusSigning.Value())).First(s.ctx)
	// 判定是否生成过合同
	if co != nil {
		// 合同处于有效期内跳过生成
		if co.ExpiresAt != nil && co.ExpiresAt.After(time.Now()) {
			skip = true
			link = *co.Link
			docId = co.DocID
		} else {
			// 否则删除原合同重新生成
			s.orm.DeleteOne(co).ExecX(s.ctx)
		}
	}

	ec := sub.Edges.City

	if !skip {
		var err error
		// 查询是否分配
		allo, _ := service.NewAllocate().QueryEffectiveSubscribeID(sub.ID)
		// 个签单电
		if sub.BrandID == nil && sub.EnterpriseID == nil {
			if allo != nil {
				// 查询分配信息是否存在, 如果存在则删除
				service.NewAllocate().SubscribeDeleteIfExists(sub.ID)
				// 存储分配信息
				allo, err = ent.Database.Allocate.Create().
					SetType(allocate.TypeBattery).
					SetSubscribe(sub).
					SetRider(r).
					SetStatus(model.AllocateStatusPending.Value()).
					SetTime(time.Now()).
					SetModel(sub.Model).
					SetRemark("骑手自主激活").
					Save(s.ctx)
				if err != nil {
					return nil, err
				}
			}
		} else {
			if allo != nil && sub.BrandID != nil && allo.StoreID == nil && allo.StationID == nil {
				return nil, errors.New("电车必须由门店或站点分配")
			}
		}

		// 定义变量
		var (
			m            = make(ar.Map)
			p            = service.NewPerson().GetNormalAuthedPerson(r)
			isEnterprise = r.EnterpriseID != nil

			// 电池型号
			bm = strings.ToUpper(sub.Model)
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

		if isEnterprise {
			// 设置团签字段
			un = service.NewContract().EnterpriseData(m, sub)
			delete(m, "payerEnt")
		} else {
			// 个签骑士卡
			un = service.NewContract().PlanData(sub)
		}

		if un == nil {
			return nil, errors.New("合同信息错误")
		}

		m["payMonth"] = un.Month

		// 电车
		if sub.BrandID != nil {
			// 查找电车分配
			bike, _ := allo.QueryEbike().WithBrand().First(s.ctx)
			if bike == nil || bike.Edges.Brand == nil {
				return nil, errors.New("未找到电车信息")
			}

			brand := bike.Edges.Brand

			// 车加电方案
			m["schemaEbike"] = true
			// 车加电方案一
			m["ebikeScheme1"] = true
			// 车辆品牌
			m["ebikeBrand"] = brand.Name
			// 车辆颜色
			m["ebikeColor"] = bike.Color
			// 车架号
			m["ebikeSN"] = bike.Sn
			// 车牌号
			m["ebikePlate"] = bike.Plate
			// 电池类型
			m["ebikeBattery"] = "时光驹电池"
			// 电池规格
			m["ebikeModel"] = bm
			// 车电方案一开始日期
			m["ebikeScheme1Start"] = now
			// 车电方案一截止日
			m["ebikeScheme1Stop"] = un.Stop
			// 车电方案一月租金
			m["ebikeScheme1Price"] = un.Price
			// 车电方案一首次缴纳月数
			m["ebikeScheme1PayMonth"] = un.Month
			// 车电方案一首次缴纳租金
			m["ebikeScheme1PayTotal"] = un.Total
		} else {
			// 单电方案
			m["schemaBattery"] = true
			// 电池规格
			m["batteryModel"] = bm
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
		}

		// 获取模版id
		cfg := ar.Config.Contract
		templateID := cfg.Template.Personal
		if isEnterprise {
			templateID = cfg.Template.Enterprise
		}

		values := make(map[string]*pb.ContractFromField)

		for k, v := range m {
			switch v.(type) {
			case bool:
				values[k] = &pb.ContractFromField{
					Value: &pb.ContractFromField_Checkbox{Checkbox: v.(bool)},
				}
			case string:
				values[k] = &pb.ContractFromField{
					Value: &pb.ContractFromField_Text{Text: v.(string)},
				}
			case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
				str, err := numberToString(v)
				if err != nil {
					return nil, err
				}
				values[k] = &pb.ContractFromField{
					Value: &pb.ContractFromField_Text{Text: str},
				}
			}
		}

		expiresAt := time.Now().Add(model.ContractExpiration * time.Minute)
		contractCreateResponse, err := rpc.Create(s.ctx, values, &definition.ContractCreateRPCReq{
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

		flowId := tools.NewUnique().NewSN28()
		// 存储合同信息
		err = ent.WithTx(s.ctx, func(tx *ent.Tx) (err error) {
			// 删除原有合同
			_, _ = tx.Contract.Delete().Where(contract.AllocateID(allo.ID)).Exec(s.ctx)
			err = tx.Contract.Create().
				SetFlowID(flowId).
				SetRiderID(r.ID).
				SetStatus(model.ContractStatusSigning.Value()).
				SetSn(sn).
				SetNillableEmployeeID(allo.EmployeeID).
				SetAllocateID(allo.ID).
				SetSubscribe(sub).
				SetRiderInfo(&model.ContractRider{
					Phone:        r.Phone,
					Name:         r.Name,
					IDCardNumber: r.IDCardNumber,
				}).
				SetLink(contractCreateResponse.Url).
				SetExpiresAt(expiresAt).
				SetDocID(docId).
				Exec(context.Background())
			if err != nil {
				return err
			}
			return sub.Update().UpdateTarget(allo.CabinetID, allo.StoreID, allo.EmployeeID).Exec(s.ctx)
		})
		if err != nil {
			return nil, err
		}
		go service.NewContract().CheckResult(flowId)
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
