// Code generated by ent, DO NOT EDIT.

package asset

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
)

const (
	// Label holds the string label denoting the asset type in the database.
	Label = "asset"
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
	// FieldBrandID holds the string denoting the brand_id field in the database.
	FieldBrandID = "brand_id"
	// FieldModelID holds the string denoting the model_id field in the database.
	FieldModelID = "model_id"
	// FieldCityID holds the string denoting the city_id field in the database.
	FieldCityID = "city_id"
	// FieldMaterialID holds the string denoting the material_id field in the database.
	FieldMaterialID = "material_id"
	// FieldType holds the string denoting the type field in the database.
	FieldType = "type"
	// FieldName holds the string denoting the name field in the database.
	FieldName = "name"
	// FieldSn holds the string denoting the sn field in the database.
	FieldSn = "sn"
	// FieldStatus holds the string denoting the status field in the database.
	FieldStatus = "status"
	// FieldEnable holds the string denoting the enable field in the database.
	FieldEnable = "enable"
	// FieldLocationsType holds the string denoting the locations_type field in the database.
	FieldLocationsType = "locations_type"
	// FieldLocationsID holds the string denoting the locations_id field in the database.
	FieldLocationsID = "locations_id"
	// FieldRtoRiderID holds the string denoting the rto_rider_id field in the database.
	FieldRtoRiderID = "rto_rider_id"
	// FieldCheckAt holds the string denoting the check_at field in the database.
	FieldCheckAt = "check_at"
	// FieldBrandName holds the string denoting the brand_name field in the database.
	FieldBrandName = "brand_name"
	// FieldSubscribeID holds the string denoting the subscribe_id field in the database.
	FieldSubscribeID = "subscribe_id"
	// FieldOrdinal holds the string denoting the ordinal field in the database.
	FieldOrdinal = "ordinal"
	// FieldRentLocationsType holds the string denoting the rent_locations_type field in the database.
	FieldRentLocationsType = "rent_locations_type"
	// FieldRentLocationsID holds the string denoting the rent_locations_id field in the database.
	FieldRentLocationsID = "rent_locations_id"
	// EdgeBrand holds the string denoting the brand edge name in mutations.
	EdgeBrand = "brand"
	// EdgeModel holds the string denoting the model edge name in mutations.
	EdgeModel = "model"
	// EdgeCity holds the string denoting the city edge name in mutations.
	EdgeCity = "city"
	// EdgeMaterial holds the string denoting the material edge name in mutations.
	EdgeMaterial = "material"
	// EdgeValues holds the string denoting the values edge name in mutations.
	EdgeValues = "values"
	// EdgeScrapDetails holds the string denoting the scrap_details edge name in mutations.
	EdgeScrapDetails = "scrap_details"
	// EdgeTransferDetails holds the string denoting the transfer_details edge name in mutations.
	EdgeTransferDetails = "transfer_details"
	// EdgeMaintenanceDetails holds the string denoting the maintenance_details edge name in mutations.
	EdgeMaintenanceDetails = "maintenance_details"
	// EdgeCheckDetails holds the string denoting the check_details edge name in mutations.
	EdgeCheckDetails = "check_details"
	// EdgeSubscribe holds the string denoting the subscribe edge name in mutations.
	EdgeSubscribe = "subscribe"
	// EdgeWarehouse holds the string denoting the warehouse edge name in mutations.
	EdgeWarehouse = "warehouse"
	// EdgeStore holds the string denoting the store edge name in mutations.
	EdgeStore = "store"
	// EdgeCabinet holds the string denoting the cabinet edge name in mutations.
	EdgeCabinet = "cabinet"
	// EdgeStation holds the string denoting the station edge name in mutations.
	EdgeStation = "station"
	// EdgeRider holds the string denoting the rider edge name in mutations.
	EdgeRider = "rider"
	// EdgeOperator holds the string denoting the operator edge name in mutations.
	EdgeOperator = "operator"
	// EdgeEbikeAllocates holds the string denoting the ebike_allocates edge name in mutations.
	EdgeEbikeAllocates = "ebike_allocates"
	// EdgeBatteryAllocates holds the string denoting the battery_allocates edge name in mutations.
	EdgeBatteryAllocates = "battery_allocates"
	// EdgeRtoRider holds the string denoting the rto_rider edge name in mutations.
	EdgeRtoRider = "rto_rider"
	// EdgeBatteryRider holds the string denoting the battery_rider edge name in mutations.
	EdgeBatteryRider = "battery_rider"
	// EdgeRentStore holds the string denoting the rent_store edge name in mutations.
	EdgeRentStore = "rent_store"
	// EdgeRentStation holds the string denoting the rent_station edge name in mutations.
	EdgeRentStation = "rent_station"
	// Table holds the table name of the asset in the database.
	Table = "asset"
	// BrandTable is the table that holds the brand relation/edge.
	BrandTable = "asset"
	// BrandInverseTable is the table name for the EbikeBrand entity.
	// It exists in this package in order to avoid circular dependency with the "ebikebrand" package.
	BrandInverseTable = "ebike_brand"
	// BrandColumn is the table column denoting the brand relation/edge.
	BrandColumn = "brand_id"
	// ModelTable is the table that holds the model relation/edge.
	ModelTable = "asset"
	// ModelInverseTable is the table name for the BatteryModel entity.
	// It exists in this package in order to avoid circular dependency with the "batterymodel" package.
	ModelInverseTable = "battery_model"
	// ModelColumn is the table column denoting the model relation/edge.
	ModelColumn = "model_id"
	// CityTable is the table that holds the city relation/edge.
	CityTable = "asset"
	// CityInverseTable is the table name for the City entity.
	// It exists in this package in order to avoid circular dependency with the "city" package.
	CityInverseTable = "city"
	// CityColumn is the table column denoting the city relation/edge.
	CityColumn = "city_id"
	// MaterialTable is the table that holds the material relation/edge.
	MaterialTable = "asset"
	// MaterialInverseTable is the table name for the Material entity.
	// It exists in this package in order to avoid circular dependency with the "material" package.
	MaterialInverseTable = "material"
	// MaterialColumn is the table column denoting the material relation/edge.
	MaterialColumn = "material_id"
	// ValuesTable is the table that holds the values relation/edge.
	ValuesTable = "asset_attribute_values"
	// ValuesInverseTable is the table name for the AssetAttributeValues entity.
	// It exists in this package in order to avoid circular dependency with the "assetattributevalues" package.
	ValuesInverseTable = "asset_attribute_values"
	// ValuesColumn is the table column denoting the values relation/edge.
	ValuesColumn = "asset_id"
	// ScrapDetailsTable is the table that holds the scrap_details relation/edge.
	ScrapDetailsTable = "asset_scrap_details"
	// ScrapDetailsInverseTable is the table name for the AssetScrapDetails entity.
	// It exists in this package in order to avoid circular dependency with the "assetscrapdetails" package.
	ScrapDetailsInverseTable = "asset_scrap_details"
	// ScrapDetailsColumn is the table column denoting the scrap_details relation/edge.
	ScrapDetailsColumn = "asset_id"
	// TransferDetailsTable is the table that holds the transfer_details relation/edge.
	TransferDetailsTable = "asset_transfer_details"
	// TransferDetailsInverseTable is the table name for the AssetTransferDetails entity.
	// It exists in this package in order to avoid circular dependency with the "assettransferdetails" package.
	TransferDetailsInverseTable = "asset_transfer_details"
	// TransferDetailsColumn is the table column denoting the transfer_details relation/edge.
	TransferDetailsColumn = "asset_id"
	// MaintenanceDetailsTable is the table that holds the maintenance_details relation/edge.
	MaintenanceDetailsTable = "asset_maintenance_details"
	// MaintenanceDetailsInverseTable is the table name for the AssetMaintenanceDetails entity.
	// It exists in this package in order to avoid circular dependency with the "assetmaintenancedetails" package.
	MaintenanceDetailsInverseTable = "asset_maintenance_details"
	// MaintenanceDetailsColumn is the table column denoting the maintenance_details relation/edge.
	MaintenanceDetailsColumn = "asset_id"
	// CheckDetailsTable is the table that holds the check_details relation/edge.
	CheckDetailsTable = "asset_check_details"
	// CheckDetailsInverseTable is the table name for the AssetCheckDetails entity.
	// It exists in this package in order to avoid circular dependency with the "assetcheckdetails" package.
	CheckDetailsInverseTable = "asset_check_details"
	// CheckDetailsColumn is the table column denoting the check_details relation/edge.
	CheckDetailsColumn = "asset_id"
	// SubscribeTable is the table that holds the subscribe relation/edge.
	SubscribeTable = "asset"
	// SubscribeInverseTable is the table name for the Subscribe entity.
	// It exists in this package in order to avoid circular dependency with the "subscribe" package.
	SubscribeInverseTable = "subscribe"
	// SubscribeColumn is the table column denoting the subscribe relation/edge.
	SubscribeColumn = "subscribe_id"
	// WarehouseTable is the table that holds the warehouse relation/edge.
	WarehouseTable = "asset"
	// WarehouseInverseTable is the table name for the Warehouse entity.
	// It exists in this package in order to avoid circular dependency with the "warehouse" package.
	WarehouseInverseTable = "warehouse"
	// WarehouseColumn is the table column denoting the warehouse relation/edge.
	WarehouseColumn = "locations_id"
	// StoreTable is the table that holds the store relation/edge.
	StoreTable = "asset"
	// StoreInverseTable is the table name for the Store entity.
	// It exists in this package in order to avoid circular dependency with the "store" package.
	StoreInverseTable = "store"
	// StoreColumn is the table column denoting the store relation/edge.
	StoreColumn = "locations_id"
	// CabinetTable is the table that holds the cabinet relation/edge.
	CabinetTable = "asset"
	// CabinetInverseTable is the table name for the Cabinet entity.
	// It exists in this package in order to avoid circular dependency with the "cabinet" package.
	CabinetInverseTable = "cabinet"
	// CabinetColumn is the table column denoting the cabinet relation/edge.
	CabinetColumn = "locations_id"
	// StationTable is the table that holds the station relation/edge.
	StationTable = "asset"
	// StationInverseTable is the table name for the EnterpriseStation entity.
	// It exists in this package in order to avoid circular dependency with the "enterprisestation" package.
	StationInverseTable = "enterprise_station"
	// StationColumn is the table column denoting the station relation/edge.
	StationColumn = "locations_id"
	// RiderTable is the table that holds the rider relation/edge.
	RiderTable = "asset"
	// RiderInverseTable is the table name for the Rider entity.
	// It exists in this package in order to avoid circular dependency with the "rider" package.
	RiderInverseTable = "rider"
	// RiderColumn is the table column denoting the rider relation/edge.
	RiderColumn = "locations_id"
	// OperatorTable is the table that holds the operator relation/edge.
	OperatorTable = "asset"
	// OperatorInverseTable is the table name for the Maintainer entity.
	// It exists in this package in order to avoid circular dependency with the "maintainer" package.
	OperatorInverseTable = "maintainer"
	// OperatorColumn is the table column denoting the operator relation/edge.
	OperatorColumn = "locations_id"
	// EbikeAllocatesTable is the table that holds the ebike_allocates relation/edge.
	EbikeAllocatesTable = "allocate"
	// EbikeAllocatesInverseTable is the table name for the Allocate entity.
	// It exists in this package in order to avoid circular dependency with the "allocate" package.
	EbikeAllocatesInverseTable = "allocate"
	// EbikeAllocatesColumn is the table column denoting the ebike_allocates relation/edge.
	EbikeAllocatesColumn = "ebike_id"
	// BatteryAllocatesTable is the table that holds the battery_allocates relation/edge.
	BatteryAllocatesTable = "allocate"
	// BatteryAllocatesInverseTable is the table name for the Allocate entity.
	// It exists in this package in order to avoid circular dependency with the "allocate" package.
	BatteryAllocatesInverseTable = "allocate"
	// BatteryAllocatesColumn is the table column denoting the battery_allocates relation/edge.
	BatteryAllocatesColumn = "battery_id"
	// RtoRiderTable is the table that holds the rto_rider relation/edge.
	RtoRiderTable = "asset"
	// RtoRiderInverseTable is the table name for the Rider entity.
	// It exists in this package in order to avoid circular dependency with the "rider" package.
	RtoRiderInverseTable = "rider"
	// RtoRiderColumn is the table column denoting the rto_rider relation/edge.
	RtoRiderColumn = "rto_rider_id"
	// BatteryRiderTable is the table that holds the battery_rider relation/edge.
	BatteryRiderTable = "asset"
	// BatteryRiderInverseTable is the table name for the Rider entity.
	// It exists in this package in order to avoid circular dependency with the "rider" package.
	BatteryRiderInverseTable = "rider"
	// BatteryRiderColumn is the table column denoting the battery_rider relation/edge.
	BatteryRiderColumn = "locations_id"
	// RentStoreTable is the table that holds the rent_store relation/edge.
	RentStoreTable = "asset"
	// RentStoreInverseTable is the table name for the Store entity.
	// It exists in this package in order to avoid circular dependency with the "store" package.
	RentStoreInverseTable = "store"
	// RentStoreColumn is the table column denoting the rent_store relation/edge.
	RentStoreColumn = "rent_locations_id"
	// RentStationTable is the table that holds the rent_station relation/edge.
	RentStationTable = "asset"
	// RentStationInverseTable is the table name for the EnterpriseStation entity.
	// It exists in this package in order to avoid circular dependency with the "enterprisestation" package.
	RentStationInverseTable = "enterprise_station"
	// RentStationColumn is the table column denoting the rent_station relation/edge.
	RentStationColumn = "rent_locations_id"
)

