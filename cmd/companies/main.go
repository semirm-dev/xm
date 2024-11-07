package main

import (
	"context"
	"flag"
	"xm/companies"
	"xm/companies/datastore"
	"xm/companies/notifiers"
	"xm/internal/mongo"
)

var (
	mongoUri = flag.String("mongo", "mongodb://localhost:27017", "Mongo URI")
	mongoDb  = flag.String("mongodb", "xm", "Mongo DB")
	grpcAddr = flag.String("addr", ":8001", "gRPC address to listen connections")
)

func main() {
	flag.Parse()

	mongoCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mongoClient := mongo.NewClient(*mongoUri)
	defer mongoClient.Disconnect(mongoCtx)

	svcCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ds := datastore.NewMongoStore(mongoClient, *mongoDb)
	svc := companies.NewService(ds, notifiers.NewNoopNotifier())
	go svc.ListenForNotifications(svcCtx)

	grpcCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	grpcServer := companies.NewGrpcServer(*grpcAddr, svc)
	grpcServer.StartListening(grpcCtx)
}
