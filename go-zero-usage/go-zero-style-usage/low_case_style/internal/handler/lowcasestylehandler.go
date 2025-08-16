package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"low_case_style/internal/logic"
	"low_case_style/internal/svc"
	"low_case_style/internal/types"
)

func Low_case_styleHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.Request
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewLow_case_styleLogic(r.Context(), svcCtx)
		resp, err := l.Low_case_style(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
