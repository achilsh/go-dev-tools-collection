package redis

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"time"

	redisV9 "github.com/redis/go-redis/v9"
	"github.com/sony/sonyflake"
)

var snowflake *sonyflake.Sonyflake = nil

func init() {
	snowflake, _ = sonyflake.New(sonyflake.Settings{})
}

// RedisOption 设置redis option的函数对象
type RedisOption func(*redisV9.Options)

// WithRedisAddrOpt 设置ip:port, eg: localhost:6379
func WithRedisAddrOpt(addr string) RedisOption {
	return func(options *redisV9.Options) {
		if options == nil {
			return
		}
		options.Addr = addr
	}
}

// WithRedisDB 选择redis db
func WithRedisDB(db int) RedisOption {
	return func(options *redisV9.Options) {
		if options == nil {
			return
		}
		if db <= 0 || db > 15 {
			db = 0
		}
		options.DB = db
	}
}

// WithRedisPasswd 设置 password
func WithRedisPasswd(pw string) RedisOption {
	return func(options *redisV9.Options) {
		if options == nil {
			return
		}
		if pw == "" {
			return
		}
		options.Password = pw
	}
}

// WithRedisDialTimeout 设置连接 耗时配置
func WithRedisDialTimeout(tm int32) RedisOption {
	if tm <= 0 {
		tm = 10
	}
	return func(options *redisV9.Options) {
		if options == nil {
			return
		}
		options.DialTimeout = time.Duration(tm) * time.Second
	}
}

// WithRedisReadTimeout 设置读数据超时配置
func WithRedisReadTimeout(tm int32) RedisOption {
	if tm <= 0 {
		tm = 30
	}
	return func(options *redisV9.Options) {
		if options == nil {
			return
		}
		options.ReadTimeout = time.Duration(tm) * time.Second
	}
}

// WithRedisWriteTimeout 设置写数据的耗时
func WithRedisWriteTimeout(tm int32) RedisOption {
	if tm <= 0 {
		tm = 30
	}
	return func(options *redisV9.Options) {
		if options == nil {
			return
		}
		options.ReadTimeout = time.Duration(tm) * time.Second
	}
}

// WithRedisPoolSize 设置client 连接redis 并发最大值
func WithRedisPoolSize(ps int) RedisOption {
	if ps <= 0 {
		ps = 50000
	}
	return func(options *redisV9.Options) {
		if options == nil {
			return
		}
		options.PoolSize = ps
	}
}

// NewRedisClient 创建一个 redis 连接句柄
func NewRedisClient(opts ...RedisOption) RedisOps {
	defaultRedisOption := &redisV9.Options{
		Addr:         ":6379",
		PoolSize:     3 * 10000,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}
	for _, opt := range opts {
		opt(defaultRedisOption)
	}

	ret := &RedisClient{
		client: redisV9.NewClient(defaultRedisOption),
	}

	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	_, err := ret.client.Ping(ctx).Result()
	if err != nil {
		fmt.Printf("ping redis fail, %v", err.Error())
		return nil
	} else {
		fmt.Printf("ping redis succ.")
	}
	return ret
}

// RedisClient redis 在 c2 的对象，封装了redis 底层版本。
type RedisClient struct {
	client *redisV9.Client
}

func (srds *RedisClient) Client() *redisV9.Client {
	return srds.client
}

func (srds *RedisClient) GetRedisRow() *redisV9.Client {
	if srds == nil {
		return nil
	}
	return srds.client
}
func (srds *RedisClient) checkParamInvalid() error {
	if srds == nil || srds.client == nil {
		return errors.New("redis client not init")
	}
	return nil
}

// HSet fieldValue is map[string]any{"field1":1, "field2": 123.123, "field3":"isfa"}
func (srds *RedisClient) HSet(key string, fieldValues map[string]any) error {
	if err := srds.checkParamInvalid(); err != nil {
		return err
	}
	if len(key) == 0 {
		return errors.New("hset key is empty")
	}
	if len(fieldValues) == 0 {
		return nil
	}

	if err := srds.client.HSet(context.Background(), key, fieldValues).Err(); err != nil {
		return fmt.Errorf("hset fail, key: %v, field value: %+v, err: %v", key, fieldValues, err)
	}
	return nil
}

