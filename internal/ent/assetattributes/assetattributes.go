// Code generated by ent, DO NOT EDIT.

package assetattributes

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
)

const (
	// Label holds the string label denoting the assetattributes type in the database.
	Label = "asset_attributes"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldCreatedAt holds the string denoting the created_at field in the database.
	FieldCreatedAt = "created_at"
	// FieldUpdatedAt holds the string denoting the updated_at field in the database.
	FieldUpdatedAt = "updated_at"
	// FieldAssetType holds the string denoting the asset_type field in the database.
	FieldAssetType = "asset_type"
	// FieldName holds the string denoting the name field in the database.
	FieldName = "name"
	// FieldKey holds the string denoting the key field in the database.
	FieldKey = "key"
	// EdgeValues holds the string denoting the values edge name in mutations.
	EdgeValues = "values"
	// Table holds the table name of the assetattributes in the database.
	Table = "asset_attributes"
	// ValuesTable is the table that holds the values relation/edge.
	ValuesTable = "asset_attribute_values"
	// ValuesInverseTable is the table name for the AssetAttributeValues entity.
	// It exists in this package in order to avoid circular dependency with the "assetattributevalues" package.
	ValuesInverseTable = "asset_attribute_values"
	// ValuesColumn is the table column denoting the values relation/edge.
	ValuesColumn = "attribute_id"
)

// Columns holds all SQL columns for assetattributes fields.
var Columns = []string{
	FieldID,
	FieldCreatedAt,
	FieldUpdatedAt,
	FieldAssetType,
	FieldName,
	FieldKey,
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

var (
	// DefaultCreatedAt holds the default value on creation for the "created_at" field.
	DefaultCreatedAt func() time.Time
	// DefaultUpdatedAt holds the default value on creation for the "updated_at" field.
	DefaultUpdatedAt func() time.Time
	// UpdateDefaultUpdatedAt holds the default value on update for the "updated_at" field.
	UpdateDefaultUpdatedAt func() time.Time
)

// OrderOption defines the ordering options for the AssetAttributes queries.
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

// ByAssetType orders the results by the asset_type field.
func ByAssetType(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldAssetType, opts...).ToFunc()
}

// ByName orders the results by the name field.
func ByName(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldName, opts...).ToFunc()
}

// ByKey orders the results by the key field.
func ByKey(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldKey, opts...).ToFunc()
}

// ByValuesCount orders the results by values count.
func ByValuesCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newValuesStep(), opts...)
	}
}

// ByValues orders the results by values terms.
func ByValues(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newValuesStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}
func newValuesStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(ValuesInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, ValuesTable, ValuesColumn),
	)
}