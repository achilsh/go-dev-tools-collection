package mivus_interface

import "github.com/milvus-io/milvus/client/v2/entity"

type MilvusField struct {
	field *entity.Field
}

func NewMiluvsField() *MilvusField {
	return &MilvusField{
		field: entity.NewField(),
	}
}

type FieldOption func(*MilvusField)

func WithSchemaFieldName(name string) FieldOption {
	return func(item *MilvusField) {
		item.field.WithName(name)
	}
}
func WithSchemaFieldIsAutoID(enable bool) FieldOption {
	return func(item *MilvusField) {
		item.field.WithIsAutoID(enable)
	}
}

// FieldTypeNone FieldType = 0 // zero value place holder
// // FieldTypeBool field type boolean
// FieldTypeBool FieldType = 1
// // FieldTypeInt8 field type int8
// FieldTypeInt8 FieldType = 2
// // FieldTypeInt16 field type int16
// FieldTypeInt16 FieldType = 3
// // FieldTypeInt32 field type int32
// FieldTypeInt32 FieldType = 4
// // FieldTypeInt64 field type int64
// FieldTypeInt64 FieldType = 5
// // FieldTypeFloat field type float
// FieldTypeFloat FieldType = 10
// // FieldTypeDouble field type double
// FieldTypeDouble FieldType = 11
// // FieldTypeString field type string
// FieldTypeString FieldType = 20
// // FieldTypeVarChar field type varchar
// FieldTypeVarChar FieldType = 21 // variable-length strings with a specified maximum length
// // FieldTypeArray field type Array
// FieldTypeArray FieldType = 22
// // FieldTypeJSON field type JSON
// FieldTypeJSON FieldType = 23
// // FieldTypeBinaryVector field type binary vector
// FieldTypeBinaryVector FieldType = 100
// // FieldTypeFloatVector field type float vector
// FieldTypeFloatVector FieldType = 101
// // FieldTypeBinaryVector field type float16 vector
// FieldTypeFloat16Vector FieldType = 102
// // FieldTypeBinaryVector field type bf16 vector
// FieldTypeBFloat16Vector FieldType = 103
// // FieldTypeBinaryVector field type sparse vector
// FieldTypeSparseVector FieldType = 104
// // FieldTypeInt8Vector field type int8 vector
// FieldTypeInt8Vector FieldType = 105
// eg: entity.FieldTypeInt64； entity.FieldTypeFloatVector, entity.FieldTypeVarChar
func WithSchemaFieldDataType(fieldType entity.FieldType) FieldOption {
	return func(item *MilvusField) {
		item.field.WithDataType(fieldType)
	}
}

// 设置 是否是主键
func WithSchemaFieldIsPrimaryKey(enable bool) FieldOption {
	return func(item *MilvusField) {
		item.field.WithIsPrimaryKey(enable)
	}
}

func WithSchemaFieldDimNums(n int) FieldOption {
	return func(item *MilvusField) {
		item.field.WithDim(int64(n))
	}
}

func WithSchemaFieldthMaxLength(maxLen int) FieldOption {
	return func(item *MilvusField) {
		item.field.WithMaxLength(int64(maxLen))
	}
}
