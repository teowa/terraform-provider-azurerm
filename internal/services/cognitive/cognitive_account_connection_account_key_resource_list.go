// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cognitive

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2026-03-01/accountconnectionresource"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2026-03-01/cognitiveservicesaccounts"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type CognitiveAccountConnectionAccountKeyListResource struct{}

var _ sdk.FrameworkListWrappedResource = new(CognitiveAccountConnectionAccountKeyListResource)

func (CognitiveAccountConnectionAccountKeyListResource) ResourceFunc() *pluginsdk.Resource {
	return sdk.WrappedResource(CognitiveAccountConnectionAccountKeyResource{})
}

func (CognitiveAccountConnectionAccountKeyListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = CognitiveAccountConnectionAccountKeyResource{}.ResourceType()
}

func (CognitiveAccountConnectionAccountKeyListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, response *list.ListResourceSchemaResponse) {
	response.Schema = cognitiveAccountConnectionListResourceConfigSchema()
}

func (CognitiveAccountConnectionAccountKeyListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.Cognitive.AccountConnectionResourceClient

	var data cognitiveAccountConnectionListModel
	diags := request.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	accountId, err := cognitiveservicesaccounts.ParseAccountID(data.CognitiveAccountId.ValueString())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, "parsing Cognitive Account ID", err)
		return
	}

	connectionsResp, err := client.AccountConnectionsListComplete(ctx, accountconnectionresource.NewAccountID(accountId.SubscriptionId, accountId.ResourceGroupName, accountId.AccountName), accountconnectionresource.DefaultAccountConnectionsListOperationOptions())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", CognitiveAccountConnectionAccountKeyResource{}.ResourceType()), err)
		return
	}

	stream.Results = func(push func(list.ListResult) bool) {
		for _, connection := range connectionsResp.Items {
			if connection.Properties == nil {
				continue
			}

			base := connection.Properties.ConnectionPropertiesV2()
			if string(base.AuthType) != string(accountconnectionresource.ConnectionAuthTypeAccountKey) {
				continue
			}

			connectionId, err := accountconnectionresource.ParseConnectionID(pointer.From(connection.Id))
			if err != nil {
				result := request.NewListResult(ctx)
				sdk.SetErrorDiagnosticAndPushListResult(result, push, "parsing Cognitive Account Connection ID", err)
				return
			}

			result := request.NewListResult(ctx)
			result.DisplayName = pointer.From(connection.Name)

			r := CognitiveAccountConnectionAccountKeyResource{}
			meta := sdk.NewResourceMetaData(metadata.Client, r)
			meta.SetID(connectionId)

			if err := r.flatten(meta, connectionId, &connection, nil, ""); err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("encoding `%s` resource data", r.ResourceType()), err)
				return
			}

			sdk.EncodeListResult(ctx, meta.ResourceData, &result)
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
