# Deployment Flow

All deployment steps are driven by Mage targets defined in `magefile.go`.

## `mage all` orchestration

```mermaid
flowchart LR
    START(["mage all"]) --> K1["Kind.CreateOlly"]
    START --> K2["Kind.CreateApps"]

    subgraph OLLY["observability cluster (sequential)"]
        K1 --> IG["Prometheus.InstallGlobal"] --> DG["Prometheus.DeployGlobal"] --> LG["LGTM.Deploy"]
    end

    subgraph APPSC["apps cluster (sequential)"]
        K2 --> IW["Prometheus.InstallWriter"] --> DR["Prometheus.DeployRemote"] --> AD["Apps.Deploy"]
    end

    LG --> DONE(["stack ready"])
    AD --> DONE
```

Both clusters are created in parallel (`mg.Deps`); provisioning inside each cluster runs
sequentially (`mg.SerialDeps`).

## Manual step-by-step order

```mermaid
flowchart TD
    A["1. mage kind:createOlly<br/>create kind-observability-stack"]
      --> B["2. mage prometheus:installglobal<br/>namespace + operator bundle + RBAC"]
      --> C["3. mage prometheus:deployglobal<br/>Prometheus CR (write receiver) + services + ServiceMonitor"]
      --> D["4. mage LGTM:deploy<br/>helm install lgtm-distributed (namespace monitoring)"]
      --> E["5. mage kind:createApps<br/>create kind-demo-apps"]
      --> F["6. mage prometheus:installwriter<br/>namespace + operator bundle + RBAC"]
      --> G["7. mage prometheus:deployremote<br/>Prometheus CR (remote write) + services + ServiceMonitor"]
      --> H["8. mage apps:deploy<br/>example-app + service + ServiceMonitor"]
```

## What each install step applies

```mermaid
flowchart LR
    subgraph INSTALL["Install (per cluster)"]
        N["prometheus/namespace.yaml"] --> O["prometheus/operator/bundle.yaml"] --> W["wait for operator pod"] --> R["prometheus/rbac.yaml"]
    end

    subgraph DEPLOY["Deploy global"]
        GI["prometheus/global/instance.yaml"] --> GS["prometheus/service.yaml"] --> GE["prometheus/global/extservice.yaml"] --> GSM["prometheus/global/servicemonitor.yaml"]
    end

    subgraph DEPLOYW["Deploy writer"]
        WI["prometheus/writer/instance.yaml"] --> WS["prometheus/service.yaml"] --> WE["prometheus/writer/extservice.yaml"] --> WSM["prometheus/writer/servicemonitor.yaml"]
    end

    subgraph APPD["Deploy apps"]
        ADp["apps/sample/deployment.yaml"] --> ASv["apps/sample/service.yaml"] --> ASM["apps/sample/servicemonitor.yaml"]
    end

    INSTALL --> DEPLOY
    INSTALL --> DEPLOYW
    DEPLOYW --> APPD
```

## Teardown

```mermaid
flowchart LR
    D["mage deleteAll"] --> DA["Kind.DeleteApps<br/>(interactive confirmation)"] --> DO["Kind.DeleteOlly<br/>(interactive confirmation)"]
```
