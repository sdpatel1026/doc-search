package main

import (
	"github.com/gin-gonic/gin"
	"github.com/sdpatel1026/doc-search/controllers"
)

var router *gin.Engine

func initializeRoutes() {
	router = gin.New()
	versionV1 := router.Group("/v1")
	versionV1.POST("/train", controllers.Train)
	versionV1.GET("/search/:text", controllers.DocSearch)
}
