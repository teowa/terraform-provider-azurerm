// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package monitor

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/monitor/2026-04-01/pipelinegroups"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var _ sdk.ResourceWithUpdate = MonitorPipelineGroupResource{}

type MonitorPipelineGroupResource struct{}

type MonitorPipelineGroupResourceModel struct {
	Name               string                                      `tfschema:"name"`
	ResourceGroupName  string                                      `tfschema:"resource_group_name"`
	Location           string                                      `tfschema:"location"`
	Exporter           []MonitorPipelineGroupExporterModel         `tfschema:"exporter"`
	Receiver           []MonitorPipelineGroupReceiverModel         `tfschema:"receiver"`
	Service            []MonitorPipelineGroupServiceModel          `tfschema:"service"`
	ExecutionPlacement []ExecutionPlacementModel                   `tfschema:"execution_placement"`
	ExtendedLocation   []ExtendedLocationModel                     `tfschema:"extended_location"`
	Processor          []MonitorPipelineGroupProcessorModel        `tfschema:"processor"`
	Replicas           int64                                       `tfschema:"replicas"`
	TlsConfiguration   []MonitorPipelineGroupTLSConfigurationModel `tfschema:"tls_configuration"`
	Tags               map[string]interface{}                      `tfschema:"tags"`
}

type ExecutionPlacementModel struct {
	Constraint              []PlacementConstraintModel `tfschema:"constraint"`
	MaximumInstancesPerHost int64                      `tfschema:"maximum_instances_per_host"`
}

type PlacementConstraintModel struct {
	Capability string   `tfschema:"capability"`
	Operator   string   `tfschema:"operator"`
	Values     []string `tfschema:"values"`
}

type ExtendedLocationModel struct {
	Name string `tfschema:"name"`
	Type string `tfschema:"type"`
}

type MonitorPipelineGroupExporterModel struct {
	Name                      string                                   `tfschema:"name"`
	Type                      string                                   `tfschema:"type"`
	AzureMonitorWorkspaceLogs []AzureMonitorWorkspaceLogsExporterModel `tfschema:"azure_monitor_workspace_logs"`
}

type AzureMonitorWorkspaceLogsExporterModel struct {
	API         []AzureMonitorWorkspaceLogsAPIConfigModel `tfschema:"api"`
	Persistence []ExporterPersistenceConfigurationModel   `tfschema:"persistence"`
}

type AzureMonitorWorkspaceLogsAPIConfigModel struct {
	DataCollectionEndpointURL string           `tfschema:"data_collection_endpoint_url"`
	DataCollectionRule        string           `tfschema:"data_collection_rule"`
	Schema                    []SchemaMapModel `tfschema:"schema"`
	Stream                    string           `tfschema:"stream"`
}

type SchemaMapModel struct {
	RecordMap   []MapValueModel `tfschema:"record_map"`
	ResourceMap []MapValueModel `tfschema:"resource_map"`
	ScopeMap    []MapValueModel `tfschema:"scope_map"`
}

type MapValueModel struct {
	From string `tfschema:"from"`
	To   string `tfschema:"to"`
}

type ExporterPersistenceConfigurationModel struct {
	MaximumStorageUsage int64 `tfschema:"maximum_storage_usage"`
	RetentionPeriod     int64 `tfschema:"retention_period"`
}

type MonitorPipelineGroupReceiverModel struct {
	Name                 string                `tfschema:"name"`
	Type                 string                `tfschema:"type"`
	OtlpEndpoint         string                `tfschema:"otlp_endpoint"`
	Syslog               []SyslogReceiverModel `tfschema:"syslog"`
	TlsConfigurationName string                `tfschema:"tls_configuration_name"`
}

type SyslogReceiverModel struct {
	AllowSkipPriHeader bool     `tfschema:"allow_skip_pri_header"`
	AllowedFormats     []string `tfschema:"allowed_formats"`
	Endpoint           string   `tfschema:"endpoint"`
	TransportProtocol  string   `tfschema:"transport_protocol"`
}

type MonitorPipelineGroupProcessorModel struct {
	Name               string                `tfschema:"name"`
	Type               string                `tfschema:"type"`
	Batch              []BatchProcessorModel `tfschema:"batch"`
	TransformStatement string                `tfschema:"transform_statement"`
}

type BatchProcessorModel struct {
	BatchSize int64 `tfschema:"batch_size"`
	Timeout   int64 `tfschema:"timeout"`
}

type MonitorPipelineGroupServiceModel struct {
	PersistentVolumeName string          `tfschema:"persistent_volume_name"`
	Pipeline             []PipelineModel `tfschema:"pipeline"`
}

type PipelineModel struct {
	Exporter  []string `tfschema:"exporter"`
	Name      string   `tfschema:"name"`
	Processor []string `tfschema:"processor"`
	Receiver  []string `tfschema:"receiver"`
	Type      string   `tfschema:"type"`
}

type MonitorPipelineGroupTLSConfigurationModel struct {
	Name           string                    `tfschema:"name"`
	Mode           string                    `tfschema:"mode"`
	ClientCA       []CertificateSourceModel  `tfschema:"client_ca"`
	TLSCertificate []CertificateWithKeyModel `tfschema:"tls_certificate"`
}

type CertificateWithKeyModel struct {
	Certificate []CertificateSourceModel `tfschema:"certificate"`
	PrivateKey  []PrivateKeySourceModel  `tfschema:"private_key"`
}

type CertificateSourceModel struct {
	Location    string `tfschema:"location"`
	SubLocation string `tfschema:"sub_location"`
	Type        string `tfschema:"type"`
}

type PrivateKeySourceModel struct {
	Location    string `tfschema:"location"`
	SubLocation string `tfschema:"sub_location"`
	Type        string `tfschema:"type"`
}

