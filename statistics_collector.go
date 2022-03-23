package main

import (
	"log"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

type Statistics struct {
	numClients         float64
	numRecoveries      float64
	clientPacketsSent  float64
	clientPacketsRecv  float64
	maxHopCount        float64
	numCallLatency     float64
	minCallLatency     float64
	avgCallLatency     float64
	maxCallLatency     float64
	numLockwaitLatency float64
	minLockwaitLatency float64
	avgLockwaitLatency float64
	maxLockwaitLatency float64
}

type StatisticsCollector struct {
	numClients         *prometheus.Desc
	numRecoveries      *prometheus.Desc
	clientPacketsSent  *prometheus.Desc
	clientPacketsRecv  *prometheus.Desc
	maxHopCount        *prometheus.Desc
	numCallLatency     *prometheus.Desc
	minCallLatency     *prometheus.Desc
	avgCallLatency     *prometheus.Desc
	maxCallLatency     *prometheus.Desc
	numLockwaitLatency *prometheus.Desc
	minLockwaitLatency *prometheus.Desc
	avgLockwaitLatency *prometheus.Desc
	maxLockwaitLatency *prometheus.Desc
}

func NewStatisticsCollector() *StatisticsCollector {
	return &StatisticsCollector{
		numClients: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "num_clients"),
			"CTDB active clients connexions", []string{"id"}, nil,
		),
		numRecoveries: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "num_recoveries"),
			"Number of recoveries since the start of ctdb or since the last statistics reset", []string{"id"}, nil,
		),
		clientPacketsSent: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "client_packets_sent"),
			"Number of packets sent to client processes via unix domain socket", []string{"id"}, nil,
		),
		clientPacketsRecv: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "client_packets_recv"),
			"Number of packets received from client processes via unix domain socket", []string{"id"}, nil,
		),
		maxHopCount: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "max_hop_count"),
			"The maximum number of hops required for a record migration request to obtain the record", []string{"id"}, nil,
		),
		numCallLatency: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "num_call_latency"),
			"Number of REQ_CALL messages from client", []string{"id"}, nil,
		),
		minCallLatency: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "min_call_latency"),
			"Minimum time (in seconds) required to process a REQ_CALL message from client", []string{"id"}, nil,
		),
		avgCallLatency: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "avg_call_latency"),
			"Average time (in seconds) required to process a REQ_CALL message from client", []string{"id"}, nil,
		),
		maxCallLatency: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "max_call_latency"),
			"Maximum time (in seconds) required to process a REQ_CALL message from client", []string{"id"}, nil,
		),
		numLockwaitLatency: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "num_lockwait_latency"),
			"Number of record locks", []string{"id"}, nil,
		),
		minLockwaitLatency: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "min_lockwait_latency"),
			"Minimum time (in seconds) required to obtain record locks", []string{"id"}, nil,
		),
		avgLockwaitLatency: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "avg_lockwait_latency"),
			"Average time (in seconds) required to obtain record locks", []string{"id"}, nil,
		),
		maxLockwaitLatency: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "max_lockwait_latency"),
			"Maximum time (in seconds) required to obtain record locks", []string{"id"}, nil,
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
	ch <- prometheus.MustNewConstMetric(c.numRecoveries, prometheus.CounterValue, statistics.numRecoveries, pnn)
	ch <- prometheus.MustNewConstMetric(c.clientPacketsSent, prometheus.CounterValue, statistics.clientPacketsSent, pnn)
	ch <- prometheus.MustNewConstMetric(c.clientPacketsRecv, prometheus.CounterValue, statistics.clientPacketsRecv, pnn)
	ch <- prometheus.MustNewConstMetric(c.maxHopCount, prometheus.GaugeValue, statistics.maxHopCount, pnn)
	ch <- prometheus.MustNewConstMetric(c.numCallLatency, prometheus.CounterValue, statistics.numCallLatency, pnn)
	ch <- prometheus.MustNewConstMetric(c.minCallLatency, prometheus.GaugeValue, statistics.minCallLatency, pnn)
	ch <- prometheus.MustNewConstMetric(c.avgCallLatency, prometheus.GaugeValue, statistics.avgCallLatency, pnn)
	ch <- prometheus.MustNewConstMetric(c.maxCallLatency, prometheus.GaugeValue, statistics.maxCallLatency, pnn)
	ch <- prometheus.MustNewConstMetric(c.numLockwaitLatency, prometheus.CounterValue, statistics.numLockwaitLatency, pnn)
	ch <- prometheus.MustNewConstMetric(c.minLockwaitLatency, prometheus.GaugeValue, statistics.minLockwaitLatency, pnn)
	ch <- prometheus.MustNewConstMetric(c.avgLockwaitLatency, prometheus.GaugeValue, statistics.avgLockwaitLatency, pnn)
	ch <- prometheus.MustNewConstMetric(c.maxLockwaitLatency, prometheus.GaugeValue, statistics.maxLockwaitLatency, pnn)
}

