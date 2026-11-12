// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cdn_frontdoor_secret

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2026-10-01/profiles"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2026-10-01/secrets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceCdnFrontDoorSecretRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Cdn.FrontDoorSecretsClient

	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := secrets.ParseSecretID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.SecretName)
	d.Set("cdn_frontdoor_profile_id", profiles.NewProfileID(id.SubscriptionId, id.ResourceGroupName, id.ProfileName).ID())
	d.Set("cdn_frontdoor_profile_name", id.ProfileName)

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			customerCertificate, err := flattenCdnFrontDoorSecretParameters(ctx, props.Parameters, meta)
			if err != nil {
				return fmt.Errorf("flattening `secret`: %+v", err)
			}

			if err := d.Set("secret", customerCertificate); err != nil {
				return fmt.Errorf("setting `secret`: %+v", err)
			}
		}
	}

	return nil
}

func flattenCdnFrontDoorSecretParameters(ctx context.Context, input secrets.SecretParameters, meta any) ([]any, error) {
	client := meta.(*clients.Client).KeyVault

	results := make([]any, 0)
	if input == nil {
		return results, nil
	}

	result := make(map[string]any)

	var customerCertificate secrets.CustomerCertificateParameters
	var chain *secrets.AfdSecretMtlsCertificateChain
	switch v := input.(type) {
	case secrets.CustomerCertificateParameters:
		customerCertificate = v
	case secrets.AfdSecretMtlsCertificateChain:
		chain = &v
		customerCertificate.SecretSource = v.SecretSource
		customerCertificate.SecretVersion = pointer.To(v.SecretVersion)
	default:
		return results, fmt.Errorf("received an unexpected type (`%T`)", input)
	}

	secretSourceId, err := parse.SecretVersionlessID(pointer.From(customerCertificate.SecretSource.Id))
	if err != nil {
		return nil, err
	}

	keyVaultId := commonids.NewKeyVaultID(secretSourceId.SubscriptionId, secretSourceId.ResourceGroup, secretSourceId.VaultName)
	keyVaultBaseUri, err := client.BaseUriForKeyVault(ctx, keyVaultId)
	if err != nil {
		return nil, fmt.Errorf("looking up Base URI for Certificate %q in %s: %+v", secretSourceId.SecretName, keyVaultId, err)
	}

	itemType := keyvault.NestedItemTypeCertificate
	if chain != nil {
		itemType = keyvault.NestedItemTypeSecret
	}
	keyVaultCertificateId, err := keyvault.NewNestedItemID(*keyVaultBaseUri, itemType, secretSourceId.SecretName, pointer.From(customerCertificate.SecretVersion))
	if err != nil {
		return nil, err
	}

	certificateID := keyVaultCertificateId.ID()
	if pointer.From(customerCertificate.UseLatestVersion) {
		certificateID = keyVaultCertificateId.VersionlessID()
	}

	if chain != nil {
		return []any{map[string]any{
			"mtls_certificate_chain": []any{map[string]any{
				"key_vault_secret_id": certificateID,
				"expiration_date":     pointer.From(chain.ExpirationDate),
			}},
		}}, nil
	}

	result["customer_certificate"] = []any{
		map[string]any{
			"expiration_date":           pointer.From(customerCertificate.ExpirationDate),
			"key_vault_certificate_id":  certificateID,
			"subject_alternative_names": pluginsdk.FlattenSlice(customerCertificate.SubjectAlternativeNames),
		},
	}
	results = append(results, result)

	return results, nil
}
