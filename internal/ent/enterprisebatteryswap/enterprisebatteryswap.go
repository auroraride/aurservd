// Code generated by ent, DO NOT EDIT.

package enterprisebatteryswap

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
)

const (
	// Label holds the string label denoting the enterprisebatteryswap type in the database.
	Label = "enterprise_battery_swap"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldCreatedAt holds the string denoting the created_at field in the database.
	FieldCreatedAt = "created_at"
	// FieldUpdatedAt holds the string denoting the updated_at field in the database.
	FieldUpdatedAt = "updated_at"
	// FieldExchangeID holds the string denoting the exchange_id field in the database.
	FieldExchangeID = "exchange_id"
	// FieldCabinetID holds the string denoting the cabinet_id field in the database.
	FieldCabinetID = "cabinet_id"
	// FieldPutinID holds the string denoting the putin_id field in the database.
	FieldPutinID = "putin_id"
	// FieldPutinSn holds the string denoting the putin_sn field in the database.
	FieldPutinSn = "putin_sn"
	// FieldPutinEnterpriseID holds the string denoting the putin_enterprise_id field in the database.
	FieldPutinEnterpriseID = "putin_enterprise_id"
	// FieldPutinStationID holds the string denoting the putin_station_id field in the database.
	FieldPutinStationID = "putin_station_id"
	// FieldPutoutID holds the string denoting the putout_id field in the database.
	FieldPutoutID = "putout_id"
	// FieldPutoutSn holds the string denoting the putout_sn field in the database.
	FieldPutoutSn = "putout_sn"
	// FieldPutoutEnterpriseID holds the string denoting the putout_enterprise_id field in the database.
	FieldPutoutEnterpriseID = "putout_enterprise_id"
	// FieldPutoutStationID holds the string denoting the putout_station_id field in the database.
	FieldPutoutStationID = "putout_station_id"
	// EdgeExchange holds the string denoting the exchange edge name in mutations.
	EdgeExchange = "exchange"
	// EdgeCabinet holds the string denoting the cabinet edge name in mutations.
	EdgeCabinet = "cabinet"
	// EdgePutin holds the string denoting the putin edge name in mutations.
	EdgePutin = "putin"
	// EdgePutinEnterprise holds the string denoting the putin_enterprise edge name in mutations.
	EdgePutinEnterprise = "putin_enterprise"
	// EdgePutinStation holds the string denoting the putin_station edge name in mutations.
	EdgePutinStation = "putin_station"
	// EdgePutout holds the string denoting the putout edge name in mutations.
	EdgePutout = "putout"
	// EdgePutoutEnterprise holds the string denoting the putout_enterprise edge name in mutations.
	EdgePutoutEnterprise = "putout_enterprise"
	// EdgePutoutStation holds the string denoting the putout_station edge name in mutations.
	EdgePutoutStation = "putout_station"
	// Table holds the table name of the enterprisebatteryswap in the database.
	Table = "enterprise_battery_swap"
	// ExchangeTable is the table that holds the exchange relation/edge.
	ExchangeTable = "enterprise_battery_swap"
	// ExchangeInverseTable is the table name for the Exchange entity.
	// It exists in this package in order to avoid circular dependency with the "exchange" package.
	ExchangeInverseTable = "exchange"
	// ExchangeColumn is the table column denoting the exchange relation/edge.
	ExchangeColumn = "exchange_id"
	// CabinetTable is the table that holds the cabinet relation/edge.
	CabinetTable = "enterprise_battery_swap"
	// CabinetInverseTable is the table name for the Cabinet entity.
	// It exists in this package in order to avoid circular dependency with the "cabinet" package.
	CabinetInverseTable = "cabinet"
	// CabinetColumn is the table column denoting the cabinet relation/edge.
	CabinetColumn = "cabinet_id"
	// PutinTable is the table that holds the putin relation/edge.
	PutinTable = "enterprise_battery_swap"
	// PutinInverseTable is the table name for the Battery entity.
	// It exists in this package in order to avoid circular dependency with the "battery" package.
	PutinInverseTable = "battery"
	// PutinColumn is the table column denoting the putin relation/edge.
	PutinColumn = "putin_id"
	// PutinEnterpriseTable is the table that holds the putin_enterprise relation/edge.
	PutinEnterpriseTable = "enterprise_battery_swap"
	// PutinEnterpriseInverseTable is the table name for the Enterprise entity.
	// It exists in this package in order to avoid circular dependency with the "enterprise" package.
	PutinEnterpriseInverseTable = "enterprise"
	// PutinEnterpriseColumn is the table column denoting the putin_enterprise relation/edge.
	PutinEnterpriseColumn = "putin_enterprise_id"
	// PutinStationTable is the table that holds the putin_station relation/edge.
	PutinStationTable = "enterprise_battery_swap"
	// PutinStationInverseTable is the table name for the EnterpriseStation entity.
	// It exists in this package in order to avoid circular dependency with the "enterprisestation" package.
	PutinStationInverseTable = "enterprise_station"
	// PutinStationColumn is the table column denoting the putin_station relation/edge.
	PutinStationColumn = "putin_station_id"
	// PutoutTable is the table that holds the putout relation/edge.
	PutoutTable = "enterprise_battery_swap"
	// PutoutInverseTable is the table name for the Battery entity.
	// It exists in this package in order to avoid circular dependency with the "battery" package.
	PutoutInverseTable = "battery"
	// PutoutColumn is the table column denoting the putout relation/edge.
	PutoutColumn = "putout_id"
	// PutoutEnterpriseTable is the table that holds the putout_enterprise relation/edge.
	PutoutEnterpriseTable = "enterprise_battery_swap"
	// PutoutEnterpriseInverseTable is the table name for the Enterprise entity.
	// It exists in this package in order to avoid circular dependency with the "enterprise" package.
	PutoutEnterpriseInverseTable = "enterprise"
	// PutoutEnterpriseColumn is the table column denoting the putout_enterprise relation/edge.
	PutoutEnterpriseColumn = "putout_enterprise_id"
	// PutoutStationTable is the table that holds the putout_station relation/edge.
	PutoutStationTable = "enterprise_battery_swap"
	// PutoutStationInverseTable is the table name for the EnterpriseStation entity.
	// It exists in this package in order to avoid circular dependency with the "enterprisestation" package.
	PutoutStationInverseTable = "enterprise_station"
	// PutoutStationColumn is the table column denoting the putout_station relation/edge.
	PutoutStationColumn = "putout_station_id"
)

