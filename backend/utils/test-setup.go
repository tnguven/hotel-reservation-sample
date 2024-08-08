package utils

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/tnguven/hotel-reservation-app/config"
	"github.com/tnguven/hotel-reservation-app/db"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewDb() *mongo.Database {
	ctx := context.Background()
	return db.New(ctx, config.New().
		WithDbUserName("admin").
		WithDbPassword("secret").
		WithDbName("test"))
}

func NewRequestWithHeader(method, target string, body io.Reader) *http.Request {
	req := httptest.NewRequest(method, target, body)
	req.Header.Add("Content-Type", "application/json")
	return req
}
