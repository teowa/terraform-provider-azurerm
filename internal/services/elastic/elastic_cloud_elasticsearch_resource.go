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
	"github.com/hashicorp/go-azure-sdk/resource-manager/elastic/2024-03-01/rules"
	"github.com/hashicorp/go-azure-sdk/resource-manager/elastic/2025-06-01/elasticmonitorresources"
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
											string(rules.TagActionExclude),
											string(rules.TagActionInclude),
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
			"monitor_properties": elasticsearchMonitorPropertiesSchema(),
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
				return fmt.Errorf("checking for existing `%s`: %+v", id, err)
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

	if err := client.MonitorsCreateCallbackThenPoll(ctx, id, body, sdk.SetIDCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	if v, ok := d.GetOk("logs"); ok {
		tagRulesClient := meta.(*clients.Client).Elastic.TagRuleClient
		tagRuleId := rules.NewTagRuleID(id.SubscriptionId, id.ResourceGroupName, id.MonitorName, "default")
		tagRule := rules.MonitoringTagRules{
			Properties: &rules.MonitoringTagRulesProperties{
				LogRules: expandTagRule(v.([]interface{})),
			},
		}
		if _, err := tagRulesClient.TagRulesCreateOrUpdate(ctx, tagRuleId, tagRule); err != nil {
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

	tagRuleId := rules.NewTagRuleID(id.SubscriptionId, id.ResourceGroupName, id.MonitorName, "default")
	rulesResp, err := logsClient.TagRulesGet(ctx, tagRuleId)
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

		if err := d.Set("monitor_properties", flattenMonitorProperties(model.Properties)); err != nil {
			return fmt.Errorf("setting `monitor_properties`: %+v", err)
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
		tagRuleId := rules.NewTagRuleID(id.SubscriptionId, id.ResourceGroupName, id.MonitorName, "default")
		tagRule := expandTagRule(d.Get("logs").([]interface{}))
		body := rules.MonitoringTagRules{
			Properties: &rules.MonitoringTagRulesProperties{
				LogRules: tagRule,
			},
		}
		if _, err := client.TagRulesCreateOrUpdate(ctx, tagRuleId, body); err != nil {
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

func elasticsearchMonitorPropertiesSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"generate_api_key": {
					Type:     pluginsdk.TypeBool,
					Computed: true,
				},
				"hosting_type": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},
				"liftr_resource_category": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},
				"liftr_resource_preference": {
					Type:     pluginsdk.TypeInt,
					Computed: true,
				},
				"monitoring_status": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},
				"plan_details": {
					Type:     pluginsdk.TypeList,
					Computed: true,
					MaxItems: 1,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"offer_id": {
								Type:     pluginsdk.TypeString,
								Computed: true,
							},
							"plan_id": {
								Type:     pluginsdk.TypeString,
								Computed: true,
							},
							"plan_name": {
								Type:     pluginsdk.TypeString,
								Computed: true,
							},
							"publisher_id": {
								Type:     pluginsdk.TypeString,
								Computed: true,
							},
							"term_id": {
								Type:     pluginsdk.TypeString,
								Computed: true,
							},
						},
					},
				},
				"project_details": {
					Type:     pluginsdk.TypeList,
					Computed: true,
					MaxItems: 1,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"configuration_type": {
								Type:     pluginsdk.TypeString,
								Computed: true,
							},
							"project_type": {
								Type:     pluginsdk.TypeString,
								Computed: true,
							},
						},
					},
				},
				"provisioning_state": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},
				"saas_azure_subscription_status": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},
				"source_campaign_id": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},
				"source_campaign_name": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},
				"subscription_state": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},
				"user_info": {
					Type:     pluginsdk.TypeList,
					Computed: true,
					MaxItems: 1,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"company_info": {
								Type:     pluginsdk.TypeList,
								Computed: true,
								MaxItems: 1,
								Elem: &pluginsdk.Resource{
									Schema: map[string]*pluginsdk.Schema{
										"business": {
											Type:     pluginsdk.TypeString,
											Computed: true,
										},
										"country": {
											Type:     pluginsdk.TypeString,
											Computed: true,
										},
										"domain": {
											Type:     pluginsdk.TypeString,
											Computed: true,
										},
										"employees_number": {
											Type:     pluginsdk.TypeString,
											Computed: true,
										},
										"state": {
											Type:     pluginsdk.TypeString,
											Computed: true,
										},
									},
								},
							},
							"company_name": {
								Type:     pluginsdk.TypeString,
								Computed: true,
							},
							"email_address": {
								Type:     pluginsdk.TypeString,
								Computed: true,
							},
							"first_name": {
								Type:     pluginsdk.TypeString,
								Computed: true,
							},
							"last_name": {
								Type:     pluginsdk.TypeString,
								Computed: true,
							},
						},
					},
				},
				"version": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},
			},
		},
	}
}

