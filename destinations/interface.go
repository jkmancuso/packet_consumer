package destinations

import "context"

type Destination interface {
	SendRecord(context.Context,
		string,
		map[string]string,
		map[string]interface{}) error
}
