package command

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/vault/api"
	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

var _ cli.Command = (*SecretsEnableCommand)(nil)
var _ cli.CommandAutocomplete = (*SecretsEnableCommand)(nil)

type SecretsEnableCommand struct {
	*BaseCommand

	flagDescription     string
	flagPath            string
	flagDefaultLeaseTTL time.Duration
	flagMaxLeaseTTL     time.Duration
	flagForceNoCache    bool
	flagPluginName      string
	flagLocal           bool
	flagSealWrap        bool
}

func (c *SecretsEnableCommand) Synopsis() string {
	return "Enable a secrets engine"
}

func (c *SecretsEnableCommand) Help() string {
	helpText := `
Usage: vault secrets enable [options] TYPE

  Enables a secrets engine. By default, secrets engines are enabled at the path
  corresponding to their TYPE, but users can customize the path using the
  -path option.

  Once enabled, Vault will route all requests which begin with the path to the
  secrets engine.

  Enable the AWS secrets engine at aws/:

      $ vault secrets enable aws

  Enable the SSH secrets engine at ssh-prod/:

      $ vault secrets enable -path=ssh-prod ssh

  Enable the database secrets engine with an explicit maximum TTL of 30m:

      $ vault secrets enable -max-lease-ttl=30m database

  Enable a custom plugin (after it is registered in the plugin registry):

      $ vault secrets enable -path=my-secrets -plugin-name=my-plugin plugin

  For a full list of secrets engines and examples, please see the documentation.

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *SecretsEnableCommand) Flags() *FlagSets {
	set := c.flagSet(FlagSetHTTP)

	f := set.NewFlagSet("Command Options")

	f.StringVar(&StringVar{
		Name:       "description",
		Target:     &c.flagDescription,
		Completion: complete.PredictAnything,
		Usage:      "Human-friendly description for the purpose of this engine.",
	})

	f.StringVar(&StringVar{
		Name:       "path",
		Target:     &c.flagPath,
		Default:    "", // The default is complex, so we have to manually document
		Completion: complete.PredictAnything,
		Usage: "Place where the secrets engine will be accessible. This must be " +
			"unique cross all secrets engines. This defaults to the \"type\" of the " +
			"secrets engine.",
	})

	f.DurationVar(&DurationVar{
		Name:       "default-lease-ttl",
		Target:     &c.flagDefaultLeaseTTL,
		Completion: complete.PredictAnything,
		Usage: "The default lease TTL for this secrets engine. If unspecified, " +
			"this defaults to the Vault server's globally configured default lease " +
			"TTL.",
	})

	f.DurationVar(&DurationVar{
		Name:       "max-lease-ttl",
		Target:     &c.flagMaxLeaseTTL,
		Completion: complete.PredictAnything,
		Usage: "The maximum lease TTL for this secrets engine. If unspecified, " +
			"this defaults to the Vault server's globally configured maximum lease " +
			"TTL.",
	})

	f.BoolVar(&BoolVar{
		Name:    "force-no-cache",
		Target:  &c.flagForceNoCache,
		Default: false,
		Usage: "Force the secrets engine to disable caching. If unspecified, this " +
			"defaults to the Vault server's globally configured cache settings. " +
			"This does not affect caching of the underlying encrypted data storage.",
	})

	f.StringVar(&StringVar{
		Name:       "plugin-name",
		Target:     &c.flagPluginName,
		Completion: complete.PredictAnything,
		Usage: "Name of the secrets engine plugin. This plugin name must already " +
			"exist in Vault's plugin catalog.",
	})

	f.BoolVar(&BoolVar{
		Name:    "local",
		Target:  &c.flagLocal,
		Default: false,
		Usage: "Mark the secrets engine as local-only. Local engines are not " +
			"replicated or removed by replication.",
	})

	f.BoolVar(&BoolVar{
		Name:    "seal-wrap",
		Target:  &c.flagSealWrap,
		Default: false,
		Usage:   "Enable seal wrapping of critical values in the secrets engine.",
	})

	return set
}

func (c *SecretsEnableCommand) AutocompleteArgs() complete.Predictor {
	return c.PredictVaultAvailableMounts()
}

func (c *SecretsEnableCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *SecretsEnableCommand) Run(args []string) int {
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

	// Get the engine type type (first arg)
	engineType := strings.TrimSpace(args[0])

	// If no path is specified, we default the path to the backend type
	// or use the plugin name if it's a plugin backend
	mountPath := c.flagPath
	if mountPath == "" {
		if engineType == "plugin" {
			mountPath = c.flagPluginName
		} else {
			mountPath = engineType
		}
	}

	// Append a trailing slash to indicate it's a path in output
	mountPath = ensureTrailingSlash(mountPath)

	// Build mount input
	mountInput := &api.MountInput{
		Type:        engineType,
		Description: c.flagDescription,
		Local:       c.flagLocal,
		SealWrap:    c.flagSealWrap,
		Config: api.MountConfigInput{
			DefaultLeaseTTL: c.flagDefaultLeaseTTL.String(),
			MaxLeaseTTL:     c.flagMaxLeaseTTL.String(),
			ForceNoCache:    c.flagForceNoCache,
			PluginName:      c.flagPluginName,
		},
	}

	if err := client.Sys().Mount(mountPath, mountInput); err != nil {
		c.UI.Error(fmt.Sprintf("Error enabling: %s", err))
		return 2
	}

	thing := engineType + " secrets engine"
	if engineType == "plugin" {
		thing = c.flagPluginName + " plugin"
	}

	c.UI.Output(fmt.Sprintf("Success! Enabled the %s at: %s", thing, mountPath))
	return 0
}
