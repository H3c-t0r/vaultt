package dbplugin

import (
	"context"
	"errors"
	"sync"

	log "github.com/hashicorp/go-hclog"
	plugin "github.com/hashicorp/go-plugin"
	"github.com/hashicorp/vault/sdk/helper/pluginutil"
)

// DatabasePluginClient embeds a databasePluginRPCClient and wraps it's Close
// method to also call Kill() on the plugin.Client.
type DatabasePluginClient struct {
	client *plugin.Client
	sync.Mutex

	Database
}

// This wraps the Close call and ensures we both close the database connection
// and kill the plugin.
func (dc *DatabasePluginClient) Close() error {
	err := dc.Database.Close()
	dc.client.Kill()

	return err
}

// NewPluginClient returns a databaseRPCClient with a connection to a running
// plugin. The client is wrapped in a DatabasePluginClient object to ensure the
// plugin is killed on call of Close().
func NewPluginClient(ctx context.Context, sys pluginutil.RunnerUtil, pluginRunner *pluginutil.PluginRunner, logger log.Logger, isMetadataMode bool) (Database, error) {
	rpcClient, err := sys.NewPluginClient(ctx, pluginRunner, logger, isMetadataMode)
	if err != nil {
		return nil, err
	}

	// Request the plugin
	raw, err := rpcClient.Dispense("database")
	if err != nil {
		return nil, err
	}

	// We should have a database type now. This feels like a normal interface
	// implementation but is in fact over an RPC connection.
	var db Database
	switch raw.(type) {
	case gRPCClient:
		db = raw.(gRPCClient)
	default:
		return nil, errors.New("unsupported client type")
	}

	// Wrap RPC implementation in DatabasePluginClient
	return &DatabasePluginClient{
		Database: db,
	}, nil
}
