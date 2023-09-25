// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"sync"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/appconfiguration/2023-03-01/configurationstores"
	"github.com/hashicorp/go-azure-sdk/resource-manager/appconfiguration/2023-03-01/replicas"
	resourcesClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/resource/client"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

var (
	ConfigurationStoreCache = map[string]ConfigurationStoreDetails{}
	keysmith                = &sync.RWMutex{}
	lock                    = map[string]*sync.RWMutex{}
)

type ConfigurationStoreDetails struct {
	configurationStoreId string
	replicaName          string
	dataPlaneEndpoint    string
}

func (c Client) AddToCache(configurationStoreId configurationstores.ConfigurationStoreId, replicaName, dataPlaneEndpoint string) {
	cacheKey := c.cacheKeyForConfigurationStore(configurationStoreId.ConfigurationStoreName, replicaName)
	keysmith.Lock()
	ConfigurationStoreCache[cacheKey] = ConfigurationStoreDetails{
		configurationStoreId: configurationStoreId.ID(),
		replicaName:          replicaName,
		dataPlaneEndpoint:    dataPlaneEndpoint,
	}
	keysmith.Unlock()
}

func (c Client) ConfigurationStoreDetailsFromEndpoint(ctx context.Context, resourcesClient *resourcesClient.Client, configurationStoreEndpoint string) (*ConfigurationStoreDetails, error) {
	configurationStoreName, err := c.parseNameFromEndpoint(configurationStoreEndpoint)
	if err != nil {
		return nil, err
	}

	cacheKey := c.cacheKeyForConfigurationStore(*configurationStoreName, "")
	keysmith.Lock()
	if lock[cacheKey] == nil {
		lock[cacheKey] = &sync.RWMutex{}
	}
	keysmith.Unlock()
	lock[cacheKey].Lock()
	defer lock[cacheKey].Unlock()

	if v, ok := ConfigurationStoreCache[cacheKey]; ok {
		return &v, nil
	}

	filter := fmt.Sprintf("resourceType eq 'Microsoft.AppConfiguration/configurationStores' and name eq '%s'", *configurationStoreName)
	result, err := resourcesClient.ResourcesClient.List(ctx, filter, "", utils.Int32(5))
	if err != nil {
		return nil, fmt.Errorf("listing resources matching %q: %+v", filter, err)
	}

	for result.NotDone() {
		for _, v := range result.Values() {
			if v.ID == nil {
				continue
			}

			id, err := configurationstores.ParseConfigurationStoreIDInsensitively(*v.ID)
			if err != nil {
				return nil, fmt.Errorf("parsing %q: %+v", *v.ID, err)
			}
			if !strings.EqualFold(id.ConfigurationStoreName, *configurationStoreName) {
				continue
			}

			resp, err := c.ConfigurationStoresClient.Get(ctx, *id)
			if err != nil {
				return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
			}
			if resp.Model == nil || resp.Model.Properties == nil || resp.Model.Properties.Endpoint == nil {
				return nil, fmt.Errorf("retrieving %s: `model.properties.Endpoint` was nil", *id)
			}

			c.AddToCache(*id, "", *resp.Model.Properties.Endpoint)

			return &ConfigurationStoreDetails{
				configurationStoreId: id.ID(),
				dataPlaneEndpoint:    *resp.Model.Properties.Endpoint,
			}, nil
		}

		if err := result.NextWithContext(ctx); err != nil {
			return nil, fmt.Errorf("iterating over results: %+v", err)
		}
	}

	// check if is a replica endpoint
	if index := strings.LastIndex(*configurationStoreName, "-"); index != -1 {
		replicaName := (*configurationStoreName)[index+1:]
		originalConfigurationStoreName := (*configurationStoreName)[:index]

		filter := fmt.Sprintf("resourceType eq 'Microsoft.AppConfiguration/configurationStores' and name eq '%s'", originalConfigurationStoreName)
		result, err := resourcesClient.ResourcesClient.List(ctx, filter, "", utils.Int32(5))
		if err != nil {
			return nil, fmt.Errorf("listing resources matching %q: %+v", filter, err)
		}

		for result.NotDone() {
			for _, v := range result.Values() {
				if v.ID == nil {
					continue
				}

				id, err := configurationstores.ParseConfigurationStoreIDInsensitively(*v.ID)
				if err != nil {
					return nil, fmt.Errorf("parsing %q: %+v", *v.ID, err)
				}
				if !strings.EqualFold(id.ConfigurationStoreName, *configurationStoreName) {
					continue
				}

				resp, err := c.ConfigurationStoresClient.Get(ctx, *id)
				if err != nil {
					return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
				}
				if resp.Model == nil || resp.Model.Properties == nil || resp.Model.Properties.Endpoint == nil {
					return nil, fmt.Errorf("retrieving %s: `model.properties.Endpoint` was nil", *id)
				}

				replicaId := replicas.NewReplicaID(id.SubscriptionId, id.ResourceGroupName, id.ConfigurationStoreName, replicaName)

				existingReplica, err := c.ReplicasClient.Get(ctx, replicaId)
				if err != nil {
					if !response.WasNotFound(existingReplica.HttpResponse) {
						return nil, fmt.Errorf("retrieving %s: %+v", replicaId, err)
					}
				}

				if response.WasNotFound(existingReplica.HttpResponse) {
					return nil, nil
				}

				if existingReplica.Model == nil || existingReplica.Model.Properties == nil || existingReplica.Model.Properties.Endpoint == nil {
					return nil, fmt.Errorf("retrieving %s: `model.properties.Endpoint` was nil", replicaId)
				}

				c.AddToCache(*id, replicaName, *existingReplica.Model.Properties.Endpoint)

				return &ConfigurationStoreDetails{
					configurationStoreId: id.ID(),
					replicaName:          replicaName,
					dataPlaneEndpoint:    *resp.Model.Properties.Endpoint,
				}, nil
			}

			if err := result.NextWithContext(ctx); err != nil {
				return nil, fmt.Errorf("iterating over results: %+v", err)
			}
		}
	}

	// we haven't found it, but Data Sources and Resources need to handle this error separately
	return nil, nil
}

