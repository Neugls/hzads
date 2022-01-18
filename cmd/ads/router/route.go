package router

import "github.com/gin-gonic/gin"

func Build() *gin.Engine {
	r := gin.New()
	return r
}
