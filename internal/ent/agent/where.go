// Code generated by ent, DO NOT EDIT.

package agent

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/auroraride/aurservd/internal/ent/predicate"
)

// ID filters vertices based on their ID field.
func ID(id uint64) predicate.Agent {
	return predicate.Agent(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldID), id))
	})
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id uint64) predicate.Agent {
	return predicate.Agent(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldID), id))
	})
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id uint64) predicate.Agent {
	return predicate.Agent(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldID), id))
	})
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...uint64) predicate.Agent {
	return predicate.Agent(func(s *sql.Selector) {
		v := make([]interface{}, len(ids))
		for i := range v {
			v[i] = ids[i]
		}
		s.Where(sql.In(s.C(FieldID), v...))
	})
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...uint64) predicate.Agent {
	return predicate.Agent(func(s *sql.Selector) {
		v := make([]interface{}, len(ids))
		for i := range v {
			v[i] = ids[i]
		}
		s.Where(sql.NotIn(s.C(FieldID), v...))
	})
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id uint64) predicate.Agent {
	return predicate.Agent(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldID), id))
	})
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id uint64) predicate.Agent {
	return predicate.Agent(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldID), id))
	})
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id uint64) predicate.Agent {
	return predicate.Agent(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldID), id))
	})
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id uint64) predicate.Agent {
	return predicate.Agent(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldID), id))
	})
}

// CreatedAt applies equality check predicate on the "created_at" field. It's identical to CreatedAtEQ.
func CreatedAt(v time.Time) predicate.Agent {
	return predicate.Agent(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldCreatedAt), v))
	})
}

// UpdatedAt applies equality check predicate on the "updated_at" field. It's identical to UpdatedAtEQ.
func UpdatedAt(v time.Time) predicate.Agent {
	return predicate.Agent(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldUpdatedAt), v))
	})
}

// DeletedAt applies equality check predicate on the "deleted_at" field. It's identical to DeletedAtEQ.
func DeletedAt(v time.Time) predicate.Agent {
	return predicate.Agent(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldDeletedAt), v))
	})
}

// Remark applies equality check predicate on the "remark" field. It's identical to RemarkEQ.
func Remark(v string) predicate.Agent {
	return predicate.Agent(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldRemark), v))
	})
}

// EnterpriseID applies equality check predicate on the "enterprise_id" field. It's identical to EnterpriseIDEQ.
func EnterpriseID(v uint64) predicate.Agent {
	return predicate.Agent(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldEnterpriseID), v))
	})
}

// Name applies equality check predicate on the "name" field. It's identical to NameEQ.
func Name(v string) predicate.Agent {
	return predicate.Agent(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldName), v))
	})
}

// Phone applies equality check predicate on the "phone" field. It's identical to PhoneEQ.
func Phone(v string) predicate.Agent {
	return predicate.Agent(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldPhone), v))
	})
}

// Password applies equality check predicate on the "password" field. It's identical to PasswordEQ.
func Password(v string) predicate.Agent {
	return predicate.Agent(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldPassword), v))
	})
}

// CreatedAtEQ applies the EQ predicate on the "created_at" field.
func CreatedAtEQ(v time.Time) predicate.Agent {
	return predicate.Agent(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldCreatedAt), v))
	})
}

// CreatedAtNEQ applies the NEQ predicate on the "created_at" field.
func CreatedAtNEQ(v time.Time) predicate.Agent {
	return predicate.Agent(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldCreatedAt), v))
	})
}

// CreatedAtIn applies the In predicate on the "created_at" field.
func CreatedAtIn(vs ...time.Time) predicate.Agent {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Agent(func(s *sql.Selector) {
		s.Where(sql.In(s.C(FieldCreatedAt), v...))
	})
}

// CreatedAtNotIn applies the NotIn predicate on the "created_at" field.
func CreatedAtNotIn(vs ...time.Time) predicate.Agent {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Agent(func(s *sql.Selector) {
		s.Where(sql.NotIn(s.C(FieldCreatedAt), v...))
	})
}

// CreatedAtGT applies the GT predicate on the "created_at" field.
func CreatedAtGT(v time.Time) predicate.Agent {
	return predicate.Agent(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldCreatedAt), v))
	})
}

// CreatedAtGTE applies the GTE predicate on the "created_at" field.
func CreatedAtGTE(v time.Time) predicate.Agent {
	return predicate.Agent(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldCreatedAt), v))
	})
}

