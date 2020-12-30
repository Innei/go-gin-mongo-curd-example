package middlewares

import (
	"clipboard/controllers"
	"clipboard/db"
	"clipboard/models"
	"clipboard/utils"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

				var id string
				for k, v := range claims {
					if k == "id" {
						id = v.(string)

					}

				}

				if len(id) == 0 {
					ctx.JSON(422, gin.H{"message": "id is empty."})
					ctx.Abort()
					return
				}

				var user models.User

				i, err := primitive.ObjectIDFromHex(id)
				if err != nil {
					log.Fatal(err)
				}
				err = db.UserCollection.FindOne(context.TODO(), bson.M{"_id": i}).Decode(&user)
				if err != nil {
					ctx.Abort()
					ctx.JSON(500, utils.ErrorFactory(err))
					return
				}
				ctx.Set("user", user)
				ctx.Next()
			} else {

				ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
				ctx.Abort()

			}

		}
	}

}
