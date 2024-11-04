package wait

import (
	"strings"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/magefile/mage/sh"
)

const PodNotFoundErrMessage = "error: no matching resources found"

func ForPodWithConstantBackoff(selector string, timeout string) error {
	err := backoff.Retry(
		func() error {
			o, err := sh.Output("kubectl", "wait", "--for=condition=Ready", "pods", "-l", "app.kubernetes.io/name=prometheus-operator", "--timeout", "120s")

			if err != nil {
				if strings.Contains(o, PodNotFoundErrMessage) {
					return backoff.Permanent(err)
				}

				return err
			} else {
				return nil
			}

		}, backoff.WithMaxRetries(backoff.NewConstantBackOff(5*time.Second), 6))

	return err
}
