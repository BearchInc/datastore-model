package db

import (
	"appengine"
	"appengine/datastore"
	"reflect"
	"errors"
)

type PageIterator struct {
	query              *datastore.Query
	context            appengine.Context
	nextCursor         datastore.Cursor
	prevCursor         datastore.Cursor
	started    bool
}

func NewPagesIterator(q *datastore.Query, c appengine.Context) *PageIterator {
	return &PageIterator{
		query:              q,
		context:            c,
		nextCursor:         datastore.Cursor{},
		prevCursor:         datastore.Cursor{},
	}
}

// TODO refactor this mess :~
// Perhaps it would be better to have a PerPageIterator and a PerItemIterator
// to avoid messing up with the iterator internal state when using
// LoadNext and LoadNextPage intermittently
func (this *PageIterator) LoadNext(slice interface{}) error {
	ErrInvalidEntityType := errors.New("Invalid entity type. Make sure your model implements appx.Entity (watch out for pointer receivers)")
	ErrInvalidSliceType  := errors.New("Invalid slice type. Make sure you pass a pointer to a slice of appx.Entity")

	sv := reflect.ValueOf(slice)
	if sv.Kind() != reflect.Ptr || sv.IsNil() || sv.Elem().Kind() != reflect.Slice {
		return ErrInvalidSliceType
	}
	sv = sv.Elem()

	elemType := sv.Type().Elem()
	if elemType.Kind() != reflect.Ptr || elemType.Elem().Kind() != reflect.Struct {
		return ErrInvalidEntityType
	}

	this.started = true
	iter := this.query.Run(this.context)
	for {
		dstValue := reflect.New(elemType.Elem())
		dst := dstValue.Interface()

		entity, ok := dst.(Entity)
		if !ok {
			return ErrInvalidEntityType
		}

		key, err := iter.Next(entity)
		if err == datastore.Done {
			cursor, cursorErr := iter.Cursor()
			if cursorErr != nil {
				return cursorErr
			}
			this.prevCursor = this.nextCursor
			this.nextCursor = cursor
			this.query = this.query.Start(this.nextCursor)
			if !this.HasNext() {
				this.prevCursor = datastore.Cursor{}
				this.nextCursor = datastore.Cursor{}
				return datastore.Done
			}
			break
		}

		if err != nil {
			return err
		}

		entity.SetKey(key)
		sv.Set(reflect.Append(sv, dstValue))
	}

	return nil
}

func (this *PageIterator) HasNext() bool {
	return !this.started || this.prevCursor.String() != this.nextCursor.String()
}

func (this *PageIterator) Cursor() string {
	return this.nextCursor.String()
}