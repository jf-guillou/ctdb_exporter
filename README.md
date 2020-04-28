# ctdb_exporter
Prometheus exporter for CTDB

## Usage
```
Usage of ./ctdb_exporter:
-web.listen-address string
The address to listen on for HTTP requests. (default ":9725")
-web.endpoint
The endpoint exposing metrics. (default "/metrics")
-ctdb.bin-path
The complete path to ctdb binary. (default "/usr/bin/ctdb")
-ctdb.sudo
Prefix ctdb commands with sudo. (default true)
```

## Requirements

- CTDB
- [sudo](#sudo) (optional)

## Sudo

This tool is most likely going to need sudo as default CTDB configuration locks ctdb.socket access to root only.
Assuming you are running the exporter with the user `prometheus`,
the easiest way to handle this would be creating a `/etc/sudoers.d/ctdb_exporter` file containing :
```
prometheus ALL=(ALL) NOPASSWD: /usr/bin/ctdb pnn,/usr/bin/ctdb recmaster,/usr/bin/ctdb status -Y,/usr/bin/ctdb statistics -Y
```

## Prometheus configuration

Minimal Prometheus scrape configuration : 

```
scrape_configs:
  - job_name: "ctdb"
    static_configs:
      - targets:
        - samba-01.example.com:9725
        - samba-02.example.com:9725
```

## Exposed metrics

`ctdb_up` will return 0 on scrape errors.

The results of `ctdb status -Y` on master node and `ctdb statistics -Y` on all nodes will be returned as gauges.

Example metrics :
```
ctdb_up 1
...
ctdb_banned{id="1",ip="0.0.0.1"} 0
ctdb_banned{id="2",ip="0.0.0.2"} 0
ctdb_disconnected{id="1",ip="0.0.0.1"} 0
ctdb_disconnected{id="2",ip="0.0.0.2"} 0
...
ctdb_num_clients{id="1"} 12
ctdb_num_clients{id="2"} 21
```
