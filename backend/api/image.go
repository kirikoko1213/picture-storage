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
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, "参数错误")
		return
	}

	pagination := utils.GetPage(req.Page, req.PageSize)

	images, total, err := imageService.GetImageListByDirectory(req.Directory, req.Tags, pagination)
	if err != nil {
		Fail(c, err.Error())
		return
	}

	imageIDs := make([]uint64, 0)
	for _, image := range images {
		imageIDs = append(imageIDs, image.ID)
	}

	Success(c, gin.H{
		"list": func() []map[string]any {
			tagMap, err := imageService.GetTagsByImageIDs(imageIDs)
			if err != nil {
				return []map[string]any{}
			}
			list := make([]map[string]any, 0)
			for _, image := range images {
				list = append(list, map[string]any{
					"id":           image.ID,
					"imageName":    image.ImageName,
					"imageCode":    image.ImageCode,
					"url":          minio.Client.GetObjectURL(image.Directory, image.ImageCode+"."+image.Ext),
					"thumbnailUrl": minio.Client.GetObjectURL("tmp-thumbnail", image.ThumbnailCode+"."+image.Ext),
					"ext":          image.Ext,
					"tags": func() []string {
						return tagMap[image.ID]
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

func (api *ImageAPI) GetRandomImage(c *gin.Context) {
	tagsParam := c.Query("tags")
	countParam := c.Query("count")
	directoryParam := c.Query("directory")

	tags := ut.Then(tagsParam != "", strings.Split(tagsParam, ","), []string{})
	count := ut.Then(countParam != "", ut.Convert(countParam).Int64Value(), 1)
	directory := ut.Then(directoryParam != "", directoryParam, "")

	image, err := imageService.GetRandomImage(directory, tags, count)
	if err != nil {
		Fail(c, err.Error())
		return
	}
	Success(c, image)
}

func (api *ImageAPI) GetTags(c *gin.Context) {
	tags, err := imageService.GetTags()
	if err != nil {
		Fail(c, err.Error())
		return
	}
	Success(c, tags)
}

func (api *ImageAPI) DeleteImages(c *gin.Context) {
	var req struct {
		IDs []int `json:"ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, "参数错误")
		return
	}

	err := imageService.DeleteImages(req.IDs)
	if err != nil {
		Fail(c, err.Error())
		return
	}

	Success(c, nil)
}

func (api *ImageAPI) AddTags(c *gin.Context) {
	var req struct {
		ImageIDs []uint64 `json:"image_ids"`
		Tags     []string `json:"tags"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, "参数错误")
		return
	}
	if err := imageService.AddTags(req.ImageIDs, req.Tags); err != nil {
		Fail(c, err.Error())
		return
	}
	Success(c, nil)
}

// 获取标签详情（包含图片数量）
func (api *ImageAPI) GetTagDetails(c *gin.Context) {
	tagDetails, err := imageService.GetTagDetails()
	if err != nil {
		Fail(c, err.Error())
		return
	}
	Success(c, gin.H{
		"list": tagDetails,
	})
}

// 创建标签
func (api *ImageAPI) CreateTag(c *gin.Context) {
	var req struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, "参数错误")
		return
	}

	err := imageService.CreateTag(req.Name)
	if err != nil {
		Fail(c, err.Error())
		return
	}

	Success(c, nil)
}

// 更新标签
func (api *ImageAPI) UpdateTag(c *gin.Context) {
	var req struct {
		OldName string `json:"old_name" binding:"required"`
		NewName string `json:"new_name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, "参数错误")
		return
	}

	err := imageService.UpdateTag(req.OldName, req.NewName)
	if err != nil {
		Fail(c, err.Error())
		return
	}

	Success(c, nil)
}

// 删除标签
func (api *ImageAPI) DeleteTag(c *gin.Context) {
	var req struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, "参数错误")
		return
	}

	err := imageService.DeleteTag(req.Name)
	if err != nil {
		Fail(c, err.Error())
		return
	}

	Success(c, nil)
}
