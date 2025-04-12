package main

import (
	"picture_storage/api"
	"picture_storage/pkg/minio"

	_ "github.com/kiririx/easy-config"
	_ "github.com/sirupsen/logrus"
	_ "github.com/tidwall/gjson"
	_ "gorm.io/driver/mysql"
	_ "gorm.io/gorm"
)

func main() {
	// 初始化 MinIO 客户端
	minio.InitMinioClient()
	router := api.InitRouter()
	router.Run(":10048")
}
