package util

type Result struct {
	FailedChan  chan int
	DoneChan    chan bool
	FailedCount int
}

func NewResult(maxLenForFail int) *Result {
	return &Result{
		FailedChan:  make(chan int, maxLenForFail),
		DoneChan:    make(chan bool),
		FailedCount: 0,
	}
}

func (r *Result) Loop() {
LOOP:
	for {
		select {
		case <-r.FailedChan:
			r.FailedCount++
		case <-r.DoneChan:
			close(r.FailedChan)
			close(r.DoneChan)
			break LOOP
		}
	}
}
