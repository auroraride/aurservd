// Code generated by ent, DO NOT EDIT.

package ebikebrandattribute

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
)

const (
	// Label holds the string label denoting the ebikebrandattribute type in the database.
	Label = "ebike_brand_attribute"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldCreatedAt holds the string denoting the created_at field in the database.
	FieldCreatedAt = "created_at"
	// FieldUpdatedAt holds the string denoting the updated_at field in the database.
	FieldUpdatedAt = "updated_at"
	// FieldDeletedAt holds the string denoting the deleted_at field in the database.
	FieldDeletedAt = "deleted_at"
	// FieldCreator holds the string denoting the creator field in the database.
	FieldCreator = "creator"
	// FieldLastModifier holds the string denoting the last_modifier field in the database.
	FieldLastModifier = "last_modifier"
	// FieldRemark holds the string denoting the remark field in the database.
	FieldRemark = "remark"
	// FieldName holds the string denoting the name field in the database.
	FieldName = "name"
	// FieldValue holds the string denoting the value field in the database.
	FieldValue = "value"
	// FieldBrandID holds the string denoting the brand_id field in the database.
	FieldBrandID = "brand_id"
	// EdgeBrand holds the string denoting the brand edge name in mutations.
	EdgeBrand = "brand"
	// Table holds the table name of the ebikebrandattribute in the database.
	Table = "ebike_brand_attribute"
	// BrandTable is the table that holds the brand relation/edge.
	BrandTable = "ebike_brand_attribute"
	// BrandInverseTable is the table name for the EbikeBrand entity.
	// It exists in this package in order to avoid circular dependency with the "ebikebrand" package.
	BrandInverseTable = "ebike_brand"
	// BrandColumn is the table column denoting the brand relation/edge.
	BrandColumn = "brand_id"
)

// Columns holds all SQL columns for ebikebrandattribute fields.
var Columns = []string{
	FieldID,
	FieldCreatedAt,
	FieldUpdatedAt,
	FieldDeletedAt,
	FieldCreator,
	FieldLastModifier,
	FieldRemark,
	FieldName,
	FieldValue,
	FieldBrandID,
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	return false
}

// Note that the variables below are initialized by the runtime
// package on the initialization of the application. Therefore,
// it should be imported in the main as follows:
//
//	import _ "github.com/auroraride/aurservd/internal/ent/runtime"
var (
	Hooks [1]ent.Hook
	// DefaultCreatedAt holds the default value on creation for the "created_at" field.
	DefaultCreatedAt func() time.Time
	// DefaultUpdatedAt holds the default value on creation for the "updated_at" field.
	DefaultUpdatedAt func() time.Time
	// UpdateDefaultUpdatedAt holds the default value on update for the "updated_at" field.
	UpdateDefaultUpdatedAt func() time.Time
)

// OrderOption defines the ordering options for the EbikeBrandAttribute queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByCreatedAt orders the results by the created_at field.
func ByCreatedAt(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldCreatedAt, opts...).ToFunc()
}

// ByUpdatedAt orders the results by the updated_at field.
func ByUpdatedAt(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldUpdatedAt, opts...).ToFunc()
}

// ByDeletedAt orders the results by the deleted_at field.
func ByDeletedAt(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldDeletedAt, opts...).ToFunc()
}

// ByRemark orders the results by the remark field.
func ByRemark(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldRemark, opts...).ToFunc()
}

// ByName orders the results by the name field.
func ByName(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldName, opts...).ToFunc()
}

// ByValue orders the results by the value field.
func ByValue(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldValue, opts...).ToFunc()
}

// ByBrandID orders the results by the brand_id field.
func ByBrandID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldBrandID, opts...).ToFunc()
}

// ByBrandField orders the results by brand field.
func ByBrandField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newBrandStep(), sql.OrderByField(field, opts...))
	}
}
func newBrandStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(BrandInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, BrandTable, BrandColumn),
	)
}