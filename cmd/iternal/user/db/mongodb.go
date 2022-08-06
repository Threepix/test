package db

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"restapi/cmd/iternal/user"
	"restapi/cmd/pkg/logging"
)

type db struct {
	collection *mongo.Collection
	logger     *logging.Logger
}

func (d *db) Create(ctx context.Context, user user.User) (string, error) {
	d.logger.Info("create user")
	result, err := d.collection.InsertOne(ctx, user)
	if err != nil {
		return "", fmt.Errorf("failed to create user %v", err)
	}
	d.logger.Debug("convert insertID")
	oid, ok := result.InsertedID.(primitive.ObjectID)
	if ok {
		return oid.Hex(), nil
	}
	d.logger.Trace(user)
	return "", fmt.Errorf("failed to convert")
}

func (d *db) FINDAll(ctx context.Context) (u []user.User, err error) {
	cursor, err := d.collection.Find(ctx, bson.M{})
	if cursor.Err() != nil {
		return u, fmt.Errorf("failed to find avv users %v", err)
	}

	if err = cursor.All(ctx, &u); err != nil {
		return u, fmt.Errorf("failed to read all documents from cursor %v", err)
	}
	return u, nil
}

func (d *db) FINDOne(ctx context.Context, id string) (u user.User, err error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return u, fmt.Errorf("failed to convert hex to id %s", id)
	}
	filter := bson.M{"_id": oid}

	result := d.collection.FindOne(ctx, filter)
	if result.Err() != nil {
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			return u, fmt.Errorf("404")
		}
		return u, fmt.Errorf("failed to find %s", id)
	}
	if err = result.Decode(&u); err != nil {
		return u, fmt.Errorf("failed to decode %s", id, err)
	}
	return u, nil
}

func (d *db) UPDATEOne(ctx context.Context, user user.User) error {
	objectiD, err := primitive.ObjectIDFromHex(user.ID)
	if err != nil {
		return fmt.Errorf("failed to convert UID = %s", user.ID)
	}
	filter := bson.M{"_id": objectiD}

	userBytes, err := bson.Marshal(user)
	if err != nil {
		return fmt.Errorf("failed ro marshal user, error: %v", err)
	}
	var updateUserObj bson.M
	err = bson.Unmarshal(userBytes, &updateUserObj)
	if err != nil {
		return fmt.Errorf("failed to unmarshal user bites error: %v", err)
	}
	delete(updateUserObj, "_id")
	update := bson.M{
		"$set": updateUserObj,
	}
	result, err := d.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update user query %v", err)
	}
	if result.MatchedCount == 0 {
		return fmt.Errorf("not found")
	}
	d.logger.Tracef("matched %d", result.MatchedCount, result.ModifiedCount)
	return nil
}

func (d *db) DELETEOne(ctx context.Context, id string) error {
	objectiD, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("failed to convert user ID to objID = %s", id)
	}
	filter := bson.M{"_id": objectiD}
	result, err := d.collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to execute query error: %v", err)
	}
	if result.DeletedCount == 0 {
		return fmt.Errorf("404")
	}
	d.logger.Tracef("matched %d", result.DeletedCount)
	return nil
}

func NewStorage(database *mongo.Database, collection string, logger *logging.Logger) user.Storage {
	return &db{
		collection: database.Collection(collection),
		logger:     logger,
	}
}
