package main

import (
	db "clipboard/db"
	"clipboard/router"
	"github.com/gin-gonic/gin"
)

func main() {
	db.InitDb()

	//if err := utils.InitTrans("zh"); err != nil {
	//	fmt.Printf("init trans failed, err:%v\n", err)
	//	return
	//}

	r := gin.Default()

	router.RegisterRoutes(r)

	r.Run(":3333")
}
