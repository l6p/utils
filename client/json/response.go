package json

import (
	"errors"
	"time"
)

type Response interface {
	D() *Data
	RedoUntil(checker func(data *Data) bool, timeout time.Duration)
}

type ResponseImpl struct {
	data *Data
	req  Request
}

func (r *ResponseImpl) D() *Data {
	return r.data
}

func (r *ResponseImpl) RedoUntil(checker func(data *Data) bool, timeout time.Duration) {
	if checker(r.D()) {
		return
	}

	t1 := time.NewTicker(2 * time.Second)
	t2 := time.NewTimer(timeout)
	for {
		select {
		case <-t1.C:
			if checker(r.req.Redo().D()) {
				return
			}
		case <-t2.C:
			panic(errors.New("redo timeout"))
		}
	}
}
