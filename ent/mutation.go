// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"storage/ent/file"
	"storage/ent/predicate"
	"sync"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
)

const (
	// Operation types.
	OpCreate    = ent.OpCreate
	OpDelete    = ent.OpDelete
	OpDeleteOne = ent.OpDeleteOne
	OpUpdate    = ent.OpUpdate
	OpUpdateOne = ent.OpUpdateOne

	// Node types.
	TypeFile = "File"
)

// FileMutation represents an operation that mutates the File nodes in the graph.
type FileMutation struct {
	config
	op            Op
	typ           string
	id            *int
	uid           *uuid.UUID
	user_id       *int
	adduser_id    *int
	filename      *string
	object_path   *string
	size          *int
	addsize       *int
	mime_type     *string
	created_at    *time.Time
	updated_at    *time.Time
	deleted_at    *time.Time
	clearedFields map[string]struct{}
	done          bool
	oldValue      func(context.Context) (*File, error)
	predicates    []predicate.File
}

var _ ent.Mutation = (*FileMutation)(nil)

// fileOption allows management of the mutation configuration using functional options.
type fileOption func(*FileMutation)

// newFileMutation creates new mutation for the File entity.
func newFileMutation(c config, op Op, opts ...fileOption) *FileMutation {
	m := &FileMutation{
		config:        c,
		op:            op,
		typ:           TypeFile,
		clearedFields: make(map[string]struct{}),
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

// withFileID sets the ID field of the mutation.
func withFileID(id int) fileOption {
	return func(m *FileMutation) {
		var (
			err   error
			once  sync.Once
			value *File
		)
		m.oldValue = func(ctx context.Context) (*File, error) {
			once.Do(func() {
				if m.done {
					err = errors.New("querying old values post mutation is not allowed")
				} else {
					value, err = m.Client().File.Get(ctx, id)
				}
			})
			return value, err
		}
		m.id = &id
	}
}

// withFile sets the old File of the mutation.
func withFile(node *File) fileOption {
	return func(m *FileMutation) {
		m.oldValue = func(context.Context) (*File, error) {
			return node, nil
		}
		m.id = &node.ID
	}
}

// Client returns a new `ent.Client` from the mutation. If the mutation was
// executed in a transaction (ent.Tx), a transactional client is returned.
func (m FileMutation) Client() *Client {
	client := &Client{config: m.config}
	client.init()
	return client
}

// Tx returns an `ent.Tx` for mutations that were executed in transactions;
// it returns an error otherwise.
func (m FileMutation) Tx() (*Tx, error) {
	if _, ok := m.driver.(*txDriver); !ok {
		return nil, errors.New("ent: mutation is not running in a transaction")
	}
	tx := &Tx{config: m.config}
	tx.init()
	return tx, nil
}

// ID returns the ID value in the mutation. Note that the ID is only available
// if it was provided to the builder or after it was returned from the database.
func (m *FileMutation) ID() (id int, exists bool) {
	if m.id == nil {
		return
	}
	return *m.id, true
}

// IDs queries the database and returns the entity ids that match the mutation's predicate.
// That means, if the mutation is applied within a transaction with an isolation level such
// as sql.LevelSerializable, the returned ids match the ids of the rows that will be updated
// or updated by the mutation.
func (m *FileMutation) IDs(ctx context.Context) ([]int, error) {
	switch {
	case m.op.Is(OpUpdateOne | OpDeleteOne):
		id, exists := m.ID()
		if exists {
			return []int{id}, nil
		}
		fallthrough
	case m.op.Is(OpUpdate | OpDelete):
		return m.Client().File.Query().Where(m.predicates...).IDs(ctx)
	default:
		return nil, fmt.Errorf("IDs is not allowed on %s operations", m.op)
	}
}

// SetUID sets the "uid" field.
func (m *FileMutation) SetUID(u uuid.UUID) {
	m.uid = &u
}

// UID returns the value of the "uid" field in the mutation.
func (m *FileMutation) UID() (r uuid.UUID, exists bool) {
	v := m.uid
	if v == nil {
		return
	}
	return *v, true
}

// OldUID returns the old "uid" field's value of the File entity.
// If the File object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *FileMutation) OldUID(ctx context.Context) (v uuid.UUID, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldUID is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldUID requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldUID: %w", err)
	}
	return oldValue.UID, nil
}

