package vault

import (
	"reflect"
	"testing"

	"github.com/hashicorp/vault/helper/strutil"
	"github.com/hashicorp/vault/logical"
)

func TestACL_Capabilities(t *testing.T) {
	// Create the root policy ACL
	policy := []*Policy{&Policy{Name: "root"}}
	acl, err := NewACL(policy)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	actual := acl.Capabilities("any/path")
	expected := []string{"root"}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("bad: got\n%#v\nexpected\n%#v\n", actual, expected)
	}

	policies, err := Parse(aclPolicy)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	acl, err = NewACL([]*Policy{policies})
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	actual = acl.Capabilities("dev")
	expected = []string{"deny"}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("bad: path:%s\ngot\n%#v\nexpected\n%#v\n", "deny", actual, expected)
	}

	actual = acl.Capabilities("dev/")
	expected = []string{"sudo", "read", "list", "update", "delete", "create"}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("bad: path:%s\ngot\n%#v\nexpected\n%#v\n", "dev/", actual, expected)
	}

	actual = acl.Capabilities("stage/aws/test")
	expected = []string{"sudo", "read", "list", "update"}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("bad: path:%s\ngot\n%#v\nexpected\n%#v\n", "stage/aws/test", actual, expected)
	}

}

func TestACL_Root(t *testing.T) {
	// Create the root policy ACL
	policy := []*Policy{&Policy{Name: "root"}}
	acl, err := NewACL(policy)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	allowed, rootPrivs, _, _ := acl.AllowOperation(logical.UpdateOperation, "sys/mount/foo")
	if !rootPrivs {
		t.Fatalf("expected root")
	}
	if !allowed {
		t.Fatalf("expected permission")
	}
}

func TestACL_Single(t *testing.T) {
	policy, err := Parse(aclPolicy)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	acl, err := NewACL([]*Policy{policy})
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Type of operation is not important here as we only care about checking
	// sudo/root
	_, rootPrivs, _, _ := acl.AllowOperation(logical.ReadOperation, "sys/mount/foo")
	if rootPrivs {
		t.Fatalf("unexpected root")
	}

	type tcase struct {
		op        logical.Operation
		path      string
		allowed   bool
		rootPrivs bool
	}
	tcases := []tcase{
		{logical.ReadOperation, "root", false, false},
		{logical.HelpOperation, "root", true, false},

		{logical.ReadOperation, "dev/foo", true, true},
		{logical.UpdateOperation, "dev/foo", true, true},

		{logical.DeleteOperation, "stage/foo", true, false},
		{logical.ListOperation, "stage/aws/foo", true, true},
		{logical.UpdateOperation, "stage/aws/foo", true, true},
		{logical.UpdateOperation, "stage/aws/policy/foo", true, true},

		{logical.DeleteOperation, "prod/foo", false, false},
		{logical.UpdateOperation, "prod/foo", false, false},
		{logical.ReadOperation, "prod/foo", true, false},
		{logical.ListOperation, "prod/foo", true, false},
		{logical.ReadOperation, "prod/aws/foo", false, false},

		{logical.ReadOperation, "foo/bar", true, true},
		{logical.ListOperation, "foo/bar", false, true},
		{logical.UpdateOperation, "foo/bar", false, true},
		{logical.CreateOperation, "foo/bar", true, true},
	}

	for _, tc := range tcases {
		allowed, rootPrivs, _, _ := acl.AllowOperation(tc.op, tc.path)
		if allowed != tc.allowed {
			t.Fatalf("bad: case %#v: %v, %v", tc, allowed, rootPrivs)
		}
		if rootPrivs != tc.rootPrivs {
			t.Fatalf("bad: case %#v: %v, %v", tc, allowed, rootPrivs)
		}
	}
}

func TestACL_Layered(t *testing.T) {
	policy1, err := Parse(aclPolicy)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	policy2, err := Parse(aclPolicy2)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	acl, err := NewACL([]*Policy{policy1, policy2})
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	testLayeredACL(t, acl, false)
}

func TestACL_Layered_With_MFA(t *testing.T) {
	policy1, err := Parse(aclPolicyWithMFA)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	policy2, err := Parse(aclPolicyWithMFA2)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	acl, err := NewACL([]*Policy{policy1, policy2})
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	testLayeredACL(t, acl, true)
}

