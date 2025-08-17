## 使用 goctl model

* 支持 mysql, postgresql, mongo 代码生成，mysql支持 从sql文件和数据库连接两种生成，postgreSQL 仅支持从数据库连接生成。

* 对mysql 使用 sql 来生成源码过程：

1. 创建类似 sql 语句，比如： sql/user.sql
2. 使用命令行创建 ```goctl model mysql ddl -s ./sql/us*.sql  -d ./model/mysql/user/ --style=go_zero --strict```
其中 -s 指定 sql 源文件名匹配； -d 指定产生go文件路径；

默认用户在建表时会创建createTime、updateTime字段(忽略大小写、下划线命名风格)且默认值均为CURRENT_TIMESTAMP，而updateTime支持ON UPDATE CURRENT_TIMESTAMP；
对于这两个字段生成insert、update时会被移除，不在赋值范畴内。
默认忽略这些字段的操作，比如源码中：

```
 userRowsExpectAutoSet   = strings.Join(stringx.Remove(userFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), ",")
 userRowsWithPlaceHolder = strings.Join(stringx.Remove(userFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), "=?,") + "=?"

```

如何不想忽略任何字段则设置 为空字符串：

```
goctl model mysql ddl -s ./sql/us*.sql -d ./model/mysql/user/ --style=go_zero -i "" --strict
```

生成源码：

```
 userRowsExpectAutoSet   = strings.Join(stringx.Remove(userFieldNames, "`id`"), ",")
 userRowsWithPlaceHolder = strings.Join(stringx.Remove(userFieldNames, "`id`"), "=?,") + "=?"

 ```

* 如何使用生成的 sql 代码， 使用不带cache mysql

1) 在配置文件增加 mysql 配置，在配置文件源码增加对应项， 在internal/config/config.go新增：

```
type Config struct {
 zrpc.RpcServerConf
 //  增加数据库连接地址配置：
 DataSource string
}
```

在配置文件etc/ 文件中新增：

```

DataSource: root:@tcp(localhost:3306)/demo_test

```

2) 在上下文 上增加对 db连接初始化， 修改文件：svc/servercontext.go

```

type ServiceContext struct {
 Config config.Config
 // 增加数据连接对象：
 UserModel user.UserModel
}

```

3) 在业务逻辑中做数据库操作，eg: 修改 gozeromodelmysql.go 源文件：

```
//  增加测试
go func() {
  time.Sleep(2 * time.Second)

  logic.AddInsert(ctx, logx.WithContext(context.Background()))
 }()
```

4) 如何定制化查询
首先要拿到model 连接句柄，参考文件：svr/servicecontext.go

```

type ServiceContext struct {
 Config config.Config
 // 增加数据连接对象：
 UserModel user.UserModel
 // model的连接句柄
 ModelCOnn sqlx.SqlConn
}

func NewServiceContext(c config.Config) *ServiceContext {
 conn := sqlx.NewMysql(c.DataSource)
 return &ServiceContext{
  Config: c,
  // 创建连接
  UserModel: user.NewUserModel(conn),
  // 如果要使用底层的查询语句，可以使用 model conn.
  ModelCOnn: conn,
 }
}
```

可以使用 ServiceContext中的 ModelConn 对象，接口是：单条记录 QueryRowPartialCtx/QueryRowPartialCtx

```
type User struct {
    Id   int64  `db:"id"`
    Name string `db:"name"`
    Age int `db:"age"`
}

var conn sqlx.SqlConn
var u User
err := conn.QueryRowPartialCtx(context.Background(), &u, "select id, name from user where id = ? limit 1", 1)
if err != nil { // err == nil
    fmt.Println(err) 
    return
}
_ = u // age is default 0
```

查询多条记录：使用接口 QueryRowsPartial/QueryRowsPartialCtx

```
type User struct {
    Id   int64  `db:"id"`
    Name string `db:"name"`
}

var conn sqlx.SqlConn
var users []*User
err := conn.QueryRowsPartial(context.Background(), &users, "select id, name from user where name = ?", "dylan")
if err != nil {
    fmt.Println(err)
    return
}
_ = users
```

5） 自定义的增，删，改操作，使用接口： ExecCtx()

* 自定义事务：

```
 var conn sqlx.SqlConn
    err := conn.TransactCtx(context.Background(), func(ctx context.Context, session sqlx.Session) error {
        r, err := session.ExecCtx(ctx, "insert into user (id, name) values (?, ?)", 1, "test")
        if err != nil {
            return err
        }
        r ,err =session.ExecCtx(ctx, "insert into user (id, name) values (?, ?)", 2, "test")
        if err != nil {
            return err
        }
    })

```

### go-0中 redis 使用

* 增加 redis 配置， 在etc 配置文件中新增：

```

RedisConf:
  # redis 服务地址 ip:port, 如果是 redis cluster 则为 ip1:port1,ip2:port2,ip3:port3...(暂不支持redis sentinel)
  Host: 0.0.0.0:6379
  # node 单节点 redis，cluster redis 集群
  Type: node 
  Pass: pw_test
  Tls: false
```

修改go 源码中的配置项， 修改文件 internal/config/config.go

```

type Config struct {
 zrpc.RpcServerConf
 //  增加数据库连接地址配置：
 DataSource string
 // 增加 redis 配置
 RedisConf redis.RedisConf
}

```

创建全局redis 句柄， 修改文件 svc/servicecontext.go

```


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


```

增加测试用例， 修改文件 gozeromodelmysql.go 源码：

```
//增加 redis op test
 go func() {
  time.Sleep(2 * time.Second)
  e := ctx.RedisConn.SetCtx(context.Background(), "key_1", "value_1")
  if e != nil {
   fmt.Println("insert redis key fail, err: ", e)
  }
 }()

```
