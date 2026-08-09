// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/comradesharf/terraform-provider-hostinger/internal/client"
	"github.com/hashicorp/terraform-plugin-framework-nettypes/iptypes"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource              = &VPSVirtualMachineNameserversResource{}
	_ resource.ResourceWithConfigure = &VPSVirtualMachineNameserversResource{}
)

func NewVPSVirtualMachineNameserversResource() resource.Resource {
	return &VPSVirtualMachineNameserversResource{}
}

type VPSVirtualMachineNameserversResource struct {
	client *client.ClientWithResponses
}

type VPSVirtualMachineNameserversResourceModel struct {
	VPSVirtualMachineNameserversModel
	VirtualMachineID types.Int64    `tfsdk:"virtual_machine_id"`
	WaitForAction    types.Bool     `tfsdk:"wait_for_action"`
	Timeouts         timeouts.Value `tfsdk:"timeouts"`
}

func (r *VPSVirtualMachineNameserversResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vps_virtual_machine_nameservers"
}

func (r *VPSVirtualMachineNameserversResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages the nameservers of a VPS virtual machine in the Hostinger account.",
		Attributes: map[string]schema.Attribute{
			"virtual_machine_id": schema.Int64Attribute{
				MarkdownDescription: "The ID of the VPS virtual machine.",
				Required:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"ns1": schema.StringAttribute{
				Required:            true,
				CustomType:          iptypes.IPAddressType{},
				MarkdownDescription: "The primary nameserver IP address for the VPS virtual machine.",
			},
			"ns2": schema.StringAttribute{
				Optional:            true,
				CustomType:          iptypes.IPAddressType{},
				MarkdownDescription: "The secondary nameserver IP address for the VPS virtual machine.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(7),
				},
			},
			"wait_for_action": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Whether to wait for the action to complete before returning. Defaults to true.",
			},
			"timeouts": timeouts.AttributesAll(ctx),
		},
	}
}

func (r *VPSVirtualMachineNameserversResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *VPSVirtualMachineNameserversResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan VPSVirtualMachineNameserversResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createTimeout, diags := plan.Timeouts.Create(ctx, 20*time.Minute)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, createTimeout)
	defer cancel()

	params := client.VPSSetNameserversV1JSONRequestBody{
		Ns1: plan.NS1.ValueString(),
		Ns2: plan.NS2.ValueStringPointer(),
	}

	response, err := r.client.VPSSetNameserversV1WithResponse(ctx, client.VirtualMachineId(plan.VirtualMachineID.ValueInt64()), params)
	if err != nil {
		resp.Diagnostics.AddError("Unable to Set VPS Nameservers", fmt.Sprintf("Got error: %s", err))
		return
	}
	if response.StatusCode() != http.StatusOK {
		resp.Diagnostics.AddError("Unable to Set VPS Nameservers", fmt.Sprintf("Unexpected status code: %d, response: %s", response.StatusCode(), string(response.Body)))
		return
	}
	if response.JSON200 == nil {
		resp.Diagnostics.AddError("Unable to Set VPS Nameservers", "Response body is nil")
		return
	}
	if response.JSON200.Id == nil {
		resp.Diagnostics.AddError("Unable to Set VPS Nameservers", "Response body does not contain an action ID")
		return
	}

	if plan.WaitForAction.IsNull() || plan.WaitForAction.IsUnknown() || plan.WaitForAction.ValueBool() {
	poll:
		for {
			response, err := r.client.VPSGetActionDetailsV1WithResponse(
				ctx,
				client.VirtualMachineId(plan.VirtualMachineID.ValueInt64()),
				*response.JSON200.Id,
			)
			if err != nil {
				resp.Diagnostics.AddError("Unable to Get VPS Actions", fmt.Sprintf("Got error: %s", err))
				return
			}
			if response.StatusCode() != http.StatusOK {
				resp.Diagnostics.AddError("Unable to Get VPS Actions", fmt.Sprintf("Unexpected status code: %d, response: %s", response.StatusCode(), string(response.Body)))
				return
			}
			if response.JSON200 == nil {
				resp.Diagnostics.AddError("Unable to Get VPS Actions", "Response body is nil")
				return
			}
			if response.JSON200.State == nil {
				resp.Diagnostics.AddError("Unable to Get VPS Actions", "Response body does not contain an action state")
				return
			}

			switch *response.JSON200.State {
			case client.VPSV1ActionActionResourceStateSuccess:
				break poll
			case client.VPSV1ActionActionResourceStateError:
				resp.Diagnostics.AddError("Unable to Set VPS Nameservers", "The action to set the VPS nameservers failed")
				break poll
			default:
				select {
				case <-ctx.Done():
					resp.Diagnostics.AddError("Unable to Set VPS Nameservers", "The action to set the VPS nameservers timed out")
					return
				case <-time.After(2 * time.Second):
					// continue polling
				}
			}
		}

		if resp.Diagnostics.HasError() {
			return
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *VPSVirtualMachineNameserversResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state VPSVirtualMachineNameserversResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Nameserver values are refreshed from the API; only the virtual machine ID is required here.
	if state.VirtualMachineID.IsUnknown() || state.VirtualMachineID.IsNull() || state.VirtualMachineID.ValueInt64() == 0 {
		resp.Diagnostics.AddError("Invalid VPS Virtual Machine ID", "The VPS virtual machine ID is unknown, null, or zero, so it cannot be read.")
	}

	if resp.Diagnostics.HasError() {
		return
	}

	readTimeout, diags := state.Timeouts.Read(ctx, 20*time.Minute)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, readTimeout)
	defer cancel()

	response, err := r.client.VPSGetVirtualMachineDetailsV1WithResponse(ctx, client.VirtualMachineId(state.VirtualMachineID.ValueInt64()))
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

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *VPSVirtualMachineNameserversResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state, plan VPSVirtualMachineNameserversResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.VirtualMachineID.IsUnknown() || plan.VirtualMachineID.IsNull() || plan.VirtualMachineID.ValueInt64() == 0 {
		resp.Diagnostics.AddError("Invalid VPS Virtual Machine ID", "The VPS virtual machine ID is unknown, null, or zero, so it cannot be updated.")
		return
	}

	if plan.NS2.Equal(state.NS2) && plan.NS1.Equal(state.NS1) {
		resp.Diagnostics.AddWarning("No Changes Detected", "The nameservers are already set to the specified values. No update is necessary.")
		resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
		return
	}

	updateTimeout, diags := plan.Timeouts.Update(ctx, 20*time.Minute)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, updateTimeout)
	defer cancel()

	param := client.VPSSetNameserversV1JSONRequestBody{
		Ns1: plan.NS1.ValueString(),
		Ns2: plan.NS2.ValueStringPointer(),
	}
	response, err := r.client.VPSSetNameserversV1WithResponse(ctx, client.VirtualMachineId(state.VirtualMachineID.ValueInt64()), param)
	if err != nil {
		resp.Diagnostics.AddError("Unable to Update VPS Virtual Machine Nameservers", fmt.Sprintf("Got error: %s", err))
		return
	}
	if response.StatusCode() == http.StatusNotFound {
		resp.State.RemoveResource(ctx)
		return
	}
	if response.StatusCode() != http.StatusOK {
		resp.Diagnostics.AddError("Unable to Update VPS Virtual Machine Nameservers", fmt.Sprintf("Unexpected status code: %d, response: %s", response.StatusCode(), string(response.Body)))
		return
	}
	if response.JSON200 == nil {
		resp.Diagnostics.AddError("Unable to Update VPS Virtual Machine Nameservers", "Response body is nil")
		return
	}
	if response.JSON200.Id == nil {
		resp.Diagnostics.AddError("Unable to Update VPS Virtual Machine Nameservers", "Response body does not contain an action ID")
		return
	}

	if plan.WaitForAction.IsNull() || plan.WaitForAction.IsUnknown() || plan.WaitForAction.ValueBool() {
	poll:
		for {
			response, err := r.client.VPSGetActionDetailsV1WithResponse(
				ctx,
				client.VirtualMachineId(state.VirtualMachineID.ValueInt64()),
				*response.JSON200.Id,
			)
			if err != nil {
				resp.Diagnostics.AddError("Unable to Get VPS Actions", fmt.Sprintf("Got error: %s", err))
				return
			}
			if response.StatusCode() != http.StatusOK {
				resp.Diagnostics.AddError("Unable to Get VPS Actions", fmt.Sprintf("Unexpected status code: %d, response: %s", response.StatusCode(), string(response.Body)))
				return
			}
			if response.JSON200 == nil {
				resp.Diagnostics.AddError("Unable to Get VPS Actions", "Response body is nil")
				return
			}
			if response.JSON200.State == nil {
				resp.Diagnostics.AddError("Unable to Get VPS Actions", "Response body does not contain an action state")
				return
			}

			switch *response.JSON200.State {
			case "success":
				break poll
			case "error":
				resp.Diagnostics.AddError("Unable to Set VPS Nameservers", "The action to set the VPS nameservers failed")
				break poll
			default:
				select {
				case <-ctx.Done():
					resp.Diagnostics.AddError("Unable to Update VPS Virtual Machine Nameservers", "The action to set the VPS nameservers timed out")
					return
				case <-time.After(2 * time.Second):
					// continue polling
				}
			}
		}

		if resp.Diagnostics.HasError() {
			return
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete resets the nameservers for the VPS virtual machine.
func (r *VPSVirtualMachineNameserversResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state VPSVirtualMachineNameserversResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.VirtualMachineID.IsUnknown() || state.VirtualMachineID.IsNull() || state.VirtualMachineID.ValueInt64() == 0 {
		resp.Diagnostics.AddError("Invalid VPS Virtual Machine ID", "The VPS virtual machine ID is unknown, null, or zero, so it cannot be updated.")
		return
	}

	deleteTimeout, diags := state.Timeouts.Delete(ctx, 20*time.Minute)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, deleteTimeout)
	defer cancel()

	ns2 := "1.1.1.1"
	param := client.VPSSetNameserversV1JSONRequestBody{
		Ns1: "153.92.2.6",
		Ns2: &ns2,
	}

	response, err := r.client.VPSSetNameserversV1WithResponse(ctx, client.VirtualMachineId(state.VirtualMachineID.ValueInt64()), param)
	if err != nil {
		resp.Diagnostics.AddError("Unable to reset VPS Virtual Machine Nameservers", fmt.Sprintf("Got error: %s", err))
		return
	}
	if response.StatusCode() == http.StatusNotFound {
		return
	}
	if response.StatusCode() != http.StatusOK {
		resp.Diagnostics.AddError("Unable to reset VPS Virtual Machine Nameservers", fmt.Sprintf("Unexpected status code: %d, response: %s", response.StatusCode(), string(response.Body)))
		return
	}
	if response.JSON200 == nil {
		resp.Diagnostics.AddError("Unable to reset VPS Virtual Machine Nameservers", "Response body is nil")
		return
	}
	if response.JSON200.Id == nil {
		resp.Diagnostics.AddError("Unable to reset VPS Virtual Machine Nameservers", "Response body does not contain an action ID")
		return
	}

	if state.WaitForAction.IsNull() || state.WaitForAction.IsUnknown() || state.WaitForAction.ValueBool() {

	poll:
		for {
			response, err := r.client.VPSGetActionDetailsV1WithResponse(
				ctx,
				client.VirtualMachineId(state.VirtualMachineID.ValueInt64()),
				*response.JSON200.Id,
			)
			if err != nil {
				resp.Diagnostics.AddError("Unable to Get VPS Actions", fmt.Sprintf("Got error: %s", err))
				return
			}
			if response.StatusCode() != http.StatusOK {
				resp.Diagnostics.AddError("Unable to Get VPS Actions", fmt.Sprintf("Unexpected status code: %d, response: %s", response.StatusCode(), string(response.Body)))
				return
			}
			if response.JSON200 == nil {
				resp.Diagnostics.AddError("Unable to Get VPS Actions", "Response body is nil")
				return
			}
			if response.JSON200.State == nil {
				resp.Diagnostics.AddError("Unable to Get VPS Actions", "Response body does not contain an action state")
				return
			}

			switch *response.JSON200.State {
			case "success":
				break poll
			case "error":
				resp.Diagnostics.AddError("Unable to Reset VPS Nameservers", "The action to reset the VPS nameservers failed")
				break poll
			default:
				select {
				case <-ctx.Done():
					resp.Diagnostics.AddError("Unable to Reset VPS Nameservers", "The action to reset the VPS nameservers timed out")
					return
				case <-time.After(2 * time.Second):
					// continue polling
				}
			}
		}

		if resp.Diagnostics.HasError() {
			return
		}
	}
}