func testLayeredACL(t *testing.T, acl *ACL, checkMFAMethods bool) {
	// Type of operation is not important here as we only care about checking
	// sudo/root
	_, rootPrivs, _, _ := acl.AllowOperation(logical.ReadOperation, "sys/mount/foo")
	if rootPrivs {
		t.Fatalf("unexpected root")
	}

	type tcase struct {
		op         logical.Operation
		path       string
		allowed    bool
		rootPrivs  bool
		mfaMethods []string
	}
	tcases := []tcase{
		{logical.ReadOperation, "root", false, false, nil},
		{logical.HelpOperation, "root", true, false, nil},

		{logical.ReadOperation, "dev/foo", true, true, nil},
		{logical.UpdateOperation, "dev/foo", true, true, nil},
		{logical.ReadOperation, "dev/hide/foo", false, false, nil},
		{logical.UpdateOperation, "dev/hide/foo", false, false, nil},

		{logical.DeleteOperation, "stage/foo", true, false, nil},
		{logical.ListOperation, "stage/aws/foo", true, true, nil},
		{logical.UpdateOperation, "stage/aws/foo", true, true, []string{"mfa1", "mfa2"}},
		{logical.UpdateOperation, "stage/aws/policy/foo", false, false, nil},

		{logical.DeleteOperation, "prod/foo", true, false, []string{"mfa3"}},
		{logical.UpdateOperation, "prod/foo", true, false, []string{"mfa3"}},
		{logical.ReadOperation, "prod/foo", true, false, []string{"mfa3"}},
		{logical.ListOperation, "prod/foo", true, false, []string{"mfa3"}},
		{logical.ReadOperation, "prod/aws/foo", false, false, nil},

		{logical.ReadOperation, "sys/status", false, false, nil},
		{logical.UpdateOperation, "sys/seal", true, true, nil},

		{logical.ReadOperation, "foo/bar", false, false, nil},
		{logical.ListOperation, "foo/bar", false, false, nil},
		{logical.UpdateOperation, "foo/bar", false, false, nil},
		{logical.CreateOperation, "foo/bar", false, false, nil},
	}

	for _, tc := range tcases {
		allowed, rootPrivs, mfaMethods, sudoMFAMethods := acl.AllowOperation(tc.op, tc.path)
		if allowed != tc.allowed {
			t.Fatalf("bad: case %#v: %v, %v", tc, allowed, rootPrivs)
		}
		if rootPrivs != tc.rootPrivs {
			t.Fatalf("bad: case %#v: %v, %v", tc, allowed, rootPrivs)
		}
		if checkMFAMethods {
			if rootPrivs {
				// These are additive so we don't actually require exactly
				// what's above, just make sure it's a subset
				if !strutil.StrListSubset(sudoMFAMethods, tc.mfaMethods) {
					t.Fatalf("bad: case %#v: %v, %v, %v", tc, allowed, rootPrivs, sudoMFAMethods)
				}
			} else {
				if !reflect.DeepEqual(mfaMethods, tc.mfaMethods) {
					t.Fatalf("bad: case %#v: %v, %v, %v", tc, allowed, rootPrivs, mfaMethods)
				}
			}
		}
	}
}

var aclPolicy = `
name = "dev"
path "dev/*" {
	policy = "sudo"
}
path "stage/*" {
	policy = "write"
}
path "stage/aws/*" {
	policy = "read"
	capabilities = ["update", "sudo"]
}
path "stage/aws/policy/*" {
	policy = "sudo"
}
path "prod/*" {
	policy = "read"
}
path "prod/aws/*" {
	policy = "deny"
}
path "sys/*" {
	policy = "deny"
}
path "foo/bar" {
	capabilities = ["read", "create", "sudo"]
}
`

var aclPolicy2 = `
name = "ops"
path "dev/hide/*" {
	policy = "deny"
}
path "stage/aws/policy/*" {
	policy = "deny"
	# This should have no effect
	capabilities = ["read", "update", "sudo"]
}
path "prod/*" {
	policy = "write"
}
path "sys/seal" {
	policy = "sudo"
}
path "foo/bar" {
	capabilities = ["deny"]
}
`

var aclPolicyWithMFA = `
name = "dev"
path "dev/*" {
	policy = "sudo"
}
path "stage/*" {
	policy = "write"
}
path "stage/aws/*" {
	policy = "read"
	capabilities = ["update", "sudo"]
}
path "stage/aws/policy/*" {
	policy = "sudo"
}
path "prod/*" {
	policy = "read"
}
path "prod/aws/*" {
	policy = "deny"
}
path "sys/*" {
	policy = "deny"
}
path "foo/bar" {
	capabilities = ["read", "create", "sudo"]
}
`

var aclPolicyWithMFA2 = `
name = "ops"
path "dev/hide/*" {
	policy = "deny"
}
path "stage/aws/*" {
	capabilities = ["update", "sudo"]
	mfa_methods = ["mfa1", "mfa2"]
}
path "stage/aws/policy/*" {
	policy = "deny"
	# This should have no effect
	capabilities = ["read", "update", "sudo"]
	mfa_methods = ["mfa1", "mfa2"]
}
path "prod/*" {
	policy = "write"
	mfa_methods = ["mfa3"]
}
path "sys/seal" {
	policy = "sudo"
}
path "foo/bar" {
	capabilities = ["deny"]
}
`
