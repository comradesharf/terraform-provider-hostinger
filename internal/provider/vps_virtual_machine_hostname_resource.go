// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/comradesharf/terraform-provider-hostinger/internal/client"
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
	_ resource.Resource              = &VPSVirtualMachineHostnameResource{}
	_ resource.ResourceWithConfigure = &VPSVirtualMachineHostnameResource{}
)

func NewVPSVirtualMachineHostnameResource() resource.Resource {
	return &VPSVirtualMachineHostnameResource{}
}

type VPSVirtualMachineHostnameResource struct {
	client *client.ClientWithResponses
}

type VPSVirtualMachineHostnameResourceModel struct {
	VPSVirtualMachineHostnameModel
	VirtualMachineID types.Int64    `tfsdk:"virtual_machine_id"`
	WaitForAction    types.Bool     `tfsdk:"wait_for_action"`
	Timeouts         timeouts.Value `tfsdk:"timeouts"`
}

func (r *VPSVirtualMachineHostnameResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vps_virtual_machine_hostname"
}

func (r *VPSVirtualMachineHostnameResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages the hostname of a VPS virtual machine in the Hostinger account.",
		Attributes: map[string]schema.Attribute{
			"virtual_machine_id": schema.Int64Attribute{
				MarkdownDescription: "The ID of the VPS virtual machine.",
				Required:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"hostname": schema.StringAttribute{
				MarkdownDescription: "The hostname of the VPS virtual machine.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"wait_for_action": schema.BoolAttribute{
				MarkdownDescription: "Whether to wait for the action to complete before returning. Defaults to true.",
				Optional:            true,
			},
			"timeouts": timeouts.AttributesAll(ctx),
		},
	}
}

func (r *VPSVirtualMachineHostnameResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *VPSVirtualMachineHostnameResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan VPSVirtualMachineHostnameResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if plan.VirtualMachineID.IsUnknown() || plan.VirtualMachineID.IsNull() || plan.VirtualMachineID.ValueInt64() <= 0 {
		resp.Diagnostics.AddError("Invalid VPS Virtual Machine ID", "The VPS virtual machine ID is unknown, null, or zero, so it cannot be updated.")
		return
	}

	timeout, diags := plan.Timeouts.Create(ctx, 20*time.Minute)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	response, err := r.client.VPSSetHostnameV1WithResponse(ctx, client.VirtualMachineId(plan.VirtualMachineID.ValueInt64()), client.VPSSetHostnameV1JSONRequestBody{
		Hostname: plan.Hostname.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Unable to Set VPS Hostname", fmt.Sprintf("Got error: %s", err))
		return
	}
	if response.StatusCode() != http.StatusOK || response.JSON200 == nil || response.JSON200.Id == nil || *response.JSON200.Id <= 0 {
		resp.Diagnostics.AddError("Unable to Set VPS Hostname", "The response did not contain a valid action ID.")
		return
	}
	if plan.WaitForAction.IsNull() || plan.WaitForAction.IsUnknown() || plan.WaitForAction.ValueBool() {
	poll:
		for {
			actionResponse, actionErr := r.client.VPSGetActionDetailsV1WithResponse(ctx, client.VirtualMachineId(plan.VirtualMachineID.ValueInt64()), *response.JSON200.Id)
			if actionErr != nil {
				resp.Diagnostics.AddError("Unable to Set VPS Hostname", fmt.Sprintf("Unable to get VPS action: %s", actionErr))
				return
			}
			if actionResponse.StatusCode() != http.StatusOK {
				resp.Diagnostics.AddError("Unable to Set VPS Hostname", fmt.Sprintf("Unexpected action status code: %d, response: %s", actionResponse.StatusCode(), string(actionResponse.Body)))
				return
			}
			if actionResponse.JSON200 == nil || actionResponse.JSON200.State == nil {
				resp.Diagnostics.AddError("Unable to Set VPS Hostname", "Action response does not contain an action state")
				return
			}
			switch *actionResponse.JSON200.State {
			case client.VPSV1ActionActionResourceStateSuccess:
				break poll
			case client.VPSV1ActionActionResourceStateError:
				resp.Diagnostics.AddError("Unable to Set VPS Hostname", "The VPS hostname action failed")
				return
			default:
				select {
				case <-ctx.Done():
					resp.Diagnostics.AddError("Unable to Set VPS Hostname", fmt.Sprintf("The VPS hostname action timed out: %s", ctx.Err()))
					return
				case <-time.After(2 * time.Second):
				}
			}
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *VPSVirtualMachineHostnameResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state VPSVirtualMachineHostnameResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if state.VirtualMachineID.IsUnknown() || state.VirtualMachineID.IsNull() || state.VirtualMachineID.ValueInt64() <= 0 {
		resp.Diagnostics.AddError("Invalid VPS Virtual Machine ID", "The VPS virtual machine ID is unknown, null, or zero, so it cannot be read.")
		return
	}

	timeout, diags := state.Timeouts.Read(ctx, 20*time.Minute)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	response, err := r.client.VPSGetVirtualMachineDetailsV1WithResponse(ctx, client.VirtualMachineId(state.VirtualMachineID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError("Unable to Read VPS Virtual Machine", fmt.Sprintf("Got error: %s", err))
		return
	}
	if response.StatusCode() == http.StatusUnprocessableEntity || response.StatusCode() == http.StatusNotFound {
		resp.State.RemoveResource(ctx)
		return
	}
	if response.StatusCode() != http.StatusOK || response.JSON200 == nil {
		resp.Diagnostics.AddError("Unable to Read VPS Virtual Machine", fmt.Sprintf("Unexpected status code: %d, response body is nil or invalid.", response.StatusCode()))
		return
	}

	state.Merge(*response.JSON200)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *VPSVirtualMachineHostnameResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state VPSVirtualMachineHostnameResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if state.VirtualMachineID.IsUnknown() || state.VirtualMachineID.IsNull() || state.VirtualMachineID.ValueInt64() <= 0 {
		resp.Diagnostics.AddError("Invalid VPS Virtual Machine ID", "The VPS virtual machine ID is unknown, null, or zero, so it cannot be updated.")
		return
	}

	timeout, diags := plan.Timeouts.Update(ctx, 20*time.Minute)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	response, err := r.client.VPSSetHostnameV1WithResponse(ctx, client.VirtualMachineId(state.VirtualMachineID.ValueInt64()), client.VPSSetHostnameV1JSONRequestBody{
		Hostname: plan.Hostname.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Unable to Update VPS Hostname", fmt.Sprintf("Got error: %s", err))
		return
	}
	if response.StatusCode() != http.StatusOK || response.JSON200 == nil || response.JSON200.Id == nil || *response.JSON200.Id <= 0 {
		resp.Diagnostics.AddError("Unable to Update VPS Hostname", "The response did not contain a valid action ID.")
		return
	}
	if plan.WaitForAction.IsNull() || plan.WaitForAction.IsUnknown() || plan.WaitForAction.ValueBool() {
	poll:
		for {
			actionResponse, actionErr := r.client.VPSGetActionDetailsV1WithResponse(ctx, client.VirtualMachineId(state.VirtualMachineID.ValueInt64()), *response.JSON200.Id)
			if actionErr != nil {
				resp.Diagnostics.AddError("Unable to Update VPS Hostname", fmt.Sprintf("Unable to get VPS action: %s", actionErr))
				return
			}
			if actionResponse.StatusCode() != http.StatusOK {
				resp.Diagnostics.AddError("Unable to Update VPS Hostname", fmt.Sprintf("Unexpected action status code: %d, response: %s", actionResponse.StatusCode(), string(actionResponse.Body)))
				return
			}
			if actionResponse.JSON200 == nil || actionResponse.JSON200.State == nil {
				resp.Diagnostics.AddError("Unable to Update VPS Hostname", "Action response does not contain an action state")
				return
			}
			switch *actionResponse.JSON200.State {
			case client.VPSV1ActionActionResourceStateSuccess:
				break poll
			case client.VPSV1ActionActionResourceStateError:
				resp.Diagnostics.AddError("Unable to Update VPS Hostname", "The VPS hostname action failed")
				return
			default:
				select {
				case <-ctx.Done():
					resp.Diagnostics.AddError("Unable to Update VPS Hostname", fmt.Sprintf("The VPS hostname action timed out: %s", ctx.Err()))
					return
				case <-time.After(2 * time.Second):
				}
			}
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *VPSVirtualMachineHostnameResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state VPSVirtualMachineHostnameResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if state.VirtualMachineID.IsUnknown() || state.VirtualMachineID.IsNull() || state.VirtualMachineID.ValueInt64() <= 0 {
		resp.Diagnostics.AddError("Invalid VPS Virtual Machine ID", "The VPS virtual machine ID is unknown, null, or zero, so it cannot be deleted.")
		return
	}

	timeout, diags := state.Timeouts.Delete(ctx, 20*time.Minute)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	response, err := r.client.VPSResetHostnameV1WithResponse(ctx, client.VirtualMachineId(state.VirtualMachineID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError("Unable to Reset VPS Hostname", fmt.Sprintf("Got error: %s", err))
		return
	}
	if response.StatusCode() == http.StatusNotFound {
		return
	}
	if response.StatusCode() != http.StatusOK || response.JSON200 == nil || response.JSON200.Id == nil || *response.JSON200.Id <= 0 {
		resp.Diagnostics.AddError("Unable to Reset VPS Hostname", "The response did not contain a valid action ID.")
		return
	}
	if state.WaitForAction.IsNull() || state.WaitForAction.IsUnknown() || state.WaitForAction.ValueBool() {
		for {
			actionResponse, actionErr := r.client.VPSGetActionDetailsV1WithResponse(ctx, client.VirtualMachineId(state.VirtualMachineID.ValueInt64()), *response.JSON200.Id)
			if actionErr != nil {
				resp.Diagnostics.AddError("Unable to Reset VPS Hostname", fmt.Sprintf("Unable to get VPS action: %s", actionErr))
				return
			}
			if actionResponse.StatusCode() != http.StatusOK {
				resp.Diagnostics.AddError("Unable to Reset VPS Hostname", fmt.Sprintf("Unexpected action status code: %d, response: %s", actionResponse.StatusCode(), string(actionResponse.Body)))
				return
			}
			if actionResponse.JSON200 == nil || actionResponse.JSON200.State == nil {
				resp.Diagnostics.AddError("Unable to Reset VPS Hostname", "Action response does not contain an action state")
				return
			}
			switch *actionResponse.JSON200.State {
			case client.VPSV1ActionActionResourceStateSuccess:
				return
			case client.VPSV1ActionActionResourceStateError:
				resp.Diagnostics.AddError("Unable to Reset VPS Hostname", "The VPS hostname action failed")
				return
			default:
				select {
				case <-ctx.Done():
					resp.Diagnostics.AddError("Unable to Reset VPS Hostname", fmt.Sprintf("The VPS hostname action timed out: %s", ctx.Err()))
					return
				case <-time.After(2 * time.Second):
				}
			}
		}
	}
}
