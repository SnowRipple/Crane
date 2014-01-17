package command

import (
	"flag"
	"github.com/SnowRipple/crane/command/executer"
	"github.com/SnowRipple/crane/constants"
	"github.com/SnowRipple/crane/container"
	ownLog "github.com/SnowRipple/crane/logger"
	"github.com/SnowRipple/crane/utils"
	flags "github.com/jessevdk/go-flags"
	"github.com/mitchellh/cli"
	"strings"
)

var logger = ownLog.GetLogger()

// PullCommand creates an example configuration file Cranefile.toml.
type PullCommand struct {
	Ui         cli.Ui
	Containers map[string]container.Container
}

func (c *PullCommand) Help() string {
	helpText := `
Usage: crane pull [options] <Image1> <Image2>

Pulls chosen images from the docker public repository.
Options:

-a (--all)  Pulls all images defined in the Cranefile from the docker public repository.`
	return strings.TrimSpace(helpText)
}

//Pull chosen containers from the docker public repository
func (c *PullCommand) Run(commandArguments []string) int {
	var (
		images  []string
		options constants.CommonFlags
	)

	logger.Debug("Entered pull command...")

	cmdFlags := flag.NewFlagSet("pull", flag.ContinueOnError)
	cmdFlags.Usage = func() { c.Ui.Output(c.Help()) }

	if len(commandArguments) == 0 {
		logger.Fatalf("Not enough arguments provided for the pull command.Please correct.")
	}

	images, err := flags.ParseArgs(&options, commandArguments)

	logger.Debug("Provided images are:%v", images)

	if err != nil {
		logger.Fatalf("Failed to parse pull flags for following CLI arguments:\n%v", commandArguments)
	}

	if options.All {
		for _, container := range c.Containers {
			logger.Debug("Container image:" + container.Image)
			images = append(images, container.Image)
		}
		logger.Debug("Images defined in Cranefile that will be downloaded:%v", images)
	}

	for _, imageName := range images {
		logger.Debug("Pulling image %q", imageName)

		dockerCommand := []string{constants.DOCKER, constants.PULL, imageName}

		outputBytes, err := executer.GetCommandOutput(dockerCommand)
		if err != nil {
			logger.Fatal("Error during \"pull\" command:", utils.ExtractContainerMessage(outputBytes, err))
		}
		utils.PrintCommandOutput(outputBytes)
	}

	return 0
}

func (c *PullCommand) Synopsis() string {
	return "Pull chosen image(s)."
}
