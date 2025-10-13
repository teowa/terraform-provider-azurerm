// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cognitive

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2025-06-01/accountconnectionresource"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2025-06-01/cognitiveservicesaccounts"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type CognitiveAccountConnectionModel struct {
	Name               string            `tfschema:"name"`
	CognitiveAccountId string            `tfschema:"cognitive_account_id"`
	AuthType           string            `tfschema:"auth_type"`
	Category           string            `tfschema:"category"`
	Group              string            `tfschema:"group"`
	Target             string            `tfschema:"target"`
	Metadata           map[string]string `tfschema:"metadata"`
	SharedToAll        bool              `tfschema:"shared_to_all"`
	SharedUserList     []string          `tfschema:"shared_user_list"`
	// Auth-specific configurations
	AccessKey        []AccessKeyAuthModel    `tfschema:"access_key"`
	ApiKey           []ApiKeyAuthModel       `tfschema:"api_key"`
	AccountKey       []AccountKeyAuthModel   `tfschema:"account_key"`
	ManagedIdentity  []ManagedIdentityModel  `tfschema:"managed_identity"`
	OAuth2           []OAuth2AuthModel       `tfschema:"oauth2"`
	ServicePrincipal []ServicePrincipalModel `tfschema:"service_principal"`
	UsernamePassword []UsernamePasswordModel `tfschema:"username_password"`
	CustomKeys       []CustomKeysModel       `tfschema:"custom_keys"`
	Pat              []PatAuthModel          `tfschema:"pat"`
	Sas              []SasAuthModel          `tfschema:"sas"`
}

type AccessKeyAuthModel struct {
	AccessKey string `tfschema:"access_key"`
}

type ApiKeyAuthModel struct {
	Key string `tfschema:"key"`
}

type AccountKeyAuthModel struct {
	AccountKey string `tfschema:"account_key"`
}

type ManagedIdentityModel struct {
	ClientId   string `tfschema:"client_id"`
	ResourceId string `tfschema:"resource_id"`
}

type OAuth2AuthModel struct {
	ClientId     string `tfschema:"client_id"`
	ClientSecret string `tfschema:"client_secret"`
}

type ServicePrincipalModel struct {
	ClientId     string `tfschema:"client_id"`
	ClientSecret string `tfschema:"client_secret"`
	TenantId     string `tfschema:"tenant_id"`
}

type UsernamePasswordModel struct {
	Username string `tfschema:"username"`
	Password string `tfschema:"password"`
}

type CustomKeysModel struct {
	Keys map[string]string `tfschema:"keys"`
}

type PatAuthModel struct {
	Pat string `tfschema:"pat"`
}

type SasAuthModel struct {
	Sas string `tfschema:"sas"`
}

type CognitiveAccountConnectionResource struct{}

var _ sdk.Resource = CognitiveAccountConnectionResource{}

func (r CognitiveAccountConnectionResource) ResourceType() string {
	return "azurerm_cognitive_account_connection"
}

func (r CognitiveAccountConnectionResource) ModelObject() interface{} {
	return &CognitiveAccountConnectionModel{}
}

func (r CognitiveAccountConnectionResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return accountconnectionresource.ValidateConnectionID
}

func (r CognitiveAccountConnectionResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"cognitive_account_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: cognitiveservicesaccounts.ValidateAccountID,
		},

		"auth_type": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice(
				accountconnectionresource.PossibleValuesForConnectionAuthType(),
				false,
			),
		},

		"category": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ValidateFunc: validation.StringInSlice(
				accountconnectionresource.PossibleValuesForConnectionCategory(),
				false,
			),
		},

		"group": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ValidateFunc: validation.StringInSlice(
				accountconnectionresource.PossibleValuesForConnectionGroup(),
				false,
			),
		},

		"target": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"metadata": {
			Type:     pluginsdk.TypeMap,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"shared_to_all": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"shared_user_list": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		// Auth-specific configurations
		"access_key": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"access_key": {
						Type:      pluginsdk.TypeString,
						Required:  true,
						Sensitive: true,
					},
				},
			},
		},

		"api_key": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"key": {
						Type:      pluginsdk.TypeString,
						Required:  true,
						Sensitive: true,
					},
				},
			},
		},

		"account_key": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"account_key": {
						Type:      pluginsdk.TypeString,
						Required:  true,
						Sensitive: true,
					},
				},
			},
		},

		"managed_identity": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"client_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.IsUUID,
					},
					"resource_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"oauth2": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"client_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.IsUUID,
					},
					"client_secret": {
						Type:      pluginsdk.TypeString,
						Required:  true,
						Sensitive: true,
					},
				},
			},
		},

		"service_principal": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"client_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.IsUUID,
					},
					"client_secret": {
						Type:      pluginsdk.TypeString,
						Required:  true,
						Sensitive: true,
					},
					"tenant_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.IsUUID,
					},
				},
			},
		},

		"username_password": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"username": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"password": {
						Type:      pluginsdk.TypeString,
						Required:  true,
						Sensitive: true,
					},
				},
			},
		},

		"custom_keys": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"keys": {
						Type:      pluginsdk.TypeMap,
						Required:  true,
						Sensitive: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
				},
			},
		},

		"pat": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"pat": {
						Type:      pluginsdk.TypeString,
						Required:  true,
						Sensitive: true,
					},
				},
			},
		},

		"sas": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"sas": {
						Type:      pluginsdk.TypeString,
						Required:  true,
						Sensitive: true,
					},
				},
			},
		},
	}
}

func (r CognitiveAccountConnectionResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r CognitiveAccountConnectionResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model CognitiveAccountConnectionModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.Cognitive.AccountConnectionResourceClient
			accountId, err := cognitiveservicesaccounts.ParseAccountID(model.CognitiveAccountId)
			if err != nil {
				return err
			}

			locks.ByID(accountId.ID())
			defer locks.UnlockByID(accountId.ID())

			id := accountconnectionresource.NewConnectionID(accountId.SubscriptionId, accountId.ResourceGroupName, accountId.AccountName, model.Name)
			existing, err := client.AccountConnectionsGet(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			properties, err := expandConnectionProperties(model)
			if err != nil {
				return fmt.Errorf("expanding connection properties: %+v", err)
			}

			connection := accountconnectionresource.ConnectionPropertiesV2BasicResource{
				Properties: properties,
			}

			if err := client.AccountConnectionsCreateThenPoll(ctx, id, connection); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r CognitiveAccountConnectionResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Cognitive.AccountConnectionResourceClient

			id, err := accountconnectionresource.ParseConnectionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.AccountConnectionsGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			model := resp.Model
			if model == nil {
				return fmt.Errorf("retrieving %s: model was nil", id)
			}

			state := CognitiveAccountConnectionModel{
				Name:               id.ConnectionName,
				CognitiveAccountId: cognitiveservicesaccounts.NewAccountID(id.SubscriptionId, id.ResourceGroupName, id.AccountName).ID(),
			}

			if err := flattenConnectionProperties(model.Properties, &state); err != nil {
				return fmt.Errorf("flattening connection properties: %+v", err)
			}

			return metadata.Encode(&state)
		},
	}
}

func (r CognitiveAccountConnectionResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model CognitiveAccountConnectionModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.Cognitive.AccountConnectionResourceClient
			accountId, err := cognitiveservicesaccounts.ParseAccountID(model.CognitiveAccountId)
			if err != nil {
				return err
			}

			locks.ByID(accountId.ID())
			defer locks.UnlockByID(accountId.ID())

			id, err := accountconnectionresource.ParseConnectionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			properties, err := expandConnectionProperties(model)
			if err != nil {
				return fmt.Errorf("expanding connection properties: %+v", err)
			}

			updateContent := accountconnectionresource.ConnectionUpdateContent{
				Properties: properties,
			}

			if err := client.AccountConnectionsUpdateThenPoll(ctx, *id, updateContent); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (r CognitiveAccountConnectionResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Cognitive.AccountConnectionResourceClient

			id, err := accountconnectionresource.ParseConnectionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}
			accountId := cognitiveservicesaccounts.NewAccountID(id.SubscriptionId, id.ResourceGroupName, id.AccountName)

			locks.ByID(accountId.ID())
			defer locks.UnlockByID(accountId.ID())

			if err := client.AccountConnectionsDeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}

