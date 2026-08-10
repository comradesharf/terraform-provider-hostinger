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
	_ action.Action              = &VPSPublicKeysAttachAction{}
	_ action.ActionWithConfigure = &VPSPublicKeysAttachAction{}
)

func NewVPSPublicKeysAttachAction() action.Action {
	return &VPSPublicKeysAttachAction{}
}

type VPSPublicKeysAttachAction struct {
	client *client.ClientWithResponses
}

type VPSPublicKeysAttachActionModel struct {
	Timeouts         timeouts.Value `tfsdk:"timeouts"`
	PublicKeyIDs     types.List     `tfsdk:"public_key_ids"`
	VirtualMachineID types.Int64    `tfsdk:"virtual_machine_id"`
	WaitForAction    types.Bool     `tfsdk:"wait_for_action"`
}

func (a *VPSPublicKeysAttachAction) Configure(ctx context.Context, req action.ConfigureRequest, resp *action.ConfigureResponse) {
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

func (a *VPSPublicKeysAttachAction) Schema(ctx context.Context, req action.SchemaRequest, resp *action.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Attaches VPS public keys to a virtual machine.",
		Attributes: map[string]schema.Attribute{
			"public_key_ids": schema.ListAttribute{
				Required:            true,
				ElementType:         types.Int64Type,
				MarkdownDescription: "The IDs of the public keys to attach.",
			},
			"virtual_machine_id": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "The ID of the virtual machine to attach the public keys to.",
			},
			"wait_for_action": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Whether to wait for the action to complete before returning. Defaults to true.",
			},
			"timeouts": timeouts.Attributes(ctx),
		},
	}
}

func (a *VPSPublicKeysAttachAction) Metadata(ctx context.Context, req action.MetadataRequest, resp *action.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vps_public_keys_attach"
}

func (a *VPSPublicKeysAttachAction) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	var config VPSPublicKeysAttachActionModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.VirtualMachineID.IsNull() || config.VirtualMachineID.IsUnknown() || config.VirtualMachineID.ValueInt64() <= 0 {
		resp.Diagnostics.AddError(
			"Invalid Virtual Machine ID",
			"The virtual_machine_id attribute is required and must be a positive integer.",
		)
	}
	if config.PublicKeyIDs.IsNull() || config.PublicKeyIDs.IsUnknown() || len(config.PublicKeyIDs.Elements()) == 0 {
		resp.Diagnostics.AddError(
			"Invalid Public Key IDs",
			"The public_key_ids attribute is required and must contain at least one positive integer.",
		)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	var publicKeyIDs []int64
	resp.Diagnostics.Append(config.PublicKeyIDs.ElementsAs(ctx, &publicKeyIDs, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ids := make([]int, 0, len(publicKeyIDs))
	for _, publicKeyID := range publicKeyIDs {
		if publicKeyID <= 0 {
			resp.Diagnostics.AddError(
				"Invalid Public Key IDs",
				"The public_key_ids attribute is required and must contain at least one positive integer.",
			)
			return
		}
		ids = append(ids, int(publicKeyID))
	}

	invokeTimeout, diags := config.Timeouts.Invoke(ctx, 20*time.Minute)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "virtual_machine_id", config.VirtualMachineID.ValueInt64())

	ctx, cancel := context.WithTimeout(ctx, invokeTimeout)
	defer cancel()

	response, err := a.client.VPSAttachPublicKeyV1WithResponse(
		ctx,
		client.VirtualMachineId(config.VirtualMachineID.ValueInt64()),
		client.VPSAttachPublicKeyV1JSONRequestBody{Ids: ids},
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Attach VPS Public Keys",
			fmt.Sprintf("Got error: %s", err),
		)
		return
	}
	if response.StatusCode() != http.StatusOK {
		resp.Diagnostics.AddError(
			"Unable to Attach VPS Public Keys",
			fmt.Sprintf("Unexpected status code: %d, response: %s", response.StatusCode(), string(response.Body)),
		)
		return
	}
	if response.JSON200 == nil || response.JSON200.Id == nil {
		resp.Diagnostics.AddError(
			"Unable to Attach VPS Public Keys",
			"Response body is nil or missing ID",
		)
		return
	}

	if config.WaitForAction.IsNull() || config.WaitForAction.IsUnknown() || config.WaitForAction.ValueBool() {
	poll:
		for {
			response, err := a.client.VPSGetActionDetailsV1WithResponse(
				ctx,
				client.VirtualMachineId(config.VirtualMachineID.ValueInt64()),
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
				resp.Diagnostics.AddError("Unable to Attach VPS Public Keys", "The public key attach action failed")
				return
			default:
				select {
				case <-ctx.Done():
					resp.Diagnostics.AddError("Unable to Attach VPS Public Keys", "The action to attach VPS public keys timed out")
					return
				case <-time.After(2 * time.Second):
					resp.SendProgress(action.InvokeProgressEvent{
						Message: "Waiting for action to complete...",
					})
				}
			}
		}
	}
}
