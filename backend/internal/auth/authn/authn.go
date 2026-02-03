package authn

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
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
) func() gin.HandlerFunc {

	//return middleware
	return func() gin.HandlerFunc {
		return func(c *gin.Context) {
			tokenStr := extractBearerToken(c.GetHeader("Authorization"))
			if tokenStr == "" {
				c.Next()
				return
			}

			claims, err := jwt_parser.Parse(tokenStr, appSecret)
			if err != nil {

				log.Warn("failed to parse token", sl.Err(err))

				ctx := context.WithValue(c.Request.Context(), errorKey, err)
				c.Request = c.Request.WithContext(ctx)
				c.Next()
				return
			}

			log.Info("user authorized", slog.Any("claims", claims))

			ctx := context.WithValue(c.Request.Context(), uidKey, claims.UserId)
			c.Request = c.Request.WithContext(ctx)
			c.Next()
		}
	}
}

func extractBearerToken(authHeader string) string {
	splitToken := strings.Split(authHeader, "Bearer")

	if len(splitToken) != 2 {
		return ""
	}

	return strings.TrimSpace(splitToken[1])
}

func UIDFromContext(ctx context.Context) (int64, bool) {
	uid, ok := ctx.Value(uidKey).(int64)
	return uid, ok

}

func ErrorFromContext(ctx context.Context) (error, bool) {
	err, ok := ctx.Value(errorKey).(error)
	return err, ok
}

func NewRequireAuth(
	log *slog.Logger,
) func() gin.HandlerFunc {

	return func() gin.HandlerFunc {
		return func(c *gin.Context) {
			if _, ok := UIDFromContext(c.Request.Context()); !ok {
				log.Warn("user is unauthorized")
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": "unauthorized",
				})
				return
			}
			log.Info("user is authorized")
			c.Next()
		}
	}
}
