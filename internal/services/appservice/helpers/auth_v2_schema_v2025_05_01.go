// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package helpers

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	webapps20250501 "github.com/hashicorp/go-azure-sdk/resource-manager/web/2025-05-01/webapps"
)

// This file contains duplicates of helper functions in auth_v2_schema.go typed against the
// 2025-05-01 webapps SDK package, scoped exclusively to `azurerm_function_app_flex_consumption`.

func expandAuthV2LoginSettingsV20250501(input []AuthV2Login) *webapps20250501.Login {
	if len(input) == 0 {
		return nil
	}
	login := input[0]
	result := &webapps20250501.Login{
		Routes: &webapps20250501.LoginRoutes{},
		TokenStore: &webapps20250501.TokenStore{
			Enabled:          pointer.To(login.TokenStoreEnabled),
			FileSystem:       &webapps20250501.FileSystemTokenStore{},
			AzureBlobStorage: &webapps20250501.BlobStorageTokenStore{},
		},
		PreserveURLFragmentsForLogins: pointer.To(login.PreserveURLFragmentsForLogins),
		Nonce: &webapps20250501.Nonce{
			ValidateNonce:           pointer.To(login.ValidateNonce),
			NonceExpirationInterval: pointer.To(login.NonceExpirationTime),
		},
		CookieExpiration: &webapps20250501.CookieExpiration{
			Convention:       pointer.ToEnum[webapps20250501.CookieExpirationConvention](login.CookieExpirationConvention),
			TimeToExpiration: pointer.To(login.CookieExpirationTime),
		},
	}

	if login.TokenFilesystemPath != "" || login.TokenBlobStorageSAS != "" {
		result.TokenStore.Enabled = pointer.To(true)
		if login.TokenFilesystemPath != "" {
			result.TokenStore.FileSystem = &webapps20250501.FileSystemTokenStore{
				Directory: pointer.To(login.TokenFilesystemPath),
			}
		}
		if login.TokenBlobStorageSAS != "" {
			result.TokenStore.AzureBlobStorage = &webapps20250501.BlobStorageTokenStore{
				SasURLSettingName: pointer.To(login.TokenBlobStorageSAS),
			}
		}
	}

	if login.LogoutEndpoint != "" {
		result.Routes = &webapps20250501.LoginRoutes{
			LogoutEndpoint: pointer.To(login.LogoutEndpoint),
		}
	}
	result.TokenStore.TokenRefreshExtensionHours = pointer.To(login.TokenRefreshExtension)
	if login.TokenFilesystemPath != "" {
		result.TokenStore.FileSystem = &webapps20250501.FileSystemTokenStore{
			Directory: pointer.To(login.TokenFilesystemPath),
		}
	}
	if login.TokenBlobStorageSAS != "" {
		result.TokenStore.AzureBlobStorage = &webapps20250501.BlobStorageTokenStore{
			SasURLSettingName: pointer.To(login.TokenBlobStorageSAS),
		}
	}
	result.AllowedExternalRedirectURLs = pointer.To(login.AllowedExternalRedirectURLs)

	return result
}

func flattenAuthV2LoginSettingsV20250501(input *webapps20250501.Login) []AuthV2Login {
	if input == nil {
		return []AuthV2Login{}
	}
	result := AuthV2Login{
		PreserveURLFragmentsForLogins: pointer.From(input.PreserveURLFragmentsForLogins),
		AllowedExternalRedirectURLs:   pointer.From(input.AllowedExternalRedirectURLs),
	}
	if routes := input.Routes; routes != nil {
		result.LogoutEndpoint = pointer.From(routes.LogoutEndpoint)
	}
	if token := input.TokenStore; token != nil {
		result.TokenStoreEnabled = pointer.From(token.Enabled)
		result.TokenRefreshExtension = pointer.From(token.TokenRefreshExtensionHours)
		if fs := token.FileSystem; fs != nil {
			result.TokenFilesystemPath = pointer.From(fs.Directory)
		}
		if bs := token.AzureBlobStorage; bs != nil {
			result.TokenBlobStorageSAS = pointer.From(bs.SasURLSettingName)
		}
	}

	if nonce := input.Nonce; nonce != nil {
		result.NonceExpirationTime = pointer.From(nonce.NonceExpirationInterval)
		result.ValidateNonce = pointer.From(nonce.ValidateNonce)
	}

	if cookie := input.CookieExpiration; cookie != nil {
		result.CookieExpirationConvention = string(pointer.From(cookie.Convention))
		result.CookieExpirationTime = pointer.From(cookie.TimeToExpiration)
	}

	return []AuthV2Login{result}
}

