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

type TestRequest struct {
	Method  string
	Target  string
	Payload io.Reader
	Token   string
}

func (t *TestRequest) NewRequestWithHeader() *http.Request {
	request := httptest.NewRequest(t.Method, t.Target, t.Payload)
	request.Header.Add("Content-Type", "application/json")
	if t.Token != "" {
		request.Header.Add("X-api-token", t.Token)
	}

	return request
}