func (r MonitorPipelineGroupResource) ResourceType() string {
	return "azurerm_monitor_pipeline_group"
}

func (r MonitorPipelineGroupResource) ModelObject() interface{} {
	return &MonitorPipelineGroupResourceModel{}
}

func (r MonitorPipelineGroupResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return pipelinegroups.ValidatePipelineGroupID
}

func (r MonitorPipelineGroupResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"exporter": monitorPipelineGroupExporterSchema(),

		"receiver": monitorPipelineGroupReceiverSchema(),

		"service": monitorPipelineGroupServiceSchema(),

		"execution_placement": monitorPipelineGroupExecutionPlacementSchema(),

		"extended_location": monitorPipelineGroupExtendedLocationSchema(),

		"processor": monitorPipelineGroupProcessorSchema(),

		"replicas": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntAtLeast(1),
		},

		"tls_configuration": monitorPipelineGroupTLSConfigurationSchema(),

		"tags": commonschema.Tags(),
	}
}

func (r MonitorPipelineGroupResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r MonitorPipelineGroupResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Monitor.PipelineGroupsClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model MonitorPipelineGroupResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := pipelinegroups.NewPipelineGroupID(subscriptionId, model.ResourceGroupName, model.Name)
			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for existing %s: %+v", id, err)
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			payload, err := expandMonitorPipelineGroup(model)
			if err != nil {
				return err
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r MonitorPipelineGroupResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Monitor.PipelineGroupsClient

			id, err := pipelinegroups.ParsePipelineGroupID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if resp.Model == nil {
				return fmt.Errorf("retrieving %s: model was nil", *id)
			}

			if resp.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: properties were nil", *id)
			}

			state := MonitorPipelineGroupResourceModel{
				Name:               id.PipelineGroupName,
				ResourceGroupName:  id.ResourceGroupName,
				Location:           location.Normalize(resp.Model.Location),
				Exporter:           flattenMonitorPipelineGroupExporters(resp.Model.Properties.Exporters),
				ExecutionPlacement: flattenMonitorPipelineGroupExecutionPlacement(resp.Model.Properties.ExecutionPlacement),
				ExtendedLocation:   flattenMonitorPipelineGroupExtendedLocation(resp.Model.ExtendedLocation),
				Processor:          flattenMonitorPipelineGroupProcessors(resp.Model.Properties.Processors),
				Receiver:           flattenMonitorPipelineGroupReceivers(resp.Model.Properties.Receivers),
				Replicas:           flattenMonitorPipelineGroupReplicas(resp.Model.Properties.Replicas),
				Service:            flattenMonitorPipelineGroupService(resp.Model.Properties.Service),
				Tags:               tags.Flatten(resp.Model.Tags),
				TlsConfiguration:   flattenMonitorPipelineGroupTLSConfigurations(resp.Model.Properties.TlsConfigurations),
			}

			return metadata.Encode(&state)
		},
	}
}

