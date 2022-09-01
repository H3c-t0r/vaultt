package plugin

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/helper/license"
	"github.com/hashicorp/vault/sdk/helper/pluginutil"
	"github.com/hashicorp/vault/sdk/helper/wrapping"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/sdk/plugin/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func newGRPCSystemView(conn *grpc.ClientConn) *gRPCSystemViewClient {
	return &gRPCSystemViewClient{
		client: pb.NewSystemViewClient(conn),
	}
}

type gRPCSystemViewClient struct {
	client pb.SystemViewClient
}

func (s *gRPCSystemViewClient) DefaultLeaseTTL() time.Duration {
	reply, err := s.client.DefaultLeaseTTL(context.Background(), &pb.Empty{})
	if err != nil {
		return 0
	}

	return time.Duration(reply.TTL)
}

func (s *gRPCSystemViewClient) MaxLeaseTTL() time.Duration {
	reply, err := s.client.MaxLeaseTTL(context.Background(), &pb.Empty{})
	if err != nil {
		return 0
	}

	return time.Duration(reply.TTL)
}

func (s *gRPCSystemViewClient) Tainted() bool {
	reply, err := s.client.Tainted(context.Background(), &pb.Empty{})
	if err != nil {
		return false
	}

	return reply.Tainted
}

func (s *gRPCSystemViewClient) CachingDisabled() bool {
	reply, err := s.client.CachingDisabled(context.Background(), &pb.Empty{})
	if err != nil {
		return false
	}

	return reply.Disabled
}

func (s *gRPCSystemViewClient) ReplicationState() consts.ReplicationState {
	reply, err := s.client.ReplicationState(context.Background(), &pb.Empty{})
	if err != nil {
		return consts.ReplicationUnknown
	}

	return consts.ReplicationState(reply.State)
}

func (s *gRPCSystemViewClient) ResponseWrapData(ctx context.Context, data map[string]interface{}, ttl time.Duration, jwt bool) (*wrapping.ResponseWrapInfo, error) {
	buf, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	reply, err := s.client.ResponseWrapData(ctx, &pb.ResponseWrapDataArgs{
		Data: string(buf[:]),
		TTL:  int64(ttl),
		JWT:  false,
	})
	if err != nil {
		return nil, err
	}
	if reply.Err != "" {
		return nil, errors.New(reply.Err)
	}

	info, err := pb.ProtoResponseWrapInfoToLogicalResponseWrapInfo(reply.WrapInfo)
	if err != nil {
		return nil, err
	}

	return info, nil
}

func (s *gRPCSystemViewClient) NewPluginClient(ctx context.Context, config pluginutil.PluginClientConfig) (pluginutil.PluginClient, error) {
	return nil, fmt.Errorf("cannot call NewPluginClient from a plugin backend")
}

func (s *gRPCSystemViewClient) LookupPlugin(_ context.Context, _ string, _ consts.PluginType) (*pluginutil.PluginRunner, error) {
	return nil, fmt.Errorf("cannot call LookupPlugin from a plugin backend")
}

func (s *gRPCSystemViewClient) LookupPluginVersion(_ context.Context, _ string, _ consts.PluginType, _ string) (*pluginutil.PluginRunner, error) {
	return nil, fmt.Errorf("cannot call LookupPluginVersion from a plugin backend")
}

func (s *gRPCSystemViewClient) ListVersionedPlugins(_ context.Context, _ consts.PluginType) ([]pluginutil.VersionedPlugin, error) {
	return nil, fmt.Errorf("cannot call ListVersionedPlugins from a plugin backend")
}

func (s *gRPCSystemViewClient) MlockEnabled() bool {
	reply, err := s.client.MlockEnabled(context.Background(), &pb.Empty{})
	if err != nil {
		return false
	}

	return reply.Enabled
}

func (s *gRPCSystemViewClient) HasFeature(feature license.Features) bool {
	// Not implemented
	return false
}

func (s *gRPCSystemViewClient) LocalMount() bool {
	reply, err := s.client.LocalMount(context.Background(), &pb.Empty{})
	if err != nil {
		return false
	}

	return reply.Local
}

func (s *gRPCSystemViewClient) EntityInfo(entityID string) (*logical.Entity, error) {
	reply, err := s.client.EntityInfo(context.Background(), &pb.EntityInfoArgs{
		EntityID: entityID,
	})
	if err != nil {
		return nil, err
	}
	if reply.Err != "" {
		return nil, errors.New(reply.Err)
	}

	return reply.Entity, nil
}

func (s *gRPCSystemViewClient) GroupsForEntity(entityID string) ([]*logical.Group, error) {
	reply, err := s.client.GroupsForEntity(context.Background(), &pb.EntityInfoArgs{
		EntityID: entityID,
	})
	if err != nil {
		return nil, err
	}
	if reply.Err != "" {
		return nil, errors.New(reply.Err)
	}

	return reply.Groups, nil
}

func (s *gRPCSystemViewClient) PluginEnv(ctx context.Context) (*logical.PluginEnvironment, error) {
	reply, err := s.client.PluginEnv(ctx, &pb.Empty{})
	if err != nil {
		return nil, err
	}

	return reply.PluginEnvironment, nil
}

