package internal

import (
	"context"
	"github.com/signalfx/golib/sfxclient"
	"time"
)

type Scheduler interface {
	ReportingDelay(time.Duration)
	AddCallback(collector sfxclient.Collector)
	Schedule(ctx context.Context) error
	ReportOnce(ctx context.Context) error
}
