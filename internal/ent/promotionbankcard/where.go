// Code generated by ent, DO NOT EDIT.

package promotionbankcard

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/auroraride/aurservd/internal/ent/predicate"
)

// ID filters vertices based on their ID field.
func ID(id uint64) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id uint64) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id uint64) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...uint64) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...uint64) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id uint64) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id uint64) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id uint64) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id uint64) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldLTE(FieldID, id))
}

// CreatedAt applies equality check predicate on the "created_at" field. It's identical to CreatedAtEQ.
func CreatedAt(v time.Time) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldEQ(FieldCreatedAt, v))
}

// UpdatedAt applies equality check predicate on the "updated_at" field. It's identical to UpdatedAtEQ.
func UpdatedAt(v time.Time) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldEQ(FieldUpdatedAt, v))
}

// DeletedAt applies equality check predicate on the "deleted_at" field. It's identical to DeletedAtEQ.
func DeletedAt(v time.Time) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldEQ(FieldDeletedAt, v))
}

// Remark applies equality check predicate on the "remark" field. It's identical to RemarkEQ.
func Remark(v string) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldEQ(FieldRemark, v))
}

// MemberID applies equality check predicate on the "member_id" field. It's identical to MemberIDEQ.
func MemberID(v uint64) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldEQ(FieldMemberID, v))
}

// CardNo applies equality check predicate on the "card_no" field. It's identical to CardNoEQ.
func CardNo(v string) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldEQ(FieldCardNo, v))
}

// Bank applies equality check predicate on the "bank" field. It's identical to BankEQ.
func Bank(v string) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldEQ(FieldBank, v))
}

// IsDefault applies equality check predicate on the "is_default" field. It's identical to IsDefaultEQ.
func IsDefault(v bool) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldEQ(FieldIsDefault, v))
}

// BankLogoURL applies equality check predicate on the "bank_logo_url" field. It's identical to BankLogoURLEQ.
func BankLogoURL(v string) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldEQ(FieldBankLogoURL, v))
}

// Province applies equality check predicate on the "province" field. It's identical to ProvinceEQ.
func Province(v string) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldEQ(FieldProvince, v))
}

// City applies equality check predicate on the "city" field. It's identical to CityEQ.
func City(v string) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldEQ(FieldCity, v))
}

// CreatedAtEQ applies the EQ predicate on the "created_at" field.
func CreatedAtEQ(v time.Time) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldEQ(FieldCreatedAt, v))
}

// CreatedAtNEQ applies the NEQ predicate on the "created_at" field.
func CreatedAtNEQ(v time.Time) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldNEQ(FieldCreatedAt, v))
}

// CreatedAtIn applies the In predicate on the "created_at" field.
func CreatedAtIn(vs ...time.Time) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldIn(FieldCreatedAt, vs...))
}

// CreatedAtNotIn applies the NotIn predicate on the "created_at" field.
func CreatedAtNotIn(vs ...time.Time) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldNotIn(FieldCreatedAt, vs...))
}

// CreatedAtGT applies the GT predicate on the "created_at" field.
func CreatedAtGT(v time.Time) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldGT(FieldCreatedAt, v))
}

// CreatedAtGTE applies the GTE predicate on the "created_at" field.
func CreatedAtGTE(v time.Time) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldGTE(FieldCreatedAt, v))
}

// CreatedAtLT applies the LT predicate on the "created_at" field.
func CreatedAtLT(v time.Time) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldLT(FieldCreatedAt, v))
}

// CreatedAtLTE applies the LTE predicate on the "created_at" field.
func CreatedAtLTE(v time.Time) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldLTE(FieldCreatedAt, v))
}

// UpdatedAtEQ applies the EQ predicate on the "updated_at" field.
func UpdatedAtEQ(v time.Time) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldEQ(FieldUpdatedAt, v))
}

