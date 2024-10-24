// Code generated by ent, DO NOT EDIT.

package ent

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ent/goods"
	"github.com/auroraride/aurservd/internal/ent/purchaseorder"
	"github.com/auroraride/aurservd/internal/ent/rider"
	"github.com/auroraride/aurservd/internal/ent/store"
)

// PurchaseOrder is the model entity for the PurchaseOrder schema.
type PurchaseOrder struct {
	config `json:"-"`
	// ID of the ent.
	ID uint64 `json:"id,omitempty"`
	// CreatedAt holds the value of the "created_at" field.
	CreatedAt time.Time `json:"created_at,omitempty"`
	// UpdatedAt holds the value of the "updated_at" field.
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	// DeletedAt holds the value of the "deleted_at" field.
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
	// 创建人
	Creator *model.Modifier `json:"creator,omitempty"`
	// 最后修改人
	LastModifier *model.Modifier `json:"last_modifier,omitempty"`
	// 管理员改动原因/备注
	Remark string `json:"remark,omitempty"`
	// 骑手ID
	RiderID uint64 `json:"rider_id,omitempty"`
	// 商品ID
	GoodsID uint64 `json:"goods_id,omitempty"`
	// 门店ID
	StoreID *uint64 `json:"store_id,omitempty"`
	// 车架号
	Sn string `json:"sn,omitempty"`
	// 状态, pending: 待支付, staging: 分期执行中, ended: 已完成, cancelled: 已取消, refunded: 已退款
	Status purchaseorder.Status `json:"status,omitempty"`
	// 合同URL
	ContractURL string `json:"contract_url,omitempty"`
	// 合同ID
	DocID string `json:"doc_id,omitempty"`
	// 是否签约
	Signed bool `json:"signed,omitempty"`
	// 当前分期阶段，从0开始
	InstallmentStage int `json:"installment_stage,omitempty"`
	// 分期总数
	InstallmentTotal int `json:"installment_total,omitempty"`
	// 分期方案
	InstallmentPlan model.GoodsPaymentPlan `json:"installment_plan,omitempty"`
	// 开始日期
	StartDate *time.Time `json:"start_date,omitempty"`
	// 下次支付日期
	NextDate *time.Time `json:"next_date,omitempty"`
	// 图片
	Images []string `json:"images,omitempty"`
	// 激活人姓名
	ActiveName string `json:"active_name,omitempty"`
	// 激活人电话
	ActivePhone string `json:"active_phone,omitempty"`
	// 车辆颜色
	Color string `json:"color,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the PurchaseOrderQuery when eager-loading is set.
	Edges        PurchaseOrderEdges `json:"edges"`
	selectValues sql.SelectValues
}

// PurchaseOrderEdges holds the relations/edges for other nodes in the graph.
type PurchaseOrderEdges struct {
	// 骑手
	Rider *Rider `json:"rider,omitempty"`
	// 商品ID
	Goods *Goods `json:"goods,omitempty"`
	// Store holds the value of the store edge.
	Store *Store `json:"store,omitempty"`
	// Payments holds the value of the payments edge.
	Payments []*PurchasePayment `json:"payments,omitempty"`
	// Follows holds the value of the follows edge.
	Follows []*PurchaseFollow `json:"follows,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [5]bool
}

// RiderOrErr returns the Rider value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e PurchaseOrderEdges) RiderOrErr() (*Rider, error) {
	if e.Rider != nil {
		return e.Rider, nil
	} else if e.loadedTypes[0] {
		return nil, &NotFoundError{label: rider.Label}
	}
	return nil, &NotLoadedError{edge: "rider"}
}

// GoodsOrErr returns the Goods value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e PurchaseOrderEdges) GoodsOrErr() (*Goods, error) {
	if e.Goods != nil {
		return e.Goods, nil
	} else if e.loadedTypes[1] {
		return nil, &NotFoundError{label: goods.Label}
	}
	return nil, &NotLoadedError{edge: "goods"}
}

