// Code generated by entc, DO NOT EDIT.

package export

import (
	"time"
)

const (
	// Label holds the string label denoting the export type in the database.
	Label = "export"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldCreatedAt holds the string denoting the created_at field in the database.
	FieldCreatedAt = "created_at"
	// FieldUpdatedAt holds the string denoting the updated_at field in the database.
	FieldUpdatedAt = "updated_at"
	// FieldDeletedAt holds the string denoting the deleted_at field in the database.
	FieldDeletedAt = "deleted_at"
	// FieldManagerID holds the string denoting the manager_id field in the database.
	FieldManagerID = "manager_id"
	// FieldTaxonomy holds the string denoting the taxonomy field in the database.
	FieldTaxonomy = "taxonomy"
	// FieldSn holds the string denoting the sn field in the database.
	FieldSn = "sn"
	// FieldStatus holds the string denoting the status field in the database.
	FieldStatus = "status"
	// FieldPath holds the string denoting the path field in the database.
	FieldPath = "path"
	// FieldMessage holds the string denoting the message field in the database.
	FieldMessage = "message"
	// FieldFinishAt holds the string denoting the finish_at field in the database.
	FieldFinishAt = "finish_at"
	// FieldDuration holds the string denoting the duration field in the database.
	FieldDuration = "duration"
	// FieldCondition holds the string denoting the condition field in the database.
	FieldCondition = "condition"
	// FieldInfo holds the string denoting the info field in the database.
	FieldInfo = "info"
	// FieldRemark holds the string denoting the remark field in the database.
	FieldRemark = "remark"
	// EdgeManager holds the string denoting the manager edge name in mutations.
	EdgeManager = "manager"
	// Table holds the table name of the export in the database.
	Table = "export"
	// ManagerTable is the table that holds the manager relation/edge.
	ManagerTable = "export"
	// ManagerInverseTable is the table name for the Manager entity.
	// It exists in this package in order to avoid circular dependency with the "manager" package.
	ManagerInverseTable = "manager"
	// ManagerColumn is the table column denoting the manager relation/edge.
	ManagerColumn = "manager_id"
)

// Columns holds all SQL columns for export fields.
var Columns = []string{
	FieldID,
	FieldCreatedAt,
	FieldUpdatedAt,
	FieldDeletedAt,
	FieldManagerID,
	FieldTaxonomy,
	FieldSn,
	FieldStatus,
	FieldPath,
	FieldMessage,
	FieldFinishAt,
	FieldDuration,
	FieldCondition,
	FieldInfo,
	FieldRemark,
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
	// DefaultStatus holds the default value on creation for the "status" field.
	DefaultStatus uint8
)