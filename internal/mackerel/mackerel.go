package mackerel

import (
	"cmp"
	"context"
	"errors"
	"os"

	mackerel "github.com/mackerelio/mackerel-client-go"
)

type mackerelClient interface {
	CreateGraphDefsContext(ctx context.Context, payloads []*mackerel.GraphDefsParam) error
	PostHostMetricValuesByHostIDContext(ctx context.Context, hostID string, metricValues []*mackerel.MetricValue) error
	FindHostByCustomIdentifierContext(ctx context.Context, customIdentifier string, param *mackerel.FindHostByCustomIdentifierParam) (*mackerel.Host, error)
}

type Mackerel struct {
	client mackerelClient
}

func New(apikey string) *Mackerel {
	baseURL := cmp.Or(os.Getenv("MACKEREL_APIBASE"), "https://api.mackerelio.com/")
	client, _ := mackerel.NewClientWithOptions(apikey, baseURL, false)
	return &Mackerel{
		client: client,
	}
}

func (m *Mackerel) CreateGraphDefs(ctx context.Context, d []*mackerel.GraphDefsParam) error {
	return m.client.CreateGraphDefsContext(ctx, d)
}

func (m *Mackerel) Send(ctx context.Context, hostID string, value []*mackerel.MetricValue) error {
	return m.client.PostHostMetricValuesByHostIDContext(ctx, hostID, value)
}

func (m *Mackerel) FindHostByCustomIdentifierContext(ctx context.Context, customIdentifier string) (string, error) {
	host, err := m.client.FindHostByCustomIdentifierContext(ctx, customIdentifier, &mackerel.FindHostByCustomIdentifierParam{
		CaseInsensitive: false,
	})
	if err != nil {
		return "", err
	}
	if host == nil {
		return "", errors.New("host not found")
	}
	return host.ID, nil
}
