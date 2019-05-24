package scheduler

import (
	"context"
	"sync_study/entity"
	"sync_study/utils"
)

type queue struct {
	workerChan chan entity.Request
	readyChan  chan chan entity.Request
	ctx        context.Context
}

func (q *queue) Ready(r chan entity.Request) {
	q.readyChan <- r
}

func (q *queue) Submit(r entity.Request) {
	q.workerChan <- r
}

func (q *queue) Run() {

	defer utils.Recover("[queue][Run]")

	q.workerChan = q.Worker()
	q.readyChan = make(chan chan entity.Request)
	go func() {
		var (
			workerQ []entity.Request
			readyQ  []chan entity.Request
		)
		for {
			var (
				w entity.Request
				r chan entity.Request
			)
			if len(workerQ) > 0 && len(readyQ) > 0 {
				w = workerQ[0]
				r = readyQ[0]
			}

			select {
			case c := <-q.workerChan:
				workerQ = append(workerQ, c)
			case c := <-q.readyChan:
				readyQ = append(readyQ, c)
			case r <- w:
				workerQ = workerQ[1:]
				readyQ = readyQ[1:]
			case <-q.ctx.Done():
				return
			}
		}
	}()
}

func (q *queue) Worker() chan entity.Request {
	return make(chan entity.Request)
}
