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
	router.GET("/api/images/random", imageAPI.GetRandomImage)
	router.GET("/api/tags", imageAPI.GetTags)
	router.GET("/api/tags/details", imageAPI.GetTagDetails)
	router.POST("/api/tags", imageAPI.CreateTag)
	router.PUT("/api/tags", imageAPI.UpdateTag)
	router.DELETE("/api/tags", imageAPI.DeleteTag)
	router.POST("/api/images/tags", imageAPI.AddTags)
	router.DELETE("/api/images", imageAPI.DeleteImages)
	return router
}
