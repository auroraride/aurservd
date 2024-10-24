// Code generated by ent, DO NOT EDIT.

package role

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
)

const (
	// Label holds the string label denoting the role type in the database.
	Label = "role"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldName holds the string denoting the name field in the database.
	FieldName = "name"
	// FieldPermissions holds the string denoting the permissions field in the database.
	FieldPermissions = "permissions"
	// FieldBuildin holds the string denoting the buildin field in the database.
	FieldBuildin = "buildin"
	// FieldSuper holds the string denoting the super field in the database.
	FieldSuper = "super"
	// FieldCreatedAt holds the string denoting the created_at field in the database.
	FieldCreatedAt = "created_at"
	// EdgeManagers holds the string denoting the managers edge name in mutations.
	EdgeManagers = "managers"
	// Table holds the table name of the role in the database.
	Table = "role"
	// ManagersTable is the table that holds the managers relation/edge.
	ManagersTable = "manager"
	// ManagersInverseTable is the table name for the Manager entity.
	// It exists in this package in order to avoid circular dependency with the "manager" package.
	ManagersInverseTable = "manager"
	// ManagersColumn is the table column denoting the managers relation/edge.
	ManagersColumn = "role_id"
)

// Columns holds all SQL columns for role fields.
var Columns = []string{
	FieldID,
	FieldName,
	FieldPermissions,
	FieldBuildin,
	FieldSuper,
	FieldCreatedAt,
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
	// DefaultBuildin holds the default value on creation for the "buildin" field.
	DefaultBuildin bool
	// DefaultSuper holds the default value on creation for the "super" field.
	DefaultSuper bool
	// DefaultCreatedAt holds the default value on creation for the "created_at" field.
	DefaultCreatedAt func() time.Time
)

// OrderOption defines the ordering options for the Role queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByName orders the results by the name field.
func ByName(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldName, opts...).ToFunc()
}

// ByBuildin orders the results by the buildin field.
func ByBuildin(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldBuildin, opts...).ToFunc()
}

// BySuper orders the results by the super field.
func BySuper(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldSuper, opts...).ToFunc()
}

// ByCreatedAt orders the results by the created_at field.
func ByCreatedAt(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldCreatedAt, opts...).ToFunc()
}

// ByManagersCount orders the results by managers count.
func ByManagersCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newManagersStep(), opts...)
	}
}

// ByManagers orders the results by managers terms.
func ByManagers(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newManagersStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}
func newManagersStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(ManagersInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, ManagersTable, ManagersColumn),
	)
}
