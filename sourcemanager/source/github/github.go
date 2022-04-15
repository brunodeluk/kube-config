package github

import (
	"context"
	"fmt"
)

type GitHub struct {
}

func (gh *GitHub) Fetch(ctx context.Context) error {
	fmt.Println("[source] cloning GitHub project ....")
	return nil
}
