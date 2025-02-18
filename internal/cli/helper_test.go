package cli

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testResponseStruct struct {
	Param1 	string	`json:"param1"`
	Param2 	bool	`json:"param2"`
}

func Test_UnmarshalResponse(t *testing.T){

	req := testResponseStruct{
		Param1: "test",
		Param2: true,
	}

	tests := []struct {
		description 				string
		res							interface{}
		retrieveCustomSchemas		bool
		expectError 				bool
	}{
		{
			description : "happy path - no custom schema retrieval",
			res	: map[string]interface{}{
				"param1" : req.Param1,
				"param2" : req.Param2,
			},
			retrieveCustomSchemas : false,
			expectError: false,
		},
		{
			description : "happy path - custom schema retrieval",
			res	: map[string]interface{}{
				"param1" : req.Param1,
				"param2" : req.Param2,
				"customSchemas" : "valid-custom-schema-structure",
			},
			retrieveCustomSchemas : true,
			expectError: false,
		},
		{
			description: "error path - nil response body",
			res : nil,
			retrieveCustomSchemas: false,
			expectError: true,
		},
	}

	for _, test := range tests {
		
		t.Run(test.description, func(t *testing.T) {
			res , cS , err := unMarshalResponse[testResponseStruct](test.res, test.retrieveCustomSchemas)

			if test.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.IsType(t, testResponseStruct{}, res)

				assert.Equal(t, req.Param1, res.Param1)
				assert.Equal(t, req.Param2, res.Param2)
				
				if test.retrieveCustomSchemas{
					assert.NotZero(t, cS)
				} 
			}

		})
	}
}

func Test_GetCustomSchemas(t *testing.T){
	tests := []struct {
		description 				string
		res							interface{}
		containsCustomSchema		bool
		expectError 				bool
	}{
		{
			description: "happy path",
			res: map[string]interface{}{
				"param1" : "test",
				"param2" : false,
				"customSchemas" : "valid-custom-schemas-structure",
			},
			containsCustomSchema: true,
			expectError: false,
		},
		{
			description: "happy path - no custom schemas",
			res: map[string]interface{}{
				"param1" : "test",
				"param2" : false,
			},
			containsCustomSchema: false,
			expectError: false,
		},
	}

	for _, test := range tests{
		t.Run(test.description, func(t *testing.T) {
			cS , err := getCustomSchemas[testResponseStruct](test.res)

			if test.expectError {
				assert.Error(t, err)
			} else if test.containsCustomSchema{
				assert.NoError(t, err)
				assert.NotZero(t, cS)
			} else {
				assert.NoError(t, err)		
				assert.Zero(t, cS)
			}

		})
	}
}

// the following tests compare() as well
func Test_ValidateCustomSchemaResponse(t *testing.T){
	
	customSchemasMarshaled, _ := json.Marshal(map[string]interface{}{
		"schema_id" : map[string]interface{} {
			"schema_attr_1" : 1,
			"schema_attr_2" : false, 
		},
	})
	
	tests := []struct {
		description 				string
		res							interface{}
		customSchemasReq 			string
		expectError 				bool
	}{
		{
			description: "happy path",
			res: map[string]interface{}{
				"param1" : "test",
				"param2" : true,
				"schema_id" : map[string]interface{} {
					"schema_attr_1" : 1,
					"schema_attr_2" : false, 
				},
			},
			customSchemasReq: string(customSchemasMarshaled),
			expectError: false,
		},
		{
			description: "error path - custom schemas request and response mismatch",
			res: map[string]interface{}{
				"param1" : "test",
				"param2" : true,
				"schema_id" : map[string]interface{} {
					"schema_attr_1" : 2,
					"schema_attr_2" : false, 
				},
			},
			customSchemasReq: string(customSchemasMarshaled),
			expectError: true,
		},
		{
			description: "error path - custom schemas request and response mismatch",
			res: map[string]interface{}{
				"param1" : "test",
				"param2" : true,
				"schema_id_mismatch" : map[string]interface{} {
					"schema_attr_1" : 1,
				},
			},
			customSchemasReq: string(customSchemasMarshaled),
			expectError: true,
		},
	}

	for _, test := range tests {
		result , err := validateCustomSchemasResponse(test.res, test.customSchemasReq)

		if test.expectError {
			assert.Error(t, err)
			assert.False(t, result)
		} else {
			assert.NoError(t, err)		
			assert.True(t, result)
		}
	}
}

