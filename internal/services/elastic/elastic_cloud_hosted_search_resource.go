// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package elastic

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/elastic/2025-06-01/elasticmonitorresources"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/elastic/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

//go:generate go run ../../tools/generator-tests resourceidentity

const (
	elasticCloudHostedSearchResourceName = "azurerm_elastic_cloud_hosted_search"
	elasticHostedDeploymentKind          = "elastic-hosted-deployment"
)

func resourceElasticCloudHostedSearch() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceElasticCloudHostedSearchCreate,
		Read:   resourceElasticCloudHostedSearchRead,
		Update: resourceElasticCloudHostedSearchUpdate,
		Delete: resourceElasticCloudHostedSearchDelete,

		Importer: pluginsdk.ImporterValidatingIdentityThen(&elasticmonitorresources.MonitorId{}, resourceElasticCloudHostedSearchImporter),

		Identity: &schema.ResourceIdentity{
			SchemaFunc: pluginsdk.GenerateIdentitySchema(&elasticmonitorresources.MonitorId{}),
		},

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(60 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(60 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(60 * time.Minute),
		},

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

func resourceElasticCloudHostedSearchImporter(ctx context.Context, d *pluginsdk.ResourceData, meta interface{}) ([]*pluginsdk.ResourceData, error) {
	client := meta.(*clients.Client).Elastic.HostedSearchMonitorClient

	id, err := elasticmonitorresources.ParseMonitorID(d.Id())
	if err != nil {
		return nil, err
	}

	resp, err := client.MonitorsGet(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	if err := elasticCloudHostedSearchValidateMonitor(*id, resp.Model); err != nil {
		return nil, err
	}

	d.SetId(id.ID())
	if err := pluginsdk.SetResourceIdentityData(d, id); err != nil {
		return nil, err
	}

	return []*pluginsdk.ResourceData{d}, nil
}

func resourceElasticCloudHostedSearchCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Elastic.HostedSearchMonitorClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := elasticmonitorresources.NewMonitorID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.MonitorsGet(ctx, id)
		if err != nil && !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for existing %s: %+v", id, err)
		}
		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError(elasticCloudHostedSearchResourceName, id.ID())
		}
	}

	monitoringStatus := elasticmonitorresources.MonitoringStatusDisabled
	if d.Get("monitoring_enabled").(bool) {
		monitoringStatus = elasticmonitorresources.MonitoringStatusEnabled
	}

	body := elasticmonitorresources.ElasticMonitorResource{
		Kind:     pointer.To(elasticHostedDeploymentKind),
		Location: location.Normalize(d.Get("location").(string)),
		Properties: &elasticmonitorresources.MonitorProperties{
			HostingType:      pointer.To(elasticmonitorresources.HostingTypeHosted),
			MonitoringStatus: pointer.To(monitoringStatus),
			ProjectDetails: &elasticmonitorresources.ProjectDetails{
				ConfigurationType: pointer.To(elasticmonitorresources.ConfigurationTypeNotApplicable),
				ProjectType:       pointer.To(elasticmonitorresources.ProjectTypeNotApplicable),
			},
			UserInfo: &elasticmonitorresources.UserInfo{
				EmailAddress: pointer.To(d.Get("elastic_cloud_email_address").(string)),
			},
		},
		Sku: &elasticmonitorresources.ResourceSku{
			Name: d.Get("sku_name").(string),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if err := client.MonitorsCreateCallbackThenPoll(ctx, id, body, sdk.SetIDAndIdentityCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	if err := pluginsdk.SetResourceIdentityData(d, &id); err != nil {
		return err
	}

	return resourceElasticCloudHostedSearchRead(d, meta)
}

func resourceElasticCloudHostedSearchRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Elastic.HostedSearchMonitorClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := elasticmonitorresources.ParseMonitorID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.MonitorsGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s was not found - removing from state", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	if err := elasticCloudHostedSearchValidateMonitor(*id, resp.Model); err != nil {
		return err
	}

	if err := elasticCloudHostedSearchSetResourceData(d, id, resp.Model); err != nil {
		return err
	}

	return pluginsdk.SetResourceIdentityData(d, id)
}

func resourceElasticCloudHostedSearchUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Elastic.HostedSearchMonitorClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := elasticmonitorresources.ParseMonitorID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.MonitorsGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("retrieving %s: %+v", *id, err)
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	if err := elasticCloudHostedSearchValidateMonitor(*id, existing.Model); err != nil {
		return err
	}

	if d.HasChange("tags") {
		body := elasticmonitorresources.ElasticMonitorResourceUpdateParameters{
			Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
		}

		if _, err := client.MonitorsUpdate(ctx, *id, body); err != nil {
			return fmt.Errorf("updating %s: %+v", *id, err)
		}
	}

	return resourceElasticCloudHostedSearchRead(d, meta)
}

func resourceElasticCloudHostedSearchDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Elastic.HostedSearchMonitorClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := elasticmonitorresources.ParseMonitorID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.MonitorsGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(existing.HttpResponse) {
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	if err := elasticCloudHostedSearchValidateMonitor(*id, existing.Model); err != nil {
		return err
	}

	if err := client.MonitorsDeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}

func elasticCloudHostedSearchSetResourceData(d *pluginsdk.ResourceData, id *elasticmonitorresources.MonitorId, model *elasticmonitorresources.ElasticMonitorResource) error {
	d.Set("name", id.MonitorName)
	d.Set("resource_group_name", id.ResourceGroupName)

	d.Set("location", location.Normalize(model.Location))

	skuName := ""
	if model.Sku != nil {
		skuName = model.Sku.Name
	}
	d.Set("sku_name", skuName)

	if props := model.Properties; props != nil {
		monitoringEnabled := false
		if props.MonitoringStatus != nil {
			monitoringEnabled = *props.MonitoringStatus == elasticmonitorresources.MonitoringStatusEnabled
		}
		d.Set("monitoring_enabled", monitoringEnabled)

		if elastic := props.ElasticProperties; elastic != nil {
			if elastic.ElasticCloudDeployment != nil {
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

	if err := tags.FlattenAndSet(d, model.Tags); err != nil {
		return err
	}

	return nil
}

func elasticCloudHostedSearchValidateMonitor(id elasticmonitorresources.MonitorId, model *elasticmonitorresources.ElasticMonitorResource) error {
	if model == nil {
		return fmt.Errorf("retrieving %s: model was nil", id)
	}

	kind := pointer.From(model.Kind)
	if kind != elasticHostedDeploymentKind {
		return fmt.Errorf("retrieving %s: expected kind `%s`, got `%s`", id, elasticHostedDeploymentKind, kind)
	}

	if props := model.Properties; props != nil && props.HostingType != nil && *props.HostingType != elasticmonitorresources.HostingTypeHosted {
		return fmt.Errorf("retrieving %s: expected `hosting_type` to be `%s`, got `%s`", id, elasticmonitorresources.HostingTypeHosted, *props.HostingType)
	}

	return nil
}
