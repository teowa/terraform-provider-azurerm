// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package elastic

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/elastic/2025-06-01/elasticmonitorresources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/elastic/2025-06-01/tagrules"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/elastic/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceElasticsearch() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceElasticsearchCreate,
		Read:   resourceElasticsearchRead,
		Update: resourceElasticsearchUpdate,
		Delete: resourceElasticsearchDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(60 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(60 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(60 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := elasticmonitorresources.ParseMonitorID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ElasticsearchName,
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"location": commonschema.Location(),

			"sku_name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"elastic_cloud_email_address": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsEmailAddress,
			},

			"monitoring_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
				ForceNew: true,
			},

			"hosting_type": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice(elasticmonitorresources.PossibleValuesForHostingType(), false),
			},

			"project_type": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice(elasticmonitorresources.PossibleValuesForProjectType(), false),
			},

			"project_configuration_type": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice(elasticmonitorresources.PossibleValuesForConfigurationType(), false),
			},

			"logs": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"filtering_tag": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"name": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"value": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"action": {
										Type:     pluginsdk.TypeString,
										Required: true,
										ValidateFunc: validation.StringInSlice([]string{
											string(tagrules.TagActionExclude),
											string(tagrules.TagActionInclude),
										}, false),
									},
								},
							},
						},

						"send_activity_logs": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Default:  false,
						},

						"send_azuread_logs": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Default:  false,
						},

						"send_subscription_logs": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Default:  false,
						},
					},
				},
			},

			"tags": commonschema.Tags(),

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

func resourceElasticsearchCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).Elastic.MonitorClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := elasticmonitorresources.NewMonitorID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.MonitorsGet(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %q: %+v", id, err)
			}
		}
		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_elastic_cloud_elasticsearch", id.ID())
		}
	}

	monitoringStatus := elasticmonitorresources.MonitoringStatusDisabled
	if d.Get("monitoring_enabled").(bool) {
		monitoringStatus = elasticmonitorresources.MonitoringStatusEnabled
	}

	body := elasticmonitorresources.ElasticMonitorResource{
		Location: location.Normalize(d.Get("location").(string)),
		Properties: &elasticmonitorresources.MonitorProperties{
			MonitoringStatus: &monitoringStatus,
			UserInfo: &elasticmonitorresources.UserInfo{
				EmailAddress: pointer.To(d.Get("elastic_cloud_email_address").(string)),
			},
		},
		Sku: &elasticmonitorresources.ResourceSku{
			Name: d.Get("sku_name").(string),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if hostingType := d.Get("hosting_type").(string); hostingType != "" {
		body.Properties.HostingType = pointer.To(elasticmonitorresources.HostingType(hostingType))
	}

	projectType := d.Get("project_type").(string)
	projectConfigurationType := d.Get("project_configuration_type").(string)
	if projectType != "" || projectConfigurationType != "" {
		body.Properties.ProjectDetails = &elasticmonitorresources.ProjectDetails{}
		if projectType != "" {
			body.Properties.ProjectDetails.ProjectType = pointer.To(elasticmonitorresources.ProjectType(projectType))
		}
		if projectConfigurationType != "" {
			body.Properties.ProjectDetails.ConfigurationType = pointer.To(elasticmonitorresources.ConfigurationType(projectConfigurationType))
		}
	}

	if err := client.MonitorsCreateCallbackThenPoll(ctx, id, body, sdk.SetIDCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	if v, ok := d.GetOk("logs"); ok {
		tagRulesClient := meta.(*clients.Client).Elastic.TagRuleClient
		tagRuleId := tagrules.NewTagRuleID(id.SubscriptionId, id.ResourceGroupName, id.MonitorName, "default")
		tagRule := tagrules.MonitoringTagRules{
			Properties: &tagrules.MonitoringTagRulesProperties{
				LogRules: expandTagRule(v.([]interface{})),
			},
		}
		if _, err := tagRulesClient.CreateOrUpdate(ctx, tagRuleId, tagRule); err != nil {
			return fmt.Errorf("updating the logs for %s: %+v", id, err)
		}
	}

	return resourceElasticsearchRead(d, meta)
}

func resourceElasticsearchRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Elastic.MonitorClient
	logsClient := meta.(*clients.Client).Elastic.TagRuleClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := elasticmonitorresources.ParseMonitorID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.MonitorsGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s was not found", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	tagRuleId := tagrules.NewTagRuleID(id.SubscriptionId, id.ResourceGroupName, id.MonitorName, "default")
	rulesResp, err := logsClient.Get(ctx, tagRuleId)
	if err != nil {
		if !response.WasNotFound(rulesResp.HttpResponse) {
			return fmt.Errorf("retrieving logs for %s: %+v", *id, err)
		}
	}

	d.Set("name", id.MonitorName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.Normalize(model.Location))

		if props := model.Properties; props != nil {
			monitoringEnabled := false
			if props.MonitoringStatus != nil {
				monitoringEnabled = *props.MonitoringStatus == elasticmonitorresources.MonitoringStatusEnabled
			}
			d.Set("monitoring_enabled", monitoringEnabled)
			d.Set("hosting_type", pointer.From(props.HostingType))
			if projectDetails := props.ProjectDetails; projectDetails != nil {
				d.Set("project_type", pointer.From(projectDetails.ProjectType))
				d.Set("project_configuration_type", pointer.From(projectDetails.ConfigurationType))
			} else {
				d.Set("project_type", "")
				d.Set("project_configuration_type", "")
			}

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

	return nil
}

func resourceElasticsearchUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := elasticmonitorresources.ParseMonitorID(d.Id())
	if err != nil {
		return err
	}

	if d.HasChange("logs") {
		client := meta.(*clients.Client).Elastic.TagRuleClient
		tagRuleId := tagrules.NewTagRuleID(id.SubscriptionId, id.ResourceGroupName, id.MonitorName, "default")
		tagRule := expandTagRule(d.Get("logs").([]interface{}))
		body := tagrules.MonitoringTagRules{
			Properties: &tagrules.MonitoringTagRulesProperties{
				LogRules: tagRule,
			},
		}
		if _, err := client.CreateOrUpdate(ctx, tagRuleId, body); err != nil {
			return fmt.Errorf("updating `logs` from %s: %+v", *id, err)
		}
	}

	if d.HasChange("tags") {
		client := meta.(*clients.Client).Elastic.MonitorClient
		body := elasticmonitorresources.ElasticMonitorResourceUpdateParameters{
			Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
		}
		if _, err := client.MonitorsUpdate(ctx, *id, body); err != nil {
			return fmt.Errorf("updating %s: %+v", *id, err)
		}
	}

	return resourceElasticsearchRead(d, meta)
}

func resourceElasticsearchDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Elastic.MonitorClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := elasticmonitorresources.ParseMonitorID(d.Id())
	if err != nil {
		return err
	}

	if err := client.MonitorsDeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}

