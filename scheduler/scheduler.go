package scheduler

import (
	"context"
	"sync_study/entity"
	"sync_study/logger"
	"sync_study/utils"
	"time"
)

// Scheduler
type (
	RunScheduler interface {
		Run()
		SubmitScheduler
	}

	SubmitScheduler interface {
		Ready(chan entity.Request)
		Submit(entity.Request)
		Worker() chan entity.Request
	}

	Scheduler struct {
		WorkerCount uint
		scheduler   RunScheduler
		ctx         context.Context
	}
)

// NewScheduler new scheduler
func NewScheduler(workerCount uint, ctx context.Context) *Scheduler {
	return &Scheduler{
		WorkerCount: workerCount,
		ctx:         ctx,
		scheduler:   &queue{ctx: ctx},
	}
}

// Submit -
func (s *Scheduler) Submit(reqs ...entity.Request) {
	for _, r := range reqs {
		s.scheduler.Submit(r)
	}
}

// Run -
func (s *Scheduler) Run() {

	defer utils.Recover("[Scheduler][Run]")

	go s.scheduler.Run()
	time.Sleep(2 * time.Second)
	out := s.scheduler.Worker()
	for i := 0; i < int(s.WorkerCount); i++ {
		s.worker(s.scheduler.Worker(), out, s.scheduler)
	}

	for {
		select {
		case <-out:
		case <-s.ctx.Done():
			return
		}
	}
}

func (s *Scheduler) worker(in chan entity.Request, out chan entity.Request, scheduler SubmitScheduler) {
	go func() {
		defer utils.Recover("[Scheduler][worker]")
		for {
			scheduler.Ready(in)
			select {

			case r := <-in:
				if r.URI == "" {
					out <- r
					logger.Info("empty uri")
					continue
				}

				// data, err := utils.HTTPGetResponse(r.URI, nil) // *http.Response
				data, err := utils.HTTPGet(r.URI, nil)
				if err != nil {
					out <- r
					logger.Errorf("request [%s] error[%s]", r.URI, err)
					continue
				}
				result := r.ParseHandler(&entity.RequestResult{
					Body:          data,
					Data:          r.Data,
					RequestConfig: r.RequestConfig,
					SourceURI:     r.URI,
				})

				if result == nil {
					out <- r
					continue
				}
				handleParseReuqest(result, scheduler)

				out <- r

			case <-s.ctx.Done():
				return
			}
		}
	}()
}

func handleParseReuqest(request *entity.ParseRequest, scheduler SubmitScheduler) {
	if request == nil {
		return
	}

	if request.Items == nil {
		return
	}

	for i := range request.Items {
		if v, ok := request.Items[i].(*entity.Request); ok {
			scheduler.Submit(*v)
		}
	}
}
