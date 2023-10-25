package client

import (
	"fmt"

	v20230904 "github.com/hashicorp/go-azure-sdk/resource-manager/redhatopenshift/2023-09-04"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	*v20230904.Client
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	client, err := v20230904.NewClientWithBaseURI(o.Environment.ResourceManager, func(c *resourcemanager.Client) {
		o.Configure(c, o.Authorizers.ResourceManager)
	})
	if err != nil {
		return nil, fmt.Errorf("building 2021-11-01 client: %+v", err)
	}

	return &Client{
		Client: client,
	}, nil
}
