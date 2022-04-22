package sourcemanager

import (
	"context"
	"github.com/brunodeluk/kube-config/internal/sourcemanager/source"
	"sync"
)

type SourceManager struct {
	sources []source.Source
}

func New() *SourceManager {
	return &SourceManager{}
}

func (sm *SourceManager) Add(src source.Source) {
	sm.sources = append(sm.sources, src)
}

func (sm *SourceManager) FetchAll(ctx context.Context) ([]string, []error) {
	var wg sync.WaitGroup
	wg.Add(len(sm.sources))

	errors := make([]error, 0)
	paths := make([]string, 0)

	for _, src := range sm.sources {
		go func(src source.Source) {
			defer wg.Done()
			path, err := src.Fetch(ctx)
			if err != nil {
				errors = append(errors, err)
			}
			paths = append(paths, path)
		}(src)
	}

	wg.Wait()

	return paths, errors
}
