// Code generated by ent, DO NOT EDIT.

package subscribereminder

import (
	"fmt"
	"time"
)

const (
	// Label holds the string label denoting the subscribereminder type in the database.
	Label = "subscribe_reminder"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldCreatedAt holds the string denoting the created_at field in the database.
	FieldCreatedAt = "created_at"
	// FieldUpdatedAt holds the string denoting the updated_at field in the database.
	FieldUpdatedAt = "updated_at"
	// FieldSubscribeID holds the string denoting the subscribe_id field in the database.
	FieldSubscribeID = "subscribe_id"
	// FieldPlanID holds the string denoting the plan_id field in the database.
	FieldPlanID = "plan_id"
	// FieldType holds the string denoting the type field in the database.
	FieldType = "type"
	// FieldPhone holds the string denoting the phone field in the database.
	FieldPhone = "phone"
	// FieldName holds the string denoting the name field in the database.
	FieldName = "name"
	// FieldSuccess holds the string denoting the success field in the database.
	FieldSuccess = "success"
	// FieldDays holds the string denoting the days field in the database.
	FieldDays = "days"
	// FieldPlanName holds the string denoting the plan_name field in the database.
	FieldPlanName = "plan_name"
	// FieldDate holds the string denoting the date field in the database.
	FieldDate = "date"
	// FieldFee holds the string denoting the fee field in the database.
	FieldFee = "fee"
	// FieldFeeFormula holds the string denoting the fee_formula field in the database.
	FieldFeeFormula = "fee_formula"
	// EdgeSubscribe holds the string denoting the subscribe edge name in mutations.
	EdgeSubscribe = "subscribe"
	// EdgePlan holds the string denoting the plan edge name in mutations.
	EdgePlan = "plan"
	// Table holds the table name of the subscribereminder in the database.
	Table = "subscribe_reminder"
	// SubscribeTable is the table that holds the subscribe relation/edge.
	SubscribeTable = "subscribe_reminder"
	// SubscribeInverseTable is the table name for the Subscribe entity.
	// It exists in this package in order to avoid circular dependency with the "subscribe" package.
	SubscribeInverseTable = "subscribe"
	// SubscribeColumn is the table column denoting the subscribe relation/edge.
	SubscribeColumn = "subscribe_id"
	// PlanTable is the table that holds the plan relation/edge.
	PlanTable = "subscribe_reminder"
	// PlanInverseTable is the table name for the Plan entity.
	// It exists in this package in order to avoid circular dependency with the "plan" package.
	PlanInverseTable = "plan"
	// PlanColumn is the table column denoting the plan relation/edge.
	PlanColumn = "plan_id"
)

// Columns holds all SQL columns for subscribereminder fields.
var Columns = []string{
	FieldID,
	FieldCreatedAt,
	FieldUpdatedAt,
	FieldSubscribeID,
	FieldPlanID,
	FieldType,
	FieldPhone,
	FieldName,
	FieldSuccess,
	FieldDays,
	FieldPlanName,
	FieldDate,
	FieldFee,
	FieldFeeFormula,
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
	// DefaultFee holds the default value on creation for the "fee" field.
	DefaultFee float64
)

// Type defines the type for the "type" enum field.
type Type string

// Type values.
const (
	TypeSms Type = "sms"
	TypeVms Type = "vms"
)

func (_type Type) String() string {
	return string(_type)
}

// TypeValidator is a validator for the "type" field enum values. It is called by the builders before save.
func TypeValidator(_type Type) error {
	switch _type {
	case TypeSms, TypeVms:
		return nil
	default:
		return fmt.Errorf("subscribereminder: invalid enum value for type field: %q", _type)
	}
}