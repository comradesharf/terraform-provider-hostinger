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
	_ action.Action              = &VPSVirtualMachineRecreateAction{}
	_ action.ActionWithConfigure = &VPSVirtualMachineRecreateAction{}
)

func NewVPSVirtualMachineRecreateAction() action.Action {
	return &VPSVirtualMachineRecreateAction{}
}

type VPSVirtualMachineRecreateAction struct {
	client *client.ClientWithResponses
}

type VPSVirtualMachineRecreateActionModel struct {
	Timeouts          timeouts.Value `tfsdk:"timeouts"`
	VirtualMachineID  types.Int64    `tfsdk:"virtual_machine_id"`
	TemplateID        types.Int64    `tfsdk:"template_id"`
	Password          types.String   `tfsdk:"password"`
	PanelPassword     types.String   `tfsdk:"panel_password"`
	PostInstallScript types.Int64    `tfsdk:"post_install_script_id"`
	WaitForAction     types.Bool     `tfsdk:"wait_for_action"`
}

func (a *VPSVirtualMachineRecreateAction) Configure(ctx context.Context, req action.ConfigureRequest, resp *action.ConfigureResponse) {
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

func (a *VPSVirtualMachineRecreateAction) Schema(ctx context.Context, req action.SchemaRequest, resp *action.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Recreates a VPS virtual machine with a new operating system template.",
		Attributes: map[string]schema.Attribute{
			"virtual_machine_id": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "The ID of the virtual machine to recreate.",
			},
			"template_id": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "The ID of the operating system template to install.",
			},
			"password": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The root password for the virtual machine. A random password is generated when omitted.",
			},
			"panel_password": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The panel password for panel-based operating system templates.",
			},
			"post_install_script_id": schema.Int64Attribute{
				Optional:            true,
				MarkdownDescription: "The ID of the post-install script to execute after recreation.",
			},
			"wait_for_action": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Whether to wait for the action to complete before returning. Defaults to true.",
			},
			"timeouts": timeouts.Attributes(ctx),
		},
	}
}

func (a *VPSVirtualMachineRecreateAction) Metadata(ctx context.Context, req action.MetadataRequest, resp *action.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vps_virtual_machine_recreate"
}

func (a *VPSVirtualMachineRecreateAction) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	var config VPSVirtualMachineRecreateActionModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.VirtualMachineID.IsNull() || config.VirtualMachineID.IsUnknown() || config.VirtualMachineID.ValueInt64() <= 0 {
		resp.Diagnostics.AddError("Invalid Virtual Machine ID", "The virtual_machine_id attribute is required and must be a positive integer.")
	}
	if config.TemplateID.IsNull() || config.TemplateID.IsUnknown() || config.TemplateID.ValueInt64() <= 0 {
		resp.Diagnostics.AddError("Invalid Template ID", "The template_id attribute is required and must be a positive integer.")
	}
	if !config.PostInstallScript.IsNull() && !config.PostInstallScript.IsUnknown() && config.PostInstallScript.ValueInt64() <= 0 {
		resp.Diagnostics.AddError("Invalid Post-Install Script ID", "The post_install_script_id attribute must be a positive integer when provided.")
	}
	if resp.Diagnostics.HasError() {
		return
	}

	invokeTimeout, diags := config.Timeouts.Invoke(ctx, 20*time.Minute)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "virtual_machine_id", config.VirtualMachineID.ValueInt64())
	ctx, cancel := context.WithTimeout(ctx, invokeTimeout)
	defer cancel()

	body := client.VPSRecreateVirtualMachineV1JSONRequestBody{
		TemplateId: int(config.TemplateID.ValueInt64()),
	}
	if !config.Password.IsNull() && !config.Password.IsUnknown() {
		body.Password = config.Password.ValueStringPointer()
	}
	if !config.PanelPassword.IsNull() && !config.PanelPassword.IsUnknown() {
		body.PanelPassword = config.PanelPassword.ValueStringPointer()
	}
	if !config.PostInstallScript.IsNull() && !config.PostInstallScript.IsUnknown() {
		postInstallScriptID := int(config.PostInstallScript.ValueInt64())
		body.PostInstallScriptId = &postInstallScriptID
	}

	response, err := a.client.VPSRecreateVirtualMachineV1WithResponse(
		ctx,
		client.VirtualMachineId(config.VirtualMachineID.ValueInt64()),
		body,
	)
	if err != nil {
		resp.Diagnostics.AddError("Unable to Recreate VPS Virtual Machine", fmt.Sprintf("Got error: %s", err))
		return
	}
	if response.StatusCode() != http.StatusOK {
		resp.Diagnostics.AddError("Unable to Recreate VPS Virtual Machine", fmt.Sprintf("Unexpected status code: %d, response: %s", response.StatusCode(), string(response.Body)))
		return
	}
	if response.JSON200 == nil || response.JSON200.Id == nil {
		resp.Diagnostics.AddError("Unable to Recreate VPS Virtual Machine", "Response body is nil or missing ID")
		return
	}

	if config.WaitForAction.IsNull() || config.WaitForAction.IsUnknown() || config.WaitForAction.ValueBool() {
	poll:
		for {
			actionResponse, actionErr := a.client.VPSGetActionDetailsV1WithResponse(
				ctx,
				client.VirtualMachineId(config.VirtualMachineID.ValueInt64()),
				*response.JSON200.Id,
			)
			if actionErr != nil {
				resp.Diagnostics.AddError("Unable to Get VPS Actions", fmt.Sprintf("Got error: %s", actionErr))
				return
			}
			if actionResponse.StatusCode() != http.StatusOK {
				resp.Diagnostics.AddError("Unable to Get VPS Actions", fmt.Sprintf("Unexpected status code: %d, response: %s", actionResponse.StatusCode(), string(actionResponse.Body)))
				return
			}
			if actionResponse.JSON200 == nil || actionResponse.JSON200.State == nil {
				resp.Diagnostics.AddError("Unable to Get VPS Actions", "Response body does not contain an action state")
				return
			}
			switch *actionResponse.JSON200.State {
			case client.VPSV1ActionActionResourceStateSuccess:
				break poll
			case client.VPSV1ActionActionResourceStateError:
				resp.Diagnostics.AddError("Unable to Recreate VPS Virtual Machine", "The virtual machine recreation action failed")
				return
			default:
				select {
				case <-ctx.Done():
					resp.Diagnostics.AddError("Unable to Recreate VPS Virtual Machine", fmt.Sprintf("The virtual machine recreation action timed out: %s", ctx.Err()))
					return
				case <-time.After(2 * time.Second):
					resp.SendProgress(action.InvokeProgressEvent{Message: "Waiting for action to complete..."})
				}
			}
		}
	}
}
