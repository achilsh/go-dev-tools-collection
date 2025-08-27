## 演示使用Jwt 来生成token 和验证token

* 编写 api 文件，提供两个接口，分别是模拟用户登录，和使用 jwt 验证 token
* 使用命令 生成http 服务：  goctl api go --api ./jwt.api --dir ./
* 在配置文件：修改jwt token 时间和密钥

```


JwtAuth:
  AccessSecret: xxxxxxxxxxxxxxxxxxxxxxxxxxxxx
  AccessExpire: 60
```

* 修改 生成token 的jwt 流程， 修改文件是：internal/logic/jwtlogic.go
可以在 类如代码中增加自定义字段：

```
accessToken, err := l.GenToken(now,
  l.svcCtx.Config.JwtAuth.AccessSecret,
  map[string]interface{}{
   "user_id": req.UserId},
  accessExpire)


  
func (l *JwtLogic) GenToken(
 iat int64,
 secretKey string,
 payloads map[string]interface{},
 seconds int64,
) (string, error) {
 claims := make(jwt.MapClaims)
 claims["exp"] = iat + seconds
 claims["iat"] = iat
 for k, v := range payloads {
  claims[k] = v
 }

 // 设置签名方法，保持和前后端一致即可
 token := jwt.New(jwt.SigningMethodHS256)
 token.Claims = claims

 return token.SignedString([]byte(secretKey))
}
```

* 在客户端调用，产生token, 比如 使用post:
url: <http://172.31.60.55:8888/user/token>
body:
{
    "user_id":"achilsh",
    "passwd": "abasdfadf"
}

* 产生 access_token:
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTYzMTIyMjAsImlhdCI6MTc1NjMxMjE2MCwidXNlcl9pZCI6ImFjaGlsc2gifQ.sUlfJ5PHyWuflehULhHmWqL9t3OFI8fn6_Dt03Bm7Ks

* 然后使用该token 访问其他接口：
比如访问其他接口， post:
<http://172.31.60.55:8888/user/info>
http header 字段 Authorization
该字段值为上面分配的 access_token值。

* 触发调用

* 如果想在业务上从token 中解析token生成时设置值，比如刚才设置的值：

```
accessToken, err := l.GenToken(now,
  l.svcCtx.Config.JwtAuth.AccessSecret,
  map[string]interface{}{
   "user_id": req.UserId},
  accessExpire)

  
```

比如上面的 user_id值，
可以从上下文中获取，因为框架已经把token中解析出的自定义字段和值放到 上下文中，
比如在 internal/logic/getuserlogic.go 文件中，例如：

```
func (l *GetUserLogic) GetUser(req *types.GetUserRequest) (resp *types.GetUserResponse, err error) {
 // todo: add your logic here and delete this line
 //获取 登录 token 内内置一些值
 tokenInfo, ok := l.ctx.Value("user_id").(string)
 if ok {
  l.Logger.Debugf("token info: %v", tokenInfo)
 }
 return
}
```
