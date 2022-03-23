package main

import (
	"log"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

type Status struct {
	id              string
	ip              string
	disconnected    float64
	banned          float64
	disabled        float64
	unhealthy       float64
	stopped         float64
	inactive        float64
	partiallyOnline float64
	thisNode        float64
}

type StatusCollector struct {
	up              *prometheus.Desc
	banned          *prometheus.Desc
	disconnected    *prometheus.Desc
	inactive        *prometheus.Desc
	partiallyOnline *prometheus.Desc
	stopped         *prometheus.Desc
	unhealthy       *prometheus.Desc
}

func NewStatusCollector() *StatusCollector {
	return &StatusCollector{
		up: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "up"),
			"Is CTDB running", nil, nil,
		),
		banned: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "banned"),
			"Is node banned", []string{"id", "ip"}, nil,
		),
		disconnected: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "disconnected"),
			"Is node disconnected", []string{"id", "ip"}, nil,
		),
		inactive: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "inactive"),
			"Is node inactive", []string{"id", "ip"}, nil,
		),
		partiallyOnline: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "partially_online"),
			"Is node partially Online", []string{"id", "ip"}, nil,
		),
		stopped: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "stopped"),
			"Is node stopped", []string{"id", "ip"}, nil,
		),
		unhealthy: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "unhealthy"),
			"Is node unhealthy", []string{"id", "ip"}, nil,
		),
	}
}

func (c *StatusCollector) Collect(ch chan<- prometheus.Metric) {
	masterNode, err := isMasterNode(runCmd)
	if err != nil {
		ch <- prometheus.MustNewConstMetric(c.up, prometheus.GaugeValue, 0)
		log.Println(err)
		// We cannot determine node master status, assume ctdb is broken/down on this node
		return
	}

	if !masterNode {
		ch <- prometheus.MustNewConstMetric(c.up, prometheus.GaugeValue, 1)
		// This node is not master, we don't need to scrape status as it is shared among all nodes
		return
	}

	status, err := scrapeStatus(runCmd)
	if err != nil {
		ch <- prometheus.MustNewConstMetric(c.up, prometheus.GaugeValue, 0)
		log.Println(err)
		// We cannot scrape status, assume ctdb is broken/down on this node
		return
	}

	ch <- prometheus.MustNewConstMetric(c.up, prometheus.GaugeValue, 1)
	for _, line := range status {
		ch <- prometheus.MustNewConstMetric(c.banned, prometheus.GaugeValue, line.banned, line.id, line.ip)
		ch <- prometheus.MustNewConstMetric(c.disconnected, prometheus.GaugeValue, line.disconnected, line.id, line.ip)
		ch <- prometheus.MustNewConstMetric(c.inactive, prometheus.GaugeValue, line.inactive, line.id, line.ip)
		ch <- prometheus.MustNewConstMetric(c.partiallyOnline, prometheus.GaugeValue, line.partiallyOnline, line.id, line.ip)
		ch <- prometheus.MustNewConstMetric(c.stopped, prometheus.GaugeValue, line.stopped, line.id, line.ip)
		ch <- prometheus.MustNewConstMetric(c.unhealthy, prometheus.GaugeValue, line.unhealthy, line.id, line.ip)
	}
}

func (c *StatusCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.up

	ch <- c.banned
	ch <- c.disconnected
	ch <- c.inactive
	ch <- c.partiallyOnline
	ch <- c.stopped
	ch <- c.unhealthy
}

func scrapeStatus(run runner) ([]Status, error) {
	status, err := run("status -Y")
	if err != nil {
		return nil, err
	}

	lines := strings.Split(status, "\n")
	table := make([]Status, 0, len(lines))
	headers := strings.Split(lines[0], "|")
	for idx, line := range lines {
		if idx < 1 {
			continue
		}
		values := Status{}
		for field, val := range strings.Split(line, "|") {
			if val == "" {
				continue
			}
			switch headers[field] {
			case "Node":
				values.id = val
			case "IP":
				values.ip = val
			case "Disconnected":
				if val == "1" {
					values.disconnected = 1
				} else {
					values.disconnected = 0
				}
			case "Banned":
				if val == "1" {
					values.banned = 1
				} else {
					values.banned = 0
				}
			case "Disabled":
				if val == "1" {
					values.disabled = 1
				} else {
					values.disabled = 0
				}
			case "Unhealthy":
				if val == "1" {
					values.unhealthy = 1
				} else {
					values.unhealthy = 0
				}
			case "Stopped":
				if val == "1" {
					values.stopped = 1
				} else {
					values.stopped = 0
				}
			case "Inactive":
				if val == "1" {
					values.inactive = 1
				} else {
					values.inactive = 0
				}
			case "PartiallyOnline":
				if val == "1" {
					values.partiallyOnline = 1
				} else {
					values.partiallyOnline = 0
				}
			case "ThisNode":
				if val == "Y" {
					values.thisNode = 1
				} else {
					values.thisNode = 0
				}
			}
		}
		table = append(table, values)
	}

	return table, nil
}
