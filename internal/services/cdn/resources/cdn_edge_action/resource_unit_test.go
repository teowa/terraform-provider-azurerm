// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cdn_edge_action

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2026-10-01/edgeactions"
	sdkclient "github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	cdnclient "github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/client"
)

func TestUnitEdgeActionUpdateClearsTagsWithoutSku(t *testing.T) {
	var patched bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set(sdkclient.SkipPollingDelayHeader, "true")
		if req.Method == http.MethodPatch {
			patched = true
			var payload map[string]interface{}
			if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
				t.Error(err)
			}
			if !reflect.DeepEqual(payload, map[string]interface{}{"tags": map[string]interface{}{}}) {
				t.Errorf("expected explicit empty tags with no SKU, got %#v", payload)
			}
		} else if req.Method != http.MethodGet {
			t.Errorf("unexpected request: %s", req.Method)
		}
		if _, err := w.Write([]byte(`{"location":"global","sku":{"name":"Standard","tier":"Standard"},"properties":{"provisioningState":"Succeeded","attachments":[]}}`)); err != nil {
			t.Error(err)
		}
	}))
	defer server.Close()
	client, err := edgeactions.NewEdgeActionsClientWithBaseURI(environments.ResourceManagerAPI(server.URL))
	if err != nil {
		t.Fatal(err)
	}
	client.Client.AuthorizeRequest = nil
	client.Client.DisableRetries = true
	metadata := sdk.NewResourceMetaData(&clients.Client{Cdn: &cdnclient.Client{EdgeActionsClient: client}}, Resource{})
	metadata.ResourceData.SetId(edgeactions.NewEdgeActionID("00000000-0000-0000-0000-000000000000", "test", "test").ID())
	if err := metadata.Encode(&Model{Name: "test", ResourceGroupName: "test", Location: "global", SkuName: "Standard", SkuTier: "Standard"}); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(t.Context(), time.Minute)
	defer cancel()
	if err := (Resource{}).Update().Func(ctx, metadata); err != nil {
		t.Fatal(err)
	}
	if !patched {
		t.Fatal("expected tags PATCH")
	}
}

func TestUnitEdgeActionOverwriteFeature(t *testing.T) {
	id := edgeactions.NewEdgeActionID("00000000-0000-0000-0000-000000000000", "test", "test")
	for _, allowOverwrite := range []bool{false, true} {
		t.Run(map[bool]string{false: "requires import", true: "overwrites"}[allowOverwrite], func(t *testing.T) {
			var operations []string
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				operations = append(operations, req.Method)
				w.Header().Set("Content-Type", "application/json")
				w.Header().Set(sdkclient.SkipPollingDelayHeader, "true")
				if _, err := w.Write([]byte(`{"location":"global","sku":{"name":"Standard","tier":"Standard"},"properties":{"provisioningState":"Succeeded","attachments":[]}}`)); err != nil {
					t.Error(err)
				}
			}))
			defer server.Close()
			client, err := edgeactions.NewEdgeActionsClientWithBaseURI(environments.ResourceManagerAPI(server.URL))
			if err != nil {
				t.Fatal(err)
			}
			client.Client.AuthorizeRequest = nil
			client.Client.DisableRetries = true
			metadata := sdk.NewResourceMetaData(&clients.Client{
				Account:  &clients.ResourceManagerAccount{SubscriptionId: id.SubscriptionId},
				Features: features.UserFeatures{SkipImportCheckOnCreateAndAllowOverwritingExistingResources: allowOverwrite},
				Cdn:      &cdnclient.Client{EdgeActionsClient: client},
			}, Resource{})
			if err := metadata.Encode(&Model{Name: "test", ResourceGroupName: "test", Location: "global", SkuName: "Standard", SkuTier: "Standard"}); err != nil {
				t.Fatal(err)
			}
			ctx, cancel := context.WithTimeout(t.Context(), time.Minute)
			defer cancel()
			err = (Resource{}).Create().Func(ctx, metadata)
			if !allowOverwrite {
				if err == nil || !strings.Contains(err.Error(), "needs to be imported") {
					t.Fatalf("expected import error, got %v", err)
				}
				if !reflect.DeepEqual(operations, []string{http.MethodGet}) {
					t.Fatalf("unexpected requests on import conflict: %v", operations)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if len(operations) == 0 || operations[0] != http.MethodPut || metadata.ResourceData.Id() != id.ID() {
				t.Fatalf("expected overwrite without import lookup and persisted ID, got %v, %q", operations, metadata.ResourceData.Id())
			}
		})
	}
}
