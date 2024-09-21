package destinations

import (
	"context"
	"time"
)

type Destination interface {
	SendRecord(context.Context,
		string,
		map[string]string,
		map[string]interface{},
		time.Time) error
}
