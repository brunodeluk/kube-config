package main

import (
	"context"
	"fmt"
	"github.com/brunodeluk/kube-config/configmanager/client"
	"github.com/brunodeluk/kube-config/sourcemanager"
	"github.com/brunodeluk/kube-config/sourcemanager/source"
	"os"
)

func main() {
	sm := sourcemanager.New()

	sm.Add(&source.Git{
		Token:  os.Getenv("git_token"),
		URL:    os.Getenv("url"),
		Branch: os.Getenv("branch"),
		Dir:    os.Getenv("dir"),
	})

	paths, errors := sm.FetchAll(context.Background())
	if len(errors) > 0 {
		for _, err := range errors {
			fmt.Println(err)
		}
	}

	for _, path := range paths {
		k8s := client.Kubernetes{}
		err := k8s.Apply(context.Background(), path)
		if err != nil {
			fmt.Println(err)
		}
	}
}