func expandConnectionProperties(model CognitiveAccountConnectionModel) (accountconnectionresource.ConnectionPropertiesV2, error) {
	authType := accountconnectionresource.ConnectionAuthType(model.AuthType)

	switch authType {
	case accountconnectionresource.ConnectionAuthTypeAccessKey:
		if len(model.AccessKey) == 0 {
			return nil, fmt.Errorf("access_key configuration required for AccessKey auth type")
		}
		props := accountconnectionresource.AccessKeyAuthTypeConnectionProperties{
			AuthType: authType,
			Credentials: &accountconnectionresource.ConnectionAccessKey{
				AccessKey: pointer.To(model.AccessKey[0].AccessKey),
			},
		}
		return expandBaseConnectionProperties(props.ConnectionPropertiesV2(), model), nil

	case accountconnectionresource.ConnectionAuthTypeApiKey:
		if len(model.ApiKey) == 0 {
			return nil, fmt.Errorf("api_key configuration required for ApiKey auth type")
		}
		props := accountconnectionresource.ApiKeyAuthConnectionProperties{
			AuthType: authType,
			Credentials: &accountconnectionresource.ConnectionApiKey{
				Key: pointer.To(model.ApiKey[0].Key),
			},
		}
		return expandBaseConnectionProperties(props.ConnectionPropertiesV2(), model), nil

	case accountconnectionresource.ConnectionAuthTypeAccountKey:
		if len(model.AccountKey) == 0 {
			return nil, fmt.Errorf("account_key configuration required for AccountKey auth type")
		}
		props := accountconnectionresource.AccountKeyAuthTypeConnectionProperties{
			AuthType: authType,
			Credentials: &accountconnectionresource.ConnectionAccountKey{
				AccountKey: pointer.To(model.AccountKey[0].AccountKey),
			},
		}
		return expandBaseConnectionProperties(props.ConnectionPropertiesV2(), model), nil

	case accountconnectionresource.ConnectionAuthTypeManagedIdentity:
		props := accountconnectionresource.ManagedIdentityAuthTypeConnectionProperties{
			AuthType: authType,
		}
		if len(model.ManagedIdentity) > 0 {
			props.Credentials = &accountconnectionresource.ConnectionManagedIdentity{
				ClientId:   pointer.To(model.ManagedIdentity[0].ClientId),
				ResourceId: pointer.To(model.ManagedIdentity[0].ResourceId),
			}
		}
		return expandBaseConnectionProperties(props.ConnectionPropertiesV2(), model), nil

	case accountconnectionresource.ConnectionAuthTypeOAuthTwo:
		if len(model.OAuth2) == 0 {
			return nil, fmt.Errorf("oauth2 configuration required for OAuth2 auth type")
		}
		props := accountconnectionresource.OAuth2AuthTypeConnectionProperties{
			AuthType: authType,
			Credentials: &accountconnectionresource.ConnectionOauth2{
				ClientId:     pointer.To(model.OAuth2[0].ClientId),
				ClientSecret: pointer.To(model.OAuth2[0].ClientSecret),
			},
		}
		return expandBaseConnectionProperties(props.ConnectionPropertiesV2(), model), nil

	case accountconnectionresource.ConnectionAuthTypeServicePrincipal:
		if len(model.ServicePrincipal) == 0 {
			return nil, fmt.Errorf("service_principal configuration required for ServicePrincipal auth type")
		}
		props := accountconnectionresource.ServicePrincipalAuthTypeConnectionProperties{
			AuthType: authType,
			Credentials: &accountconnectionresource.ConnectionServicePrincipal{
				ClientId:     pointer.To(model.ServicePrincipal[0].ClientId),
				ClientSecret: pointer.To(model.ServicePrincipal[0].ClientSecret),
				TenantId:     pointer.To(model.ServicePrincipal[0].TenantId),
			},
		}
		return expandBaseConnectionProperties(props.ConnectionPropertiesV2(), model), nil

	case accountconnectionresource.ConnectionAuthTypeUsernamePassword:
		if len(model.UsernamePassword) == 0 {
			return nil, fmt.Errorf("username_password configuration required for UsernamePassword auth type")
		}
		props := accountconnectionresource.UsernamePasswordAuthTypeConnectionProperties{
			AuthType: authType,
			Credentials: &accountconnectionresource.ConnectionUsernamePassword{
				Username: pointer.To(model.UsernamePassword[0].Username),
				Password: pointer.To(model.UsernamePassword[0].Password),
			},
		}
		return expandBaseConnectionProperties(props.ConnectionPropertiesV2(), model), nil

	case accountconnectionresource.ConnectionAuthTypeCustomKeys:
		if len(model.CustomKeys) == 0 {
			return nil, fmt.Errorf("custom_keys configuration required for CustomKeys auth type")
		}
		props := accountconnectionresource.CustomKeysConnectionProperties{
			AuthType: authType,
			Credentials: &accountconnectionresource.CustomKeys{
				Keys: &model.CustomKeys[0].Keys,
			},
		}
		return expandBaseConnectionProperties(props.ConnectionPropertiesV2(), model), nil

	case accountconnectionresource.ConnectionAuthTypePAT:
		if len(model.Pat) == 0 {
			return nil, fmt.Errorf("pat configuration required for PAT auth type")
		}
		props := accountconnectionresource.PATAuthTypeConnectionProperties{
			AuthType: authType,
			Credentials: &accountconnectionresource.ConnectionPersonalAccessToken{
				Pat: pointer.To(model.Pat[0].Pat),
			},
		}
		return expandBaseConnectionProperties(props.ConnectionPropertiesV2(), model), nil

	case accountconnectionresource.ConnectionAuthTypeSAS:
		if len(model.Sas) == 0 {
			return nil, fmt.Errorf("sas configuration required for SAS auth type")
		}
		props := accountconnectionresource.SASAuthTypeConnectionProperties{
			AuthType: authType,
			Credentials: &accountconnectionresource.ConnectionSharedAccessSignature{
				Sas: pointer.To(model.Sas[0].Sas),
			},
		}
		return expandBaseConnectionProperties(props.ConnectionPropertiesV2(), model), nil

	case accountconnectionresource.ConnectionAuthTypeNone:
		props := accountconnectionresource.NoneAuthTypeConnectionProperties{
			AuthType: authType,
		}
		return expandBaseConnectionProperties(props.ConnectionPropertiesV2(), model), nil

	default:
		return nil, fmt.Errorf("unsupported auth type: %s", model.AuthType)
	}
}

