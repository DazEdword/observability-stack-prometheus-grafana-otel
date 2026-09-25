# Architecture

The stack runs two local Kind clusters on a single UNIX host. The observability cluster
(`kind-observability-stack`) hosts the Grafana LGTM stack and the global Prometheus instance;
the apps cluster (`kind-demo-apps`) hosts a Prometheus instance in remote-write mode plus an
instrumented example app. The Kind host is reachable from every cluster at `172.19.0.1`.

## Component diagram

Operators, RBAC, and ServiceMonitor wiring are omitted here for clarity — see the resource
diagram below.

```mermaid
flowchart LR
    USER(("User / browser"))

    subgraph APPS["Cluster 2 · kind-demo-apps"]
        APP["example-app<br/>3 replicas · :8080/metrics"]
        PW["Prometheus writer<br/>remote-write mode"]
    end

    subgraph OLLY["Cluster 1 · kind-observability-stack"]
        GRAFANA["Grafana<br/>UI via port-forward :3000"]
        LOKI["Loki · logs"]
        MIMIR["Mimir · metrics<br/>default data source"]
        TEMPO["Tempo · traces"]
        NP1["NodePort Service<br/>:30900 → :9090"]
        PG["Prometheus global<br/>remote-write receiver"]
    end

    USER -->|"http://localhost:3000"| GRAFANA
    USER -->|"http://localhost:30900"| NP1
    APP -->|"scrape /metrics"| PW
    PW -->|"remote write<br/>172.19.0.1:30900/api/v1/write"| NP1
    NP1 --> PG
    GRAFANA -->|"PromQL data source<br/>172.19.0.1:30900"| NP1
    GRAFANA --- LOKI
    GRAFANA --- MIMIR
    GRAFANA --- TEMPO
```

## Kubernetes resource relationships (per cluster)

Both clusters apply the same set of manifests from `deploy/prometheus/`; only the `Prometheus`
CR and the NodePort differ (`global/` vs `writer/`).

```mermaid
flowchart LR
    subgraph INPUTS["Applied manifests"]
        B["operator/bundle.yaml<br/>CRDs + prometheus-operator"]
        RB["rbac.yaml<br/>SA · ClusterRole · Binding"]
        P["Prometheus CR<br/>global or writer instance.yaml"]
        SM["ServiceMonitors<br/>self-scrape · example-app"]
        SVC["service.yaml<br/>ClusterIP :9090"]
        EXT["extservice.yaml<br/>NodePort :30900 / :30901"]
    end

    B --> OP["prometheus-operator"]
    RB --> OP
    P -->|reconciled by| OP
    OP ==>|creates & manages| PODS["Prometheus pods<br/>+ prometheus-operated headless service"]
    SM -.->|selected via serviceMonitorSelector| PODS
    PODS -->|"pods labeled app.kubernetes.io/name=prometheus"| SVC
    PODS -->|"pods labeled app.kubernetes.io/name=prometheus"| EXT
    PODS -->|"operated-prometheus=true"| HEAD["self-scrape target<br/>:9090/metrics"]
    HEAD -.-> SM
```
