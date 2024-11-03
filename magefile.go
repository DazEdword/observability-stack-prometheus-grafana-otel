//go:build mage

package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/bitfield/script"
	"github.com/cenkalti/backoff/v4"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// Kind cluster management
type Kind mg.Namespace
type Prometheus mg.Namespace
type LGTM mg.Namespace
type Apps mg.Namespace

const PodNotFoundErrMessage = "error: no matching resources found"

func All() error {
	mg.Deps(Kind.CreateOlly, Kind.CreateApps)

	// move context to Olly cluster and install/deploy all
	if err := sh.RunV("kubectx", "kind-observability-stack"); err != nil {
		return err
	}

	mg.SerialDeps(Prometheus.Install, Prometheus.Deploy, LGTM.Deploy)

	return nil
}

func DeleteAll() error {
	mg.SerialDeps(Kind.DeleteApps, Kind.DeleteOlly)

	return nil
}

func (Kind) CreateOlly() error {
	if err := sh.RunV("kind", "create", "cluster", "--name", "observability-stack", "--config=deploy/clusterconfig.yaml"); err != nil {
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

	// pods might not be immediately available after creation, hence we attempt the wait with backoff
	err := backoff.Retry(
		func() error {
			o, err := sh.Output("kubectl", "wait", "--for=condition=Ready", "pods", "-l", "app.kubernetes.io/name=prometheus-operator", "--timeout", "120s")

			if err != nil {
				if strings.Contains(o, PodNotFoundErrMessage) {
					return backoff.Permanent(err)
				}

				return err

			} else {
				return err
			}

		}, backoff.WithMaxRetries(backoff.NewConstantBackOff(5*time.Second), 6))

	if err != nil {
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

	// pods might not be immediately available after creation, hence we attempt the wait with backoff
	err := backoff.Retry(
		func() error {
			o, err := sh.Output("kubectl", "wait", "--for=condition=Ready", "pods", "-l", "app.kubernetes.io/instance=prometheus", "--timeout", "120s")

			if err != nil {
				if strings.Contains(o, PodNotFoundErrMessage) {
					return backoff.Permanent(err)
				}

				return err

			} else {
				return err
			}

		}, backoff.WithMaxRetries(backoff.NewConstantBackOff(5*time.Second), 6))

	if err != nil {
		return err
	}

	if err := sh.RunV("kubectl", "create", "-f", "deploy/prometheus/service.yaml"); err != nil {
		return err
	}

	if err := sh.RunV("kubectl", "create", "-f", "deploy/prometheus/extservice.yaml"); err != nil {
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

func (Prometheus) Forward() error {
	// todo remove existing custom resources in each namespace
	if err := sh.RunV("kubectl", "port-forward", "svc/prometheus", "9090"); err != nil {
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

	fmt.Println("Username:")
	fmt.Printf("%s\n\n", "admin")

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
