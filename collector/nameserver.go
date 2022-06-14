// Copyright 2015 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//go:build !nomdadm
// +build !nomdadm

package collector

import (
	"fmt"
	"github.com/docker/libnetwork/resolvconf"
	"github.com/docker/libnetwork/types"
	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
)

type NameServerCollector struct {
	logger log.Logger
}

func init() {
	registerCollector("nameserver", defaultEnabled, NewNameServerCollector)
}

// NewNameServerCollector returns a new Collector exposing raid statistics.
func NewNameServerCollector(logger log.Logger) (Collector, error) {
	return &NameServerCollector{logger}, nil
}

var (
	nameServers = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "nameserver", "info"),
		"Nameserver endpoints in /etc/resolv.conf",
		[]string{"nameserver"},
		nil,
	)
	
	searchDomains = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "searchdomains", "info"),
		"searchdomains endpoints in /etc/resolv.conf",
		[]string{"searchdomains"},
		nil,
	)
)

func (c *NameServerCollector) Update(ch chan<- prometheus.Metric) error {

	conf, err := resolvconf.Get()
	if err != nil {
		return fmt.Errorf("error obtaining resolv.conf: %w", err)
	}
	nameservers := resolvconf.GetNameservers(conf.Content, types.IPv4)
	for _, server := range nameservers {
		ch <- prometheus.MustNewConstMetric(
			nameServers,
			prometheus.GaugeValue,
			1,
			server,
		)
	}
	searchdomains := resolvconf.GetSearchDomains(conf.Content)

	for _, host := range searchdomains {

		ch <- prometheus.MustNewConstMetric(
			searchDomains,
			prometheus.GaugeValue,
			1,
			host,
		)
		}

	return nil
}
