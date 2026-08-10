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
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/identityschema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource              = &VPSMonarxResource{}
	_ resource.ResourceWithConfigure = &VPSMonarxResource{}
	_ resource.ResourceWithIdentity  = &VPSMonarxResource{}
)

func NewVPSMonarxResource() resource.Resource {
	return &VPSMonarxResource{}
}

type VPSMonarxResource struct {
	client *client.ClientWithResponses
}

type VPSMonarxResourceModel struct {
	VirtualMachineID types.Int64    `tfsdk:"virtual_machine_id"`
	WaitForAction    types.Bool     `tfsdk:"wait_for_action"`
	Timeouts         timeouts.Value `tfsdk:"timeouts"`
}

func (r *VPSMonarxResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vps_monarx"
}

func (r *VPSMonarxResource) IdentitySchema(ctx context.Context, req resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = identityschema.Schema{
		Attributes: map[string]identityschema.Attribute{
			"virtual_machine_id": identityschema.Int64Attribute{
				RequiredForImport: true,
				Description:       "The ID of the VPS virtual machine.",
			},
		},
	}
}

func (r *VPSMonarxResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages the Monarx malware scanner on a VPS virtual machine.",
		Attributes: map[string]schema.Attribute{
			"virtual_machine_id": schema.Int64Attribute{
				MarkdownDescription: "The ID of the VPS virtual machine.",
				Required:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"wait_for_action": schema.BoolAttribute{
				MarkdownDescription: "Whether to wait for the Monarx action to complete before returning. Defaults to true.",
				Optional:            true,
			},
			"timeouts": timeouts.AttributesAll(ctx),
		},
	}
}

