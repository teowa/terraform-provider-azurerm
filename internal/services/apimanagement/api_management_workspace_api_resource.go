// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package apimanagement

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/api"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/workspace"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/schemaz"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type ApiManagementWorkspaceApiModel struct {
	Name                          string                                                   `tfschema:"name"`
	ApiManagementWorkspaceId      string                                                   `tfschema:"api_management_workspace_id"`
	DisplayName                   string                                                   `tfschema:"display_name"`
	Path                          string                                                   `tfschema:"path"`
	Protocols                     []string                                                 `tfschema:"protocols"`
	Revision                      string                                                   `tfschema:"revision"`
	ApiType                       string                                                   `tfschema:"api_type"`
	Contact                       []ApiManagementWorkspaceApiContactModel                  `tfschema:"contact"`
	Description                   string                                                   `tfschema:"description"`
	Import                        []ApiManagementWorkspaceApiImportModel                   `tfschema:"import"`
	License                       []ApiManagementWorkspaceApiLicenseModel                  `tfschema:"license"`
	OAuth2Authorization           []ApiManagementWorkspaceApiOAuth2AuthorizationModel      `tfschema:"oauth2_authorization"`
	OpenidAuthentication          []ApiManagementWorkspaceApiOpenidAuthenticationModel     `tfschema:"openid_authentication"`
	RevisionDescription           string                                                   `tfschema:"revision_description"`
	ServiceUrl                    string                                                   `tfschema:"service_url"`
	SourceApiId                   string                                                   `tfschema:"source_api_id"`
	SubscriptionKeyParameterNames []ApiManagementWorkspaceApiSubscriptionKeyParameterModel `tfschema:"subscription_key_parameter_names"`
	SubscriptionRequired          bool                                                     `tfschema:"subscription_required"`
	TermsOfServiceUrl             string                                                   `tfschema:"terms_of_service_url"`
	Version                       string                                                   `tfschema:"version"`
	VersionDescription            string                                                   `tfschema:"version_description"`
	VersionSetId                  string                                                   `tfschema:"version_set_id"`
	IsCurrent                     bool                                                     `tfschema:"is_current"`
	IsOnline                      bool                                                     `tfschema:"is_online"`
}

type ApiManagementWorkspaceApiContactModel struct {
	Email string `tfschema:"email"`
	Name  string `tfschema:"name"`
	Url   string `tfschema:"url"`
}

type ApiManagementWorkspaceApiImportModel struct {
	ContentFormat string                                        `tfschema:"content_format"`
	ContentValue  string                                        `tfschema:"content_value"`
	WsdlSelector  []ApiManagementWorkspaceApiImportWsdlSelector `tfschema:"wsdl_selector"`
}

type ApiManagementWorkspaceApiImportWsdlSelector struct {
	ServiceName  string `tfschema:"service_name"`
	EndpointName string `tfschema:"endpoint_name"`
}

type ApiManagementWorkspaceApiLicenseModel struct {
	Name string `tfschema:"name"`
	Url  string `tfschema:"url"`
}

type ApiManagementWorkspaceApiOAuth2AuthorizationModel struct {
	AuthorizationServerName string `tfschema:"authorization_server_name"`
	Scope                   string `tfschema:"scope"`
}

type ApiManagementWorkspaceApiOpenidAuthenticationModel struct {
	OpenidProviderName        string   `tfschema:"openid_provider_name"`
	BearerTokenSendingMethods []string `tfschema:"bearer_token_sending_methods"`
}

type ApiManagementWorkspaceApiSubscriptionKeyParameterModel struct {
	Header string `tfschema:"header"`
	Query  string `tfschema:"query"`
}

type ApiManagementWorkspaceApiResource struct{}

var (
	_ sdk.ResourceWithUpdate        = ApiManagementWorkspaceApiResource{}
	_ sdk.ResourceWithCustomizeDiff = ApiManagementWorkspaceApiResource{}
)

func (r ApiManagementWorkspaceApiResource) ResourceType() string {
	return "azurerm_api_management_workspace_api"
}

