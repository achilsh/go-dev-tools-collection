package example

import (
	mivus_interface "github.com/achilsh/go-dev-tools-collection/vector_db_usage/milvus_interface"
	"github.com/milvus-io/milvus/client/v2/entity"
)

type MilvusQuestionSearch struct {
	searchOpt *mivus_interface.MilvusSearchOption
	//
	collectName string
	targetVect  []float32 //要检索的向量
	retLimit    int       //topK
	outFields   []string
}

func NewMilvusQestionsSearch(collectName string, targetVect []float32, topK int, outFields []string) *MilvusQuestionSearch {
	return &MilvusQuestionSearch{
		collectName: collectName,
		targetVect:  targetVect,
		retLimit:    topK,
		outFields:   outFields,
		searchOpt:   mivus_interface.NewMilvusSearchOptions(),
	}
}

func (mqs *MilvusQuestionSearch) BuildSearchVectOpt() (any, bool) {
	mqs.searchOpt.Register(mivus_interface.WithSearchCollectName(mqs.collectName),
		mivus_interface.WithSearchRetTopK(mqs.retLimit),
		//支持同时检索多个向量
		mivus_interface.WithSearchTargetVector(entity.FloatVector(mqs.targetVect)),
		mivus_interface.WithSearchOutFields(mqs.outFields),
	)

	return mqs.searchOpt.BuildSearchVectOpt()
}
