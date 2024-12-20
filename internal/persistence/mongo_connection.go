package persistence

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateMongoConnection() *mongo.Client {

	host := viper.GetString("mongodb.host")
	user := viper.GetString("mongodb.user")
	password := viper.GetString("mongodb.password")

	uri := fmt.Sprintf("mongodb://%v:%v@%v:27017", user, password, host)

	clientOptions := options.Client().ApplyURI(uri)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	return client
}
