// Code generated by ent, DO NOT EDIT.

package promotionperson

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/auroraride/aurservd/internal/ent/predicate"
)

// ID filters vertices based on their ID field.
func ID(id uint64) predicate.PromotionPerson {
	return predicate.PromotionPerson(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id uint64) predicate.PromotionPerson {
	return predicate.PromotionPerson(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id uint64) predicate.PromotionPerson {
	return predicate.PromotionPerson(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...uint64) predicate.PromotionPerson {
	return predicate.PromotionPerson(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...uint64) predicate.PromotionPerson {
	return predicate.PromotionPerson(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id uint64) predicate.PromotionPerson {
	return predicate.PromotionPerson(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id uint64) predicate.PromotionPerson {
	return predicate.PromotionPerson(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id uint64) predicate.PromotionPerson {
	return predicate.PromotionPerson(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id uint64) predicate.PromotionPerson {
	return predicate.PromotionPerson(sql.FieldLTE(FieldID, id))
}

// CreatedAt applies equality check predicate on the "created_at" field. It's identical to CreatedAtEQ.
func CreatedAt(v time.Time) predicate.PromotionPerson {
	return predicate.PromotionPerson(sql.FieldEQ(FieldCreatedAt, v))
}

// UpdatedAt applies equality check predicate on the "updated_at" field. It's identical to UpdatedAtEQ.
func UpdatedAt(v time.Time) predicate.PromotionPerson {
	return predicate.PromotionPerson(sql.FieldEQ(FieldUpdatedAt, v))
}

// DeletedAt applies equality check predicate on the "deleted_at" field. It's identical to DeletedAtEQ.
func DeletedAt(v time.Time) predicate.PromotionPerson {
	return predicate.PromotionPerson(sql.FieldEQ(FieldDeletedAt, v))
}

// Remark applies equality check predicate on the "remark" field. It's identical to RemarkEQ.
func Remark(v string) predicate.PromotionPerson {
	return predicate.PromotionPerson(sql.FieldEQ(FieldRemark, v))
}

// Status applies equality check predicate on the "status" field. It's identical to StatusEQ.
func Status(v uint8) predicate.PromotionPerson {
	return predicate.PromotionPerson(sql.FieldEQ(FieldStatus, v))
}

// Name applies equality check predicate on the "name" field. It's identical to NameEQ.
func Name(v string) predicate.PromotionPerson {
	return predicate.PromotionPerson(sql.FieldEQ(FieldName, v))
}

// IDCardNumber applies equality check predicate on the "id_card_number" field. It's identical to IDCardNumberEQ.
func IDCardNumber(v string) predicate.PromotionPerson {
	return predicate.PromotionPerson(sql.FieldEQ(FieldIDCardNumber, v))
}

// Address applies equality check predicate on the "address" field. It's identical to AddressEQ.
func Address(v string) predicate.PromotionPerson {
	return predicate.PromotionPerson(sql.FieldEQ(FieldAddress, v))
}

// CreatedAtEQ applies the EQ predicate on the "created_at" field.
func CreatedAtEQ(v time.Time) predicate.PromotionPerson {
	return predicate.PromotionPerson(sql.FieldEQ(FieldCreatedAt, v))
}

// CreatedAtNEQ applies the NEQ predicate on the "created_at" field.
func CreatedAtNEQ(v time.Time) predicate.PromotionPerson {
	return predicate.PromotionPerson(sql.FieldNEQ(FieldCreatedAt, v))
}

// CreatedAtIn applies the In predicate on the "created_at" field.
func CreatedAtIn(vs ...time.Time) predicate.PromotionPerson {
	return predicate.PromotionPerson(sql.FieldIn(FieldCreatedAt, vs...))
}

// CreatedAtNotIn applies the NotIn predicate on the "created_at" field.
func CreatedAtNotIn(vs ...time.Time) predicate.PromotionPerson {
	return predicate.PromotionPerson(sql.FieldNotIn(FieldCreatedAt, vs...))
}

// CreatedAtGT applies the GT predicate on the "created_at" field.
func CreatedAtGT(v time.Time) predicate.PromotionPerson {
	return predicate.PromotionPerson(sql.FieldGT(FieldCreatedAt, v))
}

// CreatedAtGTE applies the GTE predicate on the "created_at" field.
func CreatedAtGTE(v time.Time) predicate.PromotionPerson {
	return predicate.PromotionPerson(sql.FieldGTE(FieldCreatedAt, v))
}

// CreatedAtLT applies the LT predicate on the "created_at" field.
func CreatedAtLT(v time.Time) predicate.PromotionPerson {
	return predicate.PromotionPerson(sql.FieldLT(FieldCreatedAt, v))
}

// CreatedAtLTE applies the LTE predicate on the "created_at" field.
func CreatedAtLTE(v time.Time) predicate.PromotionPerson {
	return predicate.PromotionPerson(sql.FieldLTE(FieldCreatedAt, v))
}

// UpdatedAtEQ applies the EQ predicate on the "updated_at" field.
func UpdatedAtEQ(v time.Time) predicate.PromotionPerson {
	return predicate.PromotionPerson(sql.FieldEQ(FieldUpdatedAt, v))
}

// UpdatedAtNEQ applies the NEQ predicate on the "updated_at" field.
func UpdatedAtNEQ(v time.Time) predicate.PromotionPerson {
	return predicate.PromotionPerson(sql.FieldNEQ(FieldUpdatedAt, v))
}

// UpdatedAtIn applies the In predicate on the "updated_at" field.
func UpdatedAtIn(vs ...time.Time) predicate.PromotionPerson {
	return predicate.PromotionPerson(sql.FieldIn(FieldUpdatedAt, vs...))
}

// UpdatedAtNotIn applies the NotIn predicate on the "updated_at" field.
func UpdatedAtNotIn(vs ...time.Time) predicate.PromotionPerson {
	return predicate.PromotionPerson(sql.FieldNotIn(FieldUpdatedAt, vs...))
}

// UpdatedAtGT applies the GT predicate on the "updated_at" field.
func UpdatedAtGT(v time.Time) predicate.PromotionPerson {
	return predicate.PromotionPerson(sql.FieldGT(FieldUpdatedAt, v))
}

// UpdatedAtGTE applies the GTE predicate on the "updated_at" field.
func UpdatedAtGTE(v time.Time) predicate.PromotionPerson {
	return predicate.PromotionPerson(sql.FieldGTE(FieldUpdatedAt, v))
}

// UpdatedAtLT applies the LT predicate on the "updated_at" field.
func UpdatedAtLT(v time.Time) predicate.PromotionPerson {
	return predicate.PromotionPerson(sql.FieldLT(FieldUpdatedAt, v))
}

// UpdatedAtLTE applies the LTE predicate on the "updated_at" field.
func UpdatedAtLTE(v time.Time) predicate.PromotionPerson {
	return predicate.PromotionPerson(sql.FieldLTE(FieldUpdatedAt, v))
}

// DeletedAtEQ applies the EQ predicate on the "deleted_at" field.
func DeletedAtEQ(v time.Time) predicate.PromotionPerson {
	return predicate.PromotionPerson(sql.FieldEQ(FieldDeletedAt, v))
}

// DeletedAtNEQ applies the NEQ predicate on the "deleted_at" field.
func DeletedAtNEQ(v time.Time) predicate.PromotionPerson {
	return predicate.PromotionPerson(sql.FieldNEQ(FieldDeletedAt, v))
}

// DeletedAtIn applies the In predicate on the "deleted_at" field.
func DeletedAtIn(vs ...time.Time) predicate.PromotionPerson {
	return predicate.PromotionPerson(sql.FieldIn(FieldDeletedAt, vs...))
}

// DeletedAtNotIn applies the NotIn predicate on the "deleted_at" field.
func DeletedAtNotIn(vs ...time.Time) predicate.PromotionPerson {
	return predicate.PromotionPerson(sql.FieldNotIn(FieldDeletedAt, vs...))
}

// DeletedAtGT applies the GT predicate on the "deleted_at" field.
func DeletedAtGT(v time.Time) predicate.PromotionPerson {
	return predicate.PromotionPerson(sql.FieldGT(FieldDeletedAt, v))
}

// DeletedAtGTE applies the GTE predicate on the "deleted_at" field.
func DeletedAtGTE(v time.Time) predicate.PromotionPerson {
	return predicate.PromotionPerson(sql.FieldGTE(FieldDeletedAt, v))
}

// DeletedAtLT applies the LT predicate on the "deleted_at" field.
func DeletedAtLT(v time.Time) predicate.PromotionPerson {
	return predicate.PromotionPerson(sql.FieldLT(FieldDeletedAt, v))
}

// DeletedAtLTE applies the LTE predicate on the "deleted_at" field.
func DeletedAtLTE(v time.Time) predicate.PromotionPerson {
	return predicate.PromotionPerson(sql.FieldLTE(FieldDeletedAt, v))
}

// DeletedAtIsNil applies the IsNil predicate on the "deleted_at" field.
func DeletedAtIsNil() predicate.PromotionPerson {
	return predicate.PromotionPerson(sql.FieldIsNull(FieldDeletedAt))
}

// DeletedAtNotNil applies the NotNil predicate on the "deleted_at" field.
func DeletedAtNotNil() predicate.PromotionPerson {
	return predicate.PromotionPerson(sql.FieldNotNull(FieldDeletedAt))
}

// CreatorIsNil applies the IsNil predicate on the "creator" field.
func CreatorIsNil() predicate.PromotionPerson {
	return predicate.PromotionPerson(sql.FieldIsNull(FieldCreator))
}

// CreatorNotNil applies the NotNil predicate on the "creator" field.
func CreatorNotNil() predicate.PromotionPerson {
	return predicate.PromotionPerson(sql.FieldNotNull(FieldCreator))
}

// LastModifierIsNil applies the IsNil predicate on the "last_modifier" field.
func LastModifierIsNil() predicate.PromotionPerson {
	return predicate.PromotionPerson(sql.FieldIsNull(FieldLastModifier))
}

// LastModifierNotNil applies the NotNil predicate on the "last_modifier" field.
func LastModifierNotNil() predicate.PromotionPerson {
	return predicate.PromotionPerson(sql.FieldNotNull(FieldLastModifier))
}

// RemarkEQ applies the EQ predicate on the "remark" field.
func RemarkEQ(v string) predicate.PromotionPerson {
	return predicate.PromotionPerson(sql.FieldEQ(FieldRemark, v))
}

// RemarkNEQ applies the NEQ predicate on the "remark" field.
func RemarkNEQ(v string) predicate.PromotionPerson {
	return predicate.PromotionPerson(sql.FieldNEQ(FieldRemark, v))
}

// RemarkIn applies the In predicate on the "remark" field.
func RemarkIn(vs ...string) predicate.PromotionPerson {
	return predicate.PromotionPerson(sql.FieldIn(FieldRemark, vs...))
}

// RemarkNotIn applies the NotIn predicate on the "remark" field.
func RemarkNotIn(vs ...string) predicate.PromotionPerson {
	return predicate.PromotionPerson(sql.FieldNotIn(FieldRemark, vs...))
}

// RemarkGT applies the GT predicate on the "remark" field.
func RemarkGT(v string) predicate.PromotionPerson {
	return predicate.PromotionPerson(sql.FieldGT(FieldRemark, v))
}

// RemarkGTE applies the GTE predicate on the "remark" field.
func RemarkGTE(v string) predicate.PromotionPerson {
	return predicate.PromotionPerson(sql.FieldGTE(FieldRemark, v))
}

// RemarkLT applies the LT predicate on the "remark" field.
func RemarkLT(v string) predicate.PromotionPerson {
	return predicate.PromotionPerson(sql.FieldLT(FieldRemark, v))
}

// RemarkLTE applies the LTE predicate on the "remark" field.
func RemarkLTE(v string) predicate.PromotionPerson {
	return predicate.PromotionPerson(sql.FieldLTE(FieldRemark, v))
}

// RemarkContains applies the Contains predicate on the "remark" field.
func RemarkContains(v string) predicate.PromotionPerson {
	return predicate.PromotionPerson(sql.FieldContains(FieldRemark, v))
}

// RemarkHasPrefix applies the HasPrefix predicate on the "remark" field.
func RemarkHasPrefix(v string) predicate.PromotionPerson {
	return predicate.PromotionPerson(sql.FieldHasPrefix(FieldRemark, v))
}

// RemarkHasSuffix applies the HasSuffix predicate on the "remark" field.
func RemarkHasSuffix(v string) predicate.PromotionPerson {
	return predicate.PromotionPerson(sql.FieldHasSuffix(FieldRemark, v))
}

// RemarkIsNil applies the IsNil predicate on the "remark" field.
func RemarkIsNil() predicate.PromotionPerson {
	return predicate.PromotionPerson(sql.FieldIsNull(FieldRemark))
}

// RemarkNotNil applies the NotNil predicate on the "remark" field.
func RemarkNotNil() predicate.PromotionPerson {
	return predicate.PromotionPerson(sql.FieldNotNull(FieldRemark))
}

// RemarkEqualFold applies the EqualFold predicate on the "remark" field.
func RemarkEqualFold(v string) predicate.PromotionPerson {
	return predicate.PromotionPerson(sql.FieldEqualFold(FieldRemark, v))
}

// RemarkContainsFold applies the ContainsFold predicate on the "remark" field.
func RemarkContainsFold(v string) predicate.PromotionPerson {
	return predicate.PromotionPerson(sql.FieldContainsFold(FieldRemark, v))
}

// StatusEQ applies the EQ predicate on the "status" field.
func StatusEQ(v uint8) predicate.PromotionPerson {
	return predicate.PromotionPerson(sql.FieldEQ(FieldStatus, v))
}

// StatusNEQ applies the NEQ predicate on the "status" field.
func StatusNEQ(v uint8) predicate.PromotionPerson {
	return predicate.PromotionPerson(sql.FieldNEQ(FieldStatus, v))
}

// StatusIn applies the In predicate on the "status" field.
func StatusIn(vs ...uint8) predicate.PromotionPerson {
	return predicate.PromotionPerson(sql.FieldIn(FieldStatus, vs...))
}

// StatusNotIn applies the NotIn predicate on the "status" field.
func StatusNotIn(vs ...uint8) predicate.PromotionPerson {
	return predicate.PromotionPerson(sql.FieldNotIn(FieldStatus, vs...))
}

// StatusGT applies the GT predicate on the "status" field.
func StatusGT(v uint8) predicate.PromotionPerson {
	return predicate.PromotionPerson(sql.FieldGT(FieldStatus, v))
}

// StatusGTE applies the GTE predicate on the "status" field.
func StatusGTE(v uint8) predicate.PromotionPerson {
	return predicate.PromotionPerson(sql.FieldGTE(FieldStatus, v))
}

// StatusLT applies the LT predicate on the "status" field.
func StatusLT(v uint8) predicate.PromotionPerson {
	return predicate.PromotionPerson(sql.FieldLT(FieldStatus, v))
}

// StatusLTE applies the LTE predicate on the "status" field.
func StatusLTE(v uint8) predicate.PromotionPerson {
	return predicate.PromotionPerson(sql.FieldLTE(FieldStatus, v))
}

// NameEQ applies the EQ predicate on the "name" field.
func NameEQ(v string) predicate.PromotionPerson {
	return predicate.PromotionPerson(sql.FieldEQ(FieldName, v))
}

// NameNEQ applies the NEQ predicate on the "name" field.
func NameNEQ(v string) predicate.PromotionPerson {
	return predicate.PromotionPerson(sql.FieldNEQ(FieldName, v))
}

// NameIn applies the In predicate on the "name" field.
func NameIn(vs ...string) predicate.PromotionPerson {
	return predicate.PromotionPerson(sql.FieldIn(FieldName, vs...))
}

// NameNotIn applies the NotIn predicate on the "name" field.
func NameNotIn(vs ...string) predicate.PromotionPerson {
	return predicate.PromotionPerson(sql.FieldNotIn(FieldName, vs...))
}

// NameGT applies the GT predicate on the "name" field.
func NameGT(v string) predicate.PromotionPerson {
	return predicate.PromotionPerson(sql.FieldGT(FieldName, v))
}

// NameGTE applies the GTE predicate on the "name" field.
func NameGTE(v string) predicate.PromotionPerson {
	return predicate.PromotionPerson(sql.FieldGTE(FieldName, v))
}

// NameLT applies the LT predicate on the "name" field.
func NameLT(v string) predicate.PromotionPerson {
	return predicate.PromotionPerson(sql.FieldLT(FieldName, v))
}

// NameLTE applies the LTE predicate on the "name" field.
func NameLTE(v string) predicate.PromotionPerson {
	return predicate.PromotionPerson(sql.FieldLTE(FieldName, v))
}

// NameContains applies the Contains predicate on the "name" field.
func NameContains(v string) predicate.PromotionPerson {
	return predicate.PromotionPerson(sql.FieldContains(FieldName, v))
}

// NameHasPrefix applies the HasPrefix predicate on the "name" field.
func NameHasPrefix(v string) predicate.PromotionPerson {
	return predicate.PromotionPerson(sql.FieldHasPrefix(FieldName, v))
}

// NameHasSuffix applies the HasSuffix predicate on the "name" field.
func NameHasSuffix(v string) predicate.PromotionPerson {
	return predicate.PromotionPerson(sql.FieldHasSuffix(FieldName, v))
}

// NameIsNil applies the IsNil predicate on the "name" field.
func NameIsNil() predicate.PromotionPerson {
	return predicate.PromotionPerson(sql.FieldIsNull(FieldName))
}

// NameNotNil applies the NotNil predicate on the "name" field.
func NameNotNil() predicate.PromotionPerson {
	return predicate.PromotionPerson(sql.FieldNotNull(FieldName))
}

// NameEqualFold applies the EqualFold predicate on the "name" field.
func NameEqualFold(v string) predicate.PromotionPerson {
	return predicate.PromotionPerson(sql.FieldEqualFold(FieldName, v))
}

// NameContainsFold applies the ContainsFold predicate on the "name" field.
func NameContainsFold(v string) predicate.PromotionPerson {
	return predicate.PromotionPerson(sql.FieldContainsFold(FieldName, v))
}

// IDCardNumberEQ applies the EQ predicate on the "id_card_number" field.
func IDCardNumberEQ(v string) predicate.PromotionPerson {
	return predicate.PromotionPerson(sql.FieldEQ(FieldIDCardNumber, v))
}

// IDCardNumberNEQ applies the NEQ predicate on the "id_card_number" field.
func IDCardNumberNEQ(v string) predicate.PromotionPerson {
	return predicate.PromotionPerson(sql.FieldNEQ(FieldIDCardNumber, v))
}

// IDCardNumberIn applies the In predicate on the "id_card_number" field.
func IDCardNumberIn(vs ...string) predicate.PromotionPerson {
	return predicate.PromotionPerson(sql.FieldIn(FieldIDCardNumber, vs...))
}

// IDCardNumberNotIn applies the NotIn predicate on the "id_card_number" field.
func IDCardNumberNotIn(vs ...string) predicate.PromotionPerson {
	return predicate.PromotionPerson(sql.FieldNotIn(FieldIDCardNumber, vs...))
}

// IDCardNumberGT applies the GT predicate on the "id_card_number" field.
func IDCardNumberGT(v string) predicate.PromotionPerson {
	return predicate.PromotionPerson(sql.FieldGT(FieldIDCardNumber, v))
}

// IDCardNumberGTE applies the GTE predicate on the "id_card_number" field.
func IDCardNumberGTE(v string) predicate.PromotionPerson {
	return predicate.PromotionPerson(sql.FieldGTE(FieldIDCardNumber, v))
}

// IDCardNumberLT applies the LT predicate on the "id_card_number" field.
func IDCardNumberLT(v string) predicate.PromotionPerson {
	return predicate.PromotionPerson(sql.FieldLT(FieldIDCardNumber, v))
}

// IDCardNumberLTE applies the LTE predicate on the "id_card_number" field.
func IDCardNumberLTE(v string) predicate.PromotionPerson {
	return predicate.PromotionPerson(sql.FieldLTE(FieldIDCardNumber, v))
}

// IDCardNumberContains applies the Contains predicate on the "id_card_number" field.
func IDCardNumberContains(v string) predicate.PromotionPerson {
	return predicate.PromotionPerson(sql.FieldContains(FieldIDCardNumber, v))
}

// IDCardNumberHasPrefix applies the HasPrefix predicate on the "id_card_number" field.
func IDCardNumberHasPrefix(v string) predicate.PromotionPerson {
	return predicate.PromotionPerson(sql.FieldHasPrefix(FieldIDCardNumber, v))
}

// IDCardNumberHasSuffix applies the HasSuffix predicate on the "id_card_number" field.
func IDCardNumberHasSuffix(v string) predicate.PromotionPerson {
	return predicate.PromotionPerson(sql.FieldHasSuffix(FieldIDCardNumber, v))
}

// IDCardNumberIsNil applies the IsNil predicate on the "id_card_number" field.
func IDCardNumberIsNil() predicate.PromotionPerson {
	return predicate.PromotionPerson(sql.FieldIsNull(FieldIDCardNumber))
}

// IDCardNumberNotNil applies the NotNil predicate on the "id_card_number" field.
func IDCardNumberNotNil() predicate.PromotionPerson {
	return predicate.PromotionPerson(sql.FieldNotNull(FieldIDCardNumber))
}

// IDCardNumberEqualFold applies the EqualFold predicate on the "id_card_number" field.
func IDCardNumberEqualFold(v string) predicate.PromotionPerson {
	return predicate.PromotionPerson(sql.FieldEqualFold(FieldIDCardNumber, v))
}

// IDCardNumberContainsFold applies the ContainsFold predicate on the "id_card_number" field.
func IDCardNumberContainsFold(v string) predicate.PromotionPerson {
	return predicate.PromotionPerson(sql.FieldContainsFold(FieldIDCardNumber, v))
}

// AddressEQ applies the EQ predicate on the "address" field.
func AddressEQ(v string) predicate.PromotionPerson {
	return predicate.PromotionPerson(sql.FieldEQ(FieldAddress, v))
}

// AddressNEQ applies the NEQ predicate on the "address" field.
func AddressNEQ(v string) predicate.PromotionPerson {
	return predicate.PromotionPerson(sql.FieldNEQ(FieldAddress, v))
}

// AddressIn applies the In predicate on the "address" field.
func AddressIn(vs ...string) predicate.PromotionPerson {
	return predicate.PromotionPerson(sql.FieldIn(FieldAddress, vs...))
}

// AddressNotIn applies the NotIn predicate on the "address" field.
func AddressNotIn(vs ...string) predicate.PromotionPerson {
	return predicate.PromotionPerson(sql.FieldNotIn(FieldAddress, vs...))
}

// AddressGT applies the GT predicate on the "address" field.
func AddressGT(v string) predicate.PromotionPerson {
	return predicate.PromotionPerson(sql.FieldGT(FieldAddress, v))
}

// AddressGTE applies the GTE predicate on the "address" field.
func AddressGTE(v string) predicate.PromotionPerson {
	return predicate.PromotionPerson(sql.FieldGTE(FieldAddress, v))
}

// AddressLT applies the LT predicate on the "address" field.
func AddressLT(v string) predicate.PromotionPerson {
	return predicate.PromotionPerson(sql.FieldLT(FieldAddress, v))
}

// AddressLTE applies the LTE predicate on the "address" field.
func AddressLTE(v string) predicate.PromotionPerson {
	return predicate.PromotionPerson(sql.FieldLTE(FieldAddress, v))
}

// AddressContains applies the Contains predicate on the "address" field.
func AddressContains(v string) predicate.PromotionPerson {
	return predicate.PromotionPerson(sql.FieldContains(FieldAddress, v))
}

// AddressHasPrefix applies the HasPrefix predicate on the "address" field.
func AddressHasPrefix(v string) predicate.PromotionPerson {
	return predicate.PromotionPerson(sql.FieldHasPrefix(FieldAddress, v))
}

// AddressHasSuffix applies the HasSuffix predicate on the "address" field.
func AddressHasSuffix(v string) predicate.PromotionPerson {
	return predicate.PromotionPerson(sql.FieldHasSuffix(FieldAddress, v))
}

// AddressIsNil applies the IsNil predicate on the "address" field.
func AddressIsNil() predicate.PromotionPerson {
	return predicate.PromotionPerson(sql.FieldIsNull(FieldAddress))
}

// AddressNotNil applies the NotNil predicate on the "address" field.
func AddressNotNil() predicate.PromotionPerson {
	return predicate.PromotionPerson(sql.FieldNotNull(FieldAddress))
}

// AddressEqualFold applies the EqualFold predicate on the "address" field.
func AddressEqualFold(v string) predicate.PromotionPerson {
	return predicate.PromotionPerson(sql.FieldEqualFold(FieldAddress, v))
}

// AddressContainsFold applies the ContainsFold predicate on the "address" field.
func AddressContainsFold(v string) predicate.PromotionPerson {
	return predicate.PromotionPerson(sql.FieldContainsFold(FieldAddress, v))
}

// HasMember applies the HasEdge predicate on the "member" edge.
func HasMember() predicate.PromotionPerson {
	return predicate.PromotionPerson(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, MemberTable, MemberColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasMemberWith applies the HasEdge predicate on the "member" edge with a given conditions (other predicates).
func HasMemberWith(preds ...predicate.PromotionMember) predicate.PromotionPerson {
	return predicate.PromotionPerson(func(s *sql.Selector) {
		step := newMemberStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.PromotionPerson) predicate.PromotionPerson {
	return predicate.PromotionPerson(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for _, p := range predicates {
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.PromotionPerson) predicate.PromotionPerson {
	return predicate.PromotionPerson(func(s *sql.Selector) {
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
func Not(p predicate.PromotionPerson) predicate.PromotionPerson {
	return predicate.PromotionPerson(func(s *sql.Selector) {
		p(s.Not())
	})
}