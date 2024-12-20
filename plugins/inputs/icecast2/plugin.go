package icecast2

import (
	"runtime/debug"
	"time"

	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/config"
	"github.com/influxdata/telegraf/plugins/inputs"
	"github.com/rs/zerolog/log"
	"github.com/spacepille/icecast2-telegraf-plugin"
)

type IcecastInputPlugin struct {
	Url              string          `toml:"url"`
	Username         string          `toml:"username"`
	Password         string          `toml:"password"`
	ResponseTimeout  config.Duration `toml:"response_timeout"`
	CollectListeners bool            `toml:"collect_listeners"`
	Geoip2Path       string          `toml:"geoip2_path"`

	pluginVersion   string
	iceastCollector *IceastCollector
}

func init() {
	inputs.Add("icecast2", func() telegraf.Input {
		return &IcecastInputPlugin{
			Url:              "http://localhost:8000/admin/",
			Username:         "admin",
			Password:         "hackme",
			ResponseTimeout:  config.Duration(time.Second * 5),
			CollectListeners: true,
			pluginVersion:    pluginVersion(),
		}
	})
}

// Init is for setup, and validating config.
func (input *IcecastInputPlugin) Init() error {

	log.Debug().Msgf("Init called...")

	input.iceastCollector = newIceastCollector()
	if err := input.iceastCollector.Init(input); err != nil {
		return err
	}

	return nil
}

func (input *IcecastInputPlugin) Stop() {
	log.Debug().Msg("Stop called...")
}

func (input *IcecastInputPlugin) SampleConfig() string {
	return icecast2.SampleConfig
}

func (input *IcecastInputPlugin) Description() string {
	return "Gather information from icecast2 server"
}

func (input *IcecastInputPlugin) Gather(acc telegraf.Accumulator) (err error) {
	log.Debug().Msg("Gathering metrics...")
	defer log.Debug().Msg("Done gathering metrics")

	err = input.iceastCollector.collect()
	if err != nil {
		return
	}

	return
}

func pluginVersion() string {

	if bi, ok := debug.ReadBuildInfo(); ok {
		for _, module := range bi.Deps {
			if module.Path == "github.com/spacepille/icecast2-telegraf-plugin" {
				log.Debug().Msgf("Found version %v", module.Version)
				return module.Version
			}
		}
	}

	return "(unknown)"
}
