/*
** Copyright [2013-2017] [Megam Systems]
**
** Licensed under the Apache License, Version 2.0 (the "License");
** you may not use this file except in compliance with the License.
** You may obtain a copy of the License at
**
** http://www.apache.org/licenses/LICENSE-2.0
**
** Unless required by applicable law or agreed to in writing, software
** distributed under the License is distributed on an "AS IS" BASIS,
** WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
** See the License for the specific language governing permissions and
** limitations under the License.
 */
package carton

import (
	log "github.com/Sirupsen/logrus"
	"github.com/virtengine/libgo/cmd"
)

// ReqOperator is a operator which will be used to execute the given request.
type ReqOperator struct {
	Id        string
	CartonsId string
	AccountId string
	Category  string
	Action    string
}

// NewReqOperator returns a new instance of ReqOperator
// for the operatable id (Assemblies)
func NewReqOperator(r *Requests) *ReqOperator {
	return &ReqOperator{CartonsId: r.CatId, Category: r.Category, Action: r.Action, AccountId: r.AccountId}
}

// Accept will run the given request
// it will call the given processor on the Cartons retrieved from Vertice API.
func (p *ReqOperator) Accept(r *MegdProcessor) error {
	c, err := p.Get()
	if err != nil {
		return err
	}
	md := *r
	log.Debugf(cmd.Colorfy(md.String(), "cyan", "", "bold"))
	return md.Process(c)
}

// Get will create Cartons based on the given category
// For categories `BACKUP`, `DISKS` and `SNAPSHOT`, the result list will contains a single Carton;
// For other categories, the number of Cartons will depend on the result from Vertice API call.
func (p *ReqOperator) Get() (Cartons, error) {
	switch p.Category {
	case BACKUPS:
		b, err := GetBackup(p.CartonsId, p.AccountId)
		if err != nil {
			return nil, err
		}
		c, err := b.MkCartons()
		if err != nil {
			return nil, err
		}
		return c, nil

	case DISKS:
		d, err := GetDisks(p.CartonsId, p.AccountId)
		if err != nil {
			return nil, err
		}
		c, err := d.MkCartons()
		if err != nil {
			return nil, err
		}
		return c, nil

	case SNAPSHOT:
		s, err := GetSnap(p.CartonsId, p.AccountId)
		if err != nil {
			return nil, err
		}
		c, err := s.MkCartons()
		if err != nil {
			return nil, err
		}
		return c, nil
	default:
		a, err := Get(p.CartonsId, p.AccountId)
		if err != nil {
			return nil, err
		}
		c, err := a.MkCartons(p.AccountId)
		if err != nil {
			return nil, err
		}
		return c, nil
	}
}

// MegdProcessor represents a single operation in vertice.
type MegdProcessor interface {
	Process(c Cartons) error
	String() string
}
