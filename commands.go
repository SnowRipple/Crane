package main

import (
	"github.com/SnowRipple/crane/command"
	"github.com/SnowRipple/crane/config"
	"github.com/mitchellh/cli"
	"os"
)

// Commands is the mapping of all the available Crane commands.
var Commands map[string]cli.CommandFactory

func init() {
	ui := &cli.BasicUi{Writer: os.Stdout}

	Commands = map[string]cli.CommandFactory{

		"create": func() (cli.Command, error) {
			return &command.CreateCommand{
				Ui: ui,
			}, nil
		},

		"pull": func() (cli.Command, error) {
			return &command.PullCommand{
				Ui:         ui,
				Containers: config.ReadConfig().CraneConfig.Containers,
			}, nil
		},

		"rmi": func() (cli.Command, error) {
			return &command.RemoveImageCommand{
				Ui:         ui,
				Containers: config.ReadConfig().CraneConfig.Containers,
			}, nil
		},

		"start": func() (cli.Command, error) {
			return &command.StartCommand{
				Ui:     ui,
				Config: config.ReadConfig(),
			}, nil
		},

		"destroy": func() (cli.Command, error) {
			return &command.DestroyCommand{
				Ui:     ui,
				Config: config.ReadConfig(),
			}, nil
		},

		"build": func() (cli.Command, error) {
			return &command.BuildImageCommand{
				Ui:         ui,
				Containers: config.ReadConfig().CraneConfig.Containers,
			}, nil
		},

		"run": func() (cli.Command, error) {
			return &command.RunCommand{
				Ui:     ui,
				Config: config.ReadConfig(),
			}, nil
		},

		"runall": func() (cli.Command, error) {
			return &command.RunallCommand{
				Ui:     ui,
				Config: config.ReadConfig(),
			}, nil
		},

		"enter": func() (cli.Command, error) {
			return &command.EnterCommand{
				Ui:     ui,
				Config: config.ReadConfig(),
			}, nil
		},

		"freeze": func() (cli.Command, error) {
			return &command.FreezeCommand{
				Ui:     ui,
				Config: config.ReadConfig(),
			}, nil
		},

		"version": func() (cli.Command, error) {
			return &command.VersionCommand{
				Revision: GitCommit,
				Version:  VERSION,
				Status:   STATUS,
				Ui:       ui,
			}, nil
		},
	}
}
