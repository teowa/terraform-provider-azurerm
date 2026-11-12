// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cdn_edge_action_version

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2026-10-01/edgeactionversions"
	sdkclient "github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	cdnclient "github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/client"
)

func TestUnitVersionDefaultConfiguration(t *testing.T) {
	for _, test := range []struct {
		name      string
		value     interface{}
		wantError bool
	}{
		{name: "service-selected default"},
		{name: "promote", value: true},
		{name: "cannot demote without another version", value: false, wantError: true},
	} {
		t.Run(test.name, func(t *testing.T) {
			config := map[string]interface{}{
				"name":               "v1",
				"cdn_edge_action_id": "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test/providers/Microsoft.Cdn/edgeActions/test",
				"location":           "global",
				"deployment_type":    "file",
			}
			if test.value != nil {
				config["default_version_enabled"] = test.value
			}
			diags := sdk.WrappedResource(Resource{}).Validate(terraform.NewResourceConfigRaw(config))
			if diags.HasError() != test.wantError {
				t.Fatalf("unexpected diagnostics: %+v", diags)
			}
		})
	}
}

func TestUnitVersionUpdateAppliesTagsCodeAndPromotion(t *testing.T) {
	for _, failOperation := range []string{"", "tags", "code", "promote"} {
		t.Run("failure="+failOperation, func(t *testing.T) {
			id := edgeactionversions.NewVersionID("00000000-0000-0000-0000-000000000000", "test", "test", "v2")
			var operations []string
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.Header().Set(sdkclient.SkipPollingDelayHeader, "true")
				switch {
				case req.Method == http.MethodPatch:
					operations = append(operations, "tags")
					var payload map[string]interface{}
					if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
						t.Error(err)
					}
					if !reflect.DeepEqual(payload, map[string]interface{}{"tags": map[string]interface{}{}}) {
						t.Errorf("expected tags-only PATCH with explicit empty tags, got %#v", payload)
					}
					if failOperation == "tags" {
						w.WriteHeader(http.StatusBadRequest)
						return
					}
				case req.Method == http.MethodPost && strings.HasSuffix(req.URL.Path, "/deployVersionCode"):
					operations = append(operations, "code")
					var payload edgeactionversions.VersionCode
					if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
						t.Error(err)
					}
					if payload.Name != "index.js" || payload.Content != "aGVsbG8=" {
						t.Errorf("unexpected code payload: %#v", payload)
					}
					if failOperation == "code" {
						w.WriteHeader(http.StatusBadRequest)
						if _, err := w.Write([]byte(`{"error":{"code":"InvalidCode","message":"invalid code"}}`)); err != nil {
							t.Error(err)
						}
						return
					}
				case req.Method == http.MethodPost && strings.HasSuffix(req.URL.Path, "/swapDefault"):
					operations = append(operations, "promote")
					if failOperation == "promote" {
						w.WriteHeader(http.StatusBadRequest)
						return
					}
					w.WriteHeader(http.StatusNoContent)
					return
				case req.Method != http.MethodGet:
					t.Errorf("unexpected request: %s %s", req.Method, req.URL.Path)
				}
				if _, err := w.Write([]byte(`{"location":"global","properties":{"deploymentType":"file","isDefaultVersion":"False","provisioningState":"Succeeded"}}`)); err != nil {
					t.Error(err)
				}
			}))
			defer server.Close()
			client, err := edgeactionversions.NewEdgeActionVersionsClientWithBaseURI(environments.ResourceManagerAPI(server.URL))
			if err != nil {
				t.Fatal(err)
			}
			client.Client.AuthorizeRequest = nil
			client.Client.DisableRetries = true
			metadata := sdk.NewResourceMetaData(&clients.Client{Cdn: &cdnclient.Client{EdgeActionVersionsClient: client}}, Resource{})
			r := sdk.WrappedResource(Resource{})
			config := map[string]interface{}{
				"name":               "v2",
				"location":           "global",
				"deployment_type":    "file",
				"cdn_edge_action_id": edgeactionversions.NewEdgeActionID(id.SubscriptionId, id.ResourceGroupName, id.EdgeActionName).ID(),
				"code":               []interface{}{map[string]interface{}{"file_name": "old.js", "content_base64": "b2xk"}},
				"tags":               map[string]interface{}{"environment": "old"},
			}
			old := schema.TestResourceDataRaw(t, r.Schema, config)
			old.SetId(id.ID())
			state := old.State()
			delete(config, "tags")
			config["code"] = []interface{}{map[string]interface{}{"file_name": "index.js", "content_base64": "aGVsbG8="}}
			config["default_version_enabled"] = true
			diff, err := r.Diff(t.Context(), state, terraform.NewResourceConfigRaw(config), metadata.Client)
			if err != nil {
				t.Fatal(err)
			}
			metadata.ResourceData, err = schema.InternalMap(r.Schema).Data(state, diff)
			if err != nil {
				t.Fatal(err)
			}
			ctx, cancel := context.WithTimeout(t.Context(), time.Minute)
			defer cancel()
			updateErr := (Resource{}).Update().Func(ctx, metadata)
			expectedOperations := []string{"tags", "code", "promote"}
			if failOperation != "" {
				if updateErr == nil {
					t.Fatalf("expected %s failure", failOperation)
				}
				switch failOperation {
				case "tags":
					expectedOperations = []string{"tags"}
				case "code":
					expectedOperations = []string{"tags", "code"}
				}
			} else if updateErr != nil {
				t.Fatal(updateErr)
			}
			expectedFileName := "index.js"
			if failOperation == "tags" || failOperation == "code" {
				expectedFileName = "old.js"
			}
			if actual := metadata.ResourceData.Get("code.0.file_name"); actual != expectedFileName {
				t.Fatalf("expected source %q in state for retry, got %v", expectedFileName, actual)
			}
			if !reflect.DeepEqual(operations, expectedOperations) {
				t.Fatalf("expected all updates in order, got %v", operations)
			}
		})
	}
}

