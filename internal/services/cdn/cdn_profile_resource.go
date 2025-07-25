// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cdn

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2020-09-01/cdn" // nolint: staticcheck
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceCdnProfile() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceCdnProfileCreate,
		Read:   resourceCdnProfileRead,
		Update: resourceCdnProfileUpdate,
		Delete: resourceCdnProfileDelete,

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.CdnProfileV0ToV1{},
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.ProfileID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"location": commonschema.Location(),

			"resource_group_name": commonschema.ResourceGroupName(),

			"sku": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(cdn.SkuNameStandardAkamai),
					string(cdn.SkuNameStandardChinaCdn),
					string(cdn.SkuNameStandardVerizon),
					string(cdn.SkuNameStandardMicrosoft),
					string(cdn.SkuNamePremiumVerizon),
				}, false),
			},

			"tags": tags.Schema(),
		},

		CustomizeDiff: pluginsdk.CustomizeDiffShim(func(ctx context.Context, d *pluginsdk.ResourceDiff, v interface{}) error {
			sku := d.Get("sku").(string)

			if IsCdnFullyRetired() {
				return fmt.Errorf("%s", FullyRetiredMessage)
			}

			switch sku {
			case string(cdn.SkuNameStandardAkamai):
				// The Akamai sku was retired on 31 October 2023...
				if IsCdnStandardAkamaiDeprecatedForCreation() {
					return fmt.Errorf("%s", AkamaiDeprecationMessage)
				}
			case string(cdn.SkuNameStandardVerizon), string(cdn.SkuNamePremiumVerizon):
				// The Verizon skus were retired on 15 January 2025...
				if IsCdnStandardVerizonPremiumVerizonDeprecatedForCreation() {
					return fmt.Errorf("%s", VerizonDeprecationMessage)
				}
			case string(cdn.SkuNameStandardMicrosoft), string(cdn.SkuNameStandardChinaCdn):
				// The only currently supported sku's for CDN are 'Standard_Microsoft' and 'Standard_ChinaCdn' which
				// will also be deprecated as of October 1, 2025, but can still be updated until September 30, 2027.
				// First, check to see if any of the force new fields have changed after the October 1, 2025 deprecation date,
				// if they have we need to block the update because the force new will cause the deployments to fail...
				if IsCdnDeprecatedForCreation() && d.HasChanges("name", "location", "resource_group_name", "sku") {
					return fmt.Errorf("%s", CreateDeprecationMessage)
				}
			}

			return nil
		}),
	}
}

func resourceCdnProfileCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.ProfilesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewProfileID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_cdn_profile", id.ID())
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	sku := d.Get("sku").(string)
	t := d.Get("tags").(map[string]interface{})

	cdnProfile := cdn.Profile{
		Location: &location,
		Tags:     tags.Expand(t),
		Sku: &cdn.Sku{
			Name: cdn.SkuName(sku),
		},
	}

	future, err := client.Create(ctx, id.ResourceGroup, id.Name, cdnProfile)
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the creation of %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceCdnProfileRead(d, meta)
}

func resourceCdnProfileUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.ProfilesClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	if !d.HasChange("tags") {
		return nil
	}

	id, err := parse.ProfileID(d.Id())
	if err != nil {
		return err
	}

	newTags := d.Get("tags").(map[string]interface{})

	props := cdn.ProfileUpdateParameters{
		Tags: tags.Expand(newTags),
	}

	future, err := client.Update(ctx, id.ResourceGroup, id.Name, props)
	if err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the update of %s: %+v", *id, err)
	}

	return resourceCdnProfileRead(d, meta)
}

func resourceCdnProfileRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.ProfilesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ProfileID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("making Read request on Azure CDN Profile %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("location", location.NormalizeNilable(resp.Location))

	if sku := resp.Sku; sku != nil {
		d.Set("sku", string(sku.Name))
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceCdnProfileDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.ProfilesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ProfileID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the deletion of %s: %+v", *id, err)
	}

	return err
}
