package test

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	//- compulsory fields
	ID        primitive.ObjectID
	UserIP    string
	Quota     int
	QuotaUsed int
	CreatedAt time.Time
	UpdatedAt time.Time
}
