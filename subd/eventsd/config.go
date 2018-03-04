package eventsd

import (
	"bytes"
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/virtengine/libgo/cmd"
	"github.com/virtengine/libgo/events"
	constants "github.com/virtengine/libgo/utils"
	"github.com/virtengine/vertice/meta"
)

type Config struct {
	Enabled bool    `toml:"enabled"`
	Mailgun Mailgun `toml:"mailgun"`
	Slack   Slack   `toml:"slack"`
	Infobip Infobip `toml:"infobip"`
	BillMgr BillMgr `toml:"bill"`
	Addons  Addons  `toml:"addons"`
}

func NewConfig() *Config {
	return &Config{
		Enabled: true,
	}
}

func (c Config) String() string {
	w := new(tabwriter.Writer)
	var b bytes.Buffer
	w.Init(&b, 0, 8, 0, '\t', 0)
	b.Write([]byte(cmd.Colorfy("Config:", "white", "", "bold") + "\t" +
		cmd.Colorfy("Eventsd", "cyan", "", "") + "\n"))
	b.Write([]byte(c.Mailgun.String()))
	b.Write([]byte(c.Infobip.String()))
	b.Write([]byte(c.Slack.String() + "\n"))
	b.Write([]byte(c.BillMgr.String() + "\n"))
	b.Write([]byte("---\n"))
	fmt.Fprintln(w)
	w.Flush()
	return strings.TrimSpace(b.String())
}

func (c Config) toMap() events.EventsConfigMap {
	em := make(events.EventsConfigMap)
	em[constants.MAILGUN] = c.Mailgun.toMap()
	em[constants.SLACK] = c.Slack.toMap()
	em[constants.INFOBIP] = c.Infobip.toMap()
	em[constants.BILLMGR] = c.BillMgr.toMap()
	em[constants.ADDONS] = c.Addons.toMap()
	em[constants.META] = meta.MC.ToMap()
	return em
}
