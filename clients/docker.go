package clients

import (
	"context"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/versions"
)

type Docker interface {
	ContainerCreate(
		ctx context.Context, 
		config *container.Config, 
		hostConfig *container.HostConfig, 
		networkingConfig *network.NetworkingConfig, 
		containerName string,
	) (container.ContainerCreateCreatedBody, error)
}