func (s *gRPCSystemViewClient) GeneratePasswordFromPolicy(ctx context.Context, policyName string) (password string, err error) {
	req := &pb.GeneratePasswordFromPolicyRequest{
		PolicyName: policyName,
	}
	resp, err := s.client.GeneratePasswordFromPolicy(ctx, req)
	if err != nil {
		return "", err
	}
	return resp.Password, nil
}

type gRPCSystemViewServer struct {
	pb.UnimplementedSystemViewServer

	impl logical.SystemView
}

func (s *gRPCSystemViewServer) DefaultLeaseTTL(ctx context.Context, _ *pb.Empty) (*pb.TTLReply, error) {
	ttl := s.impl.DefaultLeaseTTL()
	return &pb.TTLReply{
		TTL: int64(ttl),
	}, nil
}

func (s *gRPCSystemViewServer) MaxLeaseTTL(ctx context.Context, _ *pb.Empty) (*pb.TTLReply, error) {
	ttl := s.impl.MaxLeaseTTL()
	return &pb.TTLReply{
		TTL: int64(ttl),
	}, nil
}

func (s *gRPCSystemViewServer) Tainted(ctx context.Context, _ *pb.Empty) (*pb.TaintedReply, error) {
	tainted := s.impl.Tainted()
	return &pb.TaintedReply{
		Tainted: tainted,
	}, nil
}

func (s *gRPCSystemViewServer) CachingDisabled(ctx context.Context, _ *pb.Empty) (*pb.CachingDisabledReply, error) {
	cachingDisabled := s.impl.CachingDisabled()
	return &pb.CachingDisabledReply{
		Disabled: cachingDisabled,
	}, nil
}

func (s *gRPCSystemViewServer) ReplicationState(ctx context.Context, _ *pb.Empty) (*pb.ReplicationStateReply, error) {
	replicationState := s.impl.ReplicationState()
	return &pb.ReplicationStateReply{
		State: int32(replicationState),
	}, nil
}

func (s *gRPCSystemViewServer) ResponseWrapData(ctx context.Context, args *pb.ResponseWrapDataArgs) (*pb.ResponseWrapDataReply, error) {
	data := map[string]interface{}{}
	err := json.Unmarshal([]byte(args.Data), &data)
	if err != nil {
		return &pb.ResponseWrapDataReply{}, err
	}

	// Do not allow JWTs to be returned
	info, err := s.impl.ResponseWrapData(ctx, data, time.Duration(args.TTL), false)
	if err != nil {
		return &pb.ResponseWrapDataReply{
			Err: pb.ErrToString(err),
		}, nil
	}

	pbInfo, err := pb.LogicalResponseWrapInfoToProtoResponseWrapInfo(info)
	if err != nil {
		return &pb.ResponseWrapDataReply{}, err
	}

	return &pb.ResponseWrapDataReply{
		WrapInfo: pbInfo,
	}, nil
}

func (s *gRPCSystemViewServer) MlockEnabled(ctx context.Context, _ *pb.Empty) (*pb.MlockEnabledReply, error) {
	enabled := s.impl.MlockEnabled()
	return &pb.MlockEnabledReply{
		Enabled: enabled,
	}, nil
}

func (s *gRPCSystemViewServer) LocalMount(ctx context.Context, _ *pb.Empty) (*pb.LocalMountReply, error) {
	local := s.impl.LocalMount()
	return &pb.LocalMountReply{
		Local: local,
	}, nil
}

func (s *gRPCSystemViewServer) EntityInfo(ctx context.Context, args *pb.EntityInfoArgs) (*pb.EntityInfoReply, error) {
	entity, err := s.impl.EntityInfo(args.EntityID)
	if err != nil {
		return &pb.EntityInfoReply{
			Err: pb.ErrToString(err),
		}, nil
	}
	return &pb.EntityInfoReply{
		Entity: entity,
	}, nil
}

func (s *gRPCSystemViewServer) GroupsForEntity(ctx context.Context, args *pb.EntityInfoArgs) (*pb.GroupsForEntityReply, error) {
	groups, err := s.impl.GroupsForEntity(args.EntityID)
	if err != nil {
		return &pb.GroupsForEntityReply{
			Err: pb.ErrToString(err),
		}, nil
	}
	return &pb.GroupsForEntityReply{
		Groups: groups,
	}, nil
}

func (s *gRPCSystemViewServer) PluginEnv(ctx context.Context, _ *pb.Empty) (*pb.PluginEnvReply, error) {
	pluginEnv, err := s.impl.PluginEnv(ctx)
	if err != nil {
		return &pb.PluginEnvReply{
			Err: pb.ErrToString(err),
		}, nil
	}
	return &pb.PluginEnvReply{
		PluginEnvironment: pluginEnv,
	}, nil
}

func (s *gRPCSystemViewServer) GeneratePasswordFromPolicy(ctx context.Context, req *pb.GeneratePasswordFromPolicyRequest) (*pb.GeneratePasswordFromPolicyReply, error) {
	policyName := req.PolicyName
	if policyName == "" {
		return &pb.GeneratePasswordFromPolicyReply{}, status.Errorf(codes.InvalidArgument, "no password policy specified")
	}

	password, err := s.impl.GeneratePasswordFromPolicy(ctx, policyName)
	if err != nil {
		return &pb.GeneratePasswordFromPolicyReply{}, status.Errorf(codes.Internal, "failed to generate password")
	}

	resp := &pb.GeneratePasswordFromPolicyReply{
		Password: password,
	}
	return resp, nil
}
