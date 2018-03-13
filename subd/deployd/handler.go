package deployd

import (
	log "github.com/Sirupsen/logrus"
	"github.com/virtengine/vertice/carton"
)

type Handler struct {
	d            *Config
	EventChannel chan bool
}

func NewHandler(c *Config) *Handler {
	return &Handler{d: c}

}

func (h *Handler) serveNSQ(r *carton.Requests) error {
	p, err := carton.ParseRequest(r)
	if err != nil {
		return err
	}
	if rp := carton.NewReqOperator(r); rp != nil {
		err = rp.Accept(&p)
		if err != nil {
			log.Errorf("Error Request : %s  -  %s  : %s", r.Category, r.Action, err)
		}

		return err //error is swalled here.
	}

	return nil
}
