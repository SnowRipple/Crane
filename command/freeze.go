package command

import (
	"flag"
	"github.com/SnowRipple/crane/command/executer"
	"github.com/SnowRipple/crane/config"
	"github.com/SnowRipple/crane/constants"
	"github.com/SnowRipple/crane/container"
	"github.com/SnowRipple/crane/utils"
	flags "github.com/jessevdk/go-flags"
	"github.com/mitchellh/cli"
	"strings"
)

const (
	NUMBER_OF_PARAMS = 2
)

// FreezeCommand creates an example configuration file Cranefile.toml.
type FreezeCommand struct {
	Ui     cli.Ui
	Config config.TomlConfig
}

func (c *FreezeCommand) Help() string {
	helpText := `

  Usage: crane freeze <containerName1> <containerName2>
    
Transforms (using docker lingo "commits") containers into immutable images using names defined in Cranefile (overwrites existing images). 

Usage: crane freeze <containerName1>::<imageName1> <containerName2>::<imageName2>
    
    Transforms (using docker lingo "commits") containers into immutable images with chosen names.
    
    
    Usage: crane freeze <containername1> <containerName2>::<imageName2>

Both options can be used within a single crane command.

Available options:

-a (--all) : Freeze all containers defined in the Cranefile.Useful when you want to save all your work done on different containers.`
	return strings.TrimSpace(helpText)
}

//Transform containers into immutable images.
func (c *FreezeCommand) Run(containerNames []string) int {

	var options constants.CommonFlags

	logger.Debug("Entered freeze command...")

	cmdFlags := flag.NewFlagSet("freeze", flag.ContinueOnError)
	cmdFlags.Usage = func() { c.Ui.Output(c.Help()) }

	containerNames, err := flags.ParseArgs(&options, containerNames)

	if err != nil {
		logger.Fatalf("Failed to parse freeze flags for following CLI arguments:\n%v", containerNames)
	}

	if options.All { //Freeze all containers defined in the Cranefile
		for containerName, _ := range c.Config.CraneConfig.Containers {
			containerNames = append(containerNames, containerName)
		}
		logger.Debug("Commit all containers defined in theCranefile:\n%v", containerNames)
	} else if len(containerNames) == 0 {
		logger.Fatal("No arguments provided for the freeze command.Please correct.")
	}
	for _, chosenContainerName := range containerNames {

		containerName, imageName := extractContainerImageNames(c.Config.CraneConfig.Containers, chosenContainerName)
		containerState := utils.GetRequestedContainerState(c.Config.CraneState.StateContainers, containerName, true) //It must exist in the state file to be frozen

		logger.Debug("Committing container %q into image %q", containerName, imageName)

		dockerCommand := []string{constants.DOCKER, constants.COMMIT, containerState.ID, imageName}

		outputBytes, err := executer.GetCommandOutput(dockerCommand)
		if err != nil {
			logger.Fatal("Error during \"freeze\" command:", utils.ExtractContainerMessage(outputBytes, err))
		}

		imageId := strings.TrimSpace(string(outputBytes))
		logger.Notice("Successfully froze container %q into image %q with id %q", containerName, imageName, imageId)
	}
	return 0
}

//Extracts the name of a container to be frozen and the new name of frozen container(image).
func extractContainerImageNames(containers map[string]container.Container, chosenContainer string) (currentContainerName, imageName string) {

	if strings.Contains(chosenContainer, constants.FREEZE_DELIMITER) { //User specified image names
		imageParameters := strings.Split(chosenContainer, constants.FREEZE_DELIMITER)

		if len(imageParameters) != 2 {
			logger.Fatalf("Invalid number of image parameters provided.Expected %d, Actual %d.Parameters are: \n%v", NUMBER_OF_PARAMS, len(imageParameters), imageParameters)
		}

		currentContainerName = imageParameters[0]
		imageName = imageParameters[1]

	} else { //Get image name from the Cranefile and overwrite it

		currentContainerName = chosenContainer

		currentContainerConfig := utils.GetRequestedContainerConfig(containers, chosenContainer, true)

		imageName = currentContainerConfig.Image
	}
	return currentContainerName, imageName
}

func (c *FreezeCommand) Synopsis() string {
	return "Freeze chosen container(s)."
}
