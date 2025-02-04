package handler

import (
	"context"
	"github.com/trapajim/snapmatch-ai/server/middleware"
)

func GetSessionExpiry(ctx context.Context) string {
	sess := middleware.GetSession(ctx)
	if sess == nil {
		return "No session found"
	}
	if val := sess.Get("expire"); val != nil {
		return sess.Get("expire").(string)

	}
	return "Session expiration not set"
}
