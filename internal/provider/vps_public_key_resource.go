// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"net/http"
	"slices"
	"time"

	"github.com/comradesharf/terraform-provider-hostinger/internal/client"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/identityschema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

var (
	_ resource.Resource                = &VPSPublicKeyResource{}
	_ resource.ResourceWithConfigure   = &VPSPublicKeyResource{}
	_ resource.ResourceWithImportState = &VPSPublicKeyResource{}
	_ resource.ResourceWithIdentity    = &VPSPublicKeyResource{}
)

func NewVPSPublicKeyResource() resource.Resource {
	return &VPSPublicKeyResource{}
}

type VPSPublicKeyResource struct {
	client *client.ClientWithResponses
}

type VPSPublicKeyResourceModel struct {
	VPSPublicKeyModel
	Timeouts timeouts.Value `tfsdk:"timeouts"`
}

func (r *VPSPublicKeyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vps_public_key"
}

func (r *VPSPublicKeyResource) IdentitySchema(ctx context.Context, req resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = identityschema.Schema{
		Attributes: map[string]identityschema.Attribute{
			"id": identityschema.Int64Attribute{
				RequiredForImport: true,
			},
		},
	}
}

func (r *VPSPublicKeyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages an SSH public key in the Hostinger account.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Public key ID.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Public key name.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"key": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "SSH public key content.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"timeouts": timeouts.AttributesAll(ctx),
		},
	}
}

func (r *VPSPublicKeyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *VPSPublicKeyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan VPSPublicKeyResourceModel
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

	response, err := r.client.VPSCreatePublicKeyV1WithResponse(ctx, client.VPSCreatePublicKeyV1JSONRequestBody{
		Name: plan.Name.ValueString(),
		Key:  plan.Key.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Unable to Create VPS Public Key", fmt.Sprintf("Got error: %s", err))
		return
	}
	if response.StatusCode() != http.StatusOK {
		resp.Diagnostics.AddError("Unable to Create VPS Public Key", fmt.Sprintf("Unexpected status code: %d, response: %s", response.StatusCode(), string(response.Body)))
		return
	}
	if response.JSON200 == nil {
		resp.Diagnostics.AddError("Unable to Create VPS Public Key", "Response body is nil")
		return
	}

	plan.Merge(*response.JSON200)

	identity := VPSPublicKeyIdentity{ID: plan.ID}
	resp.Diagnostics.Append(resp.Identity.Set(ctx, &identity)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *VPSPublicKeyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state VPSPublicKeyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var identity VPSPublicKeyIdentity
	resp.Diagnostics.Append(req.Identity.Get(ctx, &identity)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.ID.IsUnknown() || state.ID.IsNull() || state.ID.ValueInt64() == 0 {
		resp.Diagnostics.AddError("Invalid Public Key ID", "The public key ID is unknown, null, or zero, so it cannot be read.")
		return
	}

	readTimeout, diags := state.Timeouts.Read(ctx, 20*time.Minute)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, readTimeout)
	defer cancel()

	var found *client.VPSV1PublicKeyPublicKeyResource
	params := client.VPSGetPublicKeysV1Params{}
	for page := 1; ; page++ {
		params.Page = &page
		response, err := r.client.VPSGetPublicKeysV1WithResponse(ctx, &params)
		if err != nil {
			resp.Diagnostics.AddError("Unable to Read VPS Public Keys", fmt.Sprintf("Got error: %s", err))
			return
		}
		if response.StatusCode() != http.StatusOK {
			resp.Diagnostics.AddError("Unable to Read VPS Public Keys", fmt.Sprintf("Unexpected status code: %d, response: %s", response.StatusCode(), string(response.Body)))
			return
		}
		if response.JSON200 == nil || response.JSON200.Data == nil {
			break
		}

		items := *response.JSON200.Data
		index := slices.IndexFunc(items, func(item client.VPSV1PublicKeyPublicKeyResource) bool {
			return item.Id != nil && int64(*item.Id) == state.ID.ValueInt64()
		})
		if index >= 0 {
			item := items[index]
			found = &item
			break
		}

		meta := response.JSON200.Meta
		if len(items) == 0 || meta == nil || meta.CurrentPage == nil || meta.PerPage == nil || meta.Total == nil || (*meta.CurrentPage)*(*meta.PerPage) >= *meta.Total {
			break
		}
	}

	if found == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	state.Merge(*found)

	identity.ID = state.ID

	resp.Diagnostics.Append(resp.Identity.Set(ctx, &identity)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *VPSPublicKeyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state VPSPublicKeyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var identity VPSPublicKeyIdentity
	resp.Diagnostics.Append(req.Identity.Get(ctx, &identity)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.Identity.Set(ctx, &identity)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *VPSPublicKeyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state VPSPublicKeyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.ID.IsUnknown() || state.ID.IsNull() || state.ID.ValueInt64() == 0 {
		resp.Diagnostics.AddError("Invalid Public Key ID", "The public key ID is unknown, null, or zero, so it cannot be deleted.")
		return
	}

	deleteTimeout, diags := state.Timeouts.Delete(ctx, 20*time.Minute)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, deleteTimeout)
	defer cancel()

	response, err := r.client.VPSDeletePublicKeyV1WithResponse(ctx, client.PublicKeyId(state.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError("Unable to Delete VPS Public Key", fmt.Sprintf("Got error: %s", err))
		return
	}
	if response.StatusCode() != http.StatusOK && response.StatusCode() != http.StatusNotFound {
		resp.Diagnostics.AddError("Unable to Delete VPS Public Key", fmt.Sprintf("Unexpected status code: %d, response: %s", response.StatusCode(), string(response.Body)))
	}
}

func (r *VPSPublicKeyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importStatePassthroughInt64Identity(ctx, path.Root("id"), path.Root("id"), req, resp)
}
