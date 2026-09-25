# Metrics Data Flow

How a sample travels from the instrumented app in the apps cluster to a Grafana dashboard
served from the observability cluster.

## Pipeline sequence

```mermaid
sequenceDiagram
    autonumber
    participant App as example-app (kind-demo-apps)
    participant Writer as prometheus-writer (kind-demo-apps)
    participant Node as Kind host 172.19.0.1:30900
    participant Global as prometheus-global (kind-observability-stack)
    participant Grafana as Grafana (kind-observability-stack)
    participant User as User browser

    Note over Writer: operator discovers ServiceMonitor example-app
    Writer->>App: GET :8080/metrics (periodic scrape)
    App-->>Writer: metric samples
    Writer->>Writer: append to local TSDB
    Writer->>Node: remote write POST /api/v1/write
    Node->>Global: NodePort forwards to :9090
    Global->>Global: append to global TSDB

    User->>Grafana: open dashboard (localhost:3000)
    Grafana->>Node: PromQL query (data source http://172.19.0.1:30900)
    Node-->>Grafana: results (writer + global series)
    Grafana-->>User: rendered panels
```

## Write path at a glance

```mermaid
flowchart LR
    APP["example-app<br/>:8080/metrics"] -->|scrape| PW["prometheus-writer"]
    PW -->|"remote write<br/>host 172.19.0.1:30900"| PG["prometheus-global<br/>write receiver"]
    PG -->|query| GF["Grafana"]
    PW -.->|query<br/>NodePort :30901| GF
```

Key configuration:

- `deploy/prometheus/writer/instance.yaml` sets `remoteWrite.url` to the global instance's
  NodePort endpoint on the Kind host (`http://172.19.0.1:30900/api/v1/write`).
- `deploy/prometheus/global/instance.yaml` sets `enableRemoteWriteReceiver: true` so the
  global instance accepts incoming remote-write traffic.
- `deploy/lgtm/values.yaml` provisions the Prometheus data source at
  `http://172.19.0.1:30900`, so Grafana queries the global instance, which holds both
  local and remotely written series.
