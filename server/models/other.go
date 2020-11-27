package models

type Pager struct {
	Page int64 `form:"page" binding:"gte=1"`
	Size int64 `form:"size" binding:"gte=1,lte=50" `
}
