package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"http_demo_server/internal/logic"
	"http_demo_server/internal/svc"
	"http_demo_server/internal/types"
)

func Http_demo_serverHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.Request
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewHttp_demo_serverLogic(r.Context(), svcCtx)
		resp, err := l.Http_demo_server(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
