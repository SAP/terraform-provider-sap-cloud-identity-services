package cli

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
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
		err := json.Unmarshal([]byte(customSchemas), &customSchemasMap)

		if err != nil {
			return false, err
		}

		resBodyMap := res.(map[string]any)

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
