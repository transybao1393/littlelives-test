package usecase

import (
	"bytes"
	"context"
	"fmt"
	"ll_test/app/logger"
	"mime/multipart"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
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
	fmt.Println("here at SaveToMinIO")
	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		fmt.Printf("error %v\n", err)
		return nil, err
	}
	fmt.Printf("minio client %#v\n", minioClient)
	return minioClient, nil
}

func SaveToMinIO(file multipart.File, contentType string, bufFile *bytes.Buffer, objectName string, fileSize int64) error {
	minioClient, _ := newMinIOClient()
	//- create new bucket
	location := "ap-southeast-1"
	createBucket(minioClient, location)

	// Upload the test file with FPutObject
	info, err := minioClient.PutObject(ctx, bucketName, objectName, bufFile, fileSize, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		fmt.Printf("Failed to upload %s\n", objectName)
		return err
	}

	log.Printf("Successfully uploaded %s of size %d\n", objectName, info.Size)
	return nil
}

func createBucket(minioClient *minio.Client, location string) {
	//- create bucket
	err := minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: location})
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, errBucketExists := minioClient.BucketExists(ctx, bucketName)

		if errBucketExists == nil && exists {
			fmt.Printf("We already own %s\n", bucketName)
			return
		} else {
			fmt.Printf("error %v", err)
			return
		}
	} else {
		fmt.Printf("Successfully created %s\n", bucketName)
	}
}
