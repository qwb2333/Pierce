package lib

import (
	"time"
	"sync/atomic"
)

func WriteCh(ch chan bool, x bool) {
	select {
	case ch<-x:
	default:

	}
}

type RunStatus int32

func RunWithTimeout(f func() (bool, error), duration time.Duration) (RunStatus, error) {
	var err error
	done := make(chan bool)
	go func() {
		var result bool
		result, err = f()
		done<-result
	}()

	select {
	case <-time.After(duration):
		return RUN_STATUS_TIMEOUT, nil

	case result:= <-done:
		if result {
			return RUN_STATUS_SUCCESS, nil
		}
		return RUN_STATUS_FAILED, err
	}
}

var maxId  uint32
func GetNextId() uint32 {
	return atomic.AddUint32(&maxId, 1)
}