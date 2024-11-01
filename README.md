# Observability Stack
Observability stack with Open Telemetry, Prometheus, LGTM (Grafana)

## Grafana

Helm chart available: https://github.com/grafana/helm-charts/releases/tag/lgtm-distributed-2.1.0
Default values: https://github.com/grafana/helm-charts/blob/main/charts/lgtm-distributed/values.yaml

# Prerequirements

## Installation

- go 1.23
- mage

## Local cluster (for self-hosted stack)

- kubectl
- kustomize
- helm
- jq

# Usage

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

```sh
# forward the application for localhost access
# access the application in localhost:3000 with user 'admin' and the password shown in stdout
mage LGTM:forward
```