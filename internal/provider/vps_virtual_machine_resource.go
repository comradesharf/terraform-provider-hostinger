// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/comradesharf/terraform-provider-hostinger/internal/client"
	"github.com/hashicorp/terraform-plugin-framework-nettypes/iptypes"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/identityschema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

var (
	_ resource.Resource                = &VPSVirtualMachineResource{}
	_ resource.ResourceWithConfigure   = &VPSVirtualMachineResource{}
	_ resource.ResourceWithImportState = &VPSVirtualMachineResource{}
	_ resource.ResourceWithIdentity    = &VPSVirtualMachineResource{}
)

func NewVPSVirtualMachineResource() resource.Resource {
	return &VPSVirtualMachineResource{}
}

type VPSVirtualMachineResource struct {
	client *client.ClientWithResponses
}

type VPSVirtualMachineResourceModel struct {
	VPSVirtualMachineModel
	Timeouts timeouts.Value `tfsdk:"timeouts"`
}

func (r *VPSVirtualMachineResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vps_virtual_machine"
}

func (r *VPSVirtualMachineResource) IdentitySchema(ctx context.Context, req resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = identityschema.Schema{
		Attributes: map[string]identityschema.Attribute{
			"id": identityschema.Int64Attribute{
				RequiredForImport: true,
			},
		},
	}
}

func (r *VPSVirtualMachineResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages an existing VPS virtual machine in the Hostinger account. This resource does not support creation, so it must be imported using the `terraform import` or `terraform query` command.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
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
				Required: true,
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
				Required:   true,
				CustomType: iptypes.IPAddressType{},
			},
			"ns2": schema.StringAttribute{
				Required:   true,
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
			"timeouts": timeouts.AttributesAll(ctx),
		},
	}
}

func (r *VPSVirtualMachineResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	c, ok := req.ProviderData.(*client.ClientWithResponses)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.ClientWithResponses, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = c
}

func (r *VPSVirtualMachineResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	resp.Diagnostics.AddError(
		"Create not yet implemented",
		"This resource does not support creation yet. Please use the import command to import an existing VPS virtual machine into Terraform management.",
	)
}

func (r *VPSVirtualMachineResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state VPSVirtualMachineResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var identity VPSVirtualMachineIdentity
	resp.Diagnostics.Append(req.Identity.Get(ctx, &identity)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.ID.IsUnknown() || state.ID.IsNull() || state.ID.ValueInt64() == 0 {
		resp.Diagnostics.AddError("Invalid VPS Virtual Machine ID", "The VPS virtual machine ID is unknown, null, or zero, so it cannot be read.")
		return
	}

	readTimeout, diags := state.Timeouts.Read(ctx, 20*time.Minute)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, readTimeout)
	defer cancel()

	response, err := r.client.VPSGetVirtualMachineDetailsV1WithResponse(ctx, client.VirtualMachineId(state.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError("Unable to Read VPS Virtual Machine", fmt.Sprintf("Got error: %s", err))
		return
	}
	if response.StatusCode() == http.StatusNotFound {
		resp.State.RemoveResource(ctx)
		return
	}
	if response.StatusCode() != http.StatusOK {
		resp.Diagnostics.AddError("Unable to Read VPS Virtual Machine", fmt.Sprintf("Unexpected status code: %d, response: %s", response.StatusCode(), string(response.Body)))
		return
	}
	if response.JSON200 == nil {
		resp.Diagnostics.AddError("Unable to Read VPS Virtual Machine", "Response body is nil")
		return
	}

	state.Merge(*response.JSON200)

	identity.ID = state.ID

	resp.Diagnostics.Append(resp.Identity.Set(ctx, &identity)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *VPSVirtualMachineResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state, config VPSVirtualMachineResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.ID.IsUnknown() || state.ID.IsNull() || state.ID.ValueInt64() == 0 {
		resp.Diagnostics.AddError("Invalid VPS Virtual Machine ID", "The VPS virtual machine ID is unknown, null, or zero, so it cannot be updated.")
		return
	}

	updateTimeout, diags := state.Timeouts.Update(ctx, 20*time.Minute)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, updateTimeout)
	defer cancel()

	var wg sync.WaitGroup

	if !state.Hostname.Equal(config.Hostname) {
		wg.Go(func() {
			param := client.VPSSetHostnameV1JSONRequestBody{
				Hostname: config.Hostname.ValueString(),
			}
			response, err := r.client.VPSSetHostnameV1WithResponse(ctx, client.VirtualMachineId(state.ID.ValueInt64()), param)
			if err != nil {
				resp.Diagnostics.AddError("Unable to Update VPS Virtual Machine Hostname", fmt.Sprintf("Got error: %s", err))
				return
			}
			if response.StatusCode() != http.StatusOK {
				resp.Diagnostics.AddError("Unable to Update VPS Virtual Machine Hostname", fmt.Sprintf("Unexpected status code: %d, response: %s", response.StatusCode(), string(response.Body)))
				return
			}

			resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("hostname"), config.Hostname)...)
		})
	}

	if !state.NS1.Equal(config.NS1) || !state.NS2.Equal(config.NS2) {
		wg.Go(func() {
			param := client.VPSSetNameserversV1JSONRequestBody{
				Ns1: config.NS1.ValueString(),
				Ns2: config.NS2.ValueStringPointer(),
			}
			response, err := r.client.VPSSetNameserversV1WithResponse(ctx, client.VirtualMachineId(state.ID.ValueInt64()), param)
			if err != nil {
				resp.Diagnostics.AddError("Unable to Update VPS Virtual Machine Nameservers", fmt.Sprintf("Got error: %s", err))
				return
			}
			if response.StatusCode() != http.StatusOK {
				resp.Diagnostics.AddError("Unable to Update VPS Virtual Machine Nameservers", fmt.Sprintf("Unexpected status code: %d, response: %s", response.StatusCode(), string(response.Body)))
				return
			}

			resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("ns1"), config.NS1)...)
			if resp.Diagnostics.HasError() {
				return
			}

			resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("ns2"), config.NS2)...)
		})
	}

	wg.Wait()

	if resp.Diagnostics.HasError() {
		return
	}

	var identity VPSVirtualMachineIdentity
	resp.Diagnostics.Append(req.Identity.Get(ctx, &identity)...)
}

// Delete disables auto-renewal for the VPS virtual machine subscription, effectively scheduling it for deletion at the end of the current billing cycle.
func (r *VPSVirtualMachineResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state VPSVirtualMachineResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.ID.IsUnknown() || state.ID.IsNull() || state.ID.ValueInt64() == 0 {
		resp.Diagnostics.AddError("Invalid VPS Virtual Machine ID", "The VPS virtual machine ID is unknown, null, or zero, so it cannot be deleted.")
		return
	}

	deleteTimeout, diags := state.Timeouts.Delete(ctx, 20*time.Minute)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, deleteTimeout)
	defer cancel()

	response, err := r.client.BillingDisableAutoRenewalV1WithResponse(ctx, state.SubscriptionID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Unable to Disable Auto Renewal for VPS Virtual Machine", fmt.Sprintf("Got error: %s", err))
		return
	}
	if response.StatusCode() != http.StatusOK && response.StatusCode() != http.StatusNotFound {
		resp.Diagnostics.AddError("Unable to Disable Auto Renewal for VPS Virtual Machine", fmt.Sprintf("Unexpected status code: %d, response: %s", response.StatusCode(), string(response.Body)))
	}
}

func (r *VPSVirtualMachineResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importStatePassthroughInt64Identity(ctx, path.Root("id"), path.Root("id"), req, resp)
}
