// Code generated by ent, DO NOT EDIT.

package city

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/auroraride/aurservd/internal/ent/predicate"
)

// ID filters vertices based on their ID field.
func ID(id uint64) predicate.City {
	return predicate.City(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id uint64) predicate.City {
	return predicate.City(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id uint64) predicate.City {
	return predicate.City(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...uint64) predicate.City {
	return predicate.City(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...uint64) predicate.City {
	return predicate.City(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id uint64) predicate.City {
	return predicate.City(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id uint64) predicate.City {
	return predicate.City(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id uint64) predicate.City {
	return predicate.City(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id uint64) predicate.City {
	return predicate.City(sql.FieldLTE(FieldID, id))
}

// CreatedAt applies equality check predicate on the "created_at" field. It's identical to CreatedAtEQ.
func CreatedAt(v time.Time) predicate.City {
	return predicate.City(sql.FieldEQ(FieldCreatedAt, v))
}

// UpdatedAt applies equality check predicate on the "updated_at" field. It's identical to UpdatedAtEQ.
func UpdatedAt(v time.Time) predicate.City {
	return predicate.City(sql.FieldEQ(FieldUpdatedAt, v))
}

// DeletedAt applies equality check predicate on the "deleted_at" field. It's identical to DeletedAtEQ.
func DeletedAt(v time.Time) predicate.City {
	return predicate.City(sql.FieldEQ(FieldDeletedAt, v))
}

// Remark applies equality check predicate on the "remark" field. It's identical to RemarkEQ.
func Remark(v string) predicate.City {
	return predicate.City(sql.FieldEQ(FieldRemark, v))
}

// Open applies equality check predicate on the "open" field. It's identical to OpenEQ.
func Open(v bool) predicate.City {
	return predicate.City(sql.FieldEQ(FieldOpen, v))
}

// Name applies equality check predicate on the "name" field. It's identical to NameEQ.
func Name(v string) predicate.City {
	return predicate.City(sql.FieldEQ(FieldName, v))
}

// Code applies equality check predicate on the "code" field. It's identical to CodeEQ.
func Code(v string) predicate.City {
	return predicate.City(sql.FieldEQ(FieldCode, v))
}

// ParentID applies equality check predicate on the "parent_id" field. It's identical to ParentIDEQ.
func ParentID(v uint64) predicate.City {
	return predicate.City(sql.FieldEQ(FieldParentID, v))
}

// Lng applies equality check predicate on the "lng" field. It's identical to LngEQ.
func Lng(v float64) predicate.City {
	return predicate.City(sql.FieldEQ(FieldLng, v))
}

// Lat applies equality check predicate on the "lat" field. It's identical to LatEQ.
func Lat(v float64) predicate.City {
	return predicate.City(sql.FieldEQ(FieldLat, v))
}

// CreatedAtEQ applies the EQ predicate on the "created_at" field.
func CreatedAtEQ(v time.Time) predicate.City {
	return predicate.City(sql.FieldEQ(FieldCreatedAt, v))
}

// CreatedAtNEQ applies the NEQ predicate on the "created_at" field.
func CreatedAtNEQ(v time.Time) predicate.City {
	return predicate.City(sql.FieldNEQ(FieldCreatedAt, v))
}

// CreatedAtIn applies the In predicate on the "created_at" field.
func CreatedAtIn(vs ...time.Time) predicate.City {
	return predicate.City(sql.FieldIn(FieldCreatedAt, vs...))
}

// CreatedAtNotIn applies the NotIn predicate on the "created_at" field.
func CreatedAtNotIn(vs ...time.Time) predicate.City {
	return predicate.City(sql.FieldNotIn(FieldCreatedAt, vs...))
}

// CreatedAtGT applies the GT predicate on the "created_at" field.
func CreatedAtGT(v time.Time) predicate.City {
	return predicate.City(sql.FieldGT(FieldCreatedAt, v))
}

// CreatedAtGTE applies the GTE predicate on the "created_at" field.
func CreatedAtGTE(v time.Time) predicate.City {
	return predicate.City(sql.FieldGTE(FieldCreatedAt, v))
}

// CreatedAtLT applies the LT predicate on the "created_at" field.
func CreatedAtLT(v time.Time) predicate.City {
	return predicate.City(sql.FieldLT(FieldCreatedAt, v))
}

// CreatedAtLTE applies the LTE predicate on the "created_at" field.
func CreatedAtLTE(v time.Time) predicate.City {
	return predicate.City(sql.FieldLTE(FieldCreatedAt, v))
}

// UpdatedAtEQ applies the EQ predicate on the "updated_at" field.
func UpdatedAtEQ(v time.Time) predicate.City {
	return predicate.City(sql.FieldEQ(FieldUpdatedAt, v))
}

// UpdatedAtNEQ applies the NEQ predicate on the "updated_at" field.
func UpdatedAtNEQ(v time.Time) predicate.City {
	return predicate.City(sql.FieldNEQ(FieldUpdatedAt, v))
}

// UpdatedAtIn applies the In predicate on the "updated_at" field.
func UpdatedAtIn(vs ...time.Time) predicate.City {
	return predicate.City(sql.FieldIn(FieldUpdatedAt, vs...))
}

// UpdatedAtNotIn applies the NotIn predicate on the "updated_at" field.
func UpdatedAtNotIn(vs ...time.Time) predicate.City {
	return predicate.City(sql.FieldNotIn(FieldUpdatedAt, vs...))
}

// UpdatedAtGT applies the GT predicate on the "updated_at" field.
func UpdatedAtGT(v time.Time) predicate.City {
	return predicate.City(sql.FieldGT(FieldUpdatedAt, v))
}

// UpdatedAtGTE applies the GTE predicate on the "updated_at" field.
func UpdatedAtGTE(v time.Time) predicate.City {
	return predicate.City(sql.FieldGTE(FieldUpdatedAt, v))
}

// UpdatedAtLT applies the LT predicate on the "updated_at" field.
func UpdatedAtLT(v time.Time) predicate.City {
	return predicate.City(sql.FieldLT(FieldUpdatedAt, v))
}

// UpdatedAtLTE applies the LTE predicate on the "updated_at" field.
func UpdatedAtLTE(v time.Time) predicate.City {
	return predicate.City(sql.FieldLTE(FieldUpdatedAt, v))
}

// DeletedAtEQ applies the EQ predicate on the "deleted_at" field.
func DeletedAtEQ(v time.Time) predicate.City {
	return predicate.City(sql.FieldEQ(FieldDeletedAt, v))
}

// DeletedAtNEQ applies the NEQ predicate on the "deleted_at" field.
func DeletedAtNEQ(v time.Time) predicate.City {
	return predicate.City(sql.FieldNEQ(FieldDeletedAt, v))
}

// DeletedAtIn applies the In predicate on the "deleted_at" field.
func DeletedAtIn(vs ...time.Time) predicate.City {
	return predicate.City(sql.FieldIn(FieldDeletedAt, vs...))
}

// DeletedAtNotIn applies the NotIn predicate on the "deleted_at" field.
func DeletedAtNotIn(vs ...time.Time) predicate.City {
	return predicate.City(sql.FieldNotIn(FieldDeletedAt, vs...))
}

// DeletedAtGT applies the GT predicate on the "deleted_at" field.
func DeletedAtGT(v time.Time) predicate.City {
	return predicate.City(sql.FieldGT(FieldDeletedAt, v))
}

// DeletedAtGTE applies the GTE predicate on the "deleted_at" field.
func DeletedAtGTE(v time.Time) predicate.City {
	return predicate.City(sql.FieldGTE(FieldDeletedAt, v))
}

// DeletedAtLT applies the LT predicate on the "deleted_at" field.
func DeletedAtLT(v time.Time) predicate.City {
	return predicate.City(sql.FieldLT(FieldDeletedAt, v))
}

// DeletedAtLTE applies the LTE predicate on the "deleted_at" field.
func DeletedAtLTE(v time.Time) predicate.City {
	return predicate.City(sql.FieldLTE(FieldDeletedAt, v))
}

// DeletedAtIsNil applies the IsNil predicate on the "deleted_at" field.
func DeletedAtIsNil() predicate.City {
	return predicate.City(sql.FieldIsNull(FieldDeletedAt))
}

// DeletedAtNotNil applies the NotNil predicate on the "deleted_at" field.
func DeletedAtNotNil() predicate.City {
	return predicate.City(sql.FieldNotNull(FieldDeletedAt))
}

// CreatorIsNil applies the IsNil predicate on the "creator" field.
func CreatorIsNil() predicate.City {
	return predicate.City(sql.FieldIsNull(FieldCreator))
}

// CreatorNotNil applies the NotNil predicate on the "creator" field.
func CreatorNotNil() predicate.City {
	return predicate.City(sql.FieldNotNull(FieldCreator))
}

// LastModifierIsNil applies the IsNil predicate on the "last_modifier" field.
func LastModifierIsNil() predicate.City {
	return predicate.City(sql.FieldIsNull(FieldLastModifier))
}

// LastModifierNotNil applies the NotNil predicate on the "last_modifier" field.
func LastModifierNotNil() predicate.City {
	return predicate.City(sql.FieldNotNull(FieldLastModifier))
}

// RemarkEQ applies the EQ predicate on the "remark" field.
func RemarkEQ(v string) predicate.City {
	return predicate.City(sql.FieldEQ(FieldRemark, v))
}

// RemarkNEQ applies the NEQ predicate on the "remark" field.
func RemarkNEQ(v string) predicate.City {
	return predicate.City(sql.FieldNEQ(FieldRemark, v))
}

// RemarkIn applies the In predicate on the "remark" field.
func RemarkIn(vs ...string) predicate.City {
	return predicate.City(sql.FieldIn(FieldRemark, vs...))
}

// RemarkNotIn applies the NotIn predicate on the "remark" field.
func RemarkNotIn(vs ...string) predicate.City {
	return predicate.City(sql.FieldNotIn(FieldRemark, vs...))
}

// RemarkGT applies the GT predicate on the "remark" field.
func RemarkGT(v string) predicate.City {
	return predicate.City(sql.FieldGT(FieldRemark, v))
}

// RemarkGTE applies the GTE predicate on the "remark" field.
func RemarkGTE(v string) predicate.City {
	return predicate.City(sql.FieldGTE(FieldRemark, v))
}

// RemarkLT applies the LT predicate on the "remark" field.
func RemarkLT(v string) predicate.City {
	return predicate.City(sql.FieldLT(FieldRemark, v))
}

// RemarkLTE applies the LTE predicate on the "remark" field.
func RemarkLTE(v string) predicate.City {
	return predicate.City(sql.FieldLTE(FieldRemark, v))
}

// RemarkContains applies the Contains predicate on the "remark" field.
func RemarkContains(v string) predicate.City {
	return predicate.City(sql.FieldContains(FieldRemark, v))
}

// RemarkHasPrefix applies the HasPrefix predicate on the "remark" field.
func RemarkHasPrefix(v string) predicate.City {
	return predicate.City(sql.FieldHasPrefix(FieldRemark, v))
}

// RemarkHasSuffix applies the HasSuffix predicate on the "remark" field.
func RemarkHasSuffix(v string) predicate.City {
	return predicate.City(sql.FieldHasSuffix(FieldRemark, v))
}

// RemarkIsNil applies the IsNil predicate on the "remark" field.
func RemarkIsNil() predicate.City {
	return predicate.City(sql.FieldIsNull(FieldRemark))
}

// RemarkNotNil applies the NotNil predicate on the "remark" field.
func RemarkNotNil() predicate.City {
	return predicate.City(sql.FieldNotNull(FieldRemark))
}

// RemarkEqualFold applies the EqualFold predicate on the "remark" field.
func RemarkEqualFold(v string) predicate.City {
	return predicate.City(sql.FieldEqualFold(FieldRemark, v))
}

// RemarkContainsFold applies the ContainsFold predicate on the "remark" field.
func RemarkContainsFold(v string) predicate.City {
	return predicate.City(sql.FieldContainsFold(FieldRemark, v))
}

// OpenEQ applies the EQ predicate on the "open" field.
func OpenEQ(v bool) predicate.City {
	return predicate.City(sql.FieldEQ(FieldOpen, v))
}

// OpenNEQ applies the NEQ predicate on the "open" field.
func OpenNEQ(v bool) predicate.City {
	return predicate.City(sql.FieldNEQ(FieldOpen, v))
}

// OpenIsNil applies the IsNil predicate on the "open" field.
func OpenIsNil() predicate.City {
	return predicate.City(sql.FieldIsNull(FieldOpen))
}

// OpenNotNil applies the NotNil predicate on the "open" field.
func OpenNotNil() predicate.City {
	return predicate.City(sql.FieldNotNull(FieldOpen))
}

// NameEQ applies the EQ predicate on the "name" field.
func NameEQ(v string) predicate.City {
	return predicate.City(sql.FieldEQ(FieldName, v))
}

// NameNEQ applies the NEQ predicate on the "name" field.
func NameNEQ(v string) predicate.City {
	return predicate.City(sql.FieldNEQ(FieldName, v))
}

// NameIn applies the In predicate on the "name" field.
func NameIn(vs ...string) predicate.City {
	return predicate.City(sql.FieldIn(FieldName, vs...))
}

// NameNotIn applies the NotIn predicate on the "name" field.
func NameNotIn(vs ...string) predicate.City {
	return predicate.City(sql.FieldNotIn(FieldName, vs...))
}

// NameGT applies the GT predicate on the "name" field.
func NameGT(v string) predicate.City {
	return predicate.City(sql.FieldGT(FieldName, v))
}

// NameGTE applies the GTE predicate on the "name" field.
func NameGTE(v string) predicate.City {
	return predicate.City(sql.FieldGTE(FieldName, v))
}

// NameLT applies the LT predicate on the "name" field.
func NameLT(v string) predicate.City {
	return predicate.City(sql.FieldLT(FieldName, v))
}

// NameLTE applies the LTE predicate on the "name" field.
func NameLTE(v string) predicate.City {
	return predicate.City(sql.FieldLTE(FieldName, v))
}

// NameContains applies the Contains predicate on the "name" field.
func NameContains(v string) predicate.City {
	return predicate.City(sql.FieldContains(FieldName, v))
}

// NameHasPrefix applies the HasPrefix predicate on the "name" field.
func NameHasPrefix(v string) predicate.City {
	return predicate.City(sql.FieldHasPrefix(FieldName, v))
}

// NameHasSuffix applies the HasSuffix predicate on the "name" field.
func NameHasSuffix(v string) predicate.City {
	return predicate.City(sql.FieldHasSuffix(FieldName, v))
}

// NameEqualFold applies the EqualFold predicate on the "name" field.
func NameEqualFold(v string) predicate.City {
	return predicate.City(sql.FieldEqualFold(FieldName, v))
}

// NameContainsFold applies the ContainsFold predicate on the "name" field.
func NameContainsFold(v string) predicate.City {
	return predicate.City(sql.FieldContainsFold(FieldName, v))
}

// CodeEQ applies the EQ predicate on the "code" field.
func CodeEQ(v string) predicate.City {
	return predicate.City(sql.FieldEQ(FieldCode, v))
}

// CodeNEQ applies the NEQ predicate on the "code" field.
func CodeNEQ(v string) predicate.City {
	return predicate.City(sql.FieldNEQ(FieldCode, v))
}

// CodeIn applies the In predicate on the "code" field.
func CodeIn(vs ...string) predicate.City {
	return predicate.City(sql.FieldIn(FieldCode, vs...))
}

// CodeNotIn applies the NotIn predicate on the "code" field.
func CodeNotIn(vs ...string) predicate.City {
	return predicate.City(sql.FieldNotIn(FieldCode, vs...))
}

// CodeGT applies the GT predicate on the "code" field.
func CodeGT(v string) predicate.City {
	return predicate.City(sql.FieldGT(FieldCode, v))
}

// CodeGTE applies the GTE predicate on the "code" field.
func CodeGTE(v string) predicate.City {
	return predicate.City(sql.FieldGTE(FieldCode, v))
}

// CodeLT applies the LT predicate on the "code" field.
func CodeLT(v string) predicate.City {
	return predicate.City(sql.FieldLT(FieldCode, v))
}

// CodeLTE applies the LTE predicate on the "code" field.
func CodeLTE(v string) predicate.City {
	return predicate.City(sql.FieldLTE(FieldCode, v))
}

// CodeContains applies the Contains predicate on the "code" field.
func CodeContains(v string) predicate.City {
	return predicate.City(sql.FieldContains(FieldCode, v))
}

// CodeHasPrefix applies the HasPrefix predicate on the "code" field.
func CodeHasPrefix(v string) predicate.City {
	return predicate.City(sql.FieldHasPrefix(FieldCode, v))
}

// CodeHasSuffix applies the HasSuffix predicate on the "code" field.
func CodeHasSuffix(v string) predicate.City {
	return predicate.City(sql.FieldHasSuffix(FieldCode, v))
}

// CodeEqualFold applies the EqualFold predicate on the "code" field.
func CodeEqualFold(v string) predicate.City {
	return predicate.City(sql.FieldEqualFold(FieldCode, v))
}

// CodeContainsFold applies the ContainsFold predicate on the "code" field.
func CodeContainsFold(v string) predicate.City {
	return predicate.City(sql.FieldContainsFold(FieldCode, v))
}

// ParentIDEQ applies the EQ predicate on the "parent_id" field.
func ParentIDEQ(v uint64) predicate.City {
	return predicate.City(sql.FieldEQ(FieldParentID, v))
}

// ParentIDNEQ applies the NEQ predicate on the "parent_id" field.
func ParentIDNEQ(v uint64) predicate.City {
	return predicate.City(sql.FieldNEQ(FieldParentID, v))
}

// ParentIDIn applies the In predicate on the "parent_id" field.
func ParentIDIn(vs ...uint64) predicate.City {
	return predicate.City(sql.FieldIn(FieldParentID, vs...))
}

// ParentIDNotIn applies the NotIn predicate on the "parent_id" field.
func ParentIDNotIn(vs ...uint64) predicate.City {
	return predicate.City(sql.FieldNotIn(FieldParentID, vs...))
}

// ParentIDIsNil applies the IsNil predicate on the "parent_id" field.
func ParentIDIsNil() predicate.City {
	return predicate.City(sql.FieldIsNull(FieldParentID))
}

// ParentIDNotNil applies the NotNil predicate on the "parent_id" field.
func ParentIDNotNil() predicate.City {
	return predicate.City(sql.FieldNotNull(FieldParentID))
}

// LngEQ applies the EQ predicate on the "lng" field.
func LngEQ(v float64) predicate.City {
	return predicate.City(sql.FieldEQ(FieldLng, v))
}

// LngNEQ applies the NEQ predicate on the "lng" field.
func LngNEQ(v float64) predicate.City {
	return predicate.City(sql.FieldNEQ(FieldLng, v))
}

// LngIn applies the In predicate on the "lng" field.
func LngIn(vs ...float64) predicate.City {
	return predicate.City(sql.FieldIn(FieldLng, vs...))
}

// LngNotIn applies the NotIn predicate on the "lng" field.
func LngNotIn(vs ...float64) predicate.City {
	return predicate.City(sql.FieldNotIn(FieldLng, vs...))
}

// LngGT applies the GT predicate on the "lng" field.
func LngGT(v float64) predicate.City {
	return predicate.City(sql.FieldGT(FieldLng, v))
}

// LngGTE applies the GTE predicate on the "lng" field.
func LngGTE(v float64) predicate.City {
	return predicate.City(sql.FieldGTE(FieldLng, v))
}

// LngLT applies the LT predicate on the "lng" field.
func LngLT(v float64) predicate.City {
	return predicate.City(sql.FieldLT(FieldLng, v))
}

// LngLTE applies the LTE predicate on the "lng" field.
func LngLTE(v float64) predicate.City {
	return predicate.City(sql.FieldLTE(FieldLng, v))
}

// LngIsNil applies the IsNil predicate on the "lng" field.
func LngIsNil() predicate.City {
	return predicate.City(sql.FieldIsNull(FieldLng))
}

// LngNotNil applies the NotNil predicate on the "lng" field.
func LngNotNil() predicate.City {
	return predicate.City(sql.FieldNotNull(FieldLng))
}

// LatEQ applies the EQ predicate on the "lat" field.
func LatEQ(v float64) predicate.City {
	return predicate.City(sql.FieldEQ(FieldLat, v))
}

// LatNEQ applies the NEQ predicate on the "lat" field.
func LatNEQ(v float64) predicate.City {
	return predicate.City(sql.FieldNEQ(FieldLat, v))
}

// LatIn applies the In predicate on the "lat" field.
func LatIn(vs ...float64) predicate.City {
	return predicate.City(sql.FieldIn(FieldLat, vs...))
}

// LatNotIn applies the NotIn predicate on the "lat" field.
func LatNotIn(vs ...float64) predicate.City {
	return predicate.City(sql.FieldNotIn(FieldLat, vs...))
}

// LatGT applies the GT predicate on the "lat" field.
func LatGT(v float64) predicate.City {
	return predicate.City(sql.FieldGT(FieldLat, v))
}

// LatGTE applies the GTE predicate on the "lat" field.
func LatGTE(v float64) predicate.City {
	return predicate.City(sql.FieldGTE(FieldLat, v))
}

// LatLT applies the LT predicate on the "lat" field.
func LatLT(v float64) predicate.City {
	return predicate.City(sql.FieldLT(FieldLat, v))
}

// LatLTE applies the LTE predicate on the "lat" field.
func LatLTE(v float64) predicate.City {
	return predicate.City(sql.FieldLTE(FieldLat, v))
}

// LatIsNil applies the IsNil predicate on the "lat" field.
func LatIsNil() predicate.City {
	return predicate.City(sql.FieldIsNull(FieldLat))
}

// LatNotNil applies the NotNil predicate on the "lat" field.
func LatNotNil() predicate.City {
	return predicate.City(sql.FieldNotNull(FieldLat))
}

// HasParent applies the HasEdge predicate on the "parent" edge.
func HasParent() predicate.City {
	return predicate.City(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, ParentTable, ParentColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasParentWith applies the HasEdge predicate on the "parent" edge with a given conditions (other predicates).
func HasParentWith(preds ...predicate.City) predicate.City {
	return predicate.City(func(s *sql.Selector) {
		step := newParentStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasChildren applies the HasEdge predicate on the "children" edge.
func HasChildren() predicate.City {
	return predicate.City(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, ChildrenTable, ChildrenColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasChildrenWith applies the HasEdge predicate on the "children" edge with a given conditions (other predicates).
func HasChildrenWith(preds ...predicate.City) predicate.City {
	return predicate.City(func(s *sql.Selector) {
		step := newChildrenStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasPlans applies the HasEdge predicate on the "plans" edge.
func HasPlans() predicate.City {
	return predicate.City(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2M, true, PlansTable, PlansPrimaryKey...),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasPlansWith applies the HasEdge predicate on the "plans" edge with a given conditions (other predicates).
func HasPlansWith(preds ...predicate.Plan) predicate.City {
	return predicate.City(func(s *sql.Selector) {
		step := newPlansStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasMaintainers applies the HasEdge predicate on the "maintainers" edge.
func HasMaintainers() predicate.City {
	return predicate.City(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2M, false, MaintainersTable, MaintainersPrimaryKey...),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasMaintainersWith applies the HasEdge predicate on the "maintainers" edge with a given conditions (other predicates).
func HasMaintainersWith(preds ...predicate.Maintainer) predicate.City {
	return predicate.City(func(s *sql.Selector) {
		step := newMaintainersStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.City) predicate.City {
	return predicate.City(sql.AndPredicates(predicates...))
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.City) predicate.City {
	return predicate.City(sql.OrPredicates(predicates...))
}

// Not applies the not operator on the given predicate.
func Not(p predicate.City) predicate.City {
	return predicate.City(sql.NotPredicates(p))
}
