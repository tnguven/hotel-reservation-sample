package utils

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/tnguven/hotel-reservation-app/internals/config"
	"github.com/tnguven/hotel-reservation-app/internals/repo"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewDb(config *config.Configs) (*mongo.Client, *mongo.Database) {
	ctx := context.Background()
	return repo.NewMongoClient(ctx, config)
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
		request.Header.Add("X-Api-Token", t.Token)
	}

	return request
}
