// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-01-01/loadbalancers"
	loadbalancers20250701 "github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/loadbalancers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	LoadBalancersClient          *loadbalancers.LoadBalancersClient
	LoadBalancersClientV20250701 *loadbalancers20250701.LoadBalancersClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	loadBalancersClient, err := loadbalancers.NewLoadBalancersClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building loadBalancers client: %+v", err)
	}
	o.Configure(loadBalancersClient.Client, o.Authorizers.ResourceManager)

	loadBalancersClientV20250701, err := loadbalancers20250701.NewLoadBalancersClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building loadBalancers client (2025-07-01): %+v", err)
	}
	o.Configure(loadBalancersClientV20250701.Client, o.Authorizers.ResourceManager)

	return &Client{
		LoadBalancersClient:          loadBalancersClient,
		LoadBalancersClientV20250701: loadBalancersClientV20250701,
	}, nil
}
