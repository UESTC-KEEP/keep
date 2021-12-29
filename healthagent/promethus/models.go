package prome

import "github.com/prometheus/client_golang/prometheus"

func AssembleModel(modelName, modelHelp string) prometheus.Gauge {
	return prometheus.NewGauge(prometheus.GaugeOpts{
		Name: modelName,
		Help: modelHelp,
	})
}
