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

func (service *ImageService) GetTagsByImageIDs(imageIDs []uint64) (map[uint64][]string, error) {
	if len(imageIDs) == 0 {
		return make(map[uint64][]string), nil
	}

	// 初始化结果映射
	result := make(map[uint64][]string)
	for _, id := range imageIDs {
		result[id] = make([]string, 0)
	}

	// 查询所有相关的图片标签关联
	var imageTagList []model.ImageTagModel
	err := db.DB.Model(&model.ImageTagModel{}).
		Where("image_id IN ?", imageIDs).
		Find(&imageTagList).Error
	if err != nil {
		return nil, err
	}

	if len(imageTagList) == 0 {
		return result, nil
	}

	// 获取所有标签ID
	tagIDs := make([]uint64, 0)
	for _, imageTag := range imageTagList {
		tagIDs = append(tagIDs, imageTag.TagID)
	}

	// 查询所有标签信息
	var tags []model.TagModel
	err = db.DB.Model(&model.TagModel{}).
		Where("id IN ?", tagIDs).
		Find(&tags).Error
	if err != nil {
		return nil, err
	}

	// 构建标签ID到标签名的映射
	tagMap := make(map[uint64]string)
	for _, tag := range tags {
		tagMap[tag.ID] = tag.TagName
	}

	// 构建最终结果
	for _, imageTag := range imageTagList {
		if tagName, exists := tagMap[imageTag.TagID]; exists {
			result[imageTag.ImageID] = append(result[imageTag.ImageID], tagName)
		}
	}

	return result, nil
}

func (service *ImageService) GetTagsByImageID(imageID uint64) ([]string, error) {
	//cacheKey := fmt.Sprintf("image_tags_%d", imageID)
	//if tags, err := cache.Get(cacheKey); err == nil {
	//	return tags.([]string), nil
	//}
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
	//cache.Set(cacheKey, tags)
	return tags, nil
}

func (service *ImageService) GetImageListByDirectory(directory string, tags []string, page model.Pagination) ([]model.ImageModel, int64, error) {
	imageList := make([]model.ImageModel, 0)
	var total int64
	if len(tags) > 0 {
		// 查询同时拥有所有指定标签的图片
		// 使用 GROUP BY 和 HAVING 来确保图片拥有所有指定的标签
		baseQuery := db.DB.Model(&model.ImageModel{}).
			Joins("JOIN image_tag ON image_tag.image_id = image.id").
			Joins("JOIN tag ON image_tag.tag_id = tag.id").
			Where("tag.tag_name IN ? AND image.directory = ?", tags, directory).
			Group("image.id").
			Having("COUNT(DISTINCT tag.id) = ?", len(tags))

		// 统计总数
		var countResult []struct {
			ID uint64
		}
		err := baseQuery.Select("image.id").Find(&countResult).Error
		if err != nil {
			return nil, 0, err
		}
		total = int64(len(countResult))

		// 查询具体数据
		result := baseQuery.
			Order("image.created_at DESC").
			Offset((page.Page - 1) * page.PageSize).
			Limit(page.PageSize).
			Select("image.*").
			Find(&imageList)

		err = result.Error
		if err != nil {
			return nil, 0, err
		}

	} else {
		// 没有标签筛选时，查询该目录下的所有图片
		baseQuery := db.DB.Where("directory = ?", directory).
			Order("created_at DESC")

		// 统计总数
		baseQuery.Model(&model.ImageModel{}).Count(&total)

		// 查询具体数据
		result := baseQuery.
			Offset((page.Page - 1) * page.PageSize).
			Limit(page.PageSize).
			Find(&imageList)
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

func (service *ImageService) DeleteImages(ids []int) error {
	tx := db.DB.Begin()

	for _, id := range ids {
		var image model.ImageModel
		if err := tx.Where("id = ?", id).First(&image).Error; err != nil {
			tx.Rollback()
			return err
		}
		if err := tx.Where("image_id = ?", id).Delete(&model.ImageTagModel{}).Error; err != nil {
			tx.Rollback()
			return err
		}
		if err := tx.Where("id = ?", id).Delete(&model.ImageModel{}).Error; err != nil {
			tx.Rollback()
			return err
		}
		// 同时删除 minio 的图片
		imageCodeWithExt := image.ImageCode + "." + image.Ext
		thumbnailCodeWithExt := image.ThumbnailCode + ".jpg"
		if err := minio.Client.DeleteFile(image.Directory, imageCodeWithExt); err != nil {
			tx.Rollback()
			return err
		}
		if err := minio.Client.DeleteFile("tmp-thumbnail", thumbnailCodeWithExt); err != nil {
			tx.Rollback()
			return err
		}
	}
	if err := tx.Commit().Error; err != nil {
		return err
	}
	return nil
}

func (service *ImageService) AddTags(imageIDs []uint64, tags []string) error {
	tx := db.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	for _, tagName := range tags {
		var tag model.TagModel
		if err := tx.Where("tag_name = ?", tagName).First(&tag).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				tag = model.TagModel{
					TagName: tagName,
				}
				if err := tx.Create(&tag).Error; err != nil {
					tx.Rollback()
					return err
				}
			} else {
				tx.Rollback()
				return err
			}
		}
		for _, imageID := range imageIDs {
			imageTag := &model.ImageTagModel{
				ImageID: imageID,
				TagID:   tag.ID,
			}
			if err := tx.Where("image_id = ? AND tag_id = ?", imageID, tag.ID).First(&model.ImageTagModel{}).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					if err := tx.Create(imageTag).Error; err != nil {
						tx.Rollback()
						return err
					}
				}
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}
	return nil
}

