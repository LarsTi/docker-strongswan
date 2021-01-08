package main
import (
	"github.com/prometheus/client_golang/prometheus"
	"strings"
	"time"
	"./viciwrapper"
)
type StrongswanCollector struct {
	wrapper			*viciwrapper.ViciWrapper
	namespace		string

	viciCntCommands		*prometheus.Desc
	viciCntErrors		*prometheus.Desc
	viciLastCommandSec	*prometheus.Desc
	viciExecLastNanoSec	*prometheus.Desc
	viciExecAvgNanoSec	*prometheus.Desc
	viciLoadedSecrets	*prometheus.Desc

	ikeCnt			*prometheus.Desc
	ikeConnCnt		*prometheus.Desc
	ikeVersion		*prometheus.Desc
	ikeState		*prometheus.Desc
	ikeInitiator		*prometheus.Desc
	ikeNatRemote		*prometheus.Desc
	ikeNatFake		*prometheus.Desc
	ikeEncKeysize		*prometheus.Desc
	ikeIntegKeysize		*prometheus.Desc
	ikeEstablishSecs	*prometheus.Desc
	ikeRekeySecs		*prometheus.Desc
	ikeReauthSecs		*prometheus.Desc
	ikeChildren		*prometheus.Desc

	saState			*prometheus.Desc
	saEncap			*prometheus.Desc
	saEncKeysize		*prometheus.Desc
	saIntegKeysize		*prometheus.Desc
	saBytesIn		*prometheus.Desc
	saPacketsIn		*prometheus.Desc
	saLastInSecs		*prometheus.Desc
	saBytesOut		*prometheus.Desc
	saPacketsOut		*prometheus.Desc
	saLastOutSecs		*prometheus.Desc
	saEstablishSecs		*prometheus.Desc
	saRekeySecs		*prometheus.Desc
	saLifetimeSecs		*prometheus.Desc
}
func NewStrongswanCollector(w *viciwrapper.ViciWrapper) *StrongswanCollector {
	ns := "strongswan_"
	return &StrongswanCollector{
		wrapper: w,
		namespace: ns,

		viciCntCommands: prometheus.NewDesc(
			ns+"vici_command_count",
			"Number of totally send commands",
			nil, nil,
		),
		viciCntErrors: prometheus.NewDesc(
			ns+"vici_error_count",
			"Number of commands which returned an error",
			nil, nil,
		),
		viciLastCommandSec: prometheus.NewDesc(
			ns+"vici_seconds_since_last_command",
			"Seconds since last command was issued",
			nil, nil,
		),
		viciExecLastNanoSec: prometheus.NewDesc(
			ns+"vici_execution_nanoseconds_last",
			"Nanoseconds the last command execution took",
			nil, nil,
		),
		viciExecAvgNanoSec: prometheus.NewDesc(
			ns+"vici_execution_nanoseconds_avg",
			"Average nanoseconds for command execution during this vici session",
			nil, nil,
		),
		viciLoadedSecrets: prometheus.NewDesc(
			ns+"vici_loaded_secrets",
			"Loaded Secrets Number",
			nil, nil,
		),

		ikeCnt: prometheus.NewDesc(
			ns+"number_of_known_ikes",
			"Number of known IKEs",
			nil, nil,
		),
		ikeConnCnt: prometheus.NewDesc(
			ns+"number_of_ikes_connected",
			"Number of temporary connected IKEs",
			nil, nil,
		),
		ikeVersion: prometheus.NewDesc(
			ns+"ike_version",
			"Version Number of this IKE",
			[]string{"name"}, nil,
		),
		ikeState: prometheus.NewDesc(
			ns+"ike_state",
			"Status of this IKE",
			[]string{"name"}, nil,
		),
		ikeInitiator: prometheus.NewDesc(
			ns+"ike_initiator",
			"Flag if the server is the initiator for this connection",
			[]string{"name"}, nil,
		),
		ikeNatRemote: prometheus.NewDesc(
			ns+"ike_nat_remote",
			"Flag if the remote server is behind nat",
			[]string{"name"}, nil,
		),
		ikeNatFake: prometheus.NewDesc(
			ns+"ike_nat_fake",
			"Flag if the NAT is faked (to float to 4500)",
			[]string{"name"}, nil,
		),
		ikeEncKeysize: prometheus.NewDesc(
			ns+"ike_encryption_keysize",
			"Keysize of the encryption algorithm",
			[]string{"name", "algorithm", "dh_group"}, nil,
		),
		ikeIntegKeysize: prometheus.NewDesc(
			ns+"ike_integrity_keysize",
			"Keysize of the integrity algorithm",
			[]string{"name", "algorithm", "dh_group"}, nil,
		),
		ikeEstablishSecs: prometheus.NewDesc(
			ns+"ike_established_second",
			"Seconds since the IKE was established",
			[]string{"name"}, nil,
		),
		ikeRekeySecs: prometheus.NewDesc(
			ns+"ike_rekey_second",
			"Second count until the IKE will be rekeyed",
			[]string{"name"}, nil,
		),
		ikeReauthSecs: prometheus.NewDesc(
			ns+"ike_reauth_second",
			"Second count until the IKE will be reauthed",
			[]string{"name"}, nil,
		),
		ikeChildren: prometheus.NewDesc(
			ns+"ike_children",
			"Count of children of this IKE",
			[]string{"name"}, nil,
		),

		saState: prometheus.NewDesc(
			ns+"sa_state",
			"Status of this child sa",
			[]string{"ike_name","child_name","localTS","remoteTS"}, nil,
		),
		saEncap: prometheus.NewDesc(
			ns+"sa_encap",
			"Forced Encapsulation in UDP Packets",
			[]string{"ike_name", "child_name"}, nil,
		),
		saEncKeysize: prometheus.NewDesc(
			ns+"sa_encryption_keysize",
			"Keysize of the encryption algorithm",
			[]string{"ike_name", "child_name", "algorithm", "dh_group"}, nil,
		),
		saIntegKeysize: prometheus.NewDesc(
			ns+"sa_integrity_keysize",
			"Keysize of the integrity algorithm",
			[]string{"ike_name", "child_name", "algorithm", "dh_group"}, nil,
		),
		saBytesIn: prometheus.NewDesc(
			ns+"sa_bytes_inbound",
			"Number of bytes coming to the local server",
			[]string{"ike_name", "child_name", "localTS", "remoteTS"}, nil,
		),
		saPacketsIn: prometheus.NewDesc(
			ns+"sa_packets_inbound",
			"Number of packets coming to the local server",
			[]string{"ike_name", "child_name", "localTS", "remoteTS"}, nil,
		),
		saLastInSecs: prometheus.NewDesc(
			ns+"sa_last_inbound_seconds",
			"Number of seconds since the last inbound packet was received",
			[]string{"ike_name", "child_name", "localTS", "remoteTS"}, nil,
		),
		saBytesOut: prometheus.NewDesc(
			ns+"sa_bytes_outbound",
			"Number of bytes going to the remote server",
			[]string{"ike_name", "child_name", "localTS", "remoteTS"}, nil,
		),
		saPacketsOut: prometheus.NewDesc(
			ns+"sa_packets_outbound",
			"Number of packets going to the remote server",
			[]string{"ike_name", "child_name", "localTS", "remoteTS"}, nil,
		),
		saLastOutSecs: prometheus.NewDesc(
			ns+"sa_last_outbound_seconds",
			"Number of seconds since the last outbound packet was sent",
			[]string{"ike_name", "child_name", "localTS", "remoteTS"}, nil,
		),
		saEstablishSecs: prometheus.NewDesc(
			ns+"sa_established_second",
			"Seconds since the child SA was established",
			[]string{"ike_name", "name"}, nil,
		),
		saRekeySecs: prometheus.NewDesc(
			ns+"sa_rekey_second",
			"Second count until the child SA will be rekeyed",
			[]string{"ike_name", "name"}, nil,
		),
		saLifetimeSecs: prometheus.NewDesc(
			ns+"sa_lifetime_second",
			"Second count until the lifetime expires",
			[]string{"ike_name", "name"}, nil,
		),


	}
}
func (c *StrongswanCollector) init(){
	prometheus.MustRegister(c)
}
func (c *StrongswanCollector) Describe (ch chan<- *prometheus.Desc){
	ch <- c.viciCntCommands
	ch <- c.viciCntErrors
	ch <- c.viciLastCommandSec
	ch <- c.viciExecLastNanoSec
	ch <- c.viciExecAvgNanoSec
	ch <- c.viciLoadedSecrets

	ch <- c.ikeCnt
	ch <- c.ikeConnCnt
	ch <- c.ikeVersion
	ch <- c.ikeState
	ch <- c.ikeInitiator
	ch <- c.ikeNatRemote
	ch <- c.ikeNatFake
	ch <- c.ikeEncKeysize
	ch <- c.ikeIntegKeysize
	ch <- c.ikeEstablishSecs
	ch <- c.ikeRekeySecs
	ch <- c.ikeReauthSecs
	ch <- c.ikeChildren

	ch <- c.saState
	ch <- c.saEncap
	ch <- c.saEncKeysize
	ch <- c.saIntegKeysize
	ch <- c.saBytesIn
	ch <- c.saPacketsIn
	ch <- c.saLastInSecs
	ch <- c.saBytesOut
	ch <- c.saPacketsOut
	ch <- c.saLastOutSecs
	ch <- c.saEstablishSecs
	ch <- c.saRekeySecs
	ch <- c.saLifetimeSecs

}
func (c *StrongswanCollector) Collect (ch chan<- prometheus.Metric) {
	ch <- prometheus.MustNewConstMetric(
		c.ikeCnt, //Description
		prometheus.GaugeValue, //Type
		float64(c.wrapper.GetIkesInSystem()), //value
	)
	c.collectViciMetrics(ch)

	data, err := c.wrapper.ListIkes()
	if err != nil {
		ch <- prometheus.MustNewConstMetric(
			c.ikeConnCnt, //Description
			prometheus.GaugeValue, //Type
			float64(0), //Value
		)
		return
	}
	ch <- prometheus.MustNewConstMetric(
		c.ikeConnCnt, //Description
		prometheus.GaugeValue, //Type
		float64(len(data)), //Value
	)
	for _,v := range data {
		c.collectIkeMetrics(v, ch)
		for _, child := range v.Children {
			c.collectSaMetrics(v.Name,child, ch)
		}
	}
}
func (c *StrongswanCollector) collectViciMetrics(ch chan<- prometheus.Metric){
	metric := c.wrapper.GetViciMetrics()
	ch <- prometheus.MustNewConstMetric(
		c.viciCntCommands, //Description
		prometheus.CounterValue, //Type
		float64(metric.CounterCommands), //Value
	)
	ch <- prometheus.MustNewConstMetric(
		c.viciCntErrors, //Description
		prometheus.CounterValue, //Type
		float64(metric.CounterErrors), //Value
	)
	ch <- prometheus.MustNewConstMetric(
		c.viciLastCommandSec, //Description
		prometheus.GaugeValue, //Type
		float64(time.Since(metric.LastCommand).Seconds()), //Value
	)
	ch <- prometheus.MustNewConstMetric(
		c.viciExecLastNanoSec, //Description
		prometheus.GaugeValue, //Type
		float64(metric.ExecDuraLast.Nanoseconds()), //Value
	)
	ch <- prometheus.MustNewConstMetric(
		c.viciExecAvgNanoSec, //Description
		prometheus.GaugeValue, //Type
		float64(metric.ExecDuraAvgNs), //Value
	)
	ch <- prometheus.MustNewConstMetric(
		c.viciLoadedSecrets, //Description
		prometheus.GaugeValue, //Type
		float64(metric.LoadedSecrets),
	)
}
func (c *StrongswanCollector) collectIkeMetrics(d viciwrapper.LoadedIKE, ch chan<- prometheus.Metric){
	ch <- prometheus.MustNewConstMetric(
		c.ikeVersion, //Description
		prometheus.GaugeValue, //Type
		float64(d.Version), //Value
		d.Name, //Labels
	)
	
	state := 0
	if d.State == "ESTABLISHED" {
		state = 1
	}
	
	ch <- prometheus.MustNewConstMetric(
		c.ikeState, //Description
		prometheus.GaugeValue, //Type
		float64(state), //Value
		d.Name, //Labels
	)
	
	ch <- prometheus.MustNewConstMetric(
		c.ikeInitiator, //Description
		prometheus.GaugeValue, //Type
		float64(viciBoolToInt(d.Initiator)), //Value
		d.Name, //Labels
	)
	ch <- prometheus.MustNewConstMetric(
		c.ikeNatRemote, //Description
		prometheus.GaugeValue, //Type
		float64(viciBoolToInt(d.NatRemote)), //Value
		d.Name, //Labels
	)
	ch <- prometheus.MustNewConstMetric(
		c.ikeNatFake, //Description
		prometheus.GaugeValue, //Type
		float64(viciBoolToInt(d.NatFake)), //Value
		d.Name, //Labels
	)
	ch <- prometheus.MustNewConstMetric(
		c.ikeEncKeysize, //Description
		prometheus.GaugeValue, //Type
		float64(d.EncKey), //Value
		d.Name, d.EncAlg, d.DHGroup,//Labels
	)
	ch <- prometheus.MustNewConstMetric(
		c.ikeIntegKeysize, //Description
		prometheus.GaugeValue, //Type
		float64(d.IntegKey), //Value
		d.Name, d.IntegAlg, d.DHGroup,//Labels
	)
	ch <- prometheus.MustNewConstMetric(
		c.ikeEstablishSecs, //Description
		prometheus.GaugeValue, //Type
		float64(d.EstablishSec), //Value
		d.Name, //Labels
	)
	ch <- prometheus.MustNewConstMetric(
		c.ikeRekeySecs, //Description
		prometheus.GaugeValue, //Type
		float64(d.RekeySec), //Value
		d.Name, //Labels
	)
	ch <- prometheus.MustNewConstMetric(
		c.ikeReauthSecs, //Description
		prometheus.GaugeValue, //Type
		float64(d.ReauthSec), //Value
		d.Name, //Labels
	)
	ch <- prometheus.MustNewConstMetric(
		c.ikeChildren,
		prometheus.GaugeValue, //Type
		float64(len(d.Children)), //Value
		d.Name, //Labels
	)
}
func (c *StrongswanCollector) collectSaMetrics(name string ,d viciwrapper.LoadedChild, ch chan<- prometheus.Metric){
	state := 0
	if d.State == "ESTABLISHED" {
		state = 1
	}
	ch <- prometheus.MustNewConstMetric(
		c.saState, //Description
		prometheus.GaugeValue, //Type
		float64(state), //Value
		name, d.Name, strings.Join(d.LocalTS, ";"), strings.Join(d.RemoteTS, ";"), //Labels
	)
	ch <- prometheus.MustNewConstMetric(
		c.saEncap, //Description
		prometheus.GaugeValue, //Type
		float64(viciBoolToInt(d.Encap)), //Value
		name, d.Name, //Labels
	)
	ch <- prometheus.MustNewConstMetric(
		c.saEncKeysize, //Description
		prometheus.GaugeValue, //Type
		float64(d.EncKey), //Value
		name, d.Name, d.EncAlg, d.DHGroup, //Labels
	)
	ch <- prometheus.MustNewConstMetric(
		c.saIntegKeysize, //Description
		prometheus.GaugeValue, //Type
		float64(d.IntegKey), //Value
		name, d.Name, d.IntegAlg, d.DHGroup, //Labels
	)
	ch <- prometheus.MustNewConstMetric(
		c.saBytesIn, //Description
		prometheus.GaugeValue, //Type
		float64(d.BytesIn), //Value
		name, d.Name, strings.Join(d.LocalTS, ";"), strings.Join(d.RemoteTS, ";"), //Labels
	)
	ch <- prometheus.MustNewConstMetric(
		c.saPacketsIn, //Description
		prometheus.GaugeValue, //Type
		float64(d.PacketsIn), //Value
		name, d.Name, strings.Join(d.LocalTS, ";"), strings.Join(d.RemoteTS, ";"), //Labels
	)
	ch <- prometheus.MustNewConstMetric(
		c.saLastInSecs, //Description
		prometheus.GaugeValue, //Type
		float64(d.LastInSec), //Value
		name, d.Name, strings.Join(d.LocalTS, ";"), strings.Join(d.RemoteTS, ";"), //Labels
	)
	ch <- prometheus.MustNewConstMetric(
		c.saBytesOut, //Description
		prometheus.GaugeValue, //Type
		float64(d.BytesOut), //Value
		name, d.Name, strings.Join(d.LocalTS, ";"), strings.Join(d.RemoteTS, ";"), //Labels
	)
	ch <- prometheus.MustNewConstMetric(
		c.saPacketsOut, //Description
		prometheus.GaugeValue, //Type
		float64(d.PacketsOut), //Value
		name, d.Name, strings.Join(d.LocalTS, ";"), strings.Join(d.RemoteTS, ";"), //Labels
	)
	ch <- prometheus.MustNewConstMetric(
		c.saLastOutSecs, //Description
		prometheus.GaugeValue, //Type
		float64(d.LastOutSec), //Value
		name, d.Name, strings.Join(d.LocalTS, ";"), strings.Join(d.RemoteTS, ";"), //Labels
	)
	
	ch <- prometheus.MustNewConstMetric(
		c.saEstablishSecs, //Description
		prometheus.GaugeValue, //Type
		float64(d.EstablishSec), //Value
		name, d.Name, //Labels
	)
	ch <- prometheus.MustNewConstMetric(
		c.saRekeySecs, //Description
		prometheus.GaugeValue, //Type
		float64(d.RekeySec), //Value
		name, d.Name, //Labels
	)
	ch <- prometheus.MustNewConstMetric(
		c.saLifetimeSecs, //Description
		prometheus.GaugeValue, //Type
		float64(d.LifetimeSec), //Value
		name, d.Name, //Labels
	)
}
func viciBoolToInt(v string) int {
	if v == "yes" {
		return 1
	}else {
		return 0
	}
}