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

	Aggregate(context.Context,
		time.Time,
		time.Time,
		chan string)
}
