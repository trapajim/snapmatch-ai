package middleware

import (
	"context"
	"github.com/trapajim/snapmatch-ai/snapmatchai"
	"log"
	"net/http"
)

var sessionKey = struct{}{}

func AuthMiddleware(sessionManager *snapmatchai.Manager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session := sessionManager.SessionStart(w, r)
			r = r.WithContext(context.WithValue(r.Context(), sessionKey, session))
			next.ServeHTTP(w, r)
		})
	}
}

func GetSession(r *http.Request) snapmatchai.Session {
	log.Println(r.Context().Value(sessionKey))
	if r.Context().Value(sessionKey) == nil {
		return nil
	}
	return r.Context().Value(sessionKey).(snapmatchai.Session)
}
