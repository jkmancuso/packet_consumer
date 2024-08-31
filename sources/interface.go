package sources

import "context"

type Consumer interface {
	GetRecord(context.Context) ([]byte, error)
	Teardown()
}
