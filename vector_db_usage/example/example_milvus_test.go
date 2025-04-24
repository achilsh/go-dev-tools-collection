package example

import (
	"context"
	"fmt"
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
	existDataVect := []float32{
		-0.17596263,
		0.00033492805,
		-0.19089417,
		0.032562274,
		0.31517065,
		-0.37880203,
		-0.09877797,
		0.3172381,
		-0.3328588,
		-0.22385845,
		-0.05447146,
		0.004956845,
		-0.1614905,
		0.04045877,
		-0.10009884,
		-0.17803007,
		-0.17320603,
		-0.08528215,
		-0.2761189,
		0.027465444,
		0.026331222,
		-0.018104509,
		0.11416896,
		-0.0938965,
		0.13461371,
		-0.20915662,
		0.105497174,
		-0.037472457,
		-0.06742171,
		-0.15770018,
		-0.24143174,
		-0.2138658,
	}
	existRetdata := GetRetriever(RETRIEVE_VECT_ON_MILVUS).Retrieve(context.Background(), "", existDataVect)
	for _, v := range existRetdata {
		logger.Infof("码农： not filter, find exist item: %v", v)
	}

	logger.Infof("--------------------------")

	filter := fmt.Sprintf("answer not in [\"码农\"]")

	withFilterRetData := GetRetriever(
		RETRIEVE_VECT_ON_MILVUS,
	).Retrieve(context.Background(), filter, existDataVect) //"answer != '码农'"
	for _, v := range withFilterRetData {
		logger.Infof("filter ===> find exist item: %v", v)
	}

	logger.Infof("--------------------------en")

	// 检索一个和里面有一定关联的,比如：中餐怎么做，有什么菜谱推荐么
	relativeDataVect := []float32{
		0.042535435, -0.49152058, 0.028050063, -0.07119928, 0.16339009, -0.25582638, 0.08335227, 0.12404634, 0.2175261, 0.06055008, -0.014508388, -0.20291796, 0.10219552, -0.06696416, -0.21065167, -0.047998138, -0.08016057, 0.24035896, -0.019809972, 0.48832887, 0.10330034, -0.07917851, -0.03753307, 0.20966962, 0.15590188, 0.03808548, 0.014025032, 0.16731831, -0.17615685, -0.09059495, -0.013741155, -0.15344673,
	}
	existRetdata = GetRetriever(RETRIEVE_VECT_ON_MILVUS).Retrieve(context.Background(), "", relativeDataVect)
	for _, v := range existRetdata {
		logger.Infof("find exist item for relate record: %v", v)
	}

	//检索库中两个记录有关联的，比如： 中国改革开放最大的经济特区
	relativeDataVectWithTwos := []float32{
		0.08636823, -0.043513015, 0.41625148, 0.35441893, 0.06558197, -0.10623358, 0.06091164, 0.25390813, 0.02298987, -0.052130103, 0.19983755, 0.11143015, -0.4499305, 0.008707536, 0.25325033, 0.10498378, -0.14010993, -0.08189524, -0.09630091, 0.10393131, -0.1456354, -0.09577467, -0.16076463, 0.09340662, 0.00061616715, -0.13346621, -0.2775887, -0.14063616, -0.13945213, 0.019816017, 0.02476591, -0.16813192,
	}
	existRetdata = GetRetriever(RETRIEVE_VECT_ON_MILVUS).Retrieve(context.Background(), "", relativeDataVectWithTwos)
	for _, v := range existRetdata {
		logger.Infof("find exist item for relate record with two: %v", v)
	}

}

func InsertDemoList() {
	handle := GetRetriever(RETRIEVE_VECT_ON_MILVUS).(*MilvusRetrieveVect)
	handle.Insert(context.Background())
}
