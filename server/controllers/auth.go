package controllers

import (
	"clipboard/db"
	"clipboard/models"
	"clipboard/utils"
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type Auth struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

var SERCRT_KEY = []byte("sddsdsdasdasddsdsdasdasddsdsdasda")

type LoginPayload struct {
	Username string             `json:"username"`
	Token    string             `json:"token"`
	Email    string             `json:"email"`
	Id       primitive.ObjectID `json:"id"`
}

func (a Auth) LoginRoute(ctx *gin.Context) {
	var body Auth
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(422, utils.ErrorFactory(err))

		return
	}

	username := body.Username
	var document models.User

	err := db.UserCollection.FindOne(context.TODO(), bson.M{"username": username}).Decode(&document)
	if err != nil {
		ctx.JSON(422, utils.ErrorFactory(err))
		return
	}

	password := body.Password
	isValid := comparePasswords(document.Password, []byte(password))
	if !isValid {
		ctx.JSON(422, gin.H{"message": "password is not correct."})
		return
	}
	payload := LoginPayload{
		Email:    document.Email,
		Id:       document.Id,
		Token:    SignToken(document.Id.Hex()),
		Username: document.Username,
	}

	ctx.JSON(200, payload)
}

func (a Auth) RegisterRoute(ctx *gin.Context) {
	var body models.User

	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(422, utils.ErrorFactory(err))
		return
	}
	body.CreatedAt = time.Now()
	body.Password = hashAndSalt([]byte(body.Password))
	res, err := db.UserCollection.InsertOne(context.TODO(), body)
	if err != nil {
		ctx.JSON(500, utils.ErrorFactory(err))
		return
	}
	token := SignToken(res.InsertedID.(primitive.ObjectID).Hex())
	ctx.JSON(200, gin.H{"data": res.InsertedID, "token": token})
}

func SignToken(id string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":        id,
		"timestamp": time.Now(),
	})
	tokenString, err := token.SignedString(SERCRT_KEY)
	if err != nil {
		log.Fatal(err)

	}
	return tokenString
}

func hashAndSalt(pwd []byte) string {

	// Use GenerateFromPassword to hash & salt pwd.
	// MinCost is just an integer constant provided by the bcrypt
	// package along with DefaultCost & MaxCost.
	// The cost can be any value you want provided it isn't lower
	// than the MinCost (4)
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	// GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it
	return string(hash)
}

func comparePasswords(hashedPwd string, plainPwd []byte) bool {
	// Since we'll be getting the hashed password from the DB it
	// will be a string so we'll need to convert it to a byte slice
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}
