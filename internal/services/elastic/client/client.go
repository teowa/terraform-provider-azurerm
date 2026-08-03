// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/elastic/2023-06-01/monitorsresource"
	"github.com/hashicorp/go-azure-sdk/resource-manager/elastic/2023-06-01/rules"
	"github.com/hashicorp/go-azure-sdk/resource-manager/elastic/2025-06-01/elasticmonitorresources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	HostedSearchMonitorClient *elasticmonitorresources.ElasticMonitorResourcesClient
	MonitorClient             *monitorsresource.MonitorsResourceClient
	TagRuleClient             *rules.RulesClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	hostedSearchMonitorClient, err := elasticmonitorresources.NewElasticMonitorResourcesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Hosted Search Monitor Client: %+v", err)
	}
	o.Configure(hostedSearchMonitorClient.Client, o.Authorizers.ResourceManager)

	monitorClient, err := monitorsresource.NewMonitorsResourceClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Monitor Client: %+v", err)
	}
	o.Configure(monitorClient.Client, o.Authorizers.ResourceManager)

	tagRuleClient, err := rules.NewRulesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building TagRule Client: %+v", err)
	}
	o.Configure(tagRuleClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		HostedSearchMonitorClient: hostedSearchMonitorClient,
		MonitorClient:             monitorClient,
		TagRuleClient:             tagRuleClient,
	}, nil
}
