package command

import (
	"flag"
	"github.com/SnowRipple/crane/command/executer"
	"github.com/SnowRipple/crane/config"
	"github.com/SnowRipple/crane/constants"
	"github.com/SnowRipple/crane/io"
	"github.com/SnowRipple/crane/utils"
	"github.com/mitchellh/cli"
	"strings"
)

// DestroyCommand creates an example configuration file Cranefile.toml.
type DestroyCommand struct {
	Ui     cli.Ui
	Config config.TomlConfig
}

func (c *DestroyCommand) Help() string {
	helpText := `
  Usage: crane destroy
  
  Kills and removes all running containers which were defined in the Cranefile beforehand.
  
  Usage: crane destroy <containerName1> <containerName2>
  
  Kills and removes all specified containers.`
	return strings.TrimSpace(helpText)
}

//Pulls all images defined in the Cranefile
func (c *DestroyCommand) Run(arguments []string) int {

	var (
		containersIdsToBeDestroyed   []string
		containersNamesToBeDestroyed []string
	)

	logger.Debug("Entered destroy command...")

	cmdFlags := flag.NewFlagSet("destroy", flag.ContinueOnError)
	cmdFlags.Usage = func() { c.Ui.Output(c.Help()) }

	dockerKillCommand := []string{constants.DOCKER, constants.KILL}
	dockerRemoveCommand := []string{constants.DOCKER, constants.REMOVE}
	killThemAll := true //Kill all containers

	if len(arguments) > 0 { //kill only specified containers
		killThemAll = false
		logger.Debug("Only user provided containers will be destroyed:\n%v", arguments)
	}

	stateContainers := c.Config.CraneState.StateContainers

	if len(stateContainers) == 0 {
		logger.Fatalf("There are no containers in the state file hence no containers will be destroyed.You can destroy only containers that were created by the crane.")
	}

	//append containers ids to the docker kill command
	for containerName, stateContainer := range stateContainers {
		if killThemAll { //Kill all
			containersIdsToBeDestroyed, containersNamesToBeDestroyed = addToBeDestroyedList(containerName, stateContainer.ID, containersIdsToBeDestroyed, containersNamesToBeDestroyed)
		} else { //Kill specific containers only
			for _, containerToBeDeleted := range arguments {
				if containerToBeDeleted == containerName {
					containersIdsToBeDestroyed, containersNamesToBeDestroyed = addToBeDestroyedList(containerName, stateContainer.ID, containersIdsToBeDestroyed, containersNamesToBeDestroyed)
					break
				}
			}
		}
	}

	dockerKillCommand = append(dockerKillCommand, containersIdsToBeDestroyed...)
	dockerRemoveCommand = append(dockerRemoveCommand, containersIdsToBeDestroyed...)

	if len(containersIdsToBeDestroyed) == 0 {
		logger.Fatal("No such containers were found in the state file  hence they cannott be deleted.Please note that you can destroy only containers created by the crane.")
	}

	logger.Debug("Following containers will be destroyed:\n%v", containersNamesToBeDestroyed)

	//First Kill, then Remove
	killContainers(dockerKillCommand)
	removeContainers(dockerRemoveCommand)

	//Remove destroyed containers from the state file.
	io.RemoveStateContainers(containersNamesToBeDestroyed)

	return 0
}

//Add container to the lists of containers to be killed and removed
func addToBeDestroyedList(containerName, containerId string, containersIdsToBeDestroyed, containersNamesToBeDestroyed []string) ([]string, []string) {

	logger.Debug("Following container will be destroyed: %q with ID %q", containerName, containerId)
	containersIdsToBeDestroyed = append(containersIdsToBeDestroyed, containerId)
	containersNamesToBeDestroyed = append(containersNamesToBeDestroyed, containerName)

	return containersIdsToBeDestroyed, containersNamesToBeDestroyed

}

//Kill running containers. If containers are not running nothing will happen.
func killContainers(killCommand []string) {

	killedContainersBytes, err := executer.GetCommandOutput(killCommand)
	if err != nil {
		logger.Fatal("Error when trying to destroy container(s):", utils.ExtractContainerMessage(killedContainersBytes, err))
	}

	killedContainers := strings.TrimSpace(string(killedContainersBytes))
	logger.Notice("Kill command output:\n%v", killedContainers)

}

//Remove containers from the system
func removeContainers(removeCommand []string) {

	removedContainersBytes, err := executer.GetCommandOutput(removeCommand)
	if err != nil {
		logger.Fatal("Error when trying to destroy container(s):", utils.ExtractContainerMessage(removedContainersBytes, err))
	}

	removedContainers := strings.TrimSpace(string(removedContainersBytes))
	logger.Notice("Remove command output:\n%v", removedContainers)
}

func (c *DestroyCommand) Synopsis() string {
	return "Kill and remove containers."
}
