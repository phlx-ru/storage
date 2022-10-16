package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

// File holds the schema definition for the File entity.
type File struct {
	ent.Schema
}

// Fields of the File.
func (File) Fields() []ent.Field {
	return []ent.Field{
		field.UUID(`uid`, uuid.UUID{}).
			Default(uuid.New).
			Comment(`unique file identifier`),

		field.Int(`user_id`).
			Comment(`user identification number`),

		field.String(`filename`).
			Comment(`filename of downloaded file`),

		field.String(`object_path`).
			Comment(`path to file object in s3 storage`),

		field.Int(`size`).
			Comment(`size of file in bytes`),

		field.String(`mime_type`).
			Comment(`file mime type`),

		field.Time(`created_at`).
			Default(time.Now).
			Immutable().
			Comment(`creation time of file`),

		field.Time(`updated_at`).
			Default(time.Now).
			UpdateDefault(time.Now).
			Annotations(
				&entsql.Annotation{
					Default: `CURRENT_TIMESTAMP`,
				},
			).
			Comment(`last update time of file`),

		field.Time(`deleted_at`).
			Optional().
			Nillable().
			Default(nil).
			Comment(`time of file deletion`),
	}
}

// Edges of the File.
func (File) Edges() []ent.Edge {
	return nil
}

func (File) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields(`uid`).Unique(),
		index.Fields(`user_id`),
		index.Fields(`deleted_at`),
		index.Fields(`filename`),
		index.Fields(`object_path`),
	}
}
