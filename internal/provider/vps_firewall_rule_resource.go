// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/comradesharf/terraform-provider-hostinger/internal/client"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/identityschema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &VPSFirewallRuleResource{}
	_ resource.ResourceWithImportState = &VPSFirewallRuleResource{}
	_ resource.ResourceWithConfigure   = &VPSFirewallRuleResource{}
	_ resource.ResourceWithIdentity    = &VPSFirewallRuleResource{}
)

func NewVPSFirewallRuleResource() resource.Resource {
	return &VPSFirewallRuleResource{}
}

// VPSFirewallRuleResource defines the resource implementation.
type VPSFirewallRuleResource struct {
	client *client.ClientWithResponses
}

// VPSFirewallRuleResourceModel describes the resource data model.
type VPSFirewallRuleResourceModel struct {
	VPSFirewallRuleModel
	FirewallID types.Int64    `tfsdk:"firewall_id"`
	Timeouts   timeouts.Value `tfsdk:"timeouts"`
}

func (r *VPSFirewallRuleResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vps_firewall_rule"
}

func (r *VPSFirewallRuleResource) IdentitySchema(ctx context.Context, req resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = identityschema.Schema{
		Attributes: map[string]identityschema.Attribute{
			"id": identityschema.Int64Attribute{
				Description:       "Firewall rule ID",
				RequiredForImport: true,
			},
			"firewall_id": identityschema.Int64Attribute{
				Description:       "ID of the firewall this rule belongs to",
				RequiredForImport: true,
			},
		},
	}
}

func (r *VPSFirewallRuleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a VPS firewall rule.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Firewall rule ID",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"firewall_id": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "ID of the firewall this rule belongs to",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"action": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Firewall rule action",
				Validators: []validator.String{
					stringvalidator.OneOf(
						"accept",
						"drop",
					),
				},
			},
			"protocol": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Firewall rule protocol",
				Validators: []validator.String{
					stringvalidator.OneOf(
						"TCP",
						"UDP",
						"ICMP",
						"GRE",
						"any",
						"ESP",
						"AH",
						"ICMPv6",
						"SSH",
						"HTTP",
						"HTTPS",
						"MySQL",
						"PostgreSQL",
					),
				},
			},
			"port": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Firewall rule destination port: single or port range",
			},
			"source": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Firewall rule source. Can be `any` or `custom`",
				Validators: []validator.String{
					stringvalidator.OneOf("any", "custom"),
				},
			},
			"source_detail": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Firewall rule source detail. Can be `any` or IP address, CIDR or range",
			},
			"timeouts": timeouts.AttributesAll(ctx),
		},
	}
}

func (r *VPSFirewallRuleResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = c
}