// ResetUID resets all changes to the "uid" field.
func (m *FileMutation) ResetUID() {
	m.uid = nil
}

// SetUserID sets the "user_id" field.
func (m *FileMutation) SetUserID(i int) {
	m.user_id = &i
	m.adduser_id = nil
}

// UserID returns the value of the "user_id" field in the mutation.
func (m *FileMutation) UserID() (r int, exists bool) {
	v := m.user_id
	if v == nil {
		return
	}
	return *v, true
}

// OldUserID returns the old "user_id" field's value of the File entity.
// If the File object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *FileMutation) OldUserID(ctx context.Context) (v int, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldUserID is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldUserID requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldUserID: %w", err)
	}
	return oldValue.UserID, nil
}

// AddUserID adds i to the "user_id" field.
func (m *FileMutation) AddUserID(i int) {
	if m.adduser_id != nil {
		*m.adduser_id += i
	} else {
		m.adduser_id = &i
	}
}

// AddedUserID returns the value that was added to the "user_id" field in this mutation.
func (m *FileMutation) AddedUserID() (r int, exists bool) {
	v := m.adduser_id
	if v == nil {
		return
	}
	return *v, true
}

// ResetUserID resets all changes to the "user_id" field.
func (m *FileMutation) ResetUserID() {
	m.user_id = nil
	m.adduser_id = nil
}

// SetFilename sets the "filename" field.
func (m *FileMutation) SetFilename(s string) {
	m.filename = &s
}

// Filename returns the value of the "filename" field in the mutation.
func (m *FileMutation) Filename() (r string, exists bool) {
	v := m.filename
	if v == nil {
		return
	}
	return *v, true
}

// OldFilename returns the old "filename" field's value of the File entity.
// If the File object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *FileMutation) OldFilename(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldFilename is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldFilename requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldFilename: %w", err)
	}
	return oldValue.Filename, nil
}

// ResetFilename resets all changes to the "filename" field.
func (m *FileMutation) ResetFilename() {
	m.filename = nil
}

// SetObjectPath sets the "object_path" field.
func (m *FileMutation) SetObjectPath(s string) {
	m.object_path = &s
}

// ObjectPath returns the value of the "object_path" field in the mutation.
func (m *FileMutation) ObjectPath() (r string, exists bool) {
	v := m.object_path
	if v == nil {
		return
	}
	return *v, true
}

// OldObjectPath returns the old "object_path" field's value of the File entity.
// If the File object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *FileMutation) OldObjectPath(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldObjectPath is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldObjectPath requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldObjectPath: %w", err)
	}
	return oldValue.ObjectPath, nil
}

// ResetObjectPath resets all changes to the "object_path" field.
func (m *FileMutation) ResetObjectPath() {
	m.object_path = nil
}

// SetSize sets the "size" field.
func (m *FileMutation) SetSize(i int) {
	m.size = &i
	m.addsize = nil
}

// Size returns the value of the "size" field in the mutation.
func (m *FileMutation) Size() (r int, exists bool) {
	v := m.size
	if v == nil {
		return
	}
	return *v, true
}

// OldSize returns the old "size" field's value of the File entity.
// If the File object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *FileMutation) OldSize(ctx context.Context) (v int, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldSize is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldSize requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldSize: %w", err)
	}
	return oldValue.Size, nil
}

// AddSize adds i to the "size" field.
func (m *FileMutation) AddSize(i int) {
	if m.addsize != nil {
		*m.addsize += i
	} else {
		m.addsize = &i
	}
}

// AddedSize returns the value that was added to the "size" field in this mutation.
func (m *FileMutation) AddedSize() (r int, exists bool) {
	v := m.addsize
	if v == nil {
		return
	}
	return *v, true
}

// ResetSize resets all changes to the "size" field.
func (m *FileMutation) ResetSize() {
	m.size = nil
	m.addsize = nil
}

// SetMimeType sets the "mime_type" field.
func (m *FileMutation) SetMimeType(s string) {
	m.mime_type = &s
}

