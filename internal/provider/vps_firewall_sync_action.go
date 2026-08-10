// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/comradesharf/terraform-provider-hostinger/internal/client"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/action/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/action/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ action.Action              = &VPSFirewallSyncAction{}
	_ action.ActionWithConfigure = &VPSFirewallSyncAction{}
)

func NewVPSFirewallSyncAction() action.Action {
	return &VPSFirewallSyncAction{}
}

type VPSFirewallSyncAction struct {
	client *client.ClientWithResponses
}

type VPSFirewallSyncActionModel struct {
	Timeouts         timeouts.Value `tfsdk:"timeouts"`
	FirewallID       types.Int64    `tfsdk:"firewall_id"`
	VirtualMachineID types.Int64    `tfsdk:"virtual_machine_id"`
	WaitForAction    types.Bool     `tfsdk:"wait_for_action"`
}

func (a *VPSFirewallSyncAction) Configure(ctx context.Context, req action.ConfigureRequest, resp *action.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	c, ok := req.ProviderData.(*client.ClientWithResponses)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Action Configure Type",
			fmt.Sprintf("Expected *client.ClientWithResponses, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	a.client = c
}

func (a *VPSFirewallSyncAction) Schema(ctx context.Context, req action.SchemaRequest, resp *action.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Synchronizes a VPS firewall with a virtual machine.",
		Attributes: map[string]schema.Attribute{
			"firewall_id": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "The ID of the firewall to synchronize.",
			},
			"virtual_machine_id": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "The ID of the virtual machine with which to synchronize the firewall.",
			},
			"wait_for_action": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Whether to wait for the action to complete before returning. Defaults to true.",
			},
			"timeouts": timeouts.Attributes(ctx),
		},
	}
}

func (a *VPSFirewallSyncAction) Metadata(ctx context.Context, req action.MetadataRequest, resp *action.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vps_firewall_sync"
}

func (a *VPSFirewallSyncAction) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	var config VPSFirewallSyncActionModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.FirewallID.IsNull() || config.FirewallID.IsUnknown() || config.FirewallID.ValueInt64() <= 0 {
		resp.Diagnostics.AddError("Invalid Firewall ID", "The firewall_id attribute is required and must be a positive integer.")
	}
	if config.VirtualMachineID.IsNull() || config.VirtualMachineID.IsUnknown() || config.VirtualMachineID.ValueInt64() <= 0 {
		resp.Diagnostics.AddError("Invalid Virtual Machine ID", "The virtual_machine_id attribute is required and must be a positive integer.")
	}
	if resp.Diagnostics.HasError() {
		return
	}

	invokeTimeout, diags := config.Timeouts.Invoke(ctx, 20*time.Minute)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "firewall_id", config.FirewallID.ValueInt64())
	ctx, cancel := context.WithTimeout(ctx, invokeTimeout)
	defer cancel()

	response, err := a.client.VPSSyncFirewallV1WithResponse(ctx, client.FirewallId(config.FirewallID.ValueInt64()), client.VirtualMachineId(config.VirtualMachineID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError("Unable to Synchronize VPS Firewall", fmt.Sprintf("Got error: %s", err))
		return
	}
	if response.StatusCode() != http.StatusOK {
		resp.Diagnostics.AddError("Unable to Synchronize VPS Firewall", fmt.Sprintf("Unexpected status code: %d, response: %s", response.StatusCode(), string(response.Body)))
		return
	}
	if response.JSON200 == nil || response.JSON200.Id == nil {
		resp.Diagnostics.AddError("Unable to Synchronize VPS Firewall", "Response body is nil or missing ID")
		return
	}

	if config.WaitForAction.IsNull() || config.WaitForAction.IsUnknown() || config.WaitForAction.ValueBool() {
	poll:
		for {
			response, err := a.client.VPSGetActionDetailsV1WithResponse(ctx, client.VirtualMachineId(config.VirtualMachineID.ValueInt64()), *response.JSON200.Id)
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
				resp.Diagnostics.AddError("Unable to Synchronize VPS Firewall", "The firewall synchronization action failed")
				return
			default:
				select {
				case <-ctx.Done():
					resp.Diagnostics.AddError("Unable to Synchronize VPS Firewall", "The firewall synchronization action timed out")
					return
				case <-time.After(2 * time.Second):
					resp.SendProgress(action.InvokeProgressEvent{Message: "Waiting for action to complete..."})
				}
			}
		}
	}
}
