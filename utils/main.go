package utils

import "github.com/gin-gonic/gin"

func ErrorFactory(err error) gin.H {
	return gin.H{"message": err.Error()}
}
