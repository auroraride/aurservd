// Code generated by ent, DO NOT EDIT.

package promotionwithdrawal

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/auroraride/aurservd/internal/ent/predicate"
)

// ID filters vertices based on their ID field.
func ID(id uint64) predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id uint64) predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id uint64) predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...uint64) predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...uint64) predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id uint64) predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id uint64) predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id uint64) predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id uint64) predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldLTE(FieldID, id))
}

// CreatedAt applies equality check predicate on the "created_at" field. It's identical to CreatedAtEQ.
func CreatedAt(v time.Time) predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldEQ(FieldCreatedAt, v))
}

// UpdatedAt applies equality check predicate on the "updated_at" field. It's identical to UpdatedAtEQ.
func UpdatedAt(v time.Time) predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldEQ(FieldUpdatedAt, v))
}

// DeletedAt applies equality check predicate on the "deleted_at" field. It's identical to DeletedAtEQ.
func DeletedAt(v time.Time) predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldEQ(FieldDeletedAt, v))
}

// Remark applies equality check predicate on the "remark" field. It's identical to RemarkEQ.
func Remark(v string) predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldEQ(FieldRemark, v))
}

// MemberID applies equality check predicate on the "member_id" field. It's identical to MemberIDEQ.
func MemberID(v uint64) predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldEQ(FieldMemberID, v))
}

// Status applies equality check predicate on the "status" field. It's identical to StatusEQ.
func Status(v uint8) predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldEQ(FieldStatus, v))
}

// ApplyAmount applies equality check predicate on the "apply_amount" field. It's identical to ApplyAmountEQ.
func ApplyAmount(v float64) predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldEQ(FieldApplyAmount, v))
}

// Amount applies equality check predicate on the "amount" field. It's identical to AmountEQ.
func Amount(v float64) predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldEQ(FieldAmount, v))
}

// Fee applies equality check predicate on the "fee" field. It's identical to FeeEQ.
func Fee(v float64) predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldEQ(FieldFee, v))
}

// Method applies equality check predicate on the "method" field. It's identical to MethodEQ.
func Method(v uint8) predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldEQ(FieldMethod, v))
}

// AccountID applies equality check predicate on the "account_id" field. It's identical to AccountIDEQ.
func AccountID(v uint64) predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldEQ(FieldAccountID, v))
}

// ApplyTime applies equality check predicate on the "apply_time" field. It's identical to ApplyTimeEQ.
func ApplyTime(v time.Time) predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldEQ(FieldApplyTime, v))
}

// ReviewTime applies equality check predicate on the "review_time" field. It's identical to ReviewTimeEQ.
func ReviewTime(v time.Time) predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldEQ(FieldReviewTime, v))
}

// CreatedAtEQ applies the EQ predicate on the "created_at" field.
func CreatedAtEQ(v time.Time) predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldEQ(FieldCreatedAt, v))
}

// CreatedAtNEQ applies the NEQ predicate on the "created_at" field.
func CreatedAtNEQ(v time.Time) predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldNEQ(FieldCreatedAt, v))
}

// CreatedAtIn applies the In predicate on the "created_at" field.
func CreatedAtIn(vs ...time.Time) predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldIn(FieldCreatedAt, vs...))
}

// CreatedAtNotIn applies the NotIn predicate on the "created_at" field.
func CreatedAtNotIn(vs ...time.Time) predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldNotIn(FieldCreatedAt, vs...))
}

// CreatedAtGT applies the GT predicate on the "created_at" field.
func CreatedAtGT(v time.Time) predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldGT(FieldCreatedAt, v))
}

// CreatedAtGTE applies the GTE predicate on the "created_at" field.
func CreatedAtGTE(v time.Time) predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldGTE(FieldCreatedAt, v))
}

// CreatedAtLT applies the LT predicate on the "created_at" field.
func CreatedAtLT(v time.Time) predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldLT(FieldCreatedAt, v))
}

// CreatedAtLTE applies the LTE predicate on the "created_at" field.
func CreatedAtLTE(v time.Time) predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldLTE(FieldCreatedAt, v))
}

// UpdatedAtEQ applies the EQ predicate on the "updated_at" field.
func UpdatedAtEQ(v time.Time) predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldEQ(FieldUpdatedAt, v))
}

// UpdatedAtNEQ applies the NEQ predicate on the "updated_at" field.
func UpdatedAtNEQ(v time.Time) predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldNEQ(FieldUpdatedAt, v))
}