func expandAppleAuthV2SettingsV20250501(input []AppleAuthV2Settings) *webapps20250501.Apple {
	if len(input) == 1 {
		apple := input[0]
		return &webapps20250501.Apple{
			Enabled: pointer.To(true),
			Registration: &webapps20250501.AppleRegistration{
				ClientId:                pointer.To(apple.ClientId),
				ClientSecretSettingName: pointer.To(apple.ClientSecretSettingName),
			},
			Login: &webapps20250501.LoginScopes{
				Scopes: pointer.To(apple.LoginScopes),
			},
		}
	}

	return &webapps20250501.Apple{
		Enabled: pointer.To(false),
	}
}

func flattenAppleAuthV2SettingsV20250501(input *webapps20250501.Apple) []AppleAuthV2Settings {
	if input == nil || !pointer.From(input.Enabled) {
		return []AppleAuthV2Settings{}
	}
	result := AppleAuthV2Settings{}

	props := *input
	if reg := props.Registration; reg != nil {
		result.ClientId = pointer.From(reg.ClientId)
		result.ClientSecretSettingName = pointer.From(reg.ClientSecretSettingName)
	}
	if loginScopes := props.Login; loginScopes != nil {
		result.LoginScopes = pointer.From(loginScopes.Scopes)
	}

	return []AppleAuthV2Settings{result}
}

func expandAadAuthV2SettingsV20250501(input []AadAuthV2Settings) *webapps20250501.AzureActiveDirectory {
	result := &webapps20250501.AzureActiveDirectory{
		Enabled: pointer.To(false),
	}

	if len(input) == 1 {
		aad := input[0]
		result = &webapps20250501.AzureActiveDirectory{
			Enabled: pointer.To(true),
			Registration: &webapps20250501.AzureActiveDirectoryRegistration{
				OpenIdIssuer: pointer.To(aad.TenantAuthURI),
				ClientId:     pointer.To(aad.ClientId),
			},
			Login: &webapps20250501.AzureActiveDirectoryLogin{
				DisableWWWAuthenticate: pointer.To(aad.DisableWWWAuth),
			},
		}

		if aad.ClientSecretSettingName != "" {
			result.Registration.ClientSecretSettingName = pointer.To(aad.ClientSecretSettingName)
		}

		if aad.ClientSecretCertificateThumbprint != "" {
			result.Registration.ClientSecretCertificateThumbprint = pointer.To(aad.ClientSecretCertificateThumbprint)
		}

		if len(aad.LoginParameters) > 0 {
			params := make([]string, 0)
			for k, v := range aad.LoginParameters {
				params = append(params, fmt.Sprintf("%s=%s", k, v))
			}
			result.Login.LoginParameters = &params
		}

		if len(aad.JWTAllowedGroups) != 0 || len(aad.JWTAllowedClientApps) != 0 {
			if result.Validation == nil {
				result.Validation = &webapps20250501.AzureActiveDirectoryValidation{}
			}
			result.Validation.JwtClaimChecks = &webapps20250501.JwtClaimChecks{}
			if len(aad.JWTAllowedGroups) != 0 {
				result.Validation.JwtClaimChecks.AllowedGroups = pointer.To(aad.JWTAllowedGroups)
			}
			if len(aad.JWTAllowedClientApps) != 0 {
				result.Validation.JwtClaimChecks.AllowedClientApplications = pointer.To(aad.JWTAllowedClientApps)
			}
		}

		if len(aad.AllowedGroups) > 0 || len(aad.AllowedIdentities) > 0 {
			if result.Validation == nil {
				result.Validation = &webapps20250501.AzureActiveDirectoryValidation{}
			}
			result.Validation.DefaultAuthorizationPolicy = &webapps20250501.DefaultAuthorizationPolicy{
				AllowedPrincipals: &webapps20250501.AllowedPrincipals{},
			}
			if len(aad.AllowedGroups) > 0 {
				result.Validation.DefaultAuthorizationPolicy.AllowedPrincipals.Groups = pointer.To(aad.AllowedGroups)
			}
			if len(aad.AllowedIdentities) > 0 {
				result.Validation.DefaultAuthorizationPolicy.AllowedPrincipals.Identities = pointer.To(aad.AllowedIdentities)
			}
		}
		if len(aad.AllowedAudiences) > 0 {
			if result.Validation == nil {
				result.Validation = &webapps20250501.AzureActiveDirectoryValidation{}
			}
			result.Validation.AllowedAudiences = pointer.To(aad.AllowedAudiences)
		}

		if len(aad.AllowedApplications) > 0 {
			if result.Validation == nil {
				result.Validation = &webapps20250501.AzureActiveDirectoryValidation{}
			}
			if result.Validation.DefaultAuthorizationPolicy == nil {
				result.Validation.DefaultAuthorizationPolicy = &webapps20250501.DefaultAuthorizationPolicy{}
			}
			result.Validation.DefaultAuthorizationPolicy.AllowedApplications = pointer.To(aad.AllowedApplications)
		}
	}

	return result
}

