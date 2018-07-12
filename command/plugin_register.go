package command

import (
	"fmt"
	"strings"

	"github.com/hashicorp/vault/api"
	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

var _ cli.Command = (*PluginRegisterCommand)(nil)
var _ cli.CommandAutocomplete = (*PluginRegisterCommand)(nil)

type PluginRegisterCommand struct {
	*BaseCommand

	flagArgs    []string
	flagCommand string
	flagSHA256  string
}

func (c *PluginRegisterCommand) Synopsis() string {
	return "Registers a new plugin in the catalog"
}

func (c *PluginRegisterCommand) Help() string {
	helpText := `
Usage: vault plugin register [options] NAME

  Registers a new plugin in the catalog. The plugin binary must exist in Vault's
  configured plugin directory.

  Register the plugin named my-custom-plugin:

      $ vault plugin register -sha256=d3f0a8b... my-custom-plugin

  Register a plugin with custom arguments:

      $ vault plugin register \
          -sha256=d3f0a8b... \
          -args=--with-glibc,--with-cgo \
          my-custom-plugin

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *PluginRegisterCommand) Flags() *FlagSets {
	set := c.flagSet(FlagSetHTTP)

	f := set.NewFlagSet("Command Options")

	f.StringSliceVar(&StringSliceVar{
		Name:       "args",
		Target:     &c.flagArgs,
		Completion: complete.PredictAnything,
		Usage: "Arguments to pass to the plugin when starting. Separate " +
			"multiple arguments with a comma.",
	})

	f.StringVar(&StringVar{
		Name:       "command",
		Target:     &c.flagCommand,
		Completion: complete.PredictAnything,
		Usage: "Command to spawn the plugin. This defaults to the name of the " +
			"plugin if unspecified.",
	})

	f.StringVar(&StringVar{
		Name:       "sha256",
		Target:     &c.flagSHA256,
		Completion: complete.PredictAnything,
		Usage:      "SHA256 of the plugin binary. This is required for all plugins.",
	})

	return set
}

func (c *PluginRegisterCommand) AutocompleteArgs() complete.Predictor {
	return c.PredictVaultPlugins()
}

func (c *PluginRegisterCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *PluginRegisterCommand) Run(args []string) int {
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
	case c.flagSHA256 == "":
		c.UI.Error("SHA256 is required for all plugins, please provide -sha256")
		return 1
	}

	client, err := c.Client()
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	pluginName := strings.TrimSpace(args[0])

	command := c.flagCommand
	if command == "" {
		command = pluginName
	}

	if err := client.Sys().RegisterPlugin(&api.RegisterPluginInput{
		Name:    pluginName,
		Args:    c.flagArgs,
		Command: command,
		SHA256:  c.flagSHA256,
	}); err != nil {
		c.UI.Error(fmt.Sprintf("Error registering plugin %s: %s", pluginName, err))
		return 2
	}

	c.UI.Output(fmt.Sprintf("Success! Registered plugin: %s", pluginName))
	return 0
}