// UpdatedAtIn applies the In predicate on the "updated_at" field.
func UpdatedAtIn(vs ...time.Time) predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldIn(FieldUpdatedAt, vs...))
}

// UpdatedAtNotIn applies the NotIn predicate on the "updated_at" field.
func UpdatedAtNotIn(vs ...time.Time) predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldNotIn(FieldUpdatedAt, vs...))
}

// UpdatedAtGT applies the GT predicate on the "updated_at" field.
func UpdatedAtGT(v time.Time) predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldGT(FieldUpdatedAt, v))
}

// UpdatedAtGTE applies the GTE predicate on the "updated_at" field.
func UpdatedAtGTE(v time.Time) predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldGTE(FieldUpdatedAt, v))
}

// UpdatedAtLT applies the LT predicate on the "updated_at" field.
func UpdatedAtLT(v time.Time) predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldLT(FieldUpdatedAt, v))
}

// UpdatedAtLTE applies the LTE predicate on the "updated_at" field.
func UpdatedAtLTE(v time.Time) predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldLTE(FieldUpdatedAt, v))
}

// DeletedAtEQ applies the EQ predicate on the "deleted_at" field.
func DeletedAtEQ(v time.Time) predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldEQ(FieldDeletedAt, v))
}

// DeletedAtNEQ applies the NEQ predicate on the "deleted_at" field.
func DeletedAtNEQ(v time.Time) predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldNEQ(FieldDeletedAt, v))
}

// DeletedAtIn applies the In predicate on the "deleted_at" field.
func DeletedAtIn(vs ...time.Time) predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldIn(FieldDeletedAt, vs...))
}

// DeletedAtNotIn applies the NotIn predicate on the "deleted_at" field.
func DeletedAtNotIn(vs ...time.Time) predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldNotIn(FieldDeletedAt, vs...))
}

// DeletedAtGT applies the GT predicate on the "deleted_at" field.
func DeletedAtGT(v time.Time) predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldGT(FieldDeletedAt, v))
}

// DeletedAtGTE applies the GTE predicate on the "deleted_at" field.
func DeletedAtGTE(v time.Time) predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldGTE(FieldDeletedAt, v))
}

// DeletedAtLT applies the LT predicate on the "deleted_at" field.
func DeletedAtLT(v time.Time) predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldLT(FieldDeletedAt, v))
}

// DeletedAtLTE applies the LTE predicate on the "deleted_at" field.
func DeletedAtLTE(v time.Time) predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldLTE(FieldDeletedAt, v))
}

// DeletedAtIsNil applies the IsNil predicate on the "deleted_at" field.
func DeletedAtIsNil() predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldIsNull(FieldDeletedAt))
}

// DeletedAtNotNil applies the NotNil predicate on the "deleted_at" field.
func DeletedAtNotNil() predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldNotNull(FieldDeletedAt))
}

// CreatorIsNil applies the IsNil predicate on the "creator" field.
func CreatorIsNil() predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldIsNull(FieldCreator))
}

// CreatorNotNil applies the NotNil predicate on the "creator" field.
func CreatorNotNil() predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldNotNull(FieldCreator))
}

// LastModifierIsNil applies the IsNil predicate on the "last_modifier" field.
func LastModifierIsNil() predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldIsNull(FieldLastModifier))
}

// LastModifierNotNil applies the NotNil predicate on the "last_modifier" field.
func LastModifierNotNil() predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldNotNull(FieldLastModifier))
}

// RemarkEQ applies the EQ predicate on the "remark" field.
func RemarkEQ(v string) predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldEQ(FieldRemark, v))
}

// RemarkNEQ applies the NEQ predicate on the "remark" field.
func RemarkNEQ(v string) predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldNEQ(FieldRemark, v))
}

// RemarkIn applies the In predicate on the "remark" field.
func RemarkIn(vs ...string) predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldIn(FieldRemark, vs...))
}

// RemarkNotIn applies the NotIn predicate on the "remark" field.
func RemarkNotIn(vs ...string) predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldNotIn(FieldRemark, vs...))
}

// RemarkGT applies the GT predicate on the "remark" field.
func RemarkGT(v string) predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldGT(FieldRemark, v))
}

// RemarkGTE applies the GTE predicate on the "remark" field.
func RemarkGTE(v string) predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldGTE(FieldRemark, v))
}

// RemarkLT applies the LT predicate on the "remark" field.
func RemarkLT(v string) predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldLT(FieldRemark, v))
}

