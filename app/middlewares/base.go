package middlewares

import (
	"net/http"
	"os"

	"chat-backend/app/lib"

	"github.com/gorilla/context"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
)

// BaseMiddleware struct
type BaseMiddleware struct {
	DB *gorm.DB
}

// SetContentTypeHeader to JSON
func SetContentTypeHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		response.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(response, request)
	})
}

// AuthJwtVerify verify token and add UserID to the request context
func AuthJwtVerify(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		var header = request.Header.Get("Authorization")

		if header == "" {
			lib.Error(response, http.StatusForbidden, "Missing Authorization Token")
			return
		}

		token, err := jwt.Parse(header, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("SECRET")), nil
		})

		if err != nil {
			lib.Error(response, http.StatusForbidden, "Invalid token, please login")
			return
		}

		claims, _ := token.Claims.(jwt.MapClaims)
		context.Set(request, "UserID", claims["UserID"])
		context.Set(request, "RoleName", claims["RoleName"])
		next.ServeHTTP(response, request)
	})
}

// OnlyHighAdmin can access
func OnlyHighAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		roleName := context.Get(request, "RoleName")

		if roleName != "High Admin" {
			lib.Error(response, http.StatusUnauthorized, "You can't access this page")
			return
		}

		next.ServeHTTP(response, request)
	})
}
