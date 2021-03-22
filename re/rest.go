package re

import (
	"errors"
	"kubeedge-banch/util"
	"log"
	"strings"
	"time"
)

type resultChan struct {
	failedChan  chan int
	doneChan    chan bool
	failedCount int
}

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func Run(rest, restMethod, msg, msgBinary string, num, concurrency int) error {
	if msg != "" && msgBinary != "" {
		log.Println("Warning: msg-binary will be prioritized.")
	}
	var msgToSend []byte
	if msg != "" {
		msgToSend = []byte(msg)
	}
	if msgBinary != "" {
		fileMsg, err := util.ReadFile(msgBinary)
		if err != nil && msg == "" {
			return err
		}
		if err == nil {
			msgToSend = fileMsg
		}
	}

	var hanler util.HttpHandler
	switch strings.ToUpper(restMethod) {
	case "POST":
		hanler = util.Post
	default:
		return errors.New("unsupport method: " + restMethod)
	}

	pool := util.NewPool(concurrency)
	startTime := time.Now()
	log.Println("Sending start:", startTime.UTC())

	resChan := &resultChan{
		failedChan:  make(chan int, concurrency),
		doneChan:    make(chan bool),
		failedCount: 0,
	}

	go countFailed(resChan)
	for i := 0; i < num; i++ {
		pool.Add(1)
		go func() {
			defer pool.Done()
			_, err := hanler(rest, msgToSend)
			if err != nil {
				log.Panic("Request failed")
				resChan.failedChan <- 1
			}
		}()
	}
	pool.Wait()
	endTime := time.Now()
	resChan.doneChan <- true
	log.Printf("Test done. Sum: %d, Faild: %d, Time-Spanding: %v",
		num, resChan.failedCount, endTime.Sub(startTime))
	return nil
}

func countFailed(r *resultChan) {
LOOP:
	for {
		select {
		case <-r.failedChan:
			r.failedCount++
		case <-r.doneChan:
			break LOOP
		}
	}
}
