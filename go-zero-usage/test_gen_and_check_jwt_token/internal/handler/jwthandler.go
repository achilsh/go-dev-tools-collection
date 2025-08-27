package handler

import (
	"net/http"

	"gen_jwt_token_and_check/internal/logic"
	"gen_jwt_token_and_check/internal/svc"
	"gen_jwt_token_and_check/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func JwtHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.JwtTokenRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewJwtLogic(r.Context(), svcCtx)
		resp, err := l.Jwt(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