// RemarkLTE applies the LTE predicate on the "remark" field.
func RemarkLTE(v string) predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldLTE(FieldRemark, v))
}

// RemarkContains applies the Contains predicate on the "remark" field.
func RemarkContains(v string) predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldContains(FieldRemark, v))
}

// RemarkHasPrefix applies the HasPrefix predicate on the "remark" field.
func RemarkHasPrefix(v string) predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldHasPrefix(FieldRemark, v))
}

// RemarkHasSuffix applies the HasSuffix predicate on the "remark" field.
func RemarkHasSuffix(v string) predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldHasSuffix(FieldRemark, v))
}

// RemarkIsNil applies the IsNil predicate on the "remark" field.
func RemarkIsNil() predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldIsNull(FieldRemark))
}

// RemarkNotNil applies the NotNil predicate on the "remark" field.
func RemarkNotNil() predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldNotNull(FieldRemark))
}

// RemarkEqualFold applies the EqualFold predicate on the "remark" field.
func RemarkEqualFold(v string) predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldEqualFold(FieldRemark, v))
}

// RemarkContainsFold applies the ContainsFold predicate on the "remark" field.
func RemarkContainsFold(v string) predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldContainsFold(FieldRemark, v))
}

// MemberIDEQ applies the EQ predicate on the "member_id" field.
func MemberIDEQ(v uint64) predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldEQ(FieldMemberID, v))
}

// MemberIDNEQ applies the NEQ predicate on the "member_id" field.
func MemberIDNEQ(v uint64) predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldNEQ(FieldMemberID, v))
}

// MemberIDIn applies the In predicate on the "member_id" field.
func MemberIDIn(vs ...uint64) predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldIn(FieldMemberID, vs...))
}

// MemberIDNotIn applies the NotIn predicate on the "member_id" field.
func MemberIDNotIn(vs ...uint64) predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldNotIn(FieldMemberID, vs...))
}

// StatusEQ applies the EQ predicate on the "status" field.
func StatusEQ(v uint8) predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldEQ(FieldStatus, v))
}

// StatusNEQ applies the NEQ predicate on the "status" field.
func StatusNEQ(v uint8) predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldNEQ(FieldStatus, v))
}

// StatusIn applies the In predicate on the "status" field.
func StatusIn(vs ...uint8) predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldIn(FieldStatus, vs...))
}

// StatusNotIn applies the NotIn predicate on the "status" field.
func StatusNotIn(vs ...uint8) predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldNotIn(FieldStatus, vs...))
}

// StatusGT applies the GT predicate on the "status" field.
func StatusGT(v uint8) predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldGT(FieldStatus, v))
}

// StatusGTE applies the GTE predicate on the "status" field.
func StatusGTE(v uint8) predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldGTE(FieldStatus, v))
}

// StatusLT applies the LT predicate on the "status" field.
func StatusLT(v uint8) predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldLT(FieldStatus, v))
}

// StatusLTE applies the LTE predicate on the "status" field.
func StatusLTE(v uint8) predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldLTE(FieldStatus, v))
}

// ApplyAmountEQ applies the EQ predicate on the "apply_amount" field.
func ApplyAmountEQ(v float64) predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldEQ(FieldApplyAmount, v))
}

// ApplyAmountNEQ applies the NEQ predicate on the "apply_amount" field.
func ApplyAmountNEQ(v float64) predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldNEQ(FieldApplyAmount, v))
}

// ApplyAmountIn applies the In predicate on the "apply_amount" field.
func ApplyAmountIn(vs ...float64) predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldIn(FieldApplyAmount, vs...))
}

// ApplyAmountNotIn applies the NotIn predicate on the "apply_amount" field.
func ApplyAmountNotIn(vs ...float64) predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldNotIn(FieldApplyAmount, vs...))
}

// ApplyAmountGT applies the GT predicate on the "apply_amount" field.
func ApplyAmountGT(v float64) predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldGT(FieldApplyAmount, v))
}

// ApplyAmountGTE applies the GTE predicate on the "apply_amount" field.
func ApplyAmountGTE(v float64) predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldGTE(FieldApplyAmount, v))
}

// ApplyAmountLT applies the LT predicate on the "apply_amount" field.
func ApplyAmountLT(v float64) predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldLT(FieldApplyAmount, v))
}

// ApplyAmountLTE applies the LTE predicate on the "apply_amount" field.
func ApplyAmountLTE(v float64) predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldLTE(FieldApplyAmount, v))
}