// Columns holds all SQL columns for asset fields.
var Columns = []string{
	FieldID,
	FieldCreatedAt,
	FieldUpdatedAt,
	FieldDeletedAt,
	FieldCreator,
	FieldLastModifier,
	FieldRemark,
	FieldBrandID,
	FieldModelID,
	FieldCityID,
	FieldMaterialID,
	FieldType,
	FieldName,
	FieldSn,
	FieldStatus,
	FieldEnable,
	FieldLocationsType,
	FieldLocationsID,
	FieldRtoRiderID,
	FieldCheckAt,
	FieldBrandName,
	FieldSubscribeID,
	FieldOrdinal,
	FieldRentLocationsType,
	FieldRentLocationsID,
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
	// DefaultStatus holds the default value on creation for the "status" field.
	DefaultStatus uint8
	// DefaultEnable holds the default value on creation for the "enable" field.
	DefaultEnable bool
)

// OrderOption defines the ordering options for the Asset queries.
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

// ByBrandID orders the results by the brand_id field.
func ByBrandID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldBrandID, opts...).ToFunc()
}

// ByModelID orders the results by the model_id field.
func ByModelID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldModelID, opts...).ToFunc()
}

// ByCityID orders the results by the city_id field.
func ByCityID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldCityID, opts...).ToFunc()
}

