// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package helpers

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	webapps20250501 "github.com/hashicorp/go-azure-sdk/resource-manager/web/2025-05-01/webapps"
)

// This file contains duplicates of helper functions in shared_schema.go typed against the
// 2025-05-01 webapps SDK package, scoped exclusively to `azurerm_function_app_flex_consumption`.
// The shared `2023-12-01` webapps client and its helpers remain untouched to avoid breaking
// other appservice resources that depend on `SiteProperties` fields removed in `2025-05-01`.

func ExpandCorsSettingsV20250501(input []CorsSetting) *webapps20250501.CorsSettings {
	if len(input) != 1 {
		return &webapps20250501.CorsSettings{}
	}
	cors := input[0]

	return &webapps20250501.CorsSettings{
		AllowedOrigins:     pointer.To(cors.AllowedOrigins),
		SupportCredentials: pointer.To(cors.SupportCredentials),
	}
}

func FlattenCorsSettingsV20250501(input *webapps20250501.CorsSettings) []CorsSetting {
	if input == nil {
		return []CorsSetting{}
	}

	cors := *input
	if len(pointer.From(cors.AllowedOrigins)) == 0 && !pointer.From(cors.SupportCredentials) {
		return []CorsSetting{}
	}

	return []CorsSetting{{
		SupportCredentials: pointer.From(cors.SupportCredentials),
		AllowedOrigins:     pointer.From(cors.AllowedOrigins),
	}}
}

func ExpandIpRestrictionsV20250501(restrictions []IpRestriction) (*[]webapps20250501.IPSecurityRestriction, error) {
	expanded := make([]webapps20250501.IPSecurityRestriction, 0)
	if len(restrictions) == 0 {
		return &expanded, nil
	}

	for _, v := range restrictions {
		if err := v.Validate(); err != nil {
			return nil, err
		}

		var restriction webapps20250501.IPSecurityRestriction
		if v.Name != "" {
			restriction.Name = pointer.To(v.Name)
		}

		if v.IpAddress != "" {
			restriction.IPAddress = pointer.To(v.IpAddress)
		}

		if v.ServiceTag != "" {
			restriction.IPAddress = pointer.To(v.ServiceTag)
			restriction.Tag = pointer.To(webapps20250501.IPFilterTagServiceTag)
		}

		if v.VnetSubnetId != "" {
			restriction.VnetSubnetResourceId = pointer.To(v.VnetSubnetId)
		}

		if v.Description != "" {
			restriction.Description = pointer.To(v.Description)
		}

		restriction.Priority = pointer.To(v.Priority)

		restriction.Action = pointer.To(v.Action)

		restriction.Headers = expandIpRestrictionHeadersV20250501(v.Headers)

		expanded = append(expanded, restriction)
	}

	return &expanded, nil
}

func expandIpRestrictionHeadersV20250501(headers []IpRestrictionHeaders) *map[string][]string {
	result := make(map[string][]string)
	if len(headers) == 0 {
		return nil
	}

	for _, v := range headers {
		if len(v.XForwardedHost) > 0 {
			result["x-forwarded-host"] = v.XForwardedHost
		}
		if len(v.XForwardedFor) > 0 {
			result["x-forwarded-for"] = v.XForwardedFor
		}
		if len(v.XAzureFDID) > 0 {
			result["x-azure-fdid"] = v.XAzureFDID
		}
		if len(v.XFDHealthProbe) > 0 {
			result["x-fd-healthprobe"] = v.XFDHealthProbe
		}
	}

	return &result
}

