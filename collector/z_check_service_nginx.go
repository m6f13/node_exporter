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

// +build !notime

package collector

import (
	//"time"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"os/exec"
)

type checkServiceNginxStruct struct {
	desc   *prometheus.Desc
	logger log.Logger
}

func init() {
	registerCollector("z_check_service_nginx", defaultEnabled, NewMyCollectorNginx)
}

func NewMyCollectorNginx(logger log.Logger) (Collector, error) {
	return &checkServiceNginxStruct{
		desc: prometheus.NewDesc(
			namespace+"_z_check_service_nginx",
			"Check if service NGINX is running.",
			nil, nil,
		),
		logger: logger,
	}, nil
}

func (c *checkServiceNginxStruct) Update(ch chan<- prometheus.Metric) error {
	cmd := exec.Command("systemctl", "check", "nginx")
	out, err := cmd.CombinedOutput()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			level.Debug(c.logger).Log("msg", "systemctl finished with non-zero", "status", exitErr)
			ch <- prometheus.MustNewConstMetric(c.desc, prometheus.CounterValue, float64(0))
		} else {
			level.Debug(c.logger).Log("msg", "failed to run systemctl", "status", err)
			ch <- prometheus.MustNewConstMetric(c.desc, prometheus.CounterValue, float64(0))
		}
	} else {
		level.Debug(c.logger).Log("msg", "NGINX status:", "status", out)
		ch <- prometheus.MustNewConstMetric(c.desc, prometheus.CounterValue, float64(1))
	}
	return nil
}
