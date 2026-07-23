// Copyright IBM Corp. 2021, 2025
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/comradesharf/terraform-provider-hostinger/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ datasource.DataSource              = &DataSourceVPSVirtualMachines{}
	_ datasource.DataSourceWithConfigure = &DataSourceVPSVirtualMachines{}
)

func NewDataSourceVPSVirtualMachines() datasource.DataSource {
	return &DataSourceVPSVirtualMachines{}
}

// DataSourceVPSVirtualMachines defines the data source implementation.
type DataSourceVPSVirtualMachines struct {
	client *client.ClientWithResponses
}

type VPSVirtualMachineIPv4Model struct {
	Address types.String `tfsdk:"address"`
	Ptr     types.String `tfsdk:"ptr"`
	Netmask types.String `tfsdk:"netmask"`
}

type VPSVirtualMachineIPv6Model struct {
	Address types.String `tfsdk:"address"`
	Ptr     types.String `tfsdk:"ptr"`
}

type VPSVirtualMachineModel struct {
	Id              types.Int64                  `tfsdk:"id"`
	Hostname        types.String                 `tfsdk:"hostname"`
	State           types.String                 `tfsdk:"state"`
	Cpus            types.Int64                  `tfsdk:"cpus"`
	Memory          types.Int64                  `tfsdk:"memory"`
	Disk            types.Int64                  `tfsdk:"disk"`
	Bandwidth       types.Int64                  `tfsdk:"bandwidth"`
	DataCenterId    types.Int64                  `tfsdk:"data_center_id"`
	FirewallGroupId types.Int64                  `tfsdk:"firewall_group_id"`
	Ipv4            []VPSVirtualMachineIPv4Model `tfsdk:"ipv4"`
	Ipv6            []VPSVirtualMachineIPv6Model `tfsdk:"ipv6"`
	OsName          types.String                 `tfsdk:"os_name"`
	CreatedAt       types.String                 `tfsdk:"created_at"`
	ActionsLock     types.String                 `tfsdk:"actions_lock"`
}

// DataSourceVPSVirtualMachinesModel describes the data source data model.
type DataSourceVPSVirtualMachinesModel struct {
	VirtualMachineId types.Int64              `tfsdk:"virtual_machine_id"`
	VirtualMachines  []VPSVirtualMachineModel `tfsdk:"virtual_machines"`
}

func (d *DataSourceVPSVirtualMachines) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vps_virtual_machines"
}

func (d *DataSourceVPSVirtualMachines) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"virtual_machine_id": schema.Int64Attribute{
				Optional: true,
			},
			"virtual_machines": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":                schema.Int64Attribute{Computed: true},
						"hostname":          schema.StringAttribute{Computed: true},
						"state":             schema.StringAttribute{Computed: true},
						"cpus":              schema.Int64Attribute{Computed: true},
						"memory":            schema.Int64Attribute{Computed: true},
						"disk":              schema.Int64Attribute{Computed: true},
						"bandwidth":         schema.Int64Attribute{Computed: true},
						"data_center_id":    schema.Int64Attribute{Computed: true},
						"firewall_group_id": schema.Int64Attribute{Computed: true},
						"ipv4": schema.ListNestedAttribute{
							Computed: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"address": schema.StringAttribute{Computed: true},
									"ptr":     schema.StringAttribute{Computed: true},
									"netmask": schema.StringAttribute{Computed: true},
								},
							},
						},
						"ipv6": schema.ListNestedAttribute{
							Computed: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"address": schema.StringAttribute{Computed: true},
									"ptr":     schema.StringAttribute{Computed: true},
								},
							},
						},
						"os_name":      schema.StringAttribute{Computed: true},
						"created_at":   schema.StringAttribute{Computed: true},
						"actions_lock": schema.StringAttribute{Computed: true},
					},
				},
			},
		},
	}
}

func (d *DataSourceVPSVirtualMachines) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	c, ok := req.ProviderData.(*client.ClientWithResponses)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.ClientWithResponses, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = c
}

