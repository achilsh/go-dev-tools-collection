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

// CreateDB 初始化， 创建 db
func (m *MilvusVectOp) CreateDB(ctx context.Context, dbName string) bool {
	if m.checkClientIsOK() == false {
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
	if m.checkClientIsOK() == false {
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
	if m.checkClientIsOK() == false {
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
	if m.checkClientIsOK() == false {
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
