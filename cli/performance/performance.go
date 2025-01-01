package performance

import (
	"fmt"
	"log"
	"main/random"
	"os/exec"
	"time"
)

const (
	minMinutes = 1
	maxMinutes = 2
)

type Testing struct {
}

func (t Testing) Start(c chan bool) {
	timestamp := time.Now()
	stopMinutes := random.Int(minMinutes, maxMinutes)

	for {
		select {
		case <-c:
			log.Println("exit start")
			return
		default:
			if (int)(time.Now().Sub(timestamp).Minutes()) >= stopMinutes {
				timestamp = time.Now()
				stopMinutes = random.Int(minMinutes, maxMinutes)
				log.Println("running k6")
				if err := t.run(); err != nil {
					log.Printf("error running tests: %s\n", err)
				}
			}
		}
	}
}

func (t Testing) run() error {
	cmd := exec.Command("k6", "run", "stress-test.js")

	if err := cmd.Start(); err != nil {
		return err
	}

	go func() {
		if err := cmd.Wait(); err != nil {
			fmt.Printf("k6 finished with error: %s\n", err)
		}
	}()

	return nil
}