func flattenAadAuthV2SettingsV20250501(input *webapps20250501.AzureActiveDirectory) []AadAuthV2Settings {
	if input == nil || !pointer.From(input.Enabled) {
		return []AadAuthV2Settings{}
	}

	result := AadAuthV2Settings{}

	if reg := input.Registration; reg != nil {
		result.TenantAuthURI = pointer.From(reg.OpenIdIssuer)
		result.ClientId = pointer.From(reg.ClientId)
		result.ClientSecretSettingName = pointer.From(reg.ClientSecretSettingName)
		result.ClientSecretCertificateThumbprint = pointer.From(reg.ClientSecretCertificateThumbprint)
	}

	if login := input.Login; login != nil {
		result.DisableWWWAuth = pointer.From(login.DisableWWWAuthenticate)
		if loginParamsRaw := login.LoginParameters; loginParamsRaw != nil {
			loginParams := make(map[string]string)
			for _, v := range *loginParamsRaw {
				parts := strings.Split(v, "=")
				if len(parts) == 2 && parts[0] != "" {
					loginParams[parts[0]] = parts[1]
				}
			}
			result.LoginParameters = loginParams
		}
	}

	if validation := input.Validation; validation != nil {
		if validation.AllowedAudiences != nil {
			result.AllowedAudiences = *validation.AllowedAudiences
		}
		if jwt := validation.JwtClaimChecks; jwt != nil {
			result.JWTAllowedGroups = pointer.From(jwt.AllowedGroups)
			result.JWTAllowedClientApps = pointer.From(jwt.AllowedClientApplications)
		}
		if defaultPolicy := validation.DefaultAuthorizationPolicy; defaultPolicy != nil {
			result.AllowedApplications = pointer.From(defaultPolicy.AllowedApplications)
			if defaultPolicy.AllowedPrincipals != nil {
				result.AllowedGroups = pointer.From(defaultPolicy.AllowedPrincipals.Groups)
				result.AllowedIdentities = pointer.From(defaultPolicy.AllowedPrincipals.Identities)
			}
		}
	}

	return []AadAuthV2Settings{result}
}

func expandStaticWebAppAuthV2SettingsV20250501(input []StaticWebAppAuthV2Settings) *webapps20250501.AzureStaticWebApps {
	if len(input) == 1 {
		swa := input[0]
		return &webapps20250501.AzureStaticWebApps{
			Enabled: pointer.To(true),
			Registration: &webapps20250501.AzureStaticWebAppsRegistration{
				ClientId: pointer.To(swa.ClientId),
			},
		}
	}

	return &webapps20250501.AzureStaticWebApps{
		Enabled: pointer.To(false),
	}
}

func flattenStaticWebAppAuthV2SettingsV20250501(input *webapps20250501.AzureStaticWebApps) []StaticWebAppAuthV2Settings {
	if input == nil || (input.Enabled != nil && !*input.Enabled) {
		return []StaticWebAppAuthV2Settings{}
	}

	result := StaticWebAppAuthV2Settings{}

	if pointer.From(input.Enabled) {
		if input.Registration != nil {
			result.ClientId = pointer.From(input.Registration.ClientId)
		}
	}

	return []StaticWebAppAuthV2Settings{result}
}

