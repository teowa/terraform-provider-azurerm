// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/elastic/2025-06-01/elasticmonitorresources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/elastic/2025-06-01/tagrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	MonitorClient *elasticmonitorresources.ElasticMonitorResourcesClient
	TagRuleClient *tagrules.TagRulesClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	monitorClient, err := elasticmonitorresources.NewElasticMonitorResourcesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Monitor Client: %+v", err)
	}
	o.Configure(monitorClient.Client, o.Authorizers.ResourceManager)

	tagRuleClient, err := tagrules.NewTagRulesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building TagRule Client: %+v", err)
	}
	o.Configure(tagRuleClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		MonitorClient: monitorClient,
		TagRuleClient: tagRuleClient,
	}, nil
}
