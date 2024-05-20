// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package azurestackhci

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automanage/2022-05-04/configurationprofilehciassignments"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automanage/2022-05-04/configurationprofiles"
	"github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2024-01-01/clusters"
	"github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2024-01-01/deploymentsettings"
	"github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2024-01-01/storagecontainers"
	"github.com/hashicorp/go-azure-sdk/resource-manager/extendedlocation/2021-08-15/customlocations"
	"github.com/hashicorp/go-azure-sdk/resource-manager/hybridcompute/2022-11-10/machines"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resourceconnector/2022-10-27/appliances"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/azurestackhci/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceArmStackHCICluster() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceArmStackHCIClusterCreate,
		Read:   resourceArmStackHCIClusterRead,
		Update: resourceArmStackHCIClusterUpdate,
		Delete: resourceArmStackHCIClusterDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(3 * time.Hour),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(3 * time.Hour),
			Delete: pluginsdk.DefaultTimeout(1 * time.Hour),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := clusters.ParseClusterID(id)
			return err
		}),

		CustomizeDiff: pluginsdk.CustomDiffWithAll(
			pluginsdk.ForceNewIfChange("deployment_setting", func(ctx context.Context, old, new, meta interface{}) bool {
				return len(old.([]interface{})) > 0 && len(old.([]interface{})) != len(new.([]interface{}))
			}),
		),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ClusterName,
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"location": commonschema.Location(),

			"client_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},

			"tenant_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},

			"automanage_configuration_id": {
				// TODO: this field should be removed in 4.0 - there's an "association" API specifically for this purpose
				// so we should be outputting this as an association resource.
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: configurationprofiles.ValidateConfigurationProfileID,
			},

			"cloud_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"deployment_setting": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"arc_resource_ids": {
							Type:     pluginsdk.TypeList,
							Required: true,
							ForceNew: true,
							Elem: &pluginsdk.Schema{
								Type:         pluginsdk.TypeString,
								ValidateFunc: machines.ValidateMachineID,
							},
						},

						"version": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ForceNew: true,
							ValidateFunc: validation.StringMatch(
								regexp.MustCompile(`^[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+$`),
								"the version must be a set of numbers separated by dots: `10.0.0.1`",
							),
						},

						"scale_unit": {
							Type:     pluginsdk.TypeList,
							Required: true,
							ForceNew: true,
							MaxItems: 1,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"adou_path": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ForceNew:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},

									"cluster": {
										Type:     pluginsdk.TypeList,
										Required: true,
										ForceNew: true,
										MaxItems: 1,
										Elem: &pluginsdk.Resource{
											Schema: map[string]*pluginsdk.Schema{
												// "name": {
												// 	Type:     pluginsdk.TypeString,
												// 	Required: true,
												// 	ForceNew: true,
												// 	ValidateFunc: validation.StringMatch(
												// 		regexp.MustCompile("^[a-zA-Z0-9-]{3,15}$"),
												// 		"the cluster name must be 3-15 characters long and contain only letters, numbers and hyphens",
												// 	),
												// },

												"azure_service_endpoint": {
													Type:         pluginsdk.TypeString,
													Required:     true,
													ForceNew:     true,
													ValidateFunc: validation.StringIsNotEmpty,
												},

												"cloud_account_name": {
													Type:         pluginsdk.TypeString,
													Required:     true,
													ForceNew:     true,
													ValidateFunc: validation.StringIsNotEmpty,
												},

												"witness_type": {
													Type:     pluginsdk.TypeString,
													Required: true,
													ForceNew: true,
													ValidateFunc: validation.StringInSlice([]string{
														"Cloud",
														"FileShare",
													}, false),
												},

												"witness_path": {
													Type:         pluginsdk.TypeString,
													Required:     true,
													ForceNew:     true,
													ValidateFunc: validation.StringIsNotEmpty,
												},
											},
										},
									},

									"domain_fqdn": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ForceNew:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},

									"host_network": {
										Type:     pluginsdk.TypeList,
										Required: true,
										ForceNew: true,
										MaxItems: 1,
										Elem: &pluginsdk.Resource{
											Schema: map[string]*pluginsdk.Schema{
												"intent": {
													Type:     pluginsdk.TypeList,
													Required: true,
													ForceNew: true,
													MinItems: 1,
													Elem: &pluginsdk.Resource{
														Schema: map[string]*pluginsdk.Schema{
															"name": {
																Type:         pluginsdk.TypeString,
																Required:     true,
																ForceNew:     true,
																ValidateFunc: validation.StringIsNotEmpty,
															},

															"adapter": {
																Type:     pluginsdk.TypeList,
																Required: true,
																ForceNew: true,
																MinItems: 1,
																Elem: &pluginsdk.Schema{
																	Type:         pluginsdk.TypeString,
																	ValidateFunc: validation.StringIsNotEmpty,
																},
															},

															"traffic_type": {
																Type:     pluginsdk.TypeList,
																Required: true,
																ForceNew: true,
																MinItems: 1,
																Elem: &pluginsdk.Schema{
																	Type: pluginsdk.TypeString,
																	ValidateFunc: validation.StringInSlice([]string{
																		"Compute",
																		"Storage",
																		"Management",
																	}, false),
																},
															},

															"override_adapter_property": {
																Type:     pluginsdk.TypeList,
																Optional: true,
																ForceNew: true,
																MaxItems: 1,
																Elem: &pluginsdk.Resource{
																	Schema: map[string]*pluginsdk.Schema{
																		"jumbo_packet": {
																			Type:         pluginsdk.TypeString,
																			Optional:     true,
																			ForceNew:     true,
																			ValidateFunc: validation.StringIsNotEmpty,
																		},

																		"network_direct": {
																			Type:         pluginsdk.TypeString,
																			Optional:     true,
																			ForceNew:     true,
																			ValidateFunc: validation.StringIsNotEmpty,
																		},

																		"network_direct_technology": {
																			Type:         pluginsdk.TypeString,
																			Optional:     true,
																			ForceNew:     true,
																			ValidateFunc: validation.StringIsNotEmpty,
																		},
																	},
																},
															},

															"override_adapter_property_enabled": {
																Type:     pluginsdk.TypeBool,
																Optional: true,
																ForceNew: true,
																Default:  false,
															},

															"override_qos_policy": {
																Type:     pluginsdk.TypeList,
																Optional: true,
																ForceNew: true,
																MaxItems: 1,
																Elem: &pluginsdk.Resource{
																	Schema: map[string]*pluginsdk.Schema{
																		"bandwidth_percentage_smb": {
																			Type:         pluginsdk.TypeString,
																			Optional:     true,
																			ForceNew:     true,
																			ValidateFunc: validation.StringIsNotEmpty,
																		},

																		"priority_value8021_action_cluster": {
																			Type:         pluginsdk.TypeString,
																			Optional:     true,
																			ForceNew:     true,
																			ValidateFunc: validation.StringIsNotEmpty,
																		},

																		"priority_value8021_action_smb": {
																			Type:         pluginsdk.TypeString,
																			Optional:     true,
																			ForceNew:     true,
																			ValidateFunc: validation.StringIsNotEmpty,
																		},
																	},
																},
															},

															"override_qos_policy_enabled": {
																Type:     pluginsdk.TypeBool,
																Optional: true,
																ForceNew: true,
																Default:  false,
															},

															"override_virtual_switch_configuration": {
																Type:     pluginsdk.TypeList,
																Optional: true,
																ForceNew: true,
																MaxItems: 1,
																Elem: &pluginsdk.Resource{
																	Schema: map[string]*pluginsdk.Schema{
																		"enable_iov": {
																			Type:         pluginsdk.TypeString,
																			Optional:     true,
																			ForceNew:     true,
																			ValidateFunc: validation.StringIsNotEmpty,
																		},

																		"load_balancing_algorithm": {
																			Type:         pluginsdk.TypeString,
																			Optional:     true,
																			ForceNew:     true,
																			ValidateFunc: validation.StringIsNotEmpty,
																		},
																	},
																},
															},

															"override_virtual_switch_configuration_enabled": {
																Type:     pluginsdk.TypeBool,
																Optional: true,
																ForceNew: true,
																Default:  false,
															},
														},
													},
												},

												"storage_network": {
													Type:     pluginsdk.TypeList,
													Required: true,
													ForceNew: true,
													MinItems: 1,
													Elem: &pluginsdk.Resource{
														Schema: map[string]*pluginsdk.Schema{
															"name": {
																Type:         pluginsdk.TypeString,
																Required:     true,
																ForceNew:     true,
																ValidateFunc: validation.StringIsNotEmpty,
															},

															"network_adapter_name": {
																Type:         pluginsdk.TypeString,
																Required:     true,
																ForceNew:     true,
																ValidateFunc: validation.StringIsNotEmpty,
															},

															"vlan_id": {
																Type:         pluginsdk.TypeString,
																Required:     true,
																ForceNew:     true,
																ValidateFunc: validation.StringIsNotEmpty,
															},
														},
													},
												},

												"storage_auto_ip_enabled": {
													Type:     pluginsdk.TypeBool,
													Optional: true,
													ForceNew: true,
													Default:  true,
												},

												"storage_connectivity_switchless_enabled": {
													Type:     pluginsdk.TypeBool,
													Optional: true,
													ForceNew: true,
													Default:  false,
												},
											},
										},
									},

									"infrastructure_network": {
										Type:     pluginsdk.TypeList,
										Required: true,
										ForceNew: true,
										MinItems: 1,
										Elem: &pluginsdk.Resource{
											Schema: map[string]*pluginsdk.Schema{
												"dns_server": {
													Type:     pluginsdk.TypeList,
													Required: true,
													ForceNew: true,
													MinItems: 1,
													Elem: &pluginsdk.Schema{
														Type:         pluginsdk.TypeString,
														ValidateFunc: validation.IsIPv4Address,
													},
												},

												"gateway": {
													Type:         pluginsdk.TypeString,
													Required:     true,
													ForceNew:     true,
													ValidateFunc: validation.IsIPv4Address,
												},

												"ip_pool": {
													Type:     pluginsdk.TypeList,
													Required: true,
													ForceNew: true,
													MinItems: 1,
													Elem: &pluginsdk.Resource{
														Schema: map[string]*pluginsdk.Schema{
															"starting_address": {
																Type:         pluginsdk.TypeString,
																Required:     true,
																ForceNew:     true,
																ValidateFunc: validation.IsIPv4Address,
															},

															"ending_address": {
																Type:         pluginsdk.TypeString,
																Required:     true,
																ForceNew:     true,
																ValidateFunc: validation.IsIPv4Address,
															},
														},
													},
												},

												"subnet_mask": {
													Type:         pluginsdk.TypeString,
													Required:     true,
													ForceNew:     true,
													ValidateFunc: validation.IsIPv4Address,
												},

												"dhcp_enabled": {
													Type:     pluginsdk.TypeBool,
													Optional: true,
													ForceNew: true,
													Default:  false,
												},
											},
										},
									},

									"naming_prefix": {
										Type:     pluginsdk.TypeString,
										Required: true,
										ForceNew: true,
										ValidateFunc: validation.StringMatch(
											regexp.MustCompile("^[a-zA-Z0-9-]{1,8}$"),
											"the naming prefix must be 1-8 characters long and contain only letters, numbers and hyphens",
										),
									},

									"optional_service": {
										Type:     pluginsdk.TypeList,
										Required: true,
										ForceNew: true,
										MaxItems: 1,
										Elem: &pluginsdk.Resource{
											Schema: map[string]*pluginsdk.Schema{
												"custom_location": {
													Type:         pluginsdk.TypeString,
													Required:     true,
													ForceNew:     true,
													ValidateFunc: validation.StringIsNotEmpty,
												},
											},
										},
									},

									"physical_node": {
										Type:     pluginsdk.TypeList,
										Required: true,
										ForceNew: true,
										MinItems: 1,
										Elem: &pluginsdk.Resource{
											Schema: map[string]*pluginsdk.Schema{
												"name": {
													Type:         pluginsdk.TypeString,
													Required:     true,
													ForceNew:     true,
													ValidateFunc: validation.StringIsNotEmpty,
												},

												"ipv4_address": {
													Type:         pluginsdk.TypeString,
													Required:     true,
													ForceNew:     true,
													ValidateFunc: validation.IsIPv4Address,
												},
											},
										},
									},

									"secrets_location": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ForceNew:     true,
										ValidateFunc: validation.IsURLWithHTTPS,
									},

									"storage": {
										Type:     pluginsdk.TypeList,
										Required: true,
										ForceNew: true,
										MaxItems: 1,
										Elem: &pluginsdk.Resource{
											Schema: map[string]*pluginsdk.Schema{
												"configuration_mode": {
													Type:     pluginsdk.TypeString,
													Required: true,
													ForceNew: true,
													ValidateFunc: validation.StringInSlice([]string{
														"Express",
														"InfraOnly",
														"KeepStorage",
													}, false),
												},
											},
										},
									},

									"streaming_data_client_enabled": {
										Type:     pluginsdk.TypeBool,
										Optional: true,
										Default:  true,
										ForceNew: true,
									},

									"eu_location_enabled": {
										Type:     pluginsdk.TypeBool,
										Optional: true,
										Default:  false,
										ForceNew: true,
									},

									"episodic_data_upload_enabled": {
										Type:     pluginsdk.TypeBool,
										Optional: true,
										Default:  true,
										ForceNew: true,
									},

									"bitlocker_boot_volume_enabled": {
										Type:     pluginsdk.TypeBool,
										Optional: true,
										Default:  true,
										ForceNew: true,
									},

									"bitlocker_data_volume_enabled": {
										Type:     pluginsdk.TypeBool,
										Optional: true,
										Default:  true,
										ForceNew: true,
									},

									"credential_guard_enabled": {
										Type:     pluginsdk.TypeBool,
										Optional: true,
										Default:  false,
										ForceNew: true,
									},

									"drift_control_enabled": {
										Type:     pluginsdk.TypeBool,
										Optional: true,
										Default:  true,
										ForceNew: true,
									},

									"drtm_protection_enabled": {
										Type:     pluginsdk.TypeBool,
										Optional: true,
										Default:  true,
										ForceNew: true,
									},

									"hvci_protection_enabled": {
										Type:     pluginsdk.TypeBool,
										Optional: true,
										Default:  true,
										ForceNew: true,
									},

									"side_channel_mitigation_enabled": {
										Type:     pluginsdk.TypeBool,
										Optional: true,
										Default:  true,
										ForceNew: true,
									},

									"smb_signing_enabled": {
										Type:     pluginsdk.TypeBool,
										Optional: true,
										Default:  true,
										ForceNew: true,
									},

									"smb_cluster_encryption_enabled": {
										Type:     pluginsdk.TypeBool,
										Optional: true,
										Default:  false,
										ForceNew: true,
									},

									"wdac_enabled": {
										Type:     pluginsdk.TypeBool,
										Optional: true,
										Default:  true,
										ForceNew: true,
									},
								},
							},
						},
					},
				},
			},

			"service_endpoint": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"resource_provider_object_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"identity": commonschema.SystemAssignedIdentityOptional(),

			"tags": commonschema.Tags(),
		},
	}
}

func resourceArmStackHCIClusterCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AzureStackHCI.Clusters
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := clusters.NewClusterID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	existing, err := client.Get(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
	}

	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_stack_hci_cluster", id.ID())
	}

	cluster := clusters.Cluster{
		Location: location.Normalize(d.Get("location").(string)),
		Properties: &clusters.ClusterProperties{
			AadClientId: utils.String(d.Get("client_id").(string)),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if v, ok := d.GetOk("identity"); ok {
		cluster.Identity = expandSystemAssigned(v.([]interface{}))
	}

	if v, ok := d.GetOk("tenant_id"); ok {
		cluster.Properties.AadTenantId = utils.String(v.(string))
	} else {
		tenantId := meta.(*clients.Client).Account.TenantId
		cluster.Properties.AadTenantId = utils.String(tenantId)
	}

	if _, err := client.Create(ctx, id, cluster); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	if v, ok := d.GetOk("automanage_configuration_id"); ok {
		configurationProfilesClient := meta.(*clients.Client).Automanage.ConfigurationProfilesClient
		hciAssignmentsClient := meta.(*clients.Client).Automanage.ConfigurationProfileHCIAssignmentsClient

		configurationProfileId, err := configurationprofiles.ParseConfigurationProfileID(v.(string))
		if err != nil {
			return err
		}

		if _, err = configurationProfilesClient.Get(ctx, *configurationProfileId); err != nil {
			return fmt.Errorf("checking for existing %s: %+v", configurationProfileId, err)
		}

		hciAssignmentId := configurationprofilehciassignments.NewConfigurationProfileAssignmentID(subscriptionId, id.ResourceGroupName, id.ClusterName, "default")
		assignmentsResp, err := hciAssignmentsClient.Get(ctx, hciAssignmentId)
		if err != nil && !response.WasNotFound(assignmentsResp.HttpResponse) {
			return fmt.Errorf("checking for existing %s: %+v", hciAssignmentId, err)
		}

		if response.WasNotFound(assignmentsResp.HttpResponse) {
			properties := configurationprofilehciassignments.ConfigurationProfileAssignment{
				Properties: &configurationprofilehciassignments.ConfigurationProfileAssignmentProperties{
					ConfigurationProfile: utils.String(configurationProfileId.ID()),
				},
			}

			if _, err := hciAssignmentsClient.CreateOrUpdate(ctx, hciAssignmentId, properties); err != nil {
				return fmt.Errorf("creating %s: %+v", hciAssignmentId, err)
			}
		}
	}

	d.SetId(id.ID())

	deploymentSettingRaw := d.Get("deployment_setting").([]interface{})
	if len(deploymentSettingRaw) > 0 {
		deploymentSetting := deploymentSettingRaw[0].(map[string]interface{})
		deploymentSettingPayload := deploymentsettings.DeploymentSetting{
			Properties: &deploymentsettings.DeploymentSettingsProperties{
				ArcNodeResourceIds: expandDeploymentSettingArcNodeResourceIds(deploymentSetting["arc_resource_ids"].([]interface{})),
				DeploymentMode:     deploymentsettings.DeploymentModeValidate,
				DeploymentConfiguration: deploymentsettings.DeploymentConfiguration{
					Version:    pointer.To(deploymentSetting["version"].(string)),
					ScaleUnits: ExpandDeploymentSettingScaleUnits(deploymentSetting["scale_unit"].([]interface{}), id.ClusterName),
				},
			},
		}

		deploymentSettingsClient := meta.(*clients.Client).AzureStackHCI.DeploymentSettings
		// the deploymentSetting can only have "default" as name
		deploymentSettingId := deploymentsettings.NewDeploymentSettingID(id.SubscriptionId, id.ResourceGroupName, id.ClusterName, "default")

		// do the validation
		if err := deploymentSettingsClient.CreateOrUpdateThenPoll(ctx, deploymentSettingId, deploymentSettingPayload); err != nil {
			return fmt.Errorf("validating %s: %+v", id, err)
		}

		// do the deployment
		deploymentSettingPayload.Properties.DeploymentMode = deploymentsettings.DeploymentModeDeploy
		if err := deploymentSettingsClient.CreateOrUpdateThenPoll(ctx, deploymentSettingId, deploymentSettingPayload); err != nil {
			return fmt.Errorf("deploying %s: %+v", id, err)
		}
	}

	return resourceArmStackHCIClusterRead(d, meta)
}

func resourceArmStackHCIClusterRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AzureStackHCI.Clusters
	hciAssignmentsClient := meta.(*clients.Client).Automanage.ConfigurationProfileHCIAssignmentsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := clusters.ParseClusterID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s was not found - removing from state!", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.ClusterName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.Normalize(model.Location))
		d.Set("identity", flattenSystemAssigned(model.Identity))

		if props := model.Properties; props != nil {
			d.Set("client_id", props.AadClientId)
			d.Set("tenant_id", props.AadTenantId)
			d.Set("cloud_id", props.CloudId)
			d.Set("service_endpoint", props.ServiceEndpoint)
			d.Set("resource_provider_object_id", props.ResourceProviderObjectId)
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	hclAssignmentId := configurationprofilehciassignments.NewConfigurationProfileAssignmentID(id.SubscriptionId, id.ResourceGroupName, id.ClusterName, "default")
	assignmentResp, err := hciAssignmentsClient.Get(ctx, hclAssignmentId)
	if err != nil && !response.WasNotFound(assignmentResp.HttpResponse) {
		return err
	}
	configId := ""
	if model := assignmentResp.Model; model != nil && model.Properties != nil && model.Properties.ConfigurationProfile != nil {
		parsed, err := configurationprofiles.ParseConfigurationProfileIDInsensitively(*model.Properties.ConfigurationProfile)
		if err != nil {
			return err
		}
		configId = parsed.ID()
	}
	d.Set("automanage_configuration_id", configId)

	deploymentSettingClient := meta.(*clients.Client).AzureStackHCI.DeploymentSettings
	// the deploymentSetting can only have "default" as name
	deploymentSettingId := deploymentsettings.NewDeploymentSettingID(id.SubscriptionId, id.ResourceGroupName, id.ClusterName, "default")
	deploymentSettingResp, err := deploymentSettingClient.Get(ctx, deploymentSettingId)
	if err != nil && !response.WasNotFound(deploymentSettingResp.HttpResponse) {
		return fmt.Errorf("retrieving %s: %+v", deploymentSettingId, err)
	}

	deploymentSetting := make([]interface{}, 0)
	if model := deploymentSettingResp.Model; model != nil && model.Properties != nil {
		deploymentSetting = []interface{}{
			map[string]interface{}{
				"arc_resource_ids": model.Properties.ArcNodeResourceIds,
				"version":          pointer.From(model.Properties.DeploymentConfiguration.Version),
				"scale_unit":       FlattenDeploymentSettingScaleUnits(model.Properties.DeploymentConfiguration.ScaleUnits),
			},
		}
	}
	d.Set("deployment_setting", deploymentSetting)

	return nil
}

func resourceArmStackHCIClusterUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AzureStackHCI.Clusters
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := clusters.ParseClusterID(d.Id())
	if err != nil {
		return err
	}

	cluster := clusters.ClusterPatch{}

	if d.HasChange("tags") {
		cluster.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	if d.HasChange("identity") {
		cluster.Identity = expandSystemAssigned(d.Get("identity").([]interface{}))
	}

	if _, err := client.Update(ctx, *id, cluster); err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	if d.HasChange("automanage_configuration_id") {
		hciAssignmentClient := meta.(*clients.Client).Automanage.ConfigurationProfileHCIAssignmentsClient
		configurationProfilesClient := meta.(*clients.Client).Automanage.ConfigurationProfilesClient
		hciAssignmentId := configurationprofilehciassignments.NewConfigurationProfileAssignmentID(id.SubscriptionId, id.ResourceGroupName, id.ClusterName, "default")

		if v, ok := d.GetOk("automanage_configuration_id"); ok {
			configurationProfileId, err := configurationprofiles.ParseConfigurationProfileID(v.(string))
			if err != nil {
				return err
			}

			if _, err = configurationProfilesClient.Get(ctx, *configurationProfileId); err != nil {
				return fmt.Errorf("checking for existing %s: %+v", configurationProfileId, err)
			}

			properties := configurationprofilehciassignments.ConfigurationProfileAssignment{
				Properties: &configurationprofilehciassignments.ConfigurationProfileAssignmentProperties{
					ConfigurationProfile: utils.String(configurationProfileId.ID()),
				},
			}

			if _, err := hciAssignmentClient.CreateOrUpdate(ctx, hciAssignmentId, properties); err != nil {
				return fmt.Errorf("creating %s: %+v", hciAssignmentId, err)
			}
		} else {
			assignmentResp, err := hciAssignmentClient.Get(ctx, hciAssignmentId)
			if err != nil && !response.WasNotFound(assignmentResp.HttpResponse) {
				return err
			}

			if !response.WasNotFound(assignmentResp.HttpResponse) {
				if _, err := hciAssignmentClient.Delete(ctx, hciAssignmentId); err != nil {
					return fmt.Errorf("deleting %s: %+v", id, err)
				}
			}
		}
	}

	if d.HasChange("deployment_setting") {
		deploymentSettingPayload := deploymentsettings.DeploymentSetting{
			Properties: &deploymentsettings.DeploymentSettingsProperties{
				ArcNodeResourceIds: expandDeploymentSettingArcNodeResourceIds(d.Get("arc_resource_ids").([]interface{})),
				DeploymentMode:     deploymentsettings.DeploymentModeValidate,
				DeploymentConfiguration: deploymentsettings.DeploymentConfiguration{
					Version:    pointer.To(d.Get("version").(string)),
					ScaleUnits: ExpandDeploymentSettingScaleUnits(d.Get("scale_unit").([]interface{}), id.ClusterName),
				},
			},
		}

		deploymentSettingsClient := meta.(*clients.Client).AzureStackHCI.DeploymentSettings
		// the deploymentSetting can only have "default" as name
		deploymentSettingId := deploymentsettings.NewDeploymentSettingID(id.SubscriptionId, id.ResourceGroupName, id.ClusterName, "default")

		// do the validation
		if err := deploymentSettingsClient.CreateOrUpdateThenPoll(ctx, deploymentSettingId, deploymentSettingPayload); err != nil {
			return fmt.Errorf("validating %s: %+v", id, err)
		}

		// do the deployment
		deploymentSettingPayload.Properties.DeploymentMode = deploymentsettings.DeploymentModeDeploy
		if err := deploymentSettingsClient.CreateOrUpdateThenPoll(ctx, deploymentSettingId, deploymentSettingPayload); err != nil {
			return fmt.Errorf("deploying %s: %+v", id, err)
		}
	}

	return resourceArmStackHCIClusterRead(d, meta)
}

func resourceArmStackHCIClusterDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AzureStackHCI.Clusters
	hciAssignmentClient := meta.(*clients.Client).Automanage.ConfigurationProfileHCIAssignmentsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := clusters.ParseClusterID(d.Id())
	if err != nil {
		return err
	}

	hciAssignmentId := configurationprofilehciassignments.NewConfigurationProfileAssignmentID(id.SubscriptionId, id.ResourceGroupName, id.ClusterName, "default")
	assignmentResp, err := hciAssignmentClient.Get(ctx, hciAssignmentId)
	if err != nil && !response.WasNotFound(assignmentResp.HttpResponse) {
		return err
	}

	if !response.WasNotFound(assignmentResp.HttpResponse) {
		if _, err := hciAssignmentClient.Delete(ctx, hciAssignmentId); err != nil {
			return fmt.Errorf("deleting %s: %+v", id, err)
		}
	}

	deploymentSettingsClient := meta.(*clients.Client).AzureStackHCI.DeploymentSettings
	deploymentSettingId := deploymentsettings.NewDeploymentSettingID(id.SubscriptionId, id.ResourceGroupName, id.ClusterName, "default")
	deploymentSettingResp, err := deploymentSettingsClient.Get(ctx, deploymentSettingId)
	if err != nil && !response.WasNotFound(deploymentSettingResp.HttpResponse) {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if !response.WasNotFound(deploymentSettingResp.HttpResponse) {
		if err := deploymentSettingsClient.DeleteThenPoll(ctx, deploymentSettingId); err != nil {
			return fmt.Errorf("deleting %s: %+v", deploymentSettingId, err)
		}
	}

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	if meta.(*clients.Client).Features.AzureStackHci.DeleteArcBridgeOnDestroy && !response.WasNotFound(deploymentSettingResp.HttpResponse) {
		applianceName := fmt.Sprintf("%s-arcbridge", id.ClusterName)
		applianceId := appliances.NewApplianceID(id.SubscriptionId, id.ResourceGroupName, applianceName)

		log.Printf("[DEBUG] delete_arc_bridge_on_destroy is enabled - removing Arc Bridge %s", applianceId)

		applianceClient := meta.(*clients.Client).ArcResourceBridge.AppliancesClient
		applianceResp, err := applianceClient.Get(ctx, applianceId)
		if err != nil && !response.WasNotFound(applianceResp.HttpResponse) {
			return fmt.Errorf("retrieving %s: %+v", applianceId, err)
		}
		if !response.WasNotFound(applianceResp.HttpResponse) {
			if err := applianceClient.DeleteThenPoll(ctx, applianceId); err != nil {
				return fmt.Errorf("deleting %s: %+v", applianceId, err)
			}
		}
	}

	if meta.(*clients.Client).Features.AzureStackHci.DeleteCustomLocationOnDestroy && !response.WasNotFound(deploymentSettingResp.HttpResponse) {
		log.Printf("[DEBUG] delete_custom_location_on_destroy is enabled - removing Custom Location and Azure Stack HCI Storage Containers in the location!")
		var customLocationName string
		if deploymentSettingResp.Model != nil && deploymentSettingResp.Model.Properties != nil &&
			len(deploymentSettingResp.Model.Properties.DeploymentConfiguration.ScaleUnits) > 0 &&
			deploymentSettingResp.Model.Properties.DeploymentConfiguration.ScaleUnits[0].DeploymentData.OptionalServices != nil {
			customLocationName = pointer.From(deploymentSettingResp.Model.Properties.DeploymentConfiguration.ScaleUnits[0].DeploymentData.OptionalServices.CustomLocation)
		}

		// try to delete the Custom Location generated during deployment
		if customLocationName != "" {
			customLocationsClient := meta.(*clients.Client).ExtendedLocation.CustomLocations
			customLocationId := customlocations.NewCustomLocationID(id.SubscriptionId, id.ResourceGroupName, customLocationName)

			log.Printf("[DEBUG] removing Custom Location %s", customLocationId)

			customLocationResp, err := customLocationsClient.Get(ctx, customLocationId)
			if err != nil && !response.WasNotFound(customLocationResp.HttpResponse) {
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			if !response.WasNotFound(customLocationResp.HttpResponse) {
				// try to delete the HCI Staroge Containers generated during deployment in the Custom Location
				storageContainerClient := meta.(*clients.Client).AzureStackHCI.StorageContainers
				resourceGroupId := commonids.NewResourceGroupID(id.SubscriptionId, id.ResourceGroupName)
				storageContainers, err := storageContainerClient.ListComplete(ctx, resourceGroupId)
				if err != nil {
					return fmt.Errorf("retrieving Stack HCI Storage Containers under %s: %+v", resourceGroupId.ID(), err)
				}

				// find all Storage Containers under the Custom Location
				storageContainerNamePattern := regexp.MustCompile(`UserStorage[0-9]+-[a-z0-9]{32}`)
				for _, v := range storageContainers.Items {
					if v.Id != nil && v.ExtendedLocation != nil && v.ExtendedLocation.Name != nil && strings.EqualFold(*v.ExtendedLocation.Name, customLocationId.ID()) && v.Name != nil && storageContainerNamePattern.Match([]byte(*v.Name)) {
						log.Printf("[DEBUG] removing Azure Stack HCI Storage Container %s", *v.Id)

						storageContainerId, err := storagecontainers.ParseStorageContainerIDInsensitively(*v.Id)
						if err != nil {
							return err
						}

						if err := storageContainerClient.DeleteThenPoll(ctx, *storageContainerId); err != nil {
							return fmt.Errorf("deleting %s: %+v", *storageContainerId, err)
						}
					}
				}

				if err := customLocationsClient.DeleteThenPoll(ctx, customLocationId); err != nil {
					return fmt.Errorf("deleting %s: %+v", customLocationId, err)
				}
			}
		}
	}

	return nil
}

// API does not accept userAssignedIdentity as in swagger https://github.com/Azure/azure-rest-api-specs/issues/28260
func expandSystemAssigned(input []interface{}) *identity.SystemAndUserAssignedMap {
	if len(input) == 0 || input[0] == nil {
		return &identity.SystemAndUserAssignedMap{
			Type: identity.TypeNone,
		}
	}

	return &identity.SystemAndUserAssignedMap{
		Type: identity.TypeSystemAssigned,
	}
}

func flattenSystemAssigned(input *identity.SystemAndUserAssignedMap) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	if input.Type == identity.TypeNone {
		return []interface{}{}
	}

	return []interface{}{
		map[string]interface{}{
			"type":         input.Type,
			"principal_id": input.PrincipalId,
			"tenant_id":    input.TenantId,
		},
	}
}

func expandDeploymentSettingArcNodeResourceIds(input []interface{}) []string {
	if len(input) == 0 || input[0] == nil {
		return make([]string, 0)
	}

	results := make([]string, 0, len(input))
	for _, item := range input {
		results = append(results, item.(string))
	}
	return results
}

func ExpandDeploymentSettingScaleUnits(input []interface{}, clusterName string) []deploymentsettings.ScaleUnits {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	v := input[0].(map[string]interface{})

	return []deploymentsettings.ScaleUnits{
		{
			DeploymentData: deploymentsettings.DeploymentData{
				AdouPath:              pointer.To(v["adou_path"].(string)),
				Cluster:               ExpandDeploymentSettingCluster(v["cluster"].([]interface{}), clusterName),
				DomainFqdn:            pointer.To(v["domain_fqdn"].(string)),
				HostNetwork:           ExpandDeploymentSettingHostNetwork(v["host_network"].([]interface{})),
				InfrastructureNetwork: ExpandDeploymentSettingInfrastructureNetwork(v["infrastructure_network"].([]interface{})),
				NamingPrefix:          pointer.To(v["naming_prefix"].(string)),
				OptionalServices:      ExpandDeploymentSettingOptionalService(v["optional_service"].([]interface{})),
				PhysicalNodes:         ExpandDeploymentSettingPhysicalNode(v["physical_node"].([]interface{})),
				SecretsLocation:       pointer.To(v["secrets_location"].(string)),
				Storage:               ExpandDeploymentSettingStorage(v["storage"].([]interface{})),
				Observability: &deploymentsettings.Observability{
					EpisodicDataUpload:  pointer.To(v["episodic_data_upload_enabled"].(bool)),
					EuLocation:          pointer.To(v["eu_location_enabled"].(bool)),
					StreamingDataClient: pointer.To(v["streaming_data_client_enabled"].(bool)),
				},
				SecuritySettings: &deploymentsettings.DeploymentSecuritySettings{
					BitlockerBootVolume:           pointer.To(v["bitlocker_boot_volume_enabled"].(bool)),
					BitlockerDataVolumes:          pointer.To(v["bitlocker_data_volume_enabled"].(bool)),
					CredentialGuardEnforced:       pointer.To(v["credential_guard_enabled"].(bool)),
					DriftControlEnforced:          pointer.To(v["drift_control_enabled"].(bool)),
					DrtmProtection:                pointer.To(v["drtm_protection_enabled"].(bool)),
					HvciProtection:                pointer.To(v["hvci_protection_enabled"].(bool)),
					SideChannelMitigationEnforced: pointer.To(v["side_channel_mitigation_enabled"].(bool)),
					SmbClusterEncryption:          pointer.To(v["smb_cluster_encryption_enabled"].(bool)),
					SmbSigningEnforced:            pointer.To(v["smb_signing_enabled"].(bool)),
					WdacEnforced:                  pointer.To(v["wdac_enabled"].(bool)),
				},
			},
		},
	}
}

func FlattenDeploymentSettingScaleUnits(input []deploymentsettings.ScaleUnits) []interface{} {
	if len(input) == 0 {
		return make([]interface{}, 0)
	}

	results := make([]interface{}, 0, len(input))
	for _, item := range input {
		result := map[string]interface{}{
			"adou_path":              pointer.From(item.DeploymentData.AdouPath),
			"cluster":                FlattenDeploymentSettingCluster(item.DeploymentData.Cluster),
			"domain_fqdn":            pointer.From(item.DeploymentData.DomainFqdn),
			"host_network":           FlattenDeploymentSettingHostNetwork(item.DeploymentData.HostNetwork),
			"infrastructure_network": FlattenDeploymentSettingInfrastructureNetwork(item.DeploymentData.InfrastructureNetwork),
			"naming_prefix":          pointer.From(item.DeploymentData.NamingPrefix),
			"optional_service":       FlattenDeploymentSettingOptionalService(item.DeploymentData.OptionalServices),
			"physical_node":          FlattenDeploymentSettingPhysicalNode(item.DeploymentData.PhysicalNodes),
			"secrets_location":       pointer.From(item.DeploymentData.SecretsLocation),
			"storage":                FlattenDeploymentSettingStorage(item.DeploymentData.Storage),
		}

		if observability := item.DeploymentData.Observability; observability != nil {
			result["episodic_data_upload_enabled"] = pointer.From(observability.EpisodicDataUpload)
			result["eu_location_enabled"] = pointer.From(observability.EuLocation)
			result["streaming_data_client_enabled"] = pointer.From(observability.StreamingDataClient)
		}

		if securitySettings := item.DeploymentData.SecuritySettings; securitySettings != nil {
			result["bitlocker_boot_volume_enabled"] = pointer.From(securitySettings.BitlockerBootVolume)
			result["bitlocker_data_volume_enabled"] = pointer.From(securitySettings.BitlockerDataVolumes)
			result["credential_guard_enabled"] = pointer.From(securitySettings.CredentialGuardEnforced)
			result["drift_control_enabled"] = pointer.From(securitySettings.DriftControlEnforced)
			result["drtm_protection_enabled"] = pointer.From(securitySettings.DrtmProtection)
			result["hvci_protection_enabled"] = pointer.From(securitySettings.HvciProtection)
			result["side_channel_mitigation_enabled"] = pointer.From(securitySettings.SideChannelMitigationEnforced)
			result["smb_cluster_encryption_enabled"] = pointer.From(securitySettings.SmbClusterEncryption)
			result["smb_signing_enabled"] = pointer.From(securitySettings.SmbSigningEnforced)
			result["wdac_enabled"] = pointer.From(securitySettings.WdacEnforced)
		}

		results = append(results, result)
	}

	return results
}

func ExpandDeploymentSettingCluster(input []interface{}, clusterName string) *deploymentsettings.DeploymentCluster {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	v := input[0].(map[string]interface{})

	return &deploymentsettings.DeploymentCluster{
		Name:                 pointer.To(clusterName),
		AzureServiceEndpoint: pointer.To(v["azure_service_endpoint"].(string)),
		CloudAccountName:     pointer.To(v["cloud_account_name"].(string)),
		WitnessType:          pointer.To(v["witness_type"].(string)),
		WitnessPath:          pointer.To(v["witness_path"].(string)),
	}
}

func FlattenDeploymentSettingCluster(input *deploymentsettings.DeploymentCluster) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	return []interface{}{
		map[string]interface{}{
			"azure_service_endpoint": pointer.From(input.AzureServiceEndpoint),
			"cloud_account_name":     pointer.From(input.CloudAccountName),
			"witness_type":           pointer.From(input.WitnessType),
			"witness_path":           pointer.From(input.WitnessPath),
		},
	}
}

