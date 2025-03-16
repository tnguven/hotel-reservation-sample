package main

import (
	"cmp"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/tnguven/hotel-reservation-app/internals/configure"
)

type Configs struct {
	configure.Common
	configure.DbConfig

	mongoDbURI  string
	mongoDbName string
	log         bool
	env         string
}

func (conf *Configs) WithMongoDbURI(dbURI string) *Configs {
	conf.mongoDbURI = dbURI
	return conf
}

func (conf *Configs) WithDbName(dbName string) *Configs {
	conf.mongoDbName = dbName
	return conf
}

func (conf *Configs) WithEnv(env string) *Configs {
	conf.env = env
	return conf
}

func (conf *Configs) Validate() *Configs {
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
		mongoDbName: cmp.Or(os.Getenv("MONGO_DATABASE"), "hotel_io_dev"),
		mongoDbURI:  cmp.Or(os.Getenv("MONGO_URI"), "mongodb://localhost:27017"),
		env:         cmp.Or(os.Getenv("ENV"), "development"),
		log:         true,
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

func (conf *Configs) GoEnv() string {
	return conf.env
}

func (conf *Configs) WithLog() bool {
	return conf.log
}
