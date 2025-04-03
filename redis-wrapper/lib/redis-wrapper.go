package lib

import (
	"time"

	"github.com/redis/go-redis/v9"
	redisV9 "github.com/redis/go-redis/v9"
)

var RETNil = redis.Nil

// RedisOps redis的操作接口
type RedisOps interface {
	GetRedisRow() *redisV9.Client
	// SetAtomic 增加一个带时间戳的key
	SetEx(key string, val any, tm time.Duration) error
	// Get 获取 key的值
	Get(key string) (any, error)
	// HSet fieldValue is map[string]any{"aaa":1, "bbb": 123.123, "ccc":"isfa"}
	HSet(key string, fieldValues map[string]any) error
	// HSetStruct fieldValues must be struct with tag:
	//type ExampleUser struct {
	//		Name string `redis:"name"`
	//		Age  int    `redis:"age"`
	//	}
	// you can use HSetStruct("aaa",ExampleUser{})
	HSetStruct(key string, fieldValues any) error
	// HGet 根据key 和 field 获取对应的value.
	HGet(key string, field string) (string, error)
	// HDel 根据key 和多个 field name 删除对应的 field 和 value.
	HDel(key string, fields ...string) error
	// HGetAll 根据 key 获取 所有的 field and value
	HGetAll(key string) (map[string]string, error)
	// SAdd 增加集合中的元素, 其中 members 可以是任何数据列表
	SAdd(key string, members ...any) error
	// SMembers 获取集合中的元素列表
	SMembers(key string) ([]string, error)
	// SRem 删除某些成员
	SRem(key string, members ...any) error
	// DeleteKey 删除redis key
	DeleteKey(key string) error
	// ExpireKey 设置key 有效期
	ExpireKey(key string, ttl time.Duration) error
	// GenUniqueKey 根据 key 产生全局唯一key. 并设置有效时间戳:
	// 假如 key 存在，则返回原有的，并更新有效时间戳，
	// 假如 key 不存在，则创建一个新的值，设置 有效时间戳
	GenUniqueKey(key string, expireTm time.Duration) (string, error)
	// GenUniqueKeyWithGenVal 创建自定义value创建函数的 key
	GenUniqueKeyWithGenVal(key string, expireTm time.Duration, genFun func() string) (string, error)
	// ZADDAtomic 增加一个元素，当member个数超过 maxItemNums 则删除最早记录，然后添加新纪录； 整个过程是原子操作
	ZADDAtomic(key string, score float64, member any, maxItemNums int32) error
	// ZRange 获取索引 index [start,stop] 之间的元素
	ZRange(key string, start, stop int64) ([]string, error)
	// ZDelMemberByScore 根据score 范围来删除 members
	ZDelMemberByScore(key string, beginScore, endScore float64) (int64, error)
	// CheckKeyExist 检查key 是否存在，不存在则创建并返回false, 存在返回true
	CheckKeyExist(key string, tm time.Duration) (bool, error)
	// CheckKeyExistCheckKeyExistFresh 检查key 是否存在，不存在则创建并返回false, 存在并刷新时间返回true
	CheckKeyExistFresh(key string, tm time.Duration) (bool, error)
	// CheckKeyExistNoCreate 检查key 是否存在
	CheckKeyExistNoCreate(key string) (bool, error)
	// Lock 获取分布式锁
	Lock(key string, ttl time.Duration, globalId string) error
	// UnLock 释放分布式锁
	UnLock(key string, globalId string) error
	// RefreshTTL 刷新持有时间
	RefreshTTL(key string, globalId string, ttl time.Duration) error
	Client() *redisV9.Client
}
