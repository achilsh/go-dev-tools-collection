package example

import (
	"context"
	"testing"

	logger "github.com/achilsh/go-dev-tools-collection/base-lib/log"
	"github.com/achilsh/go-dev-tools-collection/demo-service/service/utils/mock/mock_log"
)

func TestCallOne(t *testing.T) {
	mock_log.LoggerMock()
	logger.Infof("ssss")
}

// 测试 collection 创建，查询； insert iterms, query by vector
func TestRetrieve(t *testing.T) {
	mock_log.LoggerMock()
	//先插入向量到向量数据库中，

	//检索开始
	InitAllRetriever()
	//插入记录
	InsertDemoList()

	//开始检索：检索一个库中有的记录: 我是谁，我是一个it 码农。
	existDataVect := []float32{-0.1759626269340515, 0.00033492804504930973, -0.19089417159557343, 0.03256227448582649, 0.31517064571380615, -0.37880203127861023, -0.09877797216176987, 0.31723809242248535, -0.3328588008880615, -0.22385844588279724, -0.0544714592397213, 0.004956845194101334, -0.16149049997329712, 0.04045876860618591, -0.10009884089231491, -0.1780300736427307, -0.17320603132247925, -0.08528214693069458, -0.2761189043521881, 0.027465444058179855, 0.0263312216848135, -0.01810450851917267, 0.11416895687580109, -0.09389650076627731, 0.13461370766162872, -0.20915661752223969, 0.10549717396497726, -0.03747245669364929, -0.06742171198129654, -0.15770018100738525, -0.24143174290657043, -0.21386580169200897}
	existRetdata := GetRetriever(RETRIEVE_VECT_ON_MILVUS).Retrieve(context.Background(), existDataVect)
	for _, v := range existRetdata {
		logger.Infof("find exist item: %v", v)
	}
	// 检索一个和里面有一定关联的,比如：中餐怎么做，有什么菜谱推荐么
	relativeDataVect := []float32{
		0.042535435, -0.49152058, 0.028050063, -0.07119928, 0.16339009, -0.25582638, 0.08335227, 0.12404634, 0.2175261, 0.06055008, -0.014508388, -0.20291796, 0.10219552, -0.06696416, -0.21065167, -0.047998138, -0.08016057, 0.24035896, -0.019809972, 0.48832887, 0.10330034, -0.07917851, -0.03753307, 0.20966962, 0.15590188, 0.03808548, 0.014025032, 0.16731831, -0.17615685, -0.09059495, -0.013741155, -0.15344673,
	}
	existRetdata = GetRetriever(RETRIEVE_VECT_ON_MILVUS).Retrieve(context.Background(), relativeDataVect)
	for _, v := range existRetdata {
		logger.Infof("find exist item for relate record: %v", v)
	}

	//检索库中两个记录有关联的，比如： 中国改革开放最大的经济特区
	relativeDataVectWithTwos := []float32{
		0.08636823, -0.043513015, 0.41625148, 0.35441893, 0.06558197, -0.10623358, 0.06091164, 0.25390813, 0.02298987, -0.052130103, 0.19983755, 0.11143015, -0.4499305, 0.008707536, 0.25325033, 0.10498378, -0.14010993, -0.08189524, -0.09630091, 0.10393131, -0.1456354, -0.09577467, -0.16076463, 0.09340662, 0.00061616715, -0.13346621, -0.2775887, -0.14063616, -0.13945213, 0.019816017, 0.02476591, -0.16813192,
	}
	existRetdata = GetRetriever(RETRIEVE_VECT_ON_MILVUS).Retrieve(context.Background(), relativeDataVectWithTwos)
	for _, v := range existRetdata {
		logger.Infof("find exist item for relate record with two: %v", v)
	}

}

func InsertDemoList() {
	handle := GetRetriever(RETRIEVE_VECT_ON_MILVUS).(*MilvusRetrieveVect)
	handle.Insert(context.Background())
}
