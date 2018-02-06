package pb

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/hashicorp/vault/helper/errutil"
	"github.com/hashicorp/vault/helper/wrapping"
	"github.com/hashicorp/vault/logical"
)

const (
	ErrTypeUnknown uint32 = iota
	ErrTypeUserError
	ErrTypeInternalError
	ErrTypeCodedError
	ErrTypeStatusBadRequest
	ErrTypeUnsupportedOperation
	ErrTypeUnsupportedPath
	ErrTypeInvalidRequest
	ErrTypePermissionDenied
	ErrTypeMultiAuthzPending
)

func ProtoErrToErr(e *ProtoError) error {
	if e == nil {
		return nil
	}

	var err error
	switch e.ErrType {
	case ErrTypeUnknown:
		err = errors.New(e.ErrMsg)
	case ErrTypeUserError:
		err = errutil.UserError{Err: e.ErrMsg}
	case ErrTypeInternalError:
		err = errutil.InternalError{Err: e.ErrMsg}
	case ErrTypeCodedError:
		err = logical.CodedError(int(e.ErrCode), e.ErrMsg)
	case ErrTypeStatusBadRequest:
		err = &logical.StatusBadRequest{Err: e.ErrMsg}
	case ErrTypeUnsupportedOperation:
		err = logical.ErrUnsupportedOperation
	case ErrTypeUnsupportedPath:
		err = logical.ErrUnsupportedPath
	case ErrTypeInvalidRequest:
		err = logical.ErrInvalidRequest
	case ErrTypePermissionDenied:
		err = logical.ErrPermissionDenied
	case ErrTypeMultiAuthzPending:
		err = logical.ErrMultiAuthzPending
	}

	return err
}

func ErrToProtoErr(e error) *ProtoError {
	if e == nil {
		return nil
	}
	pbErr := &ProtoError{
		ErrMsg:  e.Error(),
		ErrType: ErrTypeUnknown,
	}

	switch e.(type) {
	case errutil.UserError:
		pbErr.ErrType = ErrTypeUserError
	case errutil.InternalError:
		pbErr.ErrType = ErrTypeInternalError
	case logical.HTTPCodedError:
		pbErr.ErrType = ErrTypeCodedError
		pbErr.ErrCode = int64(e.(logical.HTTPCodedError).Code())
	case *logical.StatusBadRequest:
		pbErr.ErrType = ErrTypeStatusBadRequest
	}

	switch {
	case e == logical.ErrUnsupportedOperation:
		pbErr.ErrType = ErrTypeUnsupportedOperation
	case e == logical.ErrUnsupportedPath:
		pbErr.ErrType = ErrTypeUnsupportedPath
	case e == logical.ErrInvalidRequest:
		pbErr.ErrType = ErrTypeInvalidRequest
	case e == logical.ErrPermissionDenied:
		pbErr.ErrType = ErrTypePermissionDenied
	case e == logical.ErrMultiAuthzPending:
		pbErr.ErrType = ErrTypeMultiAuthzPending
	}

	return pbErr
}

func ErrToString(e error) string {
	if e == nil {
		return ""
	}

	return e.Error()
}

func LogicalStorageEntryToProtoStorageEntry(e *logical.StorageEntry) *StorageEntry {
	if e == nil {
		return nil
	}

	return &StorageEntry{
		Key:      e.Key,
		Value:    e.Value,
		SealWrap: e.SealWrap,
	}
}

func ProtoStorageEntryToLogicalStorageEntry(e *StorageEntry) *logical.StorageEntry {
	if e == nil {
		return nil
	}

	return &logical.StorageEntry{
		Key:      e.Key,
		Value:    e.Value,
		SealWrap: e.SealWrap,
	}
}

