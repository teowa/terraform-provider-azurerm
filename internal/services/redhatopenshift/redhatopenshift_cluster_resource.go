package redhatopenshift

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redhatopenshift/2023-09-04/openshiftclusters"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	openShiftValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/redhatopenshift/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type RedHatOpenShiftCluster struct {
}

var _ sdk.ResourceWithUpdate = RedHatOpenShiftCluster{}

type RedHatOpenShiftClusterModel struct {
	Name             string             `tfschema:"name"`
	Location         string             `tfschema:"location"`
	ResourceGroup    string             `tfschema:"resource_group_name"`
	Version          string             `tfschema:"version"`
	ConsoleUrl       string             `tfschema:"console_url"`
	ServicePrincipal []servicePrincipal `tfschema:"service_principal"`
	ClusterProfile   []clusterProfile   `tfschema:"cluster_profile"`
	NetworkProfile   []networkProfile   `tfschema:"network_profile"`
	MainProfile      []mainProfile      `tfschema:"main_profile"`
	WorkerProfile    []workerProfile    `tfschema:"worker_profile"`
	ApiServerProfile []apiServerProfile `tfschema:"api_server_profile"`
	IngressProfile   []ingressProfile   `tfschema:"ingress_profile"`
	Tags             map[string]string  `tfschema:"tags"`
}

type servicePrincipal struct {
	ClientId     string `tfschema:"client_id"`
	ClientSecret string `tfschema:"client_secret"`
}

type clusterProfile struct {
	PullSecret  string `tfschema:"pull_secret"`
	Domain      string `tfschema:"domain"`
	FipsEnabled bool   `tfschema:"fips_enabled"`
}

type networkProfile struct {
	PodCidr     string `tfschema:"pod_cidr"`
	ServiceCidr string `tfschema:"service_cidr"`
}

type mainProfile struct {
	SubnetId                string `tfschema:"subnet_id"`
	VmSize                  string `tfschema:"vm_size"`
	EncryptionAtHostEnabled bool   `tfschema:"encryption_at_host_enabled"`
	DiskEncryptionSetId     string `tfschema:"disk_encryption_set_id"`
}

type workerProfile struct {
	VmSize                  string `tfschema:"vm_size"`
	DiskSizeGb              int32  `tfschema:"disk_size_gb"`
	NodeCount               int32  `tfschema:"node_count"`
	SubnetId                string `tfschema:"subnet_id"`
	EncryptionAtHostEnabled bool   `tfschema:"encryption_at_host_enabled"`
	DiskEncryptionSetId     string `tfschema:"disk_encryption_set_id"`
}

type ingressProfile struct {
	Visibility string `tfschema:"visibility"`
	Ip         string `tfschema:"ip"`
}

type apiServerProfile struct {
	Visibility string `tfschema:"visibility"`
	Url        string `tfschema:"url"`
}

func (r RedHatOpenShiftCluster) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"location": commonschema.Location(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"main_profile": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"subnet_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: azure.ValidateResourceID,
					},
					"vm_size": {
						Type:             pluginsdk.TypeString,
						Required:         true,
						DiffSuppressFunc: suppress.CaseDifference,
						ValidateFunc:     validation.StringIsNotEmpty,
					},
					"encryption_at_host_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},
					"disk_encryption_set_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: azure.ValidateResourceID,
					},
				},
			},
		},

		"service_principal": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"client_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.IsUUID,
					},
					"client_secret": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						Sensitive:    true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"api_server_profile": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"visibility": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(openshiftclusters.VisibilityPublic),
							string(openshiftclusters.VisibilityPrivate),
						}, false),
					},
					"url": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"cluster_profile": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"pull_secret": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"domain": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Computed:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"fips_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},
				},
			},
		},

		"ingress_profile": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			Computed: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"visibility": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(openshiftclusters.VisibilityPublic),
							string(openshiftclusters.VisibilityPrivate),
						}, false),
					},
					"ip": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"network_profile": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"pod_cidr": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Computed:     true,
						ForceNew:     true,
						ValidateFunc: validate.CIDR,
					},
					"service_cidr": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Computed:     true,
						ForceNew:     true,
						ValidateFunc: validate.CIDR,
					},
				},
			},
		},

		"worker_profile": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"vm_size": {
						Type:             pluginsdk.TypeString,
						Required:         true,
						DiffSuppressFunc: suppress.CaseDifference,
						ValidateFunc:     validation.StringIsNotEmpty,
					},
					"disk_size_gb": {
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: openShiftValidate.ValidateDiskSizeGB,
					},
					"node_count": {
						Type:     pluginsdk.TypeInt,
						Required: true,
					},
					"subnet_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: azure.ValidateResourceID,
					},

					"encryption_at_host_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},
					"disk_encryption_set_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: azure.ValidateResourceID,
					},
				},
			},
		},

		"tags": tags.Schema(),
	}
}

