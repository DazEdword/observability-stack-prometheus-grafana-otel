# Observability Stack
Observability stack with Open Telemetry, Prometheus, LGTM (Grafana).

## Grafana

Helm chart available: https://github.com/grafana/helm-charts/releases/tag/lgtm-distributed-2.1.0
Default values: https://github.com/grafana/helm-charts/blob/main/charts/lgtm-distributed/values.yaml

# Prerequirements

- `go` 1.23
- `mage` 1.15.0
- `kubectl` v1.30.1
- `kind` v0.24.0
- `kustomize` v5.4.1
- `skaffold` v2.13.2
- `helm` v3.15.0
- `jq` jq-1.7.1

# Installation

## Option 1: All-in-one (Mage)

Feel free to use this installation if you start form scratch, and prefer to have the full installation performed automatically.

> **_NOTE:_**  Make sure to have all prerequirements correctly installed before attempting the automated installation.

```sh
# create kind observability cluster
mage all
```

## Option 2: Manual (Mage)

Run every step only after successful completion of the previous one, in the order specified above.


```sh
# create kind observability cluster
mage kind:createOlly
```

```sh
# install prometheus operator
mage prometheus:install
```

```sh
# deploy prometheus
mage prometheus:deploy
```

```sh
# deploy the LGTM stack
mage LGTM:deploy
```

# Usage

After forwarding Grafana, it will be available in `http://localhost:3000`.


```sh
# forward Grafana for localhost access
# access the application in localhost:3000 with user 'admin' and the password shown in stdout
mage LGTM:forward
```

Optionally, `prometheus` can also be forwarded an accessed in `http://localhost:9090`.

```sh
# optionally forward prometheus
# no credentials are necessary
mage prometheus:forward
```

# Learn more
// TODO continue instrumentation and simple http server
- [Open Telemetry](https://opentelemetry.io/docs/languages/go/getting-started/)

// TODO Document Prometheus fwd urls
// http://localhost:9090/config
// http://localhost:9090/targets

// TODO document use prometheus INTERNAL in grafana dashboard
// http://prometheus-operated.default.svc:9090
// import dashboard 3662

// TODO custom app and cluster federation
// https://kodekloud.com/community/t/multi-cluster-monitoring-using-prometheus/401110/2