func (c *StatisticsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.numClients
	ch <- c.numRecoveries
	ch <- c.clientPacketsSent
	ch <- c.clientPacketsRecv
	ch <- c.maxHopCount
	ch <- c.numCallLatency
	ch <- c.minCallLatency
	ch <- c.avgCallLatency
	ch <- c.maxCallLatency
	ch <- c.numLockwaitLatency
	ch <- c.minLockwaitLatency
	ch <- c.avgLockwaitLatency
	ch <- c.maxLockwaitLatency
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
		case "num_recoveries":
			numRecoveries, err := strconv.Atoi(val)
			if err != nil {
				continue
			}
			statistics.numRecoveries = float64(numRecoveries)
		case "client_packets_sent":
			clientPacketsSent, err := strconv.Atoi(val)
			if err != nil {
				continue
			}
			statistics.clientPacketsSent = float64(clientPacketsSent)
		case "client_packets_recv":
			clientPacketsRecv, err := strconv.Atoi(val)
			if err != nil {
				continue
			}
			statistics.clientPacketsRecv = float64(clientPacketsRecv)
		case "max_hop_count":
			maxHopCount, err := strconv.Atoi(val)
			if err != nil {
				continue
			}
			statistics.maxHopCount = float64(maxHopCount)
		case "num_call_latency":
			numCallLatency, err := strconv.Atoi(val)
			if err != nil {
				continue
			}
			statistics.numCallLatency = float64(numCallLatency)
		case "min_call_latency":
			minCallLatency, err := strconv.ParseFloat(val, 64)
			if err != nil {
				continue
			}
			statistics.minCallLatency = minCallLatency
		case "avg_call_latency":
			avgCallLatency, err := strconv.ParseFloat(val, 64)
			if err != nil {
				continue
			}
			statistics.avgCallLatency = avgCallLatency
		case "max_call_latency":
			maxCallLatency, err := strconv.ParseFloat(val, 64)
			if err != nil {
				continue
			}
			statistics.maxCallLatency = maxCallLatency
		case "num_lockwait_latency":
			numLockwaitLatency, err := strconv.Atoi(val)
			if err != nil {
				continue
			}
			statistics.numLockwaitLatency = float64(numLockwaitLatency)
		case "min_lockwait_latency":
			minLockwaitLatency, err := strconv.ParseFloat(val, 64)
			if err != nil {
				continue
			}
			statistics.minLockwaitLatency = minLockwaitLatency
		case "avg_lockwait_latency":
			avgLockwaitLatency, err := strconv.ParseFloat(val, 64)
			if err != nil {
				continue
			}
			statistics.avgLockwaitLatency = avgLockwaitLatency
		case "max_lockwait_latency":
			maxLockwaitLatency, err := strconv.ParseFloat(val, 64)
			if err != nil {
				continue
			}
			statistics.maxLockwaitLatency = maxLockwaitLatency
		}
	}

	return &statistics, nil
}