func (r MonitorPipelineGroupResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Monitor.PipelineGroupsClient

			id, err := pipelinegroups.ParsePipelineGroupID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model MonitorPipelineGroupResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			payload, err := expandMonitorPipelineGroup(model)
			if err != nil {
				return err
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, payload); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r MonitorPipelineGroupResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Monitor.PipelineGroupsClient

			id, err := pipelinegroups.ParsePipelineGroupID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func expandMonitorPipelineGroup(input MonitorPipelineGroupResourceModel) (pipelinegroups.PipelineGroup, error) {
	exporters, exporterNames, err := expandMonitorPipelineGroupExporters(input.Exporter)
	if err != nil {
		return pipelinegroups.PipelineGroup{}, err
	}

	receivers, receiverNames, err := expandMonitorPipelineGroupReceivers(input.Receiver)
	if err != nil {
		return pipelinegroups.PipelineGroup{}, err
	}

	processors, processorNames, err := expandMonitorPipelineGroupProcessors(input.Processor)
	if err != nil {
		return pipelinegroups.PipelineGroup{}, err
	}

	tlsConfigurations, tlsConfigurationNames, err := expandMonitorPipelineGroupTLSConfigurations(input.TlsConfiguration)
	if err != nil {
		return pipelinegroups.PipelineGroup{}, err
	}

	if err := validateReceiverTLSConfigurationReferences(input.Receiver, tlsConfigurationNames); err != nil {
		return pipelinegroups.PipelineGroup{}, err
	}

	service, err := expandMonitorPipelineGroupService(input.Service, receiverNames, processorNames, exporterNames)
	if err != nil {
		return pipelinegroups.PipelineGroup{}, err
	}

	output := pipelinegroups.PipelineGroup{
		Location: location.Normalize(input.Location),
		Properties: &pipelinegroups.PipelineGroupProperties{
			Exporters:  exporters,
			Processors: processors,
			Receivers:  receivers,
			Service:    service,
		},
		Tags: tags.Expand(input.Tags),
	}

	if len(input.ExecutionPlacement) > 0 {
		output.Properties.ExecutionPlacement = expandMonitorPipelineGroupExecutionPlacement(input.ExecutionPlacement)
	}

	if len(input.ExtendedLocation) > 0 {
		output.ExtendedLocation = expandMonitorPipelineGroupExtendedLocation(input.ExtendedLocation)
	}

	if input.Replicas > 0 {
		output.Properties.Replicas = &input.Replicas
	}

	if len(tlsConfigurations) > 0 {
		output.Properties.TlsConfigurations = &tlsConfigurations
	}

	return output, nil
}

func expandMonitorPipelineGroupExecutionPlacement(input []ExecutionPlacementModel) *pipelinegroups.ExecutionPlacement {
	if len(input) == 0 {
		return nil
	}

	model := input[0]
	output := &pipelinegroups.ExecutionPlacement{}
	if len(model.Constraint) > 0 {
		constraints := make([]pipelinegroups.PlacementConstraint, 0, len(model.Constraint))
		for _, item := range model.Constraint {
			constraint := pipelinegroups.PlacementConstraint{
				Capability: item.Capability,
				Operator:   pipelinegroups.CapabilityOperator(item.Operator),
			}
			if len(item.Values) > 0 {
				values := append([]string{}, item.Values...)
				constraint.Values = &values
			}
			constraints = append(constraints, constraint)
		}
		output.Constraints = &constraints
	}

	if model.MaximumInstancesPerHost > 0 {
		distribution := &pipelinegroups.DistributionPolicy{}
		distribution.MaxInstancesPerHost = &model.MaximumInstancesPerHost
		output.Distribution = distribution
	}

	return output
}

func expandMonitorPipelineGroupExtendedLocation(input []ExtendedLocationModel) *pipelinegroups.AzureResourceManagerCommonTypesExtendedLocation {
	if len(input) == 0 {
		return nil
	}

	model := input[0]
	return &pipelinegroups.AzureResourceManagerCommonTypesExtendedLocation{
		Name: model.Name,
		Type: pipelinegroups.ExtendedLocationType(model.Type),
	}
}

func expandMonitorPipelineGroupExporters(input []MonitorPipelineGroupExporterModel) ([]pipelinegroups.Exporter, map[string]struct{}, error) {
	names := make(map[string]struct{}, len(input))
	output := make([]pipelinegroups.Exporter, 0, len(input))
	for _, item := range input {
		if err := ensureUniqueName("exporter", item.Name, names); err != nil {
			return nil, nil, err
		}

		exporter := pipelinegroups.Exporter{
			Name: item.Name,
			Type: pipelinegroups.ExporterType(item.Type),
		}

		switch exporter.Type {
		case pipelinegroups.ExporterTypeAzureMonitorWorkspaceLogs:
			if len(item.AzureMonitorWorkspaceLogs) != 1 {
				return nil, nil, fmt.Errorf("field `azure_monitor_workspace_logs` must be specified when `type` is `%s` for exporter `%s`", exporter.Type, item.Name)
			}
			exporter.AzureMonitorWorkspaceLogs = expandMonitorPipelineGroupAzureMonitorWorkspaceLogsExporter(item.AzureMonitorWorkspaceLogs[0])
		default:
			if len(item.AzureMonitorWorkspaceLogs) > 0 {
				return nil, nil, fmt.Errorf("field `azure_monitor_workspace_logs` cannot be specified when `type` is `%s` for exporter `%s`", exporter.Type, item.Name)
			}
		}

		output = append(output, exporter)
	}

	return output, names, nil
}

func expandMonitorPipelineGroupAzureMonitorWorkspaceLogsExporter(input AzureMonitorWorkspaceLogsExporterModel) *pipelinegroups.AzureMonitorWorkspaceLogsExporter {
	api := input.API[0]
	schema := api.Schema[0]
	output := &pipelinegroups.AzureMonitorWorkspaceLogsExporter{
		Api: pipelinegroups.AzureMonitorWorkspaceLogsApiConfig{
			DataCollectionEndpointURL: api.DataCollectionEndpointURL,
			DataCollectionRule:        api.DataCollectionRule,
			Schema: pipelinegroups.SchemaMap{
				RecordMap:   expandMonitorPipelineGroupRecordMaps(schema.RecordMap),
				ResourceMap: expandOptionalMonitorPipelineGroupResourceMaps(schema.ResourceMap),
				ScopeMap:    expandOptionalMonitorPipelineGroupScopeMaps(schema.ScopeMap),
			},
			Stream: api.Stream,
		},
	}

	if len(input.Persistence) > 0 {
		persistence := pipelinegroups.ExporterPersistenceConfiguration{}
		if input.Persistence[0].MaximumStorageUsage > 0 {
			persistence.MaxStorageUsage = &input.Persistence[0].MaximumStorageUsage
		}
		if input.Persistence[0].RetentionPeriod > 0 {
			persistence.RetentionPeriod = &input.Persistence[0].RetentionPeriod
		}
		output.Persistence = &persistence
	}

	return output
}

func expandMonitorPipelineGroupRecordMaps(input []MapValueModel) []pipelinegroups.RecordMap {
	output := make([]pipelinegroups.RecordMap, 0, len(input))
	for _, item := range input {
		output = append(output, pipelinegroups.RecordMap{
			From: item.From,
			To:   item.To,
		})
	}
	return output
}

func expandOptionalMonitorPipelineGroupResourceMaps(input []MapValueModel) *[]pipelinegroups.ResourceMap {
	if len(input) == 0 {
		return nil
	}

	output := make([]pipelinegroups.ResourceMap, 0, len(input))
	for _, item := range input {
		output = append(output, pipelinegroups.ResourceMap{
			From: item.From,
			To:   item.To,
		})
	}
	return &output
}

func expandOptionalMonitorPipelineGroupScopeMaps(input []MapValueModel) *[]pipelinegroups.ScopeMap {
	if len(input) == 0 {
		return nil
	}

	output := make([]pipelinegroups.ScopeMap, 0, len(input))
	for _, item := range input {
		output = append(output, pipelinegroups.ScopeMap{
			From: item.From,
			To:   item.To,
		})
	}
	return &output
}

func expandMonitorPipelineGroupReceivers(input []MonitorPipelineGroupReceiverModel) ([]pipelinegroups.Receiver, map[string]struct{}, error) {
	names := make(map[string]struct{}, len(input))
	output := make([]pipelinegroups.Receiver, 0, len(input))
	for _, item := range input {
		if err := ensureUniqueName("receiver", item.Name, names); err != nil {
			return nil, nil, err
		}

		receiver := pipelinegroups.Receiver{
			Name: item.Name,
			Type: pipelinegroups.ReceiverType(item.Type),
		}
		if item.TlsConfigurationName != "" {
			receiver.TlsConfiguration = &item.TlsConfigurationName
		}

		switch receiver.Type {
		case pipelinegroups.ReceiverTypeOTLP:
			if item.OtlpEndpoint == "" {
				return nil, nil, fmt.Errorf("field `otlp_endpoint` must be specified when `type` is `%s` for receiver `%s`", receiver.Type, item.Name)
			}
			if len(item.Syslog) > 0 {
				return nil, nil, fmt.Errorf("field `syslog` cannot be specified when `type` is `%s` for receiver `%s`", receiver.Type, item.Name)
			}
			receiver.Otlp = &pipelinegroups.OtlpReceiver{
				Endpoint: item.OtlpEndpoint,
			}
		case pipelinegroups.ReceiverTypeSyslog:
			if len(item.Syslog) != 1 {
				return nil, nil, fmt.Errorf("field `syslog` must be specified when `type` is `%s` for receiver `%s`", receiver.Type, item.Name)
			}
			if item.OtlpEndpoint != "" {
				return nil, nil, fmt.Errorf("field `otlp_endpoint` cannot be specified when `type` is `%s` for receiver `%s`", receiver.Type, item.Name)
			}
			syslog := pipelinegroups.SyslogReceiver{
				AllowSkipPriHeader: &item.Syslog[0].AllowSkipPriHeader,
				Endpoint:           item.Syslog[0].Endpoint,
			}
			if len(item.Syslog[0].AllowedFormats) > 0 {
				allowedFormats := make([]pipelinegroups.AllowedFormats, 0, len(item.Syslog[0].AllowedFormats))
				for _, format := range item.Syslog[0].AllowedFormats {
					allowedFormats = append(allowedFormats, pipelinegroups.AllowedFormats(format))
				}
				syslog.AllowedFormats = &allowedFormats
			}
			if item.Syslog[0].TransportProtocol != "" {
				transportProtocol := pipelinegroups.TransportProtocol(item.Syslog[0].TransportProtocol)
				syslog.TransportProtocol = &transportProtocol
			}
			receiver.Syslog = &syslog
		default:
			if item.OtlpEndpoint != "" || len(item.Syslog) > 0 {
				return nil, nil, fmt.Errorf("fields `otlp_endpoint` and `syslog` cannot be specified when `type` is `%s` for receiver `%s`", receiver.Type, item.Name)
			}
		}

		output = append(output, receiver)
	}

	return output, names, nil
}

func expandMonitorPipelineGroupProcessors(input []MonitorPipelineGroupProcessorModel) ([]pipelinegroups.Processor, map[string]struct{}, error) {
	names := make(map[string]struct{}, len(input))
	output := make([]pipelinegroups.Processor, 0, len(input))
	for _, item := range input {
		if err := ensureUniqueName("processor", item.Name, names); err != nil {
			return nil, nil, err
		}

		processor := pipelinegroups.Processor{
			Name: item.Name,
			Type: pipelinegroups.ProcessorType(item.Type),
		}

		switch processor.Type {
		case pipelinegroups.ProcessorTypeBatch:
			if item.TransformStatement != "" {
				return nil, nil, fmt.Errorf("field `transform_statement` cannot be specified when `type` is `%s` for processor `%s`", processor.Type, item.Name)
			}
			if len(item.Batch) > 0 {
				batch := pipelinegroups.BatchProcessor{}
				if item.Batch[0].BatchSize > 0 {
					batch.BatchSize = &item.Batch[0].BatchSize
				}
				if item.Batch[0].Timeout > 0 {
					batch.Timeout = &item.Batch[0].Timeout
				}
				processor.Batch = &batch
			}
		case pipelinegroups.ProcessorTypeTransformLanguage:
			if item.TransformStatement == "" {
				return nil, nil, fmt.Errorf("field `transform_statement` must be specified when `type` is `%s` for processor `%s`", processor.Type, item.Name)
			}
			if len(item.Batch) > 0 {
				return nil, nil, fmt.Errorf("field `batch` cannot be specified when `type` is `%s` for processor `%s`", processor.Type, item.Name)
			}
			processor.TransformLanguage = &pipelinegroups.TransformLanguageProcessor{
				TransformStatement: item.TransformStatement,
			}
		default:
			if len(item.Batch) > 0 || item.TransformStatement != "" {
				return nil, nil, fmt.Errorf("fields `batch` and `transform_statement` cannot be specified when `type` is `%s` for processor `%s`", processor.Type, item.Name)
			}
		}

		output = append(output, processor)
	}

	return output, names, nil
}

func expandMonitorPipelineGroupService(input []MonitorPipelineGroupServiceModel, receiverNames, processorNames, exporterNames map[string]struct{}) (pipelinegroups.Service, error) {
	model := input[0]
	service := pipelinegroups.Service{
		Pipelines: make([]pipelinegroups.Pipeline, 0, len(model.Pipeline)),
	}

	if model.PersistentVolumeName != "" {
		service.Persistence = &pipelinegroups.PersistenceConfigurations{
			PersistentVolumeName: model.PersistentVolumeName,
		}
	}

	pipelineNames := make(map[string]struct{}, len(model.Pipeline))
	for _, item := range model.Pipeline {
		if err := ensureUniqueName("pipeline", item.Name, pipelineNames); err != nil {
			return pipelinegroups.Service{}, err
		}
		if err := validatePipelineGroupReferenceNames("receiver", item.Receiver, receiverNames, item.Name); err != nil {
			return pipelinegroups.Service{}, err
		}
		if err := validatePipelineGroupReferenceNames("processor", item.Processor, processorNames, item.Name); err != nil {
			return pipelinegroups.Service{}, err
		}
		if err := validatePipelineGroupReferenceNames("exporter", item.Exporter, exporterNames, item.Name); err != nil {
			return pipelinegroups.Service{}, err
		}

		pipeline := pipelinegroups.Pipeline{
			Exporters: append([]string{}, item.Exporter...),
			Name:      item.Name,
			Receivers: append([]string{}, item.Receiver...),
			Type:      pipelinegroups.PipelineType(item.Type),
		}
		if len(item.Processor) > 0 {
			processors := append([]string{}, item.Processor...)
			pipeline.Processors = &processors
		}
		service.Pipelines = append(service.Pipelines, pipeline)
	}

	return service, nil
}

func expandMonitorPipelineGroupTLSConfigurations(input []MonitorPipelineGroupTLSConfigurationModel) ([]pipelinegroups.TlsConfiguration, map[string]struct{}, error) {
	names := make(map[string]struct{}, len(input))
	output := make([]pipelinegroups.TlsConfiguration, 0, len(input))
	for _, item := range input {
		if err := ensureUniqueName("tls_configuration", item.Name, names); err != nil {
			return nil, nil, err
		}

		tlsConfiguration := pipelinegroups.TlsConfiguration{
			Name: item.Name,
		}
		if item.Mode != "" {
			mode := pipelinegroups.TlsMode(item.Mode)
			tlsConfiguration.Mode = &mode
		}
		if len(item.ClientCA) > 0 {
			tlsConfiguration.ClientCa = expandMonitorPipelineGroupCertificateSource(item.ClientCA[0])
		}
		if len(item.TLSCertificate) > 0 {
			tlsCertificate, err := expandMonitorPipelineGroupCertificateWithKey(item.TLSCertificate[0])
			if err != nil {
				return nil, nil, fmt.Errorf("expanding tls configuration `%s`: %+v", item.Name, err)
			}
			tlsConfiguration.TlsCertificate = tlsCertificate
		}

		output = append(output, tlsConfiguration)
	}

	return output, names, nil
}

func expandMonitorPipelineGroupCertificateWithKey(input CertificateWithKeyModel) (*pipelinegroups.CertificateWithKey, error) {
	if len(input.Certificate) != 1 {
		return nil, errors.New("field `certificate` must be specified when `tls_certificate` is configured")
	}
	if len(input.PrivateKey) != 1 {
		return nil, errors.New("field `private_key` must be specified when `tls_certificate` is configured")
	}

	return &pipelinegroups.CertificateWithKey{
		Certificate: *expandMonitorPipelineGroupCertificateSource(input.Certificate[0]),
		PrivateKey:  *expandMonitorPipelineGroupPrivateKeySource(input.PrivateKey[0]),
	}, nil
}

func expandMonitorPipelineGroupCertificateSource(input CertificateSourceModel) *pipelinegroups.CertificateSource {
	return &pipelinegroups.CertificateSource{
		Location:    input.Location,
		SubLocation: input.SubLocation,
		Type:        pipelinegroups.CertificateSourceType(input.Type),
	}
}

func expandMonitorPipelineGroupPrivateKeySource(input PrivateKeySourceModel) *pipelinegroups.PrivateKeySource {
	return &pipelinegroups.PrivateKeySource{
		Location:    input.Location,
		SubLocation: input.SubLocation,
		Type:        pipelinegroups.PrivateKeySourceType(input.Type),
	}
}

func flattenMonitorPipelineGroupExecutionPlacement(input *pipelinegroups.ExecutionPlacement) []ExecutionPlacementModel {
	if input == nil {
		return []ExecutionPlacementModel{}
	}

	model := ExecutionPlacementModel{
		Constraint: []PlacementConstraintModel{},
	}
	if input.Constraints != nil {
		for _, item := range *input.Constraints {
			constraint := PlacementConstraintModel{
				Capability: item.Capability,
				Operator:   string(item.Operator),
				Values:     []string{},
			}
			if item.Values != nil {
				constraint.Values = append(constraint.Values, *item.Values...)
			}
			model.Constraint = append(model.Constraint, constraint)
		}
	}
	if input.Distribution != nil {
		model.MaximumInstancesPerHost = flattenMonitorPipelineGroupInt64(input.Distribution.MaxInstancesPerHost)
	}

	return []ExecutionPlacementModel{model}
}

func flattenMonitorPipelineGroupExtendedLocation(input *pipelinegroups.AzureResourceManagerCommonTypesExtendedLocation) []ExtendedLocationModel {
	if input == nil {
		return []ExtendedLocationModel{}
	}

	return []ExtendedLocationModel{
		{
			Name: input.Name,
			Type: string(input.Type),
		},
	}
}

func flattenMonitorPipelineGroupExporters(input []pipelinegroups.Exporter) []MonitorPipelineGroupExporterModel {
	output := make([]MonitorPipelineGroupExporterModel, 0, len(input))
	for _, item := range input {
		model := MonitorPipelineGroupExporterModel{
			Name:                      item.Name,
			Type:                      string(item.Type),
			AzureMonitorWorkspaceLogs: []AzureMonitorWorkspaceLogsExporterModel{},
		}
		if item.AzureMonitorWorkspaceLogs != nil {
			model.AzureMonitorWorkspaceLogs = []AzureMonitorWorkspaceLogsExporterModel{
				flattenMonitorPipelineGroupAzureMonitorWorkspaceLogsExporter(item.AzureMonitorWorkspaceLogs),
			}
		}
		output = append(output, model)
	}
	return output
}

func flattenMonitorPipelineGroupAzureMonitorWorkspaceLogsExporter(input *pipelinegroups.AzureMonitorWorkspaceLogsExporter) AzureMonitorWorkspaceLogsExporterModel {
	model := AzureMonitorWorkspaceLogsExporterModel{
		API: []AzureMonitorWorkspaceLogsAPIConfigModel{
			{
				DataCollectionEndpointURL: input.Api.DataCollectionEndpointURL,
				DataCollectionRule:        input.Api.DataCollectionRule,
				Schema: []SchemaMapModel{
					{
						RecordMap:   flattenMonitorPipelineGroupRecordMaps(input.Api.Schema.RecordMap),
						ResourceMap: flattenMonitorPipelineGroupResourceMaps(input.Api.Schema.ResourceMap),
						ScopeMap:    flattenMonitorPipelineGroupScopeMaps(input.Api.Schema.ScopeMap),
					},
				},
				Stream: input.Api.Stream,
			},
		},
		Persistence: []ExporterPersistenceConfigurationModel{},
	}
	if input.Persistence != nil {
		model.Persistence = []ExporterPersistenceConfigurationModel{
			{
				MaximumStorageUsage: flattenMonitorPipelineGroupInt64(input.Persistence.MaxStorageUsage),
				RetentionPeriod:     flattenMonitorPipelineGroupInt64(input.Persistence.RetentionPeriod),
			},
		}
	}
	return model
}

func flattenMonitorPipelineGroupRecordMaps(input []pipelinegroups.RecordMap) []MapValueModel {
	output := make([]MapValueModel, 0, len(input))
	for _, item := range input {
		output = append(output, MapValueModel{
			From: item.From,
			To:   item.To,
		})
	}
	return output
}

func flattenMonitorPipelineGroupResourceMaps(input *[]pipelinegroups.ResourceMap) []MapValueModel {
	if input == nil {
		return []MapValueModel{}
	}

	output := make([]MapValueModel, 0, len(*input))
	for _, item := range *input {
		output = append(output, MapValueModel{
			From: item.From,
			To:   item.To,
		})
	}
	return output
}

func flattenMonitorPipelineGroupScopeMaps(input *[]pipelinegroups.ScopeMap) []MapValueModel {
	if input == nil {
		return []MapValueModel{}
	}

	output := make([]MapValueModel, 0, len(*input))
	for _, item := range *input {
		output = append(output, MapValueModel{
			From: item.From,
			To:   item.To,
		})
	}
	return output
}

func flattenMonitorPipelineGroupReceivers(input []pipelinegroups.Receiver) []MonitorPipelineGroupReceiverModel {
	output := make([]MonitorPipelineGroupReceiverModel, 0, len(input))
	for _, item := range input {
		model := MonitorPipelineGroupReceiverModel{
			Name:                 item.Name,
			Type:                 string(item.Type),
			Syslog:               []SyslogReceiverModel{},
			TlsConfigurationName: flattenMonitorPipelineGroupString(item.TlsConfiguration),
		}
		if item.Otlp != nil {
			model.OtlpEndpoint = item.Otlp.Endpoint
		}
		if item.Syslog != nil {
			syslog := SyslogReceiverModel{
				AllowSkipPriHeader: flattenMonitorPipelineGroupBool(item.Syslog.AllowSkipPriHeader),
				AllowedFormats:     []string{},
				Endpoint:           item.Syslog.Endpoint,
				TransportProtocol:  flattenMonitorPipelineGroupTransportProtocol(item.Syslog.TransportProtocol),
			}
			if item.Syslog.AllowedFormats != nil {
				for _, format := range *item.Syslog.AllowedFormats {
					syslog.AllowedFormats = append(syslog.AllowedFormats, string(format))
				}
			}
			model.Syslog = []SyslogReceiverModel{syslog}
		}
		output = append(output, model)
	}
	return output
}

func flattenMonitorPipelineGroupProcessors(input []pipelinegroups.Processor) []MonitorPipelineGroupProcessorModel {
	output := make([]MonitorPipelineGroupProcessorModel, 0, len(input))
	for _, item := range input {
		model := MonitorPipelineGroupProcessorModel{
			Name:  item.Name,
			Type:  string(item.Type),
			Batch: []BatchProcessorModel{},
		}
		if item.Batch != nil {
			model.Batch = []BatchProcessorModel{
				{
					BatchSize: flattenMonitorPipelineGroupInt64(item.Batch.BatchSize),
					Timeout:   flattenMonitorPipelineGroupInt64(item.Batch.Timeout),
				},
			}
		}
		if item.TransformLanguage != nil {
			model.TransformStatement = item.TransformLanguage.TransformStatement
		}
		output = append(output, model)
	}
	return output
}

func flattenMonitorPipelineGroupService(input pipelinegroups.Service) []MonitorPipelineGroupServiceModel {
	model := MonitorPipelineGroupServiceModel{
		Pipeline: make([]PipelineModel, 0, len(input.Pipelines)),
	}
	if input.Persistence != nil {
		model.PersistentVolumeName = input.Persistence.PersistentVolumeName
	}
	for _, item := range input.Pipelines {
		pipeline := PipelineModel{
			Exporter:  append([]string{}, item.Exporters...),
			Name:      item.Name,
			Processor: flattenMonitorPipelineGroupStringList(item.Processors),
			Receiver:  append([]string{}, item.Receivers...),
			Type:      string(item.Type),
		}
		model.Pipeline = append(model.Pipeline, pipeline)
	}
	return []MonitorPipelineGroupServiceModel{model}
}

func flattenMonitorPipelineGroupTLSConfigurations(input *[]pipelinegroups.TlsConfiguration) []MonitorPipelineGroupTLSConfigurationModel {
	if input == nil {
		return []MonitorPipelineGroupTLSConfigurationModel{}
	}

	output := make([]MonitorPipelineGroupTLSConfigurationModel, 0, len(*input))
	for _, item := range *input {
		model := MonitorPipelineGroupTLSConfigurationModel{
			Name:           item.Name,
			Mode:           flattenMonitorPipelineGroupTLSMode(item.Mode),
			ClientCA:       []CertificateSourceModel{},
			TLSCertificate: []CertificateWithKeyModel{},
		}
		if item.ClientCa != nil {
			model.ClientCA = []CertificateSourceModel{flattenMonitorPipelineGroupCertificateSource(item.ClientCa)}
		}
		if item.TlsCertificate != nil {
			model.TLSCertificate = []CertificateWithKeyModel{
				{
					Certificate: []CertificateSourceModel{flattenMonitorPipelineGroupCertificateSource(&item.TlsCertificate.Certificate)},
					PrivateKey:  []PrivateKeySourceModel{flattenMonitorPipelineGroupPrivateKeySource(&item.TlsCertificate.PrivateKey)},
				},
			}
		}
		output = append(output, model)
	}
	return output
}

func flattenMonitorPipelineGroupCertificateSource(input *pipelinegroups.CertificateSource) CertificateSourceModel {
	return CertificateSourceModel{
		Location:    input.Location,
		SubLocation: input.SubLocation,
		Type:        string(input.Type),
	}
}

func flattenMonitorPipelineGroupPrivateKeySource(input *pipelinegroups.PrivateKeySource) PrivateKeySourceModel {
	return PrivateKeySourceModel{
		Location:    input.Location,
		SubLocation: input.SubLocation,
		Type:        string(input.Type),
	}
}

func flattenMonitorPipelineGroupReplicas(input *int64) int64 {
	if input == nil {
		return 0
	}
	return *input
}

func flattenMonitorPipelineGroupInt64(input *int64) int64 {
	if input == nil {
		return 0
	}
	return *input
}

func flattenMonitorPipelineGroupBool(input *bool) bool {
	if input == nil {
		return false
	}
	return *input
}

func flattenMonitorPipelineGroupString(input *string) string {
	if input == nil {
		return ""
	}
	return *input
}

func flattenMonitorPipelineGroupStringList(input *[]string) []string {
	if input == nil {
		return []string{}
	}
	return append([]string{}, *input...)
}

func flattenMonitorPipelineGroupTLSMode(input *pipelinegroups.TlsMode) string {
	if input == nil {
		return ""
	}
	return string(*input)
}

func flattenMonitorPipelineGroupTransportProtocol(input *pipelinegroups.TransportProtocol) string {
	if input == nil {
		return ""
	}
	return string(*input)
}

func validateReceiverTLSConfigurationReferences(input []MonitorPipelineGroupReceiverModel, tlsConfigurationNames map[string]struct{}) error {
	for _, receiver := range input {
		if receiver.TlsConfigurationName == "" {
			continue
		}
		if _, ok := tlsConfigurationNames[receiver.TlsConfigurationName]; !ok {
			return fmt.Errorf("receiver `%s` references undefined `tls_configuration` `%s`", receiver.Name, receiver.TlsConfigurationName)
		}
	}
	return nil
}

func validatePipelineGroupReferenceNames(kind string, input []string, valid map[string]struct{}, pipelineName string) error {
	for _, item := range input {
		if _, ok := valid[item]; !ok {
			return fmt.Errorf("pipeline `%s` references undefined `%s` `%s`", pipelineName, kind, item)
		}
	}
	return nil
}

func ensureUniqueName(kind, name string, existing map[string]struct{}) error {
	if _, ok := existing[name]; ok {
		return fmt.Errorf("%s names must be unique, duplicate `%s` was specified", kind, name)
	}
	existing[name] = struct{}{}
	return nil
}

func monitorPipelineGroupExporterSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Required: true,
		MinItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"name": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},
				"type": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringInSlice(pipelinegroups.PossibleValuesForExporterType(), false),
				},
				"azure_monitor_workspace_logs": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					MaxItems: 1,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"api": {
								Type:     pluginsdk.TypeList,
								Required: true,
								MaxItems: 1,
								Elem: &pluginsdk.Resource{
									Schema: map[string]*pluginsdk.Schema{
										"data_collection_endpoint_url": {
											Type:         pluginsdk.TypeString,
											Required:     true,
											ValidateFunc: validation.StringIsNotEmpty,
										},
										"data_collection_rule": {
											Type:         pluginsdk.TypeString,
											Required:     true,
											ValidateFunc: validation.StringIsNotEmpty,
										},
										"schema": {
											Type:     pluginsdk.TypeList,
											Required: true,
											MaxItems: 1,
											Elem: &pluginsdk.Resource{
												Schema: map[string]*pluginsdk.Schema{
													"record_map":   monitorPipelineGroupMapValueSchema(true),
													"resource_map": monitorPipelineGroupMapValueSchema(false),
													"scope_map":    monitorPipelineGroupMapValueSchema(false),
												},
											},
										},
										"stream": {
											Type:         pluginsdk.TypeString,
											Required:     true,
											ValidateFunc: validation.StringIsNotEmpty,
										},
									},
								},
							},
							"persistence": {
								Type:     pluginsdk.TypeList,
								Optional: true,
								MaxItems: 1,
								Elem: &pluginsdk.Resource{
									Schema: map[string]*pluginsdk.Schema{
										"maximum_storage_usage": {
											Type:         pluginsdk.TypeInt,
											Optional:     true,
											ValidateFunc: validation.IntAtLeast(1),
										},
										"retention_period": {
											Type:         pluginsdk.TypeInt,
											Optional:     true,
											ValidateFunc: validation.IntAtLeast(1),
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func monitorPipelineGroupMapValueSchema(required bool) *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Required: required,
		Optional: !required,
		MinItems: minItemsForRequiredList(required),
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"from": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},
				"to": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},
		},
	}
}

