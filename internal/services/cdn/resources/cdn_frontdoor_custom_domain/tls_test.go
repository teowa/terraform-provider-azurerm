// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cdn_frontdoor_custom_domain

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2026-10-01/afddomains"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestFrontDoorTlsGroups(t *testing.T) {
	for _, tc := range []struct {
		policy string
		groups []any
		valid  bool
	}{
		{policy: "Standard", valid: true},
		{policy: "Enhanced", valid: true},
		{policy: "Custom", groups: []any{"X25519", "X25519MLKEM768"}, valid: true},
		{policy: "Custom"},
		{policy: "Standard", groups: []any{"X25519"}},
	} {
		t.Run(tc.policy, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, cdnFrontdoorCustomDomainSchema(), map[string]any{
				"tls": []any{map[string]any{
					"server_tls_group_policy": tc.policy,
					"server_tls_groups":       tc.groups,
				}},
			})
			result, err := expandAfdDomainTlsParameters(d, d.Get("tls").([]any))
			if (err == nil) != tc.valid {
				t.Fatalf("unexpected error: %v", err)
			}
			if !tc.valid {
				return
			}
			if string(pointer.From(result.ServerTlsGroupPolicy)) != tc.policy {
				t.Fatalf("unexpected group policy: %v", result.ServerTlsGroupPolicy)
			}
			flattened := flattenAfdDomainHttpsParameters(result, false)
			dataSource := schema.TestResourceDataRaw(t, RegisterDataSource().Schema, nil)
			if err := dataSource.Set("tls", flattened); err != nil {
				t.Fatal(err)
			}
			if flattened[0].(map[string]any)["server_tls_group_policy"] != tc.policy {
				t.Fatalf("unexpected TLS state: %#v", flattened)
			}
			if len(flattened[0].(map[string]any)["server_tls_groups"].([]string)) != len(tc.groups) {
				t.Fatalf("unexpected TLS groups: %#v", flattened)
			}
		})
	}
	result := flattenAfdDomainHttpsParameters(&afddomains.AFDDomainHTTPSParameters{
		ServerTlsGroupPolicy: pointer.To(afddomains.AfdServerTlsGroupPolicyStandard),
		ServerTlsGroups:      pointer.To([]afddomains.AfdServerTlsGroup{"X25519"}),
	}, false)
	if len(result[0].(map[string]any)["server_tls_groups"].([]string)) != 0 {
		t.Fatal("service-managed groups must not populate user-managed state")
	}
}

func TestFrontDoorMutualTls(t *testing.T) {
	secretID := "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test/providers/Microsoft.Cdn/profiles/test/secrets/test"
	for _, scenario := range afddomains.PossibleValuesForMtlsScenarioType() {
		t.Run(scenario, func(t *testing.T) {
			input := map[string]any{"scenario": scenario}
			if scenario == "ClientCertificateRequiredAndValidated" || scenario == "ClientCertificateValidatedIfPresented" {
				input["cdn_frontdoor_secret_ids"] = []any{secretID}
				input["allowed_fqdns"] = []any{"client.example.com"}
				input["certificate_revocation_check_enabled"] = true
			}
			d := schema.TestResourceDataRaw(t, cdnFrontdoorCustomDomainSchema(), map[string]any{"mutual_tls": []any{input}})
			settings, err := expandFrontDoorCustomDomainMutualTls(d.Get("mutual_tls").([]any), "example.com")
			if err != nil {
				t.Fatal(err)
			}
			body, err := json.Marshal(afddomains.AFDDomainUpdatePropertiesParameters{MtlsSettings: settings})
			if err != nil {
				t.Fatal(err)
			}
			var response afddomains.AFDDomainUpdatePropertiesParameters
			if err := json.Unmarshal(body, &response); err != nil {
				t.Fatal(err)
			}
			flattened, err := flattenFrontDoorCustomDomainMutualTls(response.MtlsSettings)
			if err != nil {
				t.Fatal(err)
			}
			dataSource := schema.TestResourceDataRaw(t, RegisterDataSource().Schema, nil)
			if err := dataSource.Set("mutual_tls", flattened); err != nil {
				t.Fatal(err)
			}
			if err := d.Set("mutual_tls", flattened); err != nil {
				t.Fatal(err)
			}
			roundTrip, err := expandFrontDoorCustomDomainMutualTls(d.Get("mutual_tls").([]any), "example.com")
			if err != nil || !reflect.DeepEqual(settings, roundTrip) {
				t.Fatalf("mutual TLS round trip failed: %#v, %v", roundTrip, err)
			}
			if _, err := expandFrontDoorCustomDomainMutualTls(d.Get("mutual_tls").([]any), "*.example.com"); err == nil {
				t.Fatal("wildcard domains must reject mutual TLS")
			}
		})
	}
}

func TestFrontDoorMutualTlsRemoval(t *testing.T) {
	settings, err := expandFrontDoorCustomDomainMutualTls(nil, "example.com")
	if err != nil {
		t.Fatal(err)
	}
	body, err := json.Marshal(afddomains.AFDDomainUpdatePropertiesParameters{MtlsSettings: settings})
	if err != nil {
		t.Fatal(err)
	}
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(body, &fields); err != nil {
		t.Fatal(err)
	}
	if string(fields["mtlsSettings"]) != "null" {
		t.Fatalf("removing mutual TLS must send explicit null: %s", body)
	}
	var response afddomains.AFDDomainUpdatePropertiesParameters
	if err := json.Unmarshal(body, &response); err != nil {
		t.Fatal(err)
	}
	flattened, err := flattenFrontDoorCustomDomainMutualTls(response.MtlsSettings)
	if err != nil || len(flattened) != 0 {
		t.Fatalf("unexpected disabled mutual TLS state: %#v, %v", flattened, err)
	}
}

func TestFrontDoorMutualTlsValidation(t *testing.T) {
	for _, input := range []map[string]any{
		{"scenario": "ClientCertificateRequiredAndValidated"},
		{"scenario": "ClientCertificateValidatedIfPresented"},
		{"scenario": "CompleteMtlsPassthroughToOrigin", "certificate_revocation_check_enabled": true},
		{"scenario": "ClientCertificateRequiredAndOriginValidates", "allowed_fqdns": []any{"client.example.com"}},
	} {
		d := schema.TestResourceDataRaw(t, cdnFrontdoorCustomDomainSchema(), map[string]any{"mutual_tls": []any{input}})
		if _, err := expandFrontDoorCustomDomainMutualTls(d.Get("mutual_tls").([]any), "example.com"); err == nil {
			t.Fatalf("expected invalid mutual TLS settings to fail: %#v", input)
		}
	}
}