// StoreOrErr returns the Store value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e PurchaseOrderEdges) StoreOrErr() (*Store, error) {
	if e.Store != nil {
		return e.Store, nil
	} else if e.loadedTypes[2] {
		return nil, &NotFoundError{label: store.Label}
	}
	return nil, &NotLoadedError{edge: "store"}
}

// PaymentsOrErr returns the Payments value or an error if the edge
// was not loaded in eager-loading.
func (e PurchaseOrderEdges) PaymentsOrErr() ([]*PurchasePayment, error) {
	if e.loadedTypes[3] {
		return e.Payments, nil
	}
	return nil, &NotLoadedError{edge: "payments"}
}

// FollowsOrErr returns the Follows value or an error if the edge
// was not loaded in eager-loading.
func (e PurchaseOrderEdges) FollowsOrErr() ([]*PurchaseFollow, error) {
	if e.loadedTypes[4] {
		return e.Follows, nil
	}
	return nil, &NotLoadedError{edge: "follows"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*PurchaseOrder) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case purchaseorder.FieldCreator, purchaseorder.FieldLastModifier, purchaseorder.FieldInstallmentPlan, purchaseorder.FieldImages:
			values[i] = new([]byte)
		case purchaseorder.FieldSigned:
			values[i] = new(sql.NullBool)
		case purchaseorder.FieldID, purchaseorder.FieldRiderID, purchaseorder.FieldGoodsID, purchaseorder.FieldStoreID, purchaseorder.FieldInstallmentStage, purchaseorder.FieldInstallmentTotal:
			values[i] = new(sql.NullInt64)
		case purchaseorder.FieldRemark, purchaseorder.FieldSn, purchaseorder.FieldStatus, purchaseorder.FieldContractURL, purchaseorder.FieldDocID, purchaseorder.FieldActiveName, purchaseorder.FieldActivePhone, purchaseorder.FieldColor:
			values[i] = new(sql.NullString)
		case purchaseorder.FieldCreatedAt, purchaseorder.FieldUpdatedAt, purchaseorder.FieldDeletedAt, purchaseorder.FieldStartDate, purchaseorder.FieldNextDate:
			values[i] = new(sql.NullTime)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the PurchaseOrder fields.
func (po *PurchaseOrder) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case purchaseorder.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			po.ID = uint64(value.Int64)
		case purchaseorder.FieldCreatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field created_at", values[i])
			} else if value.Valid {
				po.CreatedAt = value.Time
			}
		case purchaseorder.FieldUpdatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field updated_at", values[i])
			} else if value.Valid {
				po.UpdatedAt = value.Time
			}
		case purchaseorder.FieldDeletedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field deleted_at", values[i])
			} else if value.Valid {
				po.DeletedAt = new(time.Time)
				*po.DeletedAt = value.Time
			}
		case purchaseorder.FieldCreator:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field creator", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &po.Creator); err != nil {
					return fmt.Errorf("unmarshal field creator: %w", err)
				}
			}
		case purchaseorder.FieldLastModifier:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field last_modifier", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &po.LastModifier); err != nil {
					return fmt.Errorf("unmarshal field last_modifier: %w", err)
				}
			}
		case purchaseorder.FieldRemark:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field remark", values[i])
			} else if value.Valid {
				po.Remark = value.String
			}
		case purchaseorder.FieldRiderID:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field rider_id", values[i])
			} else if value.Valid {
				po.RiderID = uint64(value.Int64)
			}
		case purchaseorder.FieldGoodsID:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field goods_id", values[i])
			} else if value.Valid {
				po.GoodsID = uint64(value.Int64)
			}
		case purchaseorder.FieldStoreID:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field store_id", values[i])
			} else if value.Valid {
				po.StoreID = new(uint64)
				*po.StoreID = uint64(value.Int64)
			}
		case purchaseorder.FieldSn:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field sn", values[i])
			} else if value.Valid {
				po.Sn = value.String
			}
		case purchaseorder.FieldStatus:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field status", values[i])
			} else if value.Valid {
				po.Status = purchaseorder.Status(value.String)
			}
		case purchaseorder.FieldContractURL:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field contract_url", values[i])
			} else if value.Valid {
				po.ContractURL = value.String
			}
		case purchaseorder.FieldDocID:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field doc_id", values[i])
			} else if value.Valid {
				po.DocID = value.String
			}
		case purchaseorder.FieldSigned:
			if value, ok := values[i].(*sql.NullBool); !ok {
				return fmt.Errorf("unexpected type %T for field signed", values[i])
			} else if value.Valid {
				po.Signed = value.Bool
			}
		case purchaseorder.FieldInstallmentStage:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field installment_stage", values[i])
			} else if value.Valid {
				po.InstallmentStage = int(value.Int64)
			}
		case purchaseorder.FieldInstallmentTotal:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field installment_total", values[i])
			} else if value.Valid {
				po.InstallmentTotal = int(value.Int64)
			}
		case purchaseorder.FieldInstallmentPlan:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field installment_plan", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &po.InstallmentPlan); err != nil {
					return fmt.Errorf("unmarshal field installment_plan: %w", err)
				}
			}
		case purchaseorder.FieldStartDate:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field start_date", values[i])
			} else if value.Valid {
				po.StartDate = new(time.Time)
				*po.StartDate = value.Time
			}
		case purchaseorder.FieldNextDate:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field next_date", values[i])
			} else if value.Valid {
				po.NextDate = new(time.Time)
				*po.NextDate = value.Time
			}
		case purchaseorder.FieldImages:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field images", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &po.Images); err != nil {
					return fmt.Errorf("unmarshal field images: %w", err)
				}
			}
		case purchaseorder.FieldActiveName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field active_name", values[i])
			} else if value.Valid {
				po.ActiveName = value.String
			}
		case purchaseorder.FieldActivePhone:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field active_phone", values[i])
			} else if value.Valid {
				po.ActivePhone = value.String
			}
		case purchaseorder.FieldColor:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field color", values[i])
			} else if value.Valid {
				po.Color = value.String
			}
		default:
			po.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the PurchaseOrder.
