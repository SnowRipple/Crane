package command

import (
	"flag"
	"github.com/SnowRipple/crane/config"
	"github.com/SnowRipple/crane/constants"
	"github.com/SnowRipple/crane/container"
	"github.com/SnowRipple/crane/utils"
	flags "github.com/jessevdk/go-flags"
	"github.com/mitchellh/cli"
	"strconv"
	"strings"
)

// RunallCommand creates an example configuration file Cranefile.toml.
type RunallCommand struct {
	Ui     cli.Ui
	Config config.TomlConfig
}

func (c *RunallCommand) Help() string {
	helpText := `
  This command is used to automate multiple run commands.

  Usage: crane runall
  
  Executes all commands defined in the Cranefile.

  Options:

	crane runall -c="<CranefileCommand1>;<CranefileCommand2>"
  Runs specified Cranefile commands in all containers.
  
  crane runall -o="<ownCommand>;<ownCommand>"
  Runs own (not defined in the Cranefile but typed on the commandline) commands in all containers.
  
  crane runall -l="<containerName1>;<containerName2>"
  Runs all Cranefile's commands in specified containers.
  
  crane runall -u

  Transforms(commits) all containers into immutable images that will replace existing images.
  
  `
	return strings.TrimSpace(helpText)
}

//Execute all commands defined in the Cranefile.
func (c *RunallCommand) Run(arguments []string) int {

	var (
		options       constants.CommonFlags
		freezeCommand = FreezeCommand{Ui: c.Ui, Config: c.Config}
	)
	logger.Debug("Entered runall command..")

	cmdFlags := flag.NewFlagSet("runall", flag.ContinueOnError)
	cmdFlags.Usage = func() { c.Ui.Output(c.Help()) }

	arguments, err := flags.ParseArgs(&options, arguments)
	if err != nil {
		logger.Fatalf("Failed to parse runall options due to error:", err)
	}

	allContainersConfig := c.Config.CraneConfig.Containers
	allContainersState := c.Config.CraneState.StateContainers
	runCommand := options.RunAllCommand
	runOwnCommand := options.RunAllOwnCommand
	runContainer := options.RunAllContainer

	//Only one runall option can be specified in a command.If you need more use scripts
	checkRunAllOptionsValidity([]string{runCommand, runOwnCommand, runContainer})

	//Get a list of containers to be run if required
	if len(runContainer) > 0 {
		allContainersConfig = extractChosenContainers(allContainersConfig, runContainer)
		logger.Debug("Containers specified by the user are:\n %v", allContainersConfig)
	}

	for containerName, containerConfig := range allContainersConfig {

		command := buildCommand(runCommand, runOwnCommand, containerConfig.Commands)

		containerState := utils.GetRequestedContainerState(allContainersState, containerName, false)

		runCommandInContainer(c.Ui, containerConfig, containerState, containerName, command, options.ForceImage)

		if options.Update { //Update images if requested
			logger.Debug("Overwriting existing image for container %q...", containerName)
			freezeCommand.Run([]string{containerName})
		}
	}

	return 0
}

//Builds an overall command that will be run in a single container
func buildCommand(runCommand, runOwnCommand string, cranefileCommands [][]string) string {

	var command string
	if len(runCommand) > 0 { //Run specified Cranefile commands across all containers
		command = extractCranefileCommands(cranefileCommands, runCommand)
	} else if len(runOwnCommand) > 0 { //Run own commands across all containers
		command = runOwnCommand
	} else { //Run all Cranefile commands
		command = extractCranefileCommands(cranefileCommands, "") //Empty string so all commands will be used
	}

	logger.Debug("Final command is %q", command)
	return command
}

//Returns containers that were specified by the user.
func extractChosenContainers(allContainers map[string]container.Container, runContainer string) map[string]container.Container {

	chosenContainersList := strings.Split(runContainer, constants.INTERNAL_DELIMITER)
	var chosenContainers = map[string]container.Container{}

	for _, chosenContainerName := range chosenContainersList {
		for containerName, container := range allContainers {
			if chosenContainerName == containerName {
				chosenContainers[chosenContainerName] = container
				break
			}
		}
	}
	return chosenContainers
}

//Extract commands specified by the user from the Cranefile.
//If userChosenCraneCommands are empty, all commands will be returned.
func extractCranefileCommands(commands [][]string, userChosenCraneCommands string) string {

	var (
		command           string
		userCraneCommands []string
	)

	if len(userChosenCraneCommands) > 0 { //Use commands specified by the user only

		userCraneCommands = strings.Split(userChosenCraneCommands, constants.INTERNAL_DELIMITER)

		for _, userCraneCommand := range userCraneCommands {
			command = appendRunAllCommand(commands, command, userCraneCommand)
		}
	} else { //No user commands specified so grab all commands

		command = appendRunAllCommand(commands, command, "") //Empty so all cranefile commands will be used
	}

	return command
}

//Finds and appends a cranefile command to a list of commands to be executed.
func appendRunAllCommand(commands [][]string, commandList, userCraneCommand string) string {

	for index, cranefileCommandPair := range commands {
		if len(cranefileCommandPair) != COMMANDS_ARGUMENT_COUNT {
			logger.Fatal("Wrong amount of command  arguments specified for the pair nr " + strconv.Itoa(index) + ": Expected " + strconv.Itoa(COMMANDS_ARGUMENT_COUNT) + ", Actual " + strconv.Itoa(len(cranefileCommandPair)) + ". Please correct(if you need to specify multiple commands please use \"<command1>;<command2>\" format.")
		}

		if len(userCraneCommand) == 0 { //All Cranefile commands
			commandList = appendCommand(commandList, cranefileCommandPair[1])
		} else if userCraneCommand == cranefileCommandPair[0] { //Specific Cranefile commands
			return appendCommand(commandList, cranefileCommandPair[1])
		}
	}
	return commandList
}

//Appends a new command to the list of commands.
func appendCommand(commandList, part string) string {

	if len(commandList) > 1 {
		commandList = commandList + constants.INTERNAL_DELIMITER //Append delimiter
	}

	return commandList + part
}

// The user can specify only one option at a time or no option a all in which case all commands in all containers are executed.
// If the user wants to use more than one option in a sequential order he/she use multiple "runall" commands.
func checkRunAllOptionsValidity(runallCommands []string) {

	optionPresent := false

	for _, command := range runallCommands {
		if len(command) > 0 && !optionPresent {
			optionPresent = true
		} else if len(command) > 0 && optionPresent {
			logger.Fatal("Can't combine multiple options with the \"runall\" command.Please use only one of the following options at a time : \"-c\",\"-o\",\"-l\".")
		}
	}
}

func (c *RunallCommand) Synopsis() string {
	return "Execute all commands."
}
