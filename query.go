package db

import (
	"appengine/datastore"
)

func From(e Entity) *Query {
	metadata := &Metadata{}
	MetadataExtractorChain{KindExtractor{metadata}}.ExtractFrom(e)
	q := &Query{datastore.NewQuery(metadata.Kind)}

	// Filter out soft deleted items in case we are
	// dealing with a model supporting soft deletes
	if _, ok := e.(SoftDeletableEntity); ok {
		q = q.Filter("Deleted=", false)
	}
	return q
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
