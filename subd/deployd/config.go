package deployd

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/virtengine/libgo/cmd"
	constants "github.com/virtengine/libgo/utils"
	"github.com/virtengine/opennebula-go/api"
	"github.com/virtengine/vertice/provision/one"
)

const (
	// Default provisioning provider for vms is OpenNebula.
	// This is just an endpoint for Megam. We could have openstack, chef, salt, puppet etc.
	DefaultProvider = "one"
	DefaultImage    = "megam"

	DefaultCpuThrottle = "1"
	// DefaultOneEndpoint is the default address that the service binds to an IaaS (OpenNebula).
	DefaultOneEndpoint = "http://localhost:2633/RPC2"

	// DefaultUserid the default userid for the IaaS service (OpenNebula).
	DefaultOneUserid = "oneadmin"

	// DefaultOnePassword is the default password for the IaaS service (OpenNebula).
	DefaultOnePassword = "password"

	// DefaultOneTemplate is the default template for the IaaS service (OpenNebula).
	DefaultOneTemplate = "megam"

	// DefaultOneZone is the default master zone for the IaaS service (OpenNebula).
	DefaultOneMasterZone = "OpenNebula"

	// DefaultOneZone is the default zone for the IaaS service (OpenNebula).
	DefaultOneZone = "africa"

	//DefaultOneCluster is the default cluster for Host in the Iaas service (OpenNebula)
	DefaultOneCluster = "cluster-a"

	//DefaultOneCluster is the default cluster for Host in the Iaas service (OpenNebula)
	DefaultOneVnetPri = "vnet-pri"

	//DefaultOneCluster is the default cluster for Host in the Iaas service (OpenNebula)
	DefaultOneVnetPub = "vnet-pub"

	ONEZONE = "zone"
)

type Config struct {
	Provider string  `json:"provider" toml:"provider"`
	One      one.One `json:"one" toml:"one"`
}

/*
type deployd struct {

}
*/

func NewConfig() *Config {
	cl := make([]one.Cluster, 2)
	rg := make([]one.Region, 2)

	c := one.Cluster{
		Enabled:       false,
		StorageType:   "hdd",
		ClusterId:     DefaultOneCluster,
		Vnet_pri_ipv4: []string{DefaultOneVnetPri},
		Vnet_pub_ipv4: []string{DefaultOneVnetPub},
		Vnet_pri_ipv6: []string{DefaultOneVnetPri},
		Vnet_pub_ipv6: []string{DefaultOneVnetPri},
	}

	r := one.Region{
		OneZone:        DefaultOneZone,
		OneEndPoint:    DefaultOneEndpoint,
		OneUserid:      DefaultOneUserid,
		OnePassword:    DefaultOnePassword,
		OneTemplate:    DefaultOneCluster,
		Datastore:      "1",
		Certificate:    "/var/lib/megam/vertice/id_rsa.pub",
		Image:          DefaultImage,
		VCPUPercentage: "",
		Clusters:       append(cl, c),
	}
	o := one.One{
		Enabled:        true,
		Regions:        append(rg, r),
		OneTemplate:    DefaultOneTemplate,
		Image:          DefaultImage,
		VCPUPercentage: DefaultCpuThrottle,
	}

	return &Config{
		Provider: DefaultProvider,
		One:      o,
	}
}

func (c Config) String() string {
	w := new(tabwriter.Writer)
	var b bytes.Buffer
	w.Init(&b, 0, 8, 0, '\t', 0)
	b.Write([]byte(cmd.Colorfy("\nConfig:", "white", "", "bold") + "\t" +
		cmd.Colorfy("Deployd", "cyan", "", "") + "\n"))
	b.Write([]byte(constants.PROVIDER + "\t" + c.Provider + "\n"))
	b.Write([]byte("enabled      " + "\t" + strconv.FormatBool(c.One.Enabled) + "\n"))
	for _, v := range c.One.Regions {
		b.Write([]byte(api.ONEZONE + "\t" + v.OneZone + "\n"))
		b.Write([]byte(api.ENDPOINT + "\t" + v.OneEndPoint + "\n"))
		b.Write([]byte(api.USERID + "    \t" + v.OneUserid + "\n"))
		b.Write([]byte(api.PASSWORD + "\t" + v.OnePassword + "\n"))
		b.Write([]byte(api.TEMPLATE + "\t" + v.OneTemplate + "\n"))
		b.Write([]byte(api.IMAGE + "    \t" + v.Image + "\n"))
		b.Write([]byte(api.VCPU_PERCENTAGE + "\t" + v.VCPUPercentage + "\n"))
		for _, k := range v.Clusters {
			if k.Enabled {
				b.Write([]byte(api.CLUSTER + "\t" + k.ClusterId + "  storage type" + k.StorageType + "\n"))
			}
		}
		b.Write([]byte("---\n"))
	}
	fmt.Fprintln(w)
	w.Flush()
	return strings.TrimSpace(b.String())
}

//convert the config to just an interface.
func (c Config) ToInterface() interface{} {
	return c.One
}
