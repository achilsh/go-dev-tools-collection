package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"http_gateway_demo/upstreams-server/userservice/internal/logic"
	"http_gateway_demo/upstreams-server/userservice/internal/svc"
	"http_gateway_demo/upstreams-server/userservice/internal/types"
)

func UserServicePostHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.PRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewUserServicePostLogic(r.Context(), svcCtx)
		resp, err := l.UserServicePost(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