func Test_CompareAttributes(t *testing.T){
	tests := []struct {
		description 				string
		key 						string
		resMap						map[string]interface{}
		customSchemasMap 			map[string]interface{}
		errMessage					string
	}{
		{
			description: "happy path",
			key: "",
			resMap: map[string]interface{}{
				"schema_attr_1" : "test",
				"schema_attr_2" : false,
				"schema_attr_3" : 12.24,
				"schema_attr_4" : map[string]interface{}{
					"schema_attr_4a" : "test",
					"schema_attr_4b" : true,
				},
			},
			customSchemasMap: map[string]interface{}{
				"schema_attr_1" : "test",
				"schema_attr_2" : false,
				"schema_attr_3" : 12.24,
				"schema_attr_4" : map[string]interface{}{
					"schema_attr_4a" : "test",
					"schema_attr_4b" : true,
				},
			},
			errMessage: "",
		},
		{
			description: "error path - attribute not found in response",
			key: "schema_id",
			resMap: map[string]interface{}{
				"schema_attr_1" : "test",
				"schema_attr_2" : false,
				"schema_attr_3" : 12.24,
			},
			customSchemasMap: map[string]interface{}{
				"schema_attr_1" : "test",
				"schema_attr_2" : false,
				"schema_attr_3" : 12.24,
				"schema_attr_4" : map[string]interface{}{
					"schema_attr_4a" : "test",
					"schema_attr_4b" : true,
				},
			},
			errMessage: "mismatch between response and request for attribute schema_id.schema_attr_4, attribute not found in response",
		},
		{
			description: "error path - mismatch in string attribute",
			key: "schema_id",
			resMap: map[string]interface{}{
				"schema_attr_1" : "test",
			},
			customSchemasMap: map[string]interface{}{
				"schema_attr_1" : "new_test",
			},
			errMessage: "mismatch between response and request in attribute schema_id.schema_attr_1, request sent: \"new_test\" but response received: \"test\"",
		},
		{
			description: "error path - mismatch in float attribute",
			key: "schema_id",
			resMap: map[string]interface{}{
				"schema_attr_1" : 1.24,
			},
			customSchemasMap: map[string]interface{}{
				"schema_attr_1" : 1.23,
			},
			errMessage: "mismatch between response and request in attribute schema_id.schema_attr_1, request sent: \"1.23\" but response received: \"1.24\"",
		},
		{
			description: "error path - mismatch in boolean attribute",
			key: "schema_id",
			resMap: map[string]interface{}{
				"schema_attr_1" : true,
			},
			customSchemasMap: map[string]interface{}{
				"schema_attr_1" : false,
			},
			errMessage: "mismatch between response and request in attribute schema_id.schema_attr_1, request sent: \"false\" but response received: \"true\"",
		},
		{
			description: "error path - mismatch in nested attribute",
			key: "schema_id",
			resMap: map[string]interface{}{
				"schema_attr_1" : map[string]interface{}{
					"schema_attr_1a" : "test",
				},
			},
			customSchemasMap: map[string]interface{}{
				"schema_attr_1" : map[string]interface{}{
					"schema_attr_1a" : "new_test",
				},
			},
			errMessage: "mismatch between response and request in attribute schema_id.schema_attr_1.schema_attr_1a, request sent: \"new_test\" but response received: \"test\"",
		},
	}

	for _, test := range tests{
		result, err := compareAttributes(test.key, test.customSchemasMap, test.resMap)

		if len(test.errMessage) > 0 {
			assert.NotZero(t, err)
			assert.Equal(t, test.errMessage, err)
			assert.False(t, result)
		} else {
			assert.Zero(t, err)
			assert.True(t, result)
		}
	}
}

