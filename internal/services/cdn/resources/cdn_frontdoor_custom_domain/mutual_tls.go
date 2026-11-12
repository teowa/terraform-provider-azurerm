// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cdn_frontdoor_custom_domain

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2026-10-01/afddomains"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2026-10-01/secrets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func frontDoorCustomDomainAdvancedTlsCustomizeDiff(_ context.Context, diff *pluginsdk.ResourceDiff, _ any) error {
	config := diff.GetRawConfig()
	if config.IsNull() || !config.IsKnown() {
		return nil
	}
	if tls := config.GetAttr("tls"); !tls.IsNull() && tls.IsWhollyKnown() {
		if blocks := diff.Get("tls").([]any); len(blocks) > 0 && blocks[0] != nil {
			block := blocks[0].(map[string]any)
			if err := validateFrontDoorCustomDomainTlsGroups(afddomains.AfdServerTlsGroupPolicy(block["server_tls_group_policy"].(string)), block["server_tls_groups"].(*pluginsdk.Set).Len()); err != nil {
				return err
			}
		}
	}
	if mtls := config.GetAttr("mutual_tls"); !mtls.IsNull() && mtls.IsWhollyKnown() && config.GetAttr("host_name").IsKnown() {
		_, err := expandFrontDoorCustomDomainMutualTls(diff.Get("mutual_tls").([]any), diff.Get("host_name").(string))
		return err
	}
	return nil
}

func frontDoorCustomDomainMutualTlsSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"scenario": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringInSlice(afddomains.PossibleValuesForMtlsScenarioType(), false),
				},
				"cdn_frontdoor_secret_ids": {
					Type:     pluginsdk.TypeSet,
					Optional: true,
					Elem: &pluginsdk.Schema{
						Type:         pluginsdk.TypeString,
						ValidateFunc: secrets.ValidateSecretID,
					},
				},
				"allowed_fqdns": {
					Type:     pluginsdk.TypeSet,
					Optional: true,
					Elem: &pluginsdk.Schema{
						Type:         pluginsdk.TypeString,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
				"certificate_revocation_check_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  false,
				},
			},
		},
	}
}

func expandFrontDoorCustomDomainMutualTls(input []any, hostName string) (afddomains.AFDDomainMtlsParameters, error) {
	if len(input) == 0 || input[0] == nil {
		return nil, nil
	}
	if strings.Contains(hostName, "*") {
		return nil, errors.New("`mutual_tls` cannot be configured for a wildcard `host_name`")
	}
	v := input[0].(map[string]any)
	scenario := afddomains.MtlsScenarioType(v["scenario"].(string))
	references := make([]afddomains.ResourceReference, 0)
	for _, raw := range v["cdn_frontdoor_secret_ids"].(*pluginsdk.Set).List() {
		id, err := secrets.ParseSecretID(raw.(string))
		if err != nil {
			return nil, err
		}
		references = append(references, afddomains.ResourceReference{Id: pointer.To(id.ID())})
	}
	fqdns := make([]string, 0)
	for _, raw := range v["allowed_fqdns"].(*pluginsdk.Set).List() {
		fqdns = append(fqdns, raw.(string))
	}
	revocationCheck := afddomains.CertificateRevocationCheckEnabledStateDisabled
	if v["certificate_revocation_check_enabled"].(bool) {
		revocationCheck = afddomains.CertificateRevocationCheckEnabledStateEnabled
	}
	switch scenario {
	case afddomains.MtlsScenarioTypeClientCertificateRequiredAndValidated, afddomains.MtlsScenarioTypeClientCertificateValidatedIfPresented:
		if len(references) == 0 {
			return nil, errors.New("`mutual_tls.cdn_frontdoor_secret_ids` must be specified when Front Door validates client certificates")
		}
		if scenario == afddomains.MtlsScenarioTypeClientCertificateRequiredAndValidated {
			return afddomains.ClientCertificateRequiredAndValidatedAdvancedSettings{
				Scenario: scenario, Secrets: references, AllowedFqdns: pointer.To(fqdns), CertificateRevocationCheck: pointer.To(revocationCheck),
			}, nil
		}
		return afddomains.ClientCertificateValidatedIfPresentedAdvancedSettings{
			Scenario: scenario, Secrets: references, AllowedFqdns: pointer.To(fqdns), CertificateRevocationCheck: pointer.To(revocationCheck),
		}, nil
	case afddomains.MtlsScenarioTypeClientCertificateRequiredAndOriginValidates, afddomains.MtlsScenarioTypeCompleteMtlsPassthroughToOrigin:
		if len(references) > 0 || len(fqdns) > 0 || revocationCheck == afddomains.CertificateRevocationCheckEnabledStateEnabled {
			return nil, errors.New("`mutual_tls.cdn_frontdoor_secret_ids`, `mutual_tls.allowed_fqdns`, and `mutual_tls.certificate_revocation_check_enabled` can only be configured when Front Door validates client certificates")
		}
		if scenario == afddomains.MtlsScenarioTypeClientCertificateRequiredAndOriginValidates {
			return afddomains.ClientCertificateRequiredAndOriginValidatesAdvancedSettings{Scenario: scenario}, nil
		}
		return afddomains.CompleteMtlsPassthroughToOriginAdvancedSettings{Scenario: scenario}, nil
	default:
		return nil, fmt.Errorf("unsupported `mutual_tls.scenario` `%s`", scenario)
	}
}

func flattenFrontDoorCustomDomainMutualTls(input afddomains.AFDDomainMtlsParameters) ([]any, error) {
	if input == nil {
		return []any{}, nil
	}
	var references []afddomains.ResourceReference
	var fqdns []string
	var revocationCheck afddomains.CertificateRevocationCheckEnabledState
	scenario := input.AFDDomainMtlsParameters().Scenario
	switch v := input.(type) {
	case afddomains.ClientCertificateRequiredAndValidatedAdvancedSettings:
		references, fqdns, revocationCheck = v.Secrets, pointer.From(v.AllowedFqdns), pointer.From(v.CertificateRevocationCheck)
	case afddomains.ClientCertificateValidatedIfPresentedAdvancedSettings:
		references, fqdns, revocationCheck = v.Secrets, pointer.From(v.AllowedFqdns), pointer.From(v.CertificateRevocationCheck)
	case afddomains.ClientCertificateRequiredAndOriginValidatesAdvancedSettings, afddomains.CompleteMtlsPassthroughToOriginAdvancedSettings:
	case afddomains.RawAFDDomainMtlsParametersImpl:
		if v.Values == nil {
			return []any{}, nil
		}
		return []any{}, fmt.Errorf("unsupported mutual TLS scenario `%s`", v.Type)
	default:
		return []any{}, fmt.Errorf("unsupported mutual TLS settings type `%T`", input)
	}
	ids := make([]string, 0, len(references))
	for _, reference := range references {
		id, err := secrets.ParseSecretIDInsensitively(pointer.From(reference.Id))
		if err != nil {
			return []any{}, err
		}
		ids = append(ids, id.ID())
	}
	return []any{map[string]any{
		"scenario":                             string(scenario),
		"cdn_frontdoor_secret_ids":             ids,
		"allowed_fqdns":                        fqdns,
		"certificate_revocation_check_enabled": revocationCheck == afddomains.CertificateRevocationCheckEnabledStateEnabled,
	}}, nil
}
