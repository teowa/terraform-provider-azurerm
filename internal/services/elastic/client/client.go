// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/elastic/2024-03-01/rules"
	"github.com/hashicorp/go-azure-sdk/resource-manager/elastic/2025-06-01/elasticmonitorresources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	MonitorClient *elasticmonitorresources.ElasticMonitorResourcesClient
	TagRuleClient *rules.RulesClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	monitorClient, err := elasticmonitorresources.NewElasticMonitorResourcesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building monitor client: %+v", err)
	}
	o.Configure(monitorClient.Client, o.Authorizers.ResourceManager)

	tagRuleClient, err := rules.NewRulesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building tag rule client: %+v", err)
	}
	o.Configure(tagRuleClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		MonitorClient: monitorClient,
		TagRuleClient: tagRuleClient,
	}, nil
}
