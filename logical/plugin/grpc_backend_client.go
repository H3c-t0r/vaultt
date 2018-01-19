package plugin

import (
	"context"
	"errors"

	"google.golang.org/grpc"

	"github.com/hashicorp/go-plugin"
	"github.com/hashicorp/vault/helper/pluginutil"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/plugin/pb"
	log "github.com/mgutz/logxi/v1"
)

var ErrPluginShutdown = errors.New("plugin is shut down")

// Validate backendGRPCPluginClient satisfies the logical.Backend interface
var _ logical.Backend = &backendGRPCPluginClient{}

// backendPluginClient implements logical.Backend and is the
// go-plugin client.
type backendGRPCPluginClient struct {
	broker       *plugin.GRPCBroker
	client       pb.BackendClient
	metadataMode bool

	system logical.SystemView
	logger log.Logger

	// server is the grpc server used for serving storage and sysview requests.
	server *grpc.Server
	// clientConn is the underlying grpc connection to the server, we store it
	// so it can be cleaned up.
	clientConn *grpc.ClientConn
	doneCtx    context.Context
}

func (b *backendGRPCPluginClient) HandleRequest(ctx context.Context, req *logical.Request) (*logical.Response, error) {
	if b.metadataMode {
		return nil, ErrClientInMetadataMode
	}

	ctx, cancel := context.WithCancel(ctx)
	quitCh := pluginutil.CtxCancelIfCanceled(cancel, b.doneCtx)
	defer close(quitCh)
	defer cancel()

	protoReq, err := pb.LogicalRequestToProtoRequest(req)
	if err != nil {
		return nil, err
	}

	reply, err := b.client.HandleRequest(ctx, &pb.HandleRequestArgs{
		Request: protoReq,
	})
	if err != nil {
		if b.doneCtx.Err() != nil {
			return nil, ErrPluginShutdown
		}

		return nil, err
	}
	resp, err := pb.ProtoResponseToLogicalResponse(reply.Response)
	if err != nil {
		return nil, err
	}
	if reply.Err != nil {
		return resp, pb.ProtoErrToErr(reply.Err)
	}

	return resp, nil
}

func (b *backendGRPCPluginClient) SpecialPaths() *logical.Paths {
	// Timeout the connection
	reply, err := b.client.SpecialPaths(b.doneCtx, &pb.Empty{})
	if err != nil {
		return nil
	}

	return &logical.Paths{
		Root:            reply.Paths.Root,
		Unauthenticated: reply.Paths.Unauthenticated,
		LocalStorage:    reply.Paths.LocalStorage,
		SealWrapStorage: reply.Paths.SealWrapStorage,
	}
}

// System returns vault's system view. The backend client stores the view during
// Setup, so there is no need to shim the system just to get it back.
func (b *backendGRPCPluginClient) System() logical.SystemView {
	return b.system
}

// Logger returns vault's logger. The backend client stores the logger during
// Setup, so there is no need to shim the logger just to get it back.
func (b *backendGRPCPluginClient) Logger() log.Logger {
	return b.logger
}

func (b *backendGRPCPluginClient) HandleExistenceCheck(ctx context.Context, req *logical.Request) (bool, bool, error) {
	if b.metadataMode {
		return false, false, ErrClientInMetadataMode
	}

	protoReq, err := pb.LogicalRequestToProtoRequest(req)
	if err != nil {
		return false, false, err
	}

	ctx, cancel := context.WithCancel(ctx)
	quitCh := pluginutil.CtxCancelIfCanceled(cancel, b.doneCtx)
	defer close(quitCh)
	defer cancel()
	reply, err := b.client.HandleExistenceCheck(ctx, &pb.HandleExistenceCheckArgs{
		Request: protoReq,
	})
	if err != nil {
		if b.doneCtx.Err() != nil {
			return false, false, ErrPluginShutdown
		}
		return false, false, err
	}
	if reply.Err != nil {
		return false, false, pb.ProtoErrToErr(reply.Err)
	}

	return reply.CheckFound, reply.Exists, nil
}

func (b *backendGRPCPluginClient) Cleanup(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	quitCh := pluginutil.CtxCancelIfCanceled(cancel, b.doneCtx)
	defer close(quitCh)
	defer cancel()

	b.client.Cleanup(ctx, &pb.Empty{})
	if b.server != nil {
		b.server.GracefulStop()
	}
	b.clientConn.Close()
}

func (b *backendGRPCPluginClient) Initialize(ctx context.Context) error {
	if b.metadataMode {
		return ErrClientInMetadataMode
	}

	ctx, cancel := context.WithCancel(ctx)
	quitCh := pluginutil.CtxCancelIfCanceled(cancel, b.doneCtx)
	defer close(quitCh)
	defer cancel()

	_, err := b.client.Initialize(ctx, &pb.Empty{})
	return err
}

func (b *backendGRPCPluginClient) InvalidateKey(ctx context.Context, key string) {
	if b.metadataMode {
		return
	}

	ctx, cancel := context.WithCancel(ctx)
	quitCh := pluginutil.CtxCancelIfCanceled(cancel, b.doneCtx)
	defer close(quitCh)
	defer cancel()

	b.client.InvalidateKey(ctx, &pb.InvalidateKeyArgs{
		Key: key,
	})
}

func (b *backendGRPCPluginClient) Setup(ctx context.Context, config *logical.BackendConfig) error {
	// Shim logical.Storage
	storageImpl := config.StorageView
	if b.metadataMode {
		storageImpl = &NOOPStorage{}
	}
	storage := &GRPCStorageServer{
		impl: storageImpl,
	}

	// Shim logical.SystemView
	sysViewImpl := config.System
	if b.metadataMode {
		sysViewImpl = &logical.StaticSystemView{}
	}
	sysView := &gRPCSystemViewServer{
		impl: sysViewImpl,
	}

	// Register the server in this closure.
	serverFunc := func(opts []grpc.ServerOption) *grpc.Server {
		s := grpc.NewServer(opts...)
		pb.RegisterSystemViewServer(s, sysView)
		pb.RegisterStorageServer(s, storage)
		b.server = s
		return s
	}
	brokerID := b.broker.NextId()
	go b.broker.AcceptAndServe(brokerID, serverFunc)

	args := &pb.SetupArgs{
		BrokerID: brokerID,
		Config:   config.Config,
	}

	ctx, cancel := context.WithCancel(ctx)
	quitCh := pluginutil.CtxCancelIfCanceled(cancel, b.doneCtx)
	defer close(quitCh)
	defer cancel()

	reply, err := b.client.Setup(ctx, args)
	if err != nil {
		return err
	}
	if reply.Err != "" {
		return errors.New(reply.Err)
	}

	// Set system and logger for getter methods
	b.system = config.System
	b.logger = config.Logger

	return nil
}

func (b *backendGRPCPluginClient) Type() logical.BackendType {
	reply, err := b.client.Type(b.doneCtx, &pb.Empty{})
	if err != nil {
		return logical.TypeUnknown
	}

	return logical.BackendType(reply.Type)
}