func (r *VPSFirewallRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan VPSFirewallRuleResourceModel
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

	p := client.VPSCreateFirewallRuleV1JSONRequestBody{
		Protocol:     client.VPSV1FirewallRulesStoreRequestProtocol(plan.Protocol.ValueString()),
		Port:         plan.Port.ValueString(),
		Source:       client.VPSV1FirewallRulesStoreRequestSource(plan.Source.ValueString()),
		SourceDetail: plan.SourceDetail.ValueString(),
	}

	response, err := r.client.VPSCreateFirewallRuleV1WithResponse(ctx, client.FirewallId(plan.FirewallID.ValueInt64()), p)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create VPS Firewall Rule",
			fmt.Sprintf("Got error: %s", err),
		)
		return
	}
	if response.StatusCode() != http.StatusOK {
		resp.Diagnostics.AddError(
			"Unable to Create VPS Firewall Rule",
			fmt.Sprintf("Unexpected status code: %d, response: %s", response.StatusCode(), string(response.Body)),
		)
		return
	}
	if response.JSON200 == nil {
		resp.Diagnostics.AddError(
			"Unable to Create VPS Firewall Rule",
			"Response body is nil",
		)
		return
	}

	plan.Merge(*response.JSON200)

	identity := VPSFirewallRuleIdentity{ID: plan.ID, FirewallID: plan.FirewallID}
	resp.Diagnostics.Append(resp.Identity.Set(ctx, &identity)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *VPSFirewallRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state VPSFirewallRuleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var identity VPSFirewallRuleIdentity
	resp.Diagnostics.Append(req.Identity.Get(ctx, &identity)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.ID.IsUnknown() {
		resp.Diagnostics.AddError(
			"Unknown ID",
			"ID is unknown, unable to read VPS firewall.",
		)
		return
	}

	if state.ID.IsNull() || state.ID.ValueInt64() == 0 {
		resp.Diagnostics.AddError(
			"Null ID",
			"ID is null or zero, unable to read VPS firewall.",
		)
		return
	}

	readTimeout, diags := state.Timeouts.Read(ctx, 20*time.Minute)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, readTimeout)
	defer cancel()

	response, err := r.client.VPSGetFirewallDetailsV1WithResponse(ctx, client.FirewallId(state.FirewallID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read VPS Firewalls",
			fmt.Sprintf("Got error: %s", err),
		)
		return
	}
	if response.StatusCode() == http.StatusNotFound {
		resp.State.RemoveResource(ctx)
		return
	}
	if response.StatusCode() != http.StatusOK {
		resp.Diagnostics.AddError(
			"Unable to Read VPS Firewalls",
			fmt.Sprintf("Unexpected status code: %d, response: %s", response.StatusCode(), string(response.Body)),
		)
		return
	}
	if response.JSON200 == nil {
		resp.Diagnostics.AddError(
			"Unable to Read VPS Firewalls",
			"Response body is nil",
		)
		return
	}
	if response.JSON200.Rules == nil || len(*response.JSON200.Rules) == 0 {
		resp.State.RemoveResource(ctx)
		return
	}

	rules := *response.JSON200.Rules

	index := slices.IndexFunc(rules, func(rule client.VPSV1FirewallFirewallRuleResource) bool {
		if rule.Id == nil {
			return false
		}
		return types.Int64Value(int64(*rule.Id)).Equal(state.ID)
	})

	if index == -1 {
		resp.State.RemoveResource(ctx)
		return
	}

	state.Merge(rules[index])

	identity.ID = state.ID
	identity.FirewallID = state.FirewallID

	resp.Diagnostics.Append(resp.Identity.Set(ctx, &identity)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *VPSFirewallRuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan VPSFirewallRuleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var identity VPSFirewallRuleIdentity
	resp.Diagnostics.Append(req.Identity.Get(ctx, &identity)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.ID.IsUnknown() {
		resp.Diagnostics.AddError(
			"Unknown ID",
			"ID is unknown, unable to read VPS firewall.",
		)
		return
	}

	if plan.ID.IsNull() || plan.ID.ValueInt64() == 0 {
		resp.Diagnostics.AddError(
			"Null ID",
			"ID is null or zero, unable to read VPS firewall.",
		)
		return
	}

	updateTimeout, diags := plan.Timeouts.Update(ctx, 20*time.Minute)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, updateTimeout)
	defer cancel()

	params := client.VPSUpdateFirewallRuleV1JSONRequestBody{
		Protocol:     client.VPSV1FirewallRulesStoreRequestProtocol(plan.Protocol.ValueString()),
		Port:         plan.Port.ValueString(),
		Source:       client.VPSV1FirewallRulesStoreRequestSource(plan.Source.ValueString()),
		SourceDetail: plan.SourceDetail.ValueString(),
	}

	response, err := r.client.VPSUpdateFirewallRuleV1WithResponse(
		ctx,
		client.FirewallId(plan.FirewallID.ValueInt64()),
		client.RuleId(plan.ID.ValueInt64()),
		params,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Update VPS Firewall Rule",
			fmt.Sprintf("Got error: %s", err),
		)
		return
	}
	if response.StatusCode() != http.StatusOK {
		resp.Diagnostics.AddError(
			"Unable to Update VPS Firewall Rule",
			fmt.Sprintf("Unexpected status code: %d, response: %s", response.StatusCode(), string(response.Body)),
		)
		return
	}
	if response.JSON200 == nil {
		resp.Diagnostics.AddError(
			"Unable to Update VPS Firewall Rule",
			"Response body is nil",
		)
		return
	}

	plan.Merge(*response.JSON200)

	identity.ID = plan.ID
	identity.FirewallID = plan.FirewallID

	resp.Diagnostics.Append(resp.Identity.Set(ctx, &identity)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *VPSFirewallRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state VPSFirewallRuleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.ID.IsUnknown() {
		resp.Diagnostics.AddError(
			"Unknown ID",
			"ID is unknown, unable to read VPS firewall.",
		)
		return
	}

	if state.ID.IsNull() || state.ID.ValueInt64() == 0 {
		resp.Diagnostics.AddError(
			"Null ID",
			"ID is null or zero, unable to read VPS firewall.",
		)
		return
	}

	deleteTimeout, diags := state.Timeouts.Delete(ctx, 20*time.Minute)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, deleteTimeout)
	defer cancel()

	response, err := r.client.VPSDeleteFirewallRuleV1WithResponse(
		ctx,
		client.FirewallId(state.FirewallID.ValueInt64()),
		client.RuleId(state.ID.ValueInt64()),
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Delete VPS Firewall Rule",
			fmt.Sprintf("Got error: %s", err),
		)
		return
	}
	if response.StatusCode() != http.StatusOK {
		resp.Diagnostics.AddError(
			"Unable to Delete VPS Firewalls",
			fmt.Sprintf("Unexpected status code: %d, response: %s", response.StatusCode(), string(response.Body)),
		)
		return
	}
}

func (r *VPSFirewallRuleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Debug(ctx, "Importing VPS Firewall Rule resource state", map[string]any{"id": req.ID})
	var identity VPSFirewallRuleIdentity

	if req.ID != "" {
		parts := strings.Split(req.ID, "/")
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			resp.Diagnostics.AddError(
				"Resource Import Passthrough Invalid ID",
				"Import ID must be in the format 'firewall_id/rule_id'",
			)
			return
		}

		firewallID, err := strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			resp.Diagnostics.AddError(
				"Resource Import Passthrough Invalid ID",
				fmt.Sprintf(
					"Failed to parse import ID as int64: %s", err.Error(),
				),
			)
			return
		}

		ruleID, err := strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			resp.Diagnostics.AddError(
				"Resource Import Passthrough Invalid ID",
				fmt.Sprintf(
					"Failed to parse import ID as int64: %s", err.Error(),
				),
			)
			return
		}

		identity.ID = types.Int64Value(ruleID)
		identity.FirewallID = types.Int64Value(firewallID)

		resp.Diagnostics.Append(resp.Identity.Set(ctx, &identity)...)
	} else {
		resp.Diagnostics.Append(resp.Identity.Get(ctx, &identity)...)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), identity.ID)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("firewall_id"), identity.FirewallID)...)
}
