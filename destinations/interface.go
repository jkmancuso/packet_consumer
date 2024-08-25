package destinations

import "context"

type Destination interface {
	SendRecord(context.Context, string) error
}
