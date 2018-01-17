package command

import (
	"fmt"
	"strings"

	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

var _ cli.Command = (*ReadCommand)(nil)
var _ cli.CommandAutocomplete = (*ReadCommand)(nil)

type ReadCommand struct {
	*BaseCommand
}

func (c *ReadCommand) Synopsis() string {
	return "Read data and retrieves secrets"
}

func (c *ReadCommand) Help() string {
	helpText := `
Usage: vault read [options] PATH

  Reads data from Vault at the given path. This can be used to read secrets,
  generate dynamic credentials, get configuration details, and more.

  Read a secret from the static secrets engine:

      $ vault read secret/my-secret

  For a full list of examples and paths, please see the documentation that
  corresponds to the secrets engine in use.

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *ReadCommand) Flags() *FlagSets {
	return c.flagSet(FlagSetHTTP | FlagSetOutputField | FlagSetOutputFormat)
}

func (c *ReadCommand) AutocompleteArgs() complete.Predictor {
	return c.PredictVaultFiles()
}

func (c *ReadCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *ReadCommand) Run(args []string) int {
	f := c.Flags()

	if err := f.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	args = f.Args()
	switch {
	case len(args) < 1:
		c.UI.Error(fmt.Sprintf("Not enough arguments (expected 1, got %d)", len(args)))
		return 1
	case len(args) > 1:
		c.UI.Error(fmt.Sprintf("Too many arguments (expected 1, got %d)", len(args)))
		return 1
	}

	client, err := c.Client()
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	path := sanitizePath(args[0])

	secret, err := client.Logical().Read(path)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error reading %s: %s", path, err))
		return 2
	}
	if secret == nil {
		c.UI.Error(fmt.Sprintf("No value found at %s", path))
		return 2
	}

	if c.flagField != "" {
		return PrintRawField(c.UI, secret, c.flagField)
	}

	return OutputSecret(c.UI, c.flagFormat, secret)
}
