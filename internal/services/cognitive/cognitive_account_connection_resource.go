// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cognitive

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2025-06-01/accountconnectionresource"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2025-06-01/cognitiveservicesaccounts"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceCognitiveAccountConnection() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceCognitiveAccountConnectionCreate,
		Read:   resourceCognitiveAccountConnectionRead,
		Update: resourceCognitiveAccountConnectionUpdate,
		Delete: resourceCognitiveAccountConnectionDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := accountconnectionresource.ParseConnectionID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"cognitive_account_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: cognitiveservicesaccounts.ValidateAccountID,
			},

			"auth_type": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"AAD",
					"AccessKey",
					"AccountKey",
					"ApiKey",
					"CustomKeys",
					"ManagedIdentity",
					"None",
					"OAuth2",
					"PAT",
					"SAS",
					"ServicePrincipal",
					"UsernamePassword",
				}, false),
			},

			"category": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"ADLSGen2",
					"AIServices",
					"AmazonMws",
					"AmazonRdsForOracle",
					"AmazonRdsForSqlServer",
					"AmazonRedshift",
					"AmazonS3Compatible",
					"ApiKey",
					"AzureBlob",
					"AzureDataExplorer",
					"AzureDatabricksDeltaLake",
					"AzureMariaDb",
					"AzureMySqlDb",
					"AzureOneLake",
					"AzureOpenAI",
					"AzurePostgresDb",
					"AzureSqlDb",
					"AzureSqlMi",
					"AzureSynapseAnalytics",
					"AzureTableStorage",
					"BingLLMSearch",
					"Cassandra",
					"CognitiveSearch",
					"CognitiveService",
					"Concur",
					"ContainerRegistry",
					"CosmosDb",
					"CosmosDbMongoDbApi",
					"Couchbase",
					"CustomKeys",
					"Db2",
					"Drill",
					"Dynamics",
					"DynamicsAx",
					"DynamicsCrm",
					"Elasticsearch",
					"Eloqua",
					"FileServer",
					"FtpServer",
					"GenericContainerRegistry",
					"GenericHttp",
					"GenericRest",
					"Git",
					"GoogleAdWords",
					"GoogleBigQuery",
					"GoogleCloudStorage",
					"Greenplum",
					"Hbase",
					"Hdfs",
					"Hive",
					"Hubspot",
					"Impala",
					"Informix",
					"Jira",
					"Magento",
					"ManagedOnlineEndpoint",
					"MariaDb",
					"Marketo",
					"MicrosoftAccess",
					"MongoDbAtlas",
					"MongoDbV2",
					"MySql",
					"Netezza",
					"ODataRest",
					"Odbc",
					"Office365",
					"OpenAI",
					"Oracle",
					"OracleCloudStorage",
					"OracleServiceCloud",
					"PayPal",
					"Phoenix",
					"Pinecone",
					"PostgreSql",
					"Presto",
					"PythonFeed",
					"QuickBooks",
					"Redis",
					"Responsys",
					"S3",
					"Salesforce",
					"SalesforceMarketingCloud",
					"SalesforceServiceCloud",
					"SapBw",
					"SapCloudForCustomer",
					"SapEcc",
					"SapHana",
					"SapOpenHub",
					"SapTable",
					"Serp",
					"Serverless",
					"ServiceNow",
					"Sftp",
					"SharePointOnlineList",
					"Shopify",
					"Snowflake",
					"Spark",
					"SqlServer",
					"Square",
					"Sybase",
					"Teradata",
					"Vertica",
					"WebTable",
					"Xero",
					"Zoho",
				}, false),
			},

			"target": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"is_shared_to_all": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"shared_user_list": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},

			"use_workspace_managed_identity": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"metadata": {
				Type:     pluginsdk.TypeMap,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"api_key": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Sensitive:    true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"api_version": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"api_base": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"username": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"password": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Sensitive:    true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"client_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsUUID,
			},

			"client_secret": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Sensitive:    true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"tenant_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsUUID,
			},

			"access_token": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Sensitive:    true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"refresh_token": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Sensitive:    true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"pat": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Sensitive:    true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"sas_token": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Sensitive:    true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"account_key": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Sensitive:    true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"subscription_key": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Sensitive:    true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"service_endpoint": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"custom_keys": {
				Type:     pluginsdk.TypeMap,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type:      pluginsdk.TypeString,
					Sensitive: true,
				},
			},
		},
	}
}

func resourceCognitiveAccountConnectionCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cognitive.AccountConnectionsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	cognitiveAccountId, err := cognitiveservicesaccounts.ParseAccountID(d.Get("cognitive_account_id").(string))
	if err != nil {
		return err
	}

	id := accountconnectionresource.NewConnectionID(cognitiveAccountId.SubscriptionId, cognitiveAccountId.ResourceGroupName, cognitiveAccountId.AccountName, d.Get("name").(string))

	existing, err := client.AccountConnectionsGet(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
	}

	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_cognitive_account_connection", id.ID())
	}

	authType := accountconnectionresource.ConnectionAuthType(d.Get("auth_type").(string))
	
	props, err := expandCognitiveAccountConnectionProperties(d, authType)
	if err != nil {
		return fmt.Errorf("expanding connection properties: %+v", err)
	}

	connection := accountconnectionresource.ConnectionPropertiesV2BasicResource{
		Properties: props,
	}

	if _, err := client.AccountConnectionsCreate(ctx, id, connection); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceCognitiveAccountConnectionRead(d, meta)
}

func resourceCognitiveAccountConnectionRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cognitive.AccountConnectionsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := accountconnectionresource.ParseConnectionID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.AccountConnectionsGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.ConnectionName)
	
	cognitiveAccountId := cognitiveservicesaccounts.NewAccountID(id.SubscriptionId, id.ResourceGroupName, id.AccountName)
	d.Set("cognitive_account_id", cognitiveAccountId.ID())

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			baseProps := props.ConnectionPropertiesV2()
			d.Set("auth_type", string(baseProps.AuthType))
			
			if baseProps.Category != nil {
				d.Set("category", string(*baseProps.Category))
			}
			
			d.Set("target", baseProps.Target)
			d.Set("is_shared_to_all", baseProps.IsSharedToAll)
			d.Set("shared_user_list", utils.FlattenStringSlice(baseProps.SharedUserList))
			d.Set("use_workspace_managed_identity", baseProps.UseWorkspaceManagedIdentity)
			
			if baseProps.Metadata != nil {
				d.Set("metadata", *baseProps.Metadata)
			}
		}
	}

	return nil
}

func resourceCognitiveAccountConnectionUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cognitive.AccountConnectionsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := accountconnectionresource.ParseConnectionID(d.Id())
	if err != nil {
		return err
	}

	authType := accountconnectionresource.ConnectionAuthType(d.Get("auth_type").(string))
	
	props, err := expandCognitiveAccountConnectionProperties(d, authType)
	if err != nil {
		return fmt.Errorf("expanding connection properties: %+v", err)
	}

	update := accountconnectionresource.ConnectionUpdateContent{
		Properties: props,
	}

	if _, err := client.AccountConnectionsUpdate(ctx, *id, update); err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	return resourceCognitiveAccountConnectionRead(d, meta)
}

func resourceCognitiveAccountConnectionDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cognitive.AccountConnectionsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := accountconnectionresource.ParseConnectionID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.AccountConnectionsDelete(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}

