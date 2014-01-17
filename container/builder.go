package container

import (
	"github.com/SnowRipple/crane/constants"
	ownLog "github.com/SnowRipple/crane/logger"
	"strconv"
	"strings"
)

/*
 Docker options
*/
const (
	INTERACTIVE_OPTION     = "-i"
	TTY_OPTION             = "-t"
	VOLUME_OPTION          = "-v="
	NAME_OPTION            = "-name="
	DNS_OPTION             = "-dns="
	DAEMONIZED_OPTION      = "-d"
	WORKING_DIR_OPTION     = "-w="
	PORT_REDIRECT_OPTION   = "-p="
	PRIVILEDGED_OPTION     = "-privileged"
	BUILD_WITH_NAME_OPTION = "-t"
	CID_OPTION             = "-cidfile="

	MOUNTPOINTS_ARGUMENT_COUNT = 3
	PORTS_ARGUMENT_COUNT       = 2
)

var logger = ownLog.GetLogger()

//Builds a docker run command used by:
//->crane start - to start daemonized containers.
//->crane enter,run and runall - to run the non-deamonized containers.
func BuildRunCommand(container Container, containerName string, needsTTY, needsCidfile bool) []string {

	logger.Debug("Starting building run command...")
	dockerCommand := []string{constants.DOCKER, constants.RUN}

	//Closure to simplify append process
	addCommandPart := func(command string) { dockerCommand = append(dockerCommand, command) }

	//Add necessary options
	addCommandPart(INTERACTIVE_OPTION)
	addCommandPart(PRIVILEDGED_OPTION)

	if needsCidfile {
		fileName := constants.ID_FILE + containerName
		addCommandPart(CID_OPTION + fileName)
		logger.Debug("Using cidfile %q to store containerId temporarily.", fileName)
	}
	//Set up DNS if needed
	if len(strings.TrimSpace(container.Dns)) > 0 {
		addCommandPart(DNS_OPTION + "[" + container.Dns + "]")
		logger.Debug("Using DNS:" + container.Dns + ".")
	}

	//Set up current working directory if needed
	if len(strings.TrimSpace(container.Cwd)) > 0 {
		addCommandPart(WORKING_DIR_OPTION + container.Cwd)
		logger.Debug("Using working directory: " + container.Cwd + ".")
	}

	//Ports redirection
	if len(container.Ports) > 0 {
		portsCommands := buildPortsCommands(container.Ports)

		for _, portCommand := range portsCommands {
			addCommandPart(portCommand)
		}

	}

	//Daemonized?
	if container.Daemonized == true {
		addCommandPart(DAEMONIZED_OPTION)
		logger.Debug("Added daemonized option.")
	}
	if !container.Daemonized && needsTTY { //Allocate tty
		addCommandPart(TTY_OPTION)
		logger.Debug("Allocated tty.")
	}

	//Mount external directories
	if len(container.Mountpoints) > 0 {
		logger.Debug("Mountpoints detected, extracting...")
		mountpointCommands := buildMountpointCommands(container.Mountpoints)

		for _, mountpointCommand := range mountpointCommands {
			addCommandPart(mountpointCommand)
		}
	}

	//Image takes precedence over the dockerfile

	if len(strings.TrimSpace(container.Image)) == 0 {
		logger.Fatal("No image was specified in the Cranefile.Please correct.")
	} else {
		//Use existing image
		addCommandPart(container.Image)
		logger.Debug("Using Image:" + container.Image + ".")

	}

	logger.Debug("Final builded run command:\n%v", dockerCommand)

	return dockerCommand
}

//Parses through all mountpoints pairs provided in the Cranefile (if any) and builds docker volume option commands.
func buildMountpointCommands(mountpoints [][]string) []string {

	mountpointsCommands := []string{}

	for index := 0; index < len(mountpoints); index++ {
		mountpointPair := mountpoints[index]
		logger.Debug("\n\nMountpointPair is %v", mountpointPair)
		if len(mountpointPair) != MOUNTPOINTS_ARGUMENT_COUNT {
			logger.Fatal("Wrong amount of mountpoints arguments specified for the pair nr " + strconv.Itoa(index) + " : Expected " + strconv.Itoa(MOUNTPOINTS_ARGUMENT_COUNT) + ", Actual " + strconv.Itoa(len(mountpointPair)) + ".Please correct")
		}

		logger.Debug("\nHost mountpoint is %v and corresponding remote will be placed in %v as a %s filesystem\n", mountpointPair[0], mountpointPair[1], mountpointPair[2])
		mountpointsCommands = append(mountpointsCommands, VOLUME_OPTION+mountpointPair[0]+constants.COMMANDS_DELIMITER+mountpointPair[1]+constants.COMMANDS_DELIMITER+mountpointPair[2])

	}

	return mountpointsCommands
}

//Parses through all ports pairs provided in the Cranefile and builds docker port option commands.
func buildPortsCommands(ports [][]int) []string {

	portsCommands := []string{}

	for index := 0; index < len(ports); index++ {
		portsPair := ports[index]
		logger.Debug("\nPorts Pair is %v", portsPair)
		if len(portsPair) != PORTS_ARGUMENT_COUNT {
			logger.Fatal("Wrong amount of port arguments specified for the pair nr %d : Expected %d, Actual %d. Please correct", index, PORTS_ARGUMENT_COUNT, len(portsPair))
		}

		logger.Debug("\nPublic port is %d and the private port is %d", portsPair[0], portsPair[1])

		portsCommands = append(portsCommands, PORT_REDIRECT_OPTION+strconv.Itoa(portsPair[0])+constants.COMMANDS_DELIMITER+strconv.Itoa(portsPair[1]))
	}
	return portsCommands
}