func ExpandDeploymentSettingHostNetwork(input []interface{}) *deploymentsettings.HostNetwork {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	v := input[0].(map[string]interface{})

	return &deploymentsettings.HostNetwork{
		Intents:                       ExpandDeploymentSettingHostNetworkIntent(v["intent"].([]interface{})),
		EnableStorageAutoIP:           pointer.To(v["storage_auto_ip_enabled"].(bool)),
		StorageConnectivitySwitchless: pointer.To(v["storage_connectivity_switchless_enabled"].(bool)),
		StorageNetworks:               ExpandDeploymentSettingHostNetworkStorageNetwork(v["storage_network"].([]interface{})),
	}
}

func FlattenDeploymentSettingHostNetwork(input *deploymentsettings.HostNetwork) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	return []interface{}{
		map[string]interface{}{
			"intent":                  FlattenDeploymentSettingHostNetworkIntent(input.Intents),
			"storage_auto_ip_enabled": pointer.From(input.EnableStorageAutoIP),
			"storage_connectivity_switchless_enabled": pointer.From(input.StorageConnectivitySwitchless),
			"storage_network":                         FlattenDeploymentSettingHostNetworkStorageNetwork(input.StorageNetworks),
		},
	}
}

func ExpandDeploymentSettingHostNetworkIntent(input []interface{}) *[]deploymentsettings.Intents {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	results := make([]deploymentsettings.Intents, 0, len(input))
	for _, item := range input {
		v := item.(map[string]interface{})
		results = append(results, deploymentsettings.Intents{
			Adapter:                             utils.ExpandStringSlice(v["adapter"].([]interface{})),
			AdapterPropertyOverrides:            ExpandHostNetworkIntentAdapterPropertyOverride(v["override_adapter_property"].([]interface{})),
			Name:                                pointer.To(v["name"].(string)),
			OverrideAdapterProperty:             pointer.To(v["override_adapter_property_enabled"].(bool)),
			OverrideQosPolicy:                   pointer.To(v["override_qos_policy_enabled"].(bool)),
			OverrideVirtualSwitchConfiguration:  pointer.To(v["override_virtual_switch_configuration_enabled"].(bool)),
			QosPolicyOverrides:                  ExpandHostNetworkIntentQosPolicyOverride(v["override_qos_policy"].([]interface{})),
			TrafficType:                         utils.ExpandStringSlice(v["traffic_type"].([]interface{})),
			VirtualSwitchConfigurationOverrides: ExpandHostNetworkIntentVirtualSwitchConfigurationOverride(v["override_virtual_switch_configuration"].([]interface{})),
		})
	}

	return &results
}

