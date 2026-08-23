// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package elastic

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/elastic/2025-06-01/elasticmonitorresources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/elastic/2025-06-01/tagrules"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/elastic/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourceElasticsearch() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceElasticsearchRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			// Required
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.ElasticsearchName,
			},

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			// Computed
			"location": commonschema.LocationComputed(),

			"sku_name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"elastic_cloud_email_address": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"monitoring_enabled": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"kind": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"hosting_type": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"generate_api_key_enabled": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"identity": commonschema.SystemAssignedIdentityComputed(),

			"logs": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"filtering_tag": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"name": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},
									"value": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},
									"action": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},
								},
							},
						},

						"send_azuread_logs": {
							Type:     pluginsdk.TypeBool,
							Computed: true,
						},

						"send_activity_logs": {
							Type:     pluginsdk.TypeBool,
							Computed: true,
						},

						"send_subscription_logs": {
							Type:     pluginsdk.TypeBool,
							Computed: true,
						},
					},
				},
			},

			"tags": commonschema.TagsDataSource(),

			"elastic_cloud_deployment_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
			"elastic_cloud_sso_default_url": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
			"elastic_cloud_user_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
			"elasticsearch_service_url": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
			"kibana_service_url": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
			"kibana_sso_uri": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceElasticsearchRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Elastic.MonitorClient
	logsClient := meta.(*clients.Client).Elastic.TagRuleClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := elasticmonitorresources.NewMonitorID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	resp, err := client.MonitorsGet(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	tagRuleId := tagrules.NewTagRuleID(id.SubscriptionId, id.ResourceGroupName, id.MonitorName, "default")
	rulesResp, err := logsClient.Get(ctx, tagRuleId)
	if err != nil {
		if !response.WasNotFound(rulesResp.HttpResponse) {
			return fmt.Errorf("retrieving logs for %s: %+v", id, err)
		}
	}

	d.Set("name", id.MonitorName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.Normalize(model.Location))
		d.Set("kind", pointer.From(model.Kind))

		if err := d.Set("identity", identity.FlattenSystemAssigned(model.Identity)); err != nil {
			return fmt.Errorf("setting `identity`: %+v", err)
		}

		if props := model.Properties; props != nil {
			monitoringEnabled := false
			if props.MonitoringStatus != nil {
				monitoringEnabled = *props.MonitoringStatus == elasticmonitorresources.MonitoringStatusEnabled
			}
			d.Set("monitoring_enabled", monitoringEnabled)
			d.Set("hosting_type", pointer.From(props.HostingType))
			d.Set("generate_api_key_enabled", pointer.From(props.GenerateApiKey))

			if elastic := props.ElasticProperties; elastic != nil {
				if elastic.ElasticCloudDeployment != nil {
					// AzureSubscriptionId is the same as the subscription deployed into, so no point exposing it
					// ElasticsearchRegion is `{Cloud}-{Region}` - so the same as location/not worth exposing for now?
					d.Set("elastic_cloud_deployment_id", elastic.ElasticCloudDeployment.DeploymentId)
					d.Set("elasticsearch_service_url", elastic.ElasticCloudDeployment.ElasticsearchServiceURL)
					d.Set("kibana_service_url", elastic.ElasticCloudDeployment.KibanaServiceURL)
					d.Set("kibana_sso_uri", elastic.ElasticCloudDeployment.KibanaSsoURL)
				}
				if elastic.ElasticCloudUser != nil {
					d.Set("elastic_cloud_user_id", elastic.ElasticCloudUser.Id)
					d.Set("elastic_cloud_email_address", elastic.ElasticCloudUser.EmailAddress)
					d.Set("elastic_cloud_sso_default_url", elastic.ElasticCloudUser.ElasticCloudSsoDefaultURL)
				}
			}
		}

		skuName := ""
		if model.Sku != nil {
			skuName = model.Sku.Name
		}
		d.Set("sku_name", skuName)

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	if err := d.Set("logs", flattenTagRule(rulesResp.Model)); err != nil {
		return fmt.Errorf("setting `logs`: %+v", err)
	}

	d.SetId(id.ID())

	return nil
}
