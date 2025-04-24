package mivus_interface

import (
	"github.com/milvus-io/milvus/client/v2/entity"
	client "github.com/milvus-io/milvus/client/v2/milvusclient"
)

type MilvusSearchOptItem struct {
	collectName    string
	outFields      []string //
	cLevel         entity.ConsistencyLevel
	topK           int
	vectors        []entity.Vector
	partitionNames []string
	offSet         int //TODO: offset diff with topK
	filterCond     string
}
type MilvusOptions func(*MilvusSearchOptItem)

func WithSearchCollectName(name string) MilvusOptions {
	return func(item *MilvusSearchOptItem) {
		item.collectName = name
	}
}
func WithSearchOutFields(fieldsName []string) MilvusOptions {
	return func(item *MilvusSearchOptItem) {
		item.outFields = fieldsName
	}
}
func WithSearchConsistLevel(level entity.ConsistencyLevel) MilvusOptions {
	return func(item *MilvusSearchOptItem) {
		item.cLevel = level
	}
}
func WithSearchRetTopK(tpk int) MilvusOptions {
	return func(item *MilvusSearchOptItem) {
		item.topK = tpk
	}
}
func WithSearchTargetVector(vect ...entity.Vector) MilvusOptions {
	return func(item *MilvusSearchOptItem) {
		item.vectors = append(item.vectors, vect...)
	}
}
func WithSearchFilterCond(cond string) MilvusOptions {
	return func(item *MilvusSearchOptItem) {
		item.filterCond = cond
	}
}

func WithSearchPartionName(name ...string) MilvusOptions {
	return func(item *MilvusSearchOptItem) {
		item.partitionNames = append(item.partitionNames, name...)
	}
}
func WithSearchTopK(tpk int) MilvusOptions {
	return func(item *MilvusSearchOptItem) {
		item.offSet = tpk
	}
}

type MilvusSearchOption struct {
	option   client.SearchOption
	searchOp MilvusSearchOptItem
}

func NewMilvusSearchOptions() *MilvusSearchOption {
	return &MilvusSearchOption{}
}

func (mso *MilvusSearchOption) Register(opts ...MilvusOptions) {
	item := &mso.searchOp
	for _, opt := range opts {
		opt(item)
	}

	searchOp := client.NewSearchOption(mso.searchOp.collectName, mso.searchOp.topK, mso.searchOp.vectors)
	if len(mso.searchOp.outFields) > 0 {
		searchOp.WithOutputFields(mso.searchOp.outFields...)
	}
	if len(mso.searchOp.partitionNames) > 0 {
		searchOp.WithPartitions(mso.searchOp.partitionNames...)
	}
	if mso.searchOp.offSet > 0 {
		searchOp.WithOffset(mso.searchOp.offSet)
	}
	if mso.searchOp.cLevel >= 0 {
		searchOp.WithConsistencyLevel(mso.searchOp.cLevel)
	}

	if len(mso.searchOp.filterCond) > 0 {
		searchOp.WithFilter(mso.searchOp.filterCond)
	}

	mso.option = searchOp
}
func (mso *MilvusSearchOption) BuildSearchVectOpt() (any, bool) {
	return mso.option, true
}