func (r RedHatOpenShiftCluster) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"version": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"console_url": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (r RedHatOpenShiftCluster) ModelObject() interface{} {
	return &RedHatOpenShiftClusterModel{}
}

func (r RedHatOpenShiftCluster) ResourceType() string {
	return "azurerm_redhat_openshift_cluster"
}

func (r RedHatOpenShiftCluster) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return openShiftValidate.ValidateClusterID
}

func (r RedHatOpenShiftCluster) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model RedHatOpenShiftClusterModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			client := metadata.Client.RedHatOpenshift.Client.OpenShiftClusters
			subscriptionId := metadata.Client.Account.SubscriptionId

			resourceGroupName := model.ResourceGroup
			name := model.Name
			id := openshiftclusters.NewProviderOpenShiftClusterID(subscriptionId, resourceGroupName, name)
			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}

			if existing.Model != nil && existing.Model.Id != nil {
				return tf.ImportAsExistsError("azurerm_redhat_openshift_cluster", *existing.Model.Id)
			}

			location := azure.NormalizeLocation(model.Location)

			parameters := openshiftclusters.OpenShiftCluster{
				Name:     &name,
				Location: location,
				Properties: &openshiftclusters.OpenShiftClusterProperties{
					ClusterProfile:          expandOpenshiftClusterProfile(model.ClusterProfile, subscriptionId),
					ConsoleProfile:          &openshiftclusters.ConsoleProfile{},
					ServicePrincipalProfile: expandOpenshiftServicePrincipalProfile(model.ServicePrincipal),
					NetworkProfile:          expandOpenshiftNetworkProfile(model.NetworkProfile),
					MasterProfile:           expandOpenshiftMasterProfile(model.MainProfile),
					WorkerProfiles:          expandOpenshiftWorkerProfiles(model.WorkerProfile),
					ApiserverProfile:        expandOpenshiftApiServerProfile(model.ApiServerProfile),
					IngressProfiles:         expandOpenshiftIngressProfiles(model.IngressProfile),
				},
				Tags: pointer.To(model.Tags),
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, parameters); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
		Timeout: 90 * time.Minute,
	}
}

func (r RedHatOpenShiftCluster) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.RedHatOpenshift.Client.OpenShiftClusters
			subscriptionId := metadata.Client.Account.SubscriptionId
			id, err := openshiftclusters.ParseProviderOpenShiftClusterID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var state RedHatOpenShiftClusterModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			parameter := openshiftclusters.OpenShiftClusterUpdate{
				Tags: pointer.To(state.Tags),
			}

			if metadata.ResourceData.HasChange("cluster_profile") {
				parameter.Properties.ClusterProfile = expandOpenshiftClusterProfile(state.ClusterProfile, subscriptionId)
			}

			if metadata.ResourceData.HasChange("main_profile") {
				parameter.Properties.MasterProfile = expandOpenshiftMasterProfile(state.MainProfile)
			}

			if metadata.ResourceData.HasChange("worker_profile") {
				parameter.Properties.WorkerProfiles = expandOpenshiftWorkerProfiles(state.WorkerProfile)
			}

			if err := client.UpdateThenPoll(ctx, *id, parameter); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},

		Timeout: 90 * time.Minute,
	}
}

