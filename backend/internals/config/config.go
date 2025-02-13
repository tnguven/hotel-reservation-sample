package config

import (
	"cmp"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/tnguven/hotel-reservation-app/internals/must"
)

type Configs struct {
	MongoDbURI   string
	MongoDbName  string
	JWTSecret    string
	TokenExpHour int64
	ListenAddr   string
	Log          bool
	Env          string
}

func (conf *Configs) WithMongoDbURI(dbURI string) *Configs {
	conf.MongoDbURI = dbURI
	return conf
}

func (conf *Configs) WithDbName(dbName string) *Configs {
	conf.MongoDbName = dbName
	return conf
}

func (conf *Configs) WithJWTSecret(secret string) *Configs {
	conf.JWTSecret = secret
	return conf
}

func (conf *Configs) WithTokenExpirationHours(hour int64) *Configs {
	conf.TokenExpHour = hour
	return conf
}

func (conf *Configs) WithListenAddr(addr string) *Configs {
	conf.ListenAddr = addr
	return conf
}

func (conf *Configs) WithEnv(env string) *Configs {
	conf.Env = env
	return conf
}

func (conf *Configs) Validate() *Configs {
	if conf.ListenAddr == "" {
		log.Fatal(errors.New("missing listen addr"))
	}

	if conf.MongoDbName == "" {
		log.Fatal(errors.New("missing mongo database name"))
	}

	if conf.MongoDbURI == "" {
		log.Fatal(errors.New("missing mongo database URI"))
	}

	return conf
}

func (conf *Configs) Debug() *Configs {
	fmt.Printf("ENVS: %+v", conf)
	return conf
}

func New() *Configs {
	return &Configs{
		MongoDbName:  cmp.Or(os.Getenv("MONGO_DATABASE"), "hotel_io_dev"),
		MongoDbURI:   cmp.Or(os.Getenv("MONGO_URI"), "mongodb://localhost:27017"),
		JWTSecret:    cmp.Or(os.Getenv("JWT_SECRET"), "top_secret"),
		TokenExpHour: must.Panic(strconv.ParseInt((cmp.Or(os.Getenv("EXPIRE_IN_HOURS"), "72")), 10, 64)),
		ListenAddr:   fmt.Sprintf(":%s", cmp.Or(os.Getenv("LISTEN_ADDR"), "5000")),
		Env:          cmp.Or(os.Getenv("ENV"), "development"),
		Log:          true,
	}
}