// UpdatedAtNEQ applies the NEQ predicate on the "updated_at" field.
func UpdatedAtNEQ(v time.Time) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldNEQ(FieldUpdatedAt, v))
}

// UpdatedAtIn applies the In predicate on the "updated_at" field.
func UpdatedAtIn(vs ...time.Time) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldIn(FieldUpdatedAt, vs...))
}

// UpdatedAtNotIn applies the NotIn predicate on the "updated_at" field.
func UpdatedAtNotIn(vs ...time.Time) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldNotIn(FieldUpdatedAt, vs...))
}

// UpdatedAtGT applies the GT predicate on the "updated_at" field.
func UpdatedAtGT(v time.Time) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldGT(FieldUpdatedAt, v))
}

// UpdatedAtGTE applies the GTE predicate on the "updated_at" field.
func UpdatedAtGTE(v time.Time) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldGTE(FieldUpdatedAt, v))
}

// UpdatedAtLT applies the LT predicate on the "updated_at" field.
func UpdatedAtLT(v time.Time) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldLT(FieldUpdatedAt, v))
}

// UpdatedAtLTE applies the LTE predicate on the "updated_at" field.
func UpdatedAtLTE(v time.Time) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldLTE(FieldUpdatedAt, v))
}

// DeletedAtEQ applies the EQ predicate on the "deleted_at" field.
func DeletedAtEQ(v time.Time) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldEQ(FieldDeletedAt, v))
}

// DeletedAtNEQ applies the NEQ predicate on the "deleted_at" field.
func DeletedAtNEQ(v time.Time) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldNEQ(FieldDeletedAt, v))
}

// DeletedAtIn applies the In predicate on the "deleted_at" field.
func DeletedAtIn(vs ...time.Time) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldIn(FieldDeletedAt, vs...))
}

// DeletedAtNotIn applies the NotIn predicate on the "deleted_at" field.
func DeletedAtNotIn(vs ...time.Time) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldNotIn(FieldDeletedAt, vs...))
}

// DeletedAtGT applies the GT predicate on the "deleted_at" field.
func DeletedAtGT(v time.Time) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldGT(FieldDeletedAt, v))
}

// DeletedAtGTE applies the GTE predicate on the "deleted_at" field.
func DeletedAtGTE(v time.Time) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldGTE(FieldDeletedAt, v))
}

// DeletedAtLT applies the LT predicate on the "deleted_at" field.
func DeletedAtLT(v time.Time) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldLT(FieldDeletedAt, v))
}

// DeletedAtLTE applies the LTE predicate on the "deleted_at" field.
func DeletedAtLTE(v time.Time) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldLTE(FieldDeletedAt, v))
}

// DeletedAtIsNil applies the IsNil predicate on the "deleted_at" field.
func DeletedAtIsNil() predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldIsNull(FieldDeletedAt))
}

// DeletedAtNotNil applies the NotNil predicate on the "deleted_at" field.
func DeletedAtNotNil() predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldNotNull(FieldDeletedAt))
}

// CreatorIsNil applies the IsNil predicate on the "creator" field.
func CreatorIsNil() predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldIsNull(FieldCreator))
}

// CreatorNotNil applies the NotNil predicate on the "creator" field.
func CreatorNotNil() predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldNotNull(FieldCreator))
}

// LastModifierIsNil applies the IsNil predicate on the "last_modifier" field.
func LastModifierIsNil() predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldIsNull(FieldLastModifier))
}

// LastModifierNotNil applies the NotNil predicate on the "last_modifier" field.
func LastModifierNotNil() predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldNotNull(FieldLastModifier))
}

// RemarkEQ applies the EQ predicate on the "remark" field.
func RemarkEQ(v string) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldEQ(FieldRemark, v))
}

// RemarkNEQ applies the NEQ predicate on the "remark" field.
func RemarkNEQ(v string) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldNEQ(FieldRemark, v))
}

