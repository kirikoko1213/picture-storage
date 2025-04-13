package minio

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"io"
	"log"
	"net/url"
	"path/filepath"
	"picture_storage/config"
	"strconv"
	"time"

	"github.com/kiririx/krutils/ut"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var Client *MinioClient

type MinioClient struct {
	client *minio.Client
}

func (m *MinioClient) GetObjectURL(bucketName, objectName string) string {
	presignedURL, err := m.client.PresignedGetObject(context.Background(), bucketName, objectName, time.Hour*24, make(url.Values))
	if err != nil {
		log.Fatalln(err)
	}

	// return "http://" + m.client.EndpointURL().Host + "/" + bucketName + "/" + objectName
	return presignedURL.String()
}

func (m *MinioClient) UploadFile(bucketName, originalFilename string, fileSize int64, fileContent io.Reader, contentType string) (string, int64, error) {
	ctx := context.Background()

	// 读取文件内容
	content, err := io.ReadAll(fileContent)
	if err != nil {
		return "", 0, err
	}

	// 计算文件内容的 MD5
	hash := md5.New()
	hash.Write(content)
	md5Hash := hex.EncodeToString(hash.Sum(nil))

	// 获取文件扩展名
	ext := filepath.Ext(originalFilename)
	// 使用 MD5 作为文件名，保留原始扩展名
	objectName := md5Hash + ext

	// 确保桶存在
	exists, err := m.client.BucketExists(ctx, bucketName)
	if err != nil {
		return "", 0, err
	}
	if !exists {
		err = m.client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return "", 0, err
		}
	}

	// 检查文件是否已存在
	info, err := m.client.StatObject(ctx, bucketName, objectName, minio.StatObjectOptions{})
	if err == nil {
		return objectName, info.Size, nil
	}

	// 上传文件
	uploadInfo, err := m.client.PutObject(ctx, bucketName, objectName, bytes.NewReader(content), fileSize, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", 0, err
	}

	return objectName, uploadInfo.Size, nil
}

func (m *MinioClient) UploadFileBytes(bucketName, originalFilename string, fileSize int64, content []byte, contentType string) (string, int64, error) {
	ctx := context.Background()

	// 计算文件内容的 MD5
	hash := md5.New()
	hash.Write(content)
	md5Hash := hex.EncodeToString(hash.Sum(nil))

	// 获取文件扩展名
	ext := filepath.Ext(originalFilename)
	// 使用 MD5 作为文件名，保留原始扩展名
	objectName := md5Hash + ext

	// 确保桶存在
	exists, err := m.client.BucketExists(ctx, bucketName)
	if err != nil {
		return "", 0, err
	}
	if !exists {
		err = m.client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return "", 0, err
		}
	}

	// 检查文件是否已存在
	info, err := m.client.StatObject(ctx, bucketName, objectName, minio.StatObjectOptions{})
	if err == nil {
		return objectName, info.Size, nil
	}

	// 上传文件
	uploadInfo, err := m.client.PutObject(ctx, bucketName, objectName, bytes.NewReader(content), fileSize, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", 0, err
	}

	return objectName, uploadInfo.Size, nil
}

func (m *MinioClient) GetDirectoryList() ([]minio.BucketInfo, error) {
	ctx := context.Background()

	// 获取目录列表
	objects, err := m.client.ListBuckets(ctx)
	if err != nil {
		return nil, err
	}

	return objects, nil
}

func NewMinioClient(endpoint, accessKeyID, secretAccessKey string, useSSL bool) *MinioClient {
	// 初始化 MinIO 客户端
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalln(err)
	}

	return &MinioClient{
		client: minioClient,
	}
}

func InitMinioClient() {
	// 初始化 MinIO 客户端
	endpoint := ut.String().DefaultIfEmpty(config.H.Get("minio.endpoint"), "localhost:9000")
	accessKeyID := ut.String().DefaultIfEmpty(config.H.Get("minio.accessKeyID"), "minioadmin")
	secretAccessKey := ut.String().DefaultIfEmpty(config.H.Get("minio.secretAccessKey"), "minioadmin")
	useSSL := ut.String().DefaultIfEmpty(config.H.Get("minio.useSSL"), "false")
	useSSLBool, _ := strconv.ParseBool(useSSL)
	Client = NewMinioClient(endpoint, accessKeyID, secretAccessKey, useSSLBool)
}