type TagDetailItem struct {
	Name  string `json:"name"`
	Count int64  `json:"count"`
}

func (service *ImageService) GetTagDetails() ([]TagDetailItem, error) {
	var tags []model.TagModel
	err := db.DB.Model(&model.TagModel{}).Order("created_at ASC").Find(&tags).Error
	if err != nil {
		return nil, err
	}

	tagDetails := make([]TagDetailItem, 0)
	for _, tag := range tags {
		// 统计每个标签的图片数量
		var count int64
		err := db.DB.Table("image_tag").
			Joins("JOIN image ON image_tag.image_id = image.id").
			Where("image_tag.tag_id = ?", tag.ID).
			Count(&count).Error
		if err != nil {
			return nil, err
		}

		tagDetails = append(tagDetails, TagDetailItem{
			Name:  tag.TagName,
			Count: count,
		})
	}

	return tagDetails, nil
}

func (service *ImageService) CreateTag(tagName string) error {
	// 检查标签是否已存在
	var existingTag model.TagModel
	err := db.DB.Where("tag_name = ?", tagName).First(&existingTag).Error
	if err == nil {
		return fmt.Errorf("标签 '%s' 已存在", tagName)
	}
	if err != gorm.ErrRecordNotFound {
		return err
	}

	// 创建新标签
	tag := model.TagModel{
		TagName: tagName,
	}
	return db.DB.Create(&tag).Error
}

func (service *ImageService) UpdateTag(oldName, newName string) error {
	// 检查旧标签是否存在
	var tag model.TagModel
	err := db.DB.Where("tag_name = ?", oldName).First(&tag).Error
	if err != nil {
		return err
	}

	// 检查新标签名是否已存在
	if oldName != newName {
		var existingTag model.TagModel
		err := db.DB.Where("tag_name = ?", newName).First(&existingTag).Error
		if err == nil {
			return fmt.Errorf("标签 '%s' 已存在", newName)
		}
		if err != gorm.ErrRecordNotFound {
			return err
		}
	}

	// 更新标签名
	return db.DB.Model(&tag).Update("tag_name", newName).Error
}

func (service *ImageService) DeleteTag(tagName string) error {
	// 查找标签
	var tag model.TagModel
	err := db.DB.Where("tag_name = ?", tagName).First(&tag).Error
	if err != nil {
		return err
	}

	// 开启事务
	tx := db.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// 删除标签关联
	if err := tx.Where("tag_id = ?", tag.ID).Delete(&model.ImageTagModel{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 删除标签
	if err := tx.Where("id = ?", tag.ID).Delete(&model.TagModel{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 提交事务
	return tx.Commit().Error
}
