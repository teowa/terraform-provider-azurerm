// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package helpers

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	webapps20250501 "github.com/hashicorp/go-azure-sdk/resource-manager/web/2025-05-01/webapps"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

// This file contains duplicates of helper functions in function_app_schema.go typed against the
// 2025-05-01 webapps SDK package, scoped exclusively to `azurerm_function_app_flex_consumption`.

func ExpandSiteConfigFunctionFlexConsumptionAppV20250501(siteConfigFlexConsumption []SiteConfigFunctionAppFlexConsumption, existing *webapps20250501.SiteConfig, metadata sdk.ResourceMetaData, storageUsesMSI bool, storageStringFlex string, storageConnStringForFCApp string) (*webapps20250501.SiteConfig, error) {
	if len(siteConfigFlexConsumption) == 0 {
		return nil, nil
	}

	expanded := &webapps20250501.SiteConfig{}
	if existing != nil {
		expanded = existing
	}

	appSettings := make([]webapps20250501.NameValuePair, 0)

	if existing != nil && existing.AppSettings != nil {
		appSettings = *existing.AppSettings
	}

	if storageStringFlex != "" {
		appSettings = updateOrAppendAppSettingsV20250501(appSettings, "AzureWebJobsStorage", storageStringFlex, false)
		if storageConnStringForFCApp != "" {
			appSettings = updateOrAppendAppSettingsV20250501(appSettings, storageConnStringForFCApp, storageStringFlex, false)
		}
	}

	FlexConsumptionSiteConfig := siteConfigFlexConsumption[0]

	v := strconv.FormatInt(FlexConsumptionSiteConfig.HealthCheckEvictionTime, 10)
	if v == "0" || FlexConsumptionSiteConfig.HealthCheckPath == "" {
		appSettings = updateOrAppendAppSettingsV20250501(appSettings, "WEBSITE_HEALTHCHECK_MAXPINGFAILURES", v, true)
	} else {
		appSettings = updateOrAppendAppSettingsV20250501(appSettings, "WEBSITE_HEALTHCHECK_MAXPINGFAILURES", v, false)
	}

	if FlexConsumptionSiteConfig.AppInsightsConnectionString == "" {
		appSettings = updateOrAppendAppSettingsV20250501(appSettings, "APPLICATIONINSIGHTS_CONNECTION_STRING", FlexConsumptionSiteConfig.AppInsightsConnectionString, true)
	} else {
		appSettings = updateOrAppendAppSettingsV20250501(appSettings, "APPLICATIONINSIGHTS_CONNECTION_STRING", FlexConsumptionSiteConfig.AppInsightsConnectionString, false)
	}

	if FlexConsumptionSiteConfig.AppInsightsInstrumentationKey == "" {
		appSettings = updateOrAppendAppSettingsV20250501(appSettings, "APPINSIGHTS_INSTRUMENTATIONKEY", FlexConsumptionSiteConfig.AppInsightsInstrumentationKey, true)
	} else {
		appSettings = updateOrAppendAppSettingsV20250501(appSettings, "APPINSIGHTS_INSTRUMENTATIONKEY", FlexConsumptionSiteConfig.AppInsightsInstrumentationKey, false)
	}

	if metadata.ResourceData.HasChange("site_config.0.api_management_api_id") {
		expanded.ApiManagementConfig = &webapps20250501.ApiManagementConfig{
			Id: pointer.To(FlexConsumptionSiteConfig.ApiManagementConfigId),
		}
	}

	if metadata.ResourceData.HasChange("site_config.0.api_definition_url") {
		expanded.ApiDefinition = &webapps20250501.ApiDefinitionInfo{
			Url: pointer.To(FlexConsumptionSiteConfig.ApiDefinition),
		}
	}

	if metadata.ResourceData.HasChange("site_config.0.app_command_line") {
		expanded.AppCommandLine = pointer.To(FlexConsumptionSiteConfig.AppCommandLine)
	}

	if metadata.ResourceData.HasChange("site_config.0.container_registry_use_managed_identity") {
		expanded.AcrUseManagedIdentityCreds = pointer.To(FlexConsumptionSiteConfig.UseManagedIdentityACR)
	}

	if metadata.ResourceData.HasChange("site_config.0.default_documents") {
		expanded.DefaultDocuments = &FlexConsumptionSiteConfig.DefaultDocuments
	}

	if metadata.ResourceData.HasChange("site_config.0.http2_enabled") {
		expanded.HTTP20Enabled = pointer.To(FlexConsumptionSiteConfig.Http2Enabled)
	}

	if metadata.ResourceData.HasChange("site_config.0.ip_restriction") {
		ipRestrictions, err := ExpandIpRestrictionsV20250501(FlexConsumptionSiteConfig.IpRestriction)
		if err != nil {
			return nil, err
		}
		expanded.IPSecurityRestrictions = ipRestrictions
	}

	if metadata.ResourceData.HasChange("site_config.0.ip_restriction_default_action") {
		expanded.IPSecurityRestrictionsDefaultAction = pointer.ToEnum[webapps20250501.DefaultAction](FlexConsumptionSiteConfig.IpRestrictionDefaultAction)
	}

	if metadata.ResourceData.HasChange("site_config.0.scm_use_main_ip_restriction") {
		expanded.ScmIPSecurityRestrictionsUseMain = pointer.To(FlexConsumptionSiteConfig.ScmUseMainIpRestriction)
	}

	if metadata.ResourceData.HasChange("site_config.0.scm_ip_restriction") {
		scmIpRestrictions, err := ExpandIpRestrictionsV20250501(FlexConsumptionSiteConfig.ScmIpRestriction)
		if err != nil {
			return nil, err
		}
		expanded.ScmIPSecurityRestrictions = scmIpRestrictions
	}

	if metadata.ResourceData.HasChange("site_config.0.scm_ip_restriction_default_action") {
		expanded.ScmIPSecurityRestrictionsDefaultAction = pointer.ToEnum[webapps20250501.DefaultAction](FlexConsumptionSiteConfig.ScmIpRestrictionDefaultAction)
	}

	if metadata.ResourceData.HasChange("site_config.0.load_balancing_mode") {
		expanded.LoadBalancing = pointer.ToEnum[webapps20250501.SiteLoadBalancing](FlexConsumptionSiteConfig.LoadBalancing)
	}

	if metadata.ResourceData.HasChange("site_config.0.managed_pipeline_mode") {
		expanded.ManagedPipelineMode = pointer.ToEnum[webapps20250501.ManagedPipelineMode](FlexConsumptionSiteConfig.ManagedPipelineMode)
	}

	if metadata.ResourceData.HasChange("site_config.0.remote_debugging_enabled") {
		expanded.RemoteDebuggingEnabled = pointer.To(FlexConsumptionSiteConfig.RemoteDebugging)
	}

	if metadata.ResourceData.HasChange("site_config.0.remote_debugging_version") {
		expanded.RemoteDebuggingVersion = pointer.To(FlexConsumptionSiteConfig.RemoteDebuggingVersion)
	}

	if metadata.ResourceData.HasChange("site_config.0.vnet_route_all_enabled") {
		expanded.VnetRouteAllEnabled = pointer.To(FlexConsumptionSiteConfig.VnetRouteAllEnabled)
	}

	if metadata.ResourceData.HasChange("site_config.0.websockets_enabled") {
		expanded.WebSocketsEnabled = pointer.To(FlexConsumptionSiteConfig.WebSockets)
	}

	if metadata.ResourceData.HasChange("site_config.0.health_check_path") {
		expanded.HealthCheckPath = pointer.To(FlexConsumptionSiteConfig.HealthCheckPath)
	}

	if metadata.ResourceData.HasChange("site_config.0.worker_count") {
		expanded.NumberOfWorkers = pointer.To(FlexConsumptionSiteConfig.WorkerCount)
	}

	if metadata.ResourceData.HasChange("site_config.0.minimum_tls_version") {
		expanded.MinTlsVersion = pointer.ToEnum[webapps20250501.SupportedTlsVersions](FlexConsumptionSiteConfig.MinTlsVersion)
	}

	if metadata.ResourceData.HasChange("site_config.0.scm_minimum_tls_version") {
		expanded.ScmMinTlsVersion = pointer.ToEnum[webapps20250501.SupportedTlsVersions](FlexConsumptionSiteConfig.ScmMinTlsVersion)
	}

	if metadata.ResourceData.HasChange("site_config.0.cors") {
		expanded.Cors = ExpandCorsSettingsV20250501(FlexConsumptionSiteConfig.Cors)
	}

	if metadata.ResourceData.HasChange("site_config.0.elastic_instance_minimum") {
		expanded.MinimumElasticInstanceCount = pointer.To(FlexConsumptionSiteConfig.ElasticInstanceMinimum)
	}

	if metadata.ResourceData.HasChange("site_config.0.runtime_scale_monitoring_enabled") {
		expanded.FunctionsRuntimeScaleMonitoringEnabled = pointer.To(FlexConsumptionSiteConfig.RuntimeScaleMonitoring)
	}

	expanded.AppSettings = &appSettings

	return expanded, nil
}