func (r ApiManagementWorkspaceApiResource) ModelObject() interface{} {
	return &ApiManagementWorkspaceApiModel{}
}

func (r ApiManagementWorkspaceApiResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return api.ValidateWorkspaceApiID
}

func (r ApiManagementWorkspaceApiResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": schemaz.SchemaApiManagementApiName(),

		"api_management_workspace_id": commonschema.ResourceIDReferenceRequiredForceNew(&workspace.WorkspaceId{}),

		"revision": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"api_type": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true,
			ValidateFunc: validation.StringInSlice(api.PossibleValuesForApiType(), false),
		},

		"contact": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MinItems: 1,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"email": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validate.EmailAddress,
					},
					"name": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"url": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.IsURLWithHTTPorHTTPS,
					},
				},
			},
		},

		"description": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"display_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"import": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"content_value": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"content_format": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(api.PossibleValuesForContentFormat(), false),
					},

					"wsdl_selector": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"service_name": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"endpoint_name": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
							},
						},
					},
				},
			},
		},

		"license": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MinItems: 1,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"url": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.IsURLWithHTTPorHTTPS,
					},
				},
			},
		},

		"oauth2_authorization": {
			Type:          pluginsdk.TypeList,
			Optional:      true,
			MaxItems:      1,
			ConflictsWith: []string{"openid_authentication"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"authorization_server_name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validate.ApiManagementChildName,
					},
					"scope": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"openid_authentication": {
			Type:          pluginsdk.TypeList,
			Optional:      true,
			MaxItems:      1,
			ConflictsWith: []string{"oauth2_authorization"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"openid_provider_name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validate.ApiManagementChildName,
					},
					"bearer_token_sending_methods": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringInSlice(api.PossibleValuesForBearerTokenSendingMethods(), false),
						},
					},
				},
			},
		},

		"path": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true,
			ValidateFunc: validate.ApiManagementApiPath,
		},

		"protocols": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringInSlice(api.PossibleValuesForProtocol(), false),
			},
		},

		"revision_description": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"service_url": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
		},

		"source_api_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validate.ApiID,
		},

		"subscription_key_parameter_names": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"header": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"query": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"subscription_required": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"terms_of_service_url": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.IsURLWithHTTPorHTTPS,
		},

		"version": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"version_description": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"version_set_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func (r ApiManagementWorkspaceApiResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"is_current": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},

		"is_online": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},
	}
}

func (r ApiManagementWorkspaceApiResource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model ApiManagementWorkspaceApiModel
			if err := metadata.DecodeDiff(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			if model.Version != "" && model.VersionSetId == "" {
				return errors.New("`version` must be set with `version_set_id`")
			}

			protocols := expandWorkspaceApiProtocols(model.Protocols)
			if model.SourceApiId == "" && (model.DisplayName == "" || protocols == nil || len(*protocols) == 0) {
				return errors.New("`display_name`, `protocols` are required when `source_api_id` is not set")
			}

			if model.ApiType == string(api.ApiTypeWebsocket) && model.ServiceUrl == "" {
				return errors.New("`service_url` is required when `api_type` is `websocket`")
			}

			return nil
		},
	}
}

func (r ApiManagementWorkspaceApiResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ApiManagement.WorkspaceApiClient

			var model ApiManagementWorkspaceApiModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			workspaceId, err := workspace.ParseWorkspaceID(model.ApiManagementWorkspaceId)
			if err != nil {
				return err
			}

			apiId := fmt.Sprintf("%s;rev=%s", model.Name, model.Revision)
			id := api.NewWorkspaceApiID(workspaceId.SubscriptionId, workspaceId.ResourceGroupName, workspaceId.ServiceName, workspaceId.WorkspaceId, apiId)

			existing, err := client.WorkspaceApiGet(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			apiType := api.ApiTypeHTTP
			if model.ApiType != "" {
				apiType = api.ApiType(model.ApiType)
			}
			soapApiType := workspaceSoapApiTypeFromApiType(apiType)

			if len(model.Import) > 0 {
				apiParams := expandWorkspaceApiImport(model.Import, apiType, soapApiType, model.Path, model.ServiceUrl, model.Version, model.VersionSetId)
				if apiParams != nil {
					if _, err := client.WorkspaceApiCreateOrUpdate(ctx, id, *apiParams, api.DefaultWorkspaceApiCreateOrUpdateOperationOptions()); err != nil {
						return fmt.Errorf("creating with import %s: %+v", id, err)
					}
				}
			}

			protocols := expandWorkspaceApiProtocols(model.Protocols)
			subscriptionKeyParameterNames := expandWorkspaceApiSubscriptionKeyParamNames(model.SubscriptionKeyParameterNames)
			authenticationSettings := expandWorkspaceApiAuthenticationSettings(model.OAuth2Authorization, model.OpenidAuthentication)
			contactInfo := expandWorkspaceApiContact(model.Contact)
			licenseInfo := expandWorkspaceApiLicense(model.License)

			params := api.ApiCreateOrUpdateParameter{
				Properties: &api.ApiCreateOrUpdateProperties{
					Type:                          pointer.To(apiType),
					ApiType:                       pointer.To(soapApiType),
					Path:                          model.Path,
					Protocols:                     protocols,
					SubscriptionKeyParameterNames: subscriptionKeyParameterNames,
					SubscriptionRequired:          pointer.To(model.SubscriptionRequired),
					AuthenticationSettings:        authenticationSettings,
					ApiRevisionDescription:        pointer.To(model.RevisionDescription),
					ApiVersionDescription:         pointer.To(model.VersionDescription),
					Contact:                       contactInfo,
					License:                       licenseInfo,
				},
			}

			if model.ServiceUrl != "" {
				params.Properties.ServiceURL = pointer.To(model.ServiceUrl)
			}

			if model.SourceApiId != "" {
				params.Properties.SourceApiId = pointer.To(model.SourceApiId)
			}

			if model.Description != "" {
				params.Properties.Description = pointer.To(model.Description)
			}

			if model.DisplayName != "" {
				params.Properties.DisplayName = pointer.To(model.DisplayName)
			}

			if model.Version != "" {
				params.Properties.ApiVersion = pointer.To(model.Version)
			}

			if model.VersionSetId != "" {
				params.Properties.ApiVersionSetId = pointer.To(model.VersionSetId)
			}

			if model.TermsOfServiceUrl != "" {
				params.Properties.TermsOfServiceURL = pointer.To(model.TermsOfServiceUrl)
			}

			if _, err := client.WorkspaceApiCreateOrUpdate(ctx, id, params, api.DefaultWorkspaceApiCreateOrUpdateOperationOptions()); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r ApiManagementWorkspaceApiResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ApiManagement.WorkspaceApiClient

			var model ApiManagementWorkspaceApiModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id, err := api.ParseWorkspaceApiID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			apiType := api.ApiTypeHTTP
			if model.ApiType != "" {
				apiType = api.ApiType(model.ApiType)
			}
			soapApiType := workspaceSoapApiTypeFromApiType(apiType)

			if metadata.ResourceData.HasChange("import") {
				if len(model.Import) > 0 {
					apiParams := expandWorkspaceApiImport(model.Import, apiType, soapApiType, model.Path, model.ServiceUrl, model.Version, model.VersionSetId)
					if apiParams != nil {
						if _, err := client.WorkspaceApiCreateOrUpdate(ctx, *id, *apiParams, api.DefaultWorkspaceApiCreateOrUpdateOperationOptions()); err != nil {
							return fmt.Errorf("updating with import %s: %+v", id, err)
						}
					}
				}
			}

			resp, err := client.WorkspaceApiGet(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			if resp.Model == nil || resp.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", id)
			}

			existing := resp.Model.Properties
			if existing.Type != nil {
				soapApiType = workspaceSoapApiTypeFromApiType(pointer.From(existing.Type))
			}

			prop := &api.ApiCreateOrUpdateProperties{
				Path:                          existing.Path,
				Protocols:                     existing.Protocols,
				ServiceURL:                    existing.ServiceURL,
				Description:                   existing.Description,
				ApiVersionDescription:         existing.ApiVersionDescription,
				ApiRevisionDescription:        existing.ApiRevisionDescription,
				SubscriptionRequired:          existing.SubscriptionRequired,
				SubscriptionKeyParameterNames: existing.SubscriptionKeyParameterNames,
				Contact:                       existing.Contact,
				License:                       existing.License,
				SourceApiId:                   existing.SourceApiId,
				DisplayName:                   existing.DisplayName,
				ApiVersion:                    existing.ApiVersion,
				ApiVersionSetId:               existing.ApiVersionSetId,
				TermsOfServiceURL:             existing.TermsOfServiceURL,
				Type:                          existing.Type,
				ApiType:                       pointer.To(soapApiType),
			}

			if v := existing.AuthenticationSettings; v != nil {
				authenticationSettings := &api.AuthenticationSettingsContract{}
				if v.OAuth2 != nil {
					authenticationSettings.OAuth2 = v.OAuth2
					prop.AuthenticationSettings = authenticationSettings
				}

				if v.Openid != nil {
					authenticationSettings.Openid = v.Openid
					prop.AuthenticationSettings = authenticationSettings
				}
			}

			if metadata.ResourceData.HasChange("path") {
				prop.Path = model.Path
			}

			if metadata.ResourceData.HasChange("protocols") {
				prop.Protocols = expandWorkspaceApiProtocols(model.Protocols)
			}

			if metadata.ResourceData.HasChange("api_type") {
				prop.Type = pointer.To(apiType)
				prop.ApiType = pointer.To(soapApiType)
			}

			if metadata.ResourceData.HasChange("service_url") {
				prop.ServiceURL = pointer.To(model.ServiceUrl)
			}

			if metadata.ResourceData.HasChange("description") {
				prop.Description = pointer.To(model.Description)
			}

			if metadata.ResourceData.HasChange("revision_description") {
				prop.ApiRevisionDescription = pointer.To(model.RevisionDescription)
			}

			if metadata.ResourceData.HasChange("version_description") {
				prop.ApiVersionDescription = pointer.To(model.VersionDescription)
			}

			if metadata.ResourceData.HasChange("subscription_required") {
				prop.SubscriptionRequired = pointer.To(model.SubscriptionRequired)
			}

			if metadata.ResourceData.HasChange("subscription_key_parameter_names") {
				prop.SubscriptionKeyParameterNames = expandWorkspaceApiSubscriptionKeyParamNames(model.SubscriptionKeyParameterNames)
			}

			if metadata.ResourceData.HasChange("oauth2_authorization") {
				authenticationSettings := &api.AuthenticationSettingsContract{}
				oAuth2AuthorizationSettings := expandWorkspaceApiOAuth2AuthenticationSettingsContract(model.OAuth2Authorization)
				authenticationSettings.OAuth2 = oAuth2AuthorizationSettings
				prop.AuthenticationSettings = authenticationSettings
			}

			if metadata.ResourceData.HasChange("openid_authentication") {
				authenticationSettings := &api.AuthenticationSettingsContract{}
				openIDAuthorizationSettings := expandWorkspaceApiOpenIDAuthenticationSettingsContract(model.OpenidAuthentication)
				authenticationSettings.Openid = openIDAuthorizationSettings
				prop.AuthenticationSettings = authenticationSettings
			}

			if metadata.ResourceData.HasChange("contact") {
				prop.Contact = expandWorkspaceApiContact(model.Contact)
			}

			if metadata.ResourceData.HasChange("license") {
				prop.License = expandWorkspaceApiLicense(model.License)
			}

			if metadata.ResourceData.HasChange("source_api_id") {
				prop.SourceApiId = pointer.To(model.SourceApiId)
			}

			if metadata.ResourceData.HasChange("display_name") {
				prop.DisplayName = pointer.To(model.DisplayName)
			}

			if metadata.ResourceData.HasChange("version") {
				prop.ApiVersion = pointer.To(model.Version)
			}

			if metadata.ResourceData.HasChange("version_set_id") {
				prop.ApiVersionSetId = pointer.To(model.VersionSetId)
			}

			if metadata.ResourceData.HasChange("terms_of_service_url") {
				prop.TermsOfServiceURL = pointer.To(model.TermsOfServiceUrl)
			}

			params := api.ApiCreateOrUpdateParameter{
				Properties: prop,
			}

			if _, err := client.WorkspaceApiCreateOrUpdate(ctx, *id, params, api.DefaultWorkspaceApiCreateOrUpdateOperationOptions()); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (r ApiManagementWorkspaceApiResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ApiManagement.WorkspaceApiClient

			id, err := api.ParseWorkspaceApiID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.WorkspaceApiGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			state := ApiManagementWorkspaceApiModel{
				Name:                     getWorkspaceApiName(id.ApiId),
				ApiManagementWorkspaceId: workspace.NewWorkspaceID(id.SubscriptionId, id.ResourceGroupName, id.ServiceName, id.WorkspaceId).ID(),
			}

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					apiType := string(pointer.From(props.Type))
					if len(apiType) == 0 {
						apiType = string(api.ApiTypeHTTP)
					}
					state.ApiType = apiType
					state.Description = pointer.From(props.Description)
					state.DisplayName = pointer.From(props.DisplayName)
					state.IsCurrent = pointer.From(props.IsCurrent)
					state.IsOnline = pointer.From(props.IsOnline)
					state.Path = props.Path
					state.ServiceUrl = pointer.From(props.ServiceURL)
					state.Revision = pointer.From(props.ApiRevision)
					state.SubscriptionRequired = pointer.From(props.SubscriptionRequired)
					state.Version = pointer.From(props.ApiVersion)
					state.VersionSetId = pointer.From(props.ApiVersionSetId)
					state.RevisionDescription = pointer.From(props.ApiRevisionDescription)
					state.VersionDescription = pointer.From(props.ApiVersionDescription)
					state.TermsOfServiceUrl = pointer.From(props.TermsOfServiceURL)

					state.Protocols = flattenWorkspaceApiProtocols(props.Protocols)
					state.SubscriptionKeyParameterNames = flattenWorkspaceApiSubscriptionKeyParamNames(props.SubscriptionKeyParameterNames)
					state.OAuth2Authorization = flattenWorkspaceApiOAuth2Authorization(props.AuthenticationSettings)
					state.OpenidAuthentication = flattenWorkspaceApiOpenIDAuthentication(props.AuthenticationSettings)
					state.Contact = flattenWorkspaceApiContact(props.Contact)
					state.License = flattenWorkspaceApiLicense(props.License)
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r ApiManagementWorkspaceApiResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ApiManagement.WorkspaceApiClient

			id, err := api.ParseWorkspaceApiID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err = client.WorkspaceApiDelete(ctx, *id, api.DefaultWorkspaceApiDeleteOperationOptions()); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}

func expandWorkspaceApiImport(importInput []ApiManagementWorkspaceApiImportModel, apiType api.ApiType, soapApiType api.SoapApiType, path, serviceUrl, version, versionSetId string) *api.ApiCreateOrUpdateParameter {
	if len(importInput) == 0 {
		return nil
	}

	importV := importInput[0]
	apiParams := api.ApiCreateOrUpdateParameter{
		Properties: &api.ApiCreateOrUpdateProperties{
			Type:    pointer.To(apiType),
			ApiType: pointer.To(soapApiType),
			Format:  pointer.To(api.ContentFormat(importV.ContentFormat)),
			Value:   pointer.To(importV.ContentValue),
			Path:    path,
		},
	}

	if len(importV.WsdlSelector) > 0 {
		wsdlSelectorV := importV.WsdlSelector[0]
		apiParams.Properties.WsdlSelector = &api.ApiCreateOrUpdatePropertiesWsdlSelector{
			WsdlServiceName:  pointer.To(wsdlSelectorV.ServiceName),
			WsdlEndpointName: pointer.To(wsdlSelectorV.EndpointName),
		}
	}

	if serviceUrl != "" {
		apiParams.Properties.ServiceURL = pointer.To(serviceUrl)
	}

	if version != "" {
		apiParams.Properties.ApiVersion = pointer.To(version)
	}

	if versionSetId != "" {
		apiParams.Properties.ApiVersionSetId = pointer.To(versionSetId)
	}

	return &apiParams
}

func expandWorkspaceApiProtocols(input []string) *[]api.Protocol {
	if len(input) == 0 {
		return nil
	}
	results := make([]api.Protocol, 0)

	for _, v := range input {
		results = append(results, api.Protocol(v))
	}

	return &results
}

func flattenWorkspaceApiProtocols(input *[]api.Protocol) []string {
	if input == nil {
		return []string{}
	}

	results := make([]string, 0)
	for _, v := range *input {
		results = append(results, string(v))
	}

	return results
}

func expandWorkspaceApiSubscriptionKeyParamNames(input []ApiManagementWorkspaceApiSubscriptionKeyParameterModel) *api.SubscriptionKeyParameterNamesContract {
	if len(input) == 0 {
		return nil
	}

	v := input[0]
	return &api.SubscriptionKeyParameterNamesContract{
		Query:  pointer.To(v.Query),
		Header: pointer.To(v.Header),
	}
}

func flattenWorkspaceApiSubscriptionKeyParamNames(paramNames *api.SubscriptionKeyParameterNamesContract) []ApiManagementWorkspaceApiSubscriptionKeyParameterModel {
	if paramNames == nil {
		return make([]ApiManagementWorkspaceApiSubscriptionKeyParameterModel, 0)
	}

	return []ApiManagementWorkspaceApiSubscriptionKeyParameterModel{
		{
			Header: pointer.From(paramNames.Header),
			Query:  pointer.From(paramNames.Query),
		},
	}
}

func expandWorkspaceApiAuthenticationSettings(oauth2 []ApiManagementWorkspaceApiOAuth2AuthorizationModel, openid []ApiManagementWorkspaceApiOpenidAuthenticationModel) *api.AuthenticationSettingsContract {
	authenticationSettings := &api.AuthenticationSettingsContract{}

	oAuth2AuthorizationSettings := expandWorkspaceApiOAuth2AuthenticationSettingsContract(oauth2)
	authenticationSettings.OAuth2 = oAuth2AuthorizationSettings

	openIDAuthorizationSettings := expandWorkspaceApiOpenIDAuthenticationSettingsContract(openid)
	authenticationSettings.Openid = openIDAuthorizationSettings

	return authenticationSettings
}

func expandWorkspaceApiOAuth2AuthenticationSettingsContract(input []ApiManagementWorkspaceApiOAuth2AuthorizationModel) *api.OAuth2AuthenticationSettingsContract {
	if len(input) == 0 {
		return nil
	}

	oAuth2AuthorizationV := input[0]
	return &api.OAuth2AuthenticationSettingsContract{
		AuthorizationServerId: pointer.To(oAuth2AuthorizationV.AuthorizationServerName),
		Scope:                 pointer.To(oAuth2AuthorizationV.Scope),
	}
}

func flattenWorkspaceApiOAuth2Authorization(input *api.AuthenticationSettingsContract) []ApiManagementWorkspaceApiOAuth2AuthorizationModel {
	if input == nil || input.OAuth2 == nil {
		return make([]ApiManagementWorkspaceApiOAuth2AuthorizationModel, 0)
	}

	return []ApiManagementWorkspaceApiOAuth2AuthorizationModel{
		{
			AuthorizationServerName: pointer.From(input.OAuth2.AuthorizationServerId),
			Scope:                   pointer.From(input.OAuth2.Scope),
		},
	}
}

func expandWorkspaceApiOpenIDAuthenticationSettingsContract(input []ApiManagementWorkspaceApiOpenidAuthenticationModel) *api.OpenIdAuthenticationSettingsContract {
	if len(input) == 0 {
		return nil
	}

	openIDAuthorizationV := input[0]
	return &api.OpenIdAuthenticationSettingsContract{
		OpenidProviderId:          pointer.To(openIDAuthorizationV.OpenidProviderName),
		BearerTokenSendingMethods: expandWorkspaceApiOpenIDAuthenticationSettingsBearerTokenSendingMethods(openIDAuthorizationV.BearerTokenSendingMethods),
	}
}

func expandWorkspaceApiOpenIDAuthenticationSettingsBearerTokenSendingMethods(input []string) *[]api.BearerTokenSendingMethods {
	if input == nil {
		return nil
	}
	results := make([]api.BearerTokenSendingMethods, 0)

	for _, v := range input {
		results = append(results, api.BearerTokenSendingMethods(v))
	}

	return &results
}

func flattenWorkspaceApiOpenIDAuthentication(input *api.AuthenticationSettingsContract) []ApiManagementWorkspaceApiOpenidAuthenticationModel {
	if input == nil || input.Openid == nil {
		return make([]ApiManagementWorkspaceApiOpenidAuthenticationModel, 0)
	}

	bearerTokenSendingMethods := make([]string, 0)
	if s := input.Openid.BearerTokenSendingMethods; s != nil {
		for _, v := range *s {
			bearerTokenSendingMethods = append(bearerTokenSendingMethods, string(v))
		}
	}

	return []ApiManagementWorkspaceApiOpenidAuthenticationModel{
		{
			OpenidProviderName:        pointer.From(input.Openid.OpenidProviderId),
			BearerTokenSendingMethods: bearerTokenSendingMethods,
		},
	}
}

func expandWorkspaceApiContact(input []ApiManagementWorkspaceApiContactModel) *api.ApiContactInformation {
	if len(input) == 0 {
		return nil
	}

	v := input[0]
	return &api.ApiContactInformation{
		Email: pointer.To(v.Email),
		Name:  pointer.To(v.Name),
		Url:   pointer.To(v.Url),
	}
}

func flattenWorkspaceApiContact(contact *api.ApiContactInformation) []ApiManagementWorkspaceApiContactModel {
	if contact == nil {
		return make([]ApiManagementWorkspaceApiContactModel, 0)
	}

	return []ApiManagementWorkspaceApiContactModel{
		{
			Email: pointer.From(contact.Email),
			Name:  pointer.From(contact.Name),
			Url:   pointer.From(contact.Url),
		},
	}
}

func expandWorkspaceApiLicense(input []ApiManagementWorkspaceApiLicenseModel) *api.ApiLicenseInformation {
	if len(input) == 0 {
		return nil
	}

	v := input[0]
	return &api.ApiLicenseInformation{
		Name: pointer.To(v.Name),
		Url:  pointer.To(v.Url),
	}
}

func flattenWorkspaceApiLicense(license *api.ApiLicenseInformation) []ApiManagementWorkspaceApiLicenseModel {
	if license == nil {
		return make([]ApiManagementWorkspaceApiLicenseModel, 0)
	}

	return []ApiManagementWorkspaceApiLicenseModel{
		{
			Name: pointer.From(license.Name),
			Url:  pointer.From(license.Url),
		},
	}
}

func workspaceSoapApiTypeFromApiType(apiType api.ApiType) api.SoapApiType {
	return map[api.ApiType]api.SoapApiType{
		api.ApiTypeGraphql:   api.SoapApiTypeGraphql,
		api.ApiTypeHTTP:      api.SoapApiTypeHTTP,
		api.ApiTypeSoap:      api.SoapApiTypeSoap,
		api.ApiTypeWebsocket: api.SoapApiTypeWebsocket,
	}[apiType]
}

func getWorkspaceApiName(apiId string) string {
	name := apiId
	if len(apiId) > 0 && apiId[len(apiId)-1:] == ";" {
		name = apiId[:len(apiId)-1]
	}
	for i := len(apiId) - 1; i >= 0; i-- {
		if apiId[i] == ';' {
			name = apiId[:i]
			break
		}
	}
	return name
}
