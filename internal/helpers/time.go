package helpers

import (
	"time"
)

func NowISO3601() string {
	return time.Now().UTC().Format(time.RFC3339)
}
