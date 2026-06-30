// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package paloalto

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	components "github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2020-02-02/componentsapis"
	metricsobjectfirewallresources "github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2025-10-08/metricsobjectfirewallresources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type NextGenerationFirewallMetricsResource struct{}

type NextGenerationFirewallMetricsModel struct {
	FirewallId                          string `tfschema:"firewall_id"`
	ApplicationInsightsConnectionString string `tfschema:"application_insights_connection_string"`
	ApplicationInsightsResourceId       string `tfschema:"application_insights_resource_id"`
}

var _ sdk.ResourceWithUpdate = NextGenerationFirewallMetricsResource{}

func (r NextGenerationFirewallMetricsResource) ModelObject() interface{} {
	return &NextGenerationFirewallMetricsModel{}
}

func (r NextGenerationFirewallMetricsResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return metricsobjectfirewallresources.ValidateFirewallID
}

func (r NextGenerationFirewallMetricsResource) ResourceType() string {
	return "azurerm_palo_alto_next_generation_firewall_metrics"
}

func (r NextGenerationFirewallMetricsResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"firewall_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: metricsobjectfirewallresources.ValidateFirewallID,
		},

		"application_insights_connection_string": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"application_insights_resource_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: components.ValidateComponentID,
		},
	}
}

func (r NextGenerationFirewallMetricsResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r NextGenerationFirewallMetricsResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PaloAlto.MetricsObjectFirewallResources

			var model NextGenerationFirewallMetricsModel
			if err := metadata.Decode(&model); err != nil {
				return err
			}

			id, err := metricsobjectfirewallresources.ParseFirewallID(model.FirewallId)
			if err != nil {
				return err
			}

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.MetricsObjectFirewallGet(ctx, *id)
				if err != nil {
					if !response.WasNotFound(existing.HttpResponse) {
						return fmt.Errorf("checking for presence of existing metrics for %s: %+v", *id, err)
					}
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), *id)
				}
			}

			payload := metricsobjectfirewallresources.MetricsObjectFirewallResource{
				Properties: metricsobjectfirewallresources.MetricsObject{
					ApplicationInsightsConnectionString: model.ApplicationInsightsConnectionString,
					ApplicationInsightsResourceId:       model.ApplicationInsightsResourceId,
				},
			}

			if err := client.MetricsObjectFirewallCreateOrUpdateCallbackThenPoll(ctx, *id, payload, metadata.SetIDCallback(id)); err != nil {
				return fmt.Errorf("creating metrics for %s: %+v", *id, err)
			}

			metadata.SetID(*id)
			return nil
		},
	}
}

func (r NextGenerationFirewallMetricsResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PaloAlto.MetricsObjectFirewallResources

			id, err := metricsobjectfirewallresources.ParseFirewallID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var state NextGenerationFirewallMetricsModel

			existing, err := client.MetricsObjectFirewallGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("reading metrics for %s: %+v", *id, err)
			}

			state.FirewallId = id.ID()
			if model := existing.Model; model != nil {
				state.ApplicationInsightsConnectionString = model.Properties.ApplicationInsightsConnectionString
				state.ApplicationInsightsResourceId = model.Properties.ApplicationInsightsResourceId
			}

			return metadata.Encode(&state)
		},
	}
}

func (r NextGenerationFirewallMetricsResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PaloAlto.MetricsObjectFirewallResources

			id, err := metricsobjectfirewallresources.ParseFirewallID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.MetricsObjectFirewallDeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting metrics for %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r NextGenerationFirewallMetricsResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PaloAlto.MetricsObjectFirewallResources

			id, err := metricsobjectfirewallresources.ParseFirewallID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model NextGenerationFirewallMetricsModel
			if err := metadata.Decode(&model); err != nil {
				return err
			}

			payload := metricsobjectfirewallresources.MetricsObjectFirewallResource{
				Properties: metricsobjectfirewallresources.MetricsObject{
					ApplicationInsightsConnectionString: model.ApplicationInsightsConnectionString,
					ApplicationInsightsResourceId:       model.ApplicationInsightsResourceId,
				},
			}

			if err := client.MetricsObjectFirewallCreateOrUpdateThenPoll(ctx, *id, payload); err != nil {
				return fmt.Errorf("updating metrics for %s: %+v", *id, err)
			}

			return nil
		},
	}
}
