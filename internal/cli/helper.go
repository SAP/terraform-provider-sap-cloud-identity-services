package cli

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/SAP/terraform-provider-sap-cloud-identity-services/internal/cli/apiObjects/generic"
)

func unMarshalResponse[I any](res any, retrieveCustomSchemas bool) (I, string, error) {
	var obj I
	if res == nil {
		return obj, "", fmt.Errorf("response is nil")
	}

	marshaledRes, err := json.Marshal(res)
	if err != nil {
		return obj, "", err
	}

	if err := json.Unmarshal(marshaledRes, &obj); err != nil {
		return obj, "", err
	}

	if !retrieveCustomSchemas {
		return obj, "", nil
	}

	customSchemaResponse, err := getCustomSchemas[I](res)
	return obj, customSchemaResponse, err
}

func getCustomSchemas[I any](res any) (string, error) {

	var customSchemas string
	var obj I

	reflectType := reflect.TypeOf(obj)
	resMap := res.(map[string]any)

	for field := range reflectType.Fields() {
		key := strings.Split(field.Tag.Get("json"), ",")[0]
		delete(resMap, key)
	}

	if len(resMap) > 0 {
		if masrhaledRes, err := json.Marshal(resMap); err == nil {
			customSchemas = string(masrhaledRes)
		} else {
			return customSchemas, err
		}
	}

	return customSchemas, nil
}

func validateCustomSchemasResponse(res any, customSchemas string) (bool, error) {

	var resBody string
	if marshaledRes, err := json.Marshal(res); err == nil {
		resBody = string(marshaledRes)
	} else {
		return false, err
	}

	//remove the beginning and ending characters from the custom schemas string
	modifiedCustomSchemas := customSchemas[1:][:len(customSchemas[1:])-1]

	//check if the custom schemas passed as a request is a substring in the response body
	if !strings.Contains(resBody, modifiedCustomSchemas) {

		var customSchemasMap map[string]any
		if err := json.Unmarshal([]byte(customSchemas), &customSchemasMap); err != nil {
			return false, err
		}

		var resBodyMap map[string]any
		if err := json.Unmarshal([]byte(resBody), &resBodyMap); err != nil {
			return false, err
		}

		// if not a substring, compare the request and response
		return compare(customSchemasMap, resBodyMap)
	}

	return true, nil
}

func compare(cS map[string]any, rB map[string]any) (bool, error) {

	for k, csValue := range cS {

		rbValue, ok := rB[k]
		if !ok {
			err := fmt.Errorf("%s not found in the returned response", k)
			return false, err
		}

		result, err := compareAttributes(k, csValue.(map[string]any), rbValue.(map[string]any))

		if !result {
			return false, fmt.Errorf("%s", err)
		}
	}

	return true, nil
}

func compareAttributes(key string, csValue map[string]any, rbValue map[string]any) (bool, string) {
	for ckey, cval := range csValue {

		rval, ok := rbValue[ckey]

		if !ok {
			err := fmt.Sprintf("mismatch between response and request for attribute %s.%s, attribute not found in response", key, ckey)
			return false, err
		}

		var result bool
		var err string

		// parse through each attribute and retrieve the respective value according to the data type of the attribute
		switch rval.(type) {

		case string:
			rRes := rval.(string)
			cRes := cval.(string)
			if result = (rRes == cRes); !result {
				err = fmt.Sprintf("mismatch between response and request in attribute %s.%s, request sent: \"%s\" but response received: \"%s\"", key, ckey, cRes, rRes)
				break
			}

		case float64:
			rRes := rval.(float64)
			cRes := cval.(float64)
			if result = (rRes == cRes); !result {
				err = fmt.Sprintf("mismatch between response and request in attribute %s.%s, request sent: \"%.2f\" but response received: \"%.2f\"", key, ckey, cRes, rRes)
				break
			}

		case bool:
			rRes := rval.(bool)
			cRes := cval.(bool)
			if result = (rRes == cRes); !result {
				err = fmt.Sprintf("mismatch between response and request in attribute %s.%s, request sent: \"%t\" but response received: \"%t\"", key, ckey, cRes, rRes)
				break
			}

		// for nested structures, call the function recursively
		// API allows only one level of nesting
		case map[string]any:
			rRes := rval.(map[string]any)
			cRes := cval.(map[string]any)
			result, err = compareAttributes(ckey, cRes, rRes)
			if !result {
				err = err[:51] + key + "." + err[51:]
				break
			}

			// TODO handling reference datatcRespe
		}

		if !result {
			return false, err
		}
	}

	return true, ""
}

