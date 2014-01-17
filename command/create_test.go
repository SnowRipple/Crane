package command

import (
	"github.com/mitchellh/cli"
  "testing"
)

func TestCreateCommand_implements(t *testing.T) {
	var _ cli.Command = &CreateCommand{}
}

