package cluster

import "context"

type Cluster interface {
	Apply(ctx context.Context, path string) error
}
