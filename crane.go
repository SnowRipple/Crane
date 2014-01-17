package main

import (
	"fmt"
	"github.com/SnowRipple/crane/constants"
	ownLog "github.com/SnowRipple/crane/logger"
	flags "github.com/jessevdk/go-flags"
	"github.com/mitchellh/cli"
	log "github.com/op/go-logging"
	"os"
	"path/filepath"
	"strings"
)

type options struct {
	DebugMode bool `short:"d" long:"debug" description:"When crane is used in the debug mode a lot of extra information is provided during program execution." `

	Version bool `short:"v" long:"version" description:"Shows the information about the crane version you are using."`
}

const LOGGER_NAME = "crane"

var logger = ownLog.GetLogger()

func main() {

	var options constants.CommonFlags

	craneArguments := os.Args[1:]
	_, err := flags.ParseArgs(&options, craneArguments)

	if options.DebugMode {
		log.SetLevel(log.DEBUG, LOGGER_NAME)
	}

	logger.Debug("Command line arguments provided: %v", craneArguments)

	if options.Version {
		craneArguments = []string{"version"}
	}

	cli := &cli.CLI{
		Args:     craneArguments,
		Commands: Commands,
	}

	clean()

	_, err = cli.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing CLI: %s\n", err.Error())
	}
}

//Remove cidfile that could remain after previous runs
func clean() {

	err := filepath.Walk(".", traverse) //current dir
	if err != nil {
		logger.Fatalf("Failed to clean the directory from the last run leftovers due to error: %s", err.Error())
	}
}

func traverse(path string, f os.FileInfo, err error) error {
	if strings.Contains(path, constants.ID_FILE) {
		err := os.Remove(path)
		if err != nil {
			return err
		}
		logger.Debug("Found orphaned file %q and removed it.", path)
	}
	return nil
}