func (r RedHatOpenShiftCluster) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.RedHatOpenshift.Client.OpenShiftClusters

			id, err := openshiftclusters.ParseProviderOpenShiftClusterID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := RedHatOpenShiftClusterModel{
				Name:          id.OpenShiftClusterName,
				ResourceGroup: id.ResourceGroupName,
			}

			if model := resp.Model; model != nil {
				state.Tags = pointer.From(model.Tags)
				state.Location = azure.NormalizeLocation(model.Location)

				if props := model.Properties; props != nil {
					state.ClusterProfile = flattenOpenShiftClusterProfile(props.ClusterProfile)
					state.ServicePrincipal = flattenOpenShiftServicePrincipalProfile(props.ServicePrincipalProfile, metadata)
					state.NetworkProfile = flattenOpenShiftNetworkProfile(props.NetworkProfile)
					state.MainProfile = flattenOpenShiftMasterProfile(props.MasterProfile)
					state.WorkerProfile = flattenOpenShiftWorkerProfiles(props.WorkerProfiles)
					state.ApiServerProfile = flattenOpenShiftAPIServerProfile(props.ApiserverProfile)
					state.IngressProfile = flattenOpenShiftIngressProfiles(props.IngressProfiles)

					if props.ClusterProfile != nil && props.ClusterProfile.Version != nil {
						state.Version = *props.ClusterProfile.Version
					}

					if props.ConsoleProfile != nil && props.ConsoleProfile.Url != nil {
						state.ConsoleUrl = *props.ConsoleProfile.Url
					}
				}
			}

			return metadata.Encode(&state)
		},
		Timeout: 5 * time.Minute,
	}
}

func (r RedHatOpenShiftCluster) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := openshiftclusters.ParseProviderOpenShiftClusterID(metadata.ResourceData.Id())

			if err != nil {
				return fmt.Errorf("while parsing resource ID: %+v", err)
			}

			client := metadata.Client.RedHatOpenshift.Client.OpenShiftClusters

			if err := client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
		Timeout: 90 * time.Minute,
	}
}

func flattenOpenShiftClusterProfile(profile *openshiftclusters.ClusterProfile) []clusterProfile {
	if profile == nil {
		return []clusterProfile{}
	}

	pullSecret := ""
	if profile.PullSecret != nil {
		pullSecret = *profile.PullSecret
	}

	clusterDomain := ""
	if profile.Domain != nil {
		clusterDomain = *profile.Domain
	}

	fipsEnabled := profile.FipsValidatedModules == string(openshiftclusters.FipsValidatedModulesEnabled)

	return []clusterProfile{
		{
			PullSecret:  pullSecret,
			Domain:      clusterDomain,
			FipsEnabled: fipsEnabled,
		},
	}

}

func flattenOpenShiftServicePrincipalProfile(profile *redhatopenshift.ServicePrincipalProfile, metadata sdk.ResourceMetaData) []servicePrincipal {
	if profile == nil {
		return []servicePrincipal{}
	}

	clientID := ""
	if profile.ClientID != nil {
		clientID = *profile.ClientID
	}

	clientSecret := ""
	if sp, ok := metadata.ResourceData.GetOk("service_principal"); ok {
		val := sp.([]interface{})

		if len(val) > 0 && val[0] != nil {
			raw := val[0].(map[string]interface{})
			clientSecret = raw["client_secret"].(string)
		}
	}

	return []servicePrincipal{
		{
			ClientId:     clientID,
			ClientSecret: clientSecret,
		},
	}
}

func flattenOpenShiftNetworkProfile(profile *redhatopenshift.NetworkProfile) []networkProfile {
	if profile == nil {
		return []networkProfile{}
	}

	podCidr := ""
	if profile.PodCidr != nil {
		podCidr = *profile.PodCidr
	}

	serviceCidr := ""
	if profile.ServiceCidr != nil {
		serviceCidr = *profile.ServiceCidr
	}

	return []networkProfile{
		{
			PodCidr:     podCidr,
			ServiceCidr: serviceCidr,
		},
	}
}

