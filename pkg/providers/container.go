package providers

import (
	"time"

	hclog "github.com/hashicorp/go-hclog"
	"github.com/shipyard-run/shipyard/pkg/clients"
	"github.com/shipyard-run/shipyard/pkg/config"
)

// Container is a provider for creating and destroying Docker containers
type Container struct {
	config     *config.Container
	client     clients.ContainerTasks
	httpClient clients.HTTP
	log        hclog.Logger
}

// NewContainer creates a new container with the given config and Docker client
func NewContainer(co *config.Container, cl clients.ContainerTasks, hc clients.HTTP, l hclog.Logger) *Container {
	return &Container{co, cl, hc, l}
}

// Create implements provider method and creates a Docker container with the given config
func (c *Container) Create() error {
	c.log.Info("Creating Container", "ref", c.config.Name)

	// pull any images needed for this container
	err := c.client.PullImage(c.config.Image, false)
	if err != nil {
		c.log.Error("Error pulling container image", "ref", c.config.Name, "image", c.config.Image.Name)

		return err
	}

	_, err = c.client.CreateContainer(c.config)

	if c.config.HealthCheck == nil {
		return err
	}

	// check the health of the container
	if hc := c.config.HealthCheck.HTTP; hc != "" {
		d, err := time.ParseDuration(c.config.HealthCheck.Timeout)
		if err != nil {
			return err
		}

		return c.httpClient.HealthCheckHTTP(hc, d)
	}

	return nil
}

// Destroy stops and removes the container
func (c *Container) Destroy() error {
	c.log.Info("Destroy Container", "ref", c.config.Name)
	ids, err := c.client.FindContainerIDs(c.config.Name, c.config.Type)

	if err != nil {
		return err
	}

	if len(ids) > 0 {
		for _, id := range ids {
			for _, n := range c.config.Networks {
				err := c.client.DetachNetwork(n.Name, id)
				if err != nil {
					c.log.Error("Unable to detach network", "ref", c.config.Name, "network", n.Name)
				}
			}

			err := c.client.RemoveContainer(id)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Lookup the ID based on the config
func (c *Container) Lookup() ([]string, error) {
	return c.client.FindContainerIDs(c.config.Name, c.config.Type)
}
