package mackerel

import "github.com/mackerelio/mackerel-client-go"

var GraphDefs = []*mackerel.GraphDefsParam{
	{
		Name:        "custom.sabapingd.rtt",
		Unit:        "milliseconds",
		DisplayName: "Ping RTT",
		Metrics: []*mackerel.GraphDefsMetric{
			{
				Name:        "custom.sabapingd.rtt.*",
				DisplayName: "%1",
			},
		},
	},
	{
		Name:        "custom.sabapingd.packetLoss",
		Unit:        "percentage",
		DisplayName: "Ping PacketLoss",
		Metrics: []*mackerel.GraphDefsMetric{
			{
				Name:        "custom.sabapingd.packetLoss.*",
				DisplayName: "%1",
			},
		},
	},
}
