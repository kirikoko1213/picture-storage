package api

import (
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	router := gin.Default()
	imageAPI := NewImageAPI()
	router.POST("/api/upload", imageAPI.UploadImage)
	router.GET("/api/directory", imageAPI.GetDirectoryList)
	router.POST("/api/images", imageAPI.GetImageList)
	router.GET("/api/tags", imageAPI.GetTags)
	router.DELETE("/api/images", imageAPI.DeleteImages)
	return router
}