func expandCustomOIDCAuthV2SettingsV20250501(input []CustomOIDCAuthV2Settings) map[string]webapps20250501.CustomOpenIdConnectProvider {
	if len(input) == 0 {
		return nil
	}
	result := make(map[string]webapps20250501.CustomOpenIdConnectProvider)
	for _, v := range input {
		if v.Name == "" {
			continue
		}
		provider := webapps20250501.CustomOpenIdConnectProvider{
			Enabled: pointer.To(true),
			Registration: &webapps20250501.OpenIdConnectRegistration{
				ClientId: pointer.To(v.ClientId),
				ClientCredential: &webapps20250501.OpenIdConnectClientCredential{
					Method:                  pointer.To(webapps20250501.MethodClientSecretPost),
					ClientSecretSettingName: pointer.To(fmt.Sprintf("%s_PROVIDER_AUTHENTICATION_SECRET", strings.ToUpper(v.Name))),
				},
				OpenIdConnectConfiguration: &webapps20250501.OpenIdConnectConfig{
					WellKnownOpenIdConfiguration: pointer.To(v.OpenIDConfigurationEndpoint),
				},
			},
			Login: &webapps20250501.OpenIdConnectLogin{
				Scopes: pointer.To(v.Scopes),
			},
		}

		if v.NameClaimType != "" {
			provider.Login.NameClaimType = pointer.To(v.NameClaimType)
		}

		result[v.Name] = provider
	}

	return result
}

func flattenCustomOIDCAuthV2SettingsV20250501(input *map[string]webapps20250501.CustomOpenIdConnectProvider) []CustomOIDCAuthV2Settings {
	if input == nil || len(*input) == 0 {
		return []CustomOIDCAuthV2Settings{}
	}

	result := make([]CustomOIDCAuthV2Settings, 0)
	for k, v := range *input {
		if !pointer.From(v.Enabled) {
			continue
		} else {
			provider := CustomOIDCAuthV2Settings{
				Name: k,
			}
			if reg := v.Registration; reg != nil {
				provider.ClientId = pointer.From(reg.ClientId)
				if reg.ClientCredential != nil {
					provider.ClientSecretSettingName = pointer.From(reg.ClientCredential.ClientSecretSettingName)
					provider.ClientCredentialMethod = string(pointer.From(reg.ClientCredential.Method))
				}
				if config := reg.OpenIdConnectConfiguration; config != nil {
					provider.OpenIDConfigurationEndpoint = pointer.From(config.WellKnownOpenIdConfiguration)
					provider.AuthorizationEndpoint = pointer.From(config.AuthorizationEndpoint)
					provider.TokenEndpoint = pointer.From(config.TokenEndpoint)
					provider.IssuerEndpoint = pointer.From(config.Issuer)
					provider.CertificationURI = pointer.From(config.CertificationUri)
				}
			}
			if login := v.Login; login != nil {
				if login.Scopes != nil {
					provider.Scopes = *login.Scopes
				}
				provider.NameClaimType = pointer.From(login.NameClaimType)
			}
			result = append(result, provider)
		}
	}

	return result
}

func expandFacebookAuthV2SettingsV20250501(input []FacebookAuthV2Settings) *webapps20250501.Facebook {
	if len(input) == 1 {
		facebook := input[0]
		result := &webapps20250501.Facebook{
			Enabled: pointer.To(true),
			Registration: &webapps20250501.AppRegistration{
				AppId:                pointer.To(facebook.AppId),
				AppSecretSettingName: pointer.To(facebook.AppSecretSettingName),
			},
		}

		result.GraphApiVersion = pointer.To(facebook.GraphAPIVersion)
		result.Login = &webapps20250501.LoginScopes{
			Scopes: pointer.To(facebook.LoginScopes),
		}

		return result
	}

	return &webapps20250501.Facebook{
		Enabled: pointer.To(false),
	}
}

func flattenFacebookAuthV2SettingsV20250501(input *webapps20250501.Facebook) []FacebookAuthV2Settings {
	if input == nil || !pointer.From(input.Enabled) {
		return []FacebookAuthV2Settings{}
	}

	result := FacebookAuthV2Settings{
		GraphAPIVersion: pointer.From(input.GraphApiVersion),
	}

	if reg := input.Registration; reg != nil {
		result.AppId = pointer.From(reg.AppId)
		result.AppSecretSettingName = pointer.From(reg.AppSecretSettingName)
	}
	if login := input.Login; login != nil {
		result.LoginScopes = pointer.From(login.Scopes)
	}

	return []FacebookAuthV2Settings{result}
}

