package db

import (
	"appengine/datastore"
)

// Note: Do not use From to build queries on models implementing SoftDeletableEntity
// where you need to select only results where Deleted=true. That won't work since
// there is already a filter by Deleted=false and apparently datastore instead of
// overriding the previous filter, it only tries to match entities where Delete is
// both, true and false
func From(e Entity) *Query {
	metadata := &Metadata{}
	MetadataExtractorChain{KindExtractor{metadata}}.ExtractFrom(e)
	q := &Query{datastore.NewQuery(metadata.Kind)}
	if _, ok := e.(SoftDeletableEntity); ok {
		q = q.Filter("Deleted=", false)
	}
	return q
}

// Use this functiong to construct queries where you need to fetch
// entities that are soft deleted
func FromSoftDeleted(e SoftDeletableEntity) *Query {
	metadata := &Metadata{}
	MetadataExtractorChain{KindExtractor{metadata}}.ExtractFrom(e)
	return (&Query{datastore.NewQuery(metadata.Kind)}).Filter("Deleted=", true)
}

type Query struct {
	*datastore.Query
}

func (this *Query) Filter(filter string, value interface{}) *Query {
	this.Query = this.Query.Filter(filter, value)
	return this
}

func (this *Query) Limit(limit int) *Query {
	this.Query = this.Query.Limit(limit)
	return this
}

func (this *Query) Project(fieldNames ...string) *Query {
	this.Query = this.Query.Project(fieldNames...)
	return this
}
