// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	privateendpoints "github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-01-01/privateendpoints"
	privateendpoints20250701 "github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/privateendpoints"
	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
)

type privateEndpointCreateOrUpdateModel20250701 struct {
	Location   *string                                                   `json:"location,omitempty"`
	Properties *privateendpoints20250701.CommonPrivateEndpointProperties `json:"properties,omitempty"`
	Tags       *map[string]string                                        `json:"tags,omitempty"`
}

type privateEndpointCreateOrUpdateOperationResponse20250701 struct {
	HttpResponse *http.Response
	Poller       pollers.Poller
}

func expandPrivateEndpointCreateOrUpdateModel20250701(input privateendpoints.PrivateEndpoint, billingSku string) (privateEndpointCreateOrUpdateModel20250701, error) {
	payload, err := json.Marshal(input)
	if err != nil {
		return privateEndpointCreateOrUpdateModel20250701{}, fmt.Errorf("marshaling private endpoint payload: %+v", err)
	}

	var model privateEndpointCreateOrUpdateModel20250701
	if err := json.Unmarshal(payload, &model); err != nil {
		return privateEndpointCreateOrUpdateModel20250701{}, fmt.Errorf("unmarshaling private endpoint payload: %+v", err)
	}

	if model.Properties != nil && billingSku != "" {
		value := privateendpoints20250701.PrivateEndpointBillingSku(billingSku)
		model.Properties.BillingSku = &value
	}

	return model, nil
}

func createOrUpdatePrivateEndpoint20250701(ctx context.Context, client20250701 *privateendpoints20250701.PrivateEndpointsClient, resourceID string, input privateEndpointCreateOrUpdateModel20250701) (result privateEndpointCreateOrUpdateOperationResponse20250701, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusCreated,
			http.StatusOK,
		},
		HttpMethod: http.MethodPut,
		Path:       resourceID,
	}

	req, err := client20250701.Client.NewRequest(ctx, opts)
	if err != nil {
		return result, err
	}

	if err = req.Marshal(input); err != nil {
		return result, err
	}

	var resp *client.Response
	resp, err = req.Execute(ctx)
	if resp != nil {
		result.HttpResponse = resp.Response
	}
	if err != nil {
		return result, err
	}

	result.Poller, err = resourcemanager.PollerFromResponse(resp, client20250701.Client)
	return result, err
}
