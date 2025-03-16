package main

import (
	"cmp"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/tnguven/hotel-reservation-app/internals/configure"
	"github.com/tnguven/hotel-reservation-app/internals/must"
)

type Configs struct {
	configure.Common
	configure.DbConfig
	configure.Server
	configure.Secrets
	configure.Session

	mongoDbURI   string
	mongoDbName  string
	jwtSecret    string
	tokenExpHour int64
	listenAddr   string
	log          bool
	env          string
}

func (conf *Configs) WithMongoDbURI(dbURI string) *Configs {
	conf.mongoDbURI = dbURI
	return conf
}

func (conf *Configs) WithDbName(dbName string) *Configs {
	conf.mongoDbName = dbName
	return conf
}

func (conf *Configs) WithJWTSecret(secret string) *Configs {
	conf.jwtSecret = secret
	return conf
}

func (conf *Configs) WithTokenExpirationHours(hour int64) *Configs {
	conf.tokenExpHour = hour
	return conf
}

func (conf *Configs) WithListenAddr(addr string) *Configs {
	conf.listenAddr = addr
	return conf
}

func (conf *Configs) WithEnv(env string) *Configs {
	conf.env = env
	return conf
}

func (conf *Configs) Validate() *Configs {
	if conf.listenAddr == "" {
		log.Fatal(errors.New("missing listen addr"))
	}

	if conf.mongoDbName == "" {
		log.Fatal(errors.New("missing mongo database name"))
	}

	if conf.mongoDbURI == "" {
		log.Fatal(errors.New("missing mongo database URI"))
	}

	return conf
}

func (conf *Configs) Debug() *Configs {
	fmt.Printf("ENVS: %+v", conf)
	return conf
}

func NewConfig() *Configs {
	return &Configs{
		mongoDbName:  cmp.Or(os.Getenv("MONGO_DATABASE"), "hotel_io_dev"),
		mongoDbURI:   cmp.Or(os.Getenv("MONGO_URI"), "mongodb://localhost:27017"),
		jwtSecret:    cmp.Or(os.Getenv("JWT_SECRET"), "top_secret"),
		tokenExpHour: must.Panic(strconv.ParseInt((cmp.Or(os.Getenv("EXPIRE_IN_HOURS"), "72")), 10, 64)),
		listenAddr:   fmt.Sprintf(":%s", cmp.Or(os.Getenv("LISTEN_ADDR"), "5000")),
		env:          cmp.Or(os.Getenv("ENV"), "development"),
		log:          true,
	}
}

func (conf *Configs) DbName() string {
	return conf.mongoDbName
}

func (conf *Configs) DbURI() string {
	return conf.mongoDbURI
}

func (conf *Configs) DbUriWithDbName() string {
	return fmt.Sprintf("%s/%s", conf.DbURI(), conf.DbName())
}

func (conf *Configs) JWTSecret() string {
	return conf.jwtSecret
}

func (conf *Configs) TokenExpHour() int64 {
	return conf.tokenExpHour
}

func (conf *Configs) ListenAddr() string {
	return conf.listenAddr
}

func (conf *Configs) GoEnv() string {
	return conf.env
}

func (conf *Configs) WithLog() bool {
	return conf.log
}
