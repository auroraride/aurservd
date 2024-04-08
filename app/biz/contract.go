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
	"github.com/auroraride/aurservd/internal/ent/cabinet"
	"github.com/auroraride/aurservd/internal/ent/contract"
	"github.com/auroraride/aurservd/internal/ent/contracttemplate"
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
func (s *Contract) Sign(r *ent.Rider, req *definition.ContractSignReq) (err error) {
	// 查找订阅
	sub, _ := ent.Database.Subscribe.QueryNotDeleted().
		Where(
			subscribe.ID(req.SubscribeID),
			subscribe.Status(model.SubscribeStatusInactive),
		).WithCity(func(query *ent.CityQuery) {
		query.WithParent()
	}).First(s.ctx)
	if sub == nil {
		return errors.New("未找到骑士卡")
	}

	// 是否免签或已签约
	if !service.NewSubscribe().NeedContract(sub) {
		return errors.New("当前订阅无需签约")
	}

	if sub.BrandID == nil && sub.EbikeID != nil {
		return errors.New("当前订阅错误")
	}

	// 查找分配信息
	allo := service.NewAllocate().QueryEffectiveSubscribeIDX(sub.ID)
	if allo == nil {
		return errors.New("未找到分配信息")
	}

	if sub.BrandID != nil && allo.StoreID == nil && allo.StationID == nil {
		return errors.New("电车必须由门店或站点分配")
	}

	person, _ := r.QueryPerson().First(s.ctx)
	if person == nil {
		return errors.New("未找到骑手信息")
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
		return errors.New("未找到合同信息")
	}

	// 请求签署合同
	url, err := rpc.Sgin(s.ctx, &pb.ContractSignRequest{
		Sn: cont.Sn,
		Sign: &pb.ContractSignEntity{
			Seal:     req.Seal,
			Name:     person.Name,
			Province: province,
			City:     city,
			Address:  person.AuthResult.Address,
			Phone:    r.Phone,
			Idcard:   person.IDCardNumber,
		},
	})
	if err != nil {
		zap.L().Error("签署合同失败", zap.Error(err))
		return err
	}

	//  更新合同状态
	err = cont.Update().SetStatus(model.ContractStatusSuccess.Value()).SetLink(url).Exec(s.ctx)

	return nil
}

// Create 骑手添加合同
func (s *Contract) Create(r *ent.Rider, req *model.BusinessCabinetReq) (*model.ContractSignRes, error) {

	sub, _ := ent.Database.Subscribe.QueryNotDeleted().Where(subscribe.ID(req.ID), subscribe.Status(model.SubscribeStatusInactive)).WithCity().First(s.ctx)
	if sub == nil {
		return nil, errors.New("未找到骑士卡")
	}

	// 是否免签或已签约
	if !service.NewSubscribe().NeedContract(sub) {
		return &model.ContractSignRes{Effective: true}, nil
	}

	if sub.BrandID == nil && sub.EbikeID != nil {
		return nil, errors.New("当前订阅错误")
	}

	// 查询电柜
	cab, _ := ent.Database.Cabinet.QueryNotDeleted().Where(cabinet.Serial(req.Serial)).First(s.ctx)
	if cab == nil {
		return nil, errors.New("未找到电柜")

	}

	var link, sn string
	skip := false
	co, _ := s.orm.QueryNotDeleted().Where(contract.SubscribeID(sub.ID), contract.LinkNotNil(), contract.Status(model.ContractStatusSigning.Value())).First(s.ctx)
	// 判定是否生成过合同
	if co != nil {
		// 合同处于有效期内跳过生成
		if co.ExpiresAt != nil && co.ExpiresAt.After(time.Now()) {
			skip = true
			link = *co.Link
			sn = co.Sn
		} else {
			// 否则删除原合同重新生成
			s.orm.DeleteOne(co).ExecX(s.ctx)
		}
	}

	// 查询分配信息是否存在, 如果存在则删除
	service.NewAllocate().SubscribeDeleteIfExists(sub.ID)
	// 存储分配信息
	allo, err := ent.Database.Allocate.Create().
		SetType(allocate.TypeBattery).
		SetSubscribe(sub).
		SetRider(r).
		SetStatus(model.AllocateStatusPending.Value()).
		SetTime(time.Now()).
		SetModel(sub.Model).
		SetCabinetID(cab.ID).
		SetRemark("骑手扫码").
		Save(s.ctx)
	if err != nil {
		return nil, err
	}

	if sub.BrandID != nil && allo.StoreID == nil && allo.StationID == nil {
		return nil, errors.New("电车必须由门店或站点分配")
	}

	ec := sub.Edges.City

	if !skip {

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
		sn = tools.NewUnique().NewSN()
		// 填充公共变量
		// 合同编号
		m["sn"] = sn
		// 骑手姓名
		m["name"] = p.Name
		// 身份证号
		m["idcard"] = p.IDCardNumber
		// 户口地址
		m["address"] = p.AuthResult.Address
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
			// 团签代缴
			m["payEnt"] = true
		} else {
			// 个签骑士卡
			un = service.NewContract().PlanData(sub)
			// 骑手缴费
			m["payRider"] = true
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

		// 查询模版
		q := ent.Database.ContractTemplate.QueryNotDeleted()
		if isEnterprise {
			q.Where(contracttemplate.UserType(1))
		} else {
			q.Where(contracttemplate.UserType(2))
		}
		if sub.BrandID != nil {
			q.Where(contracttemplate.SubType(1))
		} else {
			q.Where(contracttemplate.SubType(2))
		}
		temp, _ := q.Where(contracttemplate.Enable(true)).First(s.ctx)

		if temp == nil {
			return nil, errors.New("未找到合同模板")
		}

		// 转换ar.Map 为 map[string]string
		values := make(map[string]string)
		for k, v := range m {
			switch v.(type) {
			case bool:
				values[k] = "On"
			case string:
				values[k] = v.(string)
			case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
				str, err := numberToString(v)
				if err != nil {
					return nil, err
				}
				values[k] = str
			}
		}

		contractCreateResponse, err := rpc.Create(s.ctx, temp.Sn, values)
		if err != nil {
			return nil, err
		}

		link = contractCreateResponse.Url

		// 存储合同信息
		err = ent.WithTx(s.ctx, func(tx *ent.Tx) (err error) {
			// todo 生成合同流水号
			flowId := tools.NewUnique().NewSN28()
			expiresAt := time.Now().Add(model.ContractExpiration * time.Minute)
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
				Exec(context.Background())
			if err != nil {
				return err
			}
			return sub.Update().UpdateTarget(allo.CabinetID, allo.StoreID, allo.EmployeeID).Exec(s.ctx)
		})
		if err != nil {
			return nil, err
		}
	}

	return &model.ContractSignRes{
		Url: link,
		Sn:  sn,
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
