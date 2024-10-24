// Code generated by ent, DO NOT EDIT.

package promotionsetting

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/auroraride/aurservd/internal/ent/predicate"
)

// ID filters vertices based on their ID field.
func ID(id uint64) predicate.PromotionSetting {
	return predicate.PromotionSetting(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id uint64) predicate.PromotionSetting {
	return predicate.PromotionSetting(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id uint64) predicate.PromotionSetting {
	return predicate.PromotionSetting(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...uint64) predicate.PromotionSetting {
	return predicate.PromotionSetting(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...uint64) predicate.PromotionSetting {
	return predicate.PromotionSetting(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id uint64) predicate.PromotionSetting {
	return predicate.PromotionSetting(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id uint64) predicate.PromotionSetting {
	return predicate.PromotionSetting(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id uint64) predicate.PromotionSetting {
	return predicate.PromotionSetting(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id uint64) predicate.PromotionSetting {
	return predicate.PromotionSetting(sql.FieldLTE(FieldID, id))
}

// CreatedAt applies equality check predicate on the "created_at" field. It's identical to CreatedAtEQ.
func CreatedAt(v time.Time) predicate.PromotionSetting {
	return predicate.PromotionSetting(sql.FieldEQ(FieldCreatedAt, v))
}

// UpdatedAt applies equality check predicate on the "updated_at" field. It's identical to UpdatedAtEQ.
func UpdatedAt(v time.Time) predicate.PromotionSetting {
	return predicate.PromotionSetting(sql.FieldEQ(FieldUpdatedAt, v))
}

// Remark applies equality check predicate on the "remark" field. It's identical to RemarkEQ.
func Remark(v string) predicate.PromotionSetting {
	return predicate.PromotionSetting(sql.FieldEQ(FieldRemark, v))
}

// Title applies equality check predicate on the "title" field. It's identical to TitleEQ.
func Title(v string) predicate.PromotionSetting {
	return predicate.PromotionSetting(sql.FieldEQ(FieldTitle, v))
}

// Content applies equality check predicate on the "content" field. It's identical to ContentEQ.
func Content(v string) predicate.PromotionSetting {
	return predicate.PromotionSetting(sql.FieldEQ(FieldContent, v))
}

// Key applies equality check predicate on the "key" field. It's identical to KeyEQ.
func Key(v string) predicate.PromotionSetting {
	return predicate.PromotionSetting(sql.FieldEQ(FieldKey, v))
}

// CreatedAtEQ applies the EQ predicate on the "created_at" field.
func CreatedAtEQ(v time.Time) predicate.PromotionSetting {
	return predicate.PromotionSetting(sql.FieldEQ(FieldCreatedAt, v))
}

// CreatedAtNEQ applies the NEQ predicate on the "created_at" field.
func CreatedAtNEQ(v time.Time) predicate.PromotionSetting {
	return predicate.PromotionSetting(sql.FieldNEQ(FieldCreatedAt, v))
}

// CreatedAtIn applies the In predicate on the "created_at" field.
func CreatedAtIn(vs ...time.Time) predicate.PromotionSetting {
	return predicate.PromotionSetting(sql.FieldIn(FieldCreatedAt, vs...))
}

// CreatedAtNotIn applies the NotIn predicate on the "created_at" field.
func CreatedAtNotIn(vs ...time.Time) predicate.PromotionSetting {
	return predicate.PromotionSetting(sql.FieldNotIn(FieldCreatedAt, vs...))
}

// CreatedAtGT applies the GT predicate on the "created_at" field.
func CreatedAtGT(v time.Time) predicate.PromotionSetting {
	return predicate.PromotionSetting(sql.FieldGT(FieldCreatedAt, v))
}

// CreatedAtGTE applies the GTE predicate on the "created_at" field.
func CreatedAtGTE(v time.Time) predicate.PromotionSetting {
	return predicate.PromotionSetting(sql.FieldGTE(FieldCreatedAt, v))
}

// CreatedAtLT applies the LT predicate on the "created_at" field.
func CreatedAtLT(v time.Time) predicate.PromotionSetting {
	return predicate.PromotionSetting(sql.FieldLT(FieldCreatedAt, v))
}

// CreatedAtLTE applies the LTE predicate on the "created_at" field.
func CreatedAtLTE(v time.Time) predicate.PromotionSetting {
	return predicate.PromotionSetting(sql.FieldLTE(FieldCreatedAt, v))
}

// UpdatedAtEQ applies the EQ predicate on the "updated_at" field.
func UpdatedAtEQ(v time.Time) predicate.PromotionSetting {
	return predicate.PromotionSetting(sql.FieldEQ(FieldUpdatedAt, v))
}

// UpdatedAtNEQ applies the NEQ predicate on the "updated_at" field.
func UpdatedAtNEQ(v time.Time) predicate.PromotionSetting {
	return predicate.PromotionSetting(sql.FieldNEQ(FieldUpdatedAt, v))
}

// UpdatedAtIn applies the In predicate on the "updated_at" field.
func UpdatedAtIn(vs ...time.Time) predicate.PromotionSetting {
	return predicate.PromotionSetting(sql.FieldIn(FieldUpdatedAt, vs...))
}

// UpdatedAtNotIn applies the NotIn predicate on the "updated_at" field.
func UpdatedAtNotIn(vs ...time.Time) predicate.PromotionSetting {
	return predicate.PromotionSetting(sql.FieldNotIn(FieldUpdatedAt, vs...))
}

// UpdatedAtGT applies the GT predicate on the "updated_at" field.
func UpdatedAtGT(v time.Time) predicate.PromotionSetting {
	return predicate.PromotionSetting(sql.FieldGT(FieldUpdatedAt, v))
}

// UpdatedAtGTE applies the GTE predicate on the "updated_at" field.
func UpdatedAtGTE(v time.Time) predicate.PromotionSetting {
	return predicate.PromotionSetting(sql.FieldGTE(FieldUpdatedAt, v))
}

// UpdatedAtLT applies the LT predicate on the "updated_at" field.
func UpdatedAtLT(v time.Time) predicate.PromotionSetting {
	return predicate.PromotionSetting(sql.FieldLT(FieldUpdatedAt, v))
}

// UpdatedAtLTE applies the LTE predicate on the "updated_at" field.
func UpdatedAtLTE(v time.Time) predicate.PromotionSetting {
	return predicate.PromotionSetting(sql.FieldLTE(FieldUpdatedAt, v))
}

// CreatorIsNil applies the IsNil predicate on the "creator" field.
func CreatorIsNil() predicate.PromotionSetting {
	return predicate.PromotionSetting(sql.FieldIsNull(FieldCreator))
}

// CreatorNotNil applies the NotNil predicate on the "creator" field.
func CreatorNotNil() predicate.PromotionSetting {
	return predicate.PromotionSetting(sql.FieldNotNull(FieldCreator))
}

// LastModifierIsNil applies the IsNil predicate on the "last_modifier" field.
func LastModifierIsNil() predicate.PromotionSetting {
	return predicate.PromotionSetting(sql.FieldIsNull(FieldLastModifier))
}

// LastModifierNotNil applies the NotNil predicate on the "last_modifier" field.
func LastModifierNotNil() predicate.PromotionSetting {
	return predicate.PromotionSetting(sql.FieldNotNull(FieldLastModifier))
}

// RemarkEQ applies the EQ predicate on the "remark" field.
func RemarkEQ(v string) predicate.PromotionSetting {
	return predicate.PromotionSetting(sql.FieldEQ(FieldRemark, v))
}

// RemarkNEQ applies the NEQ predicate on the "remark" field.
func RemarkNEQ(v string) predicate.PromotionSetting {
	return predicate.PromotionSetting(sql.FieldNEQ(FieldRemark, v))
}

// RemarkIn applies the In predicate on the "remark" field.
func RemarkIn(vs ...string) predicate.PromotionSetting {
	return predicate.PromotionSetting(sql.FieldIn(FieldRemark, vs...))
}

// RemarkNotIn applies the NotIn predicate on the "remark" field.
func RemarkNotIn(vs ...string) predicate.PromotionSetting {
	return predicate.PromotionSetting(sql.FieldNotIn(FieldRemark, vs...))
}

// RemarkGT applies the GT predicate on the "remark" field.
func RemarkGT(v string) predicate.PromotionSetting {
	return predicate.PromotionSetting(sql.FieldGT(FieldRemark, v))
}

// RemarkGTE applies the GTE predicate on the "remark" field.
func RemarkGTE(v string) predicate.PromotionSetting {
	return predicate.PromotionSetting(sql.FieldGTE(FieldRemark, v))
}

// RemarkLT applies the LT predicate on the "remark" field.
func RemarkLT(v string) predicate.PromotionSetting {
	return predicate.PromotionSetting(sql.FieldLT(FieldRemark, v))
}

// RemarkLTE applies the LTE predicate on the "remark" field.
func RemarkLTE(v string) predicate.PromotionSetting {
	return predicate.PromotionSetting(sql.FieldLTE(FieldRemark, v))
}

// RemarkContains applies the Contains predicate on the "remark" field.
func RemarkContains(v string) predicate.PromotionSetting {
	return predicate.PromotionSetting(sql.FieldContains(FieldRemark, v))
}

// RemarkHasPrefix applies the HasPrefix predicate on the "remark" field.
func RemarkHasPrefix(v string) predicate.PromotionSetting {
	return predicate.PromotionSetting(sql.FieldHasPrefix(FieldRemark, v))
}

// RemarkHasSuffix applies the HasSuffix predicate on the "remark" field.
func RemarkHasSuffix(v string) predicate.PromotionSetting {
	return predicate.PromotionSetting(sql.FieldHasSuffix(FieldRemark, v))
}

// RemarkIsNil applies the IsNil predicate on the "remark" field.
func RemarkIsNil() predicate.PromotionSetting {
	return predicate.PromotionSetting(sql.FieldIsNull(FieldRemark))
}

// RemarkNotNil applies the NotNil predicate on the "remark" field.
func RemarkNotNil() predicate.PromotionSetting {
	return predicate.PromotionSetting(sql.FieldNotNull(FieldRemark))
}

// RemarkEqualFold applies the EqualFold predicate on the "remark" field.
func RemarkEqualFold(v string) predicate.PromotionSetting {
	return predicate.PromotionSetting(sql.FieldEqualFold(FieldRemark, v))
}

// RemarkContainsFold applies the ContainsFold predicate on the "remark" field.
func RemarkContainsFold(v string) predicate.PromotionSetting {
	return predicate.PromotionSetting(sql.FieldContainsFold(FieldRemark, v))
}

// TitleEQ applies the EQ predicate on the "title" field.
func TitleEQ(v string) predicate.PromotionSetting {
	return predicate.PromotionSetting(sql.FieldEQ(FieldTitle, v))
}

// TitleNEQ applies the NEQ predicate on the "title" field.
func TitleNEQ(v string) predicate.PromotionSetting {
	return predicate.PromotionSetting(sql.FieldNEQ(FieldTitle, v))
}

// TitleIn applies the In predicate on the "title" field.
func TitleIn(vs ...string) predicate.PromotionSetting {
	return predicate.PromotionSetting(sql.FieldIn(FieldTitle, vs...))
}

// TitleNotIn applies the NotIn predicate on the "title" field.
func TitleNotIn(vs ...string) predicate.PromotionSetting {
	return predicate.PromotionSetting(sql.FieldNotIn(FieldTitle, vs...))
}

// TitleGT applies the GT predicate on the "title" field.
func TitleGT(v string) predicate.PromotionSetting {
	return predicate.PromotionSetting(sql.FieldGT(FieldTitle, v))
}

// TitleGTE applies the GTE predicate on the "title" field.
func TitleGTE(v string) predicate.PromotionSetting {
	return predicate.PromotionSetting(sql.FieldGTE(FieldTitle, v))
}

// TitleLT applies the LT predicate on the "title" field.
func TitleLT(v string) predicate.PromotionSetting {
	return predicate.PromotionSetting(sql.FieldLT(FieldTitle, v))
}

// TitleLTE applies the LTE predicate on the "title" field.
func TitleLTE(v string) predicate.PromotionSetting {
	return predicate.PromotionSetting(sql.FieldLTE(FieldTitle, v))
}

// TitleContains applies the Contains predicate on the "title" field.
func TitleContains(v string) predicate.PromotionSetting {
	return predicate.PromotionSetting(sql.FieldContains(FieldTitle, v))
}

// TitleHasPrefix applies the HasPrefix predicate on the "title" field.
func TitleHasPrefix(v string) predicate.PromotionSetting {
	return predicate.PromotionSetting(sql.FieldHasPrefix(FieldTitle, v))
}

// TitleHasSuffix applies the HasSuffix predicate on the "title" field.
func TitleHasSuffix(v string) predicate.PromotionSetting {
	return predicate.PromotionSetting(sql.FieldHasSuffix(FieldTitle, v))
}

// TitleIsNil applies the IsNil predicate on the "title" field.
func TitleIsNil() predicate.PromotionSetting {
	return predicate.PromotionSetting(sql.FieldIsNull(FieldTitle))
}

// TitleNotNil applies the NotNil predicate on the "title" field.
func TitleNotNil() predicate.PromotionSetting {
	return predicate.PromotionSetting(sql.FieldNotNull(FieldTitle))
}

// TitleEqualFold applies the EqualFold predicate on the "title" field.
func TitleEqualFold(v string) predicate.PromotionSetting {
	return predicate.PromotionSetting(sql.FieldEqualFold(FieldTitle, v))
}

// TitleContainsFold applies the ContainsFold predicate on the "title" field.
func TitleContainsFold(v string) predicate.PromotionSetting {
	return predicate.PromotionSetting(sql.FieldContainsFold(FieldTitle, v))
}

// ContentEQ applies the EQ predicate on the "content" field.
func ContentEQ(v string) predicate.PromotionSetting {
	return predicate.PromotionSetting(sql.FieldEQ(FieldContent, v))
}

// ContentNEQ applies the NEQ predicate on the "content" field.
func ContentNEQ(v string) predicate.PromotionSetting {
	return predicate.PromotionSetting(sql.FieldNEQ(FieldContent, v))
}

// ContentIn applies the In predicate on the "content" field.
func ContentIn(vs ...string) predicate.PromotionSetting {
	return predicate.PromotionSetting(sql.FieldIn(FieldContent, vs...))
}

// ContentNotIn applies the NotIn predicate on the "content" field.
func ContentNotIn(vs ...string) predicate.PromotionSetting {
	return predicate.PromotionSetting(sql.FieldNotIn(FieldContent, vs...))
}

// ContentGT applies the GT predicate on the "content" field.
func ContentGT(v string) predicate.PromotionSetting {
	return predicate.PromotionSetting(sql.FieldGT(FieldContent, v))
}

// ContentGTE applies the GTE predicate on the "content" field.
func ContentGTE(v string) predicate.PromotionSetting {
	return predicate.PromotionSetting(sql.FieldGTE(FieldContent, v))
}

// ContentLT applies the LT predicate on the "content" field.
func ContentLT(v string) predicate.PromotionSetting {
	return predicate.PromotionSetting(sql.FieldLT(FieldContent, v))
}

// ContentLTE applies the LTE predicate on the "content" field.
func ContentLTE(v string) predicate.PromotionSetting {
	return predicate.PromotionSetting(sql.FieldLTE(FieldContent, v))
}

// ContentContains applies the Contains predicate on the "content" field.
func ContentContains(v string) predicate.PromotionSetting {
	return predicate.PromotionSetting(sql.FieldContains(FieldContent, v))
}

// ContentHasPrefix applies the HasPrefix predicate on the "content" field.
func ContentHasPrefix(v string) predicate.PromotionSetting {
	return predicate.PromotionSetting(sql.FieldHasPrefix(FieldContent, v))
}

// ContentHasSuffix applies the HasSuffix predicate on the "content" field.
func ContentHasSuffix(v string) predicate.PromotionSetting {
	return predicate.PromotionSetting(sql.FieldHasSuffix(FieldContent, v))
}

// ContentIsNil applies the IsNil predicate on the "content" field.
func ContentIsNil() predicate.PromotionSetting {
	return predicate.PromotionSetting(sql.FieldIsNull(FieldContent))
}

// ContentNotNil applies the NotNil predicate on the "content" field.
func ContentNotNil() predicate.PromotionSetting {
	return predicate.PromotionSetting(sql.FieldNotNull(FieldContent))
}

// ContentEqualFold applies the EqualFold predicate on the "content" field.
func ContentEqualFold(v string) predicate.PromotionSetting {
	return predicate.PromotionSetting(sql.FieldEqualFold(FieldContent, v))
}

// ContentContainsFold applies the ContainsFold predicate on the "content" field.
func ContentContainsFold(v string) predicate.PromotionSetting {
	return predicate.PromotionSetting(sql.FieldContainsFold(FieldContent, v))
}

// KeyEQ applies the EQ predicate on the "key" field.
func KeyEQ(v string) predicate.PromotionSetting {
	return predicate.PromotionSetting(sql.FieldEQ(FieldKey, v))
}

// KeyNEQ applies the NEQ predicate on the "key" field.
func KeyNEQ(v string) predicate.PromotionSetting {
	return predicate.PromotionSetting(sql.FieldNEQ(FieldKey, v))
}

// KeyIn applies the In predicate on the "key" field.
func KeyIn(vs ...string) predicate.PromotionSetting {
	return predicate.PromotionSetting(sql.FieldIn(FieldKey, vs...))
}

// KeyNotIn applies the NotIn predicate on the "key" field.
func KeyNotIn(vs ...string) predicate.PromotionSetting {
	return predicate.PromotionSetting(sql.FieldNotIn(FieldKey, vs...))
}

// KeyGT applies the GT predicate on the "key" field.
func KeyGT(v string) predicate.PromotionSetting {
	return predicate.PromotionSetting(sql.FieldGT(FieldKey, v))
}

// KeyGTE applies the GTE predicate on the "key" field.
func KeyGTE(v string) predicate.PromotionSetting {
	return predicate.PromotionSetting(sql.FieldGTE(FieldKey, v))
}

// KeyLT applies the LT predicate on the "key" field.
func KeyLT(v string) predicate.PromotionSetting {
	return predicate.PromotionSetting(sql.FieldLT(FieldKey, v))
}

// KeyLTE applies the LTE predicate on the "key" field.
func KeyLTE(v string) predicate.PromotionSetting {
	return predicate.PromotionSetting(sql.FieldLTE(FieldKey, v))
}

// KeyContains applies the Contains predicate on the "key" field.
func KeyContains(v string) predicate.PromotionSetting {
	return predicate.PromotionSetting(sql.FieldContains(FieldKey, v))
}

// KeyHasPrefix applies the HasPrefix predicate on the "key" field.
func KeyHasPrefix(v string) predicate.PromotionSetting {
	return predicate.PromotionSetting(sql.FieldHasPrefix(FieldKey, v))
}

// KeyHasSuffix applies the HasSuffix predicate on the "key" field.
func KeyHasSuffix(v string) predicate.PromotionSetting {
	return predicate.PromotionSetting(sql.FieldHasSuffix(FieldKey, v))
}

// KeyEqualFold applies the EqualFold predicate on the "key" field.
func KeyEqualFold(v string) predicate.PromotionSetting {
	return predicate.PromotionSetting(sql.FieldEqualFold(FieldKey, v))
}

// KeyContainsFold applies the ContainsFold predicate on the "key" field.
func KeyContainsFold(v string) predicate.PromotionSetting {
	return predicate.PromotionSetting(sql.FieldContainsFold(FieldKey, v))
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.PromotionSetting) predicate.PromotionSetting {
	return predicate.PromotionSetting(sql.AndPredicates(predicates...))
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.PromotionSetting) predicate.PromotionSetting {
	return predicate.PromotionSetting(sql.OrPredicates(predicates...))
}

// Not applies the not operator on the given predicate.
func Not(p predicate.PromotionSetting) predicate.PromotionSetting {
	return predicate.PromotionSetting(sql.NotPredicates(p))
}
