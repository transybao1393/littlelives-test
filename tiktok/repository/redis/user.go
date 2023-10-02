package redis

import (
	"context"
	"tiktok_api/app/logger"
	"tiktok_api/domain/dbInstance"
	"time"
)

var clientInstance = dbInstance.GetRedisInstance()
var log = logger.NewLogrusLogger()
var ctx = context.Background()

func AddSimpleUser(key string, value string) bool {
	err := clientInstance.Set(ctx, key, value, 0).Err()
	if err != nil {
		//- writing logs
		log.Fields(logger.Fields{
			"key":   key,
			"value": value,
			"error": err,
		}).Errorf(err, "Error when set new key-value into redis")
		return false
	}
	return true
}

func GetSimpleUser(key string) string {
	val, err := clientInstance.Get(ctx, key).Result()
	if err != nil {
		panic(err)
	}
	log.Fields(logger.Fields{
		"key":  key,
		"date": time.Now(),
	}).Info("Get result from redis")
	return val
}