func flattenMonitorProperties(input *elasticmonitorresources.MonitorProperties) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	generateAPIKey := false
	if input.GenerateApiKey != nil {
		generateAPIKey = *input.GenerateApiKey
	}

	hostingType := ""
	if input.HostingType != nil {
		hostingType = string(*input.HostingType)
	}

	liftrResourceCategory := ""
	if input.LiftrResourceCategory != nil {
		liftrResourceCategory = string(*input.LiftrResourceCategory)
	}

	liftrResourcePreference := 0
	if input.LiftrResourcePreference != nil {
		liftrResourcePreference = int(*input.LiftrResourcePreference)
	}

	monitoringStatus := ""
	if input.MonitoringStatus != nil {
		monitoringStatus = string(*input.MonitoringStatus)
	}

	provisioningState := ""
	if input.ProvisioningState != nil {
		provisioningState = string(*input.ProvisioningState)
	}

	sourceCampaignID := ""
	if input.SourceCampaignId != nil {
		sourceCampaignID = *input.SourceCampaignId
	}

	sourceCampaignName := ""
	if input.SourceCampaignName != nil {
		sourceCampaignName = *input.SourceCampaignName
	}

	saaSAzureSubscriptionStatus := ""
	if input.SaaSAzureSubscriptionStatus != nil {
		saaSAzureSubscriptionStatus = *input.SaaSAzureSubscriptionStatus
	}

	subscriptionState := ""
	if input.SubscriptionState != nil {
		subscriptionState = *input.SubscriptionState
	}

	version := ""
	if input.Version != nil {
		version = *input.Version
	}

	return []interface{}{
		map[string]interface{}{
			"generate_api_key":               generateAPIKey,
			"hosting_type":                   hostingType,
			"liftr_resource_category":        liftrResourceCategory,
			"liftr_resource_preference":      liftrResourcePreference,
			"monitoring_status":              monitoringStatus,
			"plan_details":                   flattenMonitorPlanDetails(input.PlanDetails),
			"project_details":                flattenMonitorProjectDetails(input.ProjectDetails),
			"provisioning_state":             provisioningState,
			"saas_azure_subscription_status": saaSAzureSubscriptionStatus,
			"source_campaign_id":             sourceCampaignID,
			"source_campaign_name":           sourceCampaignName,
			"subscription_state":             subscriptionState,
			"user_info":                      flattenMonitorUserInfo(input.UserInfo),
			"version":                        version,
		},
	}
}

func flattenMonitorPlanDetails(input *elasticmonitorresources.PlanDetails) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	offerID := ""
	if input.OfferID != nil {
		offerID = *input.OfferID
	}

	planID := ""
	if input.PlanID != nil {
		planID = *input.PlanID
	}

	planName := ""
	if input.PlanName != nil {
		planName = *input.PlanName
	}

	publisherID := ""
	if input.PublisherID != nil {
		publisherID = *input.PublisherID
	}

	termID := ""
	if input.TermID != nil {
		termID = *input.TermID
	}

	return []interface{}{
		map[string]interface{}{
			"offer_id":     offerID,
			"plan_id":      planID,
			"plan_name":    planName,
			"publisher_id": publisherID,
			"term_id":      termID,
		},
	}
}

func flattenMonitorProjectDetails(input *elasticmonitorresources.ProjectDetails) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	configurationType := ""
	if input.ConfigurationType != nil {
		configurationType = string(*input.ConfigurationType)
	}

	projectType := ""
	if input.ProjectType != nil {
		projectType = string(*input.ProjectType)
	}

	return []interface{}{
		map[string]interface{}{
			"configuration_type": configurationType,
			"project_type":       projectType,
		},
	}
}

func flattenMonitorUserInfo(input *elasticmonitorresources.UserInfo) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	companyName := ""
	if input.CompanyName != nil {
		companyName = *input.CompanyName
	}

	emailAddress := ""
	if input.EmailAddress != nil {
		emailAddress = *input.EmailAddress
	}

	firstName := ""
	if input.FirstName != nil {
		firstName = *input.FirstName
	}

	lastName := ""
	if input.LastName != nil {
		lastName = *input.LastName
	}

	return []interface{}{
		map[string]interface{}{
			"company_info":  flattenMonitorCompanyInfo(input.CompanyInfo),
			"company_name":  companyName,
			"email_address": emailAddress,
			"first_name":    firstName,
			"last_name":     lastName,
		},
	}
}

func flattenMonitorCompanyInfo(input *elasticmonitorresources.CompanyInfo) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	business := ""
	if input.Business != nil {
		business = *input.Business
	}

	country := ""
	if input.Country != nil {
		country = *input.Country
	}

	domain := ""
	if input.Domain != nil {
		domain = *input.Domain
	}

	employeesNumber := ""
	if input.EmployeesNumber != nil {
		employeesNumber = *input.EmployeesNumber
	}

	state := ""
	if input.State != nil {
		state = *input.State
	}

	return []interface{}{
		map[string]interface{}{
			"business":         business,
			"country":          country,
			"domain":           domain,
			"employees_number": employeesNumber,
			"state":            state,
		},
	}
}

func expandTagRule(input []interface{}) *rules.LogRules {
	if len(input) == 0 {
		return nil
	}

	raw := input[0].(map[string]interface{})
	filteringTags := make([]rules.FilteringTag, 0)
	for _, v := range raw["filtering_tag"].([]interface{}) {
		item := v.(map[string]interface{})

		action := rules.TagAction(item["action"].(string))
		filteringTags = append(filteringTags, rules.FilteringTag{
			Action: &action,
			Name:   pointer.To(item["name"].(string)),
			Value:  pointer.To(item["value"].(string)),
		})
	}

	sendAzureAdLogs := raw["send_azuread_logs"].(bool)
	sendActivityLogs := raw["send_activity_logs"].(bool)
	sendSubscriptionLogs := raw["send_subscription_logs"].(bool)

	return &rules.LogRules{
		FilteringTags:        &filteringTags,
		SendAadLogs:          pointer.To(sendAzureAdLogs),
		SendActivityLogs:     pointer.To(sendActivityLogs),
		SendSubscriptionLogs: pointer.To(sendSubscriptionLogs),
	}
}

func flattenTagRule(input *rules.MonitoringTagRules) []interface{} {
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
