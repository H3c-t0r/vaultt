// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package audit

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// TestOptions_WithFormat exercises WithFormat option to ensure it performs as expected.
func TestOptions_WithFormat(t *testing.T) {
	tests := map[string]struct {
		Value                string
		IsErrorExpected      bool
		ExpectedErrorMessage string
		ExpectedValue        auditFormat
	}{
		"empty": {
			Value:                "",
			IsErrorExpected:      true,
			ExpectedErrorMessage: "format cannot be empty",
		},
		"whitespace": {
			Value:                "     ",
			IsErrorExpected:      true,
			ExpectedErrorMessage: "format cannot be empty",
		},
		"invalid-test": {
			Value:                "test",
			IsErrorExpected:      true,
			ExpectedErrorMessage: "audit.(auditFormat).validate: 'test' is not a valid format: invalid parameter",
		},
		"valid-json": {
			Value:           "json",
			IsErrorExpected: false,
			ExpectedValue:   AuditFormatJSON,
		},
		"valid-jsonx": {
			Value:           "jsonx",
			IsErrorExpected: false,
			ExpectedValue:   AuditFormatJSONx,
		},
	}

	for name, tc := range tests {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			options := &AuditOptions{}
			applyOption := WithFormat(tc.Value)
			err := applyOption(options)
			switch {
			case tc.IsErrorExpected:
				require.Error(t, err)
				require.EqualError(t, err, tc.ExpectedErrorMessage)
			default:
				require.NoError(t, err)
				require.Equal(t, tc.ExpectedValue, options.withFormat)
			}
		})
	}
}

// TestOptions_WithSubtype exercises WithSubtype option to ensure it performs as expected.
func TestOptions_WithSubtype(t *testing.T) {
	tests := map[string]struct {
		Value                string
		IsErrorExpected      bool
		ExpectedErrorMessage string
		ExpectedValue        auditSubtype
	}{
		"empty": {
			Value:                "",
			IsErrorExpected:      true,
			ExpectedErrorMessage: "subtype cannot be empty",
		},
		"whitespace": {
			Value:                "     ",
			IsErrorExpected:      true,
			ExpectedErrorMessage: "subtype cannot be empty",
		},
		"valid": {
			Value:           "AuditResponse",
			IsErrorExpected: false,
			ExpectedValue:   AuditResponseType,
		},
	}

	for name, tc := range tests {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			options := &AuditOptions{}
			applyOption := WithSubtype(tc.Value)
			err := applyOption(options)
			switch {
			case tc.IsErrorExpected:
				require.Error(t, err)
				require.EqualError(t, err, tc.ExpectedErrorMessage)
			default:
				require.NoError(t, err)
				require.Equal(t, tc.ExpectedValue, options.withSubtype)
			}
		})
	}
}

// TestOptions_WithNow exercises WithNow option to ensure it performs as expected.
func TestOptions_WithNow(t *testing.T) {
	tests := map[string]struct {
		Value                time.Time
		IsErrorExpected      bool
		ExpectedErrorMessage string
		ExpectedValue        time.Time
	}{
		"default-time": {
			Value:                time.Time{},
			IsErrorExpected:      true,
			ExpectedErrorMessage: "cannot specify 'now' to be the zero time instant",
		},
		"valid-time": {
			Value:           time.Date(2023, time.July, 4, 12, 3, 0, 0, time.Local),
			IsErrorExpected: false,
			ExpectedValue:   time.Date(2023, time.July, 4, 12, 3, 0, 0, time.Local),
		},
	}

	for name, tc := range tests {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			options := &AuditOptions{}
			applyOption := WithNow(tc.Value)
			err := applyOption(options)
			switch {
			case tc.IsErrorExpected:
				require.Error(t, err)
				require.EqualError(t, err, tc.ExpectedErrorMessage)
			default:
				require.NoError(t, err)
				require.Equal(t, tc.ExpectedValue, options.WithNow)
			}
		})
	}
}

// TestOptions_WithID exercises WithID option to ensure it performs as expected.
func TestOptions_WithID(t *testing.T) {
	tests := map[string]struct {
		Value                string
		IsErrorExpected      bool
		ExpectedErrorMessage string
		ExpectedValue        string
	}{
		"empty": {
			Value:                "",
			IsErrorExpected:      true,
			ExpectedErrorMessage: "id cannot be empty",
		},
		"whitespace": {
			Value:                "     ",
			IsErrorExpected:      true,
			ExpectedErrorMessage: "id cannot be empty",
		},
		"valid": {
			Value:           "test",
			IsErrorExpected: false,
			ExpectedValue:   "test",
		},
	}

	for name, tc := range tests {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			options := &AuditOptions{}
			applyOption := WithID(tc.Value)
			err := applyOption(options)
			switch {
			case tc.IsErrorExpected:
				require.Error(t, err)
				require.EqualError(t, err, tc.ExpectedErrorMessage)
			default:
				require.NoError(t, err)
				require.Equal(t, tc.ExpectedValue, options.WithID)
			}
		})
	}
}

