package utils

import (
	"github.com/SnowRipple/crane/config"
	"github.com/SnowRipple/crane/container"
)

//Returns requested container config.
func GetRequestedContainerConfig(containers map[string]container.Container, containerName string, throwError bool) container.Container {

	requestedContainer, exists := containers[containerName]

	if throwError && !exists {
		logger.Fatalf("Chosen container:%q does not exist in the configuration file.Please correct.", containerName)
	}

	return requestedContainer

}

//Returns requested container state.
func GetRequestedContainerState(containers map[string]container.StateContainer, containerName string, throwError bool) container.StateContainer {
	requestedContainer, exists := containers[containerName]
	if throwError {
		if !exists {
			logger.Fatalf("Chosen container:%q does not exist in the state file.Please correct.", containerName)
		} else if requestedContainer.ID == "" {
			logger.Fatalf("Chosen container:%q has empty ID.", containerName)
		} else if requestedContainer.IP == "" {
			logger.Fatalf("Chosen container:%q has empty IP.", containerName)
		}
	}
	return requestedContainer
}

//Returns requested container config and state.
func GetContainerConfigAndState(config config.TomlConfig, containerName string, throwErrorConfig, throwErrorState bool) (configContainer container.Container, stateContainer container.StateContainer) {
	configContainer = GetRequestedContainerConfig(config.CraneConfig.Containers, containerName, throwErrorConfig)
	stateContainer = GetRequestedContainerState(config.CraneState.StateContainers, containerName, throwErrorState)
	return configContainer, stateContainer
}
