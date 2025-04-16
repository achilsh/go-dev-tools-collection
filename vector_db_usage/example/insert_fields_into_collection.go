package example

import (
	"reflect"

	logger "github.com/achilsh/go-dev-tools-collection/base-lib/log"
	models "github.com/achilsh/go-dev-tools-collection/vector_db_usage/example/models"
	mivus_interface "github.com/achilsh/go-dev-tools-collection/vector_db_usage/milvus_interface"
)

type VectDimTypeData struct {
	dim   int
	dType int
	data  any //slice of slice
}

type ScalarTypeData struct {
	vType reflect.Kind
	data  any //slice
}

type InsertColumnInfo struct {
	scalarColumns map[string]ScalarTypeData  //key is column
	vectColumns   map[string]VectDimTypeData //key is column
}

type MilvusInserter struct {
	lang      string
	insertOpt *mivus_interface.InsertOps
	*InsertColumnInfo
}
type OptionInsertColumn func(*InsertColumnInfo)

func WithScalarColumns(columnName string, data any) OptionInsertColumn {
	if reflect.TypeOf(data).Kind() != reflect.Slice || reflect.ValueOf(data).IsNil() {
		logger.Errorf("input data is not slice or is nil")
		return nil
	}

	if reflect.ValueOf(data).Len() <= 0 {
		logger.Errorf("data input slice is empty")
		return nil
	}

	return func(item *InsertColumnInfo) {
		_, ok := item.scalarColumns[columnName]
		if !ok {
			vType := reflect.TypeOf(reflect.ValueOf(data).Index(0)).Kind()
			item.scalarColumns[columnName] = ScalarTypeData{
				vType: vType,
				data:  data,
			}
			return
		}
	}
}

func WithVectVolumns(columnName string, dim int, vectItemType int, data any) OptionInsertColumn {
	if reflect.TypeOf(data).Kind() != reflect.Slice || reflect.ValueOf(data).IsNil() {
		logger.Errorf("input data is not slice or is nil")
		return nil
	}

	return func(item *InsertColumnInfo) {
		_, ok := item.vectColumns[columnName]
		if !ok {
			item.vectColumns[columnName] = VectDimTypeData{
				dim:   dim,
				dType: vectItemType,
				data:  data,
			}
			return
		}

	}
}

func NewMilvusInsertOption(lang string, opts ...OptionInsertColumn) *MilvusInserter {
	ret := &MilvusInserter{
		lang:      lang,
		insertOpt: mivus_interface.NewInsertOpts(),
		InsertColumnInfo: &InsertColumnInfo{
			scalarColumns: make(map[string]ScalarTypeData),
			vectColumns:   make(map[string]VectDimTypeData),
		},
	}
	for _, opt := range opts {
		opt(ret.InsertColumnInfo)
	}
	return ret
}