func ProtoLeaseOptionsToLogicalLeaseOptions(l *LeaseOptions) (logical.LeaseOptions, error) {
	if l == nil {
		return logical.LeaseOptions{}, nil
	}

	t, err := ptypes.Timestamp(l.IssueTime)
	return logical.LeaseOptions{
		TTL:       time.Duration(l.TTL),
		Renewable: l.Renewable,
		Increment: time.Duration(l.Increment),
		IssueTime: t,
	}, err
}

func LogicalLeaseOptionsToProtoLeaseOptions(l logical.LeaseOptions) (*LeaseOptions, error) {
	t, err := ptypes.TimestampProto(l.IssueTime)
	if err != nil {
		return nil, err
	}

	return &LeaseOptions{
		TTL:       int64(l.TTL),
		Renewable: l.Renewable,
		Increment: int64(l.Increment),
		IssueTime: t,
	}, err
}

func ProtoSecretToLogicalSecret(s *Secret) (*logical.Secret, error) {
	if s == nil {
		return nil, nil
	}

	data := map[string]interface{}{}
	err := json.Unmarshal([]byte(s.InternalData), &data)
	if err != nil {
		return nil, err
	}

	lease, err := ProtoLeaseOptionsToLogicalLeaseOptions(s.LeaseOptions)
	if err != nil {
		return nil, err
	}

	return &logical.Secret{
		LeaseOptions: lease,
		InternalData: data,
		LeaseID:      s.LeaseID,
	}, nil
}

func LogicalSecretToProtoSecret(s *logical.Secret) (*Secret, error) {
	if s == nil {
		return nil, nil
	}

	buf, err := json.Marshal(s.InternalData)
	if err != nil {
		return nil, err
	}

	lease, err := LogicalLeaseOptionsToProtoLeaseOptions(s.LeaseOptions)
	if err != nil {
		return nil, err
	}

	return &Secret{
		LeaseOptions: lease,
		InternalData: string(buf[:]),
		LeaseID:      s.LeaseID,
	}, err
}

func LogicalRequestToProtoRequest(r *logical.Request) (*Request, error) {
	if r == nil {
		return nil, nil
	}

	buf, err := json.Marshal(r.Data)
	if err != nil {
		return nil, err
	}

	secret, err := LogicalSecretToProtoSecret(r.Secret)
	if err != nil {
		return nil, err
	}

	auth, err := LogicalAuthToProtoAuth(r.Auth)
	if err != nil {
		return nil, err
	}

	headers := map[string]*Header{}
	for k, v := range r.Headers {
		headers[k] = &Header{v}
	}

	return &Request{
		ID:                       r.ID,
		ReplicationCluster:       r.ReplicationCluster,
		Operation:                string(r.Operation),
		Path:                     r.Path,
		Data:                     string(buf[:]),
		Secret:                   secret,
		Auth:                     auth,
		Headers:                  headers,
		ClientToken:              r.ClientToken,
		ClientTokenAccessor:      r.ClientTokenAccessor,
		DisplayName:              r.DisplayName,
		MountPoint:               r.MountPoint,
		MountType:                r.MountType,
		MountAccessor:            r.MountAccessor,
		WrapInfo:                 LogicalRequestWrapInfoToProtoRequestWrapInfo(r.WrapInfo),
		ClientTokenRemainingUses: int64(r.ClientTokenRemainingUses),
		//MFACreds: MFACreds,
		EntityID:        r.EntityID,
		PolicyOverride:  r.PolicyOverride,
		Unauthenticated: r.Unauthenticated,
	}, nil
}

