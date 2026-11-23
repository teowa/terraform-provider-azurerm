// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package api_management_api_tag_description

import (
	"strings"

	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/apitag"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func apiManagementApiTagDescriptionSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"api_tag_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: apitag.ValidateApiTagID,
			DiffSuppressFunc: func(k, old, new string, d *pluginsdk.ResourceData) bool {
				oldApiTagId, err := apitag.ParseApiTagID(old)
				if err != nil {
					return false
				}

				newApiTagId, err := apitag.ParseApiTagID(new)
				if err != nil {
					return false
				}

				return oldApiTagId.ApiId == stripApiRevision(newApiTagId.ApiId)
			},
		},

		"description": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"external_documentation_url": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.IsURLWithHTTPorHTTPS,
		},

		"external_documentation_description": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},
	}
}

// stripApiRevision strips the `;rev=` suffix from an API ID, since the service
// normalizes the ID by removing it during the creation of a tag description.
func stripApiRevision(apiId string) string {
	name := apiId
	if strings.Contains(apiId, ";") {
		name = strings.Split(apiId, ";")[0]
	}

	return name
}
