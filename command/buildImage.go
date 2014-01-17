package command

import (
	"flag"
	"fmt"
	"github.com/SnowRipple/crane/command/executer"
	"github.com/SnowRipple/crane/constants"
	"github.com/SnowRipple/crane/container"
	"github.com/SnowRipple/crane/utils"
	flags "github.com/jessevdk/go-flags"
	"github.com/mitchellh/cli"
	"strings"
)

// BuildImageCommand builds an image using Dockerfile providedin the Cranefile.toml
type BuildImageCommand struct {
	Ui         cli.Ui
	Containers map[string]container.Container
}

func (c *BuildImageCommand) Help() string {
	helpText := `
    Usage: crane build [options] <containerName1> <containerName2>
            
    Builds docker images specified by user and defined in the Cranefile.
  
  Options:

  -a(--all) Builds all images defined in the Cranefile.
    `
	return strings.TrimSpace(helpText)
}

func (c *BuildImageCommand) Run(containers []string) int {

	logger.Debug("Entered build command, containers provided are:\n%v", containers)

	var options constants.CommonFlags

	cmdFlags := flag.NewFlagSet("build", flag.ContinueOnError)
	cmdFlags.Usage = func() { c.Ui.Output(c.Help()) }

	containers, err := flags.ParseArgs(&options, containers)

	if err != nil {
		logger.Fatalf("Failed to parse options of the build command due to error: %v", err)
	}

	if options.All { //Build all images
		for containerName, _ := range c.Containers {
			containers = append(containers, containerName)
		}
		logger.Debug("All Dockerfiles defined in the Cranefile will be build for containers:\n%v", containers)
	}

	var images []string

	for _, chosenContainerName := range containers {
		chosenContainerConfig := utils.GetRequestedContainerConfig(c.Containers, chosenContainerName, true) //throw an error if container is not found

		c.buildImage(chosenContainerConfig.Dockerfile, chosenContainerConfig.Image)
		images = append(images, chosenContainerConfig.Image)
	}

	utils.PrintCommandOutput([]byte(fmt.Sprintf("Successfully build following images:\n%v", images)))
	return 0
}

//Checks if an image needs to be build from Dockerfile.If yes then the image will be build.
func (c *BuildImageCommand) BuildImageIfNeeded(container container.Container) {
	if isDockerfileBuildNeeded(container.Image) {
		c.buildImage(container.Dockerfile, container.Image)
	}

}

//Builds a docker image based on Dockerfile provided in the Cranefile.
func (c *BuildImageCommand) buildImage(dockerfilePath, imageName string) {

	logger.Debug("Building image from the Dockerfile...")

	//"-t" means build with specified name
	buildCommand := []string{constants.DOCKER, constants.BUILD, "-t", imageName, dockerfilePath}

	buildBytes, err := executer.GetCommandOutput(buildCommand)

	buildResults := utils.ExtractContainerMessage(buildBytes, err)
	logger.Debug("The building process results:\n%s", buildResults)
	if err != nil {
		logger.Fatal("Error when trying to build image %q using the Dockerfile located in %q:\n%q", imageName, dockerfilePath, buildResults)
	}
}

//Checks if the Dockerfile build is needed.
//1. Check if image exists in the host system.If yes, then Dockerfile build is not needed.
//2. Check if image exists in the docker public repository.If yes, then Dockerfile build is not needed.
//3.If 1. and 2. are false then the Dockerfile build is needed.
func isDockerfileBuildNeeded(imageName string) bool {

	//Check if image is present in the host system.

	if checkIfImageExists(imageName) || checkIfImagePresentInRepository(imageName) {
		return false //Image exists so Dockerfile build is not needed.
	}

	return true
}

//Checks if a given image is present in the docker public repository.
func checkIfImagePresentInRepository(imageName string) bool {

	logger.Debug("Checking if the image %s is present in the docker public repository...", imageName)

	checkPublicRepositoryCommand := []string{constants.DOCKER, constants.SEARCH, imageName}

	imageSearchBytes, err := executer.GetCommandOutput(checkPublicRepositoryCommand)
	searchMessage := utils.ExtractContainerMessage(imageSearchBytes, err)
	logger.Debug("Search command results:\n%s", searchMessage)
	if err != nil {
		logger.Fatal("Error when searching for the image:"+imageName+" in the docker public repository:\n%s", searchMessage)
	}

	if strings.Contains(searchMessage, "Found 0 results matching your query") {
		logger.Debug("Image: %s does not exists in the public repository", imageName)
		return false
	} else {
		logger.Debug("Image: %s exists in the docker public repository.", imageName)
	}
	return true
}

//Checks if a given image exists in the host's system.
func checkIfImageExists(imageName string) bool {

	logger.Debug("Checking if the image %s is present in the host system...", imageName)

	imagesCommand := []string{constants.SUDO, constants.DOCKER, constants.IMAGES, imageName}

	imageCommandOutputBytes, err := executer.GetCommandOutput(imagesCommand)

	imagesOutput := utils.ExtractContainerMessage(imageCommandOutputBytes, err)
	logger.Debug("Images command results:\n%s", imagesOutput)
	if err != nil {
		logger.Fatal("Error when searching for the image:"+imageName+" in the docker public repository:\n%s", imagesOutput)
	}
	if strings.Contains(string(imageCommandOutputBytes), imageName) {
		logger.Debug("Image %s exists in the host system.", imageName)
		return true
	} else {
		logger.Debug("Image %s does NOT exist in the host system.", imageName)
	}
	return false
}

func (c *BuildImageCommand) Synopsis() string {
	return "Build image(s) from chosen containers."
}