func ProtoRequestToLogicalRequest(r *Request) (*logical.Request, error) {
	if r == nil {
		return nil, nil
	}

	data := map[string]interface{}{}
	err := json.Unmarshal([]byte(r.Data), &data)
	if err != nil {
		return nil, err
	}

	secret, err := ProtoSecretToLogicalSecret(r.Secret)
	if err != nil {
		return nil, err
	}

	auth, err := ProtoAuthToLogicalAuth(r.Auth)
	if err != nil {
		return nil, err
	}

	var headers map[string][]string
	if len(r.Headers) > 0 {
		headers = make(map[string][]string, len(r.Headers))
		for k, v := range r.Headers {
			headers[k] = v.Header
		}
	}

	return &logical.Request{
		ID:                       r.ID,
		ReplicationCluster:       r.ReplicationCluster,
		Operation:                logical.Operation(r.Operation),
		Path:                     r.Path,
		Data:                     data,
		Secret:                   secret,
		Auth:                     auth,
		Headers:                  headers,
		ClientToken:              r.ClientToken,
		ClientTokenAccessor:      r.ClientTokenAccessor,
		DisplayName:              r.DisplayName,
		MountPoint:               r.MountPoint,
		MountType:                r.MountType,
		MountAccessor:            r.MountAccessor,
		WrapInfo:                 ProtoRequestWrapInfoToLogicalRequestWrapInfo(r.WrapInfo),
		ClientTokenRemainingUses: int(r.ClientTokenRemainingUses),
		//MFACreds: MFACreds,
		EntityID:        r.EntityID,
		PolicyOverride:  r.PolicyOverride,
		Unauthenticated: r.Unauthenticated,
	}, nil
}

func LogicalRequestWrapInfoToProtoRequestWrapInfo(i *logical.RequestWrapInfo) *RequestWrapInfo {
	if i == nil {
		return nil
	}

	return &RequestWrapInfo{
		TTL:      int64(i.TTL),
		Format:   i.Format,
		SealWrap: i.SealWrap,
	}
}

func ProtoRequestWrapInfoToLogicalRequestWrapInfo(i *RequestWrapInfo) *logical.RequestWrapInfo {
	if i == nil {
		return nil
	}

	return &logical.RequestWrapInfo{
		TTL:      time.Duration(i.TTL),
		Format:   i.Format,
		SealWrap: i.SealWrap,
	}
}

func ProtoResponseToLogicalResponse(r *Response) (*logical.Response, error) {
	if r == nil {
		return nil, nil
	}

	secret, err := ProtoSecretToLogicalSecret(r.Secret)
	if err != nil {
		return nil, err
	}

	auth, err := ProtoAuthToLogicalAuth(r.Auth)
	if err != nil {
		return nil, err
	}

	data := map[string]interface{}{}
	err = json.Unmarshal([]byte(r.Data), &data)
	if err != nil {
		return nil, err
	}

	wrapInfo, err := ProtoResponseWrapInfoToLogicalResponseWrapInfo(r.WrapInfo)
	if err != nil {
		return nil, err
	}

	return &logical.Response{
		Secret:   secret,
		Auth:     auth,
		Data:     data,
		Redirect: r.Redirect,
		Warnings: r.Warnings,
		WrapInfo: wrapInfo,
	}, nil
}

func ProtoResponseWrapInfoToLogicalResponseWrapInfo(i *ResponseWrapInfo) (*wrapping.ResponseWrapInfo, error) {
	if i == nil {
		return nil, nil
	}

	t, err := ptypes.Timestamp(i.CreationTime)
	if err != nil {
		return nil, err
	}

	return &wrapping.ResponseWrapInfo{
		TTL:             time.Duration(i.TTL),
		Token:           i.Token,
		Accessor:        i.Accessor,
		CreationTime:    t,
		WrappedAccessor: i.WrappedAccessor,
		WrappedEntityID: i.WrappedEntityID,
		Format:          i.Format,
		CreationPath:    i.CreationPath,
		SealWrap:        i.SealWrap,
	}, nil
}

func LogicalResponseWrapInfoToProtoResponseWrapInfo(i *wrapping.ResponseWrapInfo) (*ResponseWrapInfo, error) {
	if i == nil {
		return nil, nil
	}

	t, err := ptypes.TimestampProto(i.CreationTime)
	if err != nil {
		return nil, err
	}

	return &ResponseWrapInfo{
		TTL:             int64(i.TTL),
		Token:           i.Token,
		Accessor:        i.Accessor,
		CreationTime:    t,
		WrappedAccessor: i.WrappedAccessor,
		WrappedEntityID: i.WrappedEntityID,
		Format:          i.Format,
		CreationPath:    i.CreationPath,
		SealWrap:        i.SealWrap,
	}, nil
}

