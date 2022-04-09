package fluentkafkarouter

import "github.com/UESTC-KEEP/keep/cloud/pkg/requestDispatcher/Router/routers/kafka/metrics"

type FluentKafkaRouter struct {
	Metrics metrics.Metrics
}
