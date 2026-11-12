// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cdn_edge_action_version

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2026-10-01/edgeactionversions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity -resource-name cdn_edge_action_version -properties "name" -compare-values "subscription_id:cdn_edge_action_id,resource_group_name:cdn_edge_action_id,edge_action_name:cdn_edge_action_id"

type Resource struct{}

type Model struct {
	Name                  string            `tfschema:"name"`
	CdnEdgeActionID       string            `tfschema:"cdn_edge_action_id"`
	Location              string            `tfschema:"location"`
	DeploymentType        string            `tfschema:"deployment_type"`
	DefaultVersionEnabled bool              `tfschema:"default_version_enabled"`
	Code                  []CodeModel       `tfschema:"code"`
	Tags                  map[string]string `tfschema:"tags"`
	ValidationStatus      string            `tfschema:"validation_status"`
}

type CodeModel struct {
	FileName      string `tfschema:"file_name"`
	ContentBase64 string `tfschema:"content_base64"`
}

var (
	_ sdk.ResourceWithUpdate   = Resource{}
	_ sdk.ResourceWithIdentity = Resource{}
)

func (Resource) ResourceType() string             { return "azurerm_cdn_edge_action_version" }
func (Resource) ModelObject() interface{}         { return &Model{} }
func (Resource) Identity() resourceids.ResourceId { return &edgeactionversions.VersionId{} }
func (Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return edgeactionversions.ValidateVersionID
}

func (Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.All(validation.StringLenBetween(1, 50), validation.StringMatch(regexp.MustCompile(`^[a-zA-Z0-9]+$`), "must contain only alphanumeric characters")),
		},
		"cdn_edge_action_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: edgeactionversions.ValidateEdgeActionID,
		},
		"location": commonschema.Location(),
		"deployment_type": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(edgeactionversions.PossibleValuesForEdgeActionVersionDeploymentType(), false),
		},
		"default_version_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			// NOTE: O+C - Azure selects the first default and swaps can change other versions' default status.
			Computed: true,
			ValidateFunc: func(v interface{}, k string) ([]string, []error) {
				if !v.(bool) {
					return nil, []error{fmt.Errorf("`%s` can only be `true`; omit it on other versions and promote another version to change the default", k)}
				}
				return nil, nil
			},
		},
		"code": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
				"file_name": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},
				"content_base64": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					Sensitive:    true,
					ValidateFunc: validation.All(validation.StringIsNotEmpty, validation.StringIsBase64),
				},
			}},
		},
		"tags": commonschema.Tags(),
	}
}

func (Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"validation_status": {Type: pluginsdk.TypeString, Computed: true},
	}
}

func (Resource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			diff := metadata.ResourceDiff
			if diff.Id() != "" && diff.HasChange("code") && diff.NewValueKnown("code") && len(diff.Get("code").([]interface{})) == 0 {
				return diff.ForceNew("code")
			}
			return nil
		},
	}
}

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model Model
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}
			parent, err := edgeactionversions.ParseEdgeActionID(model.CdnEdgeActionID)
			if err != nil {
				return err
			}
			locks.ByName(parent.ID(), "azurerm_cdn_edge_action")
			defer locks.UnlockByName(parent.ID(), "azurerm_cdn_edge_action")
			id := edgeactionversions.NewVersionID(parent.SubscriptionId, parent.ResourceGroupName, parent.EdgeActionName, model.Name)
			client := metadata.Client.Cdn.EdgeActionVersionsClient
			defaultVersion := edgeactionversions.EdgeActionIsDefaultVersionFalse
			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}
			if !response.WasNotFound(existing.HttpResponse) {
				if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
				if existing.Model == nil || existing.Model.Properties == nil {
					return fmt.Errorf("retrieving %s: model or properties was nil", id)
				}
				defaultVersion = existing.Model.Properties.IsDefaultVersion
			}
			payload := edgeactionversions.EdgeActionVersion{
				Location: location.Normalize(model.Location),
				Properties: &edgeactionversions.EdgeActionVersionProperties{
					DeploymentType:   edgeactionversions.EdgeActionVersionDeploymentType(model.DeploymentType),
					IsDefaultVersion: defaultVersion,
				},
				Tags: pointer.To(model.Tags),
			}
			if err := client.CreateCallbackThenPoll(ctx, id, payload, metadata.SetIDAndIdentityCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}
			metadata.SetID(id)
			if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, &id); err != nil {
				return err
			}
			if err := deployCode(ctx, client, id, model.Code); err != nil {
				return err
			}
			if model.DefaultVersionEnabled {
				if err := promoteVersion(ctx, client, id); err != nil {
					return err
				}
			}
			return nil
		},
	}
}

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := edgeactionversions.ParseVersionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}
			resp, err := metadata.Client.Cdn.EdgeActionVersionsClient.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}
			return r.flatten(metadata, id, resp.Model)
		},
	}
}

