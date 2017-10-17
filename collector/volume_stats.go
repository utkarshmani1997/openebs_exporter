package collector

import (
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

type volumeStats []struct {
	val func(*volume) string
	vec *prometheus.GaugeVec
}

// ChannelStats creates a new stats collector which is able to
// expose the channel metrics of a openebs exporter node to Prometheus.
// The channel metrics are reported per topic.
func VolumeStats(namespace string) StatsCollector {
	labels := []string{"OpenEBS"}
	namespace += "_openebs"
	return volumeStats{
		{
			val: func(c *volume) string { return string(c.RevisionCounter) },
			vec: prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "revision_counter",
				Help:      "Number of jiva volumes",
			}, labels),
		},
		{
			val: func(c *volume) string { return string(c.ReplicaCounter) },
			vec: prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "replica_counter",
				Help:      "Number of Jiva Replicas",
			}, labels),
		},
		/*		{
					val: func(c *rest.VolumeStats) string {
						for key, value := range c.SCSIIOCount {
							return string(value)
						}
					},
					vec: prometheus.NewGaugeVec(prometheus.GaugeOpts{
						Namespace: namespace,
						Name:      "scsiio_count",
						Help:      "SCSI Input/Output count",
					}, labels),
				},
		*/{
			val: func(c *volume) string { return string(c.ReadIOPS) },
			vec: prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "read_iops",
				Help:      "Read Input/Output of volumes",
			}, labels),
		},
		{
			val: func(c *volume) string { return string(c.TotalReadTime) },
			vec: prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "total_read_time",
				Help:      "Total Read Time of volumes",
			}, labels),
		},
		{
			val: func(c *volume) string { return string(c.TotalReadBlockCount) },
			vec: prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "total_read_block_count",
				Help:      "Total Read Block Count of volumes",
			}, labels),
		},
		{
			val: func(c *volume) string { return string(c.WriteIOPS) },
			vec: prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "write_iops",
				Help:      "Write Input/Output of volumes",
			}, labels),
		},
		{
			val: func(c *volume) string { return string(c.TotalWriteTime) },
			vec: prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "total_write_time",
				Help:      "Total write time of volume",
			}, labels),
		},
		{
			val: func(c *volume) string { return string(c.TotalWriteBlockCount) },
			vec: prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "total_write_block_count",
				Help:      "Total Write Block Count of volumes",
			}, labels),
		},
		{
			val: func(c *volume) string { return string(c.UsedLogicalBlocks) },
			vec: prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "used_logical_blocks",
				Help:      "Used Logical Blocks of volumes",
			}, labels),
		},
		{
			val: func(c *volume) string { return string(c.UsedBlocks) },
			vec: prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "used_blocks",
				Help:      "Used Blocks of volumes",
			}, labels),
		},
		{
			val: func(c *volume) string { return string(c.SectorSize) },
			vec: prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "sector_size",
				Help:      "Sector Size of volumes",
			}, labels),
		},
	}
}

func (cs volumeStats) set(s *stats) {
	for _, stat := range s.VolStats {
		labels := prometheus.Labels{
			"OpenEBS": "Jiva",
		}
		for _, c := range cs {
			v, _ := strconv.ParseFloat(c.val(stat), 64)
			c.vec.With(labels).Set(v)
		}
	}
}

func (cs volumeStats) collect(out chan<- prometheus.Metric) {
	for _, c := range cs {
		c.vec.Collect(out)
	}
}

func (cs volumeStats) describe(ch chan<- *prometheus.Desc) {
	for _, c := range cs {
		c.vec.Describe(ch)
	}
}

func (cs volumeStats) reset() {
	for _, c := range cs {
		c.vec.Reset()
	}
}
