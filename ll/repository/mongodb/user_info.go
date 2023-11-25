package mongodb

import (
	"context"
	"errors"
	"fmt"
	"ll_test/domain"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	USER_COLLECTION_NAME = "user_info"
	user_collection      = mongoClient.Database(DATABASE_NAME).Collection(USER_COLLECTION_NAME)
)

func init() {
	//- TODO: Create index
	indexModels := []mongo.IndexModel{
		{
			Keys: bson.D{primitive.E{
				Key:   "userip",
				Value: 1,
			}},
		},
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
	fmt.Printf("Created index for user collection %s\n", name)
}

type UserInfo domain.User
type UserActions interface {
	AddNewUser() error
	GetUserByIP() (*UserInfo, error)
	GetQuotaByUserIP() (int, error)
	UpdateQuotaUsedByUserIP() bool
}

//- create

func (ui *UserInfo) AddNewUser() error {
	userCount := ui.countByUserIP()
	if userCount == 0 {
		ui.CreatedAt = time.Now()
		_, err := user_collection.InsertOne(ctx, ui)
		if err != nil {
			handleError(err, "Error when insert into mongodb", "fatal")
			return err
		}
	}
	return nil
}

func (ui *UserInfo) countByUserIP() int64 {
	userCount, err := user_collection.CountDocuments(ctx, bson.M{"userip": ui.UserIP})
	if err != nil {
		handleError(err, "Error when count user by IP", "fatal")
		return 0
	}
	return userCount
}

func (ui *UserInfo) GetUserByIP() (*UserInfo, error) {
	filter := bson.D{primitive.E{Key: "userip", Value: ui.UserIP}}
	// retrieve all the documents that match the filter
	cursor, err := user_collection.Find(context.TODO(), filter)
	// check for errors in the finding
	if err != nil {
		handleError(err, "Error when finding results", "fatal")
		return nil, err
	}
	// check for errors in the conversion
	if err = cursor.All(context.TODO(), &ui); err != nil {
		handleError(err, "Error when convert cursor to model", "fatal")
		return nil, err
	}
	return ui, nil
}

func (ui *UserInfo) GetQuotaByUserIP() (int, error) {
	filter := bson.D{primitive.E{Key: "userip", Value: ui.UserIP}}

	err := user_collection.FindOne(ctx, filter).Decode(&ui)
	if err != nil {
		return 0, err
	}

	if !(ui.QuotaUsed >= ui.Quota) {
		return ui.Quota, nil
	}
	return 0, errors.New("quota exceeded! Please buy more quota")
}

func (ui *UserInfo) UpdateQuotaUsedByUserIP() (bool, error) {
	//- get quota by user ip
	_, err := ui.GetQuotaByUserIP()
	if err != nil {
		handleError(err, "Quota exceeded! Please buy more quota", "fatal")
		return false, err
	}

	filter := bson.D{primitive.E{Key: "userip", Value: ui.UserIP}}
	update := bson.D{
		{
			Key: "$set",
			Value: bson.D{
				{
					Key:   "quotaUsed",
					Value: ui.QuotaUsed + 1,
				},
			},
		},
		{
			Key: "$set",
			Value: bson.D{
				{
					Key:   "updatedAt",
					Value: time.Now(),
				},
			},
		},
	}
	_, err = user_collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		handleError(err, "Error when update quota used", "fatal")
		return false, err
	}

	return true, nil
}
