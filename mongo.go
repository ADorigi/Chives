package chives

import (
	"context"
	"errors"
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

func (db *Mongodb) CreateDocument(ctx context.Context, collectionName string, document bson.D) (primitive.ObjectID, error) {

	log.Println("chives: creating document")
	collection := db.Client.Database(db.Database).Collection(collectionName)

	response, err := collection.InsertOne(ctx, document)
	if err != nil {
		log.Println("Cannot insert document")
		return primitive.NilObjectID, err
	}

	log.Println("chives: document created")
	return response.InsertedID.(primitive.ObjectID), nil
}

func (db *Mongodb) GetDocument(ctx context.Context, collectionName string, filter primitive.E) (*mongo.SingleResult, error) {

	log.Println("chives: getting object id")

	collection := db.Client.Database(db.Database).Collection(collectionName)

	singleresult := collection.FindOne(ctx, bson.D{filter})
	if singleresult.Err() != nil {
		log.Printf("chives: %s", singleresult.Err())
		return nil, errors.New("chives: no document found")
	}

	log.Println("chives: object recovered")
	return singleresult, nil
}
