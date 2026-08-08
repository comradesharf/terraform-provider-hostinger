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
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/identityschema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var (
	_ resource.Resource                = &VPSPostInstallScriptResource{}
	_ resource.ResourceWithConfigure   = &VPSPostInstallScriptResource{}
	_ resource.ResourceWithImportState = &VPSPostInstallScriptResource{}
	_ resource.ResourceWithIdentity    = &VPSPostInstallScriptResource{}
)

func NewVPSPostInstallScriptResource() resource.Resource {
	return &VPSPostInstallScriptResource{}
}

type VPSPostInstallScriptResource struct {
	client *client.ClientWithResponses
}

type VPSPostInstallScriptResourceModel struct {
	VPSPostInstallScriptModel
	Timeouts timeouts.Value `tfsdk:"timeouts"`
}

func (r *VPSPostInstallScriptResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vps_post_install_script"
}

func (r *VPSPostInstallScriptResource) IdentitySchema(ctx context.Context, req resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = identityschema.Schema{
		Attributes: map[string]identityschema.Attribute{
			"id": identityschema.Int64Attribute{
				RequiredForImport: true,
			},
		},
	}
}

func (r *VPSPostInstallScriptResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a VPS post-install script.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Post-install script ID.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Post-install script name.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"content": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Post-install script content.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"created_at": schema.StringAttribute{
				Computed:            true,
				CustomType:          timetypes.RFC3339Type{},
				MarkdownDescription: "Timestamp when the post-install script was created (RFC3339).",
			},
			"updated_at": schema.StringAttribute{
				Computed:            true,
				CustomType:          timetypes.RFC3339Type{},
				MarkdownDescription: "Timestamp when the post-install script was updated (RFC3339).",
			},
			"timeouts": timeouts.AttributesAll(ctx),
		},
	}
}

func (r *VPSPostInstallScriptResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *VPSPostInstallScriptResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state VPSPostInstallScriptResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createTimeout, diags := state.Timeouts.Create(ctx, 20*time.Minute)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, createTimeout)
	defer cancel()

	response, err := r.client.VPSCreatePostInstallScriptV1WithResponse(ctx, client.VPSCreatePostInstallScriptV1JSONRequestBody{
		Name:    state.Name.ValueString(),
		Content: state.Content.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Unable to Create VPS Post-install Script", fmt.Sprintf("Got error: %s", err))
		return
	}
	if response.StatusCode() != http.StatusOK {
		resp.Diagnostics.AddError("Unable to Create VPS Post-install Script", fmt.Sprintf("Unexpected status code: %d, response: %s", response.StatusCode(), string(response.Body)))
		return
	}
	if response.JSON200 == nil {
		resp.Diagnostics.AddError("Unable to Create VPS Post-install Script", "Response body is nil")
		return
	}

	state.Merge(*response.JSON200)

	identity := VPSPostInstallScriptIdentity{ID: state.ID}
	resp.Diagnostics.Append(resp.Identity.Set(ctx, &identity)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *VPSPostInstallScriptResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state VPSPostInstallScriptResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var identity VPSPostInstallScriptIdentity
	resp.Diagnostics.Append(req.Identity.Get(ctx, &identity)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.ID.IsUnknown() || state.ID.IsNull() || state.ID.ValueInt64() == 0 {
		resp.Diagnostics.AddError("Invalid Post-install Script ID", "The post-install script ID is unknown, null, or zero, so it cannot be read.")
		return
	}

	readTimeout, diags := state.Timeouts.Read(ctx, 20*time.Minute)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, readTimeout)
	defer cancel()

	response, err := r.client.VPSGetPostInstallScriptV1WithResponse(ctx, client.PostInstallScriptId(state.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError("Unable to Read VPS Post-install Script", fmt.Sprintf("Got error: %s", err))
		return
	}
	if response.StatusCode() == http.StatusNotFound {
		resp.State.RemoveResource(ctx)
		return
	}
	if response.StatusCode() != http.StatusOK {
		resp.Diagnostics.AddError("Unable to Read VPS Post-install Script", fmt.Sprintf("Unexpected status code: %d, response: %s", response.StatusCode(), string(response.Body)))
		return
	}
	if response.JSON200 == nil {
		resp.Diagnostics.AddError("Unable to Read VPS Post-install Script", "Response body is nil")
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

func (r *VPSPostInstallScriptResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state VPSPostInstallScriptResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var identity VPSPostInstallScriptIdentity
	resp.Diagnostics.Append(req.Identity.Get(ctx, &identity)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.ID.IsUnknown() || state.ID.IsNull() || state.ID.ValueInt64() == 0 {
		resp.Diagnostics.AddError("Invalid Post-install Script ID", "The post-install script ID is unknown, null, or zero, so it cannot be updated.")
		return
	}

	updateTimeout, diags := plan.Timeouts.Update(ctx, 20*time.Minute)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, updateTimeout)
	defer cancel()

	response, err := r.client.VPSUpdatePostInstallScriptV1WithResponse(ctx, client.PostInstallScriptId(state.ID.ValueInt64()), client.VPSUpdatePostInstallScriptV1JSONRequestBody{
		Name:    plan.Name.ValueString(),
		Content: plan.Content.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Unable to Update VPS Post-install Script", fmt.Sprintf("Got error: %s", err))
		return
	}
	if response.StatusCode() != http.StatusOK {
		resp.Diagnostics.AddError("Unable to Update VPS Post-install Script", fmt.Sprintf("Unexpected status code: %d, response: %s", response.StatusCode(), string(response.Body)))
		return
	}
	if response.JSON200 == nil {
		resp.Diagnostics.AddError("Unable to Update VPS Post-install Script", "Response body is nil")
		return
	}

	plan.Merge(*response.JSON200)

	identity.ID = plan.ID

	resp.Diagnostics.Append(resp.Identity.Set(ctx, &identity)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *VPSPostInstallScriptResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state VPSPostInstallScriptResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if state.ID.IsUnknown() || state.ID.IsNull() || state.ID.ValueInt64() == 0 {
		resp.Diagnostics.AddError("Invalid Post-install Script ID", "The post-install script ID is unknown, null, or zero, so it cannot be deleted.")
		return
	}

	deleteTimeout, diags := state.Timeouts.Delete(ctx, 20*time.Minute)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, deleteTimeout)
	defer cancel()

	response, err := r.client.VPSDeletePostInstallScriptV1WithResponse(ctx, client.PostInstallScriptId(state.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError("Unable to Delete VPS Post-install Script", fmt.Sprintf("Got error: %s", err))
		return
	}
	if response.StatusCode() != http.StatusOK && response.StatusCode() != http.StatusNotFound {
		resp.Diagnostics.AddError("Unable to Delete VPS Post-install Script", fmt.Sprintf("Unexpected status code: %d, response: %s", response.StatusCode(), string(response.Body)))
	}
}

func (r *VPSPostInstallScriptResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importStatePassthroughInt64Identity(ctx, path.Root("id"), path.Root("id"), req, resp)
}
