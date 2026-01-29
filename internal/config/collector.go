package config

import (
	"context"
	"fmt"
	"log/slog"
	"time"
)

type mackerelClient interface {
	FindHostByCustomIdentifierContext(ctx context.Context, customIdentifier string) (string, error)
}

func convertCollectors(ctx context.Context, client mackerelClient, collector []*yamlCollectorConfig, privileged bool) []*CollectorConfig {
	var cs []*CollectorConfig
	for i := range collector {
		conf, err := convertCollector(ctx, client, collector[i], privileged)
		if err != nil {
			slog.Warn("skipped because failed parse config", slog.Int("index", i), slog.String("errror", err.Error()))
			continue
		}
		cs = append(cs, conf)
	}
	return cs
}

func convertCollector(ctx context.Context, client mackerelClient, t *yamlCollectorConfig, privileged bool) (*CollectorConfig, error) {
	if t.Host == "" {
		return nil, fmt.Errorf("host is needed")
	}
	if t.HostID == "" && t.CustomIdentifier == "" {
		return nil, fmt.Errorf("host-id or custom-identifier is needed")
	}
	if t.HostID != "" && t.CustomIdentifier != "" {
		return nil, fmt.Errorf("host-id, custom-identifier is exclusive")
	}

	if t.CustomIdentifier != "" {
		var err error
		for range 3 {
			t.HostID, err = client.FindHostByCustomIdentifierContext(ctx, t.CustomIdentifier)
			if err != nil {
				slog.WarnContext(ctx, "retry: host id is invalid",
					slog.String("custom-identifier", t.CustomIdentifier),
					slog.String("error", err.Error()),
				)
				time.Sleep(100 * time.Millisecond)
				continue
			}
			break
		}
		if err != nil {
			return nil, fmt.Errorf("host id is invalid, custom-identifier: %s, error: %w", t.CustomIdentifier, err)
		}
	}

	c := &CollectorConfig{
		HostID: t.HostID,

		Host:       t.Host,
		Average:    t.Average,
		Privileged: privileged,
	}

	return c, nil
}