// CreatedAtLT applies the LT predicate on the "created_at" field.
func CreatedAtLT(v time.Time) predicate.Agent {
	return predicate.Agent(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldCreatedAt), v))
	})
}

// CreatedAtLTE applies the LTE predicate on the "created_at" field.
func CreatedAtLTE(v time.Time) predicate.Agent {
	return predicate.Agent(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldCreatedAt), v))
	})
}

// UpdatedAtEQ applies the EQ predicate on the "updated_at" field.
func UpdatedAtEQ(v time.Time) predicate.Agent {
	return predicate.Agent(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldUpdatedAt), v))
	})
}

// UpdatedAtNEQ applies the NEQ predicate on the "updated_at" field.
func UpdatedAtNEQ(v time.Time) predicate.Agent {
	return predicate.Agent(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldUpdatedAt), v))
	})
}

// UpdatedAtIn applies the In predicate on the "updated_at" field.
func UpdatedAtIn(vs ...time.Time) predicate.Agent {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Agent(func(s *sql.Selector) {
		s.Where(sql.In(s.C(FieldUpdatedAt), v...))
	})
}

// UpdatedAtNotIn applies the NotIn predicate on the "updated_at" field.
func UpdatedAtNotIn(vs ...time.Time) predicate.Agent {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Agent(func(s *sql.Selector) {
		s.Where(sql.NotIn(s.C(FieldUpdatedAt), v...))
	})
}

// UpdatedAtGT applies the GT predicate on the "updated_at" field.
func UpdatedAtGT(v time.Time) predicate.Agent {
	return predicate.Agent(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldUpdatedAt), v))
	})
}

// UpdatedAtGTE applies the GTE predicate on the "updated_at" field.
func UpdatedAtGTE(v time.Time) predicate.Agent {
	return predicate.Agent(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldUpdatedAt), v))
	})
}

// UpdatedAtLT applies the LT predicate on the "updated_at" field.
func UpdatedAtLT(v time.Time) predicate.Agent {
	return predicate.Agent(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldUpdatedAt), v))
	})
}

// UpdatedAtLTE applies the LTE predicate on the "updated_at" field.
func UpdatedAtLTE(v time.Time) predicate.Agent {
	return predicate.Agent(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldUpdatedAt), v))
	})
}

// DeletedAtEQ applies the EQ predicate on the "deleted_at" field.
func DeletedAtEQ(v time.Time) predicate.Agent {
	return predicate.Agent(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldDeletedAt), v))
	})
}

// DeletedAtNEQ applies the NEQ predicate on the "deleted_at" field.
func DeletedAtNEQ(v time.Time) predicate.Agent {
	return predicate.Agent(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldDeletedAt), v))
	})
}

// DeletedAtIn applies the In predicate on the "deleted_at" field.
func DeletedAtIn(vs ...time.Time) predicate.Agent {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Agent(func(s *sql.Selector) {
		s.Where(sql.In(s.C(FieldDeletedAt), v...))
	})
}

// DeletedAtNotIn applies the NotIn predicate on the "deleted_at" field.
func DeletedAtNotIn(vs ...time.Time) predicate.Agent {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Agent(func(s *sql.Selector) {
		s.Where(sql.NotIn(s.C(FieldDeletedAt), v...))
	})
}

// DeletedAtGT applies the GT predicate on the "deleted_at" field.
func DeletedAtGT(v time.Time) predicate.Agent {
	return predicate.Agent(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldDeletedAt), v))
	})
}

// DeletedAtGTE applies the GTE predicate on the "deleted_at" field.
func DeletedAtGTE(v time.Time) predicate.Agent {
	return predicate.Agent(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldDeletedAt), v))
	})
}

// DeletedAtLT applies the LT predicate on the "deleted_at" field.
func DeletedAtLT(v time.Time) predicate.Agent {
	return predicate.Agent(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldDeletedAt), v))
	})
}

// DeletedAtLTE applies the LTE predicate on the "deleted_at" field.
func DeletedAtLTE(v time.Time) predicate.Agent {
	return predicate.Agent(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldDeletedAt), v))
	})
}

// DeletedAtIsNil applies the IsNil predicate on the "deleted_at" field.
func DeletedAtIsNil() predicate.Agent {
	return predicate.Agent(func(s *sql.Selector) {
		s.Where(sql.IsNull(s.C(FieldDeletedAt)))
	})
}

// DeletedAtNotNil applies the NotNil predicate on the "deleted_at" field.
func DeletedAtNotNil() predicate.Agent {
	return predicate.Agent(func(s *sql.Selector) {
		s.Where(sql.NotNull(s.C(FieldDeletedAt)))
	})
}

