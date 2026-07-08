// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storagecache

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagecache/2025-07-01/autoimportjobs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storagecache/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

//go:generate go run ../../tools/generator-tests resourceidentity -resource-name managed_lustre_file_system_auto_import_job -service-package-name storagecache -properties "name" -compare-values "subscription_id:managed_lustre_file_system_id,resource_group_name:managed_lustre_file_system_id,aml_filesystem_name:managed_lustre_file_system_id"

type ManagedLustreFileSystemAutoImportJobResourceModel struct {
	Name                      string                                `tfschema:"name"`
	ManagedLustreFileSystemId string                                `tfschema:"managed_lustre_file_system_id"`
	Location                  string                                `tfschema:"location"`
	AdminStatus               autoimportjobs.AdminStatus            `tfschema:"admin_status"`
	AutoImportPrefixes        []string                              `tfschema:"auto_import_prefixes"`
	ConflictResolutionMode    autoimportjobs.ConflictResolutionMode `tfschema:"conflict_resolution_mode"`
	EnableDeletions           bool                                  `tfschema:"enable_deletions"`
	MaximumErrors             int64                                 `tfschema:"maximum_errors"`
	ProvisioningState         string                                `tfschema:"provisioning_state"`
	Tags                      map[string]string                     `tfschema:"tags"`
}

type ManagedLustreFileSystemAutoImportJobResource struct{}

var (
	_ sdk.ResourceWithIdentity = ManagedLustreFileSystemAutoImportJobResource{}
	_ sdk.ResourceWithUpdate   = ManagedLustreFileSystemAutoImportJobResource{}
)

func (r ManagedLustreFileSystemAutoImportJobResource) Identity() resourceids.ResourceId {
	return &autoimportjobs.AutoImportJobId{}
}

func (r ManagedLustreFileSystemAutoImportJobResource) ResourceType() string {
	return "azurerm_managed_lustre_file_system_auto_import_job"
}

func (r ManagedLustreFileSystemAutoImportJobResource) ModelObject() interface{} {
	return &ManagedLustreFileSystemAutoImportJobResourceModel{}
}

func (r ManagedLustreFileSystemAutoImportJobResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return autoimportjobs.ValidateAutoImportJobID
}

func (r ManagedLustreFileSystemAutoImportJobResource) Arguments() map[string]*pluginsdk.Schema {
	locationSchema := commonschema.Location()
	locationSchema.ForceNew = true

	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ManagedLustreFileSystemName,
		},

		"managed_lustre_file_system_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: autoimportjobs.ValidateAmlFilesystemID,
		},

		"location": locationSchema,

		"admin_status": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(autoimportjobs.AdminStatusEnable),
			ValidateFunc: validation.StringInSlice(autoimportjobs.PossibleValuesForAdminStatus(), false),
		},

		"auto_import_prefixes": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			Default: []interface{}{
				"/",
			},
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},

		"conflict_resolution_mode": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			Default:      string(autoimportjobs.ConflictResolutionModeSkip),
			ValidateFunc: validation.StringInSlice(autoimportjobs.PossibleValuesForConflictResolutionMode(), false),
		},

		"enable_deletions": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			ForceNew: true,
			Default:  false,
		},

		"maximum_errors": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.IntAtLeast(-1),
		},

		"tags": commonschema.Tags(),
	}
}