// AmountEQ applies the EQ predicate on the "amount" field.
func AmountEQ(v float64) predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldEQ(FieldAmount, v))
}

// AmountNEQ applies the NEQ predicate on the "amount" field.
func AmountNEQ(v float64) predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldNEQ(FieldAmount, v))
}

// AmountIn applies the In predicate on the "amount" field.
func AmountIn(vs ...float64) predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldIn(FieldAmount, vs...))
}

// AmountNotIn applies the NotIn predicate on the "amount" field.
func AmountNotIn(vs ...float64) predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldNotIn(FieldAmount, vs...))
}

// AmountGT applies the GT predicate on the "amount" field.
func AmountGT(v float64) predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldGT(FieldAmount, v))
}

// AmountGTE applies the GTE predicate on the "amount" field.
func AmountGTE(v float64) predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldGTE(FieldAmount, v))
}

// AmountLT applies the LT predicate on the "amount" field.
func AmountLT(v float64) predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldLT(FieldAmount, v))
}

// AmountLTE applies the LTE predicate on the "amount" field.
func AmountLTE(v float64) predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldLTE(FieldAmount, v))
}

// FeeEQ applies the EQ predicate on the "fee" field.
func FeeEQ(v float64) predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldEQ(FieldFee, v))
}

// FeeNEQ applies the NEQ predicate on the "fee" field.
func FeeNEQ(v float64) predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldNEQ(FieldFee, v))
}

// FeeIn applies the In predicate on the "fee" field.
func FeeIn(vs ...float64) predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldIn(FieldFee, vs...))
}

// FeeNotIn applies the NotIn predicate on the "fee" field.
func FeeNotIn(vs ...float64) predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldNotIn(FieldFee, vs...))
}

// FeeGT applies the GT predicate on the "fee" field.
func FeeGT(v float64) predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldGT(FieldFee, v))
}

// FeeGTE applies the GTE predicate on the "fee" field.
func FeeGTE(v float64) predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldGTE(FieldFee, v))
}

// FeeLT applies the LT predicate on the "fee" field.
func FeeLT(v float64) predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldLT(FieldFee, v))
}

// FeeLTE applies the LTE predicate on the "fee" field.
func FeeLTE(v float64) predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldLTE(FieldFee, v))
}

// MethodEQ applies the EQ predicate on the "method" field.
func MethodEQ(v uint8) predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldEQ(FieldMethod, v))
}

// MethodNEQ applies the NEQ predicate on the "method" field.
func MethodNEQ(v uint8) predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldNEQ(FieldMethod, v))
}

// MethodIn applies the In predicate on the "method" field.
func MethodIn(vs ...uint8) predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldIn(FieldMethod, vs...))
}

// MethodNotIn applies the NotIn predicate on the "method" field.
func MethodNotIn(vs ...uint8) predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldNotIn(FieldMethod, vs...))
}

// MethodGT applies the GT predicate on the "method" field.
func MethodGT(v uint8) predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldGT(FieldMethod, v))
}

// MethodGTE applies the GTE predicate on the "method" field.
func MethodGTE(v uint8) predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldGTE(FieldMethod, v))
}

// MethodLT applies the LT predicate on the "method" field.
func MethodLT(v uint8) predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldLT(FieldMethod, v))
}

// MethodLTE applies the LTE predicate on the "method" field.
func MethodLTE(v uint8) predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldLTE(FieldMethod, v))
}

// AccountIDEQ applies the EQ predicate on the "account_id" field.
func AccountIDEQ(v uint64) predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldEQ(FieldAccountID, v))
}

// AccountIDNEQ applies the NEQ predicate on the "account_id" field.
func AccountIDNEQ(v uint64) predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldNEQ(FieldAccountID, v))
}

// AccountIDIn applies the In predicate on the "account_id" field.
func AccountIDIn(vs ...uint64) predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldIn(FieldAccountID, vs...))
}

// AccountIDNotIn applies the NotIn predicate on the "account_id" field.
func AccountIDNotIn(vs ...uint64) predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldNotIn(FieldAccountID, vs...))
}

// AccountIDIsNil applies the IsNil predicate on the "account_id" field.
func AccountIDIsNil() predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldIsNull(FieldAccountID))
}

// AccountIDNotNil applies the NotNil predicate on the "account_id" field.
func AccountIDNotNil() predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldNotNull(FieldAccountID))
}

// ApplyTimeEQ applies the EQ predicate on the "apply_time" field.
func ApplyTimeEQ(v time.Time) predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldEQ(FieldApplyTime, v))
}