// RemarkIn applies the In predicate on the "remark" field.
func RemarkIn(vs ...string) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldIn(FieldRemark, vs...))
}

// RemarkNotIn applies the NotIn predicate on the "remark" field.
func RemarkNotIn(vs ...string) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldNotIn(FieldRemark, vs...))
}

// RemarkGT applies the GT predicate on the "remark" field.
func RemarkGT(v string) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldGT(FieldRemark, v))
}

// RemarkGTE applies the GTE predicate on the "remark" field.
func RemarkGTE(v string) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldGTE(FieldRemark, v))
}

// RemarkLT applies the LT predicate on the "remark" field.
func RemarkLT(v string) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldLT(FieldRemark, v))
}

// RemarkLTE applies the LTE predicate on the "remark" field.
func RemarkLTE(v string) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldLTE(FieldRemark, v))
}

// RemarkContains applies the Contains predicate on the "remark" field.
func RemarkContains(v string) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldContains(FieldRemark, v))
}

// RemarkHasPrefix applies the HasPrefix predicate on the "remark" field.
func RemarkHasPrefix(v string) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldHasPrefix(FieldRemark, v))
}

// RemarkHasSuffix applies the HasSuffix predicate on the "remark" field.
func RemarkHasSuffix(v string) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldHasSuffix(FieldRemark, v))
}

// RemarkIsNil applies the IsNil predicate on the "remark" field.
func RemarkIsNil() predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldIsNull(FieldRemark))
}

// RemarkNotNil applies the NotNil predicate on the "remark" field.
func RemarkNotNil() predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldNotNull(FieldRemark))
}

// RemarkEqualFold applies the EqualFold predicate on the "remark" field.
func RemarkEqualFold(v string) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldEqualFold(FieldRemark, v))
}

// RemarkContainsFold applies the ContainsFold predicate on the "remark" field.
func RemarkContainsFold(v string) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldContainsFold(FieldRemark, v))
}

// MemberIDEQ applies the EQ predicate on the "member_id" field.
func MemberIDEQ(v uint64) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldEQ(FieldMemberID, v))
}

// MemberIDNEQ applies the NEQ predicate on the "member_id" field.
func MemberIDNEQ(v uint64) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldNEQ(FieldMemberID, v))
}

// MemberIDIn applies the In predicate on the "member_id" field.
func MemberIDIn(vs ...uint64) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldIn(FieldMemberID, vs...))
}

// MemberIDNotIn applies the NotIn predicate on the "member_id" field.
func MemberIDNotIn(vs ...uint64) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldNotIn(FieldMemberID, vs...))
}

// MemberIDIsNil applies the IsNil predicate on the "member_id" field.
func MemberIDIsNil() predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldIsNull(FieldMemberID))
}

// MemberIDNotNil applies the NotNil predicate on the "member_id" field.
func MemberIDNotNil() predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldNotNull(FieldMemberID))
}

// CardNoEQ applies the EQ predicate on the "card_no" field.
func CardNoEQ(v string) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldEQ(FieldCardNo, v))
}

// CardNoNEQ applies the NEQ predicate on the "card_no" field.
func CardNoNEQ(v string) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldNEQ(FieldCardNo, v))
}

// CardNoIn applies the In predicate on the "card_no" field.
func CardNoIn(vs ...string) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldIn(FieldCardNo, vs...))
}

// CardNoNotIn applies the NotIn predicate on the "card_no" field.
func CardNoNotIn(vs ...string) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldNotIn(FieldCardNo, vs...))
}

// CardNoGT applies the GT predicate on the "card_no" field.
func CardNoGT(v string) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldGT(FieldCardNo, v))
}

// CardNoGTE applies the GTE predicate on the "card_no" field.
func CardNoGTE(v string) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldGTE(FieldCardNo, v))
}

// CardNoLT applies the LT predicate on the "card_no" field.
func CardNoLT(v string) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldLT(FieldCardNo, v))
}