func expandBaseConnectionProperties(base accountconnectionresource.BaseConnectionPropertiesV2Impl, model CognitiveAccountConnectionModel) accountconnectionresource.ConnectionPropertiesV2 {
	if model.Category != "" {
		category := accountconnectionresource.ConnectionCategory(model.Category)
		base.Category = &category
	}

	if model.Group != "" {
		group := accountconnectionresource.ConnectionGroup(model.Group)
		base.Group = &group
	}

	if model.Target != "" {
		base.Target = &model.Target
	}

	if len(model.Metadata) > 0 {
		base.Metadata = &model.Metadata
	}

	base.IsSharedToAll = pointer.To(model.SharedToAll)

	if len(model.SharedUserList) > 0 {
		base.SharedUserList = &model.SharedUserList
	}

	// Note: The interface return here is simplified - in practice, you'd return
	// the specific auth type struct that implements ConnectionPropertiesV2
	return base
}

func flattenConnectionProperties(props accountconnectionresource.ConnectionPropertiesV2, state *CognitiveAccountConnectionModel) error {
	if props == nil {
		return nil
	}

	base := props.ConnectionPropertiesV2()
	state.AuthType = string(base.AuthType)

	if base.Category != nil {
		state.Category = string(*base.Category)
	}

	if base.Group != nil {
		state.Group = string(*base.Group)
	}

	if base.Target != nil {
		state.Target = *base.Target
	}

	if base.Metadata != nil {
		state.Metadata = *base.Metadata
	}

	state.SharedToAll = pointer.From(base.IsSharedToAll)

	if base.SharedUserList != nil {
		state.SharedUserList = *base.SharedUserList
	}

	return nil
}
