package icecast2

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/gookit/goutil/dump"
	"github.com/oschwald/geoip2-golang"
	"github.com/rs/zerolog/log"
)

type IceastCollector struct {
	plugin *IcecastInputPlugin
	geoip  *geoip2.Reader
	url    *url.URL
}

func newIceastCollector() *IceastCollector {
	return &IceastCollector{}
}

func (col *IceastCollector) Init(plugin *IcecastInputPlugin) error {
	col.plugin = plugin

	if plugin.Geoip2Path != "" {
		db, err := geoip2.Open(plugin.Geoip2Path)
		if err != nil {
			//log.Fatal().Err(err).Msg("Unable to open GeopIP2 mmdb file")
			return err
		}
		col.geoip = db
	}

	url, err := url.ParseRequestURI(plugin.Url)
	if err != nil {
		return err
	}

	col.url = url

	if _, err := col.fetchStats(); err != nil {
		return err
	}

	return nil
}

func (col *IceastCollector) collect() error {
	return nil
}

func (col *IceastCollector) fetchStats() (*IcecastXmlStats, error) {

	var icecastStats IcecastXmlStats
	if data, err := col.fetchRaw("stats"); err != nil {
		return nil, err
	} else {
		if err := xml.Unmarshal(data, &icecastStats); err != nil {
			return nil, err
		}
	}

	//log.Debug().Msgf("Send request to %s", urlString)
	dump.P(icecastStats)

	return &icecastStats, nil
}

func (col *IceastCollector) fetchRaw(path string) ([]byte, error) {

	timeout := time.Duration(col.plugin.ResponseTimeout)
	ctx, cncl := context.WithTimeout(context.Background(), timeout)
	defer cncl()

	url, err := col.url.Parse(path)
	if err != nil {
		return nil, err
	}

	urlString := url.String()

	log.Debug().Msgf("Send request to %s", urlString)

	req, err := http.NewRequestWithContext(ctx, "GET", urlString, nil)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(col.plugin.Username, col.plugin.Password)

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http request failed: error %s", resp.Status)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	//bodyString := string(bodyBytes)
	//log.Info().Msg(bodyString)
	return bodyBytes, nil
}
