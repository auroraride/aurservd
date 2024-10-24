// Code generated by ent, DO NOT EDIT.

package contract

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/auroraride/aurservd/internal/ent/predicate"
)

// ID filters vertices based on their ID field.
func ID(id uint64) predicate.Contract {
	return predicate.Contract(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id uint64) predicate.Contract {
	return predicate.Contract(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id uint64) predicate.Contract {
	return predicate.Contract(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...uint64) predicate.Contract {
	return predicate.Contract(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...uint64) predicate.Contract {
	return predicate.Contract(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id uint64) predicate.Contract {
	return predicate.Contract(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id uint64) predicate.Contract {
	return predicate.Contract(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id uint64) predicate.Contract {
	return predicate.Contract(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id uint64) predicate.Contract {
	return predicate.Contract(sql.FieldLTE(FieldID, id))
}

// CreatedAt applies equality check predicate on the "created_at" field. It's identical to CreatedAtEQ.
func CreatedAt(v time.Time) predicate.Contract {
	return predicate.Contract(sql.FieldEQ(FieldCreatedAt, v))
}

// UpdatedAt applies equality check predicate on the "updated_at" field. It's identical to UpdatedAtEQ.
func UpdatedAt(v time.Time) predicate.Contract {
	return predicate.Contract(sql.FieldEQ(FieldUpdatedAt, v))
}

// DeletedAt applies equality check predicate on the "deleted_at" field. It's identical to DeletedAtEQ.
func DeletedAt(v time.Time) predicate.Contract {
	return predicate.Contract(sql.FieldEQ(FieldDeletedAt, v))
}

// Remark applies equality check predicate on the "remark" field. It's identical to RemarkEQ.
func Remark(v string) predicate.Contract {
	return predicate.Contract(sql.FieldEQ(FieldRemark, v))
}

// SubscribeID applies equality check predicate on the "subscribe_id" field. It's identical to SubscribeIDEQ.
func SubscribeID(v uint64) predicate.Contract {
	return predicate.Contract(sql.FieldEQ(FieldSubscribeID, v))
}

// EmployeeID applies equality check predicate on the "employee_id" field. It's identical to EmployeeIDEQ.
func EmployeeID(v uint64) predicate.Contract {
	return predicate.Contract(sql.FieldEQ(FieldEmployeeID, v))
}

// Status applies equality check predicate on the "status" field. It's identical to StatusEQ.
func Status(v uint8) predicate.Contract {
	return predicate.Contract(sql.FieldEQ(FieldStatus, v))
}

// RiderID applies equality check predicate on the "rider_id" field. It's identical to RiderIDEQ.
func RiderID(v uint64) predicate.Contract {
	return predicate.Contract(sql.FieldEQ(FieldRiderID, v))
}

// FlowID applies equality check predicate on the "flow_id" field. It's identical to FlowIDEQ.
func FlowID(v string) predicate.Contract {
	return predicate.Contract(sql.FieldEQ(FieldFlowID, v))
}

// Sn applies equality check predicate on the "sn" field. It's identical to SnEQ.
func Sn(v string) predicate.Contract {
	return predicate.Contract(sql.FieldEQ(FieldSn, v))
}

// Effective applies equality check predicate on the "effective" field. It's identical to EffectiveEQ.
func Effective(v bool) predicate.Contract {
	return predicate.Contract(sql.FieldEQ(FieldEffective, v))
}

// AllocateID applies equality check predicate on the "allocate_id" field. It's identical to AllocateIDEQ.
func AllocateID(v uint64) predicate.Contract {
	return predicate.Contract(sql.FieldEQ(FieldAllocateID, v))
}

// Link applies equality check predicate on the "link" field. It's identical to LinkEQ.
func Link(v string) predicate.Contract {
	return predicate.Contract(sql.FieldEQ(FieldLink, v))
}

// ExpiresAt applies equality check predicate on the "expires_at" field. It's identical to ExpiresAtEQ.
func ExpiresAt(v time.Time) predicate.Contract {
	return predicate.Contract(sql.FieldEQ(FieldExpiresAt, v))
}

// SignedAt applies equality check predicate on the "signed_at" field. It's identical to SignedAtEQ.
func SignedAt(v time.Time) predicate.Contract {
	return predicate.Contract(sql.FieldEQ(FieldSignedAt, v))
}

// DocID applies equality check predicate on the "doc_id" field. It's identical to DocIDEQ.
func DocID(v string) predicate.Contract {
	return predicate.Contract(sql.FieldEQ(FieldDocID, v))
}

// CreatedAtEQ applies the EQ predicate on the "created_at" field.
func CreatedAtEQ(v time.Time) predicate.Contract {
	return predicate.Contract(sql.FieldEQ(FieldCreatedAt, v))
}

// CreatedAtNEQ applies the NEQ predicate on the "created_at" field.
func CreatedAtNEQ(v time.Time) predicate.Contract {
	return predicate.Contract(sql.FieldNEQ(FieldCreatedAt, v))
}

// CreatedAtIn applies the In predicate on the "created_at" field.
func CreatedAtIn(vs ...time.Time) predicate.Contract {
	return predicate.Contract(sql.FieldIn(FieldCreatedAt, vs...))
}

// CreatedAtNotIn applies the NotIn predicate on the "created_at" field.
func CreatedAtNotIn(vs ...time.Time) predicate.Contract {
	return predicate.Contract(sql.FieldNotIn(FieldCreatedAt, vs...))
}

// CreatedAtGT applies the GT predicate on the "created_at" field.
func CreatedAtGT(v time.Time) predicate.Contract {
	return predicate.Contract(sql.FieldGT(FieldCreatedAt, v))
}

// CreatedAtGTE applies the GTE predicate on the "created_at" field.
func CreatedAtGTE(v time.Time) predicate.Contract {
	return predicate.Contract(sql.FieldGTE(FieldCreatedAt, v))
}

// CreatedAtLT applies the LT predicate on the "created_at" field.
func CreatedAtLT(v time.Time) predicate.Contract {
	return predicate.Contract(sql.FieldLT(FieldCreatedAt, v))
}

// CreatedAtLTE applies the LTE predicate on the "created_at" field.
func CreatedAtLTE(v time.Time) predicate.Contract {
	return predicate.Contract(sql.FieldLTE(FieldCreatedAt, v))
}

// UpdatedAtEQ applies the EQ predicate on the "updated_at" field.
func UpdatedAtEQ(v time.Time) predicate.Contract {
	return predicate.Contract(sql.FieldEQ(FieldUpdatedAt, v))
}

// UpdatedAtNEQ applies the NEQ predicate on the "updated_at" field.
func UpdatedAtNEQ(v time.Time) predicate.Contract {
	return predicate.Contract(sql.FieldNEQ(FieldUpdatedAt, v))
}

// UpdatedAtIn applies the In predicate on the "updated_at" field.
func UpdatedAtIn(vs ...time.Time) predicate.Contract {
	return predicate.Contract(sql.FieldIn(FieldUpdatedAt, vs...))
}

// UpdatedAtNotIn applies the NotIn predicate on the "updated_at" field.
func UpdatedAtNotIn(vs ...time.Time) predicate.Contract {
	return predicate.Contract(sql.FieldNotIn(FieldUpdatedAt, vs...))
}

// UpdatedAtGT applies the GT predicate on the "updated_at" field.
func UpdatedAtGT(v time.Time) predicate.Contract {
	return predicate.Contract(sql.FieldGT(FieldUpdatedAt, v))
}

// UpdatedAtGTE applies the GTE predicate on the "updated_at" field.
func UpdatedAtGTE(v time.Time) predicate.Contract {
	return predicate.Contract(sql.FieldGTE(FieldUpdatedAt, v))
}

// UpdatedAtLT applies the LT predicate on the "updated_at" field.
func UpdatedAtLT(v time.Time) predicate.Contract {
	return predicate.Contract(sql.FieldLT(FieldUpdatedAt, v))
}

// UpdatedAtLTE applies the LTE predicate on the "updated_at" field.
func UpdatedAtLTE(v time.Time) predicate.Contract {
	return predicate.Contract(sql.FieldLTE(FieldUpdatedAt, v))
}

// DeletedAtEQ applies the EQ predicate on the "deleted_at" field.
func DeletedAtEQ(v time.Time) predicate.Contract {
	return predicate.Contract(sql.FieldEQ(FieldDeletedAt, v))
}

// DeletedAtNEQ applies the NEQ predicate on the "deleted_at" field.
func DeletedAtNEQ(v time.Time) predicate.Contract {
	return predicate.Contract(sql.FieldNEQ(FieldDeletedAt, v))
}

// DeletedAtIn applies the In predicate on the "deleted_at" field.
func DeletedAtIn(vs ...time.Time) predicate.Contract {
	return predicate.Contract(sql.FieldIn(FieldDeletedAt, vs...))
}

// DeletedAtNotIn applies the NotIn predicate on the "deleted_at" field.
func DeletedAtNotIn(vs ...time.Time) predicate.Contract {
	return predicate.Contract(sql.FieldNotIn(FieldDeletedAt, vs...))
}

// DeletedAtGT applies the GT predicate on the "deleted_at" field.
func DeletedAtGT(v time.Time) predicate.Contract {
	return predicate.Contract(sql.FieldGT(FieldDeletedAt, v))
}

// DeletedAtGTE applies the GTE predicate on the "deleted_at" field.
func DeletedAtGTE(v time.Time) predicate.Contract {
	return predicate.Contract(sql.FieldGTE(FieldDeletedAt, v))
}

// DeletedAtLT applies the LT predicate on the "deleted_at" field.
func DeletedAtLT(v time.Time) predicate.Contract {
	return predicate.Contract(sql.FieldLT(FieldDeletedAt, v))
}

// DeletedAtLTE applies the LTE predicate on the "deleted_at" field.
func DeletedAtLTE(v time.Time) predicate.Contract {
	return predicate.Contract(sql.FieldLTE(FieldDeletedAt, v))
}

// DeletedAtIsNil applies the IsNil predicate on the "deleted_at" field.
func DeletedAtIsNil() predicate.Contract {
	return predicate.Contract(sql.FieldIsNull(FieldDeletedAt))
}

// DeletedAtNotNil applies the NotNil predicate on the "deleted_at" field.
func DeletedAtNotNil() predicate.Contract {
	return predicate.Contract(sql.FieldNotNull(FieldDeletedAt))
}

// CreatorIsNil applies the IsNil predicate on the "creator" field.
func CreatorIsNil() predicate.Contract {
	return predicate.Contract(sql.FieldIsNull(FieldCreator))
}

// CreatorNotNil applies the NotNil predicate on the "creator" field.
func CreatorNotNil() predicate.Contract {
	return predicate.Contract(sql.FieldNotNull(FieldCreator))
}

// LastModifierIsNil applies the IsNil predicate on the "last_modifier" field.
func LastModifierIsNil() predicate.Contract {
	return predicate.Contract(sql.FieldIsNull(FieldLastModifier))
}

// LastModifierNotNil applies the NotNil predicate on the "last_modifier" field.
func LastModifierNotNil() predicate.Contract {
	return predicate.Contract(sql.FieldNotNull(FieldLastModifier))
}

// RemarkEQ applies the EQ predicate on the "remark" field.
func RemarkEQ(v string) predicate.Contract {
	return predicate.Contract(sql.FieldEQ(FieldRemark, v))
}

// RemarkNEQ applies the NEQ predicate on the "remark" field.
func RemarkNEQ(v string) predicate.Contract {
	return predicate.Contract(sql.FieldNEQ(FieldRemark, v))
}

// RemarkIn applies the In predicate on the "remark" field.
func RemarkIn(vs ...string) predicate.Contract {
	return predicate.Contract(sql.FieldIn(FieldRemark, vs...))
}

// RemarkNotIn applies the NotIn predicate on the "remark" field.
func RemarkNotIn(vs ...string) predicate.Contract {
	return predicate.Contract(sql.FieldNotIn(FieldRemark, vs...))
}

// RemarkGT applies the GT predicate on the "remark" field.
func RemarkGT(v string) predicate.Contract {
	return predicate.Contract(sql.FieldGT(FieldRemark, v))
}

// RemarkGTE applies the GTE predicate on the "remark" field.
func RemarkGTE(v string) predicate.Contract {
	return predicate.Contract(sql.FieldGTE(FieldRemark, v))
}

// RemarkLT applies the LT predicate on the "remark" field.
func RemarkLT(v string) predicate.Contract {
	return predicate.Contract(sql.FieldLT(FieldRemark, v))
}

// RemarkLTE applies the LTE predicate on the "remark" field.
func RemarkLTE(v string) predicate.Contract {
	return predicate.Contract(sql.FieldLTE(FieldRemark, v))
}

// RemarkContains applies the Contains predicate on the "remark" field.
func RemarkContains(v string) predicate.Contract {
	return predicate.Contract(sql.FieldContains(FieldRemark, v))
}

// RemarkHasPrefix applies the HasPrefix predicate on the "remark" field.
func RemarkHasPrefix(v string) predicate.Contract {
	return predicate.Contract(sql.FieldHasPrefix(FieldRemark, v))
}

// RemarkHasSuffix applies the HasSuffix predicate on the "remark" field.
func RemarkHasSuffix(v string) predicate.Contract {
	return predicate.Contract(sql.FieldHasSuffix(FieldRemark, v))
}

// RemarkIsNil applies the IsNil predicate on the "remark" field.
func RemarkIsNil() predicate.Contract {
	return predicate.Contract(sql.FieldIsNull(FieldRemark))
}

// RemarkNotNil applies the NotNil predicate on the "remark" field.
func RemarkNotNil() predicate.Contract {
	return predicate.Contract(sql.FieldNotNull(FieldRemark))
}

// RemarkEqualFold applies the EqualFold predicate on the "remark" field.
func RemarkEqualFold(v string) predicate.Contract {
	return predicate.Contract(sql.FieldEqualFold(FieldRemark, v))
}

// RemarkContainsFold applies the ContainsFold predicate on the "remark" field.
func RemarkContainsFold(v string) predicate.Contract {
	return predicate.Contract(sql.FieldContainsFold(FieldRemark, v))
}

// SubscribeIDEQ applies the EQ predicate on the "subscribe_id" field.
func SubscribeIDEQ(v uint64) predicate.Contract {
	return predicate.Contract(sql.FieldEQ(FieldSubscribeID, v))
}

// SubscribeIDNEQ applies the NEQ predicate on the "subscribe_id" field.
func SubscribeIDNEQ(v uint64) predicate.Contract {
	return predicate.Contract(sql.FieldNEQ(FieldSubscribeID, v))
}

// SubscribeIDIn applies the In predicate on the "subscribe_id" field.
func SubscribeIDIn(vs ...uint64) predicate.Contract {
	return predicate.Contract(sql.FieldIn(FieldSubscribeID, vs...))
}

// SubscribeIDNotIn applies the NotIn predicate on the "subscribe_id" field.
func SubscribeIDNotIn(vs ...uint64) predicate.Contract {
	return predicate.Contract(sql.FieldNotIn(FieldSubscribeID, vs...))
}

// SubscribeIDIsNil applies the IsNil predicate on the "subscribe_id" field.
func SubscribeIDIsNil() predicate.Contract {
	return predicate.Contract(sql.FieldIsNull(FieldSubscribeID))
}

// SubscribeIDNotNil applies the NotNil predicate on the "subscribe_id" field.
func SubscribeIDNotNil() predicate.Contract {
	return predicate.Contract(sql.FieldNotNull(FieldSubscribeID))
}

// EmployeeIDEQ applies the EQ predicate on the "employee_id" field.
func EmployeeIDEQ(v uint64) predicate.Contract {
	return predicate.Contract(sql.FieldEQ(FieldEmployeeID, v))
}

// EmployeeIDNEQ applies the NEQ predicate on the "employee_id" field.
func EmployeeIDNEQ(v uint64) predicate.Contract {
	return predicate.Contract(sql.FieldNEQ(FieldEmployeeID, v))
}

// EmployeeIDIn applies the In predicate on the "employee_id" field.
func EmployeeIDIn(vs ...uint64) predicate.Contract {
	return predicate.Contract(sql.FieldIn(FieldEmployeeID, vs...))
}

// EmployeeIDNotIn applies the NotIn predicate on the "employee_id" field.
func EmployeeIDNotIn(vs ...uint64) predicate.Contract {
	return predicate.Contract(sql.FieldNotIn(FieldEmployeeID, vs...))
}

// EmployeeIDIsNil applies the IsNil predicate on the "employee_id" field.
func EmployeeIDIsNil() predicate.Contract {
	return predicate.Contract(sql.FieldIsNull(FieldEmployeeID))
}

// EmployeeIDNotNil applies the NotNil predicate on the "employee_id" field.
func EmployeeIDNotNil() predicate.Contract {
	return predicate.Contract(sql.FieldNotNull(FieldEmployeeID))
}

// StatusEQ applies the EQ predicate on the "status" field.
func StatusEQ(v uint8) predicate.Contract {
	return predicate.Contract(sql.FieldEQ(FieldStatus, v))
}

// StatusNEQ applies the NEQ predicate on the "status" field.
func StatusNEQ(v uint8) predicate.Contract {
	return predicate.Contract(sql.FieldNEQ(FieldStatus, v))
}

// StatusIn applies the In predicate on the "status" field.
func StatusIn(vs ...uint8) predicate.Contract {
	return predicate.Contract(sql.FieldIn(FieldStatus, vs...))
}

// StatusNotIn applies the NotIn predicate on the "status" field.
func StatusNotIn(vs ...uint8) predicate.Contract {
	return predicate.Contract(sql.FieldNotIn(FieldStatus, vs...))
}

// StatusGT applies the GT predicate on the "status" field.
func StatusGT(v uint8) predicate.Contract {
	return predicate.Contract(sql.FieldGT(FieldStatus, v))
}

// StatusGTE applies the GTE predicate on the "status" field.
func StatusGTE(v uint8) predicate.Contract {
	return predicate.Contract(sql.FieldGTE(FieldStatus, v))
}

// StatusLT applies the LT predicate on the "status" field.
func StatusLT(v uint8) predicate.Contract {
	return predicate.Contract(sql.FieldLT(FieldStatus, v))
}

// StatusLTE applies the LTE predicate on the "status" field.
func StatusLTE(v uint8) predicate.Contract {
	return predicate.Contract(sql.FieldLTE(FieldStatus, v))
}

// RiderIDEQ applies the EQ predicate on the "rider_id" field.
func RiderIDEQ(v uint64) predicate.Contract {
	return predicate.Contract(sql.FieldEQ(FieldRiderID, v))
}

// RiderIDNEQ applies the NEQ predicate on the "rider_id" field.
func RiderIDNEQ(v uint64) predicate.Contract {
	return predicate.Contract(sql.FieldNEQ(FieldRiderID, v))
}

// RiderIDIn applies the In predicate on the "rider_id" field.
func RiderIDIn(vs ...uint64) predicate.Contract {
	return predicate.Contract(sql.FieldIn(FieldRiderID, vs...))
}

// RiderIDNotIn applies the NotIn predicate on the "rider_id" field.
func RiderIDNotIn(vs ...uint64) predicate.Contract {
	return predicate.Contract(sql.FieldNotIn(FieldRiderID, vs...))
}

// FlowIDEQ applies the EQ predicate on the "flow_id" field.
func FlowIDEQ(v string) predicate.Contract {
	return predicate.Contract(sql.FieldEQ(FieldFlowID, v))
}

// FlowIDNEQ applies the NEQ predicate on the "flow_id" field.
func FlowIDNEQ(v string) predicate.Contract {
	return predicate.Contract(sql.FieldNEQ(FieldFlowID, v))
}

// FlowIDIn applies the In predicate on the "flow_id" field.
func FlowIDIn(vs ...string) predicate.Contract {
	return predicate.Contract(sql.FieldIn(FieldFlowID, vs...))
}

// FlowIDNotIn applies the NotIn predicate on the "flow_id" field.
func FlowIDNotIn(vs ...string) predicate.Contract {
	return predicate.Contract(sql.FieldNotIn(FieldFlowID, vs...))
}

// FlowIDGT applies the GT predicate on the "flow_id" field.
func FlowIDGT(v string) predicate.Contract {
	return predicate.Contract(sql.FieldGT(FieldFlowID, v))
}

// FlowIDGTE applies the GTE predicate on the "flow_id" field.
func FlowIDGTE(v string) predicate.Contract {
	return predicate.Contract(sql.FieldGTE(FieldFlowID, v))
}

// FlowIDLT applies the LT predicate on the "flow_id" field.
func FlowIDLT(v string) predicate.Contract {
	return predicate.Contract(sql.FieldLT(FieldFlowID, v))
}

// FlowIDLTE applies the LTE predicate on the "flow_id" field.
func FlowIDLTE(v string) predicate.Contract {
	return predicate.Contract(sql.FieldLTE(FieldFlowID, v))
}

// FlowIDContains applies the Contains predicate on the "flow_id" field.
func FlowIDContains(v string) predicate.Contract {
	return predicate.Contract(sql.FieldContains(FieldFlowID, v))
}

// FlowIDHasPrefix applies the HasPrefix predicate on the "flow_id" field.
func FlowIDHasPrefix(v string) predicate.Contract {
	return predicate.Contract(sql.FieldHasPrefix(FieldFlowID, v))
}

// FlowIDHasSuffix applies the HasSuffix predicate on the "flow_id" field.
func FlowIDHasSuffix(v string) predicate.Contract {
	return predicate.Contract(sql.FieldHasSuffix(FieldFlowID, v))
}

// FlowIDEqualFold applies the EqualFold predicate on the "flow_id" field.
func FlowIDEqualFold(v string) predicate.Contract {
	return predicate.Contract(sql.FieldEqualFold(FieldFlowID, v))
}

// FlowIDContainsFold applies the ContainsFold predicate on the "flow_id" field.
func FlowIDContainsFold(v string) predicate.Contract {
	return predicate.Contract(sql.FieldContainsFold(FieldFlowID, v))
}

// SnEQ applies the EQ predicate on the "sn" field.
func SnEQ(v string) predicate.Contract {
	return predicate.Contract(sql.FieldEQ(FieldSn, v))
}

// SnNEQ applies the NEQ predicate on the "sn" field.
func SnNEQ(v string) predicate.Contract {
	return predicate.Contract(sql.FieldNEQ(FieldSn, v))
}

// SnIn applies the In predicate on the "sn" field.
func SnIn(vs ...string) predicate.Contract {
	return predicate.Contract(sql.FieldIn(FieldSn, vs...))
}

// SnNotIn applies the NotIn predicate on the "sn" field.
func SnNotIn(vs ...string) predicate.Contract {
	return predicate.Contract(sql.FieldNotIn(FieldSn, vs...))
}

// SnGT applies the GT predicate on the "sn" field.
func SnGT(v string) predicate.Contract {
	return predicate.Contract(sql.FieldGT(FieldSn, v))
}

// SnGTE applies the GTE predicate on the "sn" field.
func SnGTE(v string) predicate.Contract {
	return predicate.Contract(sql.FieldGTE(FieldSn, v))
}

// SnLT applies the LT predicate on the "sn" field.
func SnLT(v string) predicate.Contract {
	return predicate.Contract(sql.FieldLT(FieldSn, v))
}

// SnLTE applies the LTE predicate on the "sn" field.
func SnLTE(v string) predicate.Contract {
	return predicate.Contract(sql.FieldLTE(FieldSn, v))
}

// SnContains applies the Contains predicate on the "sn" field.
func SnContains(v string) predicate.Contract {
	return predicate.Contract(sql.FieldContains(FieldSn, v))
}

// SnHasPrefix applies the HasPrefix predicate on the "sn" field.
func SnHasPrefix(v string) predicate.Contract {
	return predicate.Contract(sql.FieldHasPrefix(FieldSn, v))
}

// SnHasSuffix applies the HasSuffix predicate on the "sn" field.
func SnHasSuffix(v string) predicate.Contract {
	return predicate.Contract(sql.FieldHasSuffix(FieldSn, v))
}

// SnEqualFold applies the EqualFold predicate on the "sn" field.
func SnEqualFold(v string) predicate.Contract {
	return predicate.Contract(sql.FieldEqualFold(FieldSn, v))
}

// SnContainsFold applies the ContainsFold predicate on the "sn" field.
func SnContainsFold(v string) predicate.Contract {
	return predicate.Contract(sql.FieldContainsFold(FieldSn, v))
}

// FilesIsNil applies the IsNil predicate on the "files" field.
func FilesIsNil() predicate.Contract {
	return predicate.Contract(sql.FieldIsNull(FieldFiles))
}

// FilesNotNil applies the NotNil predicate on the "files" field.
func FilesNotNil() predicate.Contract {
	return predicate.Contract(sql.FieldNotNull(FieldFiles))
}

// EffectiveEQ applies the EQ predicate on the "effective" field.
func EffectiveEQ(v bool) predicate.Contract {
	return predicate.Contract(sql.FieldEQ(FieldEffective, v))
}

// EffectiveNEQ applies the NEQ predicate on the "effective" field.
func EffectiveNEQ(v bool) predicate.Contract {
	return predicate.Contract(sql.FieldNEQ(FieldEffective, v))
}

// RiderInfoIsNil applies the IsNil predicate on the "rider_info" field.
func RiderInfoIsNil() predicate.Contract {
	return predicate.Contract(sql.FieldIsNull(FieldRiderInfo))
}

// RiderInfoNotNil applies the NotNil predicate on the "rider_info" field.
func RiderInfoNotNil() predicate.Contract {
	return predicate.Contract(sql.FieldNotNull(FieldRiderInfo))
}

// AllocateIDEQ applies the EQ predicate on the "allocate_id" field.
func AllocateIDEQ(v uint64) predicate.Contract {
	return predicate.Contract(sql.FieldEQ(FieldAllocateID, v))
}

// AllocateIDNEQ applies the NEQ predicate on the "allocate_id" field.
func AllocateIDNEQ(v uint64) predicate.Contract {
	return predicate.Contract(sql.FieldNEQ(FieldAllocateID, v))
}

// AllocateIDIn applies the In predicate on the "allocate_id" field.
func AllocateIDIn(vs ...uint64) predicate.Contract {
	return predicate.Contract(sql.FieldIn(FieldAllocateID, vs...))
}

// AllocateIDNotIn applies the NotIn predicate on the "allocate_id" field.
func AllocateIDNotIn(vs ...uint64) predicate.Contract {
	return predicate.Contract(sql.FieldNotIn(FieldAllocateID, vs...))
}

// AllocateIDIsNil applies the IsNil predicate on the "allocate_id" field.
func AllocateIDIsNil() predicate.Contract {
	return predicate.Contract(sql.FieldIsNull(FieldAllocateID))
}

// AllocateIDNotNil applies the NotNil predicate on the "allocate_id" field.
func AllocateIDNotNil() predicate.Contract {
	return predicate.Contract(sql.FieldNotNull(FieldAllocateID))
}

// LinkEQ applies the EQ predicate on the "link" field.
func LinkEQ(v string) predicate.Contract {
	return predicate.Contract(sql.FieldEQ(FieldLink, v))
}

// LinkNEQ applies the NEQ predicate on the "link" field.
func LinkNEQ(v string) predicate.Contract {
	return predicate.Contract(sql.FieldNEQ(FieldLink, v))
}

// LinkIn applies the In predicate on the "link" field.
func LinkIn(vs ...string) predicate.Contract {
	return predicate.Contract(sql.FieldIn(FieldLink, vs...))
}

// LinkNotIn applies the NotIn predicate on the "link" field.
func LinkNotIn(vs ...string) predicate.Contract {
	return predicate.Contract(sql.FieldNotIn(FieldLink, vs...))
}

// LinkGT applies the GT predicate on the "link" field.
func LinkGT(v string) predicate.Contract {
	return predicate.Contract(sql.FieldGT(FieldLink, v))
}

// LinkGTE applies the GTE predicate on the "link" field.
func LinkGTE(v string) predicate.Contract {
	return predicate.Contract(sql.FieldGTE(FieldLink, v))
}

// LinkLT applies the LT predicate on the "link" field.
func LinkLT(v string) predicate.Contract {
	return predicate.Contract(sql.FieldLT(FieldLink, v))
}

// LinkLTE applies the LTE predicate on the "link" field.
func LinkLTE(v string) predicate.Contract {
	return predicate.Contract(sql.FieldLTE(FieldLink, v))
}

// LinkContains applies the Contains predicate on the "link" field.
func LinkContains(v string) predicate.Contract {
	return predicate.Contract(sql.FieldContains(FieldLink, v))
}

// LinkHasPrefix applies the HasPrefix predicate on the "link" field.
func LinkHasPrefix(v string) predicate.Contract {
	return predicate.Contract(sql.FieldHasPrefix(FieldLink, v))
}

// LinkHasSuffix applies the HasSuffix predicate on the "link" field.
func LinkHasSuffix(v string) predicate.Contract {
	return predicate.Contract(sql.FieldHasSuffix(FieldLink, v))
}

// LinkIsNil applies the IsNil predicate on the "link" field.
func LinkIsNil() predicate.Contract {
	return predicate.Contract(sql.FieldIsNull(FieldLink))
}

// LinkNotNil applies the NotNil predicate on the "link" field.
func LinkNotNil() predicate.Contract {
	return predicate.Contract(sql.FieldNotNull(FieldLink))
}

// LinkEqualFold applies the EqualFold predicate on the "link" field.
func LinkEqualFold(v string) predicate.Contract {
	return predicate.Contract(sql.FieldEqualFold(FieldLink, v))
}

// LinkContainsFold applies the ContainsFold predicate on the "link" field.
func LinkContainsFold(v string) predicate.Contract {
	return predicate.Contract(sql.FieldContainsFold(FieldLink, v))
}

// ExpiresAtEQ applies the EQ predicate on the "expires_at" field.
func ExpiresAtEQ(v time.Time) predicate.Contract {
	return predicate.Contract(sql.FieldEQ(FieldExpiresAt, v))
}

// ExpiresAtNEQ applies the NEQ predicate on the "expires_at" field.
func ExpiresAtNEQ(v time.Time) predicate.Contract {
	return predicate.Contract(sql.FieldNEQ(FieldExpiresAt, v))
}

// ExpiresAtIn applies the In predicate on the "expires_at" field.
func ExpiresAtIn(vs ...time.Time) predicate.Contract {
	return predicate.Contract(sql.FieldIn(FieldExpiresAt, vs...))
}

// ExpiresAtNotIn applies the NotIn predicate on the "expires_at" field.
func ExpiresAtNotIn(vs ...time.Time) predicate.Contract {
	return predicate.Contract(sql.FieldNotIn(FieldExpiresAt, vs...))
}

// ExpiresAtGT applies the GT predicate on the "expires_at" field.
func ExpiresAtGT(v time.Time) predicate.Contract {
	return predicate.Contract(sql.FieldGT(FieldExpiresAt, v))
}

// ExpiresAtGTE applies the GTE predicate on the "expires_at" field.
func ExpiresAtGTE(v time.Time) predicate.Contract {
	return predicate.Contract(sql.FieldGTE(FieldExpiresAt, v))
}

// ExpiresAtLT applies the LT predicate on the "expires_at" field.
func ExpiresAtLT(v time.Time) predicate.Contract {
	return predicate.Contract(sql.FieldLT(FieldExpiresAt, v))
}

// ExpiresAtLTE applies the LTE predicate on the "expires_at" field.
func ExpiresAtLTE(v time.Time) predicate.Contract {
	return predicate.Contract(sql.FieldLTE(FieldExpiresAt, v))
}

// ExpiresAtIsNil applies the IsNil predicate on the "expires_at" field.
func ExpiresAtIsNil() predicate.Contract {
	return predicate.Contract(sql.FieldIsNull(FieldExpiresAt))
}

// ExpiresAtNotNil applies the NotNil predicate on the "expires_at" field.
func ExpiresAtNotNil() predicate.Contract {
	return predicate.Contract(sql.FieldNotNull(FieldExpiresAt))
}

// SignedAtEQ applies the EQ predicate on the "signed_at" field.
func SignedAtEQ(v time.Time) predicate.Contract {
	return predicate.Contract(sql.FieldEQ(FieldSignedAt, v))
}

// SignedAtNEQ applies the NEQ predicate on the "signed_at" field.
func SignedAtNEQ(v time.Time) predicate.Contract {
	return predicate.Contract(sql.FieldNEQ(FieldSignedAt, v))
}

// SignedAtIn applies the In predicate on the "signed_at" field.
func SignedAtIn(vs ...time.Time) predicate.Contract {
	return predicate.Contract(sql.FieldIn(FieldSignedAt, vs...))
}

// SignedAtNotIn applies the NotIn predicate on the "signed_at" field.
func SignedAtNotIn(vs ...time.Time) predicate.Contract {
	return predicate.Contract(sql.FieldNotIn(FieldSignedAt, vs...))
}

// SignedAtGT applies the GT predicate on the "signed_at" field.
func SignedAtGT(v time.Time) predicate.Contract {
	return predicate.Contract(sql.FieldGT(FieldSignedAt, v))
}

// SignedAtGTE applies the GTE predicate on the "signed_at" field.
func SignedAtGTE(v time.Time) predicate.Contract {
	return predicate.Contract(sql.FieldGTE(FieldSignedAt, v))
}

// SignedAtLT applies the LT predicate on the "signed_at" field.
func SignedAtLT(v time.Time) predicate.Contract {
	return predicate.Contract(sql.FieldLT(FieldSignedAt, v))
}

// SignedAtLTE applies the LTE predicate on the "signed_at" field.
func SignedAtLTE(v time.Time) predicate.Contract {
	return predicate.Contract(sql.FieldLTE(FieldSignedAt, v))
}

// SignedAtIsNil applies the IsNil predicate on the "signed_at" field.
func SignedAtIsNil() predicate.Contract {
	return predicate.Contract(sql.FieldIsNull(FieldSignedAt))
}

// SignedAtNotNil applies the NotNil predicate on the "signed_at" field.
func SignedAtNotNil() predicate.Contract {
	return predicate.Contract(sql.FieldNotNull(FieldSignedAt))
}

// DocIDEQ applies the EQ predicate on the "doc_id" field.
func DocIDEQ(v string) predicate.Contract {
	return predicate.Contract(sql.FieldEQ(FieldDocID, v))
}

// DocIDNEQ applies the NEQ predicate on the "doc_id" field.
func DocIDNEQ(v string) predicate.Contract {
	return predicate.Contract(sql.FieldNEQ(FieldDocID, v))
}

// DocIDIn applies the In predicate on the "doc_id" field.
func DocIDIn(vs ...string) predicate.Contract {
	return predicate.Contract(sql.FieldIn(FieldDocID, vs...))
}

// DocIDNotIn applies the NotIn predicate on the "doc_id" field.
func DocIDNotIn(vs ...string) predicate.Contract {
	return predicate.Contract(sql.FieldNotIn(FieldDocID, vs...))
}

// DocIDGT applies the GT predicate on the "doc_id" field.
func DocIDGT(v string) predicate.Contract {
	return predicate.Contract(sql.FieldGT(FieldDocID, v))
}

// DocIDGTE applies the GTE predicate on the "doc_id" field.
func DocIDGTE(v string) predicate.Contract {
	return predicate.Contract(sql.FieldGTE(FieldDocID, v))
}

// DocIDLT applies the LT predicate on the "doc_id" field.
func DocIDLT(v string) predicate.Contract {
	return predicate.Contract(sql.FieldLT(FieldDocID, v))
}

// DocIDLTE applies the LTE predicate on the "doc_id" field.
func DocIDLTE(v string) predicate.Contract {
	return predicate.Contract(sql.FieldLTE(FieldDocID, v))
}

// DocIDContains applies the Contains predicate on the "doc_id" field.
func DocIDContains(v string) predicate.Contract {
	return predicate.Contract(sql.FieldContains(FieldDocID, v))
}

// DocIDHasPrefix applies the HasPrefix predicate on the "doc_id" field.
func DocIDHasPrefix(v string) predicate.Contract {
	return predicate.Contract(sql.FieldHasPrefix(FieldDocID, v))
}

// DocIDHasSuffix applies the HasSuffix predicate on the "doc_id" field.
func DocIDHasSuffix(v string) predicate.Contract {
	return predicate.Contract(sql.FieldHasSuffix(FieldDocID, v))
}

// DocIDIsNil applies the IsNil predicate on the "doc_id" field.
func DocIDIsNil() predicate.Contract {
	return predicate.Contract(sql.FieldIsNull(FieldDocID))
}

// DocIDNotNil applies the NotNil predicate on the "doc_id" field.
func DocIDNotNil() predicate.Contract {
	return predicate.Contract(sql.FieldNotNull(FieldDocID))
}

// DocIDEqualFold applies the EqualFold predicate on the "doc_id" field.
func DocIDEqualFold(v string) predicate.Contract {
	return predicate.Contract(sql.FieldEqualFold(FieldDocID, v))
}

// DocIDContainsFold applies the ContainsFold predicate on the "doc_id" field.
func DocIDContainsFold(v string) predicate.Contract {
	return predicate.Contract(sql.FieldContainsFold(FieldDocID, v))
}

// HasSubscribe applies the HasEdge predicate on the "subscribe" edge.
func HasSubscribe() predicate.Contract {
	return predicate.Contract(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, SubscribeTable, SubscribeColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasSubscribeWith applies the HasEdge predicate on the "subscribe" edge with a given conditions (other predicates).
func HasSubscribeWith(preds ...predicate.Subscribe) predicate.Contract {
	return predicate.Contract(func(s *sql.Selector) {
		step := newSubscribeStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasEmployee applies the HasEdge predicate on the "employee" edge.
func HasEmployee() predicate.Contract {
	return predicate.Contract(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, EmployeeTable, EmployeeColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasEmployeeWith applies the HasEdge predicate on the "employee" edge with a given conditions (other predicates).
func HasEmployeeWith(preds ...predicate.Employee) predicate.Contract {
	return predicate.Contract(func(s *sql.Selector) {
		step := newEmployeeStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasRider applies the HasEdge predicate on the "rider" edge.
func HasRider() predicate.Contract {
	return predicate.Contract(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, RiderTable, RiderColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasRiderWith applies the HasEdge predicate on the "rider" edge with a given conditions (other predicates).
func HasRiderWith(preds ...predicate.Rider) predicate.Contract {
	return predicate.Contract(func(s *sql.Selector) {
		step := newRiderStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasAllocate applies the HasEdge predicate on the "allocate" edge.
func HasAllocate() predicate.Contract {
	return predicate.Contract(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.O2O, true, AllocateTable, AllocateColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasAllocateWith applies the HasEdge predicate on the "allocate" edge with a given conditions (other predicates).
func HasAllocateWith(preds ...predicate.Allocate) predicate.Contract {
	return predicate.Contract(func(s *sql.Selector) {
		step := newAllocateStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.Contract) predicate.Contract {
	return predicate.Contract(sql.AndPredicates(predicates...))
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.Contract) predicate.Contract {
	return predicate.Contract(sql.OrPredicates(predicates...))
}

// Not applies the not operator on the given predicate.
func Not(p predicate.Contract) predicate.Contract {
	return predicate.Contract(sql.NotPredicates(p))
}
