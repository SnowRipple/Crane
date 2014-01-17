package io

import (
	"bufio"
	"bytes"
	"github.com/SnowRipple/crane/constants"
	"github.com/SnowRipple/crane/container"
	log "github.com/SnowRipple/crane/logger"
	"io"
	"os"
	"strings"
)

const (
	DEFAULT_CONTAINER_IP   = "Squik squik..."
	DEFAULT_CONTAINER_ID   = "Muuu...muuuu"
	DEFAULT_CONTAINER_NAME = "How...how"

	CONTAINERS_HEADER       = "[containers]"
	STATE_CONTAINERS_HEADER = "[statecontainers]"

	ID_LINE = "ID ="
	IP_LINE = "IP ="
)

var logger = log.GetLogger()

//Creates new config and state files.
func Create() {

	CreateNewConfigFile()
	CreateNewStateFile()

}

//Create a new config file. If config file already exists it will be overwritten.
func CreateNewConfigFile() {

	var cranefileTemplate = []string{
		CONTAINERS_HEADER,
		"[containers.firstContainer]",
		"IMAGE = \"orobix/sshfs_startup_key2\"",
		"DOCKERFILE = \".\"",
		"GRAPHICAL = true",
		"DAEMONIZED = false",
		"CWD = \"/home/foo\" #Leave empty if not needed",
		"DNS = \"\" #Leave empty if not needed",
		"PASSWORD = \"orobix2013\"#Leave empty if not needed",
		"USERNAME = \"root\"",
		"PORTS = [[49153, 22]]#Remember that port 22 is required for daemonized containers for SSH communication",
		"MOUNTPOINTS=[]#Insert own mountpoints here",
		"COMMANDS=[[\"init\",\"echo orobix\"]]"}

	writeLines(cranefileTemplate, constants.CONFIGURATION_FILE)

	if _, err := CheckIfFileExists(constants.CONFIGURATION_FILE); err != nil {
		logger.Fatalf("Failed to create file %q due to error:%v", constants.CONFIGURATION_FILE, err)
	}
}

//Checks if a file exists.
func CheckIfFileExists(filename string) (exists bool, err error) {

	if _, err := os.Stat(filename); err == nil {
		logger.Debug("File %q  exists", filename)
		return true, nil
	} else {
		return false, err
	}

	return true, nil //Make compiler happy
}

//Retrieve container id saved in a file and then removed the filesince it is not needed anymore
func GetContainerIdFromFile(containerName string) string {

	filename := constants.ID_FILE + containerName

	if _, err := CheckIfFileExists(filename); err != nil {
		logger.Fatalf("Failed to retrieve container id from  file %q due to error:%v", filename, err)
	}

	file, err := os.Open(filename)
	if err != nil {
		logger.Fatalf("Failed to open file %q due to error: %v", filename, err)
	}

	reader := bufio.NewReader(file)

	idBytes, _, err := reader.ReadLine()

	if err != nil {
		logger.Fatalf("Error trying to read file %q : %v", filename, err)
	}

	id := strings.TrimSpace(string(idBytes))

	logger.Debug("Successfully retrieved id  %q for container %q from file %q", id, containerName, filename)

	file.Close()

	defer os.Remove(filename)

	return id
}

//Creates a new state file.If state file already exists it will be overwritten
func CreateNewStateFile() {

	var statefileTemplate = []string{STATE_CONTAINERS_HEADER}

	writeLines(statefileTemplate, constants.STATE_FILE)

	if _, err := CheckIfFileExists(constants.STATE_FILE); err != nil {
		logger.Fatalf("Failed to create file %q due to error:%v", constants.STATE_FILE, err)
	}
}

//Builds a single state container text block.
func buildStateContainer(containerName, id, ip string) []string {

	containerLine := "[statecontainers." + containerName + "]"
	idLine := ID_LINE + " \"" + id + "\""
	ipLine := IP_LINE + " \"" + ip + "\""

	return []string{containerLine, idLine, ipLine}
}

