package model

import "time"

type ImageDTO struct {
	ImageModel
	Tags []string `json:"tags"`
}

type ImageModel struct {
	ID            uint64    `json:"id" gorm:"column:id;primary_key;auto_increment"`
	ImageName     string    `json:"image_name" gorm:"column:image_name"`
	ImageCode     string    `json:"image_code" gorm:"column:image_code"`
	ThumbnailCode string    `json:"thumbnail_code" gorm:"column:thumbnail_code"`
	Ext           string    `json:"ext" gorm:"column:ext"`
	Size          int64     `json:"size" gorm:"column:size"`
	Directory     string    `json:"directory" gorm:"column:directory"`
	CreatedAt     time.Time `json:"created_at" gorm:"column:created_at"`
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
