// Code generated by entc, DO NOT EDIT.

package exception

import (
	"time"

	"entgo.io/ent"
)

const (
	// Label holds the string label denoting the exception type in the database.
	Label = "exception"
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
	// FieldEmployeeID holds the string denoting the employee_id field in the database.
	FieldEmployeeID = "employee_id"
	// FieldStatus holds the string denoting the status field in the database.
	FieldStatus = "status"
	// FieldStoreID holds the string denoting the store_id field in the database.
	FieldStoreID = "store_id"
	// FieldName holds the string denoting the name field in the database.
	FieldName = "name"
	// FieldVoltage holds the string denoting the voltage field in the database.
	FieldVoltage = "voltage"
	// FieldNum holds the string denoting the num field in the database.
	FieldNum = "num"
	// FieldReason holds the string denoting the reason field in the database.
	FieldReason = "reason"
	// FieldDescription holds the string denoting the description field in the database.
	FieldDescription = "description"
	// FieldAttachments holds the string denoting the attachments field in the database.
	FieldAttachments = "attachments"
	// EdgeCity holds the string denoting the city edge name in mutations.
	EdgeCity = "city"
	// EdgeEmployee holds the string denoting the employee edge name in mutations.
	EdgeEmployee = "employee"
	// EdgeStore holds the string denoting the store edge name in mutations.
	EdgeStore = "store"
	// Table holds the table name of the exception in the database.
	Table = "exception"
	// CityTable is the table that holds the city relation/edge.
	CityTable = "exception"
	// CityInverseTable is the table name for the City entity.
	// It exists in this package in order to avoid circular dependency with the "city" package.
	CityInverseTable = "city"
	// CityColumn is the table column denoting the city relation/edge.
	CityColumn = "city_id"
	// EmployeeTable is the table that holds the employee relation/edge.
	EmployeeTable = "exception"
	// EmployeeInverseTable is the table name for the Employee entity.
	// It exists in this package in order to avoid circular dependency with the "employee" package.
	EmployeeInverseTable = "employee"
	// EmployeeColumn is the table column denoting the employee relation/edge.
	EmployeeColumn = "employee_id"
	// StoreTable is the table that holds the store relation/edge.
	StoreTable = "exception"
	// StoreInverseTable is the table name for the Store entity.
	// It exists in this package in order to avoid circular dependency with the "store" package.
	StoreInverseTable = "store"
	// StoreColumn is the table column denoting the store relation/edge.
	StoreColumn = "store_id"
)

// Columns holds all SQL columns for exception fields.
var Columns = []string{
	FieldID,
	FieldCreatedAt,
	FieldUpdatedAt,
	FieldDeletedAt,
	FieldCreator,
	FieldLastModifier,
	FieldRemark,
	FieldCityID,
	FieldEmployeeID,
	FieldStatus,
	FieldStoreID,
	FieldName,
	FieldVoltage,
	FieldNum,
	FieldReason,
	FieldDescription,
	FieldAttachments,
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
	// DefaultStatus holds the default value on creation for the "status" field.
	DefaultStatus uint8
)