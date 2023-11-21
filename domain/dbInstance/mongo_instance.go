package dbInstance

import (
	"context"
	"ll_test/app/logger"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	CONNECTION_STRING = "mongodb://mongo:27017"
	DB_NAME           = "db_file_management"
)

type singleMongoInstance struct {
	Conn *mongo.Client
}

var mongoClient *singleMongoInstance

func GetMongoInstance() *mongo.Client {
	initOnce.Do(func() {
		// Set client options
		clientOptions := options.Client().ApplyURI(CONNECTION_STRING)
		// Connect to MongoDB
		client, err := mongo.Connect(context.TODO(), clientOptions)
		if err != nil {
			fields := logger.Fields{
				"type":    "Mongodb",
				"status":  "FAILED",
				"message": "Cannot establish mongo instance base on PING signal not response properly",
			}
			//- this error will be effected of the flow of mongo connection => fatal error
			log.Fields(fields).Fatalf(err, "Cannot connect to mongo instance")
			ctx.Done()
			panic(err)
		}
		// Check the connection
		err = client.Ping(context.TODO(), nil)
		if err != nil {
			fields := logger.Fields{
				"type":    "Mongodb",
				"status":  "FAILED",
				"message": "Cannot establish mongo instance base on PING signal not response properly",
			}
			//- this error will be effected of the flow of mongo connection => fatal error
			log.Fields(fields).Fatalf(err, "Cannot establish mongo instance base on PING signal not response properly")
			ctx.Done()
			panic(err)
		}
		// clientInstance = client
		mongoClient = &singleMongoInstance{
			Conn: client,
		}

		//- if success
		fields := logger.Fields{
			"result": "Database connect successfully",
			"status": "SUCCESS",
		}
		log.Fields(fields).Infof("PING result from mongo instance")

	})
	return mongoClient.Conn
}
