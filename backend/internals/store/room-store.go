package store

import (
	"context"
	"log"
	"time"

	"github.com/tnguven/hotel-reservation-app/internals/repo"
	"github.com/tnguven/hotel-reservation-app/internals/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const roomCollection = "rooms"

type RoomStore interface {
	Dropper

	InsertRoom(context.Context, *types.Room) (*types.Room, error)
	GetRoomsByHotelID(context.Context, string) ([]*types.Room, error)
	GetRooms(context.Context, *types.GetRoomsRequest) ([]*types.Room, int64, string, error)
}

type MongoRoomStore struct {
	db   *mongo.Database
	coll *mongo.Collection

	HotelStore
}

func NewMongoRoomStore(db *repo.MongoDatabase, hotelStore HotelStore) *MongoRoomStore {
	return &MongoRoomStore{
		db:   db.GetDb(),
		coll: db.Coll(roomCollection),

		HotelStore: hotelStore,
	}
}

func (ms *MongoRoomStore) InsertRoom(ctx context.Context, room *types.Room) (*types.Room, error) {
	resp, err := ms.coll.InsertOne(ctx, room)
	if err != nil {
		return nil, err
	}

	room.ID = resp.InsertedID.(primitive.ObjectID)

	if err := ms.HotelStore.PutHotel(ctx, &types.UpdateHotelParams{RoomID: room.ID}, &room.HotelID); err != nil {
		return nil, err
	}

	return room, nil
}

func (ms *MongoRoomStore) GetRoomsByHotelID(ctx context.Context, hotelID string) ([]*types.Room, error) {
	oid, err := primitive.ObjectIDFromHex(hotelID)
	if err != nil {
		return nil, err
	}

	resp, err := ms.coll.Find(ctx, bson.M{"hotelID": oid})
	if err != nil {
		return nil, err
	}

	var rooms []*types.Room
	if err := resp.All(ctx, &rooms); err != nil {
		return nil, err
	}
	return rooms, nil
}

func (ms *MongoRoomStore) GetRooms(ctx context.Context, qParams *types.GetRoomsRequest) ([]*types.Room, int64, string, error) {
	now := time.Now()
	pipeline := mongo.Pipeline{}
	if qParams.LastID != "" {
		lastObjID, err := qParams.GetLastID()
		if err != nil {
			return nil, 0, "", err
		}
		pipeline = append(pipeline, bson.D{
			{Key: "$match", Value: bson.D{
				{Key: "_id", Value: bson.D{{Key: "$gt", Value: lastObjID}}},
			}},
		})
	}

	// Join bookings with rooms.
	pipeline = append(pipeline, bson.D{
		{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "bookings"},
			{Key: "localField", Value: "_id"},
			{Key: "foreignField", Value: "roomID"},
			{Key: "as", Value: "bookings"},
		}},
	})

	// Add computed "status" field.
	pipeline = append(pipeline, bson.D{
		{Key: "$addFields", Value: bson.D{
			{Key: "status", Value: bson.D{
				{Key: "$cond", Value: bson.A{
					// Check active bookings (fromDate <= now <= tillDate)
					bson.D{{Key: "$gt", Value: bson.A{
						bson.D{{Key: "$size", Value: bson.D{
							{Key: "$filter", Value: bson.D{
								{Key: "input", Value: "$bookings"},
								{Key: "as", Value: "b"},
								{Key: "cond", Value: bson.D{
									{Key: "$and", Value: bson.A{
										bson.D{{Key: "$lte", Value: bson.A{"$$b.fromDate", now}}},
										bson.D{{Key: "$gte", Value: bson.A{"$$b.tillDate", now}}},
									}},
								}},
							}},
						}}},
						0,
					}}},
					"occupied",
					// Else, check for any future booking (fromDate > now)
					bson.D{{Key: "$cond", Value: bson.A{
						bson.D{{Key: "$gt", Value: bson.A{
							bson.D{{Key: "$size", Value: bson.D{
								{Key: "$filter", Value: bson.D{
									{Key: "input", Value: "$bookings"},
									{Key: "as", Value: "b"},
									{Key: "cond", Value: bson.D{
										{Key: "$gt", Value: bson.A{"$$b.fromDate", now}},
									}},
								}},
							}}},
							0,
						}}},
						"booked",
						"available",
					}}},
				}},
			}},
		}},
	})

	// Optionally add a status filter.
	if len(qParams.Status) > 0 {
		pipeline = append(pipeline, bson.D{
			{Key: "$match", Value: bson.D{
				{Key: "status", Value: bson.D{{Key: "$in", Value: qParams.Status}}},
			}},
		})
	}

	// Use a $facet to get both the paginated data and total count.
	pipeline = append(pipeline, bson.D{
		{Key: "$facet", Value: bson.D{
			// Data facet: simply limit the page size.
			{Key: "data", Value: bson.A{
				bson.D{{Key: "$limit", Value: qParams.Limit}},
			}},
			// Total count facet.
			{Key: "totalCount", Value: bson.A{
				bson.D{{Key: "$count", Value: "count"}},
			}},
		}},
	})

	cur, err := ms.coll.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, 0, "", err
	}

	var aggResult []struct {
		Data       []*types.Room `bson:"data"`
		TotalCount []struct {
			Count int64 `bson:"count"`
		} `bson:"totalCount"`
	}

	if err := cur.All(ctx, &aggResult); err != nil {
		return nil, 0, "", err
	}

	if len(aggResult) == 0 {
		return []*types.Room{}, 0, "", nil
	}

	total := int64(0)
	if len(aggResult[0].TotalCount) > 0 {
		total = aggResult[0].TotalCount[0].Count
	}

	rooms := aggResult[0].Data
	var newLastID string
	// Get the _id of the last document in the paginated data (for the next page cursor).
	if len(rooms) > 0 {
		// Assuming the Room struct has an ID field of type primitive.ObjectID.
		newLastID = rooms[len(rooms)-1].ID.Hex()
	}

	return rooms, total, newLastID, nil
}

func (ms *MongoRoomStore) Drop(ctx context.Context) error {
	log.Printf("dropping %s collection", roomCollection)
	return ms.coll.Drop(ctx)
}