func FlattenDeploymentSettingHostNetworkIntent(input *[]deploymentsettings.Intents) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	results := make([]interface{}, 0, len(*input))
	for _, item := range *input {
		results = append(results, map[string]interface{}{
			"adapter":                           pointer.From(item.Adapter),
			"override_adapter_property":         FlattenHostNetworkIntentAdapterPropertyOverride(item.AdapterPropertyOverrides),
			"name":                              pointer.From(item.Name),
			"override_adapter_property_enabled": pointer.From(item.OverrideAdapterProperty),
			"override_qos_policy_enabled":       pointer.From(item.OverrideQosPolicy),
			"override_virtual_switch_configuration_enabled": pointer.From(item.OverrideVirtualSwitchConfiguration),
			"override_qos_policy":                           FlattenHostNetworkIntentQosPolicyOverride(item.QosPolicyOverrides),
			"traffic_type":                                  pointer.From(item.TrafficType),
			"override_virtual_switch_configuration":         FlattenHostNetworkIntentVirtualSwitchConfigurationOverride(item.VirtualSwitchConfigurationOverrides),
		})
	}

	return results
}

func ExpandHostNetworkIntentAdapterPropertyOverride(input []interface{}) *deploymentsettings.AdapterPropertyOverrides {
	if len(input) == 0 || input[0] == nil {
		return &deploymentsettings.AdapterPropertyOverrides{
			JumboPacket:             pointer.To(""),
			NetworkDirect:           pointer.To(""),
			NetworkDirectTechnology: pointer.To(""),
		}
	}

	v := input[0].(map[string]interface{})

	return &deploymentsettings.AdapterPropertyOverrides{
		JumboPacket:             pointer.To(v["jumbo_packet"].(string)),
		NetworkDirect:           pointer.To(v["network_direct"].(string)),
		NetworkDirectTechnology: pointer.To(v["network_direct_technology"].(string)),
	}
}

