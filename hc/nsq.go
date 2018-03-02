package hc

import (
	"fmt"

	nsq "github.com/crackcomm/nsqueue/producer"
	"github.com/virtengine/libgo/hc"
	"github.com/virtengine/vertice/meta"
)

func init() {
	hc.AddChecker("vertice:nsq", healthCheckNSQ)
}

func healthCheckNSQ() (interface{}, error) {
	if err := nsq.Connect(meta.MC.NSQd[0]); err == nil {
		return fmt.Sprintf("%s up", meta.MC.NSQd[0]), nil
	} else {
		return nil, err
	}
}
