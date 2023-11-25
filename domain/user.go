package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	//- compulsory fields
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	UserIP    string             `json:"userip" bson:"userip,omitempty"`
	Quota     int                `json:"quota" bson:"quota"` //- 100
	QuotaUsed int                `json:"quotaUsed" bson:"quotaUsed"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt" bson:"updatedAt"`
}