// updateOrAppendAppSettingsV20250501 is used to modify a collection of webapps20250501.NameValuePair items.
func updateOrAppendAppSettingsV20250501(input []webapps20250501.NameValuePair, name string, value string, remove bool) []webapps20250501.NameValuePair {
	for k, v := range input {
		if v.Name != nil && *v.Name == name {
			if remove {
				input[k] = input[len(input)-1]
				input[len(input)-1] = webapps20250501.NameValuePair{}
				input = input[:len(input)-1]
			} else {
				input[k] = webapps20250501.NameValuePair{
					Name:  pointer.To(name),
					Value: pointer.To(value),
				}
			}
			return input
		}
	}

	if !remove {
		input = append(input, webapps20250501.NameValuePair{
			Name:  pointer.To(name),
			Value: pointer.To(value),
		})
	}

	return input
}

func FlattenSiteConfigFunctionAppFlexConsumptionV20250501(functionAppFlexConsumptionSiteConfig *webapps20250501.SiteConfig) (*SiteConfigFunctionAppFlexConsumption, error) {
	if functionAppFlexConsumptionSiteConfig == nil {
		return nil, fmt.Errorf("flattening site config: SiteConfig was nil")
	}

	result := &SiteConfigFunctionAppFlexConsumption{
		AppCommandLine:                pointer.From(functionAppFlexConsumptionSiteConfig.AppCommandLine),
		ContainerRegistryMSI:          pointer.From(functionAppFlexConsumptionSiteConfig.AcrUserManagedIdentityID),
		Cors:                          FlattenCorsSettingsV20250501(functionAppFlexConsumptionSiteConfig.Cors),
		DetailedErrorLogging:          pointer.From(functionAppFlexConsumptionSiteConfig.DetailedErrorLoggingEnabled),
		HealthCheckPath:               pointer.From(functionAppFlexConsumptionSiteConfig.HealthCheckPath),
		IpRestrictionDefaultAction:    string(pointer.From(functionAppFlexConsumptionSiteConfig.IPSecurityRestrictionsDefaultAction)),
		ScmIpRestrictionDefaultAction: string(pointer.From(functionAppFlexConsumptionSiteConfig.ScmIPSecurityRestrictionsDefaultAction)),
		LoadBalancing:                 string(pointer.From(functionAppFlexConsumptionSiteConfig.LoadBalancing)),
		ManagedPipelineMode:           string(pointer.From(functionAppFlexConsumptionSiteConfig.ManagedPipelineMode)),
		WorkerCount:                   pointer.From(functionAppFlexConsumptionSiteConfig.NumberOfWorkers),
		ScmType:                       string(pointer.From(functionAppFlexConsumptionSiteConfig.ScmType)),
		RuntimeScaleMonitoring:        pointer.From(functionAppFlexConsumptionSiteConfig.FunctionsRuntimeScaleMonitoringEnabled),
		MinTlsVersion:                 string(pointer.From(functionAppFlexConsumptionSiteConfig.MinTlsVersion)),
		ScmMinTlsVersion:              string(pointer.From(functionAppFlexConsumptionSiteConfig.ScmMinTlsVersion)),
		WebSockets:                    pointer.From(functionAppFlexConsumptionSiteConfig.WebSocketsEnabled),
		ScmUseMainIpRestriction:       pointer.From(functionAppFlexConsumptionSiteConfig.ScmIPSecurityRestrictionsUseMain),
		UseManagedIdentityACR:         pointer.From(functionAppFlexConsumptionSiteConfig.AcrUseManagedIdentityCreds),
		RemoteDebugging:               pointer.From(functionAppFlexConsumptionSiteConfig.RemoteDebuggingEnabled),
		RemoteDebuggingVersion:        strings.ToUpper(pointer.From(functionAppFlexConsumptionSiteConfig.RemoteDebuggingVersion)),
		Http2Enabled:                  pointer.From(functionAppFlexConsumptionSiteConfig.HTTP20Enabled),
		VnetRouteAllEnabled:           pointer.From(functionAppFlexConsumptionSiteConfig.VnetRouteAllEnabled),
	}

	if v := functionAppFlexConsumptionSiteConfig.ApiDefinition; v != nil && v.Url != nil {
		result.ApiDefinition = *v.Url
	}

	if v := functionAppFlexConsumptionSiteConfig.ApiManagementConfig; v != nil && v.Id != nil {
		result.ApiManagementConfigId = *v.Id
	}

	if functionAppFlexConsumptionSiteConfig.IPSecurityRestrictions != nil {
		result.IpRestriction = FlattenIpRestrictionsV20250501(functionAppFlexConsumptionSiteConfig.IPSecurityRestrictions)
	}

	if functionAppFlexConsumptionSiteConfig.ScmIPSecurityRestrictions != nil {
		result.ScmIpRestriction = FlattenIpRestrictionsV20250501(functionAppFlexConsumptionSiteConfig.ScmIPSecurityRestrictions)
	}

	if v := functionAppFlexConsumptionSiteConfig.DefaultDocuments; v != nil {
		result.DefaultDocuments = *v
	}

	return result, nil
}