// CreatorIsNil applies the IsNil predicate on the "creator" field.
func CreatorIsNil() predicate.Agent {
	return predicate.Agent(func(s *sql.Selector) {
		s.Where(sql.IsNull(s.C(FieldCreator)))
	})
}

// CreatorNotNil applies the NotNil predicate on the "creator" field.
func CreatorNotNil() predicate.Agent {
	return predicate.Agent(func(s *sql.Selector) {
		s.Where(sql.NotNull(s.C(FieldCreator)))
	})
}

// LastModifierIsNil applies the IsNil predicate on the "last_modifier" field.
func LastModifierIsNil() predicate.Agent {
	return predicate.Agent(func(s *sql.Selector) {
		s.Where(sql.IsNull(s.C(FieldLastModifier)))
	})
}

// LastModifierNotNil applies the NotNil predicate on the "last_modifier" field.
func LastModifierNotNil() predicate.Agent {
	return predicate.Agent(func(s *sql.Selector) {
		s.Where(sql.NotNull(s.C(FieldLastModifier)))
	})
}

// RemarkEQ applies the EQ predicate on the "remark" field.
func RemarkEQ(v string) predicate.Agent {
	return predicate.Agent(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldRemark), v))
	})
}

// RemarkNEQ applies the NEQ predicate on the "remark" field.
func RemarkNEQ(v string) predicate.Agent {
	return predicate.Agent(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldRemark), v))
	})
}

// RemarkIn applies the In predicate on the "remark" field.
func RemarkIn(vs ...string) predicate.Agent {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Agent(func(s *sql.Selector) {
		s.Where(sql.In(s.C(FieldRemark), v...))
	})
}

// RemarkNotIn applies the NotIn predicate on the "remark" field.
func RemarkNotIn(vs ...string) predicate.Agent {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Agent(func(s *sql.Selector) {
		s.Where(sql.NotIn(s.C(FieldRemark), v...))
	})
}

// RemarkGT applies the GT predicate on the "remark" field.
func RemarkGT(v string) predicate.Agent {
	return predicate.Agent(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldRemark), v))
	})
}

// RemarkGTE applies the GTE predicate on the "remark" field.
func RemarkGTE(v string) predicate.Agent {
	return predicate.Agent(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldRemark), v))
	})
}

// RemarkLT applies the LT predicate on the "remark" field.
func RemarkLT(v string) predicate.Agent {
	return predicate.Agent(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldRemark), v))
	})
}

// RemarkLTE applies the LTE predicate on the "remark" field.
func RemarkLTE(v string) predicate.Agent {
	return predicate.Agent(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldRemark), v))
	})
}

// RemarkContains applies the Contains predicate on the "remark" field.
func RemarkContains(v string) predicate.Agent {
	return predicate.Agent(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldRemark), v))
	})
}

// RemarkHasPrefix applies the HasPrefix predicate on the "remark" field.
func RemarkHasPrefix(v string) predicate.Agent {
	return predicate.Agent(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldRemark), v))
	})
}

// RemarkHasSuffix applies the HasSuffix predicate on the "remark" field.
func RemarkHasSuffix(v string) predicate.Agent {
	return predicate.Agent(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldRemark), v))
	})
}

// RemarkIsNil applies the IsNil predicate on the "remark" field.
func RemarkIsNil() predicate.Agent {
	return predicate.Agent(func(s *sql.Selector) {
		s.Where(sql.IsNull(s.C(FieldRemark)))
	})
}

// RemarkNotNil applies the NotNil predicate on the "remark" field.
func RemarkNotNil() predicate.Agent {
	return predicate.Agent(func(s *sql.Selector) {
		s.Where(sql.NotNull(s.C(FieldRemark)))
	})
}

// RemarkEqualFold applies the EqualFold predicate on the "remark" field.
func RemarkEqualFold(v string) predicate.Agent {
	return predicate.Agent(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldRemark), v))
	})
}

// RemarkContainsFold applies the ContainsFold predicate on the "remark" field.
func RemarkContainsFold(v string) predicate.Agent {
	return predicate.Agent(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldRemark), v))
	})
}

// EnterpriseIDEQ applies the EQ predicate on the "enterprise_id" field.
func EnterpriseIDEQ(v uint64) predicate.Agent {
	return predicate.Agent(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldEnterpriseID), v))
	})
}