func FlattenHostNetworkIntentAdapterPropertyOverride(input *deploymentsettings.AdapterPropertyOverrides) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	jumboPacket := pointer.From(input.JumboPacket)
	networkDirect := pointer.From(input.NetworkDirect)
	networkDirectTechnology := pointer.From(input.NetworkDirectTechnology)

	// server will return the block with empty string in all fields by default
	if jumboPacket == "" && networkDirect == "" && networkDirectTechnology == "" {
		return make([]interface{}, 0)
	}

	return []interface{}{
		map[string]interface{}{
			"jumbo_packet":              jumboPacket,
			"network_direct":            networkDirect,
			"network_direct_technology": networkDirectTechnology,
		},
	}
}

func ExpandHostNetworkIntentQosPolicyOverride(input []interface{}) *deploymentsettings.QosPolicyOverrides {
	if len(input) == 0 || input[0] == nil {
		return &deploymentsettings.QosPolicyOverrides{
			BandwidthPercentageSMB:         pointer.To(""),
			PriorityValue8021ActionCluster: pointer.To(""),
			PriorityValue8021ActionSMB:     pointer.To(""),
		}
	}

	v := input[0].(map[string]interface{})

	return &deploymentsettings.QosPolicyOverrides{
		BandwidthPercentageSMB:         pointer.To(v["bandwidth_percentage_smb"].(string)),
		PriorityValue8021ActionCluster: pointer.To(v["priority_value8021_action_cluster"].(string)),
		PriorityValue8021ActionSMB:     pointer.To(v["priority_value8021_action_smb"].(string)),
	}
}

