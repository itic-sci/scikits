package scikits

import (
	"fmt"
	"github.com/minio/minio-go/v6"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"time"
)

type MyMinIo struct {
	minioClient *minio.Client
}

func NewMyMinIo(label string) *MyMinIo {
	m := new(MyMinIo)
	m.initClient(label)
	return m
}

func (m *MyMinIo) initClient(label string) {
	endpoint := MyViper.GetString(label + ".endpoint")
	accessKeyID := MyViper.GetString(label + ".accessKeyID")
	secretAccessKey := MyViper.GetString(label + ".secretAccessKey")
	useSSL := MyViper.GetBool(label + ".useSSL")

	// 初使化minio client对象。
	client, err := minio.New(endpoint, accessKeyID, secretAccessKey, useSSL)
	if err != nil {
		log.Fatalln(err)
	}

	m.minioClient = client
}

func (m *MyMinIo) CreateBucket(bucketName string) {
	// bucketName 创建存储桶的名称

	// Location is an optional argument, by default all buckets are created in US Standard Region.
	location := "cn-north-1"

	err := m.minioClient.MakeBucket(bucketName, location)
	if err != nil {
		// 检查存储桶是否已经存在。
		exists, err := m.minioClient.BucketExists(bucketName)
		if err == nil && exists {
			log.Printf("We already own %s\n", bucketName)
		} else {
			log.Fatalln(err)
		}
	}
	log.Printf("Successfully created %s\n", bucketName)
}

func (m *MyMinIo) UploadFile(bucketName, filePath, objectName string, opts minio.PutObjectOptions) {
	// filePath 要上传文件在磁盘的路径; objectName 上传后文件的名称

	// 使用FPutObject上传一个zip文件。
	//contentType := "application/zip"
	//opts := minio.PutObjectOptions{ContentType: contentType}

	// 保留文件原格式
	//opts := minio.PutObjectOptions{}

	n, err := m.minioClient.FPutObject(bucketName, objectName, filePath, opts)
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Successfully uploaded %s of size %d\n", objectName, n)
}

// 获取文件的公开访问链接
func (m *MyMinIo) FileUrlPublic(bucketName, myObject string, expirationSecond time.Duration) (*url.URL, error) {

	presignedURL, err := m.minioClient.PresignedGetObject(bucketName, myObject, time.Second*expirationSecond, nil)

	fmt.Println(presignedURL)

	return presignedURL, err

}

func GetFileContentType(out multipart.File) (string, error) {

	// Only the first 512 bytes are used to sniff the content type.
	buffer := make([]byte, 512)

	_, err := out.Read(buffer)
	if err != nil {
		return "", err
	}

	// Use the net/http package's handy DectectContentType function. Always returns a valid
	// content-type by returning "application/octet-stream" if no others seemed to match.
	contentType := http.DetectContentType(buffer)

	return contentType, nil
}

func (m *MyMinIo) UploadWebFile(bucketName, objectName string, file *multipart.FileHeader, contentType string) error {
	// filePath 要上传文件在磁盘的路径; objectName 上传后文件的名称

	// 使用FPutObject上传一个zip文件。
	//contentType := "application/zip"
	//opts := minio.PutObjectOptions{ContentType: contentType}

	src, err1 := file.Open()
	if err1 != nil {
		fmt.Println(err1)
		return err1
	}
	defer src.Close()

	n, err := m.minioClient.PutObject(bucketName, objectName, src, -1, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		return err
	}

	log.Printf("Successfully uploaded %s of size %d\n", objectName, n)

	return nil
}

func (m *MyMinIo) RemoveObject(bucketName, objectName string) error {
	err := m.minioClient.RemoveObject(bucketName, objectName)
	if err != nil {
		fmt.Println(err)
	}
	return err

}
