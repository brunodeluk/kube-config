package main

import (
	"context"
	"fmt"
	"github.com/brunodeluk/kube-config/configmanager/cluster"
	"github.com/brunodeluk/kube-config/sourcemanager"
	"github.com/brunodeluk/kube-config/sourcemanager/source"
	"k8s.io/utils/env"
)

func main() {
	sm := sourcemanager.New()

	sm.Add(&source.Git{
		Token:  env.GetString("git_token", ""),
		URL:    "https://github.com/kubernetes/examples.git",
		Branch: "master",
		Dir:    "guestbook/all-in-one",
	})

	paths, errors := sm.FetchAll(context.Background())
	if len(errors) > 0 {
		for _, err := range errors {
			fmt.Println(err)
		}
	}

	for _, path := range paths {
		k8s := cluster.Kubernetes{}
		err := k8s.Apply(context.Background(), path)
		if err != nil {
			fmt.Println(err)
		}
	}
}
