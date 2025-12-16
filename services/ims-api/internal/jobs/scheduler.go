package jobs

import (
	"context"
	"sync"
	"time"

	"github.com/edvirons/ssp/ims/internal/logging"
	"github.com/edvirons/ssp/ims/internal/store"
	"go.uber.org/zap"
)

type Scheduler struct {
	log *zap.Logger
	pg  *store.Postgres

	wg   sync.WaitGroup
	stop chan struct{}
}

func NewScheduler(log *zap.Logger, pg *store.Postgres) *Scheduler {
	return &Scheduler{log: log, pg: pg, stop: make(chan struct{})}
}

func (s *Scheduler) Start(ctx context.Context) {
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		t := time.NewTicker(60 * time.Second)
		defer t.Stop()

		for {
			select {
			case <-ctx.Done():
				s.log.Info("jobs: context done")
				return
			case <-s.stop:
				s.log.Info("jobs: stopped")
				return
			case <-t.C:
				n, err := s.pg.Incidents().MarkSLABreaches(ctx, time.Now().UTC())
				if err != nil {
					s.log.Warn("jobs: mark sla breaches failed", logging.Err(err))
					continue
				}
				if n > 0 {
					s.log.Info("jobs: sla breaches updated", zap.Int("count", n))
				}
			}
		}
	}()
}

func (s *Scheduler) Stop() {
	close(s.stop)
	s.wg.Wait()
}
