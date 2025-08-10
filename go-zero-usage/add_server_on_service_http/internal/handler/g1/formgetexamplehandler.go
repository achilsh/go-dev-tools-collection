package g1

import (
	"net/http"

	"add_server_on_service_http/internal/logic/g1"
	"add_server_on_service_http/internal/svc"
	"add_server_on_service_http/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func FormGetExampleHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.FormExampleReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := g1.NewFormGetExampleLogic(r.Context(), svcCtx)
		err := l.FormGetExample(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.Ok(w)
		}
	}
}