// CardNoLTE applies the LTE predicate on the "card_no" field.
func CardNoLTE(v string) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldLTE(FieldCardNo, v))
}

// CardNoContains applies the Contains predicate on the "card_no" field.
func CardNoContains(v string) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldContains(FieldCardNo, v))
}

// CardNoHasPrefix applies the HasPrefix predicate on the "card_no" field.
func CardNoHasPrefix(v string) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldHasPrefix(FieldCardNo, v))
}

// CardNoHasSuffix applies the HasSuffix predicate on the "card_no" field.
func CardNoHasSuffix(v string) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldHasSuffix(FieldCardNo, v))
}

// CardNoEqualFold applies the EqualFold predicate on the "card_no" field.
func CardNoEqualFold(v string) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldEqualFold(FieldCardNo, v))
}

// CardNoContainsFold applies the ContainsFold predicate on the "card_no" field.
func CardNoContainsFold(v string) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldContainsFold(FieldCardNo, v))
}

// BankEQ applies the EQ predicate on the "bank" field.
func BankEQ(v string) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldEQ(FieldBank, v))
}

// BankNEQ applies the NEQ predicate on the "bank" field.
func BankNEQ(v string) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldNEQ(FieldBank, v))
}

// BankIn applies the In predicate on the "bank" field.
func BankIn(vs ...string) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldIn(FieldBank, vs...))
}

// BankNotIn applies the NotIn predicate on the "bank" field.
func BankNotIn(vs ...string) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldNotIn(FieldBank, vs...))
}

// BankGT applies the GT predicate on the "bank" field.
func BankGT(v string) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldGT(FieldBank, v))
}

// BankGTE applies the GTE predicate on the "bank" field.
func BankGTE(v string) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldGTE(FieldBank, v))
}

// BankLT applies the LT predicate on the "bank" field.
func BankLT(v string) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldLT(FieldBank, v))
}

// BankLTE applies the LTE predicate on the "bank" field.
func BankLTE(v string) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldLTE(FieldBank, v))
}

// BankContains applies the Contains predicate on the "bank" field.
func BankContains(v string) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldContains(FieldBank, v))
}

// BankHasPrefix applies the HasPrefix predicate on the "bank" field.
func BankHasPrefix(v string) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldHasPrefix(FieldBank, v))
}

// BankHasSuffix applies the HasSuffix predicate on the "bank" field.
func BankHasSuffix(v string) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldHasSuffix(FieldBank, v))
}

// BankIsNil applies the IsNil predicate on the "bank" field.
func BankIsNil() predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldIsNull(FieldBank))
}

// BankNotNil applies the NotNil predicate on the "bank" field.
func BankNotNil() predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldNotNull(FieldBank))
}

// BankEqualFold applies the EqualFold predicate on the "bank" field.
func BankEqualFold(v string) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldEqualFold(FieldBank, v))
}

// BankContainsFold applies the ContainsFold predicate on the "bank" field.
func BankContainsFold(v string) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldContainsFold(FieldBank, v))
}

// IsDefaultEQ applies the EQ predicate on the "is_default" field.
func IsDefaultEQ(v bool) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldEQ(FieldIsDefault, v))
}

// IsDefaultNEQ applies the NEQ predicate on the "is_default" field.
func IsDefaultNEQ(v bool) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldNEQ(FieldIsDefault, v))
}

// BankLogoURLEQ applies the EQ predicate on the "bank_logo_url" field.
func BankLogoURLEQ(v string) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldEQ(FieldBankLogoURL, v))
}

// BankLogoURLNEQ applies the NEQ predicate on the "bank_logo_url" field.
func BankLogoURLNEQ(v string) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldNEQ(FieldBankLogoURL, v))
}

// BankLogoURLIn applies the In predicate on the "bank_logo_url" field.
func BankLogoURLIn(vs ...string) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldIn(FieldBankLogoURL, vs...))
}