func ExpandAuthSettingsV20250501(auth []AuthSettings) *webapps20250501.SiteAuthSettings {
	result := &webapps20250501.SiteAuthSettings{}
	if len(auth) == 0 {
		return result
	}

	props := &webapps20250501.SiteAuthSettingsProperties{}

	v := auth[0]

	props.Enabled = pointer.To(v.Enabled)

	additionalLoginParams := make([]string, 0)
	if len(v.AdditionalLoginParameters) > 0 {
		for k, s := range v.AdditionalLoginParameters {
			additionalLoginParams = append(additionalLoginParams, fmt.Sprintf("%s=%s", k, s))
		}
	}
	props.AdditionalLoginParams = &additionalLoginParams

	props.AllowedExternalRedirectURLs = &v.AllowedExternalRedirectURLs

	props.DefaultProvider = pointer.ToEnum[webapps20250501.BuiltInAuthenticationProvider](v.DefaultProvider)

	props.Issuer = pointer.To(v.Issuer)

	props.RuntimeVersion = pointer.To(v.RuntimeVersion)

	props.TokenStoreEnabled = pointer.To(v.TokenStoreEnabled)

	props.TokenRefreshExtensionHours = pointer.To(v.TokenRefreshExtensionHours)

	props.UnauthenticatedClientAction = pointer.ToEnum[webapps20250501.UnauthenticatedClientAction](v.UnauthenticatedClientAction)

	a := AadAuthSettings{}
	if len(v.AzureActiveDirectoryAuth) > 0 {
		a = v.AzureActiveDirectoryAuth[0]
	}
	props.ClientId = pointer.To(a.ClientId)

	if a.ClientSecret != "" {
		props.ClientSecret = pointer.To(a.ClientSecret)
	}

	if a.ClientSecretSettingName != "" {
		props.ClientSecretSettingName = pointer.To(a.ClientSecretSettingName)
	}

	props.AllowedAudiences = &a.AllowedAudiences

	f := FacebookAuthSettings{}
	if len(v.FacebookAuth) > 0 {
		f = v.FacebookAuth[0]
	}
	props.FacebookAppId = pointer.To(f.AppId)
	props.FacebookAppSecret = pointer.To(f.AppSecret)
	props.FacebookAppSecretSettingName = pointer.To(f.AppSecretSettingName)
	props.FacebookOAuthScopes = &f.OauthScopes

	gh := GithubAuthSettings{}
	if len(v.GithubAuth) > 0 {
		gh = v.GithubAuth[0]
	}
	props.GitHubClientId = pointer.To(gh.ClientId)
	props.GitHubClientSecret = pointer.To(gh.ClientSecret)
	props.GitHubClientSecretSettingName = pointer.To(gh.ClientSecretSettingName)
	props.GitHubOAuthScopes = &gh.OAuthScopes

	g := GoogleAuthSettings{}
	if len(v.GoogleAuth) > 0 {
		g = v.GoogleAuth[0]
	}

	props.GoogleClientId = pointer.To(g.ClientId)
	props.GoogleClientSecret = pointer.To(g.ClientSecret)
	props.GoogleClientSecretSettingName = pointer.To(g.ClientSecretSettingName)
	props.GoogleOAuthScopes = &g.OauthScopes

	m := MicrosoftAuthSettings{}
	if len(v.MicrosoftAuth) > 0 {
		m = v.MicrosoftAuth[0]
	}
	props.MicrosoftAccountClientId = pointer.To(m.ClientId)
	props.MicrosoftAccountClientSecret = pointer.To(m.ClientSecret)
	props.MicrosoftAccountClientSecretSettingName = pointer.To(m.ClientSecretSettingName)
	props.MicrosoftAccountOAuthScopes = &m.OauthScopes

	t := TwitterAuthSettings{}
	if len(v.TwitterAuth) > 0 {
		t = v.TwitterAuth[0]
	}
	props.TwitterConsumerKey = pointer.To(t.ConsumerKey)
	props.TwitterConsumerSecret = pointer.To(t.ConsumerSecret)
	props.TwitterConsumerSecretSettingName = pointer.To(t.ConsumerSecretSettingName)

	result.Properties = props

	return result
}

