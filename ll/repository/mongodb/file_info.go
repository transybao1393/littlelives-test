package mongodb

import (
	"context"
	"fmt"
	"ll_test/app/logger"
	"ll_test/domain"
	db "ll_test/domain/dbInstance"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ctx                  = context.Background()
	log                  = logger.NewLogrusLogger()
	mongoClient          = db.GetMongoInstance()
	DATABASE_NAME        = "littlelives"
	FILE_COLLECTION_NAME = "file_info"
	file_collection      = mongoClient.Database(DATABASE_NAME).Collection(FILE_COLLECTION_NAME)
)

func init() {
	//- TODO: Create index
	indexModels := []mongo.IndexModel{
		{
			Keys: bson.D{primitive.E{
				Key:   "_id",
				Value: 1,
			}},
		},
	}

	name, err := user_collection.Indexes().CreateMany(context.TODO(), indexModels)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Created index file collection %s\n", name)
}

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
	_, err := file_collection.InsertOne(ctx, fi)
	if err != nil {
		handleError(err, "Error when insert into mongodb", "fatal")
		return false
	}

	return true
}

func (fi *FileInfo) GetFileInfoByIP(IP string) *FileInfo {
	filter := bson.D{primitive.E{Key: "userip", Value: IP}}

	// retrieve all the documents that match the filter
	cursor, err := file_collection.Find(context.TODO(), filter)
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
