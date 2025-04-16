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
	idScalarList := []int64{
		1000,
		2000,
		3000,
		4000,
		5000,
		6000,
	}

	questionSliceStr := []string{
		"深圳市是中国广东省一个市，也是中国的经济特区",
		"香港是中国的特别行政区",
		"我是谁，我是一个it 码农。",
		"春天是一个播种的季节",
		"一年小学生学古诗100首。",
		"中国有八大菜系，湘菜，川菜，鲁菜等.",
	}

	vectDataList := [][]float32{
		{0.14181516, 0.053447783, 0.27755734, 0.425333, 0.12168437, 0.10610832, 0.02128352, 0.21277897, -0.118872814, -0.28497985, -0.15902191, 0.24786726, -0.056231227, 0.18140194, 0.012834778, 0.040092863, -0.35830536, -0.37157595, 0.074899994, 0.15587296, 0.14147776, 0.07073888, 0.07698055, 0.0009163933, -0.23774563, -0.05769324, 0.01973716, -0.13124369, -0.021030478, 0.13664187, 0.020552514, -0.05544399},
		{0.06335796, 0.057832558, 0.20333485, 0.5510669, 0.06773224, -0.027235635, -0.032899175, 0.19523092, -0.047127087, -0.15848699, -0.06856105, 0.17552365, -0.14706782, 0.09577366, 0.11769109, -0.037895057, -0.46376553, 0.052491333, 0.19559929, 0.044709723, -0.031909205, 0.007867944, 0.14927799, 0.060917575, -0.20020379, -0.03754972, -0.1639203, -0.21549073, -0.038240395, 0.2149382, 0.16705136, -0.13647747},
		{-0.17596263, 0.00033492805, -0.19089417, 0.032562274, 0.31517065, -0.37880203, -0.09877797, 0.3172381, -0.3328588, -0.22385845, -0.05447146, 0.004956845, -0.1614905, 0.04045877, -0.10009884, -0.17803007, -0.17320603, -0.08528215, -0.2761189, 0.027465444, 0.026331222, -0.018104509, 0.11416896, -0.0938965, 0.13461371, -0.20915662, 0.105497174, -0.037472457, -0.06742171, -0.15770018, -0.24143174, -0.2138658},
		{0.34291446, 0.015021714, 0.09629263, 0.055509117, 0.22280477, -0.28606084, 0.058550276, -0.11946945, -0.076060936, -0.21678647, -0.014277432, 0.2484145, -0.15609138, 0.007362801, 0.08137496, -0.21076818, -0.23663403, -0.029051052, -0.24406084, 0.035309434, 0.10371946, 0.14174993, 0.08905788, -0.2087194, 0.11639628, -0.05451674, 0.30552423, 0.30526814, 0.23548159, -0.05896643, 0.048370402, -0.25724986},
		{0.018243162, -0.16492714, 0.08501537, -0.19250453, 0.14263242, -0.18193917, -0.15660019, 0.08707473, 0.0025154299, -0.011242488, -0.084433384, -0.27631116, -0.30012798, 0.30173966, 0.28651837, -0.07637505, -0.089089304, 0.03988873, -0.0081982305, -0.023011006, 0.26682022, -0.15776418, -0.11863651, -0.002802228, 0.28884634, -0.26341784, 0.3717576, 0.05564724, 0.1396777, -0.093476616, -0.19411619, -0.09813254},
		{-0.37407523, -0.19863394, 0.36883816, 0.19451912, -0.026730793, -0.26035637, -0.06540082, 0.102558956, -0.2309291, 0.14975478, 0.04800632, -0.35562083, 0.13055225, -0.09326942, -0.1463881, -0.15137577, -0.033199176, 0.052308187, -0.19963148, 0.08709718, 0.15636344, -0.050250772, 0.025203317, -0.015259152, 0.40350246, 0.0031367766, -0.06658539, 0.040400125, -0.027634807, 0.01308484, -0.11290837, -0.1647178},
	}

	if ok := mrv.milvusObj.InsertColumns(ctx, infoTable, NewMilvusInsertOption(lang_val,
		WithScalarColumns("id", idScalarList), WithScalarColumns("question", questionSliceStr),
		WithVectVolumns("question_vect", VectDimOnQuestion, mivus_interface.VectColumnF32, vectDataList))); !ok {
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
				IsAutoID:     false, //主键自动递增模式
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
