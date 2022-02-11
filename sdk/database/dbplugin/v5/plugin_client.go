package dbplugin

import (
	"context"
	"errors"

	plugin "github.com/hashicorp/go-plugin"
	"github.com/hashicorp/vault/sdk/database/dbplugin/v5/proto"
	"github.com/hashicorp/vault/sdk/helper/pluginutil"
)

type DatabasePluginClient struct {
	client pluginutil.PluginClient
	Database
}

// This wraps the Close call and ensures we both close the database connection
// and kill the plugin.
func (dc *DatabasePluginClient) Close() error {
	err := dc.Database.Close()
	dc.client.Close()

	return err
}

// pluginSets is the map of plugins we can dispense.
var PluginSets = map[int]plugin.PluginSet{
	5: {
		"database": &GRPCDatabasePlugin{multiplexingSupport: false},
	},
	6: {
		"database": &GRPCDatabasePlugin{multiplexingSupport: true},
	},
}

// NewPluginClient returns a databaseRPCClient with a connection to a running
// plugin.
func NewPluginClient(ctx context.Context, sys pluginutil.RunnerUtil, config pluginutil.PluginClientConfig) (Database, error) {
	pluginClient, err := sys.NewPluginClient(ctx, config)
	if err != nil {
		return nil, err
	}

	// Request the plugin
	raw, err := pluginClient.Dispense("database")
	if err != nil {
		return nil, err
	}

	// We should have a database type now. This feels like a normal interface
	// implementation but is in fact over an RPC connection.
	var db Database
	switch c := raw.(type) {
	case gRPCClient:
		// This is an abstraction leak from go-plugin but it is necessary in
		// order to enable multiplexing on multiplexed plugins
		c.client = proto.NewDatabaseClient(pluginClient.Conn())

		db = c
	default:
		return nil, errors.New("unsupported client type")
	}

	return &DatabasePluginClient{
		client:   pluginClient,
		Database: db,
	}, nil
}
