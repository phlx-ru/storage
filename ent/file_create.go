// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"storage/ent/file"
	"time"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// FileCreate is the builder for creating a File entity.
type FileCreate struct {
	config
	mutation *FileMutation
	hooks    []Hook
}

// SetUID sets the "uid" field.
func (fc *FileCreate) SetUID(u uuid.UUID) *FileCreate {
	fc.mutation.SetUID(u)
	return fc
}

// SetNillableUID sets the "uid" field if the given value is not nil.
func (fc *FileCreate) SetNillableUID(u *uuid.UUID) *FileCreate {
	if u != nil {
		fc.SetUID(*u)
	}
	return fc
}

// SetUserID sets the "user_id" field.
func (fc *FileCreate) SetUserID(i int) *FileCreate {
	fc.mutation.SetUserID(i)
	return fc
}

// SetFilename sets the "filename" field.
func (fc *FileCreate) SetFilename(s string) *FileCreate {
	fc.mutation.SetFilename(s)
	return fc
}

// SetObjectPath sets the "object_path" field.
func (fc *FileCreate) SetObjectPath(s string) *FileCreate {
	fc.mutation.SetObjectPath(s)
	return fc
}

// SetSize sets the "size" field.
func (fc *FileCreate) SetSize(i int) *FileCreate {
	fc.mutation.SetSize(i)
	return fc
}

// SetMimeType sets the "mime_type" field.
func (fc *FileCreate) SetMimeType(s string) *FileCreate {
	fc.mutation.SetMimeType(s)
	return fc
}

// SetCreatedAt sets the "created_at" field.
func (fc *FileCreate) SetCreatedAt(t time.Time) *FileCreate {
	fc.mutation.SetCreatedAt(t)
	return fc
}

// SetNillableCreatedAt sets the "created_at" field if the given value is not nil.
func (fc *FileCreate) SetNillableCreatedAt(t *time.Time) *FileCreate {
	if t != nil {
		fc.SetCreatedAt(*t)
	}
	return fc
}

// SetUpdatedAt sets the "updated_at" field.
func (fc *FileCreate) SetUpdatedAt(t time.Time) *FileCreate {
	fc.mutation.SetUpdatedAt(t)
	return fc
}

// SetNillableUpdatedAt sets the "updated_at" field if the given value is not nil.
func (fc *FileCreate) SetNillableUpdatedAt(t *time.Time) *FileCreate {
	if t != nil {
		fc.SetUpdatedAt(*t)
	}
	return fc
}

// SetDeletedAt sets the "deleted_at" field.
func (fc *FileCreate) SetDeletedAt(t time.Time) *FileCreate {
	fc.mutation.SetDeletedAt(t)
	return fc
}

// SetNillableDeletedAt sets the "deleted_at" field if the given value is not nil.
func (fc *FileCreate) SetNillableDeletedAt(t *time.Time) *FileCreate {
	if t != nil {
		fc.SetDeletedAt(*t)
	}
	return fc
}

// Mutation returns the FileMutation object of the builder.
func (fc *FileCreate) Mutation() *FileMutation {
	return fc.mutation
}