// TestOptions_WithFacility exercises WithFacility option to ensure it performs as expected.
func TestOptions_WithFacility(t *testing.T) {
	tests := map[string]struct {
		Value         string
		ExpectedValue string
	}{
		"empty": {
			Value:         "",
			ExpectedValue: "",
		},
		"whitespace": {
			Value:         "    ",
			ExpectedValue: "",
		},
		"value": {
			Value:         "juan",
			ExpectedValue: "juan",
		},
		"spacey-value": {
			Value:         "   juan   ",
			ExpectedValue: "juan",
		},
	}

	for name, tc := range tests {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			options := &AuditOptions{}
			applyOption := WithFacility(tc.Value)
			err := applyOption(options)
			require.NoError(t, err)
			require.Equal(t, tc.ExpectedValue, options.withFacility)
		})
	}
}

// TestOptions_WithTag exercises WithTag option to ensure it performs as expected.
func TestOptions_WithTag(t *testing.T) {
	tests := map[string]struct {
		Value         string
		ExpectedValue string
	}{
		"empty": {
			Value:         "",
			ExpectedValue: "",
		},
		"whitespace": {
			Value:         "    ",
			ExpectedValue: "",
		},
		"value": {
			Value:         "juan",
			ExpectedValue: "juan",
		},
		"spacey-value": {
			Value:         "   juan   ",
			ExpectedValue: "juan",
		},
	}

	for name, tc := range tests {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			options := &AuditOptions{}
			applyOption := WithTag(tc.Value)
			err := applyOption(options)
			require.NoError(t, err)
			require.Equal(t, tc.ExpectedValue, options.withTag)
		})
	}
}

// TestOptions_WithSocketType exercises WithSocketType option to ensure it performs as expected.
func TestOptions_WithSocketType(t *testing.T) {
	tests := map[string]struct {
		Value         string
		ExpectedValue string
	}{
		"empty": {
			Value:         "",
			ExpectedValue: "",
		},
		"whitespace": {
			Value:         "    ",
			ExpectedValue: "",
		},
		"value": {
			Value:         "juan",
			ExpectedValue: "juan",
		},
		"spacey-value": {
			Value:         "   juan   ",
			ExpectedValue: "juan",
		},
	}

	for name, tc := range tests {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			options := &AuditOptions{}
			applyOption := WithSocketType(tc.Value)
			err := applyOption(options)
			require.NoError(t, err)
			require.Equal(t, tc.ExpectedValue, options.withSocketType)
		})
	}
}

// TestOptions_WithMaxDuration exercises WithMaxDuration option to ensure it performs as expected.
func TestOptions_WithMaxDuration(t *testing.T) {
	tests := map[string]struct {
		Value                string
		ExpectedValue        time.Duration
		IsErrorExpected      bool
		ExpectedErrorMessage string
	}{
		"empty-gives-default": {
			Value: "",
		},
		"whitespace-give-default": {
			Value: "    ",
		},
		"bad-value": {
			Value:                "juan",
			IsErrorExpected:      true,
			ExpectedErrorMessage: "time: invalid duration \"juan\"",
		},
		"bad-spacey-value": {
			Value:                "   juan   ",
			IsErrorExpected:      true,
			ExpectedErrorMessage: "time: invalid duration \"juan\"",
		},
		"duration-2s": {
			Value:         "2s",
			ExpectedValue: 2 * time.Second,
		},
		"duration-2m": {
			Value:         "2m",
			ExpectedValue: 2 * time.Minute,
		},
	}

	for name, tc := range tests {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			options := &AuditOptions{}
			applyOption := WithMaxDuration(tc.Value)
			err := applyOption(options)
			switch {
			case tc.IsErrorExpected:
				require.Error(t, err)
				require.EqualError(t, err, tc.ExpectedErrorMessage)
			default:
				require.NoError(t, err)
				require.Equal(t, tc.ExpectedValue, options.withMaxDuration)
			}
		})
	}
}

// TestOptions_WithFileMode exercises WithFileMode option to ensure it performs as expected.
func TestOptions_WithFileMode(t *testing.T) {
	tests := map[string]struct {
		Value                string
		IsErrorExpected      bool
		ExpectedErrorMessage string
		IsNilExpected        bool
		ExpectedValue        os.FileMode
	}{
		"empty": {
			Value:           "",
			IsErrorExpected: false,
			IsNilExpected:   true,
		},
		"whitespace": {
			Value:           "     ",
			IsErrorExpected: false,
			IsNilExpected:   true,
		},
		"nonsense": {
			Value:                "juan",
			IsErrorExpected:      true,
			ExpectedErrorMessage: "unable to parse file mode: strconv.ParseUint: parsing \"juan\": invalid syntax",
		},
		"zero": {
			Value:           "0000",
			IsErrorExpected: false,
			ExpectedValue:   os.FileMode(0o000),
		},
		"valid": {
			Value:           "0007",
			IsErrorExpected: false,
			ExpectedValue:   os.FileMode(0o007),
		},
	}

	for name, tc := range tests {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			options := &AuditOptions{}
			applyOption := WithFileMode(tc.Value)
			err := applyOption(options)
			switch {
			case tc.IsErrorExpected:
				require.Error(t, err)
				require.EqualError(t, err, tc.ExpectedErrorMessage)
			default:
				require.NoError(t, err)
				switch {
				case tc.IsNilExpected:
					// Optional option 'not supplied' (i.e. was whitespace/empty string)
					require.Nil(t, options.withFileMode)
				default:
					// Dereference the pointer, so we can examine the file mode.
					require.Equal(t, tc.ExpectedValue, *options.withFileMode)
				}
			}
		})
	}
}

