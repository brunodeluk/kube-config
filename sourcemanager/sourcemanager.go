package sourcemanager

import (
	"context"
	"github.com/brunodeluk/kube-config/sourcemanager/source"
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

func (sm *SourceManager) FetchAll(ctx context.Context) []error {
	var wg sync.WaitGroup
	wg.Add(len(sm.sources))

	errors := make([]error, len(sm.sources))

	for _, src := range sm.sources {
		go func(src source.Source) {
			defer wg.Done()
			err := src.Fetch(ctx)
			if err != nil {
				errors = append(errors, err)
			}
		}(src)
	}

	wg.Wait()

	return errors
}
