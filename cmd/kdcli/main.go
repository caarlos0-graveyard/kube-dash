package main

import (
	"errors"
	"os"

	"github.com/apex/log"
	lcli "github.com/apex/log/handlers/cli"
	"github.com/caarlos0/kube-dash/kdcli/kube"
	"github.com/urfave/cli"
)

func init() {
	log.SetHandler(lcli.Default)
}

func main() {
	app := cli.NewApp()
	app.Name = "kdcli"
	app.Usage = "Implement's kube-dash API to manage resources from command line"
	app.Version = "1.0.0"
	app.Commands = []cli.Command{
		{
			Name:        "scale-up",
			Aliases:     []string{"s"},
			Usage:       "kdcli scale-up {deployment-name}",
			Description: "Scales up the deployment by 1",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "namespace",
					Value: "default",
					Usage: "Inform a specific namespace for the scale-up",
				},
			},
			Action: func(c *cli.Context) error {
				return scale(c)
			},
		},
		{
			Name:        "scale-down",
			Aliases:     []string{"s"},
			Usage:       "kdcli scale-down {deployment-name}",
			Description: "Scales down the deployment by 1",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "namespace",
					Value: "default",
					Usage: "Inform a specific namespace for the scale-down",
				},
			},
			Action: func(c *cli.Context) error {
				return scale(c)
			},
		},
	}
	app.Run(os.Args)
}

func scale(c *cli.Context) error {
	if len(c.Args()) > 0 {
		ns := c.String("namespace")
		deployment := c.Args().Get(0)
		var action string
		if c.Command.Name == "scale-up" {
			action = "up"
		} else {
			action = "down"
		}
		err := kube.Scale(ns, deployment, action)
		if err != nil {
			log.WithError(err).Error("Failed to scale deployment")
			return cli.NewExitError("\n", 1)
		}
	} else {
		log.WithError(errors.New("Missing parameter for scale command")).
			Error("Failed to scale deployment")
		return cli.NewExitError("\n", 1)
	}
	return nil
}
