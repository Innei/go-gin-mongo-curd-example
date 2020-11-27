package controllers

import (
	"clipboard/db"
	"clipboard/models"
	"clipboard/utils"
	"context"
	"github.com/gin-gonic/gin"
	. "github.com/gobeam/mongo-go-pagination"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"time"
)

type Clip struct {
}

func (c Clip) GetClipsRoute(ctx *gin.Context) {

	var query models.Pager
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

func (c Clip) CreateClipRoute(ctx *gin.Context) {
	var body models.Clip

	if err := ctx.ShouldBindJSON(&body); err != nil {

		ctx.JSON(http.StatusUnprocessableEntity, utils.ErrorFactory(err))

		return
	}
	user, exist := ctx.Get("user")

	if !exist {
		ctx.JSON(422, gin.H{"message": "user is not exist."})
		return
	}

	model := models.Clip{
		Content:   body.Content,
		CreatedAt: time.Now(),
		IsDeleted: false,
		Type:      body.Type,
		UserId:    user.(models.User).Id,
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
