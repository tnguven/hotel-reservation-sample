package configure

type (
	DbConfig interface {
		DbName() string
		DbUri() string
		DbUriWithDbName() string
	}

	Secrets interface {
		JWTSecret() string
	}

	Session interface {
		TokenExpHour() int64
	}

	Server interface {
		ListenAddr() string
	}

	Common interface {
		GoEnv() string
		WithLog() bool
	}
)
