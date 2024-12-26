# icecast2-telegraf-plugin

A telegraph plugin to gather information from an Icecast 2.4 or 2.5 server.

## Download the binary

Simply download the binary from the [latest release](https://github.com/spacepille/icecast2-telegraf-plugin/releases).


## Build From Source

*These instructions assume you have Go installed and configured on your machine*

Clone the repository
```sh
git clone https://github.com/spacepille/icecast2-telegraf-plugin
cd icecast2-telegraf-plugin
```

Build the module into an executable

On Unix Linux/Mac:

```sh
go build -o icecast2-telegraf-plugin cmd/main.go
```

On Windows:

```sh
go build -o icecast2-telegraf-plugin.exe cmd/main.go
```

## Usage
Reference the executable and config in your `telegraf.conf` using the `execd` input

```toml
[[inputs.execd]]
  command = ["/path/to/icecast2-telegraf-plugin(.exe)", "-config", "/path/to/icecast2-telegraf-config.toml"]
```

More documentation on using Telegraf external plugins can be found [here](https://github.com/influxdata/telegraf/blob/master/docs/EXTERNAL_PLUGINS.md).

## Configuration

External plugins have their own configuration files. The path to the file is defined in the `inputs.execd` directive in `telegraf.conf`, as shown above.

```toml @sample.conf
[[inputs.icecast2]]

## Icecast2 Admin url
# url = "http://localhost:8000/admin/"

## Icecast2 Admin username
# username = "admin"

## Icecast2 Admin password
# password = "hackme"

## Maximum time to receive response
# response_timeout = "5s"

## Gather Listeners
# gather_listeners = true

## Path to MaxMind GeoLite2 or GeoIP2 city database
# geoip2_path = ""

## Maxmind include the country and city names in English, Simplified Chinese,
## Spanish, Brazilian Portuguese, Russian, Japanese, French, and German
# geoip2_language = "en"
```

## Metrics

- `icecast_server`
  - tags:
    - `admin`
    - `host`
    - `location`
    - `server_id`
  - fields:
    - `clients`
    - `listeners`
    - `sources`
    - `stats`    
    - `server_start`    
    - `client_connections` &nbsp;-&nbsp; *accumulating*
    - `file_connections` &nbsp;-&nbsp; *accumulating*
    - `listener_connections` &nbsp;-&nbsp; *accumulating*
    - `source_client_connections` &nbsp;-&nbsp; *accumulating*
    - `source_relay_connections` &nbsp;-&nbsp; *accumulating*
    - `source_total_connections` &nbsp;-&nbsp; *accumulating*
    - `stats_connections` &nbsp;-&nbsp; *accumulating*

- `icecast_source`
  - tags:
    - `host`
    - `mount`
    - `genre`
    - `listen_url`
    - `server_name`
    - `server_description`
    - `server_type`
    - `server_url`
    - `source_ip`
  - fields:
    - `listeners`
    - `listener_peak`
    - `slow_listeners`
    - `stream_start`
    - `total_bytes_read` &nbsp;-&nbsp; *accumulating*
    - `total_bytes_sent` &nbsp;-&nbsp; *accumulating*

- `icecast_listener`
  - tags:
    - `host`
    - `mount`
    - `ip`
    - `user_agent`
    - `continent_code` &nbsp;-&nbsp; *requires geoip2 db*
    - `country_code` &nbsp;-&nbsp; *requires geoip2 db*
    - `postal_code` &nbsp;-&nbsp; *requires geoip2 db*
    - `continent_name` &nbsp;-&nbsp; *requires geoip2 db and geoip2 language*
    - `country_name` &nbsp;-&nbsp; *requires geoip2 db and geoip2 language
    - `city_name` &nbsp;-&nbsp; *requires geoip2 db and geoip2 language
  - fields:
    - `connected`
    - `latitude` &nbsp;-&nbsp; *requires geoip2 db*
    - `longitude` &nbsp;-&nbsp; *requires geoip2 db*

## Example Output

```text
icecast_server,admin=icemaster@localhost,host=localhost,location=Earth,server_id=Icecast\ 2.5-beta.3\ xm-1.0 clients=92i,listeners=76i,stats=0i,client_connections=90593i,listener_connections=46646i,source_relay_connections=0i,sources=12i,file_connections=362i,source_client_connections=12i,source_total_connections=12i,stats_connections=0i 1735227248468115646
icecast_source,genre=Ambient\,\ Downbeat,host=localhost,listen_url=http://127.0.0.1:8000/ambient.aac,mount=ambient.aac,server_description=Music\ Radio,server_name=Ambient\ Station,server_type=audio/aacp,server_url=http://localhost,source_ip=127.0.0.1 total_bytes_read=530585047i,total_bytes_sent=4501043i,listeners=0i,listener_peak=1i,slow_listeners=0i 1735227248468169340
icecast_source,genre=Ambient\,\ Downbeat,host=localhost,listen_url=http://127.0.0.1:8000/ambient.ogg,mount=ambient.ogg,server_description=Music\ Radio,server_name=Ambient\ Station,server_type=audio/ogg,server_url=http://localhost,source_ip=127.0.0.1 total_bytes_read=1133264896i,total_bytes_sent=20652068461i,listeners=27i,listener_peak=30i,slow_listeners=40i 1735227248468391286
icecast_listener,city_name=Belgrade,continent_code=EU,continent_name=Europe,country_code=RS,country_name=Serbia,host=localhost,ip=188.120.1.1,mount=ambient.ogg,user_agent=GStreamer\ souphttpsrc\ 1.24.2\ libsoup/3.4.4 connected="7591",latitude=44.8046,longitude=20.4637 1735227248965060403
icecast_listener,city_name=San\ Luis\ Potosí\ City,continent_code=NA,continent_name=North\ America,country_code=MX,country_name=Mexico,host=localhost,ip=201.152.1.1,mount=ambient.ogg,postal_code=78250,user_agent=GStreamer\ souphttpsrc\ 1.16.2\ libsoup/2.70.0 longitude=-100.997,connected="7186",latitude=22.1615 1735227249000092134
```
