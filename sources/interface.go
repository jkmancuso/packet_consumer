package sources

import "context"

type Consumer interface {
	Start(context.Context) error
	GetRecord(context.Context) (string, error)
	Teardown()
}
