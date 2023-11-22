package usecase

import (
	"bytes"
	"context"
	"ll_test/app/logger"
	"ll_test/domain"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"golang.org/x/oauth2"
)

var log = logger.NewLogrusLogger()
var ctx = context.Background()
var config *oauth2.Config

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

func SaveToMinIO() error {
	endpoint := "http://minio1/data1"
	accessKeyID := "Q3AM3UQ867SPQQA43P2F"
	secretAccessKey := "zuf+tfteSlswRu7BJ86wekitnifILbZam1KYY3TG"
	useSSL := true

	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return err
	}

	log.Printf("%#v\n", minioClient)
	return nil
}

func LLVideoUploadFile(fileBuffer *bytes.Buffer, clientKey string, ytbFileUploadInfo *domain.YoutubeFileUploadInfo) (string, error) {
	//- upload video

	//- save info to mongodb

	//- save file to minio
	return "", nil
}
