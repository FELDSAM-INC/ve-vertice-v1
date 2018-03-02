package metrix

import (
	log "github.com/Sirupsen/logrus"
	"github.com/virtengine/libgo/api"
	"github.com/virtengine/libgo/events"
	"github.com/virtengine/libgo/events/alerts"
	constants "github.com/virtengine/libgo/utils"
	"github.com/virtengine/vertice/carton"
	"time"

)

func SendMetricsToScylla(metrics Sensors, hostname string) (err error) {
	started := time.Now()
	for _, m := range metrics {
			cl := api.NewClient(carton.NewArgs(m.AccountId, ""), "/sensors/content")
			if _, err := cl.Post(m); err != nil {
				log.Debugf(err.Error())
				continue
			}
	}
	log.Debugf("sent %d metrics in %.06f\n", len(metrics), time.Since(started).Seconds())
	return nil
}


func mkBalance(s *Sensor, du map[string]string) error {
	mi := make(map[string]string)
	m := s.Metrics.Totalcost(du)
	mi[constants.ACCOUNTID] = s.AccountId
	mi[constants.ASSEMBLYID] = s.AssemblyId
	mi[constants.ASSEMBLYNAME] = s.AssemblyName
	mi[constants.CONSUMED] = m
	mi[constants.START_TIME] = s.AuditPeriodBeginning
	mi[constants.END_TIME] = s.AuditPeriodEnding

	newEvent := events.NewMulti(
		[]*events.Event{
			&events.Event{
				AccountsId:  s.AccountId,
				EventAction: alerts.DEDUCT,
				EventType:   constants.EventBill,
				EventData:   alerts.EventData{M: mi},
				Timestamp:   time.Now().Local(),
			},
			&events.Event{
				AccountsId:  s.AccountId,
				EventAction: alerts.BILLEDHISTORY, //Change type to transaction
				EventType:   constants.EventBill,
				EventData:   alerts.EventData{M: mi},
				Timestamp:   time.Now().Local(),
			},
		})
	return newEvent.Write()
}
