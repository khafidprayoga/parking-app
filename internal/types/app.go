package types

import "time"

type AppConfig struct {
	ConnLifetime time.Duration
	AppVersion   string
}
