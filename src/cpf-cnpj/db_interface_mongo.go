package cpfcnpj

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	//DBMongo_Local Const used to Open Local db
	DBMongoLocal  = "mongodb://localhost:27017"
	DBMongoDocker = "mongodb://root:example@mongo:27017/"
	SetDockerRun  = true
	SetLocalRun   = false
)

//IsUsingMongoDocker If Using MongoDB in  a Docker image this var is True
var IsUsingMongoDocker bool

var collectionQuery *mongo.Collection
var ctx = context.TODO()

//CheckIsUsingMongoDocker Get If Using MongoDB in  a Docker image
func CheckIsUsingMongoDocker() bool {
	return IsUsingMongoDocker
}

//SetUsingMongoDocker set If Using MongoDB in  a Docker image
func SetUsingMongoDocker(isMongoDocker bool) {
	IsUsingMongoDocker = isMongoDocker
}

//InitDBMongo Initi MOngo Database
func InitDBMongo(isDockerRun bool) bool {
	urlMongo := DBMongoLocal
	if isDockerRun {
		urlMongo = DBMongoDocker
	}
	clientOptions := options.Client().ApplyURI(urlMongo)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
		return false
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
		return false
	}

	collectionQuery = client.Database("Querys").Collection("Querys")
	log.Println("Successfully connected to the DB MONGO!")
	return true

}
