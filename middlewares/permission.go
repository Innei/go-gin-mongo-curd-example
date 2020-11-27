package middlewares

import (
	"clipboard/controllers"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func PermissionMiddleware() gin.HandlerFunc {

	return func(ctx *gin.Context) {
		token := ctx.Request.Header.Get("authorization")
		authHeader := strings.Split(token, "bearer ")
		if len(authHeader) != 2 {
			
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": "invalid token."})
			ctx.Abort()
		} else {
			jwtToken := authHeader[1]

			token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
				}
				return []byte(controllers.SERCRT_KEY), nil
			})

			if err != nil {
				log.Fatal(err)
			}

			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				ctx.Set("user", claims)
				ctx.Next()
			} else {

				ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
				ctx.Abort()

			}

		}
	}

}
