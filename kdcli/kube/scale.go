package kube

import (
	"fmt"
	"net/http"

	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/caarlos0/kube-dash/kdcli/config"
)

func init() {
	log.SetHandler(cli.Default)
}

//Scale scales the provided deployment in the provided namespace by one
func Scale(namespace string, deployment string, action string) error {
	dashURL, err := config.Load()
	if err != nil {
		return err
	}

	log.WithField("kube-dash-url", *dashURL).
		WithField("action", action).
		WithField("namespace", namespace).
		Info("Scaling...")

	url := fmt.Sprintf("%s/api/deployments/%s/%s/%s", *dashURL, namespace, deployment, action)
	req, err := http.NewRequest("PUT", url, nil)
	if err != nil {
		return err
	}

	c := &http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		return err
	}

	resp.Body.Close()
	log.Info("Deployment successfully scaled!")
	return nil
}
