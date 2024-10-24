// Code generated by ent, DO NOT EDIT.

package manager

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/auroraride/aurservd/internal/ent/predicate"
)

// ID filters vertices based on their ID field.
func ID(id uint64) predicate.Manager {
	return predicate.Manager(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id uint64) predicate.Manager {
	return predicate.Manager(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id uint64) predicate.Manager {
	return predicate.Manager(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...uint64) predicate.Manager {
	return predicate.Manager(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...uint64) predicate.Manager {
	return predicate.Manager(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id uint64) predicate.Manager {
	return predicate.Manager(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id uint64) predicate.Manager {
	return predicate.Manager(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id uint64) predicate.Manager {
	return predicate.Manager(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id uint64) predicate.Manager {
	return predicate.Manager(sql.FieldLTE(FieldID, id))
}

// CreatedAt applies equality check predicate on the "created_at" field. It's identical to CreatedAtEQ.
func CreatedAt(v time.Time) predicate.Manager {
	return predicate.Manager(sql.FieldEQ(FieldCreatedAt, v))
}

// UpdatedAt applies equality check predicate on the "updated_at" field. It's identical to UpdatedAtEQ.
func UpdatedAt(v time.Time) predicate.Manager {
	return predicate.Manager(sql.FieldEQ(FieldUpdatedAt, v))
}

// DeletedAt applies equality check predicate on the "deleted_at" field. It's identical to DeletedAtEQ.
func DeletedAt(v time.Time) predicate.Manager {
	return predicate.Manager(sql.FieldEQ(FieldDeletedAt, v))
}

// Remark applies equality check predicate on the "remark" field. It's identical to RemarkEQ.
func Remark(v string) predicate.Manager {
	return predicate.Manager(sql.FieldEQ(FieldRemark, v))
}

// RoleID applies equality check predicate on the "role_id" field. It's identical to RoleIDEQ.
func RoleID(v uint64) predicate.Manager {
	return predicate.Manager(sql.FieldEQ(FieldRoleID, v))
}

// Phone applies equality check predicate on the "phone" field. It's identical to PhoneEQ.
func Phone(v string) predicate.Manager {
	return predicate.Manager(sql.FieldEQ(FieldPhone, v))
}

// Name applies equality check predicate on the "name" field. It's identical to NameEQ.
func Name(v string) predicate.Manager {
	return predicate.Manager(sql.FieldEQ(FieldName, v))
}

// Password applies equality check predicate on the "password" field. It's identical to PasswordEQ.
func Password(v string) predicate.Manager {
	return predicate.Manager(sql.FieldEQ(FieldPassword, v))
}

// LastSigninAt applies equality check predicate on the "last_signin_at" field. It's identical to LastSigninAtEQ.
func LastSigninAt(v time.Time) predicate.Manager {
	return predicate.Manager(sql.FieldEQ(FieldLastSigninAt, v))
}

// CreatedAtEQ applies the EQ predicate on the "created_at" field.
func CreatedAtEQ(v time.Time) predicate.Manager {
	return predicate.Manager(sql.FieldEQ(FieldCreatedAt, v))
}

// CreatedAtNEQ applies the NEQ predicate on the "created_at" field.
func CreatedAtNEQ(v time.Time) predicate.Manager {
	return predicate.Manager(sql.FieldNEQ(FieldCreatedAt, v))
}

// CreatedAtIn applies the In predicate on the "created_at" field.
func CreatedAtIn(vs ...time.Time) predicate.Manager {
	return predicate.Manager(sql.FieldIn(FieldCreatedAt, vs...))
}

// CreatedAtNotIn applies the NotIn predicate on the "created_at" field.
func CreatedAtNotIn(vs ...time.Time) predicate.Manager {
	return predicate.Manager(sql.FieldNotIn(FieldCreatedAt, vs...))
}

// CreatedAtGT applies the GT predicate on the "created_at" field.
func CreatedAtGT(v time.Time) predicate.Manager {
	return predicate.Manager(sql.FieldGT(FieldCreatedAt, v))
}

// CreatedAtGTE applies the GTE predicate on the "created_at" field.
func CreatedAtGTE(v time.Time) predicate.Manager {
	return predicate.Manager(sql.FieldGTE(FieldCreatedAt, v))
}

// CreatedAtLT applies the LT predicate on the "created_at" field.
func CreatedAtLT(v time.Time) predicate.Manager {
	return predicate.Manager(sql.FieldLT(FieldCreatedAt, v))
}

// CreatedAtLTE applies the LTE predicate on the "created_at" field.
func CreatedAtLTE(v time.Time) predicate.Manager {
	return predicate.Manager(sql.FieldLTE(FieldCreatedAt, v))
}

// UpdatedAtEQ applies the EQ predicate on the "updated_at" field.
func UpdatedAtEQ(v time.Time) predicate.Manager {
	return predicate.Manager(sql.FieldEQ(FieldUpdatedAt, v))
}

// UpdatedAtNEQ applies the NEQ predicate on the "updated_at" field.
func UpdatedAtNEQ(v time.Time) predicate.Manager {
	return predicate.Manager(sql.FieldNEQ(FieldUpdatedAt, v))
}

// UpdatedAtIn applies the In predicate on the "updated_at" field.
func UpdatedAtIn(vs ...time.Time) predicate.Manager {
	return predicate.Manager(sql.FieldIn(FieldUpdatedAt, vs...))
}

// UpdatedAtNotIn applies the NotIn predicate on the "updated_at" field.
func UpdatedAtNotIn(vs ...time.Time) predicate.Manager {
	return predicate.Manager(sql.FieldNotIn(FieldUpdatedAt, vs...))
}

// UpdatedAtGT applies the GT predicate on the "updated_at" field.
func UpdatedAtGT(v time.Time) predicate.Manager {
	return predicate.Manager(sql.FieldGT(FieldUpdatedAt, v))
}

// UpdatedAtGTE applies the GTE predicate on the "updated_at" field.
func UpdatedAtGTE(v time.Time) predicate.Manager {
	return predicate.Manager(sql.FieldGTE(FieldUpdatedAt, v))
}

// UpdatedAtLT applies the LT predicate on the "updated_at" field.
func UpdatedAtLT(v time.Time) predicate.Manager {
	return predicate.Manager(sql.FieldLT(FieldUpdatedAt, v))
}

// UpdatedAtLTE applies the LTE predicate on the "updated_at" field.
func UpdatedAtLTE(v time.Time) predicate.Manager {
	return predicate.Manager(sql.FieldLTE(FieldUpdatedAt, v))
}

// DeletedAtEQ applies the EQ predicate on the "deleted_at" field.
func DeletedAtEQ(v time.Time) predicate.Manager {
	return predicate.Manager(sql.FieldEQ(FieldDeletedAt, v))
}

// DeletedAtNEQ applies the NEQ predicate on the "deleted_at" field.
func DeletedAtNEQ(v time.Time) predicate.Manager {
	return predicate.Manager(sql.FieldNEQ(FieldDeletedAt, v))
}

// DeletedAtIn applies the In predicate on the "deleted_at" field.
func DeletedAtIn(vs ...time.Time) predicate.Manager {
	return predicate.Manager(sql.FieldIn(FieldDeletedAt, vs...))
}

// DeletedAtNotIn applies the NotIn predicate on the "deleted_at" field.
func DeletedAtNotIn(vs ...time.Time) predicate.Manager {
	return predicate.Manager(sql.FieldNotIn(FieldDeletedAt, vs...))
}

// DeletedAtGT applies the GT predicate on the "deleted_at" field.
func DeletedAtGT(v time.Time) predicate.Manager {
	return predicate.Manager(sql.FieldGT(FieldDeletedAt, v))
}

// DeletedAtGTE applies the GTE predicate on the "deleted_at" field.
func DeletedAtGTE(v time.Time) predicate.Manager {
	return predicate.Manager(sql.FieldGTE(FieldDeletedAt, v))
}

// DeletedAtLT applies the LT predicate on the "deleted_at" field.
func DeletedAtLT(v time.Time) predicate.Manager {
	return predicate.Manager(sql.FieldLT(FieldDeletedAt, v))
}

// DeletedAtLTE applies the LTE predicate on the "deleted_at" field.
func DeletedAtLTE(v time.Time) predicate.Manager {
	return predicate.Manager(sql.FieldLTE(FieldDeletedAt, v))
}

// DeletedAtIsNil applies the IsNil predicate on the "deleted_at" field.
func DeletedAtIsNil() predicate.Manager {
	return predicate.Manager(sql.FieldIsNull(FieldDeletedAt))
}

// DeletedAtNotNil applies the NotNil predicate on the "deleted_at" field.
func DeletedAtNotNil() predicate.Manager {
	return predicate.Manager(sql.FieldNotNull(FieldDeletedAt))
}

// CreatorIsNil applies the IsNil predicate on the "creator" field.
func CreatorIsNil() predicate.Manager {
	return predicate.Manager(sql.FieldIsNull(FieldCreator))
}

// CreatorNotNil applies the NotNil predicate on the "creator" field.
func CreatorNotNil() predicate.Manager {
	return predicate.Manager(sql.FieldNotNull(FieldCreator))
}

// LastModifierIsNil applies the IsNil predicate on the "last_modifier" field.
func LastModifierIsNil() predicate.Manager {
	return predicate.Manager(sql.FieldIsNull(FieldLastModifier))
}

// LastModifierNotNil applies the NotNil predicate on the "last_modifier" field.
func LastModifierNotNil() predicate.Manager {
	return predicate.Manager(sql.FieldNotNull(FieldLastModifier))
}

// RemarkEQ applies the EQ predicate on the "remark" field.
func RemarkEQ(v string) predicate.Manager {
	return predicate.Manager(sql.FieldEQ(FieldRemark, v))
}

// RemarkNEQ applies the NEQ predicate on the "remark" field.
func RemarkNEQ(v string) predicate.Manager {
	return predicate.Manager(sql.FieldNEQ(FieldRemark, v))
}

// RemarkIn applies the In predicate on the "remark" field.
func RemarkIn(vs ...string) predicate.Manager {
	return predicate.Manager(sql.FieldIn(FieldRemark, vs...))
}

// RemarkNotIn applies the NotIn predicate on the "remark" field.
func RemarkNotIn(vs ...string) predicate.Manager {
	return predicate.Manager(sql.FieldNotIn(FieldRemark, vs...))
}

// RemarkGT applies the GT predicate on the "remark" field.
func RemarkGT(v string) predicate.Manager {
	return predicate.Manager(sql.FieldGT(FieldRemark, v))
}

// RemarkGTE applies the GTE predicate on the "remark" field.
func RemarkGTE(v string) predicate.Manager {
	return predicate.Manager(sql.FieldGTE(FieldRemark, v))
}

// RemarkLT applies the LT predicate on the "remark" field.
func RemarkLT(v string) predicate.Manager {
	return predicate.Manager(sql.FieldLT(FieldRemark, v))
}

// RemarkLTE applies the LTE predicate on the "remark" field.
func RemarkLTE(v string) predicate.Manager {
	return predicate.Manager(sql.FieldLTE(FieldRemark, v))
}

// RemarkContains applies the Contains predicate on the "remark" field.
func RemarkContains(v string) predicate.Manager {
	return predicate.Manager(sql.FieldContains(FieldRemark, v))
}

// RemarkHasPrefix applies the HasPrefix predicate on the "remark" field.
func RemarkHasPrefix(v string) predicate.Manager {
	return predicate.Manager(sql.FieldHasPrefix(FieldRemark, v))
}

// RemarkHasSuffix applies the HasSuffix predicate on the "remark" field.
func RemarkHasSuffix(v string) predicate.Manager {
	return predicate.Manager(sql.FieldHasSuffix(FieldRemark, v))
}

// RemarkIsNil applies the IsNil predicate on the "remark" field.
func RemarkIsNil() predicate.Manager {
	return predicate.Manager(sql.FieldIsNull(FieldRemark))
}

// RemarkNotNil applies the NotNil predicate on the "remark" field.
func RemarkNotNil() predicate.Manager {
	return predicate.Manager(sql.FieldNotNull(FieldRemark))
}

// RemarkEqualFold applies the EqualFold predicate on the "remark" field.
func RemarkEqualFold(v string) predicate.Manager {
	return predicate.Manager(sql.FieldEqualFold(FieldRemark, v))
}

// RemarkContainsFold applies the ContainsFold predicate on the "remark" field.
func RemarkContainsFold(v string) predicate.Manager {
	return predicate.Manager(sql.FieldContainsFold(FieldRemark, v))
}

// RoleIDEQ applies the EQ predicate on the "role_id" field.
func RoleIDEQ(v uint64) predicate.Manager {
	return predicate.Manager(sql.FieldEQ(FieldRoleID, v))
}

// RoleIDNEQ applies the NEQ predicate on the "role_id" field.
func RoleIDNEQ(v uint64) predicate.Manager {
	return predicate.Manager(sql.FieldNEQ(FieldRoleID, v))
}

// RoleIDIn applies the In predicate on the "role_id" field.
func RoleIDIn(vs ...uint64) predicate.Manager {
	return predicate.Manager(sql.FieldIn(FieldRoleID, vs...))
}

// RoleIDNotIn applies the NotIn predicate on the "role_id" field.
func RoleIDNotIn(vs ...uint64) predicate.Manager {
	return predicate.Manager(sql.FieldNotIn(FieldRoleID, vs...))
}

// RoleIDIsNil applies the IsNil predicate on the "role_id" field.
func RoleIDIsNil() predicate.Manager {
	return predicate.Manager(sql.FieldIsNull(FieldRoleID))
}

// RoleIDNotNil applies the NotNil predicate on the "role_id" field.
func RoleIDNotNil() predicate.Manager {
	return predicate.Manager(sql.FieldNotNull(FieldRoleID))
}

// PhoneEQ applies the EQ predicate on the "phone" field.
func PhoneEQ(v string) predicate.Manager {
	return predicate.Manager(sql.FieldEQ(FieldPhone, v))
}

// PhoneNEQ applies the NEQ predicate on the "phone" field.
func PhoneNEQ(v string) predicate.Manager {
	return predicate.Manager(sql.FieldNEQ(FieldPhone, v))
}

// PhoneIn applies the In predicate on the "phone" field.
func PhoneIn(vs ...string) predicate.Manager {
	return predicate.Manager(sql.FieldIn(FieldPhone, vs...))
}

// PhoneNotIn applies the NotIn predicate on the "phone" field.
func PhoneNotIn(vs ...string) predicate.Manager {
	return predicate.Manager(sql.FieldNotIn(FieldPhone, vs...))
}

// PhoneGT applies the GT predicate on the "phone" field.
func PhoneGT(v string) predicate.Manager {
	return predicate.Manager(sql.FieldGT(FieldPhone, v))
}

// PhoneGTE applies the GTE predicate on the "phone" field.
func PhoneGTE(v string) predicate.Manager {
	return predicate.Manager(sql.FieldGTE(FieldPhone, v))
}

// PhoneLT applies the LT predicate on the "phone" field.
func PhoneLT(v string) predicate.Manager {
	return predicate.Manager(sql.FieldLT(FieldPhone, v))
}

// PhoneLTE applies the LTE predicate on the "phone" field.
func PhoneLTE(v string) predicate.Manager {
	return predicate.Manager(sql.FieldLTE(FieldPhone, v))
}

// PhoneContains applies the Contains predicate on the "phone" field.
func PhoneContains(v string) predicate.Manager {
	return predicate.Manager(sql.FieldContains(FieldPhone, v))
}

// PhoneHasPrefix applies the HasPrefix predicate on the "phone" field.
func PhoneHasPrefix(v string) predicate.Manager {
	return predicate.Manager(sql.FieldHasPrefix(FieldPhone, v))
}

// PhoneHasSuffix applies the HasSuffix predicate on the "phone" field.
func PhoneHasSuffix(v string) predicate.Manager {
	return predicate.Manager(sql.FieldHasSuffix(FieldPhone, v))
}

// PhoneEqualFold applies the EqualFold predicate on the "phone" field.
func PhoneEqualFold(v string) predicate.Manager {
	return predicate.Manager(sql.FieldEqualFold(FieldPhone, v))
}

// PhoneContainsFold applies the ContainsFold predicate on the "phone" field.
func PhoneContainsFold(v string) predicate.Manager {
	return predicate.Manager(sql.FieldContainsFold(FieldPhone, v))
}

// NameEQ applies the EQ predicate on the "name" field.
func NameEQ(v string) predicate.Manager {
	return predicate.Manager(sql.FieldEQ(FieldName, v))
}

// NameNEQ applies the NEQ predicate on the "name" field.
func NameNEQ(v string) predicate.Manager {
	return predicate.Manager(sql.FieldNEQ(FieldName, v))
}

// NameIn applies the In predicate on the "name" field.
func NameIn(vs ...string) predicate.Manager {
	return predicate.Manager(sql.FieldIn(FieldName, vs...))
}

// NameNotIn applies the NotIn predicate on the "name" field.
func NameNotIn(vs ...string) predicate.Manager {
	return predicate.Manager(sql.FieldNotIn(FieldName, vs...))
}

// NameGT applies the GT predicate on the "name" field.
func NameGT(v string) predicate.Manager {
	return predicate.Manager(sql.FieldGT(FieldName, v))
}

// NameGTE applies the GTE predicate on the "name" field.
func NameGTE(v string) predicate.Manager {
	return predicate.Manager(sql.FieldGTE(FieldName, v))
}

// NameLT applies the LT predicate on the "name" field.
func NameLT(v string) predicate.Manager {
	return predicate.Manager(sql.FieldLT(FieldName, v))
}

// NameLTE applies the LTE predicate on the "name" field.
func NameLTE(v string) predicate.Manager {
	return predicate.Manager(sql.FieldLTE(FieldName, v))
}

// NameContains applies the Contains predicate on the "name" field.
func NameContains(v string) predicate.Manager {
	return predicate.Manager(sql.FieldContains(FieldName, v))
}

// NameHasPrefix applies the HasPrefix predicate on the "name" field.
func NameHasPrefix(v string) predicate.Manager {
	return predicate.Manager(sql.FieldHasPrefix(FieldName, v))
}

// NameHasSuffix applies the HasSuffix predicate on the "name" field.
func NameHasSuffix(v string) predicate.Manager {
	return predicate.Manager(sql.FieldHasSuffix(FieldName, v))
}

// NameEqualFold applies the EqualFold predicate on the "name" field.
func NameEqualFold(v string) predicate.Manager {
	return predicate.Manager(sql.FieldEqualFold(FieldName, v))
}

// NameContainsFold applies the ContainsFold predicate on the "name" field.
func NameContainsFold(v string) predicate.Manager {
	return predicate.Manager(sql.FieldContainsFold(FieldName, v))
}

// PasswordEQ applies the EQ predicate on the "password" field.
func PasswordEQ(v string) predicate.Manager {
	return predicate.Manager(sql.FieldEQ(FieldPassword, v))
}

// PasswordNEQ applies the NEQ predicate on the "password" field.
func PasswordNEQ(v string) predicate.Manager {
	return predicate.Manager(sql.FieldNEQ(FieldPassword, v))
}

// PasswordIn applies the In predicate on the "password" field.
func PasswordIn(vs ...string) predicate.Manager {
	return predicate.Manager(sql.FieldIn(FieldPassword, vs...))
}

// PasswordNotIn applies the NotIn predicate on the "password" field.
func PasswordNotIn(vs ...string) predicate.Manager {
	return predicate.Manager(sql.FieldNotIn(FieldPassword, vs...))
}

// PasswordGT applies the GT predicate on the "password" field.
func PasswordGT(v string) predicate.Manager {
	return predicate.Manager(sql.FieldGT(FieldPassword, v))
}

// PasswordGTE applies the GTE predicate on the "password" field.
func PasswordGTE(v string) predicate.Manager {
	return predicate.Manager(sql.FieldGTE(FieldPassword, v))
}

// PasswordLT applies the LT predicate on the "password" field.
func PasswordLT(v string) predicate.Manager {
	return predicate.Manager(sql.FieldLT(FieldPassword, v))
}

// PasswordLTE applies the LTE predicate on the "password" field.
func PasswordLTE(v string) predicate.Manager {
	return predicate.Manager(sql.FieldLTE(FieldPassword, v))
}

// PasswordContains applies the Contains predicate on the "password" field.
func PasswordContains(v string) predicate.Manager {
	return predicate.Manager(sql.FieldContains(FieldPassword, v))
}

// PasswordHasPrefix applies the HasPrefix predicate on the "password" field.
func PasswordHasPrefix(v string) predicate.Manager {
	return predicate.Manager(sql.FieldHasPrefix(FieldPassword, v))
}

// PasswordHasSuffix applies the HasSuffix predicate on the "password" field.
func PasswordHasSuffix(v string) predicate.Manager {
	return predicate.Manager(sql.FieldHasSuffix(FieldPassword, v))
}

// PasswordEqualFold applies the EqualFold predicate on the "password" field.
func PasswordEqualFold(v string) predicate.Manager {
	return predicate.Manager(sql.FieldEqualFold(FieldPassword, v))
}

// PasswordContainsFold applies the ContainsFold predicate on the "password" field.
func PasswordContainsFold(v string) predicate.Manager {
	return predicate.Manager(sql.FieldContainsFold(FieldPassword, v))
}

// LastSigninAtEQ applies the EQ predicate on the "last_signin_at" field.
func LastSigninAtEQ(v time.Time) predicate.Manager {
	return predicate.Manager(sql.FieldEQ(FieldLastSigninAt, v))
}

// LastSigninAtNEQ applies the NEQ predicate on the "last_signin_at" field.
func LastSigninAtNEQ(v time.Time) predicate.Manager {
	return predicate.Manager(sql.FieldNEQ(FieldLastSigninAt, v))
}

// LastSigninAtIn applies the In predicate on the "last_signin_at" field.
func LastSigninAtIn(vs ...time.Time) predicate.Manager {
	return predicate.Manager(sql.FieldIn(FieldLastSigninAt, vs...))
}

// LastSigninAtNotIn applies the NotIn predicate on the "last_signin_at" field.
func LastSigninAtNotIn(vs ...time.Time) predicate.Manager {
	return predicate.Manager(sql.FieldNotIn(FieldLastSigninAt, vs...))
}

// LastSigninAtGT applies the GT predicate on the "last_signin_at" field.
func LastSigninAtGT(v time.Time) predicate.Manager {
	return predicate.Manager(sql.FieldGT(FieldLastSigninAt, v))
}

// LastSigninAtGTE applies the GTE predicate on the "last_signin_at" field.
func LastSigninAtGTE(v time.Time) predicate.Manager {
	return predicate.Manager(sql.FieldGTE(FieldLastSigninAt, v))
}

// LastSigninAtLT applies the LT predicate on the "last_signin_at" field.
func LastSigninAtLT(v time.Time) predicate.Manager {
	return predicate.Manager(sql.FieldLT(FieldLastSigninAt, v))
}

// LastSigninAtLTE applies the LTE predicate on the "last_signin_at" field.
func LastSigninAtLTE(v time.Time) predicate.Manager {
	return predicate.Manager(sql.FieldLTE(FieldLastSigninAt, v))
}

// LastSigninAtIsNil applies the IsNil predicate on the "last_signin_at" field.
func LastSigninAtIsNil() predicate.Manager {
	return predicate.Manager(sql.FieldIsNull(FieldLastSigninAt))
}

// LastSigninAtNotNil applies the NotNil predicate on the "last_signin_at" field.
func LastSigninAtNotNil() predicate.Manager {
	return predicate.Manager(sql.FieldNotNull(FieldLastSigninAt))
}

// HasRole applies the HasEdge predicate on the "role" edge.
func HasRole() predicate.Manager {
	return predicate.Manager(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, RoleTable, RoleColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasRoleWith applies the HasEdge predicate on the "role" edge with a given conditions (other predicates).
func HasRoleWith(preds ...predicate.Role) predicate.Manager {
	return predicate.Manager(func(s *sql.Selector) {
		step := newRoleStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.Manager) predicate.Manager {
	return predicate.Manager(sql.AndPredicates(predicates...))
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.Manager) predicate.Manager {
	return predicate.Manager(sql.OrPredicates(predicates...))
}

// Not applies the not operator on the given predicate.
func Not(p predicate.Manager) predicate.Manager {
	return predicate.Manager(sql.NotPredicates(p))
}