func expandGitHubAuthV2SettingsV20250501(input []GithubAuthV2Settings) *webapps20250501.GitHub {
	if len(input) == 1 {
		github := input[0]
		return &webapps20250501.GitHub{
			Enabled: pointer.To(true),
			Registration: &webapps20250501.ClientRegistration{
				ClientId:                pointer.To(github.ClientId),
				ClientSecretSettingName: pointer.To(github.ClientSecretSettingName),
			},
			Login: &webapps20250501.LoginScopes{
				Scopes: pointer.To(github.LoginScopes),
			},
		}
	}

	return &webapps20250501.GitHub{
		Enabled: pointer.To(false),
	}
}

func flattenGitHubAuthV2SettingsV20250501(input *webapps20250501.GitHub) []GithubAuthV2Settings {
	if input == nil || !pointer.From(input.Enabled) {
		return []GithubAuthV2Settings{}
	}

	result := GithubAuthV2Settings{}

	if reg := input.Registration; reg != nil {
		result.ClientId = pointer.From(reg.ClientId)
		result.ClientSecretSettingName = pointer.From(reg.ClientSecretSettingName)
	}
	if login := input.Login; login != nil && login.Scopes != nil {
		result.LoginScopes = pointer.From(login.Scopes)
	}

	return []GithubAuthV2Settings{result}
}

func expandGoogleAuthV2SettingsV20250501(input []GoogleAuthV2Settings) *webapps20250501.Google {
	if len(input) == 1 {
		google := input[0]
		return &webapps20250501.Google{
			Enabled: pointer.To(true),
			Registration: &webapps20250501.ClientRegistration{
				ClientId:                pointer.To(google.ClientId),
				ClientSecretSettingName: pointer.To(google.ClientSecretSettingName),
			},
			Validation: &webapps20250501.AllowedAudiencesValidation{
				AllowedAudiences: pointer.To(google.AllowedAudiences),
			},
			Login: &webapps20250501.LoginScopes{
				Scopes: pointer.To(google.LoginScopes),
			},
		}
	}

	return &webapps20250501.Google{
		Enabled: pointer.To(false),
	}
}

func flattenGoogleAuthV2SettingsV20250501(input *webapps20250501.Google) []GoogleAuthV2Settings {
	if input == nil || !pointer.From(input.Enabled) {
		return []GoogleAuthV2Settings{}
	}

	result := GoogleAuthV2Settings{}

	if reg := input.Registration; reg != nil {
		result.ClientId = pointer.From(reg.ClientId)
		result.ClientSecretSettingName = pointer.From(reg.ClientSecretSettingName)
	}
	if login := input.Login; login != nil && login.Scopes != nil {
		result.LoginScopes = *login.Scopes
	}
	if val := input.Validation; val != nil && val.AllowedAudiences != nil {
		result.LoginScopes = *val.AllowedAudiences
	}

	return []GoogleAuthV2Settings{result}
}

func expandMicrosoftAuthV2SettingsV20250501(input []MicrosoftAuthV2Settings) *webapps20250501.LegacyMicrosoftAccount {
	if len(input) == 1 {
		msft := input[0]
		return &webapps20250501.LegacyMicrosoftAccount{
			Enabled: pointer.To(true),
			Registration: &webapps20250501.ClientRegistration{
				ClientId:                pointer.To(msft.ClientId),
				ClientSecretSettingName: pointer.To(msft.ClientSecretSettingName),
			},
			Validation: &webapps20250501.AllowedAudiencesValidation{
				AllowedAudiences: pointer.To(msft.AllowedAudiences),
			},
			Login: &webapps20250501.LoginScopes{
				Scopes: pointer.To(msft.LoginScopes),
			},
		}
	}

	return &webapps20250501.LegacyMicrosoftAccount{
		Enabled: pointer.To(false),
	}
}

