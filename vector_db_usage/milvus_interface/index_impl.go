package mivus_interface

import (
	"fmt"

	"github.com/milvus-io/milvus/client/v2/index"
	client "github.com/milvus-io/milvus/client/v2/milvusclient"
)

const (
	IndexTypeAutoIndex          = 0  // autoIndex
	IndexTypediskANNIndex       = 1  // diskANNIndex
	IndexTypeflatIndex          = 2  // flatIndex
	IndexTypebinFlatIndex       = 3  //binFlatIndex
	IndexTypegpuBruteForceIndex = 4  //gpuBruteForceIndex
	IndexTypegpuIVFFlatIndex    = 5  //gpuIVFFlatIndex
	IndexTypegpuIVFPQIndex      = 6  //gpuIVFPQIndex
	IndexTypegpuCagra           = 7  //gpuCagra
	IndexTypehnswIndex          = 8  //hnswIndex
	IndexTypeGenericIndex       = 9  //GenericIndex
	IndexTypeivfFlatIndex       = 10 // ivfFlatIndex
	IndexTypeivfPQIndex         = 11 // ivfPQIndex

	IndexTypeivfSQ8Index = 12 // ivfSQ8Index
	IndexTypebinIvfFlat  = 13 // binIvfFlat

	IndexTypeJSONPathIndex = 14 // JSONPathIndex

	// 标量索引类型
	IndexTypeTrieIndex     = 15 // TrieIndex
	IndexTypeInvertedIndex = 16 // InvertedIndex
	IndexTypeSortedIndex   = 17 //SortedIndex
	IndexTypeBitmapIndex   = 18 // BitmapIndex
	//
	IndexTypeSCANNIndex          = 30 // SCANNIndex
	IndexTypesparseInvertedIndex = 31 // sparseInvertedIndex
	IndexTypesparseWANDIndex     = 32 // sparseWANDIndex

	IndexTypesparseAnnParam = 33 // sparseAnnParam

)

type MilvusIndexOptItem struct {
	collectName string
	fieldName   string
	indexName   string
	//
	indexItem index.Index
}

type MilvusIndexCreate func(*MilvusIndexOptItem)

func WithMilvusBasic(collName, fieldName, indexName string) MilvusIndexCreate {
	return func(opt *MilvusIndexOptItem) {
		opt.collectName = collName
		opt.fieldName = fieldName
		opt.indexName = indexName
	}
}

func WithIndexMetricType(mectricType index.MetricType, indexType int) MilvusIndexCreate {
	var indexItem index.Index = nil

	switch indexType {
	case IndexTypeAutoIndex:
		indexItem = index.NewAutoIndex(mectricType)
	case IndexTypediskANNIndex:
		indexItem = index.NewDiskANNIndex(mectricType)
	case IndexTypeflatIndex:
		indexItem = index.NewFlatIndex(mectricType)
	case IndexTypebinFlatIndex:
		indexItem = index.NewBinFlatIndex(mectricType)
	case IndexTypegpuBruteForceIndex:
		indexItem = index.NewGPUBruteForceIndex(mectricType)
	case IndexTypegpuIVFFlatIndex:
		indexItem = index.NewGPUIVPFlatIndex(mectricType)
	case IndexTypegpuIVFPQIndex:
		indexItem = index.NewGPUIVPPQIndex(mectricType)
	}

	return func(opt *MilvusIndexOptItem) {
		opt.indexItem = indexItem
	}
}

func WithIndexGPUCagraIndex(metricType index.MetricType, intermediateGD, gd int) MilvusIndexCreate {
	var indexItem index.Index = index.NewGPUCagraIndex(metricType, intermediateGD, gd)
	return func(opt *MilvusIndexOptItem) {
		opt.indexItem = indexItem
	}
}
func WithHNSWIndex(metricType index.MetricType, m, efConstruction int) MilvusIndexCreate {
	var indexItem index.Index = index.NewHNSWIndex(metricType, m, efConstruction)
	return func(opt *MilvusIndexOptItem) {
		opt.indexItem = indexItem
	}
}

func WithIvfSQ8Index(metricType index.MetricType, nlist int) MilvusIndexCreate {
	var indexItem index.Index = index.NewIvfSQ8Index(metricType, nlist)
	return func(opt *MilvusIndexOptItem) {
		opt.indexItem = indexItem
	}
}

func WithBinIvfFlatIndex(metricType index.MetricType, nlist int) MilvusIndexCreate {
	var indexItem index.Index = index.NewBinIvfFlatIndex(metricType, nlist)
	return func(opt *MilvusIndexOptItem) {
		opt.indexItem = indexItem
	}
}

