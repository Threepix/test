package mongodb

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewClient(ctx context.Context, host, port, username, password, database, authDb string) (db *mongo.Database, err error) {
	var mongoDbURL string
	var isAuth bool
	if username == "" && password == "" {
		mongoDbURL = fmt.Sprintf("mongodb://%s:%s", host, port)
	} else {
		isAuth = true
		mongoDbURL = fmt.Sprintf("mongodb://%s:%s@%s:%s", username, password, host, port)
	}
	clientOption := options.Client().ApplyURI(mongoDbURL)
	if isAuth {
		if authDb == "" {
			authDb = database
		}
		clientOption.SetAuth(options.Credential{
			AuthSource: authDb,
			Username:   username,
			Password:   password,
		})
	}

	client, err := mongo.Connect(ctx, clientOption)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to db: %v", err)
	}
	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping db: %v", err)
	}
	return client.Database(database), nil
}