package dialogue

import (
	"context"
	"sync"

	"city-of-lies/backend/internal/domain"
)

type DialogueJob struct {
	SessionID   string
	AgentID     string
	Context     domain.DialogueContext
	PlayerQuery string
	Callback    func(resp domain.DialogueResponse, err error)
}

type DialogueQueue struct {
	jobs    chan DialogueJob
	workers int
	wg      sync.WaitGroup
	ctx     context.Context
	cancel  context.CancelFunc
}

func NewDialogueQueue(workers int, capacity int) *DialogueQueue {
	if workers <= 0 {
		workers = 3
	}
	if capacity <= 0 {
		capacity = 50
	}
	ctx, cancel := context.WithCancel(context.Background())
	return &DialogueQueue{
		jobs:    make(chan DialogueJob, capacity),
		workers: workers,
		ctx:     ctx,
		cancel:  cancel,
	}
}

func (q *DialogueQueue) Start(provider interface {
	GenerateDialogue(ctx context.Context, dCtx domain.DialogueContext, playerQuery string) (domain.DialogueResponse, error)
}) {
	for i := 0; i < q.workers; i++ {
		q.wg.Add(1)
		go func() {
			defer q.wg.Done()
			for {
				select {
				case <-q.ctx.Done():
					return
				case job, ok := <-q.jobs:
					if !ok {
						return
					}
					resp, err := provider.GenerateDialogue(q.ctx, job.Context, job.PlayerQuery)
					if job.Callback != nil {
						job.Callback(resp, err)
					}
				}
			}
		}()
	}
}

func (q *DialogueQueue) Enqueue(job DialogueJob) bool {
	select {
	case q.jobs <- job:
		return true
	default:
		return false // Queue full
	}
}

func (q *DialogueQueue) Stop() {
	q.cancel()
	close(q.jobs)
	q.wg.Wait()
}