func FlattenHostNetworkIntentQosPolicyOverride(input *deploymentsettings.QosPolicyOverrides) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	bandwidthPercentageSMB := pointer.From(input.BandwidthPercentageSMB)
	priorityValue8021ActionCluster := pointer.From(input.PriorityValue8021ActionCluster)
	priorityValue8021ActionSMB := pointer.From(input.PriorityValue8021ActionSMB)

	// server will return the block with empty string in all fields by default
	if bandwidthPercentageSMB == "" && priorityValue8021ActionCluster == "" && priorityValue8021ActionSMB == "" {
		return make([]interface{}, 0)
	}

	return []interface{}{
		map[string]interface{}{
			"bandwidth_percentage_smb":          bandwidthPercentageSMB,
			"priority_value8021_action_cluster": priorityValue8021ActionCluster,
			"priority_value8021_action_smb":     priorityValue8021ActionSMB,
		},
	}
}

func ExpandHostNetworkIntentVirtualSwitchConfigurationOverride(input []interface{}) *deploymentsettings.VirtualSwitchConfigurationOverrides {
	if len(input) == 0 || input[0] == nil {
		return &deploymentsettings.VirtualSwitchConfigurationOverrides{
			EnableIov:              pointer.To(""),
			LoadBalancingAlgorithm: pointer.To(""),
		}
	}

	v := input[0].(map[string]interface{})

	return &deploymentsettings.VirtualSwitchConfigurationOverrides{
		EnableIov:              pointer.To(v["enable_iov"].(string)),
		LoadBalancingAlgorithm: pointer.To(v["load_balancing_algorithm"].(string)),
	}
}

