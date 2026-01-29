package config

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"regexp"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func Test_convert(t *testing.T) {
	tests := []struct {
		source   yamlConfig
		expected *Config
		wantErr  bool
	}{
		{
			source:  yamlConfig{},
			wantErr: true,
		},
		{
			source: yamlConfig{
				Collector: []*yamlCollectorConfig{
					{
						Host: "192.0.2.1",
					},
				},
			},
			wantErr: true,
		},
		{
			source: yamlConfig{
				ApiKey: "cat",
				Collector: []*yamlCollectorConfig{
					{
						HostID: "panda",
						Host:   "192.0.2.1",
					},
				},
			},
			expected: &Config{
				ApiKey: "cat",
				Collector: []*CollectorConfig{
					{
						HostID: "panda",
						Host:   "192.0.2.1",
					},
				},
			},
		},
	}

	opt1 := cmpopts.SortSlices(func(i, j string) bool { return i < j })

	opt2 := cmp.Comparer(func(x, y *regexp.Regexp) bool {
		if x == nil || y == nil {
			return x == y
		}

		return fmt.Sprint(x) == fmt.Sprint(y)
	})

	for _, tc := range tests {
		actual, err := convert(t.Context(), tc.source)
		if (err != nil) != tc.wantErr {
			t.Error(err)
		}

		if diff := cmp.Diff(actual, tc.expected, opt1, opt2); diff != "" {
			t.Errorf("value is mismatch (-actual +expected):%s", diff)
		}
	}
}

type clientMock struct {
	customIdentifier string
	hostId           string
	err              error
}

func (c *clientMock) FindHostByCustomIdentifierContext(_ context.Context, customIdentifier string) (string, error) {
	c.customIdentifier = customIdentifier
	return c.hostId, c.err
}

func Test_convertCollectors(t *testing.T) {
	tests := []struct {
		name       string
		source     []*yamlCollectorConfig
		expected   []*CollectorConfig
		wantErr    bool
		clientMock clientMock
	}{
		{
			name: "empty",
			// source:     []*yamlCollectorConfig{},
			clientMock: clientMock{},
			wantErr:    false,
		},
		{
			name: "host not found",
			source: []*yamlCollectorConfig{
				{
					HostID: "aaa",
				},
			},
			clientMock: clientMock{},
			wantErr:    true,
		},
		{
			name: "host-id not found",
			source: []*yamlCollectorConfig{
				{
					Host: "192.0.2.1",
				},
			},
			clientMock: clientMock{},
			wantErr:    true,
		},
		{
			name: "valid, host-id",
			source: []*yamlCollectorConfig{
				{
					HostID: "panda",
					Host:   "192.0.2.1",
				},
			},
			expected: []*CollectorConfig{
				{
					HostID: "panda",
					Host:   "192.0.2.1",
				},
			},
			clientMock: clientMock{},
		},
		{
			name: "valid, custom-identifier",
			source: []*yamlCollectorConfig{
				{
					CustomIdentifier: "cat",
					Host:             "192.0.2.1",
				},
			},
			expected: []*CollectorConfig{
				{
					HostID: "panda",
					Host:   "192.0.2.1",
				},
			},
			clientMock: clientMock{
				hostId: "panda",
			},
		},
		{
			name: "invalid, custom-identifier",
			source: []*yamlCollectorConfig{
				{
					CustomIdentifier: "cat",
					Host:             "192.0.2.1",
				},
			},
			expected: nil,
			clientMock: clientMock{
				err: fmt.Errorf("error"),
			},
			wantErr: true,
		},
		{
			name: "invalid, both custom-identifier, host-id",
			source: []*yamlCollectorConfig{
				{
					CustomIdentifier: "cat",
					Host:             "192.0.2.1",
					HostID:           "dog",
				},
			},
			expected:   nil,
			clientMock: clientMock{},
			wantErr:    true,
		},
	}

	opt1 := cmpopts.SortSlices(func(i, j string) bool { return i < j })

	opt2 := cmp.Comparer(func(x, y *regexp.Regexp) bool {
		if x == nil || y == nil {
			return x == y
		}

		return fmt.Sprint(x) == fmt.Sprint(y)
	})

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			originalLogger := slog.Default()
			defer slog.SetDefault(originalLogger)

			var buf bytes.Buffer
			slog.SetDefault(slog.New(slog.NewJSONHandler(&buf, nil)))

			actual := convertCollectors(t.Context(), &tc.clientMock, tc.source, false)
			if (buf.String() != "") != tc.wantErr {
				t.Error(buf.String())
			}

			if diff := cmp.Diff(actual, tc.expected, opt1, opt2); diff != "" {
				t.Errorf("value is mismatch (-actual +expected):%s", diff)
			}
		})
	}
}
