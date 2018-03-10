package command

import (
	"fmt"
	"strings"

	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

var _ cli.Command = (*KVMetadataDeleteCommand)(nil)
var _ cli.CommandAutocomplete = (*KVMetadataDeleteCommand)(nil)

type KVMetadataDeleteCommand struct {
	*BaseCommand
}

func (c *KVMetadataDeleteCommand) Synopsis() string {
	return "Delete secrets and configuration"
}

func (c *KVMetadataDeleteCommand) Help() string {
	helpText := `
Usage: vault delete [options] PATH

  Deletes secrets and configuration from Vault at the given path. The behavior
  of "delete" is delegated to the backend corresponding to the given path.

  Remove data in the status secret backend:

      $ vault delete secret/my-secret

  Uninstall an encryption key in the transit backend:

      $ vault delete transit/keys/my-key

  Delete an IAM role:

      $ vault delete aws/roles/ops

  For a full list of examples and paths, please see the documentation that
  corresponds to the secret backend in use.

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *KVMetadataDeleteCommand) Flags() *FlagSets {
	return c.flagSet(FlagSetHTTP)
}

func (c *KVMetadataDeleteCommand) AutocompleteArgs() complete.Predictor {
	return c.PredictVaultFiles()
}

func (c *KVMetadataDeleteCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *KVMetadataDeleteCommand) Run(args []string) int {
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
	path, err = addPrefixToVKVPath(path, "metadata")
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	if _, err := client.Logical().Delete(path); err != nil {
		c.UI.Error(fmt.Sprintf("Error deleting %s: %s", path, err))
		return 2
	}

	c.UI.Info(fmt.Sprintf("Success! Data deleted (if it existed) at: %s", path))
	return 0
}