func (r ManagedLustreFileSystemAutoImportJobResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"provisioning_state": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (r ManagedLustreFileSystemAutoImportJobResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model ManagedLustreFileSystemAutoImportJobResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.StorageCache_2025_07_01.AutoImportJobs

			amlFilesystemId, err := autoimportjobs.ParseAmlFilesystemID(model.ManagedLustreFileSystemId)
			if err != nil {
				return err
			}

			id := autoimportjobs.NewAutoImportJobID(amlFilesystemId.SubscriptionId, amlFilesystemId.ResourceGroupName, amlFilesystemId.AmlFilesystemName, model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for existing %s: %+v", id, err)
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			properties := autoimportjobs.AutoImportJob{
				Location: location.Normalize(model.Location),
				Properties: &autoimportjobs.AutoImportJobProperties{
					AdminStatus:            pointer.To(model.AdminStatus),
					AutoImportPrefixes:     pointer.To(model.AutoImportPrefixes),
					ConflictResolutionMode: pointer.To(model.ConflictResolutionMode),
					EnableDeletions:        pointer.To(model.EnableDeletions),
				},
				Tags: pointer.To(model.Tags),
			}

			if !metadata.ResourceData.GetRawConfig().GetAttr("maximum_errors").IsNull() {
				properties.Properties.MaximumErrors = pointer.To(model.MaximumErrors)
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, properties); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, &id); err != nil {
				return err
			}

			return nil
		},
	}
}

func (r ManagedLustreFileSystemAutoImportJobResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.StorageCache_2025_07_01.AutoImportJobs

			id, err := autoimportjobs.ParseAutoImportJobID(metadata.ResourceData.Id())
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

			return r.flatten(metadata, id, resp.Model)
		},
	}
}

func (r ManagedLustreFileSystemAutoImportJobResource) flatten(metadata sdk.ResourceMetaData, id *autoimportjobs.AutoImportJobId, model *autoimportjobs.AutoImportJob) error {
	state := ManagedLustreFileSystemAutoImportJobResourceModel{
		Name:                      id.AutoImportJobName,
		ManagedLustreFileSystemId: autoimportjobs.NewAmlFilesystemID(id.SubscriptionId, id.ResourceGroupName, id.AmlFilesystemName).ID(),
		AdminStatus:               autoimportjobs.AdminStatusEnable,
		AutoImportPrefixes:        []string{"/"},
		ConflictResolutionMode:    autoimportjobs.ConflictResolutionModeSkip,
	}

	if model != nil {
		state.Location = location.Normalize(model.Location)
		state.Tags = pointer.From(model.Tags)

		if properties := model.Properties; properties != nil {
			if properties.AdminStatus != nil {
				state.AdminStatus = pointer.From(properties.AdminStatus)
			}
			if properties.AutoImportPrefixes != nil {
				state.AutoImportPrefixes = pointer.From(properties.AutoImportPrefixes)
			}
			if properties.ConflictResolutionMode != nil {
				state.ConflictResolutionMode = pointer.From(properties.ConflictResolutionMode)
			}
			if properties.EnableDeletions != nil {
				state.EnableDeletions = pointer.From(properties.EnableDeletions)
			}
			if properties.MaximumErrors != nil {
				state.MaximumErrors = pointer.From(properties.MaximumErrors)
			}
			if properties.ProvisioningState != nil {
				state.ProvisioningState = pointer.FromEnum(properties.ProvisioningState)
			}
		}
	}

	if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
		return err
	}

	return metadata.Encode(&state)
}

func (r ManagedLustreFileSystemAutoImportJobResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.StorageCache_2025_07_01.AutoImportJobs

			id, err := autoimportjobs.ParseAutoImportJobID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model ManagedLustreFileSystemAutoImportJobResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			properties := autoimportjobs.AutoImportJobUpdate{}

			if metadata.ResourceData.HasChange("admin_status") {
				properties.Properties = &autoimportjobs.AutoImportJobUpdateProperties{
					AdminStatus: pointer.To(model.AdminStatus),
				}
			}

			if metadata.ResourceData.HasChange("tags") {
				properties.Tags = pointer.To(model.Tags)
			}

			if properties.Properties == nil && properties.Tags == nil {
				return nil
			}

			if err := client.UpdateThenPoll(ctx, *id, properties); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r ManagedLustreFileSystemAutoImportJobResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.StorageCache_2025_07_01.AutoImportJobs

			id, err := autoimportjobs.ParseAutoImportJobID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}
