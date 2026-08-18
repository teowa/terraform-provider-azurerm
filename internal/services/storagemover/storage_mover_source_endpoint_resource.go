// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package storagemover

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagemover/2025-07-01/endpoints"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagemover/2025-07-01/storagemovers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

//go:generate go run ../../tools/generator-tests resourceidentity -resource-name storage_mover_source_endpoint -service-package-name storagemover -properties "name" -compare-values "subscription_id:storage_mover_id,resource_group_name:storage_mover_id,storage_mover_name:storage_mover_id"

type StorageMoverSourceEndpointModel struct {
	Name                     string                                  `tfschema:"name"`
	StorageMoverId           string                                  `tfschema:"storage_mover_id"`
	Export                   string                                  `tfschema:"export"`
	Host                     string                                  `tfschema:"host"`
	NfsVersion               endpoints.NfsVersion                    `tfschema:"nfs_version"`
	AzureMultiCloudConnector []AzureMultiCloudConnectorEndpointModel `tfschema:"azure_multi_cloud_connector"`
	AzureStorageNfsFileShare []AzureStorageFileShareEndpointModel    `tfschema:"azure_storage_nfs_file_share"`
	AzureStorageSmbFileShare []AzureStorageFileShareEndpointModel    `tfschema:"azure_storage_smb_file_share"`
	SmbMount                 []SmbMountEndpointModel                 `tfschema:"smb_mount"`
	Identity                 []identity.ModelSystemAssigned          `tfschema:"identity"`
	Description              string                                  `tfschema:"description"`
}

type AzureMultiCloudConnectorEndpointModel struct {
	AwsS3BucketId         string `tfschema:"aws_s3_bucket_id"`
	MultiCloudConnectorId string `tfschema:"multi_cloud_connector_id"`
}

type AzureStorageFileShareEndpointModel struct {
	FileShareName            string `tfschema:"file_share_name"`
	StorageAccountResourceId string `tfschema:"storage_account_id"`
}

type SmbMountEndpointModel struct {
	Host        string                        `tfschema:"host"`
	ShareName   string                        `tfschema:"share_name"`
	Credentials []AzureKeyVaultSmbCredentials `tfschema:"credentials"`
}

type AzureKeyVaultSmbCredentials struct {
	PasswordUri string `tfschema:"password_uri"`
	UsernameUri string `tfschema:"username_uri"`
}

type StorageMoverSourceEndpointResource struct{}

var (
	_ sdk.ResourceWithIdentity = StorageMoverSourceEndpointResource{}
	_ sdk.ResourceWithUpdate   = StorageMoverSourceEndpointResource{}
)

func (r StorageMoverSourceEndpointResource) Identity() resourceids.ResourceId {
	return &endpoints.EndpointId{}
}

func (r StorageMoverSourceEndpointResource) ResourceType() string {
	return "azurerm_storage_mover_source_endpoint"
}

func (r StorageMoverSourceEndpointResource) ModelObject() interface{} {
	return &StorageMoverSourceEndpointModel{}
}

func (r StorageMoverSourceEndpointResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return endpoints.ValidateEndpointID
}

var sourceEndpointTypeFieldNames = []string{
	"host",
	"azure_multi_cloud_connector",
	"azure_storage_nfs_file_share",
	"azure_storage_smb_file_share",
	"smb_mount",
}

func (r StorageMoverSourceEndpointResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile(`^[0-9a-zA-Z][-_0-9a-zA-Z]{0,63}$`),
				`The name must be between 1 and 64 characters in length, begin with a letter or number, and may contain letters, numbers, dashes and underscore.`,
			),
		},

		"storage_mover_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: storagemovers.ValidateStorageMoverID,
		},

		"host": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			ExactlyOneOf: sourceEndpointTypeFieldNames,
		},

		"export": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"nfs_version": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
			Default:  string(endpoints.NfsVersionNFSauto),
			ValidateFunc: validation.StringInSlice([]string{
				string(endpoints.NfsVersionNFSauto),
				string(endpoints.NfsVersionNFSvFour),
				string(endpoints.NfsVersionNFSvThree),
			}, false),
		},

		"azure_multi_cloud_connector": {
			Type:         pluginsdk.TypeList,
			Optional:     true,
			ForceNew:     true,
			MaxItems:     1,
			ExactlyOneOf: sourceEndpointTypeFieldNames,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"aws_s3_bucket_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"multi_cloud_connector_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"azure_storage_nfs_file_share": {
			Type:         pluginsdk.TypeList,
			Optional:     true,
			ForceNew:     true,
			MaxItems:     1,
			ExactlyOneOf: sourceEndpointTypeFieldNames,
			Elem: &pluginsdk.Resource{
				Schema: azureStorageFileShareEndpointSchema(),
			},
		},

		"azure_storage_smb_file_share": {
			Type:         pluginsdk.TypeList,
			Optional:     true,
			ForceNew:     true,
			MaxItems:     1,
			ExactlyOneOf: sourceEndpointTypeFieldNames,
			Elem: &pluginsdk.Resource{
				Schema: azureStorageFileShareEndpointSchema(),
			},
		},

		"smb_mount": {
			Type:         pluginsdk.TypeList,
			Optional:     true,
			ForceNew:     true,
			MaxItems:     1,
			ExactlyOneOf: sourceEndpointTypeFieldNames,
			Elem: &pluginsdk.Resource{
				Schema: smbMountEndpointSchema(),
			},
		},

		"identity": commonschema.SystemAssignedIdentityOptional(),

		"description": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func azureStorageFileShareEndpointSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"file_share_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"storage_account_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func smbMountEndpointSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"host": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"share_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"credentials": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"password_uri": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
						AtLeastOneOf: []string{"smb_mount.0.credentials.0.password_uri", "smb_mount.0.credentials.0.username_uri"},
					},

					"username_uri": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
						AtLeastOneOf: []string{"smb_mount.0.credentials.0.password_uri", "smb_mount.0.credentials.0.username_uri"},
					},
				},
			},
		},
	}
}

func (r StorageMoverSourceEndpointResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r StorageMoverSourceEndpointResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model StorageMoverSourceEndpointModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.StorageMover.EndpointsClient
			storageMoverId, err := storagemovers.ParseStorageMoverID(model.StorageMoverId)
			if err != nil {
				return err
			}

			id := endpoints.NewEndpointID(storageMoverId.SubscriptionId, storageMoverId.ResourceGroupName, storageMoverId.StorageMoverName, model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for existing %s: %+v", id, err)
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			properties, err := expandSourceEndpointProperties(model)
			if err != nil {
				return err
			}

			payload := endpoints.Endpoint{
				Properties: properties,
			}

			if len(model.Identity) > 0 {
				identityValue, err := identity.ExpandSystemAssignedFromModel(model.Identity)
				if err != nil {
					return fmt.Errorf("expanding `identity`: %+v", err)
				}

				payload.Identity = &identity.LegacySystemAndUserAssignedMap{
					Type: identityValue.Type,
				}
			}

			if _, err := client.CreateOrUpdate(ctx, id, payload); err != nil {
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

func (r StorageMoverSourceEndpointResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.StorageMover.EndpointsClient

			id, err := endpoints.ParseEndpointID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model StorageMoverSourceEndpointModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			properties := resp.Model
			if properties == nil {
				return fmt.Errorf("retrieving %s: model was nil", *id)
			}

			if metadata.ResourceData.HasChange("description") {
				switch v := properties.Properties.(type) {
				case endpoints.NfsMountEndpointProperties:
					v.Description = pointer.To(model.Description)
					properties.Properties = v
				case endpoints.AzureMultiCloudConnectorEndpointProperties:
					v.Description = pointer.To(model.Description)
					properties.Properties = v
				case endpoints.AzureStorageNfsFileShareEndpointProperties:
					v.Description = pointer.To(model.Description)
					properties.Properties = v
				case endpoints.AzureStorageSmbFileShareEndpointProperties:
					v.Description = pointer.To(model.Description)
					properties.Properties = v
				case endpoints.SmbMountEndpointProperties:
					v.Description = pointer.To(model.Description)
					properties.Properties = v
				}
			}

			if metadata.ResourceData.HasChange("smb_mount.0.credentials") {
				if v, ok := properties.Properties.(endpoints.SmbMountEndpointProperties); ok {
					v.Credentials = expandSmbMountCredentials(model.SmbMount)
					properties.Properties = v
				}
			}

			if metadata.ResourceData.HasChange("identity") {
				if len(model.Identity) > 0 {
					identityValue, err := identity.ExpandSystemAssignedFromModel(model.Identity)
					if err != nil {
						return fmt.Errorf("expanding `identity`: %+v", err)
					}

					properties.Identity = &identity.LegacySystemAndUserAssignedMap{
						Type: identityValue.Type,
					}
				} else {
					properties.Identity = &identity.LegacySystemAndUserAssignedMap{
						Type: identity.TypeNone,
					}
				}
			}

			if _, err := client.CreateOrUpdate(ctx, *id, *properties); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r StorageMoverSourceEndpointResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.StorageMover.EndpointsClient

			id, err := endpoints.ParseEndpointID(metadata.ResourceData.Id())
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

func (r StorageMoverSourceEndpointResource) flatten(metadata sdk.ResourceMetaData, id *endpoints.EndpointId, model *endpoints.Endpoint) error {
	state := StorageMoverSourceEndpointModel{
		Name:           id.EndpointName,
		StorageMoverId: storagemovers.NewStorageMoverID(id.SubscriptionId, id.ResourceGroupName, id.StorageMoverName).ID(),
	}

	if model != nil {
		switch v := model.Properties.(type) {
		case endpoints.NfsMountEndpointProperties:
			state.Export = v.Export
			state.Host = v.Host

			if nfsVersion := v.NfsVersion; nfsVersion != nil {
				state.NfsVersion = *nfsVersion
			}

			state.Description = pointer.From(v.Description)

		case endpoints.AzureMultiCloudConnectorEndpointProperties:
			state.AzureMultiCloudConnector = []AzureMultiCloudConnectorEndpointModel{
				{
					AwsS3BucketId:         v.AwsS3BucketId,
					MultiCloudConnectorId: v.MultiCloudConnectorId,
				},
			}
			state.Description = pointer.From(v.Description)

		case endpoints.AzureStorageNfsFileShareEndpointProperties:
			state.AzureStorageNfsFileShare = []AzureStorageFileShareEndpointModel{
				{
					FileShareName:            v.FileShareName,
					StorageAccountResourceId: v.StorageAccountResourceId,
				},
			}
			state.Description = pointer.From(v.Description)

		case endpoints.AzureStorageSmbFileShareEndpointProperties:
			state.AzureStorageSmbFileShare = []AzureStorageFileShareEndpointModel{
				{
					FileShareName:            v.FileShareName,
					StorageAccountResourceId: v.StorageAccountResourceId,
				},
			}
			state.Description = pointer.From(v.Description)

		case endpoints.SmbMountEndpointProperties:
			smbMount := SmbMountEndpointModel{
				Host:      v.Host,
				ShareName: v.ShareName,
			}

			if v.Credentials != nil {
				smbMount.Credentials = []AzureKeyVaultSmbCredentials{
					{
						PasswordUri: pointer.From(v.Credentials.PasswordUri),
						UsernameUri: pointer.From(v.Credentials.UsernameUri),
					},
				}
			}

			state.SmbMount = []SmbMountEndpointModel{smbMount}
			state.Description = pointer.From(v.Description)
		}
	}

	if model != nil && model.Identity != nil {
		systemAssigned := &identity.SystemAssigned{
			Type:        model.Identity.Type,
			PrincipalId: model.Identity.PrincipalId,
			TenantId:    model.Identity.TenantId,
		}
		state.Identity = identity.FlattenSystemAssignedToModel(systemAssigned)
	}

	if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
		return err
	}

	return metadata.Encode(&state)
}

func (r StorageMoverSourceEndpointResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.StorageMover.EndpointsClient

			id, err := endpoints.ParseEndpointID(metadata.ResourceData.Id())
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

func expandSourceEndpointProperties(model StorageMoverSourceEndpointModel) (endpoints.EndpointBaseProperties, error) {
	var properties endpoints.EndpointBaseProperties

	switch {
	case model.Host != "":
		nfsVersion := model.NfsVersion
		properties = endpoints.NfsMountEndpointProperties{
			Export:     model.Export,
			Host:       model.Host,
			NfsVersion: &nfsVersion,
		}

	case len(model.AzureMultiCloudConnector) > 0:
		v := model.AzureMultiCloudConnector[0]
		properties = endpoints.AzureMultiCloudConnectorEndpointProperties{
			AwsS3BucketId:         v.AwsS3BucketId,
			MultiCloudConnectorId: v.MultiCloudConnectorId,
		}

	case len(model.AzureStorageNfsFileShare) > 0:
		v := model.AzureStorageNfsFileShare[0]
		properties = endpoints.AzureStorageNfsFileShareEndpointProperties{
			FileShareName:            v.FileShareName,
			StorageAccountResourceId: v.StorageAccountResourceId,
		}

	case len(model.AzureStorageSmbFileShare) > 0:
		v := model.AzureStorageSmbFileShare[0]
		properties = endpoints.AzureStorageSmbFileShareEndpointProperties{
			FileShareName:            v.FileShareName,
			StorageAccountResourceId: v.StorageAccountResourceId,
		}

	case len(model.SmbMount) > 0:
		v := model.SmbMount[0]
		properties = endpoints.SmbMountEndpointProperties{
			Host:        v.Host,
			ShareName:   v.ShareName,
			Credentials: expandSmbMountCredentials(model.SmbMount),
		}

	default:
		return nil, errors.New("one of `host`, `azure_multi_cloud_connector`, `azure_storage_nfs_file_share`, `azure_storage_smb_file_share` or `smb_mount` must be specified")
	}

	if model.Description != "" {
		switch v := properties.(type) {
		case endpoints.NfsMountEndpointProperties:
			v.Description = pointer.To(model.Description)
			properties = v
		case endpoints.AzureMultiCloudConnectorEndpointProperties:
			v.Description = pointer.To(model.Description)
			properties = v
		case endpoints.AzureStorageNfsFileShareEndpointProperties:
			v.Description = pointer.To(model.Description)
			properties = v
		case endpoints.AzureStorageSmbFileShareEndpointProperties:
			v.Description = pointer.To(model.Description)
			properties = v
		case endpoints.SmbMountEndpointProperties:
			v.Description = pointer.To(model.Description)
			properties = v
		}
	}

	return properties, nil
}

func expandSmbMountCredentials(input []SmbMountEndpointModel) *endpoints.AzureKeyVaultSmbCredentials {
	if len(input) == 0 || len(input[0].Credentials) == 0 {
		return nil
	}

	v := input[0].Credentials[0]
	return &endpoints.AzureKeyVaultSmbCredentials{
		PasswordUri: pointer.To(v.PasswordUri),
		UsernameUri: pointer.To(v.UsernameUri),
	}
}
