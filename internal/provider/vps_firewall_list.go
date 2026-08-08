// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"net/http"

	"github.com/comradesharf/terraform-provider-hostinger/internal/client"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/list/timeouts"
	resourcetimeouts "github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ list.ListResource              = &VPSFirewallList{}
	_ list.ListResourceWithConfigure = &VPSFirewallList{}
)

func NewVPSFirewallList() list.ListResource {
	return &VPSFirewallList{}
}

// VPSFirewallList defines the list resource implementation.
type VPSFirewallList struct {
	client *client.ClientWithResponses
}

type VPSFirewallListModel struct {
	Timeouts *timeouts.Type `tfsdk:"timeouts"`
}

func (l *VPSFirewallList) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	l.client = c
}

func (l *VPSFirewallList) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vps_firewall"
}

func (l *VPSFirewallList) ListResourceConfigSchema(ctx context.Context, req list.ListResourceSchemaRequest, resp *list.ListResourceSchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "This resource provides a list of VPS firewalls.",
		Attributes: map[string]schema.Attribute{
			"timeouts": timeouts.Attributes(ctx),
		},
	}
}

func (l *VPSFirewallList) List(ctx context.Context, req list.ListRequest, stream *list.ListResultsStream) {
	var data VPSFirewallListModel
	diags := req.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	stream.Results = func(push func(list.ListResult) bool) {
		params := &client.VPSGetFirewallListV1Params{}
		page := 1

		for {
			params.Page = &page
			ctx = tflog.SetField(ctx, "page", params.Page)

			response, err := l.client.VPSGetFirewallListV1WithResponse(ctx, params)
			if err != nil {
				result := req.NewListResult(ctx)
				result.Diagnostics.AddError(
					"Unable to Read VPS Firewalls",
					fmt.Sprintf("Got error: %s", err),
				)
				push(result)
				return
			}
			if response.StatusCode() != http.StatusOK {
				result := req.NewListResult(ctx)
				result.Diagnostics.AddError(
					"Unable to Read VPS Firewalls",
					fmt.Sprintf("Unexpected status code: %d, response: %s", response.StatusCode(), string(response.Body)),
				)
				push(result)
				return
			}
			if response.JSON200 == nil {
				result := req.NewListResult(ctx)
				result.Diagnostics.AddError(
					"Unable to Read VPS Firewalls",
					"Response body is nil",
				)
				push(result)
				return
			}
			if response.JSON200.Data == nil || len(*response.JSON200.Data) == 0 {
				break
			}

			for _, item := range *response.JSON200.Data {
				var d VPSFirewallResourceModel
				d.Merge(item)

				result := req.NewListResult(ctx)
				result.DisplayName = d.Name.ValueString()

				identity := VPSFirewallIdentity{ID: d.ID}
				result.Diagnostics.Append(result.Identity.Set(ctx, &identity)...)

				if req.IncludeResource {
					var r VPSFirewallResourceModel
					r.ID = d.ID
					r.Name = d.Name
					r.IsSynced = d.IsSynced
					r.CreatedAt = d.CreatedAt
					r.UpdatedAt = d.UpdatedAt
					r.Timeouts = resourcetimeouts.Value{
						Object: types.ObjectNull(map[string]attr.Type{
							"create": types.StringType,
							"read":   types.StringType,
							"update": types.StringType,
							"delete": types.StringType,
						}),
					}

					result.Diagnostics.Append(result.Resource.Set(ctx, &r)...)
				}

				if !push(result) {
					return
				}
			}

			meta := response.JSON200.Meta
			if meta == nil || meta.CurrentPage == nil || meta.PerPage == nil || meta.Total == nil {
				break
			}
			fetched := (*meta.CurrentPage) * (*meta.PerPage)
			if fetched >= *meta.Total {
				break
			}
			page++
		}
	}
}
