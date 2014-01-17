package config

import (
	"github.com/BurntSushi/toml"
	"github.com/SnowRipple/crane/constants"
	"github.com/SnowRipple/crane/container"
	"github.com/SnowRipple/crane/io"
	log "github.com/SnowRipple/crane/logger"
)

type TomlConfig struct {
	CraneConfig
	CraneState
}

type CraneConfig struct {
	Containers map[string]container.Container
}

type CraneState struct {
	StateContainers map[string]container.StateContainer
}

var logger = log.GetLogger()

//Reads config  and state files.
func ReadConfig() TomlConfig {

	var (
		config CraneConfig
		state  CraneState
	)

	//Create example state and config files if they do not exist already.
	if exists, _ := io.CheckIfFileExists(constants.STATE_FILE); !exists {
		logger.Debug("Creating new state file %q", constants.STATE_FILE)
		io.CreateNewStateFile()
	}

	if exists, _ := io.CheckIfFileExists(constants.CONFIGURATION_FILE); !exists {

		logger.Debug("Creating new configuration file %q", constants.CONFIGURATION_FILE)
		io.CreateNewConfigFile()
	}

	//Decode configuration file
	_, err := toml.DecodeFile(constants.CONFIGURATION_FILE, &config)
	if err != nil {
		logger.Fatalf("Failed to decode %q file due to error:", constants.CONFIGURATION_FILE, err)
	}

	//Decode state file
	_, err = toml.DecodeFile(constants.STATE_FILE, &state)
	if err != nil {
		logger.Fatalf("Failed to decode %q file due to error:", constants.STATE_FILE, err)
	}

	logger.Debug("Decoded config file:\n%v", config)
	logger.Debug("Decoded state file:\n%v", state)

	return TomlConfig{
		CraneConfig: config,
		CraneState:  state}
}

/*
func decodeTomlFile(filename string, config interface{}) interface{}{

	_, err := toml.DecodeFile(filename, &config)
	if err != nil {
    logger.Fatalf("Failed to decode %q file due to error:",filename, err)
	}

  return config
}

*/
