// Code generated by ent, DO NOT EDIT.

package promotionmembercommission

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/auroraride/aurservd/internal/ent/predicate"
)

// ID filters vertices based on their ID field.
func ID(id uint64) predicate.PromotionMemberCommission {
	return predicate.PromotionMemberCommission(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id uint64) predicate.PromotionMemberCommission {
	return predicate.PromotionMemberCommission(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id uint64) predicate.PromotionMemberCommission {
	return predicate.PromotionMemberCommission(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...uint64) predicate.PromotionMemberCommission {
	return predicate.PromotionMemberCommission(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...uint64) predicate.PromotionMemberCommission {
	return predicate.PromotionMemberCommission(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id uint64) predicate.PromotionMemberCommission {
	return predicate.PromotionMemberCommission(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id uint64) predicate.PromotionMemberCommission {
	return predicate.PromotionMemberCommission(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id uint64) predicate.PromotionMemberCommission {
	return predicate.PromotionMemberCommission(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id uint64) predicate.PromotionMemberCommission {
	return predicate.PromotionMemberCommission(sql.FieldLTE(FieldID, id))
}

// CreatedAt applies equality check predicate on the "created_at" field. It's identical to CreatedAtEQ.
func CreatedAt(v time.Time) predicate.PromotionMemberCommission {
	return predicate.PromotionMemberCommission(sql.FieldEQ(FieldCreatedAt, v))
}

// UpdatedAt applies equality check predicate on the "updated_at" field. It's identical to UpdatedAtEQ.
func UpdatedAt(v time.Time) predicate.PromotionMemberCommission {
	return predicate.PromotionMemberCommission(sql.FieldEQ(FieldUpdatedAt, v))
}

// DeletedAt applies equality check predicate on the "deleted_at" field. It's identical to DeletedAtEQ.
func DeletedAt(v time.Time) predicate.PromotionMemberCommission {
	return predicate.PromotionMemberCommission(sql.FieldEQ(FieldDeletedAt, v))
}

// CommissionID applies equality check predicate on the "commission_id" field. It's identical to CommissionIDEQ.
func CommissionID(v uint64) predicate.PromotionMemberCommission {
	return predicate.PromotionMemberCommission(sql.FieldEQ(FieldCommissionID, v))
}

// MemberID applies equality check predicate on the "member_id" field. It's identical to MemberIDEQ.
func MemberID(v uint64) predicate.PromotionMemberCommission {
	return predicate.PromotionMemberCommission(sql.FieldEQ(FieldMemberID, v))
}

// CreatedAtEQ applies the EQ predicate on the "created_at" field.
func CreatedAtEQ(v time.Time) predicate.PromotionMemberCommission {
	return predicate.PromotionMemberCommission(sql.FieldEQ(FieldCreatedAt, v))
}

// CreatedAtNEQ applies the NEQ predicate on the "created_at" field.
func CreatedAtNEQ(v time.Time) predicate.PromotionMemberCommission {
	return predicate.PromotionMemberCommission(sql.FieldNEQ(FieldCreatedAt, v))
}

// CreatedAtIn applies the In predicate on the "created_at" field.
func CreatedAtIn(vs ...time.Time) predicate.PromotionMemberCommission {
	return predicate.PromotionMemberCommission(sql.FieldIn(FieldCreatedAt, vs...))
}

// CreatedAtNotIn applies the NotIn predicate on the "created_at" field.
func CreatedAtNotIn(vs ...time.Time) predicate.PromotionMemberCommission {
	return predicate.PromotionMemberCommission(sql.FieldNotIn(FieldCreatedAt, vs...))
}

// CreatedAtGT applies the GT predicate on the "created_at" field.
func CreatedAtGT(v time.Time) predicate.PromotionMemberCommission {
	return predicate.PromotionMemberCommission(sql.FieldGT(FieldCreatedAt, v))
}

// CreatedAtGTE applies the GTE predicate on the "created_at" field.
func CreatedAtGTE(v time.Time) predicate.PromotionMemberCommission {
	return predicate.PromotionMemberCommission(sql.FieldGTE(FieldCreatedAt, v))
}

// CreatedAtLT applies the LT predicate on the "created_at" field.
func CreatedAtLT(v time.Time) predicate.PromotionMemberCommission {
	return predicate.PromotionMemberCommission(sql.FieldLT(FieldCreatedAt, v))
}

// CreatedAtLTE applies the LTE predicate on the "created_at" field.
func CreatedAtLTE(v time.Time) predicate.PromotionMemberCommission {
	return predicate.PromotionMemberCommission(sql.FieldLTE(FieldCreatedAt, v))
}

// UpdatedAtEQ applies the EQ predicate on the "updated_at" field.
func UpdatedAtEQ(v time.Time) predicate.PromotionMemberCommission {
	return predicate.PromotionMemberCommission(sql.FieldEQ(FieldUpdatedAt, v))
}

// UpdatedAtNEQ applies the NEQ predicate on the "updated_at" field.
func UpdatedAtNEQ(v time.Time) predicate.PromotionMemberCommission {
	return predicate.PromotionMemberCommission(sql.FieldNEQ(FieldUpdatedAt, v))
}

// UpdatedAtIn applies the In predicate on the "updated_at" field.
func UpdatedAtIn(vs ...time.Time) predicate.PromotionMemberCommission {
	return predicate.PromotionMemberCommission(sql.FieldIn(FieldUpdatedAt, vs...))
}

// UpdatedAtNotIn applies the NotIn predicate on the "updated_at" field.
func UpdatedAtNotIn(vs ...time.Time) predicate.PromotionMemberCommission {
	return predicate.PromotionMemberCommission(sql.FieldNotIn(FieldUpdatedAt, vs...))
}

// UpdatedAtGT applies the GT predicate on the "updated_at" field.
func UpdatedAtGT(v time.Time) predicate.PromotionMemberCommission {
	return predicate.PromotionMemberCommission(sql.FieldGT(FieldUpdatedAt, v))
}

// UpdatedAtGTE applies the GTE predicate on the "updated_at" field.
func UpdatedAtGTE(v time.Time) predicate.PromotionMemberCommission {
	return predicate.PromotionMemberCommission(sql.FieldGTE(FieldUpdatedAt, v))
}

// UpdatedAtLT applies the LT predicate on the "updated_at" field.
func UpdatedAtLT(v time.Time) predicate.PromotionMemberCommission {
	return predicate.PromotionMemberCommission(sql.FieldLT(FieldUpdatedAt, v))
}

// UpdatedAtLTE applies the LTE predicate on the "updated_at" field.
func UpdatedAtLTE(v time.Time) predicate.PromotionMemberCommission {
	return predicate.PromotionMemberCommission(sql.FieldLTE(FieldUpdatedAt, v))
}

// DeletedAtEQ applies the EQ predicate on the "deleted_at" field.
func DeletedAtEQ(v time.Time) predicate.PromotionMemberCommission {
	return predicate.PromotionMemberCommission(sql.FieldEQ(FieldDeletedAt, v))
}

// DeletedAtNEQ applies the NEQ predicate on the "deleted_at" field.
func DeletedAtNEQ(v time.Time) predicate.PromotionMemberCommission {
	return predicate.PromotionMemberCommission(sql.FieldNEQ(FieldDeletedAt, v))
}

// DeletedAtIn applies the In predicate on the "deleted_at" field.
func DeletedAtIn(vs ...time.Time) predicate.PromotionMemberCommission {
	return predicate.PromotionMemberCommission(sql.FieldIn(FieldDeletedAt, vs...))
}

// DeletedAtNotIn applies the NotIn predicate on the "deleted_at" field.
func DeletedAtNotIn(vs ...time.Time) predicate.PromotionMemberCommission {
	return predicate.PromotionMemberCommission(sql.FieldNotIn(FieldDeletedAt, vs...))
}

// DeletedAtGT applies the GT predicate on the "deleted_at" field.
func DeletedAtGT(v time.Time) predicate.PromotionMemberCommission {
	return predicate.PromotionMemberCommission(sql.FieldGT(FieldDeletedAt, v))
}

// DeletedAtGTE applies the GTE predicate on the "deleted_at" field.
func DeletedAtGTE(v time.Time) predicate.PromotionMemberCommission {
	return predicate.PromotionMemberCommission(sql.FieldGTE(FieldDeletedAt, v))
}

// DeletedAtLT applies the LT predicate on the "deleted_at" field.
func DeletedAtLT(v time.Time) predicate.PromotionMemberCommission {
	return predicate.PromotionMemberCommission(sql.FieldLT(FieldDeletedAt, v))
}

// DeletedAtLTE applies the LTE predicate on the "deleted_at" field.
func DeletedAtLTE(v time.Time) predicate.PromotionMemberCommission {
	return predicate.PromotionMemberCommission(sql.FieldLTE(FieldDeletedAt, v))
}

// DeletedAtIsNil applies the IsNil predicate on the "deleted_at" field.
func DeletedAtIsNil() predicate.PromotionMemberCommission {
	return predicate.PromotionMemberCommission(sql.FieldIsNull(FieldDeletedAt))
}

// DeletedAtNotNil applies the NotNil predicate on the "deleted_at" field.
func DeletedAtNotNil() predicate.PromotionMemberCommission {
	return predicate.PromotionMemberCommission(sql.FieldNotNull(FieldDeletedAt))
}

// CommissionIDEQ applies the EQ predicate on the "commission_id" field.
func CommissionIDEQ(v uint64) predicate.PromotionMemberCommission {
	return predicate.PromotionMemberCommission(sql.FieldEQ(FieldCommissionID, v))
}

// CommissionIDNEQ applies the NEQ predicate on the "commission_id" field.
func CommissionIDNEQ(v uint64) predicate.PromotionMemberCommission {
	return predicate.PromotionMemberCommission(sql.FieldNEQ(FieldCommissionID, v))
}

// CommissionIDIn applies the In predicate on the "commission_id" field.
func CommissionIDIn(vs ...uint64) predicate.PromotionMemberCommission {
	return predicate.PromotionMemberCommission(sql.FieldIn(FieldCommissionID, vs...))
}

// CommissionIDNotIn applies the NotIn predicate on the "commission_id" field.
func CommissionIDNotIn(vs ...uint64) predicate.PromotionMemberCommission {
	return predicate.PromotionMemberCommission(sql.FieldNotIn(FieldCommissionID, vs...))
}

// MemberIDEQ applies the EQ predicate on the "member_id" field.
func MemberIDEQ(v uint64) predicate.PromotionMemberCommission {
	return predicate.PromotionMemberCommission(sql.FieldEQ(FieldMemberID, v))
}

// MemberIDNEQ applies the NEQ predicate on the "member_id" field.
func MemberIDNEQ(v uint64) predicate.PromotionMemberCommission {
	return predicate.PromotionMemberCommission(sql.FieldNEQ(FieldMemberID, v))
}

// MemberIDIn applies the In predicate on the "member_id" field.
func MemberIDIn(vs ...uint64) predicate.PromotionMemberCommission {
	return predicate.PromotionMemberCommission(sql.FieldIn(FieldMemberID, vs...))
}

// MemberIDNotIn applies the NotIn predicate on the "member_id" field.
func MemberIDNotIn(vs ...uint64) predicate.PromotionMemberCommission {
	return predicate.PromotionMemberCommission(sql.FieldNotIn(FieldMemberID, vs...))
}

// MemberIDIsNil applies the IsNil predicate on the "member_id" field.
func MemberIDIsNil() predicate.PromotionMemberCommission {
	return predicate.PromotionMemberCommission(sql.FieldIsNull(FieldMemberID))
}

// MemberIDNotNil applies the NotNil predicate on the "member_id" field.
func MemberIDNotNil() predicate.PromotionMemberCommission {
	return predicate.PromotionMemberCommission(sql.FieldNotNull(FieldMemberID))
}

// HasCommission applies the HasEdge predicate on the "commission" edge.
func HasCommission() predicate.PromotionMemberCommission {
	return predicate.PromotionMemberCommission(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, CommissionTable, CommissionColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasCommissionWith applies the HasEdge predicate on the "commission" edge with a given conditions (other predicates).
func HasCommissionWith(preds ...predicate.PromotionCommission) predicate.PromotionMemberCommission {
	return predicate.PromotionMemberCommission(func(s *sql.Selector) {
		step := newCommissionStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasMember applies the HasEdge predicate on the "member" edge.
func HasMember() predicate.PromotionMemberCommission {
	return predicate.PromotionMemberCommission(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, MemberTable, MemberColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasMemberWith applies the HasEdge predicate on the "member" edge with a given conditions (other predicates).
func HasMemberWith(preds ...predicate.PromotionMember) predicate.PromotionMemberCommission {
	return predicate.PromotionMemberCommission(func(s *sql.Selector) {
		step := newMemberStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.PromotionMemberCommission) predicate.PromotionMemberCommission {
	return predicate.PromotionMemberCommission(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for _, p := range predicates {
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.PromotionMemberCommission) predicate.PromotionMemberCommission {
	return predicate.PromotionMemberCommission(func(s *sql.Selector) {
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
func Not(p predicate.PromotionMemberCommission) predicate.PromotionMemberCommission {
	return predicate.PromotionMemberCommission(func(s *sql.Selector) {
		p(s.Not())
	})
}