// ByMaterialID orders the results by the material_id field.
func ByMaterialID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldMaterialID, opts...).ToFunc()
}

// ByType orders the results by the type field.
func ByType(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldType, opts...).ToFunc()
}

// ByName orders the results by the name field.
func ByName(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldName, opts...).ToFunc()
}

// BySn orders the results by the sn field.
func BySn(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldSn, opts...).ToFunc()
}

// ByStatus orders the results by the status field.
func ByStatus(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldStatus, opts...).ToFunc()
}

// ByEnable orders the results by the enable field.
func ByEnable(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldEnable, opts...).ToFunc()
}

// ByLocationsType orders the results by the locations_type field.
func ByLocationsType(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldLocationsType, opts...).ToFunc()
}

// ByLocationsID orders the results by the locations_id field.
func ByLocationsID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldLocationsID, opts...).ToFunc()
}

// ByRtoRiderID orders the results by the rto_rider_id field.
func ByRtoRiderID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldRtoRiderID, opts...).ToFunc()
}

// ByCheckAt orders the results by the check_at field.
func ByCheckAt(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldCheckAt, opts...).ToFunc()
}

// ByBrandName orders the results by the brand_name field.
func ByBrandName(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldBrandName, opts...).ToFunc()
}

