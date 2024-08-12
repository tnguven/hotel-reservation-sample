package types

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Booking struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	RoomID      primitive.ObjectID `bson:"roomID,omitempty" json:"roomID,omitempty"`
	UserID      primitive.ObjectID `bson:"userID,omitempty" json:"userID,omitempty"`
	CountPerson int                `bson:"countPerson,omitempty" json:"countPerson,omitempty"`
	FromDate    time.Time          `bson:"from,omitempty" json:"fromDate,omitempty"`
	TillDate    time.Time          `bson:"tillDate,omitempty" json:"tillDate,omitempty"`
}

type BookingParam struct {
	RoomID      string             `json:"roomID,omitempty"`
	UserID      primitive.ObjectID `json:"userID,omitempty"`
	CountPerson int                `json:"countPerson,omitempty"`
	FromDate    time.Time          `json:"fromDate,omitempty"`
	TillDate    time.Time          `json:"tillDate,omitempty"`
}

func NewBookingFromParams(params *BookingParam) (*Booking, error) {
	roomOID, err := primitive.ObjectIDFromHex(params.RoomID)
	if err != nil {
		return nil, err
	}

	return &Booking{
		RoomID:      roomOID,
		UserID:      params.UserID,
		CountPerson: params.CountPerson,
		FromDate:    params.FromDate,
		TillDate:    params.TillDate,
	}, nil
}
