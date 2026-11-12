// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cdn_edge_action_execution_filter

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2026-10-01/edgeactionexecutionfilters"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2026-10-01/edgeactionversions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity -resource-name cdn_edge_action_execution_filter -properties "name" -compare-values "subscription_id:cdn_edge_action_id,resource_group_name:cdn_edge_action_id,edge_action_name:cdn_edge_action_id"

type Resource struct{}

type Model struct {
	Name                   string            `tfschema:"name"`
	CdnEdgeActionID        string            `tfschema:"cdn_edge_action_id"`
	Location               string            `tfschema:"location"`
	CdnEdgeActionVersionID string            `tfschema:"cdn_edge_action_version_id"`
	HeaderName             string            `tfschema:"header_name"`
	HeaderValue            string            `tfschema:"header_value"`
	Tags                   map[string]string `tfschema:"tags"`
}

var (
	_ sdk.ResourceWithUpdate   = Resource{}
	_ sdk.ResourceWithIdentity = Resource{}
)

func (Resource) ResourceType() string     { return "azurerm_cdn_edge_action_execution_filter" }
func (Resource) ModelObject() interface{} { return &Model{} }
func (Resource) Identity() resourceids.ResourceId {
	return &edgeactionexecutionfilters.ExecutionFilterId{}
}

func (Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return edgeactionexecutionfilters.ValidateExecutionFilterID
}

func (Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.All(validation.StringLenBetween(1, 50), validation.StringMatch(regexp.MustCompile(`^[a-zA-Z0-9]+$`), "must contain only alphanumeric characters")),
		},
		"cdn_edge_action_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: edgeactionexecutionfilters.ValidateEdgeActionID,
		},
		"location": commonschema.Location(),
		"cdn_edge_action_version_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: edgeactionversions.ValidateVersionID,
		},
		"header_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"header_value": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"tags": commonschema.Tags(),
	}
}

func (Resource) Attributes() map[string]*pluginsdk.Schema { return map[string]*pluginsdk.Schema{} }

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model Model
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}
			parent, err := edgeactionexecutionfilters.ParseEdgeActionID(model.CdnEdgeActionID)
			if err != nil {
				return err
			}
			id := edgeactionexecutionfilters.NewExecutionFilterID(parent.SubscriptionId, parent.ResourceGroupName, parent.EdgeActionName, model.Name)
			client := metadata.Client.Cdn.EdgeActionExecutionFiltersClient
			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("retrieving %s: %+v", id, err)
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}
			payload := edgeactionexecutionfilters.EdgeActionExecutionFilter{
				Location: location.Normalize(model.Location),
				Properties: &edgeactionexecutionfilters.EdgeActionExecutionFilterProperties{
					VersionId:                            model.CdnEdgeActionVersionID,
					ExecutionFilterIdentifierHeaderName:  model.HeaderName,
					ExecutionFilterIdentifierHeaderValue: model.HeaderValue,
				},
				Tags: pointer.To(model.Tags),
			}
			if err := client.CreateCallbackThenPoll(ctx, id, payload, metadata.SetIDAndIdentityCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}
			metadata.SetID(id)
			return pluginsdk.SetResourceIdentityData(metadata.ResourceData, &id)
		},
	}
}

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := edgeactionexecutionfilters.ParseExecutionFilterID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}
			resp, err := metadata.Client.Cdn.EdgeActionExecutionFiltersClient.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}
			return r.flatten(metadata, id, resp.Model)
		},
	}
}

func (Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := edgeactionexecutionfilters.ParseExecutionFilterID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}
			var model Model
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}
			updatedTags := model.Tags
			if updatedTags == nil {
				updatedTags = map[string]string{}
			}
			payload := edgeactionexecutionfilters.EdgeActionExecutionFilterUpdate{
				Properties: &edgeactionexecutionfilters.EdgeActionExecutionFilterUpdateProperties{
					VersionId:                            pointer.To(model.CdnEdgeActionVersionID),
					ExecutionFilterIdentifierHeaderName:  pointer.To(model.HeaderName),
					ExecutionFilterIdentifierHeaderValue: pointer.To(model.HeaderValue),
				},
				Tags: pointer.To(updatedTags),
			}
			if err := metadata.Client.Cdn.EdgeActionExecutionFiltersClient.UpdateThenPoll(ctx, *id, payload); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}
			return nil
		},
	}
}

func (Resource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := edgeactionexecutionfilters.ParseExecutionFilterID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}
			if err := metadata.Client.Cdn.EdgeActionExecutionFiltersClient.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}
			return nil
		},
	}
}

func (Resource) flatten(metadata sdk.ResourceMetaData, id *edgeactionexecutionfilters.ExecutionFilterId, model *edgeactionexecutionfilters.EdgeActionExecutionFilter) error {
	if model == nil || model.Properties == nil {
		return fmt.Errorf("retrieving %s: model or properties was nil", id)
	}
	version, err := edgeactionversions.ParseVersionIDInsensitively(model.Properties.VersionId)
	if err != nil {
		return err
	}
	state := Model{
		Name:                   id.ExecutionFilterName,
		CdnEdgeActionID:        edgeactionexecutionfilters.NewEdgeActionID(id.SubscriptionId, id.ResourceGroupName, id.EdgeActionName).ID(),
		Location:               location.Normalize(model.Location),
		CdnEdgeActionVersionID: version.ID(),
		HeaderName:             model.Properties.ExecutionFilterIdentifierHeaderName,
		HeaderValue:            model.Properties.ExecutionFilterIdentifierHeaderValue,
		Tags:                   pointer.From(model.Tags),
	}
	if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
		return err
	}
	return metadata.Encode(&state)
}
