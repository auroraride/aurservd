// Code generated by ent, DO NOT EDIT.

package assetcheck

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
)

const (
	// Label holds the string label denoting the assetcheck type in the database.
	Label = "asset_check"
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
	// FieldStatus holds the string denoting the status field in the database.
	FieldStatus = "status"
	// FieldBatteryNum holds the string denoting the battery_num field in the database.
	FieldBatteryNum = "battery_num"
	// FieldBatteryNumReal holds the string denoting the battery_num_real field in the database.
	FieldBatteryNumReal = "battery_num_real"
	// FieldEbikeNum holds the string denoting the ebike_num field in the database.
	FieldEbikeNum = "ebike_num"
	// FieldEbikeNumReal holds the string denoting the ebike_num_real field in the database.
	FieldEbikeNumReal = "ebike_num_real"
	// FieldOperateID holds the string denoting the operate_id field in the database.
	FieldOperateID = "operate_id"
	// FieldOperateType holds the string denoting the operate_type field in the database.
	FieldOperateType = "operate_type"
	// FieldLocationsType holds the string denoting the locations_type field in the database.
	FieldLocationsType = "locations_type"
	// FieldLocationsID holds the string denoting the locations_id field in the database.
	FieldLocationsID = "locations_id"
	// FieldStartAt holds the string denoting the start_at field in the database.
	FieldStartAt = "start_at"
	// FieldEndAt holds the string denoting the end_at field in the database.
	FieldEndAt = "end_at"
	// EdgeCheckDetails holds the string denoting the check_details edge name in mutations.
	EdgeCheckDetails = "check_details"
	// EdgeOperateManager holds the string denoting the operate_manager edge name in mutations.
	EdgeOperateManager = "operate_manager"
	// EdgeOperateStore holds the string denoting the operate_store edge name in mutations.
	EdgeOperateStore = "operate_store"
	// EdgeOperateAgent holds the string denoting the operate_agent edge name in mutations.
	EdgeOperateAgent = "operate_agent"
	// EdgeWarehouse holds the string denoting the warehouse edge name in mutations.
	EdgeWarehouse = "warehouse"
	// EdgeStore holds the string denoting the store edge name in mutations.
	EdgeStore = "store"
	// EdgeStation holds the string denoting the station edge name in mutations.
	EdgeStation = "station"
	// Table holds the table name of the assetcheck in the database.
	Table = "asset_check"
	// CheckDetailsTable is the table that holds the check_details relation/edge.
	CheckDetailsTable = "asset_check_details"
	// CheckDetailsInverseTable is the table name for the AssetCheckDetails entity.
	// It exists in this package in order to avoid circular dependency with the "assetcheckdetails" package.
	CheckDetailsInverseTable = "asset_check_details"
	// CheckDetailsColumn is the table column denoting the check_details relation/edge.
	CheckDetailsColumn = "check_id"
	// OperateManagerTable is the table that holds the operate_manager relation/edge.
	OperateManagerTable = "asset_check"
	// OperateManagerInverseTable is the table name for the AssetManager entity.
	// It exists in this package in order to avoid circular dependency with the "assetmanager" package.
	OperateManagerInverseTable = "asset_manager"
	// OperateManagerColumn is the table column denoting the operate_manager relation/edge.
	OperateManagerColumn = "operate_id"
	// OperateStoreTable is the table that holds the operate_store relation/edge.
	OperateStoreTable = "asset_check"
	// OperateStoreInverseTable is the table name for the Store entity.
	// It exists in this package in order to avoid circular dependency with the "store" package.
	OperateStoreInverseTable = "store"
	// OperateStoreColumn is the table column denoting the operate_store relation/edge.
	OperateStoreColumn = "operate_id"
	// OperateAgentTable is the table that holds the operate_agent relation/edge.
	OperateAgentTable = "asset_check"
	// OperateAgentInverseTable is the table name for the Agent entity.
	// It exists in this package in order to avoid circular dependency with the "agent" package.
	OperateAgentInverseTable = "agent"
	// OperateAgentColumn is the table column denoting the operate_agent relation/edge.
	OperateAgentColumn = "operate_id"
	// WarehouseTable is the table that holds the warehouse relation/edge.
	WarehouseTable = "asset_check"
	// WarehouseInverseTable is the table name for the Warehouse entity.
	// It exists in this package in order to avoid circular dependency with the "warehouse" package.
	WarehouseInverseTable = "warehouse"
	// WarehouseColumn is the table column denoting the warehouse relation/edge.
	WarehouseColumn = "locations_id"
	// StoreTable is the table that holds the store relation/edge.
	StoreTable = "asset_check"
	// StoreInverseTable is the table name for the Store entity.
	// It exists in this package in order to avoid circular dependency with the "store" package.
	StoreInverseTable = "store"
	// StoreColumn is the table column denoting the store relation/edge.
	StoreColumn = "locations_id"
	// StationTable is the table that holds the station relation/edge.
	StationTable = "asset_check"
	// StationInverseTable is the table name for the EnterpriseStation entity.
	// It exists in this package in order to avoid circular dependency with the "enterprisestation" package.
	StationInverseTable = "enterprise_station"
	// StationColumn is the table column denoting the station relation/edge.
	StationColumn = "locations_id"
)