// writeOnlyPaths contains PATCH operation paths whose values are never returned
// by the GET response (e.g. secrets), and must be skipped during polling.
var writeOnlyPaths = map[string]bool{
	"/oidcConfiguration/clientSecret": true,
}

// patchOpsReflected checks whether all replace/add operations are visible in the GET response.
func patchOpsReflected(ops []generic.PatchRequest, resMap map[string]any) bool {
	for _, op := range ops {
		if op.Op == "remove" || writeOnlyPaths[op.Path] {
			continue
		}
		if !pathMatchesValue(op.Path, op.Value, resMap) {
			return false
		}
	}
	return true
}

func pathMatchesValue(path string, expected any, resMap map[string]any) bool {
	segments := strings.Split(strings.TrimPrefix(path, "/"), "/")

	current := resMap
	for i, seg := range segments {
		if i == len(segments)-1 {
			expBytes, err := json.Marshal(expected)
			if err != nil {
				return false
			}
			var normalizedExpected any
			if err := json.Unmarshal(expBytes, &normalizedExpected); err != nil {
				return false
			}

			actual, ok := current[seg]
			if !ok {
				// Field absent from response: omitempty suppressed it because it is a zero
				// value. Treat as a match only when the expected value is also a zero value.
				return isJSONZeroValue(normalizedExpected)
			}
			return valuesMatch(normalizedExpected, actual)
		}
		nested, ok := current[seg]
		if !ok {
			return false
		}
		nestedMap, ok := nested.(map[string]any)
		if !ok {
			return false
		}
		current = nestedMap
	}
	return true
}

// isJSONZeroValue reports whether v is a zero value that omitempty would suppress:
// false for bools, 0 for numbers, "" for strings, nil, empty slices, and empty objects.
func isJSONZeroValue(v any) bool {
	switch val := v.(type) {
	case bool:
		return !val
	case float64:
		return val == 0
	case string:
		return val == ""
	case nil:
		return true
	case []any:
		return len(val) == 0
	case map[string]any:
		return len(val) == 0
	default:
		return false
	}
}

// valuesMatch compares two values recursively.
// Arrays are compared as unordered sets: every expected element must have a match in actual.
// Objects are compared as subsets: every expected key must exist in actual with a matching value,
// tolerating extra fields added by the server.
// All other values are compared via their JSON representations.
func valuesMatch(expected, actual any) bool {
	switch e := expected.(type) {
	case []any:
		a, ok := actual.([]any)
		if !ok || len(e) != len(a) {
			return false
		}
		matched := make([]bool, len(a))
		for _, expElem := range e {
			found := false
			for j, actElem := range a {
				if !matched[j] && valuesMatch(expElem, actElem) {
					matched[j] = true
					found = true
					break
				}
			}
			if !found {
				return false
			}
		}
		return true

	case map[string]any:
		a, ok := actual.(map[string]any)
		if !ok {
			return false
		}
		for k, expVal := range e {
			actVal, ok := a[k]
			if !ok || !valuesMatch(expVal, actVal) {
				return false
			}
		}
		return true

	default:
		expBytes, _ := json.Marshal(expected)
		actBytes, _ := json.Marshal(actual)
		return string(expBytes) == string(actBytes)
	}
}
