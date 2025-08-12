package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"http_gateway_demo/upstreams-server/userservice/internal/logic"
	"http_gateway_demo/upstreams-server/userservice/internal/svc"
)

func UserserviceUserHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewUserserviceUserLogic(r.Context(), svcCtx)
		resp, err := l.UserserviceUser()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