func (Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) (err error) {
			id, err := edgeactionversions.ParseVersionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}
			parent := edgeactionversions.NewEdgeActionID(id.SubscriptionId, id.ResourceGroupName, id.EdgeActionName)
			locks.ByName(parent.ID(), "azurerm_cdn_edge_action")
			defer locks.UnlockByName(parent.ID(), "azurerm_cdn_edge_action")
			var model Model
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}
			codeChanged := metadata.ResourceData.HasChange("code")
			codeDeployed := false
			defer func() {
				if err != nil && codeChanged && !codeDeployed {
					previous, _ := metadata.ResourceData.GetChange("code")
					if restoreErr := metadata.ResourceData.Set("code", previous); restoreErr != nil {
						err = fmt.Errorf("updating %s: %+v; restoring previous `code`: %+v", id, err, restoreErr)
					}
				}
			}()
			client := metadata.Client.Cdn.EdgeActionVersionsClient
			if metadata.ResourceData.HasChange("tags") {
				updatedTags := model.Tags
				if updatedTags == nil {
					updatedTags = map[string]string{}
				}
				if err := client.UpdateThenPoll(ctx, *id, edgeactionversions.EdgeActionVersionUpdate{Tags: pointer.To(updatedTags)}); err != nil {
					return fmt.Errorf("updating %s: %+v", id, err)
				}
			}
			if codeChanged {
				if err := deployCode(ctx, client, *id, model.Code); err != nil {
					return err
				}
				codeDeployed = true
			}
			if metadata.ResourceData.HasChange("default_version_enabled") && model.DefaultVersionEnabled {
				if err := promoteVersion(ctx, client, *id); err != nil {
					return err
				}
			}
			return nil
		},
	}
}

func (Resource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := edgeactionversions.ParseVersionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}
			parent := edgeactionversions.NewEdgeActionID(id.SubscriptionId, id.ResourceGroupName, id.EdgeActionName)
			locks.ByName(parent.ID(), "azurerm_cdn_edge_action")
			defer locks.UnlockByName(parent.ID(), "azurerm_cdn_edge_action")
			if err := metadata.Client.Cdn.EdgeActionVersionsClient.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}
			return nil
		},
	}
}

func deployCode(ctx context.Context, client *edgeactionversions.EdgeActionVersionsClient, id edgeactionversions.VersionId, code []CodeModel) error {
	if len(code) == 0 {
		return nil
	}
	if err := client.DeployVersionCodeThenPoll(ctx, id, edgeactionversions.VersionCode{Name: code[0].FileName, Content: code[0].ContentBase64}); err != nil {
		return fmt.Errorf("deploying code for %s: %+v", id, err)
	}
	return nil
}

func promoteVersion(ctx context.Context, client *edgeactionversions.EdgeActionVersionsClient, id edgeactionversions.VersionId) error {
	resp, err := client.Get(ctx, id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}
	if resp.Model == nil || resp.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: model or properties was nil", id)
	}
	if resp.Model.Properties.IsDefaultVersion == edgeactionversions.EdgeActionIsDefaultVersionTrue {
		return nil
	}
	result, err := client.SwapDefault(ctx, id)
	// The SDK attempts to build a poller even for the API's synchronous 204 response.
	if result.HttpResponse != nil && result.HttpResponse.StatusCode == http.StatusNoContent {
		return nil
	}
	if err != nil {
		return fmt.Errorf("promoting %s to the default version: %+v", id, err)
	}
	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling default-version promotion for %s: %+v", id, err)
	}
	return nil
}

func (Resource) flatten(metadata sdk.ResourceMetaData, id *edgeactionversions.VersionId, model *edgeactionversions.EdgeActionVersion) error {
	if model == nil || model.Properties == nil {
		return fmt.Errorf("retrieving %s: model or properties was nil", id)
	}
	var state Model
	if err := metadata.Decode(&state); err != nil {
		return fmt.Errorf("decoding: %+v", err)
	}
	state.Name = id.VersionName
	state.CdnEdgeActionID = edgeactionversions.NewEdgeActionID(id.SubscriptionId, id.ResourceGroupName, id.EdgeActionName).ID()
	state.Location = location.Normalize(model.Location)
	state.DeploymentType = string(model.Properties.DeploymentType)
	state.DefaultVersionEnabled = model.Properties.IsDefaultVersion == edgeactionversions.EdgeActionIsDefaultVersionTrue
	state.ValidationStatus = pointer.FromEnum(model.Properties.ValidationStatus)
	state.Tags = pointer.From(model.Tags)
	if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
		return err
	}
	return metadata.Encode(&state)
}
