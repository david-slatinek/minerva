package main

import (
	"context"
	"flag"
	"io"
	"log"
	"main/database"
	"main/docker"
	"main/performance"
	"os"
	"os/signal"
	"syscall"
	"time"
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

	db, er := database.NewMeasurements()
	if er != nil {
		log.Fatalf("error connecting to database: %s", er)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("error closing database connection: %s", err)
		}
	}()

	dockerClient, err := docker.NewDocker(db, *modePtr)
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
	go dockerClient.Stop(endC, docker.ApiContainer)
	go dockerClient.Stop(endC, docker.DbContainer)
	go new(performance.Testing).Start(endC)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)
	<-c

	for i := 0; i < 4; i++ {
		endC <- true
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	select {
	case <-ctx.Done():
		break
	}

	log.Println("shutting down")
}
