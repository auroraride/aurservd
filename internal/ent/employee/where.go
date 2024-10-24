// Code generated by ent, DO NOT EDIT.

package employee

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/auroraride/aurservd/internal/ent/predicate"
	"github.com/google/uuid"
)

// ID filters vertices based on their ID field.
func ID(id uint64) predicate.Employee {
	return predicate.Employee(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id uint64) predicate.Employee {
	return predicate.Employee(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id uint64) predicate.Employee {
	return predicate.Employee(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...uint64) predicate.Employee {
	return predicate.Employee(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...uint64) predicate.Employee {
	return predicate.Employee(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id uint64) predicate.Employee {
	return predicate.Employee(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id uint64) predicate.Employee {
	return predicate.Employee(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id uint64) predicate.Employee {
	return predicate.Employee(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id uint64) predicate.Employee {
	return predicate.Employee(sql.FieldLTE(FieldID, id))
}

// CreatedAt applies equality check predicate on the "created_at" field. It's identical to CreatedAtEQ.
func CreatedAt(v time.Time) predicate.Employee {
	return predicate.Employee(sql.FieldEQ(FieldCreatedAt, v))
}

// UpdatedAt applies equality check predicate on the "updated_at" field. It's identical to UpdatedAtEQ.
func UpdatedAt(v time.Time) predicate.Employee {
	return predicate.Employee(sql.FieldEQ(FieldUpdatedAt, v))
}

// DeletedAt applies equality check predicate on the "deleted_at" field. It's identical to DeletedAtEQ.
func DeletedAt(v time.Time) predicate.Employee {
	return predicate.Employee(sql.FieldEQ(FieldDeletedAt, v))
}

// Remark applies equality check predicate on the "remark" field. It's identical to RemarkEQ.
func Remark(v string) predicate.Employee {
	return predicate.Employee(sql.FieldEQ(FieldRemark, v))
}

// CityID applies equality check predicate on the "city_id" field. It's identical to CityIDEQ.
func CityID(v uint64) predicate.Employee {
	return predicate.Employee(sql.FieldEQ(FieldCityID, v))
}

// GroupID applies equality check predicate on the "group_id" field. It's identical to GroupIDEQ.
func GroupID(v uint64) predicate.Employee {
	return predicate.Employee(sql.FieldEQ(FieldGroupID, v))
}

// Sn applies equality check predicate on the "sn" field. It's identical to SnEQ.
func Sn(v uuid.UUID) predicate.Employee {
	return predicate.Employee(sql.FieldEQ(FieldSn, v))
}

// Name applies equality check predicate on the "name" field. It's identical to NameEQ.
func Name(v string) predicate.Employee {
	return predicate.Employee(sql.FieldEQ(FieldName, v))
}

// Phone applies equality check predicate on the "phone" field. It's identical to PhoneEQ.
func Phone(v string) predicate.Employee {
	return predicate.Employee(sql.FieldEQ(FieldPhone, v))
}

// Enable applies equality check predicate on the "enable" field. It's identical to EnableEQ.
func Enable(v bool) predicate.Employee {
	return predicate.Employee(sql.FieldEQ(FieldEnable, v))
}

// Password applies equality check predicate on the "password" field. It's identical to PasswordEQ.
func Password(v string) predicate.Employee {
	return predicate.Employee(sql.FieldEQ(FieldPassword, v))
}

// Limit applies equality check predicate on the "limit" field. It's identical to LimitEQ.
func Limit(v uint) predicate.Employee {
	return predicate.Employee(sql.FieldEQ(FieldLimit, v))
}

// DutyStoreID applies equality check predicate on the "duty_store_id" field. It's identical to DutyStoreIDEQ.
func DutyStoreID(v uint64) predicate.Employee {
	return predicate.Employee(sql.FieldEQ(FieldDutyStoreID, v))
}

// CreatedAtEQ applies the EQ predicate on the "created_at" field.
func CreatedAtEQ(v time.Time) predicate.Employee {
	return predicate.Employee(sql.FieldEQ(FieldCreatedAt, v))
}

// CreatedAtNEQ applies the NEQ predicate on the "created_at" field.
func CreatedAtNEQ(v time.Time) predicate.Employee {
	return predicate.Employee(sql.FieldNEQ(FieldCreatedAt, v))
}

// CreatedAtIn applies the In predicate on the "created_at" field.
func CreatedAtIn(vs ...time.Time) predicate.Employee {
	return predicate.Employee(sql.FieldIn(FieldCreatedAt, vs...))
}

// CreatedAtNotIn applies the NotIn predicate on the "created_at" field.
func CreatedAtNotIn(vs ...time.Time) predicate.Employee {
	return predicate.Employee(sql.FieldNotIn(FieldCreatedAt, vs...))
}

// CreatedAtGT applies the GT predicate on the "created_at" field.
func CreatedAtGT(v time.Time) predicate.Employee {
	return predicate.Employee(sql.FieldGT(FieldCreatedAt, v))
}

// CreatedAtGTE applies the GTE predicate on the "created_at" field.
func CreatedAtGTE(v time.Time) predicate.Employee {
	return predicate.Employee(sql.FieldGTE(FieldCreatedAt, v))
}

// CreatedAtLT applies the LT predicate on the "created_at" field.
func CreatedAtLT(v time.Time) predicate.Employee {
	return predicate.Employee(sql.FieldLT(FieldCreatedAt, v))
}

// CreatedAtLTE applies the LTE predicate on the "created_at" field.
func CreatedAtLTE(v time.Time) predicate.Employee {
	return predicate.Employee(sql.FieldLTE(FieldCreatedAt, v))
}

// UpdatedAtEQ applies the EQ predicate on the "updated_at" field.
func UpdatedAtEQ(v time.Time) predicate.Employee {
	return predicate.Employee(sql.FieldEQ(FieldUpdatedAt, v))
}

// UpdatedAtNEQ applies the NEQ predicate on the "updated_at" field.
func UpdatedAtNEQ(v time.Time) predicate.Employee {
	return predicate.Employee(sql.FieldNEQ(FieldUpdatedAt, v))
}

// UpdatedAtIn applies the In predicate on the "updated_at" field.
func UpdatedAtIn(vs ...time.Time) predicate.Employee {
	return predicate.Employee(sql.FieldIn(FieldUpdatedAt, vs...))
}

// UpdatedAtNotIn applies the NotIn predicate on the "updated_at" field.
func UpdatedAtNotIn(vs ...time.Time) predicate.Employee {
	return predicate.Employee(sql.FieldNotIn(FieldUpdatedAt, vs...))
}

// UpdatedAtGT applies the GT predicate on the "updated_at" field.
func UpdatedAtGT(v time.Time) predicate.Employee {
	return predicate.Employee(sql.FieldGT(FieldUpdatedAt, v))
}

// UpdatedAtGTE applies the GTE predicate on the "updated_at" field.
func UpdatedAtGTE(v time.Time) predicate.Employee {
	return predicate.Employee(sql.FieldGTE(FieldUpdatedAt, v))
}

// UpdatedAtLT applies the LT predicate on the "updated_at" field.
func UpdatedAtLT(v time.Time) predicate.Employee {
	return predicate.Employee(sql.FieldLT(FieldUpdatedAt, v))
}

// UpdatedAtLTE applies the LTE predicate on the "updated_at" field.
func UpdatedAtLTE(v time.Time) predicate.Employee {
	return predicate.Employee(sql.FieldLTE(FieldUpdatedAt, v))
}

// DeletedAtEQ applies the EQ predicate on the "deleted_at" field.
func DeletedAtEQ(v time.Time) predicate.Employee {
	return predicate.Employee(sql.FieldEQ(FieldDeletedAt, v))
}

// DeletedAtNEQ applies the NEQ predicate on the "deleted_at" field.
func DeletedAtNEQ(v time.Time) predicate.Employee {
	return predicate.Employee(sql.FieldNEQ(FieldDeletedAt, v))
}

// DeletedAtIn applies the In predicate on the "deleted_at" field.
func DeletedAtIn(vs ...time.Time) predicate.Employee {
	return predicate.Employee(sql.FieldIn(FieldDeletedAt, vs...))
}

// DeletedAtNotIn applies the NotIn predicate on the "deleted_at" field.
func DeletedAtNotIn(vs ...time.Time) predicate.Employee {
	return predicate.Employee(sql.FieldNotIn(FieldDeletedAt, vs...))
}

// DeletedAtGT applies the GT predicate on the "deleted_at" field.
func DeletedAtGT(v time.Time) predicate.Employee {
	return predicate.Employee(sql.FieldGT(FieldDeletedAt, v))
}

// DeletedAtGTE applies the GTE predicate on the "deleted_at" field.
func DeletedAtGTE(v time.Time) predicate.Employee {
	return predicate.Employee(sql.FieldGTE(FieldDeletedAt, v))
}

// DeletedAtLT applies the LT predicate on the "deleted_at" field.
func DeletedAtLT(v time.Time) predicate.Employee {
	return predicate.Employee(sql.FieldLT(FieldDeletedAt, v))
}

// DeletedAtLTE applies the LTE predicate on the "deleted_at" field.
func DeletedAtLTE(v time.Time) predicate.Employee {
	return predicate.Employee(sql.FieldLTE(FieldDeletedAt, v))
}

// DeletedAtIsNil applies the IsNil predicate on the "deleted_at" field.
func DeletedAtIsNil() predicate.Employee {
	return predicate.Employee(sql.FieldIsNull(FieldDeletedAt))
}

// DeletedAtNotNil applies the NotNil predicate on the "deleted_at" field.
func DeletedAtNotNil() predicate.Employee {
	return predicate.Employee(sql.FieldNotNull(FieldDeletedAt))
}

// CreatorIsNil applies the IsNil predicate on the "creator" field.
func CreatorIsNil() predicate.Employee {
	return predicate.Employee(sql.FieldIsNull(FieldCreator))
}

// CreatorNotNil applies the NotNil predicate on the "creator" field.
func CreatorNotNil() predicate.Employee {
	return predicate.Employee(sql.FieldNotNull(FieldCreator))
}

// LastModifierIsNil applies the IsNil predicate on the "last_modifier" field.
func LastModifierIsNil() predicate.Employee {
	return predicate.Employee(sql.FieldIsNull(FieldLastModifier))
}

// LastModifierNotNil applies the NotNil predicate on the "last_modifier" field.
func LastModifierNotNil() predicate.Employee {
	return predicate.Employee(sql.FieldNotNull(FieldLastModifier))
}

// RemarkEQ applies the EQ predicate on the "remark" field.
func RemarkEQ(v string) predicate.Employee {
	return predicate.Employee(sql.FieldEQ(FieldRemark, v))
}

// RemarkNEQ applies the NEQ predicate on the "remark" field.
func RemarkNEQ(v string) predicate.Employee {
	return predicate.Employee(sql.FieldNEQ(FieldRemark, v))
}

// RemarkIn applies the In predicate on the "remark" field.
func RemarkIn(vs ...string) predicate.Employee {
	return predicate.Employee(sql.FieldIn(FieldRemark, vs...))
}

// RemarkNotIn applies the NotIn predicate on the "remark" field.
func RemarkNotIn(vs ...string) predicate.Employee {
	return predicate.Employee(sql.FieldNotIn(FieldRemark, vs...))
}

// RemarkGT applies the GT predicate on the "remark" field.
func RemarkGT(v string) predicate.Employee {
	return predicate.Employee(sql.FieldGT(FieldRemark, v))
}

// RemarkGTE applies the GTE predicate on the "remark" field.
func RemarkGTE(v string) predicate.Employee {
	return predicate.Employee(sql.FieldGTE(FieldRemark, v))
}

// RemarkLT applies the LT predicate on the "remark" field.
func RemarkLT(v string) predicate.Employee {
	return predicate.Employee(sql.FieldLT(FieldRemark, v))
}

// RemarkLTE applies the LTE predicate on the "remark" field.
func RemarkLTE(v string) predicate.Employee {
	return predicate.Employee(sql.FieldLTE(FieldRemark, v))
}

// RemarkContains applies the Contains predicate on the "remark" field.
func RemarkContains(v string) predicate.Employee {
	return predicate.Employee(sql.FieldContains(FieldRemark, v))
}

// RemarkHasPrefix applies the HasPrefix predicate on the "remark" field.
func RemarkHasPrefix(v string) predicate.Employee {
	return predicate.Employee(sql.FieldHasPrefix(FieldRemark, v))
}

// RemarkHasSuffix applies the HasSuffix predicate on the "remark" field.
func RemarkHasSuffix(v string) predicate.Employee {
	return predicate.Employee(sql.FieldHasSuffix(FieldRemark, v))
}

// RemarkIsNil applies the IsNil predicate on the "remark" field.
func RemarkIsNil() predicate.Employee {
	return predicate.Employee(sql.FieldIsNull(FieldRemark))
}

// RemarkNotNil applies the NotNil predicate on the "remark" field.
func RemarkNotNil() predicate.Employee {
	return predicate.Employee(sql.FieldNotNull(FieldRemark))
}

// RemarkEqualFold applies the EqualFold predicate on the "remark" field.
func RemarkEqualFold(v string) predicate.Employee {
	return predicate.Employee(sql.FieldEqualFold(FieldRemark, v))
}

// RemarkContainsFold applies the ContainsFold predicate on the "remark" field.
func RemarkContainsFold(v string) predicate.Employee {
	return predicate.Employee(sql.FieldContainsFold(FieldRemark, v))
}

// CityIDEQ applies the EQ predicate on the "city_id" field.
func CityIDEQ(v uint64) predicate.Employee {
	return predicate.Employee(sql.FieldEQ(FieldCityID, v))
}

// CityIDNEQ applies the NEQ predicate on the "city_id" field.
func CityIDNEQ(v uint64) predicate.Employee {
	return predicate.Employee(sql.FieldNEQ(FieldCityID, v))
}

// CityIDIn applies the In predicate on the "city_id" field.
func CityIDIn(vs ...uint64) predicate.Employee {
	return predicate.Employee(sql.FieldIn(FieldCityID, vs...))
}

// CityIDNotIn applies the NotIn predicate on the "city_id" field.
func CityIDNotIn(vs ...uint64) predicate.Employee {
	return predicate.Employee(sql.FieldNotIn(FieldCityID, vs...))
}

// GroupIDEQ applies the EQ predicate on the "group_id" field.
func GroupIDEQ(v uint64) predicate.Employee {
	return predicate.Employee(sql.FieldEQ(FieldGroupID, v))
}

// GroupIDNEQ applies the NEQ predicate on the "group_id" field.
func GroupIDNEQ(v uint64) predicate.Employee {
	return predicate.Employee(sql.FieldNEQ(FieldGroupID, v))
}

// GroupIDIn applies the In predicate on the "group_id" field.
func GroupIDIn(vs ...uint64) predicate.Employee {
	return predicate.Employee(sql.FieldIn(FieldGroupID, vs...))
}

// GroupIDNotIn applies the NotIn predicate on the "group_id" field.
func GroupIDNotIn(vs ...uint64) predicate.Employee {
	return predicate.Employee(sql.FieldNotIn(FieldGroupID, vs...))
}

// GroupIDIsNil applies the IsNil predicate on the "group_id" field.
func GroupIDIsNil() predicate.Employee {
	return predicate.Employee(sql.FieldIsNull(FieldGroupID))
}

// GroupIDNotNil applies the NotNil predicate on the "group_id" field.
func GroupIDNotNil() predicate.Employee {
	return predicate.Employee(sql.FieldNotNull(FieldGroupID))
}

// SnEQ applies the EQ predicate on the "sn" field.
func SnEQ(v uuid.UUID) predicate.Employee {
	return predicate.Employee(sql.FieldEQ(FieldSn, v))
}

// SnNEQ applies the NEQ predicate on the "sn" field.
func SnNEQ(v uuid.UUID) predicate.Employee {
	return predicate.Employee(sql.FieldNEQ(FieldSn, v))
}

// SnIn applies the In predicate on the "sn" field.
func SnIn(vs ...uuid.UUID) predicate.Employee {
	return predicate.Employee(sql.FieldIn(FieldSn, vs...))
}

// SnNotIn applies the NotIn predicate on the "sn" field.
func SnNotIn(vs ...uuid.UUID) predicate.Employee {
	return predicate.Employee(sql.FieldNotIn(FieldSn, vs...))
}

// SnGT applies the GT predicate on the "sn" field.
func SnGT(v uuid.UUID) predicate.Employee {
	return predicate.Employee(sql.FieldGT(FieldSn, v))
}

// SnGTE applies the GTE predicate on the "sn" field.
func SnGTE(v uuid.UUID) predicate.Employee {
	return predicate.Employee(sql.FieldGTE(FieldSn, v))
}

// SnLT applies the LT predicate on the "sn" field.
func SnLT(v uuid.UUID) predicate.Employee {
	return predicate.Employee(sql.FieldLT(FieldSn, v))
}

// SnLTE applies the LTE predicate on the "sn" field.
func SnLTE(v uuid.UUID) predicate.Employee {
	return predicate.Employee(sql.FieldLTE(FieldSn, v))
}

// SnIsNil applies the IsNil predicate on the "sn" field.
func SnIsNil() predicate.Employee {
	return predicate.Employee(sql.FieldIsNull(FieldSn))
}

// SnNotNil applies the NotNil predicate on the "sn" field.
func SnNotNil() predicate.Employee {
	return predicate.Employee(sql.FieldNotNull(FieldSn))
}

// NameEQ applies the EQ predicate on the "name" field.
func NameEQ(v string) predicate.Employee {
	return predicate.Employee(sql.FieldEQ(FieldName, v))
}

// NameNEQ applies the NEQ predicate on the "name" field.
func NameNEQ(v string) predicate.Employee {
	return predicate.Employee(sql.FieldNEQ(FieldName, v))
}

// NameIn applies the In predicate on the "name" field.
func NameIn(vs ...string) predicate.Employee {
	return predicate.Employee(sql.FieldIn(FieldName, vs...))
}

// NameNotIn applies the NotIn predicate on the "name" field.
func NameNotIn(vs ...string) predicate.Employee {
	return predicate.Employee(sql.FieldNotIn(FieldName, vs...))
}

// NameGT applies the GT predicate on the "name" field.
func NameGT(v string) predicate.Employee {
	return predicate.Employee(sql.FieldGT(FieldName, v))
}

// NameGTE applies the GTE predicate on the "name" field.
func NameGTE(v string) predicate.Employee {
	return predicate.Employee(sql.FieldGTE(FieldName, v))
}

// NameLT applies the LT predicate on the "name" field.
func NameLT(v string) predicate.Employee {
	return predicate.Employee(sql.FieldLT(FieldName, v))
}

// NameLTE applies the LTE predicate on the "name" field.
func NameLTE(v string) predicate.Employee {
	return predicate.Employee(sql.FieldLTE(FieldName, v))
}

// NameContains applies the Contains predicate on the "name" field.
func NameContains(v string) predicate.Employee {
	return predicate.Employee(sql.FieldContains(FieldName, v))
}

// NameHasPrefix applies the HasPrefix predicate on the "name" field.
func NameHasPrefix(v string) predicate.Employee {
	return predicate.Employee(sql.FieldHasPrefix(FieldName, v))
}

// NameHasSuffix applies the HasSuffix predicate on the "name" field.
func NameHasSuffix(v string) predicate.Employee {
	return predicate.Employee(sql.FieldHasSuffix(FieldName, v))
}

// NameEqualFold applies the EqualFold predicate on the "name" field.
func NameEqualFold(v string) predicate.Employee {
	return predicate.Employee(sql.FieldEqualFold(FieldName, v))
}

// NameContainsFold applies the ContainsFold predicate on the "name" field.
func NameContainsFold(v string) predicate.Employee {
	return predicate.Employee(sql.FieldContainsFold(FieldName, v))
}

// PhoneEQ applies the EQ predicate on the "phone" field.
func PhoneEQ(v string) predicate.Employee {
	return predicate.Employee(sql.FieldEQ(FieldPhone, v))
}

// PhoneNEQ applies the NEQ predicate on the "phone" field.
func PhoneNEQ(v string) predicate.Employee {
	return predicate.Employee(sql.FieldNEQ(FieldPhone, v))
}

// PhoneIn applies the In predicate on the "phone" field.
func PhoneIn(vs ...string) predicate.Employee {
	return predicate.Employee(sql.FieldIn(FieldPhone, vs...))
}

// PhoneNotIn applies the NotIn predicate on the "phone" field.
func PhoneNotIn(vs ...string) predicate.Employee {
	return predicate.Employee(sql.FieldNotIn(FieldPhone, vs...))
}

// PhoneGT applies the GT predicate on the "phone" field.
func PhoneGT(v string) predicate.Employee {
	return predicate.Employee(sql.FieldGT(FieldPhone, v))
}

// PhoneGTE applies the GTE predicate on the "phone" field.
func PhoneGTE(v string) predicate.Employee {
	return predicate.Employee(sql.FieldGTE(FieldPhone, v))
}

// PhoneLT applies the LT predicate on the "phone" field.
func PhoneLT(v string) predicate.Employee {
	return predicate.Employee(sql.FieldLT(FieldPhone, v))
}

// PhoneLTE applies the LTE predicate on the "phone" field.
func PhoneLTE(v string) predicate.Employee {
	return predicate.Employee(sql.FieldLTE(FieldPhone, v))
}

// PhoneContains applies the Contains predicate on the "phone" field.
func PhoneContains(v string) predicate.Employee {
	return predicate.Employee(sql.FieldContains(FieldPhone, v))
}

// PhoneHasPrefix applies the HasPrefix predicate on the "phone" field.
func PhoneHasPrefix(v string) predicate.Employee {
	return predicate.Employee(sql.FieldHasPrefix(FieldPhone, v))
}

// PhoneHasSuffix applies the HasSuffix predicate on the "phone" field.
func PhoneHasSuffix(v string) predicate.Employee {
	return predicate.Employee(sql.FieldHasSuffix(FieldPhone, v))
}

// PhoneEqualFold applies the EqualFold predicate on the "phone" field.
func PhoneEqualFold(v string) predicate.Employee {
	return predicate.Employee(sql.FieldEqualFold(FieldPhone, v))
}

// PhoneContainsFold applies the ContainsFold predicate on the "phone" field.
func PhoneContainsFold(v string) predicate.Employee {
	return predicate.Employee(sql.FieldContainsFold(FieldPhone, v))
}

// EnableEQ applies the EQ predicate on the "enable" field.
func EnableEQ(v bool) predicate.Employee {
	return predicate.Employee(sql.FieldEQ(FieldEnable, v))
}

// EnableNEQ applies the NEQ predicate on the "enable" field.
func EnableNEQ(v bool) predicate.Employee {
	return predicate.Employee(sql.FieldNEQ(FieldEnable, v))
}

// PasswordEQ applies the EQ predicate on the "password" field.
func PasswordEQ(v string) predicate.Employee {
	return predicate.Employee(sql.FieldEQ(FieldPassword, v))
}

// PasswordNEQ applies the NEQ predicate on the "password" field.
func PasswordNEQ(v string) predicate.Employee {
	return predicate.Employee(sql.FieldNEQ(FieldPassword, v))
}

// PasswordIn applies the In predicate on the "password" field.
func PasswordIn(vs ...string) predicate.Employee {
	return predicate.Employee(sql.FieldIn(FieldPassword, vs...))
}

// PasswordNotIn applies the NotIn predicate on the "password" field.
func PasswordNotIn(vs ...string) predicate.Employee {
	return predicate.Employee(sql.FieldNotIn(FieldPassword, vs...))
}

// PasswordGT applies the GT predicate on the "password" field.
func PasswordGT(v string) predicate.Employee {
	return predicate.Employee(sql.FieldGT(FieldPassword, v))
}

// PasswordGTE applies the GTE predicate on the "password" field.
func PasswordGTE(v string) predicate.Employee {
	return predicate.Employee(sql.FieldGTE(FieldPassword, v))
}

// PasswordLT applies the LT predicate on the "password" field.
func PasswordLT(v string) predicate.Employee {
	return predicate.Employee(sql.FieldLT(FieldPassword, v))
}

// PasswordLTE applies the LTE predicate on the "password" field.
func PasswordLTE(v string) predicate.Employee {
	return predicate.Employee(sql.FieldLTE(FieldPassword, v))
}

// PasswordContains applies the Contains predicate on the "password" field.
func PasswordContains(v string) predicate.Employee {
	return predicate.Employee(sql.FieldContains(FieldPassword, v))
}

// PasswordHasPrefix applies the HasPrefix predicate on the "password" field.
func PasswordHasPrefix(v string) predicate.Employee {
	return predicate.Employee(sql.FieldHasPrefix(FieldPassword, v))
}

// PasswordHasSuffix applies the HasSuffix predicate on the "password" field.
func PasswordHasSuffix(v string) predicate.Employee {
	return predicate.Employee(sql.FieldHasSuffix(FieldPassword, v))
}

// PasswordIsNil applies the IsNil predicate on the "password" field.
func PasswordIsNil() predicate.Employee {
	return predicate.Employee(sql.FieldIsNull(FieldPassword))
}

// PasswordNotNil applies the NotNil predicate on the "password" field.
func PasswordNotNil() predicate.Employee {
	return predicate.Employee(sql.FieldNotNull(FieldPassword))
}

// PasswordEqualFold applies the EqualFold predicate on the "password" field.
func PasswordEqualFold(v string) predicate.Employee {
	return predicate.Employee(sql.FieldEqualFold(FieldPassword, v))
}

// PasswordContainsFold applies the ContainsFold predicate on the "password" field.
func PasswordContainsFold(v string) predicate.Employee {
	return predicate.Employee(sql.FieldContainsFold(FieldPassword, v))
}

// LimitEQ applies the EQ predicate on the "limit" field.
func LimitEQ(v uint) predicate.Employee {
	return predicate.Employee(sql.FieldEQ(FieldLimit, v))
}

// LimitNEQ applies the NEQ predicate on the "limit" field.
func LimitNEQ(v uint) predicate.Employee {
	return predicate.Employee(sql.FieldNEQ(FieldLimit, v))
}

// LimitIn applies the In predicate on the "limit" field.
func LimitIn(vs ...uint) predicate.Employee {
	return predicate.Employee(sql.FieldIn(FieldLimit, vs...))
}

// LimitNotIn applies the NotIn predicate on the "limit" field.
func LimitNotIn(vs ...uint) predicate.Employee {
	return predicate.Employee(sql.FieldNotIn(FieldLimit, vs...))
}

// LimitGT applies the GT predicate on the "limit" field.
func LimitGT(v uint) predicate.Employee {
	return predicate.Employee(sql.FieldGT(FieldLimit, v))
}

// LimitGTE applies the GTE predicate on the "limit" field.
func LimitGTE(v uint) predicate.Employee {
	return predicate.Employee(sql.FieldGTE(FieldLimit, v))
}

// LimitLT applies the LT predicate on the "limit" field.
func LimitLT(v uint) predicate.Employee {
	return predicate.Employee(sql.FieldLT(FieldLimit, v))
}

// LimitLTE applies the LTE predicate on the "limit" field.
func LimitLTE(v uint) predicate.Employee {
	return predicate.Employee(sql.FieldLTE(FieldLimit, v))
}

// DutyStoreIDEQ applies the EQ predicate on the "duty_store_id" field.
func DutyStoreIDEQ(v uint64) predicate.Employee {
	return predicate.Employee(sql.FieldEQ(FieldDutyStoreID, v))
}

// DutyStoreIDNEQ applies the NEQ predicate on the "duty_store_id" field.
func DutyStoreIDNEQ(v uint64) predicate.Employee {
	return predicate.Employee(sql.FieldNEQ(FieldDutyStoreID, v))
}

// DutyStoreIDIn applies the In predicate on the "duty_store_id" field.
func DutyStoreIDIn(vs ...uint64) predicate.Employee {
	return predicate.Employee(sql.FieldIn(FieldDutyStoreID, vs...))
}

// DutyStoreIDNotIn applies the NotIn predicate on the "duty_store_id" field.
func DutyStoreIDNotIn(vs ...uint64) predicate.Employee {
	return predicate.Employee(sql.FieldNotIn(FieldDutyStoreID, vs...))
}

// DutyStoreIDIsNil applies the IsNil predicate on the "duty_store_id" field.
func DutyStoreIDIsNil() predicate.Employee {
	return predicate.Employee(sql.FieldIsNull(FieldDutyStoreID))
}

// DutyStoreIDNotNil applies the NotNil predicate on the "duty_store_id" field.
func DutyStoreIDNotNil() predicate.Employee {
	return predicate.Employee(sql.FieldNotNull(FieldDutyStoreID))
}

// HasCity applies the HasEdge predicate on the "city" edge.
func HasCity() predicate.Employee {
	return predicate.Employee(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, CityTable, CityColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasCityWith applies the HasEdge predicate on the "city" edge with a given conditions (other predicates).
func HasCityWith(preds ...predicate.City) predicate.Employee {
	return predicate.Employee(func(s *sql.Selector) {
		step := newCityStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasGroup applies the HasEdge predicate on the "group" edge.
func HasGroup() predicate.Employee {
	return predicate.Employee(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, GroupTable, GroupColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasGroupWith applies the HasEdge predicate on the "group" edge with a given conditions (other predicates).
func HasGroupWith(preds ...predicate.StoreGroup) predicate.Employee {
	return predicate.Employee(func(s *sql.Selector) {
		step := newGroupStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasStore applies the HasEdge predicate on the "store" edge.
func HasStore() predicate.Employee {
	return predicate.Employee(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.O2O, false, StoreTable, StoreColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasStoreWith applies the HasEdge predicate on the "store" edge with a given conditions (other predicates).
func HasStoreWith(preds ...predicate.Store) predicate.Employee {
	return predicate.Employee(func(s *sql.Selector) {
		step := newStoreStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasAttendances applies the HasEdge predicate on the "attendances" edge.
func HasAttendances() predicate.Employee {
	return predicate.Employee(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, AttendancesTable, AttendancesColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasAttendancesWith applies the HasEdge predicate on the "attendances" edge with a given conditions (other predicates).
func HasAttendancesWith(preds ...predicate.Attendance) predicate.Employee {
	return predicate.Employee(func(s *sql.Selector) {
		step := newAttendancesStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasStocks applies the HasEdge predicate on the "stocks" edge.
func HasStocks() predicate.Employee {
	return predicate.Employee(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, StocksTable, StocksColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasStocksWith applies the HasEdge predicate on the "stocks" edge with a given conditions (other predicates).
func HasStocksWith(preds ...predicate.Stock) predicate.Employee {
	return predicate.Employee(func(s *sql.Selector) {
		step := newStocksStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasExchanges applies the HasEdge predicate on the "exchanges" edge.
func HasExchanges() predicate.Employee {
	return predicate.Employee(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, ExchangesTable, ExchangesColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasExchangesWith applies the HasEdge predicate on the "exchanges" edge with a given conditions (other predicates).
func HasExchangesWith(preds ...predicate.Exchange) predicate.Employee {
	return predicate.Employee(func(s *sql.Selector) {
		step := newExchangesStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasCommissions applies the HasEdge predicate on the "commissions" edge.
func HasCommissions() predicate.Employee {
	return predicate.Employee(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, CommissionsTable, CommissionsColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasCommissionsWith applies the HasEdge predicate on the "commissions" edge with a given conditions (other predicates).
func HasCommissionsWith(preds ...predicate.Commission) predicate.Employee {
	return predicate.Employee(func(s *sql.Selector) {
		step := newCommissionsStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasAssistances applies the HasEdge predicate on the "assistances" edge.
func HasAssistances() predicate.Employee {
	return predicate.Employee(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, AssistancesTable, AssistancesColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasAssistancesWith applies the HasEdge predicate on the "assistances" edge with a given conditions (other predicates).
func HasAssistancesWith(preds ...predicate.Assistance) predicate.Employee {
	return predicate.Employee(func(s *sql.Selector) {
		step := newAssistancesStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasStores applies the HasEdge predicate on the "stores" edge.
func HasStores() predicate.Employee {
	return predicate.Employee(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2M, true, StoresTable, StoresPrimaryKey...),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasStoresWith applies the HasEdge predicate on the "stores" edge with a given conditions (other predicates).
func HasStoresWith(preds ...predicate.Store) predicate.Employee {
	return predicate.Employee(func(s *sql.Selector) {
		step := newStoresStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasDutyStore applies the HasEdge predicate on the "duty_store" edge.
func HasDutyStore() predicate.Employee {
	return predicate.Employee(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, DutyStoreTable, DutyStoreColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasDutyStoreWith applies the HasEdge predicate on the "duty_store" edge with a given conditions (other predicates).
func HasDutyStoreWith(preds ...predicate.Store) predicate.Employee {
	return predicate.Employee(func(s *sql.Selector) {
		step := newDutyStoreStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.Employee) predicate.Employee {
	return predicate.Employee(sql.AndPredicates(predicates...))
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.Employee) predicate.Employee {
	return predicate.Employee(sql.OrPredicates(predicates...))
}

// Not applies the not operator on the given predicate.
func Not(p predicate.Employee) predicate.Employee {
	return predicate.Employee(sql.NotPredicates(p))
}