//Remove chosen containers from the state file
func RemoveStateContainers(containersToBeRemoved []string) {

	file, err := os.Open(constants.STATE_FILE)
	if err != nil {
		logger.Fatalf("Failed to open the file %q due to error: %v", constants.STATE_FILE, err)
	}

	defer file.Close()

	reader := bufio.NewReader(file)
	buffer := bytes.NewBuffer(make([]byte, 0))

	//Remove chosen containers from the state file

	var (
		//Line counter for keeping track of lines to be removed
		counter = 0
		lines   []string
	)

OUTER:
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF { //End of file, exit loop
			break
		} else if err != nil {
			logger.Fatalf("Error trying to read file %q : %v", constants.STATE_FILE, err)
		} else if counter > 0 {
			counter--
			continue //skip current line
		}

		if strings.Contains(line, "[statecontainers.") {
			for _, containerName := range containersToBeRemoved {
				if strings.Contains(line, containerName) {
					counter = 2
					continue OUTER //skip current line
				}
			}
		}
		buffer.WriteString(line)

		lines = append(lines, buffer.String())
		buffer.Reset()
	}
	writeLines(lines, constants.STATE_FILE)
}

//Updates state file. If container already exists in the state file it is updated. If it does not exist it is added to it.
//If state file does not exists it is created.
func UpdateStateFile(stateContainers map[string]container.StateContainer) {

	var (
		currentContainerIP   = DEFAULT_CONTAINER_IP
		currentContainerID   = DEFAULT_CONTAINER_ID
		currentContainerName = DEFAULT_CONTAINER_NAME
		lines                []string
		updatedContainers    = map[string]container.StateContainer{}
	)

	//Create new state file if it does not exists already.
	if exists, _ := CheckIfFileExists(constants.STATE_FILE); !exists {
		CreateNewStateFile()
	}

	file, err := os.Open(constants.STATE_FILE)
	if err != nil {
		logger.Fatalf("Failed to open the file %q due to error: %v", constants.STATE_FILE, err)
	}

	defer file.Close()

	reader := bufio.NewReader(file)
	buffer := bytes.NewBuffer(make([]byte, 0))

	//1.Update records if they exists.

	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF { //End of file, exit loop
			break
		}
		if err != nil {
			logger.Fatalf("Error trying to read file %q : %v", constants.STATE_FILE, err)
		}

		if strings.Contains(line, "[statecontainers.") {
			for containerName, stateContainer := range stateContainers {
				if strings.Contains(line, containerName) {
					currentContainerIP = stateContainer.IP
					currentContainerID = stateContainer.ID
					currentContainerName = containerName
					logger.Debug("Current container is %q with ID=%q and IP=%q", containerName, currentContainerID, currentContainerIP)
					updatedContainers[containerName] = stateContainer
					break
				}
			}
		} else if currentContainerID != DEFAULT_CONTAINER_ID && strings.Contains(line, ID_LINE) { //Update ID
			logger.Debug("For container %q: Current ID line is:\n%s", currentContainerName, line)
			line = ID_LINE + " \"" + currentContainerID + "\""
			logger.Debug("For container %q:New ID line is:\n %s", currentContainerName, line)

			//Reset current container
			currentContainerID = DEFAULT_CONTAINER_ID
		} else if currentContainerIP != DEFAULT_CONTAINER_IP && strings.Contains(line, IP_LINE) { //Update IP
			logger.Debug("For container %q: Current IP line is:\n%s", currentContainerName, line)
			line = IP_LINE + " \"" + currentContainerIP + "\""
			logger.Debug("For container %q: New IP line is:\n %s", currentContainerName, line)

			//Reset current container
			currentContainerIP = DEFAULT_CONTAINER_IP
		}

		buffer.WriteString(line)

		lines = append(lines, buffer.String())
		buffer.Reset()
	}

	//2.Add new containers to the state file if they don't exist yet.

	for containerName, stateContainer := range stateContainers {
		if _, ok := updatedContainers[containerName]; !ok {
			logger.Debug("New container %q will be added to the state file", containerName)
			//If key does not exists in the updatedContainers it means that it wasn't updated hence it must be created and added.

			newContainerLines := buildStateContainer(containerName, stateContainer.ID, stateContainer.IP)

			lines = append(lines, newContainerLines...)
		}
	}

	writeLines(lines, constants.STATE_FILE)
}

func writeLines(lines []string, filename string) {

	file, err := os.Create(filename)
	if err != nil {
		logger.Fatalf("Failed to create file %q due to error:%v", filename, err)
	}
	defer file.Close()

	for _, item := range lines {
		_, err := file.WriteString(strings.TrimSpace(item) + "\n")
		if err != nil {
			logger.Fatalf("Failed to write to a file %q due to error:", filename, err)
		}
	}
}
