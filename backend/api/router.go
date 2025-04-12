package api

import (
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	router := gin.Default()
	imageAPI := NewImageAPI()
	router.POST("/api/upload", imageAPI.UploadImage)
	router.GET("/api/directory", imageAPI.GetDirectoryList)
	router.GET("/api/images", imageAPI.GetImageList)
	router.GET("/api/tags", imageAPI.GetTags)
	return router
}
