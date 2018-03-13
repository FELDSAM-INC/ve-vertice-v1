package run

import (
	"errors"

	"github.com/virtengine/vertice/meta"
	"github.com/virtengine/vertice/storage"
	"github.com/virtengine/vertice/subd/deployd"
	"github.com/virtengine/vertice/subd/dns"
	"github.com/virtengine/vertice/subd/docker"
	"github.com/virtengine/vertice/subd/eventsd"
	"github.com/virtengine/vertice/subd/httpd"
	"github.com/virtengine/vertice/subd/marketplacesd"
	"github.com/virtengine/vertice/subd/metricsd"
	"github.com/virtengine/vertice/subd/rancher"
)

type Config struct {
	Meta         *meta.Config          `toml:"meta"`
	Deployd      *deployd.Config       `toml:"deployd"`
	HTTPD        *httpd.Config         `toml:"http"`
	Docker       *docker.Config        `toml:"docker"`
	Metrics      *metricsd.Config      `toml:"metrics"`
	DNS          *dns.Config           `toml:"dns"`
	Events       *eventsd.Config       `toml:"events"`
	Storage      *storage.Config       `toml:"storage"`
	Rancher      *rancher.Config       `toml:"rancher"`
	MarketPlaces *marketplacesd.Config `toml:"marketplaces"`
}

func (c Config) String() string {
	return ("\n" +
		c.Meta.String() +
		c.Deployd.String() + "\n" +
		c.HTTPD.String() + "\n" +
		c.Docker.String() + "\n" +
		c.Metrics.String() + "\n" +
		c.DNS.String() + "\n" +
		c.Events.String() + "\n" +
		c.Storage.String() + "\n" +
		c.MarketPlaces.String() + "\n" +
		c.Rancher.String())

}

// NewConfig returns an instance of Config with reasonable defaults.
func NewConfig() *Config {
	c := &Config{}
	c.Meta = meta.NewConfig()
	c.Deployd = deployd.NewConfig()
	c.HTTPD = httpd.NewConfig()
	c.Docker = docker.NewConfig()
	c.Metrics = metricsd.NewConfig()
	c.Events = eventsd.NewConfig()
	c.DNS = dns.NewConfig()
	c.Storage = storage.NewConfig()
	c.Rancher = rancher.NewConfig()
	c.MarketPlaces = marketplacesd.NewConfig()
	return c
}

// NewDemoConfig returns the config that runs when no config is specified.
func NewDemoConfig() (*Config, error) {
	c := NewConfig()
	return c, nil
}

// Validate returns an error if the config is invalid.
func (c *Config) Validate() error {
	if c.Meta.Dir == "" {
		return errors.New("Meta.Dir must be specified")
	}
	return nil
}
