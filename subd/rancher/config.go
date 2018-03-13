package rancher

import (
	"bytes"
	"fmt"
	"github.com/virtengine/libgo/cmd"
	constants "github.com/virtengine/libgo/utils"
	"github.com/virtengine/vertice/provision/rancher"
	"github.com/virtengine/vertice/provision/rancher/cluster"
	"github.com/virtengine/vertice/toml"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"
)

const (

	// DefaultRegistry is the default registry for docker (public)
	DefaultRegistry = "https://hub.docker.com"

	DefaultProvider = "rancher"

	DefaultRancherZone = "India"
	// DefaultNamespace is the default highlevel namespace(userid) under the registry eg: https://hub.docker.com/megam
	DefaultNamespace = "megam"

	// DefaultMemSize is the default memory size in MB used for every container launch
	DefaultMemSize = 256 * 1024 * 1024

	// DefaultSwapSize is the default memory size in MB used for every container launch
	DefaultSwapSize = 210 * 1024 * 1024

	// DefaultCPUPeriod is the default cpu period used for every container launch in ms
	DefaultCPUPeriod = 25000 * time.Millisecond

	// DefaultCPUQuota is the default cpu quota allocated for every cpu cycle for the launched container in ms
	DefaultCPUQuota = 25000 * time.Millisecond

	// DefaultOneZone is the default zone for the IaaS service.
	// Access credentials for radosgw
	DefaultAccessKey = "vertadmin"
	DefaultSecretKey = "vertadmin"

	// DefaultSwarmEndpoint is the default address that the service binds to an IaaS (Swarm).
	DefaultRancherEndpoint = "http://localhost:8080"
)

type Config struct {
	Provider string          `json:"provider" toml:"provider"`
	Rancher  rancher.Rancher `json:"container" toml:"container"`
}

func NewConfig() *Config {

	rg := make([]rancher.Region, 0)
	r := rancher.Region{
		RancherZone:     DefaultRancherZone,
		RancherEndPoint: DefaultRancherEndpoint,
		Registry:        DefaultRegistry,
		AdminId:         "info@megam.io",
		AdminAccess:     DefaultAccessKey,
		AdminSecret:     DefaultSecretKey,
		CPUPeriod:       toml.Duration(DefaultCPUPeriod),
		CPUQuota:        toml.Duration(DefaultCPUQuota),
	}

	o := rancher.Rancher{
		Enabled: true,
		Regions: append(rg, r),
	}

	return &Config{
		Provider: DefaultProvider,
		Rancher:  o,
	}
}

func (c Config) String() string {
	w := new(tabwriter.Writer)
	var b bytes.Buffer
	w.Init(&b, 0, 8, 0, '\t', 0)
	b.Write([]byte(cmd.Colorfy("Config:", "white", "", "bold") + "\t" +
		cmd.Colorfy("rancher", "cyan", "", "") + "\n"))
	b.Write([]byte(constants.PROVIDER + "\t" + c.Provider + "\n"))
	b.Write([]byte("enabled      " + "\t" + strconv.FormatBool(c.Rancher.Enabled) + "\n"))
	for _, v := range c.Rancher.Regions {
		b.Write([]byte(cluster.RANCHER_ZONE + "\t" + v.RancherZone + "\n"))
		b.Write([]byte(cluster.RANCHER_SERVER + "\t" + v.RancherEndPoint + "\n"))
		b.Write([]byte("Admin Id    \t" + v.AdminId + "\n"))
		b.Write([]byte("AdminAccess" + "    \t" + v.AdminAccess + "\n"))
		b.Write([]byte("AdminSecret" + "\t" + v.AdminSecret + "\n"))
		b.Write([]byte(cluster.RANCHER_CPUPERIOD + "    \t" + v.CPUPeriod.String() + "\n"))
		b.Write([]byte(cluster.RANCHER_CPUQUOTA + "    \t" + v.CPUQuota.String() + "\n"))
		b.Write([]byte("---\n"))
	}
	fmt.Fprintln(w)
	w.Flush()
	return strings.TrimSpace(b.String())
}

func (c Config) toInterface() interface{} {
	return c.Rancher
}