// BankLogoURLNotIn applies the NotIn predicate on the "bank_logo_url" field.
func BankLogoURLNotIn(vs ...string) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldNotIn(FieldBankLogoURL, vs...))
}

// BankLogoURLGT applies the GT predicate on the "bank_logo_url" field.
func BankLogoURLGT(v string) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldGT(FieldBankLogoURL, v))
}

// BankLogoURLGTE applies the GTE predicate on the "bank_logo_url" field.
func BankLogoURLGTE(v string) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldGTE(FieldBankLogoURL, v))
}

// BankLogoURLLT applies the LT predicate on the "bank_logo_url" field.
func BankLogoURLLT(v string) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldLT(FieldBankLogoURL, v))
}

// BankLogoURLLTE applies the LTE predicate on the "bank_logo_url" field.
func BankLogoURLLTE(v string) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldLTE(FieldBankLogoURL, v))
}

// BankLogoURLContains applies the Contains predicate on the "bank_logo_url" field.
func BankLogoURLContains(v string) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldContains(FieldBankLogoURL, v))
}

// BankLogoURLHasPrefix applies the HasPrefix predicate on the "bank_logo_url" field.
func BankLogoURLHasPrefix(v string) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldHasPrefix(FieldBankLogoURL, v))
}

// BankLogoURLHasSuffix applies the HasSuffix predicate on the "bank_logo_url" field.
func BankLogoURLHasSuffix(v string) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldHasSuffix(FieldBankLogoURL, v))
}

// BankLogoURLIsNil applies the IsNil predicate on the "bank_logo_url" field.
func BankLogoURLIsNil() predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldIsNull(FieldBankLogoURL))
}

// BankLogoURLNotNil applies the NotNil predicate on the "bank_logo_url" field.
func BankLogoURLNotNil() predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldNotNull(FieldBankLogoURL))
}

// BankLogoURLEqualFold applies the EqualFold predicate on the "bank_logo_url" field.
func BankLogoURLEqualFold(v string) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldEqualFold(FieldBankLogoURL, v))
}

// BankLogoURLContainsFold applies the ContainsFold predicate on the "bank_logo_url" field.
func BankLogoURLContainsFold(v string) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldContainsFold(FieldBankLogoURL, v))
}

// ProvinceEQ applies the EQ predicate on the "province" field.
func ProvinceEQ(v string) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldEQ(FieldProvince, v))
}

// ProvinceNEQ applies the NEQ predicate on the "province" field.
func ProvinceNEQ(v string) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldNEQ(FieldProvince, v))
}

// ProvinceIn applies the In predicate on the "province" field.
func ProvinceIn(vs ...string) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldIn(FieldProvince, vs...))
}

// ProvinceNotIn applies the NotIn predicate on the "province" field.
func ProvinceNotIn(vs ...string) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldNotIn(FieldProvince, vs...))
}

// ProvinceGT applies the GT predicate on the "province" field.
func ProvinceGT(v string) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldGT(FieldProvince, v))
}

// ProvinceGTE applies the GTE predicate on the "province" field.
func ProvinceGTE(v string) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldGTE(FieldProvince, v))
}

// ProvinceLT applies the LT predicate on the "province" field.
func ProvinceLT(v string) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldLT(FieldProvince, v))
}

// ProvinceLTE applies the LTE predicate on the "province" field.
func ProvinceLTE(v string) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldLTE(FieldProvince, v))
}

// ProvinceContains applies the Contains predicate on the "province" field.
func ProvinceContains(v string) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldContains(FieldProvince, v))
}

// ProvinceHasPrefix applies the HasPrefix predicate on the "province" field.
func ProvinceHasPrefix(v string) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldHasPrefix(FieldProvince, v))
}

// ProvinceHasSuffix applies the HasSuffix predicate on the "province" field.
func ProvinceHasSuffix(v string) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldHasSuffix(FieldProvince, v))
}

