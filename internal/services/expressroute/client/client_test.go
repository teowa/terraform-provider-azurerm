// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"net/http"
	"testing"
	"time"

	sdkclient "github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

func TestExpressRouteClientAPIVersions(t *testing.T) {
	client, err := NewClient(&common.ClientOptions{
		Authorizers: &common.Authorizers{},
		Environment: *environments.AzurePublic(),
	})
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name       string
		client     *resourcemanager.Client
		apiVersion string
	}{
		{
			name:       "circuits",
			client:     client.ExpressRouteCircuits.Client,
			apiVersion: "2026-03-01",
		},
		{
			name:       "ports",
			client:     client.ExpressRoutePorts.Client,
			apiVersion: "2025-07-01",
		},
		{
			name:       "connections",
			client:     client.ExpressRouteCircuitConnections.Client,
			apiVersion: "2025-07-01",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			defer cancel()

			request, err := test.client.NewRequest(ctx, sdkclient.RequestOptions{
				ContentType:         "application/json; charset=utf-8",
				ExpectedStatusCodes: []int{http.StatusOK},
				HttpMethod:          http.MethodGet,
				Path:                "/subscriptions/00000000-0000-0000-0000-000000000000",
			})
			if err != nil {
				t.Fatal(err)
			}
			if got := request.URL.Query().Get("api-version"); got != test.apiVersion {
				t.Fatalf("expected API version %q, got %q", test.apiVersion, got)
			}
		})
	}
}