func FlattenHostNetworkIntentVirtualSwitchConfigurationOverride(input *deploymentsettings.VirtualSwitchConfigurationOverrides) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	enableIov := pointer.From(input.EnableIov)
	loadBalancingAlgorithm := pointer.From(input.LoadBalancingAlgorithm)

	// server will return the block with empty string in all fields by default
	if enableIov == "" && loadBalancingAlgorithm == "" {
		return make([]interface{}, 0)
	}

	return []interface{}{
		map[string]interface{}{
			"enable_iov":               enableIov,
			"load_balancing_algorithm": loadBalancingAlgorithm,
		},
	}
}

func ExpandDeploymentSettingHostNetworkStorageNetwork(input []interface{}) *[]deploymentsettings.StorageNetworks {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	results := make([]deploymentsettings.StorageNetworks, 0, len(input))
	for _, item := range input {
		v := item.(map[string]interface{})
		results = append(results, deploymentsettings.StorageNetworks{
			Name:               pointer.To(v["name"].(string)),
			NetworkAdapterName: pointer.To(v["network_adapter_name"].(string)),
			VlanId:             pointer.To(v["vlan_id"].(string)),
		})
	}

	return &results
}

func FlattenDeploymentSettingHostNetworkStorageNetwork(input *[]deploymentsettings.StorageNetworks) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	results := make([]interface{}, 0, len(*input))
	for _, item := range *input {
		results = append(results, map[string]interface{}{
			"name":                 pointer.From(item.Name),
			"network_adapter_name": pointer.From(item.NetworkAdapterName),
			"vlan_id":              pointer.From(item.VlanId),
		})
	}

	return results
}

