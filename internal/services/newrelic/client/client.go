// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	newRelic20240301 "github.com/hashicorp/go-azure-sdk/resource-manager/newrelic/2024-03-01"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	*newRelic20240301.Client
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	client, err := newRelic20240301.NewClientWithBaseURI(o.Environment.ResourceManager, func(c *resourcemanager.Client) {
		o.Configure(c, o.Authorizers.ResourceManager)
	})
	if err != nil {
		return nil, fmt.Errorf("building clients for Network: %+v", err)
	}

	return &Client{
		Client: client,
	}, nil
}
