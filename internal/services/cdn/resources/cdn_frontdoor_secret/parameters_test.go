// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cdn_frontdoor_secret

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2026-10-01/secrets"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	keyvaultclient "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/client"
)

func TestFrontDoorSecretCertificateTypes(t *testing.T) {
	certificate := []any{map[string]any{
		"key_vault_certificate_id": "https://example.vault.azure.net/certificates/test",
	}}
	chain := []any{map[string]any{
		"key_vault_secret_id": "https://example.vault.azure.net/secrets/test/0123456789abcdef0123456789abcdef",
	}}
	for _, tc := range []struct {
		name   string
		secret map[string]any
		valid  bool
	}{
		{"certificate", map[string]any{"customer_certificate": certificate}, true},
		{"chain", map[string]any{"mtls_certificate_chain": chain}, true},
		{"both", map[string]any{"customer_certificate": certificate, "mtls_certificate_chain": chain}, false},
		{"neither", map[string]any{}, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			config := terraform.NewResourceConfigRaw(map[string]any{
				"name":                     "test",
				"cdn_frontdoor_profile_id": "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test/providers/Microsoft.Cdn/profiles/test",
				"secret":                   []any{tc.secret},
			})
			if diagnostics := RegisterResource().Validate(config); diagnostics.HasError() == tc.valid {
				t.Fatalf("unexpected validation result: %#v", diagnostics)
			}
		})
	}
}

func TestFrontDoorSecretParameters(t *testing.T) {
	vaultID := commonids.NewKeyVaultID("00000000-0000-0000-0000-000000000000", "test", "frontdoor-secret-unit-test")
	client := &clients.Client{
		Account:  &clients.ResourceManagerAccount{SubscriptionId: vaultID.SubscriptionId},
		KeyVault: &keyvaultclient.Client{},
	}
	client.KeyVault.AddToCache(vaultID, "https://frontdoor-secret-unit-test.vault.azure.net/")
	t.Cleanup(func() { client.KeyVault.Purge(vaultID) })
	for _, tc := range []struct {
		block string
		field string
		id    string
	}{
		{"customer_certificate", "key_vault_certificate_id", "https://frontdoor-secret-unit-test.vault.azure.net/certificates/test"},
		{"customer_certificate", "key_vault_certificate_id", "https://frontdoor-secret-unit-test.vault.azure.net/certificates/test/0123456789abcdef0123456789abcdef"},
		{"mtls_certificate_chain", "key_vault_secret_id", "https://frontdoor-secret-unit-test.vault.azure.net/secrets/test/0123456789abcdef0123456789abcdef"},
	} {
		t.Run(tc.id, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, cdnFrontdoorSecretSchema(), map[string]any{
				"secret": []any{map[string]any{tc.block: []any{map[string]any{tc.field: tc.id}}}},
			})
			expanded, err := expandCdnFrontDoorSecretParameters(context.Background(), d.Get("secret").([]any), client)
			if err != nil {
				t.Fatal(err)
			}
			body, err := json.Marshal(expanded)
			if err != nil {
				t.Fatal(err)
			}
			response, err := secrets.UnmarshalSecretParametersImplementation(body)
			if err != nil {
				t.Fatal(err)
			}
			flattened, err := flattenCdnFrontDoorSecretParameters(context.Background(), response, client)
			if err != nil {
				t.Fatal(err)
			}
			dataSource := schema.TestResourceDataRaw(t, RegisterDataSource().Schema, nil)
			if err := dataSource.Set("secret", flattened); err != nil {
				t.Fatal(err)
			}
			if err := d.Set("secret", flattened); err != nil {
				t.Fatal(err)
			}
			result, err := expandCdnFrontDoorSecretParameters(context.Background(), d.Get("secret").([]any), client)
			if err != nil || !reflect.DeepEqual(expanded, result) {
				t.Fatalf("secret round trip failed: %#v, %v", result, err)
			}
			if tc.block == "mtls_certificate_chain" {
				chain, ok := expanded.(secrets.AfdSecretMtlsCertificateChain)
				if !ok || chain.SecretVersion != "0123456789abcdef0123456789abcdef" || *chain.SecretSource.Id != vaultID.ID()+"/secrets/test" {
					t.Fatalf("unexpected certificate chain: %#v", expanded)
				}
			}
		})
	}
}
