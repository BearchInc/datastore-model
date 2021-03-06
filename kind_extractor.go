package db

import (
	"reflect"
	"strings"
)

type KindExtractor struct {
	Metadata *Metadata
}

func (this KindExtractor) Accept(f reflect.StructField) bool {
	return f.Type.Name() == reflect.TypeOf(Model{}).Name() || f.Type.Name() == reflect.TypeOf(SoftDeletableModel{}).Name()
}

func (this KindExtractor) Extract(e Entity, f reflect.StructField, v reflect.Value) error {
	elem := reflect.TypeOf(e).Elem()
	this.Metadata.Kind = elem.Name()

	kindMetadata := f.Tag.Get("db")
	values := strings.Split(kindMetadata, ",")
	if strings.TrimSpace(values[0]) != "" {
		this.Metadata.Kind = values[0]
	}

	return nil
}
