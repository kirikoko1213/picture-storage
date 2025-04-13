package service

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"mime/multipart"
	"path/filepath"
	"picture_storage/cache"
	"picture_storage/db"
	"picture_storage/model"
	"picture_storage/pkg/minio"
	"strings"

	"github.com/disintegration/imaging"
	"gorm.io/gorm"
)

type ImageService struct {
}

func NewImageService() *ImageService {
	return &ImageService{}
}

func (service *ImageService) GetDirectoryList() ([]string, error) {
	directoryList, err := minio.Client.GetDirectoryList()
	if err != nil {
		return nil, err
	}
	directoryNameList := make([]string, 0)
	for _, object := range directoryList {
		directoryNameList = append(directoryNameList, object.Name)
	}
	return directoryNameList, nil
}

func (service *ImageService) UploadImage(directory string, file *multipart.FileHeader) (string, int64, error) {
	// 打开文件
	src, err := file.Open()
	if err != nil {
		return "", 0, err
	}
	defer src.Close()

	// 上传到 MinIO
	objectName := file.Filename

	// 上传文件
	md5WithExt, size, err := minio.Client.UploadFile(directory, objectName, file.Size, src, file.Header.Get("Content-Type"))
	if err != nil {
		return "", size, err
	}

	return md5WithExt, size, nil
}

func (service *ImageService) GetTagsByImageID(imageID uint64) ([]string, error) {
	cacheKey := fmt.Sprintf("image_tags_%d", imageID)
	if tags, err := cache.Get(cacheKey); err == nil {
		return tags.([]string), nil
	}
	var tags []string
	var imageTagList []model.ImageTagModel
	err := db.DB.Model(&model.ImageTagModel{}).Where("image_id = ?", imageID).Find(&imageTagList).Error
	if err != nil {
		return nil, err
	}
	for _, imageTag := range imageTagList {
		var tag model.TagModel
		err := db.DB.Model(&model.TagModel{}).Where("id = ?", imageTag.TagID).Find(&tag).Error
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag.TagName)
	}
	cache.Set(cacheKey, tags)
	return tags, nil
}

func (service *ImageService) GetImageListByDirectory(directory string, tags []string, page model.Pagination) ([]model.ImageModel, int64, error) {
	imageList := make([]model.ImageModel, 0)
	var total int64
	if len(tags) > 0 {
		originalQuery := db.DB.Model(&model.ImageModel{}).Joins("JOIN image_tag ON image_tag.image_id = image.id").
			Joins("JOIN tag ON image_tag.tag_id = tag.id").
			Where("tag.tag_name IN ? and image.directory = ?", tags, directory)

		originalQuery.Group("image.id").Count(&total)

		result := originalQuery.
			Offset((page.Page - 1) * page.PageSize).
			Limit(page.PageSize).Distinct("image.*").Find(&imageList)

		err := result.Error
		if err != nil {
			return nil, 0, err
		}

	} else {
		originalQuery := db.DB.Where("directory = ?", directory).
			Order("created_at DESC").
			Find(&imageList)

		originalQuery.Count(&total)

		result := originalQuery.
			Offset((page.Page - 1) * page.PageSize).
			Limit(page.PageSize).Find(&imageList)
		err := result.Error
		if err != nil {
			return nil, 0, err
		}
	}
	return imageList, total, nil
}

func (service *ImageService) GetImageListByTag(directory string, tag string, page model.Pagination) ([]model.ImageDTO, int64, error) {
	// TODO: 这里需要实现从数据库获取图片列表的逻辑
	// 临时返回空数据
	return []model.ImageDTO{}, 0, nil
}