// HSetStruct fieldValues must be struct with tag:
//
//	type ExampleUser struct {
//			Name string `redis:"name"`
//			Age  int    `redis:"age"`
//		}
//
// you can use HSetStruct("aaa",ExampleUser{})
func (srds *RedisClient) HSetStruct(key string, fieldValues any) error {
	if err := srds.checkParamInvalid(); err != nil {
		return err
	}

	fvValue := reflect.ValueOf(fieldValues)
	kindFV := fvValue.Kind()
	if kindFV == reflect.Ptr {
		fvValue = fvValue.Elem()
		kindFV = fvValue.Kind()
	}

	if kindFV != reflect.Struct {
		return fmt.Errorf("%+v type not struct or pointer to struct", fieldValues)
	}

	hasRedisTag := false
	fvType := fvValue.Type()
	for i := 0; i < fvType.NumField(); i++ {
		field := fvType.Field(i)
		_, ok := field.Tag.Lookup("redis")
		if ok {
			hasRedisTag = true
			break
		}
	}

	if hasRedisTag == false {
		return fmt.Errorf("%+v type not have redis tag", fieldValues)
	}

	if err := srds.client.HSet(context.Background(), key, fieldValues).Err(); err != nil {
		return fmt.Errorf("hset key: %v, value: %+v fail, err: %v", err, key, fieldValues)
	}
	return nil
}

// HGet 根据key 和 field 获取对应的value.
func (srds *RedisClient) HGet(key string, field string) (string, error) {
	if err := srds.checkParamInvalid(); err != nil {
		return "", err
	}
	ret := srds.client.HGet(context.Background(), key, field)
	if ret == nil {
		return "", fmt.Errorf("hget key: %v, ret is nil", key)
	}

	v, err := ret.Result()
	if errors.Is(err, redisV9.Nil) {
		return "", fmt.Errorf("key: %v not exit field: %v", key, field)
	}
	return v, nil
}

// HDel 根据key 和多个 field name 删除对应的 field 和 value.
func (srds *RedisClient) HDel(key string, fields ...string) error {
	if err := srds.checkParamInvalid(); err != nil {
		return err
	}

	_, err := srds.client.HDel(context.Background(), key, fields...).Result()
	if err != nil {
		return fmt.Errorf("hdel key: %v, field: %v, fail, err: %v", key, fields, err)
	}
	return nil
}

// HGetAll 根据 key 获取 所有的 field and value
func (srds *RedisClient) HGetAll(key string) (map[string]string, error) {
	if err := srds.checkParamInvalid(); err != nil {
		return nil, err
	}

	v, err := srds.client.HGetAll(context.Background(), key).Result()
	if err != nil {
		return nil, fmt.Errorf("hgetall key: %v, err: %v", key, err)
	}
	return v, nil
}

// SAdd 增加集合中的元素, 其中 members 可以是任何数据列表
func (srds *RedisClient) SAdd(key string, members ...any) error {
	if err := srds.checkParamInvalid(); err != nil {
		return err
	}

	_, e := srds.client.SAdd(context.Background(), key, members).Result()
	if e != nil {
		return fmt.Errorf("sadd key: %v, members: %+v, fail: %v", key, members, e)
	}
	return nil
}

// SMembers 获取集合中的元素列表
func (srds *RedisClient) SMembers(key string) ([]string, error) {
	if err := srds.checkParamInvalid(); err != nil {
		return nil, err
	}
	ret, e := srds.client.SMembers(context.Background(), key).Result()
	if e != nil {
		return nil, fmt.Errorf("smembers key: %v, fail: %v", key, e)
	}
	return ret, nil
}

// SRem 删除某些成员
func (srds *RedisClient) SRem(key string, members ...any) error {
	if err := srds.checkParamInvalid(); err != nil {
		return err
	}

	_, e := srds.client.SRem(context.Background(), key, members...).Result()
	if e != nil {
		return fmt.Errorf("srem key: %v, members: %v, fail: %v", key, members, e)
	}
	return nil
}

