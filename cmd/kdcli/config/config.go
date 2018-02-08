package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/user"

	"github.com/apex/log"
	input "github.com/tcnksm/go-input"
)

//Load loads the configured kube-dash url
func Load() (*string, error) {
	var dashURL *string

	usr, err := user.Current()
	if err != nil {
		return nil, err
	}

	path := fmt.Sprintf("%s/.kube_dash", usr.HomeDir)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.WithField("path", path).
			Info("Configuration not found!")
		os.MkdirAll(path, 0755)
		dashURL, err = askDashURLToUser()
		if err != nil {
			return nil, err
		}
		err = save(dashURL, path)
		if err != nil {
			return nil, err
		}
	} else {
		filePath := fmt.Sprintf("%s/config", path)
		bt, err := ioutil.ReadFile(filePath)
		if err != nil {
			return nil, err
		}
		content := string(bt)
		dashURL = &content
	}
	return dashURL, nil
}

func save(dashURL *string, configPath string) error {
	filePath := fmt.Sprintf("%s/config", configPath)

	bytes := []byte(*dashURL)
	err := ioutil.WriteFile(filePath, bytes, 0755)
	if err != nil {
		return err
	}
	return nil
}

func askDashURLToUser() (*string, error) {
	ui := &input.UI{
		Writer: os.Stdout,
		Reader: os.Stdin,
	}

	result, err := ui.Ask("Please, inform the kube-dashboard URL", &input.Options{
		Required: true,
	})
	if err != nil {
		return nil, err
	}
	return &result, err
}