func FlattenAuthSettingsV20250501(auth *webapps20250501.SiteAuthSettings) []AuthSettings {
	if auth == nil || auth.Properties == nil || strings.ToLower(pointer.From(auth.Properties.ConfigVersion)) != "v1" {
		return []AuthSettings{}
	}

	props := *auth.Properties

	result := AuthSettings{
		DefaultProvider:             string(pointer.From(props.DefaultProvider)),
		UnauthenticatedClientAction: string(pointer.From(props.UnauthenticatedClientAction)),
	}

	if props.Enabled != nil {
		result.Enabled = *props.Enabled
	}

	if props.AdditionalLoginParams != nil {
		params := make(map[string]string)
		for _, v := range *props.AdditionalLoginParams {
			parts := strings.Split(v, "=")
			if len(parts) != 2 {
				continue
			}
			params[parts[0]] = parts[1]
		}
		result.AdditionalLoginParameters = params
	}

	result.AllowedExternalRedirectURLs = pointer.From(props.AllowedExternalRedirectURLs)

	if props.Issuer != nil {
		result.Issuer = *props.Issuer
	}

	if props.RuntimeVersion != nil {
		result.RuntimeVersion = *props.RuntimeVersion
	}

	if props.TokenRefreshExtensionHours != nil {
		result.TokenRefreshExtensionHours = *props.TokenRefreshExtensionHours
	}

	if props.TokenStoreEnabled != nil {
		result.TokenStoreEnabled = *props.TokenStoreEnabled
	}

	// AAD Auth
	if props.ClientId != nil {
		aadAuthSettings := AadAuthSettings{
			ClientId: *props.ClientId,
		}

		if props.ClientSecret != nil {
			aadAuthSettings.ClientSecret = *props.ClientSecret
		}

		if props.ClientSecretSettingName != nil {
			aadAuthSettings.ClientSecretSettingName = *props.ClientSecretSettingName
		}

		if props.AllowedAudiences != nil {
			aadAuthSettings.AllowedAudiences = *props.AllowedAudiences
		}

		result.AzureActiveDirectoryAuth = []AadAuthSettings{aadAuthSettings}
	}

	if props.FacebookAppId != nil {
		facebookAuthSettings := FacebookAuthSettings{
			AppId: *props.FacebookAppId,
		}

		if props.FacebookAppSecret != nil {
			facebookAuthSettings.AppSecret = *props.FacebookAppSecret
		}

		if props.FacebookAppSecretSettingName != nil {
			facebookAuthSettings.AppSecretSettingName = *props.FacebookAppSecretSettingName
		}

		if props.FacebookOAuthScopes != nil {
			facebookAuthSettings.OauthScopes = *props.FacebookOAuthScopes
		}

		result.FacebookAuth = []FacebookAuthSettings{facebookAuthSettings}
	}

	if props.GitHubClientId != nil {
		githubAuthSetting := GithubAuthSettings{
			ClientId: *props.GitHubClientId,
		}

		if props.GitHubClientSecret != nil {
			githubAuthSetting.ClientSecret = *props.GitHubClientSecret
		}

		if props.GitHubClientSecretSettingName != nil {
			githubAuthSetting.ClientSecretSettingName = *props.GitHubClientSecretSettingName
		}

		result.GithubAuth = []GithubAuthSettings{githubAuthSetting}
	}

	if props.GoogleClientId != nil {
		googleAuthSettings := GoogleAuthSettings{
			ClientId: *props.GoogleClientId,
		}

		if props.GoogleClientSecret != nil {
			googleAuthSettings.ClientSecret = *props.GoogleClientSecret
		}

		if props.GoogleClientSecretSettingName != nil {
			googleAuthSettings.ClientSecretSettingName = *props.GoogleClientSecretSettingName
		}

		if props.GoogleOAuthScopes != nil {
			googleAuthSettings.OauthScopes = *props.GoogleOAuthScopes
		}

		result.GoogleAuth = []GoogleAuthSettings{googleAuthSettings}
	}

	if props.MicrosoftAccountClientId != nil {
		microsoftAuthSettings := MicrosoftAuthSettings{
			ClientId: *props.MicrosoftAccountClientId,
		}

		if props.MicrosoftAccountClientSecret != nil {
			microsoftAuthSettings.ClientSecret = *props.MicrosoftAccountClientSecret
		}

		if props.MicrosoftAccountClientSecretSettingName != nil {
			microsoftAuthSettings.ClientSecretSettingName = *props.MicrosoftAccountClientSecretSettingName
		}

		if props.MicrosoftAccountOAuthScopes != nil {
			microsoftAuthSettings.OauthScopes = *props.MicrosoftAccountOAuthScopes
		}

		result.MicrosoftAuth = []MicrosoftAuthSettings{microsoftAuthSettings}
	}

	if props.TwitterConsumerKey != nil {
		twitterAuthSetting := TwitterAuthSettings{
			ConsumerKey: *props.TwitterConsumerKey,
		}
		if props.TwitterConsumerSecret != nil {
			twitterAuthSetting.ConsumerSecret = *props.TwitterConsumerSecret
		}
		if props.TwitterConsumerSecretSettingName != nil {
			twitterAuthSetting.ConsumerSecretSettingName = *props.TwitterConsumerSecretSettingName
		}

		result.TwitterAuth = []TwitterAuthSettings{twitterAuthSetting}
	}

	return []AuthSettings{result}
}

func FlattenIpRestrictionsV20250501(ipRestrictionsList *[]webapps20250501.IPSecurityRestriction) []IpRestriction {
	if ipRestrictionsList == nil {
		return []IpRestriction{}
	}

	ipRestrictions := make([]IpRestriction, 0, len(*ipRestrictionsList))
	for _, v := range *ipRestrictionsList {
		ipRestriction := IpRestriction{}

		if v.Name != nil {
			ipRestriction.Name = *v.Name
		}

		if v.IPAddress != nil {
			if *v.IPAddress == "Any" {
				continue
			}

			if v.Tag != nil && *v.Tag == webapps20250501.IPFilterTagServiceTag {
				ipRestriction.ServiceTag = *v.IPAddress
			} else {
				ipRestriction.IpAddress = *v.IPAddress
			}
		}

		if v.VnetSubnetResourceId != nil {
			ipRestriction.VnetSubnetId = *v.VnetSubnetResourceId
		}

		if v.Priority != nil {
			ipRestriction.Priority = *v.Priority
		}

		if v.Action != nil {
			ipRestriction.Action = *v.Action
		}

		ipRestriction.Headers = flattenIpRestrictionHeadersV20250501(pointer.From(v.Headers))
		if v.Description != nil {
			ipRestriction.Description = *v.Description
		}

		ipRestrictions = append(ipRestrictions, ipRestriction)
	}

	return ipRestrictions
}

