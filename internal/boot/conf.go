package boot

import (
	"github.com/khafidprayoga/parking-app/internal/types"
	"time"
)

var AppConfig = types.AppConfig{
	ConnLifetime: 10 * time.Second,
	AppVersion:   "v0.0.1",
}
