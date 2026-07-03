// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cognitive

import (
	"context"
	"errors"
	"strings"

	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourcegroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2026-03-01/cognitiveservicesaccounts"
	"github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cognitive/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type cognitiveAccountConnectionListModel struct {
	CognitiveAccountName types.String `tfsdk:"cognitive_account_name"`
	ResourceGroupName    types.String `tfsdk:"resource_group_name"`
	SubscriptionId       types.String `tfsdk:"subscription_id"`
}

func cognitiveAccountConnectionListResourceConfigSchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"cognitive_account_name": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{
						Func: validate.AccountName(),
					},
				},
			},

			"resource_group_name": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{
						Func: resourcegroups.ValidateName,
					},
				},
			},

			"subscription_id": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{
						Func: validation.IsUUID,
					},
				},
			},
		},
	}
}

func cognitiveAccountConnectionListAccounts(ctx context.Context, metadata sdk.ResourceMetadata, data cognitiveAccountConnectionListModel) ([]cognitiveservicesaccounts.Account, error) {
	subscriptionID := metadata.SubscriptionId
	if !data.SubscriptionId.IsNull() {
		subscriptionID = data.SubscriptionId.ValueString()
	}

	switch {
	case !data.CognitiveAccountName.IsNull():
		if data.ResourceGroupName.IsNull() {
			return nil, errors.New("`resource_group_name` is required when `cognitive_account_name` is specified")
		}

		id := cognitiveservicesaccounts.NewAccountID(subscriptionID, data.ResourceGroupName.ValueString(), data.CognitiveAccountName.ValueString())
		resp, err := metadata.Client.Cognitive.AccountsClient.AccountsGet(ctx, id)
		if err != nil {
			return nil, err
		}

		if resp.Model == nil {
			return []cognitiveservicesaccounts.Account{}, nil
		}

		return []cognitiveservicesaccounts.Account{*resp.Model}, nil

	case !data.ResourceGroupName.IsNull():
		resp, err := metadata.Client.Cognitive.AccountsClient.AccountsListByResourceGroupComplete(ctx, commonids.NewResourceGroupID(subscriptionID, data.ResourceGroupName.ValueString()))
		if err != nil {
			return nil, err
		}

		return resp.Items, nil

	default:
		resp, err := metadata.Client.Cognitive.AccountsClient.AccountsListComplete(ctx, commonids.NewSubscriptionID(subscriptionID))
		if err != nil {
			return nil, err
		}

		return resp.Items, nil
	}
}

func flattenAccountConnectionMetadata(priorMetadata map[string]string, apiMetadata *map[string]string) map[string]string {
	// Some Connection APIs return additional metadata fields beyond those configured (e.g. `ApiVersion`,
	// `DeploymentApiVersion`). When prior configuration is known (Read/Update) only the configured
	// keys are surfaced to avoid diffs; otherwise (import or list) all API metadata fields are returned.
	apiMetadataValues := pointer.From(apiMetadata)
	if len(priorMetadata) == 0 {
		return apiMetadataValues
	}

	filtered := make(map[string]string)
	for configKey := range priorMetadata {
		for apiKey, apiValue := range apiMetadataValues {
			if strings.EqualFold(configKey, apiKey) {
				filtered[configKey] = apiValue
				break
			}
		}
	}
	return filtered
}