func flattenIpRestrictionHeadersV20250501(headers map[string][]string) []IpRestrictionHeaders {
	if len(headers) == 0 {
		return []IpRestrictionHeaders{}
	}
	ipRestrictionHeader := IpRestrictionHeaders{}
	if xForwardFor, ok := headers["x-forwarded-for"]; ok {
		ipRestrictionHeader.XForwardedFor = xForwardFor
	}

	if xForwardedHost, ok := headers["x-forwarded-host"]; ok {
		ipRestrictionHeader.XForwardedHost = xForwardedHost
	}

	if xAzureFDID, ok := headers["x-azure-fdid"]; ok {
		ipRestrictionHeader.XAzureFDID = xAzureFDID
	}

	if xFDHealthProbe, ok := headers["x-fd-healthprobe"]; ok {
		ipRestrictionHeader.XFDHealthProbe = xFDHealthProbe
	}

	return []IpRestrictionHeaders{ipRestrictionHeader}
}

func FlattenSiteCredentialsV20250501(input *webapps20250501.User) []SiteCredential {
	var result []SiteCredential
	if input == nil || input.Properties == nil {
		return result
	}

	userProps := *input.Properties
	result = append(result, SiteCredential{
		Username: userProps.PublishingUserName,
		Password: pointer.From(userProps.PublishingPassword),
	})

	return result
}

func ExpandStickySettingsV20250501(input []StickySettings) *webapps20250501.SlotConfigNames {
	if len(input) == 0 {
		return nil
	}

	return &webapps20250501.SlotConfigNames{
		AppSettingNames:       &input[0].AppSettingNames,
		ConnectionStringNames: &input[0].ConnectionStringNames,
	}
}

func FlattenStickySettingsV20250501(input *webapps20250501.SlotConfigNames) []StickySettings {
	result := StickySettings{}
	if input == nil || (input.AppSettingNames == nil && input.ConnectionStringNames == nil) || (len(*input.AppSettingNames) == 0 && len(*input.ConnectionStringNames) == 0) {
		return []StickySettings{}
	}

	if input.AppSettingNames != nil && len(*input.AppSettingNames) > 0 {
		result.AppSettingNames = *input.AppSettingNames
	}

	if input.ConnectionStringNames != nil && len(*input.ConnectionStringNames) > 0 {
		result.ConnectionStringNames = *input.ConnectionStringNames
	}

	return []StickySettings{result}
}

func DefaultAuthSettingsPropertiesV20250501() *webapps20250501.SiteAuthSettingsProperties {
	return &webapps20250501.SiteAuthSettingsProperties{
		Enabled:                                 pointer.To(false),
		AdditionalLoginParams:                   pointer.To(make([]string, 0)),
		AllowedAudiences:                        pointer.To(make([]string, 0)),
		ClientId:                                pointer.To(""),
		ClientSecret:                            pointer.To(""),
		ClientSecretSettingName:                 pointer.To(""),
		ClientSecretCertificateThumbprint:       pointer.To(""),
		FacebookAppId:                           pointer.To(""),
		FacebookAppSecret:                       pointer.To(""),
		FacebookAppSecretSettingName:            pointer.To(""),
		FacebookOAuthScopes:                     pointer.To(make([]string, 0)),
		GitHubClientId:                          pointer.To(""),
		GitHubOAuthScopes:                       pointer.To(make([]string, 0)),
		GitHubClientSecret:                      pointer.To(""),
		GitHubClientSecretSettingName:           pointer.To(""),
		GoogleClientId:                          pointer.To(""),
		GoogleOAuthScopes:                       pointer.To(make([]string, 0)),
		GoogleClientSecret:                      pointer.To(""),
		GoogleClientSecretSettingName:           pointer.To(""),
		Issuer:                                  pointer.To(""),
		MicrosoftAccountClientId:                pointer.To(""),
		MicrosoftAccountOAuthScopes:             pointer.To(make([]string, 0)),
		MicrosoftAccountClientSecret:            pointer.To(""),
		MicrosoftAccountClientSecretSettingName: pointer.To(""),
		TokenRefreshExtensionHours:              pointer.To(72.0),
		TokenStoreEnabled:                       pointer.To(false),
		TwitterConsumerKey:                      pointer.To(""),
		TwitterConsumerSecret:                   pointer.To(""),
		TwitterConsumerSecretSettingName:        pointer.To(""),
	}
}