// MimeType returns the value of the "mime_type" field in the mutation.
func (m *FileMutation) MimeType() (r string, exists bool) {
	v := m.mime_type
	if v == nil {
		return
	}
	return *v, true
}

// OldMimeType returns the old "mime_type" field's value of the File entity.
// If the File object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *FileMutation) OldMimeType(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldMimeType is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldMimeType requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldMimeType: %w", err)
	}
	return oldValue.MimeType, nil
}

// ResetMimeType resets all changes to the "mime_type" field.
func (m *FileMutation) ResetMimeType() {
	m.mime_type = nil
}

// SetCreatedAt sets the "created_at" field.
func (m *FileMutation) SetCreatedAt(t time.Time) {
	m.created_at = &t
}

// CreatedAt returns the value of the "created_at" field in the mutation.
func (m *FileMutation) CreatedAt() (r time.Time, exists bool) {
	v := m.created_at
	if v == nil {
		return
	}
	return *v, true
}

// OldCreatedAt returns the old "created_at" field's value of the File entity.
// If the File object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *FileMutation) OldCreatedAt(ctx context.Context) (v time.Time, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldCreatedAt is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldCreatedAt requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldCreatedAt: %w", err)
	}
	return oldValue.CreatedAt, nil
}

// ResetCreatedAt resets all changes to the "created_at" field.
func (m *FileMutation) ResetCreatedAt() {
	m.created_at = nil
}

// SetUpdatedAt sets the "updated_at" field.
func (m *FileMutation) SetUpdatedAt(t time.Time) {
	m.updated_at = &t
}

// UpdatedAt returns the value of the "updated_at" field in the mutation.
func (m *FileMutation) UpdatedAt() (r time.Time, exists bool) {
	v := m.updated_at
	if v == nil {
		return
	}
	return *v, true
}

// OldUpdatedAt returns the old "updated_at" field's value of the File entity.
// If the File object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *FileMutation) OldUpdatedAt(ctx context.Context) (v time.Time, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldUpdatedAt is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldUpdatedAt requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldUpdatedAt: %w", err)
	}
	return oldValue.UpdatedAt, nil
}

// ResetUpdatedAt resets all changes to the "updated_at" field.
func (m *FileMutation) ResetUpdatedAt() {
	m.updated_at = nil
}

// SetDeletedAt sets the "deleted_at" field.
func (m *FileMutation) SetDeletedAt(t time.Time) {
	m.deleted_at = &t
}

// DeletedAt returns the value of the "deleted_at" field in the mutation.
func (m *FileMutation) DeletedAt() (r time.Time, exists bool) {
	v := m.deleted_at
	if v == nil {
		return
	}
	return *v, true
}

// OldDeletedAt returns the old "deleted_at" field's value of the File entity.
// If the File object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *FileMutation) OldDeletedAt(ctx context.Context) (v *time.Time, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldDeletedAt is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldDeletedAt requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldDeletedAt: %w", err)
	}
	return oldValue.DeletedAt, nil
}

// ClearDeletedAt clears the value of the "deleted_at" field.
func (m *FileMutation) ClearDeletedAt() {
	m.deleted_at = nil
	m.clearedFields[file.FieldDeletedAt] = struct{}{}
}

// DeletedAtCleared returns if the "deleted_at" field was cleared in this mutation.
func (m *FileMutation) DeletedAtCleared() bool {
	_, ok := m.clearedFields[file.FieldDeletedAt]
	return ok
}

// ResetDeletedAt resets all changes to the "deleted_at" field.
func (m *FileMutation) ResetDeletedAt() {
	m.deleted_at = nil
	delete(m.clearedFields, file.FieldDeletedAt)
}

// Where appends a list predicates to the FileMutation builder.
func (m *FileMutation) Where(ps ...predicate.File) {
	m.predicates = append(m.predicates, ps...)
}

// WhereP appends storage-level predicates to the FileMutation builder. Using this method,
// users can use type-assertion to append predicates that do not depend on any generated package.
func (m *FileMutation) WhereP(ps ...func(*sql.Selector)) {
	p := make([]predicate.File, len(ps))
	for i := range ps {
		p[i] = ps[i]
	}
	m.Where(p...)
}

