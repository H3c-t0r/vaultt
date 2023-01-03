package internal

import (
	"context"

	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/physical/raft"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
)

type WrappedCoreNodeStatus interface {
	GetSealStatus(ctx context.Context) (*vault.SealStatusResponse, error)
	ReplicationState() consts.ReplicationState
}

var _ WrappedCoreNodeStatus = &vault.Core{}

type WrappedCoreStandbyStates interface {
	StandbyStates() (bool, bool)
}

var _ WrappedCoreStandbyStates = &vault.Core{}

type WrappedCoreHCPToken interface {
	Sealed() bool
	CreateToken(context.Context, *logical.TokenEntry) error
	WrappedCoreStandbyStates
}

var _ WrappedCoreHCPToken = &vault.Core{}

type WrappedCoreMeta interface {
	NamespaceByID(ctx context.Context, nsID string) (*namespace.Namespace, error)
	ListNamespaces(includePath bool) []*namespace.Namespace
	ListMounts() ([]*vault.MountEntry, error)
	ListAuths() ([]*vault.MountEntry, error)
	HAEnabled() bool
	HAState() consts.HAState
	GetHAPeerNodesCached() []vault.PeerNode
	GetRaftConfiguration(ctx context.Context) (*raft.RaftConfigurationResponse, error)
	GetRaftAutopilotState(ctx context.Context) (*raft.AutopilotState, error)
	StorageType() string
}

var _ WrappedCoreMeta = &vault.Core{}

type WrappedCoreHCPLinkStatus interface {
	WrappedCoreStandbyStates
	SetHCPLinkStatus(status, name string)
	GetHCPLinkStatus() (string, string)
}

var _ WrappedCoreHCPLinkStatus = &vault.Core{}
