package configmanager

import (
	"context"
	"github.com/brunodeluk/kube-config/configmanager/cluster"
)

type ConfigManager struct {
	Cluster cluster.Cluster
}

func (cm *ConfigManager) Apply(ctx context.Context, path string) error {
	return cm.Cluster.Apply(ctx, path)
}
