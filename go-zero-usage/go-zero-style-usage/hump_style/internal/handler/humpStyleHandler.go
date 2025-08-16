package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"hump_style/internal/logic"
	"hump_style/internal/svc"
	"hump_style/internal/types"
)

func Hump_styleHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.Request
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewHump_styleLogic(r.Context(), svcCtx)
		resp, err := l.Hump_style(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
