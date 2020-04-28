package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"strconv"
	"strings"
)

type Statistics struct {
	numClients float64
}

type StatisticsCollector struct {
	numClients *prometheus.Desc
}

func NewStatisticsCollector() *StatisticsCollector {
	return &StatisticsCollector{
		numClients: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "num_clients"),
			"CTDB active clients connexions", []string{"id"}, nil,
		),
	}
}

func (c *StatisticsCollector) Collect(ch chan<- prometheus.Metric) {
	if pnn == "" {
		return
	}

	statistics, err := scrapeStatistics(runCmd)
	if err != nil {
		log.Println(err)
		return
	}

	ch <- prometheus.MustNewConstMetric(c.numClients, prometheus.GaugeValue, statistics.numClients, pnn)
}

func (c *StatisticsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.numClients
}

func scrapeStatistics(run runner) (*Statistics, error) {
	status, err := run("statistics -Y")
	if err != nil {
		return nil, err
	}

	lines := strings.Split(status, "\n")
	statistics := Statistics{}
	headers := strings.Split(lines[0], "|")

	for field, val := range strings.Split(lines[1], "|") {
		if val == "" {
			continue
		}
		switch headers[field] {
		case "num_clients":
			numClients, err := strconv.Atoi(val)
			if err != nil {
				continue
			}
			statistics.numClients = float64(numClients)
		}
	}

	return &statistics, nil
}
