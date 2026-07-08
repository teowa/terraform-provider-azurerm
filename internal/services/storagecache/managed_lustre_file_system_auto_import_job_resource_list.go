// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storagecache

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagecache/2025-07-01/autoimportjobs"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ManagedLustreFileSystemAutoImportJobListResource struct{}

type ManagedLustreFileSystemAutoImportJobListModel struct {
	ManagedLustreFileSystemId types.String `tfsdk:"managed_lustre_file_system_id"`
}

var _ sdk.FrameworkListWrappedResource = ManagedLustreFileSystemAutoImportJobListResource{}

func (ManagedLustreFileSystemAutoImportJobListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = ManagedLustreFileSystemAutoImportJobResource{}.ResourceType()
}

func (ManagedLustreFileSystemAutoImportJobListResource) ResourceFunc() *pluginsdk.Resource {
	return sdk.WrappedResource(ManagedLustreFileSystemAutoImportJobResource{})
}

func (ManagedLustreFileSystemAutoImportJobListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, response *list.ListResourceSchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"managed_lustre_file_system_id": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{Func: autoimportjobs.ValidateAmlFilesystemID},
				},
			},
		},
	}
}

func (ManagedLustreFileSystemAutoImportJobListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.StorageCache_2025_07_01.AutoImportJobs

	var data ManagedLustreFileSystemAutoImportJobListModel
	diags := request.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	amlFilesystemId, err := autoimportjobs.ParseAmlFilesystemID(data.ManagedLustreFileSystemId.ValueString())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("parsing Managed Lustre File System ID for `%s`", ManagedLustreFileSystemAutoImportJobResource{}.ResourceType()), err)
		return
	}

	resp, err := client.ListByAmlFilesystemComplete(ctx, *amlFilesystemId)
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", ManagedLustreFileSystemAutoImportJobResource{}.ResourceType()), err)
		return
	}

	resource := ManagedLustreFileSystemAutoImportJobResource{}
	stream.Results = func(push func(list.ListResult) bool) {
		for _, item := range resp.Items {
			result := request.NewListResult(ctx)
			result.DisplayName = pointer.From(item.Name)

			id, err := autoimportjobs.ParseAutoImportJobID(pointer.From(item.Id))
			if err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, "parsing Managed Lustre File System Auto Import Job ID", err)
				return
			}

			meta := sdk.NewResourceMetaData(metadata.Client, resource)
			meta.SetID(id)

			if err := resource.flatten(meta, id, &item); err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("encoding `%s` resource data", resource.ResourceType()), err)
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
