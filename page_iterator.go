package db

import (
	"appengine"
	"appengine/datastore"
	"reflect"
)

type PageIterator struct {
	query   *datastore.Query
	context appengine.Context
	nextCursor  datastore.Cursor
	prevCursor datastore.Cursor
}

func NewPagesIterator(q *datastore.Query, c appengine.Context) *PageIterator {
	return &PageIterator{
		query:   q,
		context: c,
	}
}

// TODO refactor this mess :~
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
		this.context.Infof("Epa epa epa 0")
		dstValue := reflect.New(elemType.Elem())
		dst := dstValue.Interface()

		entity, ok := dst.(Entity)
		if !ok {
			this.context.Infof("Epa epa epa 1")
			return datastore.ErrInvalidEntityType
		}

		key, err := iter.Next(entity)
		if err == datastore.Done {
			cursor, cursorErr := iter.Cursor()
			if cursorErr != nil {
				this.context.Infof("Epa epa epa 2")
				return cursorErr
			}
			this.prevCursor = this.nextCursor
			this.nextCursor = cursor
			this.query = this.query.Start(this.nextCursor)
			this.context.Infof("Epa epa epa 3")
			return err
		}

		if err != nil {
			this.context.Infof("Epa epa epa 4 - %+v", err)
			return err
		}

		entity.SetKey(key)
		sv.Set(reflect.Append(sv, dstValue))
	}

	return nil
}

func (this *PageIterator) HasNext() bool {
	return this.prevCursor.String() != this.nextCursor.String()
}

func (this *PageIterator) Cursor() string {
	return this.nextCursor.String()
}