// ProvinceIsNil applies the IsNil predicate on the "province" field.
func ProvinceIsNil() predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldIsNull(FieldProvince))
}

// ProvinceNotNil applies the NotNil predicate on the "province" field.
func ProvinceNotNil() predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldNotNull(FieldProvince))
}

// ProvinceEqualFold applies the EqualFold predicate on the "province" field.
func ProvinceEqualFold(v string) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldEqualFold(FieldProvince, v))
}

// ProvinceContainsFold applies the ContainsFold predicate on the "province" field.
func ProvinceContainsFold(v string) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldContainsFold(FieldProvince, v))
}

// CityEQ applies the EQ predicate on the "city" field.
func CityEQ(v string) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldEQ(FieldCity, v))
}

// CityNEQ applies the NEQ predicate on the "city" field.
func CityNEQ(v string) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldNEQ(FieldCity, v))
}

// CityIn applies the In predicate on the "city" field.
func CityIn(vs ...string) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldIn(FieldCity, vs...))
}

// CityNotIn applies the NotIn predicate on the "city" field.
func CityNotIn(vs ...string) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldNotIn(FieldCity, vs...))
}

// CityGT applies the GT predicate on the "city" field.
func CityGT(v string) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldGT(FieldCity, v))
}

// CityGTE applies the GTE predicate on the "city" field.
func CityGTE(v string) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldGTE(FieldCity, v))
}

// CityLT applies the LT predicate on the "city" field.
func CityLT(v string) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldLT(FieldCity, v))
}

// CityLTE applies the LTE predicate on the "city" field.
func CityLTE(v string) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldLTE(FieldCity, v))
}

// CityContains applies the Contains predicate on the "city" field.
func CityContains(v string) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldContains(FieldCity, v))
}

// CityHasPrefix applies the HasPrefix predicate on the "city" field.
func CityHasPrefix(v string) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldHasPrefix(FieldCity, v))
}

// CityHasSuffix applies the HasSuffix predicate on the "city" field.
func CityHasSuffix(v string) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldHasSuffix(FieldCity, v))
}

// CityIsNil applies the IsNil predicate on the "city" field.
func CityIsNil() predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldIsNull(FieldCity))
}

// CityNotNil applies the NotNil predicate on the "city" field.
func CityNotNil() predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldNotNull(FieldCity))
}

// CityEqualFold applies the EqualFold predicate on the "city" field.
func CityEqualFold(v string) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldEqualFold(FieldCity, v))
}

// CityContainsFold applies the ContainsFold predicate on the "city" field.
func CityContainsFold(v string) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(sql.FieldContainsFold(FieldCity, v))
}

// HasMember applies the HasEdge predicate on the "member" edge.
func HasMember() predicate.PromotionBankCard {
	return predicate.PromotionBankCard(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, MemberTable, MemberColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasMemberWith applies the HasEdge predicate on the "member" edge with a given conditions (other predicates).
func HasMemberWith(preds ...predicate.PromotionMember) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(func(s *sql.Selector) {
		step := newMemberStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasWithdrawals applies the HasEdge predicate on the "withdrawals" edge.
func HasWithdrawals() predicate.PromotionBankCard {
	return predicate.PromotionBankCard(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, WithdrawalsTable, WithdrawalsColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasWithdrawalsWith applies the HasEdge predicate on the "withdrawals" edge with a given conditions (other predicates).
func HasWithdrawalsWith(preds ...predicate.PromotionWithdrawal) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(func(s *sql.Selector) {
		step := newWithdrawalsStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.PromotionBankCard) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for _, p := range predicates {
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.PromotionBankCard) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(func(s *sql.Selector) {
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
func Not(p predicate.PromotionBankCard) predicate.PromotionBankCard {
	return predicate.PromotionBankCard(func(s *sql.Selector) {
		p(s.Not())
	})
}