func (mio *MilvusInserter) BuildInsertOption() (any, bool) {
	if mio.insertOpt == nil {
		logger.Errorf("not register collection info.")
		return nil, false
	}

	collectName := models.QuestionVectorCollection{}.TableName(mio.lang)
	logger.Infof("collection name: %v", collectName)

	if mio.InsertColumnInfo == nil {
		logger.Errorf("not init insert column info")
		return nil, false
	}
	//
	var insertOptList []mivus_interface.InsertOptions
	for columnName, item := range mio.scalarColumns {
		if len(columnName) <= 0 {
			continue
		}
		//
		switch item.vType {
		case reflect.Bool:
			itemData, ok := (item.data).([]bool)
			if ok {
				insertOptList = append(insertOptList, mivus_interface.WithScalarColumnNums(columnName, itemData))
			}

		case reflect.Int8:
			itemData, ok := (item.data).([]int8)
			if ok {
				insertOptList = append(insertOptList, mivus_interface.WithScalarColumnNums(columnName, itemData))
			}

		case reflect.Uint8:
			itemData, ok := (item.data).([]uint8)
			if ok {
				insertOptList = append(insertOptList, mivus_interface.WithScalarColumnNums(columnName, itemData))
			}

		case reflect.Int16:
			itemData, ok := (item.data).([]int16)
			if ok {
				insertOptList = append(insertOptList, mivus_interface.WithScalarColumnNums(columnName, itemData))
			}

		case reflect.Uint16:
			itemData, ok := (item.data).([]uint16)
			if ok {
				insertOptList = append(insertOptList, mivus_interface.WithScalarColumnNums(columnName, itemData))
			}

		case reflect.Int32:
			itemData, ok := (item.data).([]int32)
			if ok {
				insertOptList = append(insertOptList, mivus_interface.WithScalarColumnNums(columnName, itemData))
			}

		case reflect.Uint32:
			itemData, ok := (item.data).([]uint32)
			if ok {
				insertOptList = append(insertOptList, mivus_interface.WithScalarColumnNums(columnName, itemData))
			}
		case reflect.Int:
			itemData, ok := (item.data).([]int)
			if ok {
				insertOptList = append(insertOptList, mivus_interface.WithScalarColumnNums(columnName, itemData))
			}

		case reflect.Uint:
			itemData, ok := (item.data).([]uint)
			if ok {
				insertOptList = append(insertOptList, mivus_interface.WithScalarColumnNums(columnName, itemData))
			}

		case reflect.Uint64:
			itemData, ok := (item.data).([]uint64)
			if ok {
				insertOptList = append(insertOptList, mivus_interface.WithScalarColumnNums(columnName, itemData))
			}

		case reflect.Int64:
			itemData, ok := (item.data).([]int64)
			if ok {
				insertOptList = append(insertOptList, mivus_interface.WithScalarColumnNums(columnName, itemData))
			}
			logger.Debugf("is int64 type.")

		case reflect.Float32:
			itemData, ok := (item.data).([]float32)
			if ok {
				insertOptList = append(insertOptList, mivus_interface.WithScalarColumnNums(columnName, itemData))
			}
		case reflect.Float64:
			itemData, ok := (item.data).([]float64)
			if ok {
				insertOptList = append(insertOptList, mivus_interface.WithScalarColumnNums(columnName, itemData))
			}
		case reflect.String:
			itemData, ok := (item.data).([]string)
			if ok {
				insertOptList = append(insertOptList, mivus_interface.WithScalarColumnNums(columnName, itemData))
			}
		default:
			itemData, ok := (item.data).([]int64)
			if ok {
				insertOptList = append(insertOptList, mivus_interface.WithScalarColumnNums(columnName, itemData))
			}
		}

	}
	var insertVectOptions []mivus_interface.InsertOptions
	for columnName, item := range mio.vectColumns {
		if len(columnName) <= 0 {
			continue
		}
		if item.dType == mivus_interface.VectColumnByte {
			itemData, ok := (item.data).([][]byte)
			if ok {
				insertVectOptions = append(insertVectOptions, mivus_interface.WithVectColumnValue(columnName, item.dim, item.dType, itemData))
			}
		} else {
			itemData, ok := (item.data).([][]float32)
			if ok {
				insertVectOptions = append(insertVectOptions, mivus_interface.WithVectColumnValue(columnName, item.dim, item.dType, itemData))
			}
		}
	}

	var insertOpts = append(insertOptList, insertVectOptions...)
	mio.insertOpt.Register(collectName, insertOpts...)

	data, ok := mio.insertOpt.BuildInsertOption()
	if !ok {
		logger.Errorf("build insert record fail")
	}
	return data, ok
}

