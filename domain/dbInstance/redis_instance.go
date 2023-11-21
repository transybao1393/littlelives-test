package dbInstance

import (
	"context"
	"ll_test/app/logger"
	"sync"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()
var log = logger.NewLogrusLogger()

type singleRedisInstance struct {
	Conn *redis.Client
}

var redisClient *singleRedisInstance
var initOnce sync.Once

func GetRedisInstance() *redis.Client {
	initOnce.Do(func() {
		redisClient = &singleRedisInstance{
			Conn: redis.NewClient(&redis.Options{
				Network:  "tcp",
				Addr:     "redis:6379",
				Password: "", // no password set
				DB:       0,  // use default DB
			}),
		}

		//- check connection after create new client
		result, err := redisClient.Conn.Ping(ctx).Result()
		//- if failed
		if err != nil {
			fields := logger.Fields{
				"db-type": "redis",
				"status":  "FAILED",
			}
			//- this error will be effected of the flow of redis connection => fatal error
			log.Fields(fields).Fatalf(err, "Cannot establish redis instance base on PING signal not response properly")
			ctx.Done()
			panic(err)
		}

		//- if success
		fields := logger.Fields{
			"result": result,
			"status": "SUCCESS",
		}
		log.Fields(fields).Infof("PING result from redis instance")

	})
	return redisClient.Conn
}
