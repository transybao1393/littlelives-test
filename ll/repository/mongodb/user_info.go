package mongodb

import (
	"context"
	"fmt"
	"ll_test/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserInfo domain.User
type UserActions interface {
	AddNewUser() error
	GetUserByIP() (*UserInfo, error)
	GetQuotaByUserIP() (int, error)
	UpdateQuotaUsedByUserIP() bool
}

func (ui *UserInfo) AddNewUser() error {
	result, err := collection.InsertOne(ctx, ui)
	if err != nil {
		handleError(err, "Error when insert into mongodb", "fatal")
		return err
	}

	fmt.Println(result.InsertedID)
	return nil
}

func (ui *UserInfo) GetUserByIP() (*UserInfo, error) {
	filter := bson.D{primitive.E{Key: "userip", Value: ui.UserIP}}

	// retrieve all the documents that match the filter
	cursor, err := collection.Find(context.TODO(), filter)
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

	// retrieve all the documents that match the filter
	cursor, err := collection.Find(context.TODO(), filter)
	// check for errors in the finding
	if err != nil {
		return 0, err
	}
	// check for errors in the conversion
	if err = cursor.All(context.TODO(), &ui); err != nil {
		return 0, err
	}
	return ui.Quota, nil
}

func (ui *UserInfo) UpdateQuotaUsedByUserIP() bool {
	//- get quota by user ip
	quota, _ := ui.GetQuotaByUserIP()
	filter := bson.D{primitive.E{Key: "userip", Value: ui.UserIP}}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "quotaUsed", Value: quota + 1},
		}},
	}
	_, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		handleError(err, "Error when update quota used", "fatal")
		return false
	}

	return true
}
