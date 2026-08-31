// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package helpers

import (
	webapps20250501 "github.com/hashicorp/go-azure-sdk/resource-manager/web/2025-05-01/webapps"
)

// This file contains duplicates of helper functions in common_web_app_schema.go typed against the
// 2025-05-01 webapps SDK package, scoped exclusively to `azurerm_function_app_flex_consumption`.

func ExpandConnectionStringsV20250501(connectionStringsConfig []ConnectionString) *webapps20250501.ConnectionStringDictionary {
	result := &webapps20250501.ConnectionStringDictionary{}
	if len(connectionStringsConfig) == 0 {
		return result
	}

	connectionStrings := make(map[string]webapps20250501.ConnStringValueTypePair)
	for _, v := range connectionStringsConfig {
		connectionStrings[v.Name] = webapps20250501.ConnStringValueTypePair{
			Value: v.Value,
			Type:  webapps20250501.ConnectionStringType(v.Type),
		}
	}
	result.Properties = &connectionStrings

	return result
}

func FlattenConnectionStringsV20250501(appConnectionStrings *webapps20250501.ConnectionStringDictionary) []ConnectionString {
	if appConnectionStrings.Properties == nil || len(*appConnectionStrings.Properties) == 0 {
		return []ConnectionString{}
	}

	connectionStrings := make([]ConnectionString, 0, len(*appConnectionStrings.Properties))
	for k, v := range *appConnectionStrings.Properties {
		connectionString := ConnectionString{
			Name:  k,
			Type:  string(v.Type),
			Value: v.Value,
		}
		connectionStrings = append(connectionStrings, connectionString)
	}

	return connectionStrings
}
