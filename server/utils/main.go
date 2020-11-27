package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func ErrorFactory(err error) gin.H {
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		return ErrorMessage(err.Error())
	} else {

		return ErrorMessage(errs.Error())
	}
}

func ErrorMessage(msg string) gin.H {
	return gin.H{"message": msg}

}

func Contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
