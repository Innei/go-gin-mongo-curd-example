package controllers

import (
	"clipboard/db"
	"clipboard/models"
	"clipboard/utils"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func GetClipRoute(ctx *gin.Context) {
	ctx.JSON(200, gin.H{})
}

func CreateClipRoute(ctx *gin.Context) {
	var body models.Clip

	if err := ctx.ShouldBindJSON(&body); err != nil {

		ctx.JSON(http.StatusUnprocessableEntity, utils.ErrorFactory(err))

		return
	}
	model := models.Clip{
		Content:   body.Content,
		CreatedAt: time.Now(),
		IsDeleted: false,
		Type:      body.Type,
	}
	if res, err := db.ClipCollection.InsertOne(context.TODO(), model); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorFactory(err))
		return
	} else {

		id := res.InsertedID
		var result models.Clip
		err := db.ClipCollection.FindOne(context.TODO(), bson.D{{"_id", id}}).Decode(&result)

		if err != nil {
			ctx.JSON(400, gin.H{"message": "query error."})
			return
		}

		ctx.JSON(200, result)
	}

}

type MyTime struct {
	time.Time
}
