# Observability Stack
Observability stack with Prometheus, LGTM (Grafana), Open Telemetry.

This repository includes a multi-cluster setup serving as a practical Grafana example, with Prometheus as a data source.
The setup is based on Kind local cluster, and it's completely self contained, designed to run locally in UNIX systems.
Sucessfully tested in Linux Mint 22 (Wilma) and macOS Sequoia 15.1.

Cluster 1 (`kind-observability-stack`) has the full Grafana LGTM stack and a Prometheus instance, configured  write receiver.
It also includes a service monitoring setup and a dashboard to monitor the Prometheus installation itself.

Cluster 2 (`kind-demo-apps`) has another Prometheus instance, this time configured as a remote writer and pointing to the global instance, and
an example app generating metrics that are sent remotely.

## Grafana

Helm chart available: https://github.com/grafana/helm-charts/releases/tag/lgtm-distributed-2.1.0
Default values: https://github.com/grafana/helm-charts/blob/main/charts/lgtm-distributed/values.yaml

# Prerequirements

Pre-setup:
- `go` 1.23
- `brew` 4.4.3

Dev dependencies:
- `mage` 1.15.0
- `kubectl` v1.31.2
- `kubectx` 0.9.5
- `kind` v0.24.0
- `kustomize` v5.4.1
- `skaffold` v2.13.2
- `kubefwd` 1.22.5
- `helm` v3.15.0
- `jq` jq-1.7.1
- `golangci-lint` 1.63.4

# Installation

Dev dependencies can be installed manually, or automatically via `brew`:

```sh
# install brew dev dependencies
mage setup
```

## Option 1: All-in-one (Mage)

Feel free to use this installation if you start form scratch, and prefer to have the full installation performed automatically.

> **_NOTE:_**  Make sure to have all prerequirements correctly installed before attempting the automated installation.

```sh
# create kind observability cluster
mage all
```

## Option 2: Manual (Mage)

Run every step only after successful completion of the previous one, in the order specified below.


```sh
# create kind observability cluster
mage kind:createOlly
```

```sh
# install prometheus operator
mage prometheus:installglobal
```

```sh
# deploy prometheus
mage prometheus:deployglobal
```

```sh
# deploy the LGTM stack
mage LGTM:deploy
```

```sh
# create kind apps cluster
mage kind:createApps
```

```sh
# install prometheus operator
mage prometheus:installwriter
```

```sh
# deploy prometheus (remote writer mode)
mage prometheus:deployremote
```

```sh
# deploy instrumented apps (example app)
mage apps:deploy
```


# Usage

The stack is configured to include the `prometheus` data source and a `prometheus` example dashboard, as defined in `deploy/lgtm/values.yaml`.
After forwarding Grafana, the application will be available at `http://localhost:3000`.

> **_NOTE:_**  The Grafana user and password will be visible in the forward command's output

```sh
# forward Grafana for localhost access
# access the application in localhost:3000 with user 'admin' and the password shown in stdout
mage LGTM:forward
```

For convenience, the `prometheus` service for the global instance is exposed at `http://localhost:30900`.
Optionally, any of the `prometheus` instances can also be forwarded an accessed in `http://localhost:9090`.
These can be useful for user access via browser UI.

```sh
# optionally forward prometheus
# no credentials are necessary
mage prometheus:forward
```

Some useful URLs:
- http://localhost:30900/config
- http://localhost:30900/targets: shows global cluster targets
- http://localhost:30900/graph: metrics explorer will show metrics from both the global and writer clusters.


When configuring `prometheus` as a data source in `grafana`, the host machine IP needs to be used instead, followed by the exposed service.
Kind's host machine IP is always `172.19.0.1`, and as such it can be used reliably to point to the exposed `prometheus` instance.
`http://172.19.0.1:30900`

> **_NOTE:_**  The node's Internal IP could be use instead and it would similarly work. Find the node's internal IP with the command `kubectl get nodes -o wide`.
> **_NOTE:_**  If the Grafana and Prometheus setups live within the same cluster, the internal service can also be used: `http://prometheus-operated.default.svc:9090`.

# Learn more
- [Open Telemetry](https://opentelemetry.io/docs/languages/go/getting-started/)
- [Prometheus CRDs](https://doc.crds.dev/github.com/prometheus-operator/prometheus-operator/monitoring.coreos.com/Prometheus/v1@v0.77.1)
- [Prometheus Remote Write](https://last9.io/blog/what-is-prometheus-remote-write/)
- [Prometheus Multi Cluster](https://sysrant.com/posts/prometheus-multi-cluster/)
- [Grafana Provisioning](https://grafana.com/tutorials/provision-dashboards-and-data-sources/)
- [Grafana Data Source - Prometheus](https://grafana.com/docs/grafana/latest/datasources/prometheus/)