// This includes values selected through modifiers, order, etc.
func (po *PurchaseOrder) Value(name string) (ent.Value, error) {
	return po.selectValues.Get(name)
}

// QueryRider queries the "rider" edge of the PurchaseOrder entity.
func (po *PurchaseOrder) QueryRider() *RiderQuery {
	return NewPurchaseOrderClient(po.config).QueryRider(po)
}

// QueryGoods queries the "goods" edge of the PurchaseOrder entity.
func (po *PurchaseOrder) QueryGoods() *GoodsQuery {
	return NewPurchaseOrderClient(po.config).QueryGoods(po)
}

// QueryStore queries the "store" edge of the PurchaseOrder entity.
func (po *PurchaseOrder) QueryStore() *StoreQuery {
	return NewPurchaseOrderClient(po.config).QueryStore(po)
}

// QueryPayments queries the "payments" edge of the PurchaseOrder entity.
func (po *PurchaseOrder) QueryPayments() *PurchasePaymentQuery {
	return NewPurchaseOrderClient(po.config).QueryPayments(po)
}

// QueryFollows queries the "follows" edge of the PurchaseOrder entity.
func (po *PurchaseOrder) QueryFollows() *PurchaseFollowQuery {
	return NewPurchaseOrderClient(po.config).QueryFollows(po)
}

// Update returns a builder for updating this PurchaseOrder.
// Note that you need to call PurchaseOrder.Unwrap() before calling this method if this PurchaseOrder
// was returned from a transaction, and the transaction was committed or rolled back.
func (po *PurchaseOrder) Update() *PurchaseOrderUpdateOne {
	return NewPurchaseOrderClient(po.config).UpdateOne(po)
}