// ApplyTimeNEQ applies the NEQ predicate on the "apply_time" field.
func ApplyTimeNEQ(v time.Time) predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldNEQ(FieldApplyTime, v))
}

// ApplyTimeIn applies the In predicate on the "apply_time" field.
func ApplyTimeIn(vs ...time.Time) predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldIn(FieldApplyTime, vs...))
}

// ApplyTimeNotIn applies the NotIn predicate on the "apply_time" field.
func ApplyTimeNotIn(vs ...time.Time) predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldNotIn(FieldApplyTime, vs...))
}

// ApplyTimeGT applies the GT predicate on the "apply_time" field.
func ApplyTimeGT(v time.Time) predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldGT(FieldApplyTime, v))
}

// ApplyTimeGTE applies the GTE predicate on the "apply_time" field.
func ApplyTimeGTE(v time.Time) predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldGTE(FieldApplyTime, v))
}

// ApplyTimeLT applies the LT predicate on the "apply_time" field.
func ApplyTimeLT(v time.Time) predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldLT(FieldApplyTime, v))
}

// ApplyTimeLTE applies the LTE predicate on the "apply_time" field.
func ApplyTimeLTE(v time.Time) predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldLTE(FieldApplyTime, v))
}

// ApplyTimeIsNil applies the IsNil predicate on the "apply_time" field.
func ApplyTimeIsNil() predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldIsNull(FieldApplyTime))
}

// ApplyTimeNotNil applies the NotNil predicate on the "apply_time" field.
func ApplyTimeNotNil() predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldNotNull(FieldApplyTime))
}

// ReviewTimeEQ applies the EQ predicate on the "review_time" field.
func ReviewTimeEQ(v time.Time) predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldEQ(FieldReviewTime, v))
}

// ReviewTimeNEQ applies the NEQ predicate on the "review_time" field.
func ReviewTimeNEQ(v time.Time) predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldNEQ(FieldReviewTime, v))
}

// ReviewTimeIn applies the In predicate on the "review_time" field.
func ReviewTimeIn(vs ...time.Time) predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldIn(FieldReviewTime, vs...))
}

// ReviewTimeNotIn applies the NotIn predicate on the "review_time" field.
func ReviewTimeNotIn(vs ...time.Time) predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldNotIn(FieldReviewTime, vs...))
}

// ReviewTimeGT applies the GT predicate on the "review_time" field.
func ReviewTimeGT(v time.Time) predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldGT(FieldReviewTime, v))
}

// ReviewTimeGTE applies the GTE predicate on the "review_time" field.
func ReviewTimeGTE(v time.Time) predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldGTE(FieldReviewTime, v))
}

// ReviewTimeLT applies the LT predicate on the "review_time" field.
func ReviewTimeLT(v time.Time) predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldLT(FieldReviewTime, v))
}

// ReviewTimeLTE applies the LTE predicate on the "review_time" field.
func ReviewTimeLTE(v time.Time) predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldLTE(FieldReviewTime, v))
}

// ReviewTimeIsNil applies the IsNil predicate on the "review_time" field.
func ReviewTimeIsNil() predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldIsNull(FieldReviewTime))
}

// ReviewTimeNotNil applies the NotNil predicate on the "review_time" field.
func ReviewTimeNotNil() predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(sql.FieldNotNull(FieldReviewTime))
}

// HasMember applies the HasEdge predicate on the "member" edge.
func HasMember() predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, MemberTable, MemberColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasMemberWith applies the HasEdge predicate on the "member" edge with a given conditions (other predicates).
func HasMemberWith(preds ...predicate.PromotionMember) predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(func(s *sql.Selector) {
		step := newMemberStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasCards applies the HasEdge predicate on the "cards" edge.
func HasCards() predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, CardsTable, CardsColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasCardsWith applies the HasEdge predicate on the "cards" edge with a given conditions (other predicates).
func HasCardsWith(preds ...predicate.PromotionBankCard) predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(func(s *sql.Selector) {
		step := newCardsStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.PromotionWithdrawal) predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for _, p := range predicates {
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.PromotionWithdrawal) predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for i, p := range predicates {
			if i > 0 {
				s1.Or()
			}
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Not applies the not operator on the given predicate.
func Not(p predicate.PromotionWithdrawal) predicate.PromotionWithdrawal {
	return predicate.PromotionWithdrawal(func(s *sql.Selector) {
		p(s.Not())
	})
}