func minItemsForRequiredList(required bool) int {
	if required {
		return 1
	}

	return 0
}

func monitorPipelineGroupReceiverSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Required: true,
		MinItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"name": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},
				"type": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringInSlice(pipelinegroups.PossibleValuesForReceiverType(), false),
				},
				"otlp_endpoint": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},
				"syslog": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					MaxItems: 1,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"allow_skip_pri_header": {
								Type:     pluginsdk.TypeBool,
								Optional: true,
							},
							"allowed_formats": {
								Type:     pluginsdk.TypeList,
								Optional: true,
								Elem: &pluginsdk.Schema{
									Type:         pluginsdk.TypeString,
									ValidateFunc: validation.StringInSlice(pipelinegroups.PossibleValuesForAllowedFormats(), false),
								},
							},
							"endpoint": {
								Type:         pluginsdk.TypeString,
								Required:     true,
								ValidateFunc: validation.StringIsNotEmpty,
							},
							"transport_protocol": {
								Type:         pluginsdk.TypeString,
								Optional:     true,
								ValidateFunc: validation.StringInSlice(pipelinegroups.PossibleValuesForTransportProtocol(), false),
							},
						},
					},
				},
				"tls_configuration_name": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},
		},
	}
}

func monitorPipelineGroupProcessorSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"name": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},
				"type": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringInSlice(pipelinegroups.PossibleValuesForProcessorType(), false),
				},
				"batch": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					MaxItems: 1,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"batch_size": {
								Type:         pluginsdk.TypeInt,
								Optional:     true,
								ValidateFunc: validation.IntAtLeast(1),
							},
							"timeout": {
								Type:         pluginsdk.TypeInt,
								Optional:     true,
								ValidateFunc: validation.IntAtLeast(1),
							},
						},
					},
				},
				"transform_statement": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},
		},
	}
}

func monitorPipelineGroupServiceSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Required: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"persistent_volume_name": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},
				"pipeline": {
					Type:     pluginsdk.TypeList,
					Required: true,
					MinItems: 1,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"exporter": {
								Type:     pluginsdk.TypeList,
								Required: true,
								MinItems: 1,
								Elem: &pluginsdk.Schema{
									Type:         pluginsdk.TypeString,
									ValidateFunc: validation.StringIsNotEmpty,
								},
							},
							"name": {
								Type:         pluginsdk.TypeString,
								Required:     true,
								ValidateFunc: validation.StringIsNotEmpty,
							},
							"processor": {
								Type:     pluginsdk.TypeList,
								Optional: true,
								Elem: &pluginsdk.Schema{
									Type:         pluginsdk.TypeString,
									ValidateFunc: validation.StringIsNotEmpty,
								},
							},
							"receiver": {
								Type:     pluginsdk.TypeList,
								Required: true,
								MinItems: 1,
								Elem: &pluginsdk.Schema{
									Type:         pluginsdk.TypeString,
									ValidateFunc: validation.StringIsNotEmpty,
								},
							},
							"type": {
								Type:         pluginsdk.TypeString,
								Required:     true,
								ValidateFunc: validation.StringInSlice(pipelinegroups.PossibleValuesForPipelineType(), false),
							},
						},
					},
				},
			},
		},
	}
}

func monitorPipelineGroupExecutionPlacementSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"constraint": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"capability": {
								Type:         pluginsdk.TypeString,
								Required:     true,
								ValidateFunc: validation.StringIsNotEmpty,
							},
							"operator": {
								Type:         pluginsdk.TypeString,
								Required:     true,
								ValidateFunc: validation.StringInSlice(pipelinegroups.PossibleValuesForCapabilityOperator(), false),
							},
							"values": {
								Type:     pluginsdk.TypeList,
								Optional: true,
								Elem: &pluginsdk.Schema{
									Type:         pluginsdk.TypeString,
									ValidateFunc: validation.StringIsNotEmpty,
								},
							},
						},
					},
				},
				"maximum_instances_per_host": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					ValidateFunc: validation.IntAtLeast(1),
				},
			},
		},
	}
}

func monitorPipelineGroupExtendedLocationSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"name": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},
				"type": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringInSlice(pipelinegroups.PossibleValuesForExtendedLocationType(), false),
				},
			},
		},
	}
}

func monitorPipelineGroupTLSConfigurationSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"name": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},
				"mode": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringInSlice(pipelinegroups.PossibleValuesForTlsMode(), false),
				},
				"client_ca": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					MaxItems: 1,
					Elem: &pluginsdk.Resource{
						Schema: monitorPipelineGroupCertificateSourceSchema(),
					},
				},
				"tls_certificate": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					MaxItems: 1,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"certificate": {
								Type:     pluginsdk.TypeList,
								Required: true,
								MaxItems: 1,
								Elem: &pluginsdk.Resource{
									Schema: monitorPipelineGroupCertificateSourceSchema(),
								},
							},
							"private_key": {
								Type:     pluginsdk.TypeList,
								Required: true,
								MaxItems: 1,
								Elem: &pluginsdk.Resource{
									Schema: monitorPipelineGroupPrivateKeySourceSchema(),
								},
							},
						},
					},
				},
			},
		},
	}
}

func monitorPipelineGroupCertificateSourceSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"sub_location": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"type": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice(pipelinegroups.PossibleValuesForCertificateSourceType(), false),
		},
	}
}

func monitorPipelineGroupPrivateKeySourceSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"sub_location": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"type": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice(pipelinegroups.PossibleValuesForPrivateKeySourceType(), false),
		},
	}
}
