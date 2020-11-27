package controllers

import (
	"clipboard/db"
	"clipboard/models"
	"clipboard/utils"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	. "github.com/gobeam/mongo-go-pagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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

func (c Clip) GetClipOneRoute(ctx *gin.Context) {
	var id string = ctx.Param("id")
	if len(id) == 0 {
		ctx.JSON(422, utils.ErrorMessage("id can't be empty."))
		return
	}

	_id, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		ctx.JSON(422, utils.ErrorMessage("id must be mongo object id."))
		return
	}

	var res bson.M

	var matchStage = bson.D{{"$match", bson.D{{"_id", _id}}}}
	var lookStage = bson.D{{"$lookup", bson.D{
		{"from", "users"},
		{"let", bson.D{{"id", "$user_id"}}},
		{"pipeline", bson.A{
			bson.D{
				{"$project", models.UserProjection},
			},
		}},

		{"as", "user"},
	}}}
	var projectionStage = bson.D{{"$project",
		bson.D{
			{"id", "$_id"},
			{"_id", 0},
			{"content", 1},
			{"created_at", 1},
			{"is_deleted", 1},
			{"type", 1},
			{"user", bson.D{

				{"$arrayElemAt", bson.A{"$user", 0}},
			}},
		}}}

	cur, err := db.ClipCollection.Aggregate(context.TODO(), mongo.Pipeline{matchStage, lookStage, projectionStage})

	if err != nil {
		ctx.JSON(500, utils.ErrorFactory(err))
		return
	}

	for cur.Next(context.Background()) {
		cur.Decode(&res)
	}

	fmt.Println(res)

	ctx.JSON(200, res)

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
