// Copyright IBM Corp. 2021, 2025
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"net/http"

	"github.com/comradesharf/terraform-provider-hostinger/internal/client"
	"github.com/hashicorp/terraform-plugin-framework-nettypes/iptypes"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
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

// DataSourceVPSVirtualMachinesModel describes the data source data model.
type DataSourceVPSVirtualMachinesModel struct {
	VirtualMachines []VPSVirtualMachineModel `tfsdk:"virtual_machines"`
}

func (d *DataSourceVPSVirtualMachines) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vps_virtual_machines"
}

func (d *DataSourceVPSVirtualMachines) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"virtual_machines": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							Computed: true,
						},
						"firewall_group_id": schema.Int64Attribute{
							Computed: true,
						},
						"subscription_id": schema.StringAttribute{
							Computed: true,
						},
						"data_center_id": schema.Int64Attribute{
							Computed: true,
						},
						"plan": schema.StringAttribute{
							Computed: true,
						},
						"hostname": schema.StringAttribute{
							Computed: true,
						},
						"state": schema.StringAttribute{
							Computed: true,
						},
						"actions_lock": schema.StringAttribute{
							Computed: true,
						},
						"cpus": schema.Int64Attribute{
							Computed: true,
						},
						"memory": schema.Int64Attribute{
							Computed: true,
						},
						"disk": schema.Int64Attribute{
							Computed: true,
						},
						"bandwidth": schema.Int64Attribute{
							Computed: true,
						},
						"ns1": schema.StringAttribute{
							Computed:   true,
							CustomType: iptypes.IPAddressType{},
						},
						"ns2": schema.StringAttribute{
							Computed:   true,
							CustomType: iptypes.IPAddressType{},
						},
						"ipv4": schema.ListNestedAttribute{
							Computed: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"id": schema.Int64Attribute{
										Computed: true,
									},
									"address": schema.StringAttribute{
										Computed:   true,
										CustomType: iptypes.IPAddressType{},
									},
									"ptr": schema.StringAttribute{
										Computed: true,
									},
								},
							},
						},
						"ipv6": schema.ListNestedAttribute{
							Computed: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"id": schema.Int64Attribute{
										Computed: true,
									},
									"address": schema.StringAttribute{
										Computed:   true,
										CustomType: iptypes.IPAddressType{},
									},
									"ptr": schema.StringAttribute{
										Computed: true,
									},
								},
							},
						},
						"template": schema.SingleNestedAttribute{
							Computed: true,
							Attributes: map[string]schema.Attribute{
								"id": schema.Int64Attribute{
									Computed: true,
								},
								"name": schema.StringAttribute{
									Computed: true,
								},
								"description": schema.StringAttribute{
									Computed: true,
								},
								"documentation": schema.StringAttribute{
									Computed: true,
								},
							},
						},
						"created_at": schema.StringAttribute{
							Computed:   true,
							CustomType: timetypes.RFC3339Type{},
						},
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

	for _, item := range *response.JSON200 {
		var d VPSVirtualMachineModel
		d.Id = int64Value(item.Id)
		d.FirewallGroupId = int64Value(item.FirewallGroupId)
		d.SubscriptionID = types.StringPointerValue(item.SubscriptionId)
		d.DataCenterID = int64Value(item.DataCenterId)
		d.Plan = types.StringPointerValue(item.Plan)
		d.Hostname = types.StringPointerValue(item.Hostname)
		d.State = types.StringPointerValue((*string)(item.State))
		d.ActionsLock = types.StringPointerValue((*string)(item.ActionsLock))
		d.Cpus = int64Value(item.Cpus)
		d.Memory = int64Value(item.Memory)
		d.Disk = int64Value(item.Disk)
		d.Bandwidth = int64Value(item.Bandwidth)
		d.NS1 = iptypes.NewIPAddressPointerValue(item.Ns1)
		d.NS2 = iptypes.NewIPAddressPointerValue(item.Ns2)

		if item.Ipv4 != nil {
			v, err := item.Ipv4.AsVPSV1IPAddressIPAddressCollection()
			if err == nil {
				for _, ip := range v {
					p := VPSVirtualMachineIPAddressModel{}
					p.ID = int64Value(ip.Id)
					p.Address = iptypes.NewIPAddressPointerValue(ip.Address)
					p.Ptr = types.StringPointerValue(ip.Ptr)

					d.Ipv4 = append(d.Ipv4, p)
				}
			}
		}

		if item.Ipv6 != nil {
			v, err := item.Ipv6.AsVPSV1IPAddressIPAddressCollection()
			if err == nil {
				for _, ip := range v {
					p := VPSVirtualMachineIPAddressModel{}
					p.ID = int64Value(ip.Id)
					p.Address = iptypes.NewIPAddressPointerValue(ip.Address)
					p.Ptr = types.StringPointerValue(ip.Ptr)

					d.Ipv6 = append(d.Ipv6, p)
				}
			}
		}

		if item.Template != nil {
			v, err := item.Template.AsVPSV1TemplateTemplateResource()
			if err == nil {
				p := &VPSVirtualMachineTemplateModel{}
				p.ID = int64Value(v.Id)
				p.Name = types.StringPointerValue(v.Name)
				p.Description = types.StringPointerValue(v.Description)
				p.Documentation = types.StringPointerValue(v.Documentation)

				d.Template = p
			}
		}

		d.CreatedAt = timetypes.NewRFC3339TimePointerValue(item.CreatedAt)

		data.VirtualMachines = append(data.VirtualMachines, d)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
