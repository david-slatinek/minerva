package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	connectionString = "host=localhost user=david password=david dbname=cli port=5000 sslmode=disable TimeZone=Europe/Ljubljana"
	containerName    = "minerva-api"
	timeout          = 1
)

func configureLog() (func(), error) {
	f, err := os.OpenFile("measurements.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %s", err)
	}
	wrt := io.MultiWriter(os.Stdout, f)
	log.SetOutput(wrt)

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	return func() {
		if err := f.Close(); err != nil {
			log.Printf("error closing file: %s", err)
		}
	}, err
}

func main() {
	closer, err := configureLog()
	if err != nil {
		log.Fatalf("error configuring log: %s", err)
	}
	defer closer()

	modePtr := flag.Int("mode", 0, "Monitoring mode: 1, 2, 3")
	flag.Parse()

	if *modePtr <= 0 || *modePtr > 3 {
		log.Fatalf("invalid monitoring mode: %d", *modePtr)
	}

	db, er := NewMeasurements()
	if er != nil {
		log.Fatalf("error connecting to database: %s", er)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("error closing database connection: %s", err)
		}
	}()

	dockerClient, err := NewDocker(db, *modePtr)
	if err != nil {
		log.Fatalf("error creating docker client: %s", err)
	}
	defer func() {
		if err := dockerClient.Close(); err != nil {
			log.Printf("error closing docker client: %s", err)
		}
	}()

	endC := make(chan bool, 1)
	go dockerClient.Produce(endC)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)
	<-c

	endC <- true
	ctx2, cancel2 := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel2()

	select {
	case <-ctx2.Done():
		break
	}

	fmt.Println("shutting down")
}