func (d *DataSourceVPSVirtualMachines) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data DataSourceVPSVirtualMachinesModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.VirtualMachineId.IsUnknown() {
		resp.Diagnostics.AddError(
			"Unknown Virtual Machine ID",
			"The 'virtual_machine_id' attribute cannot be unknown.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	var virtualMachines []client.VPSV1VirtualMachineVirtualMachineResource

	if !data.VirtualMachineId.IsNull() {
		vmID := data.VirtualMachineId.ValueInt64()
		ctx = tflog.SetField(ctx, "virtual_machine_id", vmID)

		response, err := d.client.VPSGetVirtualMachineDetailsV1WithResponse(ctx, client.VirtualMachineId(vmID))
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to Read VPS Virtual Machines",
				fmt.Sprintf("Got error: %s", err),
			)
			return
		}
		if response.StatusCode() != http.StatusOK {
			resp.Diagnostics.AddError(
				"Unable to Read VPS Virtual Machines",
				fmt.Sprintf("Unexpected status code: %d", response.StatusCode()),
			)
			return
		}
		if response.JSON200 == nil {
			resp.Diagnostics.AddError(
				"Unable to Read VPS Virtual Machines",
				"Response body is nil",
			)
			return
		}

		virtualMachines = append(virtualMachines, *response.JSON200)
	} else {
		response, err := d.client.VPSGetVirtualMachinesV1WithResponse(ctx)
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to Read VPS Virtual Machines",
				fmt.Sprintf("Got error: %s", err),
			)
			return
		}
		if response.StatusCode() != http.StatusOK {
			resp.Diagnostics.AddError(
				"Unable to Read VPS Virtual Machines",
				fmt.Sprintf("Unexpected status code: %d", response.StatusCode()),
			)
			return
		}
		if response.JSON200 == nil {
			resp.Diagnostics.AddError(
				"Unable to Read VPS Virtual Machines",
				"Response body is nil",
			)
			return
		}

		virtualMachines = *response.JSON200
	}

	data.VirtualMachines = make([]VPSVirtualMachineModel, 0, len(virtualMachines))
	for _, vm := range virtualMachines {
		item, err := mapVPSVirtualMachine(vm)
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to Read VPS Virtual Machines",
				fmt.Sprintf("Unable to parse virtual machine response: %s", err),
			)
			return
		}

		data.VirtualMachines = append(data.VirtualMachines, item)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func mapVPSVirtualMachine(vm client.VPSV1VirtualMachineVirtualMachineResource) (VPSVirtualMachineModel, error) {
	ipv4, err := parseVPSVirtualMachineIPv4(vm.Ipv4)
	if err != nil {
		return VPSVirtualMachineModel{}, err
	}

	ipv6, err := parseVPSVirtualMachineIPv6(vm.Ipv6)
	if err != nil {
		return VPSVirtualMachineModel{}, err
	}

	var osName types.String
	if vm.Template != nil {
		template, err := vm.Template.AsVPSV1TemplateTemplateResource()
		if err != nil {
			return VPSVirtualMachineModel{}, err
		}
		osName = stringValueOrNull(template.Name)
	} else {
		osName = types.StringNull()
	}

	return VPSVirtualMachineModel{
		Id:              types.Int64PointerValue(intPointerToInt64Pointer(vm.Id)),
		Hostname:        stringValueOrNull(vm.Hostname),
		State:           stringValueOrNull(vpsStateToStringPointer(vm.State)),
		Cpus:            types.Int64PointerValue(intPointerToInt64Pointer(vm.Cpus)),
		Memory:          types.Int64PointerValue(intPointerToInt64Pointer(vm.Memory)),
		Disk:            types.Int64PointerValue(intPointerToInt64Pointer(vm.Disk)),
		Bandwidth:       types.Int64PointerValue(intPointerToInt64Pointer(vm.Bandwidth)),
		DataCenterId:    types.Int64PointerValue(intPointerToInt64Pointer(vm.DataCenterId)),
		FirewallGroupId: types.Int64PointerValue(intPointerToInt64Pointer(vm.FirewallGroupId)),
		Ipv4:            ipv4,
		Ipv6:            ipv6,
		OsName:          osName,
		CreatedAt:       rfc3339StringOrNull(vm.CreatedAt),
		ActionsLock:     stringValueOrNull(vpsActionsLockToStringPointer(vm.ActionsLock)),
	}, nil
}

func parseVPSVirtualMachineIPv4(ipv4 *client.VPSV1VirtualMachineVirtualMachineResource_Ipv4) ([]VPSVirtualMachineIPv4Model, error) {
	if ipv4 == nil {
		return nil, nil
	}

	raw, err := ipv4.MarshalJSON()
	if err != nil {
		return nil, err
	}

	if string(raw) == "null" {
		return nil, nil
	}

	var values []struct {
		Address *string `json:"address"`
		Ptr     *string `json:"ptr"`
		Netmask *string `json:"netmask"`
	}
	if err := json.Unmarshal(raw, &values); err != nil {
		return nil, err
	}

	result := make([]VPSVirtualMachineIPv4Model, 0, len(values))
	for _, value := range values {
		result = append(result, VPSVirtualMachineIPv4Model{
			Address: stringValueOrNull(value.Address),
			Ptr:     stringValueOrNull(value.Ptr),
			Netmask: stringValueOrNull(value.Netmask),
		})
	}

	return result, nil
}

func parseVPSVirtualMachineIPv6(ipv6 *client.VPSV1VirtualMachineVirtualMachineResource_Ipv6) ([]VPSVirtualMachineIPv6Model, error) {
	if ipv6 == nil {
		return nil, nil
	}

	raw, err := ipv6.MarshalJSON()
	if err != nil {
		return nil, err
	}

	if string(raw) == "null" {
		return nil, nil
	}

	var values []struct {
		Address *string `json:"address"`
		Ptr     *string `json:"ptr"`
	}
	if err := json.Unmarshal(raw, &values); err != nil {
		return nil, err
	}

	result := make([]VPSVirtualMachineIPv6Model, 0, len(values))
	for _, value := range values {
		result = append(result, VPSVirtualMachineIPv6Model{
			Address: stringValueOrNull(value.Address),
			Ptr:     stringValueOrNull(value.Ptr),
		})
	}

	return result, nil
}

func intPointerToInt64Pointer(value *int) *int64 {
	if value == nil {
		return nil
	}

	v := int64(*value)
	return &v
}

func stringValueOrNull(value *string) types.String {
	if value == nil {
		return types.StringNull()
	}

	return types.StringValue(*value)
}

func rfc3339StringOrNull(value *time.Time) types.String {
	if value == nil {
		return types.StringNull()
	}

	return types.StringValue(value.Format(time.RFC3339))
}

func vpsStateToStringPointer(value *client.VPSV1VirtualMachineVirtualMachineResourceState) *string {
	if value == nil {
		return nil
	}

	s := string(*value)
	return &s
}

func vpsActionsLockToStringPointer(value *client.VPSV1VirtualMachineVirtualMachineResourceActionsLock) *string {
	if value == nil {
		return nil
	}

	s := string(*value)
	return &s
}
