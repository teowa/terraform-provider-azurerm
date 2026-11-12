// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cdn_edge_action_version

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2026-10-01/edgeactionversions"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ListResource struct{}

type ListModel struct {
	CdnEdgeActionID types.String `tfsdk:"cdn_edge_action_id"`
}

var _ sdk.FrameworkListWrappedResource = ListResource{}

func (ListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = Resource{}.ResourceType()
}

func (ListResource) ResourceFunc() *pluginsdk.Resource { return sdk.WrappedResource(Resource{}) }

func (ListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, response *list.ListResourceSchemaResponse) {
	response.Schema = schema.Schema{Attributes: map[string]schema.Attribute{
		"cdn_edge_action_id": schema.StringAttribute{
			Required:   true,
			Validators: []validator.String{typehelpers.WrappedStringValidator{Func: edgeactionversions.ValidateEdgeActionID}},
		},
	}}
}

func (ListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	var config ListModel
	if diags := request.Config.Get(ctx, &config); diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}
	parent, err := edgeactionversions.ParseEdgeActionID(config.CdnEdgeActionID.ValueString())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, "parsing `cdn_edge_action_id`", err)
		return
	}
	resp, err := metadata.Client.Cdn.EdgeActionVersionsClient.ListByEdgeActionComplete(ctx, *parent)
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, "listing CDN Edge Action Versions", err)
		return
	}
	stream.Results = func(push func(list.ListResult) bool) {
		r := Resource{}
		for _, item := range resp.Items {
			result := request.NewListResult(ctx)
			result.DisplayName = pointer.From(item.Name)
			id, err := edgeactionversions.ParseVersionIDInsensitively(pointer.From(item.Id))
			if err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, "parsing Edge Action Version ID", err)
				return
			}
			rmd := sdk.NewResourceMetaData(metadata.Client, r)
			rmd.ResourceData.SetId(id.ID())
			if err := r.flatten(rmd, id, &item); err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("encoding `%s`", r.ResourceType()), err)
				return
			}
			sdk.EncodeListResult(ctx, rmd.ResourceData, &result)
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