// 根据图片类型创建缩略图
func (service *ImageService) createThumbnail(file *multipart.FileHeader, maxWidth, maxHeight int) ([]byte, error) {
	// 打开上传的文件
	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	// 读取文件内容
	fileBytes, err := io.ReadAll(src)
	if err != nil {
		return nil, err
	}

	// 解码图片
	img, format, err := image.Decode(bytes.NewReader(fileBytes))
	if err != nil {
		return nil, err
	}

	// 使用imaging库调整图片大小，保持宽高比
	resizedImg := imaging.Fit(img, maxWidth, maxHeight, imaging.Lanczos)

	// 编码为相应格式
	var buffer bytes.Buffer
	switch format {
	case "jpeg", "jpg":
		err = jpeg.Encode(&buffer, resizedImg, &jpeg.Options{Quality: 85})
	case "png":
		err = png.Encode(&buffer, resizedImg)
	default:
		// 默认使用JPEG
		err = jpeg.Encode(&buffer, resizedImg, &jpeg.Options{Quality: 85})
	}

	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

// 上传缩略图到MinIO
func (service *ImageService) uploadThumbnail(directory string, originalFilename string, thumbnailData []byte, contentType string) (string, int64, error) {
	// 计算文件内容的MD5
	fileSize := int64(len(thumbnailData))
	md5WithExt, size, err := minio.Client.UploadFileBytes(directory, originalFilename, fileSize, thumbnailData, contentType)
	if err != nil {
		return "", 0, err
	}

	return md5WithExt, size, nil
}

func (service *ImageService) SaveImage(directory string, file *multipart.FileHeader, tags []string) (uint64, error) {
	// 上传原图到 MinIO
	imageCodeWithExt, size, err := service.UploadImage(directory, file)
	if err != nil {
		return 0, err
	}

	// 生成缩略图
	thumbnailData, err := service.createThumbnail(file, 600, 600)
	if err != nil {
		return 0, err
	}

	// 上传缩略图到 MinIO
	thumbnailCodeWithExt, _, err := service.uploadThumbnail("tmp-thumbnail", file.Filename, thumbnailData, file.Header.Get("Content-Type"))
	if err != nil {
		return 0, err
	}

	// 提取文件扩展名
	extension := strings.TrimPrefix(filepath.Ext(imageCodeWithExt), ".")
	// 提取文件名（去掉扩展名）
	imageCode := strings.TrimSuffix(imageCodeWithExt, filepath.Ext(imageCodeWithExt))
	// 开启事务
	tx := db.DB.Begin()
	if tx.Error != nil {
		return 0, tx.Error
	}

	// 先查询是否存在
	var image *model.ImageModel = &model.ImageModel{}
	err = tx.Model(&model.ImageModel{}).Where("image_code = ?", imageCode).Find(image).Error
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	if image.ID != 0 {
		return image.ID, nil
	}

	// 保存图片信息
	image = &model.ImageModel{
		ImageName:     file.Filename,
		ImageCode:     imageCode,
		Directory:     directory,
		Ext:           extension,
		Size:          size,
		ThumbnailCode: strings.TrimSuffix(thumbnailCodeWithExt, filepath.Ext(thumbnailCodeWithExt)),
	}
	if err := tx.Create(image).Error; err != nil {
		tx.Rollback()
		return 0, err
	}

	// 处理标签
	if len(tags) > 0 {
		// 先删除旧的标签关联
		if err := tx.Where("image_id = ?", image.ID).Delete(&model.ImageTagModel{}).Error; err != nil {
			tx.Rollback()
			return 0, err
		}

		// 处理每个标签
		for _, tagName := range tags {
			// 查找或创建标签
			var tag model.TagModel
			if err := tx.Where("tag_name = ?", tagName).First(&tag).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					// 标签不存在，创建新标签
					tag = model.TagModel{
						TagName: tagName,
					}
					if err := tx.Create(&tag).Error; err != nil {
						tx.Rollback()
						return 0, err
					}
				} else {
					tx.Rollback()
					return 0, err
				}
			}

			// 创建标签关联
			imageTag := &model.ImageTagModel{
				ImageID: image.ID,
				TagID:   tag.ID,
			}
			if err := tx.Create(imageTag).Error; err != nil {
				tx.Rollback()
				return 0, err
			}
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return 0, err
	}

	return image.ID, nil
}

func (service *ImageService) GetTags() ([]string, error) {
	var tags []model.TagModel
	err := db.DB.Model(&model.TagModel{}).Order("created_at ASC").Find(&tags).Error
	if err != nil {
		return nil, err
	}
	tagList := make([]string, 0)
	for _, tag := range tags {
		tagList = append(tagList, tag.TagName)
	}
	return tagList, nil
}
