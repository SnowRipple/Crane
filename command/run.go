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
	"strconv"
	"strings"
)

const COMMANDS_ARGUMENT_COUNT = 2

// RunCommand executes commands per container..
type RunCommand struct {
	Ui     cli.Ui
	Config config.TomlConfig
}

func (c *RunCommand) Help() string {
	helpText := `
  Usage: crane run [options] <containerName1>:<command1>,<command2> <containerName2>

  Execute commands defined in the Cranefile.toml.
 
  Usage: crane run [options] <containerName1>:#"<command1>,<command2> <containerName2>"

  Execute own commands.
  
  Please note different delimiters in both cases.
  The user can mixture both methods(use own and Cranefile commands) within a single crane command.

  Options:

  -u (--update) Transforms (commits) all containers into immutable images that will replcae existing images.
  `
	return strings.TrimSpace(helpText)
}

//Run chosen commands (own or Cranefile commands) per container.
func (c *RunCommand) Run(commandArguments []string) int {

	var (
		options    constants.CommonFlags
		imageQueue = utils.NewQueue(0)
	)

	commandArguments, err := flags.ParseArgs(&options, commandArguments)
	if err != nil {
		logger.Fatalf("Failed to parse runall flags for following CLI arguments:\n%v", commandArguments)
	}

	logger.Debug("Entered run command...")

	cmdFlags := flag.NewFlagSet("run", flag.ContinueOnError)
	cmdFlags.Usage = func() { c.Ui.Output(c.Help()) }

	if len(commandArguments) == 0 { //only RUN
		logger.Fatal("Missing arguments.Please specify containers and corrresponding commands or did you mean runall?")
	}

	//Freeze containers if required.
	if len(options.Save) > 0 {
		imageQueue = extractNewImageNames(options.Save)
		logger.Debug("Save option detected.Following new images will be created:\n%v", imageQueue)
	}

	for index, argument := range commandArguments {

		//Extract the container name
		logger.Debug("\nThe %d argument is %s\n", index, argument)
		arguments := strings.Split(argument, constants.COMMANDS_DELIMITER)
		if len(arguments) != 2 {
			logger.Fatal("Wrong arguments format.Please correct.")
		}

		chosenContainerName := arguments[0]
		enteredCommands := arguments[1]

		//Get requested container configuration
		requestedContainerConfig, requestedContainerState := utils.GetContainerConfigAndState(c.Config, chosenContainerName, true, false) //It must be in the config but not necessarily in the state(for new not daemonized containers)
		command := buildContainerCommand(requestedContainerConfig, chosenContainerName, enteredCommands)

		runCommandInContainer(c.Ui, requestedContainerConfig, requestedContainerState, chosenContainerName, command, options.ForceImage)

		//Freeze container into image if requested (requires updated state file)
		if options.Update {
			logger.Debug("Overwriting existing image for container %q...", chosenContainerName)
			freezeCommand := FreezeCommand{Ui: c.Ui, Config: config.ReadConfig()}
			freezeCommand.Run([]string{chosenContainerName})
		} else if imageQueue.Length() > 0 {
			newImageName := imageQueue.Pop().Value
			freezeCommand := FreezeCommand{Ui: c.Ui, Config: config.ReadConfig()}
			logger.Debug("Existing container %q will be committed as an image %q", chosenContainerName, newImageName)
			freezeCommand.Run([]string{chosenContainerName + constants.FREEZE_DELIMITER + newImageName})
		}
	}
	return 0
}

//Extracts image names for images build from running containers.

func extractNewImageNames(newImagesNames string) *utils.Queue {

	var queue *utils.Queue
	if strings.Contains(newImagesNames, ",") { //multiple images
		newImagesSlice := strings.Split(newImagesNames, ",")
		queue = utils.NewQueue(len(newImagesSlice))
		for _, imageName := range newImagesSlice {
			queue.Push(&utils.Node{imageName})
		}
	} else { //single image name

		queue = utils.NewQueue(1)
		queue.Push(&utils.Node{newImagesNames})
	}

	return queue
}