func (c Client) EndpointForConfigurationStore(ctx context.Context, configurationStoreId configurationstores.ConfigurationStoreId, replicaName string) (*string, error) {
	cacheKey := c.cacheKeyForConfigurationStore(configurationStoreId.ConfigurationStoreName, replicaName)
	keysmith.Lock()
	if lock[cacheKey] == nil {
		lock[cacheKey] = &sync.RWMutex{}
	}
	keysmith.Unlock()
	lock[cacheKey].Lock()
	defer lock[cacheKey].Unlock()

	if v, ok := ConfigurationStoreCache[cacheKey]; ok {
		return &v.dataPlaneEndpoint, nil
	}

	resp, err := c.ConfigurationStoresClient.Get(ctx, configurationStoreId)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s:%+v", configurationStoreId, err)
	}

	if resp.Model == nil || resp.Model.Properties == nil || resp.Model.Properties.Endpoint == nil {
		return nil, fmt.Errorf("retrieving %s: `model.properties.Endpoint` was nil", configurationStoreId)
	}

	if replicaName == "" {
		c.AddToCache(configurationStoreId, "", *resp.Model.Properties.Endpoint)
	} else {
		replicaId := replicas.NewReplicaID(configurationStoreId.SubscriptionId, configurationStoreId.ResourceGroupName, configurationStoreId.ConfigurationStoreName, replicaName)

		existingReplica, err := c.ReplicasClient.Get(ctx, replicaId)
		if err != nil {
			if !response.WasNotFound(existingReplica.HttpResponse) {
				return nil, fmt.Errorf("retrieving %s: %+v", replicaId, err)
			}
		}

		if response.WasNotFound(existingReplica.HttpResponse) {
			return nil, nil
		}

		if existingReplica.Model == nil || existingReplica.Model.Properties == nil || existingReplica.Model.Properties.Endpoint == nil {
			return nil, fmt.Errorf("retrieving %s: `model.properties.Endpoint` was nil", replicaId)
		}

		c.AddToCache(configurationStoreId, replicaName, *existingReplica.Model.Properties.Endpoint)
	}

	return resp.Model.Properties.Endpoint, nil
}

