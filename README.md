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

For convenience, `prometheus` service for the global instance
is exposed at `http://localhost:30900`.

Optionally, any of the `prometheus` instances can also be forwarded an accessed in `http://localhost:9090`.

```sh
# optionally forward prometheus
# no credentials are necessary
mage prometheus:forward
```

When configuring `prometheus` as a data source in `grafana`, the internal node ip needs to be used instead, followed by the exposed service.
For instance: `http://172.19.0.2:30900`



# Learn more
// TODO continue instrumentation and simple http server
- [Open Telemetry](https://opentelemetry.io/docs/languages/go/getting-started/)

// TODO Document Prometheus fwd urls
// http://localhost:30900/config
// http://localhost:30900/targets

// TODO document use prometheus INTERNAL in grafana dashboard
// http://prometheus-operated.default.svc:9090

// In kind, 172.19.0.1 always refers to the host machine
// http://172.19.0.1:30900

// nodeport URI also works
// http://172.19.0.2:30900
// get internal IP
// kubectl get nodes -o wide


// import dashboard 3662

// TODO prometheus ingress in global
// TODO custom app and cluster remote_write from local to global
