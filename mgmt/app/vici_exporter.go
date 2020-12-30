package main
import (
	"github.com/prometheus/client_golang/prometheus"
	"time"
)
type ViciCollector struct {
	vici			*viciStruct
	namespace		string
	cntCommands		*prometheus.Desc
	cntErrors		*prometheus.Desc
	lastCommandSec		*prometheus.Desc
	execDuraLastNanoSec	*prometheus.Desc
	execDuraAvgNanoSec	*prometheus.Desc
}
func NewViciCollector(v *viciStruct) *ViciCollector {
	ns := "vici_"
	return &ViciCollector{
		vici: v,
		namespace: ns,
		cntCommands: prometheus.NewDesc(
			ns+"command_count",
			"Number of totally send commands",
			nil, nil,
		),
		cntErrors: prometheus.NewDesc(
			ns+"error_count",
			"Number of commands which returned an error",
			nil, nil,
		),
		lastCommandSec: prometheus.NewDesc(
			ns+"seconds_since_last_command",
			"Time ellapsed since last command was issued in seconds",
			nil, nil,
		),
		execDuraLastNanoSec: prometheus.NewDesc(
			ns+"execution_nanoseconds_last",
			"Nanoseconds the last command took to execute",
			nil, nil,
		),
		execDuraAvgNanoSec: prometheus.NewDesc(
			ns+"execution_nanoseconds_avg",
			"Nanoseconds the average command tooks to execute during this vici session",
			nil, nil,
		),
	}
}
func (c *ViciCollector) init(){
	prometheus.MustRegister(c)
}
func (c *ViciCollector) Describe (ch chan<- *prometheus.Desc){
	ch <- c.cntCommands
	ch <- c.cntErrors
	ch <- c.lastCommandSec
	ch <- c.execDuraLastNanoSec
	ch <- c.execDuraAvgNanoSec
}
func (c *ViciCollector) Collect (ch chan<- prometheus.Metric) {
	ch <- prometheus.MustNewConstMetric(
		c.cntCommands, //Description
		prometheus.CounterValue, //Type
		float64(c.vici.counterCommands), //Value
	)
	ch <- prometheus.MustNewConstMetric(
		c.cntErrors, //Description
		prometheus.CounterValue, //Type
		float64(c.vici.counterErrors), //Value
	)

	ch <- prometheus.MustNewConstMetric(
		c.lastCommandSec, //Description
		prometheus.GaugeValue, //Type
		float64(time.Since(c.vici.lastCommand).Seconds()), //Value
	)
	ch <- prometheus.MustNewConstMetric(
		c.execDuraLastNanoSec, //Description
		prometheus.GaugeValue, //Type
		float64(c.vici.execDuraLast.Nanoseconds()), //Value
	)
	ch <- prometheus.MustNewConstMetric(
		c.execDuraAvgNanoSec, //Description
		prometheus.GaugeValue, //Type
		float64(c.vici.execDuraAvgMs), //Value
	)
}