func MergeUserAppSettingsV20250501(systemSettings *[]webapps20250501.NameValuePair, userSettings map[string]string) *[]webapps20250501.NameValuePair {
	if len(userSettings) == 0 {
		return systemSettings
	}
	combined := *systemSettings
	for k, v := range userSettings {
		// Dedupe, explicit user settings take priority over enumerated, e.g. specifying KeyVault for `AzureWebJobsStorage`
		for i, x := range combined {
			if x.Name != nil && strings.EqualFold(*x.Name, k) {
				copy(combined[i:], combined[i+1:])
				combined = combined[:len(combined)-1]
			}
		}
		combined = append(combined, webapps20250501.NameValuePair{
			Name:  pointer.To(k),
			Value: pointer.To(v),
		})
	}
	return &combined
}

func ExpandFunctionAppAppServiceLogsV20250501(input []FunctionAppAppServiceLogs) webapps20250501.SiteLogsConfig {
	if len(input) == 0 {
		return webapps20250501.SiteLogsConfig{
			Properties: &webapps20250501.SiteLogsConfigProperties{
				HTTPLogs: &webapps20250501.HTTPLogsConfig{
					FileSystem: &webapps20250501.FileSystemHTTPLogsConfig{
						Enabled: pointer.To(false),
					},
				},
			},
		}
	}

	config := input[0]
	return webapps20250501.SiteLogsConfig{
		Properties: &webapps20250501.SiteLogsConfigProperties{
			HTTPLogs: &webapps20250501.HTTPLogsConfig{
				FileSystem: &webapps20250501.FileSystemHTTPLogsConfig{
					RetentionInDays: pointer.To(config.RetentionPeriodDays),
					RetentionInMb:   pointer.To(config.DiskQuotaMB),
					Enabled:         pointer.To(true),
				},
			},
		},
	}
}