// Save creates the File in the database.
func (fc *FileCreate) Save(ctx context.Context) (*File, error) {
	fc.defaults()
	return withHooks[*File, FileMutation](ctx, fc.sqlSave, fc.mutation, fc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (fc *FileCreate) SaveX(ctx context.Context) *File {
	v, err := fc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (fc *FileCreate) Exec(ctx context.Context) error {
	_, err := fc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (fc *FileCreate) ExecX(ctx context.Context) {
	if err := fc.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (fc *FileCreate) defaults() {
	if _, ok := fc.mutation.UID(); !ok {
		v := file.DefaultUID()
		fc.mutation.SetUID(v)
	}
	if _, ok := fc.mutation.CreatedAt(); !ok {
		v := file.DefaultCreatedAt()
		fc.mutation.SetCreatedAt(v)
	}
	if _, ok := fc.mutation.UpdatedAt(); !ok {
		v := file.DefaultUpdatedAt()
		fc.mutation.SetUpdatedAt(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (fc *FileCreate) check() error {
	if _, ok := fc.mutation.UID(); !ok {
		return &ValidationError{Name: "uid", err: errors.New(`ent: missing required field "File.uid"`)}
	}
	if _, ok := fc.mutation.UserID(); !ok {
		return &ValidationError{Name: "user_id", err: errors.New(`ent: missing required field "File.user_id"`)}
	}
	if _, ok := fc.mutation.Filename(); !ok {
		return &ValidationError{Name: "filename", err: errors.New(`ent: missing required field "File.filename"`)}
	}
	if _, ok := fc.mutation.ObjectPath(); !ok {
		return &ValidationError{Name: "object_path", err: errors.New(`ent: missing required field "File.object_path"`)}
	}
	if _, ok := fc.mutation.Size(); !ok {
		return &ValidationError{Name: "size", err: errors.New(`ent: missing required field "File.size"`)}
	}
	if _, ok := fc.mutation.MimeType(); !ok {
		return &ValidationError{Name: "mime_type", err: errors.New(`ent: missing required field "File.mime_type"`)}
	}
	if _, ok := fc.mutation.CreatedAt(); !ok {
		return &ValidationError{Name: "created_at", err: errors.New(`ent: missing required field "File.created_at"`)}
	}
	if _, ok := fc.mutation.UpdatedAt(); !ok {
		return &ValidationError{Name: "updated_at", err: errors.New(`ent: missing required field "File.updated_at"`)}
	}
	return nil
}

func (fc *FileCreate) sqlSave(ctx context.Context) (*File, error) {
	if err := fc.check(); err != nil {
		return nil, err
	}
	_node, _spec := fc.createSpec()
	if err := sqlgraph.CreateNode(ctx, fc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = int(id)
	fc.mutation.id = &_node.ID
	fc.mutation.done = true
	return _node, nil
}

func (fc *FileCreate) createSpec() (*File, *sqlgraph.CreateSpec) {
	var (
		_node = &File{config: fc.config}
		_spec = sqlgraph.NewCreateSpec(file.Table, sqlgraph.NewFieldSpec(file.FieldID, field.TypeInt))
	)
	if value, ok := fc.mutation.UID(); ok {
		_spec.SetField(file.FieldUID, field.TypeUUID, value)
		_node.UID = value
	}
	if value, ok := fc.mutation.UserID(); ok {
		_spec.SetField(file.FieldUserID, field.TypeInt, value)
		_node.UserID = value
	}
	if value, ok := fc.mutation.Filename(); ok {
		_spec.SetField(file.FieldFilename, field.TypeString, value)
		_node.Filename = value
	}
	if value, ok := fc.mutation.ObjectPath(); ok {
		_spec.SetField(file.FieldObjectPath, field.TypeString, value)
		_node.ObjectPath = value
	}
	if value, ok := fc.mutation.Size(); ok {
		_spec.SetField(file.FieldSize, field.TypeInt, value)
		_node.Size = value
	}
	if value, ok := fc.mutation.MimeType(); ok {
		_spec.SetField(file.FieldMimeType, field.TypeString, value)
		_node.MimeType = value
	}
	if value, ok := fc.mutation.CreatedAt(); ok {
		_spec.SetField(file.FieldCreatedAt, field.TypeTime, value)
		_node.CreatedAt = value
	}
	if value, ok := fc.mutation.UpdatedAt(); ok {
		_spec.SetField(file.FieldUpdatedAt, field.TypeTime, value)
		_node.UpdatedAt = value
	}
	if value, ok := fc.mutation.DeletedAt(); ok {
		_spec.SetField(file.FieldDeletedAt, field.TypeTime, value)
		_node.DeletedAt = &value
	}
	return _node, _spec
}

// FileCreateBulk is the builder for creating many File entities in bulk.
type FileCreateBulk struct {
	config
	builders []*FileCreate
}

// Save creates the File entities in the database.
func (fcb *FileCreateBulk) Save(ctx context.Context) ([]*File, error) {
	specs := make([]*sqlgraph.CreateSpec, len(fcb.builders))
	nodes := make([]*File, len(fcb.builders))
	mutators := make([]Mutator, len(fcb.builders))
	for i := range fcb.builders {
		func(i int, root context.Context) {
			builder := fcb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*FileMutation)
				if !ok {
					return nil, fmt.Errorf("unexpected mutation type %T", m)
				}
				if err := builder.check(); err != nil {
					return nil, err
				}
				builder.mutation = mutation
				nodes[i], specs[i] = builder.createSpec()
				var err error
				if i < len(mutators)-1 {
					_, err = mutators[i+1].Mutate(root, fcb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, fcb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{msg: err.Error(), wrap: err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
				if specs[i].ID.Value != nil {
					id := specs[i].ID.Value.(int64)
					nodes[i].ID = int(id)
				}
				mutation.done = true
				return nodes[i], nil
			})
			for i := len(builder.hooks) - 1; i >= 0; i-- {
				mut = builder.hooks[i](mut)
			}
			mutators[i] = mut
		}(i, ctx)
	}
	if len(mutators) > 0 {
		if _, err := mutators[0].Mutate(ctx, fcb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (fcb *FileCreateBulk) SaveX(ctx context.Context) []*File {
	v, err := fcb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (fcb *FileCreateBulk) Exec(ctx context.Context) error {
	_, err := fcb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (fcb *FileCreateBulk) ExecX(ctx context.Context) {
	if err := fcb.Exec(ctx); err != nil {
		panic(err)
	}
}
