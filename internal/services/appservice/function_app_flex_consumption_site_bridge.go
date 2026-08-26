// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package appservice

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/webapps"
	webapps20250501 "github.com/hashicorp/go-azure-sdk/resource-manager/web/2025-05-01/webapps"
)

// The `azurerm_function_app_flex_consumption` resource is the only App Service resource that requires the
// 2025-05-01 `webapps` API surface, since `FunctionAppConfig.SiteUpdateStrategy` was added exclusively in that
// version. Every other App Service resource in this provider, plus the unrelated `azurerm_logic_app_standard`
// resource, shares the 2023-12-01 `webapps` client, so that client cannot be bumped wholesale. These helpers
// bridge a `webapps.Site` built with the shared (2023-12-01) helpers into the `webapps20250501.Site` shape
// required to call the 2025-05-01 `WebApps` `Get`/`CreateOrUpdate` operations, since the two `Site` object
// graphs are otherwise structurally and JSON-compatible.
func bridgeSiteToV20250501(input webapps.Site) (*webapps20250501.Site, error) {
	data, err := json.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("marshalling `Site`: %+v", err)
	}

	var output webapps20250501.Site
	if err := json.Unmarshal(data, &output); err != nil {
		return nil, fmt.Errorf("unmarshalling `Site` into the 2025-05-01 API shape: %+v", err)
	}

	return &output, nil
}