// TestOptions_Default exercises getDefaultOptions to assert the default values.
func TestOptions_Default(t *testing.T) {
	opts := getDefaultOptions()
	require.NotNil(t, opts)
	require.True(t, time.Now().After(opts.WithNow))
	require.False(t, opts.WithNow.IsZero())
	require.Equal(t, "AUTH", opts.withFacility)
	require.Equal(t, "vault", opts.withTag)
	require.Equal(t, 2*time.Second, opts.withMaxDuration)
}

// TestOptions_Opts exercises GetOpts with various Option values.
func TestOptions_Opts(t *testing.T) {
	tests := map[string]struct {
		opts                 []AuditOption
		IsErrorExpected      bool
		ExpectedErrorMessage string
		ExpectedID           string
		ExpectedSubtype      auditSubtype
		ExpectedFormat       auditFormat
		IsNowExpected        bool
		ExpectedNow          time.Time
	}{
		"nil-options": {
			opts:            nil,
			IsErrorExpected: false,
			IsNowExpected:   true,
		},
		"empty-options": {
			opts:            []AuditOption{},
			IsErrorExpected: false,
			IsNowExpected:   true,
		},
		"with-multiple-valid-id": {
			opts: []AuditOption{
				WithID("qwerty"),
				WithID("juan"),
			},
			IsErrorExpected: false,
			ExpectedID:      "juan",
			IsNowExpected:   true,
		},
		"with-multiple-valid-subtype": {
			opts: []AuditOption{
				WithSubtype("AuditRequest"),
				WithSubtype("AuditResponse"),
			},
			IsErrorExpected: false,
			ExpectedSubtype: AuditResponseType,
			IsNowExpected:   true,
		},
		"with-multiple-valid-format": {
			opts: []AuditOption{
				WithFormat("json"),
				WithFormat("jsonx"),
			},
			IsErrorExpected: false,
			ExpectedFormat:  AuditFormatJSONx,
			IsNowExpected:   true,
		},
		"with-multiple-valid-now": {
			opts: []AuditOption{
				WithNow(time.Date(2023, time.July, 4, 12, 3, 0, 0, time.Local)),
				WithNow(time.Date(2023, time.July, 4, 13, 3, 0, 0, time.Local)),
			},
			IsErrorExpected: false,
			ExpectedNow:     time.Date(2023, time.July, 4, 13, 3, 0, 0, time.Local),
			IsNowExpected:   false,
		},
		"with-multiple-valid-then-invalid-now": {
			opts: []AuditOption{
				WithNow(time.Date(2023, time.July, 4, 12, 3, 0, 0, time.Local)),
				WithNow(time.Time{}),
			},
			IsErrorExpected:      true,
			ExpectedErrorMessage: "cannot specify 'now' to be the zero time instant",
		},
		"with-multiple-valid-options": {
			opts: []AuditOption{
				WithID("qwerty"),
				WithSubtype("AuditRequest"),
				WithFormat("json"),
				WithNow(time.Date(2023, time.July, 4, 12, 3, 0, 0, time.Local)),
			},
			IsErrorExpected: false,
			ExpectedID:      "qwerty",
			ExpectedSubtype: AuditRequestType,
			ExpectedFormat:  AuditFormatJSON,
			ExpectedNow:     time.Date(2023, time.July, 4, 12, 3, 0, 0, time.Local),
		},
	}

	for name, tc := range tests {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			opts, err := getOpts(tc.opts...)

			switch {
			case tc.IsErrorExpected:
				require.Error(t, err)
				require.EqualError(t, err, tc.ExpectedErrorMessage)
			default:
				require.NotNil(t, opts)
				require.NoError(t, err)
				require.Equal(t, tc.ExpectedID, opts.WithID)
				require.Equal(t, tc.ExpectedSubtype, opts.withSubtype)
				require.Equal(t, tc.ExpectedFormat, opts.withFormat)
				switch {
				case tc.IsNowExpected:
					require.True(t, time.Now().After(opts.WithNow))
					require.False(t, opts.WithNow.IsZero())
				default:
					require.Equal(t, tc.ExpectedNow, opts.WithNow)
				}

			}
		})
	}
}
