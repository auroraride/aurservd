// Code generated by ent, DO NOT EDIT.

package assetmanager

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
)

const (
	// Label holds the string label denoting the assetmanager type in the database.
	Label = "asset_manager"
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
	// FieldPhone holds the string denoting the phone field in the database.
	FieldPhone = "phone"
	// FieldPassword holds the string denoting the password field in the database.
	FieldPassword = "password"
	// FieldRoleID holds the string denoting the role_id field in the database.
	FieldRoleID = "role_id"
	// FieldMiniEnable holds the string denoting the mini_enable field in the database.
	FieldMiniEnable = "mini_enable"
	// FieldMiniLimit holds the string denoting the mini_limit field in the database.
	FieldMiniLimit = "mini_limit"
	// FieldLastSigninAt holds the string denoting the last_signin_at field in the database.
	FieldLastSigninAt = "last_signin_at"
	// FieldWarehouseID holds the string denoting the warehouse_id field in the database.
	FieldWarehouseID = "warehouse_id"
	// EdgeRole holds the string denoting the role edge name in mutations.
	EdgeRole = "role"
	// EdgeBelongWarehouses holds the string denoting the belong_warehouses edge name in mutations.
	EdgeBelongWarehouses = "belong_warehouses"
	// EdgeDutyWarehouse holds the string denoting the duty_warehouse edge name in mutations.
	EdgeDutyWarehouse = "duty_warehouse"
	// Table holds the table name of the assetmanager in the database.
	Table = "asset_manager"
	// RoleTable is the table that holds the role relation/edge.
	RoleTable = "asset_manager"
	// RoleInverseTable is the table name for the AssetRole entity.
	// It exists in this package in order to avoid circular dependency with the "assetrole" package.
	RoleInverseTable = "asset_role"
	// RoleColumn is the table column denoting the role relation/edge.
	RoleColumn = "role_id"
	// BelongWarehousesTable is the table that holds the belong_warehouses relation/edge. The primary key declared below.
	BelongWarehousesTable = "warehouse_belong_asset_managers"
	// BelongWarehousesInverseTable is the table name for the Warehouse entity.
	// It exists in this package in order to avoid circular dependency with the "warehouse" package.
	BelongWarehousesInverseTable = "warehouse"
	// DutyWarehouseTable is the table that holds the duty_warehouse relation/edge.
	DutyWarehouseTable = "asset_manager"
	// DutyWarehouseInverseTable is the table name for the Warehouse entity.
	// It exists in this package in order to avoid circular dependency with the "warehouse" package.
	DutyWarehouseInverseTable = "warehouse"
	// DutyWarehouseColumn is the table column denoting the duty_warehouse relation/edge.
	DutyWarehouseColumn = "warehouse_id"
)

// Columns holds all SQL columns for assetmanager fields.
var Columns = []string{
	FieldID,
	FieldCreatedAt,
	FieldUpdatedAt,
	FieldDeletedAt,
	FieldCreator,
	FieldLastModifier,
	FieldRemark,
	FieldName,
	FieldPhone,
	FieldPassword,
	FieldRoleID,
	FieldMiniEnable,
	FieldMiniLimit,
	FieldLastSigninAt,
	FieldWarehouseID,
}

var (
	// BelongWarehousesPrimaryKey and BelongWarehousesColumn2 are the table columns denoting the
	// primary key for the belong_warehouses relation (M2M).
	BelongWarehousesPrimaryKey = []string{"warehouse_id", "asset_manager_id"}
)

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
	// NameValidator is a validator for the "name" field. It is called by the builders before save.
	NameValidator func(string) error
	// PhoneValidator is a validator for the "phone" field. It is called by the builders before save.
	PhoneValidator func(string) error
	// DefaultMiniEnable holds the default value on creation for the "mini_enable" field.
	DefaultMiniEnable bool
	// DefaultMiniLimit holds the default value on creation for the "mini_limit" field.
	DefaultMiniLimit uint
)

// OrderOption defines the ordering options for the AssetManager queries.
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

// ByPhone orders the results by the phone field.
func ByPhone(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldPhone, opts...).ToFunc()
}

// ByPassword orders the results by the password field.
func ByPassword(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldPassword, opts...).ToFunc()
}

// ByRoleID orders the results by the role_id field.
func ByRoleID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldRoleID, opts...).ToFunc()
}

// ByMiniEnable orders the results by the mini_enable field.
func ByMiniEnable(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldMiniEnable, opts...).ToFunc()
}

// ByMiniLimit orders the results by the mini_limit field.
func ByMiniLimit(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldMiniLimit, opts...).ToFunc()
}

// ByLastSigninAt orders the results by the last_signin_at field.
func ByLastSigninAt(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldLastSigninAt, opts...).ToFunc()
}

// ByWarehouseID orders the results by the warehouse_id field.
func ByWarehouseID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldWarehouseID, opts...).ToFunc()
}

// ByRoleField orders the results by role field.
func ByRoleField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newRoleStep(), sql.OrderByField(field, opts...))
	}
}

// ByBelongWarehousesCount orders the results by belong_warehouses count.
func ByBelongWarehousesCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newBelongWarehousesStep(), opts...)
	}
}

// ByBelongWarehouses orders the results by belong_warehouses terms.
func ByBelongWarehouses(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newBelongWarehousesStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}

// ByDutyWarehouseField orders the results by duty_warehouse field.
func ByDutyWarehouseField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newDutyWarehouseStep(), sql.OrderByField(field, opts...))
	}
}
func newRoleStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(RoleInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, RoleTable, RoleColumn),
	)
}
func newBelongWarehousesStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(BelongWarehousesInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2M, true, BelongWarehousesTable, BelongWarehousesPrimaryKey...),
	)
}
func newDutyWarehouseStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(DutyWarehouseInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, DutyWarehouseTable, DutyWarehouseColumn),
	)
}