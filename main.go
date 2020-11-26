package main

import (
	db "clipboard/db"
	"clipboard/router"

	"github.com/gin-gonic/gin"
)

func main() {
	db.Db()

	r := gin.Default()

	router.RegisterRoutes(r)

	r.Run(":3333")
}
