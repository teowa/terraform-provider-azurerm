// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cdn_edge_action_execution_filter

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2026-10-01/edgeactionexecutionfilters"
	sdkclient "github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	cdnclient "github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/client"
)

func TestUnitExecutionFilterUpdateAppliesAllFieldsAndClearsTags(t *testing.T) {
	versionID := "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test/providers/Microsoft.Cdn/edgeActions/test/versions/v2"
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
			expected := map[string]interface{}{
				"properties": map[string]interface{}{
					"versionId":                            versionID,
					"executionFilterIdentifierHeaderName":  "X-Updated",
					"executionFilterIdentifierHeaderValue": "updated",
				},
				"tags": map[string]interface{}{},
			}
			if !reflect.DeepEqual(payload, expected) {
				t.Errorf("expected all changed fields and empty tags, got %#v", payload)
			}
		} else if req.Method != http.MethodGet {
			t.Errorf("unexpected request: %s", req.Method)
		}
		if _, err := w.Write([]byte(`{"location":"global","properties":{"provisioningState":"Succeeded"}}`)); err != nil {
			t.Error(err)
		}
	}))
	defer server.Close()
	client, err := edgeactionexecutionfilters.NewEdgeActionExecutionFiltersClientWithBaseURI(environments.ResourceManagerAPI(server.URL))
	if err != nil {
		t.Fatal(err)
	}
	client.Client.AuthorizeRequest = nil
	client.Client.DisableRetries = true
	metadata := sdk.NewResourceMetaData(&clients.Client{Cdn: &cdnclient.Client{EdgeActionExecutionFiltersClient: client}}, Resource{})
	metadata.ResourceData.SetId(edgeactionexecutionfilters.NewExecutionFilterID("00000000-0000-0000-0000-000000000000", "test", "test", "filter").ID())
	if err := metadata.Encode(&Model{
		Name: "filter", Location: "global", CdnEdgeActionVersionID: versionID,
		HeaderName: "X-Updated", HeaderValue: "updated",
	}); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(t.Context(), time.Minute)
	defer cancel()
	if err := (Resource{}).Update().Func(ctx, metadata); err != nil {
		t.Fatal(err)
	}
	if !patched {
		t.Fatal("expected filter PATCH")
	}
}
