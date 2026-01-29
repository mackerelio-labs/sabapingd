package mackerel

import (
	"context"
	"testing"

	"github.com/mackerelio/mackerel-client-go"
)

type mackerelClientMock struct {
	graphDef     []*mackerel.GraphDefsParam
	hostID       string
	metricValues []*mackerel.MetricValue

	returnError         error
	returnErrorGraphDef error
}

func (m *mackerelClientMock) CreateGraphDefsContext(_ context.Context, payloads []*mackerel.GraphDefsParam) error {
	m.graphDef = payloads
	return m.returnErrorGraphDef
}
func (m *mackerelClientMock) PostHostMetricValuesByHostIDContext(_ context.Context, hostID string, metricValues []*mackerel.MetricValue) error {
	m.hostID = hostID
	m.metricValues = metricValues
	return m.returnError
}
func (m *mackerelClientMock) FindHostByCustomIdentifierContext(ctx context.Context, customIdentifier string, param *mackerel.FindHostByCustomIdentifierParam) (*mackerel.Host, error) {
	return nil, nil
}

func TestSend(t *testing.T) {
	mock := &mackerelClientMock{}
	mc := &Mackerel{
		client: mock,
	}

	if err := mc.Send(t.Context(), "0987654321", nil); err != nil {
		t.Errorf("occur error %v", err)
	}

	if mock.hostID == "" {
		t.Error("invalid need hostID")
	}

}
