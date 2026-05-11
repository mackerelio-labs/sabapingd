package ticker

import (
	"context"
	"testing"
	"time"

	"github.com/mackerelio/mackerel-client-go"

	"github.com/mackerelio-labs/sabapingd/internal/collector"
)

type mockEnqueuer struct {
	hostID  string
	metrics []*mackerel.MetricValue
}

func (m *mockEnqueuer) Enqueue(hostID string, rawMetrics []*mackerel.MetricValue) {
	m.hostID = hostID
	m.metrics = rawMetrics
}

type mockCollector struct {
	result collector.Result
	err    error
}

func (m *mockCollector) Do(_ context.Context) (collector.Result, error) {
	return m.result, m.err
}

func TestTicker_Do_SkipContinuousZeroPacketLossByDefault(t *testing.T) {
	now := time.Unix(1700000000, 0)
	queue := &mockEnqueuer{}
	collector := &mockCollector{
		result: collector.Result{
			Average:    10 * time.Millisecond,
			Maximum:    12 * time.Millisecond,
			Minimum:    8 * time.Millisecond,
			PacketLoss: 0,
		},
	}
	tk := &Ticker{
		hostID:    "host-1",
		queue:     queue,
		collector: collector,
	}

	tk.do(t.Context(), now)

	if hasMetric(queue.metrics, "custom.sabapingd.packetLoss.measure") {
		t.Fatal("packetLoss metric should not be posted when 0 continues in default mode")
	}
}

func TestTicker_Do_PostZeroPacketLossWhenEnabled(t *testing.T) {
	now := time.Unix(1700000000, 0)
	queue := &mockEnqueuer{}
	collector := &mockCollector{
		result: collector.Result{
			Average:    10 * time.Millisecond,
			Maximum:    12 * time.Millisecond,
			Minimum:    8 * time.Millisecond,
			PacketLoss: 0,
		},
	}
	tk := &Ticker{
		hostID:               "host-1",
		queue:                queue,
		collector:            collector,
		alwaysSendPacketLoss: true,
	}

	tk.do(t.Context(), now)

	if !hasMetric(queue.metrics, "custom.sabapingd.packetLoss.measure") {
		t.Fatal("packetLoss metric should be posted when always-send-packetloss is enabled")
	}
}

func hasMetric(metrics []*mackerel.MetricValue, name string) bool {
	for _, m := range metrics {
		if m.Name == name {
			return true
		}
	}
	return false
}
