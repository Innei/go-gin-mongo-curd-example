package controllers

import (
	"clipboard/db"
	"clipboard/models"
	"clipboard/utils"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	. "github.com/gobeam/mongo-go-pagination"
	"go.mongodb.org/mongo-driver/bson"
)

type Pager struct {
	Page int64 `form:"page"`
	Size int64 `form:"size"`
}

func GetClipRoute(ctx *gin.Context) {

	var query Pager
	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.JSON(200, utils.ErrorFactory(err))
		return
	}
	page := query.Page

	limit := query.Size

	res, err := New(db.ClipCollection).Limit(limit).Page(page).Sort("created_at", -1).Filter(bson.M{}).Find()
	if err != nil {
		ctx.JSON(500, utils.ErrorFactory(err))
		return
	}

	var List []models.Clip

	for _, raw := range res.Data {
		var clip *models.Clip
		if marshallErr := bson.Unmarshal(raw, &clip); marshallErr == nil {
			List = append(List, *clip)
		}

	}

	ctx.JSON(200, gin.H{"data": List, "paginator": res.Pagination})
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
