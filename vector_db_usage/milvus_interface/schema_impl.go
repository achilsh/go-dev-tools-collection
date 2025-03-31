package mivus_interface

import (
	"github.com/milvus-io/milvus/client/v2/entity"
)

type MilvusSchema struct {
	sch *entity.Schema
}

func NewMilvusSchema() *MilvusSchema {
	return &MilvusSchema{
		sch: entity.NewSchema(),
	}
}

type SchemaOption func(*entity.Schema)

func WithSchemaDynamicFieldEnabled(enable bool) SchemaOption {
	return func(sch *entity.Schema) {
		sch.WithDynamicFieldEnabled(enable)
	}
}
func WithSchemaDescript(des string) SchemaOption {
	return func(sch *entity.Schema) {
		sch.WithDescription(des)
	}
}

func WithSchemaName(name string) SchemaOption {
	return func(sch *entity.Schema) {
		sch.WithName(name)
	}
}

func WithSchemaField(fieldOpts ...FieldOption) SchemaOption {
	field := NewMiluvsField()
	return func(sch *entity.Schema) {
		for _, fopt := range fieldOpts {
			fopt(field)
		}
		sch.WithField(field.field)
	}
}

// Register 中参数可以是： WithSchemaName， WithSchemaField
func (msc *MilvusSchema) Register(opts ...SchemaOption) {
	if msc.sch == nil {
		return
	}
	for _, op := range opts {
		op(msc.sch)
	}
}

// return *entity.Schema
func (msc *MilvusSchema) BuildSchema() (any, error) {
	return msc.sch, nil
}
