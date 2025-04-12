package service

import (
	"mime/multipart"
	"path/filepath"
	"picture_storage/db"
	"picture_storage/model"
	"picture_storage/pkg/minio"
	"strings"
	"time"

	"gorm.io/gorm"
)

type ImageModel struct {
	ID        uint64    `json:"id" gorm:"column:id;primary_key;auto_increment"`
	ImageName string    `json:"image_name" gorm:"column:image_name"`
	ImageCode string    `json:"image_code" gorm:"column:image_code"`
	Ext       string    `json:"ext" gorm:"column:ext"`
	Size      int64     `json:"size" gorm:"column:size"`
	Directory string    `json:"directory" gorm:"column:directory"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
}

type ImageTagModel struct {
	ID        uint64    `json:"id" gorm:"column:id;primary_key;auto_increment"`
	ImageID   uint64    `json:"image_id" gorm:"column:image_id"`
	TagID     uint64    `json:"tag_id" gorm:"column:tag_id"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
}

type TagModel struct {
	ID        uint64    `json:"id" gorm:"column:id;primary_key;auto_increment"`
	TagName   string    `json:"tag_name" gorm:"column:tag_name"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
}

func (*TagModel) TableName() string {
	return "tag"
}

func (*ImageTagModel) TableName() string {
	return "image_tag"
}

func (*ImageModel) TableName() string {
	return "image"
}

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

func (service *ImageService) GetImageListByDirectory(directory string, page model.Pagination) ([]ImageModel, int64, error) {
	imageList := make([]ImageModel, 0)
	err := db.DB.Where("directory = ?", directory).
		Offset((page.Page - 1) * page.PageSize).
		Limit(page.PageSize).
		Order("created_at DESC").
		Find(&imageList).Error
	if err != nil {
		return nil, 0, err
	}
	var total int64
	err = db.DB.Model(&ImageModel{}).Where("directory = ?", directory).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	return imageList, total, nil
}

func (service *ImageService) GetImageListByTag(directory string, tag string, page model.Pagination) ([]ImageModel, int64, error) {
	// TODO: 这里需要实现从数据库获取图片列表的逻辑
	// 临时返回空数据
	return []ImageModel{}, 0, nil
}

func (service *ImageService) SaveImage(directory string, file *multipart.FileHeader, tags []string) (uint64, error) {
	// 上传到 MinIO
	imageCodeWithExt, size, err := service.UploadImage(directory, file)
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
	var image *ImageModel = &ImageModel{}
	err = tx.Model(&ImageModel{}).Where("image_code = ?", imageCode).Find(image).Error
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	if image.ID != 0 {
		return image.ID, nil
	}

	// 保存图片信息
	image = &ImageModel{
		ImageName: file.Filename,
		ImageCode: imageCode,
		Directory: directory,
		Ext:       extension,
		Size:      size,
	}
	if err := tx.Create(image).Error; err != nil {
		tx.Rollback()
		return 0, err
	}

	// 处理标签
	if len(tags) > 0 {
		// 先删除旧的标签关联
		if err := tx.Where("image_id = ?", image.ID).Delete(&ImageTagModel{}).Error; err != nil {
			tx.Rollback()
			return 0, err
		}

		// 处理每个标签
		for _, tagName := range tags {
			// 查找或创建标签
			var tag TagModel
			if err := tx.Where("tag_name = ?", tagName).First(&tag).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					// 标签不存在，创建新标签
					tag = TagModel{
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
			imageTag := &ImageTagModel{
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
