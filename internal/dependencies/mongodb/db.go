package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DbService struct {
	Client *mongo.Client
}

func NewDbService(uri string) *DbService {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	return &DbService{Client: client}
}

func (db *DbService) Disconnect() error {
	return db.Client.Disconnect(context.TODO())
}
