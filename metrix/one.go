package metrix

import (
	"encoding/xml"
	"github.com/virtengine/opennebula-go/metrics"
	"github.com/virtengine/vertice/carton"
	"io/ioutil"
	"strconv"
	"time"
	"fmt"
)

const (
	OPENNEBULA  = "one"
)

type OpenNebula struct {
	Url          string
	Region       string
	DefaultUnits map[string]string
	RawStatus    []byte
}

func (on *OpenNebula) Prefix() string {
	return "one"
}

func (on *OpenNebula) DeductBill(c *MetricsCollection) (e error) {
	for _, mc := range c.Sensors {
			mkBalance(mc, on.DefaultUnits)
	}
	return
}

func (on *OpenNebula) Collect(c *MetricsCollection) (e error) {
	b, e := on.ReadStatus()
	if e != nil {
		fmt.Println(e)
		return
	}

	s, e := on.ParseStatus(b)
	if e != nil {
		return
	}
	on.CollectMetricsFromStats(c, s)
	e = on.DeductBill(c)
	return
}

func (on *OpenNebula) ReadStatus() (b []byte, e error) {
	if len(on.RawStatus) == 0 {
		var res []interface{}
		res, e = carton.ProvisionerMap[on.Prefix()].MetricEnvs(time.Now().Add(-MetricsInterval).Unix(), time.Now().Unix(), on.Region, ioutil.Discard)
		if e != nil {
			return
		}
		on.RawStatus = []byte(res[0].(string))
	}

	b = on.RawStatus
	return
}

func (on *OpenNebula) ParseStatus(b []byte) (ons *metrics.OpenNebulaStatus, e error) {
	ons = &metrics.OpenNebulaStatus{}
	e = xml.Unmarshal(b, ons)
	if e != nil {
		return nil, e
	}
	return ons, nil
}

//actually the NewSensor can create trypes based on the event type.
func (on *OpenNebula) CollectMetricsFromStats(mc *MetricsCollection, s *metrics.OpenNebulaStatus) {
	for _, h := range s.History_Records {
	  usage, billable := on.vmQuota(h.QuotaId(),h.VCpu(),h.Memory(), h.Disks())
		if billable {
			sc := NewSensor(ONE_VM_SENSOR)
			sc.AccountId = h.AccountsId()
			sc.System = on.Prefix()
			sc.Node = h.HostName
			sc.AssemblyId = h.AssemblyId()
			sc.AssemblyName = h.AssemblyName()
			sc.AssembliesId = h.AssembliesId()
			sc.Source = on.Prefix()
			sc.Message = "vm billing"
			sc.Status = h.State()
			sc.AuditPeriodBeginning = time.Now().Add(-MetricsInterval).Format(time.RFC3339) //time.Unix(h.PStime, 0).String()
			sc.AuditPeriodEnding = time.Now().Format(time.RFC3339) // time.Unix(h.PEtime, 0).String()
			sc.AuditPeriodDelta = h.Elapsed()
			sc.addMetric(CPU_COST, h.CpuCost(), usage[metrics.CPU], "delta")
			sc.addMetric(MEMORY_COST, h.MemoryCost(), usage[metrics.MEMORY], "delta")
			sc.addMetric(DISK_COST, h.DiskCost(),usage[metrics.DISKS] , "delta")
			sc.CreatedAt = time.Now()
			if sc.isBillable() {
					mc.Add(sc)
			}
		}

	}
	return
}

func (on *OpenNebula) vmQuota(id, cpu, ram string, disks []metrics.Disk) (map[string]string, bool) {
  usage := make(map[string]string)
	var totalsize int64
	for _,v := range disks {
		totalsize = totalsize + v.Size
	}
	usage[metrics.CPU] = cpu
	usage[metrics.MEMORY] = ram
	usage[metrics.DISKS] = strconv.FormatInt(totalsize,10)

	if len(id) > 0 {
		if len(disks) > 1 {
			usage[metrics.CPU] = "0"
			usage[metrics.MEMORY] = "0"
			usage[metrics.DISKS] = strconv.FormatInt(totalsize - disks[0].Size, 10)
			return usage, true
		}
		return usage, false
	}

  return usage, true
}
