package middleware

import (
	"context"
	"github.com/trapajim/snapmatch-ai/snapmatchai"
	"net/http"
)

var sessionKey = struct{}{}

func AuthMiddleware(sessionManager *snapmatchai.Manager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session := sessionManager.SessionStart(w, r)
			r = r.WithContext(context.WithValue(r.Context(), sessionKey, session))
			expire := session.Get("expire")
			if expire == nil {
				w.Header().Set("X-Session-Duration", "unavailable")
			} else {
				w.Header().Set("X-Session-Duration", expire.(string))
			}

			next.ServeHTTP(w, r)
		})
	}
}

func SetSession(ctx context.Context, session snapmatchai.Session) context.Context {
	return context.WithValue(ctx, sessionKey, session)
}
func GetSession(ctx context.Context) snapmatchai.Session {
	if ctx.Value(sessionKey) == nil {
		return nil
	}
	return ctx.Value(sessionKey).(snapmatchai.Session)
}
