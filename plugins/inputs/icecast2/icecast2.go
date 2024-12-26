package icecast2

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/gookit/goutil/dump"
	"github.com/influxdata/telegraf"
	"github.com/oschwald/geoip2-golang"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type IceastCollector struct {
	plugin      *IcecastInputPlugin
	geoIpReader *geoip2.Reader
	geoIpInfo   os.FileInfo
	url         *url.URL
}

func newIceastCollector() *IceastCollector {
	return &IceastCollector{}
}

func (col *IceastCollector) Init(plugin *IcecastInputPlugin) error {

	col.plugin = plugin

	url, err := url.ParseRequestURI(plugin.Url)
	if err != nil {
		return err
	}

	col.url = url

	if col.plugin.GatherListeners {
		err = col.openGeoIpDb()
		if err != nil {
			return err
		}
	}

	_, err = col.fetchStats()
	if err != nil {
		return err
	}

	return nil
}

func (col *IceastCollector) openGeoIpDb() error {

	if col.plugin.Geoip2Path == "" {
		return nil
	}

	fileInfo, err := os.Stat(col.plugin.Geoip2Path)
	if err != nil {
		return err
	}

	if col.geoIpReader != nil {
		// check if geopip file changed
		if !fileInfo.ModTime().Equal(col.geoIpInfo.ModTime()) {
			return nil
		}
		// check if change is older than 5 minutes
		if !fileInfo.ModTime().Add(time.Minute * 5).Before(time.Now()) {
			return nil
		}
		// close geoip file
		_ = col.geoIpReader.Close()
		col.geoIpReader = nil
	}

	// open geopip file
	db, err := geoip2.Open(col.plugin.Geoip2Path)
	if err != nil {
		return err
	}

	col.geoIpInfo = fileInfo
	col.geoIpReader = db

	return nil
}

func (col *IceastCollector) gather(acc telegraf.Accumulator) error {

	stats, err := col.fetchStats()
	if err != nil {
		return err
	}

	col.gatherServerMetrics(acc, stats)
	col.gatherSourceMetrics(acc, stats)

	if col.plugin.GatherListeners {

		err = col.openGeoIpDb()
		if err != nil {
			return err
		}

		for _, source := range stats.Source {

			clientList, err := col.fetchClients(source.Mount)
			if err != nil {
				return err
			}

			err = col.gatherListenerMetrics(acc, stats, source, clientList)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (col *IceastCollector) gatherListenerMetrics(
	acc telegraf.Accumulator, stats *IcecastXmlStats,
	source *IcecastXmlSource, clientList *IcecastXmlClientList,
) error {

	for _, listener := range clientList.Source.Listener {
		records := make(map[string]interface{})
		tags := make(map[string]string)

		tags["host"] = stats.Host
		tags["mount"] = source.Mount[1:] // without leading slash
		tags["ip"] = listener.IP
		tags["user_agent"] = listener.UserAgent

		// counters
		records["connected"] = listener.Connected

		if col.geoIpReader != nil {

			ip := net.ParseIP(listener.IP)

			city, err := col.geoIpReader.City(ip)
			if err != nil {
				return err
			}

			tags["continent_code"] = city.Continent.Code
			tags["country_code"] = city.Country.IsoCode
			tags["postal_code"] = city.Postal.Code

			if col.plugin.Geoip2Language != "" {
				tags["continent_name"] = city.Continent.Names[col.plugin.Geoip2Language]
				tags["country_name"] = city.Country.Names[col.plugin.Geoip2Language]
				tags["city_name"] = city.City.Names[col.plugin.Geoip2Language]
			}

			records["latitude"] = city.Location.Latitude
			records["longitude"] = city.Location.Longitude
		}

		acc.AddFields("icecast_listener", records, tags)
	}

	return nil
}

func (col *IceastCollector) gatherSourceMetrics(
	acc telegraf.Accumulator, stats *IcecastXmlStats,
) {
	for _, source := range stats.Source {
		records := make(map[string]interface{})
		tags := make(map[string]string)

		tags["host"] = stats.Host
		tags["mount"] = source.Mount[1:] // without leading slash
		tags["genre"] = source.Genre
		tags["listen_url"] = source.ListenUrl
		tags["server_name"] = source.ServerName
		tags["server_description"] = source.ServerDescription
		tags["server_type"] = source.ServerType
		tags["server_url"] = source.ServerUrl
		tags["source_ip"] = source.SourceIp

		// counters
		records["listeners"] = source.Listeners
		records["listener_peak"] = source.ListenerPeak
		records["slow_listeners"] = source.SlowListeners

		// Todo: check if it works
		startTime, _ := time.Parse(time.RFC3339, source.StreamStartIso8601)
		records["stream_start"] = startTime.UnixNano()

		// accumulating counters
		records["total_bytes_read"] = source.TotalBytesRead
		records["total_bytes_sent"] = source.TotalBytesSent

		acc.AddFields("icecast_source", records, tags)
	}
}

func (col *IceastCollector) gatherServerMetrics(
	acc telegraf.Accumulator, stats *IcecastXmlStats,
) {
	records := make(map[string]interface{})
	tags := make(map[string]string)

	tags["admin"] = stats.Admin
	tags["host"] = stats.Host
	tags["location"] = stats.Location
	tags["server_id"] = stats.ServerId

	// counters
	records["clients"] = stats.Clients
	records["listeners"] = stats.Listeners
	records["sources"] = stats.Sources
	records["stats"] = stats.Stats

	// Todo: check if it works
	startTime, _ := time.Parse(time.RFC3339, stats.ServerStartIso8601)
	records["server_start"] = startTime.UnixNano()

	// accumulating counters
	records["client_connections"] = stats.ClientConnections
	records["file_connections"] = stats.FileConnections
	records["listener_connections"] = stats.ListenerConnections
	records["source_client_connections"] = stats.SourceClientConnections
	records["source_relay_connections"] = stats.SourceRelayConnections
	records["source_total_connections"] = stats.SourceTotalConnections
	records["stats_connections"] = stats.StatsConnections

	acc.AddFields("icecast_server", records, tags)
}

func (col *IceastCollector) fetchClients(mount string) (*IcecastXmlClientList, error) {
	var icecastClientList IcecastXmlClientList

	data, err := col.fetchRaw("listclients?mount=" + mount)
	if err != nil {
		return nil, err
	}

	// convert from 2.4 format to 2.5
	data = bytes.ReplaceAll(data, []byte("ID>"), []byte("id>"))
	data = bytes.ReplaceAll(data, []byte("IP>"), []byte("ip>"))
	data = bytes.ReplaceAll(data, []byte("UserAgent>"), []byte("useragent>"))
	data = bytes.ReplaceAll(data, []byte("Connected>"), []byte("connected>"))
	data = bytes.ReplaceAll(data, []byte("Listeners>"), []byte("listeners>"))

	//log.Debug().Msg(string(data))

	err = xml.Unmarshal(data, &icecastClientList)
	if err != nil {
		return nil, err
	}

	if zerolog.GlobalLevel() == zerolog.DebugLevel {
		dump.P(icecastClientList)
	}

	return &icecastClientList, nil
}

func (col *IceastCollector) fetchStats() (*IcecastXmlStats, error) {

	var icecastStats IcecastXmlStats

	data, err := col.fetchRaw("stats")
	if err != nil {
		return nil, err
	}

	err = xml.Unmarshal(data, &icecastStats)
	if err != nil {
		return nil, err
	}

	if zerolog.GlobalLevel() == zerolog.DebugLevel {
		dump.P(icecastStats)
	}

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
