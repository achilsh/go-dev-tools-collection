package example

import (
	"context"

	logger "github.com/achilsh/go-dev-tools-collection/base-lib/log"
)

const (
	RETRIEVE_VECT_ON_MILVUS = 1
)

var (
	// 在 retriever file 中init() 设置值内容，
	// 在 项目入口 初始化出使用 各个 retriever。
	similarRetrieverMap = make(map[int]SimilarRetriever)
)

type SimilarRetriever interface {
	Retrieve(ctx context.Context, filter string, targetVect []float32) []string
	Init() bool
}

// 对外提供接口
func GetRetriever(nType int) SimilarRetriever {
	item, ok := similarRetrieverMap[nType]
	if !ok {
		panic("not implement retrieve")
	}
	return item
}

func registerRetriever(nType int, sim SimilarRetriever) {
	similarRetrieverMap[nType] = sim
}

func InitAllRetriever() {
	for k, item := range similarRetrieverMap {
		if item == nil {
			continue
		}
		if !item.Init() {
			logger.Errorf("init retrieve: %v fail.", k)
		}
	}
}
