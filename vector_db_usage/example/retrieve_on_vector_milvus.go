package example

import (
	"context"
	"fmt"
	"slices"

	models "github.com/achilsh/go-dev-tools-collection/vector_db_usage/example/models"

	logger "github.com/achilsh/go-dev-tools-collection/base-lib/log"
	milvusVect "github.com/achilsh/go-dev-tools-collection/vector_db_usage"
	mivus_interface "github.com/achilsh/go-dev-tools-collection/vector_db_usage/milvus_interface"
	"github.com/milvus-io/milvus/client/v2/column"
	"github.com/milvus-io/milvus/client/v2/entity"
	client "github.com/milvus-io/milvus/client/v2/milvusclient"
	"github.com/samber/lo"
)

const (
	VectDimOnQuestion = 32 //向量（存储）或查询传入的维度数（数组的长度）
	QuestionRetNums   = 5  //检索出来的返回个数
	//
	QuestionVectIndexName = "question_vect_index"
	QuestionVectFieldName = "question_vect"
	QuestioScalarName     = "question"
	IdFieldName           = "id"
	IdeIndexName          = "id_index"
)

var (
	_ SimilarRetriever = (*MilvusRetrieveVect)(nil)
)

func init() {
	registerRetriever(RETRIEVE_VECT_ON_MILVUS,
		&MilvusRetrieveVect{
			QuestionMultiLangCollections: make(map[string]*models.QuestionVectorCollection),
			//
			questionRetNums: QuestionRetNums,
			//
			SearchRetQuestionFieldNames: []string{
				IdFieldName, QuestioScalarName,
			},
		})
}

type MilvusRetrieveVect struct {
	milvusObj milvusVect.VectDBOper
	//
	// key is language.
	QuestionMultiLangCollections map[string]*models.QuestionVectorCollection
	questionRetNums              int      //初始化一次就不需要再变更的 值
	SearchRetQuestionFieldNames  []string //初始化一次就不需要再变更的 值
}

func (mrv *MilvusRetrieveVect) Insert(ctx context.Context) {

	lang_val := "en"
	infoTable := models.QuestionVectorCollection{}.TableName(lang_val)
	if ok := mrv.milvusObj.InsertColumns(ctx, infoTable, NewMilvusInsertOption(lang_val)); !ok {
		logger.Errorf("insert fail")
	} else {
		logger.Infof("insert succ.")
	}
}

// 查询 collecion中的数据
func (mrv *MilvusRetrieveVect) Retrieve(ctx context.Context, targetVect []float32) []string {
	//使用特定的collection
	lang_val := "en"
	infoTable := models.QuestionVectorCollection{}.TableName(lang_val)

	//TODO： 需要在初始化时全部加载
	if loadCollectionRet := mrv.milvusObj.LoadCollection(ctx, infoTable, false); !loadCollectionRet {
		logger.Errorf("load collection: %v fail", infoTable)
		return nil
	}

	defer func() {
		go func() {
			//释放和加载是配套调用
			mrv.milvusObj.ReleaseCollection(ctx, infoTable)
		}()
	}()

	// 返回的是 []ResultSet
	searchRetList, ok := mrv.milvusObj.SearchVector(ctx, infoTable, NewMilvusQestionsSearch(
		infoTable, targetVect, mrv.questionRetNums, mrv.SearchRetQuestionFieldNames))
	if !ok {
		logger.Errorf("search vector fail.")
		return nil
	}

	retrieveRets := mrv.parseResponse(searchRetList.([]client.ResultSet))
	var retQuestions []string
	for _, result := range retrieveRets {
		_, q, _ := lo.Unpack3(result)
		retQuestions = append(retQuestions, q)
	}
	return retQuestions
}

func (mrv *MilvusRetrieveVect) parseResponse(data []client.ResultSet) []lo.Tuple3[int64, string, float32] {
	//返回 question and id

	var matchQuestion []lo.Tuple3[int64, string, float32] //id, quesiton, score

	for sRIndex := 0; sRIndex < len(data); sRIndex++ {
		result := data[sRIndex]
		var idColumn *column.ColumnInt64
		var questionScalarColumn *column.ColumnVarChar

		for _, field := range result.Fields {
			if field.Name() == IdFieldName {
				c, ok := field.(*column.ColumnInt64)
				if ok {
					idColumn = c
				}
				continue
			}
			if field.Name() == QuestioScalarName {
				c, ok := field.(*column.ColumnVarChar)
				if ok {
					questionScalarColumn = c
				}
				continue
			}
		}

		for rIndex := 0; rIndex < result.ResultCount; rIndex++ {
			id, err := idColumn.GetAsInt64(rIndex)
			if err != nil {
				logger.Errorf("get id fail,err: %v, rIndex: %v", err, rIndex)
				continue
			}

			question, err := questionScalarColumn.GetAsString(rIndex)
			if err != nil {
				logger.Errorf("get question  fail, err: %v, i: %v", err, rIndex)
				continue
			}

			matchQuestion = append(matchQuestion, lo.T3(id, question, result.Scores[rIndex]))
		}
	}
	logger.Debugf("retrieve similar result: %+v", matchQuestion)
	return matchQuestion

}

