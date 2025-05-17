package log_context

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

const (
	CtxRequestID string = "X-Request-Id"
	CtxUserID    string = "ctx_user_id"
)

func UUID() string {
	uid, _ := uuid.NewRandom()
	return strings.ReplaceAll(uid.String(), "-", "")
}

func GetUserID(ctx context.Context) string {
	userIdStr, ok := ctx.Value(CtxUserID).(string)
	if !ok {
		return ""
	}
	return userIdStr

}
func WithCtxUserID(ctx context.Context, userId string) context.Context {
	if len(userId) <= 0 {
		return ctx
	}
	return context.WithValue(ctx, any(CtxUserID), any(userId))
}

func GetRequestID(ctx context.Context) string {
	reqIdStr, ok := ctx.Value(CtxRequestID).(string)
	if !ok {
		return ""
	}
	return reqIdStr
}
func WithRequestID(ctx context.Context, id ...string) context.Context {
	var idStr = ""
	if len(id) <= 0 || len(id[0]) <= 0 {
		idStr = UUID()
	} else {
		idStr = id[0]
	}

	return context.WithValue(ctx, CtxRequestID, idStr)
}
func GetCtxInfo(ctx context.Context) string {
	id := GetRequestID(ctx)
	userId := GetUserID(ctx)
	if id == "" && userId == "" {
		return ""
	}
	return fmt.Sprintf("%v|%v|", id, userId)
}
