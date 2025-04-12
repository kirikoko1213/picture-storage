package api

import (
	"picture_storage/pkg/minio"
	"picture_storage/service"
	"picture_storage/utils"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kiririx/krutils/ut"
)

type ImageAPI struct{}

func NewImageAPI() *ImageAPI {
	return &ImageAPI{}
}

var imageService = service.NewImageService()

type UploadRequest struct {
	Directory string `json:"directory" form:"directory"`
	Tags      string `json:"tags" form:"tags"`
}

func (api *ImageAPI) UploadImage(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		Fail(c, "获取文件失败")
		return
	}

	var req UploadRequest
	if err := c.ShouldBind(&req); err != nil {
		Fail(c, "参数错误")
		return
	}

	// 保存到数据库
	tags := ut.Then(len(req.Tags) > 0, strings.Split(req.Tags, ","), []string{})
	imageID, err := imageService.SaveImage(req.Directory, file, tags)
	if err != nil {
		Fail(c, err.Error())
		return
	}

	Success(c, gin.H{
		"id": imageID,
	})
}

func (api *ImageAPI) GetDirectoryList(c *gin.Context) {
	directoryList, err := imageService.GetDirectoryList()
	if err != nil {
		Fail(c, err.Error())
		return
	}
	Success(c, directoryList)
}

type ImageListRequest struct {
	Directory string   `json:"directory" form:"directory"`
	Page      int      `json:"page" form:"page"`
	PageSize  int      `json:"page_size" form:"page_size"`
	Tags      []string `json:"tags" form:"tags"`
}

type ImageListItem struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	URL       string `json:"url"`
	Size      int64  `json:"size"`
	CreatedAt string `json:"created_at"`
	Directory string `json:"directory"`
}

func (api *ImageAPI) GetImageList(c *gin.Context) {
	var req ImageListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		Fail(c, "参数错误")
		return
	}

	pagination := utils.GetPage(req.Page, req.PageSize)

	images, total, err := imageService.GetImageListByDirectory(req.Directory, pagination)
	if err != nil {
		Fail(c, err.Error())
		return
	}

	Success(c, gin.H{
		"list": func() []map[string]any {
			list := make([]map[string]any, 0)
			for _, image := range images {
				list = append(list, map[string]any{
					"id":        image.ID,
					"imageName": image.ImageName,
					"imageCode": image.ImageCode,
					"url":       minio.Client.GetObjectURL(image.Directory, image.ImageCode+"."+image.Ext),
					"ext":       image.Ext,
					"tags": func() []string {
						tags, err := imageService.GetTagsByImageID(image.ID)
						if err != nil {
							return []string{}
						}
						return tags
					}(),
					"size":      image.Size,
					"directory": image.Directory,
					"createdAt": image.CreatedAt,
				})
			}
			return list
		}(),
		"total": total,
	})
}

func (api *ImageAPI) GetTags(c *gin.Context) {
	tags, err := imageService.GetTags()
	if err != nil {
		Fail(c, err.Error())
		return
	}
	Success(c, tags)
}