//Extract command line commands
func extractCLCommands(initialCommand string) []string {

	var craneCommands []string

	if strings.Contains(initialCommand, ",") { //multiple commands
		craneCommands = strings.Split(initialCommand, ",")
	} else {
		craneCommands = append(craneCommands, initialCommand)
	}

	logger.Debug("\nExtracted following crane commands from the user input:\n%v", craneCommands)

	return craneCommands
}

//Build a command that will be executed inside a container.
func buildContainerCommand(requestedContainerConfig container.Container, requestedContainerName, initialCommand string) string {

	//Extract commands
	if strings.Index(initialCommand, constants.OWN_COMMANDS_DELIMITER) == 0 { //own commands
		logger.Debug("User's own commands detected.")
		return strings.TrimLeft(initialCommand, constants.OWN_COMMANDS_DELIMITER)
	} else { //Cranefile commands
		logger.Debug("Cranefile commands detected.")

		containerCommands := requestedContainerConfig.Commands
		if len(containerCommands) == 0 {
			logger.Fatalf("Either the container %q does not exist in the Cranefile or it has not any commands defined in the Cranefile.Please correct.", requestedContainerName)
		}
		craneCommands := extractCLCommands(initialCommand)

		return buildCommandList(craneCommands, containerCommands)
	}

	return "Unreachable reached...End of the world approaching..." //This is unreachable but go compiler requires it...

}

//Builds a list of commands to be executed per a container. Commands come from the Cranefile.
func buildCommandList(craneCommands []string, containerCommands [][]string) string {

	var commandList string

	for _, craneCommand := range craneCommands {
		for secondIndex, cranefileCommandPair := range containerCommands {
			if len(cranefileCommandPair) != COMMANDS_ARGUMENT_COUNT {
				logger.Fatal("Wrong amount of command  arguments specified for the pair nr " + strconv.Itoa(secondIndex) + ": Expected " + strconv.Itoa(COMMANDS_ARGUMENT_COUNT) + ", Actual " + strconv.Itoa(len(cranefileCommandPair)) + ". Please correct(if you need to specify multiple commands please use \"<command1>;<command2>\" format.")
			}
			if craneCommand == cranefileCommandPair[0] { //Found chosen command
				if len(commandList) > 1 {
					commandList = commandList + ";" //Append delimiter
				}
				commandList = commandList + cranefileCommandPair[1]
				break
			}
			logger.Debug("\nFinal commands after iteration %d is %s\n", secondIndex, commandList)
		}
	}

	return commandList
}

//Run a specified command in a specified container.Updates the state file.
func runCommandInContainer(ui cli.Ui, containerConfig container.Container, containerState container.StateContainer, containerName, command string, useHostImage bool) {

	if containerConfig.Daemonized {
		ssh.SshConnect(containerState.IP, containerConfig.Username, containerConfig.Password, constants.SHELL_COMMAND+" "+constants.SHELL_STRING_OPTION+" \""+command+"\"")
	} else { //Not daemonized

		dockerCommand := container.BuildRunCommand(containerConfig, containerName, false, true) //Needs cidfile to store id

		if !useHostImage {
			buildImageCommand := BuildImageCommand{Ui: ui}
			buildImageCommand.BuildImageIfNeeded(containerConfig)

		} else {
			logger.Debug("Force option detected, will use host system image.")
		}
		//Append "/bin/bash -c" and command
		dockerCommand = append(dockerCommand, constants.SHELL_COMMAND)
		dockerCommand = append(dockerCommand, constants.SHELL_STRING_OPTION)
		dockerCommand = append(dockerCommand, command)

		outputBytes, err := executer.GetCommandOutput(dockerCommand)
		if err != nil {
			logger.Fatal("Error during \"run\" command with non daemonized container:", utils.ExtractContainerMessage(outputBytes, err))
		}

		utils.PrintCommandOutput(outputBytes)

		id := io.GetContainerIdFromFile(containerName)

		//Update the state file
		io.UpdateStateFile(map[string]container.StateContainer{containerName: {ID: id, IP: constants.NOT_DAEMONIZED_IP}}) //non-daemonized have no ip
	}
}

func (c *RunCommand) Synopsis() string {
	return "Execute commands per container."
}
