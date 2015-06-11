package db

import (
	"appengine/datastore"
	"reflect"
)

type Metadata struct {
	Kind          string
	StringID      string
	IntID         int64
	HasParent     bool
	Parent        *datastore.Key
	CacheStringID string
}

func KeyMetadata(key *datastore.Key) *Metadata {
	return &Metadata{
		Kind:      key.Kind(),
		StringID:  key.StringID(),
		IntID:     key.IntID(),
		HasParent: key.Parent() != nil,
		Parent:    key.Parent(),
		CacheStringID: key.Encode(),
	}
}

// IsAutoGenerated tells whether or not a resolved key
// is auto generated by datastore
//
// Keys are auto generated if no struct field is tagged with db:"id"
func (this *Metadata) IsAutoGenerated() bool {
	return this.IntID == 0 && this.StringID == ""
}

type MetadataExtractor interface {
	Accept(reflect.StructField) bool
	Extract(Entity, reflect.StructField, reflect.Value) error
}
