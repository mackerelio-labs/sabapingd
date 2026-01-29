package config

import (
	"cmp"
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/mackerelio-labs/sabapingd/internal/mackerel"

	"gopkg.in/yaml.v3"
)

type yamlCollectorConfig struct {
	HostID           string `yaml:"host-id"`
	CustomIdentifier string `yaml:"custom-identifier"`

	// for ping
	Host    string `yaml:"host"`
	Average bool   `yaml:"average"`
}

type yamlDiskCache struct {
	Directory string `yaml:"directory"`
	Size      Size   `yaml:"size"`
}

type yamlConfig struct {
	ApiKey     string `yaml:"x-api-key"`
	Privileged bool   `yaml:"privileged"`

	Collector []*yamlCollectorConfig `yaml:"collector"`

	DiskCache *yamlDiskCache `yaml:"disk-cache"`
}

type CollectorConfig struct {
	HostID string

	// for ping
	Host       string
	Average    bool
	Privileged bool
}

func (conf *CollectorConfig) CollectorID() string {
	return fmt.Sprintf("host=%s,hostID=%s", conf.Host, conf.HostID)
}

type DiskCache struct {
	Directory string
	Size      Size
}

type Config struct {
	ApiKey string

	Collector []*CollectorConfig
	DiskCache *DiskCache
}

func Init(ctx context.Context, filename string) (*Config, error) {
	f, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var t yamlConfig
	err = yaml.Unmarshal(f, &t)
	if err != nil {
		return nil, err
	}
	return convert(ctx, t)
}

func convert(ctx context.Context, t yamlConfig) (*Config, error) {
	apiKey := cmp.Or(os.Getenv("MACKEREL_APIKEY"), t.ApiKey)
	if apiKey == "" {
		return nil, fmt.Errorf("x-api-key is needed")
	}

	client := mackerel.New(apiKey)
	cs := convertCollectors(ctx, client, t.Collector, t.Privileged)

	var dc *DiskCache
	if t.DiskCache != nil {
		var err error
		dc, err = diskcacheValidate(t.DiskCache)
		if err != nil {
			slog.Warn("disable disk-cache because failed parse config", slog.String("error", err.Error()))
		}
	}

	return &Config{
		ApiKey:    apiKey,
		Collector: cs,
		DiskCache: dc,
	}, nil
}
