package example

import (
	"fmt"

	models "github.com/achilsh/go-dev-tools-collection/vector_db_usage/example/models"

	logger "github.com/achilsh/go-dev-tools-collection/base-lib/log"
	mivus_interface "github.com/achilsh/go-dev-tools-collection/vector_db_usage/milvus_interface"
	"github.com/milvus-io/milvus/client/v2/entity"
)

type MilvusQuestionSchema struct {
	collectInfo *models.QuestionVectorCollection
	lang        string
	schema      *mivus_interface.MilvusSchema
}

func NewMilvusQuestionSchemas(lang string, collInfo *models.QuestionVectorCollection) *MilvusQuestionSchema {
	return &MilvusQuestionSchema{
		lang:        lang,
		collectInfo: collInfo,
		schema:      mivus_interface.NewMilvusSchema(),
	}
}
func (msch *MilvusQuestionSchema) BuildSchema() (any, error) {
	if msch.collectInfo == nil {
		logger.Errorf("not register collection info.")
		return nil, fmt.Errorf("not register collection info.")
	}
	if msch.collectInfo.EnableDynamicField {
		msch.schema.Register(
			mivus_interface.WithSchemaDynamicFieldEnabled(msch.collectInfo.EnableDynamicField),
		)
	}
	//
	msch.schema.Register(
		mivus_interface.WithSchemaField(mivus_interface.WithSchemaFieldName(msch.collectInfo.IdField.FieldName),
			mivus_interface.WithSchemaFieldIsAutoID(msch.collectInfo.IdField.IsAutoID),
			mivus_interface.WithSchemaFieldDataType(entity.FieldType(msch.collectInfo.IdField.DataType)),
			mivus_interface.WithSchemaFieldIsPrimaryKey(msch.collectInfo.IdField.IsPrimaryKey),
		),
		mivus_interface.WithSchemaField(mivus_interface.WithSchemaFieldName(
			msch.collectInfo.QuestionVectField.FieldName),
			mivus_interface.WithSchemaFieldDataType(entity.FieldType(msch.collectInfo.QuestionVectField.DataType)),
			mivus_interface.WithSchemaFieldDimNums(msch.collectInfo.QuestionVectField.Dim),
		),
		mivus_interface.WithSchemaField(mivus_interface.WithSchemaFieldName(
			msch.collectInfo.QuestionStrField.FieldName),
			mivus_interface.WithSchemaFieldDataType(entity.FieldType(msch.collectInfo.QuestionStrField.DataType)),
			mivus_interface.WithSchemaFieldthMaxLength(msch.collectInfo.QuestionStrField.MaxLen),
		))

	ret, err := msch.schema.BuildSchema()
	if err != nil {
		logger.Errorf("create schema fail, err: %v", err)
	}
	return ret, err
}