func flattenMicrosoftAuthV2SettingsV20250501(input *webapps20250501.LegacyMicrosoftAccount) []MicrosoftAuthV2Settings {
	if input == nil || !pointer.From(input.Enabled) {
		return []MicrosoftAuthV2Settings{}
	}

	result := MicrosoftAuthV2Settings{}

	if reg := input.Registration; reg != nil {
		result.ClientId = pointer.From(reg.ClientId)
		result.ClientSecretSettingName = pointer.From(reg.ClientSecretSettingName)
	}
	if login := input.Login; login != nil && login.Scopes != nil {
		result.LoginScopes = *login.Scopes
	}
	if val := input.Validation; val != nil && val.AllowedAudiences != nil {
		result.LoginScopes = *val.AllowedAudiences
	}

	return []MicrosoftAuthV2Settings{result}
}

func expandTwitterAuthV2SettingsV20250501(input []TwitterAuthV2Settings) *webapps20250501.Twitter {
	if len(input) == 1 {
		twitter := input[0]
		return &webapps20250501.Twitter{
			Enabled: pointer.To(true),
			Registration: &webapps20250501.TwitterRegistration{
				ConsumerKey:               pointer.To(twitter.ConsumerKey),
				ConsumerSecretSettingName: pointer.To(twitter.ConsumerSecretSettingName),
			},
		}
	}

	return &webapps20250501.Twitter{
		Enabled: pointer.To(false),
	}
}

func flattenTwitterAuthV2SettingsV20250501(input *webapps20250501.Twitter) []TwitterAuthV2Settings {
	if input == nil || !pointer.From(input.Enabled) {
		return []TwitterAuthV2Settings{}
	}

	if pointer.From(input.Enabled) {
		result := TwitterAuthV2Settings{}
		if reg := input.Registration; reg != nil {
			result.ConsumerKey = pointer.From(reg.ConsumerKey)
			result.ConsumerSecretSettingName = pointer.From(reg.ConsumerSecretSettingName)
		}
		return []TwitterAuthV2Settings{result}
	}

	return nil
}

func ExpandAuthV2SettingsV20250501(input []AuthV2Settings) *webapps20250501.SiteAuthSettingsV2 {
	result := &webapps20250501.SiteAuthSettingsV2{}
	if len(input) != 1 {
		return result
	}

	settings := input[0]

	props := &webapps20250501.SiteAuthSettingsV2Properties{
		Platform: &webapps20250501.AuthPlatform{
			Enabled:        pointer.To(settings.AuthEnabled),
			RuntimeVersion: pointer.To(settings.RuntimeVersion),
		},
		GlobalValidation: &webapps20250501.GlobalValidation{
			RequireAuthentication:       pointer.To(settings.RequireAuth),
			UnauthenticatedClientAction: pointer.ToEnum[webapps20250501.UnauthenticatedClientActionV2](settings.UnauthenticatedAction),
			ExcludedPaths:               pointer.To(settings.ExcludedPaths),
		},
		IdentityProviders: &webapps20250501.IdentityProviders{
			AzureActiveDirectory:         expandAadAuthV2SettingsV20250501(settings.AzureActiveDirectoryAuth),
			Facebook:                     expandFacebookAuthV2SettingsV20250501(settings.FacebookAuth),
			GitHub:                       expandGitHubAuthV2SettingsV20250501(settings.GithubAuth),
			Google:                       expandGoogleAuthV2SettingsV20250501(settings.GoogleAuth),
			Twitter:                      expandTwitterAuthV2SettingsV20250501(settings.TwitterAuth),
			CustomOpenIdConnectProviders: pointer.To(expandCustomOIDCAuthV2SettingsV20250501(settings.CustomOIDCAuth)),
			LegacyMicrosoftAccount:       expandMicrosoftAuthV2SettingsV20250501(settings.MicrosoftAuth),
			Apple:                        expandAppleAuthV2SettingsV20250501(settings.AppleAuth),
			AzureStaticWebApps:           expandStaticWebAppAuthV2SettingsV20250501(settings.AzureStaticWebAuth),
		},
		Login: expandAuthV2LoginSettingsV20250501(settings.Login),
		HTTPSettings: &webapps20250501.HTTPSettings{
			RequireHTTPS: pointer.To(settings.RequireHTTPS),
			Routes: &webapps20250501.HTTPSettingsRoutes{
				ApiPrefix: pointer.To(settings.HttpRoutesAPIPrefix),
			},
			ForwardProxy: &webapps20250501.ForwardProxy{
				Convention: pointer.ToEnum[webapps20250501.ForwardProxyConvention](settings.ForwardProxyConvention),
			},
		},
	}

	// Platform
	if settings.ConfigFilePath != "" {
		props.Platform.ConfigFilePath = pointer.To(settings.ConfigFilePath)
	}

	// Global
	if settings.DefaultAuthProvider != "" {
		props.GlobalValidation.RedirectToProvider = pointer.To(settings.DefaultAuthProvider)
	}

	// HTTP
	if settings.ForwardProxyCustomHostHeaderName != "" {
		props.HTTPSettings.ForwardProxy.CustomHostHeaderName = pointer.To(settings.ForwardProxyCustomHostHeaderName)
	}
	if settings.ForwardProxyCustomSchemeHeaderName != "" {
		props.HTTPSettings.ForwardProxy.CustomProtoHeaderName = pointer.To(settings.ForwardProxyCustomSchemeHeaderName)
	}

	result.Properties = props

	return result
}

