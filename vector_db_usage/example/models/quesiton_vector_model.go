package models

import "fmt"

// 标识一个collection 属性字段
type QuestionVectorCollection struct {
	QuestionVectIndex VectorIndex
	IdIndex           ScalarIndex
	//WithDynamicFieldEnabled
	EnableDynamicField bool // 是否支持动态字段插
	//所有的字段信息
	IdField           FieldProperty
	QuestionVectField FieldProperty //向量
	QuestionStrField  FieldProperty //标量
	//加一个anwer field
	AnswerStrField FieldProperty
	//
	IsDynamicSchema bool
}

// 返回的是collection name
func (QuestionVectorCollection) TableName(lang string) string {
	return fmt.Sprintf("%v_%v", "faq_question_vect_collection", lang)
}

type FieldProperty struct {
	FieldName    string
	IsAutoID     bool //自动产生
	DataType     int32
	IsPrimaryKey bool
	Description  string
	Dim          int // 只有向量才有
	MaxLen       int // 只有标量的字符串才有
}
type VectorIndex struct {
	VectorFieldName string
	VectorIndexName string
	IsAutoIndex     int
	MetricType      string
	//add others...
}
type ScalarIndex struct {
	ScalarFieldName string
	IsSortedIndex   int
	IndexName       string
	//add others..
}
