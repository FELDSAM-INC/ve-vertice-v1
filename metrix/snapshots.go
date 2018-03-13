package metrix

import (
	"time"

	"github.com/virtengine/vertice/carton"
)

const (
	SNAPSHOTS       = "snapshots"
	SNAPSHOT_SENSOR = "instance.snapshots.exists"
)

type Snapshots struct {
	DefaultUnits map[string]string
	RawStatus    []byte
}

func (r *Snapshots) Prefix() string {
	return SNAPSHOTS
}

func (r *Snapshots) DeductBill(c *MetricsCollection) (e error) {
	for _, mc := range c.Sensors {
		mkBalance(mc, r.DefaultUnits)
	}
	return
}

func (s *Snapshots) Collect(c *MetricsCollection) (e error) {
	snp := carton.Snaps{}
	snps, e := snp.GetBox()
	if e != nil {
		return
	}

	s.CollectMetricsFromStats(c, snps)
	e = s.DeductBill(c)
	return
}

func (c *Snapshots) ReadUsers() ([]*carton.Account, error) {
	act := new(carton.Account)
	res, e := act.GetUsers()
	if e != nil {
		return nil, e
	}
	return res, nil
}

//actually the NewSensor can create trypes based on the event type.
func (c *Snapshots) CollectMetricsFromStats(mc *MetricsCollection, snps []carton.Snaps) {
	for _, a := range snps {
		if !a.IsQuota() && a.IsAlive() {
			sc := NewSensor(SNAPSHOT_SENSOR)
			sc.AccountId = a.AccountId
			sc.AssemblyId = a.AssemblyId
			sc.System = c.Prefix()
			sc.Node = ""
			sc.AssemblyName = a.Name
			sc.AssembliesId = a.Id
			sc.Source = c.Prefix()
			sc.Message = "snapshot billing"
			sc.Status = "health-ok"
			sc.AuditPeriodBeginning = time.Now().Add(-MetricsInterval).Format(time.RFC3339)
			sc.AuditPeriodEnding = time.Now().Format(time.RFC3339)
			sc.AuditPeriodDelta = ""
			sc.addMetric(STORAGE_COST, c.DefaultUnits[STORAGE_COST_PER_HOUR], a.Sizeof(), "delta")
			sc.CreatedAt = time.Now()
			mc.Add(sc)
		}
	}

	return
}