func FlattenAuthV2SettingsV20250501(input webapps20250501.SiteAuthSettingsV2) []AuthV2Settings {
	if input.Properties == nil {
		return []AuthV2Settings{}
	}

	settings := *input.Properties

	result := AuthV2Settings{}

	if platform := settings.Platform; platform != nil {
		result.AuthEnabled = pointer.From(platform.Enabled)
		result.RuntimeVersion = pointer.From(platform.RuntimeVersion)
		result.ConfigFilePath = pointer.From(platform.ConfigFilePath)
	}

	if global := settings.GlobalValidation; global != nil {
		result.RequireAuth = pointer.From(global.RequireAuthentication)
		result.UnauthenticatedAction = string(pointer.From(global.UnauthenticatedClientAction))
		result.DefaultAuthProvider = pointer.From(global.RedirectToProvider)
		result.ExcludedPaths = pointer.From(global.ExcludedPaths)
	}

	if http := settings.HTTPSettings; http != nil {
		result.RequireHTTPS = pointer.From(http.RequireHTTPS)
		if http.Routes != nil {
			result.HttpRoutesAPIPrefix = pointer.From(http.Routes.ApiPrefix)
		}
		if fp := http.ForwardProxy; fp != nil {
			result.ForwardProxyConvention = string(pointer.From(fp.Convention))
			result.ForwardProxyCustomHostHeaderName = pointer.From(fp.CustomHostHeaderName)
			result.ForwardProxyCustomSchemeHeaderName = pointer.From(fp.CustomProtoHeaderName)
		}
	}

	if login := settings.Login; login != nil {
		result.Login = flattenAuthV2LoginSettingsV20250501(login)
	}

	if authProviders := settings.IdentityProviders; authProviders != nil {
		result.AppleAuth = flattenAppleAuthV2SettingsV20250501(authProviders.Apple)
		result.AzureActiveDirectoryAuth = flattenAadAuthV2SettingsV20250501(authProviders.AzureActiveDirectory)
		result.AzureStaticWebAuth = flattenStaticWebAppAuthV2SettingsV20250501(authProviders.AzureStaticWebApps)
		result.CustomOIDCAuth = flattenCustomOIDCAuthV2SettingsV20250501(authProviders.CustomOpenIdConnectProviders)
		result.FacebookAuth = flattenFacebookAuthV2SettingsV20250501(authProviders.Facebook)
		result.GithubAuth = flattenGitHubAuthV2SettingsV20250501(authProviders.GitHub)
		result.GoogleAuth = flattenGoogleAuthV2SettingsV20250501(authProviders.Google)
		result.MicrosoftAuth = flattenMicrosoftAuthV2SettingsV20250501(authProviders.LegacyMicrosoftAccount)
		result.TwitterAuth = flattenTwitterAuthV2SettingsV20250501(authProviders.Twitter)
	}

	return []AuthV2Settings{result}
}

