// Copyright IBM Corp. 2021, 2025
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/comradesharf/terraform-provider-hostinger/internal/client"
	"github.com/hashicorp/terraform-plugin-framework-nettypes/iptypes"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ datasource.DataSource              = &VPSVirtualMachinesDataSource{}
	_ datasource.DataSourceWithConfigure = &VPSVirtualMachinesDataSource{}
)

func NewVPSVirtualMachinesDataSource() datasource.DataSource {
	return &VPSVirtualMachinesDataSource{}
}

// VPSVirtualMachinesDataSource defines the data source implementation.
type VPSVirtualMachinesDataSource struct {
	client *client.ClientWithResponses
}

// VPSVirtualMachinesDataSourceModel describes the data source data model.
type VPSVirtualMachinesDataSourceModel struct {
	VirtualMachines []VPSVirtualMachineModel `tfsdk:"virtual_machines"`
	Timeouts        timeouts.Value           `tfsdk:"timeouts"`
}

func (d *VPSVirtualMachinesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vps_virtual_machines"
}

func (d *VPSVirtualMachinesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
			"timeouts": timeouts.Attributes(ctx),
		},
	}
}

func (d *VPSVirtualMachinesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *VPSVirtualMachinesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data VPSVirtualMachinesDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	readTimeout, diags := data.Timeouts.Read(ctx, 20*time.Minute)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, readTimeout)
	defer cancel()

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
			fmt.Sprintf("Unexpected status code: %d, response: %s", response.StatusCode(), string(response.Body)),
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
		d.Merge(item)
		data.VirtualMachines = append(data.VirtualMachines, d)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
