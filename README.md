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

The stack is configured to include the `prometheus` data source and a `prometheus` example dashboard, as defined in `deploy/lgtm/values.yaml`.
After forwarding Grafana, the application will be available at `http://localhost:3000`.

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

