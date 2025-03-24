package types

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RoomType string

const (
	FamilyRoomType     RoomType = "family"
	FamilySuitRoomType RoomType = "family_suit"
	SuiteRoomType      RoomType = "suit"
	HoneyMoonRoomType  RoomType = "honey_moon"
	KingRoomType       RoomType = "king"
)

type Room struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Type      RoomType           `bson:"type" json:"type"`
	BasePrice float64            `bson:"basePrice" json:"basePrice"`
	Price     float64            `bson:"price" json:"price"`
	HotelID   primitive.ObjectID `bson:"hotelID" json:"hotelID"`
	Status    RoomStatus         `bson:"status" json:"status"`
}

type RoomStatus string

const (
	OccupiedRoom  RoomStatus = "occupied"
	BookedRoom    RoomStatus = "booked"
	AvailableRoom RoomStatus = "available"
)

type GetRoomsRequest struct {
	Status []RoomStatus

	MongoPaginateWithIDQuery
}
