package icecast2

import "encoding/xml"

type IcecastXmlStats struct {
	XMLName                 xml.Name            `xml:"icestats"`
	Admin                   string              `xml:"admin"`
	ClientConnections       int                 `xml:"client_connections"`
	Clients                 int                 `xml:"clients"`
	Connections             int                 `xml:"connections"`
	FileConnections         int                 `xml:"file_connections"`
	Host                    string              `xml:"host"`
	ListenerConnections     int                 `xml:"listener_connections"`
	Listeners               int                 `xml:"listeners"`
	Location                string              `xml:"location"`
	ServerId                string              `xml:"server_id"`
	ServerStart             string              `xml:"server_start"`
	ServerStartIso8601      string              `xml:"server_start_iso8601"`
	SourceClientConnections int                 `xml:"source_client_connections"`
	SourceRelayConnections  int                 `xml:"source_relay_connections"`
	SourceTotalConnections  int                 `xml:"source_total_connections"`
	Sources                 int                 `xml:"sources"`
	Stats                   int                 `xml:"stats"`
	StatsConnections        int                 `xml:"stats_connections"`
	Source                  []*IcecastXmlSource `xml:"source"`
}

type IcecastXmlSource struct {
	Mount              string `xml:"mount,attr"`
	Bitrate            int    `xml:"bitrate"`
	Genre              string `xml:"genre"`
	ListenerPeak       int    `xml:"listener_peak"`
	Listeners          int    `xml:"listeners"`
	ListenUrl          string `xml:"listenurl"`
	MaxListeners       string `xml:"max_listeners"`
	Public             int    `xml:"public"`
	ServerDescription  string `xml:"server_description"`
	ServerName         string `xml:"server_name"`
	ServerType         string `xml:"server_type"`
	ServerUrl          string `xml:"server_url"`
	SlowListeners      int    `xml:"slow_listeners"`
	SourceIp           string `xml:"source_ip"`
	StreamStart        string `xml:"stream_start"`
	StreamStartIso8601 string `xml:"stream_start_iso8601"`
	TotalBytesRead     int    `xml:"total_bytes_read"`
	TotalBytesSent     int    `xml:"total_bytes_sent"`
	AudioBitrate       int    `xml:"audio_bitrate,omitempty"`
	AudioChannels      int    `xml:"audio_channels,omitempty"`
	AudioSamplerate    int    `xml:"audio_samplerate,omitempty"`
	IceBitrate         int    `xml:"ice-bitrate,omitempty"`
	Subtype            string `xml:"subtype,omitempty"`
}

/*
type IcecastXmlClientList_2_4 struct {
	XMLName xml.Name `xml:"icestats"`
	Source  struct {
		Mount     string                    `xml:"mount,attr"`
		Listeners int                       `xml:"Listeners"`
		Listener  []*IcecastXmlListener_2_4 `xml:"listener"`
	} `xml:"source"`
}

type IcecastXmlListener_2_4 struct {
	IP        string `xml:"IP"`
	UserAgent string `xml:"UserAgent"`
	Connected int    `xml:"Connected"`
	ID        int    `xml:"ID"`
}
*/

type IcecastXmlClientList struct {
	XMLName xml.Name `xml:"icestats"`
	Source  struct {
		Mount     string                `xml:"mount,attr"`
		Listeners int                   `xml:"listeners"`
		Listener  []*IcecastXmlListener `xml:"listener"`
	} `xml:"source"`
}

type IcecastXmlListener struct {
	ID        string `xml:"id"`
	IP        string `xml:"ip"`
	UserAgent string `xml:"useragent"`
	Connected string `xml:"connected"`
	// 	skip 2.5 fields
	/*
		Host      string `xml:"host,omitempty"`
		Role      string `xml:"role,omitempty"`
		Acl       string `xml:"acl,omitempty"`
		Tls       string `xml:"tls,omitempty"`
		Protocol  string `xml:"protocol,omitempty"`
		Referer   string `xml:"referer,omitempty"`
	*/
}
