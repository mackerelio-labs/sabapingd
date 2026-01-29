package collector

import (
	"context"
	"time"

	probing "github.com/prometheus-community/pro-bing"

	"github.com/mackerelio-labs/sabapingd/internal/config"
)

type collector struct {
	conf *config.CollectorConfig
}

func New(conf *config.CollectorConfig) *collector {
	return &collector{conf: conf}
}

type Result struct {
	Average, Maximum, Minimum time.Duration

	PacketLoss  float64
	AverageOnly bool
}

func (c *collector) Do(ctx context.Context) (Result, error) {
	pinger, err := probing.NewPinger(c.conf.Host)
	if err != nil {
		return Result{}, err
	}
	pinger.Count = 3
	pinger.Timeout = 5 * time.Second // TODO
	pinger.SetPrivileged(c.conf.Privileged)

	if err = pinger.RunWithContext(ctx); err != nil {
		return Result{}, err
	}
	stats := pinger.Statistics()

	return Result{
		Maximum:    stats.MaxRtt,
		Minimum:    stats.MinRtt,
		Average:    stats.AvgRtt,
		PacketLoss: stats.PacketLoss,

		AverageOnly: c.conf.Average,
	}, nil
}