// EnterpriseIDNEQ applies the NEQ predicate on the "enterprise_id" field.
func EnterpriseIDNEQ(v uint64) predicate.Agent {
	return predicate.Agent(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldEnterpriseID), v))
	})
}

// EnterpriseIDIn applies the In predicate on the "enterprise_id" field.
func EnterpriseIDIn(vs ...uint64) predicate.Agent {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Agent(func(s *sql.Selector) {
		s.Where(sql.In(s.C(FieldEnterpriseID), v...))
	})
}

// EnterpriseIDNotIn applies the NotIn predicate on the "enterprise_id" field.
func EnterpriseIDNotIn(vs ...uint64) predicate.Agent {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Agent(func(s *sql.Selector) {
		s.Where(sql.NotIn(s.C(FieldEnterpriseID), v...))
	})
}

// NameEQ applies the EQ predicate on the "name" field.
func NameEQ(v string) predicate.Agent {
	return predicate.Agent(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldName), v))
	})
}

// NameNEQ applies the NEQ predicate on the "name" field.
func NameNEQ(v string) predicate.Agent {
	return predicate.Agent(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldName), v))
	})
}

// NameIn applies the In predicate on the "name" field.
func NameIn(vs ...string) predicate.Agent {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Agent(func(s *sql.Selector) {
		s.Where(sql.In(s.C(FieldName), v...))
	})
}

// NameNotIn applies the NotIn predicate on the "name" field.
func NameNotIn(vs ...string) predicate.Agent {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Agent(func(s *sql.Selector) {
		s.Where(sql.NotIn(s.C(FieldName), v...))
	})
}

// NameGT applies the GT predicate on the "name" field.
func NameGT(v string) predicate.Agent {
	return predicate.Agent(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldName), v))
	})
}

// NameGTE applies the GTE predicate on the "name" field.
func NameGTE(v string) predicate.Agent {
	return predicate.Agent(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldName), v))
	})
}

// NameLT applies the LT predicate on the "name" field.
func NameLT(v string) predicate.Agent {
	return predicate.Agent(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldName), v))
	})
}

// NameLTE applies the LTE predicate on the "name" field.
func NameLTE(v string) predicate.Agent {
	return predicate.Agent(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldName), v))
	})
}

// NameContains applies the Contains predicate on the "name" field.
func NameContains(v string) predicate.Agent {
	return predicate.Agent(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldName), v))
	})
}

// NameHasPrefix applies the HasPrefix predicate on the "name" field.
func NameHasPrefix(v string) predicate.Agent {
	return predicate.Agent(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldName), v))
	})
}

// NameHasSuffix applies the HasSuffix predicate on the "name" field.
func NameHasSuffix(v string) predicate.Agent {
	return predicate.Agent(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldName), v))
	})
}

// NameEqualFold applies the EqualFold predicate on the "name" field.
func NameEqualFold(v string) predicate.Agent {
	return predicate.Agent(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldName), v))
	})
}

// NameContainsFold applies the ContainsFold predicate on the "name" field.
func NameContainsFold(v string) predicate.Agent {
	return predicate.Agent(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldName), v))
	})
}

// PhoneEQ applies the EQ predicate on the "phone" field.
func PhoneEQ(v string) predicate.Agent {
	return predicate.Agent(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldPhone), v))
	})
}

// PhoneNEQ applies the NEQ predicate on the "phone" field.
func PhoneNEQ(v string) predicate.Agent {
	return predicate.Agent(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldPhone), v))
	})
}

// PhoneIn applies the In predicate on the "phone" field.
func PhoneIn(vs ...string) predicate.Agent {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Agent(func(s *sql.Selector) {
		s.Where(sql.In(s.C(FieldPhone), v...))
	})
}

// PhoneNotIn applies the NotIn predicate on the "phone" field.
func PhoneNotIn(vs ...string) predicate.Agent {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Agent(func(s *sql.Selector) {
		s.Where(sql.NotIn(s.C(FieldPhone), v...))
	})
}

// PhoneGT applies the GT predicate on the "phone" field.
func PhoneGT(v string) predicate.Agent {
	return predicate.Agent(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldPhone), v))
	})
}

// PhoneGTE applies the GTE predicate on the "phone" field.
func PhoneGTE(v string) predicate.Agent {
	return predicate.Agent(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldPhone), v))
	})
}

// PhoneLT applies the LT predicate on the "phone" field.
func PhoneLT(v string) predicate.Agent {
	return predicate.Agent(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldPhone), v))
	})
}

