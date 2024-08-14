package store

import (
	"context"
	"fmt"
	"log"

	"github.com/tnguven/hotel-reservation-app/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const userCollection = "users"

type UserStore interface {
	Dropper

	GetByID(context.Context, string) (*types.User, error)
	GetUserByEmail(context.Context, string) (*types.User, error)
	GetUsers(context.Context) ([]*types.User, error)
	InsertUser(context.Context, *types.User) (*types.User, error)
	DeleteUser(context.Context, string) error
	PutUser(context.Context, *types.UpdateUserParams, string) error
}

type MongoUserStore struct {
	db   *mongo.Database
	coll *mongo.Collection
}

func NewMongoUserStore(db *mongo.Database) *MongoUserStore {
	return &MongoUserStore{
		db:   db,
		coll: db.Collection(userCollection),
	}
}

func (ms *MongoUserStore) GetByID(ctx context.Context, id string) (*types.User, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var user *types.User
	if err := ms.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&user); err != nil {
		return nil, err
	}

	return user, nil
}

func (ms *MongoUserStore) GetUserByEmail(ctx context.Context, email string) (*types.User, error) {
	var user types.User

	if err := ms.coll.FindOne(ctx, bson.M{"email": email}).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (ms *MongoUserStore) GetUsers(ctx context.Context) ([]*types.User, error) {
	cur, err := ms.coll.Find(ctx, bson.M{}) // TODO add limit
	if err != nil {
		return nil, err
	}

	var users []*types.User

	if err := cur.All(ctx, &users); err != nil {
		return nil, err
	}

	return users, nil
}

func (ms *MongoUserStore) InsertUser(ctx context.Context, user *types.User) (*types.User, error) {
	res, err := ms.coll.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}

	user.ID = res.InsertedID.(primitive.ObjectID) // casting
	return user, nil
}

func (ms *MongoUserStore) DeleteUser(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = ms.coll.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return err
	}

	log.Printf("user deleted success id: %s deleted", id)

	return nil
}

func (ms *MongoUserStore) PutUser(ctx context.Context, params *types.UpdateUserParams, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	result, err := ms.coll.UpdateOne(ctx, bson.M{"_id": oid}, bson.D{{
		Key: "$set", Value: params.ToBsonMap(),
	}})
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("no user found with id %s", id)
	}

	return nil
}

func (ms *MongoUserStore) Drop(ctx context.Context) error {
	log.Printf("dropping %s collection", userCollection)
	return ms.coll.Drop(ctx)
}
