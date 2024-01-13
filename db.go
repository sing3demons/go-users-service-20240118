package main

import (
	"context"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func NewDatabase(dbName, collectionName string) *mongo.Collection {

	dsn := os.Getenv("MONGO_URI")
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dsn))
	if err != nil {
		log.WithFields(log.Fields{
			"error":          err.Error(),
			"type":           "database",
			"func":           "NewDatabase",
			"file":           "db.go",
			"tag":            "error",
			"host":           dsn,
			"dbName":         dbName,
			"collectionName": collectionName,
		}).Error("error connecting to database")
		os.Exit(1)
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		log.WithFields(log.Fields{
			"error":          err.Error(),
			"type":           "database",
			"func":           "NewDatabase",
			"file":           "db.go",
			"tag":            "error",
			"host":           dsn,
			"dbName":         dbName,
			"collectionName": collectionName,
		}).Error("error connecting to database")
		os.Exit(1)
	}

	log.WithFields(log.Fields{
		"type":           "database",
		"host":           dsn,
		"dbName":         dbName,
		"collectionName": collectionName,
		"status":         "OK",
	}).Info("connected to database")

	return client.Database(dbName).Collection(collectionName)
}
