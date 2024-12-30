package main

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"time"
)

type Docker struct {
	apiClient *client.Client
}

func NewDocker() (*Docker, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}
	return &Docker{apiClient: cli}, nil
}

func (d Docker) GetContainerId(name string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	containers, err := d.apiClient.ContainerList(ctx, container.ListOptions{
		Filters: filters.NewArgs(filters.Arg("name", name)),
	})

	if err != nil {
		return "", err
	}

	if len(containers) == 0 {
		return "", fmt.Errorf("container %s not found", name)
	}

	return containers[0].ID, nil
}