func LogicalResponseToProtoResponse(r *logical.Response) (*Response, error) {
	if r == nil {
		return nil, nil
	}

	secret, err := LogicalSecretToProtoSecret(r.Secret)
	if err != nil {
		return nil, err
	}

	auth, err := LogicalAuthToProtoAuth(r.Auth)
	if err != nil {
		return nil, err
	}

	buf, err := json.Marshal(r.Data)
	if err != nil {
		return nil, err
	}

	wrapInfo, err := LogicalResponseWrapInfoToProtoResponseWrapInfo(r.WrapInfo)
	if err != nil {
		return nil, err
	}

	return &Response{
		Secret:   secret,
		Auth:     auth,
		Data:     string(buf[:]),
		Redirect: r.Redirect,
		Warnings: r.Warnings,
		WrapInfo: wrapInfo,
	}, nil
}

func LogicalAliasToProtoAlias(a *logical.Alias) *Alias {
	if a == nil {
		return nil
	}

	return &Alias{
		MountType:     a.MountType,
		MountAccessor: a.MountAccessor,
		Name:          a.Name,
	}
}

func ProtoAliasToLogicalAlias(a *Alias) *logical.Alias {
	if a == nil {
		return nil
	}

	return &logical.Alias{
		MountType:     a.MountType,
		MountAccessor: a.MountAccessor,
		Name:          a.Name,
	}
}

func LogicalAuthToProtoAuth(a *logical.Auth) (*Auth, error) {
	if a == nil {
		return nil, nil
	}

	buf, err := json.Marshal(a.InternalData)
	if err != nil {
		return nil, err
	}

	groupAliases := make([]*Alias, len(a.GroupAliases))
	for i, al := range a.GroupAliases {
		groupAliases[i] = LogicalAliasToProtoAlias(al)
	}

	lo, err := LogicalLeaseOptionsToProtoLeaseOptions(a.LeaseOptions)
	if err != nil {
		return nil, err
	}

	return &Auth{
		LeaseOptions: lo,
		InternalData: string(buf[:]),
		DisplayName:  a.DisplayName,
		Policies:     a.Policies,
		Metadata:     a.Metadata,
		ClientToken:  a.ClientToken,
		Accessor:     a.Accessor,
		Period:       int64(a.Period),
		NumUses:      int64(a.NumUses),
		EntityID:     a.EntityID,
		Alias:        LogicalAliasToProtoAlias(a.Alias),
		GroupAliases: groupAliases,
	}, nil
}

func ProtoAuthToLogicalAuth(a *Auth) (*logical.Auth, error) {
	if a == nil {
		return nil, nil
	}

	data := map[string]interface{}{}
	err := json.Unmarshal([]byte(a.InternalData), &data)
	if err != nil {
		return nil, err
	}

	groupAliases := make([]*logical.Alias, len(a.GroupAliases))
	for i, al := range a.GroupAliases {
		groupAliases[i] = ProtoAliasToLogicalAlias(al)
	}

	lo, err := ProtoLeaseOptionsToLogicalLeaseOptions(a.LeaseOptions)
	if err != nil {
		return nil, err
	}

	return &logical.Auth{
		LeaseOptions: lo,
		InternalData: data,
		DisplayName:  a.DisplayName,
		Policies:     a.Policies,
		Metadata:     a.Metadata,
		ClientToken:  a.ClientToken,
		Accessor:     a.Accessor,
		Period:       time.Duration(a.Period),
		NumUses:      int(a.NumUses),
		EntityID:     a.EntityID,
		Alias:        ProtoAliasToLogicalAlias(a.Alias),
		GroupAliases: groupAliases,
	}, nil
}