func (mio *MilvusInserter) BuildInsertOptionV2() (any, bool) {
	if mio.insertOpt == nil {
		logger.Errorf("not register collection info.")
		return nil, false
	}

	collectName := models.QuestionVectorCollection{}.TableName(mio.lang)
	logger.Infof("collection name: %v", collectName)

	mio.insertOpt.Register(collectName,
		// mivus_interface.WithScalarColumnNums("id", []int64{
		// 	1000,
		// 	2000,
		// 	3000,
		// 	4000,
		// 	5000,
		// 	6000,
		// }),
		mivus_interface.WithScalarColumnNums("question", []string{
			"深圳市是中国广东省一个市，也是中国的经济特区",
			"香港是中国的特别行政区",
			"我是谁，我是一个it 码农。",
			"春天是一个播种的季节",
			"一年小学生学古诗100首。",
			"中国有八大菜系，湘菜，川菜，鲁菜等.",
		}),
		mivus_interface.WithVectColumnValue("question_vect", VectDimOnQuestion, mivus_interface.VectColumnF32, [][]float32{
			{0.14181516, 0.053447783, 0.27755734, 0.425333, 0.12168437, 0.10610832, 0.02128352, 0.21277897, -0.118872814, -0.28497985, -0.15902191, 0.24786726, -0.056231227, 0.18140194, 0.012834778, 0.040092863, -0.35830536, -0.37157595, 0.074899994, 0.15587296, 0.14147776, 0.07073888, 0.07698055, 0.0009163933, -0.23774563, -0.05769324, 0.01973716, -0.13124369, -0.021030478, 0.13664187, 0.020552514, -0.05544399},
			{0.06335796, 0.057832558, 0.20333485, 0.5510669, 0.06773224, -0.027235635, -0.032899175, 0.19523092, -0.047127087, -0.15848699, -0.06856105, 0.17552365, -0.14706782, 0.09577366, 0.11769109, -0.037895057, -0.46376553, 0.052491333, 0.19559929, 0.044709723, -0.031909205, 0.007867944, 0.14927799, 0.060917575, -0.20020379, -0.03754972, -0.1639203, -0.21549073, -0.038240395, 0.2149382, 0.16705136, -0.13647747},
			{-0.17596263, 0.00033492805, -0.19089417, 0.032562274, 0.31517065, -0.37880203, -0.09877797, 0.3172381, -0.3328588, -0.22385845, -0.05447146, 0.004956845, -0.1614905, 0.04045877, -0.10009884, -0.17803007, -0.17320603, -0.08528215, -0.2761189, 0.027465444, 0.026331222, -0.018104509, 0.11416896, -0.0938965, 0.13461371, -0.20915662, 0.105497174, -0.037472457, -0.06742171, -0.15770018, -0.24143174, -0.2138658},
			{0.34291446, 0.015021714, 0.09629263, 0.055509117, 0.22280477, -0.28606084, 0.058550276, -0.11946945, -0.076060936, -0.21678647, -0.014277432, 0.2484145, -0.15609138, 0.007362801, 0.08137496, -0.21076818, -0.23663403, -0.029051052, -0.24406084, 0.035309434, 0.10371946, 0.14174993, 0.08905788, -0.2087194, 0.11639628, -0.05451674, 0.30552423, 0.30526814, 0.23548159, -0.05896643, 0.048370402, -0.25724986},
			{0.018243162, -0.16492714, 0.08501537, -0.19250453, 0.14263242, -0.18193917, -0.15660019, 0.08707473, 0.0025154299, -0.011242488, -0.084433384, -0.27631116, -0.30012798, 0.30173966, 0.28651837, -0.07637505, -0.089089304, 0.03988873, -0.0081982305, -0.023011006, 0.26682022, -0.15776418, -0.11863651, -0.002802228, 0.28884634, -0.26341784, 0.3717576, 0.05564724, 0.1396777, -0.093476616, -0.19411619, -0.09813254},
			{-0.37407523, -0.19863394, 0.36883816, 0.19451912, -0.026730793, -0.26035637, -0.06540082, 0.102558956, -0.2309291, 0.14975478, 0.04800632, -0.35562083, 0.13055225, -0.09326942, -0.1463881, -0.15137577, -0.033199176, 0.052308187, -0.19963148, 0.08709718, 0.15636344, -0.050250772, 0.025203317, -0.015259152, 0.40350246, 0.0031367766, -0.06658539, 0.040400125, -0.027634807, 0.01308484, -0.11290837, -0.1647178},
		}),
	)

	data, ok := mio.insertOpt.BuildInsertOption()
	if !ok {
		logger.Errorf("build insert record fail")
	}
	return data, ok

}