func (c Client) Exists(ctx context.Context, configurationStoreId configurationstores.ConfigurationStoreId, replicaName string) (bool, error) {
	cacheKey := c.cacheKeyForConfigurationStore(configurationStoreId.ConfigurationStoreName, replicaName)
	keysmith.Lock()
	if lock[cacheKey] == nil {
		lock[cacheKey] = &sync.RWMutex{}
	}
	keysmith.Unlock()
	lock[cacheKey].Lock()
	defer lock[cacheKey].Unlock()

	if _, ok := ConfigurationStoreCache[cacheKey]; ok {
		return true, nil
	}

	resp, err := c.ConfigurationStoresClient.Get(ctx, configurationStoreId)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return false, nil
		}
		return false, fmt.Errorf("retrieving %s: %+v", configurationStoreId, err)
	}

	if resp.Model == nil || resp.Model.Properties == nil || resp.Model.Properties.Endpoint == nil {
		return false, fmt.Errorf("retrieving %s: `model.properties.Endpoint` was nil", configurationStoreId)
	}

	if replicaName == "" {
		c.AddToCache(configurationStoreId, "", *resp.Model.Properties.Endpoint)
	} else {
		replicaId := replicas.NewReplicaID(configurationStoreId.SubscriptionId, configurationStoreId.ResourceGroupName, configurationStoreId.ConfigurationStoreName, replicaName)

		existingReplica, err := c.ReplicasClient.Get(ctx, replicaId)
		if err != nil {
			if !response.WasNotFound(existingReplica.HttpResponse) {
				return false, fmt.Errorf("retrieving %s: %+v", replicaId, err)
			}
		}

		if response.WasNotFound(existingReplica.HttpResponse) {
			return false, nil
		}

		if existingReplica.Model == nil || existingReplica.Model.Properties == nil || existingReplica.Model.Properties.Endpoint == nil {
			return false, fmt.Errorf("retrieving %s: `model.properties.Endpoint` was nil", replicaId)
		}

		c.AddToCache(configurationStoreId, replicaName, *existingReplica.Model.Properties.Endpoint)
	}

	return true, nil
}

func (c Client) RemoveFromCache(configurationStoreId configurationstores.ConfigurationStoreId, replicaName string) {
	cacheKey := c.cacheKeyForConfigurationStore(configurationStoreId.ConfigurationStoreName, replicaName)
	keysmith.Lock()
	if lock[cacheKey] == nil {
		lock[cacheKey] = &sync.RWMutex{}
	}
	keysmith.Unlock()
	lock[cacheKey].Lock()
	delete(ConfigurationStoreCache, cacheKey)
	lock[cacheKey].Unlock()
}

func (c Client) cacheKeyForConfigurationStore(configurationStoreName, replicaName string) string {
	if replicaName == "" {
		return strings.ToLower(configurationStoreName)
	}
	return strings.ToLower(fmt.Sprintf("%s-%s", configurationStoreName, replicaName))
}

func (c Client) parseNameFromEndpoint(input string) (*string, error) {
	uri, err := url.ParseRequestURI(input)
	if err != nil {
		return nil, err
	}

	// https://the-appconfiguration.azconfig.io

	segments := strings.Split(uri.Host, ".")
	if len(segments) < 3 || segments[1] != "azconfig" || segments[2] != "io" {
		return nil, fmt.Errorf("expected a URI in the format `https://the-appconfiguration.azconfig.io` but got %q", uri.Host)
	}
	return &segments[0], nil
}
