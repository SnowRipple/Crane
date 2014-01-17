package command

import (
	"flag"
	"github.com/SnowRipple/crane/command/executer"
	"github.com/SnowRipple/crane/config"
	"github.com/SnowRipple/crane/constants"
	"github.com/SnowRipple/crane/container"
	"github.com/SnowRipple/crane/io"
	"github.com/SnowRipple/crane/ssh"
	"github.com/SnowRipple/crane/utils"
	flags "github.com/jessevdk/go-flags"
	"github.com/mitchellh/cli"
	"strings"
)

type EnterCommand struct {
	Ui     cli.Ui
	Config config.TomlConfig
}

type enterOptions struct {
	ForceImage bool `short:"f" long:"force" description:"If chosen, crane will assume that the chosen image already exists in the host system(useful for offline mode)" `
}

func (c *EnterCommand) Help() string {
	helpText := `

  Presents the user with the interactive command line prompt inside a chosen container (you can enter only one container at a time).

  Usage: crane enter <containerName>
  In case of daemonized containers it is necessary to "start" them first before trying to enter them.

  `
	return strings.TrimSpace(helpText)
}

//Enter the container and present the user with the interactive shell
//For daemonized containers we ssh into the container
//For not deamonized containers we run (docker run) them with shell command. Effectively this command should be used to start non-daemonized containers.
func (c *EnterCommand) Run(arguments []string) int {

	logger.Debug("Entered enter command...")

	cmdFlags := flag.NewFlagSet("enter", flag.ContinueOnError)
	cmdFlags.Usage = func() { c.Ui.Output(c.Help()) }

	var options enterOptions

	arguments, err := flags.ParseArgs(&options, arguments)
	if err != nil {
		logger.Debug("Failed to extract start command flags for following arguments:\n%v", arguments)
	}

	checkArgumentsValidity(arguments)

	requestedContainerName := arguments[0]

	//Find the requested container config and state
	requestedContainerConfig, requestedContainerState := utils.GetContainerConfigAndState(c.Config, requestedContainerName, true, false) //it might not be present in the state file since in case of non-daemonized containers we might have to create them first

	if requestedContainerConfig.Daemonized { //ssh into it and provide the user with an interactive shell
		ssh.SshConnect(requestedContainerState.IP, requestedContainerConfig.Username, requestedContainerConfig.Password, constants.SHELL_COMMAND)
	} else { //run the container and provide the user with an interactive shell
		//Needs tty allocated
		dockerCommand := container.BuildRunCommand(requestedContainerConfig, requestedContainerName, true, true)

		if !options.ForceImage {
			buildImageCommand := BuildImageCommand{Ui: c.Ui}
			buildImageCommand.BuildImageIfNeeded(requestedContainerConfig)
		} else {
			logger.Debug("Force Image option detected. Will use host's system image.")
		}
		dockerCommand = append(dockerCommand, constants.SHELL_COMMAND)

		executer.ExecuteCommand(dockerCommand)

		id := io.GetContainerIdFromFile(requestedContainerName)

		//Update the state file
		io.UpdateStateFile(map[string]container.StateContainer{requestedContainerName: container.StateContainer{ID: id, IP: constants.NOT_DAEMONIZED_IP}})
	}

	return 0
}

//Checks if provided arguments are valid.
func checkArgumentsValidity(arguments []string) {

	argumentCount := len(arguments)

	if argumentCount == 0 {
		logger.Fatal("No container name provided.Please correct")
	} else if argumentCount > 1 {
		logger.Fatal("Too many arguments provided. You can enter only one container at a time.Please correct.")
	}
}

func (c *EnterCommand) Synopsis() string {
	return "Enter chosen container."
}
