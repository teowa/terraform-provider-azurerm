package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/discovery/2026-06-01/bookshelves"
	"github.com/hashicorp/go-azure-sdk/resource-manager/discovery/2026-06-01/chatmodeldeployments"
	"github.com/hashicorp/go-azure-sdk/resource-manager/discovery/2026-06-01/nodepools"
	"github.com/hashicorp/go-azure-sdk/resource-manager/discovery/2026-06-01/projects"
	"github.com/hashicorp/go-azure-sdk/resource-manager/discovery/2026-06-01/storageassets"
	"github.com/hashicorp/go-azure-sdk/resource-manager/discovery/2026-06-01/storagecontainers"
	"github.com/hashicorp/go-azure-sdk/resource-manager/discovery/2026-06-01/supercomputers"
	"github.com/hashicorp/go-azure-sdk/resource-manager/discovery/2026-06-01/tools"
	"github.com/hashicorp/go-azure-sdk/resource-manager/discovery/2026-06-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	BookshelvesClient          *bookshelves.BookshelvesClient
	ChatModelDeploymentsClient *chatmodeldeployments.ChatModelDeploymentsClient
	NodePoolsClient            *nodepools.NodePoolsClient
	ProjectsClient             *projects.ProjectsClient
	StorageAssetsClient        *storageassets.StorageAssetsClient
	StorageContainersClient    *storagecontainers.StorageContainersClient
	SupercomputersClient       *supercomputers.SupercomputersClient
	ToolsClient                *tools.ToolsClient
	WorkspacesClient           *workspaces.WorkspacesClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	bookshelvesClient, err := bookshelves.NewBookshelvesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Bookshelves client: %+v", err)
	}
	o.Configure(bookshelvesClient.Client, o.Authorizers.ResourceManager)

	chatModelDeploymentsClient, err := chatmodeldeployments.NewChatModelDeploymentsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building ChatModelDeployments client: %+v", err)
	}
	o.Configure(chatModelDeploymentsClient.Client, o.Authorizers.ResourceManager)

	nodePoolsClient, err := nodepools.NewNodePoolsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building NodePools client: %+v", err)
	}
	o.Configure(nodePoolsClient.Client, o.Authorizers.ResourceManager)

	projectsClient, err := projects.NewProjectsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Projects client: %+v", err)
	}
	o.Configure(projectsClient.Client, o.Authorizers.ResourceManager)

	storageAssetsClient, err := storageassets.NewStorageAssetsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building StorageAssets client: %+v", err)
	}
	o.Configure(storageAssetsClient.Client, o.Authorizers.ResourceManager)

	storageContainersClient, err := storagecontainers.NewStorageContainersClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building StorageContainers client: %+v", err)
	}
	o.Configure(storageContainersClient.Client, o.Authorizers.ResourceManager)

	supercomputersClient, err := supercomputers.NewSupercomputersClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Supercomputers client: %+v", err)
	}
	o.Configure(supercomputersClient.Client, o.Authorizers.ResourceManager)

	toolsClient, err := tools.NewToolsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Tools client: %+v", err)
	}
	o.Configure(toolsClient.Client, o.Authorizers.ResourceManager)

	workspacesClient, err := workspaces.NewWorkspacesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Workspaces client: %+v", err)
	}
	o.Configure(workspacesClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		BookshelvesClient:          bookshelvesClient,
		ChatModelDeploymentsClient: chatModelDeploymentsClient,
		NodePoolsClient:            nodePoolsClient,
		ProjectsClient:             projectsClient,
		StorageAssetsClient:        storageAssetsClient,
		StorageContainersClient:    storageContainersClient,
		SupercomputersClient:       supercomputersClient,
		ToolsClient:                toolsClient,
		WorkspacesClient:           workspacesClient,
	}, nil
}