// BySubscribeID orders the results by the subscribe_id field.
func BySubscribeID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldSubscribeID, opts...).ToFunc()
}

// ByOrdinal orders the results by the ordinal field.
func ByOrdinal(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldOrdinal, opts...).ToFunc()
}

// ByRentLocationsType orders the results by the rent_locations_type field.
func ByRentLocationsType(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldRentLocationsType, opts...).ToFunc()
}

// ByRentLocationsID orders the results by the rent_locations_id field.
func ByRentLocationsID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldRentLocationsID, opts...).ToFunc()
}

// ByBrandField orders the results by brand field.
func ByBrandField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newBrandStep(), sql.OrderByField(field, opts...))
	}
}

// ByModelField orders the results by model field.
func ByModelField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newModelStep(), sql.OrderByField(field, opts...))
	}
}

// ByCityField orders the results by city field.
func ByCityField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newCityStep(), sql.OrderByField(field, opts...))
	}
}

// ByMaterialField orders the results by material field.
func ByMaterialField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newMaterialStep(), sql.OrderByField(field, opts...))
	}
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

// ByTransferDetailsCount orders the results by transfer_details count.
func ByTransferDetailsCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newTransferDetailsStep(), opts...)
	}
}

// ByTransferDetails orders the results by transfer_details terms.
func ByTransferDetails(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newTransferDetailsStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}

// ByMaintenanceDetailsCount orders the results by maintenance_details count.
func ByMaintenanceDetailsCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newMaintenanceDetailsStep(), opts...)
	}
}

