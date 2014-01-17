package constants

/*
  Crane and Docker Commands
*/

const (
	CREATE    = "create"
	PULL      = "pull"
	PULLALL   = "pullall"
	START     = "start"
	RUN       = "run"
	RUNALL    = "runall"
	ENTER     = "enter"
	DESTROY   = "destroy"
	BUILDALL  = "buildall"
	REMOVEALL = "rmiall"
)

/*
   Docker Commands
*/

const (
	DOCKER       = "docker"
	SUDO         = "sudo"
	INSPECT      = "inspect"
	KILL         = "kill"
	REMOVE       = "rm"
	HISTORY      = "history"
	SEARCH       = "search"
	BUILD        = "build"
	IMAGES       = "images"
	REMOVE_IMAGE = "rmi"
	COMMIT       = "commit"
	//"run","pull","create" are the same as for crane
)

/*
  Common Commands and Constants
*/
const (
	COMMANDS_DELIMITER     = ":"
	INTERNAL_DELIMITER     = ";"
	OWN_COMMANDS_DELIMITER = "#"
	FREEZE_DELIMITER       = "::"

	SSHD_COMMAND        = "/usr/sbin/sshd -D"
	SHELL_COMMAND       = "/bin/bash"
	SHELL_STRING_OPTION = "-c"

	GREP = "grep"

	LOGGER_NAME = "crane"

	CONFIGURATION_FILE = "Cranefile.toml"
	STATE_FILE         = ".crane"
	ID_FILE            = ".cidfile"

	NOT_DAEMONIZED_IP = "not_deamonized_has_no_ip"
)
