package mivus_interface

import (
	"context"
	"log"

	dbUsage "github.com/achilsh/go-dev-tools-collection/vector_db_usage"
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

// ops: WithAddress()
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

// 插入接口
func (m *MilvusVectOp) InsertColumns(ctx context.Context, collectName string, insertOpter dbUsage.InsertOptionBuilder) bool {
	if !m.checkClientIsOK() {
		return false
	}
	opts, ok := insertOpter.BuildInsertOption()
	if !ok {
		log.Printf("build insert option fail")
		return false
	}

	insertOpt, ok := opts.(client.InsertOption)
	if !ok || insertOpt == nil {
		log.Printf("get build insert option is nil")
		return false
	}

	retInsert, err := m.cli.Insert(ctx, insertOpt)
	if err != nil {
		log.Printf("insert to vec db fail, err: %v, collectName: %v", err, collectName)
		return false
	}
	log.Printf("insert succ, insert nums: %v", retInsert.InsertCount)
	return true
}

func (m *MilvusVectOp) QueryVector(ctx context.Context, collectName string, searchOpt dbUsage.SearchVectOpter) (any, bool) {
	if !m.checkClientIsOK() {
		return nil, false
	}

	optRet, succ := searchOpt.BuildSearchVectOpt()
	if !succ {
		log.Printf("query vector fail")
		return nil, false
	}

	opt, ok := (optRet).(client.SearchOption)
	if !ok {
		log.Printf("search option not get")
		return nil, false
	}

	searchRet, err := m.cli.Search(ctx, opt)
	if err != nil {
		log.Printf("query fail, err: %v", err)
		return nil, false
	}
	return searchRet, true
}
