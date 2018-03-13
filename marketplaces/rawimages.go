package marketplaces

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/virtengine/libgo/api"
	"github.com/virtengine/libgo/cmd"
	"github.com/virtengine/libgo/pairs"
	"github.com/virtengine/libgo/utils"
	lw "github.com/virtengine/libgo/writer"
	"github.com/virtengine/vertice/provision"
	"gopkg.in/yaml.v2"
)

const (
	APIRAWIMAGES        = "/rawimages"
	APIRAWIMAGES_UPDATE = "/rawimages/update"
)

type apiRawImages struct {
	JsonClaz string      `json:"json_claz"`
	Results  []RawImages `json:"results"`
}

type RawImages struct {
	Id         string          `json:"id"`
	Name       string          `json:"name"`
	AccountId  string          `json:"account_id"`
	OrgId      string          `json:"org_id"`
	Inputs     pairs.JsonPairs `json:"inputs"`
	Outputs    pairs.JsonPairs `json:"outputs"`
	Repository string          `json:"repos"`
	Repo       *Repo           `json:"-"`
	Status     string          `json:"status"`
	JsonClaz   string          `json:"json_claz"`
	CreatedAt  string          `json:"created_at"`
	UpdatedAt  string          `json:"updated_at"`
}

type Repo struct {
	Source    string `json:"source"`
	PublicUrl string `json:"public_url"`
	//	Properties []string `json:"properties"`
}

func (r Repo) toMap() map[string]string {
	m := make(map[string]string) //r.Properties.ToMap()
	m["source"] = r.Source
	m["public_url"] = r.PublicUrl
	return m
}

// marketplaces json string
func (s *RawImages) String() string {
	if d, err := yaml.Marshal(s); err != nil {
		return err.Error()
	} else {
		return string(d)
	}
}

func (r *RawImages) Get() (*RawImages, error) {
	return r.get()
}

/** A public function which pulls the snapshot for disk save as image.
and any others we do. **/
func (r *RawImages) get() (*RawImages, error) {
	cl := api.NewClient(newArgs(r.AccountId, ""), APIRAWIMAGES+"/"+r.Id)
	response, err := cl.Get()
	if err != nil {
		return nil, err
	}
	res := &apiRawImages{}
	err = json.Unmarshal(response, res)
	if err != nil {
		return nil, err
	}
	raw := &res.Results[0]
	repo := &Repo{}
	err = json.Unmarshal([]byte(raw.Repository), repo)
	if err != nil {
		return nil, err
	}
	raw.Repo = repo
	return raw, nil
}

func (r *RawImages) Update() error {
	return r.update()
}

func (r *RawImages) update() error {
	cl := api.NewClient(newArgs(r.AccountId, ""), APIRAWIMAGES_UPDATE)
	_, err := cl.Post(r)
	if err != nil {
		return err
	}
	return nil
}

// Deploy runs a deployment of an application. It will first try to run an
// image based deploy, and then fallback to the Git based deployment.
func (r *RawImages) create() error {
	var outBuffer bytes.Buffer
	start := time.Now()
	box := r.mkBox()
	log.Debugf("box  %v", box)
	logWriter := lw.LogWriter{Box: box}
	logWriter.Async()
	defer logWriter.Close()
	writer := io.MultiWriter(&outBuffer, &logWriter)
	err := r.deployToProvisioner(box, writer)
	elapsed := time.Since(start)
	if err != nil {
		return err
	}
	log.Debugf("%s in (%s)\n%s", cmd.Colorfy(r.Name, "cyan", "", "bold"),
		cmd.Colorfy(elapsed.String(), "green", "", "bold"),
		cmd.Colorfy(outBuffer.String(), "yellow", "", ""))
	return nil
}

func (r *RawImages) deployToProvisioner(box *provision.Box, writer io.Writer) error {
	if deployer, ok := ProvisionerMap[utils.PROVIDER_ONE].(provision.RawImageAccess); ok {
		return deployer.ISODeploy(box, writer)
	}
	return fmt.Errorf("cannot provision rawimages")
}

func (r *RawImages) mkBox() *provision.Box {
	return &provision.Box{
		CartonId:   r.Id,
		AccountId:  r.AccountId,
		Name:       r.Name,
		CartonName: r.Name,
		Region:     r.Region(),
		Provider:   r.provider(),
		PublicUrl:  r.publicUrl(),
	}
}

func (s *RawImages) Region() string {
	return s.Inputs.Match(utils.REGION)
}

func (a *RawImages) provider() string {
	return a.Inputs.Match(utils.PROVIDER)
}

func (a *RawImages) ImageId() string {
	return a.Outputs.Match(utils.RAW_IMAGE_ID)
}

func (r *RawImages) publicUrl() string {
	// if r.Repo != nil {
	// 	return r.Repo.PublicUrl
	// }
	return r.Repo.PublicUrl
}