func flattenOpenShiftMasterProfile(profile *redhatopenshift.MasterProfile) []mainProfile {
	if profile == nil {
		return []mainProfile{}
	}

	mainProfiles := make([]mainProfile, 0)

	subnetId := ""
	if profile.SubnetID != nil {
		subnetId = *profile.SubnetID
	}

	encryptionAtHostEnabled := profile.EncryptionAtHost == redhatopenshift.EncryptionAtHostEnabled

	diskEncryptionSetId := ""
	if profile.DiskEncryptionSetID != nil {
		diskEncryptionSetId = *profile.DiskEncryptionSetID
	}

	profileRaw := mainProfile{
		SubnetId:                subnetId,
		EncryptionAtHostEnabled: encryptionAtHostEnabled,
		DiskEncryptionSetId:     diskEncryptionSetId,
	}

	if profile.VMSize != nil {
		profileRaw.VmSize = *profile.VMSize
	}

	mainProfiles = append(mainProfiles, profileRaw)
	return mainProfiles
}

func flattenOpenShiftWorkerProfiles(profiles *[]redhatopenshift.WorkerProfile) []workerProfile {
	if profiles == nil {
		return []workerProfile{}
	}

	results := make([]workerProfile, 0)

	for _, profile := range *profiles {

		worker := workerProfile{
			EncryptionAtHostEnabled: profile.EncryptionAtHost == redhatopenshift.EncryptionAtHostEnabled,
		}
		if profile.DiskSizeGB != nil {
			worker.DiskSizeGb = *profile.DiskSizeGB
		}
		if profile.Count != nil {
			worker.NodeCount = *profile.Count
		}
		if profile.VMSize != nil {
			worker.VmSize = *profile.VMSize
		}
		if profile.SubnetID != nil {
			worker.SubnetId = *profile.SubnetID
		}
		if profile.DiskEncryptionSetID != nil {
			worker.DiskEncryptionSetId = *profile.DiskEncryptionSetID
		}
		results = append(results, worker)
	}
	return results
}

func flattenOpenShiftAPIServerProfile(profile *redhatopenshift.APIServerProfile) []apiServerProfile {
	if profile == nil {
		return []apiServerProfile{}
	}

	results := make([]apiServerProfile, 0)

	apiServerProfile := apiServerProfile{
		Visibility: string(profile.Visibility),
	}

	if profile.URL != nil {
		apiServerProfile.Url = *profile.URL
	}

	results = append(results, apiServerProfile)
	return results
}

func flattenOpenShiftIngressProfiles(profiles *[]redhatopenshift.IngressProfile) []ingressProfile {
	if profiles == nil {
		return []ingressProfile{}
	}

	results := make([]ingressProfile, 0)

	for _, profile := range *profiles {
		ingress := ingressProfile{
			Visibility: string(profile.Visibility),
		}
		if profile.IP != nil {
			ingress.Ip = *profile.IP
		}
		results = append(results, ingress)
	}

	return results
}

func expandOpenshiftClusterProfile(input []clusterProfile, subscriptionId string) *redhatopenshift.ClusterProfile {
	resourceGroupName := fmt.Sprintf("aro-%s", GenerateRandomDomainName())
	resourceGroupId := ResourceGroupID(subscriptionId, resourceGroupName)

	if len(input) == 0 {
		return &redhatopenshift.ClusterProfile{
			ResourceGroupID:      utils.String(resourceGroupId),
			Domain:               utils.String(GenerateRandomDomainName()),
			FipsValidatedModules: redhatopenshift.FipsValidatedModulesDisabled,
		}
	}

	config := input[0]
	domain := config.Domain
	if domain == "" {
		domain = GenerateRandomDomainName()
	}

	fipsValidatedModules := redhatopenshift.FipsValidatedModulesDisabled
	fipsEnabled := config.FipsEnabled
	if fipsEnabled {
		fipsValidatedModules = redhatopenshift.FipsValidatedModulesEnabled
	}

	return &redhatopenshift.ClusterProfile{
		ResourceGroupID:      utils.String(resourceGroupId),
		Domain:               utils.String(domain),
		PullSecret:           utils.String(config.PullSecret),
		FipsValidatedModules: fipsValidatedModules,
	}
}

