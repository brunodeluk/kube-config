package client

import "context"

type Client interface {
	Apply(ctx context.Context, path string) error
}
