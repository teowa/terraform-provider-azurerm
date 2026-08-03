// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package elastic

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/elastic/2025-06-01/elasticmonitorresources"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ElasticCloudHostedSearchListResource struct{}

var _ sdk.FrameworkListWrappedResource = new(ElasticCloudHostedSearchListResource)

func (r ElasticCloudHostedSearchListResource) ResourceFunc() *pluginsdk.Resource {
	return resourceElasticCloudHostedSearch()
}

func (r ElasticCloudHostedSearchListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = elasticCloudHostedSearchResourceName
}

func (r ElasticCloudHostedSearchListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.Elastic.HostedSearchMonitorClient

	var data sdk.DefaultListModel
	diags := request.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	subscriptionID := metadata.SubscriptionId
	if !data.SubscriptionId.IsNull() {
		subscriptionID = data.SubscriptionId.ValueString()
	}

	results := make([]elasticmonitorresources.ElasticMonitorResource, 0)

	switch {
	case !data.ResourceGroupName.IsNull():
		resp, err := client.MonitorsListByResourceGroupComplete(ctx, commonids.NewResourceGroupID(subscriptionID, data.ResourceGroupName.ValueString()))
		if err != nil {
			sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s` by resource group", elasticCloudHostedSearchResourceName), err)
			return
		}
		results = resp.Items
	default:
		resp, err := client.MonitorsListComplete(ctx, commonids.NewSubscriptionID(subscriptionID))
		if err != nil {
			sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s` by subscription", elasticCloudHostedSearchResourceName), err)
			return
		}
		results = resp.Items
	}

	deadline, ok := ctx.Deadline()
	if !ok {
		sdk.SetResponseErrorDiagnostic(stream, "internal-error", fmt.Errorf("context had no deadline"))
		return
	}

	stream.Results = func(push func(list.ListResult) bool) {
		ctx, cancel := context.WithDeadline(context.Background(), deadline)
		defer cancel()

		for _, item := range results {
			itemID, err := elasticmonitorresources.ParseMonitorIDInsensitively(pointer.From(item.Id))
			if err != nil {
				result := request.NewListResult(ctx)
				sdk.SetErrorDiagnosticAndPushListResult(result, push, "parsing Elastic Cloud Hosted Search ID", err)
				return
			}

			resp, err := client.MonitorsGet(ctx, *itemID)
			if err != nil {
				result := request.NewListResult(ctx)
				sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("retrieving `%s`", elasticCloudHostedSearchResourceName), err)
				return
			}
			if resp.Model == nil || resp.Model.Kind == nil || pointer.From(resp.Model.Kind) != elasticHostedDeploymentKind {
				continue
			}

			if err := elasticCloudHostedSearchValidateMonitor(*itemID, resp.Model); err != nil {
				result := request.NewListResult(ctx)
				sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("validating `%s`", elasticCloudHostedSearchResourceName), err)
				return
			}

			result := request.NewListResult(ctx)
			result.DisplayName = pointer.From(resp.Model.Name)

			rd := resourceElasticCloudHostedSearch().Data(&terraform.InstanceState{})
			rd.SetId(itemID.ID())

			if err := elasticCloudHostedSearchSetResourceData(rd, itemID, resp.Model); err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("encoding `%s` resource data", elasticCloudHostedSearchResourceName), err)
				return
			}

			sdk.EncodeListResult(ctx, rd, &result)
			if result.Diagnostics.HasError() {
				push(result)
				return
			}

			if !push(result) {
				return
			}
		}
	}
}

// End of File
