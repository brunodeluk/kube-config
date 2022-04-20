package configmanager

import (
	"context"
	"github.com/brunodeluk/kube-config/configmanager/client"
)

type ConfigManager struct {
	Client client.Client
}

func (cm *ConfigManager) Apply(ctx context.Context, path string) error {
	return cm.Client.Apply(ctx, path)
}
