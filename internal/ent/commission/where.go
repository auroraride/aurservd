// Code generated by ent, DO NOT EDIT.

package commission

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/auroraride/aurservd/internal/ent/predicate"
)

// ID filters vertices based on their ID field.
func ID(id uint64) predicate.Commission {
	return predicate.Commission(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id uint64) predicate.Commission {
	return predicate.Commission(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id uint64) predicate.Commission {
	return predicate.Commission(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...uint64) predicate.Commission {
	return predicate.Commission(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...uint64) predicate.Commission {
	return predicate.Commission(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id uint64) predicate.Commission {
	return predicate.Commission(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id uint64) predicate.Commission {
	return predicate.Commission(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id uint64) predicate.Commission {
	return predicate.Commission(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id uint64) predicate.Commission {
	return predicate.Commission(sql.FieldLTE(FieldID, id))
}

// CreatedAt applies equality check predicate on the "created_at" field. It's identical to CreatedAtEQ.
func CreatedAt(v time.Time) predicate.Commission {
	return predicate.Commission(sql.FieldEQ(FieldCreatedAt, v))
}

// UpdatedAt applies equality check predicate on the "updated_at" field. It's identical to UpdatedAtEQ.
func UpdatedAt(v time.Time) predicate.Commission {
	return predicate.Commission(sql.FieldEQ(FieldUpdatedAt, v))
}

// DeletedAt applies equality check predicate on the "deleted_at" field. It's identical to DeletedAtEQ.
func DeletedAt(v time.Time) predicate.Commission {
	return predicate.Commission(sql.FieldEQ(FieldDeletedAt, v))
}

// Remark applies equality check predicate on the "remark" field. It's identical to RemarkEQ.
func Remark(v string) predicate.Commission {
	return predicate.Commission(sql.FieldEQ(FieldRemark, v))
}

// BusinessID applies equality check predicate on the "business_id" field. It's identical to BusinessIDEQ.
func BusinessID(v uint64) predicate.Commission {
	return predicate.Commission(sql.FieldEQ(FieldBusinessID, v))
}

// SubscribeID applies equality check predicate on the "subscribe_id" field. It's identical to SubscribeIDEQ.
func SubscribeID(v uint64) predicate.Commission {
	return predicate.Commission(sql.FieldEQ(FieldSubscribeID, v))
}

// PlanID applies equality check predicate on the "plan_id" field. It's identical to PlanIDEQ.
func PlanID(v uint64) predicate.Commission {
	return predicate.Commission(sql.FieldEQ(FieldPlanID, v))
}

// RiderID applies equality check predicate on the "rider_id" field. It's identical to RiderIDEQ.
func RiderID(v uint64) predicate.Commission {
	return predicate.Commission(sql.FieldEQ(FieldRiderID, v))
}

// OrderID applies equality check predicate on the "order_id" field. It's identical to OrderIDEQ.
func OrderID(v uint64) predicate.Commission {
	return predicate.Commission(sql.FieldEQ(FieldOrderID, v))
}

// Amount applies equality check predicate on the "amount" field. It's identical to AmountEQ.
func Amount(v float64) predicate.Commission {
	return predicate.Commission(sql.FieldEQ(FieldAmount, v))
}

// Status applies equality check predicate on the "status" field. It's identical to StatusEQ.
func Status(v uint8) predicate.Commission {
	return predicate.Commission(sql.FieldEQ(FieldStatus, v))
}

// EmployeeID applies equality check predicate on the "employee_id" field. It's identical to EmployeeIDEQ.
func EmployeeID(v uint64) predicate.Commission {
	return predicate.Commission(sql.FieldEQ(FieldEmployeeID, v))
}

// CreatedAtEQ applies the EQ predicate on the "created_at" field.
func CreatedAtEQ(v time.Time) predicate.Commission {
	return predicate.Commission(sql.FieldEQ(FieldCreatedAt, v))
}

// CreatedAtNEQ applies the NEQ predicate on the "created_at" field.
func CreatedAtNEQ(v time.Time) predicate.Commission {
	return predicate.Commission(sql.FieldNEQ(FieldCreatedAt, v))
}

// CreatedAtIn applies the In predicate on the "created_at" field.
func CreatedAtIn(vs ...time.Time) predicate.Commission {
	return predicate.Commission(sql.FieldIn(FieldCreatedAt, vs...))
}

// CreatedAtNotIn applies the NotIn predicate on the "created_at" field.
func CreatedAtNotIn(vs ...time.Time) predicate.Commission {
	return predicate.Commission(sql.FieldNotIn(FieldCreatedAt, vs...))
}

// CreatedAtGT applies the GT predicate on the "created_at" field.
func CreatedAtGT(v time.Time) predicate.Commission {
	return predicate.Commission(sql.FieldGT(FieldCreatedAt, v))
}

// CreatedAtGTE applies the GTE predicate on the "created_at" field.
func CreatedAtGTE(v time.Time) predicate.Commission {
	return predicate.Commission(sql.FieldGTE(FieldCreatedAt, v))
}

// CreatedAtLT applies the LT predicate on the "created_at" field.
func CreatedAtLT(v time.Time) predicate.Commission {
	return predicate.Commission(sql.FieldLT(FieldCreatedAt, v))
}

// CreatedAtLTE applies the LTE predicate on the "created_at" field.
func CreatedAtLTE(v time.Time) predicate.Commission {
	return predicate.Commission(sql.FieldLTE(FieldCreatedAt, v))
}

// UpdatedAtEQ applies the EQ predicate on the "updated_at" field.
func UpdatedAtEQ(v time.Time) predicate.Commission {
	return predicate.Commission(sql.FieldEQ(FieldUpdatedAt, v))
}

// UpdatedAtNEQ applies the NEQ predicate on the "updated_at" field.
func UpdatedAtNEQ(v time.Time) predicate.Commission {
	return predicate.Commission(sql.FieldNEQ(FieldUpdatedAt, v))
}

// UpdatedAtIn applies the In predicate on the "updated_at" field.
func UpdatedAtIn(vs ...time.Time) predicate.Commission {
	return predicate.Commission(sql.FieldIn(FieldUpdatedAt, vs...))
}

// UpdatedAtNotIn applies the NotIn predicate on the "updated_at" field.
func UpdatedAtNotIn(vs ...time.Time) predicate.Commission {
	return predicate.Commission(sql.FieldNotIn(FieldUpdatedAt, vs...))
}

// UpdatedAtGT applies the GT predicate on the "updated_at" field.
func UpdatedAtGT(v time.Time) predicate.Commission {
	return predicate.Commission(sql.FieldGT(FieldUpdatedAt, v))
}

// UpdatedAtGTE applies the GTE predicate on the "updated_at" field.
func UpdatedAtGTE(v time.Time) predicate.Commission {
	return predicate.Commission(sql.FieldGTE(FieldUpdatedAt, v))
}

// UpdatedAtLT applies the LT predicate on the "updated_at" field.
func UpdatedAtLT(v time.Time) predicate.Commission {
	return predicate.Commission(sql.FieldLT(FieldUpdatedAt, v))
}

// UpdatedAtLTE applies the LTE predicate on the "updated_at" field.
func UpdatedAtLTE(v time.Time) predicate.Commission {
	return predicate.Commission(sql.FieldLTE(FieldUpdatedAt, v))
}

// DeletedAtEQ applies the EQ predicate on the "deleted_at" field.
func DeletedAtEQ(v time.Time) predicate.Commission {
	return predicate.Commission(sql.FieldEQ(FieldDeletedAt, v))
}

// DeletedAtNEQ applies the NEQ predicate on the "deleted_at" field.
func DeletedAtNEQ(v time.Time) predicate.Commission {
	return predicate.Commission(sql.FieldNEQ(FieldDeletedAt, v))
}

// DeletedAtIn applies the In predicate on the "deleted_at" field.
func DeletedAtIn(vs ...time.Time) predicate.Commission {
	return predicate.Commission(sql.FieldIn(FieldDeletedAt, vs...))
}

// DeletedAtNotIn applies the NotIn predicate on the "deleted_at" field.
func DeletedAtNotIn(vs ...time.Time) predicate.Commission {
	return predicate.Commission(sql.FieldNotIn(FieldDeletedAt, vs...))
}

// DeletedAtGT applies the GT predicate on the "deleted_at" field.
func DeletedAtGT(v time.Time) predicate.Commission {
	return predicate.Commission(sql.FieldGT(FieldDeletedAt, v))
}

// DeletedAtGTE applies the GTE predicate on the "deleted_at" field.
func DeletedAtGTE(v time.Time) predicate.Commission {
	return predicate.Commission(sql.FieldGTE(FieldDeletedAt, v))
}

// DeletedAtLT applies the LT predicate on the "deleted_at" field.
func DeletedAtLT(v time.Time) predicate.Commission {
	return predicate.Commission(sql.FieldLT(FieldDeletedAt, v))
}

// DeletedAtLTE applies the LTE predicate on the "deleted_at" field.
func DeletedAtLTE(v time.Time) predicate.Commission {
	return predicate.Commission(sql.FieldLTE(FieldDeletedAt, v))
}

// DeletedAtIsNil applies the IsNil predicate on the "deleted_at" field.
func DeletedAtIsNil() predicate.Commission {
	return predicate.Commission(sql.FieldIsNull(FieldDeletedAt))
}

// DeletedAtNotNil applies the NotNil predicate on the "deleted_at" field.
func DeletedAtNotNil() predicate.Commission {
	return predicate.Commission(sql.FieldNotNull(FieldDeletedAt))
}

// CreatorIsNil applies the IsNil predicate on the "creator" field.
func CreatorIsNil() predicate.Commission {
	return predicate.Commission(sql.FieldIsNull(FieldCreator))
}

// CreatorNotNil applies the NotNil predicate on the "creator" field.
func CreatorNotNil() predicate.Commission {
	return predicate.Commission(sql.FieldNotNull(FieldCreator))
}

// LastModifierIsNil applies the IsNil predicate on the "last_modifier" field.
func LastModifierIsNil() predicate.Commission {
	return predicate.Commission(sql.FieldIsNull(FieldLastModifier))
}

// LastModifierNotNil applies the NotNil predicate on the "last_modifier" field.
func LastModifierNotNil() predicate.Commission {
	return predicate.Commission(sql.FieldNotNull(FieldLastModifier))
}

// RemarkEQ applies the EQ predicate on the "remark" field.
func RemarkEQ(v string) predicate.Commission {
	return predicate.Commission(sql.FieldEQ(FieldRemark, v))
}

// RemarkNEQ applies the NEQ predicate on the "remark" field.
func RemarkNEQ(v string) predicate.Commission {
	return predicate.Commission(sql.FieldNEQ(FieldRemark, v))
}

// RemarkIn applies the In predicate on the "remark" field.
func RemarkIn(vs ...string) predicate.Commission {
	return predicate.Commission(sql.FieldIn(FieldRemark, vs...))
}

// RemarkNotIn applies the NotIn predicate on the "remark" field.
func RemarkNotIn(vs ...string) predicate.Commission {
	return predicate.Commission(sql.FieldNotIn(FieldRemark, vs...))
}

// RemarkGT applies the GT predicate on the "remark" field.
func RemarkGT(v string) predicate.Commission {
	return predicate.Commission(sql.FieldGT(FieldRemark, v))
}

// RemarkGTE applies the GTE predicate on the "remark" field.
func RemarkGTE(v string) predicate.Commission {
	return predicate.Commission(sql.FieldGTE(FieldRemark, v))
}

// RemarkLT applies the LT predicate on the "remark" field.
func RemarkLT(v string) predicate.Commission {
	return predicate.Commission(sql.FieldLT(FieldRemark, v))
}

// RemarkLTE applies the LTE predicate on the "remark" field.
func RemarkLTE(v string) predicate.Commission {
	return predicate.Commission(sql.FieldLTE(FieldRemark, v))
}

// RemarkContains applies the Contains predicate on the "remark" field.
func RemarkContains(v string) predicate.Commission {
	return predicate.Commission(sql.FieldContains(FieldRemark, v))
}

// RemarkHasPrefix applies the HasPrefix predicate on the "remark" field.
func RemarkHasPrefix(v string) predicate.Commission {
	return predicate.Commission(sql.FieldHasPrefix(FieldRemark, v))
}

// RemarkHasSuffix applies the HasSuffix predicate on the "remark" field.
func RemarkHasSuffix(v string) predicate.Commission {
	return predicate.Commission(sql.FieldHasSuffix(FieldRemark, v))
}

// RemarkIsNil applies the IsNil predicate on the "remark" field.
func RemarkIsNil() predicate.Commission {
	return predicate.Commission(sql.FieldIsNull(FieldRemark))
}

// RemarkNotNil applies the NotNil predicate on the "remark" field.
func RemarkNotNil() predicate.Commission {
	return predicate.Commission(sql.FieldNotNull(FieldRemark))
}

// RemarkEqualFold applies the EqualFold predicate on the "remark" field.
func RemarkEqualFold(v string) predicate.Commission {
	return predicate.Commission(sql.FieldEqualFold(FieldRemark, v))
}

// RemarkContainsFold applies the ContainsFold predicate on the "remark" field.
func RemarkContainsFold(v string) predicate.Commission {
	return predicate.Commission(sql.FieldContainsFold(FieldRemark, v))
}

// BusinessIDEQ applies the EQ predicate on the "business_id" field.
func BusinessIDEQ(v uint64) predicate.Commission {
	return predicate.Commission(sql.FieldEQ(FieldBusinessID, v))
}

// BusinessIDNEQ applies the NEQ predicate on the "business_id" field.
func BusinessIDNEQ(v uint64) predicate.Commission {
	return predicate.Commission(sql.FieldNEQ(FieldBusinessID, v))
}

// BusinessIDIn applies the In predicate on the "business_id" field.
func BusinessIDIn(vs ...uint64) predicate.Commission {
	return predicate.Commission(sql.FieldIn(FieldBusinessID, vs...))
}

// BusinessIDNotIn applies the NotIn predicate on the "business_id" field.
func BusinessIDNotIn(vs ...uint64) predicate.Commission {
	return predicate.Commission(sql.FieldNotIn(FieldBusinessID, vs...))
}

// BusinessIDIsNil applies the IsNil predicate on the "business_id" field.
func BusinessIDIsNil() predicate.Commission {
	return predicate.Commission(sql.FieldIsNull(FieldBusinessID))
}

// BusinessIDNotNil applies the NotNil predicate on the "business_id" field.
func BusinessIDNotNil() predicate.Commission {
	return predicate.Commission(sql.FieldNotNull(FieldBusinessID))
}

// SubscribeIDEQ applies the EQ predicate on the "subscribe_id" field.
func SubscribeIDEQ(v uint64) predicate.Commission {
	return predicate.Commission(sql.FieldEQ(FieldSubscribeID, v))
}

// SubscribeIDNEQ applies the NEQ predicate on the "subscribe_id" field.
func SubscribeIDNEQ(v uint64) predicate.Commission {
	return predicate.Commission(sql.FieldNEQ(FieldSubscribeID, v))
}

// SubscribeIDIn applies the In predicate on the "subscribe_id" field.
func SubscribeIDIn(vs ...uint64) predicate.Commission {
	return predicate.Commission(sql.FieldIn(FieldSubscribeID, vs...))
}

// SubscribeIDNotIn applies the NotIn predicate on the "subscribe_id" field.
func SubscribeIDNotIn(vs ...uint64) predicate.Commission {
	return predicate.Commission(sql.FieldNotIn(FieldSubscribeID, vs...))
}

// SubscribeIDIsNil applies the IsNil predicate on the "subscribe_id" field.
func SubscribeIDIsNil() predicate.Commission {
	return predicate.Commission(sql.FieldIsNull(FieldSubscribeID))
}

// SubscribeIDNotNil applies the NotNil predicate on the "subscribe_id" field.
func SubscribeIDNotNil() predicate.Commission {
	return predicate.Commission(sql.FieldNotNull(FieldSubscribeID))
}

// PlanIDEQ applies the EQ predicate on the "plan_id" field.
func PlanIDEQ(v uint64) predicate.Commission {
	return predicate.Commission(sql.FieldEQ(FieldPlanID, v))
}

// PlanIDNEQ applies the NEQ predicate on the "plan_id" field.
func PlanIDNEQ(v uint64) predicate.Commission {
	return predicate.Commission(sql.FieldNEQ(FieldPlanID, v))
}

// PlanIDIn applies the In predicate on the "plan_id" field.
func PlanIDIn(vs ...uint64) predicate.Commission {
	return predicate.Commission(sql.FieldIn(FieldPlanID, vs...))
}

// PlanIDNotIn applies the NotIn predicate on the "plan_id" field.
func PlanIDNotIn(vs ...uint64) predicate.Commission {
	return predicate.Commission(sql.FieldNotIn(FieldPlanID, vs...))
}

// PlanIDIsNil applies the IsNil predicate on the "plan_id" field.
func PlanIDIsNil() predicate.Commission {
	return predicate.Commission(sql.FieldIsNull(FieldPlanID))
}

// PlanIDNotNil applies the NotNil predicate on the "plan_id" field.
func PlanIDNotNil() predicate.Commission {
	return predicate.Commission(sql.FieldNotNull(FieldPlanID))
}

// RiderIDEQ applies the EQ predicate on the "rider_id" field.
func RiderIDEQ(v uint64) predicate.Commission {
	return predicate.Commission(sql.FieldEQ(FieldRiderID, v))
}

// RiderIDNEQ applies the NEQ predicate on the "rider_id" field.
func RiderIDNEQ(v uint64) predicate.Commission {
	return predicate.Commission(sql.FieldNEQ(FieldRiderID, v))
}

// RiderIDIn applies the In predicate on the "rider_id" field.
func RiderIDIn(vs ...uint64) predicate.Commission {
	return predicate.Commission(sql.FieldIn(FieldRiderID, vs...))
}

// RiderIDNotIn applies the NotIn predicate on the "rider_id" field.
func RiderIDNotIn(vs ...uint64) predicate.Commission {
	return predicate.Commission(sql.FieldNotIn(FieldRiderID, vs...))
}

// RiderIDIsNil applies the IsNil predicate on the "rider_id" field.
func RiderIDIsNil() predicate.Commission {
	return predicate.Commission(sql.FieldIsNull(FieldRiderID))
}

// RiderIDNotNil applies the NotNil predicate on the "rider_id" field.
func RiderIDNotNil() predicate.Commission {
	return predicate.Commission(sql.FieldNotNull(FieldRiderID))
}

// OrderIDEQ applies the EQ predicate on the "order_id" field.
func OrderIDEQ(v uint64) predicate.Commission {
	return predicate.Commission(sql.FieldEQ(FieldOrderID, v))
}

// OrderIDNEQ applies the NEQ predicate on the "order_id" field.
func OrderIDNEQ(v uint64) predicate.Commission {
	return predicate.Commission(sql.FieldNEQ(FieldOrderID, v))
}

// OrderIDIn applies the In predicate on the "order_id" field.
func OrderIDIn(vs ...uint64) predicate.Commission {
	return predicate.Commission(sql.FieldIn(FieldOrderID, vs...))
}

// OrderIDNotIn applies the NotIn predicate on the "order_id" field.
func OrderIDNotIn(vs ...uint64) predicate.Commission {
	return predicate.Commission(sql.FieldNotIn(FieldOrderID, vs...))
}

// AmountEQ applies the EQ predicate on the "amount" field.
func AmountEQ(v float64) predicate.Commission {
	return predicate.Commission(sql.FieldEQ(FieldAmount, v))
}

// AmountNEQ applies the NEQ predicate on the "amount" field.
func AmountNEQ(v float64) predicate.Commission {
	return predicate.Commission(sql.FieldNEQ(FieldAmount, v))
}

// AmountIn applies the In predicate on the "amount" field.
func AmountIn(vs ...float64) predicate.Commission {
	return predicate.Commission(sql.FieldIn(FieldAmount, vs...))
}

// AmountNotIn applies the NotIn predicate on the "amount" field.
func AmountNotIn(vs ...float64) predicate.Commission {
	return predicate.Commission(sql.FieldNotIn(FieldAmount, vs...))
}

// AmountGT applies the GT predicate on the "amount" field.
func AmountGT(v float64) predicate.Commission {
	return predicate.Commission(sql.FieldGT(FieldAmount, v))
}

// AmountGTE applies the GTE predicate on the "amount" field.
func AmountGTE(v float64) predicate.Commission {
	return predicate.Commission(sql.FieldGTE(FieldAmount, v))
}

// AmountLT applies the LT predicate on the "amount" field.
func AmountLT(v float64) predicate.Commission {
	return predicate.Commission(sql.FieldLT(FieldAmount, v))
}

// AmountLTE applies the LTE predicate on the "amount" field.
func AmountLTE(v float64) predicate.Commission {
	return predicate.Commission(sql.FieldLTE(FieldAmount, v))
}

// StatusEQ applies the EQ predicate on the "status" field.
func StatusEQ(v uint8) predicate.Commission {
	return predicate.Commission(sql.FieldEQ(FieldStatus, v))
}

// StatusNEQ applies the NEQ predicate on the "status" field.
func StatusNEQ(v uint8) predicate.Commission {
	return predicate.Commission(sql.FieldNEQ(FieldStatus, v))
}

// StatusIn applies the In predicate on the "status" field.
func StatusIn(vs ...uint8) predicate.Commission {
	return predicate.Commission(sql.FieldIn(FieldStatus, vs...))
}

// StatusNotIn applies the NotIn predicate on the "status" field.
func StatusNotIn(vs ...uint8) predicate.Commission {
	return predicate.Commission(sql.FieldNotIn(FieldStatus, vs...))
}

// StatusGT applies the GT predicate on the "status" field.
func StatusGT(v uint8) predicate.Commission {
	return predicate.Commission(sql.FieldGT(FieldStatus, v))
}

// StatusGTE applies the GTE predicate on the "status" field.
func StatusGTE(v uint8) predicate.Commission {
	return predicate.Commission(sql.FieldGTE(FieldStatus, v))
}

// StatusLT applies the LT predicate on the "status" field.
func StatusLT(v uint8) predicate.Commission {
	return predicate.Commission(sql.FieldLT(FieldStatus, v))
}

// StatusLTE applies the LTE predicate on the "status" field.
func StatusLTE(v uint8) predicate.Commission {
	return predicate.Commission(sql.FieldLTE(FieldStatus, v))
}

// EmployeeIDEQ applies the EQ predicate on the "employee_id" field.
func EmployeeIDEQ(v uint64) predicate.Commission {
	return predicate.Commission(sql.FieldEQ(FieldEmployeeID, v))
}

// EmployeeIDNEQ applies the NEQ predicate on the "employee_id" field.
func EmployeeIDNEQ(v uint64) predicate.Commission {
	return predicate.Commission(sql.FieldNEQ(FieldEmployeeID, v))
}

// EmployeeIDIn applies the In predicate on the "employee_id" field.
func EmployeeIDIn(vs ...uint64) predicate.Commission {
	return predicate.Commission(sql.FieldIn(FieldEmployeeID, vs...))
}

// EmployeeIDNotIn applies the NotIn predicate on the "employee_id" field.
func EmployeeIDNotIn(vs ...uint64) predicate.Commission {
	return predicate.Commission(sql.FieldNotIn(FieldEmployeeID, vs...))
}

// EmployeeIDIsNil applies the IsNil predicate on the "employee_id" field.
func EmployeeIDIsNil() predicate.Commission {
	return predicate.Commission(sql.FieldIsNull(FieldEmployeeID))
}

// EmployeeIDNotNil applies the NotNil predicate on the "employee_id" field.
func EmployeeIDNotNil() predicate.Commission {
	return predicate.Commission(sql.FieldNotNull(FieldEmployeeID))
}

// HasBusiness applies the HasEdge predicate on the "business" edge.
func HasBusiness() predicate.Commission {
	return predicate.Commission(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, BusinessTable, BusinessColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasBusinessWith applies the HasEdge predicate on the "business" edge with a given conditions (other predicates).
func HasBusinessWith(preds ...predicate.Business) predicate.Commission {
	return predicate.Commission(func(s *sql.Selector) {
		step := newBusinessStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasSubscribe applies the HasEdge predicate on the "subscribe" edge.
func HasSubscribe() predicate.Commission {
	return predicate.Commission(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, SubscribeTable, SubscribeColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasSubscribeWith applies the HasEdge predicate on the "subscribe" edge with a given conditions (other predicates).
func HasSubscribeWith(preds ...predicate.Subscribe) predicate.Commission {
	return predicate.Commission(func(s *sql.Selector) {
		step := newSubscribeStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasPlan applies the HasEdge predicate on the "plan" edge.
func HasPlan() predicate.Commission {
	return predicate.Commission(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, PlanTable, PlanColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasPlanWith applies the HasEdge predicate on the "plan" edge with a given conditions (other predicates).
func HasPlanWith(preds ...predicate.Plan) predicate.Commission {
	return predicate.Commission(func(s *sql.Selector) {
		step := newPlanStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasRider applies the HasEdge predicate on the "rider" edge.
func HasRider() predicate.Commission {
	return predicate.Commission(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, RiderTable, RiderColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasRiderWith applies the HasEdge predicate on the "rider" edge with a given conditions (other predicates).
func HasRiderWith(preds ...predicate.Rider) predicate.Commission {
	return predicate.Commission(func(s *sql.Selector) {
		step := newRiderStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasOrder applies the HasEdge predicate on the "order" edge.
func HasOrder() predicate.Commission {
	return predicate.Commission(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.O2O, true, OrderTable, OrderColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasOrderWith applies the HasEdge predicate on the "order" edge with a given conditions (other predicates).
func HasOrderWith(preds ...predicate.Order) predicate.Commission {
	return predicate.Commission(func(s *sql.Selector) {
		step := newOrderStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasEmployee applies the HasEdge predicate on the "employee" edge.
func HasEmployee() predicate.Commission {
	return predicate.Commission(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, EmployeeTable, EmployeeColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasEmployeeWith applies the HasEdge predicate on the "employee" edge with a given conditions (other predicates).
func HasEmployeeWith(preds ...predicate.Employee) predicate.Commission {
	return predicate.Commission(func(s *sql.Selector) {
		step := newEmployeeStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.Commission) predicate.Commission {
	return predicate.Commission(sql.AndPredicates(predicates...))
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.Commission) predicate.Commission {
	return predicate.Commission(sql.OrPredicates(predicates...))
}

// Not applies the not operator on the given predicate.
func Not(p predicate.Commission) predicate.Commission {
	return predicate.Commission(sql.NotPredicates(p))
}