func (r *VPSMonarxResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *VPSMonarxResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan VPSMonarxResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if plan.VirtualMachineID.IsUnknown() || plan.VirtualMachineID.IsNull() || plan.VirtualMachineID.ValueInt64() <= 0 {
		resp.Diagnostics.AddError("Invalid VPS Virtual Machine ID", "The VPS virtual machine ID is unknown, null, or zero, so Monarx cannot be installed.")
		return
	}

	timeout, diags := plan.Timeouts.Create(ctx, 20*time.Minute)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	response, err := r.client.VPSInstallMonarxV1WithResponse(ctx, client.VirtualMachineId(plan.VirtualMachineID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError("Unable to Install Monarx", fmt.Sprintf("Got error: %s", err))
		return
	}
	if response.StatusCode() != http.StatusOK || response.JSON200 == nil || response.JSON200.Id == nil || *response.JSON200.Id <= 0 {
		resp.Diagnostics.AddError("Unable to Install Monarx", "The response did not contain a valid action ID.")
		return
	}

	if plan.WaitForAction.IsNull() || plan.WaitForAction.IsUnknown() || plan.WaitForAction.ValueBool() {
		actionID := *response.JSON200.Id
		for {
			actionResponse, actionErr := r.client.VPSGetActionDetailsV1WithResponse(ctx, client.VirtualMachineId(plan.VirtualMachineID.ValueInt64()), actionID)
			if actionErr != nil {
				resp.Diagnostics.AddError("Unable to Install Monarx", fmt.Sprintf("Unable to get Monarx action: %s", actionErr))
				return
			}
			if actionResponse.StatusCode() != http.StatusOK || actionResponse.JSON200 == nil || actionResponse.JSON200.State == nil {
				resp.Diagnostics.AddError("Unable to Install Monarx", fmt.Sprintf("Unexpected action response: status code %d", actionResponse.StatusCode()))
				return
			}
			switch *actionResponse.JSON200.State {
			case client.VPSV1ActionActionResourceStateSuccess:
				goto installed
			case client.VPSV1ActionActionResourceStateError:
				resp.Diagnostics.AddError("Unable to Install Monarx", "The Monarx installation action failed")
				return
			default:
				select {
				case <-ctx.Done():
					resp.Diagnostics.AddError("Unable to Install Monarx", fmt.Sprintf("The Monarx installation action timed out: %s", ctx.Err()))
					return
				case <-time.After(2 * time.Second):
				}
			}
		}
	}

installed:
	identity := VPSMonarxIdentity{VirtualMachineID: plan.VirtualMachineID}
	resp.Diagnostics.Append(resp.Identity.Set(ctx, &identity)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *VPSMonarxResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state VPSMonarxResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if state.VirtualMachineID.IsUnknown() || state.VirtualMachineID.IsNull() || state.VirtualMachineID.ValueInt64() <= 0 {
		resp.Diagnostics.AddError("Invalid VPS Virtual Machine ID", "The VPS virtual machine ID is unknown, null, or zero, so Monarx cannot be read.")
		return
	}

	timeout, diags := state.Timeouts.Read(ctx, 20*time.Minute)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	response, err := r.client.VPSGetScanMetricsV1WithResponse(ctx, client.VirtualMachineId(state.VirtualMachineID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError("Unable to Read Monarx", fmt.Sprintf("Got error: %s", err))
		return
	}
	if response.StatusCode() == http.StatusUnprocessableEntity || response.StatusCode() == http.StatusNotFound {
		resp.State.RemoveResource(ctx)
		return
	}
	if response.StatusCode() != http.StatusOK || response.JSON200 == nil {
		resp.Diagnostics.AddError("Unable to Read Monarx", fmt.Sprintf("Unexpected status code: %d, response body is nil or invalid.", response.StatusCode()))
		return
	}

	identity := VPSMonarxIdentity{VirtualMachineID: state.VirtualMachineID}
	resp.Diagnostics.Append(resp.Identity.Set(ctx, &identity)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *VPSMonarxResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan VPSMonarxResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	identity := VPSMonarxIdentity{VirtualMachineID: plan.VirtualMachineID}
	resp.Diagnostics.Append(resp.Identity.Set(ctx, &identity)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *VPSMonarxResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state VPSMonarxResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if state.VirtualMachineID.IsUnknown() || state.VirtualMachineID.IsNull() || state.VirtualMachineID.ValueInt64() <= 0 {
		resp.Diagnostics.AddError("Invalid VPS Virtual Machine ID", "The VPS virtual machine ID is unknown, null, or zero, so Monarx cannot be uninstalled.")
		return
	}

	timeout, diags := state.Timeouts.Delete(ctx, 20*time.Minute)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	response, err := r.client.VPSUninstallMonarxV1WithResponse(ctx, client.VirtualMachineId(state.VirtualMachineID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError("Unable to Uninstall Monarx", fmt.Sprintf("Got error: %s", err))
		return
	}
	if response.StatusCode() == http.StatusUnprocessableEntity || response.StatusCode() == http.StatusNotFound {
		return
	}
	if response.StatusCode() != http.StatusOK || response.JSON200 == nil || response.JSON200.Id == nil || *response.JSON200.Id <= 0 {
		resp.Diagnostics.AddError("Unable to Uninstall Monarx", "The response did not contain a valid action ID.")
		return
	}

	if state.WaitForAction.IsNull() || state.WaitForAction.IsUnknown() || state.WaitForAction.ValueBool() {
		actionID := *response.JSON200.Id
		for {
			actionResponse, actionErr := r.client.VPSGetActionDetailsV1WithResponse(ctx, client.VirtualMachineId(state.VirtualMachineID.ValueInt64()), actionID)
			if actionErr != nil {
				resp.Diagnostics.AddError("Unable to Uninstall Monarx", fmt.Sprintf("Unable to get Monarx action: %s", actionErr))
				return
			}
			if actionResponse.StatusCode() != http.StatusOK || actionResponse.JSON200 == nil || actionResponse.JSON200.State == nil {
				resp.Diagnostics.AddError("Unable to Uninstall Monarx", fmt.Sprintf("Unexpected action response: status code %d", actionResponse.StatusCode()))
				return
			}
			switch *actionResponse.JSON200.State {
			case client.VPSV1ActionActionResourceStateSuccess:
				return
			case client.VPSV1ActionActionResourceStateError:
				resp.Diagnostics.AddError("Unable to Uninstall Monarx", "The Monarx uninstallation action failed")
				return
			default:
				select {
				case <-ctx.Done():
					resp.Diagnostics.AddError("Unable to Uninstall Monarx", fmt.Sprintf("The Monarx uninstallation action timed out: %s", ctx.Err()))
					return
				case <-time.After(2 * time.Second):
				}
			}
		}
	}
}
