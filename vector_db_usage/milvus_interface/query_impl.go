package mivus_interface

import (
	"github.com/milvus-io/milvus/client/v2/column"
	"github.com/milvus-io/milvus/client/v2/entity"
	client "github.com/milvus-io/milvus/client/v2/milvusclient"
)

type queryTemplateParam struct {
	key  string
	data any
}
type MilvusQueryOptItem struct {
	collectName    string
	filterCond     string
	tempParam      *queryTemplateParam
	offset         int
	limit          int //limit and offset diff.
	outFieldName   []string
	cLevel         entity.ConsistencyLevel
	partitionNames []string
	ids            column.Column //主键查询
}

type MilvusQueryOption func(*MilvusQueryOptItem)

func WithQueryCollectName(collName string) MilvusQueryOption {
	return func(qci *MilvusQueryOptItem) {
		qci.collectName = collName
	}
}

func WithQueryFilterCond(cond string) MilvusQueryOption {
	return func(qci *MilvusQueryOptItem) {
		qci.filterCond = cond
	}
}

func WithQueryTemplateParm(key string, data any) MilvusQueryOption {
	return func(qci *MilvusQueryOptItem) {
		qci.tempParam = &queryTemplateParam{
			key:  key,
			data: data,
		}
	}
}

func WithQueryOffset(offset int) MilvusQueryOption {
	return func(qci *MilvusQueryOptItem) {
		qci.offset = offset
	}
}

func WithQuerylimit(limit int) MilvusQueryOption {
	return func(qci *MilvusQueryOptItem) {
		qci.limit = limit
	}
}

func WithQueryOutField(fieldName ...string) MilvusQueryOption {
	return func(qci *MilvusQueryOptItem) {
		qci.outFieldName = append(qci.outFieldName, fieldName...)
	}
}
func WithQueryConsistenLevel(cLevel int) MilvusQueryOption {
	return func(qci *MilvusQueryOptItem) {
		qci.cLevel = entity.ConsistencyLevel(cLevel)
	}
}

func WithQueryPartitionNames(partitionName ...string) MilvusQueryOption {
	return func(qci *MilvusQueryOptItem) {
		qci.partitionNames = append(qci.partitionNames, partitionName...)
	}
}

// eg: column.NewColumnInt64("id", []int64{1, 2, 3})
func WithQueryIds(ids column.Column) MilvusQueryOption {
	return func(qci *MilvusQueryOptItem) {
		qci.ids = ids
	}
}

type MilvusQueryOptionImpl struct {
	option  client.QueryOption
	optItem *MilvusQueryOptItem
}

func NewMilvusQueryOption() *MilvusQueryOptionImpl {
	return &MilvusQueryOptionImpl{
		optItem: new(MilvusQueryOptItem),
	}
}
func (moi *MilvusQueryOptionImpl) Register(opts ...MilvusQueryOption) {
	for _, opt := range opts {
		opt(moi.optItem)
	}

	queryOpt := client.NewQueryOption(moi.optItem.collectName)
	if moi.optItem.filterCond != "" {
		queryOpt.WithFilter(moi.optItem.filterCond)
	}

	if moi.optItem.tempParam != nil {
		queryOpt.WithTemplateParam(moi.optItem.tempParam.key, moi.optItem.tempParam.data)
	}

	if moi.optItem.offset > 0 {
		queryOpt.WithOffset(moi.optItem.offset)
	}

	if moi.optItem.limit > 0 {
		queryOpt.WithLimit(moi.optItem.limit)
	}

	if moi.optItem.offset > 0 {
		queryOpt.WithOffset(moi.optItem.offset)
	}

	if len(moi.optItem.outFieldName) > 0 {
		queryOpt.WithOutputFields(moi.optItem.outFieldName...)
	}

	if moi.optItem.cLevel >= 0 {
		queryOpt.WithConsistencyLevel(moi.optItem.cLevel)
	}

	if len(moi.optItem.partitionNames) > 0 {
		queryOpt.WithPartitions(moi.optItem.partitionNames...)
	}

	if moi.optItem.ids != nil {
		queryOpt.WithIDs(moi.optItem.ids)
	}

	moi.option = queryOpt
}

func (mso *MilvusQueryOptionImpl) BuildQueryScalarOpt() (any, bool) {
	return mso.option, true
}