func expandCognitiveAccountConnectionProperties(d *pluginsdk.ResourceData, authType accountconnectionresource.ConnectionAuthType) (accountconnectionresource.ConnectionPropertiesV2, error) {
	baseProps := accountconnectionresource.BaseConnectionPropertiesV2Impl{
		AuthType:                    authType,
		Target:                      utils.String(d.Get("target").(string)),
		IsSharedToAll:               utils.Bool(d.Get("is_shared_to_all").(bool)),
		SharedUserList:              utils.ExpandStringSlice(d.Get("shared_user_list").([]interface{})),
		UseWorkspaceManagedIdentity: utils.Bool(d.Get("use_workspace_managed_identity").(bool)),
	}

	if v, ok := d.GetOk("category"); ok {
		category := accountconnectionresource.ConnectionCategory(v.(string))
		baseProps.Category = &category
	}

	if v, ok := d.GetOk("metadata"); ok {
		metadata := v.(map[string]interface{})
		metadataMap := make(map[string]string)
		for k, v := range metadata {
			metadataMap[k] = v.(string)
		}
		baseProps.Metadata = &metadataMap
	}

	switch authType {
	case accountconnectionresource.ConnectionAuthTypeApiKey:
		credentials := &accountconnectionresource.ConnectionApiKey{}
		if v, ok := d.GetOk("api_key"); ok {
			credentials.Key = utils.String(v.(string))
		}
		return accountconnectionresource.ApiKeyAuthConnectionProperties{
			BaseConnectionPropertiesV2Impl: baseProps,
			Credentials:                     credentials,
		}, nil

	case accountconnectionresource.ConnectionAuthTypeNone:
		return accountconnectionresource.NoneAuthTypeConnectionProperties{
			BaseConnectionPropertiesV2Impl: baseProps,
		}, nil

	case accountconnectionresource.ConnectionAuthTypeManagedIdentity:
		return accountconnectionresource.ManagedIdentityAuthTypeConnectionProperties{
			BaseConnectionPropertiesV2Impl: baseProps,
		}, nil

	case accountconnectionresource.ConnectionAuthTypeUsernamePassword:
		credentials := &accountconnectionresource.ConnectionUsernamePassword{}
		if v, ok := d.GetOk("username"); ok {
			credentials.Username = utils.String(v.(string))
		}
		if v, ok := d.GetOk("password"); ok {
			credentials.Password = utils.String(v.(string))
		}
		return accountconnectionresource.UsernamePasswordAuthTypeConnectionProperties{
			BaseConnectionPropertiesV2Impl: baseProps,
			Credentials:                     credentials,
		}, nil

	case accountconnectionresource.ConnectionAuthTypeServicePrincipal:
		credentials := &accountconnectionresource.ConnectionServicePrincipal{}
		if v, ok := d.GetOk("client_id"); ok {
			credentials.ClientId = utils.String(v.(string))
		}
		if v, ok := d.GetOk("client_secret"); ok {
			credentials.ClientSecret = utils.String(v.(string))
		}
		if v, ok := d.GetOk("tenant_id"); ok {
			credentials.TenantId = utils.String(v.(string))
		}
		return accountconnectionresource.ServicePrincipalAuthTypeConnectionProperties{
			BaseConnectionPropertiesV2Impl: baseProps,
			Credentials:                     credentials,
		}, nil

	case accountconnectionresource.ConnectionAuthTypePAT:
		credentials := &accountconnectionresource.ConnectionAccessKey{}
		if v, ok := d.GetOk("pat"); ok {
			credentials.AccessKey = utils.String(v.(string))
		}
		return accountconnectionresource.PATAuthTypeConnectionProperties{
			BaseConnectionPropertiesV2Impl: baseProps,
			Credentials:                     credentials,
		}, nil

	case accountconnectionresource.ConnectionAuthTypeSAS:
		credentials := &accountconnectionresource.ConnectionSharedAccessSignature{}
		if v, ok := d.GetOk("sas_token"); ok {
			credentials.Sas = utils.String(v.(string))
		}
		return accountconnectionresource.SASAuthTypeConnectionProperties{
			BaseConnectionPropertiesV2Impl: baseProps,
			Credentials:                     credentials,
		}, nil

	case accountconnectionresource.ConnectionAuthTypeAccountKey:
		credentials := &accountconnectionresource.ConnectionAccessKey{}
		if v, ok := d.GetOk("account_key"); ok {
			credentials.AccessKey = utils.String(v.(string))
		}
		return accountconnectionresource.AccountKeyAuthTypeConnectionProperties{
			BaseConnectionPropertiesV2Impl: baseProps,
			Credentials:                     credentials,
		}, nil

	case accountconnectionresource.ConnectionAuthTypeAccessKey:
		credentials := &accountconnectionresource.ConnectionAccessKey{}
		if v, ok := d.GetOk("subscription_key"); ok {
			credentials.AccessKey = utils.String(v.(string))
		}
		return accountconnectionresource.AccessKeyAuthTypeConnectionProperties{
			BaseConnectionPropertiesV2Impl: baseProps,
			Credentials:                     credentials,
		}, nil

	case accountconnectionresource.ConnectionAuthTypeCustomKeys:
		customKeys := &accountconnectionresource.CustomKeys{}
		if v, ok := d.GetOk("custom_keys"); ok {
			keysMap := v.(map[string]interface{})
			keys := make(map[string]string)
			for k, v := range keysMap {
				keys[k] = v.(string)
			}
			customKeys.Keys = &keys
		}
		return accountconnectionresource.CustomKeysConnectionProperties{
			BaseConnectionPropertiesV2Impl: baseProps,
			Credentials:                     customKeys,
		}, nil

	case accountconnectionresource.ConnectionAuthTypeOAuthTwo:
		return accountconnectionresource.OAuth2AuthTypeConnectionProperties{
			BaseConnectionPropertiesV2Impl: baseProps,
		}, nil

	case accountconnectionresource.ConnectionAuthTypeAAD:
		return accountconnectionresource.AADAuthTypeConnectionProperties{
			BaseConnectionPropertiesV2Impl: baseProps,
		}, nil

	default:
		return nil, fmt.Errorf("unsupported auth type: %s", authType)
	}
}
