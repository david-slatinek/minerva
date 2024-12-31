package main

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"log"
	"time"
)

type Docker struct {
	apiClient *client.Client
	db        *MeasurementsTable
}

func NewDocker(db *MeasurementsTable) (*Docker, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}

	return &Docker{
		apiClient: cli,
		db:        db,
	}, nil
}

func (d *Docker) getContainerState() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	containers, err := d.apiClient.ContainerList(ctx, container.ListOptions{
		Filters: filters.NewArgs(filters.Arg("name", containerName)),
	})

	if err != nil {
		return "", err
	}

	if len(containers) == 0 {
		return "", fmt.Errorf("container %s not found", containerName)
	}

	return containers[0].State, nil
}

func (d *Docker) Produce(c chan bool) {
	timestamp := time.Now()

	for {
		select {
		case <-c:
			return
		default:
			timestamp = d.write(timestamp)
		}
	}
}

func (d *Docker) write(timestamp time.Time) time.Time {
	if time.Now().Sub(timestamp).Seconds() < timeout {
		return timestamp
	}

	state, err := d.getContainerState()
	if err != nil {
		log.Printf("error getting container state: %s", err)
		state = "error"
	}

	_, err = d.db.CreateMeasurement(Measurement{
		Date:   time.Now().Format("2006-01-02"),
		Time:   time.Now().Format("15:04:05"),
		Status: state,
		Mode:   1,
	})
	if err != nil {
		log.Printf("error creating measurement: %s", err)
	}

	return time.Now()
}

func (d *Docker) Close() error {
	return d.apiClient.Close()
}