// ByMaintenanceDetails orders the results by maintenance_details terms.
func ByMaintenanceDetails(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newMaintenanceDetailsStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
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

// BySubscribeField orders the results by subscribe field.
func BySubscribeField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newSubscribeStep(), sql.OrderByField(field, opts...))
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

// ByCabinetField orders the results by cabinet field.
func ByCabinetField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newCabinetStep(), sql.OrderByField(field, opts...))
	}
}

// ByStationField orders the results by station field.
func ByStationField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newStationStep(), sql.OrderByField(field, opts...))
	}
}

// ByRiderField orders the results by rider field.
func ByRiderField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newRiderStep(), sql.OrderByField(field, opts...))
	}
}

// ByOperatorField orders the results by operator field.
func ByOperatorField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newOperatorStep(), sql.OrderByField(field, opts...))
	}
}

// ByEbikeAllocatesCount orders the results by ebike_allocates count.
func ByEbikeAllocatesCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newEbikeAllocatesStep(), opts...)
	}
}

// ByEbikeAllocates orders the results by ebike_allocates terms.
func ByEbikeAllocates(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newEbikeAllocatesStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}

// ByBatteryAllocatesCount orders the results by battery_allocates count.
func ByBatteryAllocatesCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newBatteryAllocatesStep(), opts...)
	}
}

// ByBatteryAllocates orders the results by battery_allocates terms.
func ByBatteryAllocates(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newBatteryAllocatesStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}

// ByRtoRiderField orders the results by rto_rider field.
func ByRtoRiderField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newRtoRiderStep(), sql.OrderByField(field, opts...))
	}
}

// ByBatteryRiderField orders the results by battery_rider field.
func ByBatteryRiderField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newBatteryRiderStep(), sql.OrderByField(field, opts...))
	}
}

// ByRentStoreField orders the results by rent_store field.
func ByRentStoreField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newRentStoreStep(), sql.OrderByField(field, opts...))
	}
}

// ByRentStationField orders the results by rent_station field.
func ByRentStationField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newRentStationStep(), sql.OrderByField(field, opts...))
	}
}
func newBrandStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(BrandInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2O, false, BrandTable, BrandColumn),
	)
}
func newModelStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(ModelInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2O, false, ModelTable, ModelColumn),
	)
}
func newCityStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(CityInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2O, false, CityTable, CityColumn),
	)
}
func newMaterialStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(MaterialInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2O, false, MaterialTable, MaterialColumn),
	)
}
func newValuesStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(ValuesInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, ValuesTable, ValuesColumn),
	)
}
func newScrapDetailsStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(ScrapDetailsInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, ScrapDetailsTable, ScrapDetailsColumn),
	)
}
func newTransferDetailsStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(TransferDetailsInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, TransferDetailsTable, TransferDetailsColumn),
	)
}
func newMaintenanceDetailsStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(MaintenanceDetailsInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, MaintenanceDetailsTable, MaintenanceDetailsColumn),
	)
}
func newCheckDetailsStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(CheckDetailsInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, CheckDetailsTable, CheckDetailsColumn),
	)
}
func newSubscribeStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(SubscribeInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, SubscribeTable, SubscribeColumn),
	)
}
func newWarehouseStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(WarehouseInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, WarehouseTable, WarehouseColumn),
	)
}
func newStoreStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(StoreInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, StoreTable, StoreColumn),
	)
}
func newCabinetStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(CabinetInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, CabinetTable, CabinetColumn),
	)
}
func newStationStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(StationInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, StationTable, StationColumn),
	)
}
func newRiderStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(RiderInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, RiderTable, RiderColumn),
	)
}
func newOperatorStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(OperatorInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, OperatorTable, OperatorColumn),
	)
}
func newEbikeAllocatesStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(EbikeAllocatesInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, EbikeAllocatesTable, EbikeAllocatesColumn),
	)
}
func newBatteryAllocatesStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(BatteryAllocatesInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, BatteryAllocatesTable, BatteryAllocatesColumn),
	)
}
func newRtoRiderStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(RtoRiderInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2O, false, RtoRiderTable, RtoRiderColumn),
	)
}
func newBatteryRiderStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(BatteryRiderInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.O2O, true, BatteryRiderTable, BatteryRiderColumn),
	)
}
func newRentStoreStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(RentStoreInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, RentStoreTable, RentStoreColumn),
	)
}
func newRentStationStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(RentStationInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, RentStationTable, RentStationColumn),
	)
}
