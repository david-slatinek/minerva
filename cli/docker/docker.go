package docker

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"log"
	"main/database"
	"main/random"
	"time"
)

const (
	writeTimeout  = 5
	containerName = "minerva-api"
	minMinutes    = 10
	maxMinutes    = 15
)

type Docker struct {
	apiClient *client.Client
	db        *database.MeasurementsTable
	mode      int
}

func NewDocker(db *database.MeasurementsTable, mode int) (*Docker, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}

	return &Docker{
		apiClient: cli,
		db:        db,
		mode:      mode,
	}, nil
}

func (d Docker) getContainerInfo() (string, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	containers, err := d.apiClient.ContainerList(ctx, container.ListOptions{
		Filters: filters.NewArgs(filters.Arg("name", containerName)),
	})

	if err != nil {
		return "", "", err
	}

	if len(containers) == 0 {
		return "", "", fmt.Errorf("container %s not found", containerName)
	}

	return containers[0].ID, containers[0].State, nil
}

func (d Docker) Produce(c chan bool) {
	timestamp := time.Now()

	for {
		select {
		case <-c:
			log.Println("exit produce")
			return
		default:
			timestamp = d.write(timestamp)
		}
	}
}

func (d Docker) write(timestamp time.Time) time.Time {
	if time.Now().Sub(timestamp).Seconds() < writeTimeout {
		return timestamp
	}

	_, state, err := d.getContainerInfo()
	if err != nil {
		log.Printf("error getting container state: %s", err)
		state = "error"
	}

	_, err = d.db.CreateMeasurement(database.Measurement{
		Date:   time.Now().Format("2006-01-02"),
		Time:   time.Now().Format("15:04:05"),
		Status: state,
		Mode:   d.mode,
	})
	if err != nil {
		log.Printf("error creating measurement: %s", err)
	}

	return time.Now()
}

func (d Docker) Close() error {
	return d.apiClient.Close()
}

func (d Docker) Stop(c chan bool) {
	timestamp := time.Now()
	stopMinutes := random.Int(minMinutes, maxMinutes)

	for {
		select {
		case <-c:
			log.Println("exit stop")
			return
		default:
			if (int)(time.Now().Sub(timestamp).Minutes()) >= stopMinutes {
				timestamp = time.Now()
				stopMinutes = random.Int(minMinutes, maxMinutes)

				if err := d.stop(); err == nil {
					log.Println("stop successful")
				} else {
					log.Printf("error with stop: %s", err)
				}
			}
		}
	}
}

func (d Docker) stop() error {
	id, _, err := d.getContainerInfo()

	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 12*time.Second)
	defer cancel()

	err = d.apiClient.ContainerStop(ctx, id, container.StopOptions{})
	if err != nil {
		return fmt.Errorf("error stopping container: %s", err)
	}

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = d.apiClient.ContainerRemove(ctx, id, container.RemoveOptions{})
	if err != nil {
		log.Printf("error removing container: %s", err)
		return fmt.Errorf("error removing container: %s", err)
	}

	return nil
}