// Columns holds all SQL columns for enterprisebatteryswap fields.
var Columns = []string{
	FieldID,
	FieldCreatedAt,
	FieldUpdatedAt,
	FieldExchangeID,
	FieldCabinetID,
	FieldPutinID,
	FieldPutinSn,
	FieldPutinEnterpriseID,
	FieldPutinStationID,
	FieldPutoutID,
	FieldPutoutSn,
	FieldPutoutEnterpriseID,
	FieldPutoutStationID,
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

// OrderOption defines the ordering options for the EnterpriseBatterySwap queries.
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

// ByExchangeID orders the results by the exchange_id field.
func ByExchangeID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldExchangeID, opts...).ToFunc()
}

// ByCabinetID orders the results by the cabinet_id field.
func ByCabinetID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldCabinetID, opts...).ToFunc()
}

// ByPutinID orders the results by the putin_id field.
func ByPutinID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldPutinID, opts...).ToFunc()
}

// ByPutinSn orders the results by the putin_sn field.
func ByPutinSn(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldPutinSn, opts...).ToFunc()
}

// ByPutinEnterpriseID orders the results by the putin_enterprise_id field.
func ByPutinEnterpriseID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldPutinEnterpriseID, opts...).ToFunc()
}

// ByPutinStationID orders the results by the putin_station_id field.
func ByPutinStationID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldPutinStationID, opts...).ToFunc()
}

// ByPutoutID orders the results by the putout_id field.
func ByPutoutID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldPutoutID, opts...).ToFunc()
}

// ByPutoutSn orders the results by the putout_sn field.
func ByPutoutSn(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldPutoutSn, opts...).ToFunc()
}

// ByPutoutEnterpriseID orders the results by the putout_enterprise_id field.
func ByPutoutEnterpriseID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldPutoutEnterpriseID, opts...).ToFunc()
}

// ByPutoutStationID orders the results by the putout_station_id field.
func ByPutoutStationID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldPutoutStationID, opts...).ToFunc()
}

// ByExchangeField orders the results by exchange field.
func ByExchangeField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newExchangeStep(), sql.OrderByField(field, opts...))
	}
}

// ByCabinetField orders the results by cabinet field.
func ByCabinetField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newCabinetStep(), sql.OrderByField(field, opts...))
	}
}

// ByPutinField orders the results by putin field.
func ByPutinField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newPutinStep(), sql.OrderByField(field, opts...))
	}
}

// ByPutinEnterpriseField orders the results by putin_enterprise field.
func ByPutinEnterpriseField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newPutinEnterpriseStep(), sql.OrderByField(field, opts...))
	}
}

// ByPutinStationField orders the results by putin_station field.
func ByPutinStationField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newPutinStationStep(), sql.OrderByField(field, opts...))
	}
}

// ByPutoutField orders the results by putout field.
func ByPutoutField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newPutoutStep(), sql.OrderByField(field, opts...))
	}
}

// ByPutoutEnterpriseField orders the results by putout_enterprise field.
func ByPutoutEnterpriseField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newPutoutEnterpriseStep(), sql.OrderByField(field, opts...))
	}
}

// ByPutoutStationField orders the results by putout_station field.
func ByPutoutStationField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newPutoutStationStep(), sql.OrderByField(field, opts...))
	}
}
func newExchangeStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(ExchangeInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2O, false, ExchangeTable, ExchangeColumn),
	)
}
func newCabinetStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(CabinetInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2O, false, CabinetTable, CabinetColumn),
	)
}
func newPutinStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(PutinInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2O, false, PutinTable, PutinColumn),
	)
}
func newPutinEnterpriseStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(PutinEnterpriseInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, PutinEnterpriseTable, PutinEnterpriseColumn),
	)
}
func newPutinStationStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(PutinStationInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, PutinStationTable, PutinStationColumn),
	)
}
func newPutoutStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(PutoutInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2O, false, PutoutTable, PutoutColumn),
	)
}
func newPutoutEnterpriseStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(PutoutEnterpriseInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, PutoutEnterpriseTable, PutoutEnterpriseColumn),
	)
}
func newPutoutStationStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(PutoutStationInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, PutoutStationTable, PutoutStationColumn),
	)
}
