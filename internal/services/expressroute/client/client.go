// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	network_2025_07_01 "github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2026-03-01/expressroutecircuits"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	*network_2025_07_01.Client

	ExpressRouteCircuits *expressroutecircuits.ExpressRouteCircuitsClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	circuitsClient, err := expressroutecircuits.NewExpressRouteCircuitsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building ExpressRoute Circuits client: %+v", err)
	}
	o.Configure(circuitsClient.Client, o.Authorizers.ResourceManager)

	client, err := network_2025_07_01.NewClientWithBaseURI(o.Environment.ResourceManager, func(c *resourcemanager.Client) {
		o.Configure(c, o.Authorizers.ResourceManager)
	})
	if err != nil {
		return nil, fmt.Errorf("building clients for ExpressRoute: %+v", err)
	}

	return &Client{
		Client:               client,
		ExpressRouteCircuits: circuitsClient,
	}, nil
}
