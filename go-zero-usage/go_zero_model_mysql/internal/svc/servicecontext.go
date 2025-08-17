package svc

import (
	"go_zero_model_mysql/internal/config"
	"go_zero_model_mysql/model/mysql/user"

	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config config.Config
	// 增加数据连接对象：
	UserModel user.UserModel
	// model的连接句柄
	ModelCOnn sqlx.SqlConn
	// redis 连接句柄
	RedisConn *redis.Redis
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.DataSource)
	// 连接redis
	redisConn := redis.MustNewRedis(c.RedisConf)
	return &ServiceContext{
		Config: c,
		// 创建连接
		UserModel: user.NewUserModel(conn),
		// 如果要使用底层的查询语句，可以使用 model conn.
		ModelCOnn: conn,
		// redis 连接对象
		RedisConn: redisConn,
	}
}
