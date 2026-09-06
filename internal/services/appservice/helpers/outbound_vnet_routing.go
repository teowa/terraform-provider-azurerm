// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package helpers

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2025-05-01/webapps"
)

// ExpandOutboundVnetRouting builds the `OutboundVnetRouting` block that replaced the top level
// `vnetRouteAllEnabled`, `vnetImagePullEnabled` and `vnetBackupRestoreEnabled` properties on `SiteProperties`.
func ExpandOutboundVnetRouting(routeAllEnabled, imagePullEnabled, backupRestoreEnabled *bool) *webapps.OutboundVnetRouting {
	if routeAllEnabled == nil && imagePullEnabled == nil && backupRestoreEnabled == nil {
		return nil
	}

	return &webapps.OutboundVnetRouting{
		ApplicationTraffic:   routeAllEnabled,
		ImagePullTraffic:     imagePullEnabled,
		BackupRestoreTraffic: backupRestoreEnabled,
	}
}

func FlattenVnetImagePullEnabled(input *webapps.OutboundVnetRouting) *bool {
	if input == nil {
		return nil
	}
	return input.ImagePullTraffic
}

func FlattenVnetBackupRestoreEnabled(input *webapps.OutboundVnetRouting) *bool {
	if input == nil {
		return nil
	}
	return input.BackupRestoreTraffic
}

func FlattenVnetContentShareEnabled(input *webapps.OutboundVnetRouting) *bool {
	if input == nil {
		return nil
	}
	return input.ContentShareTraffic
}