// DeleteKey 删除一个hash item
func (srds *RedisClient) DeleteKey(key string) error {
	if err := srds.checkParamInvalid(); err != nil {
		return err
	}

	if key == "" {
		return nil
	}
	err := srds.client.Del(context.Background(), key).Err()
	if err != nil {
		return fmt.Errorf("fail delete key: %v, err: %v", key, err)
	}
	return nil
}
func (srds *RedisClient) ExpireKey(key string, ttl time.Duration) error {
	if err := srds.checkParamInvalid(); err != nil {
		return err
	}
	if key == "" {
		return nil
	}
	if ttl <= 0 {
		return nil
	}
	err := srds.client.Expire(context.Background(), key, ttl).Err()
	if err != nil {
		return fmt.Errorf("expire key: %v, ttl: %v, err: %v", key, ttl, err)
	}
	return nil
}

func (srds *RedisClient) ZADDAtomic(key string, score float64, member any, maxItemNums int32) error {
	if err := srds.checkParamInvalid(); err != nil {
		return err
	}

	// Lua脚本，原子增加元素并删除最早的元素（分数最低)
	luaScript := `
	local zsetKey = KEYS[1]
	local element = ARGV[1]
	local score = tonumber(ARGV[2])
	local maxElements = tonumber(ARGV[3])

	-- 添加元素到有序集合
	redis.call('ZADD', zsetKey, score, element)

	-- 获取有序集合的元素数量
	local currentSize = redis.call('ZCARD', zsetKey)

	-- 如果元素数量超过限制，删除得分最低的元素
	if currentSize > maxElements then
		redis.call('ZREMRANGEBYRANK', zsetKey, 0, 0)
	end

	return currentSize
	`
	// 执行Lua脚本
	_, err := srds.client.Eval(context.Background(), luaScript, []string{key}, member, score, maxItemNums).Result()
	return err
}
func (srds *RedisClient) ZRange(key string, start, stop int64) ([]string, error) {
	if err := srds.checkParamInvalid(); err != nil {
		return nil, err
	}
	members, err := srds.client.ZRange(context.Background(), key, start, stop).Result()
	if err != nil {
		return nil, fmt.Errorf("get member fail, err: %v", err)
	}
	return members, nil
}
func (srds *RedisClient) ZDelMemberByScore(key string, beginScore, endScore float64) (int64, error) {
	if err := srds.checkParamInvalid(); err != nil {
		return 0, err
	}
	minScore := fmt.Sprintf("%f", beginScore)
	maxScore := fmt.Sprintf("%f", endScore)

	removed, err := srds.client.ZRemRangeByScore(context.Background(), key, minScore, maxScore).Result()
	if err != nil {
		return 0, fmt.Errorf("remove by score fail, err: %v, beginScore: %v, endScore: %v, key: %v",
			err, minScore, maxScore, key)
	}
	return removed, nil
}
func (srds *RedisClient) GenUniqueKey(key string, expireTm time.Duration) (string, error) {
	luaScript := `
		if redis.call("EXISTS", KEYS[1]) == 1 then
			redis.call("PEXPIRE", KEYS[1], ARGV[2])
			return redis.call("GET", KEYS[1])
		else 
			local uniqueKey = ARGV[1]
			redis.call("SET", KEYS[1], uniqueKey)
			redis.call("PEXPIRE", KEYS[1], ARGV[2])
			return uniqueKey
		end
		`
	uniqueKey, _ := snowflake.NextID()
	expiryMilliseconds := int64(expireTm / time.Millisecond)

	cmd := srds.client.Eval(context.Background(), luaScript, []string{key} /* keys[1] */, uniqueKey, expiryMilliseconds)
	if cmd.Err() != nil {
		return "", cmd.Err()
	}
	return cmd.Val().(string), nil
}

