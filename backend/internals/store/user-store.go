package store

import (
	"context"
	"log"

	"github.com/tnguven/hotel-reservation-app/internals/repo"
	"github.com/tnguven/hotel-reservation-app/internals/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const userCollection = "users"

type UserStore interface {
	Dropper

	GetByID(context.Context, string) (*types.User, error)
	GetUserByEmail(context.Context, string) (*types.User, error)
	GetUsers(context.Context, *types.QueryNumericPaginate) ([]*types.User, int64, error)
	InsertUser(context.Context, *types.User) (*types.User, error)
	DeleteUser(context.Context, string) error
	PutUser(context.Context, *types.UpdateUserParams, string) (int64, error)
}

type MongoUserStore struct {
	db   *mongo.Database
	coll *mongo.Collection
}

func NewMongoUserStore(mongodb *repo.MongoDatabase) *MongoUserStore {
	return &MongoUserStore{
		db:   mongodb.GetDb(),
		coll: mongodb.Coll(userCollection),
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

func (ms *MongoUserStore) GetUsers(
	ctx context.Context,
	pagination *types.QueryNumericPaginate,
) ([]*types.User, int64, error) {
	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: bson.D{}}},
		bson.D{{Key: "$facet", Value: bson.D{
			{Key: "data", Value: bson.A{
				bson.D{{Key: "$skip", Value: pagination.Page}},
				bson.D{{Key: "$limit", Value: pagination.Limit}},
			}},
			{Key: "totalCount", Value: bson.A{
				bson.D{{Key: "$count", Value: "count"}},
			}},
		}}},
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "data", Value: 1},
			{Key: "totalCount", Value: bson.D{
				{Key: "$arrayElemAt", Value: bson.A{"$totalCount.count", 0}},
			}},
		}}},
	}
	cur, err := ms.coll.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, 0, err
	}

	var aggResult []struct {
		Data       []*types.User `bson:"data"`
		TotalCount int64         `bson:"totalCount"`
	}
	if err := cur.All(ctx, &aggResult); err != nil {
		return nil, 0, err
	}

	if len(aggResult) == 0 {
		return []*types.User{}, 0, nil
	}

	return aggResult[0].Data, aggResult[0].TotalCount, nil
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

func (ms *MongoUserStore) PutUser(
	ctx context.Context,
	params *types.UpdateUserParams,
	id string,
) (int64, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return 0, err
	}

	result, err := ms.coll.UpdateOne(ctx, bson.M{"_id": oid}, bson.D{{
		Key: "$set", Value: params.ToBsonMap(),
	}})
	if err != nil {
		return 0, err
	}

	return result.MatchedCount, nil
}

func (ms *MongoUserStore) Drop(ctx context.Context) error {
	log.Printf("dropping %s collection", userCollection)
	return ms.coll.Drop(ctx)
}
