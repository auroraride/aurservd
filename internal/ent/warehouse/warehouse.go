// Code generated by ent, DO NOT EDIT.

package warehouse

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
)

const (
	// Label holds the string label denoting the warehouse type in the database.
	Label = "warehouse"
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
	// FieldCityID holds the string denoting the city_id field in the database.
	FieldCityID = "city_id"
	// FieldName holds the string denoting the name field in the database.
	FieldName = "name"
	// FieldLng holds the string denoting the lng field in the database.
	FieldLng = "lng"
	// FieldLat holds the string denoting the lat field in the database.
	FieldLat = "lat"
	// FieldAddress holds the string denoting the address field in the database.
	FieldAddress = "address"
	// FieldSn holds the string denoting the sn field in the database.
	FieldSn = "sn"
	// EdgeCity holds the string denoting the city edge name in mutations.
	EdgeCity = "city"
	// EdgeBelongAssetManagers holds the string denoting the belong_asset_managers edge name in mutations.
	EdgeBelongAssetManagers = "belong_asset_managers"
	// EdgeDutyAssetManagers holds the string denoting the duty_asset_managers edge name in mutations.
	EdgeDutyAssetManagers = "duty_asset_managers"
	// EdgeAsset holds the string denoting the asset edge name in mutations.
	EdgeAsset = "asset"
	// Table holds the table name of the warehouse in the database.
	Table = "warehouse"
	// CityTable is the table that holds the city relation/edge.
	CityTable = "warehouse"
	// CityInverseTable is the table name for the City entity.
	// It exists in this package in order to avoid circular dependency with the "city" package.
	CityInverseTable = "city"
	// CityColumn is the table column denoting the city relation/edge.
	CityColumn = "city_id"
	// BelongAssetManagersTable is the table that holds the belong_asset_managers relation/edge. The primary key declared below.
	BelongAssetManagersTable = "warehouse_belong_asset_managers"
	// BelongAssetManagersInverseTable is the table name for the AssetManager entity.
	// It exists in this package in order to avoid circular dependency with the "assetmanager" package.
	BelongAssetManagersInverseTable = "asset_manager"
	// DutyAssetManagersTable is the table that holds the duty_asset_managers relation/edge.
	DutyAssetManagersTable = "asset_manager"
	// DutyAssetManagersInverseTable is the table name for the AssetManager entity.
	// It exists in this package in order to avoid circular dependency with the "assetmanager" package.
	DutyAssetManagersInverseTable = "asset_manager"
	// DutyAssetManagersColumn is the table column denoting the duty_asset_managers relation/edge.
	DutyAssetManagersColumn = "warehouse_id"
	// AssetTable is the table that holds the asset relation/edge.
	AssetTable = "asset"
	// AssetInverseTable is the table name for the Asset entity.
	// It exists in this package in order to avoid circular dependency with the "asset" package.
	AssetInverseTable = "asset"
	// AssetColumn is the table column denoting the asset relation/edge.
	AssetColumn = "locations_id"
)

// Columns holds all SQL columns for warehouse fields.
var Columns = []string{
	FieldID,
	FieldCreatedAt,
	FieldUpdatedAt,
	FieldDeletedAt,
	FieldCreator,
	FieldLastModifier,
	FieldRemark,
	FieldCityID,
	FieldName,
	FieldLng,
	FieldLat,
	FieldAddress,
	FieldSn,
}

var (
	// BelongAssetManagersPrimaryKey and BelongAssetManagersColumn2 are the table columns denoting the
	// primary key for the belong_asset_managers relation (M2M).
	BelongAssetManagersPrimaryKey = []string{"warehouse_id", "asset_manager_id"}
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
)

// OrderOption defines the ordering options for the Warehouse queries.
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

// ByCityID orders the results by the city_id field.
func ByCityID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldCityID, opts...).ToFunc()
}

// ByName orders the results by the name field.
func ByName(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldName, opts...).ToFunc()
}

// ByLng orders the results by the lng field.
func ByLng(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldLng, opts...).ToFunc()
}

// ByLat orders the results by the lat field.
func ByLat(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldLat, opts...).ToFunc()
}

// ByAddress orders the results by the address field.
func ByAddress(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldAddress, opts...).ToFunc()
}

// BySn orders the results by the sn field.
func BySn(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldSn, opts...).ToFunc()
}

// ByCityField orders the results by city field.
func ByCityField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newCityStep(), sql.OrderByField(field, opts...))
	}
}

// ByBelongAssetManagersCount orders the results by belong_asset_managers count.
func ByBelongAssetManagersCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newBelongAssetManagersStep(), opts...)
	}
}

// ByBelongAssetManagers orders the results by belong_asset_managers terms.
func ByBelongAssetManagers(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newBelongAssetManagersStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}

// ByDutyAssetManagersCount orders the results by duty_asset_managers count.
func ByDutyAssetManagersCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newDutyAssetManagersStep(), opts...)
	}
}

// ByDutyAssetManagers orders the results by duty_asset_managers terms.
func ByDutyAssetManagers(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newDutyAssetManagersStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}

// ByAssetCount orders the results by asset count.
func ByAssetCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newAssetStep(), opts...)
	}
}

// ByAsset orders the results by asset terms.
func ByAsset(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newAssetStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}
func newCityStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(CityInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2O, false, CityTable, CityColumn),
	)
}
func newBelongAssetManagersStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(BelongAssetManagersInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2M, false, BelongAssetManagersTable, BelongAssetManagersPrimaryKey...),
	)
}
func newDutyAssetManagersStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(DutyAssetManagersInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, DutyAssetManagersTable, DutyAssetManagersColumn),
	)
}
func newAssetStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(AssetInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, AssetTable, AssetColumn),
	)
}