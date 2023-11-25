package usecase

import (
	"bytes"
	"context"
	"fmt"
	"ll_test/app/logger"
	"mime/multipart"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	mongoRepository "ll_test/ll/repository/mongodb"
)

var log = logger.NewLogrusLogger()
var ctx = context.Background()
var bucketName = "testbucket"

func handleError(err error, message string, errorType string) {
	fields := logger.Fields{
		"service": "littlelives",
		"message": message,
	}
	switch errorType {
	case "fatal":
		log.Fields(fields).Fatalf(err, message)
	case "error":
		log.Fields(fields).Errorf(err, message)
	case "warn":
		log.Fields(fields).Warnf(message)
	case "info":
		log.Fields(fields).Infof(message)
	case "debug":
		log.Fields(fields).Debugf(message)
	}
}

func newMinIOClient() (*minio.Client, error) {
	endpoint := "minio:9000"
	accessKeyID := "Q3AM3UQ867SPQQA43P2F"
	secretAccessKey := "zuf+tfteSlswRu7BJ86wekitnifILbZam1KYY3TG"
	useSSL := false

	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		handleError(err, "Error when initialize minio client", "fatal")
		return nil, err
	}
	return minioClient, nil
}

func SaveToMinIO(file multipart.File, contentType string, bufFile *bytes.Buffer, objectName string, fileSize int64, IP string) error {
	minioClient, _ := newMinIOClient()

	//- create new bucket
	location := "ap-southeast-1"
	createBucket(minioClient, location)

	// Upload the test file with PutObject
	info, err := minioClient.PutObject(ctx, bucketName, objectName, bufFile, fileSize, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		handleError(err, "Error when uploading file to bucket", "fatal")
		return err
	}

	//- save to mongodb
	fileInfo := &mongoRepository.FileInfo{
		UserIP:         IP,
		FileName:       objectName,
		FileSize:       fileSize,
		FileType:       contentType,
		FileBucketPath: fmt.Sprintf("http://%s:%s/%s/%s", "localhost", "9000", bucketName, objectName),
	}

	if !fileInfo.AddFileInfo() {
		handleError(err, "Error when add file information", "fatal")
		return nil
	}

	handleError(err, fmt.Sprintf("Successfully uploaded %s of size %d\n", objectName, info.Size), "info")
	return nil
}

func createBucket(minioClient *minio.Client, location string) {
	//- create bucket
	err := minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: location})
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, errBucketExists := minioClient.BucketExists(ctx, bucketName)
		if errBucketExists == nil && exists {
			//- bucket list exist
			return
		} else {
			handleError(err, "Bucket is not exist", "fatal")
			return
		}
	} else {
		handleError(err, fmt.Sprintf("Successfully created new bucket name %s", bucketName), "info")

	}
}