// GenUniqueKeyWithGenVal 用带有value 产生函数 产生 id
func (srds *RedisClient) GenUniqueKeyWithGenVal(key string, expireTm time.Duration, genFunc func() string) (string, error) {
	luaScript := `
		if redis.call("EXISTS", KEYS[1]) == 1 then
			redis.call("PEXPIRE", KEYS[1], ARGV[2])
			return redis.call("GET", KEYS[1])
		else 
			local uniqueKey = ARGV[1]
			redis.call("SET", KEYS[1], uniqueKey)
			redis.call("PEXPIRE", KEYS[1], ARGV[2])
			return uniqueKey
		end`
	if genFunc == nil {
		return "", fmt.Errorf("genFunc is nil")
	}

	uniqueKeyStr := genFunc()
	if uniqueKeyStr == "" {
		return "", fmt.Errorf("defined gen call return empty")
	}
	expiryMilliseconds := int64(expireTm / time.Millisecond)

	cmd := srds.client.Eval(context.Background(), luaScript, []string{key} /* keys[1] */, uniqueKeyStr, expiryMilliseconds)
	if cmd.Err() != nil {
		return "", cmd.Err()
	}
	return cmd.Val().(string), nil
}

// CheckKeyExist 检查是否存在key, 存在则返回true, 不存在则创建并设置有效时间
func (srds *RedisClient) CheckKeyExist(key string, expireTm time.Duration) (bool, error) {
	luaScript := `
		if redis.call("EXISTS", KEYS[1]) == 1 then
			return 1
		else
			redis.call("set", KEYS[1], ARGV[1])
			 redis.call("PEXPIRE", KEYS[1], tonumber(ARGV[2]))
			return 0
		end`
	expiryMilliseconds := int64(expireTm / time.Millisecond)
	cmd := srds.client.Eval(context.Background(), luaScript, []string{key} /* keys[1] */, 1, expiryMilliseconds)
	if cmd.Err() != nil {
		return false, cmd.Err()
	}
	if cmd.Val().(int64) > 0 {
		return true, nil
	}
	return false, nil
}

// CheckKeyExistFresh 检查是否存在key, 存在则返回true，并刷新时间, 不存在则创建并设置有效时间
func (srds *RedisClient) CheckKeyExistFresh(key string, expireTm time.Duration) (bool, error) {
	luaScript := `
		if redis.call("EXISTS", KEYS[1]) == 1 then
			redis.call("PEXPIRE", KEYS[1], tonumber(ARGV[2]))
			return 1
		else
			redis.call("set", KEYS[1], ARGV[1])
			 redis.call("PEXPIRE", KEYS[1], tonumber(ARGV[2]))
			return 0
		end`
	expiryMilliseconds := int64(expireTm / time.Millisecond)
	cmd := srds.client.Eval(context.Background(), luaScript, []string{key} /* keys[1] */, 1, expiryMilliseconds)
	if cmd.Err() != nil {
		return false, cmd.Err()
	}
	if cmd.Val().(int64) > 0 {
		return true, nil
	}
	return false, nil
}

// CheckKeyExistNoCreate return false: not exist, true: exist
func (srds *RedisClient) CheckKeyExistNoCreate(key string) (bool, error) {
	r, err := srds.client.Exists(context.Background(), key).Result()
	if err != nil {
		return false, err
	}
	if r == 0 {
		return false, nil
	}
	return true, nil
}

// SetEx 增加一个带时间戳的key
func (srds *RedisClient) SetEx(key string, val any, tm time.Duration) error {
	if err := srds.checkParamInvalid(); err != nil {
		return err
	}
	if len(key) == 0 {
		return errors.New("key is nil")
	}
	if tm < 0 {
		tm = 5 * time.Second
	}
	if tm == 0 {
		if err := srds.client.Set(context.Background(), key, val, 0).Err(); err != nil {
			return fmt.Errorf("set fail, key: %v, value value: %+v, err: %v", key, val, err)
		}
		return nil
	}

	if err := srds.client.SetEx(context.Background(), key, val, tm).Err(); err != nil {
		return fmt.Errorf("set ex fail, key: %v, value value: %+v, err: %v", key, val, err)
	}

	return nil
}

