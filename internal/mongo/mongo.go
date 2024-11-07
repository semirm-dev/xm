package mongo

import (
	"context"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var ErrNotExists = mongo.ErrNoDocuments

type Client struct {
	libClient *mongo.Client
}

func NewClient(uri string) *Client {
	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		logrus.Fatal(err)
	}

	return &Client{
		libClient: client,
	}
}

func (c *Client) Disconnect(ctx context.Context) {
	if err := c.libClient.Disconnect(ctx); err != nil {
		logrus.Fatal(err)
	}
}

func (c *Client) Collection(database, collection string) *mongo.Collection {
	return c.libClient.Database(database).Collection(collection)
}
