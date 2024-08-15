package config

import (
	"errors"
	"log"
)

const (
	dbUri = "mongodb://localhost:27017"
)

type Configs struct {
	DbURI        string
	DbName       string
	DbUserName   string
	DbPassword   string
	CreateIndex  bool
	JWTSecret    string
	TokenExpHour int64
}

func (conf *Configs) WithDbURI(dbURI string) *Configs {
	conf.DbURI = dbURI
	return conf
}

func (conf *Configs) WithDbName(dbName string) *Configs {
	conf.DbName = dbName
	return conf
}

func (conf *Configs) WithDbUserName(username string) *Configs {
	conf.DbUserName = username
	return conf
}

func (conf *Configs) WithDbPassword(password string) *Configs {
	conf.DbPassword = password
	return conf
}

func (conf *Configs) WithDbCreateIndex(withIndex bool) *Configs {
	conf.CreateIndex = withIndex
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

func (conf *Configs) Validate() *Configs {
	if conf.DbName == "" {
		log.Fatal(errors.New("missing database name"))
	}

	if conf.DbURI == "" {
		log.Fatal(errors.New("missing database URI"))
	}
	return conf
}

func New() *Configs {
	return &Configs{
		DbURI:        dbUri,
		CreateIndex:  false,
		JWTSecret:    "top_secret",
		TokenExpHour: 72,
	}
}
