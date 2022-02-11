package dbplugin

import (
	"fmt"

	"github.com/hashicorp/go-plugin"
	"github.com/hashicorp/vault/sdk/helper/pluginutil"
)

// Serve is called from within a plugin and wraps the provided
// Database implementation in a databasePluginRPCServer object and starts a
// RPC server.
func Serve(db Database) {
	plugin.Serve(ServeConfig(db))
}

func ServeConfig(db Database) *plugin.ServeConfig {
	err := pluginutil.OptionallyEnableMlock()
	if err != nil {
		fmt.Println(err)
		return nil
	}

	// pluginSets is the map of plugins we can dispense.
	pluginSets := map[int]plugin.PluginSet{
		5: {
			"database": &GRPCDatabasePlugin{
				Impl:                db,
				multiplexingSupport: false,
			},
		},
	}

	conf := &plugin.ServeConfig{
		HandshakeConfig:  HandshakeConfig,
		VersionedPlugins: pluginSets,
		GRPCServer:       plugin.DefaultGRPCServer,
	}

	return conf
}

func ServeMultiplex(factory Factory) {
	plugin.Serve(ServeConfigMultiplex(factory))
}

func ServeConfigMultiplex(factory Factory) *plugin.ServeConfig {
	err := pluginutil.OptionallyEnableMlock()
	if err != nil {
		fmt.Println(err)
		return nil
	}

	// pluginSets is the map of plugins we can dispense.
	pluginSets := map[int]plugin.PluginSet{
		6: {
			"database": &GRPCDatabasePlugin{
				FactoryFunc:         factory,
				multiplexingSupport: true,
			},
		},
	}

	conf := &plugin.ServeConfig{
		HandshakeConfig:  HandshakeConfig,
		VersionedPlugins: pluginSets,
		GRPCServer:       plugin.DefaultGRPCServer,
	}

	return conf
}
