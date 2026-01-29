package ticker

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/mackerelio/mackerel-client-go"

	"github.com/mackerelio-labs/sabapingd/internal/collector"
	"github.com/mackerelio-labs/sabapingd/internal/config"
)

type enqueuer interface {
	Enqueue(hostID string, rawMetrics []*mackerel.MetricValue)
}

type collectorIface interface {
	Do(ctx context.Context) (collector.Result, error)
}

type Ticker struct {
	mu sync.RWMutex

	collectorID string
	hostID      string
	queue       enqueuer
	collector   collectorIface

	previousPacketLoss float64
}

func New(conf *config.CollectorConfig, q enqueuer) *Ticker {
	return &Ticker{
		collectorID: conf.CollectorID(),
		hostID:      conf.HostID,
		queue:       q,
		collector:   collector.New(conf),
	}
}

func (t *Ticker) Tick(ctx context.Context) {
	t.do(ctx, time.Now())
}

func (t *Ticker) do(ctx context.Context, now time.Time) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	result, err := t.collector.Do(ctx)
	if err != nil {
		slog.WarnContext(ctx, "failed exec collector.Do()", slog.String("error", err.Error()))
		return
	}

	var queue = make([]*mackerel.MetricValue, 0)

	// 全てパケットロスした場合、avgRtt = 0 となるため
	if result.PacketLoss < 100 {
		queue = append(queue, &mackerel.MetricValue{
			Name:  "custom.sabapingd.rtt.Average",
			Time:  now.Unix(),
			Value: millisecond(result.Average),
		})
		if !result.AverageOnly {
			queue = append(queue, &mackerel.MetricValue{
				Name:  "custom.sabapingd.rtt.Maximum",
				Time:  now.Unix(),
				Value: millisecond(result.Maximum),
			}, &mackerel.MetricValue{
				Name:  "custom.sabapingd.rtt.Minimum",
				Time:  now.Unix(),
				Value: millisecond(result.Minimum),
			})
		}
	}

	// 前回パケットロスがあった場合で、回復した時に 0 を投稿したいため
	if t.previousPacketLoss > 0 || result.PacketLoss > 0 {
		queue = append(queue, &mackerel.MetricValue{
			Name:  "custom.sabapingd.packetLoss.measure", // TODO
			Time:  now.Unix(),
			Value: result.PacketLoss,
		})
		t.previousPacketLoss = result.PacketLoss
	}
	t.queue.Enqueue(t.hostID, queue)
}

func millisecond(d time.Duration) float64 {
	sec := d / time.Millisecond
	nsec := d % time.Millisecond
	return float64(sec) + float64(nsec)/1e6
}

func (t *Ticker) Reload(conf *config.CollectorConfig) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.hostID = conf.HostID
	t.collector = collector.New(conf)
}

func (t *Ticker) CollectorID() string {
	return t.collectorID
}
