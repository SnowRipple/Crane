package command

import (
	"fmt"
	"github.com/mitchellh/cli"
	"strings"
)

// VersionCommand presents the user with the Crane version.
type VersionCommand struct {
	Revision string
	Version  string
	Status   string
	Ui       cli.Ui
}

func (c *VersionCommand) Help() string {
	helpText := `
Usage: crane version

Usage: crane -v

Usage: crane --version

Presents the user with the Crane version information.`

	return strings.TrimSpace(helpText)
}

//Presents the user with the Crane version information
func (c *VersionCommand) Run(imageNames []string) int {

	logger.Debug("Entered version command...")

	version := fmt.Sprintf("Crane v%s.%s", c.Version, c.Status)

	if c.Revision != "" {
		version = version + " (" + c.Revision + ")"
	}

	c.Ui.Output(version)

	return 0
}

func (c *VersionCommand) Synopsis() string {
	return "Shows Crane version."
}
