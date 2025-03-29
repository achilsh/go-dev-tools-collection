package vector_db_usage

import "context"

type DbDetailer interface {
	GetName() string
	GetDetails() map[string]string
}
type VectOper interface {
	//Init 比如创建client, 建立连接，
	Connect(ctx context.Context) bool
	// CreateDB 初始化， 创建 db
	CreateDB(ctx context.Context, dbName string) bool
	// 列举该主机上的 db
	ListDB(ctx context.Context) []string
	// UsingDB 选择某个库
	UsingDB(ctx context.Context, dbName string) bool
	// GetDBDetail 获取库的信息
	GetDBDetail(ctx context.Context, dbName string) (DbDetailer, bool)
	// 创建 一个 collection
	CreateCollection() bool
	// 列举db下的所有 collection
	ListCollection() bool
	// 选择某个集合
	UsingCollection(string) bool
	//
	LoadCollection() bool
	//
	ReleaseCollection() bool

	//
	InsertVectors() bool
	//
	QueryVector() bool
	//
}