// PhoneLTE applies the LTE predicate on the "phone" field.
func PhoneLTE(v string) predicate.Agent {
	return predicate.Agent(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldPhone), v))
	})
}

// PhoneContains applies the Contains predicate on the "phone" field.
func PhoneContains(v string) predicate.Agent {
	return predicate.Agent(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldPhone), v))
	})
}

// PhoneHasPrefix applies the HasPrefix predicate on the "phone" field.
func PhoneHasPrefix(v string) predicate.Agent {
	return predicate.Agent(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldPhone), v))
	})
}

// PhoneHasSuffix applies the HasSuffix predicate on the "phone" field.
func PhoneHasSuffix(v string) predicate.Agent {
	return predicate.Agent(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldPhone), v))
	})
}

// PhoneEqualFold applies the EqualFold predicate on the "phone" field.
func PhoneEqualFold(v string) predicate.Agent {
	return predicate.Agent(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldPhone), v))
	})
}

// PhoneContainsFold applies the ContainsFold predicate on the "phone" field.
func PhoneContainsFold(v string) predicate.Agent {
	return predicate.Agent(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldPhone), v))
	})
}

// PasswordEQ applies the EQ predicate on the "password" field.
func PasswordEQ(v string) predicate.Agent {
	return predicate.Agent(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldPassword), v))
	})
}

// PasswordNEQ applies the NEQ predicate on the "password" field.
func PasswordNEQ(v string) predicate.Agent {
	return predicate.Agent(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldPassword), v))
	})
}

// PasswordIn applies the In predicate on the "password" field.
func PasswordIn(vs ...string) predicate.Agent {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Agent(func(s *sql.Selector) {
		s.Where(sql.In(s.C(FieldPassword), v...))
	})
}

// PasswordNotIn applies the NotIn predicate on the "password" field.
func PasswordNotIn(vs ...string) predicate.Agent {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Agent(func(s *sql.Selector) {
		s.Where(sql.NotIn(s.C(FieldPassword), v...))
	})
}

// PasswordGT applies the GT predicate on the "password" field.
func PasswordGT(v string) predicate.Agent {
	return predicate.Agent(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldPassword), v))
	})
}

// PasswordGTE applies the GTE predicate on the "password" field.
func PasswordGTE(v string) predicate.Agent {
	return predicate.Agent(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldPassword), v))
	})
}

// PasswordLT applies the LT predicate on the "password" field.
func PasswordLT(v string) predicate.Agent {
	return predicate.Agent(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldPassword), v))
	})
}

// PasswordLTE applies the LTE predicate on the "password" field.
func PasswordLTE(v string) predicate.Agent {
	return predicate.Agent(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldPassword), v))
	})
}

// PasswordContains applies the Contains predicate on the "password" field.
func PasswordContains(v string) predicate.Agent {
	return predicate.Agent(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldPassword), v))
	})
}

// PasswordHasPrefix applies the HasPrefix predicate on the "password" field.
func PasswordHasPrefix(v string) predicate.Agent {
	return predicate.Agent(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldPassword), v))
	})
}

// PasswordHasSuffix applies the HasSuffix predicate on the "password" field.
func PasswordHasSuffix(v string) predicate.Agent {
	return predicate.Agent(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldPassword), v))
	})
}

// PasswordEqualFold applies the EqualFold predicate on the "password" field.
func PasswordEqualFold(v string) predicate.Agent {
	return predicate.Agent(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldPassword), v))
	})
}

// PasswordContainsFold applies the ContainsFold predicate on the "password" field.
func PasswordContainsFold(v string) predicate.Agent {
	return predicate.Agent(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldPassword), v))
	})
}

// HasEnterprise applies the HasEdge predicate on the "enterprise" edge.
func HasEnterprise() predicate.Agent {
	return predicate.Agent(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(EnterpriseTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, EnterpriseTable, EnterpriseColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasEnterpriseWith applies the HasEdge predicate on the "enterprise" edge with a given conditions (other predicates).
func HasEnterpriseWith(preds ...predicate.Enterprise) predicate.Agent {
	return predicate.Agent(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(EnterpriseInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, EnterpriseTable, EnterpriseColumn),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.Agent) predicate.Agent {
	return predicate.Agent(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for _, p := range predicates {
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.Agent) predicate.Agent {
	return predicate.Agent(func(s *sql.Selector) {
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
func Not(p predicate.Agent) predicate.Agent {
	return predicate.Agent(func(s *sql.Selector) {
		p(s.Not())
	})
}