func expandOpenshiftServicePrincipalProfile(input []servicePrincipal) *redhatopenshift.ServicePrincipalProfile {
	if len(input) == 0 {
		return nil
	}

	config := input[0]

	return &redhatopenshift.ServicePrincipalProfile{
		ClientID:     utils.String(config.ClientId),
		ClientSecret: utils.String(config.ClientSecret),
	}
}

func expandOpenshiftNetworkProfile(input []networkProfile) *redhatopenshift.NetworkProfile {
	if len(input) == 0 {
		return &redhatopenshift.NetworkProfile{
			PodCidr:     utils.String("10.128.0.0/14"),
			ServiceCidr: utils.String("172.30.0.0/16"),
		}
	}

	config := input[0]

	return &redhatopenshift.NetworkProfile{
		PodCidr:     utils.String(config.PodCidr),
		ServiceCidr: utils.String(config.ServiceCidr),
	}
}

func expandOpenshiftMasterProfile(input []mainProfile) *redhatopenshift.MasterProfile {
	if len(input) == 0 {
		return nil
	}

	config := input[0]

	encryptionAtHost := redhatopenshift.EncryptionAtHostDisabled
	enableEncryptionAtHost := config.EncryptionAtHostEnabled
	if enableEncryptionAtHost {
		encryptionAtHost = redhatopenshift.EncryptionAtHostEnabled
	}

	return &redhatopenshift.MasterProfile{
		VMSize:              utils.String(config.VmSize),
		SubnetID:            utils.String(config.SubnetId),
		EncryptionAtHost:    encryptionAtHost,
		DiskEncryptionSetID: utils.String(config.DiskEncryptionSetId),
	}
}

func expandOpenshiftWorkerProfiles(inputs []workerProfile) *[]redhatopenshift.WorkerProfile {
	if len(inputs) == 0 {
		return nil
	}

	profiles := make([]redhatopenshift.WorkerProfile, 0)
	config := inputs[0]

	encryptionAtHost := redhatopenshift.EncryptionAtHostDisabled
	enableEncryptionAtHost := config.EncryptionAtHostEnabled
	if enableEncryptionAtHost {
		encryptionAtHost = redhatopenshift.EncryptionAtHostEnabled
	}

	profile := redhatopenshift.WorkerProfile{
		Name:                utils.String("worker"),
		VMSize:              utils.String(config.VmSize),
		DiskSizeGB:          utils.Int32(config.DiskSizeGb),
		SubnetID:            utils.String(config.SubnetId),
		Count:               utils.Int32(config.NodeCount),
		EncryptionAtHost:    encryptionAtHost,
		DiskEncryptionSetID: utils.String(config.DiskEncryptionSetId),
	}

	profiles = append(profiles, profile)

	return &profiles
}

func expandOpenshiftApiServerProfile(input []apiServerProfile) *redhatopenshift.APIServerProfile {
	if len(input) == 0 {
		return &redhatopenshift.APIServerProfile{
			Visibility: redhatopenshift.VisibilityPublic,
		}
	}

	config := input[0]

	return &redhatopenshift.APIServerProfile{
		Visibility: redhatopenshift.Visibility(config.Visibility),
	}
}

func expandOpenshiftIngressProfiles(inputs []ingressProfile) *[]redhatopenshift.IngressProfile {
	profiles := make([]redhatopenshift.IngressProfile, 0)

	name := "default"
	visibility := string(redhatopenshift.VisibilityPublic)

	if len(inputs) > 0 {
		visibility = inputs[0].Visibility
	}

	profile := redhatopenshift.IngressProfile{
		Name:       utils.String(name),
		Visibility: redhatopenshift.Visibility(visibility),
	}

	profiles = append(profiles, profile)

	return &profiles
}