func expandTagRule(input []interface{}) *tagrules.LogRules {
	if len(input) == 0 {
		return nil
	}

	raw := input[0].(map[string]interface{})
	filteringTags := make([]tagrules.FilteringTag, 0)
	for _, v := range raw["filtering_tag"].([]interface{}) {
		item := v.(map[string]interface{})

		action := tagrules.TagAction(item["action"].(string))
		filteringTags = append(filteringTags, tagrules.FilteringTag{
			Action: &action,
			Name:   pointer.To(item["name"].(string)),
			Value:  pointer.To(item["value"].(string)),
		})
	}

	sendAzureAdLogs := raw["send_azuread_logs"].(bool)
	sendActivityLogs := raw["send_activity_logs"].(bool)
	sendSubscriptionLogs := raw["send_subscription_logs"].(bool)

	return &tagrules.LogRules{
		FilteringTags:        &filteringTags,
		SendAadLogs:          pointer.To(sendAzureAdLogs),
		SendActivityLogs:     pointer.To(sendActivityLogs),
		SendSubscriptionLogs: pointer.To(sendSubscriptionLogs),
	}
}

func flattenTagRule(input *tagrules.MonitoringTagRules) []interface{} {
	if input == nil || input.Properties == nil || input.Properties.LogRules == nil {
		return []interface{}{}
	}

	rules := input.Properties.LogRules

	filteringTags := make([]interface{}, 0)
	if rules.FilteringTags != nil {
		for _, v := range *rules.FilteringTags {
			action := ""
			if v.Action != nil {
				action = string(*v.Action)
			}
			name := ""
			if v.Name != nil {
				name = *v.Name
			}
			value := ""
			if v.Value != nil {
				value = *v.Value
			}

			filteringTags = append(filteringTags, map[string]interface{}{
				"action": action,
				"name":   name,
				"value":  value,
			})
		}
	}

	sendActivityLogs := false
	if rules.SendActivityLogs != nil {
		sendActivityLogs = *rules.SendActivityLogs
	}
	sendAzureAdLogs := false
	if rules.SendAadLogs != nil {
		sendAzureAdLogs = *rules.SendAadLogs
	}
	sendSubscriptionLogs := false
	if rules.SendSubscriptionLogs != nil {
		sendSubscriptionLogs = *rules.SendSubscriptionLogs
	}

	return []interface{}{
		map[string]interface{}{
			"filtering_tag":          filteringTags,
			"send_activity_logs":     sendActivityLogs,
			"send_azuread_logs":      sendAzureAdLogs,
			"send_subscription_logs": sendSubscriptionLogs,
		},
	}
}
