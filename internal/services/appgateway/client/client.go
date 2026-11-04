// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/webapplicationfirewallpolicies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2026-03-01/applicationgateways"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	ApplicationGateways            *applicationgateways.ApplicationGatewaysClient
	WebApplicationFirewallPolicies *webapplicationfirewallpolicies.WebApplicationFirewallPoliciesClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	applicationGatewaysClient, err := applicationgateways.NewApplicationGatewaysClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Application Gateway client: %+v", err)
	}
	o.Configure(applicationGatewaysClient.Client, o.Authorizers.ResourceManager)

	webApplicationFirewallPoliciesClient, err := webapplicationfirewallpolicies.NewWebApplicationFirewallPoliciesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Web Application Firewall Policy client: %+v", err)
	}
	o.Configure(webApplicationFirewallPoliciesClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		ApplicationGateways:            applicationGatewaysClient,
		WebApplicationFirewallPolicies: webApplicationFirewallPoliciesClient,
	}, nil
}
