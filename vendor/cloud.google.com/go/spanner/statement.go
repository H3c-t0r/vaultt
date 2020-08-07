/*
Copyright 2017 Google LLC

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package spanner

import (
	"errors"
	"fmt"

	proto3 "github.com/golang/protobuf/ptypes/struct"
	structpb "github.com/golang/protobuf/ptypes/struct"
	sppb "google.golang.org/genproto/googleapis/spanner/v1"
	"google.golang.org/grpc/codes"
)

// A Statement is a SQL query with named parameters.
//
// A parameter placeholder consists of '@' followed by the parameter name.
// Parameter names consist of any combination of letters, numbers, and
// underscores. Names may be entirely numeric (e.g., "WHERE m.id = @5").
// Parameters may appear anywhere that a literal value is expected. The same
// parameter name may be used more than once.  It is an error to execute a
// statement with unbound parameters. On the other hand, it is allowable to
// bind parameter names that are not used.
//
// See the documentation of the Row type for how Go types are mapped to Cloud
// Spanner types.
type Statement struct {
	SQL    string
	Params map[string]interface{}
}

// NewStatement returns a Statement with the given SQL and an empty Params map.
func NewStatement(sql string) Statement {
	return Statement{SQL: sql, Params: map[string]interface{}{}}
}

var (
	errNilParam = errors.New("use T(nil), not nil")
	errNoType   = errors.New("no type information")
)

// convertParams converts a statement's parameters into proto Param and
// ParamTypes.
func (s *Statement) convertParams() (*structpb.Struct, map[string]*sppb.Type, error) {
	params := &proto3.Struct{
		Fields: map[string]*proto3.Value{},
	}
	paramTypes := map[string]*sppb.Type{}
	for k, v := range s.Params {
		if v == nil {
			return nil, nil, errBindParam(k, v, errNilParam)
		}
		val, t, err := encodeValue(v)
		if err != nil {
			return nil, nil, errBindParam(k, v, err)
		}
		if t == nil { // should not happen, because of nil check above
			return nil, nil, errBindParam(k, v, errNoType)
		}
		params.Fields[k] = val
		paramTypes[k] = t
	}

	return params, paramTypes, nil
}

// errBindParam returns error for not being able to bind parameter to query
// request.
func errBindParam(k string, v interface{}, err error) error {
	if err == nil {
		return nil
	}
	var se *Error
	if !errorAs(err, &se) {
		return spannerErrorf(codes.InvalidArgument, "failed to bind query parameter(name: %q, value: %v), error = <%v>", k, v, err)
	}
	se.decorate(fmt.Sprintf("failed to bind query parameter(name: %q, value: %v)", k, v))
	return se
}
