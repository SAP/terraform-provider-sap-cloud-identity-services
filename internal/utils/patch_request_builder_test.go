package utils

import (
	"reflect"
	"testing"

	"github.com/SAP/terraform-provider-sap-cloud-identity-services/internal/cli/apiObjects/generic"
	"github.com/stretchr/testify/assert"
)

type testStruct struct {
	FieldOne string `json:"fieldOne"`
	FieldTwo bool   `json:"fieldTwo"`
}

func TestGenerateReplacePatchRequest(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		value    any
		expected generic.PatchRequest
	}{
		{
			name:  "string value",
			path:  "/users/123",
			value: "john",
			expected: generic.PatchRequest{
				Op:    "replace",
				Path:  "/users/123",
				Value: "john",
			},
		},
		{
			name:  "integer value",
			path:  "/age",
			value: 25,
			expected: generic.PatchRequest{
				Op:    "replace",
				Path:  "/age",
				Value: 25,
			},
		},
		{
			name:  "boolean value",
			path:  "/active",
			value: true,
			expected: generic.PatchRequest{
				Op:    "replace",
				Path:  "/active",
				Value: true,
			},
		},
		{
			name:  "nil value",
			path:  "/optional",
			value: nil,
			expected: generic.PatchRequest{
				Op:    "replace",
				Path:  "/optional",
				Value: nil,
			},
		},
		{
			name:  "slice value",
			path:  "/roles",
			value: []string{"admin", "user"},
			expected: generic.PatchRequest{
				Op:    "replace",
				Path:  "/roles",
				Value: []string{"admin", "user"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GenerateReplacePatchRequest(tt.path, tt.value)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGenerateAddPatchRequest(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		value    any
		expected generic.PatchRequest
	}{
		{
			name:  "string value",
			path:  "/emails",
			value: "test@example.com",
			expected: generic.PatchRequest{
				Op:    "add",
				Path:  "/emails",
				Value: "test@example.com",
			},
		},
		{
			name:  "map value",
			path:  "/metadata",
			value: map[string]string{"key": "value"},
			expected: generic.PatchRequest{
				Op:    "add",
				Path:  "/metadata",
				Value: map[string]string{"key": "value"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GenerateAddPatchRequest(tt.path, tt.value)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGenerateDeletePatchRequest(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected generic.PatchRequest
	}{
		{
			name: "delete user field",
			path: "/phoneNumbers/0",
			expected: generic.PatchRequest{
				Op:   "remove",
				Path: "/phoneNumbers/0",
			},
		},
		{
			name: "delete optional field",
			path: "/fax",
			expected: generic.PatchRequest{
				Op:   "remove",
				Path: "/fax",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GenerateDeletePatchRequest(tt.path)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetAttributeTag(t *testing.T) {
	tests := []struct {
		name           string
		attrName       string
		argsType       reflect.Type
		expectedTag    string
		expectedError  bool
		expectedErrMsg string
	}{
		{
			name:          "valid field with json tag",
			attrName:      "FieldOne",
			argsType:      reflect.TypeFor[testStruct](),
			expectedTag:   "fieldOne",
			expectedError: false,
		},
		{
			name:          "valid field with bool",
			attrName:      "FieldTwo",
			argsType:      reflect.TypeFor[testStruct](),
			expectedTag:   "fieldTwo",
			expectedError: false,
		},
		{
			name:           "field not found",
			attrName:       "NonExistent",
			argsType:       reflect.TypeFor[testStruct](),
			expectedTag:    "",
			expectedError:  true,
			expectedErrMsg: "field 'NonExistent' not found in type",
		},
		{
			name:     "field without json tag",
			attrName: "FieldWithNoTag",
			argsType: reflect.TypeOf(struct {
				FieldWithNoTag string
			}{}),
			expectedTag:    "",
			expectedError:  true,
			expectedErrMsg: "field 'FieldWithNoTag' has no json tag",
		},
		{
			name:     "field with omitempty",
			attrName: "OptionalField",
			argsType: reflect.TypeOf(struct {
				OptionalField string `json:"optionalField,omitempty"`
			}{}),
			expectedTag:   "optionalField,omitempty",
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tag, diags := GetAttributeTag(tt.attrName, tt.argsType)

			if tt.expectedError {
				assert.True(t, diags.HasError())
				assert.Equal(t, tt.expectedErrMsg, diags[0].Detail())
				assert.Empty(t, tag)
			} else {
				assert.False(t, diags.HasError())
				assert.Equal(t, tt.expectedTag, tag)
			}
		})
	}
}

func TestGetPatchRequest(t *testing.T) {
	tests := []struct {
		name           string
		attrName       string
		path           string
		value          any
		argsType       reflect.Type
		expected       generic.PatchRequest
		expectedError  bool
		expectedErrMsg string
	}{
		{
			name:     "without prefix path",
			attrName: "FieldOne",
			path:     "",
			value:    "test-value",
			argsType: reflect.TypeFor[testStruct](),
			expected: generic.PatchRequest{Op: "replace", Path: "/fieldOne", Value: "test-value"},
		},
		{
			name:     "with prefix path",
			attrName: "FieldOne",
			path:     "/schemas/id",
			value:    "test-value",
			argsType: reflect.TypeFor[testStruct](),
			expected: generic.PatchRequest{Op: "replace", Path: "//schemas/id/fieldOne", Value: "test-value"},
		},
		{
			name:           "attribute not found",
			attrName:       "NonExistent",
			path:           "",
			value:          "test-value",
			argsType:       reflect.TypeFor[testStruct](),
			expected:       generic.PatchRequest{},
			expectedError:  true,
			expectedErrMsg: "field 'NonExistent' not found in type",
		},
		{
			name:     "with nested path",
			attrName: "FieldTwo",
			path:     "/custom/nested",
			value:    true,
			argsType: reflect.TypeFor[testStruct](),
			expected: generic.PatchRequest{Op: "replace", Path: "//custom/nested/fieldTwo", Value: true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, diags := GetPatchRequest(tt.attrName, tt.path, tt.value, tt.argsType)

			if tt.expectedError {
				assert.True(t, diags.HasError())
				assert.Equal(t, tt.expectedErrMsg, diags[0].Detail())
			} else {
				assert.False(t, diags.HasError())
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestGetScimPatchRequest(t *testing.T) {
	tests := []struct {
		name           string
		attrName       string
		path           string
		value          any
		argsType       reflect.Type
		expected       generic.PatchRequest
		expectedError  bool
		expectedErrMsg string
	}{
		{
			name:     "without scim path",
			attrName: "FieldOne",
			path:     "",
			value:    "test-value",
			argsType: reflect.TypeFor[testStruct](),
			expected: generic.PatchRequest{Op: "replace", Path: "fieldOne", Value: "test-value"},
		},
		{
			name:     "with scim path",
			attrName: "FieldOne",
			path:     "schemas.schema_id",
			value:    "test-value",
			argsType: reflect.TypeFor[testStruct](),
			expected: generic.PatchRequest{Op: "replace", Path: "schemas.schema_id:fieldOne", Value: "test-value"},
		},
		{
			name:           "attribute not found",
			attrName:       "NonExistent",
			path:           "",
			value:          "test-value",
			argsType:       reflect.TypeFor[testStruct](),
			expected:       generic.PatchRequest{},
			expectedError:  true,
			expectedErrMsg: "field 'NonExistent' not found in type",
		},
		{
			name:     "with omitempty tag",
			attrName: "OptionalField",
			path:     "",
			value:    "value",
			argsType: reflect.TypeOf(struct {
				OptionalField string `json:"optionalField,omitempty"`
			}{}),
			expected: generic.PatchRequest{Op: "replace", Path: "optionalField", Value: "value"},
		},
		{
			name:     "complex scim path",
			attrName: "FieldTwo",
			path:     "emails.home",
			value:    true,
			argsType: reflect.TypeFor[testStruct](),
			expected: generic.PatchRequest{Op: "replace", Path: "emails.home:fieldTwo", Value: true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, diags := GetScimPatchRequest(tt.attrName, tt.path, tt.value, tt.argsType)

			if tt.expectedError {
				assert.True(t, diags.HasError())
				assert.Equal(t, tt.expectedErrMsg, diags[0].Detail())
			} else {
				assert.False(t, diags.HasError())
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}
