package main

import (
	"github.com/podtato-head/podtato-head-app/pkg/podtatoserver"
	urcli "github.com/urfave/cli/v3"
	"log"
	"os"
)

type config struct {
	Component string `json:"component"`
	Port      string `json:"port"`
}

var c config

func main() {
	app := &urcli.App{
		Name: "podtato-head",
		Flags: []urcli.Flag{
			&urcli.StringFlag{
				Name:        "component",
				Value:       "all",
				Usage:       "all, leftArm, rightArm, leftLeg, rightLeg",
				Destination: &c.Component,
			},
			&urcli.StringFlag{
				Name:        "port",
				Value:       "8080",
				Destination: &c.Port,
			},
		},
		Action: func(*urcli.Context) error {
			execute()
			return nil
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func execute() {
	podtatoserver.Run(c.Component, c.Port)
}
