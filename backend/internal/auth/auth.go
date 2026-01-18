package auth

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/rwrrioe/pythia/backend/internal/lib/jwt_parser"
	"github.com/rwrrioe/pythia/backend/internal/lib/logger/sl"
)

const (
	errorKey = "middleware auth error"
	uidKey   = "user_id"
)

var (
	ErrInvalidToken = errors.New("invalid token")
)

func New(
	log *slog.Logger,
	appSecret string,
) func(next http.Handler) http.Handler {

	//return middleware
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenStr := extractBearerToken(r)
			if tokenStr == "" {
				next.ServeHTTP(w, r)
				return
			}

			claims, err := jwt_parser.Parse(tokenStr, appSecret)
			if err != nil {

				log.Warn("failed to parse token", sl.Err(err))

				ctx := context.WithValue(r.Context(), errorKey, ErrInvalidToken)
				next.ServeHTTP(w, r.WithContext(ctx))

				return
			}

			log.Info("user authorized", slog.Any("claims", claims))

			ctx := context.WithValue(r.Context(), uidKey, claims.UserId)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func extractBearerToken(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	splitToken := strings.Split(authHeader, "Bearer")

	if len(splitToken) != 2 {
		return ""
	}

	return splitToken[1]
}

func UIDFromContext(ctx context.Context) (int64, bool) {
	uid, ok := ctx.Value(uidKey).(int64)
	return uid, ok
}

func ErrorFromContext(ctx context.Context) (error, bool) {
	err, ok := ctx.Value(errorKey).(error)
	return err, ok
}
