// Code generated by entc, DO NOT EDIT.

package business

import (
	"fmt"
	"time"

	"entgo.io/ent"
)

const (
	// Label holds the string label denoting the business type in the database.
	Label = "business"
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
	// FieldRiderID holds the string denoting the rider_id field in the database.
	FieldRiderID = "rider_id"
	// FieldCityID holds the string denoting the city_id field in the database.
	FieldCityID = "city_id"
	// FieldSubscribeID holds the string denoting the subscribe_id field in the database.
	FieldSubscribeID = "subscribe_id"
	// FieldEmployeeID holds the string denoting the employee_id field in the database.
	FieldEmployeeID = "employee_id"
	// FieldStoreID holds the string denoting the store_id field in the database.
	FieldStoreID = "store_id"
	// FieldPlanID holds the string denoting the plan_id field in the database.
	FieldPlanID = "plan_id"
	// FieldEnterpriseID holds the string denoting the enterprise_id field in the database.
	FieldEnterpriseID = "enterprise_id"
	// FieldStationID holds the string denoting the station_id field in the database.
	FieldStationID = "station_id"
	// FieldType holds the string denoting the type field in the database.
	FieldType = "type"
	// EdgeRider holds the string denoting the rider edge name in mutations.
	EdgeRider = "rider"
	// EdgeCity holds the string denoting the city edge name in mutations.
	EdgeCity = "city"
	// EdgeSubscribe holds the string denoting the subscribe edge name in mutations.
	EdgeSubscribe = "subscribe"
	// EdgeEmployee holds the string denoting the employee edge name in mutations.
	EdgeEmployee = "employee"
	// EdgeStore holds the string denoting the store edge name in mutations.
	EdgeStore = "store"
	// EdgePlan holds the string denoting the plan edge name in mutations.
	EdgePlan = "plan"
	// EdgeEnterprise holds the string denoting the enterprise edge name in mutations.
	EdgeEnterprise = "enterprise"
	// EdgeStation holds the string denoting the station edge name in mutations.
	EdgeStation = "station"
	// Table holds the table name of the business in the database.
	Table = "business"
	// RiderTable is the table that holds the rider relation/edge.
	RiderTable = "business"
	// RiderInverseTable is the table name for the Rider entity.
	// It exists in this package in order to avoid circular dependency with the "rider" package.
	RiderInverseTable = "rider"
	// RiderColumn is the table column denoting the rider relation/edge.
	RiderColumn = "rider_id"
	// CityTable is the table that holds the city relation/edge.
	CityTable = "business"
	// CityInverseTable is the table name for the City entity.
	// It exists in this package in order to avoid circular dependency with the "city" package.
	CityInverseTable = "city"
	// CityColumn is the table column denoting the city relation/edge.
	CityColumn = "city_id"
	// SubscribeTable is the table that holds the subscribe relation/edge.
	SubscribeTable = "business"
	// SubscribeInverseTable is the table name for the Subscribe entity.
	// It exists in this package in order to avoid circular dependency with the "subscribe" package.
	SubscribeInverseTable = "subscribe"
	// SubscribeColumn is the table column denoting the subscribe relation/edge.
	SubscribeColumn = "subscribe_id"
	// EmployeeTable is the table that holds the employee relation/edge.
	EmployeeTable = "business"
	// EmployeeInverseTable is the table name for the Employee entity.
	// It exists in this package in order to avoid circular dependency with the "employee" package.
	EmployeeInverseTable = "employee"
	// EmployeeColumn is the table column denoting the employee relation/edge.
	EmployeeColumn = "employee_id"
	// StoreTable is the table that holds the store relation/edge.
	StoreTable = "business"
	// StoreInverseTable is the table name for the Store entity.
	// It exists in this package in order to avoid circular dependency with the "store" package.
	StoreInverseTable = "store"
	// StoreColumn is the table column denoting the store relation/edge.
	StoreColumn = "store_id"
	// PlanTable is the table that holds the plan relation/edge.
	PlanTable = "business"
	// PlanInverseTable is the table name for the Plan entity.
	// It exists in this package in order to avoid circular dependency with the "plan" package.
	PlanInverseTable = "plan"
	// PlanColumn is the table column denoting the plan relation/edge.
	PlanColumn = "plan_id"
	// EnterpriseTable is the table that holds the enterprise relation/edge.
	EnterpriseTable = "business"
	// EnterpriseInverseTable is the table name for the Enterprise entity.
	// It exists in this package in order to avoid circular dependency with the "enterprise" package.
	EnterpriseInverseTable = "enterprise"
	// EnterpriseColumn is the table column denoting the enterprise relation/edge.
	EnterpriseColumn = "enterprise_id"
	// StationTable is the table that holds the station relation/edge.
	StationTable = "business"
	// StationInverseTable is the table name for the EnterpriseStation entity.
	// It exists in this package in order to avoid circular dependency with the "enterprisestation" package.
	StationInverseTable = "enterprise_station"
	// StationColumn is the table column denoting the station relation/edge.
	StationColumn = "station_id"
)

// Columns holds all SQL columns for business fields.
var Columns = []string{
	FieldID,
	FieldCreatedAt,
	FieldUpdatedAt,
	FieldDeletedAt,
	FieldCreator,
	FieldLastModifier,
	FieldRemark,
	FieldRiderID,
	FieldCityID,
	FieldSubscribeID,
	FieldEmployeeID,
	FieldStoreID,
	FieldPlanID,
	FieldEnterpriseID,
	FieldStationID,
	FieldType,
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
//
var (
	Hooks [1]ent.Hook
	// DefaultCreatedAt holds the default value on creation for the "created_at" field.
	DefaultCreatedAt func() time.Time
	// DefaultUpdatedAt holds the default value on creation for the "updated_at" field.
	DefaultUpdatedAt func() time.Time
	// UpdateDefaultUpdatedAt holds the default value on update for the "updated_at" field.
	UpdateDefaultUpdatedAt func() time.Time
)

// Type defines the type for the "type" enum field.
type Type string

// Type values.
const (
	TypeActive      Type = "active"
	TypePause       Type = "pause"
	TypeContinue    Type = "continue"
	TypeUnsubscribe Type = "unsubscribe"
)

func (_type Type) String() string {
	return string(_type)
}

// TypeValidator is a validator for the "type" field enum values. It is called by the builders before save.
func TypeValidator(_type Type) error {
	switch _type {
	case TypeActive, TypePause, TypeContinue, TypeUnsubscribe:
		return nil
	default:
		return fmt.Errorf("business: invalid enum value for type field: %q", _type)
	}
}