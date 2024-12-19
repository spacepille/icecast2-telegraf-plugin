package icecast2

import (
	"runtime/debug"

	"golang.org/x/exp/slices"

	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/inputs"
	"github.com/rs/zerolog/log"
)

type Icecast2InputPlugin struct {
	pluginVersion string
}

func init() {
	inputs.Add("icecast2", func() telegraf.Input {
		return &Icecast2InputPlugin{
			pluginVersion: PluginVersion(),
		}
	})
}

func (input *Icecast2InputPlugin) Init() error {
	return nil
}

func (input *Icecast2InputPlugin) Stop() {
	//shmem.UnlockMutex()
}

func (input *Icecast2InputPlugin) SampleConfig() string {
	return `
[[inputs.icecast2]]
	# no config
`
}

func (input *Icecast2InputPlugin) Description() string {
	return "Gather information from icecast2 server"
}

func (input *Icecast2InputPlugin) Gather(a telegraf.Accumulator) error {
	log.Debug().Msg("Gathering metrics...")

	log.Debug().Msg("Done gathering metrics")
	return nil
}

func PluginVersion() string {
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		return "unknown"
	}

	i := slices.IndexFunc(bi.Deps, func(module *debug.Module) bool {
		return module.Path == "github.com/spacepille/icecast2-telegraf-plugin"
	})
	if i == -1 {
		return "unknown"
	}
	return bi.Deps[i].Version
}
