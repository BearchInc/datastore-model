package db

import (
	"appengine"
	"appengine/datastore"
	"reflect"
)

type PageIterator struct {
	query              *datastore.Query
	context            appengine.Context
	nextCursor         datastore.Cursor
	prevCursor         datastore.Cursor
	doneProcessingPage bool
}

func NewPagesIterator(q *datastore.Query, c appengine.Context) *PageIterator {
	return &PageIterator{
		query:              q,
		context:            c,
		nextCursor:         datastore.Cursor{},
		prevCursor:         datastore.Cursor{},
		doneProcessingPage: false,
	}
}

// TODO refactor this mess :~
// Perhaps it would be better to have a PerPageIterator and a PerItemIterator
// to avoid messing up with the iterator internal state when using
// LoadNext and LoadNextPage intermittently
func (this *PageIterator) LoadNext(slice interface{}) error {
	sv := reflect.ValueOf(slice)
	if sv.Kind() != reflect.Ptr || sv.IsNil() || sv.Elem().Kind() != reflect.Slice {
		return datastore.ErrInvalidEntityType
	}
	sv = sv.Elem()

	elemType := sv.Type().Elem()
	if elemType.Kind() != reflect.Ptr || elemType.Elem().Kind() != reflect.Struct {
		return datastore.ErrInvalidEntityType
	}

	iter := this.query.Run(this.context)
	for {
		dstValue := reflect.New(elemType.Elem())
		dst := dstValue.Interface()
		entity, ok := dst.(Entity)
		if !ok {
			return datastore.ErrInvalidEntityType
		}

		key, err := iter.Next(entity)
		this.prevCursor = this.nextCursor

		cursor, cursorErr:= iter.Cursor()
		if cursorErr != nil {
			return cursorErr
		}
		this.nextCursor = cursor
		entity.SetKey(key)
		sv.Set(reflect.Append(sv, dstValue))

		if err == datastore.Done {
			this.doneProcessingPage = true
			this.query = this.query.Start(cursor)
			return err
		}
		if err != nil {
			return nil
		}
	}

	return nil
}

func (this *PageIterator) HasNext() bool {
	return !this.doneProcessingPage || this.prevCursor != this.nextCursor
}

func (this *PageIterator) Cursor() string {
	return this.nextCursor.String()
}