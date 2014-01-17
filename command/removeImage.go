package command

import (
	"flag"
	"github.com/SnowRipple/crane/command/executer"
	"github.com/SnowRipple/crane/constants"
	"github.com/SnowRipple/crane/container"
	"github.com/SnowRipple/crane/utils"
	flags "github.com/jessevdk/go-flags"
	"github.com/mitchellh/cli"
	"strings"
)

// RemoveImageCommand creates an example configuration file Cranefile.toml.
type RemoveImageCommand struct {
	Ui         cli.Ui
	Containers map[string]container.Container
}

func (c *RemoveImageCommand) Help() string {
	helpText := `
	Removes docker images specified by the user.

  Usage: crane rmi <containerName1> <containerName2>

  Removes docker images specified in the Cranefile for chosen containers

  Available options:

  -a (--all) : Remove all images defined in the Cranefile.
  `
	return strings.TrimSpace(helpText)
}

//Removes images specified by the user from the host system.
func (c *RemoveImageCommand) Run(imageNames []string) int {

	logger.Debug("Entered rmi command...")

	var options constants.CommonFlags

	cmdFlags := flag.NewFlagSet("rmi", flag.ContinueOnError)
	cmdFlags.Usage = func() { c.Ui.Output(c.Help()) }

	imageNames, err := flags.ParseArgs(&options, imageNames)

	if err != nil {
		logger.Fatalf("Failed to parse rmi flags for following CLI arguments:\n%v", imageNames)
	}

	if options.All { //Remove all images
		for _, container := range c.Containers {
			imageNames = append(imageNames, container.Image)
		}
	} else if len(imageNames) == 0 {
		logger.Fatalf("No arguments detected in the rmi command.Please correct.")
	}

	dockerCommand := []string{constants.DOCKER, constants.REMOVE_IMAGE}
	dockerCommand = append(dockerCommand, imageNames...)

	removedImagesBytes, err := executer.GetCommandOutput(dockerCommand)
	if err != nil {
		logger.Fatal("Error when trying to remove images:\n%v\nError message:\n%v", imageNames, utils.ExtractContainerMessage(removedImagesBytes, err))
	}

	utils.PrintCommandOutput(removedImagesBytes)
	return 0
}

func (c *RemoveImageCommand) Synopsis() string {
	return "Remove chosen image(s)."
}
