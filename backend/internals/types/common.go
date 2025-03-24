package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type (
	PaginateWithID[T any] interface {
		GetLastID() (T, error)
	}

	QueryNumericPaginate struct {
		Limit int64 `query:"limit" validate:"numeric,max=100,omitempty"`
		Page  int   `query:"page" validate:"numeric,omitempty"`
		Skip  int64 `query:"skip" validate:"numeric,omitempty"`
	}

	QueryCursorPaginate[T any] struct {
		LastID string `query:"limit" validate:"string"`
		Limit  int64  `query:"limit" validate:"numeric,max=100,omitempty"`
		PaginateWithID[T]
	}
)

func NewQueryNumericPaginate(limit int, page int) QueryNumericPaginate {
	return QueryNumericPaginate{
		Limit: int64(limit),
		Page:  page,
		Skip:  int64((page - 1) * limit),
	}
}

func NewQueryCursorPaginate(lastID string, limit int) QueryCursorPaginate[primitive.ObjectID] {
	return QueryCursorPaginate[primitive.ObjectID]{
		LastID: lastID,
		Limit:  int64(limit),
	}
}

func (p QueryCursorPaginate[T]) GetLastID() (primitive.ObjectID, error) {
	return primitive.ObjectIDFromHex(p.LastID)
}