// Lock 获取分布式锁
func (srds *RedisClient) Lock(key string, ttl time.Duration, globalId string) error {
	var setLockScript = `
    if redis.call('SET', KEYS[1], ARGV[1], 'NX', 'PX', ARGV[2]) then
        return true
    else
        return false
    end
    `

	result, err := srds.client.Eval(context.Background(), setLockScript, []string{key}, globalId, int(ttl.Milliseconds())).Result()
	if err != nil || result == false {
		return fmt.Errorf("failed to acquire lock")
	}
	return nil
}

// UnLock 释放分布式锁
func (srds *RedisClient) UnLock(key string, globalId string) error {
	var releaseLockScript = `
    if redis.call('GET', KEYS[1]) == ARGV[1] then
        return redis.call('DEL', KEYS[1])
    else
        return 0
    end
    `
	result, err := srds.client.Eval(context.Background(), releaseLockScript, []string{key}, globalId).Result()
	if err != nil {
		return err
	}
	if result.(int64) == 0 {
		return fmt.Errorf("failed to release lock, result: %v", result)
	}
	return nil
}

type KeyNoExistError struct {
}

func (k *KeyNoExistError) Error() string {
	return "key not exist"
}

// RefreshTTL
func (srds *RedisClient) RefreshTTL(key string, globalId string, ttl time.Duration) error {
	var refreshLockScript = `
    if redis.call('GET', KEYS[1]) == ARGV[1] then
        return redis.call('PEXPIRE', KEYS[1], ARGV[2])
    else
        return 0
    end
    `
	result, err := srds.client.Eval(context.Background(), refreshLockScript, []string{key}, globalId, int(ttl.Milliseconds())).Result()
	if err != nil {
		return err
	}
	if result.(int64) == 0 {
		return new(KeyNoExistError)
	}
	return nil
}

// Get 获取 key的值
func (srds *RedisClient) Get(key string) (any, error) {
	if err := srds.checkParamInvalid(); err != nil {
		return nil, err
	}
	if len(key) == 0 {
		return nil, errors.New("key is nil")
	}

	ret, err := srds.client.Get(context.Background(), key).Result()
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func getMapItem(key string, mp map[string]string) string {
	if v, ok := mp[key]; ok {
		return v
	}
	return ""
}

func parseS2I(s string) int64 {
	if s == "" {
		return 0
	}
	v, e := strconv.ParseInt(s, 10, 64)
	if e != nil {
		return 0
	}
	return v
}
func parseS2F32(s string) float32 {
	if s == "" {
		return 0.0
	}
	v, e := strconv.ParseFloat(s, 32)
	if e != nil {
		return 0.0
	}
	return float32(v)
}
func parseS2F64(s string) float64 {
	if s == "" {
		return 0.0
	}
	v, e := strconv.ParseFloat(s, 64)
	if e != nil {
		return 0.0
	}
	return v
}

// TransMapToStruct map数据转化 user.TokenItem, map key is redis tag; eg: *user.TokenItem
func TransMapToStruct[T any](src map[string]string) (*T, error) {
	if len(src) == 0 {
		return nil, fmt.Errorf("param is nil")
	}

	var ret = new(T)
	objValue := reflect.ValueOf(ret).Elem()
	objType := objValue.Type()
	for i := 0; i < objValue.NumField(); i++ {
		field := objValue.Field(i)
		fieldType := objType.Field(i)

		tag := fieldType.Tag.Get("redis")
		if tag != "" && tag != "-" {
			itemValStr := getMapItem(tag, src)
			switch field.Kind() {
			case reflect.String:
				field.SetString(itemValStr)

			case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
				reflect.Int,
				reflect.Uint,
				reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				field.SetInt(parseS2I(itemValStr))

			case reflect.Float32:
				field.SetFloat(float64(parseS2F32(itemValStr)))
			case reflect.Float64:
				field.SetFloat(parseS2F64(itemValStr))

			case reflect.Bool:
				if parseS2I(itemValStr) == 0 {
					field.SetBool(false)
				} else {
					field.SetBool(true)
				}
			}
		}
	}
	return ret, nil
}