func ExpandDeploymentSettingInfrastructureNetwork(input []interface{}) *[]deploymentsettings.InfrastructureNetwork {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	results := make([]deploymentsettings.InfrastructureNetwork, 0, len(input))
	for _, item := range input {
		v := item.(map[string]interface{})
		results = append(results, deploymentsettings.InfrastructureNetwork{
			DnsServers: utils.ExpandStringSlice(v["dns_server"].([]interface{})),
			Gateway:    pointer.To(v["gateway"].(string)),
			IPPools:    ExpandDeploymentSettingInfrastructureNetworkIpPool(v["ip_pool"].([]interface{})),
			SubnetMask: pointer.To(v["subnet_mask"].(string)),
			UseDhcp:    pointer.To(v["dhcp_enabled"].(bool)),
		})
	}

	return &results
}

func FlattenDeploymentSettingInfrastructureNetwork(input *[]deploymentsettings.InfrastructureNetwork) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	results := make([]interface{}, 0, len(*input))
	for _, item := range *input {
		results = append(results, map[string]interface{}{
			"dhcp_enabled": pointer.From(item.UseDhcp),
			"dns_server":   pointer.From(item.DnsServers),
			"gateway":      pointer.From(item.Gateway),
			"ip_pool":      FlattenDeploymentSettingInfrastructureNetworkIpPool(item.IPPools),
			"subnet_mask":  pointer.From(item.SubnetMask),
		})
	}

	return results
}

func ExpandDeploymentSettingInfrastructureNetworkIpPool(input []interface{}) *[]deploymentsettings.IPPools {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	results := make([]deploymentsettings.IPPools, 0, len(input))
	for _, item := range input {
		v := item.(map[string]interface{})
		results = append(results, deploymentsettings.IPPools{
			EndingAddress:   pointer.To(v["ending_address"].(string)),
			StartingAddress: pointer.To(v["starting_address"].(string)),
		})
	}

	return &results
}

func FlattenDeploymentSettingInfrastructureNetworkIpPool(input *[]deploymentsettings.IPPools) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	results := make([]interface{}, 0, len(*input))
	for _, item := range *input {
		results = append(results, map[string]interface{}{
			"ending_address":   pointer.From(item.EndingAddress),
			"starting_address": pointer.From(item.StartingAddress),
		})
	}

	return results
}

func ExpandDeploymentSettingOptionalService(input []interface{}) *deploymentsettings.OptionalServices {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	v := input[0].(map[string]interface{})

	return &deploymentsettings.OptionalServices{
		CustomLocation: pointer.To(v["custom_location"].(string)),
	}
}

func FlattenDeploymentSettingOptionalService(input *deploymentsettings.OptionalServices) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	return []interface{}{
		map[string]interface{}{
			"custom_location": pointer.From(input.CustomLocation),
		},
	}
}

func ExpandDeploymentSettingPhysicalNode(input []interface{}) *[]deploymentsettings.PhysicalNodes {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	results := make([]deploymentsettings.PhysicalNodes, 0, len(input))
	for _, item := range input {
		v := item.(map[string]interface{})
		results = append(results, deploymentsettings.PhysicalNodes{
			IPv4Address: pointer.To(v["ipv4_address"].(string)),
			Name:        pointer.To(v["name"].(string)),
		})
	}

	return &results
}

func FlattenDeploymentSettingPhysicalNode(input *[]deploymentsettings.PhysicalNodes) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	results := make([]interface{}, 0, len(*input))
	for _, item := range *input {
		results = append(results, map[string]interface{}{
			"ipv4_address": pointer.From(item.IPv4Address),
			"name":         pointer.From(item.Name),
		})
	}

	return results
}

func ExpandDeploymentSettingStorage(input []interface{}) *deploymentsettings.Storage {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	v := input[0].(map[string]interface{})

	return &deploymentsettings.Storage{
		ConfigurationMode: pointer.To(v["configuration_mode"].(string)),
	}
}

func FlattenDeploymentSettingStorage(input *deploymentsettings.Storage) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	return []interface{}{
		map[string]interface{}{
			"configuration_mode": pointer.From(input.ConfigurationMode),
		},
	}
}
