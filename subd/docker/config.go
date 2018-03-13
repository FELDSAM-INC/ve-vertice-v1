package docker

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/virtengine/libgo/cmd"
	constants "github.com/virtengine/libgo/utils"
	"github.com/virtengine/vertice/provision/docker"
	"github.com/virtengine/vertice/provision/docker/cluster"
	"github.com/virtengine/vertice/toml"
)

const (

	// DefaultRegistry is the default registry for docker (public)
	DefaultRegistry = "https://hub.docker.com"

	DefaultProvider = "docker"

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
	DefaultDockerZone = "africa"
	//DefaultOneZone is the default bridge for the docker.
	DefaultBridgeName = "eth0"

	DefaultName = "eth0"

	DefaultGulpPort = ":6666"
	DefaultNetType  = "cluster-a"

	// DefaultSwarmEndpoint is the default address that the service binds to an IaaS (Swarm).
	DefaultSwarmEndpoint = "tcp://localhost:2375"
)

type Config struct {
	Provider string        `json:"provider" toml:"provider"`
	Docker   docker.Docker `json:"docker" toml:"docker"`
}

func NewConfig() *Config {
	rg := make([]docker.Region, 0)
	r := docker.Region{
		DockerZone:    DefaultDockerZone,
		SwarmEndPoint: DefaultSwarmEndpoint,
		Registry:      DefaultRegistry,
		CPUPeriod:     toml.Duration(DefaultCPUPeriod),
		CPUQuota:      toml.Duration(DefaultCPUQuota),
	}

	o := docker.Docker{
		Enabled: true,
		Regions: append(rg, r),
	}
	return &Config{
		Provider: DefaultProvider,
		Docker:   o,
	}
}

func (c Config) String() string {
	w := new(tabwriter.Writer)
	var b bytes.Buffer
	w.Init(&b, 0, 8, 0, '\t', 0)
	b.Write([]byte(cmd.Colorfy("Config:", "white", "", "bold") + "\t" +
		cmd.Colorfy("docker", "cyan", "", "") + "\n"))
	b.Write([]byte(constants.PROVIDER + "\t" + c.Provider + "\n"))
	b.Write([]byte("enabled      " + "\t" + strconv.FormatBool(c.Docker.Enabled) + "\n"))
	for _, v := range c.Docker.Regions {
		b.Write([]byte(cluster.DOCKER_ZONE + "\t" + v.DockerZone + "\n"))
		b.Write([]byte(cluster.DOCKER_SWARM + "\t" + v.SwarmEndPoint + "\n"))
		b.Write([]byte(cluster.DOCKER_CPUPERIOD + "    \t" + v.CPUPeriod.String() + "\n"))
		b.Write([]byte(cluster.DOCKER_CPUQUOTA + "    \t" + v.CPUQuota.String() + "\n"))
		b.Write([]byte("---\n"))
	}
	fmt.Fprintln(w)
	w.Flush()
	return strings.TrimSpace(b.String())
}

func (c Config) toInterface() interface{} {
	return c.Docker
}
