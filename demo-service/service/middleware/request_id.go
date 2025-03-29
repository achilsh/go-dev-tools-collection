package middleware

import (
	"context"
	"math/rand"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	CtxRequestID = "X-Request-Id"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var DefaultRand = rand.New(rand.NewSource(time.Now().UnixNano()))

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func RandStringRunes(n int) string {
	buf := make([]rune, n)
	for i := range buf {
		buf[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(buf)
}

func UUID() string {
	uid, _ := uuid.NewRandom()
	return strings.ReplaceAll(uid.String(), "-", "")
}

func WithRequestID(ctx context.Context) context.Context {
	if _, ok := ctx.Value(CtxRequestID).(string); ok {
		return ctx
	}
	return context.WithValue(ctx, CtxRequestID, UUID())
}

func GetRequestID(ctx context.Context) string {
	requestID, ok := ctx.Value(CtxRequestID).(string)
	if ok {
		return requestID
	}
	return ""
}

func SetRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, CtxRequestID, requestID)
}

func GenerateRequestID() string {
	return UUID()
}
