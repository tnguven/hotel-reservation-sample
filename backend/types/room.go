package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type RoomType string

const (
	FamilyRoomType     RoomType = "family"
	FamilySuitRoomType RoomType = "family suit"
	SuiteRoomType      RoomType = "suit"
	HoneyMoonRoomType  RoomType = "honey moon"
	KingRoomType       RoomType = "king"
)

type Room struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Type      RoomType           `bson:"type" json:"type"`
	BasePrice float64            `bson:"basePrice" json:"basePrice"`
	Price     float64            `bson:"price" json:"price"`
	HotelID   primitive.ObjectID `bson:"hotelID" json:"hotelID"`
}
