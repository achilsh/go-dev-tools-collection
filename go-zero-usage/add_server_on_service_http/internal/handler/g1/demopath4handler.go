package g1

import (
	"net/http"

	"add_server_on_service_http/internal/logic/g1"
	"add_server_on_service_http/internal/svc"
	"add_server_on_service_http/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func DemoPath4Handler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DemoPath4Req
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := g1.NewDemoPath4Logic(r.Context(), svcCtx)
		resp, err := l.DemoPath4(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
