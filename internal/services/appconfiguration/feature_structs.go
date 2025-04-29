// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package appconfiguration

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/mitchellh/mapstructure"
)

const (
	TargetingFilterName  = "Microsoft.Targeting"
	TimewindowFilterName = "Microsoft.TimeWindow"
	PercentageFilterName = "Microsoft.Percentage"
)

type ClientFilter struct {
	Filters []interface{}
}

func (p *ClientFilter) UnmarshalJSON(b []byte) error {
	var tempIntf []interface{}

	if err := json.Unmarshal(b, &tempIntf); err != nil {
		return err
	}

	filtersOut := make([]interface{}, 0)
	for _, filterRawIntf := range tempIntf {
		filterRaw, ok := filterRawIntf.(map[string]interface{})
		if !ok {
			return fmt.Errorf("wtf")
		}
		nameRaw, ok := filterRaw["name"]
		if !ok {
			return fmt.Errorf("missing name ...")
		}

		name := nameRaw.(string)
		switch strings.ToLower(name) {
		case "microsoft.targeting":
			{
				var out TargetingFeatureFilter
				mpc := mapstructure.DecoderConfig{TagName: "json", Result: &out}
				mpd, err := mapstructure.NewDecoder(&mpc)
				if err != nil {
					return err
				}
				err = mpd.Decode(filterRaw)
				if err != nil {
					return err
				}
				filtersOut = append(filtersOut, out)
			}
		case "microsoft.timewindow":
			{
				var out TimewindowFeatureFilter
				mpc := mapstructure.DecoderConfig{TagName: "json", Result: &out}
				mpd, err := mapstructure.NewDecoder(&mpc)
				if err != nil {
					return err
				}
				err = mpd.Decode(filterRaw)
				if err != nil {
					return err
				}
				filtersOut = append(filtersOut, out)
			}
		case "microsoft.percentage":
			{
				var out PercentageFeatureFilter
				mpc := mapstructure.DecoderConfig{TagName: "json", Result: &out}
				mpd, err := mapstructure.NewDecoder(&mpc)
				if err != nil {
					return err
				}
				err = mpd.Decode(filterRaw)
				if err != nil {
					return err
				}
				filtersOut = append(filtersOut, out)
			}

		default:
			return fmt.Errorf("unknown type %q", name)
		}
	}

	p.Filters = filtersOut
	return nil
}

func (p ClientFilter) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.Filters)
}

type CustomFilter struct {
	Name  string `tfschema:"name"`
	Value string `tfschema:"value"`
}

type PercentageFilterParameters struct {
	Value float64 `json:"Value"`
}

type PercentageFeatureFilter struct {
	Name       string                     `json:"name"`
	Parameters PercentageFilterParameters `json:"parameters"`
}

type TargetingGroupParameter struct {
	Name              string `json:"Name"              tfschema:"name"`
	RolloutPercentage int64  `json:"RolloutPercentage" tfschema:"rollout_percentage"`
}

type TargetingFilterParameters struct {
	Audience TargetingFilterAudience `json:"Audience"`
}

type TargetingFilterAudience struct {
	DefaultRolloutPercentage int64                      `json:"DefaultRolloutPercentage" tfschema:"default_rollout_percentage"`
	Exclusion                []TargetingFilterExclusion `json:"Exclusion"                tfschema:"exclusion"`
	Groups                   []TargetingGroupParameter  `json:"Groups"                   tfschema:"groups"`
	Users                    []string                   `json:"Users"                    tfschema:"users"`
}

type TargetingFilterExclusion struct {
	Groups []string `json:"Groups"                   tfschema:"groups"`
	Users  []string `json:"Users"                    tfschema:"users"`
}

type TargetingFeatureFilter struct {
	Name       string                    `json:"name"`
	Parameters TargetingFilterParameters `json:"parameters"`
}

type TimewindowFilterParameters struct {
	Start      string                       `json:"Start"      tfschema:"start"`
	End        string                       `json:"End"        tfschema:"end"`
	Recurrence []TimewindowFilterRecurrence `json:"Recurrence" tfschema:"recurrence"`
}

type TimewindowFilterRecurrence struct {
	Daily   []RecurrenceDaily  `tfschema:"daily"`
	Weekly  []RecurrenceWeekly `tfschema:"weekly"`
	EndDate string             `tfschema:"end_date"`
}

type RecurrenceDaily struct {
	Interval int64 `json:"interval" tfschema:"interval"`
}

type RecurrenceWeekly struct {
	Interval       int64    `json:"interval"          tfschema:"interval"`
	FirstDayOfWeek string   `json:"first_day_of_week" tfschema:"first_day_of_week"`
	DaysOfWeek     []string `json:"days_of_week"      tfschema:"days_of_week"`
}

type TimewindowFeatureFilter struct {
	Name       string                     `json:"name"`
	Parameters TimewindowFilterParameters `json:"parameters"`
}

type Conditions struct {
	ClientFilters ClientFilter `json:"client_filters"`
}

type FeatureValue struct {
	ID          string     `json:"id"`
	Description string     `json:"description"`
	Enabled     bool       `json:"enabled"`
	Conditions  Conditions `json:"conditions"`
}

type FeatureAllocation struct {
	DefaultVariantWhenDisabled string                           `tfschema:"default_variant_when_disabled"`
	DefaultVariantWhenEnabled  string                           `tfschema:"default_variant_when_enabled"`
	GroupOverride              []FeatureAllocationGroupOverride `tfschema:"group_override"`
	Percentile                 []FeatureAllocationPercentile    `tfschema:"percentile"`
	Seed                       string                           `tfschema:"seed"`
	UserOverride               []FeatureAllocationUserOverride  `tfschema:"user_override"`
}

type FeatureAllocationGroupOverride struct {
	Groups  []string `tfschema:"groups"`
	Variant string   `tfschema:"variant"`
}

type FeatureAllocationUserOverride struct {
	Users   []string `tfschema:"users"`
	Variant string   `tfschema:"variant"`
}

type FeatureAllocationPercentile struct {
	Variant string `json:"variant" tfschema:"variant"`
	From    int64  `json:"from"    tfschema:"from"`
	To      int64  `json:"to"      tfschema:"to"`
}

type FeatureVariant struct {
	Name  string `tfschema:"name"`
	Value string `tfschema:"value"`
}
