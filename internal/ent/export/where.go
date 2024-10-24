// Code generated by ent, DO NOT EDIT.

package export

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/auroraride/aurservd/internal/ent/predicate"
)

// ID filters vertices based on their ID field.
func ID(id uint64) predicate.Export {
	return predicate.Export(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id uint64) predicate.Export {
	return predicate.Export(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id uint64) predicate.Export {
	return predicate.Export(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...uint64) predicate.Export {
	return predicate.Export(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...uint64) predicate.Export {
	return predicate.Export(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id uint64) predicate.Export {
	return predicate.Export(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id uint64) predicate.Export {
	return predicate.Export(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id uint64) predicate.Export {
	return predicate.Export(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id uint64) predicate.Export {
	return predicate.Export(sql.FieldLTE(FieldID, id))
}

// CreatedAt applies equality check predicate on the "created_at" field. It's identical to CreatedAtEQ.
func CreatedAt(v time.Time) predicate.Export {
	return predicate.Export(sql.FieldEQ(FieldCreatedAt, v))
}

// UpdatedAt applies equality check predicate on the "updated_at" field. It's identical to UpdatedAtEQ.
func UpdatedAt(v time.Time) predicate.Export {
	return predicate.Export(sql.FieldEQ(FieldUpdatedAt, v))
}

// DeletedAt applies equality check predicate on the "deleted_at" field. It's identical to DeletedAtEQ.
func DeletedAt(v time.Time) predicate.Export {
	return predicate.Export(sql.FieldEQ(FieldDeletedAt, v))
}

// ManagerID applies equality check predicate on the "manager_id" field. It's identical to ManagerIDEQ.
func ManagerID(v uint64) predicate.Export {
	return predicate.Export(sql.FieldEQ(FieldManagerID, v))
}

// Taxonomy applies equality check predicate on the "taxonomy" field. It's identical to TaxonomyEQ.
func Taxonomy(v string) predicate.Export {
	return predicate.Export(sql.FieldEQ(FieldTaxonomy, v))
}

// Sn applies equality check predicate on the "sn" field. It's identical to SnEQ.
func Sn(v string) predicate.Export {
	return predicate.Export(sql.FieldEQ(FieldSn, v))
}

// Status applies equality check predicate on the "status" field. It's identical to StatusEQ.
func Status(v uint8) predicate.Export {
	return predicate.Export(sql.FieldEQ(FieldStatus, v))
}

// Path applies equality check predicate on the "path" field. It's identical to PathEQ.
func Path(v string) predicate.Export {
	return predicate.Export(sql.FieldEQ(FieldPath, v))
}

// Message applies equality check predicate on the "message" field. It's identical to MessageEQ.
func Message(v string) predicate.Export {
	return predicate.Export(sql.FieldEQ(FieldMessage, v))
}

// FinishAt applies equality check predicate on the "finish_at" field. It's identical to FinishAtEQ.
func FinishAt(v time.Time) predicate.Export {
	return predicate.Export(sql.FieldEQ(FieldFinishAt, v))
}

// Duration applies equality check predicate on the "duration" field. It's identical to DurationEQ.
func Duration(v int64) predicate.Export {
	return predicate.Export(sql.FieldEQ(FieldDuration, v))
}

// Condition applies equality check predicate on the "condition" field. It's identical to ConditionEQ.
func Condition(v string) predicate.Export {
	return predicate.Export(sql.FieldEQ(FieldCondition, v))
}

// Remark applies equality check predicate on the "remark" field. It's identical to RemarkEQ.
func Remark(v string) predicate.Export {
	return predicate.Export(sql.FieldEQ(FieldRemark, v))
}

// CreatedAtEQ applies the EQ predicate on the "created_at" field.
func CreatedAtEQ(v time.Time) predicate.Export {
	return predicate.Export(sql.FieldEQ(FieldCreatedAt, v))
}

// CreatedAtNEQ applies the NEQ predicate on the "created_at" field.
func CreatedAtNEQ(v time.Time) predicate.Export {
	return predicate.Export(sql.FieldNEQ(FieldCreatedAt, v))
}

// CreatedAtIn applies the In predicate on the "created_at" field.
func CreatedAtIn(vs ...time.Time) predicate.Export {
	return predicate.Export(sql.FieldIn(FieldCreatedAt, vs...))
}

// CreatedAtNotIn applies the NotIn predicate on the "created_at" field.
func CreatedAtNotIn(vs ...time.Time) predicate.Export {
	return predicate.Export(sql.FieldNotIn(FieldCreatedAt, vs...))
}

// CreatedAtGT applies the GT predicate on the "created_at" field.
func CreatedAtGT(v time.Time) predicate.Export {
	return predicate.Export(sql.FieldGT(FieldCreatedAt, v))
}

// CreatedAtGTE applies the GTE predicate on the "created_at" field.
func CreatedAtGTE(v time.Time) predicate.Export {
	return predicate.Export(sql.FieldGTE(FieldCreatedAt, v))
}

// CreatedAtLT applies the LT predicate on the "created_at" field.
func CreatedAtLT(v time.Time) predicate.Export {
	return predicate.Export(sql.FieldLT(FieldCreatedAt, v))
}

// CreatedAtLTE applies the LTE predicate on the "created_at" field.
func CreatedAtLTE(v time.Time) predicate.Export {
	return predicate.Export(sql.FieldLTE(FieldCreatedAt, v))
}

// UpdatedAtEQ applies the EQ predicate on the "updated_at" field.
func UpdatedAtEQ(v time.Time) predicate.Export {
	return predicate.Export(sql.FieldEQ(FieldUpdatedAt, v))
}

// UpdatedAtNEQ applies the NEQ predicate on the "updated_at" field.
func UpdatedAtNEQ(v time.Time) predicate.Export {
	return predicate.Export(sql.FieldNEQ(FieldUpdatedAt, v))
}

// UpdatedAtIn applies the In predicate on the "updated_at" field.
func UpdatedAtIn(vs ...time.Time) predicate.Export {
	return predicate.Export(sql.FieldIn(FieldUpdatedAt, vs...))
}

// UpdatedAtNotIn applies the NotIn predicate on the "updated_at" field.
func UpdatedAtNotIn(vs ...time.Time) predicate.Export {
	return predicate.Export(sql.FieldNotIn(FieldUpdatedAt, vs...))
}

// UpdatedAtGT applies the GT predicate on the "updated_at" field.
func UpdatedAtGT(v time.Time) predicate.Export {
	return predicate.Export(sql.FieldGT(FieldUpdatedAt, v))
}

// UpdatedAtGTE applies the GTE predicate on the "updated_at" field.
func UpdatedAtGTE(v time.Time) predicate.Export {
	return predicate.Export(sql.FieldGTE(FieldUpdatedAt, v))
}

// UpdatedAtLT applies the LT predicate on the "updated_at" field.
func UpdatedAtLT(v time.Time) predicate.Export {
	return predicate.Export(sql.FieldLT(FieldUpdatedAt, v))
}

// UpdatedAtLTE applies the LTE predicate on the "updated_at" field.
func UpdatedAtLTE(v time.Time) predicate.Export {
	return predicate.Export(sql.FieldLTE(FieldUpdatedAt, v))
}

// DeletedAtEQ applies the EQ predicate on the "deleted_at" field.
func DeletedAtEQ(v time.Time) predicate.Export {
	return predicate.Export(sql.FieldEQ(FieldDeletedAt, v))
}

// DeletedAtNEQ applies the NEQ predicate on the "deleted_at" field.
func DeletedAtNEQ(v time.Time) predicate.Export {
	return predicate.Export(sql.FieldNEQ(FieldDeletedAt, v))
}

// DeletedAtIn applies the In predicate on the "deleted_at" field.
func DeletedAtIn(vs ...time.Time) predicate.Export {
	return predicate.Export(sql.FieldIn(FieldDeletedAt, vs...))
}

// DeletedAtNotIn applies the NotIn predicate on the "deleted_at" field.
func DeletedAtNotIn(vs ...time.Time) predicate.Export {
	return predicate.Export(sql.FieldNotIn(FieldDeletedAt, vs...))
}

// DeletedAtGT applies the GT predicate on the "deleted_at" field.
func DeletedAtGT(v time.Time) predicate.Export {
	return predicate.Export(sql.FieldGT(FieldDeletedAt, v))
}

// DeletedAtGTE applies the GTE predicate on the "deleted_at" field.
func DeletedAtGTE(v time.Time) predicate.Export {
	return predicate.Export(sql.FieldGTE(FieldDeletedAt, v))
}

// DeletedAtLT applies the LT predicate on the "deleted_at" field.
func DeletedAtLT(v time.Time) predicate.Export {
	return predicate.Export(sql.FieldLT(FieldDeletedAt, v))
}

// DeletedAtLTE applies the LTE predicate on the "deleted_at" field.
func DeletedAtLTE(v time.Time) predicate.Export {
	return predicate.Export(sql.FieldLTE(FieldDeletedAt, v))
}

// DeletedAtIsNil applies the IsNil predicate on the "deleted_at" field.
func DeletedAtIsNil() predicate.Export {
	return predicate.Export(sql.FieldIsNull(FieldDeletedAt))
}

// DeletedAtNotNil applies the NotNil predicate on the "deleted_at" field.
func DeletedAtNotNil() predicate.Export {
	return predicate.Export(sql.FieldNotNull(FieldDeletedAt))
}

// ManagerIDEQ applies the EQ predicate on the "manager_id" field.
func ManagerIDEQ(v uint64) predicate.Export {
	return predicate.Export(sql.FieldEQ(FieldManagerID, v))
}

// ManagerIDNEQ applies the NEQ predicate on the "manager_id" field.
func ManagerIDNEQ(v uint64) predicate.Export {
	return predicate.Export(sql.FieldNEQ(FieldManagerID, v))
}

// ManagerIDIn applies the In predicate on the "manager_id" field.
func ManagerIDIn(vs ...uint64) predicate.Export {
	return predicate.Export(sql.FieldIn(FieldManagerID, vs...))
}

// ManagerIDNotIn applies the NotIn predicate on the "manager_id" field.
func ManagerIDNotIn(vs ...uint64) predicate.Export {
	return predicate.Export(sql.FieldNotIn(FieldManagerID, vs...))
}

// TaxonomyEQ applies the EQ predicate on the "taxonomy" field.
func TaxonomyEQ(v string) predicate.Export {
	return predicate.Export(sql.FieldEQ(FieldTaxonomy, v))
}

// TaxonomyNEQ applies the NEQ predicate on the "taxonomy" field.
func TaxonomyNEQ(v string) predicate.Export {
	return predicate.Export(sql.FieldNEQ(FieldTaxonomy, v))
}

// TaxonomyIn applies the In predicate on the "taxonomy" field.
func TaxonomyIn(vs ...string) predicate.Export {
	return predicate.Export(sql.FieldIn(FieldTaxonomy, vs...))
}

// TaxonomyNotIn applies the NotIn predicate on the "taxonomy" field.
func TaxonomyNotIn(vs ...string) predicate.Export {
	return predicate.Export(sql.FieldNotIn(FieldTaxonomy, vs...))
}

// TaxonomyGT applies the GT predicate on the "taxonomy" field.
func TaxonomyGT(v string) predicate.Export {
	return predicate.Export(sql.FieldGT(FieldTaxonomy, v))
}

// TaxonomyGTE applies the GTE predicate on the "taxonomy" field.
func TaxonomyGTE(v string) predicate.Export {
	return predicate.Export(sql.FieldGTE(FieldTaxonomy, v))
}

// TaxonomyLT applies the LT predicate on the "taxonomy" field.
func TaxonomyLT(v string) predicate.Export {
	return predicate.Export(sql.FieldLT(FieldTaxonomy, v))
}

// TaxonomyLTE applies the LTE predicate on the "taxonomy" field.
func TaxonomyLTE(v string) predicate.Export {
	return predicate.Export(sql.FieldLTE(FieldTaxonomy, v))
}

// TaxonomyContains applies the Contains predicate on the "taxonomy" field.
func TaxonomyContains(v string) predicate.Export {
	return predicate.Export(sql.FieldContains(FieldTaxonomy, v))
}

// TaxonomyHasPrefix applies the HasPrefix predicate on the "taxonomy" field.
func TaxonomyHasPrefix(v string) predicate.Export {
	return predicate.Export(sql.FieldHasPrefix(FieldTaxonomy, v))
}

// TaxonomyHasSuffix applies the HasSuffix predicate on the "taxonomy" field.
func TaxonomyHasSuffix(v string) predicate.Export {
	return predicate.Export(sql.FieldHasSuffix(FieldTaxonomy, v))
}

// TaxonomyEqualFold applies the EqualFold predicate on the "taxonomy" field.
func TaxonomyEqualFold(v string) predicate.Export {
	return predicate.Export(sql.FieldEqualFold(FieldTaxonomy, v))
}

// TaxonomyContainsFold applies the ContainsFold predicate on the "taxonomy" field.
func TaxonomyContainsFold(v string) predicate.Export {
	return predicate.Export(sql.FieldContainsFold(FieldTaxonomy, v))
}

// SnEQ applies the EQ predicate on the "sn" field.
func SnEQ(v string) predicate.Export {
	return predicate.Export(sql.FieldEQ(FieldSn, v))
}

// SnNEQ applies the NEQ predicate on the "sn" field.
func SnNEQ(v string) predicate.Export {
	return predicate.Export(sql.FieldNEQ(FieldSn, v))
}

// SnIn applies the In predicate on the "sn" field.
func SnIn(vs ...string) predicate.Export {
	return predicate.Export(sql.FieldIn(FieldSn, vs...))
}

// SnNotIn applies the NotIn predicate on the "sn" field.
func SnNotIn(vs ...string) predicate.Export {
	return predicate.Export(sql.FieldNotIn(FieldSn, vs...))
}

// SnGT applies the GT predicate on the "sn" field.
func SnGT(v string) predicate.Export {
	return predicate.Export(sql.FieldGT(FieldSn, v))
}

// SnGTE applies the GTE predicate on the "sn" field.
func SnGTE(v string) predicate.Export {
	return predicate.Export(sql.FieldGTE(FieldSn, v))
}

// SnLT applies the LT predicate on the "sn" field.
func SnLT(v string) predicate.Export {
	return predicate.Export(sql.FieldLT(FieldSn, v))
}

// SnLTE applies the LTE predicate on the "sn" field.
func SnLTE(v string) predicate.Export {
	return predicate.Export(sql.FieldLTE(FieldSn, v))
}

// SnContains applies the Contains predicate on the "sn" field.
func SnContains(v string) predicate.Export {
	return predicate.Export(sql.FieldContains(FieldSn, v))
}

// SnHasPrefix applies the HasPrefix predicate on the "sn" field.
func SnHasPrefix(v string) predicate.Export {
	return predicate.Export(sql.FieldHasPrefix(FieldSn, v))
}

// SnHasSuffix applies the HasSuffix predicate on the "sn" field.
func SnHasSuffix(v string) predicate.Export {
	return predicate.Export(sql.FieldHasSuffix(FieldSn, v))
}

// SnEqualFold applies the EqualFold predicate on the "sn" field.
func SnEqualFold(v string) predicate.Export {
	return predicate.Export(sql.FieldEqualFold(FieldSn, v))
}

// SnContainsFold applies the ContainsFold predicate on the "sn" field.
func SnContainsFold(v string) predicate.Export {
	return predicate.Export(sql.FieldContainsFold(FieldSn, v))
}

// StatusEQ applies the EQ predicate on the "status" field.
func StatusEQ(v uint8) predicate.Export {
	return predicate.Export(sql.FieldEQ(FieldStatus, v))
}

// StatusNEQ applies the NEQ predicate on the "status" field.
func StatusNEQ(v uint8) predicate.Export {
	return predicate.Export(sql.FieldNEQ(FieldStatus, v))
}

// StatusIn applies the In predicate on the "status" field.
func StatusIn(vs ...uint8) predicate.Export {
	return predicate.Export(sql.FieldIn(FieldStatus, vs...))
}

// StatusNotIn applies the NotIn predicate on the "status" field.
func StatusNotIn(vs ...uint8) predicate.Export {
	return predicate.Export(sql.FieldNotIn(FieldStatus, vs...))
}

// StatusGT applies the GT predicate on the "status" field.
func StatusGT(v uint8) predicate.Export {
	return predicate.Export(sql.FieldGT(FieldStatus, v))
}

// StatusGTE applies the GTE predicate on the "status" field.
func StatusGTE(v uint8) predicate.Export {
	return predicate.Export(sql.FieldGTE(FieldStatus, v))
}

// StatusLT applies the LT predicate on the "status" field.
func StatusLT(v uint8) predicate.Export {
	return predicate.Export(sql.FieldLT(FieldStatus, v))
}

// StatusLTE applies the LTE predicate on the "status" field.
func StatusLTE(v uint8) predicate.Export {
	return predicate.Export(sql.FieldLTE(FieldStatus, v))
}

// PathEQ applies the EQ predicate on the "path" field.
func PathEQ(v string) predicate.Export {
	return predicate.Export(sql.FieldEQ(FieldPath, v))
}

// PathNEQ applies the NEQ predicate on the "path" field.
func PathNEQ(v string) predicate.Export {
	return predicate.Export(sql.FieldNEQ(FieldPath, v))
}

// PathIn applies the In predicate on the "path" field.
func PathIn(vs ...string) predicate.Export {
	return predicate.Export(sql.FieldIn(FieldPath, vs...))
}

// PathNotIn applies the NotIn predicate on the "path" field.
func PathNotIn(vs ...string) predicate.Export {
	return predicate.Export(sql.FieldNotIn(FieldPath, vs...))
}

// PathGT applies the GT predicate on the "path" field.
func PathGT(v string) predicate.Export {
	return predicate.Export(sql.FieldGT(FieldPath, v))
}

// PathGTE applies the GTE predicate on the "path" field.
func PathGTE(v string) predicate.Export {
	return predicate.Export(sql.FieldGTE(FieldPath, v))
}

// PathLT applies the LT predicate on the "path" field.
func PathLT(v string) predicate.Export {
	return predicate.Export(sql.FieldLT(FieldPath, v))
}

// PathLTE applies the LTE predicate on the "path" field.
func PathLTE(v string) predicate.Export {
	return predicate.Export(sql.FieldLTE(FieldPath, v))
}

// PathContains applies the Contains predicate on the "path" field.
func PathContains(v string) predicate.Export {
	return predicate.Export(sql.FieldContains(FieldPath, v))
}

// PathHasPrefix applies the HasPrefix predicate on the "path" field.
func PathHasPrefix(v string) predicate.Export {
	return predicate.Export(sql.FieldHasPrefix(FieldPath, v))
}

// PathHasSuffix applies the HasSuffix predicate on the "path" field.
func PathHasSuffix(v string) predicate.Export {
	return predicate.Export(sql.FieldHasSuffix(FieldPath, v))
}

// PathIsNil applies the IsNil predicate on the "path" field.
func PathIsNil() predicate.Export {
	return predicate.Export(sql.FieldIsNull(FieldPath))
}

// PathNotNil applies the NotNil predicate on the "path" field.
func PathNotNil() predicate.Export {
	return predicate.Export(sql.FieldNotNull(FieldPath))
}

// PathEqualFold applies the EqualFold predicate on the "path" field.
func PathEqualFold(v string) predicate.Export {
	return predicate.Export(sql.FieldEqualFold(FieldPath, v))
}

// PathContainsFold applies the ContainsFold predicate on the "path" field.
func PathContainsFold(v string) predicate.Export {
	return predicate.Export(sql.FieldContainsFold(FieldPath, v))
}

// MessageEQ applies the EQ predicate on the "message" field.
func MessageEQ(v string) predicate.Export {
	return predicate.Export(sql.FieldEQ(FieldMessage, v))
}

// MessageNEQ applies the NEQ predicate on the "message" field.
func MessageNEQ(v string) predicate.Export {
	return predicate.Export(sql.FieldNEQ(FieldMessage, v))
}

// MessageIn applies the In predicate on the "message" field.
func MessageIn(vs ...string) predicate.Export {
	return predicate.Export(sql.FieldIn(FieldMessage, vs...))
}

// MessageNotIn applies the NotIn predicate on the "message" field.
func MessageNotIn(vs ...string) predicate.Export {
	return predicate.Export(sql.FieldNotIn(FieldMessage, vs...))
}

// MessageGT applies the GT predicate on the "message" field.
func MessageGT(v string) predicate.Export {
	return predicate.Export(sql.FieldGT(FieldMessage, v))
}

// MessageGTE applies the GTE predicate on the "message" field.
func MessageGTE(v string) predicate.Export {
	return predicate.Export(sql.FieldGTE(FieldMessage, v))
}

// MessageLT applies the LT predicate on the "message" field.
func MessageLT(v string) predicate.Export {
	return predicate.Export(sql.FieldLT(FieldMessage, v))
}

// MessageLTE applies the LTE predicate on the "message" field.
func MessageLTE(v string) predicate.Export {
	return predicate.Export(sql.FieldLTE(FieldMessage, v))
}

// MessageContains applies the Contains predicate on the "message" field.
func MessageContains(v string) predicate.Export {
	return predicate.Export(sql.FieldContains(FieldMessage, v))
}

// MessageHasPrefix applies the HasPrefix predicate on the "message" field.
func MessageHasPrefix(v string) predicate.Export {
	return predicate.Export(sql.FieldHasPrefix(FieldMessage, v))
}

// MessageHasSuffix applies the HasSuffix predicate on the "message" field.
func MessageHasSuffix(v string) predicate.Export {
	return predicate.Export(sql.FieldHasSuffix(FieldMessage, v))
}

// MessageIsNil applies the IsNil predicate on the "message" field.
func MessageIsNil() predicate.Export {
	return predicate.Export(sql.FieldIsNull(FieldMessage))
}

// MessageNotNil applies the NotNil predicate on the "message" field.
func MessageNotNil() predicate.Export {
	return predicate.Export(sql.FieldNotNull(FieldMessage))
}

// MessageEqualFold applies the EqualFold predicate on the "message" field.
func MessageEqualFold(v string) predicate.Export {
	return predicate.Export(sql.FieldEqualFold(FieldMessage, v))
}

// MessageContainsFold applies the ContainsFold predicate on the "message" field.
func MessageContainsFold(v string) predicate.Export {
	return predicate.Export(sql.FieldContainsFold(FieldMessage, v))
}

// FinishAtEQ applies the EQ predicate on the "finish_at" field.
func FinishAtEQ(v time.Time) predicate.Export {
	return predicate.Export(sql.FieldEQ(FieldFinishAt, v))
}

// FinishAtNEQ applies the NEQ predicate on the "finish_at" field.
func FinishAtNEQ(v time.Time) predicate.Export {
	return predicate.Export(sql.FieldNEQ(FieldFinishAt, v))
}

// FinishAtIn applies the In predicate on the "finish_at" field.
func FinishAtIn(vs ...time.Time) predicate.Export {
	return predicate.Export(sql.FieldIn(FieldFinishAt, vs...))
}

// FinishAtNotIn applies the NotIn predicate on the "finish_at" field.
func FinishAtNotIn(vs ...time.Time) predicate.Export {
	return predicate.Export(sql.FieldNotIn(FieldFinishAt, vs...))
}

// FinishAtGT applies the GT predicate on the "finish_at" field.
func FinishAtGT(v time.Time) predicate.Export {
	return predicate.Export(sql.FieldGT(FieldFinishAt, v))
}

// FinishAtGTE applies the GTE predicate on the "finish_at" field.
func FinishAtGTE(v time.Time) predicate.Export {
	return predicate.Export(sql.FieldGTE(FieldFinishAt, v))
}

// FinishAtLT applies the LT predicate on the "finish_at" field.
func FinishAtLT(v time.Time) predicate.Export {
	return predicate.Export(sql.FieldLT(FieldFinishAt, v))
}

// FinishAtLTE applies the LTE predicate on the "finish_at" field.
func FinishAtLTE(v time.Time) predicate.Export {
	return predicate.Export(sql.FieldLTE(FieldFinishAt, v))
}

// FinishAtIsNil applies the IsNil predicate on the "finish_at" field.
func FinishAtIsNil() predicate.Export {
	return predicate.Export(sql.FieldIsNull(FieldFinishAt))
}

// FinishAtNotNil applies the NotNil predicate on the "finish_at" field.
func FinishAtNotNil() predicate.Export {
	return predicate.Export(sql.FieldNotNull(FieldFinishAt))
}

// DurationEQ applies the EQ predicate on the "duration" field.
func DurationEQ(v int64) predicate.Export {
	return predicate.Export(sql.FieldEQ(FieldDuration, v))
}

// DurationNEQ applies the NEQ predicate on the "duration" field.
func DurationNEQ(v int64) predicate.Export {
	return predicate.Export(sql.FieldNEQ(FieldDuration, v))
}

// DurationIn applies the In predicate on the "duration" field.
func DurationIn(vs ...int64) predicate.Export {
	return predicate.Export(sql.FieldIn(FieldDuration, vs...))
}

// DurationNotIn applies the NotIn predicate on the "duration" field.
func DurationNotIn(vs ...int64) predicate.Export {
	return predicate.Export(sql.FieldNotIn(FieldDuration, vs...))
}

// DurationGT applies the GT predicate on the "duration" field.
func DurationGT(v int64) predicate.Export {
	return predicate.Export(sql.FieldGT(FieldDuration, v))
}

// DurationGTE applies the GTE predicate on the "duration" field.
func DurationGTE(v int64) predicate.Export {
	return predicate.Export(sql.FieldGTE(FieldDuration, v))
}

// DurationLT applies the LT predicate on the "duration" field.
func DurationLT(v int64) predicate.Export {
	return predicate.Export(sql.FieldLT(FieldDuration, v))
}

// DurationLTE applies the LTE predicate on the "duration" field.
func DurationLTE(v int64) predicate.Export {
	return predicate.Export(sql.FieldLTE(FieldDuration, v))
}

// DurationIsNil applies the IsNil predicate on the "duration" field.
func DurationIsNil() predicate.Export {
	return predicate.Export(sql.FieldIsNull(FieldDuration))
}

// DurationNotNil applies the NotNil predicate on the "duration" field.
func DurationNotNil() predicate.Export {
	return predicate.Export(sql.FieldNotNull(FieldDuration))
}

// ConditionEQ applies the EQ predicate on the "condition" field.
func ConditionEQ(v string) predicate.Export {
	return predicate.Export(sql.FieldEQ(FieldCondition, v))
}

// ConditionNEQ applies the NEQ predicate on the "condition" field.
func ConditionNEQ(v string) predicate.Export {
	return predicate.Export(sql.FieldNEQ(FieldCondition, v))
}

// ConditionIn applies the In predicate on the "condition" field.
func ConditionIn(vs ...string) predicate.Export {
	return predicate.Export(sql.FieldIn(FieldCondition, vs...))
}

// ConditionNotIn applies the NotIn predicate on the "condition" field.
func ConditionNotIn(vs ...string) predicate.Export {
	return predicate.Export(sql.FieldNotIn(FieldCondition, vs...))
}

// ConditionGT applies the GT predicate on the "condition" field.
func ConditionGT(v string) predicate.Export {
	return predicate.Export(sql.FieldGT(FieldCondition, v))
}

// ConditionGTE applies the GTE predicate on the "condition" field.
func ConditionGTE(v string) predicate.Export {
	return predicate.Export(sql.FieldGTE(FieldCondition, v))
}

// ConditionLT applies the LT predicate on the "condition" field.
func ConditionLT(v string) predicate.Export {
	return predicate.Export(sql.FieldLT(FieldCondition, v))
}

// ConditionLTE applies the LTE predicate on the "condition" field.
func ConditionLTE(v string) predicate.Export {
	return predicate.Export(sql.FieldLTE(FieldCondition, v))
}

// ConditionContains applies the Contains predicate on the "condition" field.
func ConditionContains(v string) predicate.Export {
	return predicate.Export(sql.FieldContains(FieldCondition, v))
}

// ConditionHasPrefix applies the HasPrefix predicate on the "condition" field.
func ConditionHasPrefix(v string) predicate.Export {
	return predicate.Export(sql.FieldHasPrefix(FieldCondition, v))
}

// ConditionHasSuffix applies the HasSuffix predicate on the "condition" field.
func ConditionHasSuffix(v string) predicate.Export {
	return predicate.Export(sql.FieldHasSuffix(FieldCondition, v))
}

// ConditionEqualFold applies the EqualFold predicate on the "condition" field.
func ConditionEqualFold(v string) predicate.Export {
	return predicate.Export(sql.FieldEqualFold(FieldCondition, v))
}

// ConditionContainsFold applies the ContainsFold predicate on the "condition" field.
func ConditionContainsFold(v string) predicate.Export {
	return predicate.Export(sql.FieldContainsFold(FieldCondition, v))
}

// InfoIsNil applies the IsNil predicate on the "info" field.
func InfoIsNil() predicate.Export {
	return predicate.Export(sql.FieldIsNull(FieldInfo))
}

// InfoNotNil applies the NotNil predicate on the "info" field.
func InfoNotNil() predicate.Export {
	return predicate.Export(sql.FieldNotNull(FieldInfo))
}

// RemarkEQ applies the EQ predicate on the "remark" field.
func RemarkEQ(v string) predicate.Export {
	return predicate.Export(sql.FieldEQ(FieldRemark, v))
}

// RemarkNEQ applies the NEQ predicate on the "remark" field.
func RemarkNEQ(v string) predicate.Export {
	return predicate.Export(sql.FieldNEQ(FieldRemark, v))
}

// RemarkIn applies the In predicate on the "remark" field.
func RemarkIn(vs ...string) predicate.Export {
	return predicate.Export(sql.FieldIn(FieldRemark, vs...))
}

// RemarkNotIn applies the NotIn predicate on the "remark" field.
func RemarkNotIn(vs ...string) predicate.Export {
	return predicate.Export(sql.FieldNotIn(FieldRemark, vs...))
}

// RemarkGT applies the GT predicate on the "remark" field.
func RemarkGT(v string) predicate.Export {
	return predicate.Export(sql.FieldGT(FieldRemark, v))
}

// RemarkGTE applies the GTE predicate on the "remark" field.
func RemarkGTE(v string) predicate.Export {
	return predicate.Export(sql.FieldGTE(FieldRemark, v))
}

// RemarkLT applies the LT predicate on the "remark" field.
func RemarkLT(v string) predicate.Export {
	return predicate.Export(sql.FieldLT(FieldRemark, v))
}

// RemarkLTE applies the LTE predicate on the "remark" field.
func RemarkLTE(v string) predicate.Export {
	return predicate.Export(sql.FieldLTE(FieldRemark, v))
}

// RemarkContains applies the Contains predicate on the "remark" field.
func RemarkContains(v string) predicate.Export {
	return predicate.Export(sql.FieldContains(FieldRemark, v))
}

// RemarkHasPrefix applies the HasPrefix predicate on the "remark" field.
func RemarkHasPrefix(v string) predicate.Export {
	return predicate.Export(sql.FieldHasPrefix(FieldRemark, v))
}

// RemarkHasSuffix applies the HasSuffix predicate on the "remark" field.
func RemarkHasSuffix(v string) predicate.Export {
	return predicate.Export(sql.FieldHasSuffix(FieldRemark, v))
}

// RemarkEqualFold applies the EqualFold predicate on the "remark" field.
func RemarkEqualFold(v string) predicate.Export {
	return predicate.Export(sql.FieldEqualFold(FieldRemark, v))
}

// RemarkContainsFold applies the ContainsFold predicate on the "remark" field.
func RemarkContainsFold(v string) predicate.Export {
	return predicate.Export(sql.FieldContainsFold(FieldRemark, v))
}

// HasManager applies the HasEdge predicate on the "manager" edge.
func HasManager() predicate.Export {
	return predicate.Export(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, ManagerTable, ManagerColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasManagerWith applies the HasEdge predicate on the "manager" edge with a given conditions (other predicates).
func HasManagerWith(preds ...predicate.Manager) predicate.Export {
	return predicate.Export(func(s *sql.Selector) {
		step := newManagerStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.Export) predicate.Export {
	return predicate.Export(sql.AndPredicates(predicates...))
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.Export) predicate.Export {
	return predicate.Export(sql.OrPredicates(predicates...))
}

// Not applies the not operator on the given predicate.
func Not(p predicate.Export) predicate.Export {
	return predicate.Export(sql.NotPredicates(p))
}
