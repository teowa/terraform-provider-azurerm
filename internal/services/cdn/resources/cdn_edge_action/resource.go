// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cdn_edge_action

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
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2026-10-01/edgeactions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity

type Resource struct{}

type Model struct {
	Name              string            `tfschema:"name"`
	ResourceGroupName string            `tfschema:"resource_group_name"`
	Location          string            `tfschema:"location"`
	SkuName           string            `tfschema:"sku_name"`
	SkuTier           string            `tfschema:"sku_tier"`
	Tags              map[string]string `tfschema:"tags"`
}

var (
	_ sdk.ResourceWithUpdate   = Resource{}
	_ sdk.ResourceWithIdentity = Resource{}
)

func (Resource) ResourceType() string             { return "azurerm_cdn_edge_action" }
func (Resource) ModelObject() interface{}         { return &Model{} }
func (Resource) Identity() resourceids.ResourceId { return &edgeactions.EdgeActionId{} }
func (Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return edgeactions.ValidateEdgeActionID
}

func (Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.All(validation.StringLenBetween(1, 50), validation.StringMatch(regexp.MustCompile(`^[a-zA-Z0-9]+$`), "must contain only alphanumeric characters")),
		},
		"resource_group_name": commonschema.ResourceGroupName(),
		"location":            commonschema.Location(),
		"sku_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"sku_tier": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
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
			client := metadata.Client.Cdn.EdgeActionsClient
			id := edgeactions.NewEdgeActionID(metadata.Client.Account.SubscriptionId, model.ResourceGroupName, model.Name)
			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("retrieving %s: %+v", id, err)
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}
			payload := edgeactions.EdgeAction{
				Location: location.Normalize(model.Location),
				Sku:      edgeactions.SkuType{Name: model.SkuName, Tier: model.SkuTier},
				Tags:     pointer.To(model.Tags),
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
			id, err := edgeactions.ParseEdgeActionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}
			resp, err := metadata.Client.Cdn.EdgeActionsClient.Get(ctx, *id)
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
			id, err := edgeactions.ParseEdgeActionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}
			var model Model
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}
			payload := edgeactions.EdgeActionUpdate{Tags: pointer.To(model.Tags)}
			if model.Tags == nil {
				payload.Tags = pointer.To(map[string]string{})
			}
			if err := metadata.Client.Cdn.EdgeActionsClient.UpdateThenPoll(ctx, *id, payload); err != nil {
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
			id, err := edgeactions.ParseEdgeActionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}
			if err := metadata.Client.Cdn.EdgeActionsClient.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}
			return nil
		},
	}
}

func (Resource) flatten(metadata sdk.ResourceMetaData, id *edgeactions.EdgeActionId, model *edgeactions.EdgeAction) error {
	if model == nil {
		return fmt.Errorf("retrieving %s: model was nil", id)
	}
	state := Model{
		Name:              id.EdgeActionName,
		ResourceGroupName: id.ResourceGroupName,
		Location:          location.Normalize(model.Location),
		SkuName:           model.Sku.Name,
		SkuTier:           model.Sku.Tier,
		Tags:              pointer.From(model.Tags),
	}
	if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
		return err
	}
	return metadata.Encode(&state)
}
