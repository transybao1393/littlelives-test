package mongodb

import (
	"context"
	"fmt"
	"ll_test/app/logger"
	"ll_test/domain"
	db "ll_test/domain/dbInstance"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ctx             = context.Background()
	log             = logger.NewLogrusLogger()
	mongoClient     = db.GetMongoInstance()
	collection      = mongoClient.Database(DATABASE_NAME).Collection(COLLECTION_NAME)
	DATABASE_NAME   = "littlelives"
	COLLECTION_NAME = "file_info"
)

type FileInfo domain.FileInfo
type FileActions interface {
	AddFileInfo() bool
	UpdateFileInfo() bool
	GetFileInfoByIP(IP string) *FileInfo
}

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

func (fi *FileInfo) AddFileInfo() bool {
	result, err := collection.InsertOne(ctx, fi)
	if err != nil {
		handleError(err, "Error when insert into mongodb", "fatal")
		return false
	}

	fmt.Println(result.InsertedID)

	return true
}

func (fi *FileInfo) GetFileInfoByIP(IP string) *FileInfo {
	filter := bson.D{primitive.E{Key: "userip", Value: IP}}

	// retrieve all the documents that match the filter
	cursor, err := collection.Find(context.TODO(), filter)
	// check for errors in the finding
	if err != nil {
		panic(err)
	}

	// check for errors in the conversion
	if err = cursor.All(context.TODO(), &fi); err != nil {
		panic(err)
	}

	return fi
}

func (fi *FileInfo) UpdateFileInfo() bool {
	return true
}
