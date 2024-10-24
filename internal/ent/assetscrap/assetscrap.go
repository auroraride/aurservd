// Code generated by ent, DO NOT EDIT.

package assetscrap

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
)

const (
	// Label holds the string label denoting the assetscrap type in the database.
	Label = "asset_scrap"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldCreatedAt holds the string denoting the created_at field in the database.
	FieldCreatedAt = "created_at"
	// FieldUpdatedAt holds the string denoting the updated_at field in the database.
	FieldUpdatedAt = "updated_at"
	// FieldCreator holds the string denoting the creator field in the database.
	FieldCreator = "creator"
	// FieldLastModifier holds the string denoting the last_modifier field in the database.
	FieldLastModifier = "last_modifier"
	// FieldRemark holds the string denoting the remark field in the database.
	FieldRemark = "remark"
	// FieldReasonType holds the string denoting the reason_type field in the database.
	FieldReasonType = "reason_type"
	// FieldScrapAt holds the string denoting the scrap_at field in the database.
	FieldScrapAt = "scrap_at"
	// FieldOperateID holds the string denoting the operate_id field in the database.
	FieldOperateID = "operate_id"
	// FieldOperateRoleType holds the string denoting the operate_role_type field in the database.
	FieldOperateRoleType = "operate_role_type"
	// FieldSn holds the string denoting the sn field in the database.
	FieldSn = "sn"
	// FieldNum holds the string denoting the num field in the database.
	FieldNum = "num"
	// EdgeManager holds the string denoting the manager edge name in mutations.
	EdgeManager = "manager"
	// EdgeEmployee holds the string denoting the employee edge name in mutations.
	EdgeEmployee = "employee"
	// EdgeMaintainer holds the string denoting the maintainer edge name in mutations.
	EdgeMaintainer = "maintainer"
	// EdgeAgent holds the string denoting the agent edge name in mutations.
	EdgeAgent = "agent"
	// EdgeScrapDetails holds the string denoting the scrap_details edge name in mutations.
	EdgeScrapDetails = "scrap_details"
	// Table holds the table name of the assetscrap in the database.
	Table = "asset_scrap"
	// ManagerTable is the table that holds the manager relation/edge.
	ManagerTable = "asset_scrap"
	// ManagerInverseTable is the table name for the AssetManager entity.
	// It exists in this package in order to avoid circular dependency with the "assetmanager" package.
	ManagerInverseTable = "asset_manager"
	// ManagerColumn is the table column denoting the manager relation/edge.
	ManagerColumn = "operate_id"
	// EmployeeTable is the table that holds the employee relation/edge.
	EmployeeTable = "asset_scrap"
	// EmployeeInverseTable is the table name for the Employee entity.
	// It exists in this package in order to avoid circular dependency with the "employee" package.
	EmployeeInverseTable = "employee"
	// EmployeeColumn is the table column denoting the employee relation/edge.
	EmployeeColumn = "operate_id"
	// MaintainerTable is the table that holds the maintainer relation/edge.
	MaintainerTable = "asset_scrap"
	// MaintainerInverseTable is the table name for the Maintainer entity.
	// It exists in this package in order to avoid circular dependency with the "maintainer" package.
	MaintainerInverseTable = "maintainer"
	// MaintainerColumn is the table column denoting the maintainer relation/edge.
	MaintainerColumn = "operate_id"
	// AgentTable is the table that holds the agent relation/edge.
	AgentTable = "asset_scrap"
	// AgentInverseTable is the table name for the Agent entity.
	// It exists in this package in order to avoid circular dependency with the "agent" package.
	AgentInverseTable = "agent"
	// AgentColumn is the table column denoting the agent relation/edge.
	AgentColumn = "operate_id"
	// ScrapDetailsTable is the table that holds the scrap_details relation/edge.
	ScrapDetailsTable = "asset_scrap_details"
	// ScrapDetailsInverseTable is the table name for the AssetScrapDetails entity.
	// It exists in this package in order to avoid circular dependency with the "assetscrapdetails" package.
	ScrapDetailsInverseTable = "asset_scrap_details"
	// ScrapDetailsColumn is the table column denoting the scrap_details relation/edge.
	ScrapDetailsColumn = "scrap_id"
)

// Columns holds all SQL columns for assetscrap fields.
var Columns = []string{
	FieldID,
	FieldCreatedAt,
	FieldUpdatedAt,
	FieldCreator,
	FieldLastModifier,
	FieldRemark,
	FieldReasonType,
	FieldScrapAt,
	FieldOperateID,
	FieldOperateRoleType,
	FieldSn,
	FieldNum,
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

// OrderOption defines the ordering options for the AssetScrap queries.
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

// ByRemark orders the results by the remark field.
func ByRemark(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldRemark, opts...).ToFunc()
}

// ByReasonType orders the results by the reason_type field.
func ByReasonType(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldReasonType, opts...).ToFunc()
}

// ByScrapAt orders the results by the scrap_at field.
func ByScrapAt(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldScrapAt, opts...).ToFunc()
}

// ByOperateID orders the results by the operate_id field.
func ByOperateID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldOperateID, opts...).ToFunc()
}

// ByOperateRoleType orders the results by the operate_role_type field.
func ByOperateRoleType(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldOperateRoleType, opts...).ToFunc()
}

// BySn orders the results by the sn field.
func BySn(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldSn, opts...).ToFunc()
}

// ByNum orders the results by the num field.
func ByNum(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldNum, opts...).ToFunc()
}

// ByManagerField orders the results by manager field.
func ByManagerField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newManagerStep(), sql.OrderByField(field, opts...))
	}
}

// ByEmployeeField orders the results by employee field.
func ByEmployeeField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newEmployeeStep(), sql.OrderByField(field, opts...))
	}
}

// ByMaintainerField orders the results by maintainer field.
func ByMaintainerField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newMaintainerStep(), sql.OrderByField(field, opts...))
	}
}

// ByAgentField orders the results by agent field.
func ByAgentField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newAgentStep(), sql.OrderByField(field, opts...))
	}
}

// ByScrapDetailsCount orders the results by scrap_details count.
func ByScrapDetailsCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newScrapDetailsStep(), opts...)
	}
}

// ByScrapDetails orders the results by scrap_details terms.
func ByScrapDetails(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newScrapDetailsStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}
func newManagerStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(ManagerInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2O, false, ManagerTable, ManagerColumn),
	)
}
func newEmployeeStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(EmployeeInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2O, false, EmployeeTable, EmployeeColumn),
	)
}
func newMaintainerStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(MaintainerInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2O, false, MaintainerTable, MaintainerColumn),
	)
}
func newAgentStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(AgentInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2O, false, AgentTable, AgentColumn),
	)
}
func newScrapDetailsStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(ScrapDetailsInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, ScrapDetailsTable, ScrapDetailsColumn),
	)
}
