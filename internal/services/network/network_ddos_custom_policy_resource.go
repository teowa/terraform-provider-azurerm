// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-05-01/ddoscustompolicies"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

//go:generate go run ../../tools/generator-tests resourceidentity -resource-name network_ddos_custom_policy -service-package-name network -test-name basicConfigIdentity

func resourceNetworkDDoSCustomPolicy() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceNetworkDDoSCustomPolicyCreate,
		Read:   resourceNetworkDDoSCustomPolicyRead,
		Update: resourceNetworkDDoSCustomPolicyUpdate,
		Delete: resourceNetworkDDoSCustomPolicyDelete,

		Importer: pluginsdk.ImporterValidatingIdentity(&ddoscustompolicies.DdosCustomPolicyId{}),

		Identity: &schema.ResourceIdentity{
			SchemaFunc: pluginsdk.GenerateIdentitySchema(&ddoscustompolicies.DdosCustomPolicyId{}),
		},

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"location": commonschema.Location(),

			"resource_group_name": commonschema.ResourceGroupName(),

			"detection_rule": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"detection_mode": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice(ddoscustompolicies.PossibleValuesForDdosDetectionMode(), false),
						},

						"traffic_detection_rule": {
							Type:     pluginsdk.TypeList,
							Required: true,
							MaxItems: 1,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"packets_per_second": {
										Type:         pluginsdk.TypeInt,
										Required:     true,
										ValidateFunc: validation.IntAtLeast(1),
									},

									"traffic_type": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: validation.StringInSlice(ddoscustompolicies.PossibleValuesForDdosTrafficType(), false),
									},
								},
							},
						},
					},
				},
			},

			"frontend_ip_configuration_ids": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"resource_guid": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"tags": commonschema.Tags(),
		},
	}
}

func resourceNetworkDDoSCustomPolicyCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.DdosCustomPoliciesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := ddoscustompolicies.NewDdosCustomPolicyID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_network_ddos_custom_policy", id.ID())
		}
	}

	payload := ddoscustompolicies.DdosCustomPolicy{
		Location: pointer.To(location.Normalize(d.Get("location").(string))),
		Properties: &ddoscustompolicies.DdosCustomPolicyPropertiesFormat{
			DetectionRules: expandNetworkDDoSCustomPolicyDetectionRules(d.Get("detection_rule").([]interface{})),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, payload, sdk.SetIDAndIdentityCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	if err := pluginsdk.SetResourceIdentityData(d, &id); err != nil {
		return err
	}

	return resourceNetworkDDoSCustomPolicyRead(d, meta)
}

func resourceNetworkDDoSCustomPolicyUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.DdosCustomPoliciesClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := ddoscustompolicies.ParseDdosCustomPolicyID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.Get(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", id)
	}

	payload := existing.Model
	if payload.Properties == nil {
		payload.Properties = &ddoscustompolicies.DdosCustomPolicyPropertiesFormat{}
	}

	if d.HasChange("detection_rule") {
		payload.Properties.DetectionRules = expandNetworkDDoSCustomPolicyDetectionRules(d.Get("detection_rule").([]interface{}))
	}

	if d.HasChange("tags") {
		payload.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	if err := client.CreateOrUpdateThenPoll(ctx, *id, *payload); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourceNetworkDDoSCustomPolicyRead(d, meta)
}

func resourceNetworkDDoSCustomPolicyRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.DdosCustomPoliciesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := ddoscustompolicies.ParseDdosCustomPolicyID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return resourceNetworkDDoSCustomPolicyFlatten(d, id, resp.Model)
}

func resourceNetworkDDoSCustomPolicyFlatten(d *pluginsdk.ResourceData, id *ddoscustompolicies.DdosCustomPolicyId, model *ddoscustompolicies.DdosCustomPolicy) error {
	d.Set("name", id.DdosCustomPolicyName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model != nil {
		d.Set("location", location.NormalizeNilable(model.Location))

		if props := model.Properties; props != nil {
			if err := d.Set("detection_rule", flattenNetworkDDoSCustomPolicyDetectionRules(props.DetectionRules)); err != nil {
				return fmt.Errorf("setting `detection_rule`: %+v", err)
			}
			if err := d.Set("frontend_ip_configuration_ids", flattenNetworkDDoSCustomPolicySubResourceIds(props.FrontEndIPConfiguration)); err != nil {
				return fmt.Errorf("setting `frontend_ip_configuration_ids`: %+v", err)
			}
			d.Set("resource_guid", pointer.From(props.ResourceGuid))
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	return pluginsdk.SetResourceIdentityData(d, id)
}

func resourceNetworkDDoSCustomPolicyDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.DdosCustomPoliciesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := ddoscustompolicies.ParseDdosCustomPolicyID(d.Id())
	if err != nil {
		return err
	}

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}

func expandNetworkDDoSCustomPolicyDetectionRules(input []interface{}) *[]ddoscustompolicies.DdosDetectionRule {
	if len(input) == 0 {
		return nil
	}

	output := make([]ddoscustompolicies.DdosDetectionRule, 0, len(input))

	for _, item := range input {
		data := item.(map[string]interface{})
		rule := ddoscustompolicies.DdosDetectionRule{
			Name: pointer.To(data["name"].(string)),
			Properties: &ddoscustompolicies.DdosDetectionRulePropertiesFormat{
				DetectionMode: pointer.To(ddoscustompolicies.DdosDetectionMode(data["detection_mode"].(string))),
			},
		}

		if trafficDetectionRules := data["traffic_detection_rule"].([]interface{}); len(trafficDetectionRules) > 0 {
			trafficDetectionRule := trafficDetectionRules[0].(map[string]interface{})
			rule.Properties.TrafficDetectionRule = &ddoscustompolicies.TrafficDetectionRule{
				PacketsPerSecond: pointer.To(int64(trafficDetectionRule["packets_per_second"].(int))),
				TrafficType:      pointer.To(ddoscustompolicies.DdosTrafficType(trafficDetectionRule["traffic_type"].(string))),
			}
		}

		output = append(output, rule)
	}

	return &output
}

func flattenNetworkDDoSCustomPolicyDetectionRules(input *[]ddoscustompolicies.DdosDetectionRule) []interface{} {
	output := make([]interface{}, 0)
	if input == nil {
		return output
	}

	for _, rule := range *input {
		trafficDetectionRules := make([]interface{}, 0)
		detectionMode := ""

		if props := rule.Properties; props != nil {
			detectionMode = string(pointer.From(props.DetectionMode))

			if trafficDetectionRule := props.TrafficDetectionRule; trafficDetectionRule != nil {
				trafficDetectionRules = append(trafficDetectionRules, map[string]interface{}{
					"packets_per_second": int(pointer.From(trafficDetectionRule.PacketsPerSecond)),
					"traffic_type":       string(pointer.From(trafficDetectionRule.TrafficType)),
				})
			}
		}

		output = append(output, map[string]interface{}{
			"detection_mode":         detectionMode,
			"name":                   pointer.From(rule.Name),
			"traffic_detection_rule": trafficDetectionRules,
		})
	}

	return output
}

func flattenNetworkDDoSCustomPolicySubResourceIds(input *[]ddoscustompolicies.SubResource) []string {
	output := make([]string, 0)
	if input == nil {
		return output
	}

	for _, item := range *input {
		output = append(output, pointer.From(item.Id))
	}

	return output
}
