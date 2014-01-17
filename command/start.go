package command

import (
	"flag"
	"github.com/SnowRipple/crane/command/executer"
	"github.com/SnowRipple/crane/config"
	"github.com/SnowRipple/crane/constants"
	"github.com/SnowRipple/crane/container"
	"github.com/SnowRipple/crane/io"
	"github.com/SnowRipple/crane/utils"
	flags "github.com/jessevdk/go-flags"
	"github.com/mitchellh/cli"
	"os/exec"
	"strings"
	"time"
)

// StartCommand initializes all daemonized containers.
type StartCommand struct {
	Ui     cli.Ui
	Config config.TomlConfig
}

func (c *StartCommand) Help() string {
	helpText := `
    Usage: crane start [options] <containerName1> <containerName2>
      
    Initialize daemonized containers.
Options:

  -a(--all) : Starts all daemonized containers defined in the Cranefile.
    -f(--force) : Crane assumes that the image already exists in the host system and does not attempt to download it from the docker public repository (useful for offline mode).
    `

	return strings.TrimSpace(helpText)
}

//Initailize all daemonized containers. The initialization process includes creating mountpoints and starting sshd process to listen for incoming ssh connections.
func (c *StartCommand) Run(chosenContainers []string) int {
	var (
		stateContainers = map[string]container.StateContainer{}
		dockerCommand   []string
		options         constants.CommonFlags
	)

	logger.Debug("Entered start command..")

	cmdFlags := flag.NewFlagSet("start", flag.ContinueOnError)
	cmdFlags.Usage = func() { c.Ui.Output(c.Help()) }

	if len(chosenContainers) == 0 {
		logger.Fatalf("Not enough arguments provided for the start command.Please correct.")
	}
	chosenContainers, err := flags.ParseArgs(&options, chosenContainers)
	if err != nil {
		logger.Fatalf("Failed to extract flags for the start command for the following arguments:\n%v", chosenContainers)
	}

	for containerName, containerConfig := range c.Config.CraneConfig.Containers {

		if containerConfig.Daemonized == false {
			continue //Start only daemonized containers
		}

		if !options.All && !isThisContainerChosen(containerName, chosenContainers) {
			continue //Start only chosen containers
		}

		dockerCommand = container.BuildRunCommand(containerConfig, containerName, false, false) //no need for cidfile since in case of daemonized it is returned automatically through stdout

		if !options.ForceImage {
			buildImageCommand := BuildImageCommand{Ui: c.Ui}
			buildImageCommand.BuildImageIfNeeded(containerConfig)
		} else {
			logger.Debug("Force Image option detected. Will use host's image only")
		}
		//When the container is daemonized we need to be able to access it through the ssh.
		//Hence we need to start sshd process to listen for the incoming ssh connections.

		dockerCommand = append(dockerCommand, constants.SHELL_COMMAND)
		dockerCommand = append(dockerCommand, constants.SHELL_STRING_OPTION)
		dockerCommand = append(dockerCommand, constants.SSHD_COMMAND)
		//Run the container

		containerIdBytes, err := executer.GetCommandOutput(dockerCommand)
		if err != nil {
			logger.Fatal("Error starting daemonized container:", utils.ExtractContainerMessage(containerIdBytes, err))
		} else {
			logger.Notice("Successfully started container %q...", containerName)
		}

		time.Sleep(1 * time.Second) //wait for the docker to do it's magic.

		//Get Container ID
		containerId := strings.TrimSpace(string(containerIdBytes))
		logger.Debug("Container %q ID is %q", containerName, containerId)
		//Get Container IP address
		ipAddress := getContainerIP(containerId)
		logger.Debug("Container %q IP is %q", containerName, ipAddress)

		stateContainers[containerName] = container.StateContainer{ID: containerId, IP: ipAddress}

	}
	if len(stateContainers) == 0 {
		logger.Notice("No containers in the Cranefile match provided criteria hence no containers were started.")
	} else {
		io.UpdateStateFile(stateContainers)
	}

	return 0
}

//Extracts containers's ip address using docker inspect command.
func getContainerIP(containerID string) string {

	//We have to pipe multiple commands in order to get IP address of a container
	inspectCommand := exec.Command(constants.SUDO, constants.DOCKER, constants.INSPECT, containerID)
	grepCommand := exec.Command(constants.GREP, "IPAddress")
	cutCommand := exec.Command("cut", "-d\"", "-f4")

	ipAddress := strings.TrimSpace(executer.PipeCommands(inspectCommand, grepCommand, cutCommand))
	if len(ipAddress) == 0 {
		logger.Fatal("Failed to obtain the IP Address for the container " + containerID + ". Invalid commands? Container is not able to run commands? Please investigate.")
	}
	logger.Debug(" Container %q has IP %q", containerID, ipAddress)

	return ipAddress
}

func (c *StartCommand) Synopsis() string {
	return "Initialize daemonized containers."
}

//Checks if a containerName is an element of the chosenContainers slice.
func isThisContainerChosen(containerName string, chosenContainers []string) bool {

	for _, chosenContainerName := range chosenContainers {
		if chosenContainerName == containerName {
			return true
		}
	}
	return false
}