func DefaultAuthV2SettingsPropertiesV20250501() *webapps20250501.SiteAuthSettingsV2Properties {
	return &webapps20250501.SiteAuthSettingsV2Properties{
		Platform: &webapps20250501.AuthPlatform{
			Enabled:        pointer.To(false),
			RuntimeVersion: pointer.To("~1"),
			ConfigFilePath: pointer.To(""),
		},
		GlobalValidation: &webapps20250501.GlobalValidation{
			RequireAuthentication:       pointer.To(false),
			UnauthenticatedClientAction: pointer.To(webapps20250501.UnauthenticatedClientActionV2RedirectToLoginPage),
			ExcludedPaths:               pointer.To([]string{}),
			RedirectToProvider:          pointer.To(""),
		},
		Login: &webapps20250501.Login{
			Routes: &webapps20250501.LoginRoutes{},
			TokenStore: &webapps20250501.TokenStore{
				Enabled:                    pointer.To(false),
				TokenRefreshExtensionHours: pointer.To(72.0),
				FileSystem:                 &webapps20250501.FileSystemTokenStore{},
				AzureBlobStorage:           &webapps20250501.BlobStorageTokenStore{},
			},
			PreserveURLFragmentsForLogins: pointer.To(false),
			Nonce: &webapps20250501.Nonce{
				ValidateNonce:           pointer.To(true),
				NonceExpirationInterval: pointer.To("00:05:00"),
			},
			CookieExpiration: &webapps20250501.CookieExpiration{
				Convention:       pointer.To(webapps20250501.CookieExpirationConventionFixedTime),
				TimeToExpiration: pointer.To("08:00:00"),
			},
			AllowedExternalRedirectURLs: pointer.To([]string{}),
		},
		HTTPSettings: &webapps20250501.HTTPSettings{
			RequireHTTPS: pointer.To(true),
			Routes: &webapps20250501.HTTPSettingsRoutes{
				ApiPrefix: pointer.To("/.auth"),
			},
			ForwardProxy: &webapps20250501.ForwardProxy{
				Convention: pointer.To(webapps20250501.ForwardProxyConventionNoProxy),
			},
		},
		IdentityProviders: &webapps20250501.IdentityProviders{
			AzureActiveDirectory: &webapps20250501.AzureActiveDirectory{
				Enabled:      pointer.To(false),
				Registration: &webapps20250501.AzureActiveDirectoryRegistration{},
				Login: &webapps20250501.AzureActiveDirectoryLogin{
					DisableWWWAuthenticate: pointer.To(false),
				},
				Validation: &webapps20250501.AzureActiveDirectoryValidation{
					JwtClaimChecks: &webapps20250501.JwtClaimChecks{},
					DefaultAuthorizationPolicy: &webapps20250501.DefaultAuthorizationPolicy{
						AllowedPrincipals:   &webapps20250501.AllowedPrincipals{},
						AllowedApplications: pointer.To([]string{}),
					},
				},
			},
			Facebook: &webapps20250501.Facebook{
				Enabled:      pointer.To(false),
				Registration: &webapps20250501.AppRegistration{},
				Login:        &webapps20250501.LoginScopes{},
			},
			GitHub: &webapps20250501.GitHub{
				Enabled:      pointer.To(false),
				Registration: &webapps20250501.ClientRegistration{},
				Login:        &webapps20250501.LoginScopes{},
			},
			Google: &webapps20250501.Google{
				Enabled:      pointer.To(false),
				Registration: &webapps20250501.ClientRegistration{},
				Login:        &webapps20250501.LoginScopes{},
				Validation:   &webapps20250501.AllowedAudiencesValidation{},
			},
			Twitter: &webapps20250501.Twitter{
				Enabled:      pointer.To(false),
				Registration: &webapps20250501.TwitterRegistration{},
			},
			CustomOpenIdConnectProviders: pointer.To(map[string]webapps20250501.CustomOpenIdConnectProvider{}),
			LegacyMicrosoftAccount: &webapps20250501.LegacyMicrosoftAccount{
				Enabled:      pointer.To(false),
				Registration: &webapps20250501.ClientRegistration{},
				Login:        &webapps20250501.LoginScopes{},
				Validation:   &webapps20250501.AllowedAudiencesValidation{},
			},
			Apple: &webapps20250501.Apple{
				Enabled:      pointer.To(false),
				Registration: &webapps20250501.AppleRegistration{},
				Login:        &webapps20250501.LoginScopes{},
			},
			AzureStaticWebApps: &webapps20250501.AzureStaticWebApps{
				Enabled:      pointer.To(false),
				Registration: &webapps20250501.AzureStaticWebAppsRegistration{},
			},
		},
	}
}