// Op returns the operation name.
func (m *FileMutation) Op() Op {
	return m.op
}

// SetOp allows setting the mutation operation.
func (m *FileMutation) SetOp(op Op) {
	m.op = op
}

// Type returns the node type of this mutation (File).
func (m *FileMutation) Type() string {
	return m.typ
}

// Fields returns all fields that were changed during this mutation. Note that in
// order to get all numeric fields that were incremented/decremented, call
// AddedFields().
func (m *FileMutation) Fields() []string {
	fields := make([]string, 0, 9)
	if m.uid != nil {
		fields = append(fields, file.FieldUID)
	}
	if m.user_id != nil {
		fields = append(fields, file.FieldUserID)
	}
	if m.filename != nil {
		fields = append(fields, file.FieldFilename)
	}
	if m.object_path != nil {
		fields = append(fields, file.FieldObjectPath)
	}
	if m.size != nil {
		fields = append(fields, file.FieldSize)
	}
	if m.mime_type != nil {
		fields = append(fields, file.FieldMimeType)
	}
	if m.created_at != nil {
		fields = append(fields, file.FieldCreatedAt)
	}
	if m.updated_at != nil {
		fields = append(fields, file.FieldUpdatedAt)
	}
	if m.deleted_at != nil {
		fields = append(fields, file.FieldDeletedAt)
	}
	return fields
}

// Field returns the value of a field with the given name. The second boolean
// return value indicates that this field was not set, or was not defined in the
// schema.
func (m *FileMutation) Field(name string) (ent.Value, bool) {
	switch name {
	case file.FieldUID:
		return m.UID()
	case file.FieldUserID:
		return m.UserID()
	case file.FieldFilename:
		return m.Filename()
	case file.FieldObjectPath:
		return m.ObjectPath()
	case file.FieldSize:
		return m.Size()
	case file.FieldMimeType:
		return m.MimeType()
	case file.FieldCreatedAt:
		return m.CreatedAt()
	case file.FieldUpdatedAt:
		return m.UpdatedAt()
	case file.FieldDeletedAt:
		return m.DeletedAt()
	}
	return nil, false
}

// OldField returns the old value of the field from the database. An error is
// returned if the mutation operation is not UpdateOne, or the query to the
// database failed.
func (m *FileMutation) OldField(ctx context.Context, name string) (ent.Value, error) {
	switch name {
	case file.FieldUID:
		return m.OldUID(ctx)
	case file.FieldUserID:
		return m.OldUserID(ctx)
	case file.FieldFilename:
		return m.OldFilename(ctx)
	case file.FieldObjectPath:
		return m.OldObjectPath(ctx)
	case file.FieldSize:
		return m.OldSize(ctx)
	case file.FieldMimeType:
		return m.OldMimeType(ctx)
	case file.FieldCreatedAt:
		return m.OldCreatedAt(ctx)
	case file.FieldUpdatedAt:
		return m.OldUpdatedAt(ctx)
	case file.FieldDeletedAt:
		return m.OldDeletedAt(ctx)
	}
	return nil, fmt.Errorf("unknown File field %s", name)
}

// SetField sets the value of a field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *FileMutation) SetField(name string, value ent.Value) error {
	switch name {
	case file.FieldUID:
		v, ok := value.(uuid.UUID)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetUID(v)
		return nil
	case file.FieldUserID:
		v, ok := value.(int)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetUserID(v)
		return nil
	case file.FieldFilename:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetFilename(v)
		return nil
	case file.FieldObjectPath:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetObjectPath(v)
		return nil
	case file.FieldSize:
		v, ok := value.(int)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetSize(v)
		return nil
	case file.FieldMimeType:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetMimeType(v)
		return nil
	case file.FieldCreatedAt:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetCreatedAt(v)
		return nil
	case file.FieldUpdatedAt:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetUpdatedAt(v)
		return nil
	case file.FieldDeletedAt:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetDeletedAt(v)
		return nil
	}
	return fmt.Errorf("unknown File field %s", name)
}

