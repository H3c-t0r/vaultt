package command

import (
	"fmt"
	"sort"
	"strings"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

var (
	_ cli.Command             = (*PluginListCommand)(nil)
	_ cli.CommandAutocomplete = (*PluginListCommand)(nil)
)

type PluginListCommand struct {
	*BaseCommand

	flagDetailed bool
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

  List all available plugins with detailed output:

      $ vault plugin list -detailed

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *PluginListCommand) Flags() *FlagSets {
	set := c.flagSet(FlagSetHTTP | FlagSetOutputFormat)

	f := set.NewFlagSet("Command Options")

	f.BoolVar(&BoolVar{
		Name:    "detailed",
		Target:  &c.flagDetailed,
		Default: false,
		Usage: "Print detailed plugin information such as plugin type, " +
			"version, and deprecation status for each plugin. This option " +
			"is only applicable to table-formatted output.",
	})

	return set
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
	if resp == nil {
		c.UI.Error("No response from server when listing plugins")
		return 2
	}

	switch Format(c.UI) {
	case "table":
		if c.flagDetailed {
			c.UI.Output(tableOutput(c.detailedResponse(resp), nil))
			return 0
		}
		c.UI.Output(tableOutput(c.simpleResponse(resp), nil))
		return 0
	default:
		res := make(map[string]interface{})
		for k, v := range resp.PluginsByType {
			res[k.String()] = v
		}
		return OutputData(c.UI, res)
	}
}

func (c *PluginListCommand) simpleResponse(plugins *api.ListPluginsResponse) []string {
	var flattenedNames []string
	namesAdded := make(map[string]bool)
	for _, names := range plugins.PluginsByType {
		for _, name := range names {
			if ok := namesAdded[name]; !ok {
				flattenedNames = append(flattenedNames, name)
				namesAdded[name] = true
			}
		}
		sort.Strings(flattenedNames)
	}
	list := append([]string{"Plugins"}, flattenedNames...)
	return list
}

func (c *PluginListCommand) detailedResponse(plugins *api.ListPluginsResponse) []string {
	out := []string{"Name | Type | Version | Deprecation Status"}
	for _, plugin := range plugins.Details {
		out = append(out, fmt.Sprintf("%s | %s | %s | %s", plugin.Name, plugin.Type, plugin.Version, plugin.DeprecationStatus))
	}

	return out
}
