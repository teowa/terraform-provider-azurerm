// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cdn_frontdoor_custom_domain

import (
	"errors"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2026-10-01/afddomains"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2026-10-01/secrets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func expandAfdDomainTlsParameters(d *pluginsdk.ResourceData, input []any) (*afddomains.AFDDomainHTTPSParameters, error) {
	// NOTE: With the Frontdoor service, they do not treat an empty object like an empty object
	// if it is not nil they assume it is fully defined and then end up throwing errors when they
	// attempt to get a value from one of the fields.
	if len(input) == 0 || input[0] == nil {
		return nil, nil
	}

	v := input[0].(map[string]any)

	certType := v["certificate_type"].(string)
	secretRaw := v["cdn_frontdoor_secret_id"].(string)
	secretWasConfigured := false

	if rawConfig := d.GetRawConfig(); !rawConfig.IsNull() {
		tlsConfig := rawConfig.GetAttr("tls")
		if !tlsConfig.IsNull() && tlsConfig.LengthInt() > 0 {
			tlsBlock := tlsConfig.AsValueSlice()[0]
			if !tlsBlock.IsNull() {
				secretWasConfigured = !tlsBlock.GetAttr("cdn_frontdoor_secret_id").IsNull()
			}
		}
	}

	minTlsVersion := v["minimum_version"].(string)

	cipherSuiteRaw := v["cipher_suite"].([]any)

	tls := afddomains.AFDDomainHTTPSParameters{
		CertificateType: afddomains.AfdCertificateType(certType),
	}

	groupPolicy := afddomains.AfdServerTlsGroupPolicy(v["server_tls_group_policy"].(string))
	rawGroups := v["server_tls_groups"].(*pluginsdk.Set)
	if err := validateFrontDoorCustomDomainTlsGroups(groupPolicy, rawGroups.Len()); err != nil {
		return nil, err
	}
	groups := make([]afddomains.AfdServerTlsGroup, 0)
	for _, group := range rawGroups.List() {
		groups = append(groups, afddomains.AfdServerTlsGroup(group.(string)))
	}
	tls.ServerTlsGroupPolicy = pointer.To(groupPolicy)
	tls.ServerTlsGroups = pointer.To(groups)

	if certType == string(afddomains.AfdCertificateTypeCustomerCertificate) && secretRaw == "" {
		return nil, errors.New("the `cdn_frontdoor_secret_id` field must be set if the `certificate_type` is `CustomerCertificate`")
	} else if certType == string(afddomains.AfdCertificateTypeManagedCertificate) {
		// Ignore computed `cdn_frontdoor_secret_id` for managed certs unless the
		// user explicitly configured it.
		if secretRaw != "" && secretWasConfigured {
			return nil, errors.New("the `cdn_frontdoor_secret_id` field is not supported if the `certificate_type` is `ManagedCertificate`")
		}
		if !secretWasConfigured {
			secretRaw = ""
		}
	}

	// NOTE: Secret always needs to be passed if it is defined else you will
	// receive a 500 Internal Server Error
	if secretRaw != "" {
		secret, err := secrets.ParseSecretID(secretRaw)
		if err != nil {
			return nil, err
		}

		tls.Secret = expandAfdResourceReference(secret.ID())
	}

	if minTlsVersion != "" {
		tls.MinimumTlsVersion = pointer.ToEnum[afddomains.AfdMinimumTlsVersion](minTlsVersion)
	}

	if len(cipherSuiteRaw) > 0 && cipherSuiteRaw[0] != nil {
		cipherSuite := cipherSuiteRaw[0].(map[string]any)

		if cipherSuiteType := cipherSuite["type"].(string); cipherSuiteType != "" {
			tls.CipherSuiteSetType = pointer.ToEnum[afddomains.AfdCipherSuiteSetType](cipherSuiteType)
		}

		if customCiphersRaw := cipherSuite["custom_ciphers"].([]any); len(customCiphersRaw) > 0 {
			tls.CustomizedCipherSuiteSet = expandAfdCustomizedCipherSuiteSet(customCiphersRaw)
		}
	}

	return &tls, nil
}

func validateFrontDoorCustomDomainTlsGroups(policy afddomains.AfdServerTlsGroupPolicy, count int) error {
	if policy == afddomains.AfdServerTlsGroupPolicyCustom && count == 0 {
		return errors.New("`tls.server_tls_groups` must be specified when `tls.server_tls_group_policy` is `Custom`")
	}
	if policy != afddomains.AfdServerTlsGroupPolicyCustom && count != 0 {
		return fmt.Errorf("`tls.server_tls_groups` cannot be specified when `tls.server_tls_group_policy` is `%s`", policy)
	}
	return nil
}

func expandAfdCustomizedCipherSuiteSet(input []any) *afddomains.AFDDomainHTTPSCustomizedCipherSuiteSet {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	v := input[0].(map[string]any)

	result := &afddomains.AFDDomainHTTPSCustomizedCipherSuiteSet{}

	if tls12Raw := v["tls12"].(*pluginsdk.Set); tls12Raw.Len() > 0 {
		tls12Suites := make([]afddomains.AfdCustomizedCipherSuiteForTls12, 0)
		for _, suite := range tls12Raw.List() {
			tls12Suites = append(tls12Suites, afddomains.AfdCustomizedCipherSuiteForTls12(suite.(string)))
		}
		result.CipherSuiteSetForTls12 = &tls12Suites
	}

	if tls13Raw := v["tls13"].(*pluginsdk.Set); tls13Raw.Len() > 0 {
		tls13Suites := make([]afddomains.AfdCustomizedCipherSuiteForTls13, 0)
		for _, suite := range tls13Raw.List() {
			tls13Suites = append(tls13Suites, afddomains.AfdCustomizedCipherSuiteForTls13(suite.(string)))
		}
		result.CipherSuiteSetForTls13 = &tls13Suites
	}

	return result
}

func expandAfdResourceReference(id string) *afddomains.ResourceReference {
	return &afddomains.ResourceReference{
		Id: pointer.To(id),
	}
}
