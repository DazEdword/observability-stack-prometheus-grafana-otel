//go:build mage

package main

import (
	"fmt"

	"github.com/bitfield/script"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// Kind cluster management
type Kind mg.Namespace
type Prometheus mg.Namespace
type LGTM mg.Namespace
type Apps mg.Namespace

func (Kind) CreateOlly() error {
	if err := sh.RunV("kind", "create", "cluster", "--name", "observability-stack"); err != nil {
		return err
	}

	return nil
}

func (Kind) DeleteOlly() error {
	// TODO confirmation prompt
	if err := sh.Run("kind", "delete", "cluster", "--name", "observability-stack"); err != nil {
		return err
	}

	return nil
}

func (Kind) CreateApps() error {
	if err := sh.RunV("kind", "create", "cluster", "--name", "demo-apps"); err != nil {
		return err
	}

	return nil
}

func (Kind) DeleteApps() error {
	// TODO confirmation prompt
	if err := sh.RunV("kind", "delete", "cluster", "--name", "demo-apps"); err != nil {
		return err
	}

	return nil
}

func (Prometheus) Install() error {
	if err := sh.RunV("kubectl", "apply", "-f", "deploy/prometheus/namespace.yaml"); err != nil {
		return err
	}

	if err := sh.RunV("kubectl", "create", "-f", "deploy/prometheus/operator/bundle.yaml"); err != nil {
		return err
	}

	if err := sh.RunV("kubectl", "wait", "--for=condition=Ready", "pods", "-l", "app.kubernetes.io/name=prometheus-operator"); err != nil {
		return err
	}

	if err := sh.RunV("kubectl", "create", "-f", "deploy/prometheus/rbac.yaml"); err != nil {
		return err
	}

	return nil
}

func (Prometheus) Deploy() error {
	if err := sh.RunV("kubectl", "create", "-f", "deploy/prometheus/instance.yaml"); err != nil {
		return err
	}

	if err := sh.RunV("kubectl", "wait", "--for=condition=Ready", "pods", "-l", "app.kubernetes.io/instance=prometheus"); err != nil {
		return err
	}

	if err := sh.RunV("kubectl", "create", "-f", "deploy/prometheus/service.yaml"); err != nil {
		return err
	}

	if err := sh.RunV("kubectl", "create", "-f", "deploy/prometheus/servicemonitor.yaml"); err != nil {
		return err
	}

	return nil
}

func (Prometheus) Remove() error {
	// todo remove existing custom resources in each namespace
	if err := sh.RunV("kubectl", "delete", "-f", "deploy/prometheus/rbac.yaml"); err != nil {
		return err
	}

	if err := sh.RunV("kubectl", "delete", "-f", "deploy/prometheus/operator/bundle.yaml"); err != nil {
		return err
	}

	return nil
}

func (LGTM) Deploy() error {
	if err := sh.RunV("helm", "upgrade", "-f", "deploy/lgtm/values.yaml", "observability-stack", "grafana/lgtm-distributed", "--create-namespace", "--namespace", "monitoring", "--install"); err != nil {
		return err
	}

	return nil
}

func (LGTM) Forward() error {
	password, err := script.Exec(`kubectl get secret --namespace monitoring observability-stack-grafana -o jsonpath="{.data.admin-password}"`).Exec("base64 --decode").String()
	if err != nil {
		return err
	}

	fmt.Println("Admin password:")
	fmt.Printf("%s\n\n", password)

	podName, err := script.Exec(`kubectl get pods --namespace monitoring -l "app.kubernetes.io/name=grafana,app.kubernetes.io/instance=observability-stack" -o jsonpath="{.items[0].metadata.name}"`).String()
	if err != nil {
		return err
	}

	fmt.Printf("Forwarding pod %s\n\n", podName)

	if err := sh.RunV("kubectl", "--namespace", "monitoring", "port-forward", podName, "3000"); err != nil {
		return err
	}

	return nil
}

// TODO fw all
// sudo -E /home/linuxbrew/.linuxbrew/bin/kubefwd svc

// TODO continue here (app, service monitor)
// https://grafana.com/docs/grafana-cloud/monitor-infrastructure/kubernetes-monitoring/configuration/config-other-methods/prometheus/prometheus-operator/

// TODO continue instrumentation and simple http server
// https://opentelemetry.io/docs/languages/go/getting-started/

func (Apps) Deploy() error {
	if err := sh.RunV("kubectl", "apply", "-f", "deploy/apps/sample/deployment.yaml"); err != nil {
		return err
	}

	if err := sh.RunV("kubectl", "apply", "-f", "deploy/apps/sample/service.yaml"); err != nil {
		return err
	}

	if err := sh.RunV("kubectl", "apply", "-f", "deploy/apps/sample/servicemonitor.yaml"); err != nil {
		return err
	}

	return nil
}
