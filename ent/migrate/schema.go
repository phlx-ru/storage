// Code generated by ent, DO NOT EDIT.

package migrate

import (
	"entgo.io/ent/dialect/sql/schema"
	"entgo.io/ent/schema/field"
)

var (
	// FilesColumns holds the columns for the "files" table.
	FilesColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "uid", Type: field.TypeUUID},
		{Name: "user_id", Type: field.TypeInt},
		{Name: "filename", Type: field.TypeString},
		{Name: "object_path", Type: field.TypeString},
		{Name: "size", Type: field.TypeInt},
		{Name: "mime_type", Type: field.TypeString},
		{Name: "created_at", Type: field.TypeTime},
		{Name: "updated_at", Type: field.TypeTime, Default: "CURRENT_TIMESTAMP"},
		{Name: "deleted_at", Type: field.TypeTime, Nullable: true},
	}
	// FilesTable holds the schema information for the "files" table.
	FilesTable = &schema.Table{
		Name:       "files",
		Columns:    FilesColumns,
		PrimaryKey: []*schema.Column{FilesColumns[0]},
		Indexes: []*schema.Index{
			{
				Name:    "file_uid",
				Unique:  true,
				Columns: []*schema.Column{FilesColumns[1]},
			},
			{
				Name:    "file_user_id",
				Unique:  false,
				Columns: []*schema.Column{FilesColumns[2]},
			},
			{
				Name:    "file_deleted_at",
				Unique:  false,
				Columns: []*schema.Column{FilesColumns[9]},
			},
			{
				Name:    "file_filename",
				Unique:  false,
				Columns: []*schema.Column{FilesColumns[3]},
			},
			{
				Name:    "file_object_path",
				Unique:  false,
				Columns: []*schema.Column{FilesColumns[4]},
			},
		},
	}
	// Tables holds all the tables in the schema.
	Tables = []*schema.Table{
		FilesTable,
	}
)

func init() {
}