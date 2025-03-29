package user

type TokenItem struct {
	Uid      string `json:"uid" redis:"uid"`
	Akey     string `json:"akey" redis:"akey"`
	ExpireAt int    `json:"expireAt" redis:"expireAt"` // 失效时间点（单位为秒，绝对时间）
}
