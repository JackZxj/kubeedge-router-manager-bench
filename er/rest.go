package er

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func RunRestSub(rest string, num int) error {
	bucket := make(chan bool, 1024)
	var once sync.Once
	var startTime time.Time
	timeoutInSecond := 300
	timeout := time.Second * time.Duration(timeoutInSecond)

	handler := func(w http.ResponseWriter, r *http.Request) {
		defer fmt.Fprintf(w, "ok\n")
		once.Do(func() {
			startTime = time.Now()
		})
		_, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("read body err, %v\n", err)
			return
		}
		bucket <- true
	}
	if strings.HasPrefix(rest, "https://") {
		log.Println("Warning: HTTPS is not supported and will be ignored.")
		rest = strings.Replace(rest, "https://", "", 1)
	} else if strings.HasPrefix(rest, "http://") {
		rest = strings.Replace(rest, "http://", "", 1)
	}

	url := strings.SplitN(rest, "/", 2)
	receiveCount := 0

	go func() {
		http.HandleFunc("/"+url[1], handler)
		if err := http.ListenAndServe(url[0], nil); err != nil {
			log.Panic(err)
		}
	}()

LOOP:
	for {
		select {
		case <-time.After(timeout):
			endTime := time.Now()
			log.Printf("It timed out %d seconds.\nSum: %d, Faild: %d, Time-Spending: %v",
				timeoutInSecond, num, num-receiveCount, endTime.Sub(startTime.Add(timeout)))
			break LOOP
		case <-bucket:
			receiveCount++
			if receiveCount >= num {
				endTime := time.Now()
				log.Printf("Test done.\nSum: %d, Faild: %d, Time-Spending: %v",
					num, num-receiveCount, endTime.Sub(startTime))
				break LOOP
			}
		}
	}

	return nil
}
