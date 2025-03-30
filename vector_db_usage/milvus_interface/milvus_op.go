package mivus_interface

import (
	"context"
	"log"
	"reflect"

	dbUsage "github.com/achilsh/go-dev-tools-collection/vector_db_usage"
	"github.com/milvus-io/milvus/client/v2/column"
	"github.com/milvus-io/milvus/client/v2/entity"
	client "github.com/milvus-io/milvus/client/v2/milvusclient"
)

type MilvusDbDetail struct {
	dbDetail *entity.Database
}

func (dDetail *MilvusDbDetail) GetName() string {
	return dDetail.dbDetail.Name
}
func (dDetail *MilvusDbDetail) GetDetails() map[string]string {
	return dDetail.dbDetail.Properties
}

type Options func(op *MilvusVectOption)

type MilvusVectOption struct {
	client.ClientConfig
}

func WithAddress(addr string) Options {
	return func(op *MilvusVectOption) {
		op.Address = addr
	}
}

type MilvusVectOp struct {
	// 配置项，静态的
	cfg *MilvusVectOption
	//连接的客户端
	cli *client.Client
}

func NewVectMilvusOpInst(ops ...Options) *MilvusVectOp {
	config := &MilvusVectOption{}
	for _, op := range ops {
		op(config)
	}
	return &MilvusVectOp{
		cfg: config,
	}
}
func (m *MilvusVectOp) checkClientIsOK() bool {
	if m == nil || m.cli == nil {
		return false
	}
	return true
}

// Connect 初始化，比如创建 client, 建立连接， 创建 db
func (m *MilvusVectOp) Connect(ctx context.Context) bool {
	if m.cfg == nil {
		return false
	}

	cli, err := client.New(ctx, &m.cfg.ClientConfig)
	if err != nil {
		log.Fatalf("create client to milvus fail, err: %v", err)
		return false
	}
	if cli == nil {
		log.Fatalf("create client to milvus, ret client is nil")
		return false
	}
	m.cli = cli
	return true
}
func (m *MilvusVectOp) DisConnect(ctx context.Context) bool {
	if !m.checkClientIsOK() {
		return false
	}
	if err := m.cli.Close(ctx); err != nil {
		return false
	}
	return true
}

// CreateDB 初始化， 创建 db
func (m *MilvusVectOp) CreateDB(ctx context.Context, dbName string) bool {
	if !m.checkClientIsOK() {
		return false
	}
	err := m.cli.CreateDatabase(ctx, client.NewCreateDatabaseOption(dbName))
	if err != nil {
		log.Fatalf("create db fail, dbName: %v, err: %v", dbName, err)
		return false
	}

	return true
}

// 列举该连接的主机的 已创建的 db
func (m *MilvusVectOp) ListDB(ctx context.Context) []string {
	if !m.checkClientIsOK() {
		return nil
	}

	dbNameLists, err := m.cli.ListDatabase(ctx, client.NewListDatabaseOption())
	if err != nil {
		log.Printf("list db list on this connect fail: %v", err)
		return nil
	}
	return dbNameLists
}

// UsingDB 选择某个库
func (m *MilvusVectOp) UsingDB(ctx context.Context, dbName string) bool {
	if !m.checkClientIsOK() {
		return false
	}
	err := m.cli.UseDatabase(ctx, client.NewUseDatabaseOption(dbName))
	if err != nil {
		log.Printf("use db: %v fail, err: %v", dbName, err)
		return false
	}
	return true
}

// GetDBDetail 获取库的信息
func (m *MilvusVectOp) GetDBDetail(ctx context.Context, dbName string) (dbUsage.DbDetailer, bool) {
	if !m.checkClientIsOK() {
		return nil, false
	}

	dbDetail, err := m.cli.DescribeDatabase(ctx, client.NewDescribeDatabaseOption(dbName))
	if err != nil {
		log.Printf("get db: %v detail fail: %v", dbName, err)
		return nil, false
	}
	if dbDetail == nil {
		log.Printf("get db detail ret is nil")
		return nil, false
	}

	return &MilvusDbDetail{
		dbDetail: dbDetail,
	}, true
}

func (m *MilvusVectOp) CreateCollection(ctx context.Context, collectName string, idx dbUsage.Indexer, sch dbUsage.Schemaer) bool {

	if !m.checkClientIsOK() {
		return false
	}
	//
	indexOptions, err := idx.BuildCreateIndexOptions()
	if err != nil {
		log.Printf("get index option fail, err: %v", err)
		return false
	}

	indexOptionsMilvus, ok := indexOptions.([]client.CreateIndexOption)
	if !ok {
		log.Printf("build not list of client.CreateIndexOption")
		return false
	}
	//
	schemaItemTmp, err := sch.BuildSchema()
	if err != nil {
		log.Printf("get schema item fail, err: %v", err)
		return false
	}
	//
	schemaItem, ok := schemaItemTmp.(*entity.Schema)
	if !ok || nil == schemaItem {
		log.Printf("is not entity schema or is empty")
		return false
	}

	collectOptions := client.NewCreateCollectionOption(collectName, schemaItem)
	if len(indexOptionsMilvus) > 0 {
		err = m.cli.CreateCollection(ctx, collectOptions.WithIndexOptions(indexOptionsMilvus...))
	} else {
		err = m.cli.CreateCollection(ctx, collectOptions)
	}

	if err != nil {
		log.Printf("create collection fail, err: %v", err)
		return false
	}
	return true
}