func WithGenericIndex(metricType index.MetricType, name string, params map[string]string) MilvusIndexCreate {
	item := index.NewGenericIndex(name, params)
	item.WithMetricType(metricType)
	return func(opt *MilvusIndexOptItem) {
		opt.indexItem = item
	}
}

func WithIvfFlatIndex(metricType index.MetricType, nlist int) MilvusIndexCreate {
	var indexItem index.Index = index.NewIvfFlatIndex(metricType, nlist)
	return func(opt *MilvusIndexOptItem) {
		opt.indexItem = indexItem
	}
}

func WithIvfPQIndex(metricType index.MetricType, nlist, m, nbits int) MilvusIndexCreate {
	var indexItem index.Index = index.NewIvfPQIndex(metricType, nlist, m, nbits)
	return func(opt *MilvusIndexOptItem) {
		opt.indexItem = indexItem
	}
}
func WithScalarIndex(indexType int) MilvusIndexCreate {
	var indexItem index.Index = nil
	switch indexType {
	case IndexTypeTrieIndex:
		indexItem = index.NewTrieIndex()
	case IndexTypeInvertedIndex:
		indexItem = index.NewInvertedIndex()
	case IndexTypeSortedIndex:
		indexItem = index.NewSortedIndex()
	case IndexTypeBitmapIndex:
		indexItem = index.NewBitmapIndex()
	default:
		//

	}

	return func(opt *MilvusIndexOptItem) {
		opt.indexItem = indexItem
	}
}

func WithSparseIndex(metricType index.MetricType, dRation float64, indexType int) MilvusIndexCreate {
	var indexItem index.Index = nil
	switch indexType {
	case IndexTypesparseWANDIndex:
		indexItem = index.NewSparseInvertedIndex(metricType, dRation)
	case IndexTypesparseAnnParam:
		indexItem = index.NewSparseWANDIndex(metricType, dRation)
	default:
		//
	}
	return func(opt *MilvusIndexOptItem) {
		opt.indexItem = indexItem

	}
}
func WithNewSCANNIndex(metricType index.MetricType, nlist int, withRawData bool) MilvusIndexCreate {
	var indexItem index.Index = index.NewSCANNIndex(metricType, nlist, withRawData)
	return func(opt *MilvusIndexOptItem) {
		opt.indexItem = indexItem
	}
}

// 一列建立的索引需要的数据
type milvusCreateIndexOption struct {
	// 一列字段
	optionDetails *MilvusIndexOptItem
	// 一列索引创建多次过程
	createOps []MilvusIndexCreate
}

// 一个 collection的多列索引
type MilvusCreateIndexOptionList struct {
	//多个索引和创建过程
	optionList []*milvusCreateIndexOption
}

func NewMilvusCreateIndexOptionList() *MilvusCreateIndexOptionList {
	ret := &MilvusCreateIndexOptionList{}
	return ret
}

// RegisterCreateIndexOpt 注册 集合名，索引列名，索引名和创建索引动作； 用于该列索引的创建。
// 其中 createIndexOpts 有： WithIndexMetricType 等
func (mcList *MilvusCreateIndexOptionList) RegisterCreateIndexOpt(collectName string, fieldName string, indexName string, createIndexopts ...MilvusIndexCreate) {
	item := &milvusCreateIndexOption{
		optionDetails: &MilvusIndexOptItem{
			collectName: collectName,
			fieldName:   fieldName,
			indexName:   indexName,
		},
	}

	item.createOps = append(item.createOps, WithMilvusBasic(collectName, fieldName, indexName))
	item.createOps = append(item.createOps, createIndexopts...)
	//
	mcList.optionList = append(mcList.optionList, item)
}

// BuildCreateIndexOptions 返回多个索引的创建过程： []client.CreateIndexOption
func (mcList *MilvusCreateIndexOptionList) BuildCreateIndexOptions() (any, error) {
	if mcList == nil {
		return nil, fmt.Errorf("obj is nil")
	}
	var retCreateIndexOpt []client.CreateIndexOption
	for _, createIndexOp := range mcList.optionList {
		if createIndexOp == nil {
			continue
		}

		for _, opt := range createIndexOp.createOps {
			if opt == nil {
				continue
			}
			opt(createIndexOp.optionDetails)
		}

		createIndexOptRet := client.NewCreateIndexOption(
			createIndexOp.optionDetails.collectName,
			createIndexOp.optionDetails.fieldName,
			createIndexOp.optionDetails.indexItem,
		).WithIndexName(createIndexOp.optionDetails.indexName)

		retCreateIndexOpt = append(retCreateIndexOpt, createIndexOptRet)
	}
	return retCreateIndexOpt, nil
}
