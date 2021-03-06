package main

import (
	"context"
	"fmt"
	"github.com/brunodeluk/kube-config/internal/configmanager/client"
	"github.com/brunodeluk/kube-config/internal/sourcemanager"
	"github.com/brunodeluk/kube-config/internal/sourcemanager/source"
	"os"
)

func main() {
	fmt.Printf("[main][INFO] executing kube-config\n")
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
			fmt.Printf("[main][sourcemanager][ERROR] %v\n", err)
		}
		os.Exit(1)
	}

	for _, path := range paths {
		k8s := client.Kubernetes{}
		err := k8s.Apply(context.Background(), path)
		if err != nil {
			fmt.Printf("[main][kube-client][ERROR] %v\n", err)
			os.Exit(1)
		}
	}
}
