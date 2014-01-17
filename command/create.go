package command

import (
	"flag"
	"github.com/SnowRipple/crane/io"
	"github.com/mitchellh/cli"
	"strings"
)

// CreateCommand creates an example configuration file Cranefile.toml.
type CreateCommand struct {
	Ui cli.Ui
}

func (c *CreateCommand) Help() string {
	helpText := `
  Usage: crane create
  
  Creates an example configuration file.`

	return strings.TrimSpace(helpText)
}
func (c *CreateCommand) Run(_ []string) int {

	cmdFlags := flag.NewFlagSet("create", flag.ContinueOnError)
	cmdFlags.Usage = func() { c.Ui.Output(c.Help()) }

	io.Create()
	return 0
}

func (c *CreateCommand) Synopsis() string {
	return "Create an example configuration file."
}
