package vector_db_usage

import "context"

type DbDetailer interface {
	GetName() string
	GetDetails() map[string]string
}

type Indexer interface {
	BuildCreateIndexOptions() (any, error)
}

type Schemaer interface {
	BuildSchema() (any, error)
}

type InsertOptionBuilder interface {
	BuildInsertOption() (any, bool)
}

type VectDBOper interface {
	//Init 比如创建client, 建立连接，
	Connect(ctx context.Context) bool
	//
	DisConnect(ctx context.Context) bool
	// CreateDB 初始化， 创建 db
	CreateDB(ctx context.Context, dbName string) bool
	// 列举该主机上的 db
	ListDB(ctx context.Context) []string
	// UsingDB 选择某个库
	UsingDB(ctx context.Context, dbName string) bool
	// GetDBDetail 获取库的信息
	GetDBDetail(ctx context.Context, dbName string) (DbDetailer, bool)
	// 创建 一个 collection
	CreateCollection(ctx context.Context, collectName string, idx Indexer, sch Schemaer) bool
	// 列举db下的所有 collection
	ListCollection(ctx context.Context) []string
	//
	LoadCollection(ctx context.Context, collName string, async bool) bool
	//
	ReleaseCollection(ctx context.Context, collName string) bool

	//
	InsertColumns(cctx context.Context, collectName string, insertOpter InsertOptionBuilder) bool
	//
	QueryVector(ctx context.Context, collectName string, searchOpt SearchVectOpter) bool
}

type SearchVectOpter interface {
	BuildSearchVectOpt() (any, bool)
}
