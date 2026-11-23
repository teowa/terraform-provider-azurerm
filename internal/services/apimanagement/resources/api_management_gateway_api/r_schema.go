// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package api_management_gateway_api

import (
	"strings"

	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/api"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func apiManagementGatewayApiSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"api_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: api.ValidateApiID,
			DiffSuppressFunc: func(k, old, new string, d *pluginsdk.ResourceData) bool {
				oldApiId, err := api.ParseApiID(old)
				if err != nil {
					return false
				}

				newApiId, err := api.ParseApiID(new)
				if err != nil {
					return false
				}

				return oldApiId.ApiId == stripApiRevision(newApiId.ApiId)
			},
		},
		"gateway_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: azure.ValidateResourceID,
		},
	}
}

// stripApiRevision strips the `;rev=` suffix from an API ID, since the service
// normalizes the ID by removing it during the creation of a gateway API link.
func stripApiRevision(apiId string) string {
	name := apiId
	if strings.Contains(apiId, ";") {
		name = strings.Split(apiId, ";")[0]
	}

	return name
}