// Columns holds all SQL columns for assetcheck fields.
var Columns = []string{
	FieldID,
	FieldCreatedAt,
	FieldUpdatedAt,
	FieldDeletedAt,
	FieldCreator,
	FieldLastModifier,
	FieldRemark,
	FieldStatus,
	FieldBatteryNum,
	FieldBatteryNumReal,
	FieldEbikeNum,
	FieldEbikeNumReal,
	FieldOperateID,
	FieldOperateType,
	FieldLocationsType,
	FieldLocationsID,
	FieldStartAt,
	FieldEndAt,
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

// OrderOption defines the ordering options for the AssetCheck queries.
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

// ByStatus orders the results by the status field.
func ByStatus(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldStatus, opts...).ToFunc()
}

// ByBatteryNum orders the results by the battery_num field.
func ByBatteryNum(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldBatteryNum, opts...).ToFunc()
}

// ByBatteryNumReal orders the results by the battery_num_real field.
func ByBatteryNumReal(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldBatteryNumReal, opts...).ToFunc()
}

// ByEbikeNum orders the results by the ebike_num field.
func ByEbikeNum(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldEbikeNum, opts...).ToFunc()
}

// ByEbikeNumReal orders the results by the ebike_num_real field.
func ByEbikeNumReal(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldEbikeNumReal, opts...).ToFunc()
}

// ByOperateID orders the results by the operate_id field.
func ByOperateID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldOperateID, opts...).ToFunc()
}

// ByOperateType orders the results by the operate_type field.
func ByOperateType(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldOperateType, opts...).ToFunc()
}

// ByLocationsType orders the results by the locations_type field.
func ByLocationsType(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldLocationsType, opts...).ToFunc()
}

// ByLocationsID orders the results by the locations_id field.
func ByLocationsID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldLocationsID, opts...).ToFunc()
}

// ByStartAt orders the results by the start_at field.
func ByStartAt(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldStartAt, opts...).ToFunc()
}

// ByEndAt orders the results by the end_at field.
func ByEndAt(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldEndAt, opts...).ToFunc()
}

// ByCheckDetailsCount orders the results by check_details count.
func ByCheckDetailsCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newCheckDetailsStep(), opts...)
	}
}

// ByCheckDetails orders the results by check_details terms.
func ByCheckDetails(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newCheckDetailsStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}

// ByOperateManagerField orders the results by operate_manager field.
func ByOperateManagerField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newOperateManagerStep(), sql.OrderByField(field, opts...))
	}
}

// ByOperateStoreField orders the results by operate_store field.
func ByOperateStoreField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newOperateStoreStep(), sql.OrderByField(field, opts...))
	}
}

// ByOperateAgentField orders the results by operate_agent field.
func ByOperateAgentField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newOperateAgentStep(), sql.OrderByField(field, opts...))
	}
}

// ByWarehouseField orders the results by warehouse field.
func ByWarehouseField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newWarehouseStep(), sql.OrderByField(field, opts...))
	}
}

// ByStoreField orders the results by store field.
func ByStoreField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newStoreStep(), sql.OrderByField(field, opts...))
	}
}

// ByStationField orders the results by station field.
func ByStationField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newStationStep(), sql.OrderByField(field, opts...))
	}
}
func newCheckDetailsStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(CheckDetailsInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, CheckDetailsTable, CheckDetailsColumn),
	)
}
func newOperateManagerStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(OperateManagerInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2O, false, OperateManagerTable, OperateManagerColumn),
	)
}
func newOperateStoreStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(OperateStoreInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2O, false, OperateStoreTable, OperateStoreColumn),
	)
}
func newOperateAgentStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(OperateAgentInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2O, false, OperateAgentTable, OperateAgentColumn),
	)
}
func newWarehouseStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(WarehouseInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2O, false, WarehouseTable, WarehouseColumn),
	)
}
func newStoreStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(StoreInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2O, false, StoreTable, StoreColumn),
	)
}
func newStationStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(StationInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2O, false, StationTable, StationColumn),
	)
}