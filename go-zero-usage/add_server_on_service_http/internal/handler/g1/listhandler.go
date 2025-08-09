package g1

import (
	"net/http"

	"add_server_on_service_http/internal/logic/g1"
	"add_server_on_service_http/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func ListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := g1.NewListLogic(r.Context(), svcCtx)
		resp, err := l.List()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