// AddedFields returns all numeric fields that were incremented/decremented during
// this mutation.
func (m *FileMutation) AddedFields() []string {
	var fields []string
	if m.adduser_id != nil {
		fields = append(fields, file.FieldUserID)
	}
	if m.addsize != nil {
		fields = append(fields, file.FieldSize)
	}
	return fields
}

// AddedField returns the numeric value that was incremented/decremented on a field
// with the given name. The second boolean return value indicates that this field
// was not set, or was not defined in the schema.
func (m *FileMutation) AddedField(name string) (ent.Value, bool) {
	switch name {
	case file.FieldUserID:
		return m.AddedUserID()
	case file.FieldSize:
		return m.AddedSize()
	}
	return nil, false
}

// AddField adds the value to the field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *FileMutation) AddField(name string, value ent.Value) error {
	switch name {
	case file.FieldUserID:
		v, ok := value.(int)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.AddUserID(v)
		return nil
	case file.FieldSize:
		v, ok := value.(int)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.AddSize(v)
		return nil
	}
	return fmt.Errorf("unknown File numeric field %s", name)
}

// ClearedFields returns all nullable fields that were cleared during this
// mutation.
func (m *FileMutation) ClearedFields() []string {
	var fields []string
	if m.FieldCleared(file.FieldDeletedAt) {
		fields = append(fields, file.FieldDeletedAt)
	}
	return fields
}

// FieldCleared returns a boolean indicating if a field with the given name was
// cleared in this mutation.
func (m *FileMutation) FieldCleared(name string) bool {
	_, ok := m.clearedFields[name]
	return ok
}

// ClearField clears the value of the field with the given name. It returns an
// error if the field is not defined in the schema.
func (m *FileMutation) ClearField(name string) error {
	switch name {
	case file.FieldDeletedAt:
		m.ClearDeletedAt()
		return nil
	}
	return fmt.Errorf("unknown File nullable field %s", name)
}

// ResetField resets all changes in the mutation for the field with the given name.
// It returns an error if the field is not defined in the schema.
func (m *FileMutation) ResetField(name string) error {
	switch name {
	case file.FieldUID:
		m.ResetUID()
		return nil
	case file.FieldUserID:
		m.ResetUserID()
		return nil
	case file.FieldFilename:
		m.ResetFilename()
		return nil
	case file.FieldObjectPath:
		m.ResetObjectPath()
		return nil
	case file.FieldSize:
		m.ResetSize()
		return nil
	case file.FieldMimeType:
		m.ResetMimeType()
		return nil
	case file.FieldCreatedAt:
		m.ResetCreatedAt()
		return nil
	case file.FieldUpdatedAt:
		m.ResetUpdatedAt()
		return nil
	case file.FieldDeletedAt:
		m.ResetDeletedAt()
		return nil
	}
	return fmt.Errorf("unknown File field %s", name)
}

// AddedEdges returns all edge names that were set/added in this mutation.
func (m *FileMutation) AddedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// AddedIDs returns all IDs (to other nodes) that were added for the given edge
// name in this mutation.
func (m *FileMutation) AddedIDs(name string) []ent.Value {
	return nil
}

// RemovedEdges returns all edge names that were removed in this mutation.
func (m *FileMutation) RemovedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// RemovedIDs returns all IDs (to other nodes) that were removed for the edge with
// the given name in this mutation.
func (m *FileMutation) RemovedIDs(name string) []ent.Value {
	return nil
}

// ClearedEdges returns all edge names that were cleared in this mutation.
func (m *FileMutation) ClearedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// EdgeCleared returns a boolean which indicates if the edge with the given name
// was cleared in this mutation.
func (m *FileMutation) EdgeCleared(name string) bool {
	return false
}

// ClearEdge clears the value of the edge with the given name. It returns an error
// if that edge is not defined in the schema.
func (m *FileMutation) ClearEdge(name string) error {
	return fmt.Errorf("unknown File unique edge %s", name)
}

// ResetEdge resets all changes to the edge with the given name in this mutation.
// It returns an error if the edge is not defined in the schema.
func (m *FileMutation) ResetEdge(name string) error {
	return fmt.Errorf("unknown File edge %s", name)
}
