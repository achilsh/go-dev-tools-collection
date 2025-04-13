package example

import (
	"fmt"

	models "github.com/achilsh/go-dev-tools-collection/vector_db_usage/example/models"

	logger "github.com/achilsh/go-dev-tools-collection/base-lib/log"
	mivus_interface "github.com/achilsh/go-dev-tools-collection/vector_db_usage/milvus_interface"
	"github.com/milvus-io/milvus/client/v2/index"
)

type MilvusQuestionIndexer struct {
	collectInfo *models.QuestionVectorCollection
	lang        string
	indexer     *mivus_interface.MilvusCreateIndexOptionList
}

func NewMilvusQuestionIndexer(lang string, collInfo *models.QuestionVectorCollection) *MilvusQuestionIndexer {
	return &MilvusQuestionIndexer{
		lang:        lang,
		collectInfo: collInfo,
		indexer:     mivus_interface.NewMilvusCreateIndexOptionList(),
	}
}

func (it *MilvusQuestionIndexer) BuildCreateIndexOptions() (any, error) {
	if it.collectInfo == nil {
		logger.Errorf("not register collection info.")
		return nil, fmt.Errorf("not register collection info.")
	}

	it.indexer.RegisterCreateIndexOpt(
		it.collectInfo.TableName(it.lang),
		//
		it.collectInfo.QuestionVectIndex.VectorFieldName,
		//
		it.collectInfo.QuestionVectIndex.VectorIndexName,
		//
		mivus_interface.WithIndexMetricType(
			index.MetricType(it.collectInfo.QuestionVectIndex.MetricType),
			it.collectInfo.QuestionVectIndex.IsAutoIndex),
	)

	//两个索引
	it.indexer.RegisterCreateIndexOpt(
		it.collectInfo.TableName(it.lang),
		it.collectInfo.IdIndex.ScalarFieldName,
		it.collectInfo.IdIndex.IndexName,
		mivus_interface.WithScalarIndex(it.collectInfo.IdIndex.IsSortedIndex),
	)

	createIndex, err := it.indexer.BuildCreateIndexOptions()
	if err != nil {
		logger.Errorf("create index fail, err: %v", err)
	}
	return createIndex, err
}
