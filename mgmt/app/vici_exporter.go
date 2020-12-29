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
	execDuraLastMilliSec	*prometheus.Desc
	execDuraAvgMilliSec	*prometheus.Desc
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
		execDuraLastMilliSec: prometheus.NewDesc(
			ns+"execution_milliseconds_last",
			"Milliseconds the last command took to execute",
			nil, nil,
		),
		execDuraAvgMilliSec: prometheus.NewDesc(
			ns+"execution_milliseconds_avg",
			"Milliseconds the average command tooks to execute during this vici session",
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
	ch <- c.execDuraLastMilliSec
	ch <- c.execDuraAvgMilliSec
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
		c.execDuraLastMilliSec, //Description
		prometheus.GaugeValue, //Type
		float64(c.vici.execDuraLast.Milliseconds()), //Value
	)
	ch <- prometheus.MustNewConstMetric(
		c.execDuraAvgMilliSec, //Description
		prometheus.GaugeValue, //Type
		float64(c.vici.execDuraAvgMs), //Value
	)
}

