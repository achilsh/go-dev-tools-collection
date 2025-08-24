package middleware

import (
	"http_gateway_demo/rewrites"
	"net/http"

	"github.com/zeromicro/go-zero/core/logc"
)

func PathRewriteMiddleware(rules *rewrites.RewritHttpUrlPathConf) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			oldPath := r.URL.Path

			// for _, rule := range rules.ReWrites {
			// 	if rule.Method != r.Method {
			// 		continue
			// 	}
			// 	if rule.From == oldPath {
			// 		newPath := rule.To
			// 		u := *r.URL
			// 		u.Path = newPath
			// 		r.URL = &u
			// 		logc.Infof(r.Context(), "Path rewrite: %s => %s", oldPath, newPath)
			// 		break
			// 	}
			// }

			for _, rule := range rules.RewriteRules {
				if rule.Method != r.Method {
					continue
				}
				logc.Infof(r.Context(), "||||| input path: %s, from: %s", oldPath, rule.From.String())
				if rule.From.MatchString(oldPath) {
					newPath := rule.From.ReplaceAllString(oldPath, rule.To)
					u := *r.URL
					u.Path = newPath
					r.URL = &u
					logc.Infof(r.Context(), "Path rewrite: %s => %s", oldPath, newPath)
					break
				}
			}

			next(w, r)
		}
	}
}

// func RewriteUrlInGatewayFM() func(http.HandlerFunc) http.HandlerFunc {

// 	return func(next http.HandlerFunc) http.HandlerFunc {
// 		return func(w http.ResponseWriter, r *http.Request) {
// 			oldPath := r.URL.Path

// 			for _, rule := range rules.RewriteRules {
// 				if rule.Method != r.Method {
// 					continue
// 				}

// 				if rule.From.MatchString(oldPath) {
// 					newPath := rule.From.ReplaceAllString(oldPath, rule.To)
// 					u := *r.URL
// 					u.Path = newPath
// 					r.URL = &u
// 					logc.Infof(r.Context(), "Path rewrite: %s => %s", oldPath, newPath)
// 					break
// 				}
// 			}

// 			next(w, r)
// 		}
// 	}

// }
