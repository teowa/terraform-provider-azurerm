// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	network_2025_07_01 "github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2026-03-01/webapplicationfirewallpolicies"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	*network_2025_07_01.Client

	WebApplicationFirewallPoliciesClient *webapplicationfirewallpolicies.WebApplicationFirewallPoliciesClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	webApplicationFirewallPoliciesClient, err := webapplicationfirewallpolicies.NewWebApplicationFirewallPoliciesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Web Application Firewall Policies client: %+v", err)
	}
	o.Configure(webApplicationFirewallPoliciesClient.Client, o.Authorizers.ResourceManager)

	client, err := network_2025_07_01.NewClientWithBaseURI(o.Environment.ResourceManager, func(c *resourcemanager.Client) {
		o.Configure(c, o.Authorizers.ResourceManager)
	})
	if err != nil {
		return nil, fmt.Errorf("building clients for Application Gateway: %+v", err)
	}

	return &Client{
		Client:                               client,
		WebApplicationFirewallPoliciesClient: webApplicationFirewallPoliciesClient,
	}, nil
}
