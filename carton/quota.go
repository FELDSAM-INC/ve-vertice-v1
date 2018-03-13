package carton

import (
	"encoding/json"
	"github.com/virtengine/libgo/api"
	"github.com/virtengine/libgo/pairs"
)

type Quota struct {
	Id          string          `json:"id" cql:"id"`
	AccountId   string          `json:"account_id" cql:"account_id"`
	Name        string          `json:"name" cql:"name"`
	JsonClaz    string          `json:"json_claz" cql:"json_claz"`
	Allowed     pairs.JsonPairs `json:"allowed" cql:"allowed"`
	AllocatedTo string          `json:"allocated_to" cql:"allocated_to"`
	QuotaType   string          `json:"quota_type" cql:"quota_type"`
	Status      string          `json:"status" cql:"status"`
	Inputs      pairs.JsonPairs `json:"inputs" cql:"inputs"`
}

type ApiQuota struct {
	JsonClaz string  `json:"json_claz"`
	Results  []Quota `json:"results"`
}

func (q *Quota) Update() error {
	return q.update(newArgs(q.AccountId, ""))
}

func (q *Quota) update(args api.ApiArgs) error {
	cl := api.NewClient(args, "/quotas/update")
	_, err := cl.Post(q)
	if err != nil {
		return err
	}
	return nil
}

func NewQuota(accountid, id string) (*Quota, error) {
	q := new(Quota)
	q.AccountId = accountid
	q.Id = id
	return q.get(newArgs(accountid, ""))
}

func (q *Quota) get(args api.ApiArgs) (*Quota, error) {
	cl := api.NewClient(args, "/quotas/"+q.Id)
	response, err := cl.Get()
	if err != nil {
		return nil, err
	}
	ac := &ApiQuota{}
	//log.Debugf("Response %s :  (%s)",cmd.Colorfy("[Body]", "green", "", "bold"),string(response))
	err = json.Unmarshal(response, ac)
	if err != nil {
		return nil, err
	}

	return &ac.Results[0], nil
}

func (q *Quota) ContainerQuota() (bool, error) {
	asm, err := NewAssembly(q.AllocatedTo, q.AccountId, "")
	if err != nil {
		return true, err
	}
	return !(len(asm.quotaID()) > 0), nil
}

func (q *Quota) AllowedSnaps() string {
	return q.Allowed.Match("no_of_units")
}
