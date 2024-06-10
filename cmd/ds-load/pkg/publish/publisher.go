package publish

import (
	"context"
	"io"
)

type Publisher interface {
	Publish(ctx context.Context, reader io.Reader) error
}
