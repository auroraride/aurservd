// Code generated by ent, DO NOT EDIT.

package assetmanager

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/auroraride/aurservd/internal/ent/predicate"
)

// ID filters vertices based on their ID field.
func ID(id uint64) predicate.AssetManager {
	return predicate.AssetManager(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id uint64) predicate.AssetManager {
	return predicate.AssetManager(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id uint64) predicate.AssetManager {
	return predicate.AssetManager(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...uint64) predicate.AssetManager {
	return predicate.AssetManager(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...uint64) predicate.AssetManager {
	return predicate.AssetManager(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id uint64) predicate.AssetManager {
	return predicate.AssetManager(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id uint64) predicate.AssetManager {
	return predicate.AssetManager(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id uint64) predicate.AssetManager {
	return predicate.AssetManager(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id uint64) predicate.AssetManager {
	return predicate.AssetManager(sql.FieldLTE(FieldID, id))
}

// CreatedAt applies equality check predicate on the "created_at" field. It's identical to CreatedAtEQ.
func CreatedAt(v time.Time) predicate.AssetManager {
	return predicate.AssetManager(sql.FieldEQ(FieldCreatedAt, v))
}

// UpdatedAt applies equality check predicate on the "updated_at" field. It's identical to UpdatedAtEQ.
func UpdatedAt(v time.Time) predicate.AssetManager {
	return predicate.AssetManager(sql.FieldEQ(FieldUpdatedAt, v))
}

// DeletedAt applies equality check predicate on the "deleted_at" field. It's identical to DeletedAtEQ.
func DeletedAt(v time.Time) predicate.AssetManager {
	return predicate.AssetManager(sql.FieldEQ(FieldDeletedAt, v))
}

// Remark applies equality check predicate on the "remark" field. It's identical to RemarkEQ.
func Remark(v string) predicate.AssetManager {
	return predicate.AssetManager(sql.FieldEQ(FieldRemark, v))
}

// Name applies equality check predicate on the "name" field. It's identical to NameEQ.
func Name(v string) predicate.AssetManager {
	return predicate.AssetManager(sql.FieldEQ(FieldName, v))
}

// Phone applies equality check predicate on the "phone" field. It's identical to PhoneEQ.
func Phone(v string) predicate.AssetManager {
	return predicate.AssetManager(sql.FieldEQ(FieldPhone, v))
}

// Password applies equality check predicate on the "password" field. It's identical to PasswordEQ.
func Password(v string) predicate.AssetManager {
	return predicate.AssetManager(sql.FieldEQ(FieldPassword, v))
}

// RoleID applies equality check predicate on the "role_id" field. It's identical to RoleIDEQ.
func RoleID(v uint64) predicate.AssetManager {
	return predicate.AssetManager(sql.FieldEQ(FieldRoleID, v))
}

// MiniEnable applies equality check predicate on the "mini_enable" field. It's identical to MiniEnableEQ.
func MiniEnable(v bool) predicate.AssetManager {
	return predicate.AssetManager(sql.FieldEQ(FieldMiniEnable, v))
}

// MiniLimit applies equality check predicate on the "mini_limit" field. It's identical to MiniLimitEQ.
func MiniLimit(v uint) predicate.AssetManager {
	return predicate.AssetManager(sql.FieldEQ(FieldMiniLimit, v))
}

// LastSigninAt applies equality check predicate on the "last_signin_at" field. It's identical to LastSigninAtEQ.
func LastSigninAt(v time.Time) predicate.AssetManager {
	return predicate.AssetManager(sql.FieldEQ(FieldLastSigninAt, v))
}

// WarehouseID applies equality check predicate on the "warehouse_id" field. It's identical to WarehouseIDEQ.
func WarehouseID(v uint64) predicate.AssetManager {
	return predicate.AssetManager(sql.FieldEQ(FieldWarehouseID, v))
}

// CreatedAtEQ applies the EQ predicate on the "created_at" field.
func CreatedAtEQ(v time.Time) predicate.AssetManager {
	return predicate.AssetManager(sql.FieldEQ(FieldCreatedAt, v))
}

// CreatedAtNEQ applies the NEQ predicate on the "created_at" field.
func CreatedAtNEQ(v time.Time) predicate.AssetManager {
	return predicate.AssetManager(sql.FieldNEQ(FieldCreatedAt, v))
}

// CreatedAtIn applies the In predicate on the "created_at" field.
func CreatedAtIn(vs ...time.Time) predicate.AssetManager {
	return predicate.AssetManager(sql.FieldIn(FieldCreatedAt, vs...))
}

// CreatedAtNotIn applies the NotIn predicate on the "created_at" field.
func CreatedAtNotIn(vs ...time.Time) predicate.AssetManager {
	return predicate.AssetManager(sql.FieldNotIn(FieldCreatedAt, vs...))
}

// CreatedAtGT applies the GT predicate on the "created_at" field.
func CreatedAtGT(v time.Time) predicate.AssetManager {
	return predicate.AssetManager(sql.FieldGT(FieldCreatedAt, v))
}

// CreatedAtGTE applies the GTE predicate on the "created_at" field.
func CreatedAtGTE(v time.Time) predicate.AssetManager {
	return predicate.AssetManager(sql.FieldGTE(FieldCreatedAt, v))
}

// CreatedAtLT applies the LT predicate on the "created_at" field.
func CreatedAtLT(v time.Time) predicate.AssetManager {
	return predicate.AssetManager(sql.FieldLT(FieldCreatedAt, v))
}

// CreatedAtLTE applies the LTE predicate on the "created_at" field.
func CreatedAtLTE(v time.Time) predicate.AssetManager {
	return predicate.AssetManager(sql.FieldLTE(FieldCreatedAt, v))
}

// UpdatedAtEQ applies the EQ predicate on the "updated_at" field.
func UpdatedAtEQ(v time.Time) predicate.AssetManager {
	return predicate.AssetManager(sql.FieldEQ(FieldUpdatedAt, v))
}

// UpdatedAtNEQ applies the NEQ predicate on the "updated_at" field.
func UpdatedAtNEQ(v time.Time) predicate.AssetManager {
	return predicate.AssetManager(sql.FieldNEQ(FieldUpdatedAt, v))
}

// UpdatedAtIn applies the In predicate on the "updated_at" field.
func UpdatedAtIn(vs ...time.Time) predicate.AssetManager {
	return predicate.AssetManager(sql.FieldIn(FieldUpdatedAt, vs...))
}

// UpdatedAtNotIn applies the NotIn predicate on the "updated_at" field.
func UpdatedAtNotIn(vs ...time.Time) predicate.AssetManager {
	return predicate.AssetManager(sql.FieldNotIn(FieldUpdatedAt, vs...))
}

// UpdatedAtGT applies the GT predicate on the "updated_at" field.
func UpdatedAtGT(v time.Time) predicate.AssetManager {
	return predicate.AssetManager(sql.FieldGT(FieldUpdatedAt, v))
}

// UpdatedAtGTE applies the GTE predicate on the "updated_at" field.
func UpdatedAtGTE(v time.Time) predicate.AssetManager {
	return predicate.AssetManager(sql.FieldGTE(FieldUpdatedAt, v))
}

// UpdatedAtLT applies the LT predicate on the "updated_at" field.
func UpdatedAtLT(v time.Time) predicate.AssetManager {
	return predicate.AssetManager(sql.FieldLT(FieldUpdatedAt, v))
}

// UpdatedAtLTE applies the LTE predicate on the "updated_at" field.
func UpdatedAtLTE(v time.Time) predicate.AssetManager {
	return predicate.AssetManager(sql.FieldLTE(FieldUpdatedAt, v))
}

// DeletedAtEQ applies the EQ predicate on the "deleted_at" field.
func DeletedAtEQ(v time.Time) predicate.AssetManager {
	return predicate.AssetManager(sql.FieldEQ(FieldDeletedAt, v))
}

// DeletedAtNEQ applies the NEQ predicate on the "deleted_at" field.
func DeletedAtNEQ(v time.Time) predicate.AssetManager {
	return predicate.AssetManager(sql.FieldNEQ(FieldDeletedAt, v))
}

// DeletedAtIn applies the In predicate on the "deleted_at" field.
func DeletedAtIn(vs ...time.Time) predicate.AssetManager {
	return predicate.AssetManager(sql.FieldIn(FieldDeletedAt, vs...))
}

// DeletedAtNotIn applies the NotIn predicate on the "deleted_at" field.
func DeletedAtNotIn(vs ...time.Time) predicate.AssetManager {
	return predicate.AssetManager(sql.FieldNotIn(FieldDeletedAt, vs...))
}

// DeletedAtGT applies the GT predicate on the "deleted_at" field.
func DeletedAtGT(v time.Time) predicate.AssetManager {
	return predicate.AssetManager(sql.FieldGT(FieldDeletedAt, v))
}

// DeletedAtGTE applies the GTE predicate on the "deleted_at" field.
func DeletedAtGTE(v time.Time) predicate.AssetManager {
	return predicate.AssetManager(sql.FieldGTE(FieldDeletedAt, v))
}

// DeletedAtLT applies the LT predicate on the "deleted_at" field.
func DeletedAtLT(v time.Time) predicate.AssetManager {
	return predicate.AssetManager(sql.FieldLT(FieldDeletedAt, v))
}

// DeletedAtLTE applies the LTE predicate on the "deleted_at" field.
func DeletedAtLTE(v time.Time) predicate.AssetManager {
	return predicate.AssetManager(sql.FieldLTE(FieldDeletedAt, v))
}

// DeletedAtIsNil applies the IsNil predicate on the "deleted_at" field.
func DeletedAtIsNil() predicate.AssetManager {
	return predicate.AssetManager(sql.FieldIsNull(FieldDeletedAt))
}

// DeletedAtNotNil applies the NotNil predicate on the "deleted_at" field.
func DeletedAtNotNil() predicate.AssetManager {
	return predicate.AssetManager(sql.FieldNotNull(FieldDeletedAt))
}

// CreatorIsNil applies the IsNil predicate on the "creator" field.
func CreatorIsNil() predicate.AssetManager {
	return predicate.AssetManager(sql.FieldIsNull(FieldCreator))
}

// CreatorNotNil applies the NotNil predicate on the "creator" field.
func CreatorNotNil() predicate.AssetManager {
	return predicate.AssetManager(sql.FieldNotNull(FieldCreator))
}

// LastModifierIsNil applies the IsNil predicate on the "last_modifier" field.
func LastModifierIsNil() predicate.AssetManager {
	return predicate.AssetManager(sql.FieldIsNull(FieldLastModifier))
}

// LastModifierNotNil applies the NotNil predicate on the "last_modifier" field.
func LastModifierNotNil() predicate.AssetManager {
	return predicate.AssetManager(sql.FieldNotNull(FieldLastModifier))
}

// RemarkEQ applies the EQ predicate on the "remark" field.
func RemarkEQ(v string) predicate.AssetManager {
	return predicate.AssetManager(sql.FieldEQ(FieldRemark, v))
}

// RemarkNEQ applies the NEQ predicate on the "remark" field.
func RemarkNEQ(v string) predicate.AssetManager {
	return predicate.AssetManager(sql.FieldNEQ(FieldRemark, v))
}

// RemarkIn applies the In predicate on the "remark" field.
func RemarkIn(vs ...string) predicate.AssetManager {
	return predicate.AssetManager(sql.FieldIn(FieldRemark, vs...))
}

// RemarkNotIn applies the NotIn predicate on the "remark" field.
func RemarkNotIn(vs ...string) predicate.AssetManager {
	return predicate.AssetManager(sql.FieldNotIn(FieldRemark, vs...))
}

// RemarkGT applies the GT predicate on the "remark" field.
func RemarkGT(v string) predicate.AssetManager {
	return predicate.AssetManager(sql.FieldGT(FieldRemark, v))
}

// RemarkGTE applies the GTE predicate on the "remark" field.
func RemarkGTE(v string) predicate.AssetManager {
	return predicate.AssetManager(sql.FieldGTE(FieldRemark, v))
}

// RemarkLT applies the LT predicate on the "remark" field.
func RemarkLT(v string) predicate.AssetManager {
	return predicate.AssetManager(sql.FieldLT(FieldRemark, v))
}

// RemarkLTE applies the LTE predicate on the "remark" field.
func RemarkLTE(v string) predicate.AssetManager {
	return predicate.AssetManager(sql.FieldLTE(FieldRemark, v))
}

// RemarkContains applies the Contains predicate on the "remark" field.
func RemarkContains(v string) predicate.AssetManager {
	return predicate.AssetManager(sql.FieldContains(FieldRemark, v))
}

// RemarkHasPrefix applies the HasPrefix predicate on the "remark" field.
func RemarkHasPrefix(v string) predicate.AssetManager {
	return predicate.AssetManager(sql.FieldHasPrefix(FieldRemark, v))
}

// RemarkHasSuffix applies the HasSuffix predicate on the "remark" field.
func RemarkHasSuffix(v string) predicate.AssetManager {
	return predicate.AssetManager(sql.FieldHasSuffix(FieldRemark, v))
}

// RemarkIsNil applies the IsNil predicate on the "remark" field.
func RemarkIsNil() predicate.AssetManager {
	return predicate.AssetManager(sql.FieldIsNull(FieldRemark))
}

// RemarkNotNil applies the NotNil predicate on the "remark" field.
func RemarkNotNil() predicate.AssetManager {
	return predicate.AssetManager(sql.FieldNotNull(FieldRemark))
}

// RemarkEqualFold applies the EqualFold predicate on the "remark" field.
func RemarkEqualFold(v string) predicate.AssetManager {
	return predicate.AssetManager(sql.FieldEqualFold(FieldRemark, v))
}

// RemarkContainsFold applies the ContainsFold predicate on the "remark" field.
func RemarkContainsFold(v string) predicate.AssetManager {
	return predicate.AssetManager(sql.FieldContainsFold(FieldRemark, v))
}

// NameEQ applies the EQ predicate on the "name" field.
func NameEQ(v string) predicate.AssetManager {
	return predicate.AssetManager(sql.FieldEQ(FieldName, v))
}

// NameNEQ applies the NEQ predicate on the "name" field.
func NameNEQ(v string) predicate.AssetManager {
	return predicate.AssetManager(sql.FieldNEQ(FieldName, v))
}

// NameIn applies the In predicate on the "name" field.
func NameIn(vs ...string) predicate.AssetManager {
	return predicate.AssetManager(sql.FieldIn(FieldName, vs...))
}

// NameNotIn applies the NotIn predicate on the "name" field.
func NameNotIn(vs ...string) predicate.AssetManager {
	return predicate.AssetManager(sql.FieldNotIn(FieldName, vs...))
}

// NameGT applies the GT predicate on the "name" field.
func NameGT(v string) predicate.AssetManager {
	return predicate.AssetManager(sql.FieldGT(FieldName, v))
}

// NameGTE applies the GTE predicate on the "name" field.
func NameGTE(v string) predicate.AssetManager {
	return predicate.AssetManager(sql.FieldGTE(FieldName, v))
}

// NameLT applies the LT predicate on the "name" field.
func NameLT(v string) predicate.AssetManager {
	return predicate.AssetManager(sql.FieldLT(FieldName, v))
}

// NameLTE applies the LTE predicate on the "name" field.
func NameLTE(v string) predicate.AssetManager {
	return predicate.AssetManager(sql.FieldLTE(FieldName, v))
}

// NameContains applies the Contains predicate on the "name" field.
func NameContains(v string) predicate.AssetManager {
	return predicate.AssetManager(sql.FieldContains(FieldName, v))
}

// NameHasPrefix applies the HasPrefix predicate on the "name" field.
func NameHasPrefix(v string) predicate.AssetManager {
	return predicate.AssetManager(sql.FieldHasPrefix(FieldName, v))
}

// NameHasSuffix applies the HasSuffix predicate on the "name" field.
func NameHasSuffix(v string) predicate.AssetManager {
	return predicate.AssetManager(sql.FieldHasSuffix(FieldName, v))
}

// NameEqualFold applies the EqualFold predicate on the "name" field.
func NameEqualFold(v string) predicate.AssetManager {
	return predicate.AssetManager(sql.FieldEqualFold(FieldName, v))
}

// NameContainsFold applies the ContainsFold predicate on the "name" field.
func NameContainsFold(v string) predicate.AssetManager {
	return predicate.AssetManager(sql.FieldContainsFold(FieldName, v))
}

// PhoneEQ applies the EQ predicate on the "phone" field.
func PhoneEQ(v string) predicate.AssetManager {
	return predicate.AssetManager(sql.FieldEQ(FieldPhone, v))
}

// PhoneNEQ applies the NEQ predicate on the "phone" field.
func PhoneNEQ(v string) predicate.AssetManager {
	return predicate.AssetManager(sql.FieldNEQ(FieldPhone, v))
}

// PhoneIn applies the In predicate on the "phone" field.
func PhoneIn(vs ...string) predicate.AssetManager {
	return predicate.AssetManager(sql.FieldIn(FieldPhone, vs...))
}

// PhoneNotIn applies the NotIn predicate on the "phone" field.
func PhoneNotIn(vs ...string) predicate.AssetManager {
	return predicate.AssetManager(sql.FieldNotIn(FieldPhone, vs...))
}

// PhoneGT applies the GT predicate on the "phone" field.
func PhoneGT(v string) predicate.AssetManager {
	return predicate.AssetManager(sql.FieldGT(FieldPhone, v))
}

// PhoneGTE applies the GTE predicate on the "phone" field.
func PhoneGTE(v string) predicate.AssetManager {
	return predicate.AssetManager(sql.FieldGTE(FieldPhone, v))
}

// PhoneLT applies the LT predicate on the "phone" field.
func PhoneLT(v string) predicate.AssetManager {
	return predicate.AssetManager(sql.FieldLT(FieldPhone, v))
}

// PhoneLTE applies the LTE predicate on the "phone" field.
func PhoneLTE(v string) predicate.AssetManager {
	return predicate.AssetManager(sql.FieldLTE(FieldPhone, v))
}

// PhoneContains applies the Contains predicate on the "phone" field.
func PhoneContains(v string) predicate.AssetManager {
	return predicate.AssetManager(sql.FieldContains(FieldPhone, v))
}

// PhoneHasPrefix applies the HasPrefix predicate on the "phone" field.
func PhoneHasPrefix(v string) predicate.AssetManager {
	return predicate.AssetManager(sql.FieldHasPrefix(FieldPhone, v))
}

// PhoneHasSuffix applies the HasSuffix predicate on the "phone" field.
func PhoneHasSuffix(v string) predicate.AssetManager {
	return predicate.AssetManager(sql.FieldHasSuffix(FieldPhone, v))
}

// PhoneEqualFold applies the EqualFold predicate on the "phone" field.
func PhoneEqualFold(v string) predicate.AssetManager {
	return predicate.AssetManager(sql.FieldEqualFold(FieldPhone, v))
}

// PhoneContainsFold applies the ContainsFold predicate on the "phone" field.
func PhoneContainsFold(v string) predicate.AssetManager {
	return predicate.AssetManager(sql.FieldContainsFold(FieldPhone, v))
}

// PasswordEQ applies the EQ predicate on the "password" field.
func PasswordEQ(v string) predicate.AssetManager {
	return predicate.AssetManager(sql.FieldEQ(FieldPassword, v))
}

// PasswordNEQ applies the NEQ predicate on the "password" field.
func PasswordNEQ(v string) predicate.AssetManager {
	return predicate.AssetManager(sql.FieldNEQ(FieldPassword, v))
}

// PasswordIn applies the In predicate on the "password" field.
func PasswordIn(vs ...string) predicate.AssetManager {
	return predicate.AssetManager(sql.FieldIn(FieldPassword, vs...))
}

// PasswordNotIn applies the NotIn predicate on the "password" field.
func PasswordNotIn(vs ...string) predicate.AssetManager {
	return predicate.AssetManager(sql.FieldNotIn(FieldPassword, vs...))
}

// PasswordGT applies the GT predicate on the "password" field.
func PasswordGT(v string) predicate.AssetManager {
	return predicate.AssetManager(sql.FieldGT(FieldPassword, v))
}

// PasswordGTE applies the GTE predicate on the "password" field.
func PasswordGTE(v string) predicate.AssetManager {
	return predicate.AssetManager(sql.FieldGTE(FieldPassword, v))
}

// PasswordLT applies the LT predicate on the "password" field.
func PasswordLT(v string) predicate.AssetManager {
	return predicate.AssetManager(sql.FieldLT(FieldPassword, v))
}

// PasswordLTE applies the LTE predicate on the "password" field.
func PasswordLTE(v string) predicate.AssetManager {
	return predicate.AssetManager(sql.FieldLTE(FieldPassword, v))
}

// PasswordContains applies the Contains predicate on the "password" field.
func PasswordContains(v string) predicate.AssetManager {
	return predicate.AssetManager(sql.FieldContains(FieldPassword, v))
}

// PasswordHasPrefix applies the HasPrefix predicate on the "password" field.
func PasswordHasPrefix(v string) predicate.AssetManager {
	return predicate.AssetManager(sql.FieldHasPrefix(FieldPassword, v))
}

// PasswordHasSuffix applies the HasSuffix predicate on the "password" field.
func PasswordHasSuffix(v string) predicate.AssetManager {
	return predicate.AssetManager(sql.FieldHasSuffix(FieldPassword, v))
}

// PasswordEqualFold applies the EqualFold predicate on the "password" field.
func PasswordEqualFold(v string) predicate.AssetManager {
	return predicate.AssetManager(sql.FieldEqualFold(FieldPassword, v))
}

// PasswordContainsFold applies the ContainsFold predicate on the "password" field.
func PasswordContainsFold(v string) predicate.AssetManager {
	return predicate.AssetManager(sql.FieldContainsFold(FieldPassword, v))
}

// RoleIDEQ applies the EQ predicate on the "role_id" field.
func RoleIDEQ(v uint64) predicate.AssetManager {
	return predicate.AssetManager(sql.FieldEQ(FieldRoleID, v))
}

// RoleIDNEQ applies the NEQ predicate on the "role_id" field.
func RoleIDNEQ(v uint64) predicate.AssetManager {
	return predicate.AssetManager(sql.FieldNEQ(FieldRoleID, v))
}

// RoleIDIn applies the In predicate on the "role_id" field.
func RoleIDIn(vs ...uint64) predicate.AssetManager {
	return predicate.AssetManager(sql.FieldIn(FieldRoleID, vs...))
}

// RoleIDNotIn applies the NotIn predicate on the "role_id" field.
func RoleIDNotIn(vs ...uint64) predicate.AssetManager {
	return predicate.AssetManager(sql.FieldNotIn(FieldRoleID, vs...))
}

// RoleIDIsNil applies the IsNil predicate on the "role_id" field.
func RoleIDIsNil() predicate.AssetManager {
	return predicate.AssetManager(sql.FieldIsNull(FieldRoleID))
}

// RoleIDNotNil applies the NotNil predicate on the "role_id" field.
func RoleIDNotNil() predicate.AssetManager {
	return predicate.AssetManager(sql.FieldNotNull(FieldRoleID))
}

// MiniEnableEQ applies the EQ predicate on the "mini_enable" field.
func MiniEnableEQ(v bool) predicate.AssetManager {
	return predicate.AssetManager(sql.FieldEQ(FieldMiniEnable, v))
}

// MiniEnableNEQ applies the NEQ predicate on the "mini_enable" field.
func MiniEnableNEQ(v bool) predicate.AssetManager {
	return predicate.AssetManager(sql.FieldNEQ(FieldMiniEnable, v))
}

// MiniLimitEQ applies the EQ predicate on the "mini_limit" field.
func MiniLimitEQ(v uint) predicate.AssetManager {
	return predicate.AssetManager(sql.FieldEQ(FieldMiniLimit, v))
}

// MiniLimitNEQ applies the NEQ predicate on the "mini_limit" field.
func MiniLimitNEQ(v uint) predicate.AssetManager {
	return predicate.AssetManager(sql.FieldNEQ(FieldMiniLimit, v))
}

// MiniLimitIn applies the In predicate on the "mini_limit" field.
func MiniLimitIn(vs ...uint) predicate.AssetManager {
	return predicate.AssetManager(sql.FieldIn(FieldMiniLimit, vs...))
}

// MiniLimitNotIn applies the NotIn predicate on the "mini_limit" field.
func MiniLimitNotIn(vs ...uint) predicate.AssetManager {
	return predicate.AssetManager(sql.FieldNotIn(FieldMiniLimit, vs...))
}

// MiniLimitGT applies the GT predicate on the "mini_limit" field.
func MiniLimitGT(v uint) predicate.AssetManager {
	return predicate.AssetManager(sql.FieldGT(FieldMiniLimit, v))
}

// MiniLimitGTE applies the GTE predicate on the "mini_limit" field.
func MiniLimitGTE(v uint) predicate.AssetManager {
	return predicate.AssetManager(sql.FieldGTE(FieldMiniLimit, v))
}

// MiniLimitLT applies the LT predicate on the "mini_limit" field.
func MiniLimitLT(v uint) predicate.AssetManager {
	return predicate.AssetManager(sql.FieldLT(FieldMiniLimit, v))
}

// MiniLimitLTE applies the LTE predicate on the "mini_limit" field.
func MiniLimitLTE(v uint) predicate.AssetManager {
	return predicate.AssetManager(sql.FieldLTE(FieldMiniLimit, v))
}

// LastSigninAtEQ applies the EQ predicate on the "last_signin_at" field.
func LastSigninAtEQ(v time.Time) predicate.AssetManager {
	return predicate.AssetManager(sql.FieldEQ(FieldLastSigninAt, v))
}

// LastSigninAtNEQ applies the NEQ predicate on the "last_signin_at" field.
func LastSigninAtNEQ(v time.Time) predicate.AssetManager {
	return predicate.AssetManager(sql.FieldNEQ(FieldLastSigninAt, v))
}

// LastSigninAtIn applies the In predicate on the "last_signin_at" field.
func LastSigninAtIn(vs ...time.Time) predicate.AssetManager {
	return predicate.AssetManager(sql.FieldIn(FieldLastSigninAt, vs...))
}

// LastSigninAtNotIn applies the NotIn predicate on the "last_signin_at" field.
func LastSigninAtNotIn(vs ...time.Time) predicate.AssetManager {
	return predicate.AssetManager(sql.FieldNotIn(FieldLastSigninAt, vs...))
}

// LastSigninAtGT applies the GT predicate on the "last_signin_at" field.
func LastSigninAtGT(v time.Time) predicate.AssetManager {
	return predicate.AssetManager(sql.FieldGT(FieldLastSigninAt, v))
}

// LastSigninAtGTE applies the GTE predicate on the "last_signin_at" field.
func LastSigninAtGTE(v time.Time) predicate.AssetManager {
	return predicate.AssetManager(sql.FieldGTE(FieldLastSigninAt, v))
}

// LastSigninAtLT applies the LT predicate on the "last_signin_at" field.
func LastSigninAtLT(v time.Time) predicate.AssetManager {
	return predicate.AssetManager(sql.FieldLT(FieldLastSigninAt, v))
}

// LastSigninAtLTE applies the LTE predicate on the "last_signin_at" field.
func LastSigninAtLTE(v time.Time) predicate.AssetManager {
	return predicate.AssetManager(sql.FieldLTE(FieldLastSigninAt, v))
}

// LastSigninAtIsNil applies the IsNil predicate on the "last_signin_at" field.
func LastSigninAtIsNil() predicate.AssetManager {
	return predicate.AssetManager(sql.FieldIsNull(FieldLastSigninAt))
}

// LastSigninAtNotNil applies the NotNil predicate on the "last_signin_at" field.
func LastSigninAtNotNil() predicate.AssetManager {
	return predicate.AssetManager(sql.FieldNotNull(FieldLastSigninAt))
}

// WarehouseIDEQ applies the EQ predicate on the "warehouse_id" field.
func WarehouseIDEQ(v uint64) predicate.AssetManager {
	return predicate.AssetManager(sql.FieldEQ(FieldWarehouseID, v))
}

// WarehouseIDNEQ applies the NEQ predicate on the "warehouse_id" field.
func WarehouseIDNEQ(v uint64) predicate.AssetManager {
	return predicate.AssetManager(sql.FieldNEQ(FieldWarehouseID, v))
}

// WarehouseIDIn applies the In predicate on the "warehouse_id" field.
func WarehouseIDIn(vs ...uint64) predicate.AssetManager {
	return predicate.AssetManager(sql.FieldIn(FieldWarehouseID, vs...))
}

// WarehouseIDNotIn applies the NotIn predicate on the "warehouse_id" field.
func WarehouseIDNotIn(vs ...uint64) predicate.AssetManager {
	return predicate.AssetManager(sql.FieldNotIn(FieldWarehouseID, vs...))
}

// WarehouseIDIsNil applies the IsNil predicate on the "warehouse_id" field.
func WarehouseIDIsNil() predicate.AssetManager {
	return predicate.AssetManager(sql.FieldIsNull(FieldWarehouseID))
}

// WarehouseIDNotNil applies the NotNil predicate on the "warehouse_id" field.
func WarehouseIDNotNil() predicate.AssetManager {
	return predicate.AssetManager(sql.FieldNotNull(FieldWarehouseID))
}

// HasRole applies the HasEdge predicate on the "role" edge.
func HasRole() predicate.AssetManager {
	return predicate.AssetManager(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, RoleTable, RoleColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasRoleWith applies the HasEdge predicate on the "role" edge with a given conditions (other predicates).
func HasRoleWith(preds ...predicate.AssetRole) predicate.AssetManager {
	return predicate.AssetManager(func(s *sql.Selector) {
		step := newRoleStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasBelongWarehouses applies the HasEdge predicate on the "belong_warehouses" edge.
func HasBelongWarehouses() predicate.AssetManager {
	return predicate.AssetManager(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2M, true, BelongWarehousesTable, BelongWarehousesPrimaryKey...),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasBelongWarehousesWith applies the HasEdge predicate on the "belong_warehouses" edge with a given conditions (other predicates).
func HasBelongWarehousesWith(preds ...predicate.Warehouse) predicate.AssetManager {
	return predicate.AssetManager(func(s *sql.Selector) {
		step := newBelongWarehousesStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasDutyWarehouse applies the HasEdge predicate on the "duty_warehouse" edge.
func HasDutyWarehouse() predicate.AssetManager {
	return predicate.AssetManager(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, DutyWarehouseTable, DutyWarehouseColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasDutyWarehouseWith applies the HasEdge predicate on the "duty_warehouse" edge with a given conditions (other predicates).
func HasDutyWarehouseWith(preds ...predicate.Warehouse) predicate.AssetManager {
	return predicate.AssetManager(func(s *sql.Selector) {
		step := newDutyWarehouseStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.AssetManager) predicate.AssetManager {
	return predicate.AssetManager(sql.AndPredicates(predicates...))
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.AssetManager) predicate.AssetManager {
	return predicate.AssetManager(sql.OrPredicates(predicates...))
}

// Not applies the not operator on the given predicate.
func Not(p predicate.AssetManager) predicate.AssetManager {
	return predicate.AssetManager(sql.NotPredicates(p))
}