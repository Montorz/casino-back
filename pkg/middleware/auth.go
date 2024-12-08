package middleware

import (
	"casino-back/internal/app/handler/dto"
	"casino-back/pkg/token"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const AuthHeaderName = "Authorization"
const contextUserId = "userId"

type TokenDecoder interface {
	Valid(tokenString string) (bool, error)
	Decode(tokenString string, toClaims token.Claims) (*jwt.Token, *token.Claims, error)
}

type ErrorResponseMessage struct {
	Status      int         `json:"status"`
	Code        string      `json:"string"`
	Description interface{} `json:"description"`
}

func JwtAuth(tokenDecoder TokenDecoder) gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		authHeader := ginCtx.GetHeader(AuthHeaderName)
		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 {
			ginCtx.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponseMessage{
				Status:      http.StatusUnauthorized,
				Code:        "",
				Description: "auth header is wrong",
			})
			return
		}

		if headerParts[0] != "Bearer" {
			ginCtx.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponseMessage{
				Status:      http.StatusUnauthorized,
				Code:        "",
				Description: "auth header is wrong",
			})
			return
		}

		valid, err := tokenDecoder.Valid(headerParts[1])
		if err != nil {
			ginCtx.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponseMessage{
				Status:      http.StatusUnauthorized,
				Code:        "",
				Description: err.Error(),
			})
			return
		}

		if !valid {
			ginCtx.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponseMessage{
				Status:      http.StatusUnauthorized,
				Code:        "",
				Description: "token is not valid",
			})
			return
		}

		decodeToken, _, _ := tokenDecoder.Decode(
			headerParts[1],
			&dto.JwtUserClaims{},
		)
		if claims, ok := decodeToken.Claims.(*dto.JwtUserClaims); ok {
			ginCtx.Set(contextUserId, claims.UserId)
		}

		ginCtx.Next()
	}
}

func JwtWebSocketAuth(tokenDecoder TokenDecoder) gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		query := ginCtx.DefaultQuery("token", "")
		if query == "" {
			ginCtx.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponseMessage{
				Status:      http.StatusUnauthorized,
				Code:        "",
				Description: "auth message is wrong",
			})
			return
		}

		valid, err := tokenDecoder.Valid(query)
		if err != nil || !valid {
			ginCtx.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponseMessage{
				Status:      http.StatusUnauthorized,
				Code:        "",
				Description: err.Error(),
			})
			return
		}

		if !valid {
			ginCtx.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponseMessage{
				Status:      http.StatusUnauthorized,
				Code:        "",
				Description: "token is not valid",
			})
			return
		}

		decodeToken, _, _ := tokenDecoder.Decode(
			query,
			&dto.JwtUserClaims{},
		)
		if claims, ok := decodeToken.Claims.(*dto.JwtUserClaims); ok {
			ginCtx.Set(contextUserId, claims.UserId)
		}

		ginCtx.Next()
	}
}