func (mrv *MilvusRetrieveVect) InitAllLanguageCollection() {
	registeredLanguages := []string{
		// "de",
		"en",
		// "es",
		// "fr",
		// "it",
		// "ja",
		//TODO: add others.
	}
	//填充每种语言 milvus collection创建的 index, schema 信息
	for _, lang := range registeredLanguages {
		mrv.QuestionMultiLangCollections[lang] = &models.QuestionVectorCollection{
			QuestionVectIndex: models.VectorIndex{
				VectorFieldName: QuestionVectFieldName,
				VectorIndexName: QuestionVectIndexName,
				IsAutoIndex:     mivus_interface.IndexTypeAutoIndex, //设置索引类型
				MetricType:      string(entity.COSINE),
			},
			IdIndex: models.ScalarIndex{
				ScalarFieldName: IdFieldName,
				IsSortedIndex:   mivus_interface.IndexTypeSortedIndex,
				IndexName:       IdeIndexName,
			},
			//
			EnableDynamicField: true,

			//主键
			IdField: models.FieldProperty{
				FieldName:    IdFieldName,
				IsAutoID:     true, //主键自动递增模式
				DataType:     int32(entity.FieldTypeInt64),
				IsPrimaryKey: true,
				Description:  "this is id primary key",
			},
			//向量
			QuestionVectField: models.FieldProperty{
				FieldName:   QuestionVectFieldName,
				DataType:    int32(entity.FieldTypeFloatVector),
				Dim:         VectDimOnQuestion,
				Description: "this is question vector",
			},
			// 标量
			QuestionStrField: models.FieldProperty{
				FieldName:   QuestioScalarName,
				DataType:    int32(entity.FieldTypeVarChar),
				MaxLen:      512,
				Description: "this is question detail",
			},
			IsDynamicSchema: true,
		}
	}
}
func (mrv *MilvusRetrieveVect) getNoCreatedQuestionCollection() []string {
	noCreatedCollection := []string{}

	collectionList := mrv.milvusObj.ListCollection(context.Background())
	logger.Debugf("list collections: %v", collectionList)

	// 检查代码注册中所有的集合 是否都被创建
	for lang, _ := range mrv.QuestionMultiLangCollections {
		collectName := models.QuestionVectorCollection{}.TableName(lang)

		if !slices.Contains(collectionList, collectName) {
			logger.Errorf("has not create collection: %v", collectName)
			noCreatedCollection = append(noCreatedCollection, collectName)
		}
	}
	logger.Infof("need to create collection list: %+v", noCreatedCollection)
	return noCreatedCollection
}

// 创建项目本应该有的，但是还未创建的 collection.
func (mrv *MilvusRetrieveVect) CreateCollection(toCreateTabNames []string) error {
	for lang, collectTableInfo := range mrv.QuestionMultiLangCollections {

		infoTable := models.QuestionVectorCollection{}.TableName(lang)
		if slices.Contains(toCreateTabNames, infoTable) {

			createRet := mrv.milvusObj.CreateCollection(context.Background(), infoTable,
				NewMilvusQuestionIndexer(lang, collectTableInfo),
				NewMilvusQuestionSchemas(lang, collectTableInfo))

			if !createRet {
				logger.Errorf("create collection fail, collectName: %v", infoTable)
				//
				panic("create collection fail.")
				return fmt.Errorf("create collection: %v fail", infoTable)
			}
		}
	}
	return nil
}

func (mrv *MilvusRetrieveVect) Init() bool {
	addrHost := "localhost:19530"
	dbName := "test_demo_milvus"

	mrv.milvusObj = mivus_interface.NewVectMilvusOpInst(
		mivus_interface.WithAddress(addrHost),
	)

	if !mrv.milvusObj.Connect(context.Background()) {
		logger.Errorf("connect milvus database fail, remote host: %v", addrHost)
		return false
	}

	dbCreatedList := mrv.milvusObj.ListDB(context.Background())
	if !slices.Contains(dbCreatedList, dbName) {
		logger.Infof("dbName: %v not create, now  create db: %v", dbName)

		if !mrv.milvusObj.CreateDB(context.Background(), dbName) {
			logger.Errorf("create db fail, dbName: %v", dbName)
			return false
		}
	}

	if !mrv.milvusObj.UsingDB(context.Background(), dbName) {
		logger.Errorf("use dbname: %v fail", dbName)
		return false
	}
	logger.Infof("connect milvus: %v and use db: %v succ", addrHost, dbName)

	//初始化
	mrv.InitAllLanguageCollection()
	toCreateCollections := mrv.getNoCreatedQuestionCollection()
	if err := mrv.CreateCollection(toCreateCollections); err != nil {
		logger.Errorf("create collection fail, err: %v", err)
		return false
	}
	return true
}
