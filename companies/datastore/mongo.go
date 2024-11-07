package datastore

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"xm/companies"
	"xm/internal/mongo"
)

const (
	companiesCollection = "companies"
)

type MongoStore struct {
	client   *mongo.Client
	database string
}

func NewMongoStore(client *mongo.Client, database string) *MongoStore {
	return &MongoStore{
		client:   client,
		database: database,
	}
}

func (s *MongoStore) Save(ctx context.Context, company companies.Company) (companies.Company, error) {
	filter := bson.M{"id": company.ID}
	update := bson.M{"$set": company}
	opts := options.Update().SetUpsert(true)

	_, err := s.client.Collection(s.database, companiesCollection).UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return companies.Company{}, err
	}

	return company, nil
}

func (s *MongoStore) Delete(ctx context.Context, id string) error {
	filter := bson.M{"id": id}
	_, err := s.client.Collection(s.database, companiesCollection).DeleteOne(ctx, filter)
	return err
}

func (s *MongoStore) FindByID(ctx context.Context, id string) (companies.Company, error) {
	filter := bson.M{"id": id}

	var company companies.Company
	err := s.client.Collection(s.database, companiesCollection).FindOne(ctx, filter).Decode(&company)
	if err != nil {
		if errors.Is(err, mongo.ErrNotExists) {
			return companies.Company{}, nil
		}
		return companies.Company{}, err
	}

	return company, nil
}

func (s *MongoStore) FindByName(ctx context.Context, name string) (companies.Company, error) {
	filter := bson.M{"name": name}

	var company companies.Company
	err := s.client.Collection(s.database, companiesCollection).FindOne(ctx, filter).Decode(&company)
	if err != nil {
		if errors.Is(err, mongo.ErrNotExists) {
			return companies.Company{}, nil
		}
		return companies.Company{}, err
	}

	return company, nil
}
