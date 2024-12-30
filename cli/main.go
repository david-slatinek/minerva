package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

const (
	connectionString = "host=localhost user=david password=david dbname=cli port=5000 sslmode=disable TimeZone=Europe/Ljubljana"
	containerName    = "minerva-api"
)

func configureLog() (func(), error) {
	f, err := os.OpenFile("measurements.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	wrt := io.MultiWriter(os.Stdout, f)
	log.SetOutput(wrt)

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	return func() { f.Close() }, err
}

func main() {
	closer, err := configureLog()
	if err != nil {
		log.Fatalf("error configuring log: %v", err)
	}
	defer closer()

	db, er := NewMeasurements(connectionString)
	if er != nil {
		log.Fatalf("error connecting to database: %s", er)
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Fatalf("error closing database connection: %s", err)
		}
	}()

	measurement, err := db.CreateMeasurement(Measurement{
		Date:   time.Now().Format("2006-01-02"),
		Time:   time.Now().Format("15:04:05"),
		Status: "running",
		Mode:   1,
	})
	if err != nil {
		log.Fatalf("error creating measurement: %s", err)
	}

	fmt.Println(measurement)

	dockerClient, err := NewDocker()
	if err != nil {
		log.Fatalf("error creating docker client: %s", err)
	}

	id, err := dockerClient.GetContainerId(containerName)
	if err != nil {
		log.Fatalf("error getting container id: %s", err)
	}

	fmt.Println(id)
}