func TestUnitVersionPromotionFailure(t *testing.T) {
	id := edgeactionversions.NewVersionID("00000000-0000-0000-0000-000000000000", "test", "test", "v2")
	for _, status := range []string{"Failed", "Canceled"} {
		t.Run(status, func(t *testing.T) {
			var server *httptest.Server
			server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.Header().Set(sdkclient.SkipPollingDelayHeader, "true")
				if strings.HasSuffix(req.URL.Path, "/swapDefault") {
					w.Header().Set("Azure-AsyncOperation", server.URL+"/operation")
					w.WriteHeader(http.StatusAccepted)
					return
				}
				body := `{"location":"global","properties":{"deploymentType":"file","isDefaultVersion":"False"}}`
				if req.URL.Path == "/operation" {
					body = `{"status":"` + status + `"}`
				}
				if _, err := w.Write([]byte(body)); err != nil {
					t.Error(err)
				}
			}))
			defer server.Close()
			client, err := edgeactionversions.NewEdgeActionVersionsClientWithBaseURI(environments.ResourceManagerAPI(server.URL))
			if err != nil {
				t.Fatal(err)
			}
			client.Client.AuthorizeRequest = nil
			client.Client.DisableRetries = true
			ctx, cancel := context.WithTimeout(t.Context(), time.Minute)
			defer cancel()
			if err := promoteVersion(ctx, client, id); err == nil || !strings.Contains(err.Error(), "polling default-version promotion") {
				t.Fatalf("expected terminal %s promotion to return a polling error, got %v", status, err)
			}
		})
	}
}

func TestUnitVersionCodeRemovalRequiresReplacement(t *testing.T) {
	r := sdk.WrappedResource(Resource{})
	config := map[string]interface{}{
		"name":               "v1",
		"cdn_edge_action_id": "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test/providers/Microsoft.Cdn/edgeActions/test",
		"location":           "global",
		"deployment_type":    "file",
		"code":               []interface{}{map[string]interface{}{"file_name": "index.js", "content_base64": "aGVsbG8="}},
	}
	data := schema.TestResourceDataRaw(t, r.Schema, config)
	data.SetId(config["cdn_edge_action_id"].(string) + "/versions/v1")
	state := data.State()
	delete(config, "code")
	diff, err := r.Diff(t.Context(), state, terraform.NewResourceConfigRaw(config), &clients.Client{})
	if err != nil {
		t.Fatal(err)
	}

	if diff == nil || !diff.RequiresNew() {
		t.Fatal("removing deployed code must replace the version")
	}
}

func TestUnitVersionOverwritePreservesDefault(t *testing.T) {
	id := edgeactionversions.NewVersionID("00000000-0000-0000-0000-000000000000", "test", "test", "v1")
	var overwritten bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set(sdkclient.SkipPollingDelayHeader, "true")
		if req.Method == http.MethodPut {
			overwritten = true
			var payload edgeactionversions.EdgeActionVersion
			if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
				t.Error(err)
			}
			if payload.Properties == nil || payload.Properties.IsDefaultVersion != edgeactionversions.EdgeActionIsDefaultVersionTrue {
				t.Errorf("overwrite must preserve the existing default flag, got %#v", payload.Properties)
			}
		} else if req.Method != http.MethodGet {
			t.Errorf("unexpected request: %s %s", req.Method, req.URL.Path)
		}
		if _, err := w.Write([]byte(`{"location":"global","properties":{"deploymentType":"file","isDefaultVersion":"True","provisioningState":"Succeeded"}}`)); err != nil {
			t.Error(err)
		}
	}))
	defer server.Close()
	client, err := edgeactionversions.NewEdgeActionVersionsClientWithBaseURI(environments.ResourceManagerAPI(server.URL))
	if err != nil {
		t.Fatal(err)
	}
	client.Client.AuthorizeRequest = nil
	client.Client.DisableRetries = true
	metadata := sdk.NewResourceMetaData(&clients.Client{
		Account:  &clients.ResourceManagerAccount{SubscriptionId: id.SubscriptionId},
		Features: features.UserFeatures{SkipImportCheckOnCreateAndAllowOverwritingExistingResources: true},
		Cdn:      &cdnclient.Client{EdgeActionVersionsClient: client},
	}, Resource{})
	if err := metadata.Encode(&Model{
		Name: "v1", Location: "global", DeploymentType: "file", DefaultVersionEnabled: true,
		CdnEdgeActionID: edgeactionversions.NewEdgeActionID(id.SubscriptionId, id.ResourceGroupName, id.EdgeActionName).ID(),
	}); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(t.Context(), time.Minute)
	defer cancel()
	if err := (Resource{}).Create().Func(ctx, metadata); err != nil {
		t.Fatal(err)
	}
	if !overwritten || metadata.ResourceData.Id() != id.ID() {
		t.Fatalf("expected overwritten version with persisted ID, got overwritten=%t, ID=%q", overwritten, metadata.ResourceData.Id())
	}
}