func (m *MilvusVectOp) ListCollection(ctx context.Context) []string {

	if !m.checkClientIsOK() {
		return nil
	}
	ret, err := m.cli.ListCollections(ctx, client.NewListCollectionOption())
	if err != nil {
		log.Printf("list collections fail, err: %v", err)
		return nil
	}
	return ret
}

func (m *MilvusVectOp) LoadCollection(ctx context.Context, collName string, async bool) bool {
	if !m.checkClientIsOK() {
		return false
	}

	lTask, err := m.cli.LoadCollection(ctx, client.NewLoadCollectionOption(collName))
	if err != nil {
		log.Printf("load collection: %v fail, err: %v", collName, err)
		return false
	}

	if !async {
		if err := lTask.Await(ctx); err != nil {
			log.Printf("wait load collection fail, err: %v", err)
			return false
		}
		return true
	}
	return true
}

func (m *MilvusVectOp) ReleaseCollection(ctx context.Context, collName string) bool {
	if !m.checkClientIsOK() {
		return false
	}

	if err := m.cli.ReleaseCollection(ctx, client.NewReleaseCollectionOption(collName)); err != nil {
		return false
	}
	return true
}

type columnItem[T any | string] struct {
	columnName string
	columnData []T
}

func NewcolumnItem[T any]() *columnItem[T] {
	return &columnItem[T]{}
}

type columnVecItem[T any] struct {
	columnName string
	columnData [][]T
	dim        int
}

func NewcolumnVectItem[T any]() *columnVecItem[T] {
	return &columnVecItem[T]{}
}

type miluvsInsertOption struct {
	collectName string
	//
	columnScar map[int]any // value: []*columnItem[T], T: int64, int32, int16, int8, bool, string
	//
	// columnInt64   []*columnItem[int64]
	// columnVarChar []*columnItem[string]
	// columnInt32   []*columnItem[int32]
	// columnInt16   []*columnItem[int16]
	// columnInt8    []*columnItem[int8]
	// columnBool    []*columnItem[bool]

	columnVect map[int]any /// []*columnVecItem[T]; T:float32, byte, int8
	// columnVectF32  []*columnVecItem[float32]
	// columnVectF16  []*columnVecItem[float32]
	// columnVectBF16 []*columnVecItem[float32]
	// columnVectBin  []*columnVecItem[byte]
	// columnVectI8   []*columnVecItem[int8]
}
type InsertOptions func(*miluvsInsertOption)

func WithCollectName(collName string) InsertOptions {
	return func(item *miluvsInsertOption) {
		item.collectName = collName
	}
}

func WithColumnNums[T any](columnName string, data []T) InsertOptions {
	return func(item *miluvsInsertOption) {
		var colItem *columnItem[T] = &columnItem[T]{
			columnName: columnName,
			columnData: data,
		}
		//
		tmp := new(T)
		vType := int(reflect.ValueOf(tmp).Elem().Kind())
		item.columnScar[vType] = colItem
	}

}

// func WithColumnI64(columnName string, data []int64) InsertOptions {
// 	return func(item *miluvsInsertOption) {
// 		var colItem *columnItem[int64] = &columnItem[int64]{
// 			columnName: columnName,
// 			columnData: data,
// 		}

//			item.columnInt64 = append(item.columnInt64, colItem)
//		}
//	}
func WithColumnVect[T any](columnName string, dim int, dataType int, data [][]T) InsertOptions {
	return func(item *miluvsInsertOption) {
		var colItem *columnVecItem[T] = &columnVecItem[T]{
			columnName: columnName,
			columnData: data,
			dim:        dim,
		}
		item.columnVect[dataType] = colItem
	}
}

type InsertOps struct {
	options      *miluvsInsertOption
	insertOption column.Column
}

func (iops *InsertOps) BuildInsertOption() (any, bool) {
	for k, v := range iops.options.columnScar {
		iops.insertOption.
	}
	iops.insertOption.
}

func (iops *InsertOps) Register(collectName string, ops ...InsertOptions) {
	WithCollectName(collectName)(iops.options)
	for _, op := range ops {
		if op == nil {
			continue
		}
		op(iops.options)
	}
	iops.insertOption = client.NewColumnBasedInsertOption(iops.options.collectName)
}

func (m *MilvusVectOp) InsertColumns(ctx context.Context, collectName string) bool {
	if !m.checkClientIsOK() {
		return false
	}
	//

}

func (m *MilvusVectOp) QueryVector() bool {

}
