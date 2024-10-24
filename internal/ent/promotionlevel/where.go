// Code generated by ent, DO NOT EDIT.

package promotionlevel

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/auroraride/aurservd/internal/ent/predicate"
)

// ID filters vertices based on their ID field.
func ID(id uint64) predicate.PromotionLevel {
	return predicate.PromotionLevel(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id uint64) predicate.PromotionLevel {
	return predicate.PromotionLevel(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id uint64) predicate.PromotionLevel {
	return predicate.PromotionLevel(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...uint64) predicate.PromotionLevel {
	return predicate.PromotionLevel(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...uint64) predicate.PromotionLevel {
	return predicate.PromotionLevel(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id uint64) predicate.PromotionLevel {
	return predicate.PromotionLevel(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id uint64) predicate.PromotionLevel {
	return predicate.PromotionLevel(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id uint64) predicate.PromotionLevel {
	return predicate.PromotionLevel(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id uint64) predicate.PromotionLevel {
	return predicate.PromotionLevel(sql.FieldLTE(FieldID, id))
}

// CreatedAt applies equality check predicate on the "created_at" field. It's identical to CreatedAtEQ.
func CreatedAt(v time.Time) predicate.PromotionLevel {
	return predicate.PromotionLevel(sql.FieldEQ(FieldCreatedAt, v))
}

// UpdatedAt applies equality check predicate on the "updated_at" field. It's identical to UpdatedAtEQ.
func UpdatedAt(v time.Time) predicate.PromotionLevel {
	return predicate.PromotionLevel(sql.FieldEQ(FieldUpdatedAt, v))
}

// DeletedAt applies equality check predicate on the "deleted_at" field. It's identical to DeletedAtEQ.
func DeletedAt(v time.Time) predicate.PromotionLevel {
	return predicate.PromotionLevel(sql.FieldEQ(FieldDeletedAt, v))
}

// Remark applies equality check predicate on the "remark" field. It's identical to RemarkEQ.
func Remark(v string) predicate.PromotionLevel {
	return predicate.PromotionLevel(sql.FieldEQ(FieldRemark, v))
}

// Level applies equality check predicate on the "level" field. It's identical to LevelEQ.
func Level(v uint64) predicate.PromotionLevel {
	return predicate.PromotionLevel(sql.FieldEQ(FieldLevel, v))
}

// GrowthValue applies equality check predicate on the "growth_value" field. It's identical to GrowthValueEQ.
func GrowthValue(v uint64) predicate.PromotionLevel {
	return predicate.PromotionLevel(sql.FieldEQ(FieldGrowthValue, v))
}

// CommissionRatio applies equality check predicate on the "commission_ratio" field. It's identical to CommissionRatioEQ.
func CommissionRatio(v float64) predicate.PromotionLevel {
	return predicate.PromotionLevel(sql.FieldEQ(FieldCommissionRatio, v))
}

// CreatedAtEQ applies the EQ predicate on the "created_at" field.
func CreatedAtEQ(v time.Time) predicate.PromotionLevel {
	return predicate.PromotionLevel(sql.FieldEQ(FieldCreatedAt, v))
}

// CreatedAtNEQ applies the NEQ predicate on the "created_at" field.
func CreatedAtNEQ(v time.Time) predicate.PromotionLevel {
	return predicate.PromotionLevel(sql.FieldNEQ(FieldCreatedAt, v))
}

// CreatedAtIn applies the In predicate on the "created_at" field.
func CreatedAtIn(vs ...time.Time) predicate.PromotionLevel {
	return predicate.PromotionLevel(sql.FieldIn(FieldCreatedAt, vs...))
}

// CreatedAtNotIn applies the NotIn predicate on the "created_at" field.
func CreatedAtNotIn(vs ...time.Time) predicate.PromotionLevel {
	return predicate.PromotionLevel(sql.FieldNotIn(FieldCreatedAt, vs...))
}

// CreatedAtGT applies the GT predicate on the "created_at" field.
func CreatedAtGT(v time.Time) predicate.PromotionLevel {
	return predicate.PromotionLevel(sql.FieldGT(FieldCreatedAt, v))
}

// CreatedAtGTE applies the GTE predicate on the "created_at" field.
func CreatedAtGTE(v time.Time) predicate.PromotionLevel {
	return predicate.PromotionLevel(sql.FieldGTE(FieldCreatedAt, v))
}

// CreatedAtLT applies the LT predicate on the "created_at" field.
func CreatedAtLT(v time.Time) predicate.PromotionLevel {
	return predicate.PromotionLevel(sql.FieldLT(FieldCreatedAt, v))
}

// CreatedAtLTE applies the LTE predicate on the "created_at" field.
func CreatedAtLTE(v time.Time) predicate.PromotionLevel {
	return predicate.PromotionLevel(sql.FieldLTE(FieldCreatedAt, v))
}

// UpdatedAtEQ applies the EQ predicate on the "updated_at" field.
func UpdatedAtEQ(v time.Time) predicate.PromotionLevel {
	return predicate.PromotionLevel(sql.FieldEQ(FieldUpdatedAt, v))
}

// UpdatedAtNEQ applies the NEQ predicate on the "updated_at" field.
func UpdatedAtNEQ(v time.Time) predicate.PromotionLevel {
	return predicate.PromotionLevel(sql.FieldNEQ(FieldUpdatedAt, v))
}

// UpdatedAtIn applies the In predicate on the "updated_at" field.
func UpdatedAtIn(vs ...time.Time) predicate.PromotionLevel {
	return predicate.PromotionLevel(sql.FieldIn(FieldUpdatedAt, vs...))
}

// UpdatedAtNotIn applies the NotIn predicate on the "updated_at" field.
func UpdatedAtNotIn(vs ...time.Time) predicate.PromotionLevel {
	return predicate.PromotionLevel(sql.FieldNotIn(FieldUpdatedAt, vs...))
}

// UpdatedAtGT applies the GT predicate on the "updated_at" field.
func UpdatedAtGT(v time.Time) predicate.PromotionLevel {
	return predicate.PromotionLevel(sql.FieldGT(FieldUpdatedAt, v))
}

// UpdatedAtGTE applies the GTE predicate on the "updated_at" field.
func UpdatedAtGTE(v time.Time) predicate.PromotionLevel {
	return predicate.PromotionLevel(sql.FieldGTE(FieldUpdatedAt, v))
}

// UpdatedAtLT applies the LT predicate on the "updated_at" field.
func UpdatedAtLT(v time.Time) predicate.PromotionLevel {
	return predicate.PromotionLevel(sql.FieldLT(FieldUpdatedAt, v))
}

// UpdatedAtLTE applies the LTE predicate on the "updated_at" field.
func UpdatedAtLTE(v time.Time) predicate.PromotionLevel {
	return predicate.PromotionLevel(sql.FieldLTE(FieldUpdatedAt, v))
}

// DeletedAtEQ applies the EQ predicate on the "deleted_at" field.
func DeletedAtEQ(v time.Time) predicate.PromotionLevel {
	return predicate.PromotionLevel(sql.FieldEQ(FieldDeletedAt, v))
}

// DeletedAtNEQ applies the NEQ predicate on the "deleted_at" field.
func DeletedAtNEQ(v time.Time) predicate.PromotionLevel {
	return predicate.PromotionLevel(sql.FieldNEQ(FieldDeletedAt, v))
}

// DeletedAtIn applies the In predicate on the "deleted_at" field.
func DeletedAtIn(vs ...time.Time) predicate.PromotionLevel {
	return predicate.PromotionLevel(sql.FieldIn(FieldDeletedAt, vs...))
}

// DeletedAtNotIn applies the NotIn predicate on the "deleted_at" field.
func DeletedAtNotIn(vs ...time.Time) predicate.PromotionLevel {
	return predicate.PromotionLevel(sql.FieldNotIn(FieldDeletedAt, vs...))
}

// DeletedAtGT applies the GT predicate on the "deleted_at" field.
func DeletedAtGT(v time.Time) predicate.PromotionLevel {
	return predicate.PromotionLevel(sql.FieldGT(FieldDeletedAt, v))
}

// DeletedAtGTE applies the GTE predicate on the "deleted_at" field.
func DeletedAtGTE(v time.Time) predicate.PromotionLevel {
	return predicate.PromotionLevel(sql.FieldGTE(FieldDeletedAt, v))
}

// DeletedAtLT applies the LT predicate on the "deleted_at" field.
func DeletedAtLT(v time.Time) predicate.PromotionLevel {
	return predicate.PromotionLevel(sql.FieldLT(FieldDeletedAt, v))
}

// DeletedAtLTE applies the LTE predicate on the "deleted_at" field.
func DeletedAtLTE(v time.Time) predicate.PromotionLevel {
	return predicate.PromotionLevel(sql.FieldLTE(FieldDeletedAt, v))
}

// DeletedAtIsNil applies the IsNil predicate on the "deleted_at" field.
func DeletedAtIsNil() predicate.PromotionLevel {
	return predicate.PromotionLevel(sql.FieldIsNull(FieldDeletedAt))
}

// DeletedAtNotNil applies the NotNil predicate on the "deleted_at" field.
func DeletedAtNotNil() predicate.PromotionLevel {
	return predicate.PromotionLevel(sql.FieldNotNull(FieldDeletedAt))
}

// CreatorIsNil applies the IsNil predicate on the "creator" field.
func CreatorIsNil() predicate.PromotionLevel {
	return predicate.PromotionLevel(sql.FieldIsNull(FieldCreator))
}

// CreatorNotNil applies the NotNil predicate on the "creator" field.
func CreatorNotNil() predicate.PromotionLevel {
	return predicate.PromotionLevel(sql.FieldNotNull(FieldCreator))
}

// LastModifierIsNil applies the IsNil predicate on the "last_modifier" field.
func LastModifierIsNil() predicate.PromotionLevel {
	return predicate.PromotionLevel(sql.FieldIsNull(FieldLastModifier))
}

// LastModifierNotNil applies the NotNil predicate on the "last_modifier" field.
func LastModifierNotNil() predicate.PromotionLevel {
	return predicate.PromotionLevel(sql.FieldNotNull(FieldLastModifier))
}

// RemarkEQ applies the EQ predicate on the "remark" field.
func RemarkEQ(v string) predicate.PromotionLevel {
	return predicate.PromotionLevel(sql.FieldEQ(FieldRemark, v))
}

// RemarkNEQ applies the NEQ predicate on the "remark" field.
func RemarkNEQ(v string) predicate.PromotionLevel {
	return predicate.PromotionLevel(sql.FieldNEQ(FieldRemark, v))
}

// RemarkIn applies the In predicate on the "remark" field.
func RemarkIn(vs ...string) predicate.PromotionLevel {
	return predicate.PromotionLevel(sql.FieldIn(FieldRemark, vs...))
}

// RemarkNotIn applies the NotIn predicate on the "remark" field.
func RemarkNotIn(vs ...string) predicate.PromotionLevel {
	return predicate.PromotionLevel(sql.FieldNotIn(FieldRemark, vs...))
}

// RemarkGT applies the GT predicate on the "remark" field.
func RemarkGT(v string) predicate.PromotionLevel {
	return predicate.PromotionLevel(sql.FieldGT(FieldRemark, v))
}

// RemarkGTE applies the GTE predicate on the "remark" field.
func RemarkGTE(v string) predicate.PromotionLevel {
	return predicate.PromotionLevel(sql.FieldGTE(FieldRemark, v))
}

// RemarkLT applies the LT predicate on the "remark" field.
func RemarkLT(v string) predicate.PromotionLevel {
	return predicate.PromotionLevel(sql.FieldLT(FieldRemark, v))
}

// RemarkLTE applies the LTE predicate on the "remark" field.
func RemarkLTE(v string) predicate.PromotionLevel {
	return predicate.PromotionLevel(sql.FieldLTE(FieldRemark, v))
}

// RemarkContains applies the Contains predicate on the "remark" field.
func RemarkContains(v string) predicate.PromotionLevel {
	return predicate.PromotionLevel(sql.FieldContains(FieldRemark, v))
}

// RemarkHasPrefix applies the HasPrefix predicate on the "remark" field.
func RemarkHasPrefix(v string) predicate.PromotionLevel {
	return predicate.PromotionLevel(sql.FieldHasPrefix(FieldRemark, v))
}

// RemarkHasSuffix applies the HasSuffix predicate on the "remark" field.
func RemarkHasSuffix(v string) predicate.PromotionLevel {
	return predicate.PromotionLevel(sql.FieldHasSuffix(FieldRemark, v))
}

// RemarkIsNil applies the IsNil predicate on the "remark" field.
func RemarkIsNil() predicate.PromotionLevel {
	return predicate.PromotionLevel(sql.FieldIsNull(FieldRemark))
}

// RemarkNotNil applies the NotNil predicate on the "remark" field.
func RemarkNotNil() predicate.PromotionLevel {
	return predicate.PromotionLevel(sql.FieldNotNull(FieldRemark))
}

// RemarkEqualFold applies the EqualFold predicate on the "remark" field.
func RemarkEqualFold(v string) predicate.PromotionLevel {
	return predicate.PromotionLevel(sql.FieldEqualFold(FieldRemark, v))
}

// RemarkContainsFold applies the ContainsFold predicate on the "remark" field.
func RemarkContainsFold(v string) predicate.PromotionLevel {
	return predicate.PromotionLevel(sql.FieldContainsFold(FieldRemark, v))
}

// LevelEQ applies the EQ predicate on the "level" field.
func LevelEQ(v uint64) predicate.PromotionLevel {
	return predicate.PromotionLevel(sql.FieldEQ(FieldLevel, v))
}

// LevelNEQ applies the NEQ predicate on the "level" field.
func LevelNEQ(v uint64) predicate.PromotionLevel {
	return predicate.PromotionLevel(sql.FieldNEQ(FieldLevel, v))
}

// LevelIn applies the In predicate on the "level" field.
func LevelIn(vs ...uint64) predicate.PromotionLevel {
	return predicate.PromotionLevel(sql.FieldIn(FieldLevel, vs...))
}

// LevelNotIn applies the NotIn predicate on the "level" field.
func LevelNotIn(vs ...uint64) predicate.PromotionLevel {
	return predicate.PromotionLevel(sql.FieldNotIn(FieldLevel, vs...))
}

// LevelGT applies the GT predicate on the "level" field.
func LevelGT(v uint64) predicate.PromotionLevel {
	return predicate.PromotionLevel(sql.FieldGT(FieldLevel, v))
}

// LevelGTE applies the GTE predicate on the "level" field.
func LevelGTE(v uint64) predicate.PromotionLevel {
	return predicate.PromotionLevel(sql.FieldGTE(FieldLevel, v))
}

// LevelLT applies the LT predicate on the "level" field.
func LevelLT(v uint64) predicate.PromotionLevel {
	return predicate.PromotionLevel(sql.FieldLT(FieldLevel, v))
}

// LevelLTE applies the LTE predicate on the "level" field.
func LevelLTE(v uint64) predicate.PromotionLevel {
	return predicate.PromotionLevel(sql.FieldLTE(FieldLevel, v))
}

// GrowthValueEQ applies the EQ predicate on the "growth_value" field.
func GrowthValueEQ(v uint64) predicate.PromotionLevel {
	return predicate.PromotionLevel(sql.FieldEQ(FieldGrowthValue, v))
}

// GrowthValueNEQ applies the NEQ predicate on the "growth_value" field.
func GrowthValueNEQ(v uint64) predicate.PromotionLevel {
	return predicate.PromotionLevel(sql.FieldNEQ(FieldGrowthValue, v))
}

// GrowthValueIn applies the In predicate on the "growth_value" field.
func GrowthValueIn(vs ...uint64) predicate.PromotionLevel {
	return predicate.PromotionLevel(sql.FieldIn(FieldGrowthValue, vs...))
}

// GrowthValueNotIn applies the NotIn predicate on the "growth_value" field.
func GrowthValueNotIn(vs ...uint64) predicate.PromotionLevel {
	return predicate.PromotionLevel(sql.FieldNotIn(FieldGrowthValue, vs...))
}

// GrowthValueGT applies the GT predicate on the "growth_value" field.
func GrowthValueGT(v uint64) predicate.PromotionLevel {
	return predicate.PromotionLevel(sql.FieldGT(FieldGrowthValue, v))
}

// GrowthValueGTE applies the GTE predicate on the "growth_value" field.
func GrowthValueGTE(v uint64) predicate.PromotionLevel {
	return predicate.PromotionLevel(sql.FieldGTE(FieldGrowthValue, v))
}

// GrowthValueLT applies the LT predicate on the "growth_value" field.
func GrowthValueLT(v uint64) predicate.PromotionLevel {
	return predicate.PromotionLevel(sql.FieldLT(FieldGrowthValue, v))
}

// GrowthValueLTE applies the LTE predicate on the "growth_value" field.
func GrowthValueLTE(v uint64) predicate.PromotionLevel {
	return predicate.PromotionLevel(sql.FieldLTE(FieldGrowthValue, v))
}

// CommissionRatioEQ applies the EQ predicate on the "commission_ratio" field.
func CommissionRatioEQ(v float64) predicate.PromotionLevel {
	return predicate.PromotionLevel(sql.FieldEQ(FieldCommissionRatio, v))
}

// CommissionRatioNEQ applies the NEQ predicate on the "commission_ratio" field.
func CommissionRatioNEQ(v float64) predicate.PromotionLevel {
	return predicate.PromotionLevel(sql.FieldNEQ(FieldCommissionRatio, v))
}

// CommissionRatioIn applies the In predicate on the "commission_ratio" field.
func CommissionRatioIn(vs ...float64) predicate.PromotionLevel {
	return predicate.PromotionLevel(sql.FieldIn(FieldCommissionRatio, vs...))
}

// CommissionRatioNotIn applies the NotIn predicate on the "commission_ratio" field.
func CommissionRatioNotIn(vs ...float64) predicate.PromotionLevel {
	return predicate.PromotionLevel(sql.FieldNotIn(FieldCommissionRatio, vs...))
}

// CommissionRatioGT applies the GT predicate on the "commission_ratio" field.
func CommissionRatioGT(v float64) predicate.PromotionLevel {
	return predicate.PromotionLevel(sql.FieldGT(FieldCommissionRatio, v))
}

// CommissionRatioGTE applies the GTE predicate on the "commission_ratio" field.
func CommissionRatioGTE(v float64) predicate.PromotionLevel {
	return predicate.PromotionLevel(sql.FieldGTE(FieldCommissionRatio, v))
}

// CommissionRatioLT applies the LT predicate on the "commission_ratio" field.
func CommissionRatioLT(v float64) predicate.PromotionLevel {
	return predicate.PromotionLevel(sql.FieldLT(FieldCommissionRatio, v))
}

// CommissionRatioLTE applies the LTE predicate on the "commission_ratio" field.
func CommissionRatioLTE(v float64) predicate.PromotionLevel {
	return predicate.PromotionLevel(sql.FieldLTE(FieldCommissionRatio, v))
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.PromotionLevel) predicate.PromotionLevel {
	return predicate.PromotionLevel(sql.AndPredicates(predicates...))
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.PromotionLevel) predicate.PromotionLevel {
	return predicate.PromotionLevel(sql.OrPredicates(predicates...))
}

// Not applies the not operator on the given predicate.
func Not(p predicate.PromotionLevel) predicate.PromotionLevel {
	return predicate.PromotionLevel(sql.NotPredicates(p))
}
