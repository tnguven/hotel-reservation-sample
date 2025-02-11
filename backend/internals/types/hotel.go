package types

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Hotel struct {
	ID       primitive.ObjectID   `bson:"_id,omitempty" json:"id,omitempty"`
	Name     string               `bson:"name" json:"name"`
	Location string               `bson:"location" json:"location"`
	Rooms    []primitive.ObjectID `bson:"rooms" json:"rooms"`
	Rating   int                  `bson:"rating" json:"rating"`
}

type UpdateHotelParams struct {
	Name     string             `json:"name,omitempty"`
	Location string             `json:"location,omitempty"`
	RoomID   primitive.ObjectID `json:"roomId,omitempty"`
}

func (p *UpdateHotelParams) ToBsonMap() bson.M {
	update := bson.M{}
	setValues := bson.M{}

	if len(p.Name) > 0 {
		setValues["name"] = p.Name
	}
	if len(p.Location) > 0 {
		setValues["location"] = p.Location
	}
	if len(setValues) > 0 {
		update["$set"] = setValues
	}
	if len(p.RoomID) > 0 {
		update["$push"] = bson.M{"rooms": p.RoomID}
	}

	return update
}

type HotelQueryParams struct {
	Rooms  bool `query:"room"`
	Rating int  `query:"rating"`

	*PaginationQuery
}

func NewHotelQueryParam(rooms bool, rating int, page int64, limit int64) HotelQueryParams {
	return HotelQueryParams{
		Rooms:  false,
		Rating: 0,
		PaginationQuery: &PaginationQuery{
			Limit: page,
			Page:  (page - 1) * limit,
		},
	}
}
