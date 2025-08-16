package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"snake_case_style/internal/logic"
	"snake_case_style/internal/svc"
	"snake_case_style/internal/types"
)

func Snake_case_styleHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.Request
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewSnake_case_styleLogic(r.Context(), svcCtx)
		resp, err := l.Snake_case_style(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
