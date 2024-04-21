package chives

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Mongo interface {
	Connect(context.Context)
	GetObjectID(context.Context, primitive.E) (primitive.ObjectID, error)
	// GetDocumentByID(context.Context, primitive.ObjectID) (any, error)
	CreateDocument(context.Context, bson.D) (primitive.ObjectID, string, error)
}

type Mongodb struct {
	URI      string
	Database string
	Client   *mongo.Client
}

func NewMongodb(URI string, Database string) *Mongodb {
	log.Println("chives: creating mongodb instance")
	return &Mongodb{
		URI:      URI,
		Database: Database,
	}
}

func (db *Mongodb) Connect(ctx context.Context) {
	log.Println("chives: connecting to mongodb database")
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().ApplyURI(db.URI).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		panic(err)
	}
	db.Client = client

	log.Println("chives: mongodb database connected")
}

// everything after this step
// needs to be done in the local library