// Unwrap unwraps the PurchaseOrder entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (po *PurchaseOrder) Unwrap() *PurchaseOrder {
	_tx, ok := po.config.driver.(*txDriver)
	if !ok {
		panic("ent: PurchaseOrder is not a transactional entity")
	}
	po.config.driver = _tx.drv
	return po
}

// String implements the fmt.Stringer.
func (po *PurchaseOrder) String() string {
	var builder strings.Builder
	builder.WriteString("PurchaseOrder(")
	builder.WriteString(fmt.Sprintf("id=%v, ", po.ID))
	builder.WriteString("created_at=")
	builder.WriteString(po.CreatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("updated_at=")
	builder.WriteString(po.UpdatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	if v := po.DeletedAt; v != nil {
		builder.WriteString("deleted_at=")
		builder.WriteString(v.Format(time.ANSIC))
	}
	builder.WriteString(", ")
	builder.WriteString("creator=")
	builder.WriteString(fmt.Sprintf("%v", po.Creator))
	builder.WriteString(", ")
	builder.WriteString("last_modifier=")
	builder.WriteString(fmt.Sprintf("%v", po.LastModifier))
	builder.WriteString(", ")
	builder.WriteString("remark=")
	builder.WriteString(po.Remark)
	builder.WriteString(", ")
	builder.WriteString("rider_id=")
	builder.WriteString(fmt.Sprintf("%v", po.RiderID))
	builder.WriteString(", ")
	builder.WriteString("goods_id=")
	builder.WriteString(fmt.Sprintf("%v", po.GoodsID))
	builder.WriteString(", ")
	if v := po.StoreID; v != nil {
		builder.WriteString("store_id=")
		builder.WriteString(fmt.Sprintf("%v", *v))
	}
	builder.WriteString(", ")
	builder.WriteString("sn=")
	builder.WriteString(po.Sn)
	builder.WriteString(", ")
	builder.WriteString("status=")
	builder.WriteString(fmt.Sprintf("%v", po.Status))
	builder.WriteString(", ")
	builder.WriteString("contract_url=")
	builder.WriteString(po.ContractURL)
	builder.WriteString(", ")
	builder.WriteString("doc_id=")
	builder.WriteString(po.DocID)
	builder.WriteString(", ")
	builder.WriteString("signed=")
	builder.WriteString(fmt.Sprintf("%v", po.Signed))
	builder.WriteString(", ")
	builder.WriteString("installment_stage=")
	builder.WriteString(fmt.Sprintf("%v", po.InstallmentStage))
	builder.WriteString(", ")
	builder.WriteString("installment_total=")
	builder.WriteString(fmt.Sprintf("%v", po.InstallmentTotal))
	builder.WriteString(", ")
	builder.WriteString("installment_plan=")
	builder.WriteString(fmt.Sprintf("%v", po.InstallmentPlan))
	builder.WriteString(", ")
	if v := po.StartDate; v != nil {
		builder.WriteString("start_date=")
		builder.WriteString(v.Format(time.ANSIC))
	}
	builder.WriteString(", ")
	if v := po.NextDate; v != nil {
		builder.WriteString("next_date=")
		builder.WriteString(v.Format(time.ANSIC))
	}
	builder.WriteString(", ")
	builder.WriteString("images=")
	builder.WriteString(fmt.Sprintf("%v", po.Images))
	builder.WriteString(", ")
	builder.WriteString("active_name=")
	builder.WriteString(po.ActiveName)
	builder.WriteString(", ")
	builder.WriteString("active_phone=")
	builder.WriteString(po.ActivePhone)
	builder.WriteString(", ")
	builder.WriteString("color=")
	builder.WriteString(po.Color)
	builder.WriteByte(')')
	return builder.String()
}

// PurchaseOrders is a parsable slice of PurchaseOrder.
type PurchaseOrders []*PurchaseOrder
