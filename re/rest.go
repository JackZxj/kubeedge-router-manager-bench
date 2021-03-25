package re

import (
	"errors"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/JackZxj/kubeedge-router-manager-bench/util"
)

type resultChan struct {
	failedChan  chan int
	doneChan    chan bool
	failedCount int
}

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func RunRestPub(rest, restMethod, msg, msgBinary string, num, concurrency int) error {
	msgToSend, err := util.GetMsgToSend(msg, msgBinary)
	if err != nil {
		return err
	}

	var hanler util.HttpHandler
	switch strings.ToUpper(restMethod) {
	case "POST":
		hanler = util.Post
	default:
		return errors.New("unsupport method: " + restMethod)
	}

	_, err = url.Parse(rest)
	if err != nil {
		return err
	}

	pool := util.NewPool(concurrency)
	resChan := util.NewResult(concurrency)
	startTime := time.Now()
	log.Println("Sending start:", startTime.UTC())

	go resChan.Loop()
	for i := 0; i < num; i++ {
		pool.Add(1)
		go func() {
			defer pool.Done()
			_, err := hanler(rest, msgToSend)
			if err != nil {
				resChan.FailedChan <- 1
				// Panic 会 crash 掉整个进程
				log.Print("Request failed:", err)
			}
		}()
	}
	pool.Wait()
	endTime := time.Now()
	log.Println("Sending end:", endTime.UTC())
	resChan.DoneChan <- true
	log.Printf("Test done.\nSum: %d, Faild: %d, Time-Spending: %v",
		num, resChan.FailedCount, endTime.Sub(startTime))
	return nil
}
