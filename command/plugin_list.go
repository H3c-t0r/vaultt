package command

import (
	"fmt"
	"strings"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/helper/consts"
	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

var _ cli.Command = (*PluginListCommand)(nil)
var _ cli.CommandAutocomplete = (*PluginListCommand)(nil)

type PluginListCommand struct {
	*BaseCommand
}

func (c *PluginListCommand) Synopsis() string {
	return "Lists available plugins"
}

func (c *PluginListCommand) Help() string {
	helpText := `
Usage: vault plugin list [options] [TYPE]

  Lists available plugins registered in the catalog. This does not list whether
  plugins are in use, but rather just their availability. The last argument of
  type takes "auth", "database", or "secret".

  List all available plugins in the catalog:

      $ vault plugin list

  List all available database plugins in the catalog:

      $ vault plugin list database

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *PluginListCommand) Flags() *FlagSets {
	return c.flagSet(FlagSetHTTP | FlagSetOutputFormat)
}

func (c *PluginListCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *PluginListCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *PluginListCommand) Run(args []string) int {
	f := c.Flags()

	if err := f.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	args = f.Args()
	switch {
	case len(args) > 1:
		c.UI.Error(fmt.Sprintf("Too many arguments (expected 0 or 1, got %d)", len(args)))
		return 1
	}

	pluginType := consts.PluginTypeUnknown
	if len(args) > 0 {
		pluginTypeStr := strings.TrimSpace(args[0])
		if pluginTypeStr != "" {
			var err error
			pluginType, err = consts.ParsePluginType(pluginTypeStr)
			if err != nil {
				c.UI.Error(fmt.Sprintf("Error parsing type: %s", err))
				return 2
			}
		}
	}

	client, err := c.Client()
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	resp, err := client.Sys().ListPlugins(&api.ListPluginsInput{
		Type: pluginType,
	})
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error listing available plugins: %s", err))
		return 2
	}

	return OutputData(c.UI, formatListPluginsResponse(resp))
}

func formatListPluginsResponse(resp *api.ListPluginsResponse) map[string]interface{} {
	res := make(map[string]interface{})
	for k, v := range resp.PluginsByType {
		res[k.String()] = v
	}
	return res
}
