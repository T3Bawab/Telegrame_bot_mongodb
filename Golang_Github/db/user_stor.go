package db

import (
	"T3B/bot_settings"
	"T3B/types"
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserStore interface {
	GetUserByID(context.Context, string) (*types.User, error)
	CreateUser(context.Context, *types.User) (*types.User, error)
	CheckUsername(context.Context, string) (bool, error)
	CheckTeleID(context.Context, int64) (bool, error)
	DeleteUser(context.Context, string) error
}

type MongoUserStore struct {
	bot_settings.Bot
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoUserStore(client *mongo.Client) *MongoUserStore {
	return &MongoUserStore{
		client: client,
		coll:   client.Database(DBName).Collection(DBName)}
}

func (s *MongoUserStore) GetUserByID(ctx context.Context, id string) (*types.User, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("ID Is Not Correct")
	}
	var user types.User

	if err := s.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *MongoUserStore) CreateUser(ctx context.Context, user *types.User) (*types.User, error) {

	res, err := s.coll.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}

	user.ID = res.InsertedID.(primitive.ObjectID)
	return user, nil

}

func (s *MongoUserStore) CheckUsername(ctx context.Context, u string) (bool, error) {
	var user *types.User

	err := s.coll.FindOne(ctx, bson.M{"user": u}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return true, nil
		}
		return true, err
	}

	return false, fmt.Errorf("username %s is unavailable", u)
}

func (s *MongoUserStore) CheckTeleID(ctx context.Context, TeleId int64) (bool, error) {
	var user *types.User

	// To check if the user is already in the database
	err := s.coll.FindOne(ctx, bson.M{"teleid": TeleId}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return true, nil
		}
		return true, err
	}

	return false, fmt.Errorf("you already sign in %d", TeleId)
}

func (s *MongoUserStore) DeleteUser(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	// TODO: Maybe its a good idea to handle if we did not delete any user.
	// maybe log it or something??
	e, err := s.coll.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil || e.DeletedCount == 0 {
		return fmt.Errorf("there is no user with id %s", id)
	}
	return nil
}
