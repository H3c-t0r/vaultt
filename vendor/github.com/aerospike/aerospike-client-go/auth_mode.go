// Copyright 2013-2020 Aerospike, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package aerospike

// AuthMode determines authentication mode when user/password is defined.
type AuthMode int

const (
	// AuthModeInternal uses internal authentication only.  Hashed password is stored on the server.
	// Do not send clear password. This is the default.
	AuthModeInternal AuthMode = iota

	// AuthModeExternal uses external authentication (like LDAP).  Specific external authentication is
	// configured on server.  If TLSConfig is defined, sends clear password on node login via TLS.
	// Will return an error if TLSConfig is not defined.
	AuthModeExternal
)
