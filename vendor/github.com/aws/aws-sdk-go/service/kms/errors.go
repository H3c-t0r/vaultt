// Code generated by private/model/cli/gen-api/main.go. DO NOT EDIT.

package kms

const (

	// ErrCodeAlreadyExistsException for service response error code
	// "AlreadyExistsException".
	//
	// The request was rejected because it attempted to create a resource that already
	// exists.
	ErrCodeAlreadyExistsException = "AlreadyExistsException"

	// ErrCodeCloudHsmClusterInUseException for service response error code
	// "CloudHsmClusterInUseException".
	//
	// The request was rejected because the specified AWS CloudHSM cluster is already
	// associated with a custom key store or it shares a backup history with a cluster
	// that is associated with a custom key store. Each custom key store must be
	// associated with a different AWS CloudHSM cluster.
	//
	// Clusters that share a backup history have the same cluster certificate. To
	// view the cluster certificate of a cluster, use the DescribeClusters (http://docs.aws.amazon.com/cloudhsm/latest/APIReference/API_DescribeClusters.html)
	// operation.
	ErrCodeCloudHsmClusterInUseException = "CloudHsmClusterInUseException"

	// ErrCodeCloudHsmClusterInvalidConfigurationException for service response error code
	// "CloudHsmClusterInvalidConfigurationException".
	//
	// The request was rejected because the associated AWS CloudHSM cluster did
	// not meet the configuration requirements for a custom key store. The cluster
	// must be configured with private subnets in at least two different Availability
	// Zones in the Region. Also, it must contain at least as many HSMs as the operation
	// requires.
	//
	// For the CreateCustomKeyStore, UpdateCustomKeyStore, and CreateKey operations,
	// the AWS CloudHSM cluster must have at least two active HSMs, each in a different
	// Availability Zone. For the ConnectCustomKeyStore operation, the AWS CloudHSM
	// must contain at least one active HSM.
	//
	// For information about creating a private subnet for a AWS CloudHSM cluster,
	// see Create a Private Subnet (http://docs.aws.amazon.com/cloudhsm/latest/userguide/create-subnets.html)
	// in the AWS CloudHSM User Guide. To add HSMs, use the AWS CloudHSM CreateHsm
	// (http://docs.aws.amazon.com/cloudhsm/latest/APIReference/API_CreateHsm.html)
	// operation.
	ErrCodeCloudHsmClusterInvalidConfigurationException = "CloudHsmClusterInvalidConfigurationException"

	// ErrCodeCloudHsmClusterNotActiveException for service response error code
	// "CloudHsmClusterNotActiveException".
	//
	// The request was rejected because the AWS CloudHSM cluster that is associated
	// with the custom key store is not active. Initialize and activate the cluster
	// and try the command again. For detailed instructions, see Getting Started
	// (http://docs.aws.amazon.com/cloudhsm/latest/userguide/getting-started.html)
	// in the AWS CloudHSM User Guide.
	ErrCodeCloudHsmClusterNotActiveException = "CloudHsmClusterNotActiveException"

	// ErrCodeCloudHsmClusterNotFoundException for service response error code
	// "CloudHsmClusterNotFoundException".
	//
	// The request was rejected because AWS KMS cannot find the AWS CloudHSM cluster
	// with the specified cluster ID. Retry the request with a different cluster
	// ID.
	ErrCodeCloudHsmClusterNotFoundException = "CloudHsmClusterNotFoundException"

	// ErrCodeCloudHsmClusterNotRelatedException for service response error code
	// "CloudHsmClusterNotRelatedException".
	//
	// The request was rejected because the specified AWS CloudHSM cluster has a
	// different cluster certificate than the original cluster. You cannot use the
	// operation to specify an unrelated cluster.
	//
	// Specify a cluster that shares a backup history with the original cluster.
	// This includes clusters that were created from a backup of the current cluster,
	// and clusters that were created from the same backup that produced the current
	// cluster.
	//
	// Clusters that share a backup history have the same cluster certificate. To
	// view the cluster certificate of a cluster, use the DescribeClusters (http://docs.aws.amazon.com/cloudhsm/latest/APIReference/API_DescribeClusters.html)
	// operation.
	ErrCodeCloudHsmClusterNotRelatedException = "CloudHsmClusterNotRelatedException"

	// ErrCodeCustomKeyStoreHasCMKsException for service response error code
	// "CustomKeyStoreHasCMKsException".
	//
	// The request was rejected because the custom key store contains AWS KMS customer
	// master keys (CMKs). After verifying that you do not need to use the CMKs,
	// use the ScheduleKeyDeletion operation to delete the CMKs. After they are
	// deleted, you can delete the custom key store.
	ErrCodeCustomKeyStoreHasCMKsException = "CustomKeyStoreHasCMKsException"

	// ErrCodeCustomKeyStoreInvalidStateException for service response error code
	// "CustomKeyStoreInvalidStateException".
	//
	// The request was rejected because of the ConnectionState of the custom key
	// store. To get the ConnectionState of a custom key store, use the DescribeCustomKeyStores
	// operation.
	//
	// This exception is thrown under the following conditions:
	//
	//    * You requested the CreateKey or GenerateRandom operation in a custom
	//    key store that is not connected. These operations are valid only when
	//    the custom key store ConnectionState is CONNECTED.
	//
	//    * You requested the UpdateCustomKeyStore or DeleteCustomKeyStore operation
	//    on a custom key store that is not disconnected. This operation is valid
	//    only when the custom key store ConnectionState is DISCONNECTED.
	//
	//    * You requested the ConnectCustomKeyStore operation on a custom key store
	//    with a ConnectionState of DISCONNECTING or FAILED. This operation is valid
	//    for all other ConnectionState values.
	ErrCodeCustomKeyStoreInvalidStateException = "CustomKeyStoreInvalidStateException"

	// ErrCodeCustomKeyStoreNameInUseException for service response error code
	// "CustomKeyStoreNameInUseException".
	//
	// The request was rejected because the specified custom key store name is already
	// assigned to another custom key store in the account. Try again with a custom
	// key store name that is unique in the account.
	ErrCodeCustomKeyStoreNameInUseException = "CustomKeyStoreNameInUseException"

	// ErrCodeCustomKeyStoreNotFoundException for service response error code
	// "CustomKeyStoreNotFoundException".
	//
	// The request was rejected because AWS KMS cannot find a custom key store with
	// the specified key store name or ID.
	ErrCodeCustomKeyStoreNotFoundException = "CustomKeyStoreNotFoundException"

	// ErrCodeDependencyTimeoutException for service response error code
	// "DependencyTimeoutException".
	//
	// The system timed out while trying to fulfill the request. The request can
	// be retried.
	ErrCodeDependencyTimeoutException = "DependencyTimeoutException"

	// ErrCodeDisabledException for service response error code
	// "DisabledException".
	//
	// The request was rejected because the specified CMK is not enabled.
	ErrCodeDisabledException = "DisabledException"

	// ErrCodeExpiredImportTokenException for service response error code
	// "ExpiredImportTokenException".
	//
	// The request was rejected because the provided import token is expired. Use
	// GetParametersForImport to get a new import token and public key, use the
	// new public key to encrypt the key material, and then try the request again.
	ErrCodeExpiredImportTokenException = "ExpiredImportTokenException"

	// ErrCodeIncorrectKeyMaterialException for service response error code
	// "IncorrectKeyMaterialException".
	//
	// The request was rejected because the provided key material is invalid or
	// is not the same key material that was previously imported into this customer
	// master key (CMK).
	ErrCodeIncorrectKeyMaterialException = "IncorrectKeyMaterialException"

	// ErrCodeIncorrectTrustAnchorException for service response error code
	// "IncorrectTrustAnchorException".
	//
	// The request was rejected because the trust anchor certificate in the request
	// is not the trust anchor certificate for the specified AWS CloudHSM cluster.
	//
	// When you initialize the cluster (http://docs.aws.amazon.com/cloudhsm/latest/userguide/initialize-cluster.html#sign-csr),
	// you create the trust anchor certificate and save it in the customerCA.crt
	// file.
	ErrCodeIncorrectTrustAnchorException = "IncorrectTrustAnchorException"

	// ErrCodeInternalException for service response error code
	// "KMSInternalException".
	//
	// The request was rejected because an internal exception occurred. The request
	// can be retried.
	ErrCodeInternalException = "KMSInternalException"

	// ErrCodeInvalidAliasNameException for service response error code
	// "InvalidAliasNameException".
	//
	// The request was rejected because the specified alias name is not valid.
	ErrCodeInvalidAliasNameException = "InvalidAliasNameException"

	// ErrCodeInvalidArnException for service response error code
	// "InvalidArnException".
	//
	// The request was rejected because a specified ARN was not valid.
	ErrCodeInvalidArnException = "InvalidArnException"

	// ErrCodeInvalidCiphertextException for service response error code
	// "InvalidCiphertextException".
	//
	// The request was rejected because the specified ciphertext, or additional
	// authenticated data incorporated into the ciphertext, such as the encryption
	// context, is corrupted, missing, or otherwise invalid.
	ErrCodeInvalidCiphertextException = "InvalidCiphertextException"

	// ErrCodeInvalidGrantIdException for service response error code
	// "InvalidGrantIdException".
	//
	// The request was rejected because the specified GrantId is not valid.
	ErrCodeInvalidGrantIdException = "InvalidGrantIdException"

	// ErrCodeInvalidGrantTokenException for service response error code
	// "InvalidGrantTokenException".
	//
	// The request was rejected because the specified grant token is not valid.
	ErrCodeInvalidGrantTokenException = "InvalidGrantTokenException"

	// ErrCodeInvalidImportTokenException for service response error code
	// "InvalidImportTokenException".
	//
	// The request was rejected because the provided import token is invalid or
	// is associated with a different customer master key (CMK).
	ErrCodeInvalidImportTokenException = "InvalidImportTokenException"

	// ErrCodeInvalidKeyUsageException for service response error code
	// "InvalidKeyUsageException".
	//
	// The request was rejected because the specified KeySpec value is not valid.
	ErrCodeInvalidKeyUsageException = "InvalidKeyUsageException"

	// ErrCodeInvalidMarkerException for service response error code
	// "InvalidMarkerException".
	//
	// The request was rejected because the marker that specifies where pagination
	// should next begin is not valid.
	ErrCodeInvalidMarkerException = "InvalidMarkerException"

	// ErrCodeInvalidStateException for service response error code
	// "KMSInvalidStateException".
	//
	// The request was rejected because the state of the specified resource is not
	// valid for this request.
	//
	// For more information about how key state affects the use of a CMK, see How
	// Key State Affects Use of a Customer Master Key (http://docs.aws.amazon.com/kms/latest/developerguide/key-state.html)
	// in the AWS Key Management Service Developer Guide.
	ErrCodeInvalidStateException = "KMSInvalidStateException"

	// ErrCodeKeyUnavailableException for service response error code
	// "KeyUnavailableException".
	//
	// The request was rejected because the specified CMK was not available. The
	// request can be retried.
	ErrCodeKeyUnavailableException = "KeyUnavailableException"

	// ErrCodeLimitExceededException for service response error code
	// "LimitExceededException".
	//
	// The request was rejected because a limit was exceeded. For more information,
	// see Limits (http://docs.aws.amazon.com/kms/latest/developerguide/limits.html)
	// in the AWS Key Management Service Developer Guide.
	ErrCodeLimitExceededException = "LimitExceededException"

	// ErrCodeMalformedPolicyDocumentException for service response error code
	// "MalformedPolicyDocumentException".
	//
	// The request was rejected because the specified policy is not syntactically
	// or semantically correct.
	ErrCodeMalformedPolicyDocumentException = "MalformedPolicyDocumentException"

	// ErrCodeNotFoundException for service response error code
	// "NotFoundException".
	//
	// The request was rejected because the specified entity or resource could not
	// be found.
	ErrCodeNotFoundException = "NotFoundException"

	// ErrCodeTagException for service response error code
	// "TagException".
	//
	// The request was rejected because one or more tags are not valid.
	ErrCodeTagException = "TagException"

	// ErrCodeUnsupportedOperationException for service response error code
	// "UnsupportedOperationException".
	//
	// The request was rejected because a specified parameter is not supported or
	// a specified resource is not valid for this operation.
	ErrCodeUnsupportedOperationException = "UnsupportedOperationException"
)
