package source

import "context"

type Source interface {
	Fetch(ctx context.Context) (string, error)
}
