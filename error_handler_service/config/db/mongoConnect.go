// Copyright 2025 Naveen R
//
//	Licensed under the Apache License, Version 2.0 (the "License");
//	you may not use this file except in compliance with the License.
//	You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
//	Unless required by applicable law or agreed to in writing, software
//	distributed under the License is distributed on an "AS IS" BASIS,
//	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//	See the License for the specific language governing permissions and
//	limitations under the License.

package db

import (
	"context"
	"log"
	"project_chimera/error_handle_service/config"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	logger "project_chimera/error_handle_service/pkg/logger"
)

var DB *mongo.Client

// ConnectDB establishes a connection to the database.
func ConnectDB() {
	mongoURI := config.Env.MongoDBURI
	clientOptions := options.Client().ApplyURI(mongoURI)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var err error
	DB, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Ping MongoDB
	err = DB.Ping(ctx, nil)
	if err != nil {
		logger.LogFatal("Could not connect to MongoDB: " + err.Error())
	} else {
		logger.LogInfo("Connected to MongoDB!")
	}
}

// GetCollection returns a MongoDB collection for a specific database and collection name.
func GetCollection(databaseName, collectionName string) *mongo.Collection {
	if DB == nil {
		logger.LogFatal("Database connection is not initialized")
	}

	return DB.Database(databaseName).Collection(collectionName)
}
