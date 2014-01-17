package executer

import (
	log "github.com/SnowRipple/crane/logger"
	"github.com/SnowRipple/crane/utils"
	"io"
	"os"
	"os/exec"
)

const SUDO = "sudo"

var logger = log.GetLogger()

func GetCommandOutput(command []string) ([]byte, error) {

	logger.Debug("\nFinal docker command: %v\n", command)

	//CombinedOutput is needed so errors can be returned as well
	return exec.Command(SUDO, command...).CombinedOutput()
}

func ExecuteCommand(command []string) {

	logger.Debug("\nFinal docker command: %v\n", command)

	cmd := exec.Command(SUDO, command...)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		logger.Fatal("Failed to set up command's stdin", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		logger.Fatal("Failed to set up command's stderr", err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		logger.Fatal("Failed to set up command's stdout", err)
	}

	go io.Copy(os.Stdout, stdout)
	go io.Copy(stdin, os.Stdin)
	go io.Copy(os.Stderr, stderr)

	if err = cmd.Start(); err != nil {
		logger.Fatal("Failed to start the command", err)
	}
	defer cmd.Wait()
}

//Pipe multiple commands in a unix fashion
func PipeCommands(commands ...*exec.Cmd) string {

	//Connect command's stdout with the next command's stdin
	for index, command := range commands[:len(commands)-1] {
		stdout, err := command.StdoutPipe()
		if err != nil {
			logger.Fatal("Piping commands:Failed to set up command's stdout", err)
		}
		command.Start()
		commands[index+1].Stdin = stdout
	}

	output, err := commands[len(commands)-1].CombinedOutput()
	if err != nil {
		logger.Fatal("Piping commands:Failed to get output of the final command", utils.ExtractContainerMessage(output, err))
	}
	return string(output)
}
