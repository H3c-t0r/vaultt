// Copyright © 2019, Oracle and/or its affiliates.
package ociauth

import "github.com/oracle/oci-go-sdk/common"

// Do not edit this file. This is based on standard OCI GO SDK format

// Stores a list of claims of a Principal
type Claim struct {
	Key    *string `json:"key"`
	Value  *string `json:"value"`
	Issuer *string `json:"issuer"`
}

// Prints the values of pointers in Claim,
// producing a human friendly string for an struct with pointers. Useful when debugging the values of a struct.
func (m Claim) String() string {
	return common.PointerString(m)
}

// Stores the details about a Principal
type Principal struct {
	TenantId  *string `json:"tenantId"`
	SubjectId *string `json:"subjectId"`
	Claims    []Claim `json:"claims"`
}

// Prints the values of pointers in Principal,
// producing a human friendly string for an struct with pointers. Useful when debugging the values of a struct.
func (m Principal) String() string {
	return common.PointerString(m)
}
