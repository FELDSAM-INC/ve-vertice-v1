package metricsd

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/virtengine/libgo/cmd"
	"github.com/virtengine/vertice/toml"
)

const (
	DefaultCollectInterval = 10 * time.Minute
)

type Config struct {
	Enabled         bool          `toml:"enabled"`
	CollectInterval toml.Duration `toml:"collect_interval"`
}

func NewConfig() *Config {
	return &Config{
		Enabled:         false,
		CollectInterval: toml.Duration(DefaultCollectInterval),
	}
}

func (c Config) String() string {
	w := new(tabwriter.Writer)
	var b bytes.Buffer
	w.Init(&b, 0, 8, 0, '\t', 0)
	b.Write([]byte(cmd.Colorfy("Config:", "white", "", "bold") + "\t" +
		cmd.Colorfy("Metricsd", "cyan", "", "") + "\n"))
	b.Write([]byte("enabled" + "\t" + strconv.FormatBool(c.Enabled) + "\n"))
	b.Write([]byte("collect_interval" + "\t" + c.CollectInterval.String() + "\n"))
	b.Write([]byte("---\n"))
	fmt.Fprintln(w)
	w.Flush()
	return strings.TrimSpace(b.String())
}
