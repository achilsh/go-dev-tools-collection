package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"http_server_etcd/internal/logic"
	"http_server_etcd/internal/svc"
	"http_server_etcd/internal/types"
)

func Http_server_etcdHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.Request
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewHttp_server_etcdLogic(r.Context(), svcCtx)
		resp, err := l.Http_server_etcd(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
