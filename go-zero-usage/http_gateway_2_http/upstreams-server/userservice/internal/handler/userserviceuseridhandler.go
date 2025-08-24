package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"http_gateway_demo/upstreams-server/userservice/internal/logic"
	"http_gateway_demo/upstreams-server/userservice/internal/svc"
	"http_gateway_demo/upstreams-server/userservice/internal/types"
)

func UserserviceUserIdHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.Request
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewUserserviceUserIdLogic(r.Context(), svcCtx)
		resp, err := l.UserserviceUserId(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
