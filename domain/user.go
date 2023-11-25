package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	//- compulsory fields
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	UserIP    string             `json:"userip" bson:"userip,omitempty"`
	Quota     int                `json:"quota" bson:"quota"` //- 100
	QuotaUsed int                `json:"quotaUsed" bson:"quotaUsed"`
	UpdateAt  primitive.DateTime `json:"updateAt" bson:"updateAt"`
}
